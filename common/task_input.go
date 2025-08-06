package common

import (
	"os"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
)

type TimeTaskInput struct {
	Task     string  `csv:"task"`
	WorkType string  `csv:"work_type"`
	Desc     string  `csv:"desc" json:"description"`
	Billable string  `csv:"billable" json:"billable"`
	Mon      float32 `csv:"mon"`
	Tue      float32 `csv:"tue"`
	Wed      float32 `csv:"wed"`
	Thu      float32 `csv:"thu"`
	Fri      float32 `csv:"fri"`
	Sat      float32 `csv:"sat"`
	Sun      float32 `csv:"sun"`

	// From remote source

	Id         string `json:"taskid"`
	Title      string
	ProjectId  string `json:"projectid"`
	WorkTypeId string `json:"worktypeid"`
}

func formatTotalHours(h float32) string {
	c := color.New(color.Bold)
	if h == 8 {
		c.Add(color.FgHiGreen)
	} else if h > 8 {
		c.Add(color.FgYellow)
	} else if h < 8 {
		c.Add(color.FgHiRed)
	}
	return c.Sprintf("%.2f", h)
}

func formatWeekTotalHours(h float32) string {
	c := color.New(color.Bold)
	if h >= 40 {
		c.Add(color.FgHiGreen)
	} else if h <= 40 {
		c.Add(color.FgHiRed)
	}
	return c.Sprintf("TOTAL: %.2f", h)
}

// Convert the day hour properties to day indexed array
func (i *TimeTaskInput) Hours() []float32 {
	return []float32{i.Mon, i.Tue, i.Wed, i.Thu, i.Fri, i.Sat, i.Sun}
}

// Total hours for a task from input
func (i *TimeTaskInput) TotalHours() float32 {
	return i.Mon + i.Tue + i.Wed + i.Thu + i.Fri + i.Sat + i.Sun
}

func CalcMaxFieldsLen(tasks []*TimeTaskInput) (int, int) {
	maxTitleLength := 0
	maxWorkTypeLength := 0
	for _, input := range tasks {
		if input == nil {
			continue
		}

		if len(input.Title) > maxTitleLength {
			maxTitleLength = len(input.Desc)
		}

		if len(input.Desc) > maxTitleLength {
			maxTitleLength = len(input.Desc)
		}

		if len(input.WorkType) > maxWorkTypeLength {
			maxWorkTypeLength = len(input.WorkType)
		}
	}
	return maxTitleLength, maxWorkTypeLength
}

func PrintTimeTasks(tasks []*TimeTaskInput) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{
		"Task",
		"Desc",
		"Billable",
		"WorkType",
		"Mon",
		"Tue",
		"Wed",
		"Thu",
		"Fri",
		"Sat",
		"Sun",
	})

	dailySum := make([]float32, 7)
	weekSum := float32(0.0)
	for _, task := range tasks {
		if task == nil {
			t.AppendRow([]interface{}{"-"})
		} else {
			t.AppendRow([]interface{}{
				task.Task,
				task.Desc,
				task.Billable,
				task.WorkType,
				task.Mon,
				task.Tue,
				task.Wed,
				task.Thu,
				task.Fri,
				task.Sat,
				task.Sun,
			})
			dailySum[0] += task.Mon
			dailySum[1] += task.Tue
			dailySum[2] += task.Wed
			dailySum[3] += task.Thu
			dailySum[4] += task.Fri
			dailySum[5] += task.Sat
			dailySum[6] += task.Sun
			weekSum += task.Mon + task.Tue + task.Wed + task.Thu + task.Fri + task.Sat + task.Sun
		}
		t.AppendSeparator()
	}
	t.AppendRow([]interface{}{
		"",
		"",
		"",
		formatWeekTotalHours(weekSum),
		formatTotalHours(dailySum[0]),
		formatTotalHours(dailySum[1]),
		formatTotalHours(dailySum[2]),
		formatTotalHours(dailySum[3]),
		formatTotalHours(dailySum[4]),
		formatTotalHours(dailySum[5]),
		formatTotalHours(dailySum[6]),
	})
	t.Render()
}
