package main

import (
	def "c8-emulation/pkg/c8-def"
	"fmt"
)

func main() {
	chip8 := def.NewChip8()
	fmt.Printf("Initialized Chip8: %+v\n", chip8)
}
