package host

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/TKMAX777/RDPRelativeInput/winapi"
	"github.com/lxn/win"
	"golang.org/x/sys/windows"
)

var cursorUser32 = windows.NewLazySystemDLL("user32.dll")
var getCursorInfo = cursorUser32.NewProc("GetCursorInfo")
var copyCursorIcon = cursorUser32.NewProc("CopyIcon")
var setThreadDPI = cursorUser32.NewProc("SetThreadDpiAwarenessContext")
var unregisterCursorClass = cursorUser32.NewProc("UnregisterClassW")

type cursorInfo struct {
	Size, Flags uint32
	Cursor      win.HCURSOR
	Position    win.POINT
}

// All overlay state belongs to the host's locked UI thread.
type cursorOverlay struct {
	className *uint16
	hwnd      win.HWND
	icon      win.HICON
	brush     win.HBRUSH
	mode      int
}

func cursorVisible(mode int, flags uint32) bool {
	return mode == 1 || (mode == 0 && flags&1 != 0)
}

func usePhysicalCursorCoordinates() func() {
	if setThreadDPI.Find() == nil {
		// DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2 = -4.
		previous, _, _ := setThreadDPI.Call(^uintptr(3))
		if previous != 0 {
			return func() { setThreadDPI.Call(previous) }
		}
	}
	return func() {}
}

func newCursorOverlay() (*cursorOverlay, error) {
	o := &cursorOverlay{brush: win.HBRUSH(winapi.CreateSolidBrush(0x00ff00ff))}
	if o.brush == 0 {
		return nil, fmt.Errorf("create cursor background brush")
	}
	name := winapi.MustUTF16PtrFromString("RDPRelativeInput.CursorOverlay")
	wc := win.WNDCLASSEX{HInstance: win.GetModuleHandle(nil), LpszClassName: name,
		LpfnWndProc: syscall.NewCallback(o.windowProc)}
	wc.CbSize = uint32(unsafe.Sizeof(wc))
	if win.RegisterClassEx(&wc) == 0 {
		o.close()
		return nil, fmt.Errorf("register cursor overlay class")
	}
	o.className = name
	o.hwnd = win.CreateWindowEx(win.WS_EX_LAYERED|win.WS_EX_TRANSPARENT|win.WS_EX_NOACTIVATE|win.WS_EX_TOOLWINDOW|win.WS_EX_TOPMOST,
		name, name, win.WS_POPUP, 0, 0, 1, 1, 0, 0, wc.HInstance, nil)
	if o.hwnd == 0 {
		o.close()
		return nil, fmt.Errorf("create cursor overlay window")
	}
	if err := winapi.SetLayeredWindowAttributes(o.hwnd, 0x00ff00ff, 255, winapi.LWA_COLORKEY); err != nil {
		o.close()
		return nil, err
	}
	if win.SetTimer(o.hwnd, 1, 16, 0) == 0 {
		o.close()
		return nil, fmt.Errorf("start cursor overlay timer")
	}
	return o, nil
}

func (o *cursorOverlay) close() {
	if o.hwnd != 0 {
		win.KillTimer(o.hwnd, 1)
		win.DestroyWindow(o.hwnd)
		o.hwnd = 0
	}
	if o.icon != 0 {
		win.DestroyIcon(o.icon)
		o.icon = 0
	}
	if o.brush != 0 {
		win.DeleteObject(win.HGDIOBJ(o.brush))
		o.brush = 0
	}
	if o.className != nil {
		// The library wrapper omits the required HINSTANCE argument.
		unregisterCursorClass.Call(uintptr(unsafe.Pointer(o.className)), uintptr(win.GetModuleHandle(nil)))
		o.className = nil
	}
}

func (o *cursorOverlay) hide() { win.ShowWindow(o.hwnd, win.SW_HIDE) }

func (o *cursorOverlay) refresh() {
	if o.mode == 2 {
		o.hide()
		return
	}
	var ci cursorInfo
	ci.Size = uint32(unsafe.Sizeof(ci))
	ok, _, _ := getCursorInfo.Call(uintptr(unsafe.Pointer(&ci)))
	if ok == 0 || ci.Cursor == 0 || !cursorVisible(o.mode, ci.Flags) {
		o.hide()
		return
	}
	// Copy each frame: cursor handles can be reused or their image changed in place.
	icon, _, _ := copyCursorIcon.Call(uintptr(ci.Cursor))
	if icon == 0 {
		o.hide()
		return
	}
	var ii win.ICONINFO
	if !win.GetIconInfo(win.HICON(icon), &ii) {
		win.DestroyIcon(win.HICON(icon))
		o.hide()
		return
	}
	if ii.HbmColor != 0 {
		defer win.DeleteObject(win.HGDIOBJ(ii.HbmColor))
	}
	if ii.HbmMask != 0 {
		defer win.DeleteObject(win.HGDIOBJ(ii.HbmMask))
	}
	bitmap := ii.HbmColor
	if bitmap == 0 {
		bitmap = ii.HbmMask
	}
	var bm win.BITMAP
	n := win.GetObject(win.HGDIOBJ(bitmap), unsafe.Sizeof(bm), unsafe.Pointer(&bm))
	if ii.HbmColor == 0 {
		bm.BmHeight /= 2
	}
	if n == 0 || bm.BmWidth <= 0 || bm.BmHeight <= 0 || bm.BmWidth > 512 || bm.BmHeight > 512 || ii.XHotspot >= uint32(bm.BmWidth) || ii.YHotspot >= uint32(bm.BmHeight) {
		win.DestroyIcon(win.HICON(icon))
		o.hide()
		return
	}
	if o.icon != 0 {
		win.DestroyIcon(o.icon)
	}
	o.icon = win.HICON(icon)
	win.SetWindowPos(o.hwnd, win.HWND_TOPMOST, ci.Position.X-int32(ii.XHotspot), ci.Position.Y-int32(ii.YHotspot), bm.BmWidth, bm.BmHeight, win.SWP_NOACTIVATE|win.SWP_SHOWWINDOW)
	win.InvalidateRect(o.hwnd, nil, false)
	win.UpdateWindow(o.hwnd)
}

func (o *cursorOverlay) windowProc(hwnd win.HWND, msg uint32, wp, lp uintptr) uintptr {
	switch msg {
	case win.WM_TIMER:
		o.refresh()
		return 0
	case win.WM_NCHITTEST:
		return ^uintptr(0) // HTTRANSPARENT
	case win.WM_MOUSEACTIVATE:
		return 3 // MA_NOACTIVATE
	case win.WM_PAINT:
		var ps win.PAINTSTRUCT
		dc := win.BeginPaint(hwnd, &ps)
		var rect win.RECT
		win.GetClientRect(hwnd, &rect)
		winapi.FillRect(dc, rect, o.brush)
		if o.icon != 0 {
			win.DrawIconEx(dc, 0, 0, o.icon, rect.Right, rect.Bottom, 0, 0, 3)
		}
		win.EndPaint(hwnd, &ps)
		return 0
	}
	return win.DefWindowProc(hwnd, msg, wp, lp)
}
