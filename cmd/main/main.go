package main

import (
	def "c8-emulation/pkg/c8-def"
	"fmt"
)

// MEMORY MAP for the CHIP-8 system:
// 0-79 --> fontset(0-F chars)--> load after initializing the Chip8 struct
// 80-511 --> Interpreter/Resrved
// 512-4095 --> Program/ROM/Stack/General purpose registers/Index register/Timers/Keypad/Display

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
	//fmt.Printf("Memory after loading ROM: %v\n", chip8.Memory)
	//fetchedInst := chip8.Fetch()
	//b2 gets lower 8 bits--> 0x1234--> b2=0x34 and b1 gets higher 8 bits--> 0x12
	//b1 := byte(fetchedInst)
	//shifts to right by 8 bits--> 0x1234 >> 8 = 0x12
	//b2 := byte(fetchedInst >> 8) //get second inst byte (higher 8 bits)
	//fmt.Printf("Fetched instruction(In HexCode): 0x%X\nIn Bytes: %v %v\n", fetchedInst, b1, b2)
	chip8.Render() //render the display (for debugging purposes)
}
