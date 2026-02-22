package c8def

import (
	"math/rand"
)

func RandNumGen() byte {
	return byte(rand.Intn(256)) // Gen a random byte (0-255)
	//need this for the CXNN instruction
	//used for enemies in games to move randomly or for random events
	//doesnt need to be truly random, just good enough for games
	//this should be gud enough
}

// read the inst at current PC, return it as a 16-bit value and incr the PC by 2
func (c *Chip8) Fetch() uint16 {
	//ROM is stored starting at 512, PC is initialize at 512 as well
	//so we start reading immediately from the start of the ROM
	//like the fontset, each inst is made up of multiple bytes(5 for fonts, 2 for instructions)
	//each inst is 16 bits long, stored in 2 consecutive bytes in mem, so read 2 bytes and combine them to get the full inst
	//each byte is 8 bits
	curr := c.PC //get the current PC value(512 at start)
	//basically: (1st byte) OR (2nd byte) to combine them into a single 16-bit value
	//0x12 and 0x34 are the two inst, if we do A | B only, it does basic OR operation
	//00010010(0x12) OR 00110100(0x34) = 00110110(0x36) --> still 8 bits and not the full combined 16-bit inst we need
	//use <<8, it shifts 0x12 to left by 8 bits, making it 0x1200(00010010 00000000 in binary)
	//now when we do A | B ---> 0x1200 OR 0x3:
	//00010010 00000000 (0x1200) OR 00110100(0x34) = 00010010 00110100(0x1234)
	//this way we combine the two bytes to form the full inst
	inst := uint16(c.Memory[curr])<<8 | uint16(c.Memory[curr+1])
	c.PC += 2 //incr the PC by 2 to point to the next inst
	return inst
}
