package common

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hxhieu/b1-timetask-cli-go/console"
)

type CalendarTaskParser struct {
	subjectTemplate *string
	defaultWorkType *string
}

func NewCalendarTaskParser(subjectTemplate *string, defaultWorkType *string) *CalendarTaskParser {
	return &CalendarTaskParser{
		subjectTemplate: subjectTemplate,
		defaultWorkType: defaultWorkType,
	}
}

func (p *CalendarTaskParser) ParseEvents(events *[]OutLookCalendarEvent) ([]*TimeTaskInput, error) {
	taskMap := make(map[string]*TimeTaskInput)

	for _, event := range *events {
		code, desc, err := event.fromTitle(*p.subjectTemplate)
		if err != nil {
			printError(fmt.Sprintf("parsing subject: '%s'", event.Subject), err)
			continue
		}
		workType, billable, err := event.fromCategory(p.defaultWorkType)
		if err != nil {
			printError("parsing category", err)
			continue
		}

		// Group by code + worktype + desc
		groupingKey := code + "_" + strings.ToLower(workType) + "_" + strings.ToLower(desc)
		if taskMap[groupingKey] == nil {
			taskMap[groupingKey] = &TimeTaskInput{
				Task:     code,
				Desc:     desc,
				WorkType: workType,
				Billable: billable,
			}
		}

		hours, weekDay, err := event.fromTime()
		if err != nil {
			printError("parsing time", err)
			continue
		}

		switch weekDay {
		case time.Sunday:
			taskMap[groupingKey].Sun += hours
		case time.Monday:
			taskMap[groupingKey].Mon += hours
		case time.Tuesday:
			taskMap[groupingKey].Tue += hours
		case time.Wednesday:
			taskMap[groupingKey].Wed += hours
		case time.Thursday:
			taskMap[groupingKey].Thu += hours
		case time.Friday:
			taskMap[groupingKey].Fri += hours
		case time.Saturday:
			taskMap[groupingKey].Sat += hours
		default:
			printError("parsing time", errors.New("invalid week day"))
			continue
		}
		console.Success("[DONE] ")
		console.InfoLn(event.Subject)
	}

	result := []*TimeTaskInput{}

	for _, task := range taskMap {
		result = append(result, task)
	}

	return result, nil
}

// Template example: "[TT:{task_code}]{task_desc}"
// Input example: "[TT:11238]Improve documentation"
// Returns: task_code="11238", task_desc="Improve documentation", err=nil
func (e *OutLookCalendarEvent) fromTitle(template string) (taskCode, taskDesc string, err error) {
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
	matches := re.FindStringSubmatch(e.Subject)
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

func (e *OutLookCalendarEvent) fromCategory(defaultWorkType *string) (workType string, billable string, err error) {
	billable = "t"              // default to true
	workType = *defaultWorkType // default work type

	// Event has assigned categories
	// then we need to parse it
	if len(e.Categories) > 0 {
		workTypes := make([]string, 0)
		// For each of the event categories, find it is Nonbillable or candidate work types
		for _, category := range e.Categories {
			maybeWorkType := strings.TrimSpace(category)
			if strings.ToLower(maybeWorkType) != "nonbillable" {
				workTypes = append(workTypes, maybeWorkType)
			} else {
				billable = "f"
			}
		}

		// Sometime we dont want to enter a work type, default to -None-
		if len(workTypes) == 0 {
			workType = "-None-"
		} else {
			// First assigned category will be the work type
			workType = workTypes[0]
		}
	}

	return workType, billable, nil
}

func (e *OutLookCalendarEvent) fromTime() (hours float32, weekDay time.Weekday, err error) {
	diff := e.End.DateTime.Sub(e.Start.DateTime)
	hours = float32(diff.Hours())
	startLocal := e.Start.DateTime.In(time.Local)
	return hours, startLocal.Weekday(), nil
}

func printError(msg string, err error) {
	if err == nil {
		return
	}
	console.ErrorLn(fmt.Sprintf("[ERRO] %s, skipped event. %v", msg, err))
}
