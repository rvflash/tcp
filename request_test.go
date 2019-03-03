package tcp_test

import (
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/rvflash/tcp"
)

func TestNewRequest(t *testing.T) {
	var (
		are = is.New(t)
		dt  = []struct {
			seg     string
			body    io.Reader
			segment string
			resp    []byte
		}{
			{segment: tcp.SYN},
			{seg: tcp.ACK, segment: tcp.ACK},
			{seg: tcp.SYN, segment: tcp.SYN},
			{seg: tcp.FIN, segment: tcp.FIN},
			{seg: "NOP", segment: "NOP"},
			{seg: tcp.ACK, segment: tcp.ACK, body: strings.NewReader(msg), resp: []byte(msg)},
		}
		req *tcp.Request
		b   []byte
		err error
	)
	for i, tt := range dt {
		t.Run("#"+strconv.Itoa(i), func(t *testing.T) {
			req = tcp.NewRequest(tt.seg, tt.body)
			are.Equal(req.Segment, tt.segment)
			if tt.body == nil {
				are.True(req.Body == nil)
			} else {
				b, err = ioutil.ReadAll(req.Body)
				are.NoErr(err)
				are.Equal(b, tt.resp)
			}
		})
	}
}

func TestRequest_Context(t *testing.T) {
	is.New(t).True(tcp.NewRequest(tcp.SYN, nil).Context() != nil)
}
