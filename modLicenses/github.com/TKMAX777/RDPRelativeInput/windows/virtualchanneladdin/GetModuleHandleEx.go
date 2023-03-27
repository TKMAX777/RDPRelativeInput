package main

import "C"

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var _ unsafe.Pointer

var (
	modKernel32           = windows.NewLazySystemDLL("Kernel32.dll")
	procGetModuleHandleEx = modKernel32.NewProc("GetModuleHandleExW")
)

func GetModuleHandleExW(dwFlags uint32, lpModuleName string, phModule *HANDLE) error {
	str, _ := syscall.UTF16PtrFromString(lpModuleName)
	r0, _, err := syscall.SyscallN(procGetModuleHandleEx.Addr(), uintptr(dwFlags), uintptr(unsafe.Pointer(str)), uintptr(unsafe.Pointer(phModule)))
	if r0 != 0 {
		return err
	}
	return nil
}
