package channelbuffer

import (
	"bytes"
	"io"
	"testing"
	"time"
)

func TestWrapAndOrder(t *testing.T) {
	b := New()
	first := bytes.Repeat([]byte{1}, Capacity-10)
	b.Write(first)
	out := make([]byte, Capacity-20)
	if _, err := io.ReadFull(b, out); err != nil {
		t.Fatal(err)
	}
	second := bytes.Repeat([]byte{2}, 100)
	b.Write(second)
	out = make([]byte, 110)
	if _, err := io.ReadFull(b, out); err != nil {
		t.Fatal(err)
	}
	want := append(bytes.Repeat([]byte{1}, 10), second...)
	if !bytes.Equal(out, want) {
		t.Fatal("ring buffer reordered data")
	}
}
func TestConcurrentReadWrite(t *testing.T) {
	b := New()
	want := bytes.Repeat([]byte("input"), 5000)
	done := make(chan error, 1)
	go func() {
		for _, v := range want {
			if _, err := b.Write([]byte{v}); err != nil {
				done <- err
				return
			}
		}
		done <- nil
	}()
	got := make([]byte, len(want))
	if _, err := io.ReadFull(b, got); err != nil {
		t.Fatal(err)
	}
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatal("stream corrupted")
	}
}
func TestOverflowCancels(t *testing.T) {
	b := New()
	if _, err := b.Write(make([]byte, Capacity)); err != nil {
		t.Fatal(err)
	}
	if _, err := b.Write([]byte{1}); err != ErrOverflow {
		t.Fatal(err)
	}
	if _, err := b.Read(make([]byte, 1)); err != ErrOverflow {
		t.Fatal(err)
	}
	select {
	case <-b.Context().Done():
	default:
		t.Fatal("overflow did not cancel session")
	}
}
func TestCloseUnblocksRead(t *testing.T) {
	b := New()
	done := make(chan error, 1)
	go func() { _, err := b.Read(make([]byte, 1)); done <- err }()
	b.Fail(io.EOF)
	b.Fail(ErrOverflow)
	select {
	case err := <-done:
		if err != io.EOF {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("reader stuck")
	}
	if _, err := b.Write([]byte{1}); err != io.EOF {
		t.Fatal(err)
	}
}
