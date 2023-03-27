package winrt

import (
	"syscall"
	"unsafe"

	"github.com/go-ole/go-ole"
	"github.com/lxn/win"
)

// IGraphicsCptureItemInterop
// https://learn.microsoft.com/en-us/windows/win32/api/windows.graphics.capture.interop/nn-windows-graphics-capture-interop-igraphicscaptureiteminterop

var IGraphicsCaptureItemInteropID = ole.NewGUID("{3628E81B-3CAC-4C60-B7F4-23CE0E0C3356}")

type IGraphicsCaptureItemInterop struct {
	ole.IUnknown
}

type IGraphicsCaptureItemInteropVtabl struct {
	ole.IUnknownVtbl
	CreateForWindow  uintptr
	CreateForMonitor uintptr
}

func (v *IGraphicsCaptureItemInterop) VTable() *IGraphicsCaptureItemInteropVtabl {
	return (*IGraphicsCaptureItemInteropVtabl)(unsafe.Pointer(v.RawVTable))
}

func (v *IGraphicsCaptureItemInterop) CreateForWindow(window win.HWND, guid *ole.GUID, result **ole.IInspectable) error {
	r1, _, _ := syscall.SyscallN(v.VTable().CreateForWindow, uintptr(unsafe.Pointer(v)), uintptr(window), uintptr(unsafe.Pointer(&guid.Data1)), uintptr(unsafe.Pointer(result)))
	if r1 != win.S_OK {
		return ole.NewError(r1)
	}
	return nil
}

func (v *IGraphicsCaptureItemInterop) CreateForMonitor(monitor win.HMONITOR, guid *ole.GUID, result **ole.IInspectable) error {
	r1, _, _ := syscall.SyscallN(v.VTable().CreateForMonitor, uintptr(unsafe.Pointer(v)), uintptr(monitor), uintptr(unsafe.Pointer(&guid.Data1)), uintptr(unsafe.Pointer(result)))
	if r1 != win.S_OK {
		return ole.NewError(r1)
	}
	return nil
}
