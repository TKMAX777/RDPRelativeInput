package winapi

const NULL uintptr = 0

//go:generate go run golang.org/x/sys/windows/mkwinsyscall -output windows_generate.go windows.go
//sys clipCursor(rect uintptr)(ok int, err error) = user32.ClipCursor
//sys showCursor(state bool) (counter int) = user32.ShowCursor
//sys enumDesktopWindows(hDesktop uintptr, lpEnumFunc uintptr, lParam uintptr) (err error) = user32.EnumDesktopWindows
//sys setLayeredWindowAttributes(hwnd uintptr, color uint32, bAlpha byte, dwFlags uint32) (err error) = user32.SetLayeredWindowAttributes
//sys fillRect(hdc uintptr, rect uintptr, hbr uintptr) (err error) = user32.FillRect
//sys setWindowRgn(hwnd uintptr, hRgn uintptr, bRedraw bool) (err error) = user32.SetWindowRgn
//sys updateLayeredWindow(hwnd uintptr, hdcDst uintptr, pptDst uintptr, psize uintptr, hdcSrc uintptr, pptSrc uintptr, crKey uint32, pblend uintptr, dwFlags uint32) (ok bool) = user32.UpdateLayeredWindow
//sys findWindowEx(hwndParent uintptr, hwndChildAfter uintptr, lpszClass *uint16, lpszWindow *uint16) (hwnd uintptr) = user32.FindWindowExW
//sys getWindowText(hwnd uintptr, lpString uintptr, nMax int) (length int) = user32.GetWindowTextW
//sys getClassName(hwnd uintptr, lpClassName uintptr, nMax int) (length int) = user32.GetClassNameW
//sys setWindowText(hwnd uintptr, lpString *uint16) (err error) = user32.SetWindowTextW
//sys invalidateRect(hwnd uintptr, rect uintptr, bErase bool) (err error) = user32.InvalidateRect
//sys mapVirtualKey(uCode uint32, uMapType uint32) (code uint32) = user32.MapVirtualKeyW
//sys registerClassEx(windowClass uintptr) (atom uint16, err error) = user32.RegisterClassExW

//sys createSolidBrush(color uint32) (hbrush uintptr) = Gdi32.CreateSolidBrush
//sys createPen(iStyle int, cWidth int, color uint32) (hpen uintptr) = Gdi32.CreatePen
//sys polyDraw(hdc uintptr, apt uintptr, aj uintptr, cpt int) (err error) = Gdi32.PolyDraw
//sys createRectRgnIndirect(rect uintptr) (rgn uintptr) = Gdi32.CreateRectRgnIndirect
//sys createDIBSection(hdc uintptr, pbmi uintptr, usage uint, ppvBits uintptr, hSection uintptr, offset uint32) (hBitMap uintptr) = Gdi32.CreateDIBSection
//sys extFloodFill(hdc uintptr, x int, y int, color uint32, opType uint32) (err error) = Gdi32.ExtFloodFill

//sys activateAudioInterfaceAsync(deviceInterfacePath *uint16, riid uintptr, activationParams uintptr, completionHandler uintptr, createAsync uintptr) (hresult int32) = Mmdevapi.ActivateAudioInterfaceAsync

//sys wtsOpenServerExW(pServerName *uint16) (handle uintptr) = Wtsapi32.WTSOpenServerExW
//sys wtsCloseServerExW(hServer uintptr) = Wtsapi32.WTSCloseServer
//sys wtsEnumerateSessionsEx(hServer uintptr, pLevel *uint32, Filter uint32, ppSessionInfo uintptr, pCount *uint32) (err error) = Wtsapi32.WTSEnumerateSessionsExW
//sys wtsVirtualChannelOpen(hServer uintptr, SessionId uint32, pVirtualName *uint16) (handle uintptr, err error) = Wtsapi32.WTSVirtualChannelOpen
//sys wtsVirtualChannelOpenEx(SessionId uint32, pVirtualName *byte, flags uint32) (handle uintptr, err error) = Wtsapi32.WTSVirtualChannelOpenEx
//sys wtsVirtualChannelClose(hChannelHandle uintptr) (err error) = Wtsapi32.WTSVirtualChannelClose
//sys wtsVirtualChannelWrite(hChannelHandle uintptr, Buffer uintptr, Length uint32, pBytesWritten *uint32) (err error)  = Wtsapi32.WTSVirtualChannelWrite
//sys wtsVirtualChannelRead(hChannelHandle uintptr, TimeOut uint32, Buffer uintptr, BufferSize uint32, pBytesRead *uint32) (err error) = Wtsapi32.WTSVirtualChannelRead
//sys wtsFreeMemoryEx(wtsTypeClass uintptr, pMemory uintptr, NumberOfEntries uint32) (err error) = Wtsapi32.WTSFreeMemoryExW
