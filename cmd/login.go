package cmd

import (
	"github.com/hxhieu/b1-timetask-cli-go/common"
)

type loginCmd struct {
	Token string `required:"" help:"The user token to setup the CLI. Refer to this link on how to get the token https://www.myintervals.com/api/authentication.php" short:"t"`
}

func (c *loginCmd) Run() error {
	return common.SaveUserToken(c.Token)
}
