package winapi

import "testing"

func TestChannelBufferValidation(t *testing.T) {
	for _, n := range []int{-1, 2} {
		var count uint32 = 42
		if err := WTSVirtualChannelWrite(0, make([]byte, 1), n, &count); err == nil || count != 0 {
			t.Fatalf("write: %v, %d", err, count)
		}
		count = 42
		if err := WTSVirtualChannelRead(0, 0, make([]byte, 1), n, &count); err == nil || count != 0 {
			t.Fatalf("read: %v, %d", err, count)
		}
	}
	var count uint32 = 42
	if err := WTSVirtualChannelWrite(0, nil, 0, &count); err != nil || count != 0 {
		t.Fatalf("empty write: %v, %d", err, count)
	}
	if err := WTSVirtualChannelRead(0, 0, nil, 0, &count); err != nil || count != 0 {
		t.Fatalf("empty read: %v, %d", err, count)
	}
	if err := WTSVirtualChannelRead(0, 0, nil, 0, nil); err == nil {
		t.Fatal("accepted nil count")
	}
	if err := WTSVirtualChannelWrite(0, nil, 0, nil); err == nil {
		t.Fatal("accepted nil count")
	}
}

func TestChannelNameValidation(t *testing.T) {
	for _, name := range []string{"", "RELINP\x00other"} {
		if _, err := WTSVirtualChannelOpenEx(0, name, 0); err == nil {
			t.Fatalf("accepted %q", name)
		}
	}
}
