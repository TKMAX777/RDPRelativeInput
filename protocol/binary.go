package protocol

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/TKMAX777/RDPRelativeInput/keymap"
)

// Opcodes determine payload size; remote peers cannot request allocations.
const (
	frameRelative16 byte = 0x80 + iota
	frameRelative32
	frameAbsolute
	frameKeyDown
	frameKeyUp
	frameMouseDown
	frameMouseUp
	frameWheel
	frameClose
)

const MaxFrameSize = 9

var ErrInvalidEvent = errors.New("invalid input event")
var ErrClosed = errors.New("peer closed input session")

func EncodeEvent(dst *[MaxFrameSize]byte, e Event) (int, error) {
	if err := e.Validate(); err != nil {
		return 0, err
	}
	switch e.Type {
	case keymap.EV_TYPE_MOUSE_MOVE:
		if e.Input == 1 {
			dst[0] = frameAbsolute
		} else if e.Value1 >= -32768 && e.Value1 <= 32767 && e.Value2 >= -32768 && e.Value2 <= 32767 {
			dst[0] = frameRelative16
		} else {
			dst[0] = frameRelative32
			binary.LittleEndian.PutUint32(dst[1:5], uint32(e.Value1))
			binary.LittleEndian.PutUint32(dst[5:9], uint32(e.Value2))
			return 9, nil
		}
		binary.LittleEndian.PutUint16(dst[1:3], uint16(e.Value1))
		binary.LittleEndian.PutUint16(dst[3:5], uint16(e.Value2))
		return 5, nil
	case keymap.EV_TYPE_KEY:
		dst[0], dst[1] = frameKeyDown+byte(e.Value1), byte(e.Input)
		return 2, nil
	case keymap.EV_TYPE_MOUSE:
		dst[0], dst[1] = frameMouseDown+byte(e.Value1), byte(e.Input)
		return 2, nil
	case keymap.EV_TYPE_WHEEL:
		dst[0] = frameWheel
		binary.LittleEndian.PutUint16(dst[1:3], uint16(e.Input))
		return 3, nil
	}
	return 0, ErrInvalidEvent
}

// Keep the handshake's reader so coalesced input bytes are not lost.
func (r *Reader) ReadEvent() (Event, error) {
	var e Event
	op, err := r.reader.ReadByte()
	if err != nil {
		return e, err
	}
	n := 0
	switch op {
	case frameRelative16, frameAbsolute:
		n = 4
	case frameRelative32:
		n = 8
	case frameKeyDown, frameKeyUp, frameMouseDown, frameMouseUp:
		n = 1
	case frameWheel:
		n = 2
	case frameClose:
		return e, ErrClosed
	default:
		return e, ErrInvalidEvent
	}
	b := r.frame[:n]
	if _, err := io.ReadFull(r.reader, b); err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return e, err
	}
	switch op {
	case frameRelative16:
		e = Event{keymap.EV_TYPE_MOUSE_MOVE, 0, int32(int16(binary.LittleEndian.Uint16(b))), int32(int16(binary.LittleEndian.Uint16(b[2:])))}
	case frameRelative32:
		e = Event{keymap.EV_TYPE_MOUSE_MOVE, 0, int32(binary.LittleEndian.Uint32(b)), int32(binary.LittleEndian.Uint32(b[4:]))}
	case frameAbsolute:
		e = Event{keymap.EV_TYPE_MOUSE_MOVE, 1, int32(binary.LittleEndian.Uint16(b)), int32(binary.LittleEndian.Uint16(b[2:]))}
	case frameKeyDown, frameKeyUp:
		e = Event{keymap.EV_TYPE_KEY, uint32(b[0]), int32(op - frameKeyDown), 0}
	case frameMouseDown, frameMouseUp:
		e = Event{keymap.EV_TYPE_MOUSE, uint32(b[0]), int32(op - frameMouseDown), 0}
	case frameWheel:
		delta := int32(int16(binary.LittleEndian.Uint16(b)))
		state := int32(0)
		if delta > 0 {
			state = 1
		}
		e = Event{keymap.EV_TYPE_WHEEL, uint32(delta), state, 0}
	}
	return e, e.Validate()
}

func EncodeClose(dst *[MaxFrameSize]byte) int { dst[0] = frameClose; return 1 }
