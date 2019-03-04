# TCP

[![GoDoc](https://godoc.org/github.com/rvflash/tcp?status.svg)](https://godoc.org/github.com/rvflash/tcp)
[![Build Status](https://img.shields.io/travis/rvflash/tcp.svg)](https://travis-ci.org/rvflash/tcp)
[![Code Coverage](https://img.shields.io/codecov/c/github/rvflash/tcp.svg)](http://codecov.io/github/rvflash/tcp?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/rvflash/tcp)](https://goreportcard.com/report/github.com/rvflash/tcp)

TCP provides interfaces to create a TCP server.


### Installation
    
To install it, you need to install Go and set your Go workspace first.
Download and install it:

```bash
$ go get -u github.com/rvflash/tcp
```    
Import it in your code:
    
```go
import "github.com/rvflash/tcp"
```

### Prerequisite

`tcp` uses the Go modules that required Go 1.11 or later.


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
