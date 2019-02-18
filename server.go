// Package tcp provides interfaces to create a TCP server.
package tcp

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// Handler responds to a TCP request.
type Handler interface {
	ServeTCP(ResponseWriter, *Request)
}

// HandlerFunc defines the handler interface used as return value.
type HandlerFunc func(c *Context)

// Router is implemented by the Server.
type Router interface {
	// Any registers a route that matches one of supported segment
	Any(segment string, handler ...HandlerFunc) Router
	// Use adds middleware fo any context: start and end of connection and message.
	Use(handler ...HandlerFunc) Router
	// ACK is a shortcut for Any("ACK", ...HandlerFunc).
	ACK(handler ...HandlerFunc) Router
	// FIN is a shortcut for Any("FIN", ...HandlerFunc).
	FIN(handler ...HandlerFunc) Router
	// SYN is a shortcut for Any("SYN", ...HandlerFunc).
	SYN(handler ...HandlerFunc) Router
}

// List of supported segments.
const (
	ANY = ""
	ACK = "ACK"
	FIN = "FIN"
	SYN = "SYN"
)

// Default returns an instance of TCP server with a Logger and a Recover on panic attached.
func Default() *Server {
	// Adds a logger.
	l := logrus.New()
	l.Formatter = &logrus.TextFormatter{DisableTimestamp: true}
	f := logrus.Fields{
		Latency:        0,
		Hostname:       "",
		RemoteAddr:     "",
		RequestLength:  0,
		ResponseLength: 0,
	}
	h := New()
	h.Use(Logger(l, f))
	return h
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
	return &Context{srv: s}
}

func (s *Server) computeHandlers(segment string) []HandlerFunc {
	m := make([]HandlerFunc, len(s.handlers[ANY])+len(s.handlers[segment]))
	copy(m, s.handlers[ANY])
	copy(m[len(s.handlers[ANY]):], s.handlers[segment])
	return m
}

// Any attaches handlers on the given segment.
func (s *Server) Any(segment string, f ...HandlerFunc) Router {
	switch segment {
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

// Use adds middleware(s) on all segments.
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
			return err
		}
		rwc := s.newConn(ctx, c)
		go rwc.serve()
	}
}

func newConn(l net.Listener, to time.Duration) (net.Conn, error) {
	c, err := l.Accept()
	if err != nil {
		return nil, err
	}
	if to == 0 {
		return c, err
	}
	err = c.SetReadDeadline(time.Now().Add(to))
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Server) newConn(ctx context.Context, c net.Conn) *conn {
	return &conn{
		addr: c.RemoteAddr().String(),
		ctx:  ctx,
		srv:  s,
		rwc:  c,
	}
}

func (s *Server) put(c *Context) {
	s.pool.Put(c)
}

func (s *Server) get() *Context {
	return s.pool.Get().(*Context)
}
