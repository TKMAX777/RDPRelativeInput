package client

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	relative_input "github.com/TKMAX777/RDPRelativeInput"
	"github.com/TKMAX777/RDPRelativeInput/debug"
	"github.com/TKMAX777/RDPRelativeInput/remote_send"
	"github.com/TKMAX777/RDPRelativeInput/winapi"
	windowselecter "github.com/TKMAX777/RDPRelativeInput/windows/client/window_selecter"
	"github.com/lxn/win"
)

func StartClient() {
	defer os.Stderr.Write([]byte("CLOSE\n"))

	debug.Debugln("==== START CLIENT APPLICATION ====")
	debug.Debugln("ServerProtocolVersion:", relative_input.PROTOCOL_VERSION)

	debug.Debugf("Wait for client headers...")

	var scanner = bufio.NewScanner(os.Stdin)

	for {
		scanner.Scan()
		var line = scanner.Text()

		if strings.HasSuffix(line, strconv.Itoa(relative_input.PROTOCOL_VERSION)) {
			break
		}

		debug.Debugln("error!")
		debug.Debugln("Get: ", line)
		debug.Debugln("SendStatus:INVALID_PROTOCOL_VERSION")

		// send invalid protocol version response
		fmt.Printf("Status:INVALID_PROTOCOL_VERSION\n")
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

	var rdHwnd win.HWND
	for {
		rdHwnd = winapi.FindWindowEx(0, rdHwnd, winapi.MustUTF16PtrFromString("TscShellContainerClass"), nil)
		if rdHwnd == 0 {
			var err error
			rdHwnd, err = windowselecter.Dialog()
			if err == windowselecter.ErrorQuitted {
				return
			}
			if err != nil {
				win.MessageBox(0, winapi.MustUTF16PtrFromString("Could not open window selecter"), winapi.MustUTF16PtrFromString("RDP Relative Input"), win.MB_ICONERROR)
				debug.Debugln("Window Selecter Error: ", err)
				return
			}
			break
		}
		var name = winapi.GetWindowTextString(rdHwnd)
		if strings.Contains(name, os.Getenv("SERVER_NAME")) {
			debug.Debugln("Client window found: ", name)
			break
		}
	}

	_, err := wHandler.StartClient(rdHwnd)
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(os.Stderr, "Ready for sending messages")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	wHandler.Close()
}
