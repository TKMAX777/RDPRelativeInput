package client

import "github.com/pkg/errors"

func (h *Handler) SetToggleKey(k string) error {
	if k == "" {
		return errors.New("NotSpecified")
	}
	h.options.toggleKey = k
	return nil
}
func (h *Handler) SetToggleType(t ToggleType) {
	h.options.toggleType = t
}
