package common

import (
	"os"

	"github.com/gocarina/gocsv"
	"github.com/jedib0t/go-pretty/v6/table"
)

type TimeTaskInput struct {
	Id       string  `csv:"task"`
	Desc     string  `csv:"desc"`
	Billable string  `csv:"billable"`
	Mon      float32 `csv:"mon"`
	Tue      float32 `csv:"tue"`
	Wed      float32 `csv:"wed"`
	Thu      float32 `csv:"thu"`
	Fri      float32 `csv:"fri"`
	Sat      float32 `csv:"sat"`
	Sun      float32 `csv:"sun"`

	// From remote source

	TaskDesc string
}

type TaskCsvParser struct {
	Tasks []*TimeTaskInput
}

// TODO: CLI params?
var TASK_CSV_FILE = "tasks.csv"

func NewTaskParser() (*TaskCsvParser, error) {
	parser := &TaskCsvParser{}

	// Read the file
	buffer, err := os.ReadFile(TASK_CSV_FILE)
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
	t.AppendHeader(table.Row{"Task ID", "Desc", "Billable", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"})
	for _, task := range p.Tasks {
		if task == nil {
			t.AppendRow([]interface{}{"-"})
		} else {
			t.AppendRow([]interface{}{task.Id, task.Desc, task.Billable, task.Mon, task.Tue, task.Wed, task.Thu, task.Fri, task.Sat, task.Sun})
			// console.InfoLn(fmt.Sprintf(
			// 	"Task id: %s, Desc: %s, Billable: %s, Mon: %f, Tue: %f, Wed: %f, Thu: %f, Fri: %f, Sat: %f, Sun: %f",
			// 	task.Id,
			// 	task.Desc,
			// 	task.Billable,
			// 	task.Mon,
			// 	task.Tue,
			// 	task.Wed,
			// 	task.Thu,
			// 	task.Fri,
			// 	task.Sat,
			// 	task.Sun,
			// ))
		}
	}
	t.Render()
}
