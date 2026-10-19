package protocol

import (
	"bytes"
	"fmt"
	"github.com/TKMAX777/RDPRelativeInput/keymap"
	"io"
	"testing"
	"testing/iotest"
)

func TestBinaryStream(t *testing.T) {
	events := []Event{
		{keymap.EV_TYPE_MOUSE_MOVE, 0, -32768, 32767},
		{keymap.EV_TYPE_MOUSE_MOVE, 0, -2147483648, 2147483647},
		{keymap.EV_TYPE_MOUSE_MOVE, 1, 65535, 0},
		{keymap.EV_TYPE_KEY, 65, 0, 0}, {keymap.EV_TYPE_KEY, 254, 1, 0},
		{keymap.EV_TYPE_MOUSE, 1, 0, 0}, {keymap.EV_TYPE_MOUSE, 4, 1, 0},
		{keymap.EV_TYPE_WHEEL, 4294967176, 0, 0}, {keymap.EV_TYPE_WHEEL, 240, 1, 0},
	}
	var stream bytes.Buffer
	stream.WriteString("RDPRelativeInput\nStatus:OK\n")
	var frame [MaxFrameSize]byte
	for _, e := range events {
		n, err := EncodeEvent(&frame, e)
		if err != nil {
			t.Fatal(err)
		}
		stream.Write(frame[:n])
	}
	stream.Write(frame[:EncodeClose(&frame)])
	for _, fragmented := range []bool{false, true} {
		var source io.Reader = bytes.NewReader(stream.Bytes())
		if fragmented {
			source = iotest.OneByteReader(source)
		}
		r := NewReader(source)
		for _, want := range []string{"RDPRelativeInput", "Status:OK"} {
			got, err := r.ReadLine()
			if err != nil || got != want {
				t.Fatalf("handshake %q %v", got, err)
			}
		}
		for _, want := range events {
			got, err := r.ReadEvent()
			if err != nil || got != want {
				t.Fatalf("got %+v %v want %+v", got, err, want)
			}
		}
		if _, err := r.ReadEvent(); err != ErrClosed {
			t.Fatal(err)
		}
	}
}

func TestBinaryInvalid(t *testing.T) {
	for _, data := range [][]byte{{0xff}, {frameKeyDown, 0}, {frameMouseDown, 3}, {frameWheel, 0, 0}} {
		if _, err := NewReader(bytes.NewReader(data)).ReadEvent(); err != ErrInvalidEvent {
			t.Fatalf("%x: %v", data, err)
		}
	}
	for _, op := range []byte{frameRelative16, frameRelative32, frameAbsolute, frameKeyDown, frameKeyUp, frameMouseDown, frameMouseUp, frameWheel} {
		size := 2
		switch op {
		case frameRelative16, frameAbsolute:
			size = 5
		case frameRelative32:
			size = 9
		case frameWheel:
			size = 3
		}
		for n := 1; n < size; n++ {
			data := make([]byte, n)
			data[0] = op
			if _, err := NewReader(bytes.NewReader(data)).ReadEvent(); err != io.ErrUnexpectedEOF {
				t.Fatalf("truncation %x: %v", data, err)
			}
		}
	}
	var frame [MaxFrameSize]byte
	if _, err := EncodeEvent(&frame, Event{Type: keymap.EV_TYPE_KEY, Input: 256}); err == nil {
		t.Fatal("invalid key accepted")
	}
}

func TestBinaryWireBytes(t *testing.T) {
	var frame [MaxFrameSize]byte
	n, err := EncodeEvent(&frame, Event{keymap.EV_TYPE_MOUSE_MOVE, 0, -2, 256})
	if err != nil || !bytes.Equal(frame[:n], []byte{0x80, 0xfe, 0xff, 0, 1}) {
		t.Fatalf("wire: %x %v", frame[:n], err)
	}
}

func FuzzBinaryReader(f *testing.F) {
	f.Add([]byte{0x80, 1, 0, 2, 0})
	f.Fuzz(func(t *testing.T, data []byte) {
		r := NewReader(bytes.NewReader(data))
		for {
			if _, err := r.ReadEvent(); err != nil {
				return
			}
		}
	})
}

func BenchmarkEncodeText(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmt.Fprintf(io.Discard, "%d %d %d %d\n", keymap.EV_TYPE_MOUSE_MOVE, 0, int32(-12), int32(34))
	}
}
func BenchmarkEncodeBinary(b *testing.B) {
	var frame [MaxFrameSize]byte
	e := Event{keymap.EV_TYPE_MOUSE_MOVE, 0, -12, 34}
	for i := 0; i < b.N; i++ {
		n, _ := EncodeEvent(&frame, e)
		io.Discard.Write(frame[:n])
	}
}
func BenchmarkDecodeText(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseEvent("1 0 -12 34")
	}
}
func BenchmarkDecodeBinary(b *testing.B) {
	var frame [MaxFrameSize]byte
	n, _ := EncodeEvent(&frame, Event{keymap.EV_TYPE_MOUSE_MOVE, 0, -12, 34})
	source := bytes.NewReader(frame[:n])
	r := NewReader(source)
	for i := 0; i < b.N; i++ {
		source.Reset(frame[:n])
		r.ReadEvent()
	}
}
