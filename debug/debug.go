package debug

import (
	"fmt"
	"os"
)

var DEBUG bool
var DEBUG_PATH string

func init() {
	DEBUG = os.Getenv("RELATIVE_INPUT_DEBUG") == "ON"
	DEBUG_PATH = os.Getenv("RELATIVE_INPUT_DEBUG_PATH")
}

func Debugf(format string, aug ...interface{}) (n int, err error) {
	if DEBUG {
		f, _ := os.OpenFile(DEBUG_PATH, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		defer f.Close()
		return fmt.Fprintf(f, format, aug...)
	}
	return
}

func Debugln(aug ...interface{}) (n int, err error) {
	if DEBUG {
		f, _ := os.OpenFile(DEBUG_PATH, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		defer f.Close()
		return fmt.Fprintln(f, aug...)
	}
	return
}
