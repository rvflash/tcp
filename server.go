// Package tcp provides interfaces to create a TCP server.
package tcp

import (
	"context"
	"net"

	"github.com/sirupsen/logrus"
)

// HandlerFunc defines the handler interface.
type HandlerFunc func(c Conn)

// Logger must be implemented by any logger.
type Logger interface {
	Errorf(format string, v ...interface{})
	Printf(format string, v ...interface{})
}

// Default returns an instance of TCP server with a Logger attached.
func Default() *Server {
	return &Server{log: logrus.New()}
}

// New returns a new instance of a TCP server.
func New() *Server {
	return &Server{}
}

// Server is the TCP server. It contains
type Server struct {
	in, out, msg []HandlerFunc
	log          Logger
}

func (s *Server) errorf(format string, a ...interface{}) {
	if s.log == nil {
		return
	}
	s.log.Errorf(format, a...)
}

func (s *Server) printf(format string, a ...interface{}) {
	if s.log == nil {
		return
	}
	s.log.Printf(format, a...)
}

// SYN allows to handle each new connection / client.
func (s *Server) SYN(f ...HandlerFunc) {
	if f != nil {
		s.in = append(s.in, f...)
	}
}

// ACK allows to handle each new message.
func (s *Server) ACK(f ...HandlerFunc) {
	if f != nil {
		s.msg = append(s.msg, f...)
	}
}

// FIN allows to handle when the client connection is closed.
func (s *Server) FIN(f ...HandlerFunc) {
	if f != nil {
		s.out = append(s.out, f...)
	}
}

// Run starts listening on TCP address.
// This method will block the calling goroutine indefinitely unless an error happens.
func (s *Server) Run(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer func() {
		if err := l.Close(); err != nil {
			s.errorf("tcp closing failed with %q", err)
		}
	}()
	s.printf("tcp listens on %s", addr)

	ctx := context.Background()
	for {
		// new connection
		c, err := l.Accept()
		if err != nil {
			return err
		}
		x := &client{
			conn: c,
			srv:  s,
		}
		go x.handle(ctx)
	}
}
