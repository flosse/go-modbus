package modbus

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_Error(t *testing.T) {

	Convey("Given an modbus error struct", t, func() {
		err := &Error{129, 0x01}

		Convey("it should print a human readable message", func() {
			So(err.Error(), ShouldEqual, "Error 129 (Function 1); Exception 1 ('ILLEGAL FUNCTION')")
		})
	})

	Convey("Given an exception nr", t, func() {
		Convey("it should print the exception message", func() {
			So(getExceptionMessage(1), ShouldEqual, "ILLEGAL FUNCTION")
			So(getExceptionMessage(2), ShouldEqual, "ILLEGAL DATA ADDRESS")
			So(getExceptionMessage(3), ShouldEqual, "ILLEGAL DATA VALUE")
			So(getExceptionMessage(4), ShouldEqual, "SERVER DEVICE FAILURE")
			So(getExceptionMessage(5), ShouldEqual, "ACKNOWLEDGE")
			So(getExceptionMessage(6), ShouldEqual, "SERVER DEVICE BUSY")
			So(getExceptionMessage(8), ShouldEqual, "MEMORY PARITY ERROR")
			So(getExceptionMessage(10), ShouldEqual, "GATEWAY PATH UNAVAILABLE")
			So(getExceptionMessage(11), ShouldEqual, "GATEWAY TARGET DEVICE FAILED TO RESPOND")
		})
	})
}
