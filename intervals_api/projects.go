package intervals_api

import (
	"encoding/json"

	"github.com/hxhieu/b1-timetask-cli-go/debug"
)

type ProjectWorkType struct {
	Id         string `json:"id"`
	ProjectId  string `json:"projectid"`
	WorkTypeId string `json:"worktypeid"`
	WorkType   string `json:"worktype"`
	Active     string `json:"active"`
}

type ProjectWorkTypeResponse struct {
	ProjectWorkType []ProjectWorkType `json:"projectworktype"`
}

func (c *Client) FetchProjectWorkTypes(projects string, active ...string) (*[]ProjectWorkType, error) {
	debugFile := ".debug_project-worktype.json"
	if c.debug {
		if debugData := debug.LoadDataFile[[]ProjectWorkType](debugFile); debugData != nil {
			return debugData, nil
		}
	}

	activeFlag := "t"
	if len(active) > 0 {
		activeFlag = active[0]
	}
	body, err := c.get("projectworktype?active=" + activeFlag + "&projectid=" + projects)
	if err != nil {
		return nil, err
	}

	res := ProjectWorkTypeResponse{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	result := res.ProjectWorkType

	if c.debug {
		debug.WriteDataFile(debugFile, result)
	}

	return &result, nil
}
