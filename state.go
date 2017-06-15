package main

import (
	"fmt"
	"io"
	"os"
)

const (
	programMaxLength = 512
)

type registerOutOfRange uint8

func (r registerOutOfRange) Error() string {
	return fmt.Sprintf("out of registers range: %x", uint8(r))
}

type memoryOutOfRange uint32

func (m memoryOutOfRange) Error() string {
	return fmt.Sprintf("out of memory range: %x", uint32(m))
}

type state struct {
	registers []uint32
	memory    []uint8
	progBegin uint32
	progLen   uint32
	eip       uint32
}

func newState(size uint32, eip uint32, espVal uint32, progBegin uint32) *state {
	s := &state{
		registers: make([]uint32, registersSize),
		memory:    make([]uint8, size),
		progBegin: progBegin,
		progLen:   0,
		eip:       eip,
	}
	s.registers[esp] = espVal
	return s
}

func (s *state) dumpRegisters() {
	for i := uint8(0); i < registersSize; i++ {
		fmt.Fprintf(os.Stderr, "%s = %08x\n", registerName(i), s.registers[i])
	}
	fmt.Fprintf(os.Stderr, "EIP = %08x\n", s.eip)
}

func (s *state) Write(p []byte) (n int, err error) {
	for i, v := range p {
		j := s.progBegin + s.progLen
		if s.progLen >= programMaxLength || j >= s.memorySize() {
			return i, io.EOF
		}
		s.memory[j] = v
		s.progLen++
	}
	return len(p), nil
}

func (s *state) registersSize() uint8 {
	return uint8(len(s.registers))
}

func (s *state) memorySize() uint32 {
	return uint32(len(s.memory))
}

func (s *state) getRegister(r uint8) (uint32, error) {
	if r >= s.registersSize() {
		return 0, registerOutOfRange(r)
	}
	return s.registers[r], nil
}

func (s *state) setRegister(r uint8, v uint32) error {
	if r >= s.registersSize() {
		return registerOutOfRange(r)
	}
	s.registers[r] = v
	return nil
}

func (s *state) getUint8(offset int) (uint8, error) {
	return s.getUint8Addr(s.eip + uint32(offset))
}

func (s *state) getInt8(offset int) (int8, error) {
	return s.getInt8Addr(s.eip + uint32(offset))
}

func (s *state) getUint32(offset int) (uint32, error) {
	return s.getUint32Addr(s.eip + uint32(offset))
}

func (s *state) getInt32(offset int) (int32, error) {
	return s.getInt32Addr(s.eip + uint32(offset))
}

func (s *state) getUint8Addr(addr uint32) (uint8, error) {
	if addr >= s.memorySize() {
		return 0, memoryOutOfRange(addr)
	}
	return s.memory[addr], nil
}

func (s *state) getInt8Addr(addr uint32) (int8, error) {
	v, err := s.getUint8Addr(addr)
	if err != nil {
		return 0, err
	}
	return int8(v), nil
}

func (s *state) setUint8Addr(addr uint32, v uint8) error {
	if addr >= s.memorySize() {
		return memoryOutOfRange(addr)
	}
	s.memory[addr] = v
	return nil
}

func (s *state) getUint32Addr(addr uint32) (uint32, error) {
	r := uint32(0)
	for i := uint32(0); i < 4; i++ {
		v, err := s.getUint8Addr(addr + i)
		if err != nil {
			return 0, err
		}
		r |= uint32(v) << (i * 8)
	}
	return r, nil
}

func (s *state) getInt32Addr(addr uint32) (int32, error) {
	v, err := s.getUint32Addr(addr)
	if err != nil {
		return 0, err
	}
	return int32(v), nil
}

func (s *state) setUint32Addr(addr uint32, v uint32) error {
	for i := uint32(0); i < 4; i++ {
		v8 := (v >> (i * 8)) & 0xff
		if err := s.setUint8Addr(addr+i, uint8(v8)); err != nil {
			return err
		}
	}
	return nil
}

func (s *state) push(v uint32) error {
	r, err := s.getRegister(esp)
	if err != nil {
		return err
	}
	addr := r - 4

	if err := s.setRegister(esp, addr); err != nil {
		return err
	}
	if err := s.setUint32Addr(addr, v); err != nil {
		return err
	}
	return nil
}

func (s *state) pop() (uint32, error) {
	addr, err := s.getRegister(esp)
	if err != nil {
		return 0, err
	}
	v, err := s.getUint32Addr(addr)
	if err != nil {
		return 0, err
	}

	if err := s.setRegister(esp, addr+4); err != nil {
		return 0, err
	}
	return v, nil
}
