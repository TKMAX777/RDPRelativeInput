package host

import (
	"os"
	"runtime"
	"syscall"
	"unsafe"

	"github.com/TKMAX777/RDPRelativeInput/debug"
	"github.com/TKMAX777/RDPRelativeInput/winapi"
	"github.com/lxn/win"
	"github.com/pkg/errors"
)

// Create window on remote desktop client
func (h Handler) StartHostButton() (win.HWND, error) {
	type resultAttr struct {
		hwnd win.HWND
		err  error
	}

	var result = make(chan resultAttr)

	go func() {
		const windowName = "RDP Input Wrapper"

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
			LpfnWndProc:   syscall.NewCallback(h.getWindowProc()),
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
			0,
			className,
			windowNameUTF16,
			0,
			50, 50, int32(50), int32(50),
			win.HWND(winapi.NULL), win.HMENU(winapi.NULL), hInstance, unsafe.Pointer(nil),
		)

		if hwnd == win.HWND(winapi.NULL) {
			debug.Debugln("err")
			result <- resultAttr{hwnd, errors.Errorf("CreateWindowEx: Failed to make window")}
			return
		}
		debug.Debugln("ok")

		var style uint32 = uint32(win.WS_POPUP | win.WS_BORDER)

		win.SetWindowLong(hwnd, win.GWL_STYLE, *(*int32)(unsafe.Pointer(&style)))
		win.SetWindowPos(hwnd, 0, 50, 50, 50, 50, 0)
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

	return res.hwnd, res.err
}
