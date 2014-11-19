/**
 * Copyright (C) 2014, Markus Kohlhase <mail@markus-kohlhase.de>
 *
 * Modbus API
 */

package modbus

type Client interface {

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
	ReadInputRegisters(address, quantity uint16) (readRegisters []int16, err error)

	// Function Code 3
	ReadHoldingRegisters(address, quantity uint16) (readRegisters []int16, err error)

	// Function Code 6
	WriteSingleRegister(address uint16, value int16) error

	// Function Code 16
	WriteMultipleRegisters(address, quantity uint16, values []int16) error

	// Function Code 23
	ReadWriteMultipleRegisters(readAddress, readQuantity, writeAddress, writeQuantity uint16, values []int16) (readRegisters []int16, err error)

	// Function Code 22
	MaskWriteRegister(address uint16, andMask, orMask int16) error

	// Function Code 24
	ReadFifoQueue(address uint16) (count uint16, values []int16, err error)

	/**********************
	 * File record access *
	 **********************/

	// TODO: specify methods

	/***************
	 * Diagnostics *
	 **************/

	// TODO: specify methods

}

type Handler interface {
	Handle(req *Pdu) (res *Pdu)
}

type Server interface {
	SetHandler(h *Handler)
	Start() error
	Stop() error
}
