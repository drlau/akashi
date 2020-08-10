package utils

import (
	"io"
	"os"

	"github.com/mattn/go-colorable"
	"github.com/mgutz/ansi"
)

var (
	Red        = ansi.ColorFunc("red")
	RedBold    = ansi.ColorFunc("red+b")
	Yellow     = ansi.ColorFunc("yellow")
	YellowBold = ansi.ColorFunc("yellow+b")
	Green      = ansi.ColorFunc("green")
	Bold       = ansi.ColorFunc("default+b")
)

func NewOutput(noColor bool) io.Writer {
	if noColor {
		return colorable.NewNonColorable(os.Stdout)
	}

	return colorable.NewColorable(os.Stdout)
}
