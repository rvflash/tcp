package tcp

import (
	"io"
	"strconv"
	"testing"

	"github.com/matryer/is"
)

const (
	eol        = "\n"
	msg        = "hello world"
	msgWithEol = msg + eol
)

func TestContext_Write(t *testing.T) {
	var (
		dt = []struct {
			w    *ResponseRecorder
			in   []byte
			out  string
			size int
			err  error
		}{
			{err: io.EOF},
			{w: NewRecorder()},
			{w: NewRecorder(), in: []byte(msg), size: len(msg), out: msg},
			{w: NewRecorder(), in: []byte(msgWithEol), size: len(msgWithEol), out: msgWithEol},
		}
		are = is.New(t)
	)
	for i, tt := range dt {
		tt := tt
		t.Run("#"+strconv.Itoa(i), func(t *testing.T) {
			c := newContextWithWriter(tt.w)
			size, err := c.Write(tt.in)
			are.Equal(err, tt.err)   // mismatch error
			are.Equal(size, tt.size) // mismatch size
			if tt.w != nil {
				are.Equal(tt.w.Body.String(), tt.out) // mismatch result
			}
		})
	}
}

func TestContext_String(t *testing.T) {
	var (
		dt = []struct {
			w *ResponseRecorder
			in,
			out,
			err string
		}{
			{err: io.EOF.Error()},
			{w: NewRecorder(), out: eol},
			{w: NewRecorder(), in: msg, out: msgWithEol},
			{w: NewRecorder(), in: msgWithEol, out: msgWithEol},
		}
		are = is.New(t)
	)
	for i, tt := range dt {
		tt := tt
		t.Run("#"+strconv.Itoa(i), func(t *testing.T) {
			c := newContextWithWriter(tt.w)
			c.String(tt.in)
			are.Equal(c.Err().Error(), tt.err) // mismatch error
			if tt.w != nil {
				are.Equal(tt.w.Body.String(), tt.out) // mismatch result
			}
		})
	}
}

func newContextWithWriter(w io.WriteCloser) *Context {
	c := &Context{
		Request: NewRequest(ACK, nil),
		writer:  responseWriter{ResponseWriter: w},
	}
	if w != nil {
		c.reset()
	}
	return c
}
