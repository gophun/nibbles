//                         G o   N i b b l e s
//
// Original QBasic version: Copyright (C) Microsoft Corporation 1990
//
// Nibbles is a game for one or two players. Navigate your snakes
// around the game board trying to eat up numbers while avoiding
// running into walls or other snakes. The more numbers you eat up,
// the more points you gain and the longer your snake becomes.

package main

import . "github.com/gophun/nibbles/internal/basic"

type snakeBody struct {
	row, col int
}

// This type defines the player's snake.
type snakeType struct {
	head      int
	length    int
	row       int
	col       int
	direction int
	lives     int
	score     int
	color     int
	alive     bool
}

// This type is used to represent the playing screen in memory.
// It is used to simulate graphics in text mode, and has some interesting,
// and slightly advanced methods to increasing the speed of operation.
// Instead of the normal 80x25 text graphics using "█", we will be
// using "▄" and "▀" and "█" to mimic an 80x50 pixel screen.
// Check out functions Set and PointIsThere to see how this is implemented.
type arenaType struct {
	realRow int // Maps the 80x50 point into the real 80x25.
	color   int // Stores the current color of the point.
	// Each char has 2 points in it.
	// sister is -1 if sister point is above, +1 if below.
	sister int
}

const MaxSnakeLength = 1000

// Parameters to Level function
const (
	StartOver = 1 + iota
	SameLevel
	NextLevel
)

var (
	arena      [][]arenaType
	curLevel   int
	colorTable []int
)

func main() {
	Randomize(Timer())
	Intro()
	defer Reset()
	numPlayers, speed, diff, monitor := GetInputs()
	SetColors(monitor)
	DrawScreen()
	for {
		PlayNibbles(numPlayers, speed, diff)
		if !StillWantsToPlay() {
			break
		}
	}
	Color(15, 0)
	Cls()
}

func SetColors(monitor string) {
	if monitor == "M" {
		colorTable = mono
		return
	}
	colorTable = normal
}

var (
	// {snake1, snake2, Walls, Background, Dialogs-Fore, Back}
	mono   = []int{15, 7, 7, 0, 15, 0}
	normal = []int{14, 13, 12, 1, 15, 4}
)

// Center centers text on given row.
func Center(row int, text string) {
	Locate(row, CInt(41-float64(Len(text))/2))
	Print(text)
}

// DrawScreen draws the playing field.
func DrawScreen() {
	// initialize screen
	Color(colorTable[0], colorTable[3])
	Cls()

	// Print title & message
	Center(1, "Nibbles!")
	Center(11, "Initializing Playing Field...")

	// Initialize arena array
	arena = make([][]arenaType, 50)
	for row := 1; row <= len(arena); row++ {
		arena[row-1] = make([]arenaType, 80)
		for col := 1; col <= 80; col++ {
			arena[row-1][col-1].realRow = (row + 1) / 2
			arena[row-1][col-1].sister = (row%2)*2 - 1
		}
	}
	Sleep(1) // Adds authenticity
}

// EraseSnake erases snake to facilitate moving through playing field.
func EraseSnake(snake []snakeType, snakeBod [][2]snakeBody, snakeNum int) {
	for c := 0; c < 10; c++ {
		for b := snake[snakeNum].length - c; b >= 0; b -= 10 {
			tail := (snake[snakeNum].head + MaxSnakeLength - b) % MaxSnakeLength
			Set(snakeBod[tail][snakeNum].row, snakeBod[tail][snakeNum].col, colorTable[3])
		}
		SleepMillis(20)
	}
}

// GetInputs gets player inputs.
func GetInputs() (numPlayers, speed int, diff, monitor string) {
	Color(7, 0)
	Cls()

	for numPlayers != 1 && numPlayers != 2 {
		Locate(5, 4)
		Print(Space(34))
		Locate(5, 20)
		numPlayers = Val(Input("How many players (1 or 2)"))
	}

	Locate(8, 21)
	Print("Skill level (1 to 100)")
	Locate(9, 22)
	Print("1   = Novice")
	Locate(10, 22)
	Print("90  = Expert")
	Locate(11, 22)
	Print("100 = Twiddle Fingers")
	Locate(12, 15)
	Print("(Computer speed may affect your skill level)")
	for speed < 1 || speed > 100 {
		Locate(8, 44)
		Print(Space(35))
		Locate(8, 43)
		speed = Val(Input(""))
	}
	speed = int(float64(100-speed)*2) + 1

	for diff != "Y" && diff != "N" {
		Locate(15, 56)
		Print(Space(25))
		Locate(15, 15)
		diff = UCase(Input("Increase game speed during play (Y or N)"))
	}

	for monitor != "M" && monitor != "C" {
		Locate(17, 46)
		Print(Space(34))
		Locate(17, 17)
		monitor = UCase(Input("Monochrome or color monitor (M or C)"))
	}

	return numPlayers, speed, diff, monitor
}

// InitColors initializes playing field colors.
func InitColors() {
	for row := range arena {
		for col := range arena[row] {
			arena[row][col].color = colorTable[3]
		}
	}
	Cls()

	// Set (turn on) pixels for screen border
	for col := 1; col <= 80; col++ {
		Set(3, col, colorTable[2])
		Set(50, col, colorTable[2])
	}
	for row := 4; row <= 49; row++ {
		Set(row, 1, colorTable[2])
		Set(row, 80, colorTable[2])
	}
}

// Intro displays the game introduction.
func Intro() {
	Screen(0)
	Width(80, 25)
	Color(15, 0)
	Cls()

	//Center(4, "Q B a s i c   N i b b l e s")
	Center(4, "G o   N i b b l e s")
	ColorFg(7)
	//Center(6, "Copyright (C) Microsoft Corporation 1990")
	Center(6, "(Translated from QBasic Nibbles)")
	Center(8, "Nibbles is a game for one or two players.  Navigate your snakes")
	Center(9, "around the game board trying to eat up numbers while avoiding")
	Center(10, "running into walls or other snakes.  The more numbers you eat up,")
	Center(11, "the more points you gain and the longer your snake becomes.")
	Center(13, " Game Controls ")
	Center(15, "  General             Player 1               Player 2    ")
	Center(16, "                        (Up)                   (Up)      ")
	Center(17, "P - Pause                ↑                      W       ")
	Center(18, "                     (Left) ←   → (Right)   (Left) A   D (Right)  ")
	Center(19, "                         ↓                      S       ")
	Center(20, "                       (Down)                 (Down)     ")
	Center(24, "Press any key to continue")

	Play("MBT160O1L8CDEDCDL4ECC")
	SparklePause()
}

// Level sets the game level.
func Level(whatToDo int, sammy []snakeType) {
	switch whatToDo {
	case StartOver:
		curLevel = 1
	case NextLevel:
		curLevel++
	}

	// Initialize Snakes
	sammy[0].head = 1
	sammy[0].length = 2
	sammy[0].alive = true
	sammy[1].head = 1
	sammy[1].length = 2
	sammy[1].alive = true

	InitColors()

	switch curLevel {
	case 1:
		sammy[0].row = 25
		sammy[1].row = 25
		sammy[0].col = 50
		sammy[1].col = 30
		sammy[0].direction = 4
		sammy[1].direction = 3
	case 2:
		for i := 20; i <= 60; i++ {
			Set(25, i, colorTable[2])
		}
		sammy[0].row = 7
		sammy[1].row = 43
		sammy[0].col = 60
		sammy[1].col = 20
		sammy[0].direction = 3
		sammy[1].direction = 4
	case 3:
		for i := 10; i <= 40; i++ {
			Set(i, 20, colorTable[2])
			Set(i, 60, colorTable[2])
		}
		sammy[0].row = 25
		sammy[1].row = 25
		sammy[0].col = 50
		sammy[1].col = 30
		sammy[0].direction = 1
		sammy[1].direction = 2
	case 4:
		for i := 4; i <= 30; i++ {
			Set(i, 20, colorTable[2])
			Set(53-i, 60, colorTable[2])
		}
		for i := 2; i <= 40; i++ {
			Set(38, i, colorTable[2])
			Set(15, 81-i, colorTable[2])
		}
		sammy[0].row = 7
		sammy[1].row = 43
		sammy[0].col = 60
		sammy[1].col = 20
		sammy[0].direction = 3
		sammy[1].direction = 4
	case 5:
		for i := 13; i <= 39; i++ {
			Set(i, 21, colorTable[2])
			Set(i, 59, colorTable[2])
		}
		for i := 23; i <= 57; i++ {
			Set(11, i, colorTable[2])
			Set(41, i, colorTable[2])
		}
		sammy[0].row = 25
		sammy[1].row = 25
		sammy[0].col = 50
		sammy[1].col = 30
		sammy[0].direction = 1
		sammy[1].direction = 2
	case 6:
		for i := 4; i <= 49; i++ {
			if i > 30 || i < 23 {
				Set(i, 10, colorTable[2])
				Set(i, 20, colorTable[2])
				Set(i, 30, colorTable[2])
				Set(i, 40, colorTable[2])
				Set(i, 50, colorTable[2])
				Set(i, 60, colorTable[2])
				Set(i, 70, colorTable[2])
			}
		}
		sammy[0].row = 7
		sammy[1].row = 43
		sammy[0].col = 65
		sammy[1].col = 15
		sammy[0].direction = 2
		sammy[1].direction = 1
	case 7:
		for i := 4; i <= 49; i += 2 {
			Set(i, 40, colorTable[2])
		}
		sammy[0].row = 7
		sammy[1].row = 43
		sammy[0].col = 65
		sammy[1].col = 15
		sammy[0].direction = 2
		sammy[1].direction = 1
	case 8:
		for i := 4; i <= 40; i++ {
			Set(i, 10, colorTable[2])
			Set(53-i, 20, colorTable[2])
			Set(i, 30, colorTable[2])
			Set(53-i, 40, colorTable[2])
			Set(i, 50, colorTable[2])
			Set(53-i, 60, colorTable[2])
			Set(i, 70, colorTable[2])
		}
		sammy[0].row = 7
		sammy[1].row = 43
		sammy[0].col = 65
		sammy[1].col = 15
		sammy[0].direction = 2
		sammy[1].direction = 1
	case 9:
		for i := 6; i <= 47; i++ {
			Set(i, i, colorTable[2])
			Set(i, i+28, colorTable[2])
		}
		sammy[0].row = 40
		sammy[1].row = 15
		sammy[0].col = 75
		sammy[1].col = 5
		sammy[0].direction = 1
		sammy[1].direction = 2
	default:
		for i := 4; i <= 49; i += 2 {
			Set(i, 10, colorTable[2])
			Set(i+1, 20, colorTable[2])
			Set(i, 30, colorTable[2])
			Set(i+1, 40, colorTable[2])
			Set(i, 50, colorTable[2])
			Set(i+1, 60, colorTable[2])
			Set(i, 70, colorTable[2])
		}
		sammy[0].row = 7
		sammy[1].row = 43
		sammy[0].col = 65
		sammy[1].col = 15
		sammy[0].direction = 2
		sammy[1].direction = 1
	}
}

// PlayNibbles is the main routine that controls game play.
func PlayNibbles(numPlayers, speed int, diff string) {

	// Initialize snakes
	sammyBody := make([][2]snakeBody, MaxSnakeLength)
	sammy := make([]snakeType, 2)
	sammy[0].lives = 5
	sammy[0].score = 0
	sammy[0].color = colorTable[0]
	sammy[1].lives = 5
	sammy[1].score = 0
	sammy[1].color = colorTable[1]

	Level(StartOver, sammy)

	curSpeed := speed

	// Play Nibbles until finished

	SpacePause("     Level" + Str(curLevel) + ",  Push Space")

	for sammy[0].lives > 0 && sammy[1].lives > 0 {
		// Play next round, until either of snake's lives have run out.

		if numPlayers == 1 {
			sammy[1].row = 0
		}

		number := 1   // Current number that snakes are trying to run into
		noNum := true // noNum = true if a number is not on the screen
		var numberRow, numberCol int
		playerDied := false

		PrintScore(numPlayers, sammy[0].score, sammy[1].score, sammy[0].lives, sammy[1].lives)
		Play("T160O1>L20CDEDCDL10ECC")

		for !playerDied {
			// Print number if no number exists
			if noNum {
				for {
					numberRow = int(Rnd(1)*47 + 3)
					numberCol = int(Rnd(1)*78 + 2)
					sisterRow := numberRow + arena[numberRow-1][numberCol-1].sister
					if !PointIsThere(numberRow, numberCol, colorTable[3]) && !PointIsThere(sisterRow, numberCol, colorTable[3]) {
						break
					}
				}
				numberRow = arena[numberRow-1][numberCol-1].realRow
				noNum = false
				Color(colorTable[0], colorTable[3])
				Locate(numberRow, numberCol)
				Print(Right(Str(number), 1))
			}

			// Delay game
			SleepMillis(curSpeed)

			// Get keyboard input & change direction accordingly
			switch InKey() {
			case "w", "W":
				if sammy[1].direction != 2 {
					sammy[1].direction = 1
				}
			case "s", "S":
				if sammy[1].direction != 1 {
					sammy[1].direction = 2
				}
			case "a", "A":
				if sammy[1].direction != 4 {
					sammy[1].direction = 3
				}
			case "d", "D":
				if sammy[1].direction != 3 {
					sammy[1].direction = 4
				}
			case "\x00H":
				if sammy[0].direction != 2 {
					sammy[0].direction = 1
				}
			case "\x00P":
				if sammy[0].direction != 1 {
					sammy[0].direction = 2
				}
			case "\x00K":
				if sammy[0].direction != 4 {
					sammy[0].direction = 3
				}
			case "\x00M":
				if sammy[0].direction != 3 {
					sammy[0].direction = 4
				}
			case "p", "P":
				SpacePause(" Game Paused ... Push Space  ")
			}

			for a := 0; a < numPlayers; a++ {
				// Move snake
				switch sammy[a].direction {
				case 1:
					sammy[a].row--
				case 2:
					sammy[a].row++
				case 3:
					sammy[a].col--
				case 4:
					sammy[a].col++
				}

				// If snake hits number, respond accordingly
				if numberRow == (sammy[a].row+1)/2 && numberCol == sammy[a].col {
					Play("MBO0L16>CCCE")
					if sammy[a].length < (MaxSnakeLength - 30) {
						sammy[a].length = sammy[a].length + number*4
					}
					sammy[a].score = sammy[a].score + number
					PrintScore(numPlayers, sammy[0].score, sammy[1].score, sammy[0].lives, sammy[1].lives)
					number++
					if number == 10 {
						EraseSnake(sammy, sammyBody, 0)
						EraseSnake(sammy, sammyBody, 1)
						Locate(numberRow, numberCol)
						Print(" ")
						Level(NextLevel, sammy)
						PrintScore(numPlayers, sammy[0].score, sammy[1].score, sammy[0].lives, sammy[1].lives)
						SpacePause("     Level" + Str(curLevel) + ",  Push Space")
						if numPlayers == 1 {
							sammy[1].row = 0
						}
						number = 1
						if diff == "Y" {
							speed -= 10
						}
						curSpeed = speed
					}
					noNum = true
					if curSpeed < 1 {
						curSpeed = 1
					}
				}
			}

			for a := 0; a < numPlayers; a++ {
				// If player runs into any point, or the head of the other snake, it dies.
				if PointIsThere(sammy[a].row, sammy[a].col, colorTable[3]) || (sammy[0].row == sammy[1].row && sammy[0].col == sammy[1].col) {
					Play("MBO0L32EFGEFDC")
					ColorBg(colorTable[3])
					Locate(numberRow, numberCol)
					Print(" ")

					playerDied = true
					sammy[a].alive = false
					sammy[a].lives = sammy[a].lives - 1

					// Otherwise, move the snake, and erase the tail
				} else {
					sammy[a].head = (sammy[a].head + 1) % MaxSnakeLength
					sammyBody[sammy[a].head][a].row = sammy[a].row
					sammyBody[sammy[a].head][a].col = sammy[a].col
					tail := (sammy[a].head + MaxSnakeLength - sammy[a].length) % MaxSnakeLength
					Set(sammyBody[tail][a].row, sammyBody[tail][a].col, colorTable[3])
					sammyBody[tail][a].row = 0
					Set(sammy[a].row, sammy[a].col, sammy[a].color)
				}
			}
		}

		curSpeed = speed // Reset speed to initial value

		for a := 0; a < numPlayers; a++ {
			EraseSnake(sammy, sammyBody, a)

			// If dead, then erase snake in really cool way
			if !sammy[a].alive {
				// Update score
				sammy[a].score -= 10
				PrintScore(numPlayers, sammy[0].score, sammy[1].score, sammy[0].lives, sammy[1].lives)

				if a == 0 {
					SpacePause(" Sammy Dies! Push Space! --->")
				} else {
					SpacePause(" <---- Jake Dies! Push Space ")
				}
			}
		}

		Level(SameLevel, sammy)
		PrintScore(numPlayers, sammy[0].score, sammy[1].score, sammy[0].lives, sammy[1].lives)
	}
}

// PointIsThere checks the global arena array to see if the boolean flag is set.
func PointIsThere(row, col, color int) bool {
	if row == 0 {
		return false
	}
	return arena[row-1][col-1].color != color
}

// PrintScore prints players scores and number of lives remaining.
func PrintScore(numPlayers, score1, score2, lives1, lives2 int) {
	Color(15, colorTable[3])
	if numPlayers == 2 {
		Locate(1, 1)
		PrintUsing("%7d00  Lives: %d  <--JAKE", score2, lives2)
	}
	Locate(1, 49)
	PrintUsing("SAMMY-->  Lives: %d     %7d00", lives1, score1)
}

// Set sets row and column on playing field to given color to facilitate moving
// of snakes around the field.
func Set(row, col, color int) {
	if row == 0 {
		return
	}
	// assign color to arena
	arena[row-1][col-1].color = color
	// Get real row of pixel
	realRow := arena[row-1][col-1].realRow
	// Deduce whether pixel is on top▀, or bottom▄
	topFlag := (arena[row-1][col-1].sister+1)/2 != 0
	// Get arena row of sister
	sisterRow := row + arena[row-1][col-1].sister
	// Determine sister's color
	sisterColor := arena[sisterRow-1][col-1].color

	Locate(realRow, col)

	if color == sisterColor {
		// If both points are same
		Color(color, color)
		Print("█")
	} else {
		// Since you cannot have bright backgrounds determine
		// the best combo to use.
		if topFlag {
			if color > 7 {
				Color(color, sisterColor)
				Print("▀")
			} else {
				Color(sisterColor, color)
				Print("▄")
			}
		} else {
			if color > 7 {
				Color(color, sisterColor)
				Print("▄")
			} else {
				Color(sisterColor, color)
				Print("▀")
			}
		}
	}
}

// SpacePause pauses game play and waits for space bar to be pressed before
// continuing.
func SpacePause(text string) {
	Color(colorTable[4], colorTable[5])
	Center(11, "█▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀█")
	Center(12, "█ "+Left(text+Space(29), 29)+" █")
	Center(13, "█▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄█")
	for InKey() != "" {
	}
	for InKey() != " " {
	}
	Color(15, colorTable[3])

	for i := 21; i <= 26; i++ { // Restore the screen background
		for j := 24; j <= 56; j++ {
			Set(i, j, arena[i-1][j-1].color)
		}
	}
}

// SparklePause creates a flashing border for the intro screen.
func SparklePause() {
	Color(4, 0)
	s := "*    *    *    *    *    *    *    *    *    *    *    *    *    *    *    *    *    "

	// Clear keyboard buffer
	for InKey() != "" {
	}

	for InKey() == "" {
		for a := 1; a <= 5; a++ {
			// Print horizontal sparkles
			Locate(1, 1)
			Print(Mid(s, a, 80))
			Locate(22, 1)
			Print(Mid(s, 6-a, 80))

			// Print vertical sparkles
			for b := 2; b <= 21; b++ {
				c := (a + b) % 5
				if c == 1 {
					Locate(b, 80)
					Print("*")
					Locate(23-b, 1)
					Print("*")
				} else {
					Locate(b, 80)
					Print(" ")
					Locate(23-b, 1)
					Print(" ")
				}
			}
			SleepMillis(50)
		}
	}
}

// StillWantsToPlay determines if users want to play game again.
func StillWantsToPlay() bool {
	Color(colorTable[4], colorTable[5])
	Center(10, "█▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀█")
	Center(11, "█       G A M E   O V E R       █")
	Center(12, "█                               █")
	Center(13, "█      Play Again?   (Y/N)      █")
	Center(14, "█▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄█")

	for InKey() != "" {
	}
	kbd := ""
	for kbd != "Y" && kbd != "N" {
		kbd = UCase(InKey())
	}

	Color(15, colorTable[3])
	Center(10, "                                 ")
	Center(11, "                                 ")
	Center(12, "                                 ")
	Center(13, "                                 ")
	Center(14, "                                 ")

	if kbd == "Y" {
		return true
	}

	Color(7, 0)
	Cls()
	return false
}
