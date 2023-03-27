package client

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
// rdClientHwnd must be remote desktop client hwnd, and toggleKey is a keyname for toggle wrapper mode
func (h Handler) StartClient(rdClientHwnd win.HWND) (win.HWND, error) {
	if rdClientHwnd == win.HWND(winapi.NULL) {
		return win.HWND(winapi.NULL), errors.New("NilWindowHandler")
	}

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

		// get window proc
		var windowProc = h.getWindowProc(rdClientHwnd)

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
			LpfnWndProc:   syscall.NewCallback(windowProc),
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
			win.WS_EX_OVERLAPPEDWINDOW|win.WS_EX_TOPMOST|win.WS_EX_LAYERED,
			className,
			windowNameUTF16,
			win.WS_OVERLAPPEDWINDOW,
			win.CW_USEDEFAULT, win.CW_USEDEFAULT, int32(100), int32(100),
			win.HWND(winapi.NULL), win.HMENU(winapi.NULL), hInstance, unsafe.Pointer(nil),
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
			// debug.Debugf("disp: %+v \n", msg)
		}
	}()

	var res = <-result
	close(result)

	return res.hwnd, res.err
}
