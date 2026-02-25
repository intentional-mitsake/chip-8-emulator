package main

import (
	"c8-emulation/pkg/game"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 64 * 10 // scale each CHIP-8 pixel by 10
	screenHeight = 32 * 10
)

// MEMORY MAP for the CHIP-8 system:
// 0-79 --> fontset(0-F chars)--> load after initializing the Chip8 struct
// 80-511 --> Interpreter/Resrved
// 512-4095 --> Program/ROM/Stack/General purpose registers/Index register/Timers/Keypad/Display

func main() {
	//fmt.Printf("Memory after loading ROM: %v\n", chip8.Memory)
	//fetchedInst := chip8.Fetch()
	//b2 gets lower 8 bits--> 0x1234--> b2=0x34 and b1 gets higher 8 bits--> 0x12
	//b1 := byte(fetchedInst)
	//shifts to right by 8 bits--> 0x1234 >> 8 = 0x12
	//b2 := byte(fetchedInst >> 8) //get second inst byte (higher 8 bits)
	//fmt.Printf("Fetched instruction(In HexCode): 0x%X\nIn Bytes: %v %v\n", fetchedInst, b1, b2)
	fmt.Print("Starting Emulator...\n")
	//disp
	game := game.NewGame()
	game.Chip8.LoadFontset() //load the fontset first before loading the ROM
	//%v prints value in a default format
	//%+v prints struct with field names
	//fmt.Printf("Initialized Chip8: %+v\n", chip8)
	/*
		romPath := "./pkg/assets/INVADERS"
		err := game.Chip8.LoadROM(romPath)
		if err != nil {
			fmt.Printf("Error loading ROM: %v\n", err)
		}
	*/
	//this function starts the game and handles the game loop
	ebiten.SetWindowSize(1000, 640) //to get a larger disp

	ebiten.RunGame(game)

}
