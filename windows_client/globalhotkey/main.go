package main

import "C"

import (
	"github.com/TKMAX777/RDPRelativeInput/debug"
	"github.com/TKMAX777/RDPRelativeInput/winapi"
)

type HANDLE uintptr

func init() {
	debug.Debugln("====LOGGING START====")
	debug.Debugf("INITIALIZING...")

	// Prevent from　DLL unloading
	var nh HANDLE
	err := GetModuleHandleExW(1, "", &nh)
	if err != nil {
		debug.Debugln("error")
		debug.Debugf("GetModuleHandleExW: %v\n", err)
	}
	debug.Debugln("ok")
}

// for building
// This function is not an entry point of this program.
// The entry point is VirtualChannelEntry function.
func main() {}

var hHk winapi.HHOOK
