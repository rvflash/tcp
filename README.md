# TCP

[![GoDoc](https://godoc.org/github.com/rvflash/tcp?status.svg)](https://godoc.org/github.com/rvflash/tcp)
[![Build Status](https://img.shields.io/travis/rvflash/tcp.svg)](https://travis-ci.org/rvflash/tcp)
[![Code Coverage](https://img.shields.io/codecov/c/github/rvflash/tcp.svg)](http://codecov.io/github/rvflash/tcp?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/rvflash/tcp)](https://goreportcard.com/report/github.com/rvflash/tcp)

The package `tcp` provides interfaces to create a TCP server.


### Installation
    
To install it, you need to install Go and set your Go workspace first.
Then, download and install it:

```bash
$ go get -u github.com/rvflash/tcp
```    
Import it in your code:
    
```go
import "github.com/rvflash/tcp"
```

### Prerequisite

`tcp` uses the Go modules that required Go 1.11 or later.


## Features

### TLS support

By using the `RunTLS` method instead of `Run`, you can specify a certificate and
a X509 key to create an TCP/TLS connection.


### Handler

Just as Gin, a well done web framework whose provides functions based on HTTP methods,
`tcp` provides functions based on TCP segments.

Thus, it exposes a method for each of its segments:
* `ACK` to handle each new message. 
* `FIN` to handle when the connection is closed.
* `SYN` to handle each new connection.

More functions are available, see the [godoc](https://godoc.org/github.com/rvflash/tcp) for more details.

Each of these methods take as parameter the HandlerFunc interface: `func(c *Context)`.
You must implement this interface to create your own handler.


> By analogy with the awesome standard HTTP package, `tcp` exposes and implements
the Handler interface `ServeTCP(ResponseWriter, *Request)`.


### Middleware

By using the `Default` method instead of the `New` to initiate a TCP server,
2 middlewares are defined on each segment.
The first allows to recover on panic, and the second enables logs.
 

### Custom Middleware

The `Next` method on the `Context` should only be used inside middleware. Its allows to pass to the pending handlers. 
See the `Recovery` or `Logger` methods as sample code.


## Quick start

Assuming the following code that runs a server on port 9090:

```go
package main

import (
	"log"

	"github.com/rvflash/tcp"
)

func main() {
	// creates a server with a logger and a recover on panic as middlewares.
	r := tcp.Default()
	r.ACK(func(c *tcp.Context) {
		// new message received
		// gets the request body
		buf, err := c.ReadAll()
		if err != nil {
			c.Error(err)
			return
		}
		// writes something as response
		c.String(string(buf))
	})
	err := r.Run(":9090") // listen and serve on 0.0.0.0:9090
	if err != nil {
		log.Fatalf("listen: %s", err)
	}
}
```