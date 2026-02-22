package main

import (
	def "c8-emulation/pkg/c8-def"
	"fmt"
)

//MEMORY MAP for the CHIP-8 system:
//0-79 --> fontset(0-F chars)--> load after initializing the Chip8 struct
//80-511 --> Interpreter/Resrved
//512-4095 --> Program/ROM/Stack/General purpose registers/Index register/Timers/Keypad/Display

func main() {
	chip8 := def.NewChip8()
	chip8.LoadFontset() //load the fontset first before loading the ROM
	//%v prints value in a default format
	//%+v prints struct with field names
	//fmt.Printf("Initialized Chip8: %+v\n", chip8)
	romPath := "./pkg/assets/PONG"
	err := chip8.LoadROM(romPath)
	if err != nil {
		fmt.Printf("Error loading ROM: %v\n", err)
	}
	fmt.Printf("Memory after loading ROM: %v\n", chip8.Memory)
}
