package client

import "github.com/lxn/win"

type rawMouseMotion struct {
	device        win.HANDLE
	x, y          int32
	width, height int32
	virtual       bool
	valid         bool
}

func (m *rawMouseMotion) delta(device win.HANDLE, data win.RAWMOUSE, width, height int32) (int32, int32) {
	if data.UsFlags&win.MOUSE_MOVE_ABSOLUTE == 0 {
		m.valid = false
		return data.LLastX, data.LLastY
	}
	// Tablets and some remote/virtual devices report normalized absolute data.
	// Never interpret those coordinates as a huge relative motion.
	if width <= 0 || height <= 0 || data.LLastX < 0 || data.LLastX > 65535 || data.LLastY < 0 || data.LLastY > 65535 {
		m.valid = false
		return 0, 0
	}
	x := int32(int64(data.LLastX) * int64(width-1) / 65535)
	y := int32(int64(data.LLastY) * int64(height-1) / 65535)
	virtual := data.UsFlags&win.MOUSE_VIRTUAL_DESKTOP != 0
	dx, dy := int32(0), int32(0)
	if m.valid && m.device == device && m.width == width && m.height == height && m.virtual == virtual {
		dx, dy = x-m.x, y-m.y
	}
	*m = rawMouseMotion{device: device, x: x, y: y, width: width, height: height, virtual: virtual, valid: true}
	return dx, dy
}
