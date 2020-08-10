package utils

import (
	"io"
	"os"

	"github.com/mattn/go-colorable"
	"github.com/mgutz/ansi"
)

var (
	Red    = ansi.ColorFunc("red")
	Yellow = ansi.ColorFunc("yellow")
	Green  = ansi.ColorFunc("green")
	Bold   = ansi.ColorFunc("default+b")
)

func NewOutput(noColor bool) io.Writer {
	if noColor {
		return colorable.NewNonColorable(os.Stdout)
	}

	return colorable.NewColorable(os.Stdout)
}
