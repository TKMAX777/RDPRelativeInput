package host

import (
	"errors"
	"runtime"
	"syscall"
	"unsafe"

	"github.com/TKMAX777/RDPRelativeInput/winapi/winrt"
	"github.com/go-ole/go-ole"
	"github.com/lxn/win"
)

// Direct3D11CaptureFramePool

// Protect from gabage collecter
var generatedDirect3D11CaptureFramePool = map[uintptr]*Direct3D11CaptureFramePoolVtbl{}

type Direct3D11CaptureFramePool struct {
	ole.IUnknown
}

type Direct3D11CaptureFramePoolVtbl struct {
	ole.IUnknownVtbl
	Invoke  uintptr
	counter int
}

func NewDirect3D11CaptureFramePool(invoke winrt.Direct3D11CaptureFramePoolFrameArrivedProcType) *Direct3D11CaptureFramePool {
	var v = &Direct3D11CaptureFramePoolVtbl{
		Invoke:  syscall.NewCallback(invoke),
		counter: 0,
	}

	var newV = new(Direct3D11CaptureFramePool)
	newV.RawVTable = (*interface{})(unsafe.Pointer(v))

	v.QueryInterface = syscall.NewCallback(newV.queryInterface)
	v.AddRef = syscall.NewCallback(newV.addRef)
	v.Release = syscall.NewCallback(newV.release)

	generatedDirect3D11CaptureFramePool[uintptr(unsafe.Pointer(newV))] = v

	return newV
}

func (v *Direct3D11CaptureFramePool) VTable() *Direct3D11CaptureFramePoolVtbl {
	return (*Direct3D11CaptureFramePoolVtbl)(unsafe.Pointer(v.RawVTable))
}

func (v *Direct3D11CaptureFramePool) Invoke(sender *winrt.IDirect3D11CaptureFramePool, args *ole.IInspectable) error {
	r1, _, _ := syscall.SyscallN(v.VTable().Invoke, uintptr(unsafe.Pointer(sender)), uintptr(unsafe.Pointer(args)))
	return ole.NewError(r1)
}

// QueryInterface(vp *Direct3D11CaptureFramePool, riid ole.GUID, lppvObj **ole.Inspectable)
func (v *Direct3D11CaptureFramePool) queryInterface(lpMyObj *uintptr, riid *uintptr, lppvObj **uintptr) uintptr {
	// Validate input
	if lpMyObj == nil {
		return win.E_INVALIDARG
	}

	var V = new(Direct3D11CaptureFramePool)

	var err error
	// Check dereferencability
	func() {
		defer func() {
			if recover() != nil {
				err = errors.New("InvalidObject")
			}
		}()
		// if object cannot be dereferenced, then panic occurs
		*V = *(*Direct3D11CaptureFramePool)(unsafe.Pointer(lpMyObj))
		V.VTable()
	}()
	if err != nil {
		return win.E_INVALIDARG
	}

	*lppvObj = nil
	var id = new(ole.GUID)
	*id = *(*ole.GUID)(unsafe.Pointer(riid))

	// Convert
	switch id.String() {
	case ole.IID_IUnknown.String(), winrt.ITypedEventHandlerID.String():
		V.AddRef()
		*lppvObj = (*uintptr)(unsafe.Pointer(V))

		return win.S_OK
	default:
		return win.E_NOINTERFACE
	}
}

func (v *Direct3D11CaptureFramePool) addRef(lpMyObj *uintptr) uintptr {
	// Validate input
	if lpMyObj == nil {
		return 0
	}

	var V = (*Direct3D11CaptureFramePool)(unsafe.Pointer(lpMyObj))
	V.VTable().counter++

	return uintptr(V.VTable().counter)
}

func (v *Direct3D11CaptureFramePool) release(lpMyObj *uintptr) uintptr {
	// Validate input
	if lpMyObj == nil {
		return 0
	}

	var V = (*Direct3D11CaptureFramePool)(unsafe.Pointer(lpMyObj))
	V.VTable().counter--

	if V.VTable().counter == 0 {
		_, ok := generatedDirect3D11CaptureFramePool[uintptr(unsafe.Pointer(lpMyObj))]
		if ok {
			delete(generatedDirect3D11CaptureFramePool, uintptr(unsafe.Pointer(lpMyObj)))
			runtime.GC()
		}
	}

	return uintptr(V.VTable().counter)
}
