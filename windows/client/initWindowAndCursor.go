package client

import (
	"github.com/TKMAX777/RDPRelativeInput/winapi"
	"github.com/lxn/win"
	"github.com/pkg/errors"
)

func (h Handler) initWindowAndCursor(hwnd, rdClientHwnd win.HWND) error {
	// Set window on RDP client window
	if !win.SetForegroundWindow(hwnd) {
		return errors.New("SetForegroundWindow: failed to get foreground permission")
	}

	var crectAbs win.RECT
	if !winapi.GetWindowRect(rdClientHwnd, &crectAbs) {
		return errors.New("GetWindowRectError")
	}

	if !win.SetWindowPos(hwnd, 0,
		crectAbs.Left, crectAbs.Top, crectAbs.Right-crectAbs.Left, crectAbs.Bottom-crectAbs.Top,
		win.SWP_SHOWWINDOW,
	) {
		return errors.New("SetWindowPos: failed to set window pos")
	}

	if err := h.clipRelativeCursor(rdClientHwnd); err != nil {
		return err
	}

	winapi.ShowCursor(false)

	return nil
}

// Raw motion keeps reporting device deltas while the OS cursor is stationary.
// Keep the cursor in the wrapper client area without warping it for each event.
func (h Handler) clipRelativeCursor(rdClientHwnd win.HWND) error {
	var rect win.RECT
	if !winapi.GetWindowRect(rdClientHwnd, &rect) {
		return errors.New("GetWindowRect failed")
	}
	center := h.getWindowCenterPos(rect)
	clip := win.RECT{Left: center.X, Top: center.Y, Right: center.X + 1, Bottom: center.Y + 1}
	if _, err := winapi.ClipCursor(&clip); err != nil {
		return err
	}
	return nil
}
