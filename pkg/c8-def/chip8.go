package c8def

const (
	MemSize       = 4096 // 4KB memory
	RegCount      = 16   // 16 registers (V0 to VF)
	DisplayWidth  = 64
	DisplayHeight = 32
)

type Chip8 struct {
	Memory [MemSize]byte // 4KB memory

	V  [RegCount]byte // 16 8-bit registers--> V0 to VF (VF is often used as a flag)
	I  uint16         // Index register-->Used for memory operations
	PC uint16         // Program counter-->Points to the current instruction in memory

	Stack [16]uint16 // Stack for addresses--> Used for subroutine/function calls
	SP    uint16     // Stack pointer--> Points to the top of the stack

	DT byte // Delay timer--> Decrements at a rate of 60Hz
	ST byte // Sound timer--> produces a sound when >0

	Keypad  [16]byte                           // Hex-Keypad(0-F) state--> Represents the state of the 16 keys
	Display [DisplayWidth * DisplayHeight]byte // Display state (64x32 pixels)--> Each pixel can be on (1) or off (0)S
}

func NewChip8() *Chip8 {
	return &Chip8{
		PC: 0x200, // Programs start at memory location 0x200
	}
}

func (c *Chip8) LoadFontset() {
	fontset := [5 * 16]byte{
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
	for i := 0; i < len(fontset); i++ {
		c.Memory[i] = fontset[i]
	}
}
