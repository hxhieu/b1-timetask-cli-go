package cmd

import (
	"fmt"

	"github.com/hxhieu/b1-timetask-cli-go/common"
	"github.com/hxhieu/b1-timetask-cli-go/console"
	"github.com/hxhieu/b1-timetask-cli-go/intervals_api"
)

type runCmd struct {
}

func (c *runCmd) Run() error {
	token, err := common.GetUserToken()
	if err != nil {
		return err
	}

	client := intervals_api.New(token)

	// Fetch the user
	console.Info("Fetching user...")
	me, err := client.Me()
	if err != nil {
		return err
	}
	console.Header(fmt.Sprintf("Found user: %s %s <%s>", me.FirstName, me.LastName, me.Email))

	// Parse the task inputs
	taskParser, err := common.NewTaskParser()
	if err != nil {
		return err
	}

	taskParser.DebugPrint()

	return nil
}
