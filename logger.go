package tcp

import (
	"bytes"
	"io/ioutil"
	"math"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	remoteAddr = "addr"
	reqLength  = "req_size"
	respLength = "resp_size"
	latency    = "latency"
	hostname   = "server"
)

// Logger returns a middleware to log each TCP request.
func Logger(log *logrus.Logger, fields logrus.Fields) HandlerFunc {
	return func(c *Context) {
		// Initiates the timer
		m := newMessage(c.Request)
		// Processes the request
		c.Next()
		// Logs it.
		entry := logrus.NewEntry(log).WithFields(m.fields(c.ResponseWriter, fields))
		if e := c.Err(); e == nil {
			entry.Info(m.String())
		} else if e.Recovered() {
			entry.Errorf("%s %s", m, e)
		} else {
			entry.Warnf("%s %s", m, e)
		}
	}
}

func newMessage(req *Request) *message {
	// starts the UTC timer.
	m := &message{
		start: time.Now().UTC(),
		req:   req,
	}
	// reads the request body without closing it to get its size.
	if req.Body != nil {
		buf, _ := ioutil.ReadAll(req.Body)
		req.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
		m.reqSize = len(buf)
	}
	return m
}

type message struct {
	latency time.Duration
	req     *Request
	reqSize int
	start   time.Time
}

func (m *message) fields(w ResponseWriter, f logrus.Fields) logrus.Fields {
	d := make(logrus.Fields)
	for k := range f {
		switch k {
		case remoteAddr:
			d[k] = m.req.RemoteAddr
		case reqLength:
			d[k] = w.Size()
		case respLength:
			d[k] = m.reqSize
		case latency:
			m.latency = time.Since(m.start)
			d[k] = int(math.Ceil(float64(m.latency.Nanoseconds()) / 1000.0))
		case hostname:
			d[k], _ = os.Hostname()
		}
	}
	return d
}

// String implements the fmt.Stringer interface.
func (m *message) String() string {
	return "[TCP] " + m.start.Format(time.RFC3339) + " | " + m.req.Segment
}
