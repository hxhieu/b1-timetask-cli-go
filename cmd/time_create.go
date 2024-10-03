package cmd

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hxhieu/b1-timetask-cli-go/common"
	"github.com/hxhieu/b1-timetask-cli-go/console"
	"github.com/hxhieu/b1-timetask-cli-go/intervals_api"
	"github.com/jedib0t/go-pretty/v6/progress"
)

type createTimePrepResult struct {
	userId string
	tasks  []*common.TimeTaskInput
}

func createTimePrepSteps(ctx CLIContext, inputFile *string) (*createTimePrepResult, *intervals_api.Client, error) {
	// Shared vars between steps
	result := &createTimePrepResult{}
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
		client = intervals_api.New(token, ctx.Debug)

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
							// Truncate long title
							if len(localTask.Title) > 50 {
								localTask.Title = localTask.Title[:50]
							}
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

	if !ctx.Force {
		console.Header("Press ENTER to process, or CTRL+C to terminate.")
		fmt.Scanln()
	}

	return result, client, nil
}

func createTimeExecSteps(ctx CLIContext, prepResult *createTimePrepResult, client *intervals_api.Client) error {
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
			createTimeHours := inputHours[i]

			// Only process valid time task, i.e. hours > 0
			if createTimeHours <= 0 {
				continue
			}

			// Create time task request
			createTime := &intervals_api.TimeEntry{
				PersonId: prepResult.userId,
				Date:     common.DateToString(d),
				// Need a string for remote payload
				Time: fmt.Sprintf("%f", createTimeHours),
			}

			createTime.LoadFromInput(input)
			// Reset this to avoid creation error, where remote server is not expecting this
			createTime.WorkTypeRemote = ""

			// Run each as a goroutine
			go func() {
				job := newJobTrack(pw, fmt.Sprintf(
					"Creating time task: %s | %s | %s | %.2f hour(s) |",
					createTime.Date,
					createTime.Description,
					createTime.WorkType,
					createTimeHours,
				))
				if err := client.CreateTime(createTime); err == nil {
					setJobSuccess(job, "Created")
				} else {
					setJobError(job, err)
				}
			}()
		}
	}

	// Render all jobs, until all done
	time.Sleep(time.Millisecond * 100)
	for pw.IsRenderInProgress() {
	}

	if !ctx.Force {
		console.Header("All DONE! Press ENTER to exit.")
		fmt.Scanln()
	}

	return nil
}

func (c *timeCreateCmd) Run(ctx CLIContext) error {
	// Prep checks
	prepResult, client, err := createTimePrepSteps(ctx, c.InputFile)
	if err != nil {
		return err
	}

	// Real work
	err = createTimeExecSteps(ctx, prepResult, client)
	if err != nil {
		return err
	}

	return nil
}
