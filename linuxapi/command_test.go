package linuxapi

import (
	"reflect"
	"testing"
)

func TestCommandRejectsInjection(t *testing.T) {
	x := NewXdotool(":0")
	for _, key := range []string{"a; touch /tmp/pwn", "$(id)", "a exec sh", "--clearmodifiers", "a\nexec", "a++b", ""} {
		if _, err := x.command("keydown", "--", key); err == nil {
			t.Fatalf("accepted %q", key)
		}
	}
	for _, coordinate := range []string{"1;id", "$(id)", "1 exec sh", "2147483648"} {
		if _, err := x.command("mousemove", "--", coordinate, "0"); err == nil {
			t.Fatalf("accepted %q", coordinate)
		}
	}
	if _, err := x.command("exec", "sh"); err == nil {
		t.Fatal("unlisted command accepted")
	}
}
func TestCommandUsesSeparatedArguments(t *testing.T) {
	cmd, err := NewXdotool(":0").command("keydown", "--", "ctrl+alt+Delete")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"/usr/bin/xdotool", "keydown", "--", "ctrl+alt+Delete"}
	if !reflect.DeepEqual(cmd.Args, want) {
		t.Fatalf("%q", cmd.Args)
	}
	if _, err := NewXdotool(":0").command("mousemove_relative", "--", "-10", "20"); err != nil {
		t.Fatal(err)
	}
}
