package client

import (
	"github.com/lxn/win"
	"testing"
)

func TestUniqueRDPWindow(t *testing.T) {
	for _, c := range [][]win.HWND{nil, {1, 2}} {
		if h, err := uniqueRDPWindow(c); err == nil || h != 0 {
			t.Fatal("ambiguous window accepted")
		}
	}
	if h, err := uniqueRDPWindow([]win.HWND{42}); err != nil || h != 42 {
		t.Fatalf("%d %v", h, err)
	}
}
