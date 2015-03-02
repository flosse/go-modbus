package modbus

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_Pdu(t *testing.T) {

	Convey("Given a pdu struct", t, func() {
		data := []byte{7, 8}
		pdu := &Pdu{4, data}

		Convey("When we pack it", func() {
			bin := pdu.pack()

			Convey("the length of the binary array should be correct", func() {
				So(len(bin), ShouldEqual, 3)
			})

			Convey("the function code should be encoded", func() {
				So(bin[0], ShouldEqual, 4)
			})

			Convey("the data should be added", func() {
				So(bin[2], ShouldEqual, 8)
			})
		})
	})

	Convey("Given a valid binary pdu", t, func() {
		bin := []byte{3, 7, 8}

		Convey("When we unpack it", func() {
			pdu, err := unpackPdu(bin)

			Convey("we should not get an error", func() {
				So(err, ShouldBeNil)
			})

			Convey("the function code should be decoded", func() {
				So(pdu.Function, ShouldEqual, 3)
			})

			Convey("the data field should be corret", func() {
				So(len(pdu.Data), ShouldEqual, 2)
				So(pdu.Data[0], ShouldEqual, 7)
			})
		})
	})

	Convey("Given an invalid binary pdu", t, func() {
		bin := []byte{}

		Convey("When we unpack it", func() {
			_, err := unpackPdu(bin)

			Convey("we should get an error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}
