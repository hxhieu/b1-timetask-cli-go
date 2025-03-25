package cmd

import "embed"

var APP_TITLE = "B1 TimeTask Tools"

// Sub commands

type guiCmd struct {
}

type loginCmd struct {
	Token string `required:"" help:"The user token to setup the CLI. Refer to this link on how to get the token https://www.myintervals.com/api/authentication.php" short:"t"`
}

type timeCreateCmd struct {
	WeekOffset int     `help:"Number of weeks different from current week, negative offsets mean previous weeks." default:"0"`
	InputFile  *string `optional:"" help:"The input CSV file. Optional: Default to 'tasks.csv'"`
}

type timeClearCmd struct {
	WeekOffset int `help:"Number of weeks different from current week, negative offsets mean previous weeks." default:"0"`
}

type timeCmd struct {
	Create timeCreateCmd `cmd:"" help:"Create remote time tasks, from the input, for the current week."`
	Clear  timeClearCmd  `cmd:"" help:"Clean up the remote time tasks, for the current week."`
}

type CLIContext struct {
	Debug        bool
	Force        bool
	Experimental bool
	GuiAssets    *embed.FS
}

// Command scaffolding

type CLI struct {
	Debug        bool `help:"Enable debug mode." short:"d" env:"DEBUG"`
	Force        bool `help:"Supress all prompts." short:"f"`
	Experimental bool `help:"Run with experimental features." short:"x" env:"X_MODE"`

	Gui   guiCmd   `cmd:"" help:"Launch a GUI app *EXPERIMENTAL*" default:"1"`
	Login loginCmd `cmd:"" help:"Initialise the CLI, by logging in with an user token."`
	Time  timeCmd  `cmd:"" help:"Time related sub commands"`
}
