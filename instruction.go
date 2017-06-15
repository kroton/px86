package main

import (
	"fmt"
)

type instruction func(*state) error

var instructions = map[uint8]instruction{}

func movR32Imm32(s *state) error {
	r, err := s.getUint8(0)
	if err != nil {
		return err
	}

	v, err := s.getUint32(1)
	if err != nil {
		return err
	}

	if err := s.setRegister(r-0xb8, v); err != nil {
		return err
	}

	s.eip += 5
	return nil
}

func movRm32Imm32(s *state) error {
	s.eip += 1

	m, err := parseModrm(s)
	if err != nil {
		return err
	}

	v, err := s.getUint32(0)
	if err != nil {
		return err
	}
	s.eip += 4

	return setRm32(s, m, v)
}

func movRm32R32(s *state) error {
	s.eip += 1

	m, err := parseModrm(s)
	if err != nil {
		return err
	}

	r32, err := getR32(s, m)
	if err != nil {
		return err
	}

	return setRm32(s, m, r32)
}

func movR32Rm32(s *state) error {
	s.eip += 1

	m, err := parseModrm(s)
	if err != nil {
		return err
	}

	rm32, err := getRm32(s, m)
	if err != nil {
		return err
	}

	return setR32(s, m, rm32)
}

func shortJump(s *state) error {
	d, err := s.getInt8(1)
	if err != nil {
		return err
	}

	s.eip += 2
	s.eip += uint32(int32(d))
	return nil
}

func nearJump(s *state) error {
	d, err := s.getInt32(1)
	if err != nil {
		return err
	}

	s.eip += 5
	s.eip += uint32(d)
	return nil
}

func addRm32R32(s *state) error {
	s.eip += 1

	m, err := parseModrm(s)
	if err != nil {
		return err
	}

	r32, err := getR32(s, m)
	if err != nil {
		return err
	}
	rm32, err := getRm32(s, m)
	if err != nil {
		return err
	}

	return setRm32(s, m, rm32+r32)
}

func addRm32Imm8(s *state, m modrm) error {
	rm32, err := getRm32(s, m)
	if err != nil {
		return err
	}
	imm8, err := s.getInt8(0)
	if err != nil {
		return err
	}

	s.eip += 1
	return setRm32(s, m, rm32+uint32(int32(imm8)))
}

func subRm32Imm8(s *state, m modrm) error {
	rm32, err := getRm32(s, m)
	if err != nil {
		return err
	}
	imm8, err := s.getInt8(0)
	if err != nil {
		return err
	}
	s.eip += 1

	return setRm32(s, m, rm32-uint32(int32(imm8)))
}

func code83(s *state) error {
	s.eip += 1

	m, err := parseModrm(s)
	if err != nil {
		return err
	}

	switch m.reg {
	case 0:
		return addRm32Imm8(s, m)
	case 5:
		return subRm32Imm8(s, m)
	}

	return fmt.Errorf("not implemented: 83 /%d", m.reg)
}

func incRm32(s *state, m modrm) error {
	v, err := getRm32(s, m)
	if err != nil {
		return err
	}
	return setRm32(s, m, v+1)
}

func codeFF(s *state) error {
	s.eip += 1

	m, err := parseModrm(s)
	if err != nil {
		return err
	}

	switch m.reg {
	case 0:
		return incRm32(s, m)
	}

	return fmt.Errorf("not implemented: ff /%d", m.reg)
}

func pushR32(s *state) error {
	r, err := s.getUint8(0)
	if err != nil {
		return err
	}
	v, err := s.getRegister(r - 0x50)
	if err != nil {
		return err
	}

	if err := s.push(v); err != nil {
		return err
	}
	s.eip += 1
	return nil
}

func pushImm32(s *state) error {
	v, err := s.getUint32(1)
	if err != nil {
		return err
	}
	if err := s.push(v); err != nil {
		return err
	}

	s.eip += 5
	return nil
}

func pushImm8(s *state) error {
	v, err := s.getUint8(1)
	if err != nil {
		return err
	}
	if err := s.push(uint32(v)); err != nil {
		return err
	}

	s.eip += 2
	return nil
}

func popR32(s *state) error {
	r, err := s.getUint8(0)
	if err != nil {
		return err
	}
	v, err := s.pop()
	if err != nil {
		return err
	}

	if err := s.setRegister(r-0x58, v); err != nil {
		return err
	}
	s.eip += 1
	return nil
}

func callRel32(s *state) error {
	d, err := s.getInt32(1)
	if err != nil {
		return err
	}

	if err := s.push(s.eip + 5); err != nil {
		return err
	}

	s.eip += 5
	s.eip += uint32(d)
	return nil
}

func ret(s *state) error {
	v, err := s.pop()
	if err != nil {
		return err
	}

	s.eip = v
	return nil
}

func leave(s *state) error {
	ebpVal, err := s.getRegister(ebp)
	if err != nil {
		return err
	}
	if err := s.setRegister(esp, ebpVal); err != nil {
		return err
	}

	v, err := s.pop()
	if err != nil {
		return err
	}
	if err := s.setRegister(ebp, v); err != nil {
		return err
	}

	s.eip += 1
	return nil
}

func init() {
	instructions[0x01] = addRm32R32
	for i := uint8(0); i < registersSize; i++ {
		instructions[0x50+i] = pushR32
	}
	for i := uint8(0); i < registersSize; i++ {
		instructions[0x58+i] = popR32
	}
	instructions[0x68] = pushImm32
	instructions[0x6a] = pushImm8
	instructions[0x83] = code83
	instructions[0x89] = movRm32R32
	instructions[0x8b] = movR32Rm32
	for i := uint8(0); i < registersSize; i++ {
		instructions[0xb8+i] = movR32Imm32
	}
	instructions[0xc3] = ret
	instructions[0xc9] = leave
	instructions[0xc7] = movRm32Imm32
	instructions[0xe8] = callRel32
	instructions[0xe9] = nearJump
	instructions[0xeb] = shortJump
	instructions[0xff] = codeFF
}
