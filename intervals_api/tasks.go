package intervals_api

import (
	"encoding/json"

	"github.com/hxhieu/b1-timetask-cli-go/debug"
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
		if debugData := debug.LoadDataFile[[]Task](debugFile); debugData != nil {
			return debugData, nil
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
		debug.WriteDataFile(debugFile, result)
	}

	return &result, nil
}
