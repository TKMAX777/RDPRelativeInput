package main

import (
	"unsafe"

	"github.com/TKMAX777/RDPRelativeInput/debug"
	"github.com/TKMAX777/RDPRelativeInput/winapi"
	"github.com/lxn/win"
	"github.com/pkg/errors"
)
import "C"

type HANDLE uintptr

var ChannelHandle HANDLE

var EntryPoints CHANNEL_ENTRY_POINTS

//export VirtualChannelEntry
func VirtualChannelEntry(pEntryPoints *uintptr) bool {
	defer func() {
		var rec = recover()
		if rec != nil {
			err := errors.Errorf("Panic: VirtualChannelEntry: %v", rec)
			win.MessageBox(win.HWND(winapi.NULL), winapi.MustUTF16PtrFromString(err.Error()), winapi.MustUTF16PtrFromString("Relative Input"), win.MB_ICONERROR|win.MB_OK)
			debug.Debugln(err)
		}
	}()

	EntryPoints = *(*CHANNEL_ENTRY_POINTS)(unsafe.Pointer(pEntryPoints))

	debug.Debugf("EntryPoints: %+v\n", EntryPoints)

	var cd = []CHANNEL_DEF{
		{
			name:    ChannelName,
			options: 0,
		},
	}
	debug.Debugf("cd: %+v\n", cd)

	err := EntryPoints.VirtualChannelInit.Call(
		&ChannelHandle, cd, 1, 1, ChannelInitEventFn,
	)
	if err != nil {
		err = errors.Wrap(err, "VirtualChannelInit")
		debug.Debugln(err.Error())
		win.MessageBox(win.HWND(winapi.NULL), winapi.MustUTF16PtrFromString(err.Error()), winapi.MustUTF16PtrFromString("Relative Input"), win.MB_ICONERROR|win.MB_OK)
		return false
	}

	return true
}
