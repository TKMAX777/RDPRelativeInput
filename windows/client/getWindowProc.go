package client

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/TKMAX777/RDPRelativeInput/debug"
	"github.com/TKMAX777/RDPRelativeInput/keymap"
	"github.com/TKMAX777/RDPRelativeInput/remote_send"
	"github.com/TKMAX777/RDPRelativeInput/winapi"
	"github.com/lxn/win"
)

func (h Handler) getWindowProc(rdClientHwnd win.HWND) func(hwnd win.HWND, uMsg uint32, wParam uintptr, lParam uintptr) uintptr {
	var isRelativeMode = true
	var rawMouse rawMouseMotion
	var rawPacket win.RAWINPUTMOUSE
	debug.Debugln("ToggleKey: ", h.options.toggleKey)
	debug.Debugln("ToggleType: ", h.options.toggleType)
	debug.Debugln("ExitKey: ", h.options.exitKey)

	toggleKey, err := keymap.GetWindowsKeyDetailFromEventInput(h.options.toggleKey)
	if err != nil {
		toggleKey, _ = keymap.GetWindowsKeyDetailFromEventInput("F8")
	}

	exitKey, err := keymap.GetWindowsKeyDetailFromEventInput(h.options.exitKey)
	if err != nil {
		exitKey, _ = keymap.GetWindowsKeyDetailFromEventInput("F12")
	}

	return func(hwnd win.HWND, uMsg uint32, wParam uintptr, lParam uintptr) uintptr {
		var send = func(evType keymap.EV_TYPE, key uint32, state remote_send.InputType) {
			if evType == keymap.EV_TYPE_KEY && key == exitKey.Value {
				switch state {
				case remote_send.KeyDown:
					win.SetTimer(hwnd, 1, 500, 0)
				case remote_send.KeyUp:
					win.KillTimer(hwnd, 1)
				}
			}

			if evType == keymap.EV_TYPE_KEY && key == toggleKey.Value {
				rawMouse = rawMouseMotion{}
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

		switch uMsg {
		case win.WM_CREATE:
			device := win.RAWINPUTDEVICE{UsUsagePage: 1, UsUsage: 2, HwndTarget: hwnd}
			// Foreground only. Keep legacy button/wheel messages, but ignore
			// legacy movement to avoid sending each movement twice.
			if !win.RegisterRawInputDevices(&device, 1, uint32(unsafe.Sizeof(device))) {
				debug.Debugln("RegisterRawInputDevices failed")
				return ^uintptr(0)
			}
			winapi.SetLayeredWindowAttributes(hwnd, 0x0000FF, byte(1), winapi.LWA_COLORKEY)
			h.initWindowAndCursor(hwnd, rdClientHwnd)
			win.UpdateWindow(hwnd)

			return winapi.NULL
		case win.WM_DESTROY:
			h.Close()
			h.remote.SendExit()
			os.Exit(0)
			return winapi.NULL
		case win.WM_PAINT:
			var ps = new(win.PAINTSTRUCT)
			var hdc = win.BeginPaint(hwnd, ps)
			var hBrush = winapi.CreateSolidBrush(0x000000FF)

			oldBrush := win.SelectObject(hdc, win.HGDIOBJ(hBrush))
			winapi.ExtFloodFill(hdc, 1, 1, 0xFFFFFF, winapi.FLOODFILLSURFACE)
			win.SelectObject(hdc, oldBrush)
			win.DeleteObject(hBrush)
			win.EndPaint(hwnd, ps)

			return winapi.NULL
		case win.WM_TIMER:
			h.Close()
			h.remote.SendExit()
			os.Exit(0)
			return winapi.NULL
		case win.WM_INPUT:
			size := uint32(unsafe.Sizeof(rawPacket))
			n := win.GetRawInputData(win.HRAWINPUT(lParam), win.RID_INPUT, unsafe.Pointer(&rawPacket), &size, uint32(unsafe.Sizeof(rawPacket.Header)))
			if n != ^uint32(0) && n >= uint32(unsafe.Sizeof(rawPacket)) && rawPacket.Header.DwType == win.RIM_TYPEMOUSE {
				if isRelativeMode && win.GetForegroundWindow() == hwnd {
					width, height := int32(0), int32(0)
					if rawPacket.Data.UsFlags&win.MOUSE_MOVE_ABSOLUTE != 0 {
						width, height = win.GetSystemMetrics(win.SM_CXSCREEN), win.GetSystemMetrics(win.SM_CYSCREEN)
						if rawPacket.Data.UsFlags&win.MOUSE_VIRTUAL_DESKTOP != 0 {
							width, height = win.GetSystemMetrics(win.SM_CXVIRTUALSCREEN), win.GetSystemMetrics(win.SM_CYVIRTUALSCREEN)
						}
					}
					dx, dy := rawMouse.delta(rawPacket.Header.HDevice, rawPacket.Data, width, height)
					if dx != 0 || dy != 0 {
						h.remote.SendRelativeCursor(dx, dy)
					}
				} else {
					rawMouse = rawMouseMotion{}
				}
			}
			// Required by WM_INPUT for foreground packet cleanup, even if ignored.
			win.DefWindowProc(hwnd, uMsg, wParam, lParam)
			return 0
		case win.WM_KILLFOCUS:
			rawMouse = rawMouseMotion{}
			winapi.ClipCursor(nil)
			return win.DefWindowProc(hwnd, uMsg, wParam, lParam)
		case win.WM_SETFOCUS:
			rawMouse = rawMouseMotion{}
			if isRelativeMode {
				h.clipRelativeCursor(rdClientHwnd)
			}
			return win.DefWindowProc(hwnd, uMsg, wParam, lParam)
		case win.WM_MOUSEMOVE:
			return 0
		case win.WM_SYSKEYDOWN:
			if lParam>>30&1 == 0 {
				send(keymap.EV_TYPE_KEY, uint32(wParam), remote_send.KeyDown)
			}
			return winapi.NULL
		case win.WM_SYSKEYUP:
			send(keymap.EV_TYPE_KEY, uint32(wParam), remote_send.KeyUp)
			return winapi.NULL
		case win.WM_KEYDOWN:
			if lParam>>30&1 == 0 || (uint32(wParam) != toggleKey.Value && uint32(wParam) != exitKey.Value) {
				send(keymap.EV_TYPE_KEY, uint32(wParam), remote_send.KeyDown)
			}
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
