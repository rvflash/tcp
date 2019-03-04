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
	// LogRemoteAddr is the name of the log's field for the remote address.
	LogRemoteAddr = "addr"
	// LogRequestSize is the name of the log's field for the request size.
	LogRequestSize = "req_size"
	// LogResponseSize is the name of the log's field for the response size.
	LogResponseSize = "resp_size"
	// LogLatency is the name of the log's field with the response duration.
	LogLatency = "latency"
	// LogServerHostname is the name of the log's field with the server hostname.
	LogServerHostname = "server"
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
		err := c.Err()
		switch {
		case err == nil:
			entry.Info(m.String())
		case err.Recovered():
			entry.Errorf("%s %s", m, err)
		default:
			entry.Warnf("%s %s", m, err)
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
	for k, v := range f {
		switch k {
		case LogRemoteAddr:
			d[k] = m.req.RemoteAddr
		case LogRequestSize:
			d[k] = m.reqSize
		case LogResponseSize:
			d[k] = w.Size()
		case LogLatency:
			m.latency = time.Since(m.start)
			d[k] = int(math.Ceil(float64(m.latency.Nanoseconds()) / float64(time.Millisecond)))
		case LogServerHostname:
			d[k], _ = os.Hostname()
		default:
			// allows to logs statics data
			d[k] = v
		}
	}
	return d
}

// String implements the fmt.Stringer interface.
func (m *message) String() string {
	if m.req == nil {
		// unexpected segment
		return ""
	}
	return "[TCP] " + m.start.Format(time.RFC3339) + " | " + m.req.Segment
}
