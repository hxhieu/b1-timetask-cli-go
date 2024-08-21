package intervals_api

import (
	"encoding/json"
	"errors"

	"github.com/hxhieu/b1-timetask-cli-go/common"
)

type CreateTimeRequest struct {
	Billable    string  `json:"billable"`
	Date        string  `json:"date"`
	Time        float32 `json:"time"`
	PersonId    string  `json:"personid"`
	TaskId      string  `json:"taskid"`
	ProjectId   string  `json:"projectid"`
	WorkTypeId  string  `json:"worktypeid"`
	Description string  `json:"description"`

	// For display
	WorkType string `json:"-"`
}

func (c *Client) CreateTime(createTime *CreateTimeRequest) error {
	if createTime == nil {
		return errors.New("cannot create nil time task")
	}
	_, err := c.post("time", *createTime)
	if err != nil {
		return err
	}
	return nil
}

// Map from CSV input
func (t *CreateTimeRequest) ParseInput(input *common.TimeTaskInput) error {
	if input == nil {
		return errors.New("the source input is nil")
	}

	// Could use https://github.com/mitchellh/mapstructure
	// but is is archived and probably overkill for our purpose.
	// So we just use a crude JSON de/serialising
	src, err := json.Marshal(*input)
	if err != nil {
		return err
	}
	err = json.Unmarshal(src, t)

	// Description, default to task title
	if len(t.Description) == 0 {
		t.Description = input.Title
	}

	t.WorkType = input.WorkType

	return nil
}
