/**
 * Copyright (C) 2014, Markus Kohlhase <mail@markus-kohlhase.de>
 */

package modbus

const (
	ADU_LENGTH    = 260
	HEADER_LENGTH = 7
)

type Header struct {

	// Transaction Identifier
	transaction uint16

	// Protocol Identifier
	protocol uint16

	// PDU Length
	length uint16

	// Unit Identifier
	unit uint8
}

type Adu struct {
	header Header
	pdu    Pdu
}
