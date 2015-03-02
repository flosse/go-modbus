/**
 * Copyright (C) 2014 - 2015, Markus Kohlhase <mail@markus-kohlhase.de>
 */

package modbus

import (
	"fmt"
)

/* Modbus Error */

type Error struct {

	// Error Code
	Code uint8

	// Exception Code
	Exception uint8
}

func getExceptionMessage(nr uint8) string {
	switch nr {
	case 0x01:
		return "ILLEGAL FUNCTION"
	case 0x02:
		return "ILLEGAL DATA ADDRESS"
	case 0x03:
		return "ILLEGAL DATA VALUE"
	case 0x04:
		return "SERVER DEVICE FAILURE"
	case 0x05:
		return "ACKNOWLEDGE"
	case 0x06:
		return "SERVER DEVICE BUSY"
	case 0x08:
		return "MEMORY PARITY ERROR"
	case 0x0A:
		return "GATEWAY PATH UNAVAILABLE"
	case 0x0B:
		return "GATEWAY TARGET DEVICE FAILED TO RESPOND"

	default:
		return "UNKNOWN EXCEPTION"
	}
}

func (e Error) Error() string {
	return fmt.Sprintf("Error %d (Function %d); Exception %d ('%s')", e.Code, (e.Code - 128), e.Exception, getExceptionMessage(e.Exception))
}
