package client

import (
	"github.com/TKMAX777/RDPRelativeInput/remote_send"

	"github.com/lxn/win"
)

type Handler struct {
	metrics SystemMetrics
	options option
	remote  *remote_send.Handler
}

type ToggleType int

const (
	ToggleTypeOnce ToggleType = iota + 1
	ToggleTypeAlive
)

type option struct {
	toggleKey  string
	toggleType ToggleType
}

type SystemMetrics struct {
	FrameWidthX int32
	FrameWidthY int32
	TitleHeight int32
}

func New(r *remote_send.Handler) *Handler {
	return &Handler{
		remote: r,
		options: option{
			// set default options
			toggleKey:  "F8",
			toggleType: ToggleTypeAlive,
		},
		metrics: SystemMetrics{
			FrameWidthX: win.GetSystemMetrics(win.SM_CXSIZEFRAME),
			FrameWidthY: win.GetSystemMetrics(win.SM_CYSIZEFRAME),
			TitleHeight: win.GetSystemMetrics(win.SM_CYCAPTION),
		},
	}
}
