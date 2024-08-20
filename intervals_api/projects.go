package intervals_api

import "encoding/json"

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

	return &res.ProjectWorkType, nil
}
