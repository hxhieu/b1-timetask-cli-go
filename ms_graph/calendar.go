package ms_graph

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/hxhieu/b1-timetask-cli-go/common"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/users"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

type MsGraphClient struct {
	debug         bool
	useDeviceCode bool
}

func NewMsGraphClient(debug bool, useDeviceCode bool) *MsGraphClient {
	return &MsGraphClient{
		debug:         debug,
		useDeviceCode: useDeviceCode,
	}
}

func (c *MsGraphClient) GetMyCalendarEvents(weekOffset int) (*[]common.OutLookCalendarEvent, error) {
	weekDays := common.GetWeekRange(time.Now(), weekOffset)
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

	// Extra authenticity checks
	filter += " and isDraft eq false"
	filter += " and isCancelled eq false"
	filter += " and isAllDay eq false"

	var top int32 = 100

	cred, err := common.GetMsGraphCreds(c.useDeviceCode)
	if err != nil {
		return nil, err
	}

	var client *msgraphsdk.GraphServiceClient

	if browserCreds, ok := cred.(*azidentity.InteractiveBrowserCredential); ok {
		// Use interactive browser credential
		client, err = msgraphsdk.NewGraphServiceClientWithCredentials(browserCreds, []string{})
		if err != nil {
			return nil, err
		}
	} else if deviceCodeCreds, ok := cred.(*azidentity.DeviceCodeCredential); ok {

		client, err = msgraphsdk.NewGraphServiceClientWithCredentials(deviceCodeCreds, []string{})
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("unsupported credential type: %T", cred)
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

	var calendarEvents []common.OutLookCalendarEvent

	for _, e := range events {
		if e.GetRecurrence() != nil {
			continue // TODO: skip recurrence events for now
		}
		calendarEvent := fromRemote(e)
		calendarEvents = append(calendarEvents, calendarEvent)
	}

	return &calendarEvents, nil
}

func fromRemote(event models.Eventable) common.OutLookCalendarEvent {
	calendarEvent := common.OutLookCalendarEvent{}
	// Get subject
	if subject := event.GetSubject(); subject != nil {
		calendarEvent.Subject = *subject
	}

	// Get categories
	if categories := event.GetCategories(); len(categories) > 0 {
		calendarEvent.Categories = categories
	}

	// Get start time
	if start := event.GetStart(); start != nil {
		if startTime := start.GetDateTime(); startTime != nil {
			// Parse the time string to time.Time if needed
			if parsedTime, err := time.Parse(time.RFC3339, *startTime); err == nil {
				calendarEvent.Start.DateTime = parsedTime
			}
		}
	}

	// Get end time
	if end := event.GetEnd(); end != nil {
		if endTime := end.GetDateTime(); endTime != nil {
			// Parse the time string to time.Time if needed
			if parsedTime, err := time.Parse(time.RFC3339, *endTime); err == nil {
				calendarEvent.End.DateTime = parsedTime
			}
		}
	}
	return calendarEvent
}
