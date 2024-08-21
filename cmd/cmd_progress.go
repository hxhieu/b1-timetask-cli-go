package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/progress"
)

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
	p.Style().Visibility.Value = false
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
