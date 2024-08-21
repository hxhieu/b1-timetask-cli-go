package common

import (
	"os"

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

func NewTaskParser(inputFile *string) (*TaskCsvParser, error) {
	f := TASK_CSV_FILE
	if inputFile != nil {
		f = *inputFile
	}
	parser := &TaskCsvParser{}

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

	parser.Tasks = tasks

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
		"WorkType",
		"WorkType ID",
		"Project ID",
		"Billable",
		"Mon",
		"Tue",
		"Wed",
		"Thu",
		"Fri",
		"Sat",
		"Sun",
	})
	for _, task := range p.Tasks {
		if task == nil {
			t.AppendRow([]interface{}{"-"})
		} else {
			t.AppendRow([]interface{}{
				task.Task,
				task.Id,
				task.Desc,
				task.Title,
				task.WorkType,
				task.WorkTypeId,
				task.ProjectId,
				task.Billable,
				task.Mon,
				task.Tue,
				task.Wed,
				task.Thu,
				task.Fri,
				task.Sat,
				task.Sun,
			})
		}
		t.AppendSeparator()
	}
	t.Render()
}

// Convert the day hour properties to day indexed array
func (i *TimeTaskInput) Hours() []float32 {
	return []float32{i.Mon, i.Tue, i.Wed, i.Thu, i.Fri, i.Sat, i.Sun}
}
