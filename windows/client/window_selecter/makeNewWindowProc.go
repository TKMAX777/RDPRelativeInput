package windowselecter

import (
	"unsafe"

	"github.com/TKMAX777/RDPRelativeInput/winapi"
	"github.com/lxn/win"
)

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
			case 2:
				index := int(win.SendMessage(hListBox, win.LB_GETCURSEL, 0, 0))
				if index < 0 {
					return win.DefWindowProc(hwnd, uMsg, wParam, lParam)
				}
				indexChannel <- index
				win.DestroyWindow(hwnd)
			}
		}
		return win.DefWindowProc(hwnd, uMsg, wParam, lParam)
	}
}
