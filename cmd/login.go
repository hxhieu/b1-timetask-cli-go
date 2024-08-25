package cmd

import (
	"github.com/hxhieu/b1-timetask-cli-go/common"
)

func (c *loginCmd) Run() error {
	return common.SaveUserToken(c.Token)
}
