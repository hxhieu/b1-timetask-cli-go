package cmd

import (
	"errors"
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
			Notation:         " jobs",
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

func runPrepSteps() (*intervals_api.Client, *common.TaskCsvParser, error) {
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
		job = newJobTrack(pw, "Check task inputs", 3)

		// Concat IDs, to pass to the remoter server
		var tasks string
		var projects string

		if parser, err := common.NewTaskParser(); err == nil {
			job.Increment(1)
			for _, t := range parser.Tasks {
				if t != nil {
					tasks += t.Task + ","
				}
			}
			tasks = strings.TrimSuffix(tasks, ",")

			// Fetch needed details from remote

			if remoteTasks, err := client.FetchTasks(tasks); err == nil {
				job.Increment(1)
				// TODO: Optimise this? nested loops here
				for _, remoteTask := range *remoteTasks {
					// Also build the project IDs list
					projects += remoteTask.ProjectId + ","
					for _, localTask := range parser.Tasks {
						if localTask != nil && localTask.Task == remoteTask.LocalId {
							localTask.ProjectId = remoteTask.ProjectId
							localTask.Id = remoteTask.Id
							localTask.Title = remoteTask.Title
						}
					}
				}
				projects = strings.TrimSuffix(projects, ",")

				// Fetch work types, because they are setup per project
				if remoteWorkTypes, err := client.FetchProjectWorkTypes(projects); err == nil {
					job.Increment(1)
					// TODO: Optimise this? nested loops here
					for _, localTask := range parser.Tasks {
						var defaultWorkType *string
						// Walk all work types of the same project
						for _, remoteWorkType := range *remoteWorkTypes {
							if localTask != nil && localTask.ProjectId == remoteWorkType.ProjectId {
								localTask.WorkTypeId = remoteWorkType.WorkTypeId
								defaultWorkType = &remoteWorkType.WorkType
								// Finally found the match by work type name
								if localTask.WorkType == remoteWorkType.WorkType {
									defaultWorkType = nil
									break
								}
							}
						}
						// Input work type not found, using the last one we found from remote server
						if defaultWorkType != nil {
							localTask.WorkType = *defaultWorkType
						}
					}
				} else {
					setJobError(job, err)
				}

				// All done
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

	if job.IsErrored() {
		return nil, nil, errors.New("one or more steps throwing errors")
	}

	taskParser.DebugPrint()
	console.Header("Press ENTER to process, or CTRL+C to terminate.")
	fmt.Scanln()

	return client, taskParser, nil
}

func runExecSteps(client *intervals_api.Client, tasks []*common.TimeTaskInput) error {
	return nil
}

func (c *runCmd) Run() error {
	// Prep checks
	client, taskParser, err := runPrepSteps()
	if err != nil {
		return err
	}

	// Real work
	err = runExecSteps(client, taskParser.Tasks)
	if err != nil {
		return err
	}

	return nil
}
