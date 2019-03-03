package tcp

import (
	"context"
	"io"
	"io/ioutil"
)

// Request represents an TCP request.
type Request struct {
	// Segment specifies the TCP segment (SYN, ACK, FIN).
	Segment string
	// Body is the request's body.
	Body io.ReadCloser
	// remoteAddr returns the remote network address.
	RemoteAddr string

	// Context of the request.
	ctx    context.Context
	cancel context.CancelFunc
}

// Cancel closes the request.
func (r *Request) Cancel() {
	if r.cancel != nil {
		r.cancel()
	}
}

// Canceled listens the context of the request until its closing.
func (r *Request) Canceled() <-chan struct{} {
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

// NewRequest returns a new instance of request.
// A segment is mandatory as input. If empty, a SYN segment is used.
func NewRequest(segment string, body io.Reader) *Request {
	if segment == "" {
		// by default, we use the SYN segment.
		segment = SYN
	}
	req := &Request{Segment: segment}
	if body != nil {
		rc, ok := body.(io.ReadCloser)
		if !ok {
			rc = ioutil.NopCloser(body)
		}
		req.Body = rc
	}
	return req
}
