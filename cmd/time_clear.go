package cmd

import (
	"fmt"
	"time"

	"github.com/hxhieu/b1-timetask-cli-go/common"
	"github.com/hxhieu/b1-timetask-cli-go/console"
	"github.com/hxhieu/b1-timetask-cli-go/intervals_api"
)

func clearTimePrepSteps(ctx CLIContext, weekOffset int) (*[]intervals_api.TimeEntry, *intervals_api.Client, error) {
	weekDays := common.GetWeekRange(time.Now(), weekOffset)

	// instantiate a Progress Writer and set up the options
	progress := common.NewProgressTracker(ctx.Debug)
	progress.Start()

	// Fetch time entries job
	job := progress.AddNewTrack("Fetch week time tasks")

	token, err := common.GetUserToken()
	if err != nil {
		job.SetError(err)
		progress.RenderUntilAllDone()
		return nil, nil, fmt.Errorf("failed to get user token: %w", err)
	}

	// API client
	client := intervals_api.New(token, ctx.Debug)

	// Fetch tasks
	tasks, err := client.GetTimeEntries(weekDays[0], weekDays[len(weekDays)-1])
	if err != nil {
		job.SetError(err)
		progress.RenderUntilAllDone()
		return nil, nil, fmt.Errorf("failed to fetch time entries: %w", err)
	}

	job.SetSuccess(fmt.Sprintf("Found %d task(s)", len(*tasks)))

	// Render all jobs, until all done
	progress.RenderUntilAllDone()

	console.PrintWeekRange(common.DateToString(weekDays[0]), common.DateToString(weekDays[6]), weekOffset)

	if !ctx.Force {
		console.Header("This is destructive and irreversable! Press ENTER to process, or CTRL+C to terminate.")
		fmt.Scanln()
	}

	return tasks, client, nil
}

func clearTimeExecSteps(ctx CLIContext, tasks *[]intervals_api.TimeEntry, client *intervals_api.Client) error {
	// instantiate a Progress Writer and set up the options
	progress := common.NewProgressTracker(ctx.Debug)
	progress.Start()

	// Calculate max text length for columns padding
	maxTitleLength, maxWorkTypeLength := intervals_api.CalcMaxFieldsLen(tasks)

	for _, t := range *tasks {
		job := progress.AddNewTrack(fmt.Sprintf(
			"Deleting %s",
			t.PaddedTitle(maxTitleLength, maxWorkTypeLength),
		))
		if err := client.DeleteTimeEntry(t.Id); err == nil {
			job.SetSuccess("Deleted")
		} else {
			job.SetError(err)
		}
	}

	// Render all jobs, until all done
	progress.RenderUntilAllDone()

	if !ctx.Force {
		console.Header("All DONE! Press ENTER to exit.")
		fmt.Scanln()
	}

	return nil
}

func (c *timeClearCmd) Run(ctx CLIContext) error {
	// Prep checks
	tasks, client, err := clearTimePrepSteps(ctx, c.WeekOffset)
	if err != nil {
		return err
	}

	// Real work
	err = clearTimeExecSteps(ctx, tasks, client)
	if err != nil {
		return err
	}

	return nil
}
