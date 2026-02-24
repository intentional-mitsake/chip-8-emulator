package game

import (
	def "c8-emulation/pkg/c8-def"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	DisplayWidth  = 64
	DisplayHeight = 32
)

type Game struct {
	Chip8   *def.Chip8    //first we need to create a chip8 struct
	display *ebiten.Image // then a display image diff from the disp attr of chip8
}

func NewGame() *Game {
	c8 := def.NewChip8()
	return &Game{
		Chip8: c8,
		display: ebiten.NewImage(
			DisplayWidth,
			DisplayHeight,
		),
	}
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
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	scaling := &ebiten.DrawImageOptions{}
	scaling.GeoM.Scale(10, 10)
	screen.DrawImage(g.display, scaling)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 320
}
