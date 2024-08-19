package cmd

type CLI struct {
	Debug bool `help:"Enable debug mode."`

	Login loginCmd `cmd:"" help:"Initialise the CLI, by logging in with an user token."`
	Run   runCmd   `cmd:"" help:"Run the CLI."`
}
