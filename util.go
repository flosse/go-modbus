package modbus

import (
	"encoding/binary"
	"math"
)

func wordsToByteArray(words ...uint16) []byte {
	array := make([]byte, 2*len(words))
	for i, v := range words {
		binary.BigEndian.PutUint16(array[i*2:], v)
	}
	return array
}

func bytesToWordArray(bytes ...byte) []uint16 {
	l := len(bytes)
	n := int(math.Ceil(float64(l) / 2))
	array := make([]uint16, n)
	for i := 0; i < n; i++ {
		j := i * 2
		if j+2 > l {
			array[i] = uint16(bytes[j])
		} else {
			array[i] = binary.BigEndian.Uint16(bytes[j : j+2])
		}
	}
	return array
}

func filterNullChar(a []byte) []byte {
	x := filter(a, func(c byte) bool {
		return c != 0
	})
	return x
}

func filter(s []byte, fn func(byte) bool) []byte {
	var p []byte // == nil
	for _, v := range s {
		if fn(v) {
			p = append(p, v)
		}
	}
	return p
}
