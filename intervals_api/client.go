package intervals_api

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	baseUrl string
	token   string
}

var httpClient = &http.Client{}

func (c *Client) setAuth(req *http.Request) {
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:x", c.token)))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", auth))
	req.Header.Set("Accept", "application/json")
}

func (c *Client) get(endpoint string) ([]byte, error) {
	req, err := http.NewRequest("GET", c.baseUrl+endpoint, nil)
	if err != nil {
		return nil, err
	}
	c.setAuth(req)

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 300 {
		return nil, fmt.Errorf("%s - %s", res.Status, string(body))
	}

	return body, nil
}

func New(token string) *Client {
	api := Client{
		baseUrl: "https://api.myintervals.com/",
		token:   token,
	}
	return &api
}
