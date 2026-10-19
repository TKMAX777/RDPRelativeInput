package linuxapi

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// xdotool has command chaining of its own; validate arguments as well as
// avoiding the shell. The executable path never comes from remote input.
func (x XdotoolHandler) command(event string, args ...string) (*exec.Cmd, error) {
	invalid := fmt.Errorf("invalid xdotool arguments")
	switch event {
	case "mousemove", "mousemove_relative":
		if len(args) != 3 || args[0] != "--" {
			return nil, invalid
		}
		for _, v := range args[1:] {
			if _, err := strconv.ParseInt(v, 10, 32); err != nil {
				return nil, invalid
			}
		}
	case "keydown", "keyup":
		if len(args) != 2 || args[0] != "--" || !validKey(args[1]) {
			return nil, invalid
		}
	case "mousedown", "mouseup":
		if len(args) != 2 || args[0] != "--" {
			return nil, invalid
		}
		n, err := strconv.Atoi(args[1])
		if err != nil || n < 1 || n > 5 {
			return nil, invalid
		}
	case "click":
		if len(args) != 1 || (args[0] != "4" && args[0] != "5") {
			return nil, invalid
		}
	case "getmouselocation":
		if len(args) != 1 || args[0] != "--shell" {
			return nil, invalid
		}
	case "getdisplaygeometry":
		if len(args) != 0 {
			return nil, invalid
		}
	default:
		return nil, invalid
	}
	return exec.Command("/usr/bin/xdotool", append([]string{event}, args...)...), nil
}
func validKey(key string) bool {
	if len(key) == 0 || len(key) > 128 {
		return false
	}
	for _, part := range strings.Split(key, "+") {
		if part == "" {
			return false
		}
		for _, c := range part {
			if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '_') {
				return false
			}
		}
	}
	return true
}
