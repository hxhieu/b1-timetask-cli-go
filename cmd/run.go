package cmd

import (
	"errors"
	"fmt"

	"github.com/hxhieu/b1-timetask-cli-go/common"
)

type runCmd struct {
}

func (c *runCmd) Run() error {
	token, err := common.GetUserToken()

	if err != nil {
		return err
	}

	if token == nil || len(*token) == 0 {
		return errors.New("Could not get the user token. Did we forget to initialise the CLI?")
	}

	fmt.Println(*token)

	return nil
}
