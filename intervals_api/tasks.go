package intervals_api

import (
	"encoding/json"
)

type Task struct {
	Id        string `json:"id"`
	LocalId   string `json:"localid"`
	Title     string `json:"title"`
	ProjectId string `json:"projectid"`
}

type TasksReponse struct {
	Task []Task `json:"task"`
}

func (c *Client) FetchTasks(tasks string) (*[]Task, error) {
	body, err := c.get("task?localid=" + tasks)
	if err != nil {
		return nil, err
	}

	res := TasksReponse{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	return &res.Task, nil
}
