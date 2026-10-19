package client

import (
	"fmt"
	"github.com/TKMAX777/RDPRelativeInput/winapi"
	"github.com/lxn/win"
)

func findRDPWindow(pid uint32) (win.HWND, error) {
	var candidates []win.HWND
	var hwnd win.HWND
	for {
		hwnd = winapi.FindWindowEx(0, hwnd, winapi.MustUTF16PtrFromString("TscShellContainerClass"), nil)
		if hwnd == 0 {
			break
		}
		var owner uint32
		if win.GetWindowThreadProcessId(hwnd, &owner) != 0 && owner == pid {
			candidates = append(candidates, hwnd)
		}
	}
	return uniqueRDPWindow(candidates)
}
func uniqueRDPWindow(candidates []win.HWND) (win.HWND, error) {
	if len(candidates) != 1 {
		return 0, fmt.Errorf("Cannot safely identify the RDP window (%d candidates). Use one connection per mstsc process.", len(candidates))
	}
	return candidates[0], nil
}
