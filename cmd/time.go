package cmd

type timeCreateCmd struct {
	InputFile *string `optional:"" help:"The input CSV file. Optional: Default to 'tasks.csv'"`
}

type timeClearCmd struct {
}

type timeCmd struct {
	Create timeCreateCmd `cmd:"" help:"Create remote time tasks, from the input, for the current week."`
	Clear  timeClearCmd  `cmd:"" help:"Clean up the remote time tasks, for the current week."`
}
