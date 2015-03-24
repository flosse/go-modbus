package modbus

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type dummyTransporter struct {

	// define a dummy response data
	resData []byte

	// cache current request
	req *Pdu

	// define a dummy response method
	send func(pdu *Pdu) (*Pdu, error)
}

func (t *dummyTransporter) Connect() error {
	return nil
}

func (t *dummyTransporter) Close() error {
	return nil
}

func (t *dummyTransporter) Send(pdu *Pdu) (resp *Pdu, err error) {

	t.req = pdu

	if t.send != nil {
		resp, err = t.send(pdu)
		t.resData = resp.Data
		return
	}

	return &Pdu{pdu.Function, t.resData}, nil

}

func getClient(respData []byte, send func(pdu *Pdu) (*Pdu, error)) (Client, *dummyTransporter) {
	t := &dummyTransporter{respData, nil, send}
	return &mbClient{t}, t
}

func getSerialClient(respData []byte, send func(pdu *Pdu) (*Pdu, error)) (SerialClient, *dummyTransporter) {
	t := &dummyTransporter{respData, nil, send}
	return &mbClient{t}, t
}

func getIoClient(respData []byte, send func(pdu *Pdu) (*Pdu, error)) (IoClient, *dummyTransporter) {
	t := &dummyTransporter{respData, nil, send}
	return &mbClient{t}, t
}

func Test_Client(t *testing.T) {

	Convey("Given a client", t, func() {
		c, d := getClient(nil, nil)

		// FUNCTION NR 1
		Convey("when reading coils", func() {

			c, d := getClient([]byte{0x03, 0xcd, 0x6b, 0x05}, nil)

			result, _ := c.ReadCoils(19, 19)

			Convey("the function nr should be 1", func() {
				So(d.req.Function, ShouldEqual, 1)
			})

			Convey("the first two data byte should encode the address", func() {
				So(d.req.Data[0], ShouldEqual, 0x00) // Hi
				So(d.req.Data[1], ShouldEqual, 0x13) // Lo
			})

			Convey("the following two data byte should encode the quantity", func() {
				So(d.req.Data[2], ShouldEqual, 0x00) // Hi
				So(d.req.Data[3], ShouldEqual, 0x13) // Lo
			})

			Convey("the response should be an array of coil states", func() {
				So(len(result), ShouldEqual, 19)
				So(result[11], ShouldEqual, true)  // register 36
				So(result[12], ShouldEqual, false) // register 37
				So(result[13], ShouldEqual, true)  // register 38
			})

		})

		// FUNCTION NR 2
		Convey("when reading discrete inputs", func() {

			c, d = getClient([]byte{0x03, 0xac, 0xdb, 0x35}, nil)

			result, _ := c.ReadDiscreteInputs(196, 22)

			Convey("the function nr should be 2", func() {
				So(d.req.Function, ShouldEqual, 2)
			})

			Convey("the response should be an array of input states", func() {
				So(len(result), ShouldEqual, 22)
				So(result[16], ShouldEqual, true)
				So(result[17], ShouldEqual, false)
				So(result[18], ShouldEqual, true)
				So(result[19], ShouldEqual, false)
				So(result[20], ShouldEqual, true)
				So(result[21], ShouldEqual, true)
			})

		})

		// FUNCTION NR 3
		Convey("when reading holding registers", func() {

			c, d = getClient([]byte{0x06, 0x02, 0x2b, 0x00, 0x00, 0x00, 0x64}, nil)

			result, _ := c.ReadHoldingRegisters(107, 3)

			Convey("the function nr should be 3", func() {
				So(d.req.Function, ShouldEqual, 3)
			})

			Convey("the response should be an array of register values", func() {
				So(len(result), ShouldEqual, 3)
				So(result[0], ShouldEqual, 555)
				So(result[1], ShouldEqual, 0)
				So(result[2], ShouldEqual, 100)
			})

		})

		// FUNCTION NR 4
		Convey("when reading input registers", func() {

			c, d = getClient([]byte{0x02, 0x00, 0x0a}, nil)

			result, _ := c.ReadInputRegisters(8, 1)

			Convey("the function nr should be 4", func() {
				So(d.req.Function, ShouldEqual, 4)
			})

			Convey("the response should be an array of register values", func() {
				So(len(result), ShouldEqual, 1)
				So(result[0], ShouldEqual, 10)
			})

		})

		// FUNCTION NR 5
		Convey("when writing a single coil", func() {

			c, d = getClient([]byte{0x00, 0xac, 0xff, 0x00}, nil)

			c.WriteSingleCoil(172, true)

			Convey("the function nr should be 5", func() {
				So(d.req.Function, ShouldEqual, 5)
			})

			Convey("the last word of the request should be 0xff00", func() {
				So(d.req.Data[2], ShouldEqual, 0xff) // Hi
				So(d.req.Data[3], ShouldEqual, 0x00) // Lo
			})

		})

		// FUNCTION NR 6
		Convey("when writing a single register", func() {

			c, d = getClient([]byte{0x00, 0x01, 0x00, 0x03}, nil)

			c.WriteSingleRegister(1, 3)

			Convey("the function nr should be 6", func() {
				So(d.req.Function, ShouldEqual, 6)
			})

			Convey("the register values should be encoded as big endian binary", func() {
				So(d.req.Data[2], ShouldEqual, 0x00) // Hi
				So(d.req.Data[3], ShouldEqual, 0x03) // Lo
			})
		})
	})

	Convey("Given a serial client", t, func() {

		// FUNCTION NR 7 (serial line only)
		Convey("when reading the exception status", func() {

			c, d := getSerialClient([]byte{0x6d}, nil)

			result, _ := c.ReadExceptionStatus()

			Convey("the function nr should be 7", func() {
				So(d.req.Function, ShouldEqual, 7)
			})

			Convey("the request data should be nil", func() {
				So(d.req.Data, ShouldBeNil)
			})

			Convey("response should be an array of boolean states", func() {
				So(result[0], ShouldEqual, true)
				So(result[1], ShouldEqual, false)
				So(result[2], ShouldEqual, true)
				So(result[3], ShouldEqual, true)
				So(result[4], ShouldEqual, false)
				So(result[5], ShouldEqual, true)
				So(result[6], ShouldEqual, true)
				So(result[7], ShouldEqual, false)
			})

		})

		// FUNCTION NR 8 (serial line only)
		Convey("when reading the diagnostics", func() {

			c, d := getSerialClient(nil, func(pdu *Pdu) (*Pdu, error) {
				return pdu, nil
			})

			result, _ := c.Diagnostics(0, []uint16{0xa537})

			Convey("the function nr should be 8", func() {
				So(d.req.Function, ShouldEqual, 8)
			})

			Convey("the request data of sub-function 0 should exist", func() {
				So(d.req.Data[0], ShouldEqual, 0)
				So(d.req.Data[2], ShouldEqual, 0xa5)
				So(d.req.Data[3], ShouldEqual, 0x37)
			})

			Convey("response of sub-function 0 should echo the request data", func() {
				So(result[0], ShouldEqual, 42295)
			})

		})

		// FUNCTION NR 11 (serial line only)
		Convey("when receiving the comm event counter", func() {

			c, d := getSerialClient([]byte{0xff, 0xff, 0x01, 0x08}, nil)

			state, count, _ := c.GetCommEventCounter()

			Convey("the function nr should be 11", func() {
				So(d.req.Function, ShouldEqual, 11)
			})

			Convey("the state and count values should be decoded", func() {
				So(state, ShouldEqual, true)
				So(count, ShouldEqual, 264)
			})

		})
	})

	Convey("Given an io client", t, func() {

		Convey("when creating a coil", func() {
			io, d := getIoClient([]byte{0x01, 0x01}, nil)
			coil := io.Coil(3)

			Convey("the test method should use function nr 1", func() {
				_, err := coil.Test()
				So(d.req.Function, ShouldEqual, 1)
				So(err, ShouldBeNil)
			})

			Convey("the test method should read the coil state", func() {
				x, _ := coil.Test()
				So(d.req.Data[1], ShouldEqual, 0x03)
				So(d.req.Data[2], ShouldEqual, 0x00)
				So(d.req.Data[3], ShouldEqual, 0x01)
				So(x, ShouldEqual, true)
			})

			Convey("the set method should write 0xff00", func() {
				coil.Set()
				So(d.req.Function, ShouldEqual, 5)
				So(d.req.Data[1], ShouldEqual, 0x03)
				So(d.req.Data[2], ShouldEqual, 0xff)
				So(d.req.Data[3], ShouldEqual, 0x00)
			})

			Convey("the clear method should write 0x0000", func() {
				coil.Clear()
				So(d.req.Function, ShouldEqual, 5)
				So(d.req.Data[2], ShouldEqual, 0x00)
			})

			Convey("the toggle method should invert the state", func() {
				io, d = getIoClient([]byte{0x01, 0x01}, nil)
				coil = io.Coil(0)
				coil.Toggle()
				So(d.req.Function, ShouldEqual, 5)
				So(d.req.Data[2], ShouldEqual, 0x00)
				io, d = getIoClient([]byte{0x01, 0x00}, nil)
				coil = io.Coil(0)
				coil.Toggle()
				So(d.req.Data[2], ShouldEqual, 0xff)
			})

		})

		Convey("when creating a discrete input", func() {
			io, d := getIoClient([]byte{0x02, 0xdf}, nil)
			di := io.DiscreteInput(3)

			Convey("the test method should use function nr 2", func() {
				_, err := di.Test()
				So(d.req.Function, ShouldEqual, 2)
				So(err, ShouldBeNil)
			})

			Convey("the test method should read the coil state", func() {
				x, _ := di.Test()
				So(x, ShouldEqual, true)
			})

		})

		Convey("when creating a holding register", func() {
			io, d := getIoClient([]byte{0x02, 0xda, 0x45}, nil)
			reg := io.HoldingRegister(0)

			Convey("the read method should use function nr 3", func() {
				_, err := reg.Read()
				So(d.req.Function, ShouldEqual, 3)
				So(err, ShouldBeNil)
			})

			Convey("the read method should read the value", func() {
				v, _ := reg.Read()
				So(d.req.Function, ShouldEqual, 3)
				So(v, ShouldEqual, 55877)
			})

			Convey("the write method should use function nr 6", func() {
				err := reg.Write(0)
				So(d.req.Function, ShouldEqual, 6)
				So(err, ShouldBeNil)
			})

			Convey("the write method should write the value", func() {
				reg.Write(1234)
				So(d.req.Data[2], ShouldEqual, 0x04)
				So(d.req.Data[3], ShouldEqual, 0xd2)
			})
		})

		Convey("when creating an input register", func() {
			io, d := getIoClient([]byte{0x02, 0xda, 0x45}, nil)
			reg := io.InputRegister(7)

			Convey("the read method should use function nr 4", func() {
				x, err := reg.Read()
				So(d.req.Function, ShouldEqual, 4)
				So(err, ShouldBeNil)
				So(x, ShouldEqual, 55877)
			})
		})

		Convey("when creating a multi input registers", func() {
			io, d := getIoClient([]byte{0x02, 0x00, 0x05, 0x00, 0x03}, nil)
			reg := io.InputRegisters(0x1000, 2)

			Convey("the read method should use function nr 4", func() {
				x, err := reg.Read()
				So(d.req.Function, ShouldEqual, 4)
				So(err, ShouldBeNil)
				So(x[0], ShouldEqual, 5)
				So(x[1], ShouldEqual, 3)
			})
		})

		Convey("when creating a multi holting registers", func() {
			io, d := getIoClient([]byte{0x02, 0x00, 0x05, 0x00, 0x03}, nil)
			reg := io.HoldingRegisters(0x1000, 2)

			Convey("the write method should use function nr x", func() {
				err := reg.Write([]uint16{3, 5})
				So(d.req.Function, ShouldEqual, 16)
				So(err, ShouldBeNil)
			})
		})
	})
}
