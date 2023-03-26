package host

type Handler struct {
	isCapturing    bool
	captureHandler *CaptureHandler
}

func New() *Handler {
	return &Handler{}
}
