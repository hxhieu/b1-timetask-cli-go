package intervals_api

import (
	"encoding/json"
	"os"
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
		if buf, err := os.ReadFile(debugFile); err == nil {
			debugData := []ProjectWorkType{}
			if err = json.Unmarshal(buf, &debugData); err == nil {
				return &debugData, nil
			}
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
		if debugData, err := json.MarshalIndent(result, "", "\t"); err == nil {
			os.WriteFile(debugFile, debugData, 0644)
		}
	}

	return &result, nil
}
