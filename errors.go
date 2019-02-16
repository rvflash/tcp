package tcp

import (
	"strings"
)

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
	msg   string
	cause error
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e.cause == nil {
		return "tcp: " + e.msg
	}
	return "tcp: " + e.msg + ": " + e.cause.Error()
}

// Errors contains the list of errors occurred during the request.
type Errors []error

// Error implements the error interface.
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
