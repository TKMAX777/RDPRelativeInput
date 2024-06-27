package windowselecter

import (
	"fmt"
	"os"
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
				debug.Debugln("Quit")
				return
			case -1:
				os.Stderr.Write([]byte("CLOSE\n"))
				os.Exit(0)
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

	return windows[index].Hwnd, nil
}

func newEnumWindowProc(windows *[]WindowInfo) func(hwnd win.HWND, lParam *uintptr) uintptr {
	return func(hwnd win.HWND, lParam *uintptr) uintptr {
		if !win.IsWindowVisible(hwnd) {
			return 1
		}
		var title = winapi.GetWindowTextString(hwnd)
		if title == "" {
			return 1
		}

		*windows = append(*windows, WindowInfo{hwnd, title})

		return 1
	}
}

func makeNewWindowProc(windows []WindowInfo, indexChannel chan int) winapi.Dlgproc {
	var hListBox win.HWND
	var ps win.PAINTSTRUCT

	return func(hwnd win.HWND, uMsg uint32, wParam uintptr, lParam uintptr) uintptr {
		switch uMsg {
		case win.WM_PAINT:
			hdc := win.BeginPaint(hwnd, &ps)
			const windowText = "Please select the target window..."
			win.TextOut(hdc, dialogBaseXPos, dialogBaseYPos, winapi.MustUTF16PtrFromString(windowText), int32(len(windowText)))
			win.EndPaint(hwnd, &ps)
		case win.WM_CREATE:
			var ycoordinate int32 = dialogBaseYPos + 20
			var height = int32(float32(dialogHeight-win.GetSystemMetrics(win.SM_CYCAPTION))*0.9) - 30

			hListBox = win.CreateWindowEx(
				win.WS_EX_CLIENTEDGE,
				winapi.MustUTF16PtrFromString("LISTBOX"),
				nil,
				win.WS_CHILD|win.WS_VISIBLE|win.WS_VSCROLL|win.ES_AUTOVSCROLL|win.LBS_NOTIFY,
				dialogBaseXPos, ycoordinate, dialogWidth*0.9, height,
				hwnd, 1, win.GetModuleHandle(nil),
				nil,
			)

			ycoordinate += height - 10

			for _, w := range windows {
				win.SendMessage(hListBox, win.LB_ADDSTRING, 0, uintptr(unsafe.Pointer(winapi.MustUTF16PtrFromString(w.WindowTitle))))
			}

			win.CreateWindowEx(0, winapi.MustUTF16PtrFromString("BUTTON"), winapi.MustUTF16PtrFromString("OK"),
				win.WS_CHILD|win.WS_VISIBLE|win.BS_FLAT,
				dialogBaseXPos, ycoordinate, dialogWidth*0.9, 20,
				hwnd, 2, win.GetModuleHandle(nil), nil)
			win.UpdateWindow(hwnd)
		case win.WM_COMMAND:
			switch win.LOWORD(uint32(wParam)) {
			// case win.LBN_SELCHANGE:
			case 2:
				indexChannel <- int(win.SendMessage(hListBox, win.LB_GETCURSEL, 0, 0))
				win.DestroyWindow(hwnd)
			}
		}
		return win.DefWindowProc(hwnd, uMsg, wParam, lParam)
	}
}
