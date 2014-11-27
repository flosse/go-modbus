/**
 * Copyright (C) 2014, Markus Kohlhase <mail@markus-kohlhase.de>
 */

package modbus

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
	"strconv"
)

const (
	ADU_LENGTH      = 260
	HEADER_LENGTH   = 7
	TCP_PROTOCOL_ID = 0
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
	buff := bytes.NewBuffer([]byte{})
	binary.Write(buff, binary.BigEndian, h)
	return buff.Bytes()
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

type TcpTransporter struct {
	host        string
	port        uint
	connection  net.Conn
	transaction uint16
	id          uint8
}

func (t *TcpTransporter) Connect() error {
	conn, err := net.Dial("tcp", t.host+":"+strconv.Itoa(int(t.port)))
	if err != nil {
		return err
	}
	t.connection = conn
	return nil
}

func (t *TcpTransporter) Close() error {
	if t.connection != nil {
		return t.connection.Close()
	}
	return errors.New("Not connected")
}

func (t *TcpTransporter) Send(pdu *Pdu) (*Pdu, error) {
	if t.connection == nil {
		if err := t.Connect(); err != nil {
			return nil, err
		}
	}
	t.transaction++
	header := &Header{t.transaction, TCP_PROTOCOL_ID, uint16(len(pdu.data) + 1), t.id}
	adu := &Adu{header, pdu}
	if _, err := t.connection.Write(adu.Pack()); err != nil {
		return nil, errors.New("Could not write data")
	}
	buff := make([]byte, ADU_LENGTH)
	l, err := t.connection.Read(buff)
	if err != nil {
		return nil, errors.New("Could not read data")
	}
	res, err := UnpackAdu(buff[:l])
	if err != nil {
		return nil, errors.New("Could receive PDU")
	}
	return res.pdu, nil
}

func NewTcpClient(host string, port uint) (Client, error) {
	t := &TcpTransporter{host, port, nil, 0, 0}
	return &MbClient{t}, nil
}
