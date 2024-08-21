package cmd

type CLIContext struct {
	Debug bool
	Force bool
}

type CLI struct {
	Debug bool `help:"Enable debug mode." short:"d"`
	Force bool `help:"Supress all prompts." short:"f"`

	Gui   guiCmd   `cmd:"" help:"Launch a GUI app *NOT YET IMPLEMENTED*" default:"1"`
	Login loginCmd `cmd:"" help:"Initialise the CLI, by logging in with an user token."`
	Time  timeCmd  `cmd:"" help:"Time related sub commands"`
}
