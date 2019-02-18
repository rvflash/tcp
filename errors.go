package tcp

import (
	"strings"
)

// Err represents a TCP error.
type Err interface {
	// Recovered returns true if the error comes from a panic recovering.
	Recovered() bool
	error
}

// ErrRequest is returned if the request is invalid.
var ErrRequest = NewError("invalid request")

// NewError returns a new Error based of the given cause.
func NewError(msg string, cause ...error) error {
	if cause == nil {
		return &Error{msg: msg}
	}
	return &Error{msg: msg, cause: cause[0]}
}

// Error represents a error message.
// It can wraps another error, its cause.
type Error struct {
	msg     string
	cause   error
	recover bool
}

// Error implements the Err interface.
func (e *Error) Error() string {
	if e.cause == nil {
		return "tcp: " + e.msg
	}
	return "tcp: " + e.msg + ": " + e.cause.Error()
}

// Recovered implements the Err interface.
func (e *Error) Recovered() bool {
	return e.recover
}

// Errors contains the list of errors occurred during the request.
type Errors []error

// Error implements the Err interface.
func (e Errors) Error() string {
	var (
		b   strings.Builder
		err error
	)
	for i, r := range e {
		if i > 0 {
			if _, err = b.WriteString(", "); err != nil {
				return err.Error()
			}
		}
		if _, err = b.WriteString(r.Error()); err != nil {
			return err.Error()
		}
	}
	return b.String()
}

// Recovered implements the Err interface.
func (e Errors) Recovered() (ok bool) {
	for _, r := range e {
		_, ok = r.(*Error)
		if ok {
			return
		}
	}
	return
}
