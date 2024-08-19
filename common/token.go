package common

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
)

var TOKEN_SAVED_FILE = ".timetask-token"
var forgetLoginReminder = "Did you forget initialsing the CLI, with `login`?"

func getSaveTokenPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	saveTokenPath := path.Join(homeDir, TOKEN_SAVED_FILE)
	if len(saveTokenPath) == 0 {
		return "", errors.New("cannot determine save token path, path is nil")
	}

	return saveTokenPath, nil
}

func SaveUserToken(token string) error {

	saveTokenPath, err := getSaveTokenPath()
	if err != nil {
		return err
	}

	f, err := os.Create(saveTokenPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(token)

	if err != nil {
		return err
	}

	fmt.Printf("User token saved to '%s'\n", saveTokenPath)

	return nil
}

func GetUserToken() (string, error) {
	saveTokenPath, err := getSaveTokenPath()
	if err != nil {
		return "", err
	}

	f, err := os.Open(saveTokenPath)
	if err != nil {
		errMsg := err.Error()
		if _, ok := err.(*os.PathError); ok {
			errMsg = "cannot open token file"

		}
		return "", fmt.Errorf("%s. %s", errMsg, forgetLoginReminder)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	line, _, err := r.ReadLine()

	if err != nil {
		return "", err
	}

	token := string(line)
	if len(strings.Trim(token, " ")) == 0 {
		return "", fmt.Errorf("token is invalid. %s", forgetLoginReminder)
	}
	return token, nil
}
