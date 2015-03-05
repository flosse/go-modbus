package modbus

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_Util(t *testing.T) {

	Convey("Given some uint16 words", t, func() {
		words := []uint16{5, 2, 3}

		Convey("when converting them into an array of bytes", func() {
			bytes := wordsToByteArray(words...)

			Convey("the array length should be twice the amount of words", func() {
				So(len(bytes), ShouldEqual, len(words)*2)
			})

			Convey("the words should be encoded as big endian binary", func() {
				So(bytes[0], ShouldEqual, 0)
				So(bytes[1], ShouldEqual, 5)
			})
		})
	})

	Convey("Given an even amount of bytes", t, func() {
		evenBytes := []byte{0x04, 0xd2, 0, 0xf}

		Convey("when converting them into an array of uint16 words", func() {
			evenWords := bytesToWordArray(evenBytes...)

			Convey("the array length should be halve the amount of words", func() {
				So(len(evenWords), ShouldEqual, len(evenBytes)/2)
			})

			Convey("the words should be decoded as big endian binary", func() {
				So(evenWords[0], ShouldEqual, 1234)
				So(evenWords[1], ShouldEqual, 15)
			})
		})
	})

	Convey("Given an odd amount of bytes", t, func() {
		oddBytes := []byte{0x04, 0xd2, 0x0a}

		Convey("when converting them into an array of uint16 words", func() {
			oddWords := bytesToWordArray(oddBytes...)

			Convey("the array length should be halve the amount of words plus one", func() {
				So(len(oddWords), ShouldEqual, 1+len(oddBytes)/2)
			})

			Convey("last byte will be translated as uint16 too", func() {
				So(oddWords[0], ShouldEqual, 1234)
				So(oddWords[1], ShouldEqual, 10)
			})
		})
	})
}
