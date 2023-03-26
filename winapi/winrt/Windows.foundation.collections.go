package winrt

import (
	"unsafe"

	"github.com/go-ole/go-ole"
)

// ITypedEventHandler
var ITypedEventHandlerID = ole.NewGUID("{51A947F7-79CF-5A3E-A3A5-1289CFA6DFE8}")

type ITypedEventHandler struct {
	ole.IUnknown
}

type ITypedEventHandlerVtbl struct {
	ole.IUnknownVtbl
}

func (v *ITypedEventHandler) VTable() *ITypedEventHandlerVtbl {
	return (*ITypedEventHandlerVtbl)(unsafe.Pointer(v.RawVTable))
}
