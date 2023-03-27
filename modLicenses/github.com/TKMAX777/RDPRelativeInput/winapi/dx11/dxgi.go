package dx11

import (
	"unsafe"

	"github.com/go-ole/go-ole"
)

var IDXGIObjectID = ole.NewGUID("{aec22fb8-76f3-4639-9be0-28eb43a67a2e}")

type IDXGIAdapter uintptr

type IDXGIObject struct {
	ole.IUnknown
}

type IDXGIObjectVtbl struct {
	ole.IUnknownVtbl
	SetPrivateData          uintptr
	SetPrivateDataInterface uintptr
	GetPrivateData          uintptr
	GetParent               uintptr
}

func (v *IDXGIObject) VTable() *IDXGIObjectVtbl {
	return (*IDXGIObjectVtbl)(unsafe.Pointer(v.RawVTable))
}

var IDXGIDeviceID = ole.NewGUID("{54ec77fa-1377-44e6-8c32-88fd5f44c84c}")

type IDXGIDevice struct {
	IDXGIObject
}

type IDXGIDeviceVtbl struct {
	IDXGIObjectVtbl
	GetAdapter             uintptr
	CreateSurface          uintptr
	QueryResourceResidency uintptr
	SetGPUThreadPriority   uintptr
	GetGPUThreadPriority   uintptr
}

func (v *IDXGIDevice) VTable() *IDXGIDeviceVtbl {
	return (*IDXGIDeviceVtbl)(unsafe.Pointer(v.RawVTable))
}
