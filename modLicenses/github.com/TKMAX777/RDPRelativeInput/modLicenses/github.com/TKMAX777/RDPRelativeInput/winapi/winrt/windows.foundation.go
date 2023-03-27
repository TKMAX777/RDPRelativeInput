package winrt

import (
	"syscall"
	"unsafe"

	"github.com/go-ole/go-ole"
	"github.com/lxn/win"
)

var IClosableID = ole.NewGUID("{30D5A829-7FA4-4026-83BB-D75BAE4EA99E}")

type IClosable struct {
	ole.IInspectable
}

type IClosableVtbl struct {
	ole.IInspectableVtbl
	Close uintptr
}

func (v *IClosable) VTable() *IClosableVtbl {
	return (*IClosableVtbl)(unsafe.Pointer(v.RawVTable))
}

func (v *IClosable) Close() error {
	r1, _, _ := syscall.SyscallN(v.VTable().Close, uintptr(unsafe.Pointer(v)))
	if r1 != win.S_OK {
		return ole.NewError(r1)
	}
	return nil
}
