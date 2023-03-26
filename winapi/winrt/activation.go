package winrt

import (
	"syscall"
	"unsafe"

	"github.com/go-ole/go-ole"
	"github.com/lxn/win"
)

var IActivationFactoryID = ole.NewGUID("{00000035-0000-0000-C000-000000000046}")

type IActivationFactory struct {
	ole.IInspectable
}

type IActivationFactoryVtbl struct {
	ole.IInspectableVtbl
	ActivateInstance uintptr
}

func (v *IActivationFactory) VTable() *IActivationFactoryVtbl {
	return (*IActivationFactoryVtbl)(unsafe.Pointer(v.RawVTable))
}

func (v *IActivationFactory) ActivateInstance() (*ole.IInspectable, error) {
	var ins *ole.IInspectable

	r1, _, _ := syscall.SyscallN(v.VTable().ActivateInstance, uintptr(unsafe.Pointer(v)), uintptr(unsafe.Pointer(&ins)))
	if r1 != win.S_OK {
		return nil, ole.NewError(r1)
	}
	return ins, nil
}
