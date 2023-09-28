// Credit: https://github.com/ranna-go/ranna/blob/master/pkg/cappedbuffer/cappedbuffer.go
package cappedbuffer

import (
	"bytes"
	"errors"
)

type CappedBuffer struct {
	*bytes.Buffer
	cap int
}

func New(buf []byte, cap int) *CappedBuffer {
	return &CappedBuffer{
		Buffer: bytes.NewBuffer(buf),
		cap:    cap,
	}
}

func (cb *CappedBuffer) Write(p []byte) (n int, err error) {
	if cb.cap > 0 && cb.Len()+len(p) > cb.cap {
		return 0, errors.New("buffer overflow")
	}
	return cb.Buffer.Write(p)
}
