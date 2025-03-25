package console

import (
	"fmt"

	"github.com/fatih/color"
)

func Error(err string) {
	c := color.New(color.FgRed).Add(color.Bold)
	fmt.Print(c.Sprint(err))
}

func ErrorLn(err string) {
	Error(fmt.Sprintf("%s\n", err))
}

func Success(msg string) {
	c := color.New(color.FgHiGreen).Add(color.BgWhite).Add(color.Bold)
	fmt.Print(c.Sprint(msg))
}

func SuccessLn(msg string) {
	Success(fmt.Sprintf("%s\n", msg))
}

func Info(msg string) {
	c := color.New(color.FgWhite)
	fmt.Print(c.Sprint(msg))
}

func InfoLn(msg string) {
	Info(fmt.Sprintf("%s\n", msg))
}

func Warn(msg string) {
	c := color.New(color.FgHiYellow)
	fmt.Print(c.Sprint(msg))
}

func WarnLn(msg string) {
	Warn(fmt.Sprintf("%s\n", msg))
}

func Header(header string) {
	c := color.New(color.FgWhite).Add(color.Bold)
	fmt.Print(c.Sprintf("%s\n", header))
}

func PrintWeekRange(start, end string, weekOffset int) {
	c := color.New(color.FgYellow).Add(color.Bold)
	warnColor := color.New(color.FgGreen).Add(color.Bold)
	warnText := "For Current week"
	if weekOffset != 0 {
		warnColor = color.New(color.FgRed).Add(color.Bold)
		warnText = fmt.Sprintf("For %d week(s) ", weekOffset)
		if weekOffset > 0 {
			warnText += "from now"
		} else {
			warnText += "in the past"
		}
	}
	fmt.Print(c.Sprintf("%s %s - %s\n", warnColor.Sprint(warnText), c.Sprint(start), c.Sprint(end)))
}
