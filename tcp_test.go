package modbus

import (
	"encoding/binary"
	. "github.com/smartystreets/goconvey/convey"
	"net"
	"testing"
	"time"
)

type dummyAddr struct{}

func (a *dummyAddr) Network() string {
	return ""
}

func (a *dummyAddr) String() string {
	return ""
}

type dummyConn struct{}

func (c *dummyConn) Read(b []byte) (n int, err error) {
	return 0, nil
}

func (c *dummyConn) Write(b []byte) (n int, err error) {
	return 0, nil

}

func (c *dummyConn) Close() error {
	return nil
}

func (c *dummyConn) LocalAddr() net.Addr {
	return &dummyAddr{}
}

func (c *dummyConn) RemoteAddr() net.Addr {
	return &dummyAddr{}
}

func (c *dummyConn) SetDeadline(t time.Time) error {
	return nil
}

func (c *dummyConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *dummyConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func Test_Tcp(t *testing.T) {
	Convey("Given a header struct", t, func() {
		header := &header{1234, 99, 42, 9}

		Convey("When we pack it", func() {
			bin := header.pack()

			Convey("The length of the binary array should be 7", func() {
				So(len(bin), ShouldEqual, 7)
			})

			Convey("The transaction number should be encoded as BigEndian uint16", func() {
				So(binary.BigEndian.Uint16(bin[0:2]), ShouldEqual, 1234)
			})

			Convey("The protocol id should be encoded as BigEndian uint16", func() {
				So(binary.BigEndian.Uint16(bin[2:4]), ShouldEqual, 99)
			})

			Convey("The pdu length should be encoded as BigEndian uint16", func() {
				So(binary.BigEndian.Uint16(bin[4:6]), ShouldEqual, 42)
			})

			Convey("The uni id should be the last byte", func() {
				So(bin[6], ShouldEqual, 9)
			})
		})
	})

	Convey("Given a invalid binary header", t, func() {
		header := []byte{0, 0, 0, 0, 0, 0}

		Convey("When we unpack it", func() {
			_, err := unpackHeader(header)

			Convey("we should get an error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})

	Convey("Given a valid binary header", t, func() {
		header := []byte{0xff, 0xff, 0, 5, 0, 3, 9}

		Convey("When we unpack it", func() {
			h, err := unpackHeader(header)

			Convey("we should not get an error", func() {
				So(err, ShouldBeNil)
			})

			Convey("the transaction id should be decoded", func() {
				So(h.transaction, ShouldEqual, 65535)
			})

			Convey("the protocol id should be decoded", func() {
				So(h.protocol, ShouldEqual, 5)
			})

			Convey("the pdu length should be correct", func() {
				So(h.length, ShouldEqual, 3)
			})

			Convey("the unit id should be decoded", func() {
				So(h.unit, ShouldEqual, 9)
			})
		})
	})

	Convey("Given an adu struct", t, func() {
		adu := &adu{&header{1, 2, 3, 4}, &Pdu{6, []byte{2, 4}}}

		Convey("When we pack it", func() {
			bin, _ := adu.pack()

			Convey("the byte array length should correct", func() {
				So(len(bin), ShouldEqual, 10)
			})

			Convey("the header should be encoded correctly", func() {
				So(bin[6], ShouldEqual, 4)
				So(bin[3], ShouldEqual, 2)
			})

			Convey("the function code should be correct", func() {
				So(bin[7], ShouldEqual, 6)
			})

			Convey("the data code should be included", func() {
				So(bin[9], ShouldEqual, 4)
			})
		})
	})

	Convey("Given an invalid binary adu", t, func() {
		bin := []byte{0, 0, 0, 0, 0, 0, 0}

		Convey("When we unpack it", func() {
			_, err := unpackAdu(bin)

			Convey("we should get an error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})

	Convey("Given a valid binary adu", t, func() {
		bin := []byte{0, 0x0f, 0, 5, 0, 3, 9, 4, 2}

		Convey("When we unpack it", func() {
			adu, err := unpackAdu(bin)

			Convey("we should not get an error", func() {
				So(err, ShouldBeNil)
			})

			Convey("the header should be unpacked", func() {
				So(adu.header.transaction, ShouldEqual, 15)
			})

			Convey("the pdu should be unpacked", func() {
				So(adu.pdu.Function, ShouldEqual, 4)
				So(adu.pdu.Data[0], ShouldEqual, 2)
			})
		})
	})

	Convey("Given a tcpTransporter", t, func() {
		tr := &tcpTransporter{"foo", 502, nil, 0, 0}
		tr.connection = &dummyConn{}

		Convey("when sending a pdu", func() {
			req := &Pdu{3, nil}

			Convey("the transaction id should be incremented", func() {
				tr.Send(req)
				So(tr.transaction, ShouldEqual, 1)
			})

		})
	})
}
