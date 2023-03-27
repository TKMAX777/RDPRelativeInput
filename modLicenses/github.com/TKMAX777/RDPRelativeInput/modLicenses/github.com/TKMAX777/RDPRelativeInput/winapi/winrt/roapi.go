package winrt

import (
	"github.com/go-ole/go-ole"
	"github.com/lxn/win"
	"golang.org/x/sys/windows"
)

var modcombase = windows.NewLazySystemDLL("combase.dll")

// RO_INIT_TYPE enumeration (roapi.h)
// https://learn.microsoft.com/en-us/windows/win32/api/roapi/ne-roapi-ro_init_type

type RO_INIT_TYPE uint32

const (
	RO_INIT_SINGLETHREADED RO_INIT_TYPE = 0
	RO_INIT_MULTITHREADED  RO_INIT_TYPE = 1
)

var pRoInitialize = modcombase.NewProc("RoInitialize")

func RoInitialize(thread_type RO_INIT_TYPE) error {
	r1, _, _ := pRoInitialize.Call(uintptr(thread_type))
	if r1 != win.S_OK {
		return ole.NewError(r1)
	}
	return nil
}

var pRoUninitialize = modcombase.NewProc("RoUninitialize")

func RoUninitialize() {
	pRoUninitialize.Call()
}
