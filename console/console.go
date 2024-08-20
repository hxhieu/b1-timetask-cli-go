package console

import (
	"fmt"

	"github.com/fatih/color"
)

func Error(err string) {
	c := color.New(color.FgRed).Add(color.Bold)
	c.Print(err)
}

func ErrorLn(err string) {
	Error(fmt.Sprintf("%s\n", err))
}

func Success(msg string) {
	c := color.New(color.FgHiGreen).Add(color.BgWhite).Add(color.Bold)
	c.Print(msg)
}

func SuccessLn(msg string) {
	Success(fmt.Sprintf("%s\n", msg))
}

func Info(msg string) {
	c := color.New(color.FgWhite)
	c.Print(msg)
}

func InfoLn(msg string) {
	Info(fmt.Sprintf("%s\n", msg))
}

func Warn(msg string) {
	c := color.New(color.FgHiYellow)
	c.Print(msg)
}

func WarnLn(msg string) {
	Warn(fmt.Sprintf("%s\n", msg))
}

func Header(header string) {
	c := color.New(color.FgWhite).Add(color.Bold)
	c.Printf(" %s \n", header)
}
