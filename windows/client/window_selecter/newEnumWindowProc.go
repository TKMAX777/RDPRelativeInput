package windowselecter

import (
	"github.com/TKMAX777/RDPRelativeInput/winapi"
	"github.com/lxn/win"
)

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
