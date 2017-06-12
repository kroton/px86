package main

import (
	"fmt"
	"io"
	"os"
)

type codeNotImplemented uint8

func (c codeNotImplemented) Error() string {
	return fmt.Sprintf("Not Implemented: %x", uint8(c))
}

type emulator struct {
	state        *state
	instructions map[uint8]instruction
}

func newEmulator() *emulator {
	return &emulator{
		state:        newState(1024*1024, 0x7c00, 0x7c00, 0x7c00),
		instructions: instructions,
	}
}

func (e *emulator) load(r io.Reader) error {
	_, err := io.Copy(e.state, r)
	return err
}

func (e *emulator) canStep() bool {
	return 0 <= e.state.eip && e.state.eip < len(e.state.memory)
}

func (e *emulator) step() error {
	code, err := e.state.getUint8(0)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "EIP = %x, Code = %02x\n", e.state.eip, code)

	ins, ok := e.instructions[code]
	if !ok {
		return codeNotImplemented(code)
	}
	return ins(e.state)
}

func (e *emulator) isEnd() bool {
	return e.state.eip == 0
}

func (e *emulator) eval() {
	for e.canStep() {
		if err := e.step(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			break
		}

		if e.isEnd() {
			fmt.Fprintf(os.Stderr, "end of program.\n\n")
			break
		}
	}
	e.state.dumpRegisters()
}
