package winrt

import (
	"unsafe"

	"github.com/go-ole/go-ole"
)

var IAgileObjectID = ole.NewGUID("{94ea2b94-e9cc-49e0-c0ff-ee64ca8f5b90}")

type IAgileObject struct {
	ole.IUnknown
}

type IAgileObjectVtbl struct {
	ole.IUnknownVtbl
}

func (v *IAgileObject) VTable() *IAgileObjectVtbl {
	return (*IAgileObjectVtbl)(unsafe.Pointer(v.RawVTable))
}
