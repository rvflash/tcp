# TCP

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

import "github.com/rvflash/tcp"

func main() {
	r := tcp.Default()
	r.ACK(func(c tcp.Conn) {
		// new message received
		b := c.RawData()
		// ...
	})
	_ = r.Run(":9090") // listen and serve on 0.0.0.0:9090
}
```
