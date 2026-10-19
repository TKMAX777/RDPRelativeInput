// Package channelbuffer provides bounded, nonblocking callback writes and blocking reads.
package channelbuffer

import (
	"context"
	"errors"
	"io"
	"sync"
)

const Capacity = 64 * 1024

var ErrOverflow = errors.New("virtual channel receive buffer overflow")

type Buffer struct {
	mu         sync.Mutex
	ready      *sync.Cond
	data       [Capacity]byte
	head, size int
	err        error
	ctx        context.Context
	cancel     context.CancelFunc
}

func New() *Buffer {
	ctx, cancel := context.WithCancel(context.Background())
	b := &Buffer{ctx: ctx, cancel: cancel}
	b.ready = sync.NewCond(&b.mu)
	return b
}
func (b *Buffer) Context() context.Context { return b.ctx }
func (b *Buffer) Fail(err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failLocked(err)
}
func (b *Buffer) failLocked(err error) {
	if b.err != nil {
		return
	}
	if err == nil {
		err = io.EOF
	}
	b.err = err
	b.size = 0
	b.cancel()
	b.ready.Broadcast()
}
func (b *Buffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.err != nil {
		return 0, b.err
	}
	if len(p) > Capacity-b.size {
		b.failLocked(ErrOverflow)
		return 0, ErrOverflow
	}
	tail := (b.head + b.size) % Capacity
	n := copy(b.data[tail:], p)
	copy(b.data[:], p[n:])
	b.size += len(p)
	b.ready.Signal()
	return len(p), nil
}
func (b *Buffer) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	for b.size == 0 && b.err == nil {
		b.ready.Wait()
	}
	if b.err != nil {
		return 0, b.err
	}
	n := len(p)
	if n > b.size {
		n = b.size
	}
	first := n
	if first > Capacity-b.head {
		first = Capacity - b.head
	}
	copy(p[:first], b.data[b.head:b.head+first])
	copy(p[first:n], b.data[:n-first])
	b.head = (b.head + n) % Capacity
	b.size -= n
	return n, nil
}
