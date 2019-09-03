package tcp

import (
	"strings"
)

// Err represents a TCP error.
type Err interface {
	error
	// Recovered returns true if the error comes from a panic recovering.
	Recovered() bool
}

// List of common errors
var (
	// ErrRequest is returned if the request is invalid.
	ErrRequest = NewError("invalid request")
)

// NewError returns a new Error based of the given cause.
func NewError(msg string, cause ...error) *Error {
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
	const prefix = "tcp: "
	if e.cause == nil {
		return prefix + e.msg
	}
	return prefix + e.msg + ": " + e.cause.Error()
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
	var err Err
	for _, r := range e {
		err, ok = r.(Err)
		if ok && err.Recovered() {
			return
		}
	}
	return
}
