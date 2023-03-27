package winapi

import (
	"syscall"
	"unsafe"

	"github.com/lxn/win"
)

func MustUTF16PtrFromString(str string) *uint16 {
	ptr, err := syscall.UTF16PtrFromString(str)
	if err != nil {
		panic(err)
	}
	return ptr
}

func MustUTF16FromString(str string) []uint16 {
	ptr, err := syscall.UTF16FromString(str)
	if err != nil {
		panic(err)
	}
	return ptr
}

func UTF16PtrToString(p *uint16) string {
	if p == nil {
		return ""
	}

	var char uint16
	var chars = []uint16{}

	for i := 0; ; i++ {
		char = *(*uint16)(unsafe.Pointer(unsafe.Add(unsafe.Pointer(p), unsafe.Sizeof(uint16(0))*uintptr(i))))
		chars = append(chars, char)
		// null char
		if char == 0 {
			break
		}
	}
	return syscall.UTF16ToString(chars)
}

func UTF8PtrToString(p *byte) string {
	if p == nil {
		return ""
	}

	var char byte
	var chars = []byte{}

	for i := 0; ; i++ {
		char = *(*byte)(unsafe.Pointer(unsafe.Add(unsafe.Pointer(p), unsafe.Sizeof(byte(0))*uintptr(i))))
		// null char
		if char == 0 {
			break
		}
		chars = append(chars, char)
	}

	return string(chars)
}

func EnumChildWindows(hwnd win.HWND) []win.HWND {
	var res = []win.HWND{}

	var chwnd win.HWND
	for chwnd = FindWindowEx(hwnd, chwnd, nil, nil); chwnd != win.HWND(NULL); chwnd = FindWindowEx(hwnd, chwnd, nil, nil) {
		res = append(res, chwnd)
	}

	return res
}

func FindChildWindowsFromWindowText(parentHWND win.HWND, lpszClass *uint16, lpszWindow *uint16, windowText string) win.HWND {
	var chwnd win.HWND
	for chwnd = FindWindowEx(parentHWND, chwnd, nil, nil); chwnd != win.HWND(NULL); chwnd = FindWindowEx(parentHWND, chwnd, nil, nil) {
		var UTF16name = make([]uint16, 1000)
		GetWindowText(chwnd, UTF16name, 1000)
		if syscall.UTF16ToString(UTF16name) == windowText {
			return chwnd
		}
	}
	return win.HWND(NULL)
}

func SendMessage(hwnd win.HWND, msg uint32, wParam uintptr, lParam uintptr) uintptr {
	return win.SendMessage(hwnd, msg, wParam, lParam)
}

func GetWindowTextString(hwnd win.HWND) string {
	var name = make([]uint16, 1000)
	GetWindowText(hwnd, name, 1000)
	return syscall.UTF16ToString(name)
}
