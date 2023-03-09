package main

import "C"

import (
	"os"

	"github.com/TKMAX777/RDPRelativeInput/debug"
	"github.com/TKMAX777/RDPRelativeInput/winapi"
	"github.com/lxn/win"
)

// Prevent for being collected by GC
var hWnd win.HWND

func hookFunc(code int, wParam uintptr, lParam uintptr) uintptr {
	if code < 0 || code == 3 {
		return winapi.CallNextHookEx(hHk, code, wParam, lParam)
	}

	switch {
	case lParam&(1<<31) == 0:
		win.PostMessage(hWnd, win.WM_KEYDOWN, wParam, 1<<31)
		// case lParam&(1<<31) == 1:
		// 	win.PostMessage(hWnd, win.WM_KEYUP, wParam, 0)
	}

	return winapi.CallNextHookEx(hHk, code, wParam, lParam)
}

//export StartHook
func StartHook(hwnd uintptr, handle uintptr) {
	win.MessageBox(0, winapi.MustUTF16PtrFromString("Start Hook"), winapi.MustUTF16PtrFromString("Toggle Hook Error"), win.MB_ICONERROR)

	debug.Debugln("Start Hook")
	hWnd = win.HWND(hwnd)

	var err error
	hHk, err = winapi.SetWindowHookEx(winapi.WH_KEYBOARD, winapi.HOOKPROC(hookFunc), win.HINSTANCE(handle), 0)
	if err != nil {
		win.MessageBox(0, winapi.MustUTF16PtrFromString(err.Error()), winapi.MustUTF16PtrFromString("Toggle Hook Error"), win.MB_ICONERROR)
		os.Exit(1)
	}
	debug.Debugln("SetWindowHookEx...done")
}
