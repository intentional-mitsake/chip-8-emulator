package main

import (
	def "c8-emulation/pkg/c8-def"
	"fmt"
)

func main() {
	chip8 := def.NewChip8()
	//fmt.Printf("Initialized Chip8: %+v\n", chip8)
	romPath := "./pkg/assets/PONG"
	err := chip8.LoadROM(romPath)
	if err != nil {
		fmt.Printf("Error loading ROM: %v\n", err)
	}
	fmt.Printf("Memory after loading ROM: %v\n", chip8.Memory)
}
