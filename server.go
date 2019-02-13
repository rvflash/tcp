// Package tcp provides interfaces to create a TCP server.
package tcp

import (
	"context"
	"net"
	"sync"
	"time"
)

// HandlerFunc defines the handler interface.
type HandlerFunc func(c *Context)

// Router is implemented by the Server.
type Router interface {
	// Any registers a route that matches one of supported method
	Any(method string, handler ...HandlerFunc) Router
	// Use adds middleware fo any context: start and end of connection and message.
	Use(handler ...HandlerFunc) Router
	// ACK is a shortcut for Any("ACK", ...HandlerFunc).
	ACK(handler ...HandlerFunc) Router
	// FIN is a shortcut for Any("FIN", ...HandlerFunc).
	FIN(handler ...HandlerFunc) Router
	// SYN is a shortcut for Any("SYN", ...HandlerFunc).
	SYN(handler ...HandlerFunc) Router
}

// List of supported "methods".
const (
	ANY = ""
	ACK = "ACK"
	FIN = "FIN"
	SYN = "SYN"
)

// Default returns an instance of TCP server with a Logger and a Recover on panic attached.
func Default() *Server {
	// todo :)
	return New()
}

// New returns a new instance of a TCP server.
func New() *Server {
	s := &Server{
		handlers: map[string][]HandlerFunc{},
	}
	s.pool.New = func() interface{} {
		return s.allocateContext()
	}
	return s
}

// Server is the TCP server. It contains
type Server struct {
	// ReadTimeout is the maximum duration for reading the entire request, including the body.
	// A zero value for t means Read will not time out.
	ReadTimeout time.Duration

	handlers map[string][]HandlerFunc
	pool     sync.Pool
}

func (s *Server) allocateContext() *Context {
	return &Context{}
}

func (s *Server) computeHandlers(handlers []HandlerFunc) []HandlerFunc {
	m := make([]HandlerFunc, len(s.handlers[ANY])+len(handlers))
	copy(m, s.handlers[ANY])
	copy(m[len(s.handlers[ANY]):], handlers)
	return m
}

// Any
func (s *Server) Any(method string, f ...HandlerFunc) Router {
	switch method {
	case ACK:
		return s.ACK(f...)
	case FIN:
		return s.FIN(f...)
	case SYN:
		return s.SYN(f...)
	default:
		return s.Use(f...)
	}
}

// ACK allows to handle each new message.
func (s *Server) ACK(f ...HandlerFunc) Router {
	s.handlers[ACK] = append(s.handlers[ACK], f...)
	return s
}

// FIN allows to handle when the connection is closed.
func (s *Server) FIN(f ...HandlerFunc) Router {
	s.handlers[FIN] = append(s.handlers[FIN], f...)
	return s
}

// SYN allows to handle each new connection.
func (s *Server) SYN(f ...HandlerFunc) Router {
	s.handlers[SYN] = append(s.handlers[SYN], f...)
	return s
}

// Use adds middleware(s).
func (s *Server) Use(f ...HandlerFunc) Router {
	s.handlers[ANY] = append(s.handlers[ANY], f...)
	return s
}

// Run starts listening on TCP address.
// This method will block the calling goroutine indefinitely unless an error happens.
func (s *Server) Run(addr string) (err error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return
	}
	defer func() {
		if err == nil {
			err = l.Close()
		}
	}()
	ctx := context.Background()
	for {
		c, err := newConn(l, s.ReadTimeout)
		if err != nil {
			return
		}
		r, err := newRequest(c, ctx)
		if err != nil {
			return
		}
		w := newResponseWriter(c)
		// flag request as SYN and repeat the action...
		go s.ServeTCP(w, r)
	}
}

func (s *Server) handle(c *Context) {
	if c.handlers == nil {
		return
	}
	c.Next()
}

func newConn(l net.Listener, to time.Duration) (net.Conn, error) {
	c, err := l.Accept()
	if err != nil {
		return nil, err
	}
	if to == 0 {
		// no read deadline required.
		return c, err
	}
	err = c.SetReadDeadline(time.Now().Add(to))
	if err != nil {
		return nil, err
	}
	return c, nil
}

func newRequest(c net.Conn, ctx context.Context) (*Request, error) {
	req, err := NewRequest(SYN, nil)
	if err != nil || c == nil || ctx == nil {
		return nil, err
	}
	// Retrieves the remote address of the client.
	req.RemoteAddr = c.RemoteAddr().String()

	// Initiates the connection with a context by cancellation.
	return req.WithCancel(ctx), nil
}

// ServeTCP ...
func (s *Server) ServeTCP(w ResponseWriter, req *Request) {
	c := s.pool.Get().(*Context)
	c.writer.rebase(w)
	c.Request = req
	c.reset()
	s.handle(c)
	s.pool.Put(c)
}
