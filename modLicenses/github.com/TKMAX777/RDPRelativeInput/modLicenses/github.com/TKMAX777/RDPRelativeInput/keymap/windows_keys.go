package keymap

import "github.com/pkg/errors"

func GetWindowsKeyDetail(windowsVirtualKeyCode uint32) (*WindowsKey, error) {
	w, ok := windowsKeys[windowsVirtualKeyCode]
	if !ok {
		return nil, errors.New("NotFound")
	}
	return &w, nil
}

func GetWindowsKeyDetailFromEventInput(eventInput string) (*WindowsKey, error) {
	w, ok := windowsKeysFromEventInput[eventInput]
	if !ok {
		return nil, errors.New("NotFound")
	}
	return &w, nil
}
