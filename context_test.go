package tcp_test

import (
	"context"
	"io"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/rvflash/tcp"
)

func TestContext_Close(t *testing.T) {
	is.New(t).NoErr(newContext(nil).Close())
}

func TestContext_Get(t *testing.T) {
	var (
		dt = []struct {
			ctx    *tcp.Context
			key    string
			value  interface{}
			exists bool
		}{
			{ctx: newContext(newDefaultRequest()), key: unknownKeyName},
			{ctx: newContext(newDefaultRequest()), key: boolKeyName, value: boolKeyValue, exists: true},
			{ctx: newContext(newDefaultRequest()), key: boolBlankKeyName, value: boolBlankKeyValue, exists: true},
			{ctx: newContext(newDefaultRequest()), key: durationKeyName, value: durationKeyValue, exists: true},
			{ctx: newContext(newDefaultRequest()), key: durationBlankKeyName, value: durationBlankKeyValue, exists: true},
			{ctx: newContext(newDefaultRequest()), key: float64KeyName, value: float64KeyValue, exists: true},
			{ctx: newContext(newDefaultRequest()), key: float64BlankKeyName, value: float64BlankKeyValue, exists: true},
			{ctx: newContext(newDefaultRequest()), key: intKeyName, value: intKeyValue, exists: true},
			{ctx: newContext(newDefaultRequest()), key: intBlankKeyName, value: intBlankKeyValue, exists: true},
			{ctx: newContext(newDefaultRequest()), key: int64KeyName, value: int64KeyValue, exists: true},
			{ctx: newContext(newDefaultRequest()), key: int64BlankKeyName, value: int64BlankKeyValue, exists: true},
			{ctx: newContext(newDefaultRequest()), key: stringKeyName, value: stringKeyValue, exists: true},
			{ctx: newContext(newDefaultRequest()), key: stringBlankKeyName, value: stringBlankKeyValue, exists: true},
			{ctx: newContext(newDefaultRequest()), key: unknownKeyName + intKeyName},
			{
				ctx:    newContext(newRequestWithValue(contextPrefix+intKeyName, intKeyValue)),
				key:    contextPrefix + intKeyName,
				value:  intKeyValue,
				exists: true,
			},
		}
		are = is.New(t)
	)
	for i, tt := range dt {
		tt := tt
		t.Run("#"+strconv.Itoa(i), func(t *testing.T) {
			value, exists := tt.ctx.Get(tt.key)
			are.Equal(value, tt.value)
			are.Equal(exists, tt.exists)
		})
	}
}

func TestContext_GetBool(t *testing.T) {
	var (
		dt = []struct {
			ctx   *tcp.Context
			key   string
			value bool
		}{
			{ctx: newContext(newDefaultRequest()), key: unknownKeyName},
			{ctx: newContext(newDefaultRequest()), key: boolKeyName, value: boolKeyValue},
			{ctx: newContext(newDefaultRequest()), key: boolBlankKeyName, value: boolBlankKeyValue},
			{ctx: newContext(newDefaultRequest()), key: unknownKeyName + boolKeyName},
			{
				ctx:   newContext(newRequestWithValue(contextPrefix+boolKeyName, boolKeyValue)),
				key:   contextPrefix + boolKeyName,
				value: boolKeyValue,
			},
		}
		are = is.New(t)
	)
	for i, tt := range dt {
		tt := tt
		t.Run("#"+strconv.Itoa(i), func(t *testing.T) {
			value := tt.ctx.GetBool(tt.key)
			are.Equal(value, tt.value)
		})
	}
}

func TestContext_GetDuration(t *testing.T) {
	var (
		dt = []struct {
			ctx   *tcp.Context
			key   string
			value time.Duration
		}{
			{ctx: newContext(newDefaultRequest()), key: unknownKeyName},
			{ctx: newContext(newDefaultRequest()), key: durationKeyName, value: durationKeyValue},
			{ctx: newContext(newDefaultRequest()), key: durationBlankKeyName, value: durationBlankKeyValue},
			{ctx: newContext(newDefaultRequest()), key: unknownKeyName + durationKeyName},
			{
				ctx:   newContext(newRequestWithValue(contextPrefix+durationKeyName, durationKeyValue)),
				key:   contextPrefix + durationKeyName,
				value: durationKeyValue,
			},
		}
		are = is.New(t)
	)
	for i, tt := range dt {
		tt := tt
		t.Run("#"+strconv.Itoa(i), func(t *testing.T) {
			value := tt.ctx.GetDuration(tt.key)
			are.Equal(value, tt.value)
		})
	}
}

func TestContext_GetInt(t *testing.T) {
	var (
		dt = []struct {
			ctx   *tcp.Context
			key   string
			value int
		}{
			{ctx: newContext(newDefaultRequest()), key: unknownKeyName},
			{ctx: newContext(newDefaultRequest()), key: intKeyName, value: intKeyValue},
			{ctx: newContext(newDefaultRequest()), key: intBlankKeyName, value: intBlankKeyValue},
			{ctx: newContext(newDefaultRequest()), key: unknownKeyName + intKeyName},
			{
				ctx:   newContext(newRequestWithValue(contextPrefix+intKeyName, intKeyValue)),
				key:   contextPrefix + intKeyName,
				value: intKeyValue,
			},
		}
		are = is.New(t)
	)
	for i, tt := range dt {
		tt := tt
		t.Run("#"+strconv.Itoa(i), func(t *testing.T) {
			value := tt.ctx.GetInt(tt.key)
			are.Equal(value, tt.value)
		})
	}
}

func TestContext_GetInt64(t *testing.T) {
	var (
		dt = []struct {
			ctx   *tcp.Context
			key   string
			value int64
		}{
			{ctx: newContext(newDefaultRequest()), key: unknownKeyName},
			{ctx: newContext(newDefaultRequest()), key: int64KeyName, value: int64KeyValue},
			{ctx: newContext(newDefaultRequest()), key: int64BlankKeyName, value: int64BlankKeyValue},
			{ctx: newContext(newDefaultRequest()), key: unknownKeyName + int64KeyName},
			{
				ctx:   newContext(newRequestWithValue(contextPrefix+int64KeyName, int64KeyValue)),
				key:   contextPrefix + int64KeyName,
				value: int64KeyValue,
			},
		}
		are = is.New(t)
	)
	for i, tt := range dt {
		tt := tt
		t.Run("#"+strconv.Itoa(i), func(t *testing.T) {
			value := tt.ctx.GetInt64(tt.key)
			are.Equal(value, tt.value)
		})
	}
}

func TestContext_GetFloat64(t *testing.T) {
	var (
		dt = []struct {
			ctx   *tcp.Context
			key   string
			value float64
		}{
			{ctx: newContext(newDefaultRequest()), key: unknownKeyName},
			{ctx: newContext(newDefaultRequest()), key: float64KeyName, value: float64KeyValue},
			{ctx: newContext(newDefaultRequest()), key: float64BlankKeyName, value: float64BlankKeyValue},
			{ctx: newContext(newDefaultRequest()), key: unknownKeyName + float64KeyName},
			{
				ctx:   newContext(newRequestWithValue(contextPrefix+float64KeyName, float64KeyValue)),
				key:   contextPrefix + float64KeyName,
				value: float64KeyValue,
			},
		}
		are = is.New(t)
	)
	for i, tt := range dt {
		tt := tt
		t.Run("#"+strconv.Itoa(i), func(t *testing.T) {
			value := tt.ctx.GetFloat64(tt.key)
			are.Equal(value, tt.value)
		})
	}
}

func TestContext_GetString(t *testing.T) {
	var (
		dt = []struct {
			ctx *tcp.Context
			key,
			value string
		}{
			{ctx: newContext(newDefaultRequest()), key: unknownKeyName},
			{ctx: newContext(newDefaultRequest()), key: stringKeyName, value: stringKeyValue},
			{ctx: newContext(newDefaultRequest()), key: stringBlankKeyName, value: stringBlankKeyValue},
			{ctx: newContext(newDefaultRequest()), key: unknownKeyName + stringKeyName},
			{
				ctx:   newContext(newRequestWithValue(contextPrefix+stringKeyName, stringKeyValue)),
				key:   contextPrefix + stringKeyName,
				value: stringKeyValue,
			},
		}
		are = is.New(t)
	)
	for i, tt := range dt {
		tt := tt
		t.Run("#"+strconv.Itoa(i), func(t *testing.T) {
			value := tt.ctx.GetString(tt.key)
			are.Equal(value, tt.value)
		})
	}
}

func TestContext_ReadAll(t *testing.T) {
	const data01 = "hello world"
	var (
		dt = []struct {
			req *tcp.Request
			out []byte
			err error
		}{
			{err: tcp.ErrRequest},
			{req: tcp.NewRequest(tcp.SYN, nil), err: io.EOF},
			{req: tcp.NewRequest(tcp.ACK, strings.NewReader(data01)), out: []byte(data01)},
		}
		are = is.New(t)
	)
	for i, tt := range dt {
		tt := tt
		t.Run("#"+strconv.Itoa(i), func(t *testing.T) {
			c := newContext(tt.req)
			out, err := c.ReadAll()
			are.Equal(err, tt.err)
			are.Equal(out, tt.out)
		})
	}
}

func TestContext_Canceled(t *testing.T) {
	const timeOut = time.Millisecond * 50
	var (
		ctx = context.Background()
		req = newDefaultRequest()
	)
	// Without request
	c := newContext(nil)
	select {
	case <-c.Canceled():
		t.Fatal("unexpected cancellation")
	case <-time.After(timeOut):
		// next!
	}
	// Without cancellation
	c = newContext(req.WithContext(ctx))
	select {
	case <-c.Canceled():
		t.Fatal("unexpected cancellation")
	case <-time.After(timeOut):
		// next!
	}
	// With cancellation
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	c = newContext(req.WithContext(ctx))
	select {
	case <-c.Canceled():
		// well done
	case <-time.After(timeOut):
		t.Fatal("expected cancellation")
	}
}

const (
	// names of shared data
	contextPrefix        = "ctx_"
	unknownKeyName       = "unknown"
	boolKeyName          = "boolean_ok"
	boolBlankKeyName     = "boolean_bk"
	durationKeyName      = "duration_ok"
	durationBlankKeyName = "duration_bk"
	float64KeyName       = "float64_ok"
	float64BlankKeyName  = "float64_bk"
	intKeyName           = "int_ok"
	intBlankKeyName      = "int_bk"
	int64KeyName         = "int64_ok"
	int64BlankKeyName    = "int64_bk"
	stringKeyName        = "string_ok"
	stringBlankKeyName   = "string_bk"

	// values of shared data
	boolKeyValue          bool          = true
	boolBlankKeyValue     bool          = false
	durationKeyValue                    = time.Second
	durationBlankKeyValue time.Duration = 0
	float64KeyValue       float64       = 3.14
	float64BlankKeyValue  float64       = 0
	intKeyValue           int           = 314
	intBlankKeyValue      int           = 0
	int64KeyValue         int64         = 314
	int64BlankKeyValue    int64         = 0
	stringKeyValue        string        = "hello"
	stringBlankKeyValue   string        = ""
)

func newContext(req *tcp.Request) *tcp.Context {
	return &tcp.Context{
		Request: req,
		Shared: tcp.M{
			boolKeyName:          boolKeyValue,
			boolBlankKeyName:     boolBlankKeyValue,
			durationKeyName:      durationKeyValue,
			durationBlankKeyName: durationBlankKeyValue,
			float64KeyName:       float64KeyValue,
			float64BlankKeyName:  float64BlankKeyValue,
			intKeyName:           intKeyValue,
			intBlankKeyName:      intBlankKeyValue,
			int64KeyName:         int64KeyValue,
			int64BlankKeyName:    int64BlankKeyValue,
			stringKeyName:        stringKeyValue,
			stringBlankKeyName:   stringBlankKeyValue,
		},
	}
}

func newDefaultRequest() *tcp.Request {
	return tcp.NewRequest(tcp.ACK, nil)
}

func newRequestWithValue(key, val interface{}) *tcp.Request {
	req := newDefaultRequest()
	return req.WithContext(context.WithValue(req.Context(), key, val))
}
