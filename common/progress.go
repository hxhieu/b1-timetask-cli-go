package common

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/progress"
)

// ProgressTracker is a utility to track progress of jobs with a visual representation.
type ProgressTracker struct {
	pw    progress.Writer
	debug bool
}

// ProgressTrack represents a single job track in the progress tracker.
type ProgressTrack struct {
	track *progress.Tracker
}

func setDefaultStyles(pw *progress.Writer) {
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

// NewProgressTracker creates a new ProgressTracker instance.
func NewProgressTracker(debug bool) *ProgressTracker {
	pw := progress.NewWriter()
	setDefaultStyles(&pw)
	return &ProgressTracker{
		pw:    pw,
		debug: debug,
	}
}

// Start initializes the progress tracker and starts rendering.
func (p *ProgressTracker) Start() {
	go p.pw.Render()
}

// AddNewTrack creates a new track for a job with the given title.
func (p *ProgressTracker) AddNewTrack(title string) *ProgressTrack {
	if p.debug {
		title = fmt.Sprintf("%s %s", color.YellowString("[DEBUG ONLY]"), title)
	}

	job := progress.Tracker{
		Message: title,
		Units: progress.Units{
			Notation:         " jobs",
			NotationPosition: progress.UnitsNotationPositionAfter,
		},
		DeferStart: false,
	}

	p.pw.AppendTracker(&job)

	return &ProgressTrack{
		track: &job,
	}
}

// SetSuccess updates the track message to indicate success and marks it as done.
func (t *ProgressTrack) SetSuccess(message string) {
	t.track.UpdateMessage(fmt.Sprintf("%s %s", t.track.Message, color.HiGreenString(message)))
	t.track.MarkAsDone()
}

// SetError updates the track message to indicate an error and marks it as errored.
func (t *ProgressTrack) SetError(err error) {
	t.track.UpdateMessage(fmt.Sprintf("%s %s", t.track.Message, color.RedString(err.Error())))
	t.track.MarkAsErrored()
}

func (t *ProgressTrack) IsErrored() bool {
	return t.track.IsErrored()
}

// RenderUntilAllDone waits for all jobs to be done before returning.
func (p *ProgressTracker) RenderUntilAllDone() {
	// Render all jobs, until all done
	time.Sleep(time.Millisecond * 100)
	for p.pw.IsRenderInProgress() {
	}
}
