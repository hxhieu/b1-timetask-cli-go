package ms_graph

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/hxhieu/b1-timetask-cli-go/common"
	"github.com/hxhieu/b1-timetask-cli-go/debug"
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
	debugFile := ".debug_calendar-events.json"
	if c.debug {
		if debugData := debug.LoadDataFile[[]common.OutLookCalendarEvent](debugFile); debugData != nil {
			return debugData, nil
		}
	}

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

	if c.debug {
		debug.WriteDataFile(debugFile, calendarEvents)
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
			if parsedTime, err := common.ParseGraphDateTime(*startTime); err == nil {
				calendarEvent.Start.DateTime = parsedTime
			} else {
				fmt.Printf("Error parsing start time '%s': %v\n", *startTime, err)
			}
		}
	}

	// Get end time
	if end := event.GetEnd(); end != nil {
		if endTime := end.GetDateTime(); endTime != nil {
			if parsedTime, err := common.ParseGraphDateTime(*endTime); err == nil {
				calendarEvent.End.DateTime = parsedTime
			} else {
				fmt.Printf("Error parsing end time '%s': %v\n", *endTime, err)
			}
		}
	}
	return calendarEvent
}
