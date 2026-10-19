package host

import (
	"github.com/TKMAX777/RDPRelativeInput/remote_send"
	"github.com/TKMAX777/RDPRelativeInput/winapi"
	"github.com/lxn/win"
)

func makeKeyboardInput(key uint32, state remote_send.InputType) (win.KEYBD_INPUT, bool) {
	input := win.KEYBD_INPUT{Type: win.INPUT_KEYBOARD}
	if key == 0 || key > 254 || (state != remote_send.KeyDown && state != remote_send.KeyUp) {
		return input, false
	}
	scan := winapi.MapVirtualKey(key, winapi.MAPVK_VK_TO_VSC_EX)
	// KEYEVENTF_EXTENDEDKEY represents E0, not E1 (Pause). Use virtual-key
	// injection for E1 and keys that have no scan-code mapping.
	if scan == 0 || scan>>8 == 0xe1 {
		input.Ki.WVk = uint16(key)
	} else {
		input.Ki.WScan = uint16(scan & 0xff)
		input.Ki.DwFlags = win.KEYEVENTF_SCANCODE
		if scan>>8 == 0xe0 || isExtendedVirtualKey(key) {
			input.Ki.DwFlags |= win.KEYEVENTF_EXTENDEDKEY
		}
	}
	if state == remote_send.KeyUp {
		input.Ki.DwFlags |= win.KEYEVENTF_KEYUP
	}
	return input, true
}

// Windows can return the unprefixed keypad scan code for navigation VKs,
// even with MAPVK_VK_TO_VSC_EX. V3 transmits virtual keys, so navigation VKs
// must explicitly select the dedicated navigation block.
func isExtendedVirtualKey(key uint32) bool {
	switch key {
	case 0x21, 0x22, 0x23, 0x24, // Page Up, Page Down, End, Home
		0x25, 0x26, 0x27, 0x28, // arrows
		0x2c, 0x2d, 0x2e, // Print Screen, Insert, Delete
		0x5b, 0x5c, 0x5d, // Windows / Apps
		0x6f, 0x90, // keypad Divide / NumLock
		0xa3, 0xa5: // right Ctrl / Alt
		return true
	}
	return false
}
