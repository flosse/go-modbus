/**
 * Copyright (C) 2014 - 2015, Markus Kohlhase <mail@markus-kohlhase.de>
 *
 * Modbus Master (Client) implementation
 */

package modbus

import (
	"encoding/binary"
	"fmt"
	"math"
)

type mbClient struct {
	transport Transporter
}

type io struct {
	master  Client
	address uint16
	count   uint16
}

type roBit io
type rwBit io

type roRegister io
type roRegisters io
type rwRegister io
type rwRegisters io

func (c *mbClient) request(f uint8, addr uint16, data []byte) (pdu *Pdu, err error) {
	pdu, err = c.transport.Send(&Pdu{f, append(wordsToByteArray(addr), data...)})
	if err != nil {
		return
	}
	if errFn := f + 0x80; errFn == pdu.Function {
		return nil, Error{errFn, pdu.Data[0]}
	}
	return
}

func (c *mbClient) readRegisters(fn uint8, addr, count uint16) (values []uint16, err error) {
	res, err := c.request(fn, addr, wordsToByteArray(count))
	if err != nil {
		return
	}
	values = bytesToWordArray(res.Data[1:]...)
	return
}

func (c *mbClient) ReadDiscreteInputs(addr, count uint16) (result []bool, err error) {
	resp, err := c.request(2, addr, wordsToByteArray(count))
	if err != nil {
		return
	}
	result = make([]bool, count)

	inputs := resp.Data[1:]
	byteNr := uint(0)
	bitNr := uint(0)

	for i := 0; i < int(count); i++ {
		result[i] = inputs[byteNr]&(1<<bitNr) != 0
		if (i+1)%8 == 0 {
			bitNr = 0
			byteNr++
		} else {
			bitNr++
		}
	}
	return
}

func (c *mbClient) ReadCoils(addr, count uint16) (coils []bool, err error) {
	res, err := c.request(1, addr, wordsToByteArray(count))
	if err != nil {
		return
	}
	byteCount := int(res.Data[0])
	coilStates := res.Data[1:]
	r := make([]bool, byteCount*8)
	for i := 0; i < byteCount; i++ {
		for j := 0; j < 8; j++ {
			r[i*8+j] = bool((int(coilStates[i]) & (1 << uint(j))) > 0)
		}
	}
	return r[0:count], nil
}

func (c *mbClient) WriteSingleCoil(addr uint16, value bool) (err error) {
	var set uint8
	if value {
		set = 0xff
	}
	_, err = c.request(5, addr, []byte{set, uint8(0)})
	return
}

func (c *mbClient) WriteMultipleCoils(addr uint16, values []bool) (err error) {
	count := len(values)
	byteCount := uint(math.Ceil(float64(count) / 8))
	data := make([]byte, 3+byteCount)

	binary.BigEndian.PutUint16(data, uint16(count))
	data[2] = uint8(byteCount)

	byteNr := uint(0)
	bitNr := uint8(0)
	byteVal := uint8(0)

	for v := 0; v < count; v++ {
		if v == count-1 {
			data[byteNr+3] = byteVal
			break
		}
		if values[v] {
			byteVal |= 1 << bitNr
		}
		if bitNr > 6 {
			data[byteNr+3] = byteVal
			byteVal = 0
			bitNr = 0
			byteNr++
		} else {
			bitNr++
		}
	}
	res, err := c.request(15, addr, data)
	if err != nil {
		return
	}
	if binary.BigEndian.Uint16(res.Data[2:]) != uint16(count) {
		return fmt.Errorf("%d coils were forced instead of %d", count)
	}
	return
}

func (c *mbClient) ReadInputRegisters(addr, count uint16) ([]uint16, error) {
	return c.readRegisters(4, addr, count)
}

func (c *mbClient) ReadHoldingRegisters(addr, count uint16) ([]uint16, error) {
	return c.readRegisters(3, addr, count)
}

func (c *mbClient) WriteMultipleRegisters(addr uint16, values []uint16) (err error) {

	regCount := len(values)
	byteCount := regCount * 2
	data := make([]byte, byteCount+3)
	data[0] = uint8(regCount >> 8)
	data[1] = uint8(regCount & 0xff)
	data[2] = uint8(byteCount)

	for i := 0; i < regCount; i++ {
		data[i*2+3] = uint8(values[i] >> 8)
		data[i*2+4] = uint8(values[i] & 0xff)
	}
	_, err = c.request(16, addr, data)
	return
}

func (c *mbClient) WriteSingleRegister(addr uint16, value uint16) (err error) {
	_, err = c.request(6, addr, []byte{uint8(value >> 8), uint8(value & 0xff)})
	return
}

func (c *mbClient) ReadWriteMultipleRegisters(readAddress, readQuantity, writeAddress uint16, vals []uint16) (values []uint16, err error) {
	writeQuantity := len(vals)
	data := wordsToByteArray(readQuantity, writeAddress, uint16(writeQuantity))
	data = append(data, uint8(writeQuantity*2))
	data = append(data, wordsToByteArray(vals...)...)
	resp, err := c.request(23, readAddress, data)
	if err != nil {
		return
	}
	byteCount := resp.Data[0]
	if byteCount > 0 {
		values = bytesToWordArray(resp.Data[1:]...)
	} else {
		values = []uint16{}
	}
	return
}

func (c *mbClient) MaskWriteRegister(addr, and, or uint16) (err error) {
	_, err = c.request(22, addr, wordsToByteArray(and, or))
	return
}

func (c *mbClient) ReadFifoQueue(addr uint16) (fifoValues []uint16, err error) {
	resp, err := c.request(24, addr, nil)
	fifoValues = bytesToWordArray(resp.Data[4:]...)
	return
}

func (c *mbClient) ReadExceptionStatus() (states []bool, err error) {
	pdu, err := c.transport.Send(&Pdu{7, nil})
	if err != nil {
		return
	}

	states = make([]bool, 8)

	for bit := 0; bit < 8; bit++ {
		states[bit] = bool((pdu.Data[0] & (1 << uint(bit))) > 0)
	}
	return
}

func (c *mbClient) Diagnostics(subfunction uint16, data []uint16) (result []uint16, err error) {
	resp, err := c.request(8, subfunction, wordsToByteArray(data...))
	if err != nil {
		return
	}
	return bytesToWordArray(resp.Data[2:]...), nil
}

func (c *mbClient) GetCommEventCounter() (status bool, count uint16, err error) {
	pdu, err := c.transport.Send(&Pdu{11, nil})
	if err != nil {
		return
	}
	res := bytesToWordArray(pdu.Data...)
	return (res[0] > 0), res[1], err
}

func (c *mbClient) ReportServerId() (response []byte, err error) {
	pdu, err := c.transport.Send(&Pdu{17, nil})
	if err != nil {
		return
	}
	return pdu.Data, nil
}

func (c *mbClient) DiscreteInput(addr uint16) DiscreteInput {
	return &roBit{master: c, address: addr}
}

func (c *mbClient) Coil(addr uint16) Coil {
	return &rwBit{master: c, address: addr}
}

func (c *mbClient) InputRegister(addr uint16) InputRegister {
	return &roRegister{master: c, address: addr}
}

func (c *mbClient) InputRegisters(addr, count uint16) InputRegisters {
	return &roRegisters{c, addr, count}
}

func (c *mbClient) HoldingRegister(addr uint16) HoldingRegister {
	return &rwRegister{master: c, address: addr}
}

func (c *mbClient) HoldingRegisters(addr, count uint16) HoldingRegisters {
	return &rwRegisters{c, addr, count}
}

func (io *roBit) Test() (result bool, err error) {
	res, err := io.master.ReadDiscreteInputs(io.address, 1)
	if err != nil {
		return
	}
	return res[0], nil
}

func (io *rwBit) Test() (result bool, err error) {
	res, err := io.master.ReadCoils(io.address, 1)
	if err != nil {
		return
	}
	return res[0], err
}

func (io *rwBit) Set() (err error) {
	return io.master.WriteSingleCoil(io.address, true)
}

func (io *rwBit) Clear() (err error) {
	return io.master.WriteSingleCoil(io.address, false)
}

func (io *rwBit) Toggle() (err error) {
	if res, err := io.master.ReadCoils(io.address, 1); err == nil {
		return io.master.WriteSingleCoil(io.address, !res[0])
	}
	return
}

func (io *roRegister) Read() (value uint16, err error) {
	if res, err := io.master.ReadInputRegisters(io.address, 1); err == nil {
		value = res[0]
	}
	return
}

func (io *rwRegister) Read() (value uint16, err error) {
	if res, err := io.master.ReadHoldingRegisters(io.address, 1); err == nil {
		value = res[0]
	}
	return
}

func (io *rwRegister) Write(value uint16) (err error) {
	return io.master.WriteSingleRegister(io.address, value)
}

func (io *roRegisters) Read() (values []uint16, err error) {
	return io.master.ReadInputRegisters(io.address, io.count)
}

func (io *rwRegisters) Read() (values []uint16, err error) {
	return io.master.ReadInputRegisters(io.address, io.count)
}

func (io *rwRegisters) Write(values []uint16) (err error) {
	return io.master.WriteMultipleRegisters(io.address, values)
}
