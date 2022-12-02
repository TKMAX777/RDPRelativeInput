package main

import (
	"syscall"
	"unsafe"

	"github.com/pkg/errors"
)
import "C"

const (
	CHANNEL_NAME_LEN  = 7
	CHANNEL_MAX_COUNT = 30
)

type CHANNEL_DEF struct {
	name    [CHANNEL_NAME_LEN + 1]byte
	options uint32
}

type PCHANNEL_DEF *CHANNEL_DEF
type PPCHANNEL_DEF **CHANNEL_DEF

type CHANNEL_INIT_EVENT_TYPE uint

const (
	CHANNEL_EVENT_INITIALIZED CHANNEL_INIT_EVENT_TYPE = iota
	CHANNEL_EVENT_CONNECTED
	CHANNEL_EVENT_V1_CONNECTED
	CHANNEL_EVENT_DISCONNECTED
	CHANNEL_EVENT_TERMINATED
	CHANNEL_EVENT_REMOTE_CONTROL_START
	CHANNEL_EVENT_REMOTE_CONTROL_STOP
)

type CHANNEL_EVENT_TYPE uint

const (
	CHANNEL_EVENT_DATA_RECEIVED CHANNEL_EVENT_TYPE = 10 + iota
	CHANNEL_EVENT_WRITE_COMPLETE
	CHANNEL_EVENT_WRITE_CANCELLED
)

type CHANNEL_FLAG_TYPE uint32

const (
	CHANNEL_FLAG_FIRST         CHANNEL_FLAG_TYPE = 0x01
	CHANNEL_FLAG_LAST          CHANNEL_FLAG_TYPE = 0x02
	CHANNEL_FLAG_ONLY          CHANNEL_FLAG_TYPE = (CHANNEL_FLAG_FIRST | CHANNEL_FLAG_LAST)
	CHANNEL_FLAG_MIDDLE        CHANNEL_FLAG_TYPE = 0
	CHANNEL_FLAG_FAIL          CHANNEL_FLAG_TYPE = 0x100
	CHANNEL_FLAG_SHOW_PROTOCOL CHANNEL_FLAG_TYPE = 0x10
	CHANNEL_FLAG_SUSPEND       CHANNEL_FLAG_TYPE = 0x20
	CHANNEL_FLAG_RESUME        CHANNEL_FLAG_TYPE = 0x40
)

type PVIRTUALCHANNELINIT uintptr

func (p *PVIRTUALCHANNELINIT) Call(ppInitHandle *HANDLE, pChannel []CHANNEL_DEF, channelCount int, verginRequested uint32, pChannelInitEventProc CHANNEL_INIT_EVENT_FN) error {
	r1, _, _ := syscall.SyscallN(uintptr(unsafe.Pointer(p)), uintptr(unsafe.Pointer(ppInitHandle)), uintptr(unsafe.Pointer(&pChannel[0])), uintptr(channelCount), uintptr(verginRequested), syscall.NewCallback(pChannelInitEventProc))
	if r1 != 0 {
		str, ok := CHANNEL_RETURN_CODES[CHANNEL_RETURN_CODE(r1)]
		if !ok {
			return errors.Errorf("Error: %d", r1)
		}
		return errors.New(str)
	}
	return nil
}

type PVIRTUALCHANNELOPEN uintptr

func (p *PVIRTUALCHANNELOPEN) Call(pInitHandle *HANDLE, pOpenHandle *uint32, pChannelName string, pChannelOpenEventProc CHANNEL_OPEN_EVENT_FN) error {
	var name = []byte(pChannelName)

	r1, _, _ := syscall.SyscallN(uintptr(unsafe.Pointer(p)), uintptr(unsafe.Pointer(pInitHandle)), uintptr(unsafe.Pointer(pOpenHandle)), uintptr(unsafe.Pointer(&name[0])), syscall.NewCallback(pChannelOpenEventProc))
	if r1 != 0 {
		str, ok := CHANNEL_RETURN_CODES[CHANNEL_RETURN_CODE(r1)]
		if !ok {
			return errors.Errorf("Error: %d", r1)
		}
		return errors.New(str)
	}
	return nil
}

type PVIRTUALCHANNELCLOSE uintptr

func (p *PVIRTUALCHANNELCLOSE) Call(openHandle uint32) error {
	r1, _, _ := syscall.SyscallN(uintptr(unsafe.Pointer(p)), uintptr(openHandle))
	if r1 != 0 {
		str, ok := CHANNEL_RETURN_CODES[CHANNEL_RETURN_CODE(r1)]
		if !ok {
			return errors.Errorf("Error: %d", r1)
		}
		return errors.New(str)
	}
	return nil
}

type PVIRTUALCHANNELWRITE uintptr

func (p *PVIRTUALCHANNELWRITE) Call(openHandle uint32, pData []byte, dataLength int, pUserData uintptr) error {
	r1, _, _ := syscall.SyscallN(uintptr(unsafe.Pointer(p)), uintptr(openHandle), uintptr(unsafe.Pointer(&pData[0])), uintptr(dataLength), pUserData)
	if r1 != 0 {
		str, ok := CHANNEL_RETURN_CODES[CHANNEL_RETURN_CODE(r1)]
		if !ok {
			return errors.Errorf("Error: %d", r1)
		}
		return errors.New(str)
	}
	return nil
}

type CHANNEL_INIT_EVENT_FN func(ppInitHandle *HANDLE, ChannelEvent CHANNEL_INIT_EVENT_TYPE, data *byte, dataLength uint32) uintptr
type CHANNEL_OPEN_EVENT_FN func(openHandle uint32, event CHANNEL_EVENT_TYPE, pData *byte, dataLength uint32, totalLength uint32, dataFlags CHANNEL_FLAG_TYPE) uintptr

type CHANNEL_ENTRY_POINTS struct {
	cbSize              uint32
	protocolVersion     uint32
	VirtualChannelInit  *PVIRTUALCHANNELINIT
	VirtualChannelOpen  *PVIRTUALCHANNELOPEN
	VirtualChannelClose *PVIRTUALCHANNELCLOSE
	VirtualChannelWrite *PVIRTUALCHANNELWRITE
}

type CHANNEL_OPTION uint

const (
	CHANNEL_OPTION_INITIALIZED   CHANNEL_OPTION = 0x80000000
	CHANNEL_OPTION_ENCRYPT_RDP   CHANNEL_OPTION = 0x40000000
	CHANNEL_OPTION_ENCRYPT_SC    CHANNEL_OPTION = 0x20000000
	CHANNEL_OPTION_ENCRYPT_CS    CHANNEL_OPTION = 0x10000000
	CHANNEL_OPTION_PRI_HIGH      CHANNEL_OPTION = 0x08000000
	CHANNEL_OPTION_PRI_MED       CHANNEL_OPTION = 0x04000000
	CHANNEL_OPTION_PRI_LOW       CHANNEL_OPTION = 0x02000000
	CHANNEL_OPTION_COMPRESS_RDP  CHANNEL_OPTION = 0x00800000
	CHANNEL_OPTION_COMPRESS      CHANNEL_OPTION = 0x00400000
	CHANNEL_OPTION_SHOW_PROTOCOL CHANNEL_OPTION = 0x00200000
	REMOTE_CONTROL_PERSISTENT    CHANNEL_OPTION = 0x00100000
)

type CHANNEL_RETURN_CODE uintptr

var CHANNEL_RETURN_CODES = map[CHANNEL_RETURN_CODE]string{
	0:  "CHANNEL_RC_OK",
	1:  "CHANNEL_RC_ALREADY_INITIALIZED",
	2:  "CHANNEL_RC_NOT_INITIALIZED",
	3:  "CHANNEL_RC_ALREADY_CONNECTED",
	4:  "CHANNEL_RC_NOT_CONNECTED",
	5:  "CHANNEL_RC_TOO_MANY_CHANNELS",
	6:  "CHANNEL_RC_BAD_CHANNEL",
	7:  "CHANNEL_RC_BAD_CHANNEL_HANDLE",
	8:  "CHANNEL_RC_NO_BUFFER",
	9:  "CHANNEL_RC_BAD_INIT_HANDLE",
	10: "CHANNEL_RC_NOT_OPEN",
	11: "CHANNEL_RC_BAD_PROC",
	12: "CHANNEL_RC_NO_MEMORY",
	13: "CHANNEL_RC_UNKNOWN_CHANNEL_NAME",
	14: "CHANNEL_RC_ALREADY_OPEN",
	15: "CHANNEL_RC_NOT_IN_VIRTUALCHANNELENTRY",
	16: "CHANNEL_RC_NULL_DATA",
	17: "CHANNEL_RC_ZERO_LENGTH",
}
