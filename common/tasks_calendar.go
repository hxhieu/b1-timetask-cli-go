package common

import (
	"fmt"
	"regexp"
	"strings"
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
	result := []*TimeTaskInput{}

	return result, nil
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
