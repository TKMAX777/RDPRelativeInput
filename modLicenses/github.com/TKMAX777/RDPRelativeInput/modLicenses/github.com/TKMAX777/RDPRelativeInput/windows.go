//go:build windows
// +build windows

package relative_input

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	"github.com/TKMAX777/RDPRelativeInput/debug"
	"github.com/TKMAX777/RDPRelativeInput/keymap"
	"github.com/TKMAX777/RDPRelativeInput/remote_send"
	"github.com/TKMAX777/RDPRelativeInput/winapi"
	"github.com/TKMAX777/RDPRelativeInput/windows_client/client"
	"github.com/lxn/win"
)

const PROTOCOL_VERSION = 1

func StartClient() {
	defer os.Stderr.Write([]byte("CLOSE\n"))

	debug.Debugln("==== START CLIENT APPLICATION ====")
	debug.Debugln("ServerProtocolVersion:", PROTOCOL_VERSION)

	debug.Debugf("Wait for client headers...")

	var scanner = bufio.NewScanner(os.Stdin)

restart:
	scanner.Scan()
	var line = scanner.Text()

	if !strings.HasSuffix(line, strconv.Itoa(PROTOCOL_VERSION)) {
		debug.Debugln("error!")
		debug.Debugln("Get: ", line)
		debug.Debugln("SendStatus:INVALID_PROTOCOL_VERSION")
		fmt.Printf("Status:INVALID_PROTOCOL_VERSION\n")
		goto restart
	}

	debug.Debugln("ok")

	// response
	fmt.Printf("RDPRelativeInput\n")
	fmt.Printf("Status:OK\n")

	debug.Debugln("SendStatus:OK")

	var rHandler = remote_send.New(os.Stdout)
	var wHandler = client.New(rHandler)

	var toggleKey = os.Getenv("RELATIVE_INPUT_TOGGLE_KEY")
	if toggleKey == "" {
		toggleKey = "F8"
	}

	var toggleType = os.Getenv("RELATIVE_INPUT_TOGGLE_TYPE")
	switch toggleType {
	case "ONCE":
		wHandler.SetToggleType(client.ToggleTypeOnce)
	default:
		wHandler.SetToggleType(client.ToggleTypeAlive)
	}

	wHandler.SetToggleKey(toggleKey)

	var rdHwnd win.HWND
	for {
		rdHwnd = winapi.FindWindowEx(0, rdHwnd, winapi.MustUTF16PtrFromString("TscShellContainerClass"), nil)
		if rdHwnd == 0 {
			win.MessageBox(0, winapi.MustUTF16PtrFromString("Could not find window"), winapi.MustUTF16PtrFromString("RDP Relative Input"), win.MB_ICONERROR)
			debug.Debugln("Window not found error")
			return
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

func StartServer() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	fmt.Println("Start Server")

	debug.Debugln("==== LOGGING START ====")
	debug.Debugf("OpenVirtualChannel...")

	rw, err := winapi.OpenWTSVirtualChannel(winapi.WTS_CURRENT_SESSION, CHANNEL_NAME, 0)
	if err != nil {
		debug.Debugln("error!")
		debug.Debugln(err)
		panic(err)
	}
	defer rw.Close()
	rw.Timeout = 0xffffffff

	debug.Debugln("ok")

	var MustReadln = func() string {
		var b = make([]byte, 1600)
		var res = make([]byte, 0, 100)

		for {
			n, _ := rw.Read(b)
			res = append(res, b[:n]...)
			if bytes.Contains(b[:n], []byte{'\n'}) {
				break
			}
		}

		res = bytes.TrimSuffix(res, []byte{'\n'})

		return string(res)
	}

	// Header
	debug.Debugf("SendHeaders...")

	fmt.Fprintln(rw, "RDPRelativeInput")
	fmt.Fprintf(rw, "ProtocolVersion:%d\n", PROTOCOL_VERSION)

	debug.Debugln("ok")

	debug.Debugf("Wait for host response...")

restart:
	// Response
	var line = MustReadln()
	if line != "RDPRelativeInput" {
		fmt.Println("Get: ", line)
		goto restart
	}

	line = MustReadln()
	if !strings.HasPrefix(line, "Status:") {
		debug.Debugln("error!")
		fmt.Println("Server respond with illigal format: ", line)
		return
	}

	if line != "Status:OK" {
		debug.Debugln("error!")
		debug.Debugln("ProtocolError:", line)
		fmt.Println("Server Protocol Error: ", line)
		return
	}

	debug.Debugln("ok")

	debug.Debugln("Start Connection")
	fmt.Println("Start Connection")

	for {
		text := MustReadln()

		var augs = strings.Split(text, " ")
		if len(augs) < 4 {
			continue
		}

		eventType, err := strconv.ParseUint(augs[0], 10, 32)
		if err != nil {
			continue
		}
		eventInput, err := strconv.ParseUint(augs[1], 10, 32)
		if err != nil {
			continue
		}
		eventValue1, err := strconv.ParseInt(augs[2], 10, 32)
		if err != nil {
			continue
		}
		eventValue2, err := strconv.ParseInt(augs[3], 10, 32)
		if err != nil {
			continue
		}

		switch keymap.EV_TYPE(eventType) {
		case keymap.EV_TYPE_MOUSE_MOVE:
			switch uint32(eventInput) {
			case uint32(remote_send.MouseMoveTypeRelative):
				var dx = eventValue1
				var dy = eventValue2

				var mouseInput = win.MOUSE_INPUT{
					Type: win.INPUT_MOUSE,
					Mi: win.MOUSEINPUT{
						Dx:      int32(dx),
						Dy:      int32(dy),
						DwFlags: win.MOUSEEVENTF_MOVE,
					},
				}

				win.SendInput(1, unsafe.Pointer(&mouseInput), int32(unsafe.Sizeof(win.MOUSE_INPUT{})))
				debug.Debugf("SendInput: MouseREL: dx: %d dy: %d\n", dx, dy)
			case uint32(remote_send.MouseMoveTypeAbsolute):
				var x = eventValue1
				var y = eventValue2

				var mouseInput = win.MOUSE_INPUT{
					Type: win.INPUT_MOUSE,
					Mi: win.MOUSEINPUT{
						Dx:      int32(x),
						Dy:      int32(y),
						DwFlags: win.MOUSEEVENTF_MOVE | win.MOUSEEVENTF_ABSOLUTE,
					},
				}

				win.SendInput(1, unsafe.Pointer(&mouseInput), int32(unsafe.Sizeof(win.MOUSE_INPUT{})))
				debug.Debugf("SendInput: MouseREL: x: %d y: %d\n", x, y)
			}
		case keymap.EV_TYPE_MOUSE:
			switch eventInput {
			case 0x02:
				var mouseInput = win.MOUSE_INPUT{
					Type: win.INPUT_MOUSE,
				}

				switch uint32(eventValue1) {
				case uint32(remote_send.KeyDown):
					mouseInput.Mi = win.MOUSEINPUT{
						DwFlags: win.MOUSEEVENTF_RIGHTDOWN,
					}
					debug.Debugln("SendInput: MouseRightDown")
				case uint32(remote_send.KeyUp):
					mouseInput.Mi = win.MOUSEINPUT{
						DwFlags: win.MOUSEEVENTF_RIGHTUP,
					}
					debug.Debugln("SendInput: MouseRightUp")
				}

				win.SendInput(1, unsafe.Pointer(&mouseInput), int32(unsafe.Sizeof(win.MOUSE_INPUT{})))
			case 0x01:
				var mouseInput = win.MOUSE_INPUT{
					Type: win.INPUT_MOUSE,
				}

				switch uint32(eventValue1) {
				case uint32(remote_send.KeyDown):
					mouseInput.Mi = win.MOUSEINPUT{
						DwFlags: win.MOUSEEVENTF_LEFTDOWN,
					}
					debug.Debugln("SendInput: MouseLeftDown")
				case uint32(remote_send.KeyUp):
					mouseInput.Mi = win.MOUSEINPUT{
						DwFlags: win.MOUSEEVENTF_LEFTUP,
					}
					debug.Debugln("SendInput: MouseLeftUp")
				}

				win.SendInput(1, unsafe.Pointer(&mouseInput), int32(unsafe.Sizeof(win.MOUSE_INPUT{})))
			case 0x04:
				var mouseInput = win.MOUSE_INPUT{
					Type: win.INPUT_MOUSE,
				}

				switch uint32(eventValue1) {
				case uint32(remote_send.KeyDown):
					mouseInput.Mi = win.MOUSEINPUT{
						DwFlags: win.MOUSEEVENTF_MIDDLEDOWN,
					}
					debug.Debugln("SendInput: MouseMiddleDown")
				case uint32(remote_send.KeyUp):
					mouseInput.Mi = win.MOUSEINPUT{
						DwFlags: win.MOUSEEVENTF_MIDDLEUP,
					}
					debug.Debugln("SendInput: MouseMiddleUp")
				}

				win.SendInput(1, unsafe.Pointer(&mouseInput), int32(unsafe.Sizeof(win.MOUSE_INPUT{})))
			}
		case keymap.EV_TYPE_WHEEL:
			var mouseInput = win.MOUSE_INPUT{
				Type: win.INPUT_MOUSE,
				Mi: win.MOUSEINPUT{
					DwFlags: win.MOUSEEVENTF_WHEEL,
				},
			}
			switch uint32(eventValue1) {
			case uint32(remote_send.KeyDown):
				mouseInput.Mi.MouseData = ^uint32(120) + 1
				debug.Debugln("SendInput: MouseMiddleDown")
			case uint32(remote_send.KeyUp):
				mouseInput.Mi.MouseData = 120
				debug.Debugln("SendInput: MouseMiddleUp")
			}

			win.SendInput(1, unsafe.Pointer(&mouseInput), int32(unsafe.Sizeof(win.MOUSE_INPUT{})))

		case keymap.EV_TYPE_KEY:
			var mappedKey = winapi.MapVirtualKey(uint32(eventInput), winapi.MAPVK_VK_TO_VSC)

			var keyInput = win.KEYBD_INPUT{
				Type: win.INPUT_KEYBOARD,
				Ki: win.KEYBDINPUT{
					WScan: uint16(mappedKey),
				},
			}

			switch uint32(eventValue1) {
			case uint32(remote_send.KeyDown):
				keyInput.Ki.DwFlags = win.KEYEVENTF_SCANCODE
				win.SendInput(1, unsafe.Pointer(&keyInput), int32(unsafe.Sizeof(win.KEYBD_INPUT{})))
			case uint32(remote_send.KeyUp):
				keyInput.Ki.DwFlags = win.KEYEVENTF_KEYUP | win.KEYEVENTF_SCANCODE
				win.SendInput(1, unsafe.Pointer(&keyInput), int32(unsafe.Sizeof(win.KEYBD_INPUT{})))
			}
		}
	}
}
