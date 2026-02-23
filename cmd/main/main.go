package main

import (
	def "c8-emulation/pkg/c8-def"
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 64 * 10 // scale each CHIP-8 pixel by 10
	screenHeight = 32 * 10
)

type Game struct {
	Chip8 *def.Chip8
	Scale float32
}

func (g *Game) Update() error {
	// Run CPU cycles
	g.Chip8.CoreLoop()
	g.Chip8.UpdateTimers()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw the CHIP-8 display
	colOn := color.White
	colOff := color.Black

	for y := 0; y < 32; y++ {
		for x := 0; x < 64; x++ {
			c := colOff
			if g.Chip8.Display[y*64+x] != 0 {
				c = colOn
			}
			vector.FillRect(screen, float32(x)*g.Scale, float32(y)*g.Scale, g.Scale, g.Scale, c, false)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 64 * int(g.Scale), 32 * int(g.Scale)
}

// MEMORY MAP for the CHIP-8 system:
// 0-79 --> fontset(0-F chars)--> load after initializing the Chip8 struct
// 80-511 --> Interpreter/Resrved
// 512-4095 --> Program/ROM/Stack/General purpose registers/Index register/Timers/Keypad/Display

func main() {
	def.SetupTerminal()
	defer def.CleanupTerminal()
	chip8 := def.NewChip8()
	chip8.LoadFontset() //load the fontset first before loading the ROM
	//%v prints value in a default format
	//%+v prints struct with field names
	//fmt.Printf("Initialized Chip8: %+v\n", chip8)
	romPath := "./pkg/assets/PONG"
	err := chip8.LoadROM(romPath)
	if err != nil {
		fmt.Printf("Error loading ROM: %v\n", err)
	}
	//fmt.Printf("Memory after loading ROM: %v\n", chip8.Memory)
	//fetchedInst := chip8.Fetch()
	//b2 gets lower 8 bits--> 0x1234--> b2=0x34 and b1 gets higher 8 bits--> 0x12
	//b1 := byte(fetchedInst)
	//shifts to right by 8 bits--> 0x1234 >> 8 = 0x12
	//b2 := byte(fetchedInst >> 8) //get second inst byte (higher 8 bits)
	//fmt.Printf("Fetched instruction(In HexCode): 0x%X\nIn Bytes: %v %v\n", fetchedInst, b1, b2)
	fmt.Print("Starting Emulator...\n")
	//disp code is gpt-->i hate graphics prog, fuck sdl2 and ebiton
	//for gameboy will do myself(hope)
	game := &Game{
		Chip8: chip8,
		Scale: 10,
	}

	ebiten.SetWindowSize(64*10, 32*10)
	ebiten.SetWindowTitle("CHIP-8 Emulator")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
