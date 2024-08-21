package intervals_api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	baseUrl string
	token   string
	debug   bool
}

// Single instance of the http/net client
// and it is thread safe to do so
var httpClient = &http.Client{}

func (c *Client) setAuth(req *http.Request) {
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:x", c.token)))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", auth))
	req.Header.Set("Accept", "application/json")
}

func (c *Client) doRequest(req *http.Request) (*[]byte, error) {
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

	return &body, nil
}

func (c *Client) get(endpoint string) (*[]byte, error) {
	req, err := http.NewRequest("GET", c.baseUrl+endpoint, nil)
	if err != nil {
		return nil, err
	}
	c.setAuth(req)

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *Client) post(endpoint string, data any) (*[]byte, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.baseUrl+endpoint, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	c.setAuth(req)
	req.Header.Set("Content-Type", "application/json")

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func New(token string, debugMode ...bool) *Client {
	debug := false
	if len(debugMode) > 0 {
		debug = debugMode[0]
	}
	api := Client{
		baseUrl: "https://api.myintervals.com/",
		token:   token,
		debug:   debug,
	}
	return &api
}
