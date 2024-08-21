package cmd

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/hxhieu/b1-timetask-cli-go/common"
	"github.com/hxhieu/b1-timetask-cli-go/console"
	"github.com/hxhieu/b1-timetask-cli-go/intervals_api"
	"github.com/jedib0t/go-pretty/v6/progress"
)

type runPrepResult struct {
	userId string
	tasks  []*common.TimeTaskInput
}

func newJobTrack(pw progress.Writer, title string) *progress.Tracker {
	job := progress.Tracker{
		Message: title,
		Units: progress.Units{
			Notation:         " jobs",
			NotationPosition: progress.UnitsNotationPositionAfter,
		},
		DeferStart: false,
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
	p.SetAutoStop(true)
	// p.SetUpdateFrequency(time.Millisecond * 10)
	// p.Style().Visibility.ETA = false
	// p.Style().Visibility.ETAOverall = false
	// p.Style().Visibility.Percentage = false
	// p.Style().Visibility.Speed = true
	// p.Style().Visibility.SpeedOverall = false
	// p.Style().Visibility.Time = true
	// p.Style().Visibility.TrackerOverall = true
	// p.Style().Visibility.Value = true
	// p.Style().Visibility.Pinned = false
}

// Also mark tracker as error
func setJobError(job *progress.Tracker, err error) {
	job.UpdateMessage(fmt.Sprintf("%s -> %s", job.Message, color.RedString(err.Error())))
	job.MarkAsErrored()
}

func setJobSuccess(job *progress.Tracker, message string) {
	job.UpdateMessage(fmt.Sprintf("%s -> %s", job.Message, color.HiGreenString(message)))
	job.MarkAsDone()
}

func runPrepSteps(debug bool, inputFile *string) (*runPrepResult, *intervals_api.Client, error) {
	// Shared vars between steps
	result := &runPrepResult{}
	var client *intervals_api.Client
	var taskParser *common.TaskCsvParser

	// instantiate a Progress Writer and set up the options
	pw := progress.NewWriter()
	setDefaultProgress(&pw)

	go pw.Render()

	// Check and fetch user job
	job := newJobTrack(pw, "Check user")

	if token, err := common.GetUserToken(); err == nil {

		// API client
		client = intervals_api.New(token, debug)

		// Fetch the user
		if me, err := client.Me(); err == nil {
			result.userId = me.Id
			setJobSuccess(job, fmt.Sprintf("Found user: %s %s <%s>", me.FirstName, me.LastName, me.Email))
		} else {
			setJobError(job, err)
		}
	} else {
		setJobError(job, err)
	}

	// Parse and fetch tasks from CSV
	if !job.IsErrored() {
		job = newJobTrack(pw, "Check task inputs")

		// Concat IDs, to pass to the remoter server
		var tasks string
		var projects string

		if parser, err := common.NewTaskParser(inputFile); err == nil {
			for _, t := range parser.Tasks {
				if t != nil {
					tasks += t.Task + ","
				}
			}
			tasks = strings.TrimSuffix(tasks, ",")

			// Fetch needed details from remote

			if remoteTasks, err := client.FetchTasks(tasks); err == nil {
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
	time.Sleep(time.Millisecond * 100)
	for pw.IsRenderInProgress() {
	}

	if job.IsErrored() {
		return nil, nil, errors.New("one or more steps throwing errors")
	}

	taskParser.DebugPrint()
	result.tasks = taskParser.Tasks
	console.Header("Press ENTER to process, or CTRL+C to terminate.")
	fmt.Scanln()

	return result, client, nil
}

func runExecSteps(prepResult *runPrepResult, client *intervals_api.Client) error {
	// instantiate a Progress Writer and set up the options
	pw := progress.NewWriter()
	setDefaultProgress(&pw)

	go pw.Render()

	weekDays := common.GetWeekRange(time.Now())

	for i, d := range weekDays {
		for _, input := range prepResult.tasks {
			if input == nil {
				continue
			}
			inputHours := input.Hours()

			// Create time task request
			createTime := &intervals_api.TimeEntry{
				PersonId: prepResult.userId,
				Date:     common.DateToString(d),
				Time:     inputHours[i],
			}
			// Load from input
			createTime.ParseInput(input)

			// Only process valid time task, i.e. hours > 0
			if createTime.Time > 0 {
				// Run each as a goroutine
				go func() {
					job := newJobTrack(pw, fmt.Sprintf(
						"Creating time task: %s | %s | %s | %.2f hour(s) |",
						createTime.Date,
						createTime.Description,
						createTime.WorkType,
						createTime.Time,
					))
					if err := client.CreateTime(createTime); err == nil {
						setJobSuccess(job, "Created")
					} else {
						setJobError(job, err)
					}
				}()
			}
		}
	}

	// Render all jobs, until all done
	time.Sleep(time.Millisecond * 100)
	for pw.IsRenderInProgress() {
	}

	console.Header("All DONE! Press ENTER to exit.")
	fmt.Scanln()

	return nil
}

func (c *timeCreateCmd) Run(ctx CLIContext) error {
	// Prep checks
	prepResult, client, err := runPrepSteps(ctx.Debug, c.InputFile)
	if err != nil {
		return err
	}

	// Real work
	err = runExecSteps(prepResult, client)
	if err != nil {
		return err
	}

	return nil
}
