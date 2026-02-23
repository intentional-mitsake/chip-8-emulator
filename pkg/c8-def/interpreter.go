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
	X := (inst & 0x0F00) >> 8 // in second byte usually, shift right by 8 to get the value of X as a single byte (0-15)
	Y := (inst & 0x00F0) >> 4 // in third byte usually; we get smth like 0x00A0, >>4 gives 0x000A which is just 0x0A
	N := (inst & 0x000F)      // last 4 bits--from here on out we need exact 4,8,12 bits so no shifting
	NN := (inst & 0x00FF)     // last 8 bits
	NNN := (inst & 0x0FFF)    // last 12 bits
	InstSet := NewInstructionSet()
	switch opcode {
	case 0x0000:
		//directly test for the two 00E0 and 00EE instructions since they both have the same opcode
		if inst == InstSet.Set["Clear"] {
			c.ClearDisplay()
		} else if inst == InstSet.Set["Return"] {
			//return from subroutine: remove the last addr from the stack, set PC to that addr
			c.Stack[c.SP] = 0    //clear the top of the stack
			c.SP--               //move stack pointer down
			c.PC = c.Stack[c.SP] //set PC to the addr at the top of the stack
		}
	case 0x1000:
		//JUMP
		c.PC = NNN //set the PC to the address NNN to jump there
	case 0x2000:
		//CALL: push curr PC onto the stack, incr SP, then set PC to NNN to jump to the subroutine
		c.SP++               //incr stack pointer to point to the next empty slot
		c.Stack[c.SP] = c.PC //push current PC onto the stack and incr SP
		c.PC = NNN           //set PC to NNN to jump to the subroutine
	case 0x3000:
		//SKIPEQL: skip next inst if VX == NN
		if c.V[X] == byte(NN) {
			c.PC += 2 //skip next inst by incr PC by 2
		}

	case 0x4000:
		//SKIPNE: skip next inst if VX != NN
		if c.V[X] != byte(NN) {
			c.PC += 2
		}

	case 0x5000:
		//SKIPEQV: skip next inst if VX == VY
		if c.V[X] == c.V[Y] {
			c.PC += 2
		}
	case 0x6000:
		//SET: set VX to NN
		c.V[X] = byte(NN)
	case 0x7000:
		//ADD
		c.V[X] += byte(NN) //add NN to VX, no need to worry about overflow since it wraps around in byte
	case 0x8000:
		//compare op8 value with diff 8000 instructions to determine which one it is
		switch op8 {
		case InstSet.Set["Mov"]:
			//MOV: set VX = VY
			c.V[X] = c.V[Y]
		case InstSet.Set["Or"]:
			//OR: VX = VX OR VY
			c.V[X] |= c.V[Y]
		case InstSet.Set["And"]:
			//AND: bitwise AND, VX = VX AND VY
			c.V[X] &= c.V[Y]
		case InstSet.Set["Xor"]:
			//XOR
			c.V[X] ^= c.V[Y]
		case InstSet.Set["AddV"]:
			//ADDV
			sum := uint16(c.V[X]) + uint16(c.V[Y]) //calculate the sum as uint16 to check for overflow
			c.V[0xF] = 0                           //reset VF before addition
			c.V[X] = byte(sum)                     //set VX to the sum (wraps around if >255)
			if sum > 255 {
				c.V[X] = byte(sum % 256) //wrap around if overflow occurs
				c.V[0xF] = 1             //set VF to 1 to indicate carry
			}
		case InstSet.Set["Sub"]:
			//SUB
			diff := int16(c.V[X]) - int16(c.V[Y]) //calculate the difference as int16(need signed int) to check for borrow
			if diff < 0 {
				c.V[X] = byte((diff + 256) % 256) //wrap around if borrow occurs
				c.V[0xF] = 0                      //set VF to 0 to indicate borrow
			} else {
				c.V[X] = byte(diff) //no borrow, just set VX to the difference
				c.V[0xF] = 1        //set VF to 1 to indicate no borrow
			}
		case InstSet.Set["Shr"]:
			//SHR: if the bit that is shifted out is 1, set VF to 1, else set VF to 0, then shift VX right by 1 bit
			c.V[0xF] = c.V[X] & 0x1 //the bit that is shifted out is the last bit of VX, so we AND it with 0x1 to get that bit and set VF accordingly
			c.V[X] >>= 1            //shift VX right by 1 bit
		case InstSet.Set["SubN"]:
			//SUBN
			diff := int16(c.V[Y]) - int16(c.V[X]) //calculate the difference as int16(need signed int) to check for borrow
			if diff < 0 {
				c.V[X] = byte((diff + 256) % 256) //wrap around if borrow occurs
				c.V[0xF] = 0                      //set VF to 0 to indicate borrow
			} else {
				c.V[X] = byte(diff) //no borrow, just set VX to the difference
				c.V[0xF] = 1        //set VF to 1 to indicate no borrow
			}
		case InstSet.Set["Shl"]:
			//SHL: shift vx left by 1 bit, if the bit that is shifted out is 1, set VF to 1, else set VF to 0
			c.V[0xF] = (c.V[X] & 0x80)
			c.V[X] <<= 1
		}
	case 0x9000:
		//SKIPNEV: skip next inst if VX != VY
		if c.V[X] != c.V[Y] {
			c.PC += 2
		}
	case 0xA000:
		//SETI: set I = NNN
		c.I = NNN
	case 0xB000:
		//JUMP0
		c.PC = NNN + uint16(c.V[0]) //jump to address NNN + V0
	case 0xC000:
		//RAND
		randByte := RandNumGen()     //generate a random byte
		c.V[X] = randByte & byte(NN) //set VX to the result of random byte AND NN
	case 0xD000:
		//SPRITE
		x := c.V[X] % DisplayWidth //get the x coordinate from VX, wrap around if it exceeds display width
		y := c.V[Y] % DisplayHeight
		c.DrawSprites(x, y, byte(N))
	case 0xE000:
		//compare with E09E and E0A1 to determine which one it is
		if inst == InstSet.Set["KeyEq"] {
			//KEYEQ: skip next inst if key with the value of VX is pressed,
			//we can represent the state of the keys in the Keypad array, where each index corresponds to a key (0-F) and the value is 1 if pressed and 0 if not
			key := c.V[X] //get the value of VX to determine which key we are checking
			if c.Keypad[key] {
				c.PC += 2 //skip next inst if the key is pressed
			}
		} else if inst == InstSet.Set["KeyNe"] {
			//KEYNE: skip next inst if key with the value of VX is not pressed
			key := c.V[X]
			if !c.Keypad[key] {
				c.PC += 2 //skip next inst if the key is not pressed
			}
		}
	case 0xF000:
		switch opF {
		case InstSet.Set["GetDelay"]:
			//GETDELAY:
			c.V[X] = c.DT
		case InstSet.Set["WaitKey"]:
			//WAITKEY: wait for a key press and store the value of the key in VX, we can check the Keypad array for any key that is currently pressed, if multiple keys are pressed we can just take the first one we find
			for i, pressed := range c.Keypad {
				if pressed {
					c.V[X] = byte(i) //store the value of the key in VX
					break            //exit the loop after finding the first pressed key
				}
			}

		case InstSet.Set["SetDelay"]:
			//SETDELAY
			c.DT = c.V[X] //set the delay timer to the value of VX
		case InstSet.Set["SetBuzzer"]:
			//SETBUZZER
			c.ST = c.V[X] //set the sound timer to the value of VX
		case InstSet.Set["AddI"]:
			//ADDI
			c.I += uint16(c.V[X]) //add the value of VX to I
		case InstSet.Set["Hex"]:
			//HEX
			hexV := c.V[X] & 0x0F //we only need the last nibble(4bits) of vx for this
			//fontset is stored in 0-79(80 bytes for 16 chars * 5 bytes each)
			//say hexV is 0x2, we want to set I to the location of the sprite for the char '2' in the fontset,
			// which starts at index 2*5=10 (since each char takes 5 bytes)
			c.I = uint16(hexV) * 5 //set I to the location of the sprite for the char in the fontset
		case InstSet.Set["Bcd"]:
			//BCD: store bcd rep of num in VX in memory locations I, I+1, and I+2
			value := uint16(c.V[X])                   //get the value of VX to convert to BCD
			c.Memory[c.I] = byte(value / 100)         //hundreds digit
			c.Memory[c.I+1] = byte((value / 10) % 10) //tens digit
			c.Memory[c.I+2] = byte(value % 10)        //ones digit
			//9C--> 156 in decimal--> 1 in hundreds place, 5 in tens place, 6 in ones place
			//so we store 1 at memory[I], 5 at memory[I+1], and 6 at memory[I+2]
		case InstSet.Set["Save"]:
			//SAVE: store V0 to VX in memory starting at address I
			for i := 0; i <= int(X); i++ {
				c.Memory[c.I+uint16(i)] = c.V[i] //store the value of each register from V0 to VX in memory starting at address I
			}
		case InstSet.Set["Load"]:
			//LOAD
			for i := 0; i <= int(X); i++ {
				c.V[i] = c.Memory[c.I+uint16(i)] //load the value of each register from memory starting at address I
			}
		}

	}

}

func (c *Chip8) CoreLoop() {
	inst := c.Fetch()
	c.DecodeAndExc(inst)
}
