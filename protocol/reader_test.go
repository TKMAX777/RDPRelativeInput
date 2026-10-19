package protocol

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"
)

func TestReadLines(t *testing.T) {
	for _, fragmented := range []bool{false, true} {
		var source io.Reader = strings.NewReader("RDPRelativeInput\nStatus:OK\r\n1 0 -10 20\nCLOSE\n")
		if fragmented {
			source = iotest.OneByteReader(source)
		}
		r := NewReader(source)
		for _, want := range []string{"RDPRelativeInput", "Status:OK", "1 0 -10 20", "CLOSE"} {
			got, err := r.ReadLine()
			if err != nil || got != want {
				t.Fatalf("got %q, %v; want %q", got, err, want)
			}
		}
		if _, err := r.ReadLine(); err != io.EOF {
			t.Fatalf("EOF: %v", err)
		}
	}
}

func TestReadLineRejectsIncompleteAndOversized(t *testing.T) {
	for _, s := range []string{strings.Repeat("x", MaxLineSize) + "\n", strings.Repeat("x", 10000), "CLOSE"} {
		if _, err := NewReader(strings.NewReader(s)).ReadLine(); err == nil {
			t.Fatal("accepted invalid line")
		}
	}
	if _, err := NewReader(strings.NewReader("CLOSE")).ReadLine(); !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatal(err)
	}
	line := strings.Repeat("x", MaxLineSize-1)
	if got, err := NewReader(strings.NewReader(line + "\n")).ReadLine(); err != nil || got != line {
		t.Fatalf("maximum valid line: %q, %v", got, err)
	}
	readError := errors.New("disconnected")
	if _, err := NewReader(iotest.ErrReader(readError)).ReadLine(); !errors.Is(err, readError) {
		t.Fatalf("lost read error: %v", err)
	}
}

func TestReadVersion(t *testing.T) {
	if err := ReadVersion(strings.NewReader("ProtocolVersion:2\n"), 3); err == nil {
		t.Fatal("v3 accepted a v2 peer")
	}
	if err := ReadVersion(strings.NewReader("ProtocolVersion:3\n"), 2); err == nil {
		t.Fatal("v2 accepted a v3 peer")
	}
	for _, s := range []string{"ProtocolVersion:12\n", "junk2\n", "", "ProtocolVersion:2", "ProtocolVersion:2 extra\n"} {
		if err := ReadVersion(strings.NewReader(s), 2); err == nil {
			t.Fatalf("accepted %q", s)
		}
	}
	if err := ReadVersion(strings.NewReader("ProtocolVersion:2\n"), 2); err != nil {
		t.Fatal(err)
	}
}
