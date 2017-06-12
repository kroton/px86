package main

import (
	"fmt"
	"io"
	"os"
)

const (
	programMaxLength = 512
)

type registerOutOfRange int

func (r registerOutOfRange) Error() string {
	return fmt.Sprintf("out of registers range: %x", r)
}

type memoryOutOfRange int

func (m memoryOutOfRange) Error() string {
	return fmt.Sprintf("out of memory range: %x", m)
}

type state struct {
	registers []uint32
	memory    []uint8
	progBegin int
	progLen   int
	eip       int
}

func newState(size int, eip int, espVal uint32, progBegin int) *state {
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
	for i := 0; i < registersSize; i++ {
		fmt.Fprintf(os.Stderr, "%s = %08x\n", register(i), s.registers[i])
	}
	fmt.Fprintf(os.Stderr, "EIP = %08x\n", s.eip)
}

func (s *state) Write(p []byte) (n int, err error) {
	for i, v := range p {
		j := s.progBegin + s.progLen
		if s.progLen >= programMaxLength || j >= len(s.memory) {
			return i, io.EOF
		}
		s.memory[j] = v
		s.progLen++
	}
	return len(p), nil
}

func (s *state) getEIP() int {
	return s.eip
}

func (s *state) advanceEIP(d int) {
	s.eip += d
}

func (s *state) setRegister(r int, v uint32) error {
	if r < 0 || r >= len(s.registers) {
		return registerOutOfRange(r)
	}
	s.registers[r] = v
	return nil
}

func (s *state) getUint8(offset int) (uint8, error) {
	i := s.eip + offset
	if i < 0 || i >= len(s.memory) {
		return 0, memoryOutOfRange(i)
	}
	return s.memory[i], nil
}

func (s *state) getInt8(offset int) (int8, error) {
	v, err := s.getUint8(offset)
	if err != nil {
		return 0, err
	}
	return int8(v), nil
}

func (s *state) getUint32(offset int) (uint32, error) {
	r := uint32(0)
	for i := 0; i < 4; i++ {
		v, err := s.getUint8(offset + i)
		if err != nil {
			return 0, err
		}
		r |= uint32(v) << uint(i*8)
	}
	return r, nil
}

func (s *state) getInt32(offset int) (int32, error) {
	v, err := s.getUint32(offset)
	if err != nil {
		return 0, err
	}
	return int32(v), err
}

func (s *state) hasNext() bool {
	return 0 <= s.eip && s.eip < len(s.memory)
}

func (s *state) isEnd() bool {
	return s.eip == 0
}
