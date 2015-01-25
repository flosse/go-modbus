package modbus

import (
	"fmt"
)

/* Modbus Error */

type Error struct {

	// Error Code
	code uint8

	// Exception Code
	exception uint8
}

func (e Error) Error() string {
	return fmt.Sprintf("Modbus error: %d; Exception: %d", e.code, e.exception)
}
