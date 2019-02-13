package tcp

import (
	"context"
	"io"
	"io/ioutil"
)

// Request represents an TCP request.
type Request struct {
	// Method specifies the TCP step (SYN, ACK, FIN).
	Method string

	// Body is the request's body.
	Body io.ReadCloser

	// RemoteAddr returns the remote network address.
	RemoteAddr string

	ctx    context.Context
	cancel context.CancelFunc
}

// Close implements the io.Closer interface.
func (r *Request) Cancel() {
	if r.cancel != nil {
		r.cancel()
	}
}

// Closed implements the Conn interface.
func (r *Request) Closed() <-chan struct{} {
	return r.Context().Done()
}

// Context returns the request's context.
func (r *Request) Context() context.Context {
	if r.ctx != nil {
		return r.ctx
	}
	return context.Background()
}

// WithCancel returns a shallow copy of the given request with its context changed to ctx.
func (r *Request) WithCancel(ctx context.Context) *Request {
	if ctx == nil {
		// awkward: nothing to do
		return r
	}
	r2 := new(Request)
	*r2 = *r
	r2.ctx, r2.cancel = context.WithCancel(ctx)
	return r2
}

// NewRequest ...
func NewRequest(method string, body io.Reader) (*Request, error) {
	if method == "" {
		return nil, ErrRequest
	}
	rc, ok := body.(io.ReadCloser)
	if !ok && body != nil {
		rc = ioutil.NopCloser(body)
	}
	req := &Request{
		Method: method,
		Body:   rc,
	}
	return req, nil
}
