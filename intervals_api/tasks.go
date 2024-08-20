package intervals_api

import (
	"encoding/json"
	"strings"

	"github.com/hxhieu/b1-timetask-cli-go/common"
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

func (c *Client) FetchTasks(tasks []*common.TimeTaskInput) error {
	var ids string
	for _, t := range tasks {
		if t != nil {
			ids += t.Task + ","
		}
	}
	ids = strings.TrimSuffix(ids, ",")

	body, err := c.get("task?localid=" + ids)
	if err != nil {
		return err
	}

	res := TasksReponse{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return err
	}

	// TODO: Optimise this? nested loops here
	for _, remoteTask := range res.Task {
		for _, localTask := range tasks {
			if localTask != nil && localTask.Task == remoteTask.LocalId {
				localTask.ProjectId = remoteTask.ProjectId
				localTask.Id = remoteTask.Id
				localTask.Title = remoteTask.Title
			}
		}
	}

	return nil
}
