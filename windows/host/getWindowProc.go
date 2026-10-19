package host

import (
	"os"

	"github.com/TKMAX777/RDPRelativeInput/winapi"
	"github.com/lxn/win"
)

func (h *Handler) getWindowProc() func(hwnd win.HWND, uMsg uint32, wParam uintptr, lParam uintptr) uintptr {
	return func(hwnd win.HWND, uMsg uint32, wParam uintptr, lParam uintptr) uintptr {
		switch uMsg {
		case win.WM_CREATE:
			h.button = win.CreateWindowEx(0, winapi.MustUTF16PtrFromString("BUTTON"), winapi.MustUTF16PtrFromString("Cursor: Auto"),
				win.WS_CHILD|win.WS_VISIBLE|win.BS_FLAT,
				0, 0, 178, 38,
				hwnd, 1, win.GetModuleHandle(nil), nil)
			win.UpdateWindow(hwnd)

			return winapi.NULL
		case win.WM_DESTROY:
			if h.overlay != nil {
				h.overlay.close()
			}
			os.Exit(0)
			return winapi.NULL
		case win.WM_COMMAND:
			switch win.LOWORD(uint32(wParam)) {
			case 1:
				if h.overlay != nil {
					h.overlay.mode = (h.overlay.mode + 1) % 3
					labels := [...]string{"Cursor: Auto", "Cursor: Always", "Cursor: Off"}
					winapi.SetWindowText(h.button, winapi.MustUTF16PtrFromString(labels[h.overlay.mode]))
					h.overlay.refresh()
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
