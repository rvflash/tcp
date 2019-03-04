package tcp

import (
	"strings"
	"testing"
	"time"

	"github.com/matryer/is"
)

func TestMessage_String(t *testing.T) {
	const blog = "[TCP] 0001-01-01T00:00:00Z | ACK"
	msg := &message{}
	req := NewRequest(ACK, nil)
	are := is.New(t)
	are.True(msg.String() == "") // invalid request
	msg.req = req
	are.Equal(msg.String(), blog) // invalid starting date
	msg = newMessage(req)
	are.True(strings.HasSuffix(msg.String(), ACK))
	are.True(strings.Contains(msg.String(), time.Now().UTC().Format("2006-01-02T15:04")))
}
