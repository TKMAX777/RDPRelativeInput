// Package protocol implements the bounded handshake and binary input protocol.
package protocol

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

const MaxLineSize = 256

type Reader struct {
	reader *bufio.Reader
	frame  [8]byte
}

func NewReader(r io.Reader) *Reader {
	// WTS requires a receive buffer of at least CHANNEL_CHUNK_LENGTH (1600).
	return &Reader{reader: bufio.NewReaderSize(r, 4096)}
}

func (r *Reader) ReadLine() (string, error) {
	line, err := r.reader.ReadSlice('\n')
	if len(line) > MaxLineSize || err == bufio.ErrBufferFull {
		return "", fmt.Errorf("protocol line exceeds %d bytes", MaxLineSize)
	}
	if err != nil {
		if err == io.EOF && len(line) != 0 {
			return "", io.ErrUnexpectedEOF
		}
		return "", err
	}
	return strings.TrimSuffix(string(line[:len(line)-1]), "\r"), nil
}

func ReadVersion(r io.Reader, version int) error {
	line, err := NewReader(r).ReadLine()
	if err != nil {
		return err
	}
	if line != fmt.Sprintf("ProtocolVersion:%d", version) {
		return fmt.Errorf("unsupported protocol version")
	}
	return nil
}
