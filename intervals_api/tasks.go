package intervals_api

import (
	"encoding/json"
	"os"
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
	debugFile := ".debug_tasks.json"
	if c.debug {
		if buf, err := os.ReadFile(debugFile); err == nil {
			debugData := []Task{}
			if err = json.Unmarshal(buf, &debugData); err == nil {
				return &debugData, nil
			}
		}
	}

	body, err := c.get("task?localid=" + tasks)
	if err != nil {
		return nil, err
	}

	res := TasksReponse{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	result := res.Task
	if c.debug {
		if debugData, err := json.MarshalIndent(result, "", "\t"); err == nil {
			os.WriteFile(debugFile, debugData, 0644)
		}
	}

	return &result, nil
}
