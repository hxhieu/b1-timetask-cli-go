package common

import (
	"context"
	"fmt"
	"regexp"
	"strings"
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

func GetCalendarEvents(weekOffset int, subjectTemplate *string) (*[]TimeTaskInput, error) {
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

	// timeTaskMap := make(map[string]*TimeTaskInput)

	for _, e := range events {
		if e.GetRecurrence() != nil {
			continue // TODO: skip recurrence events for now
		}
		// Get subject
		if subject := e.GetSubject(); subject != nil {
			code, desc, err := parseTaskFromSubject(*subject, *subjectTemplate)
			if err != nil {
				fmt.Printf("SKIPPED '%s': %v\n", *subject, err)
				continue // Skip this event if parsing fails
			}

			fmt.Printf("Parsed task code: %s, description: %s\n", code, desc)
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

// ParseTaskFromSubject parses a subject string using the provided template to extract task_code and task_desc
// Template example: "[TT:{task_code}]{task_desc}"
// Input example: "[TT:11238]Improve documentation"
// Returns: task_code="11238", task_desc="Improve documentation", err=nil
func parseTaskFromSubject(subject, template string) (taskCode, taskDesc string, err error) {
	// Escape special regex characters in the template, but keep our placeholders
	escapedTemplate := regexp.QuoteMeta(template)

	// Replace our placeholders with regex capture groups
	// {task_code} becomes (.+?) for non-greedy capture
	// {task_desc} becomes (.+?) for non-greedy capture
	regexPattern := escapedTemplate
	regexPattern = strings.ReplaceAll(regexPattern, `\{task_code\}`, `(.+?)`)
	regexPattern = strings.ReplaceAll(regexPattern, `\{task_desc\}`, `(.+)`) // greedy for task_desc as it's usually at the end

	// Add start and end anchors to match the entire string
	regexPattern = "^" + regexPattern + "$"

	// Compile the regex
	re, err := regexp.Compile(regexPattern)
	if err != nil {
		return "", "", err
	}

	// Find matches
	matches := re.FindStringSubmatch(subject)
	if len(matches) < 3 { // matches[0] is the full match, matches[1] and matches[2] are the groups
		return "", "", fmt.Errorf("subject does not match template pattern")
	}

	// Determine which capture group corresponds to which placeholder
	// We need to find the order of placeholders in the original template
	taskCodeIndex := strings.Index(template, "{task_code}")
	taskDescIndex := strings.Index(template, "{task_desc}")

	if taskCodeIndex < taskDescIndex {
		// task_code comes first
		return strings.TrimSpace(matches[1]), strings.TrimSpace(matches[2]), nil
	} else {
		// task_desc comes first
		return strings.TrimSpace(matches[2]), strings.TrimSpace(matches[1]), nil
	}
}
