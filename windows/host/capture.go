package host

import (
	"fmt"
	"os"
	"runtime"
	"unsafe"

	"github.com/TKMAX777/RDPRelativeInput/winapi/dx11"
	"github.com/TKMAX777/RDPRelativeInput/winapi/winrt"
	"github.com/go-ole/go-ole"
	"github.com/lxn/win"
	"github.com/pkg/errors"
)

type CaptureHandler struct {
	device                 *winrt.IDirect3DDevice
	deviceDx               *dx11.ID3D11Device
	graphicsCaptureItem    *winrt.IGraphicsCaptureItem
	framePool              *winrt.IDirect3D11CaptureFramePool
	graphicsCaptureSession *winrt.IGraphicsCaptureSession
	framePoolToken         *winrt.EventRegistrationToken
}

func (c *CaptureHandler) StartCapture(hwnd win.HWND) error {
	type resultAttr struct {
		err error
	}

	var result = make(chan resultAttr)

	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		// Initialize Windows Runtime
		err := winrt.RoInitialize(winrt.RO_INIT_MULTITHREADED)
		if err != nil {
			result <- resultAttr{errors.Wrap(err, "RoInitialize")}
			return
		}
		defer winrt.RoUninitialize()

		// Create capture device
		var featureLevels = []dx11.D3D_FEATURE_LEVEL{
			dx11.D3D_FEATURE_LEVEL_11_0,
			dx11.D3D_FEATURE_LEVEL_10_1,
			dx11.D3D_FEATURE_LEVEL_10_0,
			dx11.D3D_FEATURE_LEVEL_9_3,
			dx11.D3D_FEATURE_LEVEL_9_2,
			dx11.D3D_FEATURE_LEVEL_9_1,
		}

		err = dx11.D3D11CreateDevice(
			nil, dx11.D3D_DRIVER_TYPE_HARDWARE, 0, dx11.D3D11_CREATE_DEVICE_BGRA_SUPPORT|dx11.D3D11_CREATE_DEVICE_DEBUG,
			&featureLevels[0], len(featureLevels),
			dx11.D3D11_SDK_VERSION, &c.deviceDx, nil, nil,
		)
		if err != nil {
			result <- resultAttr{errors.Wrap(err, "D3DCreateDevice")}
			return
		}

		// Query interface of DXGIDevice
		var dxgiDevice *dx11.IDXGIDevice
		err = c.deviceDx.PutQueryInterface(dx11.IDXGIDeviceID, &dxgiDevice)
		if err != nil {
			result <- resultAttr{errors.Wrap(err, "PutQueryInterface")}
			return
		}

		var deviceRT *ole.IInspectable

		// convert D3D11Device(Dx11) to Direct3DDevice(WinRT)
		err = dx11.CreateDirect3D11DeviceFromDXGIDevice(dxgiDevice, &deviceRT)
		if err != nil {
			result <- resultAttr{errors.Wrap(err, "CreateDirect3D11DeviceFromDXGIDevice")}
			return
		}

		// Query interface of IDirect3DDevice
		err = deviceRT.PutQueryInterface(winrt.IDirect3DDeviceID, &c.device)
		if err != nil {
			result <- resultAttr{errors.Wrap(err, "QueryInterface: IDirect3DDeviceID")}
			return
		}

		// Create Capture Settings
		factory, err := ole.RoGetActivationFactory(winrt.GraphicsCaptureItemClass, winrt.IGraphicsCaptureItemInteropID)
		if err != nil {
			result <- resultAttr{errors.Wrap(err, "RoGetActivationFactory: IGraphicsCaptureItemID")}
			return
		}

		var interop *winrt.IGraphicsCaptureItemInterop
		err = factory.PutQueryInterface(winrt.IGraphicsCaptureItemInteropID, &interop)
		if err != nil {
			result <- resultAttr{errors.Wrap(err, "QueryInterface: IGraphicsCaptureItemInteropID")}
			return
		}

		var captureItemDispatch *ole.IInspectable

		// Capture for the window specified
		err = interop.CreateForWindow(hwnd, winrt.IGraphicsCaptureItemID, &captureItemDispatch)
		if err != nil {
			result <- resultAttr{errors.Wrap(err, "CreateForWindow")}
			return
		}

		// Get Interface of IGraphicsCaptureItem
		err = captureItemDispatch.PutQueryInterface(winrt.IGraphicsCaptureItemID, &c.graphicsCaptureItem)
		if err != nil {
			result <- resultAttr{errors.Wrap(err, "PutQueryInterface captureItemDispatch")}
			return
		}

		// Get Capture objects size
		size, err := c.graphicsCaptureItem.Size()
		if err != nil {
			result <- resultAttr{errors.Wrap(err, "Size")}
			return
		}

		// Get object of Direct3D11CaptureFramePoolClass
		ins, err := ole.RoGetActivationFactory(winrt.Direct3D11CaptureFramePoolClass, winrt.IDirect3D11CaptureFramePoolStaticsID)
		if err != nil {
			result <- resultAttr{errors.Wrap(err, "RoGetActivationFactory: IDirect3D11CaptureFramePoolStatics Class Instance")}
			return
		}
		defer ins.Release()

		// Get Interface of Direct3D11CaptureFramePoolClass
		var framePoolStatic *winrt.IDirect3D11CaptureFramePoolStatics
		err = ins.PutQueryInterface(winrt.IDirect3D11CaptureFramePoolStaticsID, &framePoolStatic)
		if err != nil {
			result <- resultAttr{errors.Wrap(err, "PutQueryInterface: IDirect3D11CaptureFramePoolStaticsID")}
			return
		}

		// Create frame pool
		c.framePool, err = framePoolStatic.Create(c.device, winrt.DirectXPixelFormat_B8G8R8A8UIntNormalized, 1, size)
		if err != nil {
			result <- resultAttr{errors.Wrap(err, "CreateFramePool")}
			return
		}

		// Set frame settings
		var eventObject = NewDirect3D11CaptureFramePool(c.onFrameArrived)

		c.framePoolToken, err = c.framePool.AddFrameArrived(unsafe.Pointer(eventObject))
		if err != nil {
			result <- resultAttr{errors.Wrap(err, "AddFrameArrived")}
			return
		}

		c.graphicsCaptureSession, err = c.framePool.CreateCaptureSession(c.graphicsCaptureItem)
		if err != nil {
			result <- resultAttr{errors.Wrap(err, "CreateCaptureSession")}
			return
		}

		// Start capturing
		err = c.graphicsCaptureSession.StartCapture()
		if err != nil {
			result <- resultAttr{errors.Wrap(err, "StartCapture")}
			return
		}
		result <- resultAttr{nil}

		for {
			var msg win.MSG
			switch win.GetMessage(&msg, 0, 0, 0) {
			case 0:
				fmt.Println("EXIT")
				return
			case -1:
				fmt.Println("ERROR")
				os.Exit(0)
				return
			default:
				win.TranslateMessage(&msg)
				win.DispatchMessage(&msg)
			}
		}
	}()

	var res = <-result

	return res.err
}

func (c *CaptureHandler) onFrameArrived(this_ *uintptr, sender *winrt.IDirect3D11CaptureFramePool, args *ole.IInspectable) uintptr {
	_ = (*Direct3D11CaptureFramePool)(unsafe.Pointer(this_))

	_, err := sender.TryGetNextFrame()
	if err != nil {
		os.Stderr.Write([]byte("Error: TryGetNextFrame: " + err.Error()))
		return 0
	}

	return 0
}

func (c *CaptureHandler) Close() error {
	if c.framePool != nil {
		err := c.framePool.RemoveFrameArrived(c.framePoolToken)
		if err != nil {
			return errors.Wrap(err, "RemoveFrameArrived")
		}

		var closable *winrt.IClosable
		err = c.framePool.PutQueryInterface(winrt.IClosableID, &closable)
		if err != nil {
			return errors.Wrap(err, "PutQueryInterface: graphicsCaptureSession")
		}
		defer closable.Release()

		closable.Close()

		c.framePool = nil
	}

	var closable *winrt.IClosable
	err := c.graphicsCaptureSession.PutQueryInterface(winrt.IClosableID, &closable)
	if err != nil {
		return errors.Wrap(err, "PutQueryInterface: graphicsCaptureSession")
	}
	defer closable.Release()

	closable.Close()

	c.graphicsCaptureItem = nil

	return nil
}
