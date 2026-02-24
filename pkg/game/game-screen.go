package game

import (
	def "c8-emulation/pkg/c8-def"
	"fmt"
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/basicfont"
)

const (
	DisplayWidth  = 64
	DisplayHeight = 32
	visible       = 10
)

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
		DrawOpts: &text.DrawOptions{
			// Embedded fields:

		},
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
	pixels := make([]byte, DisplayWidth*DisplayHeight*4) // RGBA

	for i := 0; i < DisplayWidth*DisplayHeight; i++ {
		if g.Chip8.Display[i] == 1 {
			pixels[i*4+0] = 255 // R
			pixels[i*4+1] = 255 // G
			pixels[i*4+2] = 255 // B
			pixels[i*4+3] = 255 // A
		} else {
			pixels[i*4+0] = 0
			pixels[i*4+1] = 0
			pixels[i*4+2] = 0
			pixels[i*4+3] = 255
		}
	}

	g.display.ReplacePixels(pixels)
}

func (g *Game) Update() error {
	switch g.GameState {
	case true:
		//running
	case false:
		//not running - menu
		g.Menu()
	}
	/*for i := 0; i < 10; i++ {
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
	*/
	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {
	switch g.GameState {
	case true:
		//running
		screen.Fill(color.Black)
		scaling := &ebiten.DrawImageOptions{}
		scaling.GeoM.Scale(10, 10)
		screen.DrawImage(g.display, scaling)
	case false:
		//not running
		screen.Fill(color.Black)
		g.DrawMenu(screen)
	}
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 320
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

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		fmt.Println(g.Roms[g.Selected])
		//g.Chip8.LoadROM()
		//g.GameState = true //run game
	}
}

func (g *Game) DrawMenu(screen *ebiten.Image) {
	for i, rom := range g.Roms {
		//looping thru roms and printing
		label := rom
		//this kept on printing all the roms on the screen again and agiain
		//thought that was a problem, but the Draw func draws the screen at 60fps
		//this was printing the name at same time as the game is running
		//not a problem
		//fmt.Println(label)
		if g.Selected < offset {
			offset = g.Selected
		} else if g.Selected >= offset+visible {
			offset = g.Selected - visible + 1
		}
		for j := 0; j < visible; j++ {
			romIndx := j + offset
			if romIndx >= len(g.Roms) {
				break
			}
			if i == g.Selected {
				label = ">> " + rom
			}
			op := &text.DrawOptions{}
			op.GeoM = g.DrawOpts.GeoM
			lineSpacing := 24.0
			op.GeoM.Translate(0, float64(i)*lineSpacing)
			text.Draw(screen, label, g.FontFace, op)

		}
	}
}
