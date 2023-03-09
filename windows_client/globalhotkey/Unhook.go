package main

import "C"
import (
	"github.com/TKMAX777/RDPRelativeInput/debug"
	"github.com/TKMAX777/RDPRelativeInput/winapi"
	"github.com/lxn/win"
)

//export Unhook
func Unhook() {
	win.MessageBox(0, winapi.MustUTF16PtrFromString("Unload"), winapi.MustUTF16PtrFromString("Toggle Hook Error"), win.MB_ICONERROR)
	winapi.UnhookWindowsHookEx(hHk)
	debug.Debugln("Unhook...done")
}
