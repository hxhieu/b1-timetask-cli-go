package cmd

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hxhieu/b1-timetask-cli-go/common"
	"github.com/hxhieu/b1-timetask-cli-go/console"
	"github.com/hxhieu/b1-timetask-cli-go/intervals_api"
	"github.com/hxhieu/b1-timetask-cli-go/ms_graph"
	"github.com/jedib0t/go-pretty/v6/progress"
)

type createTimePrepResult struct {
	userId string
	tasks  []*common.TimeTaskInput
}

func createTimePrepSteps(ctx CLIContext, inputFile *string, weekOffset int, subjectTemplate *string, inputType *string, useDeviceCode bool) (*createTimePrepResult, *intervals_api.Client, error) {
	if !(*inputType == "csv" || *inputType == "calendar") {
		return nil, nil, fmt.Errorf("invalid input type: %s, expected 'csv' or 'calendar'", *inputType)
	}

	// Shared vars between steps
	result := &createTimePrepResult{}
	var timeIntervalClient *intervals_api.Client
	var tasks []*common.TimeTaskInput

	// instantiate a Progress Writer and set up the options
	pw := progress.NewWriter()
	setDefaultProgress(&pw)

	go pw.Render()

	// Check and fetch user job
	job := newJobTrack(pw, "Check user", ctx)

	if token, err := common.GetUserToken(); err == nil {
		// API client
		timeIntervalClient = intervals_api.New(token, ctx.Debug)

		// Fetch the user
		if me, err := timeIntervalClient.Me(); err == nil {
			result.userId = me.Id
			setJobSuccess(job, fmt.Sprintf("Found user: %s %s <%s>", me.FirstName, me.LastName, me.Email))
		} else {
			setJobError(job, err)
		}
	} else {
		setJobError(job, err)
	}

	if !job.IsErrored() {
		job = newJobTrack(pw, "Prepare task inputs", ctx)

		switch *inputType {
		case "csv":
			// Parse the CSV file
			csvParser := common.NewCsvTaskParser(inputFile)
			var err error = nil
			tasks, err = csvParser.GetTasks()
			if err != nil {
				setJobError(job, err)
				return nil, nil, err
			}
		case "calendar":
			msGraphClient := ms_graph.NewMsGraphClient(ctx.Debug, useDeviceCode)
			// Get raw events from calendar
			events, err := msGraphClient.GetMyCalendarEvents(weekOffset)
			if err != nil {
				setJobError(job, err)
				return nil, nil, err
			}
			// Parse the events to tasks
			eventsParser := common.NewCalendarTaskParser(subjectTemplate)
			tasks, err = eventsParser.ParseEvents(events)
			if err != nil {
				setJobError(job, err)
				return nil, nil, err
			}
		default:
			return nil, nil, fmt.Errorf("unknown input type: %s", *inputType)
		}

		// Concat IDs, to pass to the remoter server
		var taskValues string
		var projectValues string
		for _, t := range tasks {
			if t != nil {
				taskValues += t.Task + ","
			}
		}
		taskValues = strings.TrimSuffix(taskValues, ",")

		// Fetch needed details from remote
		if remoteTasks, err := timeIntervalClient.FetchTasks(taskValues); err == nil {
			// TODO: Optimise this? nested loops here
			for _, remoteTask := range *remoteTasks {
				// Also build the project IDs list
				projectValues += remoteTask.ProjectId + ","
				for _, localTask := range tasks {
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
			projectValues = strings.TrimSuffix(projectValues, ",")

			// Fetch work types, because they are setup per project
			if remoteWorkTypes, err := timeIntervalClient.FetchProjectWorkTypes(projectValues); err == nil {
				// TODO: Optimise this? nested loops here
				for _, localTask := range tasks {
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

	// Print the tasks table
	common.PrintTimeTasks(tasks)

	return result, timeIntervalClient, nil
}

func createTimeExecSteps(ctx CLIContext, prepResult *createTimePrepResult, client *intervals_api.Client, weekOffset int) error {
	// instantiate a Progress Writer and set up the options
	pw := progress.NewWriter()
	setDefaultProgress(&pw)

	go pw.Render()

	weekDays := common.GetWeekRange(time.Now(), weekOffset)

	console.PrintWeekRange(common.DateToString(weekDays[0]), common.DateToString(weekDays[6]), weekOffset)

	if !ctx.Force {
		console.Header("Press ENTER to process, or CTRL+C to terminate.")
		fmt.Scanln()
	}

	// Calculate max text length for columns padding
	maxTitleLength, maxWorkTypeLength := common.CalcMaxFieldsLen(prepResult.tasks)

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

			job := newJobTrack(pw, fmt.Sprintf(
				"Creating %s",
				createTime.PaddedTitle(maxTitleLength, maxWorkTypeLength),
			), ctx)
			if err := client.CreateTime(createTime); err == nil {
				setJobSuccess(job, "Created")
			} else {
				setJobError(job, err)
			}
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
	prepResult, client, err := createTimePrepSteps(
		ctx,
		c.InputFile,
		c.WeekOffset,
		c.CalendarSubjectTemplate,
		c.InputType,
		c.UseDeviceCode,
	)
	if err != nil {
		return err
	}

	// Real work
	err = createTimeExecSteps(ctx, prepResult, client, c.WeekOffset)
	if err != nil {
		return err
	}

	return nil
}
