package winapi

import (
	"syscall"
	"unsafe"

	"github.com/lxn/win"
	"github.com/pkg/errors"
)

type WTS_CONNECTSTATE_CLASS int32

const (
	WTS_CURRENT_SESSION = 0xffffffff
)

const (
	WTSCONNECTSTATEActive WTS_CONNECTSTATE_CLASS = iota
	WTSCONNECTSTATEConnected
	WTSCONNECTSTATEConnectQuery
	WTSCONNECTSTATEShadow
	WTSCONNECTSTATEDisconnected
	WTSCONNECTSTATEIdle
	WTSCONNECTSTATEListen
	WTSCONNECTSTATEReset
	WTSCONNECTSTATEDown
	WTSCONNECTSTATEInit
)

type WTS_SESSION_INFO_1 struct {
	ExecEnvID   uint32
	State       WTS_CONNECTSTATE_CLASS
	SessionID   uint32
	SessionName string
	HostName    string
	UserName    string
	DomainName  string
	FarmName    string
}

type wts_SESSION_INFO_1 struct {
	ExecEnvID    uint32
	State        WTS_CONNECTSTATE_CLASS
	SessionID    uint32
	pSessionName *uint16
	pHostName    *uint16
	pUserName    *uint16
	pDomainName  *uint16
	pFarmName    *uint16
}

type WTS_CHANNEL_OPTION_DYNAMIC uint32

const (
	WTS_CHANNEL_OPTION_DYNAMIC_PRI_LOW     WTS_CHANNEL_OPTION_DYNAMIC = 0
	WTS_CHANNEL_OPTION_DYNAMIC_DYNAMIC     WTS_CHANNEL_OPTION_DYNAMIC = 1
	WTS_CHANNEL_OPTION_DYNAMIC_PRI_MED     WTS_CHANNEL_OPTION_DYNAMIC = 2
	WTS_CHANNEL_OPTION_DYNAMIC_PRI_HIGH    WTS_CHANNEL_OPTION_DYNAMIC = 4
	WTS_CHANNEL_OPTION_DYNAMIC_PRI_REAL    WTS_CHANNEL_OPTION_DYNAMIC = 6
	WTS_CHANNEL_OPTION_DYNAMIC_NO_COMPRESS WTS_CHANNEL_OPTION_DYNAMIC = 8
)

type WTS_TYPE_CLASS int

const (
	WTSTypeProcessInfoLevel0 WTS_TYPE_CLASS = iota
	WTSTypeProcessInfoLevel1
	WTSTypeSessionInfoLevel1
)

func WTSFreeMemoryExW(typeClass WTS_TYPE_CLASS, pMemory uintptr, NumEntries int) (err error) {
	return wtsFreeMemoryEx(uintptr(typeClass), pMemory, uint32(NumEntries))
}

type MemoryFreeFunc func() error

func WTSOpenServerExW(pServerName string) win.HANDLE {
	ptr, err := syscall.UTF16PtrFromString(pServerName)
	if err != nil {
		panic(err)
	}
	return win.HANDLE(wtsOpenServerExW(ptr))
}

func WTSCloseServer(hServer win.HANDLE) {
	wtsCloseServerExW(uintptr(hServer))
}

// WTSEnumerateSessionsEx gets sesion list
func WTSEnumerateSessionsEx(hServer win.HANDLE, pLevel *uint32, Filter uint32) (ppSessionInfo []WTS_SESSION_INFO_1, err error) {
	var pCount uint32 = 0
	var sessions = new(*wts_SESSION_INFO_1)

	err = wtsEnumerateSessionsEx(uintptr(hServer), pLevel, Filter, uintptr(unsafe.Pointer(sessions)), &pCount)
	if err != nil {
		return nil, errors.Wrap(err, "wtsEnumerateSessionsEx")
	}

	ppSessionInfo = make([]WTS_SESSION_INFO_1, pCount)

	for i := range ppSessionInfo {
		var s = *(*wts_SESSION_INFO_1)(unsafe.Pointer(
			uintptr(unsafe.Pointer(*sessions)) + uintptr(unsafe.Sizeof(wts_SESSION_INFO_1{}))*uintptr(i),
		))
		ppSessionInfo[i] = WTS_SESSION_INFO_1{
			ExecEnvID:   s.ExecEnvID,
			State:       s.State,
			SessionID:   s.SessionID,
			SessionName: UTF16PtrToString(s.pSessionName),
			HostName:    UTF16PtrToString(s.pHostName),
			UserName:    UTF16PtrToString(s.pUserName),
			DomainName:  UTF16PtrToString(s.pDomainName),
			FarmName:    UTF16PtrToString(s.pFarmName),
		}
	}

	return ppSessionInfo, err
}

func WTSVirtualChannelOpenEx(SessionId uint32, pVirtualName string, flags WTS_CHANNEL_OPTION_DYNAMIC) (handle win.HANDLE, err error) {
	name := []byte(pVirtualName)
	h, err := wtsVirtualChannelOpenEx(SessionId, &name[0], uint32(flags))
	return win.HANDLE(h), err
}

func WTSVirtualChannelWrite(hChannelHandle win.HANDLE, Buffer []byte, Length int, pBytesWritten *uint32) (err error) {
	if Length == 0 {
		return nil
	}
	return wtsVirtualChannelWrite(uintptr(hChannelHandle), uintptr(unsafe.Pointer(&Buffer[0])), uint32(Length), pBytesWritten)
}

func WTSVirtualChannelRead(hChannelHandle win.HANDLE, TimeOut uint32, Buffer []byte, BufferSize int, pBytesRead *uint32) (err error) {
	if BufferSize == 0 {
		return nil
	}
	return wtsVirtualChannelRead(uintptr(hChannelHandle), TimeOut, uintptr(unsafe.Pointer(&Buffer[0])), uint32(BufferSize), pBytesRead)
}

func WTSVirtualChannelClose(hChannelHandle win.HANDLE) (err error) {
	return wtsVirtualChannelClose(uintptr(hChannelHandle))
}

type WTSVirtualChannelReadCloser struct {
	ChannelHandle win.HANDLE

	// 0xFFFFFFFF will be infinite
	Timeout uint32
}

func OpenWTSVirtualChannel(sessionid uint32, VirtualChannelName string, flags WTS_CHANNEL_OPTION_DYNAMIC) (*WTSVirtualChannelReadCloser, error) {
	handle, err := WTSVirtualChannelOpenEx(sessionid, VirtualChannelName, flags)
	return &WTSVirtualChannelReadCloser{ChannelHandle: handle}, err
}

// buffer should have more than 1600 byte lengths.
func (rw WTSVirtualChannelReadCloser) Read(b []byte) (n int, err error) {
	var count uint32
	err = WTSVirtualChannelRead(rw.ChannelHandle, rw.Timeout, b, len(b), &count)
	return int(count), err
}

func (rw WTSVirtualChannelReadCloser) Write(b []byte) (n int, err error) {
	var count uint32
	err = WTSVirtualChannelWrite(rw.ChannelHandle, b, len(b), &count)
	return int(count), err
}

func (rw WTSVirtualChannelReadCloser) Close() error {
	return WTSVirtualChannelClose(rw.ChannelHandle)
}
