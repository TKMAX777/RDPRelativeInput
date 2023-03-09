package client

import (
	"os"
	"path/filepath"
	"syscall"

	"github.com/TKMAX777/RDPRelativeInput/remote_send"
	"github.com/TKMAX777/RDPRelativeInput/winapi"

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

var StartHook *syscall.Proc
var Unhook *syscall.Proc

func init() {
	var dllpath = filepath.Join(os.Getenv("ProgramW6432"), "RDPRelativeInput", "globalhotkey.dll")

	globalhotkey, err := syscall.LoadDLL(dllpath)
	if err != nil {
		win.MessageBox(0, winapi.MustUTF16PtrFromString(err.Error()), winapi.MustUTF16PtrFromString("Load DLL Error"), 0)
	}

	StartHook = globalhotkey.MustFindProc("StartHook")
	Unhook = globalhotkey.MustFindProc("Unhook")
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
