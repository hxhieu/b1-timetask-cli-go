package cmd

type CLIContext struct {
	Debug bool
}

type CLI struct {
	Debug bool `help:"Enable debug mode."`

	Login loginCmd `cmd:"" help:"Initialise the CLI, by logging in with an user token."`
	Time  timeCmd  `cmd:"" help:"Time related sub commands"`
}
