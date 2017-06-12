package main

import (
	"fmt"
)

type modrm struct {
	mod    uint8
	reg    uint8
	rm     uint8
	sib    uint8
	disp8  uint8
	disp32 uint32
}

type addressingNotImplemented modrm

func (a addressingNotImplemented) Error() string {
	return fmt.Sprintf("not implemented ModRM mod = %v, rm = %v", a.mod, a.rm)
}

func parseModrm(s *state) (modrm, error) {
	var m modrm

	code, err := s.getUint8(0)
	if err != nil {
		return modrm{}, err
	}
	s.advanceEIP(1)

	m.mod = (code & 0xc0) >> 6
	m.reg = (code & 0x38) >> 3
	m.rm = code & 0x07

	if m.mod != 3 && m.rm == 4 {
		m.sib, err = s.getUint8(0)
		if err != nil {
			return modrm{}, err
		}
		s.advanceEIP(1)
	}

	if (m.mod == 0 && m.rm == 5) || m.mod == 2 {
		m.disp32, err = s.getUint32(0)
		if err != nil {
			return modrm{}, err
		}
		s.advanceEIP(4)
	} else if m.mod == 1 {
		m.disp8, err = s.getUint8(0)
		if err != nil {
			return modrm{}, err
		}
		s.advanceEIP(1)
	}

	return m, nil
}

func calcAddress(s *state, m modrm) (uint32, error) {
	if m.mod == 0 {
		if m.rm == 4 {
			return 0, addressingNotImplemented(m)
		}
		if m.rm == 5 {
			return m.disp32, nil
		}
		return s.getRegister(int(m.rm))
	}

	if m.mod == 1 {
		if m.rm == 4 {
			return 0, addressingNotImplemented(m)
		}

		v, err := s.getRegister(int(m.rm))
		if err != nil {
			return 0, err
		}
		return v + uint32(m.disp8), nil
	}

	if m.mod == 2 {
		if m.rm == 4 {
			return 0, addressingNotImplemented(m)
		}

		v, err := s.getRegister(int(m.rm))
		if err != nil {
			return 0, err
		}
		return v + m.disp32, nil
	}

	return 0, addressingNotImplemented(m)
}

func getRm32(s *state, m modrm) (uint32, error) {
	if m.mod == 3 {
		return s.getRegister(int(m.rm))
	}

	addr, err := calcAddress(s, m)
	if err != nil {
		return 0, err
	}

	return s.getUint32Addr(int(addr))
}

func setRm32(s *state, m modrm, v uint32) error {
	if m.mod == 3 {
		return s.setRegister(int(m.rm), v)
	}

	addr, err := calcAddress(s, m)
	if err != nil {
		return err
	}

	return s.setUint32Addr(int(addr), v)
}

func getR32(s *state, m modrm) (uint32, error) {
	return s.getRegister(int(m.reg))
}

func setR32(s *state, m modrm, v uint32) error {
	return s.setRegister(int(m.reg), v)
}
