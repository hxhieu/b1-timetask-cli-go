package main

import (
	"github.com/alecthomas/kong"
	"github.com/hxhieu/b1-timetask-cli-go/cmd"
)

type context struct {
	Debug bool
}

func main() {
	cli := cmd.CLI{}
	ctx := kong.Parse(
		&cli,
		kong.Name("b1-timetask-cli"),
		kong.Description("A CLI tool to semi-automate the creation of time task."),
	)
	err := ctx.Run(&context{Debug: cli.Debug})
	ctx.FatalIfErrorf(err)
}
