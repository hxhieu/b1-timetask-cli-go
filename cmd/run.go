package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/hxhieu/b1-timetask-cli-go/common"
	"github.com/hxhieu/b1-timetask-cli-go/console"
	"github.com/hxhieu/b1-timetask-cli-go/intervals_api"
	"github.com/jedib0t/go-pretty/v6/progress"
)

type runCmd struct {
}

func newJobTrack(pw progress.Writer, title string, taskCount int) *progress.Tracker {
	job := progress.Tracker{
		Message: title,
		Total:   int64(taskCount),
		Units: progress.Units{
			Notation:         " tasks",
			NotationPosition: progress.UnitsNotationPositionAfter,
		},
		DeferStart: true,
	}
	pw.AppendTracker(&job)
	return &job
}

func setDefaultProgress(pw *progress.Writer) {
	if pw == nil {
		return
	}
	p := (*pw)
	p.SetTrackerPosition(progress.PositionRight)
}

// Also mark tracker as error
func setJobError(job *progress.Tracker, err error) {
	job.UpdateMessage(fmt.Sprintf("%s -> %s", job.Message, color.RedString(err.Error())))
	job.MarkAsErrored()
}

// Also increase the tracker count
func setJobSuccess(job *progress.Tracker, message string) {
	job.UpdateMessage(fmt.Sprintf("%s -> %s", job.Message, color.HiGreenString(message)))
	job.Increment(1)
}

func runPrepSteps() (*intervals_api.Client, *common.TaskCsvParser) {
	// Shared vars between steps
	var client *intervals_api.Client
	var taskParser *common.TaskCsvParser

	// instantiate a Progress Writer and set up the options
	pw := progress.NewWriter()
	setDefaultProgress(&pw)

	go pw.Render()

	// Check and fetch user job
	job := newJobTrack(pw, "Check user", 2)

	if token, err := common.GetUserToken(); err == nil {
		job.Increment(1)

		// API client
		client = intervals_api.New(token)

		// Fetch the user
		if me, err := client.Me(); err == nil {
			setJobSuccess(job, fmt.Sprintf("Found user: %s %s <%s>", me.FirstName, me.LastName, me.Email))
		} else {
			setJobError(job, err)
		}
	} else {
		setJobError(job, err)
	}

	// Parse and fetch tasks from CSV
	if !job.IsErrored() {
		job = newJobTrack(pw, "Check task inputs", 2)

		if parser, err := common.NewTaskParser(); err == nil {
			job.Increment(1)
			var ids string
			for _, t := range parser.Tasks {
				if t != nil {
					ids += t.Task + ","
				}
			}
			ids = strings.TrimSuffix(ids, ",")
			if remoteTasks, err := client.FetchTasks(ids); err == nil {
				// TODO: Optimise this? nested loops here
				for _, remoteTask := range *remoteTasks {
					for _, localTask := range parser.Tasks {
						if localTask != nil && localTask.Task == remoteTask.LocalId {
							localTask.ProjectId = remoteTask.ProjectId
							localTask.Id = remoteTask.Id
							localTask.Title = remoteTask.Title
						}
					}
				}
				setJobSuccess(job, "Found below task(s)")
				taskParser = parser
			} else {
				setJobError(job, err)
			}
		} else {
			setJobError(job, err)
		}
	}

	// Render all jobs, until all done
	for pw.IsRenderInProgress() {
		if pw.LengthActive() == 0 {
			pw.Stop()
		}
	}

	taskParser.DebugPrint()

	console.Header("Press ENTER to process, or CTRL+C to terminate.")
	fmt.Scanln()

	return client, taskParser
}

func (c *runCmd) Run() error {
	client, taskParser := runPrepSteps()

	// If we cannot get the variables from the prep steps then they have failed
	if client == nil || taskParser == nil {
		return nil
	}

	return nil
}
