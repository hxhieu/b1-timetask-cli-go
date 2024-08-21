package intervals_api

import (
	"encoding/json"
	"fmt"
	"os"
)

type Me struct {
	Id        string `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"username"`
}

type MeResponse struct {
	Me []Me `json:"me"`
}

func (c *Client) Me() (*Me, error) {
	debugFile := ".debug_me.json"
	if c.debug {
		if buf, err := os.ReadFile(debugFile); err == nil {
			me := Me{}
			if err = json.Unmarshal(buf, &me); err == nil {
				return &me, nil
			}
		}
	}

	body, err := c.get("me")
	if err != nil {
		return nil, err
	}

	res := MeResponse{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	if len(res.Me) == 0 {
		return nil, fmt.Errorf("cannot find the match user with the provided token")
	}

	result := res.Me[0]

	if c.debug {
		if debugData, err := json.MarshalIndent(result, "", "\t"); err == nil {
			os.WriteFile(debugFile, debugData, 0644)
		}
	}

	return &result, nil
}
