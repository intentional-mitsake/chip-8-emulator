package c8def

import (
	"fmt"
)

const (
	Scale = 10
	Green = "\033[32m"
	Reset = "\033[0m"
)

func SetupTerminal() {
	// Hide cursor: \033[?25l
	// Clear screen: \033[2J
	fmt.Print("\033[2J\033[?25l")
}

// Call this when the program exits
func CleanupTerminal() {
	// Show cursor: \033[?25h
	fmt.Print("\033[?25h")
}

func (c *Chip8) Render() {
	fmt.Print("\033[H") // move cursor to top-left
	for y := 0; y < 32; y++ {
		row := ""
		for x := 0; x < 64; x++ {
			if c.Display[y*64+x] == 1 {
				row += "██"
			} else {
				row += " "
			}
		}
		fmt.Println(row)
	}
}

func (c *Chip8) ClearDisplay() {
	for i := range c.Display {
		c.Display[i] = 0 // Set all pixels to off (0)
	}
}

func (c *Chip8) DrawSprites(xStart, yStart, height byte) {
	c.V[0xF] = 0
	c.DrawFlag = true
	//fmt.Print(spriteByte)
	for i := byte(0); i < height; i++ {
		//a spritebyte is 8 pixels wide(rows) always
		//the height of the sprite is diff tho, say heighgt is 4 then the spriteByte is like this:
		//[PIXEL] [PIXEL] [PIXEL] [PIXEL] [OFF] [OFF] [OFF] [OFF]-->1ST row--11110000
		//[PIXEL] [OFF] [OFF] [PIXEL] [PIXEL] [OFF] [PIXEL] [OFF]
		//[OFF] [PIXEL] [OFF] [OFF] [OFF] [PIXEL] [OFF] [PIXEL]
		//[PIXEL] [PIXEL] [OFF] [PIXEL] [OFF] [OFF] [PIEXEL] [OFF]-->4th row
		//hence we need to loop thru the entire heigth to get all rows(8pixels each)
		spriteByte := c.Memory[c.I+uint16(i)] //from 0-height
		for j := byte(0); j < 8; j++ {
			//check if the bit at j is set to 1
			//(0x80>>j) creates a mask by shifting to right by 8 bits: 10000000(j=0), 01000000(j=1), etc
			//doing AND with that mask and the current row of spriteByte we can get each bit
			//1st row--> 11110000 AND 100000(j=0)--> 10000000--> we got the first pixel of spriteBytes first row
			//looping thru 0-8 (j) we can get all 8 bits of each row of each spriteByte. thus we loop thrut the entire sprite
			//on the first row eg, result was 1, so we draw the pixel, for 0 we skip
			if (spriteByte & (0x80 >> j)) != 0 { //draw if 1
				targetX := (xStart + j) % DisplayWidth // x= (startPos + currPixel) mod 64 for wrap
				targetY := (yStart + i) % DisplayHeight
				index := uint16(targetY)*64 + uint16(targetX) //where on disp[64 * 32]

				//collision: if pixel at indx is already 1
				if c.Display[index] == 1 {
					c.V[0xF] = 1
				}
				//flip the pixel-XOR
				c.Display[index] ^= 1
			}
		}
	}
}
