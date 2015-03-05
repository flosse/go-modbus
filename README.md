# go-modbus

[![Build Status](https://travis-ci.org/flosse/go-modbus.svg?branch=master)](https://travis-ci.org/flosse/go-modbus)
[![GoDoc](https://godoc.org/github.com/flosse/go-modbus?status.svg)](https://godoc.org/github.com/flosse/go-modbus)

a free [Modbus](http://en.wikipedia.org/wiki/Modbus) library
for [Go](http://golang.org/).

## Usage

### Modbus Master (Client)

```go
package main

import (
  "fmt"
  "github.com/flosse/go-modbus"
)

func main(){
  master := modbus.NewTcpClient("127.0.0.1", 502)
}
```

#### High Level API

```go
/* read-only */
di := master.DiscreteInput(7)
state, err := di.Test()

/* read-write */
coil := master.Coil(2)
state, err = coil.Test()
err = coil.Set()
err = coil.Clear()
err = coil.Toggle()

/* read-only */
roRegister := master.InputRegister(5)
value, err = roRegister.Read()

/* read-write */
register := master.HoldingRegister(9)
value, err = register.Read()
err = register.Write(0x435)
```

## Run Tests

    go get github.com/smartystreets/goconvey
    go test

or run

    $GOPATH/bin/goconvey

and open `http://localhost:8080` in your browser

## License

This library is licensed under the MIT license
