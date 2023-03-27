package winrt

import (
	"unsafe"

	"github.com/go-ole/go-ole"
)

// IDirect3DDevice
// https://learn.microsoft.com/en-us/uwp/api/windows.graphics.directx.direct3d11.idirect3ddevice?view=winrt-22621

var IDirect3DDeviceID = ole.NewGUID("{A37624AB-8D5F-4650-9D3E-9EAE3D9BC670}")
var IDirect3DDeviceClass = "Windows.Graphics.DirectX.Direct3D11.IDirect3DDevice"

type IDirect3DDevice struct {
	ole.IInspectable
}

type IDirect3DDeviceVtbl struct {
	ole.IInspectableVtbl
	Trim uintptr
}

func (v *IDirect3DDevice) VTable() *IDirect3DDeviceVtbl {
	return (*IDirect3DDeviceVtbl)(unsafe.Pointer(v.RawVTable))
}
