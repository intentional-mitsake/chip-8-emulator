package c8def

type Chip8 struct {
	Memory [4096]byte // 4KB memory

	V  [16]byte // 16 8-bit registers
	I  uint16   // Index register
	PC uint16   // Program counter

	Stack [16]uint16 // Stack for addresses
	SP    uint16     // Stack pointer

	DT byte // Delay timer
	ST byte // Sound timer

	Keypad  [16]byte      // Keypad state
	Display [64 * 32]byte // Display state (64x32 pixels)
}
