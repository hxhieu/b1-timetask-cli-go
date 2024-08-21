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

type runCmd struct {
}

type runPrepResult struct {
	userId     string
	client     *intervals_api.Client
	taskParser *common.TaskCsvParser
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

func runPrepSteps() (*runPrepResult, error) {
	// Shared vars between steps
	result := &runPrepResult{}

	// instantiate a Progress Writer and set up the options
	pw := progress.NewWriter()
	setDefaultProgress(&pw)

	go pw.Render()

	// Check and fetch user job
	job := newJobTrack(pw, "Check user", 2)

	if token, err := common.GetUserToken(); err == nil {
		job.Increment(1)

		// API client
		result.client = intervals_api.New(token)

		// Fetch the user
		if me, err := result.client.Me(); err == nil {
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

			if remoteTasks, err := result.client.FetchTasks(tasks); err == nil {
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
				if remoteWorkTypes, err := result.client.FetchProjectWorkTypes(projects); err == nil {
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
				result.taskParser = parser

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
		return nil, errors.New("one or more steps throwing errors")
	}

	result.taskParser.DebugPrint()
	console.Header("Press ENTER to process, or CTRL+C to terminate.")
	fmt.Scanln()

	return result, nil
}

func runExecSteps(prepResult *runPrepResult) error {
	weekDays := common.GetWeekRange(time.Now())

	for i, d := range weekDays {
		for _, input := range prepResult.taskParser.Tasks {
			if input == nil {
				continue
			}
			inputHours := input.Hours()

			// Create time task request
			createTime := &intervals_api.CreateTimeRequest{
				PersonId: prepResult.userId,
				Date:     d.Format("2006-01-02"),
				Time:     inputHours[i],
			}
			// Load from input
			createTime.ParseInput(input)

			// Only process valid time task, i.e. hours > 0
			if createTime.Time != 0 {
				prepResult.client.CreateTime(createTime)
			}
		}
	}

	return nil
}

func (c *runCmd) Run() error {
	// Prep checks
	prepResult, err := runPrepSteps()
	if err != nil {
		return err
	}

	// Real work
	err = runExecSteps(prepResult)
	if err != nil {
		return err
	}

	return nil
}
