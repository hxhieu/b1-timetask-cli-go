package cmd

import (
	"github.com/hxhieu/b1-timetask-cli-go/common"
	"github.com/hxhieu/b1-timetask-cli-go/console"
)

type runCmd struct {
}

func (c *runCmd) Run() error {
	token, err := common.GetUserToken()

	if err != nil {
		return err
	}

	console.Header(token)

	return nil
}
