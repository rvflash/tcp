// Package tcp provides interfaces to create a TCP server.
package tcp

import (
	"context"
	"crypto/tls"
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
	f := logrus.Fields{
		LogLatency:        0,
		LogServerHostname: "",
		LogRemoteAddr:     "",
		LogRequestSize:    0,
		LogResponseSize:   0,
	}
	l := logrus.New()
	l.Formatter = &logrus.TextFormatter{DisableTimestamp: true}
	h := New()
	h.Use(Logger(l, f))
	h.Use(Recovery())
	return h
}

// New returns a new instance of a TCP server.
func New() *Server {
	s := &Server{
		handlers: map[string][]HandlerFunc{},
		shutdown: make(chan struct{}),
		closed:   make(chan struct{}),
	}
	s.pool.New = func() interface{} {
		return s.allocateContext()
	}
	return s
}

func (s *Server) allocateContext() *Context {
	return &Context{srv: s}
}

// Server is the TCP server. It contains
type Server struct {
	// ReadTimeout is the maximum duration for reading the entire request, including the body.
	// A zero value for t means Read will not time out.
	ReadTimeout time.Duration

	cancel   context.CancelFunc
	handlers map[string][]HandlerFunc
	pool     sync.Pool
	closed,
	shutdown chan struct{}
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

const network = "tcp"

// Run starts listening on TCP address.
// This method will block the calling goroutine indefinitely unless an error happens.
func (s *Server) Run(addr string) error {
	l, err := net.Listen(network, addr)
	if err != nil {
		return err
	}
	return s.serve(l)
}

// RunTLS acts identically to the Run method, except that it uses the TLS protocol.
// This method will block the calling goroutine indefinitely unless an error happens.
func (s *Server) RunTLS(addr, certFile, keyFile string) error {
	c, err := tlsConfig(certFile, keyFile)
	if err != nil {
		return err
	}
	l, err := tls.Listen(network, addr, c)
	if err != nil {
		return err
	}
	return s.serve(l)
}

func (s *Server) serve(l net.Listener) (err error) {
	var (
		w8  sync.WaitGroup
		ctx context.Context
	)
	ctx, s.cancel = context.WithCancel(context.Background())
	defer func() {
		s.cancel()
		cErr := l.Close()
		if err != nil {
			err = cErr
		}
	}()
	for {
		select {
		case <-s.shutdown:
			// Stops listening but does not interrupt any active connections.
			// See the Shutdown method to gracefully shuts down the server.
			w8.Wait()
			close(s.closed)
			return
		default:
		}
		var c net.Conn
		c, err = read(l, s.ReadTimeout)
		if err != nil {
			return
		}
		rwc := s.newConn(c)
		w8.Add(1)
		go func() {
			defer w8.Done()
			rwc.serve(ctx)
		}()
	}
}

func (s *Server) newConn(c net.Conn) *conn {
	return &conn{
		addr: c.RemoteAddr().String(),
		srv:  s,
		rwc:  c,
	}
}

// ServeTCP implements the Handler interface;
func (s *Server) ServeTCP(w ResponseWriter, req *Request) {
	ctx := s.pool.Get().(*Context)
	ctx.writer.rebase(w)
	ctx.Request = req
	ctx.reset()
	s.handle(ctx)
	s.pool.Put(ctx)
}

func (s *Server) handle(ctx *Context) {
	ctx.handlers = s.computeHandlers(ctx.Request.Segment)
	if len(ctx.handlers) == 0 {
		return
	}
	ctx.Next()
}

func (s *Server) computeHandlers(segment string) []HandlerFunc {
	m := make([]HandlerFunc, len(s.handlers[ANY])+len(s.handlers[segment]))
	copy(m, s.handlers[ANY])
	copy(m[len(s.handlers[ANY]):], s.handlers[segment])
	return m
}

// Shutdown gracefully shuts down the server without interrupting any
// active connections. Shutdown works by first closing all open listeners and
// then waiting indefinitely for connections to return to idle and then shut down.
// If the provided context expires before the shutdown is complete,
// Shutdown returns the context's error.
func (s *Server) Shutdown(ctx context.Context) error {
	if s.shutdown == nil {
		// Nothing to do
		return nil
	}
	// Stops listening.
	close(s.shutdown)

	// Stops all.
	for {
		select {
		case <-ctx.Done():
			// Forces closing of actives connections.
			s.cancel()
			return ctx.Err()
		case <-s.closed:
			return nil
		}
	}
}

func tlsConfig(certFile, keyFile string) (*tls.Config, error) {
	var err error
	c := make([]tls.Certificate, 1)
	c[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	return &tls.Config{Certificates: c}, err
}

func read(l net.Listener, to time.Duration) (net.Conn, error) {
	c, err := l.Accept()
	if err != nil {
		return nil, err
	}
	if to == 0 {
		return c, nil
	}
	err = c.SetReadDeadline(time.Now().Add(to))
	if err != nil {
		return nil, err
	}
	return c, nil
}
