package modbus

import (
	"testing"
)

func Test_Pack(t *testing.T) {
  data := []byte{7,8}
	pdu  := &Pdu{4,data}
  bin  := pdu.Pack()

  if len(bin) != 3 {
		t.Error("invalid binary length")
	}
  if bin[0] != 4 {
		t.Error("invalid function code field")
	}
  if bin[2] != 8 {
		t.Error("invalid data field")
	}
}

func Test_UnpackPdu(t *testing.T) {

  pdu, _ := UnpackPdu([]byte{3,7,8})

  if pdu.function != 3 {
		t.Error("function code should be 3")
	}

  if len(pdu.data) != 2 {
		t.Error("PDU data length should be 2")
	}

  if pdu, _ := UnpackPdu([]byte{14}); pdu.function != 14 {
		t.Error("PDU function code should be 14")
	}

  if _, err := UnpackPdu([]byte{}); err == nil {
		t.Error("an error should be returned")
	}
}
