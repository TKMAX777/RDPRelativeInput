package windowselecter

import (
	"fmt"
	"runtime"
	"syscall"
	"unsafe"

	"github.com/TKMAX777/RDPRelativeInput/debug"
	"github.com/TKMAX777/RDPRelativeInput/winapi"
	"github.com/lxn/win"
	"github.com/pkg/errors"
)

type WindowInfo struct {
	Hwnd        win.HWND
	WindowTitle string
}

var ErrorQuitted = fmt.Errorf("Quitted")

const dialogBaseXPos = 10
const dialogBaseYPos = 10
const dialogWidth = 400
const dialogHeight = 400

func Dialog() (win.HWND, error) {
	var windows = make([]WindowInfo, 0)
	err := winapi.EnumDesktopWindows(0, syscall.NewCallback(newEnumWindowProc(&windows)), uintptr(unsafe.Pointer(&windows)))
	if err != nil {
		return 0, fmt.Errorf("EnumDesktopWindows: %w", err)
	}

	type resultAttr struct {
		hwnd win.HWND
		err  error
	}

	var result = make(chan resultAttr)
	var indexChannel = make(chan int)

	go func() {
		const windowName = "Window Selecter"

		// make win main function
		debug.Debugf("GetModuleHandle...")
		var hInstance = win.GetModuleHandle(nil)
		if hInstance == win.HINSTANCE(0) {
			debug.Debugln("err")
			result <- resultAttr{win.HWND(winapi.NULL), errors.Errorf("GetModuleHandle: Failed to get handler: %d\n", win.GetLastError())}
			return
		}
		debug.Debugln("ok")

		var className = winapi.MustUTF16PtrFromString(windowName)

		// lock os thread to avoid hanging GetMessage
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		var wClass win.WNDCLASSEX
		wClass = win.WNDCLASSEX{
			CbSize:    uint32(unsafe.Sizeof(wClass)),
			HInstance: hInstance,

			// redraw when the window size changed
			Style: win.CS_HREDRAW | win.CS_VREDRAW,

			// set window background color to white
			HbrBackground: win.HBRUSH(win.GetStockObject(win.WHITE_BRUSH)),
			LpszClassName: className,
			LpfnWndProc:   syscall.NewCallback(makeNewWindowProc(windows, indexChannel)),
			LpszMenuName:  nil,

			CbClsExtra: 0,
			CbWndExtra: 0,

			HIcon:   win.HICON(winapi.NULL),
			HCursor: win.HCURSOR(winapi.NULL),
			HIconSm: win.HICON(winapi.NULL),
		}

		debug.Debugf("RegisterClassEx...")
		if a := win.RegisterClassEx(&wClass); a == 0 {
			debug.Debugln("err")
			result <- resultAttr{win.HWND(winapi.NULL), errors.Errorf("RegisterClassEx: Failed to make window class \n")}
			return
		}
		debug.Debugln("ok")

		var windowNameUTF16 = winapi.MustUTF16PtrFromString(windowName)
		debug.Debugf("CreateWindowEx...")
		var hwnd = win.CreateWindowEx(
			win.WS_EX_OVERLAPPEDWINDOW,
			className,
			windowNameUTF16,
			win.WS_OVERLAPPEDWINDOW,
			win.CW_USEDEFAULT, win.CW_USEDEFAULT, int32(dialogWidth), int32(dialogHeight),
			win.HWND(winapi.NULL), win.HMENU(winapi.NULL), hInstance, nil,
		)

		if hwnd == win.HWND(winapi.NULL) {
			debug.Debugln("err")
			result <- resultAttr{hwnd, errors.Errorf("CreateWindowEx: Failed to make window")}
			return
		}
		debug.Debugln("ok")

		winapi.ShowWindow(hwnd, win.SW_SHOW)
		winapi.UpdateWindow(hwnd)

		result <- resultAttr{hwnd, nil}

		for {
			var msg win.MSG
			switch win.GetMessage(&msg, hwnd, 0, 0) {
			case 0:
				indexChannel <- -2
				return
			case -1:
				indexChannel <- -2
				return
			}

			win.TranslateMessage(&msg)
			win.DispatchMessage(&msg)
		}
	}()

	var res = <-result
	close(result)

	if res.err != nil {
		return 0, fmt.Errorf("WindowError: %w", res.err)
	}

	var index = <-indexChannel
	if index < 0 {
		return 0, ErrorQuitted
	}

	return windows[index].Hwnd, nil
}
