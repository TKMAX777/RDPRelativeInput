package winapi

import (
	"unsafe"

	"github.com/lxn/win"
)

const (
	FLOODFILLBORDER uint32 = iota
	FLOODFILLSURFACE
)

func CreateRectRgnIndirect(rect win.RECT) win.HRGN {
	return win.HRGN(createRectRgnIndirect(uintptr(unsafe.Pointer(&rect.Left))))
}

func ExtFloodFill(hdc win.HDC, x int, y int, color uint32, opType uint32) error {
	return extFloodFill(uintptr(hdc), x, y, color, opType)
}

func CreateSolidBrush(color uint32) win.HGDIOBJ {
	return win.HGDIOBJ(createSolidBrush(color))
}

func CreatePen(iStyle int, cWidth int, color uint32) win.HPEN {
	return win.HPEN(createPen(iStyle, cWidth, color))
}

func PolyDraw(hdc win.HDC, apt win.POINT, aj byte, cpt int) error {
	return polyDraw(uintptr(hdc), uintptr(unsafe.Pointer(&apt.X)), uintptr(aj), cpt)
}

func CreateDIBSection(hdc win.HDC, pbmi *win.BITMAPINFO, usage uint, ppvBits uintptr, hSection win.HANDLE, offset uint32) win.HBITMAP {
	return win.HBITMAP(createDIBSection(uintptr(hdc), uintptr(unsafe.Pointer(&pbmi.BmiHeader.BiSize)), usage, ppvBits, uintptr(hSection), offset))
}
