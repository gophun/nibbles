package basic

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
)

var (
	screen         tcell.Screen
	fg, bg         int
	locRow, locCol int
)

func Screen(mode int) {
	if mode != 0 {
		panic("unsupported screen mode")
	}
	var err error
	screen, err = tcell.NewScreen()
	if err != nil {
		log.Fatal("could not create text mode screen")
	}
	err = screen.Init()
	if err != nil {
		log.Fatal("could not initialize text mode")
	}
}

func Reset() {
	screen.Fini()
}

func Width(columns, rows int) {
	// not implemented
}

func Cls() {
	clrRect(1, 1, 80, 25)
	Locate(1, 1)
}

func Color(foreground, background int) {
	ColorFg(foreground)
	ColorBg(background)
}

func ColorFg(foreground int) {
	fg = foreground
}

func ColorBg(background int) {
	bg = background
}

func Locate(row, column int) {
	locCol = column
	locRow = row
	screen.Show()
}

func style() tcell.Style {
	var s tcell.Style
	s = s.Background(color(bg))
	s = s.Foreground(color(fg))
	return s
}

func color(color int) tcell.Color {
	switch color {
	case 0:
		return tcell.ColorBlack
	case 1:
		return tcell.ColorNavy
	case 2:
		return tcell.ColorGreen
	case 3:
		return tcell.ColorTeal
	case 4:
		return tcell.ColorMaroon
	case 5:
		return tcell.ColorPurple
	case 6:
		return tcell.ColorOlive
	case 7:
		return tcell.ColorSilver
	case 8:
		return tcell.ColorGray
	case 9:
		return tcell.ColorBlue
	case 10:
		return tcell.ColorLime
	case 11:
		return tcell.ColorAqua
	case 12:
		return tcell.ColorRed
	case 13:
		return tcell.ColorFuchsia
	case 14:
		return tcell.ColorYellow
	case 15:
		return tcell.ColorWhite
	}
	return tcell.ColorDefault
}

func clrRect(x1, y1, x2, y2 int) {
	st := style()
	for y := y1; y <= y2; y++ {
		for x := x1; x <= x2; x++ {
			screen.SetContent(x-1, y-1, ' ', nil, st)
		}
	}
}

func Print(s string) {
	st := style()
	for _, ch := range s {
		screen.SetContent(locCol-1, locRow-1, ch, nil, st)
		locCol++
	}
	Locate(locRow, locCol)
}

func PrintUsing(format string, a ...interface{}) {
	Print(fmt.Sprintf(format, a...))
}

func Input(prompt string) string {
	Print(prompt + "? ")
	var buf []rune
	for {
		screen.ShowCursor(locCol-1, locRow-1)
		screen.Show()
		ch := readKey()
		if ch == rune(tcell.KeyEnter) {
			break
		}
		if ch == 127 {
			// Backspace
			if len(buf) == 0 {
				continue
			}
			buf = buf[:len(buf)-1]
			if locCol == 1 {
				continue
			}
			locCol--
			Print(" ")
			Locate(locRow, locCol-1)
			continue
		}
		buf = append(buf, ch)
		Print(string(ch))
	}
	screen.HideCursor()
	return string(buf)
}

func InKey() string {
	if !screen.HasPendingEvent() {
		return ""
	}
	event := screen.PollEvent()
	switch ev := event.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyRune:
			return string(ev.Rune())
		case tcell.KeyUp:
			return "\x00H"
		case tcell.KeyDown:
			return "\x00P"
		case tcell.KeyRight:
			return "\x00M"
		case tcell.KeyLeft:
			return "\x00K"
		}
		return ev.Name()
	}
	return ""
}

func readKey() rune {
	for {
		event := screen.PollEvent()
		switch ev := event.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyRune {
				return ev.Rune()
			}
			return rune(ev.Key())
		}
	}
}
