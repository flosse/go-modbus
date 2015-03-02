/**
 * Copyright (C) 2014 - 2015, Markus Kohlhase <mail@markus-kohlhase.de>
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
	aduLength     = 260
	headerLength  = 7
	tcpProtocolId = 0
)

type header struct {

	// Transaction Identifier
	transaction uint16

	// Protocol Identifier
	protocol uint16

	// PDU Length
	length uint16

	// Unit Identifier
	unit uint8
}

type adu struct {
	header *header
	pdu    *Pdu
}

func (adu *adu) pack() []byte {
	return append(adu.header.pack(), adu.pdu.pack()...)
}

func (h *header) pack() []byte {
	buff := bytes.NewBuffer([]byte{})
	binary.Write(buff, binary.BigEndian, h)
	return buff.Bytes()
}

func unpackHeader(data []byte) (*header, error) {
	if len(data) < headerLength {
		return nil, errors.New("Invalid header length")
	}
	return &header{
		binary.BigEndian.Uint16(data[0:2]),
		binary.BigEndian.Uint16(data[2:4]),
		binary.BigEndian.Uint16(data[4:6]),
		data[6],
	}, nil
}

func unpackAdu(data []byte) (*adu, error) {
	if len(data) < 8 {
		return nil, errors.New("Invalid ADU length")
	}
	pdu, err := unpackPdu(data[headerLength:])
	if err != nil {
		return nil, err
	}
	header, err := unpackHeader(data[0:headerLength])
	if err != nil {
		return nil, err
	}
	return &adu{header, pdu}, nil
}

type tcpTransporter struct {
	host        string
	port        uint
	connection  net.Conn
	transaction uint16
	id          uint8
}

func (t *tcpTransporter) Connect() error {
	conn, err := net.Dial("tcp", t.host+":"+strconv.Itoa(int(t.port)))
	if err != nil {
		return err
	}
	t.connection = conn
	return nil
}

func (t *tcpTransporter) Close() error {
	if t.connection != nil {
		return t.connection.Close()
	}
	return errors.New("Not connected")
}

func (t *tcpTransporter) Send(pdu *Pdu) (*Pdu, error) {
	if t.connection == nil {
		if err := t.Connect(); err != nil {
			return nil, err
		}
	}
	t.transaction++
	header := &header{t.transaction, tcpProtocolId, uint16(len(pdu.Data) + 2), t.id}
	adu := &adu{header, pdu}
	if _, err := t.connection.Write(adu.pack()); err != nil {
		return nil, errors.New("Could not write data")
	}
	buff := make([]byte, aduLength)
	l, err := t.connection.Read(buff)
	if err != nil {
		return nil, errors.New("Could not receive data")
	}
	res, err := unpackAdu(buff[:l])
	if err != nil {
		return nil, errors.New("Could not read PDU")
	}
	return res.pdu, nil
}

func NewTcpClient(host string, port uint) (Client, error) {
	t := &tcpTransporter{host, port, nil, 0, 0}
	return &mbClient{t}, nil
}
