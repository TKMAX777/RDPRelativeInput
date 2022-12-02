package linuxapi

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	XdotoolMouseClickLeft = 1 + iota
	XdotoolMouseClickMiddle
	XdotoolMouseClickRight
	XdotoolMouseClickWheelUp
	XdotoolMouseClickWheelDown
)

type XdotoolHandler struct {
	display string
}

func NewXdotool(display string) *XdotoolHandler {
	return &XdotoolHandler{display: display}
}

func (x XdotoolHandler) exec(event string, value ...string) string {
	var cmd = exec.Command("/bin/bash", "-c",
		fmt.Sprintf("xdotool %s %s", event, strings.Join(value, " ")),
	)
	cmd.Env = append(os.Environ(), "DISPLAY="+x.display)

	var eout = new(bytes.Buffer)
	var sout = new(bytes.Buffer)

	cmd.Stderr = eout
	cmd.Stdout = sout

	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to input: %s", err.Error())
	}

	return sout.String()
}

func (x XdotoolHandler) MouseMoveRelative(X, Y string) {
	go x.exec("mousemove_relative", "--", X, Y)
}

func (x XdotoolHandler) MouseMoveAbsolute(X, Y string) {
	go x.exec("mousemove", "--", X, Y)
}

func (x XdotoolHandler) MouseDown(Button int) {
	go x.exec("mousedown", "--", strconv.Itoa(Button))
}

func (x XdotoolHandler) MouseUp(Button int) {
	go x.exec("mouseup", "--", strconv.Itoa(Button))
}

func (x XdotoolHandler) WheelUp() {
	go x.exec("click", "4")
}

func (x XdotoolHandler) WheelDown() {
	go x.exec("click", "5")
}

func (x XdotoolHandler) GetPosition() (X, Y int) {
	var out = x.exec("getmouselocation", "--shell")
	var outs = strings.Split(out, "\n")

	if len(outs) < 2 || !strings.Contains(outs[0], "=") {
		return
	}
	X, _ = strconv.Atoi(strings.Split(outs[0], "=")[1])
	Y, _ = strconv.Atoi(strings.Split(outs[1], "=")[1])

	return X, Y
}

func (x XdotoolHandler) GetWindowGeometry() (width, height int) {
	var out = x.exec("getdisplaygeometry")
	var outs = strings.Split(out, " ")

	width, _ = strconv.Atoi(strings.Split(outs[0], "=")[1])
	height, _ = strconv.Atoi(strings.Split(outs[1], "=")[1])

	return width, height
}

func (x XdotoolHandler) KeyDown(key string) {
	x.exec("keydown", key)
}

func (x XdotoolHandler) KeyUp(key string) {
	x.exec("keyup", key)
}
