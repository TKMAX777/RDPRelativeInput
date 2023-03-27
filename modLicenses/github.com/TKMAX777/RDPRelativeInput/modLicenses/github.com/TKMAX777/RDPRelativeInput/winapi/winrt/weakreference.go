package winrt

import (
	"syscall"
	"unsafe"

	"github.com/go-ole/go-ole"
	"github.com/lxn/win"
)

var IWeakReferenceSourceID = ole.NewGUID("{00000038-0000-0000-C000-000000000046}")

type IWeakReferenceSource struct {
	ole.IInspectable
}

type IWeakReferenceSourceVtbl struct {
	ole.IInspectableVtbl
	GetWeakReference uintptr
}

func (v *IWeakReferenceSource) VTable() *IWeakReferenceSourceVtbl {
	return (*IWeakReferenceSourceVtbl)(unsafe.Pointer(v.RawVTable))
}

func (v *IWeakReferenceSource) GetWeakReference() (*IWeakReference, error) {
	var ret *IWeakReference
	r1, _, _ := syscall.SyscallN(v.VTable().GetWeakReference, uintptr(unsafe.Pointer(v)), uintptr(unsafe.Pointer(&ret)))
	if r1 != win.S_OK {
		return ret, ole.NewError(r1)
	}
	return ret, nil
}

var IWeakReferenceID = ole.NewGUID("{00000037-0000-0000-C000-000000000046}")

type IWeakReference struct {
	ole.IUnknown
}

type IWeakReferenceVtbl struct {
	ole.IUnknownVtbl
	Resolve uintptr
}

func (v *IWeakReference) VTable() *IWeakReferenceVtbl {
	return (*IWeakReferenceVtbl)(unsafe.Pointer(v.RawVTable))
}

func (v *IWeakReference) GetWeakReference(guid *ole.GUID) (*ole.IInspectable, error) {
	var objectReference *ole.IInspectable
	r1, _, _ := syscall.SyscallN(v.VTable().Resolve, uintptr(unsafe.Pointer(v)), uintptr(unsafe.Pointer(&guid.Data1)), uintptr(unsafe.Pointer(&objectReference)))
	if r1 != win.S_OK {
		return nil, ole.NewError(r1)
	}
	return objectReference, nil
}
