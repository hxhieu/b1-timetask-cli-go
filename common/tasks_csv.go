package common

import (
	"os"

	"github.com/fatih/color"
	"github.com/gocarina/gocsv"
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

type TaskCsvParser struct {
	Tasks []*TimeTaskInput
}

// Default input file
var TASK_CSV_FILE = "tasks.csv"

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

func NewTaskParser(inputFile *string) (*TaskCsvParser, error) {
	f := TASK_CSV_FILE
	if inputFile != nil {
		f = *inputFile
	}
	parser := &TaskCsvParser{
		Tasks: make([]*TimeTaskInput, 0),
	}

	// Read the file
	buffer, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}

	tasks := []*TimeTaskInput{}
	err = gocsv.UnmarshalBytes(buffer, &tasks)
	if err != nil {
		return nil, err
	}

	// Only take tasks having hours
	for _, t := range tasks {
		if t.TotalHours() > 0 {
			parser.Tasks = append(parser.Tasks, t)
		}
	}

	return parser, nil
}

func (p *TaskCsvParser) DebugPrint() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{
		"Task",
		"ID",
		"Desc",
		"Title",
		"Billable",
		"WorkType ID",
		"Project ID",
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
	for _, task := range p.Tasks {
		if task == nil {
			t.AppendRow([]interface{}{"-"})
		} else {
			t.AppendRow([]interface{}{
				task.Task,
				task.Id,
				task.Desc,
				task.Title,
				task.Billable,
				task.WorkTypeId,
				task.ProjectId,
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
		"",
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

// Convert the day hour properties to day indexed array
func (i *TimeTaskInput) Hours() []float32 {
	return []float32{i.Mon, i.Tue, i.Wed, i.Thu, i.Fri, i.Sat, i.Sun}
}

// Total hours for a task from input
func (i *TimeTaskInput) TotalHours() float32 {
	return i.Mon + i.Tue + i.Wed + i.Thu + i.Fri + i.Sat + i.Sun
}
