package main

import (
	"bytes"
	"io"
	"runtime"
	"time"
	"unsafe"

	"github.com/TKMAX777/RDPRelativeInput/debug"
	"github.com/TKMAX777/RDPRelativeInput/winapi"
	"github.com/lxn/win"
	"github.com/pkg/errors"
)
import "C"

// var OpenChannel uint32
var OpenChannel *VirtualChannelReadWriteCloser
var ServerName string

func ChannelInitEventFn(ppInitHandle *HANDLE, ChannelEvent CHANNEL_INIT_EVENT_TYPE, data *byte, dataLength uint32) uintptr {
	defer func() {
		var rec = recover()
		if rec != nil {
			err := errors.Errorf("Panic: ChannelInitEventFn: %v", rec)
			win.MessageBox(win.HWND(winapi.NULL), winapi.MustUTF16PtrFromString(err.Error()), winapi.MustUTF16PtrFromString("Relative Input"), win.MB_ICONERROR|win.MB_OK)
			debug.Debugln(err)
		}
	}()

	switch ChannelEvent {
	case CHANNEL_EVENT_INITIALIZED:
		debug.Debugln("INITALIZED")
	case CHANNEL_EVENT_CONNECTED:
		ServerName = winapi.UTF16PtrToString((*uint16)(unsafe.Pointer(data)))
		debug.Debugf("Server: %s\n", ServerName)

		var err error
		OpenChannel, err = NewChannel(EntryPoints, ppInitHandle, CHANNEL_NAME)
		if err != nil {
			err = errors.Wrap(err, "VirtualChannelOpen")
			debug.Debugln(err.Error())
			win.MessageBox(win.HWND(winapi.NULL), winapi.MustUTF16PtrFromString(err.Error()), winapi.MustUTF16PtrFromString("Relative Input"), win.MB_ICONERROR|win.MB_OK)
			return 0
		}
		debug.Debugln("Channel Opened")

		var MustReadln = func() string {
			var b = make([]byte, 1)
			var res = make([]byte, 0, 100)
			// var err error

			for {
				n, err := OpenChannel.Read(b)
				if err != nil && err != io.EOF {
					debug.Debugln("READ ERROR: ", err)
				}
				if err == io.EOF {
					time.Sleep(time.Second)
				}
				if bytes.Contains(b[:n], []byte{'\n'}) {
					break
				}

				res = append(res, b[:n]...)
			}

			res = bytes.TrimSuffix(res, []byte{'\n'})

			return string(res)
		}

		go func() {
			// get a client header
			debug.Debugln("Scan started")

			// Response
			for {
				var line = MustReadln()
				if line != "RDPRelativeInput" {
					debug.Debugln(line)
					continue
				}
				debug.Debugln("Run Application")
				StartApplication(OpenChannel, ServerName)
				debug.Debugln("==== CLIENT Application CLOSED ====")
			}
		}()
	case CHANNEL_EVENT_V1_CONNECTED:
		err := errors.Errorf("Non-Windows2000TerminalServerError")
		debug.Debugln(err.Error())
	case CHANNEL_EVENT_REMOTE_CONTROL_START:
		debug.Debugln("Start Remote Session")
	case CHANNEL_EVENT_DISCONNECTED:
		debug.Debugln("CHANNEL_EVENT_DISCONNECTED")
	case CHANNEL_EVENT_TERMINATED:
		debug.Debugln("Terminated Remote Session")
		runtime.GC()
	}

	return 0
}
