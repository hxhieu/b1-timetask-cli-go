package common

import (
	"os"

	"github.com/gocarina/gocsv"
)

type TaskCsvParser struct {
	inputFile *string
	Tasks     []*TimeTaskInput
}

// Default input file
var DEFAULT_TASK_CSV_FILE = "tasks.csv"

func NewCsvTaskParser(inputFile *string) *TaskCsvParser {
	return &TaskCsvParser{
		inputFile: inputFile,
	}
}

func (p *TaskCsvParser) GetTasks() ([]*TimeTaskInput, error) {
	f := DEFAULT_TASK_CSV_FILE
	if p.inputFile != nil {
		f = *p.inputFile
	}

	// Read the file
	buffer, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}

	allTasks := []*TimeTaskInput{}
	err = gocsv.UnmarshalBytes(buffer, &allTasks)
	if err != nil {
		return nil, err
	}

	tasks := []*TimeTaskInput{}

	// Only take tasks having hours
	for _, t := range allTasks {
		if t.TotalHours() > 0 {
			tasks = append(tasks, t)
		}
	}

	return tasks, nil
}
