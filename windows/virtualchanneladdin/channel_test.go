package main

import (
	"github.com/TKMAX777/RDPRelativeInput/internal/channelbuffer"
	"io"
	"testing"
)

func TestCallbackAndReaderConcurrently(t *testing.T) {
	rw := &VirtualChannelReadWriteCloser{received: channelbuffer.New()}
	callback := rw.getVirtualChannelOpenEventFn()
	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := 0; i < 1000; i++ {
			b := []byte{byte(i)}
			callback(0, CHANNEL_EVENT_DATA_RECEIVED, &b[0], 1, 1, CHANNEL_FLAG_ONLY)
		}
	}()
	data := make([]byte, 1000)
	if _, err := io.ReadFull(rw, data); err != nil {
		t.Fatal(err)
	}
	<-done
	for i, v := range data {
		if v != byte(i) {
			t.Fatalf("corruption at %d", i)
		}
	}
}

func TestReceiveRejectsOversizedNativePacket(t *testing.T) {
	rw := &VirtualChannelReadWriteCloser{received: channelbuffer.New()}
	// The pointer must never be dereferenced when size validation fails.
	rw.getVirtualChannelOpenEventFn()(0, CHANNEL_EVENT_DATA_RECEIVED, nil, channelbuffer.Capacity+1, channelbuffer.Capacity+1, CHANNEL_FLAG_ONLY)
	if _, err := rw.Read(make([]byte, 1)); err != channelbuffer.ErrOverflow {
		t.Fatal(err)
	}
}
func TestWriteCompletionDoesNotBlockOrLeak(t *testing.T) {
	w := &pendingWrite{data: []byte("input"), done: make(chan bool, 1)}
	rw := &VirtualChannelReadWriteCloser{pending: map[uintptr]*pendingWrite{1: w}}
	rw.complete(1, true)
	rw.complete(1, true)
	rw.complete(999, false)
	if len(rw.pending) != 0 {
		t.Fatal("pending write leaked")
	}
	select {
	case ok := <-w.done:
		if !ok {
			t.Fatal("completion failed")
		}
	default:
		t.Fatal("missing completion")
	}
}
func TestClosedChannelWrite(t *testing.T) {
	rw := &VirtualChannelReadWriteCloser{received: channelbuffer.New()}
	rw.received.Fail(io.EOF)
	if _, err := rw.Write([]byte("input")); err != io.ErrClosedPipe {
		t.Fatal(err)
	}
}
