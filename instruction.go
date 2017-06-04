package main

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

	s.setRegister(int(r-0xb8), v)
	s.advanceEIP(5)
	return nil
}

func shortJump(s *state) error {
	d, err := s.getInt8(1)
	if err != nil {
		return err
	}

	s.advanceEIP(int(d) + 2)
	return nil
}

func init() {
	for i := 0; i < registersSize; i++ {
		instructions[uint8(0xb8+i)] = movR32Imm32
	}
	instructions[uint8(0xeb)] = shortJump
}
