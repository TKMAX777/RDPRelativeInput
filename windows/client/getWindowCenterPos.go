package client

import "github.com/lxn/win"

func (h Handler) getWindowCenterPos(rect win.RECT) win.POINT {
	var windowCenterPosition win.POINT
	windowCenterPosition.X = int32(rect.Left+rect.Right) / 2
	windowCenterPosition.Y = int32(rect.Top+rect.Bottom) / 2

	return windowCenterPosition
}
