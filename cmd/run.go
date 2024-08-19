package cmd

import (
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

	fmt.Println(token)

	return nil
}
