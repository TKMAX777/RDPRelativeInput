package winapi

import (
	"syscall"
	"unsafe"

	"github.com/lxn/win"
)

const (
	LWA_COLORKEY uint32 = 1 + iota
	LWA_ALPHA
)

const (
	MAPVK_VK_TO_VSC uint32 = iota
	MAPVK_VSC_TO_VK
	MAPVK_VK_TO_CHAR
	MAPVK_VSC_TO_VK_EX
	MAPVK_VK_TO_VSC_EX
)

type CWPRETSTRUCT struct {
	LResult uintptr
	LParam  uintptr
	WParam  uintptr
	Message uint32
	Hwnd    win.HWND
}

type TYPE_HOOK_ID int

const (
	// Installs a hook procedure that monitors messages before the system sends them to the destination window procedure. For more information, see the CallWndProc hook procedure.
	WH_CALLWNDPROC TYPE_HOOK_ID = 4

	// Installs a hook procedure that monitors messages after they have been processed by the destination window procedure. For more information, see the CallWndRetProc hook procedure.
	WH_CALLWNDPROCRET TYPE_HOOK_ID = 12

	// Installs a hook procedure that receives notifications useful to a CBT application. For more information, see the CBTProc hook procedure.
	WH_CBT TYPE_HOOK_ID = 5

	// Installs a hook procedure useful for debugging other hook procedures. For more information, see the DebugProc hook procedure.
	WH_DEBUG TYPE_HOOK_ID = 9

	// Installs a hook procedure that will be called when the application's foreground thread is about to become idle. This hook is useful for performing low priority tasks during idle time. For more information, see the ForegroundIdleProc hook procedure.
	WH_FOREGROUNDIDLE TYPE_HOOK_ID = 11

	// Installs a hook procedure that monitors messages posted to a message queue. For more information, see the GetMsgProc hook procedure.
	WH_GETMESSAGE TYPE_HOOK_ID = 3

	//  Warning
	// Journaling Hooks APIs are unsupported starting in Windows 11 and will be removed in a future release. Because of this, we highly recommend calling the SendInput TextInput API instead.
	// Installs a hook procedure that posts messages previously recorded by a WH_JOURNALRECORD hook procedure. For more information, see the JournalPlaybackProc hook procedure.
	WH_JOURNALPLAYBACK TYPE_HOOK_ID = 1

	//  Warning
	// Journaling Hooks APIs are unsupported starting in Windows 11 and will be removed in a future release. Because of this, we highly recommend calling the SendInput TextInput API instead.
	// Installs a hook procedure that records input messages posted to the system message queue. This hook is useful for recording macros. For more information, see the JournalRecordProc hook procedure.
	WH_JOURNALRECORD TYPE_HOOK_ID = 0

	// Installs a hook procedure that monitors keystroke messages. For more information, see the KeyboardProc hook procedure.
	WH_KEYBOARD TYPE_HOOK_ID = 2

	// Installs a hook procedure that monitors low-level keyboard input events. For more information, see the LowLevelKeyboardProc hook procedure.
	WH_KEYBOARD_LL TYPE_HOOK_ID = 13

	// Installs a hook procedure that monitors mouse messages. For more information, see the MouseProc hook procedure.
	WH_MOUSE TYPE_HOOK_ID = 7

	// Installs a hook procedure that monitors low-level mouse input events. For more information, see the LowLevelMouseProc hook procedure.
	WH_MOUSE_LL TYPE_HOOK_ID = 14

	WH_MSGFILTER TYPE_HOOK_ID = -1

	// Installs a hook procedure that monitors messages generated as a result of an input event in a dialog box, message box, menu, or scroll bar. For more information, see the MessageProc hook procedure.
	// Installs a hook procedure that receives notifications useful to shell applications. For more information, see the ShellProc hook procedure.
	WH_SHELL TYPE_HOOK_ID = 10

	// Installs a hook procedure that monitors messages generated as a result of an input event in a dialog box, message box, menu, or scroll bar. The hook procedure monitors these messages for all applications in the same desktop as the calling thread. For more information, see the SysMsgProc hook procedure.
	WH_SYSMSGFILTER TYPE_HOOK_ID = 6
)

type KBDLLHOOKSTRUCT struct {
	VKCode      uint32
	ScanCode    uint32
	Flags       uint32
	Time        uint32
	DwExtraInfo uintptr
}

func ClipCursor(rect *win.RECT) (ok int, err error) {
	if rect == nil {
		return clipCursor(NULL)
	}
	return clipCursor(uintptr(unsafe.Pointer(&rect.Left)))
}

func EnumDesktopWindows(hDesktop win.HANDLE, lpEnumFunc uintptr, lParam uintptr) error {
	return enumDesktopWindows(uintptr(hDesktop), lpEnumFunc, lParam)
}

func FillRect(hdc win.HDC, rect win.RECT, hbr win.HBRUSH) error {
	return fillRect(uintptr(hdc), uintptr(unsafe.Pointer(&rect.Left)), uintptr(hbr))
}

func FindWindow(lpClassName, lpWindowName *uint16) win.HWND {
	return win.FindWindow(lpClassName, lpWindowName)
}

func FindWindowEx(hwndParent win.HWND, hwndChildAfter win.HWND, lpszClass *uint16, lpszWindow *uint16) (hwnd win.HWND) {
	return win.HWND(findWindowEx(uintptr(hwndParent), uintptr(hwndChildAfter), lpszClass, lpszWindow))
}

func GetClassName(hwnd win.HWND, lpClassName uintptr, nMax int) (length int) {
	return getClassName(uintptr(hwnd), lpClassName, nMax)
}

func GetWindowText(hwnd win.HWND, lpString []uint16, nMax int) (length int) {
	return getWindowText(uintptr(hwnd), uintptr(unsafe.Pointer(&lpString[0])), nMax)
}

func InvalidateRect(hwnd win.HWND, rect win.RECT, bErase bool) error {
	return invalidateRect(uintptr(hwnd), uintptr(unsafe.Pointer(&rect.Left)), bErase)
}

func SetLayeredWindowAttributes(hwnd win.HWND, color uint32, bAlpha byte, dwFlags uint32) error {
	return setLayeredWindowAttributes(uintptr(hwnd), color, bAlpha, dwFlags)
}

func SetWindowRgn(hwnd win.HWND, hRgn win.HRGN, bRedraw bool) error {
	return setWindowRgn(uintptr(hwnd), uintptr(hRgn), bRedraw)
}

func SetWindowText(hwnd win.HWND, lpString *uint16) error {
	return setWindowText(uintptr(hwnd), lpString)
}

func RegisterClassEx(windowClass *win.WNDCLASSEX) (win.ATOM, error) {
	a, err := registerClassEx(uintptr(unsafe.Pointer(&windowClass)))
	return win.ATOM(a), err
}

func ShowCursor(state bool) (counter int) {
	return showCursor(state)
}

func ShowWindow(hWnd win.HWND, nCmdShow int32) bool {
	return win.ShowWindow(hWnd, nCmdShow)
}

func UpdateLayeredWindow(hwnd win.HWND, hdcDst win.HDC, pptDst win.POINT, psize uintptr, hdcSrc win.HDC, pptSrc win.POINT, crKey uint32, pblend win.BLENDFUNCTION, dwFlags uint32) (ok bool) {
	return updateLayeredWindow(uintptr(hwnd), uintptr(hdcDst), uintptr(unsafe.Pointer(&pptDst.X)), psize, uintptr(hdcSrc), uintptr(unsafe.Pointer(&pptSrc.X)), crKey, uintptr(unsafe.Pointer(&pblend.BlendOp)), dwFlags)
}

func UpdateWindow(hwnd win.HWND) bool {
	return win.UpdateWindow(hwnd)
}

func GetWindowRect(hwnd win.HWND, rect *win.RECT) bool {
	return win.GetWindowRect(hwnd, rect)
}

func GetCursorPos(lpPoint *win.POINT) bool {
	return win.GetCursorPos(lpPoint)
}

func SetForegroundWindow(hWnd win.HWND) bool {
	return win.SetForegroundWindow(hWnd)
}

func MapVirtualKey(uCode uint32, uMapType uint32) (code uint32) {
	return mapVirtualKey(uCode, uMapType)
}

type HOOKPROC func(
	code int, wParam uintptr, lParam uintptr,
) uintptr

type HHOOK uintptr

func SetWindowHookEx(idHook TYPE_HOOK_ID, lpfn HOOKPROC, hmod win.HINSTANCE, dwThreadId uint32) (hhk HHOOK, err error) {
	h, err := setWindowHookEx(int(idHook), syscall.NewCallback(lpfn), uintptr(hmod), dwThreadId)
	return HHOOK(h), err
}

func UnhookWindowsHookEx(hhk HHOOK) error {
	return unhookWindowsHookEx(uintptr(hhk))
}

func CallNextHookEx(hhk HHOOK, nCode int, wParam uintptr, lParam uintptr) uintptr {
	return callNextHookEx(hhk, nCode, wParam, lParam)
}

type MOD_KEY uint32

const (
	// Either ALT key must be held down.
	MOD_ALT MOD_KEY = 0x0001

	// Either CTRL key must be held down.
	MOD_CONTROL MOD_KEY = 0x0002

	// Changes the hotkey behavior so that the keyboard auto-repeat does not yield multiple hotkey notifications.
	MOD_NOREPEAT MOD_KEY = 0x4000

	// Either SHIFT key must be held down.
	MOD_SHIFT MOD_KEY = 0x0004

	// Either WINDOWS key was held down. These keys are labeled with the Windows logo. Keyboard shortcuts that involve the WINDOWS key are reserved for use by the operating system.
	MOD_WIN MOD_KEY = 0x0008
)

func RegisterHotKey(hWnd win.HWND, id int, fsModifiers MOD_KEY, vk uint32) error {
	return registerHotKey(uintptr(hWnd), id, uint32(fsModifiers), vk)
}
