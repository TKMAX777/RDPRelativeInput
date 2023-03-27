package keymap

type EV_TYPE uint32

const (
	EV_TYPE_MOUSE EV_TYPE = iota
	EV_TYPE_MOUSE_MOVE
	EV_TYPE_WHEEL
	EV_TYPE_KEY
)
