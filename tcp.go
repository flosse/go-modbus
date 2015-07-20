/**
 * Copyright (C) 2014 - 2015, Markus Kohlhase <mail@markus-kohlhase.de>
 */

package modbus

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"
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

func (adu *adu) pack() (bin []byte, err error) {
	binPdu, err := adu.pdu.pack()
	if err != nil {
		return
	}
	bin = append(adu.header.pack(), binPdu...)
	return
}

func (h *header) pack() []byte {
	buff := bytes.NewBuffer([]byte{})
	binary.Write(buff, binary.BigEndian, h)
	return buff.Bytes()
}

func unpackHeader(data []byte) (*header, error) {
	if l := len(data); l < headerLength {
		return nil, fmt.Errorf("Invalid header length: %d byte", l)
	}
	return &header{
		binary.BigEndian.Uint16(data[0:2]),
		binary.BigEndian.Uint16(data[2:4]),
		binary.BigEndian.Uint16(data[4:6]),
		data[6],
	}, nil
}

func unpackAdu(data []byte) (*adu, error) {
	if l := len(data); l < 8 {
		return nil, fmt.Errorf("Invalid ADU length: %d byte", l)
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
	conn        net.Conn
	transaction uint16
	id          uint8
	timeout     time.Duration
}

func (t *tcpTransporter) Connect() error {
	address := t.host + ":" + strconv.Itoa(int(t.port))
	if t.timeout > 0 {
		conn, err := net.DialTimeout("tcp", address, t.timeout)
		t.conn = conn
		return err
	} else {
		conn, err := net.Dial("tcp", address)
		t.conn = conn
		return err
	}
}

func (t *tcpTransporter) Close() (err error) {
	if t.conn != nil {
		if err = t.conn.Close(); err != nil {
			return
		}
		t.conn = nil
		return
	}
	return errors.New("Not connected")
}

func (t *tcpTransporter) Send(pdu *Pdu) (*Pdu, error) {
	if t.conn == nil {
		if err := t.Connect(); err != nil {
			return nil, err
		}
	}
	t.transaction++
	header := &header{t.transaction, tcpProtocolId, uint16(len(pdu.Data) + 2), t.id}
	binAdu, err := (&adu{header, pdu}).pack()
	if err != nil {
		return nil, err
	}
	if _, err := t.conn.Write(binAdu); err != nil {
		return nil, fmt.Errorf("Could not write data: %s", err)
	}
	buff := make([]byte, aduLength)
	l, err := t.conn.Read(buff)
	if err != nil {
		return nil, fmt.Errorf("Could not receive data: %s", err)
	}
	res, err := unpackAdu(buff[:l])
	if err != nil {
		return nil, fmt.Errorf("Could not read PDU: %s", err)
	}
	if i := res.header.transaction; i != t.transaction {
		return nil, fmt.Errorf("Invalid transaction id: %d instead of %d", i, t.transaction)
	}
	return res.pdu, nil
}

func NewTcpClient(host string, port uint) IoClient {
	return &mbClient{&tcpTransporter{host: host, port: port}}
}

func NewTcpClientTimeout(host string, port uint, timeout time.Duration) IoClient {
	return &mbClient{&tcpTransporter{host: host, port: port, timeout: timeout}}
}
