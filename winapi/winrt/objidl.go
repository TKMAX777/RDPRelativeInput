package winrt

import (
	"syscall"
	"unsafe"

	"github.com/go-ole/go-ole"
)

var IMarshalID = ole.NewGUID("{00000003-0000-0000-C000-000000000046}")

type IMarshal struct {
	ole.IUnknown
}

type IMarshalVtbl struct {
	ole.IUnknownVtbl
	GetUnmarshalClass  uintptr
	GetMarshalSizeMax  uintptr
	MarshalInterface   uintptr
	UnmarshalInterface uintptr
	ReleaseMarshalData uintptr
	DisconnectObject   uintptr
}

func (v *IMarshal) VTable() *IMarshalVtbl {
	return (*IMarshalVtbl)(unsafe.Pointer(v.RawVTable))
}

func (v *IMarshal) GetUnmarshalClass(
	riid *ole.GUID,
	pv *uintptr,
	dwDestContext uint32,
	pvDestContext *uintptr,
	mshlflags uint32,
	pCid *uintptr,
) error {
	r1, _, _ := syscall.SyscallN(
		v.VTable().GetUnmarshalClass, uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(riid)),
		uintptr(unsafe.Pointer(pv)),
		uintptr(dwDestContext),
		uintptr(unsafe.Pointer(pvDestContext)),
		uintptr(mshlflags),
		uintptr(unsafe.Pointer(pCid)),
	)
	if r1 != 0 {
		return ole.NewError(r1)
	}
	return nil
}

func (v *IMarshal) GetMarshalSizeMax(
	riid *ole.GUID,
	pv *uintptr,
	dwDestContext uint32,
	pvDestContext *uintptr,
	mshlflags uint32,
) (size uint32, err error) {
	r1, _, _ := syscall.SyscallN(
		v.VTable().GetUnmarshalClass, uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(riid)),
		uintptr(unsafe.Pointer(pv)),
		uintptr(dwDestContext),
		uintptr(unsafe.Pointer(pvDestContext)),
		uintptr(mshlflags),
		uintptr(unsafe.Pointer(&size)),
	)
	if r1 != 0 {
		return 0, ole.NewError(r1)
	}
	return size, nil
}

// TODO: implement other methods
