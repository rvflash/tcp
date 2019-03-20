package tcp_test

import (
	"errors"
	"io"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/matryer/is"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"

	"github.com/rvflash/tcp"
)

func TestLogger(t *testing.T) {
	var (
		are = is.New(t)

		// default values
		vField        = "version"
		vValue        = "v0.0.0"
		defaultFields = logrus.Fields{
			tcp.LogLatency:        0,
			tcp.LogServerHostname: "",
			tcp.LogRemoteAddr:     "",
			tcp.LogRequestSize:    0,
			tcp.LogResponseSize:   0,
		}
		customFields = logrus.Fields{
			tcp.LogLatency:     0,
			tcp.LogRequestSize: 0,
			vField:             vValue,
		}

		// test cases
		dt = []struct {
			body io.Reader
			bodySize,
			dataSize int
			fields      logrus.Fields
			handler     []tcp.HandlerFunc
			level       logrus.Level
			minDuration int
		}{
			{
				body:     strings.NewReader("hello world"),
				bodySize: 11,
				dataSize: 5,
				fields:   defaultFields,
				level:    logrus.InfoLevel,
			},
			{
				dataSize: 5,
				fields:   defaultFields,
				handler:  []tcp.HandlerFunc{stumble},
				level:    logrus.WarnLevel,
			},
			{
				dataSize:    5,
				fields:      defaultFields,
				handler:     []tcp.HandlerFunc{sleep},
				level:       logrus.InfoLevel,
				minDuration: 100,
			},
			{
				dataSize:    3,
				fields:      customFields,
				handler:     []tcp.HandlerFunc{stumble, sleep},
				level:       logrus.WarnLevel,
				minDuration: 100,
			},
			{
				dataSize: 3,
				fields:   customFields,
				handler:  []tcp.HandlerFunc{oops},
				level:    logrus.ErrorLevel,
			},
		}

		entry *logrus.Entry
		v     interface{}
		ok    bool
	)

	log, hook := test.NewNullLogger()
	log.Formatter = &logrus.TextFormatter{DisableTimestamp: true}

	for i, tt := range dt {
		tt := tt
		t.Run("#"+strconv.Itoa(i), func(t *testing.T) {
			// launches the server
			srv := tcp.New()
			srv.Use(tcp.Logger(log, tt.fields))
			srv.Use(tcp.Recovery())
			srv.SYN(tt.handler...)
			// serves the request
			srv.ServeTCP(tcp.NewRecorder(), tcp.NewRequest(tcp.SYN, tt.body))
			// checks the log's message
			entry = hook.LastEntry()
			are.Equal(entry.Level, tt.level)                            // level mismatch
			are.Equal(len(entry.Data), tt.dataSize)                     // fields size mismatch
			are.True(entry.Data[tcp.LogLatency].(int) > tt.minDuration) // min duration
			if _, ok = tt.fields[vField]; ok {
				v, ok = entry.Data[vField]
				are.True(ok)         // version required
				are.Equal(v, vValue) // version mismatch
			}
			are.Equal(entry.Data[tcp.LogRequestSize].(int), tt.bodySize) //  request size mismatch
		})
	}
}

func sleep(_ *tcp.Context) {
	time.Sleep(time.Millisecond * 100)
}

func stumble(c *tcp.Context) {
	c.Error(errors.New("my bad, sorry"))
}
