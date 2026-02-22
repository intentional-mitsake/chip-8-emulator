package c8def

import (
	"fmt"
	"strings"
)

const (
	Scale = 10
	Green = "\033[32m"
	Reset = "\033[0m"
)

func (c *Chip8) RenderTerminal() {
	var builder strings.Builder
	for y := 0; y < 32; y++ {
		for x := 0; x < 64; x++ {
			if c.Display[y*64+x] == 1 {
				builder.WriteString(Green + "█" + Reset)
			} else {
				builder.WriteString(" ")
			}
		}
		builder.WriteString("\n")
	}
	fmt.Print(builder.String())
}

func (c *Chip8) ClearDisplay() {
	for i := range c.Display {
		c.Display[i] = 0 // Set all pixels to off (0)
	}
}

func (c *Chip8) DrawSprites(x, y byte, N byte) {
	for i := 0; i < int(N); i++ {

	}
}
