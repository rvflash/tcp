package tcp_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/matryer/is"
	"github.com/rvflash/tcp"
)

const (
	hiWorld = "hello world"
	prefix  = "tcp: "
)

func TestNewError(t *testing.T) {
	var (
		dt = []struct {
			in  string
			err error
			out string
		}{
			{out: prefix},
			{in: "hi!", out: prefix + "hi!"},
			{in: "hello", err: errors.New("world"), out: prefix + "hello: world"},
		}
		are = is.New(t)
		err error
	)
	for i, tt := range dt {
		t.Run("#"+strconv.Itoa(i), func(t *testing.T) {
			err = tcp.NewError(tt.in, tt.err)
			are.Equal(err.Error(), tt.out)
		})
	}
}

func TestError_Recovered(t *testing.T) {
	e := &tcp.Error{}
	is.New(t).True(!e.Recovered())
}

func TestErrors_Error(t *testing.T) {
	var err tcp.Errors
	err = append(err, errors.New(hiWorld))
	err = append(err, tcp.NewError(hiWorld))
	are := is.New(t)
	are.Equal(err.Error(), hiWorld+", "+prefix+hiWorld)
}

func TestErrors_Recovered(t *testing.T) {
	var (
		dt = []struct {
			err tcp.Errors
			ok  bool
		}{
			{err: tcp.Errors{}},
			{err: tcp.Errors{tcp.NewError(hiWorld)}, ok: true},
		}
		are = is.New(t)
	)
	for i, tt := range dt {
		t.Run("#"+strconv.Itoa(i), func(t *testing.T) {
			are.Equal(tt.err.Recovered(), tt.ok)
		})
	}
}
