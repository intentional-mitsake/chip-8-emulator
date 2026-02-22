package c8def

type InstructionSet struct {
	// Map instruction name -> base opcode
	Set map[string]uint16
}

func NewInstructionSet() *InstructionSet {
	return &InstructionSet{
		Set: map[string]uint16{
			"Clear":     0x00E0, // Clear the display
			"Return":    0x00EE, // Return from a subroutine
			"Jump":      0x1000, // Jump to address NNN
			"Call":      0x2000, // Call subroutine at NNN
			"SkipEq":    0x3000, // Skip next instruction if VX == NN
			"SkipNe":    0x4000, // Skip next instruction if VX != NN
			"SkipEqV":   0x5000, // Skip next instruction if VX == VY
			"Set":       0x6000, // Set VX = NN
			"Add":       0x7000, // Add NN to VX
			"Mov":       0x8000, // Set VX = VY
			"Or":        0x8001, // VX |= VY
			"And":       0x8002, // VX &= VY
			"Xor":       0x8003, // VX ^= VY
			"AddV":      0x8004, // VX += VY, VF = carry
			"Sub":       0x8005, // VX -= VY, VF = NOT borrow
			"Shr":       0x8006, // VX >>= 1, VF = least significant bit
			"SubN":      0x8007, // VX = VY - VX, VF = NOT borrow
			"Shl":       0x800E, // VX <<= 1, VF = most significant bit
			"SkipNeV":   0x9000, // Skip next instruction if VX != VY
			"SetI":      0xA000, // Set I = NNN
			"Jump0":     0xB000, // Jump to NNN + V0
			"Rand":      0xC000, // VX = random byte AND NN
			"Sprite":    0xD000, // Draw sprite at (VX,VY), VF = collision
			"KeyEq":     0xE09E, // Skip next instruction if key VX pressed
			"KeyNe":     0xE0A1, // Skip next instruction if key VX not pressed
			"GetDelay":  0xF007, // Set VX = delay timer value
			"WaitKey":   0xF00A, // Wait for keypress, store in VX
			"SetDelay":  0xF015, // Set delay timer = VX
			"SetBuzzer": 0xF018, // Set sound timer = VX
			"AddI":      0xF01E, // I += VX
			"Hex":       0xF029, // Set I = location of sprite for digit VX
			"Bcd":       0xF033, // Store BCD of VX in memory[I..I+2]
			"Save":      0xF055, // Store V0..VX in memory starting at I
			"Load":      0xF065, // Read V0..VX from memory starting at I
		},
	}
}
