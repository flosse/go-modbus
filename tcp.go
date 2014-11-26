/**
 * Copyright (C) 2014, Markus Kohlhase <mail@markus-kohlhase.de>
 */

package modbus

import (
	"errors"
	"encoding/binary"
)

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
	header *Header
	pdu    *Pdu
}

func (adu *Adu) Pack() []byte {
	return append(adu.header.Pack(), adu.pdu.Pack()...)
}

func (h *Header) Pack() []byte {
	bin := make([]byte,7)
	bin[0] = uint8(h.transaction>>8)
	bin[1] = uint8(h.transaction&0xff)
	bin[2] = uint8(h.protocol>>8)
	bin[3] = uint8(h.protocol&0xff)
	bin[4] = uint8(h.length>>8)
	bin[5] = uint8(h.length&0xff)
	bin[6] = h.unit
	return bin
}

func UnpackHeader(data []byte) (*Header, error) {
	if len(data) < 7 {
		return nil, errors.New("Invalid header length")
	}
	return &Header{
		binary.BigEndian.Uint16(data[0:2]),
		binary.BigEndian.Uint16(data[2:4]),
		binary.BigEndian.Uint16(data[4:6]),
		data[6],
	}, nil
}

func UnpackAdu(data []byte) (*Adu, error) {
	if len(data) < 8 {
		return nil, errors.New("Invalid ADU length")
	}
	pdu, err := UnpackPdu(data[7:])
	if err != nil {
		return nil, err
	}
	header, err := UnpackHeader(data[0:7])
	if err != nil {
		return nil, err
	}
	return &Adu{header, pdu}, nil
}
