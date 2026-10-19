package host

import (
	"testing"

	"github.com/TKMAX777/RDPRelativeInput/remote_send"
	"github.com/lxn/win"
)

// Exercise Windows' real mapping without injecting keys or changing NumLock.
// Each pair shares its scan byte; the E0 flag must distinguish navigation keys
// from the numeric keypad on both press and release.
func TestNavigationAndKeypadAreDistinct(t *testing.T) {
	for _, pair := range []struct {
		navigation, keypad uint32
		scan               uint16
	}{
		{0x24, 0x67, 0x47}, // Home / numpad 7
		{0x26, 0x68, 0x48}, // Up / numpad 8
		{0x21, 0x69, 0x49}, // Page Up / numpad 9
		{0x25, 0x64, 0x4b}, // Left / numpad 4
		{0x27, 0x66, 0x4d}, // Right / numpad 6
		{0x23, 0x61, 0x4f}, // End / numpad 1
		{0x28, 0x62, 0x50}, // Down / numpad 2
		{0x22, 0x63, 0x51}, // Page Down / numpad 3
		{0x2d, 0x60, 0x52}, // Insert / numpad 0
		{0x2e, 0x6e, 0x53}, // Delete / decimal
	} {
		for _, state := range []remote_send.InputType{remote_send.KeyDown, remote_send.KeyUp} {
			for _, key := range []uint32{pair.navigation, pair.keypad} {
				input, ok := makeKeyboardInput(key, state)
				flags := uint32(win.KEYEVENTF_SCANCODE)
				if key == pair.navigation {
					flags |= win.KEYEVENTF_EXTENDEDKEY
				}
				if state == remote_send.KeyUp {
					flags |= win.KEYEVENTF_KEYUP
				}
				if !ok || input.Ki.WScan != pair.scan || input.Ki.DwFlags != flags || input.Ki.WVk != 0 {
					t.Fatalf("key %x state %d: %+v", key, state, input)
				}
			}
		}
	}
}

func TestKeyboardSpecialAndInvalidKeys(t *testing.T) {
	for _, key := range []uint32{0x2c, 0x5b, 0x5c, 0x5d, 0x6f, 0x90, 0xa3, 0xa5} {
		for _, state := range []remote_send.InputType{remote_send.KeyDown, remote_send.KeyUp} {
			input, ok := makeKeyboardInput(key, state)
			if !ok || input.Ki.DwFlags&win.KEYEVENTF_EXTENDEDKEY == 0 {
				t.Fatalf("extended key %x state %d: %+v", key, state, input)
			}
		}
	}
	for _, key := range []uint32{0, 255, 0xffffffff} {
		if _, ok := makeKeyboardInput(key, remote_send.KeyDown); ok {
			t.Fatalf("accepted %x", key)
		}
	}
	if _, ok := makeKeyboardInput(0x41, 2); ok {
		t.Fatal("accepted invalid state")
	}
	for _, state := range []remote_send.InputType{remote_send.KeyDown, remote_send.KeyUp} {
		pause, ok := makeKeyboardInput(0x13, state)
		if !ok || pause.Ki.WVk != 0x13 || pause.Ki.DwFlags&(win.KEYEVENTF_SCANCODE|win.KEYEVENTF_EXTENDEDKEY) != 0 {
			t.Fatalf("Pause: %+v", pause)
		}
		letter, ok := makeKeyboardInput(0x41, state)
		if !ok || letter.Ki.WScan == 0 || letter.Ki.DwFlags&win.KEYEVENTF_EXTENDEDKEY != 0 {
			t.Fatalf("A: %+v", letter)
		}
	}
}
