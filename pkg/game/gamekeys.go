package game

type GameMeta struct {
	Title    string
	Controls []string
}

var GameLibrary = map[string]GameMeta{
	"15PUZZLE": {
		Title: "15 PUZZLE",
		Controls: []string{
			"2, Q, W, 3 [HEX 1, 4, 5, 2] -> MOVE TILES",
		},
	},
	"BLINKY": {
		Title: "BLINKY (PAC-MAN)",
		Controls: []string{
			"3 [HEX 2] -> UP    /  E [HEX 6] -> DOWN",
			"Q [HEX 4] -> LEFT  /  W [HEX 5] -> RIGHT",
		},
	},
	"BLITZ": {
		Title: "BLITZ (BOMBER)",
		Controls: []string{
			"W [HEX 5] -> DROP BOMB",
		},
	},
	"BRIX": {
		Title: "BRIX (BREAKOUT)",
		Controls: []string{
			"Q [HEX 4] -> MOVE LEFT",
			"E [HEX 6] -> MOVE RIGHT",
		},
	},
	"CONNECT4": {
		Title: "CONNECT 4",
		Controls: []string{
			"Q [HEX 4] -> MOVE LEFT",
			"E [HEX 6] -> MOVE RIGHT",
			"W [HEX 5] -> DROP PIECE",
		},
	},
	"GUESS": {
		Title: "GUESS THE NUMBER",
		Controls: []string{
			"HEX 0-F -> INPUT YOUR GUESS",
		},
	},
	"HIDDEN": {
		Title: "HIDDEN (MEMORY)",
		Controls: []string{
			"HEX 0-F -> SELECT CARD PAIRS",
		},
	},
	"INVADERS": {
		Title: "SPACE INVADERS",
		Controls: []string{
			"Q [HEX 4] -> MOVE LEFT",
			"E [HEX 6] -> MOVE RIGHT",
			"W [HEX 5] -> FIRE CANNON",
		},
	},
	"KALEID": {
		Title: "KALEIDOSCOPE",
		Controls: []string{
			"2, Q, W, 3 [HEX 1, 4, 5, 2] -> CHANGE PATTERN",
		},
	},
	"MAZE": {
		Title: "MAZE GENERATOR",
		Controls: []string{
			"ANY KEY -> GENERATE NEW MAZE",
		},
	},
	"MERLIN": {
		Title: "MERLIN (SIMON SAYS)",
		Controls: []string{
			"Q, W, A, S [HEX 4, 5, 7, 8] -> FOLLOW LIGHTS",
		},
	},
	"MISSILE": {
		Title: "MISSILE COMMAND",
		Controls: []string{
			"S [HEX 8] -> FIRE MISSILE",
		},
	},
	"PONG": {
		Title: "PONG (1 OR 2 PLAYER)",
		Controls: []string{
			"2 [HEX 1] -> UP    /  Q [HEX 4] -> DOWN",
		},
	},
	"PONG2": {
		Title: "PONG 2 (VS AI)",
		Controls: []string{
			"2 [HEX 1] -> UP    /  Q [HEX 4] -> DOWN",
		},
	},
	"PUZZLE": {
		Title: "PUZZLE SLIDER",
		Controls: []string{
			"2, Q, W, 3 [HEX 1, 4, 5, 2] -> SLIDE TILES",
		},
	},
	"SYZYGY": {
		Title: "SYZYGY (SNAKE)",
		Controls: []string{
			"3 [HEX 2] -> UP    /  E [HEX 6] -> DOWN",
			"Q [HEX 4] -> LEFT  /  W [HEX 5] -> RIGHT",
		},
	},
	"TANK": {
		Title: "TANK BATTLE",
		Controls: []string{
			"3 [HEX 2] -> UP    /  E [HEX 6] -> DOWN",
			"Q [HEX 4] -> LEFT  /  W [HEX 5] -> RIGHT",
			"W [HEX 5] -> SHOOT",
		},
	},
	"TETRIS": {
		Title: "TETRIS",
		Controls: []string{
			"W [HEX 5] -> MOVE LEFT",
			"E [HEX 6] -> MOVE RIGHT",
			"Q [HEX 4] -> ROTATE PIECE",
		},
	},
	"TICTAC": {
		Title: "TIC-TAC-TOE",
		Controls: []string{
			"1-3 / Q-E / A-D -> SELECT GRID SQUARE",
		},
	},
	"UFO": {
		Title: "UFO (INTERCEPT)",
		Controls: []string{
			"Q [HEX 4] -> SHOOT LEFT",
			"W [HEX 5] -> SHOOT UP",
			"E [HEX 6] -> SHOOT RIGHT",
		},
	},
	"VBRIX": {
		Title: "VERTICAL BRIX",
		Controls: []string{
			"2 [HEX 1] -> PADDLE UP",
			"Q [HEX 4] -> PADDLE DOWN",
		},
	},
	"VERS": {
		Title: "VERS (TRON)",
		Controls: []string{
			"P1: 2, Q, W, 3 [HEX 1, 4, 5, 2]",
			"P2: A, S, D, F [HEX 7, 8, 9, E]",
		},
	},
	"WIPEOFF": {
		Title: "WIPEOFF (BREAKOUT)",
		Controls: []string{
			"Q [HEX 4] -> MOVE LEFT",
			"E [HEX 6] -> MOVE RIGHT",
		},
	},
}
