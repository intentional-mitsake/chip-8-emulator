package c8def

import (
	"fmt"
)

func (c *Chip8) Render() {
	for y := 0; y < DisplayHeight; y++ {
		for x := 0; x < DisplayWidth; x++ {
			pixel := c.Display[y*DisplayWidth+x]
			if pixel == 1 {
				fmt.Print("█") // On pixel
			} else {
				fmt.Print(" ") // Off pixel
			}
		}
		fmt.Println() // New line after each row
	}
}

func (c *Chip8) ClearDisplay() {
	for i := range c.Display {
		c.Display[i] = 0 // Set all pixels to off (0)
	}
}

func (c *Chip8) DrawSprite(x, y byte, sprite []byte) {}
