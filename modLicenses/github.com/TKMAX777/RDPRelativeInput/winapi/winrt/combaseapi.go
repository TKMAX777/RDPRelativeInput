package winrt

import (
	"unsafe"

	"github.com/go-ole/go-ole"
	"github.com/lxn/win"
	"golang.org/x/sys/windows"
)

var modole32 = windows.NewLazySystemDLL("Ole32.dll")

var pCoCreateFreeThreadedMarshaler = modole32.NewProc("CoCreateFreeThreadedMarshaler")

func CoCreateFreeThreadedMarshaler(punkOuter *ole.IUnknown) (*ole.IUnknown, error) {
	var res *ole.IUnknown
	r0, _, _ := pCoCreateFreeThreadedMarshaler.Call(uintptr(unsafe.Pointer(punkOuter)), uintptr(unsafe.Pointer(&res)))
	if r0 != win.S_OK {
		return nil, ole.NewError(r0)
	}
	return res, nil
}
