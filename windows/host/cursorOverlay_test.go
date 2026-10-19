package host

import (
	"runtime"
	"testing"

	"github.com/lxn/win"
)

func TestCursorVisibilityModes(t *testing.T) {
	for _, tc := range []struct {
		mode  int
		flags uint32
		want  bool
	}{
		{0, 0, false}, {0, 1, true}, {0, 2, false},
		{1, 0, true}, {1, 1, true}, {1, 2, true},
		{2, 0, false}, {2, 1, false}, {2, 2, false},
	} {
		if got := cursorVisible(tc.mode, tc.flags); got != tc.want {
			t.Errorf("mode=%d flags=%d: got %v, want %v", tc.mode, tc.flags, got, tc.want)
		}
	}
}

// Create only a hidden window: exercise Win32 setup and resource teardown without
// showing an overlay, changing focus, or injecting input on the test machine.
func TestCursorOverlayLifecycle(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	restore := usePhysicalCursorCoordinates()
	defer restore()
	for i := 0; i < 2; i++ {
		o, err := newCursorOverlay()
		if err != nil {
			t.Fatal(err)
		}
		style := uint32(win.GetWindowLong(o.hwnd, win.GWL_EXSTYLE))
		want := uint32(win.WS_EX_LAYERED | win.WS_EX_TRANSPARENT | win.WS_EX_NOACTIVATE | win.WS_EX_TOOLWINDOW)
		if style&want != want {
			t.Errorf("missing overlay styles: %x", style)
		}
		if win.IsWindowVisible(o.hwnd) {
			t.Error("overlay became visible before first update")
		}
		o.mode = 2
		o.refresh()
		o.close()
		o.close()
	}
}
