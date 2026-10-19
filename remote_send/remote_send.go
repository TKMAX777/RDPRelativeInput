package remote_send

import (
	"io"
	"sync"

	"github.com/TKMAX777/RDPRelativeInput/keymap"
	"github.com/TKMAX777/RDPRelativeInput/protocol"
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

// Handler serializes frames and remembers the first transport failure. Continuing
// after a partial write would make the receiver interpret a corrupted stream.
type Handler struct {
	writer io.Writer
	mu     sync.Mutex
	frame  [protocol.MaxFrameSize]byte
	err    error
}

func New(w io.Writer) *Handler { return &Handler{writer: w} }

func (h *Handler) send(e protocol.Event, closeSession bool) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.err != nil {
		return h.err
	}
	n, err := protocol.EncodeEvent(&h.frame, e)
	if closeSession {
		n, err = protocol.EncodeClose(&h.frame), nil
	}
	if err != nil {
		return err
	}
	written, err := h.writer.Write(h.frame[:n])
	if err == nil && written != n {
		err = io.ErrShortWrite
	}
	h.err = err
	if closeSession && err == nil {
		h.err = io.ErrClosedPipe
	}
	return err
}

func (h *Handler) SendRelativeCursor(x, y int32) error {
	return h.send(protocol.Event{Type: keymap.EV_TYPE_MOUSE_MOVE, Value1: x, Value2: y}, false)
}

func (h *Handler) SendAbsoluteCursor(x, y int32) error {
	return h.send(protocol.Event{Type: keymap.EV_TYPE_MOUSE_MOVE, Input: 1, Value1: x, Value2: y}, false)
}

func (h *Handler) SendInput(eventType keymap.EV_TYPE, keyValue uint32, state InputType) error {
	return h.send(protocol.Event{Type: eventType, Input: keyValue, Value1: int32(state)}, false)
}

func (h *Handler) SendExit() error { return h.send(protocol.Event{}, true) }
