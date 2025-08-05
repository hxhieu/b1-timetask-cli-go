package common

import (
	"context"
	"fmt"
	"time"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/users"
)

type outLookCalendarEventTime struct {
	DateTime time.Time `json:"dateTime"`
}

type outLookCalendarEvent struct {
	Categories []string                 `json:"categories"`
	Subject    string                   `json:"subject"`
	Start      outLookCalendarEventTime `json:"start"`
	End        outLookCalendarEventTime `json:"end"`
}

func GetCalendarEvents(weekOffset int) (*[]TimeTaskInput, error) {
	weekDays := GetWeekRange(time.Now(), weekOffset)
	monday := weekDays[0]

	// Convert to UTC for Microsoft Graph API
	mondayUtc := monday.UTC()
	nextMondayUtc := mondayUtc.AddDate(0, 0, 7)

	// One week range
	filter := fmt.Sprintf(
		"start/dateTime ge '%s' and start/dateTime lt '%s'",
		mondayUtc.Format(time.RFC3339),
		nextMondayUtc.Format(time.RFC3339),
	)

	fmt.Println("Fetching calendar events with filter:", filter)

	// Extra authenticity checks
	filter += " and isDraft eq false"
	filter += " and isCancelled eq false"
	filter += " and isAllDay eq false"

	var top int32 = 100

	cred, err := getCreds()
	if err != nil {
		return nil, err
	}

	client, err := msgraphsdk.NewGraphServiceClientWithCredentials(cred, []string{})
	if err != nil {
		return nil, err
	}

	result, err := client.Me().Calendar().Events().Get(context.Background(), &users.ItemCalendarEventsRequestBuilderGetRequestConfiguration{
		QueryParameters: &users.ItemCalendarEventsRequestBuilderGetQueryParameters{
			Select: []string{"categories", "subject", "start", "end"},
			Filter: &filter,
			Top:    &top,
		},
	})
	if err != nil {
		return nil, err
	}

	// Get the events from the result
	events := result.GetValue()

	// Process each event
	resultEvents := make([]TimeTaskInput, 0)
	for _, event := range events {
		if event.GetRecurrence() != nil {
			continue // TODO: skip recurrence events for now
		}

	}

	return &resultEvents, nil
}

func (e *outLookCalendarEvent) fromRemote(event models.Eventable) {
	// Get subject
	if subject := event.GetSubject(); subject != nil {
		e.Subject = *subject
	}

	// Get categories
	if categories := event.GetCategories(); len(categories) > 0 {
		e.Categories = categories
	}

	// Get start time
	if start := event.GetStart(); start != nil {
		if startTime := start.GetDateTime(); startTime != nil {
			// Parse the time string to time.Time if needed
			if parsedTime, err := time.Parse(time.RFC3339, *startTime); err == nil {
				e.Start.DateTime = parsedTime
			}
		}
	}

	// Get end time
	if end := event.GetEnd(); end != nil {
		if endTime := end.GetDateTime(); endTime != nil {
			// Parse the time string to time.Time if needed
			if parsedTime, err := time.Parse(time.RFC3339, *endTime); err == nil {
				e.End.DateTime = parsedTime
			}
		}
	}
}
