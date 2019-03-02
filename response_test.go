package tcp_test

import (
	"testing"

	"github.com/matryer/is"
	"github.com/rvflash/tcp"
)

const (
	msg     = "hello world\n"
	msgSize = 12
)

func TestResponseRecorder_Close(t *testing.T) {
	is.New(t).True(tcp.NewRecorder().Close() == nil)
}

func TestResponseRecorder_Write(t *testing.T) {
	w := tcp.NewRecorder()
	n, err := w.Write([]byte(msg))
	are := is.New(t)
	are.NoErr(err)
	are.Equal(n, msgSize)
	are.Equal(w.Size(), msgSize)
}

func TestResponseRecorder_WriteString(t *testing.T) {
	w := tcp.NewRecorder()
	n, err := w.WriteString(msg)
	are := is.New(t)
	are.NoErr(err)
	are.Equal(n, msgSize)
	are.Equal(w.Size(), msgSize)
}

func TestResponseRecorder_Size(t *testing.T) {
	are := is.New(t)
	// no body
	w := &tcp.ResponseRecorder{}
	are.Equal(w.Size(), -1)
	// with a valid body
	w = tcp.NewRecorder()
	are.Equal(w.Size(), 0)

}
