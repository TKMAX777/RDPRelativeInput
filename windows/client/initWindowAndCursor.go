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

	// get remote desktop client rect
	var rect win.RECT
	if !winapi.GetWindowRect(rdClientHwnd, &rect) {
		return errors.New("GetWindowRectError")
	}

	// get remote desktop client center position
	var windowCenterPosition = h.getWindowCenterPos(rect)

	// set cursot to center of the remote desktop client window
	if !win.SetCursorPos(windowCenterPosition.X, windowCenterPosition.Y) {
		return errors.New("Error in set cursor position")
	}

	// then clip cursor
	_, err := winapi.ClipCursor(&rect)
	if err != nil {
		return errors.Errorf("Error in clip cursor: code: %d\n", err)
	}

	winapi.ShowCursor(false)

	return nil
}
