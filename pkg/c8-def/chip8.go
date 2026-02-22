package c8def

import "os"

const (
	MemSize       = 4096 // 4KB memory
	RegCount      = 16   // 16 registers (V0 to VF)
	DisplayWidth  = 64
	DisplayHeight = 32
)

//btw to clear confusion later on
//uint8 and byte are the same thing in GO
//this file will define the Chip8 struct and its initialization, as well as loading the fontset and ROM into memory

type Chip8 struct {
	Memory [MemSize]byte // 4KB memory

	V  [RegCount]byte // 16 8-bit registers--> V0 to VF (VF is often used as a flag)
	I  uint16         // Index register-->Used for memory operations
	PC uint16         // Program counter-->Points to the current instruction in memory

	Stack [16]uint16 // Stack for addresses--> Used for subroutine/function calls
	SP    uint16     // Stack pointer--> Points to the top of the stack

	DT byte // Delay timer--> Decrements at a rate of 60Hz
	ST byte // Sound timer--> produces a sound when >0

	Keypad  [16]bool                           // Hex-Keypad(0-F) state--> Represents the state of the 16 keys
	Display [DisplayWidth * DisplayHeight]byte // Display state (64x32 pixels)--> Each pixel can be on (1) or off (0)S
}

func NewChip8() *Chip8 {
	return &Chip8{
		//from 0x200 to 0xFFF is used for program and data storage
		//that is from 512 to 4095 in decimalS
		//so we initialize the program counter to 0x200 (512) to point to the start of the program memory
		//from 0x4f to 0x200 is used for C8-VM(79 to 512 in decimal)
		//finally from 0x00 to 0x4F (0 to 79) is reserved for the fontset
		//fontset size is 16 * 5 = 80 so 0-79
		PC: 0x200,
		I:  0x000,
		SP: 0x00,
	}
}

func (c *Chip8) LoadFontset() {
	//chip8 has 16 characters(0-F) in its fontset
	//each char is rep by a 5 byte sprite(5 rows of 8 bits)
	//hence the size of the fontset is 16 chars * 5 bytes/char = 80 bytes
	//each row below is a single char from 0 to F, rep by the 5 bytes
	fontset := [5 * 16]byte{
		//to understand how these bytes rep the chars,
		// we can convert each byte to binary and visualize it
		// for example, the char '0' is rep by the bytes 0xF0, 0x90, 0x90, 0x90, 0xF0
		// in binary, these bytes are:
		// 0xF0 = 11110000	->firsst row of the char '0'
		// 0x90 = 10010000
		// 0x90 = 10010000
		// 0x90 = 10010000
		// 0xF0 = 11110000  _>fifth and last row of the char '0'
		//visualizing this, 1--> * and 0--> space;
		// **** // -> 0xF0(first row)
		// *  * // -> 0x90(second row)
		// *  * // -> 0x90(third row)
		// *  * // -> 0x90(fourth row)
		// **** // -> 0xF0(fifth row)
		//thus 0
		0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
		0x20, 0x60, 0x20, 0x20, 0x70, // 1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
		0x90, 0x90, 0xF0, 0x10, 0x10, // 4
		0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
		0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
		0xF0, 0x10, 0x20, 0x40, 0x40, // 7
		0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
		0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
		0xF0, 0x90, 0xF0, 0x90, 0x90, // A
		0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
		0xF0, 0x80, 0x80, 0x80, 0xF0, // C
		0xE0, 0x90, 0x90, 0x90, 0xE0, // D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
		0xF0, 0x80, 0xF0, 0x80, 0x80, // F
	}
	//load the fontset(char thru 0-F) into the first 80 bytes of mem(0x00 to 0x4F or 0 to 79)
	for i := 0; i < len(fontset); i++ {
		c.Memory[i] = fontset[i]
	}
}

func (c *Chip8) LoadROM(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	for i := 0; i < len(data); i++ {
		//from 0x200 to 0xFFF is used for program data and code
		//ROM is loaded starting at 0x200 (512 in decimal)
		c.Memory[0x200+i] = data[i]
		//printing for debugging purposes
		//fmt.Printf("%d", data[i])
	}
	return nil
}
