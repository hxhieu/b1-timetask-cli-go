package common

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
)

var TOKEN_SAVED_FILE = ".b1-timetask-cli-token"

func getSaveTokenPath() (*string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	saveTokenPath := path.Join(homeDir, TOKEN_SAVED_FILE)
	if len(saveTokenPath) == 0 {
		return nil, errors.New("Cannot determine save token path, path is nil.")
	}

	return &saveTokenPath, nil
}

func SaveUserToken(token string) error {

	saveTokenPath, err := getSaveTokenPath()
	if err != nil {
		return err
	}

	f, err := os.Create(*saveTokenPath)
	defer f.Close()

	if err != nil {
		panic(err)
	}

	_, err = f.WriteString(token)

	if err != nil {
		return err
	}

	fmt.Printf("User token saved to '%s'\n", saveTokenPath)

	return nil
}

func GetUserToken() (*string, error) {
	saveTokenPath, err := getSaveTokenPath()
	if err != nil {
		return nil, err
	}

	f, err := os.Open(*saveTokenPath)
	defer f.Close()

	r := bufio.NewReader(f)
	line, _, err := r.ReadLine()

	if err != nil {
		return nil, err
	}

	token := string(line)
	return &token, nil
}
