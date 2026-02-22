package c8def

type InstructionSet struct {
	//map for key->vlaue pairs of instruction name
	Set map[string]string //All 34 instructions in the CHIP-8 instruction set
}

func NewInstructionSet() *InstructionSet {
	return &InstructionSet{
		Set: map[string]string{
			"Clear":     "00E0", // Clear the display
			"Return":    "00EE", // Return from a subroutine
			"Jump":      "1NNN", // Jump to address NNN
			"Call":      "2NNN", // Call subroutine at NNN
			"SkipEq":    "3XNN", // Skip next instruction if VX == NN
			"SkipNe":    "4XNN", // Skip next instruction if VX != NN
			"SkipEqV":   "5XY0", // Skip next instruction if VX == VY
			"Set":       "6XNN", // Set VX = NN
			"Add":       "7XNN", // Add NN to VX
			"Mov":       "8XY0", // Set VX = VY
			"Or":        "8XY1", // VX |= VY
			"And":       "8XY2", // VX &= VY
			"Xor":       "8XY3", // VX ^= VY
			"AddV":      "8XY4", // VX += VY, VF = carry
			"Sub":       "8XY5", // VX -= VY, VF = NOT borrow
			"Shr":       "8XY6", // VX >>= 1, VF = least significant bit
			"SubN":      "8XY7", // VX = VY - VX, VF = NOT borrow
			"Shl":       "8XYE", // VX <<= 1, VF = most significant bit
			"SkipNeV":   "9XY0", // Skip next instruction if VX != VY
			"SetI":      "ANNN", // Set I = NNN
			"Jump0":     "BNNN", // Jump to NNN + V0
			"Rand":      "CXNN", // VX = random byte AND NN
			"Sprite":    "DXYN", // Draw sprite at (VX,VY), VF = collision
			"KeyEq":     "EX9E", // Skip next instruction if key VX pressed
			"KeyNe":     "EXA1", // Skip next instruction if key VX not pressed
			"GetDelay":  "FX07", // Set VX = delay timer value
			"WaitKey":   "FX0A", // Wait for keypress, store in VX
			"SetDelay":  "FX15", // Set delay timer = VX
			"SetBuzzer": "FX18", // Set sound timer = VX
			"AddI":      "FX1E", // I += VX
			"Hex":       "FX29", // Set I = location of sprite for digit VX
			"Bcd":       "FX33", // Store BCD of VX in memory[I..I+2]
			"Save":      "FX55", // Store V0..VX in memory starting at I
			"Load":      "FX65", // Read V0..VX from memory starting at I
		},
	}
}
