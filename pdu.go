/**
 * Copyright (C) 2014 - 2015, Markus Kohlhase <mail@markus-kohlhase.de>
 */

package modbus

import (
	"errors"
)

type Pdu struct {

	// Function Code
	Function uint8

	// PDU data
	Data []byte
}

const pduLength = 253

func (pdu *Pdu) pack() []byte {
	buff := make([]byte, 1, pduLength)
	buff[0] = pdu.Function
	return append(buff, pdu.Data...)
}

func unpackPdu(data []byte) (*Pdu, error) {
	if len(data) < 1 {
		return nil, errors.New("Invalid PDU length")
	}
	return &Pdu{data[0], data[1:]}, nil
}
