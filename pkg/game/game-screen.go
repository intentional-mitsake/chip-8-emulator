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

var statChnge = true
var KEYPADCOLOR = color.RGBA{55, 65, 50, 255}
var ACTIVEKEY = color.RGBA{255, 204, 0, 255}
var KEYBG = color.RGBA{30, 35, 50, 255}      // Deep Slate Blue (Dim)
var KEYBORDER = color.RGBA{60, 90, 160, 255} // Bold Steel Blue
var STATPADCOLOR = color.RGBA{10, 15, 20, 255}
var STATSCOLOR = color.RGBA{0, 255, 255, 255}
var BORDERCOLOR = color.RGBA{0, 40, 60, 255}
var GRIDCOLOR = color.RGBA{5, 15, 40, 255}
var GRIDBORDER = color.RGBA{40, 70, 140, 255}
var SPECIALINACTIVE = color.RGBA{30, 40, 60, 255} // Deep Midnight Blue

var offset int = 0 //indx of first visibile rom

type Game struct {
	//handling game state
	GameState bool     //running(true) or paused(false)
	Roms      []string //arrat of all available roms
	Selected  int      //index of currently selected rom
	RomState  bool
	//Chip-8 def and ebiten image
	Chip8   *def.Chip8    //first we need to create a chip8 struct
	display *ebiten.Image // then a display image diff from the disp attr of chip8

	//for menu
	FontFace text.Face
	DrawOpts *text.DrawOptions
	//stats
	updtC int
	stats map[string]any //keys are string and values are any
	//for key press effect for all 20 keys(hex keys + 4 additional ones)
	keyTime [20]int
	//help screen
	HelpState bool
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
		stats: map[string]any{
			"PC":     0,
			"I":      0,
			"SP":     0,
			"DT":     0,
			"ST":     0,
			"Stack":  0,
			"V":      0,
			"Opcode": 0,
		},
		HelpState: false,
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
		//if the gaem is running run the core loop, if paused do nothing
		if g.RomState {
			for i := 0; i < 10; i++ {
				//reason for 10 is: chip8 ran around 500-1000 inst per second
				//ebiten runs at 60fps, with the i<10, we run the core loop 10 times each frame
				//chip 8 runs at appx 600Hz, so we run the core loop 10 times each frame, 60 frames per second
				//so 600 inst per second
				g.Chip8.CoreLoop()
			}
		}
		//input
		g.HandleKeyPad()

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
		//without this check, stats will be updated every frame
		//so it flickers really fast and we cant see shit
		//enitengine runs the game at 60 fps so witout this check, stats will be updated every frame
		//with this check, stats will be updated every 10 frames(update counter>10)
		//we can incr the num in updC>num to slow it down even more
		//**NOTE** due to this we only see a snap/part of the changing stats, not whole chagens
		//for that we remove this and then cant see shit
		//so its  a tradeoff to be able to see some chagne or not seeing anyting dur to speed
		g.updtC++
		if g.updtC >= 10 {
			g.stats["PC"] = g.Chip8.PC
			g.stats["I"] = g.Chip8.I
			g.stats["SP"] = g.Chip8.SP
			g.stats["DT"] = g.Chip8.DT
			g.stats["ST"] = g.Chip8.ST
			g.stats["Stack"] = g.Chip8.Stack
			g.stats["V"] = g.Chip8.V
			g.stats["Opcode"] = g.Chip8.Opcode
			g.updtC = 0
		}
	case false:
		//not running - menu
		g.Menu()
		//input
		g.HandleKeyPad()
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
		if g.HelpState {
			g.Help(screen)
		}
	case false:
		//not running
		screen.Fill(color.RGBA{15, 20, 30, 255})
		g.DrawMenu(screen)
		g.DrawSubscreen(screen)
		if g.HelpState {
			g.Help(screen)
		}
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
		g.RomState = true  //rom is runnign
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
	footerWidth := 640  //statpad should cover same width as game but below it
	footerHeight := 320 //all the height left below the game/ tinkered to figure out
	gridWidht := 200    //tinkered to figure out
	padding := 80       //tinkered to figure out
	lineHeight := 30
	//for V, I, PC, SP,..
	//also it might be better to pass footerW and footer Widht here before 0,320
	vector.FillRect(screen, 0, 320, float32(footerWidth), float32(footerHeight), STATPADCOLOR, false)
	g.DrawStatpad(screen, footerHeight, gridWidht, padding, lineHeight)
	//RIGHT SUBSCREEN--HEX KEYPAD
	//right side of game screen rect-->starts at 640 on x and 0 on y
	rightWidth := 360  //1000 is widht of disp, 640 is game-->right side = 1000 - 640= 360
	rightHeight := 640 //right subscreen for keypad covers the whole height
	//for hex keypad
	vector.FillRect(screen, 640, 0, float32(rightWidth), float32(rightHeight), KEYPADCOLOR, false)
	g.DrawKeypad(screen, rightWidth, rightHeight, 20, 640, 0)
	//BORDER-->if we do this before, the border will be covered by subscreens
	//-->window size(1000, 640)-->line needs to be betwn top half(game) and statpad
	//game(top half)-->640,320-->so line needs to be at 320 on Y and needs to cover whole X .i.e 0 to 1000 on x
	//(1000, 320)-(0, 320)= (1000, 320) covers whole X at 320 Y
	g.DrawBorder(screen, 640, 320, 0, 320, 8) //tills 640, 320 now(changed subscreen div form beofre)
	//vertical border
	g.DrawBorder(screen, 640, 640, 640, 0, 8) //starts at(640, 0), ends at (640, 640) a line that covers 0-640 on Y at X=640
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
		statCnfg.GeoM.Reset() //translate stacks(adds) for each stat	so we need to clear pos input by prev stat
		statCnfg.GeoM.Translate(x, y)
		boxWidth := 100
		if statName == "Stack" || statName == "V" {
			boxWidth = 400 //smtimes needs 400
		}
		if statName == "Opcode" {
			statCnfg.ColorScale.ScaleWithColor(STATSCOLOR)
			vector.FillRect(screen, float32(x-5), float32(y-15), float32(boxWidth+10), 60, GRIDBORDER, false)
			vector.FillRect(screen, float32(x), float32(y-10), float32(boxWidth), 50, GRIDCOLOR, false)
			text.Draw(screen, fmt.Sprintf("%v: %X ", statName, stat), g.FontFace, statCnfg) //for hex
			return
		}
		statCnfg.ColorScale.ScaleWithColor(STATSCOLOR)
		vector.FillRect(screen, float32(x-5), float32(y-15), float32(boxWidth+10), 60, GRIDBORDER, false)
		vector.FillRect(screen, float32(x), float32(y-10), float32(boxWidth), 50, GRIDCOLOR, false)
		text.Draw(screen, fmt.Sprintf("%v: %v ", statName, stat), g.FontFace, statCnfg)
	}
	//PC-->GRID 1
	x, y := getPos(0, 0)
	//fmt.Print(getPos(0, 0))
	printStat(x, y, "PC", g.stats["PC"])
	//I-->GRID 2
	x, y = getPos(0, 1) // same col, next row
	printStat(x, y, "I", g.stats["I"])
	//DT-->GRID 4
	x, y = getPos(0, 2)
	printStat(x, y, "DT", g.stats["DT"])
	//ST-->GRID 5
	x, y = getPos(0, 3)
	printStat(x, y, "ST", g.stats["ST"])
	//SP-->GRID 3
	x, y = getPos(1, 0)
	printStat(x, y, "SP", g.stats["SP"])
	//OPCODE-->GRID 8
	x, y = getPos(1, 2)
	printStat(x, y, "Opcode", g.stats["Opcode"])
	//Stack-->GRID 6
	x, y = getPos(1, 1)
	printStat(x, y, "Stack", g.stats["Stack"])
	//V-->GRID 7
	x, y = getPos(1, 3)
	printStat(x, y, "V", g.stats["V"])

}

func (g *Game) DrawKeypad(screen *ebiten.Image, keypadWidth, keypadHeight, padding, startX, startY int) {
	titleCnfg := &text.DrawOptions{}
	titleCnfg.Filter = ebiten.FilterNearest
	centerX := startX + keypadWidth/2
	titleCnfg.PrimaryAlign = text.AlignCenter //to center the text on the given pos
	titleCnfg.GeoM.Translate(float64(centerX), float64(startY+40))
	titleCnfg.ColorScale.ScaleWithColor(ACTIVEKEY)
	text.Draw(screen, "---Hex Keypad---", g.FontFace, titleCnfg)
	//key text config
	keyCnfg := &text.DrawOptions{}
	for i := 0; i < 20; i++ {
		keyCnfg.GeoM.Reset()
		//inverted this to rotate the keypad
		//if colNUm = i/4 and rowNum = i%4
		//then first row of keypad: 1 Q A Z
		//when use this inverted we get first row of keypad: 1 2 3 4
		colNum := i % 4
		rowNum := i / 4
		//fig the best thru tiral and erro
		gridLength := keypadWidth / 4
		x := startX + padding - 10 + colNum*gridLength //without the 8, the keypad shifts a lil more to left and last keys touch border
		y := startY + padding + 80 + rowNum*gridLength //to push the keypad down even more than padding
		keyCnfg.Filter = ebiten.FilterNearest
		keyCnfg.ColorScale.ScaleWithColor(color.Black)
		keyCnfg.GeoM.Translate(float64(x+(gridLength-padding)/2), float64(y+(gridLength-padding)/2))
		keyN, byteVal := GetKey(i)
		keyCnfg.PrimaryAlign = text.AlignCenter
		keyCnfg.SecondaryAlign = text.AlignCenter // Centers vertically too!
		//check if the effect should be active
		effectActive := g.keyTime[i] > 0
		if i > 15 && effectActive {
			//the -8 and +8 is adjusted to give the feel of key being presesd or unpressed
			//prettry much shifting the pos of either the first rect or the second rect for htis
			vector.FillRect(screen, float32(x-8), float32(y-8), float32(gridLength-padding+8), float32(gridLength-padding+8), KEYBG, false)
			vector.FillRect(screen, float32(x), float32(y), float32(gridLength-padding), float32(gridLength-padding), ACTIVEKEY, false)
			text.Draw(screen, fmt.Sprintf("%s", keyN), g.FontFace, keyCnfg)
		} else if effectActive {
			vector.FillRect(screen, float32(x), float32(y), float32(gridLength-padding+8), float32(gridLength-padding+8), KEYBG, false)
			vector.FillRect(screen, float32(x), float32(y), float32(gridLength-padding), float32(gridLength-padding), ACTIVEKEY, false)
			text.Draw(screen, fmt.Sprintf("%s-->%X", keyN, byteVal), g.FontFace, keyCnfg)
		} else if i > 15 {
			vector.FillRect(screen, float32(x), float32(y), float32(gridLength-padding+8), float32(gridLength-padding+8), KEYBG, false)
			vector.FillRect(screen, float32(x), float32(y), float32(gridLength-padding), float32(gridLength-padding), SPECIALINACTIVE, false)
			keyCnfg.ColorScale.ScaleWithColor(color.White)
			text.Draw(screen, fmt.Sprintf("%s", keyN), g.FontFace, keyCnfg)
		} else {
			vector.FillRect(screen, float32(x), float32(y), float32(gridLength-padding+8), float32(gridLength-padding+8), KEYBG, false)
			vector.FillRect(screen, float32(x), float32(y), float32(gridLength-padding), float32(gridLength-padding), KEYBORDER, false)
			text.Draw(screen, fmt.Sprintf("%s-->%X", keyN, byteVal), g.FontFace, keyCnfg)
		}

	}
}

func (g *Game) HandleKeyPad() {
	//continously reset keypad
	for i := range g.Chip8.Keypad {
		g.Chip8.Keypad[i] = false
	}
	//using inpoututil.IsKeyJustPressed cuz it returns true for first frame only so no multi clicks on one press effect
	for k, v := range keyMap {
		if v > 15 {
			//for special keys will continue using inpututil cuz we need just first frame and its enough
			if inpututil.IsKeyJustPressed(k) {
				g.Chip8.Keypad[v] = true
				print(fmt.Sprintf("Pressed: %v\n", k))
				switch k {
				case ebiten.KeyBackspace:
					g.Chip8.Reset()     //reset chip and clear disp
					g.GameState = false //change gamestate to go back to the menu
				case ebiten.KeyP:
					//no need to return to menu and reset chip8, just toggle pause
					if g.RomState {
						g.RomState = false
					} else {
						g.RomState = true
					}
				case ebiten.KeyN:
					g.Chip8.Reset()
					g.RomState = true //if press N while paused, romstate is false, so need to press P again
					//but with thsi, even if paused no need to press P again after pressing N
					// a new rom isnt loaded in the reset func, only fontset is loaded
					romPath := "./pkg/assets/" + g.Roms[g.Selected]
					g.Chip8.LoadROM(romPath)
				case ebiten.KeyH:
					g.HelpState = !g.HelpState //toggle it
					if g.HelpState {
						//first pause the game then show the help screen
						g.RomState = false
					} else {
						g.RomState = true
						g.HelpState = false
					}
				}
				g.keyTime[v] = 10
			}
		} else {
			//switched to ebiten.IsKeyPressed cuz in actual game, just first frame result isnt enough, need all frames key press result
			if ebiten.IsKeyPressed(k) {
				g.Chip8.Keypad[v] = true
				print(fmt.Sprintf("Pressed: %v\n", k))
				fmt.Println(v)
				//if key was pressed put a timer to hold the effect
				//byte from 0x0 to 0xF to 20 its pretty much 0-20
				g.keyTime[v] = 10
			}
		}
	}

	for i := 0; i < 20; i++ {
		//if key isnt active and time is up, reset
		if !g.Chip8.Keypad[i] && g.keyTime[i] > 0 {
			g.keyTime[i]--
		}
	}
}

func (g *Game) Help(screen *ebiten.Image) {

	//fill the entire screen
	screen.Fill(color.RGBA{5, 20, 10, 230})
	vector.StrokeRect(screen, 50, 20, 900, 560, 4, color.RGBA{0, 200, 80, 255}, false)
	//HEADER
	textCnfg := &text.DrawOptions{}
	textCnfg.Filter = ebiten.FilterNearest
	//for each line
	textCnfg.GeoM.Reset()
	textCnfg.GeoM.Translate(500, 60) //middle horizontally and and lil down from top
	textCnfg.PrimaryAlign = text.AlignCenter
	textCnfg.ColorScale.ScaleWithColor(ACTIVEKEY)
	text.Draw(screen, "---- HOW TO USE THE EMULATOR ----", g.FontFace, textCnfg)

	//GAME KEYS (Hex)
	// Shift to Left Align for the list
	textCnfg.PrimaryAlign = text.AlignStart
	textCnfg.ColorScale.Reset()
	textCnfg.ColorScale.ScaleWithColor(color.RGBA{140, 255, 170, 255})

	gameKeys := []string{
		"> GAME KEYS (0-F ONLY)",
		"KEYBOARD ->  HEX VALUE",
		"1 2 3 4  ->  HEX 1 2 3 C",
		"Q W E R  ->  HEX 4 5 6 D",
		"A S D F  ->  HEX 7 8 9 E",
		"Z X C V  ->  HEX A 0 B F",
	}

	for i, line := range gameKeys {
		textCnfg.GeoM.Reset()
		//each line is 120 pixels from left edge of screen
		//each line is 130 + i*30 pixels from top i.e 130, 160, 190, 220
		//so new line for each
		textCnfg.GeoM.Translate(120, float64(130+(i*30)))
		text.Draw(screen, line, g.FontFace, textCnfg)
	}

	//SYSTEM KEYS
	specKeys := []string{
		"> SYSTEM CONTROLS",
		"[BKSP]          ->  Hard Reset: Clears Memory and Registers, Chip 8 in Initial State",
		"[P]             ->  Pause/Resume: Stops the Fetch-Exc Cycle of the Chip or resumes it.",
		"[N]             ->  New Game: Resets Chip8 to initial state and loads ROM again.",
		"[H]             ->  Help Manual.",
		"[ESC] [Ctrl+C]  ->  Exit App.",
	}

	for i, line := range specKeys {
		textCnfg.GeoM.Reset()
		textCnfg.GeoM.Translate(120, float64(320+(i*30)))
		text.Draw(screen, line, g.FontFace, textCnfg)
	}

	//footer
	textCnfg.GeoM.Reset()
	textCnfg.GeoM.Translate(500, 540)
	textCnfg.PrimaryAlign = text.AlignCenter
	textCnfg.ColorScale.Reset()
	textCnfg.ColorScale.ScaleWithColor(color.RGBA{200, 50, 50, 255})
	text.Draw(screen, "CHIP8 EMULATION INTENTIONAL-MITSAKE", g.FontFace, textCnfg)
}
