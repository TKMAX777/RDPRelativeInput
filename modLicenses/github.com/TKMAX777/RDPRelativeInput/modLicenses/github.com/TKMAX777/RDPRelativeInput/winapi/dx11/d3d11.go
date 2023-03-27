package dx11

import (
	"syscall"
	"unsafe"

	"github.com/go-ole/go-ole"
	"github.com/lxn/win"
	"golang.org/x/sys/windows"
)

var (
	d3d11DLL = windows.NewLazySystemDLL("d3d11.dll")
)

const D3D11_SDK_VERSION = 7

type D3D11_CREATE_DEVICE_FLAG uint32

const (
	D3D11_CREATE_DEVICE_SINGLETHREADED                                D3D11_CREATE_DEVICE_FLAG = 0x1
	D3D11_CREATE_DEVICE_DEBUG                                         D3D11_CREATE_DEVICE_FLAG = 0x2
	D3D11_CREATE_DEVICE_SWITCH_TO_REF                                 D3D11_CREATE_DEVICE_FLAG = 0x4
	D3D11_CREATE_DEVICE_PREVENT_INTERNAL_THREADING_OPTIMIZATIONS      D3D11_CREATE_DEVICE_FLAG = 0x8
	D3D11_CREATE_DEVICE_BGRA_SUPPORT                                  D3D11_CREATE_DEVICE_FLAG = 0x20
	D3D11_CREATE_DEVICE_DEBUGGABLE                                    D3D11_CREATE_DEVICE_FLAG = 0x40
	D3D11_CREATE_DEVICE_PREVENT_ALTERING_LAYER_SETTINGS_FROM_REGISTRY D3D11_CREATE_DEVICE_FLAG = 0x80
	D3D11_CREATE_DEVICE_DISABLE_GPU_TIMEOUT                           D3D11_CREATE_DEVICE_FLAG = 0x100
	D3D11_CREATE_DEVICE_VIDEO_SUPPORT                                 D3D11_CREATE_DEVICE_FLAG = 0x800
)

var ID3D11DeviceID = ole.NewGUID("{db6f6ddb-ac77-4e88-8253-819df9bbf140}")

type ID3D11Device struct {
	ole.IUnknown
}

type ID3D11DeviceVtbl struct {
	ole.IUnknownVtbl
	CreateBuffer                         uintptr
	CreateTexture1D                      uintptr
	CreateTexture2D                      uintptr
	CreateTexture3D                      uintptr
	CreateShaderResourceView             uintptr
	CreateUnorderedAccessView            uintptr
	CreateRenderTargetView               uintptr
	CreateDepthStencilView               uintptr
	CreateInputLayout                    uintptr
	CreateVertexShader                   uintptr
	CreateGeometryShader                 uintptr
	CreateGeometryShaderWithStreamOutput uintptr
	CreatePixelShader                    uintptr
	CreateHullShader                     uintptr
	CreateDomainShader                   uintptr
	CreateComputeShader                  uintptr
	CreateClassLinkage                   uintptr
	CreateBlendState                     uintptr
	CreateDepthStencilState              uintptr
	CreateRasterizerState                uintptr
	CreateSamplerState                   uintptr
	CreateQuery                          uintptr
	CreatePredicate                      uintptr
	CreateCounter                        uintptr
	CreateDeferredContext                uintptr
	OpenSharedResource                   uintptr
	CheckFormatSupport                   uintptr
	CheckMultisampleQualityLevels        uintptr
	CheckCounterInfo                     uintptr
	CheckCounter                         uintptr
	CheckFeatureSupport                  uintptr
	GetPrivateData                       uintptr
	SetPrivateData                       uintptr
	SetPrivateDataInterface              uintptr
	GetFeatureLevel                      uintptr
	GetCreationFlags                     uintptr
	GetDeviceRemovedReason               uintptr
	GetImmediateContext                  uintptr
	SetExceptionMode                     uintptr
	GetExceptionMode                     uintptr
}

func (v *ID3D11Device) VTable() *ID3D11DeviceVtbl {
	return (*ID3D11DeviceVtbl)(unsafe.Pointer(v.RawVTable))
}

func (v *ID3D11Device) GetImmediateContext() (pImmediateContext *ID3D11DeviceContext) {
	syscall.SyscallN(v.VTable().GetImmediateContext, uintptr(unsafe.Pointer(v)), uintptr(unsafe.Pointer(&pImmediateContext)))
	return pImmediateContext
}

var ID3D11DeviceContextID = ole.NewGUID("{c0bfa96c-e089-44fb-8eaf-26f8796190da}")

type ID3D11DeviceContext struct {
	ole.IUnknown
}

type ID3D11DeviceContextVtbl struct {
	ole.IUnknownVtbl
	GetDevice                                 uintptr
	GetPrivateData                            uintptr
	SetPrivateData                            uintptr
	SetPrivateDataInterface                   uintptr
	VSSetConstantBuffers                      uintptr
	PSSetShaderResources                      uintptr
	PSSetShader                               uintptr
	PSSetSamplers                             uintptr
	VSSetShader                               uintptr
	DrawIndexed                               uintptr
	Draw                                      uintptr
	Map                                       uintptr
	Unmap                                     uintptr
	PSSetConstantBuffers                      uintptr
	IASetInputLayout                          uintptr
	IASetVertexBuffers                        uintptr
	IASetIndexBuffer                          uintptr
	DrawIndexedInstanced                      uintptr
	DrawInstanced                             uintptr
	GSSetConstantBuffers                      uintptr
	GSSetShader                               uintptr
	IASetPrimitiveTopology                    uintptr
	VSSetShaderResources                      uintptr
	VSSetSamplers                             uintptr
	Begin                                     uintptr
	End                                       uintptr
	GetData                                   uintptr
	SetPredication                            uintptr
	GSSetShaderResources                      uintptr
	GSSetSamplers                             uintptr
	OMSetRenderTargets                        uintptr
	OMSetRenderTargetsAndUnorderedAccessViews uintptr
	OMSetBlendState                           uintptr
	OMSetDepthStencilState                    uintptr
	SOSetTargets                              uintptr
	DrawAuto                                  uintptr
	DrawIndexedInstancedIndirect              uintptr
	DrawInstancedIndirect                     uintptr
	Dispatch                                  uintptr
	DispatchIndirect                          uintptr
	RSSetState                                uintptr
	RSSetViewports                            uintptr
	RSSetScissorRects                         uintptr
	CopySubresourceRegion                     uintptr
	CopyResource                              uintptr
	UpdateSubresource                         uintptr
	CopyStructureCount                        uintptr
	ClearRenderTargetView                     uintptr
	ClearUnorderedAccessViewUint              uintptr
	ClearUnorderedAccessViewFloat             uintptr
	ClearDepthStencilView                     uintptr
	GenerateMips                              uintptr
	SetResourceMinLOD                         uintptr
	GetResourceMinLOD                         uintptr
	ResolveSubresource                        uintptr
	ExecuteCommandList                        uintptr
	HSSetShaderResources                      uintptr
	HSSetShader                               uintptr
	HSSetSamplers                             uintptr
	HSSetConstantBuffers                      uintptr
	DSSetShaderResources                      uintptr
	DSSetShader                               uintptr
	DSSetSamplers                             uintptr
	DSSetConstantBuffers                      uintptr
	CSSetShaderResources                      uintptr
	CSSetUnorderedAccessViews                 uintptr
	CSSetShader                               uintptr
	CSSetSamplers                             uintptr
	CSSetConstantBuffers                      uintptr
	VSGetConstantBuffers                      uintptr
	PSGetShaderResources                      uintptr
	PSGetShader                               uintptr
	PSGetSamplers                             uintptr
	VSGetShader                               uintptr
	PSGetConstantBuffers                      uintptr
	IAGetInputLayout                          uintptr
	IAGetVertexBuffers                        uintptr
	IAGetIndexBuffer                          uintptr
	GSGetConstantBuffers                      uintptr
	GSGetShader                               uintptr
	IAGetPrimitiveTopology                    uintptr
	VSGetShaderResources                      uintptr
	VSGetSamplers                             uintptr
	GetPredication                            uintptr
	GSGetShaderResources                      uintptr
	GSGetSamplers                             uintptr
	OMGetRenderTargets                        uintptr
	OMGetRenderTargetsAndUnorderedAccessViews uintptr
	OMGetBlendState                           uintptr
	OMGetDepthStencilState                    uintptr
	SOGetTargets                              uintptr
	RSGetState                                uintptr
	RSGetViewports                            uintptr
	RSGetScissorRects                         uintptr
	HSGetShaderResources                      uintptr
	HSGetShader                               uintptr
	HSGetSamplers                             uintptr
	HSGetConstantBuffers                      uintptr
	DSGetShaderResources                      uintptr
	DSGetShader                               uintptr
	DSGetSamplers                             uintptr
	DSGetConstantBuffers                      uintptr
	CSGetShaderResources                      uintptr
	CSGetUnorderedAccessViews                 uintptr
	CSGetShader                               uintptr
	CSGetSamplers                             uintptr
	CSGetConstantBuffers                      uintptr
	ClearState                                uintptr
	Flush                                     uintptr
	GetType                                   uintptr
	GetContextFlags                           uintptr
	FinishCommandList                         uintptr
}

func (v *ID3D11DeviceContext) VTable() *ID3D11DeviceContextVtbl {
	return (*ID3D11DeviceContextVtbl)(unsafe.Pointer(v.RawVTable))
}

var pD3DCreateDevice = d3d11DLL.NewProc("D3D11CreateDevice")

// CreateDevice
// https://learn.microsoft.com/en-us/windows/win32/api/d3d11/nf-d3d11-d3d11createdevice
func D3D11CreateDevice(
	pAdapter *IDXGIAdapter,
	DriverType D3D_DRIVER_TYPE,
	Software win.HMODULE,
	Flags D3D11_CREATE_DEVICE_FLAG,
	pFeatureLevels *D3D_FEATURE_LEVEL,
	FeatureLevels int,
	SDKVersion uint32,
	ppDevice **ID3D11Device,
	pFeatureLevel *D3D_FEATURE_LEVEL,
	ppImmediateContext **ID3D11DeviceContext,
) error {
	r1, _, _ := pD3DCreateDevice.Call(
		uintptr(unsafe.Pointer(pAdapter)),
		uintptr(DriverType),
		uintptr(Software),
		uintptr(Flags),
		uintptr(unsafe.Pointer(pFeatureLevels)),
		uintptr(FeatureLevels),
		uintptr(SDKVersion),
		uintptr(unsafe.Pointer(ppDevice)),
		uintptr(unsafe.Pointer(pFeatureLevel)),
		uintptr(unsafe.Pointer(ppImmediateContext)),
	)
	if r1 != win.S_OK {
		return ole.NewError(r1)
	}
	return nil
}

var pCreateDirect3D11DeviceFromDXGIDevice = d3d11DLL.NewProc("CreateDirect3D11DeviceFromDXGIDevice")

func CreateDirect3D11DeviceFromDXGIDevice(dxgiDevice *IDXGIDevice, graphicsDevice **ole.IInspectable) error {
	r1, _, err := pCreateDirect3D11DeviceFromDXGIDevice.Call(uintptr(unsafe.Pointer(dxgiDevice)), uintptr(unsafe.Pointer(graphicsDevice)))
	if r1 != win.S_OK {
		return err
	}
	return nil
}
