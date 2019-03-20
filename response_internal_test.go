package tcp

import (
	"testing"

	"github.com/matryer/is"
)

func TestResponseWriter_Write(t *testing.T) {
	var (
		are = is.New(t)
		dt  = []struct {
			msg       string
			len, size int
		}{
			{msg: "hi", len: 2, size: 2},
			{msg: " ", len: 1, size: 3},
			{msg: "world!", len: 6, size: 9},
			{msg: "\n", len: 1, size: 10},
		}
		w   = newWriter(NewRecorder())
		n   int
		err error
	)
	// buffer not used yet
	are.Equal(w.Size(), noWritten)
	// writes into
	for _, tt := range dt {
		tt := tt
		t.Run(tt.msg, func(t *testing.T) {
			n, err = w.Write([]byte(tt.msg))
			are.NoErr(err)
			are.Equal(n, tt.len)         // len mismatch
			are.Equal(w.Size(), tt.size) // size mismatch
		})
	}
	// closes it.
	are.NoErr(w.Close())
}

func TestResponseWriter_WriteString(t *testing.T) {
	var (
		are = is.New(t)
		dt  = []struct {
			msg       string
			len, size int
		}{
			{msg: "hi", len: 2, size: 2},
			{msg: " ", len: 1, size: 3},
			{msg: "world!", len: 6, size: 9},
			{msg: "\n", len: 1, size: 10},
		}
		w   = newWriter(NewRecorder())
		n   int
		err error
	)
	// buffer not used yet
	are.Equal(w.Size(), noWritten)
	// writes into
	for _, tt := range dt {
		tt := tt
		t.Run(tt.msg, func(t *testing.T) {
			n, err = w.WriteString(tt.msg)
			are.NoErr(err)
			are.Equal(n, tt.len)         // len mismatch
			are.Equal(w.Size(), tt.size) // size mismatch
		})
	}
	// closes it.
	are.NoErr(w.Close())
}
