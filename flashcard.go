package pelican

import (
	"bytes"
	"fmt"
	"time"

	"github.com/albatross-org/go-albatross/entries"
)

// Flashcard represents a question and answer pair.
type Flashcard struct {
	ID   string // UUID that uniquely identifies a Flashcard.
	Path string // The path to the flashcard, used to make updates later.

	Perfect []time.Time // Dates of times it has been completed "perfectly".
	Minor   []time.Time // Dates of times it has been completed and marked as a minor mistake.
	Major   []time.Time // Dates of times it has been completed and marked as a major mistake.

	QuestionImg string // Path to question image.
	AnswerImg   string // Path to answer image.

	Parent *Set // Pointer to the flashcard set it's part of.

	Notes string // Any additional information about the question.

	Entry *entries.Entry `json:"-"` // Pointer to the underlying entry.
}

// Name returns a friendlier name for the question.
func (f *Flashcard) Name() string {
	names := []string{f.ID}

	// Go up the list of parents, appending to the list of names.
	// Afterwards, the list will be in reverse order ("Question jD...", "Ex10A", "Chapter 10", "...")
	parent := f.Parent
	for parent != nil {
		names = append(names, parent.Name)
		parent = parent.Parent
	}

	// Go backwards through the list of names and add it to the overall name.
	var out bytes.Buffer
	for i := len(names) - 1; i >= 0; i-- {
		if i != len(names)-1 {
			out.WriteString(" -> ")
		}

		out.WriteString(names[i])
	}

	return out.String()
}

// isFlashcard will return true if the entry given is a flashcard.
// This is used in the setFromTree function to determine if we should parse it as a set or a flashcard.
// It's a very rough check -- we just test for the precense of a "minor" key in the metadata.
func isFlashcard(entry *entries.Entry) bool {
	if entry == nil {
		return false
	}

	_, ok := entry.Metadata["minor"]
	return ok
}

// flashcardFromEntry returns a flashcard from an *albatross.Entry. Some values still need to be set, like the parent.
// It returns an error if the entry is malformed.
func flashcardFromEntry(entry *entries.Entry) (*Flashcard, error) {
	if entry.Title == "" {
		return nil, fmt.Errorf("entry has no name")
	}

	perfect, err := interfaceToTimeSlice(entry.Metadata["perfect"], "2006-01-02 15:04") // TODO: don't hard code
	if err != nil {
		return nil, fmt.Errorf("'perfect' key was not a list of dates: %w", err)
	}

	minor, err := interfaceToTimeSlice(entry.Metadata["minor"], "2006-01-02 15:04")
	if err != nil {
		return nil, fmt.Errorf("'minor' key was not a list of dates")
	}

	major, err := interfaceToTimeSlice(entry.Metadata["major"], "2006-01-02 15:04")
	if err != nil {
		return nil, fmt.Errorf("'major' key was not a list of dates")
	}

	questionImg := ""
	answerImg := ""

	for _, attachment := range entry.Attachments {
		// TODO: handle different file formats more effictively (what if both are present?)
		if attachment.Name == "question.png" || attachment.Name == "question.jpg" {
			questionImg = attachment.AbsPath
		}

		if attachment.Name == "answer.png" || attachment.Name == "answer.jpg" {
			answerImg = attachment.AbsPath
		}
	}

	if questionImg == "" {
		return nil, fmt.Errorf("entry has no question image attachment")
	}

	if answerImg == "" {
		return nil, fmt.Errorf("entry has no answer image attachment")
	}

	return &Flashcard{
		ID:   entry.Title,
		Path: entry.Path,

		Perfect: perfect,
		Minor:   minor,
		Major:   major,

		QuestionImg: questionImg,
		AnswerImg:   answerImg,

		Notes: entry.Contents,

		Parent: nil, // This still needs setting by the caller.
		Entry:  entry,
	}, nil
}

// interfaceToTimeSlice converts an interface{} type to a slice type.
func interfaceToTimeSlice(data interface{}, dateFormat string) ([]time.Time, error) {
	slice, ok := data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("not a slice")
	}

	times := []time.Time{}

	for _, elem := range slice {
		str, ok := elem.(string)
		if !ok {
			return nil, fmt.Errorf("contained a non-string value")
		}

		t, err := time.Parse(dateFormat, str)
		if err != nil {
			return nil, fmt.Errorf("couldn't parse date: %w", err)
		}

		times = append(times, t)
	}

	return times, nil
}
