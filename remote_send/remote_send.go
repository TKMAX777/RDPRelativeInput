package remote_send

import (
	"fmt"
	"io"

	"github.com/TKMAX777/RDPRelativeInput/keymap"
)

type MouseMoveType uint32

const (
	MouseMoveTypeRelative MouseMoveType = iota
	MouseMoveTypeAbsolute
)

type InputType uint32

const (
	KeyDown InputType = iota
	KeyUp
)

type Handler struct {
	writer io.Writer
}

func New(w io.Writer) *Handler {
	return &Handler{w}
}

func (h Handler) SendRelativeCursor(x, y int32) {
	fmt.Fprintf(h.writer, "%d %d %d %d\n", keymap.EV_TYPE_MOUSE_MOVE, MouseMoveTypeRelative, x, y)
}

func (h Handler) SendAbsoluteCursor(x, y int32) {
	fmt.Fprintf(h.writer, "%d %d %d %d\n", keymap.EV_TYPE_MOUSE_MOVE, MouseMoveTypeAbsolute, x, y)
}

func (h Handler) SendInput(eventType keymap.EV_TYPE, keyValue uint32, state InputType) {
	fmt.Fprintf(h.writer, "%d %d %d %d\n", eventType, keyValue, state, 0)
}
