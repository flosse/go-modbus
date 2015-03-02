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

	/********************
	 * Abstract Objects *
	 ********************/

	// Discrete input
	DiscreteInput(address uint16) DiscreteInput

	// Coil
	Coil(address uint16) Coil

	// Input Register
	InputRegister(address uint16) InputRegister

	// Holding Register
	HoldingRegister(address uint16) HoldingRegister

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

	/***************
	 * Diagnostics *
	 **************/

	// TODO: specify methods

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

type HoldingRegister interface {
	InputRegister
	Write(uint16) error
}

type Handler interface {
	Handle(req *Pdu) (res *Pdu)
}

type Server interface {
	SetHandler(h *Handler)
	Start() error
	Stop() error
}
