package intervals

import (
	"encoding/json"
	"fmt"

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

	// Default lmit is 10, so we need to override it with something bigger
	limit := 100
	body, err := c.get(fmt.Sprintf("task?localid=%s&limit=%d", tasks, limit))
	if err != nil {
		return nil, err
	}

	res := TasksReponse{}
	err = json.Unmarshal(*body, &res)
	if err != nil {
		return nil, err
	}

	result := res.Task

	if c.debug {
		debug.WriteDataFile(debugFile, result)
	}

	return &result, nil
}
