package tcp

import (
	"io"
)

// ResponseWriter interface is used by a TCP handler to write the response.
type ResponseWriter interface {
	io.WriteCloser
	// Size returns the number of bytes already written into the response body.
	// -1: not already written
	Size() int
}

type responseWriter struct {
	ResponseWriter io.WriteCloser
	size           int
}

const noWritten = -1

// Close implements the ResponseWriter interface.
func (r *responseWriter) Close() error {
	return r.ResponseWriter.Close()
}

// Size implements the ResponseWriter interface.
func (r *responseWriter) Size() int {
	return r.size
}

// Write implements the ResponseWriter interface.
func (r *responseWriter) Write(p []byte) (n int, err error) {
	n, err = r.ResponseWriter.Write(p)
	r.incr(n)
	return
}

// WriteString allows to directly write string.
func (r *responseWriter) WriteString(s string) (n int, err error) {
	n, err = io.WriteString(r.ResponseWriter, s)
	r.incr(n)
	return
}

func (r *responseWriter) incr(n int) {
	if n == noWritten {
		n = 0
	}
	r.size += n
}

func (r *responseWriter) rebase(w ResponseWriter) {
	r.ResponseWriter = w
	r.size = noWritten
}
