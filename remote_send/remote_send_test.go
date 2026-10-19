package remote_send

import (
	"bytes"
	"github.com/TKMAX777/RDPRelativeInput/protocol"
	"io"
	"sync"
	"testing"
)

type shortWriter struct{ calls int }

func (w *shortWriter) Write(p []byte) (int, error) { w.calls++; return len(p) - 1, nil }
func TestShortWriteStopsStream(t *testing.T) {
	w := new(shortWriter)
	h := New(w)
	if err := h.SendRelativeCursor(1, 2); err != io.ErrShortWrite {
		t.Fatal(err)
	}
	if err := h.SendRelativeCursor(3, 4); err != io.ErrShortWrite || w.calls != 1 {
		t.Fatalf("%v %d", err, w.calls)
	}
}
func TestConcurrentFrames(t *testing.T) {
	var out bytes.Buffer
	h := New(&out)
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				if err := h.SendRelativeCursor(1, -1); err != nil {
					t.Error(err)
				}
			}
		}()
	}
	wg.Wait()
	if err := h.SendExit(); err != nil {
		t.Fatal(err)
	}
	r := protocol.NewReader(&out)
	for i := 0; i < 800; i++ {
		e, err := r.ReadEvent()
		if err != nil || e.Value1 != 1 || e.Value2 != -1 {
			t.Fatalf("%+v %v", e, err)
		}
	}
	if _, err := r.ReadEvent(); err != protocol.ErrClosed {
		t.Fatal(err)
	}
	if err := h.SendRelativeCursor(1, 2); err != io.ErrClosedPipe {
		t.Fatal(err)
	}
}
