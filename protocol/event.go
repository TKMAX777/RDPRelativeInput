package protocol

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/TKMAX777/RDPRelativeInput/keymap"
)

type Event struct {
	Type           keymap.EV_TYPE
	Input          uint32
	Value1, Value2 int32
}

func ParseEvent(line string) (Event, error) {
	var e Event
	invalid := fmt.Errorf("invalid input event")
	if len(line) > MaxLineSize {
		return e, invalid
	}
	fields := strings.Fields(line)
	if len(fields) != 4 {
		return e, invalid
	}
	t, err := strconv.ParseUint(fields[0], 10, 32)
	if err != nil {
		return e, invalid
	}
	k, err := strconv.ParseUint(fields[1], 10, 32)
	if err != nil {
		return e, invalid
	}
	x, err := strconv.ParseInt(fields[2], 10, 32)
	if err != nil {
		return e, invalid
	}
	y, err := strconv.ParseInt(fields[3], 10, 32)
	if err != nil {
		return e, invalid
	}
	e = Event{keymap.EV_TYPE(t), uint32(k), int32(x), int32(y)}
	return e, e.Validate()
}

func (e Event) Validate() error {
	invalid := ErrInvalidEvent
	k, x, y := e.Input, e.Value1, e.Value2
	if e.Type == keymap.EV_TYPE_MOUSE_MOVE {
		if k > 1 || (k == 1 && (x < 0 || x > 65535 || y < 0 || y > 65535)) {
			return invalid
		}
		return nil
	}
	if (x != 0 && x != 1) || y != 0 {
		return invalid
	}
	switch e.Type {
	case keymap.EV_TYPE_MOUSE:
		if k != 1 && k != 2 && k != 4 {
			return invalid
		}
	case keymap.EV_TYPE_KEY:
		if k == 0 || k > 254 {
			return invalid
		}
	case keymap.EV_TYPE_WHEEL:
		// Wheel deltas are signed 16-bit values transmitted as uint32.
		delta := int32(k)
		if delta == 0 || delta < -32768 || delta > 32767 {
			return invalid
		}
	default:
		return invalid
	}
	return nil
}
