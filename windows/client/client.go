package client

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	relative_input "github.com/TKMAX777/RDPRelativeInput"
	"github.com/TKMAX777/RDPRelativeInput/debug"
	"github.com/TKMAX777/RDPRelativeInput/protocol"
	"github.com/TKMAX777/RDPRelativeInput/remote_send"
	"github.com/TKMAX777/RDPRelativeInput/winapi"
	"github.com/lxn/win"
)

func StartClient() {
	// The add-in supplies both the process ID and a pipe for the handshake.
	// A GUI executable otherwise waits invisibly when started from Explorer.
	stdinInfo, stdinErr := os.Stdin.Stat()
	parentPID, pidErr := strconv.ParseUint(os.Getenv("RDP_CLIENT_PID"), 10, 32)
	if pidErr != nil || parentPID == 0 || stdinErr != nil || stdinInfo.Mode()&os.ModeNamedPipe == 0 {
		message := "This client is started automatically by the Remote Desktop add-in.\n\n" +
			"1. Install RelativeInputClient.exe and RelativeInput.dll together using install.bat (Run as administrator).\n" +
			"2. Close all Remote Desktop Connection (mstsc.exe) windows, then reconnect.\n" +
			"3. Run RelativeInputServer.exe inside the remote PC's session.\n\n" +
			"Do not start RelativeInputClient.exe manually."
		if stdinErr == nil && stdinInfo.Mode()&os.ModeNamedPipe != 0 && (pidErr != nil || parentPID == 0) {
			message = "The client could not identify the Remote Desktop process.\n" +
				"An older add-in DLL may still be installed or loaded.\n\n" +
				"Close all mstsc.exe windows, install the matching RelativeInputClient.exe and RelativeInput.dll together, then reconnect and run RelativeInputServer.exe on the remote PC."
		}
		win.MessageBox(0,
			winapi.MustUTF16PtrFromString(message),
			winapi.MustUTF16PtrFromString("RDP Relative Input - Startup help"), win.MB_OK|win.MB_ICONINFORMATION)
		return
	}
	defer os.Stderr.Write([]byte("CLOSE\n"))

	debug.Debugln("==== START CLIENT APPLICATION ====")
	debug.Debugln("ServerProtocolVersion:", relative_input.PROTOCOL_VERSION)

	debug.Debugf("Wait for client headers...")

	if err := protocol.ReadVersion(os.Stdin, relative_input.PROTOCOL_VERSION); err != nil {
		debug.Debugln("Read protocol version:", err)
		fmt.Print("RDPRelativeInput\nStatus:INVALID_PROTOCOL_VERSION\n")
		return
	}

	debug.Debugln("ok")

	// response
	fmt.Printf("RDPRelativeInput\n")
	fmt.Printf("Status:OK\n")

	debug.Debugln("SendStatus:OK")

	var rHandler = remote_send.New(os.Stdout)
	var wHandler = New(rHandler)

	var toggleKey = os.Getenv("RELATIVE_INPUT_TOGGLE_KEY")
	if toggleKey == "" {
		toggleKey = "F8"
	}

	var toggleType = os.Getenv("RELATIVE_INPUT_TOGGLE_TYPE")
	switch toggleType {
	case "ONCE":
		wHandler.SetToggleType(ToggleTypeOnce)
	default:
		wHandler.SetToggleType(ToggleTypeAlive)
	}

	wHandler.SetToggleKey(toggleKey)

	rdHwnd, err := findRDPWindow(uint32(parentPID))
	if err != nil {
		win.MessageBox(0, winapi.MustUTF16PtrFromString(err.Error()), winapi.MustUTF16PtrFromString("RDP Relative Input"), win.MB_ICONERROR)
		return
	}

	_, err = wHandler.StartClient(rdHwnd)
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(os.Stderr, "Ready for sending messages")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	wHandler.Close()
}
