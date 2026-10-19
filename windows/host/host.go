//go:build windows
// +build windows

package host

import (
	"fmt"
	"runtime"
	"strings"
	"unsafe"

	relative_input "github.com/TKMAX777/RDPRelativeInput"
	"github.com/TKMAX777/RDPRelativeInput/debug"
	"github.com/TKMAX777/RDPRelativeInput/keymap"
	"github.com/TKMAX777/RDPRelativeInput/protocol"
	"github.com/TKMAX777/RDPRelativeInput/remote_send"
	"github.com/TKMAX777/RDPRelativeInput/winapi"
	"github.com/lxn/win"
)

func StartServer() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	fmt.Println("Start Server")

	debug.Debugln("==== LOGGING START ====")
	debug.Debugf("OpenVirtualChannel...")

	rw, err := winapi.OpenWTSVirtualChannel(winapi.WTS_CURRENT_SESSION, relative_input.CHANNEL_NAME, 0)
	if err != nil {
		debug.Debugln("error!")
		debug.Debugln(err)
		panic(err)
	}
	defer rw.Close()
	rw.Timeout = 0xffffffff

	debug.Debugln("ok")

	reader := protocol.NewReader(rw)

	// Header
	debug.Debugf("SendHeaders...")

	if _, err := fmt.Fprintf(rw, "RDPRelativeInput\nProtocolVersion:%d\n", relative_input.PROTOCOL_VERSION); err != nil {
		debug.Debugln("Send headers:", err)
		return
	}

	debug.Debugln("ok")

	debug.Debugf("Wait for client response...")

restart:
	// Response
	line, err := reader.ReadLine()
	if err != nil {
		debug.Debugln("Read header:", err)
		return
	}
	if line != "RDPRelativeInput" {
		fmt.Println("Get: ", line)
		goto restart
	}

	line, err = reader.ReadLine()
	if err != nil {
		debug.Debugln("Read status:", err)
		return
	}
	if !strings.HasPrefix(line, "Status:") {
		debug.Debugln("error!")
		fmt.Println("Client respond with illigal format: ", line)
		return
	}

	if line != "Status:OK" {
		debug.Debugln("error!")
		debug.Debugln("ProtocolError:", line)
		fmt.Println("Client Protocol Error: ", line)
		return
	}

	debug.Debugln("ok")

	var handler = New()
	if _, err := handler.StartHostButton(); err != nil {
		debug.Debugln("Start host button:", err)
		return
	}

	debug.Debugln("Start Connection")
	fmt.Println("Start Connection")

	for {
		event, err := reader.ReadEvent()
		if err != nil {
			if err != protocol.ErrClosed {
				debug.Debugln("Read event:", err)
			}
			return
		}
		eventType, eventInput := event.Type, event.Input
		eventValue1, eventValue2 := event.Value1, event.Value2

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
			// Mouse Right
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

			// Mouse Left
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

			// Mouse Middle
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
			mouseInput.Mi.MouseData = eventInput

			win.SendInput(1, unsafe.Pointer(&mouseInput), int32(unsafe.Sizeof(win.MOUSE_INPUT{})))

		case keymap.EV_TYPE_KEY:
			keyInput, ok := makeKeyboardInput(eventInput, remote_send.InputType(eventValue1))
			if ok {
				win.SendInput(1, unsafe.Pointer(&keyInput), int32(unsafe.Sizeof(win.KEYBD_INPUT{})))
			}
		}
	}
}
