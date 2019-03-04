package tcp

import (
	"bytes"
	"io"
)

// ResponseWriter interface is used by a TCP handler to write the response.
type ResponseWriter interface {
	// Size returns the number of bytes already written into the response body.
	// -1: not already written
	Size() int
	io.WriteCloser
}

func newWriter(wc io.WriteCloser) *responseWriter {
	return &responseWriter{
		ResponseWriter: wc,
		size:           noWritten,
	}
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
	if r.size == noWritten {
		r.size = 0
	}
	r.size += n
}

func (r *responseWriter) rebase(w io.WriteCloser) {
	r.ResponseWriter = w
	r.size = noWritten
}

// ResponseRecorder is an implementation of http.ResponseWriter that records its changes.
type ResponseRecorder struct {
	// Body is the buffer to which the Handler's Write calls are sent.
	Body *bytes.Buffer
}

// NewRecorder returns an initialized writer to record the response.
func NewRecorder() *ResponseRecorder {
	return &ResponseRecorder{
		Body: new(bytes.Buffer),
	}
}

// Close implements the ResponseWriter interface.
func (r *ResponseRecorder) Close() error {
	return nil
}

// Size implements the ResponseWriter interface.
func (r *ResponseRecorder) Size() int {
	if r.Body == nil {
		return noWritten
	}
	return r.Body.Len()
}

// Write implements the ResponseWriter interface.
func (r *ResponseRecorder) Write(p []byte) (n int, err error) {
	if r.Body == nil {
		return 0, io.EOF
	}
	n, err = r.Body.Write(p)
	return
}

// WriteString allows to directly write string.
func (r *ResponseRecorder) WriteString(s string) (n int, err error) {
	if r.Body == nil {
		return 0, io.EOF
	}
	n, err = r.Body.WriteString(s)
	return
}
