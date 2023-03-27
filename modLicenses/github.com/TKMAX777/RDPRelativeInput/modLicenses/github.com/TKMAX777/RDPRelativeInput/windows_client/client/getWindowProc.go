package client

import (
	"fmt"
	"log"
	"os"

	"github.com/TKMAX777/RDPRelativeInput/debug"
	"github.com/TKMAX777/RDPRelativeInput/keymap"
	"github.com/TKMAX777/RDPRelativeInput/remote_send"
	"github.com/TKMAX777/RDPRelativeInput/winapi"
	"github.com/lxn/win"
)

func (h Handler) getWindowProc(rdClientHwnd win.HWND) func(hwnd win.HWND, uMsg uint32, wParam uintptr, lParam uintptr) uintptr {
	// get remote desktop client rect
	var rect win.RECT
	if !winapi.GetWindowRect(rdClientHwnd, &rect) {
		fmt.Fprintf(os.Stderr, "getWindowProc: GetWindowRectError")
	}

	var isRelativeMode = true

	// get remote desktop client center position
	var windowCenterPosition = h.getWindowCenterPos(rect)

	var currentPosition = windowCenterPosition
	debug.Debugln("ToggleKey: ", h.options.toggleKey)
	debug.Debugln("ToggleType: ", h.options.toggleType)

	toggleKey, err := keymap.GetWindowsKeyDetailFromEventInput(h.options.toggleKey)
	if err != nil {
		toggleKey, _ = keymap.GetWindowsKeyDetailFromEventInput("F8")
	}

	return func(hwnd win.HWND, uMsg uint32, wParam uintptr, lParam uintptr) uintptr {
		// var send = func(key keymap.WindowsKey, state InputType) {
		var send = func(evType keymap.EV_TYPE, key uint32, state remote_send.InputType) {
			// debug.Debugf("Send: key: %s state: %+v\n", key.EventInput, state)

			if key == toggleKey.Value {
				// toggle window mode
				switch state {
				case remote_send.KeyDown:
					isRelativeMode = !isRelativeMode || h.options.toggleType == ToggleTypeOnce
				case remote_send.KeyUp:
					if h.options.toggleType == ToggleTypeOnce {
						isRelativeMode = false
					}
				}

				debug.Debugf("Set: isRelativeMode: %t\n", isRelativeMode)

				if isRelativeMode {
					h.initWindowAndCursor(hwnd, rdClientHwnd)
					debug.Debugln("Called: initWindowAndCursor")
				} else {
					var crectAbs win.RECT
					if !winapi.GetWindowRect(rdClientHwnd, &crectAbs) {
						fmt.Fprintf(os.Stderr, "GetWindowRectError")
					}

					// show window title only
					if !win.SetWindowPos(hwnd, 0,
						crectAbs.Left, crectAbs.Top, crectAbs.Right-crectAbs.Left, h.metrics.TitleHeight+h.metrics.FrameWidthY*2,
						win.SWP_SHOWWINDOW,
					) {
						fmt.Fprintf(os.Stderr, "SetWindowPos: failed to set window pos")
					}

					debug.Debugln("Called: SetWindowPos")
					winapi.ShowCursor(true)
					winapi.ClipCursor(nil)
				}
				debug.Debugln("ModeChangeDone")
			}
			if isRelativeMode {
				h.remote.SendInput(evType, key, state)
			}
		}

		// debug.Debugf("Msg: %X(%d)\n", uMsg, uMsg)

		switch uMsg {
		case win.WM_CREATE:
			winapi.SetLayeredWindowAttributes(hwnd, 0x0000FF, byte(1), winapi.LWA_COLORKEY)
			h.initWindowAndCursor(hwnd, rdClientHwnd)
			win.UpdateWindow(hwnd)

			return winapi.NULL
		case win.WM_DESTROY:
			os.Stderr.Write([]byte("CLOSE\n"))
			os.Exit(0)
			return winapi.NULL
		case win.WM_PAINT:
			var ps = new(win.PAINTSTRUCT)
			var hdc = win.BeginPaint(hwnd, ps)
			var hBrush = winapi.CreateSolidBrush(0x000000FF)

			win.SelectObject(hdc, hBrush)
			winapi.ExtFloodFill(hdc, 1, 1, 0xFFFFFF, winapi.FLOODFILLSURFACE)
			win.DeleteObject(hBrush)
			win.EndPaint(hwnd, ps)

			if isRelativeMode {
				winapi.SetLayeredWindowAttributes(hwnd, 0x0000FF, byte(1), winapi.LWA_COLORKEY)

				// get remote desktop client rect
				var rect win.RECT
				if !winapi.GetWindowRect(rdClientHwnd, &rect) {
					fmt.Fprintf(os.Stderr, "getWindowProc: GetWindowRectError")
				}

				// get remote desktop client center position
				windowCenterPosition = h.getWindowCenterPos(rect)
			}
			return winapi.NULL
		case win.WM_MOUSEMOVE:
			if isRelativeMode {
				ok := winapi.GetCursorPos(&currentPosition)
				if !ok {
					log.Printf("Error in get cursor position\n")
				}
				// Relative position mode
				var point = win.POINT{X: currentPosition.X - windowCenterPosition.X, Y: currentPosition.Y - windowCenterPosition.Y}
				if point.X != 0 || point.Y != 0 {
					debug.Debugf("X: %4d Y: %4d\n", point.X, point.Y)
					h.remote.SendRelativeCursor(point.X, point.Y)

					ok = win.SetCursorPos(windowCenterPosition.X, windowCenterPosition.Y)
					if !ok {
						log.Printf("Error in set cursor position")
					}
				}
			}

			return winapi.NULL
		case win.WM_SYSKEYDOWN:
			if lParam>>30&1 == 0 {
				send(keymap.EV_TYPE_KEY, uint32(wParam), remote_send.KeyDown)
			}
			return winapi.NULL
		case win.WM_SYSKEYUP:
			send(keymap.EV_TYPE_KEY, uint32(wParam), remote_send.KeyUp)
			return winapi.NULL
		case win.WM_KEYDOWN:
			send(keymap.EV_TYPE_KEY, uint32(wParam), remote_send.KeyDown)
			return winapi.NULL
		case win.WM_KEYUP:
			send(keymap.EV_TYPE_KEY, uint32(wParam), remote_send.KeyUp)
			return winapi.NULL
		case win.WM_MOUSEWHEEL:
			if int16(wParam>>16) > 0 {
				send(keymap.EV_TYPE_WHEEL, uint32(int16(wParam>>16)), remote_send.KeyUp)
			} else {
				send(keymap.EV_TYPE_WHEEL, uint32(int16(wParam>>16)), remote_send.KeyDown)
			}
			return winapi.NULL
		case win.WM_LBUTTONUP:
			send(keymap.EV_TYPE_MOUSE, 0x01, remote_send.KeyUp)
			return winapi.NULL
		case win.WM_LBUTTONDOWN:
			send(keymap.EV_TYPE_MOUSE, 0x01, remote_send.KeyDown)
			return winapi.NULL
		case win.WM_RBUTTONUP:
			send(keymap.EV_TYPE_MOUSE, 0x02, remote_send.KeyUp)
			return winapi.NULL
		case win.WM_RBUTTONDOWN:
			send(keymap.EV_TYPE_MOUSE, 0x02, remote_send.KeyDown)
			return winapi.NULL
		case win.WM_MBUTTONUP:
			send(keymap.EV_TYPE_MOUSE, 0x04, remote_send.KeyUp)
			return winapi.NULL
		case win.WM_MBUTTONDOWN:
			send(keymap.EV_TYPE_MOUSE, 0x04, remote_send.KeyDown)
			return winapi.NULL
		default:
			return win.DefWindowProc(hwnd, uMsg, wParam, lParam)
		}
	}
}

func (h Handler) Close() {
	winapi.ClipCursor(nil)
	winapi.ShowCursor(true)
}
