package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

var keyMap = map[ebiten.Key]byte{
	//all 16(0-F) keys mapped to qwety keys
	ebiten.Key1: 0x0, ebiten.Key2: 0x1, ebiten.Key3: 0x2, ebiten.Key4: 0x3,
	ebiten.KeyQ: 0x4, ebiten.KeyW: 0x5, ebiten.KeyE: 0x6, ebiten.KeyR: 0x7,
	ebiten.KeyA: 0x8, ebiten.KeyS: 0x9, ebiten.KeyD: 0xA, ebiten.KeyF: 0xB,
	ebiten.KeyZ: 0xC, ebiten.KeyX: 0xD, ebiten.KeyC: 0xE, ebiten.KeyV: 0xF,
	ebiten.KeyBackspace: 16, ebiten.KeyP: 17, ebiten.KeyN: 18, ebiten.KeyH: 19,
}

func GetKey(index int) (qwerty string, hex byte) {
	if index < 16 {
		keySet := map[int]string{
			0: "1", 1: "2", 2: "3", 3: "4",
			4: "Q", 5: "W", 6: "E", 7: "R",
			8: "A", 9: "S", 10: "D", 11: "F",
			12: "Z", 13: "X", 14: "C", 15: "V",
		}
		//%X gives hex value of given number
		//format: 0xHexValue \n [Key]
		return keySet[index], byte(index)
	}
	//diff labels for special keys
	specialSet := map[int]string{
		16: "<--[BKSP]",
		17: "PAUSE[P]",
		18: "NGAME[N]",
		19: "HELP[H]",
	}
	return specialSet[index], byte(index)
}
