package main

import (
	"fmt"
	"sync"

	relative_input "github.com/TKMAX777/RDPRelativeInput"
	"github.com/TKMAX777/RDPRelativeInput/debug"
	"github.com/TKMAX777/RDPRelativeInput/protocol"
)
import "C"

var connections = struct {
	sync.Mutex
	channels map[*HANDLE]*VirtualChannelReadWriteCloser
}{channels: make(map[*HANDLE]*VirtualChannelReadWriteCloser)}

func ChannelInitEventFn(handle *HANDLE, event CHANNEL_INIT_EVENT_TYPE, _ *byte, _ uint32) uintptr {
	switch event {
	case CHANNEL_EVENT_INITIALIZED:
		debug.Debugln("INITIALIZED")
	case CHANNEL_EVENT_CONNECTED:
		rw, err := NewChannel(EntryPoints, handle, CHANNEL_NAME)
		if err != nil {
			debug.Debugln("Open channel:", err)
			return 0
		}
		connections.Lock()
		old := connections.channels[handle]
		connections.channels[handle] = rw
		connections.Unlock()
		if old != nil {
			old.received.Fail(fmt.Errorf("connection replaced"))
		}
		debug.Debugln("Channel Opened")
		go serveChannel(rw)
	case CHANNEL_EVENT_DISCONNECTED, CHANNEL_EVENT_TERMINATED:
		connections.Lock()
		rw := connections.channels[handle]
		delete(connections.channels, handle)
		connections.Unlock()
		if rw != nil {
			rw.received.Fail(fmt.Errorf("connection disconnected"))
		}
	}
	return 0
}

func serveChannel(rw *VirtualChannelReadWriteCloser) {
	defer rw.Close()
	// Closing the channel cancels a partial handshake as well as a running child.
	reader := protocol.NewReader(rw)
	for {
		header, err := reader.ReadLine()
		if err != nil {
			debug.Debugln("Read header:", err)
			return
		}
		if header != "RDPRelativeInput" {
			debug.Debugln("Invalid channel header")
			return
		}
		version, err := reader.ReadLine()
		if err != nil {
			debug.Debugln("Read version:", err)
			return
		}
		if version != fmt.Sprintf("ProtocolVersion:%d", relative_input.PROTOCOL_VERSION) {
			_, _ = rw.Write([]byte("RDPRelativeInput\nStatus:INVALID_PROTOCOL_VERSION\n"))
			return
		}
		debug.Debugln("Run Application")
		if err := StartApplication(rw, version); err != nil {
			debug.Debugln("Client application:", err)
			return
		}
	}
}
