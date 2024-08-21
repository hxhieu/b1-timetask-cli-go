package cmd

import "embed"

type CLIContext struct {
	Debug        bool
	Force        bool
	Experimental bool
	GuiAssets    *embed.FS
}

type CLI struct {
	Debug        bool `help:"Enable debug mode." short:"d" env:"DEBUG"`
	Force        bool `help:"Supress all prompts." short:"f"`
	Experimental bool `help:"Run with experimental features." short:"x" env:"X_MODE"`

	Gui   guiCmd   `cmd:"" help:"Launch a GUI app *NOT YET IMPLEMENTED*" default:"1"`
	Login loginCmd `cmd:"" help:"Initialise the CLI, by logging in with an user token."`
	Time  timeCmd  `cmd:"" help:"Time related sub commands"`
}
