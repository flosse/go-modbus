package modbus

import (
	"encoding/binary"
	"testing"
)

func Test_Header_Pack(t *testing.T) {

	bin := (&Header{1234, 99, 42, 9}).Pack()

	if len(bin) != 7 {
		t.Error("invalid binary length")
	}
	if binary.BigEndian.Uint16(bin[0:2]) != 1234 {
		t.Error("invalid transaction number")
	}
	if binary.BigEndian.Uint16(bin[2:4]) != 99 {
		t.Error("invalid protocol id")
	}
	if binary.BigEndian.Uint16(bin[4:6]) != 42 {
		t.Error("invalid pdu lengh")
	}
	if bin[6] != 9 {
		t.Error("invalid uni id")
	}
}

func Test_UnpackHeader(t *testing.T) {

	if _, err := UnpackHeader([]byte{0, 0, 0, 0, 0, 0}); err == nil {
		t.Error("an error should be returned")
	}

	h, _ := UnpackHeader([]byte{0xff, 0xff, 0, 5, 0, 3, 9})
	if h.transaction != 65535 {
		t.Error("invalid transaction id")
	}
	if h.protocol != 5 {
		t.Error("invalid protocol id")
	}
	if h.length != 3 {
		t.Error("invalid pdu length")
	}
	if h.unit != 9 {
		t.Error("invalid unit id")
	}
}

func Test_Adu_Pack(t *testing.T) {
	h := &Header{1, 2, 3, 4}
	pdu := &Pdu{6, []byte{2, 4}}
	adu := &Adu{h, pdu}
	bin := adu.Pack()
	if len(bin) != 10 {
		t.Error("invalid binary length")
	}
	if bin[6] != 4 {
		t.Error("invalid header")
	}
	if bin[7] != 6 {
		t.Error("invalid function code field")
	}
	if bin[9] != 4 {
		t.Error("invalid data field")
	}
}

func Test_UnpackAdu(t *testing.T) {

	if _, err := UnpackAdu([]byte{0, 0, 0, 0, 0, 0, 0}); err == nil {
		t.Error("an error should be returned")
	}

	adu, _ := UnpackAdu([]byte{0, 0x0f, 0, 5, 0, 3, 9, 4, 2})

	if adu.header.transaction != 15 {
		t.Error("invalid transaction id")
	}
	if adu.pdu.function != 4 {
		t.Error("invalid function id")
	}
	if adu.pdu.data[0] != 2 {
		t.Error("invalid data field")
	}

}
