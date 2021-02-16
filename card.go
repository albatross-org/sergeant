package sergeant

import (
	"path/filepath"
	"time"
)

// Card is the basic unit of the program. It's an abstraction over an Albatross entry and represents a question-answer pair.
type Card struct {
	ID   string
	Path string
	Tags []string

	CompletionsPerfect []Completion
	CompletionsMinor   []Completion
	CompletionsMajor   []Completion

	QuestionPath string
	AnswerPath   string
}

// Completion is a mark specifying that a card was completed at a certain date in a certain amount of time.
type Completion struct {
	Date     time.Time
	Duration time.Duration
}

// PathParent returns the path to the parent of the card. You could think of this as the card category.
func (card *Card) PathParent() string {
	return filepath.Dir(card.Path)
}

// QuestionImage returns the data URI of the question image, created by converting the file to base64.
func (card *Card) QuestionImage() (string, error) {
	return encodeAsDataURI(card.QuestionPath)
}

// AnswerImage returns the data URI of the answer image, created by converting the file to base64.
func (card *Card) AnswerImage() (string, error) {
	return encodeAsDataURI(card.AnswerPath)
}
