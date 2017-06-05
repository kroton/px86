package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: px86 filename")
		return
	}

	fileName := os.Args[1]

	f, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%sファイルが開けません: %v", fileName, err)
		return
	}
	defer f.Close()

	emu := newEmulator()
	if err := emu.load(f); err != nil {
		fmt.Fprintf(os.Stderr, "%sファイルが読み込めません: %v", fileName, err)
		return
	}

	emu.eval()
}
