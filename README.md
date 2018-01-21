# go-modbus

[![Build Status](https://travis-ci.org/flosse/go-modbus.svg?branch=master)](https://travis-ci.org/flosse/go-modbus)
[![GoDoc](https://godoc.org/github.com/flosse/go-modbus?status.svg)](https://godoc.org/github.com/flosse/go-modbus)
[![Coverage Status](https://coveralls.io/repos/flosse/go-modbus/badge.svg?branch=master)](https://coveralls.io/r/flosse/go-modbus?branch=master)

**DONT USE IT**: This was a proof-of-concept!

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
roRegister := master.InputRegister(0x2000)
value, err := roRegister.Read()

/* multiple ro registers */
multRoRegister := master.InputRegisters(0x1000,7)
values, err    := roRegister.Read()
myString, err  := multRoRegister.ReadString()

/* read-write */
register   := master.HoldingRegister(0x0900)
value, err := register.Read()
err = register.Write(0x435)

/* multiple rw registers */
multRwRegisters := master.HoldingRegisters(0x9000, 3)
values, err     := multRwRegisters.Read()
aString, err    := multRwRegisters.ReadString()
err := multRwRegisters.Write(uint16{3,2,1})
err =  multRwRegisters.WriteString("foo")
```

#### Low Level API

```go
/* Bit access */

// read three read-only bits
res, err := master.ReadDiscreteInputs(0x0800, 3)
// res could be [true, false, false]

// read 5 read-write bits
res, err = master.ReadCoils(0x02, 2)
// res could be [false, true]

// set the coil at address 0x0734
err = master.WriteSingleCoil(0x734, true)

// set/clear multiple coils at address 0x0002
err = master.WriteMultipleCoils(2, []bool{false, true, true})

/* 16 bits access */

// read three read-only registers
res, err = master.ReadInputRegisters(0x12, 3)
// res could be [334, 912, 0]

// read two read-write registers
res, err = master.ReadHoldingRegisters(0x00, 2)
// res could be [9, 42]

// write a value to a single register
err = master.WriteSingleRegister(0x07, 9923)

// write values to multiple registers
err = master.WriteMultipleRegisters(0x03, []uint16{9,0,66})

// read two and write three values within one transaction
res, err = master.ReadWriteMultipleRegisters(0x0065, 2, 0x0800, []uint16{0,7,33})
// res could be [0, 88]
```

## Run Tests

    go get github.com/smartystreets/goconvey
    go test

or run

    $GOPATH/bin/goconvey

and open `http://localhost:8080` in your browser

## License

This library is licensed under the MIT license

## Credits

This library is inspired by [this modbus library](https://github.com/goburrow/modbus).
