/**
 * Copyright (C) 2014 - 2015, Markus Kohlhase <mail@markus-kohlhase.de>
 */

package modbus

import (
	"fmt"
)

type Pdu struct {

	// Function Code
	Function uint8

	// PDU data
	Data []byte
}

const pduLength = 253

func (pdu *Pdu) pack() (bin []byte, err error) {
	if pdu.Function < 1 {
		return nil, fmt.Errorf("Invalid function code %d", pdu.Function)
	}
	if l := len(pdu.Data); l > pduLength-1 {
		return nil, fmt.Errorf("Invalid length of data (%d instead of max. %d bytes)", l, pduLength-1)
	}
	return append([]byte{pdu.Function}, pdu.Data...), nil
}

func unpackPdu(data []byte) (*Pdu, error) {
	if l := len(data); l < 1 {
		return nil, fmt.Errorf("Invalid PDU length (%d bytes)", l)
	}
	return &Pdu{data[0], data[1:]}, nil
}
