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

func (c *Chip8) DecodeAndExc(inst uint16) {
	//16 bit inst--> 0x1NNN, 0x6XNN, etc
	//to decode, we need to extract the different parts of the inst based on the instruction format
	//in 1NNN, first 4 bits(nibble) are the opcode(1), next 12 bits are the address(NNN)
	//so 1 --> opcode, NNN --> address, JUMP to address NNN
	opcode := inst & 0xF000 //get first 4 bits for opcode--> 0x1NNN AND 0xF000 = 0x1000
	op8 := inst & 0xF00F    //get first 4 bits and last 4 bits for 8XY0, 8XY1, etc--> 0x8XY0 AND 0xF00F = 0x8000
	opF := inst & 0xF0FF    //for 0xF000 instructions where last 8 bits vary

	//X register, Y register, 4th bit(some inst need it), last byte(NN) and last 12 bits(NNN)
	X := (inst & 0x0F00)   // in second byte usually
	Y := (inst & 0x00F0)   // in third byte usually
	N := (inst & 0x000F)   // last 4 bits
	NN := (inst & 0x00FF)  // last 8 bits
	NNN := (inst & 0x0FFF) // last 12 bits
	InstSet := NewInstructionSet()
	switch opcode {
	case 0x0000:
		//directly test for the two 00E0 and 00EE instructions since they both have the same opcode
		if inst == InstSet.Set["Clear"] {
			c.ClearDisplay()
		} else if inst == InstSet.Set["Return"] {
			//return code here
		}
	case 0x1000:
		//JUMP
	case 0x2000:
		//CALL
	case 0x3000:
		//SKIPEQ
	case 0x4000:
		//SKIPNE
	case 0x5000:
		//SKIPEQV
	case 0x6000:
		//SET
	case 0x7000:
		//ADD
	case 0x8000:
		//compare op8 value with diff 8000 instructions to determine which one it is
		switch op8 {
		case InstSet.Set["Mov"]:
			//MOV
		case InstSet.Set["Or"]:
			//OR
		case InstSet.Set["And"]:
			//AND
		case InstSet.Set["Xor"]:
			//XOR
		case InstSet.Set["AddV"]:
			//ADDV
		case InstSet.Set["Sub"]:
			//SUB
		case InstSet.Set["Shr"]:
			//SHR
		case InstSet.Set["SubN"]:
			//SUBN
		case InstSet.Set["Shl"]:
			//SHL
		}
	case 0x9000:
		//SKIPNEV
	case 0xA000:
		//SETI
	case 0xB000:
		//JUMP0
	case 0xC000:
		//RAND
	case 0xD000:
		//SPRITE
	case 0xE000:
		//compare with E09E and E0A1 to determine which one it is
		if inst == InstSet.Set["KeyEq"] {
			//KEYEQ
		} else if inst == InstSet.Set["KeyNe"] {
			//KEYNE
		}
	case 0xF000:
		switch opF {
		case InstSet.Set["GetDelay"]:
			//GETDELAY
		case InstSet.Set["WaitKey"]:
			//WAITKEY
		case InstSet.Set["SetDelay"]:
			//SETDELAY
		case InstSet.Set["SetBuzzer"]:
			//SETBUZZER
		case InstSet.Set["AddI"]:
			//ADDI
		case InstSet.Set["Hex"]:
			//HEX
		case InstSet.Set["Bcd"]:
			//BCD
		case InstSet.Set["Save"]:
			//SAVE
		case InstSet.Set["Load"]:
			//LOAD
		}

	}

}
