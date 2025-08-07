package common

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hxhieu/b1-timetask-cli-go/console"
)

type CalendarTaskParser struct {
	subjectTemplate *string
}

func NewCalendarTaskParser(subjectTemplate *string) *CalendarTaskParser {
	return &CalendarTaskParser{
		subjectTemplate: subjectTemplate,
	}
}

func (p *CalendarTaskParser) ParseEvents(events *[]OutLookCalendarEvent) ([]*TimeTaskInput, error) {
	taskMap := make(map[string]*TimeTaskInput)

	for _, event := range *events {
		code, desc, err := event.fromTitle(*p.subjectTemplate)
		if err != nil {
			console.ErrorLn(fmt.Sprintf("SKIPPED ERR parsing subject: '%s'. %s", event.Subject, err.Error()))
			continue
		}
		workType, billable, err := event.fromCategory()
		if err != nil {
			console.WarnLn(fmt.Sprintf("SKIPPED parsing category: '%s'", err.Error()))
		}

		if taskMap[code] == nil {
			taskMap[code] = &TimeTaskInput{
				Task:     code,
				Desc:     desc,
				WorkType: workType,
				Billable: billable,
			}
		}

		hours, weekDay, err := event.fromTime()
		if err != nil {
			console.WarnLn(fmt.Sprintf("SKIPPED parsing time: '%s'", err.Error()))
		}

		switch weekDay {
		case time.Sunday:
			taskMap[code].Sun += hours
		case time.Monday:
			taskMap[code].Mon += hours
		case time.Tuesday:
			taskMap[code].Tue += hours
		case time.Wednesday:
			taskMap[code].Wed += hours
		case time.Thursday:
			taskMap[code].Thu += hours
		case time.Friday:
			taskMap[code].Fri += hours
		case time.Saturday:
			taskMap[code].Sat += hours
		default:
			console.WarnLn(fmt.Sprintf("SKIPPED parsing week day: '%s'", "invalid week day"))
		}
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

func (e *OutLookCalendarEvent) fromCategory() (workType string, billable string, err error) {
	return "", "f", fmt.Errorf("not implemented")
}

func (e *OutLookCalendarEvent) fromTime() (hours float32, weekDay time.Weekday, err error) {
	return 0, time.Sunday, fmt.Errorf("not implemented")
}
