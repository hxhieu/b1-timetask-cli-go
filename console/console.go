package console

import (
	"github.com/fatih/color"
)

func Error(err string) {
	c := color.New(color.FgRed).Add(color.Bold)
	c.Println(err)
}

func Success(msg string) {
	c := color.New(color.FgHiGreen).Add(color.BgWhite).Add(color.Bold)
	c.Println(msg)
}

func Header(header string) {
	c := color.New(color.FgBlack).Add(color.BgWhite).Add(color.Bold)
	c.Printf(" %s \n", header)
}
