package keymap

import "github.com/pkg/errors"

func GetLinuxKeyDetail(linuxVirtualKeyCode uint32) (*LinuxKey, error) {
	l, ok := linuxKeys[linuxVirtualKeyCode]
	if !ok {
		return nil, errors.New("NotFound")
	}
	return &l, nil
}

func GetLinuxKeyDetailFromEventInput(eventInput string) (*LinuxKey, error) {
	l, ok := linuxKeysFromEventInput[eventInput]
	if !ok {
		return nil, errors.New("NotFound")
	}
	return &l, nil
}
