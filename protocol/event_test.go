package protocol

import "testing"

func TestParseEvent(t *testing.T) {
	for _, s := range []string{"1 0 -2147483648 2147483647", "1 1 0 65535", "0 1 0 0", "0 4 1 0", "3 65 1 0", "2 240 1 0", "2 4294967176 0 0"} {
		if _, err := ParseEvent(s); err != nil {
			t.Fatalf("rejected %q: %v", s, err)
		}
	}
	for _, s := range []string{"", "3 65 0", "3 65 0 0 extra", "4 0 0 0", "1 2 0 0", "1 1 -1 0", "1 1 65536 0", "1 0 2147483648 0", "3 256 0 0", "3 0 0 0", "3 65 2 0", "0 3 0 0", "0 1 0 1", "2 0 0 0", "2 32768 0 0", "2 4294934527 0 0", "3 -1 0 0"} {
		if _, err := ParseEvent(s); err == nil {
			t.Fatalf("accepted %q", s)
		}
	}
}

func FuzzParseEvent(f *testing.F) {
	f.Add("1 0 -10 20")
	f.Add("3 65 0 0")
	f.Fuzz(func(t *testing.T, line string) { ParseEvent(line) })
}
