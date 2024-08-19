package intervals_api

import (
	"encoding/json"
	"fmt"
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

	return &res.Me[0], nil
}
