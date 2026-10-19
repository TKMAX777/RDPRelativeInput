package main

import (
	"io"
	"sync"
	"unsafe"

	"github.com/TKMAX777/RDPRelativeInput/internal/channelbuffer"
)
import "C"

type pendingWrite struct {
	data []byte
	done chan bool
}

type VirtualChannelReadWriteCloser struct {
	channelHandle uint32
	entryPoints   CHANNEL_ENTRY_POINTS
	received      *channelbuffer.Buffer
	writeMu       sync.Mutex
	mu            sync.Mutex
	nextID        uintptr
	pending       map[uintptr]*pendingWrite
	closeOnce     sync.Once
	closeErr      error
}

func NewChannel(entryPoints CHANNEL_ENTRY_POINTS, ppInitHandle *HANDLE, name string) (*VirtualChannelReadWriteCloser, error) {
	rw := &VirtualChannelReadWriteCloser{entryPoints: entryPoints, received: channelbuffer.New(), pending: make(map[uintptr]*pendingWrite)}
	if err := entryPoints.VirtualChannelOpen.Call(ppInitHandle, &rw.channelHandle, name, rw.getVirtualChannelOpenEventFn()); err != nil {
		rw.received.Fail(err)
		return nil, err
	}
	return rw, nil
}
func (rw *VirtualChannelReadWriteCloser) Read(p []byte) (int, error) { return rw.received.Read(p) }
func (rw *VirtualChannelReadWriteCloser) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	if len(p) > 32*1024 {
		return 0, io.ErrShortBuffer
	}
	rw.writeMu.Lock()
	defer rw.writeMu.Unlock()
	select {
	case <-rw.received.Context().Done():
		return 0, io.ErrClosedPipe
	default:
	}
	// Own the data until the native completion/cancellation notification. Closing
	// wakes the writer but must not release memory still in use by RDP.
	w := &pendingWrite{data: append([]byte(nil), p...), done: make(chan bool, 1)}
	rw.mu.Lock()
	rw.nextID++
	id := rw.nextID
	rw.pending[id] = w
	rw.mu.Unlock()
	err := rw.entryPoints.VirtualChannelWrite.Call(rw.channelHandle, w.data, len(w.data), id)
	if err != nil {
		rw.complete(id, false)
		return 0, err
	}
	select {
	case ok := <-w.done:
		if !ok {
			return 0, io.ErrClosedPipe
		}
		return len(p), nil
	case <-rw.received.Context().Done():
		return 0, io.ErrClosedPipe
	}
}
func (rw *VirtualChannelReadWriteCloser) complete(id uintptr, ok bool) {
	rw.mu.Lock()
	w := rw.pending[id]
	delete(rw.pending, id)
	rw.mu.Unlock()
	if w != nil {
		w.done <- ok
	}
}
func (rw *VirtualChannelReadWriteCloser) Close() error {
	rw.received.Fail(io.EOF)
	rw.closeOnce.Do(func() { rw.closeErr = rw.entryPoints.VirtualChannelClose.Call(rw.channelHandle) })
	return rw.closeErr
}
func (rw *VirtualChannelReadWriteCloser) getVirtualChannelOpenEventFn() CHANNEL_OPEN_EVENT_FN {
	return func(_ uint32, event CHANNEL_EVENT_TYPE, pData *byte, length, total uint32, flags CHANNEL_FLAG_TYPE) uintptr {
		switch event {
		case CHANNEL_EVENT_DATA_RECEIVED:
			if flags&(CHANNEL_FLAG_SUSPEND|CHANNEL_FLAG_RESUME) != 0 {
				return 0
			}
			if length > channelbuffer.Capacity || total > channelbuffer.Capacity || length > total || (length > 0 && pData == nil) {
				rw.received.Fail(channelbuffer.ErrOverflow)
				return 0
			}
			_, _ = rw.received.Write(unsafe.Slice(pData, int(length)))
		case CHANNEL_EVENT_WRITE_CANCELLED:
			rw.complete(uintptr(unsafe.Pointer(pData)), false)
		case CHANNEL_EVENT_WRITE_COMPLETE:
			rw.complete(uintptr(unsafe.Pointer(pData)), true)
		}
		return 0
	}
}
