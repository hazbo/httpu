package printer

import (
	"fmt"
)

const (
	ColorBlack int = iota + 30
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite
)

func Color(s string, c int) string {
	return fmt.Sprintf("\033[%d;1m%s\033[0m", c, s)
}
