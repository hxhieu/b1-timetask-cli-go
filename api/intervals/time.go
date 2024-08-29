package intervals

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hxhieu/b1-timetask-cli-go/common"
	"github.com/hxhieu/b1-timetask-cli-go/debug"
)

type TimeEntry struct {
	Billable    string `json:"billable"`
	Date        string `json:"date"`
	Time        string `json:"time"`
	PersonId    string `json:"personid"`
	TaskId      string `json:"taskid"`
	ProjectId   string `json:"projectid"`
	WorkTypeId  string `json:"worktypeid"`
	Description string `json:"description"`

	// From the GETs
	Id             string `json:"id,omitempty"`
	WorkTypeRemote string `json:"worktype,omitempty"`

	// For local display
	WorkType string `json:"-"`
}

type GetTimeResponse struct {
	Time []TimeEntry `json:"time"`
}

func (c *Client) CreateTime(createTime *TimeEntry) error {
	if createTime == nil {
		return errors.New("cannot create nil time task")
	}
	_, err := c.post("time", *createTime)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetTimeEntries(start time.Time, end time.Time) (*[]TimeEntry, error) {
	debugFile := ".debug_time-entries.json"
	if c.debug {
		if debugData := debug.LoadDataFile[[]TimeEntry](debugFile); debugData != nil {
			return debugData, nil
		}
	}

	// Default lmit is 10, so we need to override it with something bigger
	limit := 100
	body, err := c.get(fmt.Sprintf(
		"time?datebegin=%s&dateend=%s&limit=%d",
		common.DateToString(start),
		common.DateToString(end), limit,
	))
	if err != nil {
		return nil, err
	}

	res := GetTimeResponse{}
	err = json.Unmarshal(*body, &res)
	if err != nil {
		return nil, err
	}

	result := res.Time

	if c.debug {
		debug.WriteDataFile(debugFile, result)
	}

	return &result, nil
}

func (c *Client) DeleteTimeEntry(id string) error {

	err := c.delete("time/" + id)
	if err != nil {
		return err
	}

	return nil
}

// Map from CSV input
func (t *TimeEntry) LoadFromInput(input *common.TimeTaskInput) error {
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
