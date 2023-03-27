package winrt

import (
	"syscall"
	"unsafe"

	"github.com/go-ole/go-ole"
	"github.com/lxn/win"
)

// GraphicsCaptureItemClass
// https://learn.microsoft.com/en-us/uwp/api/windows.graphics.capture.graphicscaptureitem?view=winrt-22621

var GraphicsCaptureItemClass = "Windows.Graphics.Capture.GraphicsCaptureItem"

// IGraphicsCaptureItem

var IGraphicsCaptureItemID = ole.NewGUID("{79c3f95b-31f7-4ec2-a464-632ef5d30760}")
var IGraphicsCaptureItemClass = "Windows.Graphics.Capture.IGraphicsCaptureItem"

type IGraphicsCaptureItem struct {
	ole.IInspectable
}

type IGraphicsCaptureItemVtbl struct {
	ole.IInspectableVtbl
	DisplayName   uintptr
	Size          uintptr
	add_Closed    uintptr
	remove_Closed uintptr
}

func (v *IGraphicsCaptureItem) VTable() *IGraphicsCaptureItemVtbl {
	return (*IGraphicsCaptureItemVtbl)(unsafe.Pointer(v.RawVTable))
}

func (v *IGraphicsCaptureItem) DisplayName() (string, error) {
	var hRet ole.HString

	r1, _, _ := syscall.SyscallN(v.VTable().DisplayName, uintptr(unsafe.Pointer(v)), uintptr(unsafe.Pointer(&hRet)))
	if r1 != 0 {
		return "", ole.NewError(r1)
	}

	var ret = hRet.String()
	ole.DeleteHString(hRet)

	return ret, nil
}

func (v *IGraphicsCaptureItem) Size() (*SizeInt32, error) {
	var size SizeInt32
	r1, _, _ := syscall.SyscallN(v.VTable().Size, uintptr(unsafe.Pointer(v)), uintptr(unsafe.Pointer(&size.Width)))
	if r1 != 0 {
		return nil, ole.NewError(r1)
	}
	return &size, nil
}

// IGraphicsCaptureItemStatics

var IGraphicsCaptureItemStaticsID = ole.NewGUID("{a87ebea5-457c-5788-ab47-0cf1d3637e74}")
var IGraphicsCaptureItemStaticsClass = "Windows.Graphics.Capture.IGraphicsCaptureItemStatics"

type IGraphicsCaptureItemStatics struct {
	ole.IInspectable
}

type IGraphicsCaptureItemStaticsVtbl struct {
	ole.IInspectableVtbl
	CreateFromVisual uintptr
}

func (v *IGraphicsCaptureItemStatics) VTable() *IGraphicsCaptureItemStaticsVtbl {
	return (*IGraphicsCaptureItemStaticsVtbl)(unsafe.Pointer(v.RawVTable))
}

// IGraphicsCaptureItemStatics2

var IGraphicsCaptureItemStatics2ID = ole.NewGUID("{3b92acc9-e584-5862-bf5c-9c316c6d2dbb}")
var IGraphicsCaptureItemStatics2Class = "Windows.Graphics.Capture.IGraphicsCaptureItemStatics2"

type IGraphicsCaptureItemStatics2 struct {
	ole.IInspectable
}

type IGraphicsCaptureItemStatics2Vtbl struct {
	ole.IInspectableVtbl
	TryCreateFromWindowId  uintptr
	TryCreateFromDisplayId uintptr
}

func (v *IGraphicsCaptureItemStatics2) VTable() *IGraphicsCaptureItemStatics2Vtbl {
	return (*IGraphicsCaptureItemStatics2Vtbl)(unsafe.Pointer(v.RawVTable))
}

// Direct3D11CaptureFramePool
// https://learn.microsoft.com/en-us/uwp/api/windows.graphics.capture.direct3d11captureframepool?view=winrt-22621
var Direct3D11CaptureFramePoolClass = "Windows.Graphics.Capture.Direct3D11CaptureFramePool"

// IDirect3D11CaptureFramePool

var IDirect3D11CaptureFramePoolID = ole.NewGUID("{24EB6D22-1975-422E-82E7-780DBD8DDF24}")

type IDirect3D11CaptureFramePool struct {
	ole.IInspectable
}

type IDirect3D11CaptureFramePoolVtbl struct {
	ole.IInspectableVtbl
	Recreate             uintptr
	TryGetNextFrame      uintptr
	add_FrameArrived     uintptr
	remove_FrameArrived  uintptr
	CreateCaptureSession uintptr
	get_DispatcherQueue  uintptr
}

func (v *IDirect3D11CaptureFramePool) VTable() *IDirect3D11CaptureFramePoolVtbl {
	return (*IDirect3D11CaptureFramePoolVtbl)(unsafe.Pointer(v.RawVTable))
}

// type Direct3D11CaptureFramePoolFrameArrivedProcType func(this *Direct3D11CaptureFramePool, sender *IDirect3D11CaptureFramePool, args *ole.IInspectable) uintptr
type Direct3D11CaptureFramePoolFrameArrivedProcType func(this *uintptr, sender *IDirect3D11CaptureFramePool, args *ole.IInspectable) uintptr

/*
eventHandler:
	interface {
		IUnknown
		Invoke(sender *IDirect3D11CaptureFramePool, args *ole.IInspectable) uintptr
	}
*/
func (v *IDirect3D11CaptureFramePool) AddFrameArrived(eventHandler unsafe.Pointer) (*EventRegistrationToken, error) {
	var token EventRegistrationToken
	r1, _, _ := syscall.SyscallN(v.VTable().add_FrameArrived, uintptr(unsafe.Pointer(v)), uintptr(eventHandler), uintptr(unsafe.Pointer(&token.value)))
	if r1 != win.S_OK {
		return nil, ole.NewError(r1)
	}
	return &token, nil
}

func (v *IDirect3D11CaptureFramePool) RemoveFrameArrived(token *EventRegistrationToken) error {
	r1, _, _ := syscall.SyscallN(v.VTable().remove_FrameArrived, uintptr(unsafe.Pointer(v)), uintptr(unsafe.Pointer(&token.value)))
	if r1 != win.S_OK {
		return ole.NewError(r1)
	}
	return nil
}

func (v *IDirect3D11CaptureFramePool) CreateCaptureSession(item *IGraphicsCaptureItem) (*IGraphicsCaptureSession, error) {
	var session *IGraphicsCaptureSession
	r1, _, _ := syscall.SyscallN(v.VTable().CreateCaptureSession, uintptr(unsafe.Pointer(v)), uintptr(unsafe.Pointer(item)), uintptr(unsafe.Pointer(&session)))
	if r1 != win.S_OK {
		return nil, ole.NewError(r1)
	}

	return session, nil
}

func (v *IDirect3D11CaptureFramePool) TryGetNextFrame() (*IDirect3D11CaptureFrame, error) {
	var frame IDirect3D11CaptureFrame
	r1, _, _ := syscall.SyscallN(v.VTable().TryGetNextFrame, uintptr(unsafe.Pointer(v)), uintptr(unsafe.Pointer(&frame)))
	if r1 != win.S_OK {
		return nil, ole.NewError(r1)
	}

	return &frame, nil
}

// IDirect3D11CaptureFramePoolStatics

var IDirect3D11CaptureFramePoolStaticsID = ole.NewGUID("{7784056A-67AA-4D53-AE54-1088D5A8CA21}")

type IDirect3D11CaptureFramePoolStatics struct {
	ole.IInspectable
}

type IDirect3D11CaptureFramePoolStaticsVtbl struct {
	ole.IInspectableVtbl
	Create uintptr
}

func (v *IDirect3D11CaptureFramePoolStatics) VTable() *IDirect3D11CaptureFramePoolStaticsVtbl {
	return (*IDirect3D11CaptureFramePoolStaticsVtbl)(unsafe.Pointer(v.RawVTable))
}

func (v *IDirect3D11CaptureFramePoolStatics) Create(device *IDirect3DDevice, pixelFormat DirectXPixelFormat, numberOfBuffers int32, size *SizeInt32) (*IDirect3D11CaptureFramePool, error) {
	var ret *IDirect3D11CaptureFramePool
	r1, _, _ := syscall.SyscallN(
		v.VTable().Create, uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(device)), uintptr(pixelFormat), uintptr(numberOfBuffers), uintptr(size.Width)<<32+uintptr(size.Height), uintptr(unsafe.Pointer(&ret)),
	)
	if r1 != win.S_OK {
		return nil, ole.NewError(r1)
	}

	return ret, nil
}

// IDirect3D11CaptureFramePoolStatics2

var IDirect3D11CaptureFramePoolStatics2ID = ole.NewGUID("{589B103F-6BBC-5DF5-A991-02E28B3B66D5}")

type IDirect3D11CaptureFramePoolStatics2 struct {
	ole.IInspectable
}

type IDirect3D11CaptureFramePoolStatics2Vtbl struct {
	ole.IInspectableVtbl
	CreateFreeThreaded uintptr
}

func (v *IDirect3D11CaptureFramePoolStatics2) VTable() *IDirect3D11CaptureFramePoolStatics2Vtbl {
	return (*IDirect3D11CaptureFramePoolStatics2Vtbl)(unsafe.Pointer(v.RawVTable))
}

func (v *IDirect3D11CaptureFramePoolStatics2) CreateFreeThreaded(device *IDirect3DDevice, pixelFormat DirectXPixelFormat, numberOfBuffers int32, size *SizeInt32) (*IDirect3D11CaptureFramePool, error) {
	var ret *IDirect3D11CaptureFramePool
	r1, _, _ := syscall.SyscallN(
		v.VTable().CreateFreeThreaded, uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(device)), uintptr(pixelFormat), uintptr(numberOfBuffers), uintptr(size.Width)<<32+uintptr(size.Height), uintptr(unsafe.Pointer(&ret)),
	)
	if r1 != win.S_OK {
		return nil, ole.NewError(r1)
	}

	return ret, nil
}

type IGraphicsCaptureSession struct {
	ole.IInspectable
}

type IGraphicsCaptureSessionVtbl struct {
	ole.IInspectableVtbl
	StartCapture uintptr
}

var IGraphicsCaptureSessionID = ole.NewGUID("{814E42A9-F70F-4AD7-939B-FDDCC6EB880D}")

func (v *IGraphicsCaptureSession) VTable() *IGraphicsCaptureSessionVtbl {
	return (*IGraphicsCaptureSessionVtbl)(unsafe.Pointer(v.RawVTable))
}

func (v *IGraphicsCaptureSession) StartCapture() error {
	r1, _, _ := syscall.SyscallN(v.VTable().StartCapture, uintptr(unsafe.Pointer(v)))
	if r1 != win.S_OK {
		return ole.NewError(r1)
	}

	return nil
}

type IGraphicsCaptureSession2 struct {
	ole.IInspectable
}

type IGraphicsCaptureSession2Vtbl struct {
	ole.IInspectableVtbl
	get_IsCursorCaptureEnabled uintptr
	put_IsCursorCaptureEnabled uintptr
}

var IGraphicsCaptureSession2ID = ole.NewGUID("{2C39AE40-7D2E-5044-804E-8B6799D4CF9E}")

func (v *IGraphicsCaptureSession2) VTable() *IGraphicsCaptureSession2Vtbl {
	return (*IGraphicsCaptureSession2Vtbl)(unsafe.Pointer(v.RawVTable))
}

type IGraphicsCaptureSession3 struct {
	ole.IInspectable
}

type IGraphicsCaptureSession3Vtbl struct {
	ole.IInspectableVtbl
	get_IsBorderRequired uintptr
	put_IsBorderRequired uintptr
}

var IGraphicsCaptureSession3ID = ole.NewGUID("{F2CDD966-22AE-5EA1-9596-3A289344C3BE}")

func (v *IGraphicsCaptureSession3) VTable() *IGraphicsCaptureSession3Vtbl {
	return (*IGraphicsCaptureSession3Vtbl)(unsafe.Pointer(v.RawVTable))
}

type IGraphicsCaptureSessionStatics struct {
	ole.IInspectable
}

type IGraphicsCaptureSessionStaticsVtbl struct {
	ole.IInspectableVtbl
	IsSupported uintptr
}

var IGraphicsCaptureSessionStaticsID = ole.NewGUID("{2224A540-5974-49AA-B232-0882536F4CB5}")

func (v *IGraphicsCaptureSessionStatics) VTable() *IGraphicsCaptureSessionStaticsVtbl {
	return (*IGraphicsCaptureSessionStaticsVtbl)(unsafe.Pointer(v.RawVTable))
}

func (v *IGraphicsCaptureSessionStatics) IsSupported() (ok bool) {
	syscall.SyscallN(v.VTable().IsSupported, uintptr(unsafe.Pointer(v)), uintptr(unsafe.Pointer(&ok)))
	return ok
}

var IDirect3D11CaptureFrameClass = "Windows.Graphics.Capture.IDirect3D11CaptureFrame"
var IDirect3D11CaptureFrameID = "{FA50C623-38DA-4B32-ACF3-FA9734AD800E}"

type IDirect3D11CaptureFrame struct {
	ole.IInspectable
}

type IDirect3D11CaptureFrameVtbl struct {
	ole.IInspectableVtbl
	Surface            uintptr
	SystemRelativeTime uintptr
	ContentSize        uintptr
}

func (v *IDirect3D11CaptureFrame) VTable() *IDirect3D11CaptureFrameVtbl {
	return (*IDirect3D11CaptureFrameVtbl)(unsafe.Pointer(v.RawVTable))
}
