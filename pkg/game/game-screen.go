package game

import (
	def "c8-emulation/pkg/c8-def"
	"fmt"
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

const (
	DisplayWidth  = 64
	DisplayHeight = 32
	visible       = 10
)

var KEYPADCOLOR = color.RGBA{55, 65, 50, 255}
var INACTIVEKEY = color.RGBA{45, 50, 55, 255}
var ACTIVEKEY = color.RGBA{255, 204, 0, 255}
var STATPADCOLOR = color.RGBA{10, 15, 20, 255}
var STATSCOLOR = color.RGBA{0, 255, 255, 255}
var BORDERCOLOR = color.RGBA{0, 40, 60, 255}
var BORDEREDGE = color.RGBA{255, 180, 6, 180}

var offset int = 0 //indx of first visibile rom

type Game struct {
	//handling game state
	GameState bool     //running(true) or paused(false)
	Roms      []string //arrat of all available roms
	Selected  int      //index of currently selected rom
	//Chip-8 def and ebiten image
	Chip8   *def.Chip8    //first we need to create a chip8 struct
	display *ebiten.Image // then a display image diff from the disp attr of chip8

	//for menu
	FontFace text.Face
	DrawOpts *text.DrawOptions
}

func NewGame() *Game {
	c8 := def.NewChip8()
	roms := GetRoms()
	settings := &text.DrawOptions{}
	settings.GeoM.Translate(0, 0)
	return &Game{
		GameState: false, //not running at start
		Roms:      roms,
		Chip8:     c8,
		Selected:  0, //first rom selected at start menu
		display: ebiten.NewImage(
			DisplayWidth,
			DisplayHeight,
		),
		FontFace: text.NewGoXFace(basicfont.Face7x13),
		DrawOpts: settings,
	}
}

func GetRoms() []string {
	var roms []string
	files, err := os.ReadDir("./pkg/assets")
	if err != nil {
		fmt.Println(err)
	}
	for _, file := range files {
		//if its a dir not file, skip
		if file.IsDir() {
			continue
		}
		//fmt.Print(file.Name())
		roms = append(roms, file.Name())
	}
	return roms
}

func (g *Game) rebuildDisplay() {
	//this is callded every frame to redraw the display
	pixels := make([]byte, DisplayWidth*DisplayHeight*4) // RGBA

	for i := 0; i < DisplayWidth*DisplayHeight; i++ {
		if g.Chip8.Display[i] == 1 {
			//		pixels are this color
			pixels[i*4+0] = 155 // R
			pixels[i*4+1] = 188 // G
			pixels[i*4+2] = 015 // B
			pixels[i*4+3] = 255 // A
		} else {
			//background is this color
			pixels[i*4+0] = 15
			pixels[i*4+1] = 25
			pixels[i*4+2] = 15
			pixels[i*4+3] = 255
		}
	}
	//ebiten creates a window
	//we draw g.display over it
	//here we aare replacing individual pixels on the dispaly image we drew
	g.display.ReplacePixels(pixels)
}

func (g *Game) Update() error {
	switch g.GameState {
	case true:
		//running
		for i := 0; i < 10; i++ {
			//reason for 10 is: chip8 ran around 500-1000 inst per second
			g.Chip8.CoreLoop()
		}
		//timers
		if g.Chip8.DT > 0 {
			g.Chip8.DT--
		}
		if g.Chip8.ST > 0 {
			g.Chip8.ST--
			//sound trigger here
		}
		//draw screen
		if g.Chip8.DrawFlag {
			g.rebuildDisplay()
			g.Chip8.DrawFlag = false
		}

	case false:
		//not running - menu
		g.Menu()
	}
	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {
	switch g.GameState {
	case true:
		//running
		//screen.Fill(color.RGBA{15, 30, 15, 255})
		scaling := &ebiten.DrawImageOptions{}
		//same to make pixels square, if diff, rectangle
		//pixels is waht thsi effects nto display/screen size
		scaling.GeoM.Scale(10, 10) //this scaling effects whats drawn on teh screen
		//filter det how image is scaled
		//nearest looked best(pixelated was similar), linear was worst
		scaling.Filter = ebiten.FilterNearest
		//fading effect on pixels
		scaling.ColorScale.SetA(1.0)
		//top-left
		scaling.GeoM.Translate(0, 0)
		screen.DrawImage(g.display, scaling) //takes img as arg and draws it on the screen image

		//SUBSCreen
		g.DrawSubscreen(screen)
	case false:
		//not running
		screen.Fill(color.RGBA{15, 20, 30, 255})
		g.DrawMenu(screen)
		g.DrawSubscreen(screen)
	}
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	//this scaling defines the size of the window nto the size of the piexls on screen
	//this way we can det the size of the window
	return outsideWidth, outsideHeight
}

func (g *Game) Menu() {
	//if arrow down, select next rom
	//if we have reached the last rom and user presses arrow down, go back to first
	//prev used ebiten.IsKeyPressed, but game runs at 60fps and holding key for even half secnod \
	//triggers g.Selected++ multiple times, inpututil.IsKeyJustPressed only returns true on first frame of input
	//so tirgger onl =y once
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		g.Selected++
		if g.Selected >= len(g.Roms) {
			g.Selected = 0
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		g.Selected--
		if g.Selected < 0 {
			g.Selected = len(g.Roms) - 1
		}

	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		g.Selected -= 5 //calculted this
		if g.Selected < 0 {
			if g.Selected+25 >= len(g.Roms) {
				//each row has 5 roms, total 23 roms, so last two rows have 4 roms only, for those row last item is 4th
				g.Selected += 20
			} else {
				g.Selected += 20 + 5 //negate 5 subtracts 5 to make it positive
			}
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		g.Selected += 5
		if g.Selected >= len(g.Roms) {
			if g.Selected-25 < 0 {
				//each row has 5 roms, total 23 roms, so last two rows have 4 roms only, for those row last item is 4th
				g.Selected -= 20
			} else {
				g.Selected -= 20 + 5 //negate 5 subtracts 5 to make it positive
			}
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		//fmt.Println(g.Roms[g.Selected])
		romPath := "./pkg/assets/" + g.Roms[g.Selected]
		g.Chip8.LoadROM(romPath)
		g.GameState = true //run game
	}
}

func (g *Game) DrawMenu(screen *ebiten.Image) {
	const (
		colWidth    = 100 //space btwn each cols
		lineHeight  = 50  //space btwn each row
		itemsPerCol = 5   //each col has 5 items
		//tinkered with thse a bit to get the right feel
		startX = 90 //horizontal dist from left edge of screen
		startY = 50 // vertical dist from top edge of screen
	)
	headerOp := &text.DrawOptions{}
	headerOp.ColorScale.ScaleWithColor(color.RGBA{150, 160, 200, 255})

	// Position the header at the top center
	headerOp.PrimaryAlign = text.AlignCenter
	headerOp.GeoM.Translate(320, 10) // center in 640

	text.Draw(screen, "--- CHIP-8 ROMS ---", g.FontFace, headerOp)
	for i, rom := range g.Roms {
		//looping thru roms and printing
		label := rom
		//this kept on printing all the roms on the screen again and agiain
		//thought that was a problem, but the Draw func draws the screen at 60fps
		//this was printing the name at same time as the game is running
		//not a problem
		//fmt.Println(label
		//5 per col, when reach indx of new col colnum value changes(0, 1, 2, 3, 4)
		colNum := i / itemsPerCol //0/5 = 0, 1/5 = 0, 2/5 = 0.... 5/5 = 1, 6/5 = 1, 7/5 = 1.... 10/5 = 2, 11/5 = 2...
		rowNum := i % itemsPerCol // 0%5 = 0, 1%5 = 1, 2%5 = 2

		op := &text.DrawOptions{}
		//based on init offset values we can push text up and down like this
		x := startX + (float64(colNum) * colWidth) //for 2nd col--> 90(startX) + (1 * 100(colWidth))--> 190
		//for 3rd col--> 190(startX) + (2 * 100(colWidth))--> 290
		y := startY + (float64(rowNum) * lineHeight) //for 2nd row--> 40(startY) + (1 * 50(lineHeight))--> 90
		op.GeoM.Translate(x, y)
		op.Filter = ebiten.FilterPixelated
		if i == g.Selected {
			label = ">>" + rom + "<<"
			op.ColorScale.ScaleWithColor(color.RGBA{218, 165, 32, 255})
			text.Draw(screen, label, g.FontFace, op)
		} else {
			op.ColorScale.ScaleWithColor(color.RGBA{150, 160, 180, 255})
			op.ColorScale.Scale(0, 1, 1, 1)
			text.Draw(screen, label, g.FontFace, op)
		}

	}
}

func (g *Game) DrawSubscreen(screen *ebiten.Image) {
	//FOOTER--STATPAD
	footerWidth := 1000 //whole window should be covered-->window stats set in main.go
	footerHeight := 320 //all the height left below the game/ tinkered to figure out
	gridWidht := 80     //tinkered to figure out
	padding := 30       //tinkered to figure out
	lineHeight := 20
	//for V, I, PC, SP,..
	//also it might be better to pass footerW and footer Widht here before 0,320
	vector.FillRect(screen, 0, 320, float32(footerWidth), float32(footerHeight), STATPADCOLOR, false)
	/*stats := map[string]any{
		"V":     g.Chip8.V,
		"I":     g.Chip8.I,
		"PC":    g.Chip8.PC,
		"SP":    g.Chip8.SP,
		"DT":    g.Chip8.DT,
		"ST":    g.Chip8.ST,
		"Stack": g.Chip8.Stack[0:8], //not all 16, just 8
	}*/
	g.DrawStatpad(screen, footerHeight, gridWidht, padding, lineHeight)
	//RIGHT SUBSCREEN--HEX KEYPAD
	//right side of game screen rect-->starts at 640 on x and 0 on y
	rightWidth := 360  //1000 is widht of disp, 640 is game-->right side = 1000 - 640= 360
	rightHeight := 320 //same heigt as game disp
	//for hex keypad
	vector.FillRect(screen, 640, 0, float32(rightWidth), float32(rightHeight), KEYPADCOLOR, false)
	//BORDER-->if we do this before, the border will be covered by subscreens
	//-->window size(1000, 640)-->line needs to be betwn top half(game) and statpad
	//game(top half)-->640,320-->so line needs to be at 320 on Y and needs to cover whole X .i.e 0 to 1000 on x
	//(1000, 320)-(0, 320)= (1000, 320) covers whole X at 320 Y
	g.DrawBorder(screen, 1000, 320, 0, 320, 8)
	//vertical border
	g.DrawBorder(screen, 640, 320, 640, 0, 8) //starts at(640, 0), ends at (640, 320) a line that covers 0-320 on Y at X=640
}

func (g *Game) DrawBorder(screen *ebiten.Image, borderLength, borderHeight, startX, startY, borderWidth int) {
	vector.StrokeLine(screen, float32(borderLength), float32(borderHeight), float32(startX), float32(startY), float32(borderWidth), BORDERCOLOR, false)

}

func (g *Game) DrawStatpad(screen *ebiten.Image, footerHeight, gridWidht, lineHeight, padding int) {
	statCnfg := &text.DrawOptions{}
	statCnfg.Filter = ebiten.FilterNearest
	getPos := func(col, row int) (float64, float64) {
		return float64(padding) + (float64(col) * float64(gridWidht)), (float64(footerHeight) + float64(padding)) + (float64(row) * float64(lineHeight))
	}
	printStat := func(x, y float64, statName string, stat any) {
		statCnfg.GeoM.Reset()
		statCnfg.GeoM.Translate(x, y)
		text.Draw(screen, fmt.Sprintf("%v: %v", statName, stat), g.FontFace, statCnfg)
	}
	//PC-->GRID 1
	x, y := getPos(0, 0)
	//fmt.Print(getPos(0, 0))
	printStat(x, y, "PC", g.Chip8.PC)
	//I-->GRID 2
	x, y = getPos(0, 1) // same col, next row
	printStat(x, y, "I", g.Chip8.I)
	//SP-->GRID 3
	x, y = getPos(0, 2)
	printStat(x, y, "SP", g.Chip8.SP)
	//DT-->GRID 4
	x, y = getPos(0, 3)
	printStat(x, y, "DT", g.Chip8.DT)
	//ST-->GRID 5
	x, y = getPos(0, 4)
	printStat(x, y, "ST", g.Chip8.ST)
	//Stack-->GRID 6
	x, y = getPos(1, 0)
	printStat(x, y, "Stack", g.Chip8.Stack[0:8])
}
