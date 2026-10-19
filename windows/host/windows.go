package host

import "github.com/lxn/win"

type Handler struct {
	overlay *cursorOverlay
	button  win.HWND
}

func New() *Handler {
	return &Handler{}
}
