package main

type instruction func(*state) error

var instructions = map[uint8]instruction{}

func movR32Imm32(s *state) error {
	r, err := s.GetUint8(0)
	if err != nil {
		return err
	}

	v, err := s.GetUint32(1)
	if err != nil {
		return err
	}

	s.SetRegister(int(r-0xb8), v)
	s.AdvanceEIP(5)
	return nil
}

func shortJump(s *state) error {
	d, err := s.GetInt8(1)
	if err != nil {
		return err
	}

	s.AdvanceEIP(int(d) + 2)
	return nil
}

func init() {
	for i := 0; i < registersSize; i++ {
		instructions[uint8(0xb8+i)] = movR32Imm32
	}
	instructions[uint8(0xeb)] = shortJump
}
