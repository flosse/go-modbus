/**
 * Copyright (C) 2014, Markus Kohlhase <mail@markus-kohlhase.de>
 */

package modbus

import (
	"errors"
)

const PDU_LENGTH = 253

type Pdu struct {

	// Function Code
	function uint8

	// PDU data
	data []byte
}

func (pdu *Pdu) Pack() []byte {
	buff := make([]byte, 1, PDU_LENGTH)
	buff[0] = pdu.function
	return append(buff, pdu.data...)
}

func UnpackPdu(data []byte) (*Pdu, error) {
	if len(data) < 1 {
		return nil, errors.New("Invalid PDU length")
	}
	return &Pdu{data[0], data[1:]}, nil
}
