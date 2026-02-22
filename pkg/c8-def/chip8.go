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
