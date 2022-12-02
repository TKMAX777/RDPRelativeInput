package main

import "C"
import "github.com/TKMAX777/RDPRelativeInput/debug"

const CHANNEL_NAME = "RELINP"

var ChannelName [CHANNEL_NAME_LEN + 1]byte

func init() {
	debug.Debugln("====LOGGING START====")
	debug.Debugf("INITIALIZING...")
	for i, char := range CHANNEL_NAME {
		ChannelName[i] = byte(char)
		ChannelName[i+1] = 0
	}

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
func main() {

}
