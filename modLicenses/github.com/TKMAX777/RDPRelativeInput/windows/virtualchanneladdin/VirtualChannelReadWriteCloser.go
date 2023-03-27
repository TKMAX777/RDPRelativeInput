package main

import (
	"bytes"
	"unsafe"

	"github.com/TKMAX777/RDPRelativeInput/debug"
	"github.com/TKMAX777/RDPRelativeInput/winapi"
	"github.com/lxn/win"
	"github.com/pkg/errors"
)
import "C"

const RWBufferSize = 1e3

type VirtualChannelReadWriteCloser struct {
	channelHandle uint32
	entryPoints   CHANNEL_ENTRY_POINTS

	recievedBuffer *bytes.Buffer

	lastWriteMessageID uint32
	chanWriteComplete  map[uint32]chan bool
}

type RWTaskType int

type RWRecievedDataAttributes struct {
	Data []byte
}

func NewChannel(entryPoints CHANNEL_ENTRY_POINTS, ppInitHandle *HANDLE, ChannelName string) (*VirtualChannelReadWriteCloser, error) {
	var rw = &VirtualChannelReadWriteCloser{
		chanWriteComplete: make(map[uint32]chan bool),
		recievedBuffer:    new(bytes.Buffer),
		entryPoints:       entryPoints,
	}

	err := entryPoints.VirtualChannelOpen.Call(ppInitHandle, &rw.channelHandle, ChannelName, rw.getVirtualChannelOpenEventFn())
	if err != nil {
		err = errors.Wrap(err, "VirtualChannelOpen")
		return nil, err
	}

	return rw, nil
}

func (rw *VirtualChannelReadWriteCloser) Write(b []byte) (n int, err error) {
	var messageID = rw.lastWriteMessageID
	rw.lastWriteMessageID = (rw.lastWriteMessageID + 1) % 0xffff

	rw.chanWriteComplete[messageID] = make(chan bool)

	err = rw.entryPoints.VirtualChannelWrite.Call(rw.channelHandle, b, len(b), uintptr(messageID))
	if err != nil {
		delete(rw.chanWriteComplete, messageID)
		return 0, errors.Wrap(err, "VirtualChannelWrite")
	}

	if !<-rw.chanWriteComplete[messageID] {
		return 0, errors.New("Client Cancelled")
	}

	return len(b), nil
}

func (rw *VirtualChannelReadWriteCloser) Read(b []byte) (n int, err error) {
	return rw.recievedBuffer.Read(b)
}

func (rw *VirtualChannelReadWriteCloser) Close() error {
	return rw.entryPoints.VirtualChannelClose.Call(rw.channelHandle)
}

func (rw *VirtualChannelReadWriteCloser) getVirtualChannelOpenEventFn() CHANNEL_OPEN_EVENT_FN {
	var buffer []byte

	var eventFunction CHANNEL_OPEN_EVENT_FN = func(openHandle uint32, event CHANNEL_EVENT_TYPE, pData *byte, dataLength uint32, totalLength uint32, dataFlags CHANNEL_FLAG_TYPE) uintptr {
		defer func() {
			var rec = recover()
			if rec != nil {
				err := errors.Errorf("Panic: virtualChannelOpenEventFn: %v", rec)
				debug.Debugln(err)
				win.MessageBox(win.HWND(winapi.NULL), winapi.MustUTF16PtrFromString(err.Error()), winapi.MustUTF16PtrFromString("Relative Input"), win.MB_ICONERROR|win.MB_OK)
			}
		}()

		switch event {
		case CHANNEL_EVENT_DATA_RECEIVED:
			buffer = unsafe.Slice(pData, dataLength)
			rw.recievedBuffer.Write(buffer)
			// switch dataFlags {
			// case CHANNEL_FLAG_ONLY:
			// 	rw.recievedBuffer.Write(buffer)
			// case CHANNEL_FLAG_FIRST:
			// 	rw.recievedBuffer.Write(buffer)
			// case CHANNEL_FLAG_LAST:
			// 	rw.recievedBuffer.Write(buffer)
			// case CHANNEL_FLAG_MIDDLE:
			// 	rw.recievedBuffer.Write(buffer)
			// }
		case CHANNEL_EVENT_WRITE_CANCELLED:
			rw.chanWriteComplete[uint32(uintptr(unsafe.Pointer(pData)))] <- false
		case CHANNEL_EVENT_WRITE_COMPLETE:
			rw.chanWriteComplete[uint32(uintptr(unsafe.Pointer(pData)))] <- true
		}
		return 0
	}

	return eventFunction
}
