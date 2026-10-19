package client

import (
	"testing"
	"unsafe"

	"github.com/lxn/win"
)

func TestRawRelativeMotion(t *testing.T) {
	var motion rawMouseMotion
	for _, p := range []struct{ x, y int32 }{{0, 0}, {-1, 2}, {32768, -32769}, {-2147483648, 2147483647}} {
		x, y := motion.delta(1, win.RAWMOUSE{LLastX: p.x, LLastY: p.y}, 0, 0)
		if x != p.x || y != p.y {
			t.Fatalf("got %d,%d want %+v", x, y, p)
		}
	}
}

func TestRawAbsoluteMotion(t *testing.T) {
	var motion rawMouseMotion
	data := win.RAWMOUSE{UsFlags: win.MOUSE_MOVE_ABSOLUTE, LLastX: 32768, LLastY: 32768}
	if x, y := motion.delta(1, data, 1921, 1081); x != 0 || y != 0 {
		t.Fatal("first packet jumped")
	}
	data.LLastX, data.LLastY = 65535, 0
	if x, y := motion.delta(1, data, 1921, 1081); x != 960 || y != -540 {
		t.Fatalf("delta %d,%d", x, y)
	}
	// Device, display size, and coordinate-system changes establish a new baseline.
	if x, y := motion.delta(2, data, 1921, 1081); x != 0 || y != 0 {
		t.Fatal("device change jumped")
	}
	if x, y := motion.delta(2, data, 3841, 1081); x != 0 || y != 0 {
		t.Fatal("display change jumped")
	}
	data.UsFlags |= win.MOUSE_VIRTUAL_DESKTOP
	if x, y := motion.delta(2, data, 3841, 1081); x != 0 || y != 0 {
		t.Fatal("coordinate change jumped")
	}
	motion.delta(2, win.RAWMOUSE{LLastX: 1}, 0, 0)
	if x, y := motion.delta(2, data, 3841, 1081); x != 0 || y != 0 {
		t.Fatal("relative-to-absolute change jumped")
	}
	data.LLastX = -1
	if x, y := motion.delta(2, data, 3841, 1081); x != 0 || y != 0 || motion.valid {
		t.Fatal("invalid coordinate accepted")
	}
}

func TestRawMouseMemoryLayout(t *testing.T) {
	// GetRawInputData writes native Win32 structures into this fixed buffer.
	var packet win.RAWINPUTMOUSE
	if unsafe.Sizeof(packet.Data) != 24 || unsafe.Offsetof(packet.Data.LLastX) != 12 || unsafe.Offsetof(packet.Data.LLastY) != 16 {
		t.Fatal("unexpected RAWMOUSE ABI")
	}
	if unsafe.Offsetof(packet.Data) != unsafe.Sizeof(packet.Header) {
		t.Fatal("unexpected RAWINPUT ABI")
	}
}
