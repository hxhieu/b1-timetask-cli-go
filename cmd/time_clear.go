package cmd

import (
	"fmt"
	"time"

	"github.com/hxhieu/b1-timetask-cli-go/common"
	"github.com/hxhieu/b1-timetask-cli-go/console"
	"github.com/hxhieu/b1-timetask-cli-go/intervals_api"
	"github.com/jedib0t/go-pretty/v6/progress"
)

func clearTimePrepSteps(ctx CLIContext) (*[]intervals_api.TimeEntry, *intervals_api.Client, error) {
	var client *intervals_api.Client
	var result *[]intervals_api.TimeEntry

	// instantiate a Progress Writer and set up the options
	pw := progress.NewWriter()
	setDefaultProgress(&pw)

	go pw.Render()

	// Fetch time entries job
	job := newJobTrack(pw, "Fetch current week time tasks")

	if token, err := common.GetUserToken(); err == nil {

		// API client
		client = intervals_api.New(token, ctx.Debug)

		// Fetch tasks
		weekDays := common.GetWeekRange(time.Now())
		if tasks, err := client.GetTimeEntries(weekDays[0], weekDays[len(weekDays)-1]); err == nil {
			result = tasks
			setJobSuccess(job, fmt.Sprintf("Found %d task(s)", len(*tasks)))
		} else {
			setJobError(job, err)
		}
	} else {
		setJobError(job, err)
	}

	// Render all jobs, until all done
	time.Sleep(time.Millisecond * 100)
	for pw.IsRenderInProgress() {
	}

	if !ctx.Force {
		console.Header("This is destructive and irreversable! Press ENTER to process, or CTRL+C to terminate.")
		fmt.Scanln()
	}

	return result, client, nil
}

func clearTimeExecSteps(ctx CLIContext, tasks *[]intervals_api.TimeEntry, client *intervals_api.Client) error {
	// instantiate a Progress Writer and set up the options
	pw := progress.NewWriter()
	setDefaultProgress(&pw)

	go pw.Render()

	for _, t := range *tasks {
		// Run each as a goroutine
		go func() {
			job := newJobTrack(pw, fmt.Sprintf(
				"Deleting time task: %s | %s | %s | %s hour(s) |",
				t.Date,
				t.Description,
				t.WorkTypeRemote,
				t.Time,
			))
			if err := client.DeleteTimeEntry(t.Id); err == nil {
				setJobSuccess(job, "Deleted")
			} else {
				setJobError(job, err)
			}
		}()
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

func (c *timeClearCmd) Run(ctx CLIContext) error {
	// Prep checks
	tasks, client, err := clearTimePrepSteps(ctx)
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
