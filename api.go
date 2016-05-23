/**
 * Copyright (C) 2014 - 2015, Markus Kohlhase <mail@markus-kohlhase.de>
 *
 * Modbus API
 */

package modbus

type Transporter interface {
	Connect() error
	Close() error
	Send(pdu *Pdu) (*Pdu, error)
}

type Client interface {

	/* Access to transporter layer */

	Transporter() Transporter

	/**************
	 * Bit access *
	 **************/

	/* Physical Discrete Inputs */

	// Function Code 2
	ReadDiscreteInputs(address, quantity uint16) (inputs []bool, err error)

	/* Internal Bits Or Physical Coils */

	// Function Code 1
	ReadCoils(address, quantity uint16) (coils []bool, err error)

	// Function Code 5
	WriteSingleCoil(address uint16, coil bool) error

	// Function Code 15
	WriteMultipleCoils(address uint16, coils []bool) error

	/******************
	 * 16 bits access *
	 ******************/

	/* Physical Input Registers */

	// Function Code 4
	ReadInputRegisters(address, quantity uint16) (readRegisters []uint16, err error)

	// Function Code 3
	ReadHoldingRegisters(address, quantity uint16) (readRegisters []uint16, err error)

	// Function Code 6
	WriteSingleRegister(address, value uint16) error

	// Function Code 16
	WriteMultipleRegisters(address uint16, values []uint16) error

	// Function Code 23
	ReadWriteMultipleRegisters(readAddress, readQuantity, writeAddress uint16, values []uint16) (readRegisters []uint16, err error)

	// Function Code 22
	MaskWriteRegister(address, andMask, orMask uint16) error

	// Function Code 24
	ReadFifoQueue(address uint16) (values []uint16, err error)

	/**********************
	 * File record access *
	 **********************/

	// TODO: specify methods

}

type SerialClient interface {

	// Embed general client API
	Client

	/***************
	 * Diagnostics *
	 **************/

	// Function Code 7
	ReadExceptionStatus() (states []bool, err error)

	// Function Code 8
	Diagnostics(subfunction uint16, data []uint16) (response []uint16, err error)

	// Function Code 11
	GetCommEventCounter() (status bool, count uint16, err error)

	// TODO: specify method
	// Function Code 12
	// GetCommEventLog

	// Function Code 17
	ReportServerId() (response []byte, err error)

	// TODO: specify method
	// Function Code 43
	// ReadDeviceIdentification

}

type IoClient interface {

	/********************
	 * Abstract Objects *
	 ********************/

	// Embed general client API
	Client

	// Discrete input
	DiscreteInput(address uint16) DiscreteInput

	// Coil
	Coil(address uint16) Coil

	// Input Register
	InputRegister(address uint16) InputRegister

	// Input Registers
	InputRegisters(address, count uint16) InputRegisters

	// Holding Register
	HoldingRegister(address uint16) HoldingRegister

	// Holding Registers
	HoldingRegisters(address, count uint16) HoldingRegisters
}

type DiscreteInput interface {
	Test() (bool, error)
}

type Coil interface {
	DiscreteInput
	Set() error
	Clear() error
	Toggle() error
}

type InputRegister interface {
	Read() (uint16, error)
}

type InputRegisters interface {
	Read() ([]uint16, error)
	ReadString() (string, error)
}

type HoldingRegister interface {
	InputRegister
	Write(uint16) error
}

type HoldingRegisters interface {
	InputRegisters
	Write([]uint16) error
	WriteString(s string) error
}

type Handler interface {
	Handle(req *Pdu) (res *Pdu)
}

type Server interface {
	SetHandler(h *Handler)
	Start() error
	Stop() error
}
