package host

import (
	"fmt"
	"os"

	"github.com/TKMAX777/RDPRelativeInput/winapi"
	"github.com/lxn/win"
)

func (h *Handler) getWindowProc() func(hwnd win.HWND, uMsg uint32, wParam uintptr, lParam uintptr) uintptr {
	return func(hwnd win.HWND, uMsg uint32, wParam uintptr, lParam uintptr) uintptr {
		switch uMsg {
		case win.WM_CREATE:
			win.CreateWindowEx(0, winapi.MustUTF16PtrFromString("BUTTON"), winapi.MustUTF16PtrFromString("o"),
				win.WS_CHILD|win.WS_VISIBLE|win.BS_FLAT,
				0, 0, 50, 50,
				hwnd, 1, win.GetModuleHandle(nil), nil)
			win.UpdateWindow(hwnd)

			return winapi.NULL
		case win.WM_DESTROY:
			os.Exit(0)
			return winapi.NULL
		case win.WM_COMMAND:
			switch win.LOWORD(uint32(wParam)) {
			case 1:
				if h.isCapturing {
					h.captureHandler.Close()
					h.isCapturing = false
				} else {
					h.captureHandler = &CaptureHandler{}

					err := h.captureHandler.StartCapture(hwnd)
					if err != nil {
						fmt.Println(err)
					}
					h.isCapturing = true
				}
			}
			return win.DefWindowProc(hwnd, uMsg, wParam, lParam)
		case win.WM_PAINT:
			return win.DefWindowProc(hwnd, uMsg, wParam, lParam)
		default:
			return win.DefWindowProc(hwnd, uMsg, wParam, lParam)
		}
	}
}

func (h Handler) Close() {
	winapi.ClipCursor(nil)
	winapi.ShowCursor(true)
}
