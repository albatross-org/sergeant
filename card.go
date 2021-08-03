package sergeant

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/albatross-org/go-albatross/entries"
	"gopkg.in/yaml.v3"
)

// Card is the basic unit of the program. It's an abstraction over an Albatross entry and represents a question-answer pair.
type Card struct {
	ID   string
	Path string

	Date  time.Time
	Tags  []string
	Notes string

	CompletionsPerfect []Completion
	CompletionsMinor   []Completion
	CompletionsMajor   []Completion

	QuestionPath string
	AnswerPath   string
}

// Completion is a mark specifying that a card was completed at a certain date in a certain amount of time,
// by a certain user.
type Completion struct {
	Date     time.Time     `yaml:"date"`
	Duration time.Duration `yaml:"time"`
	User     string        `yaml:"user"`
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

// TotalCompletions returns the total number of completions for this card.
func (card *Card) TotalCompletions() int {
	return len(card.CompletionsMajor) + len(card.CompletionsMinor) + len(card.CompletionsPerfect)
}

// UserPerfect returns the total number of perfect completions for a user.
func (card *Card) UserPerfect(user string) int {
	if user == "" {
		return len(card.CompletionsPerfect)
	}

	count := 0

	for _, completion := range card.CompletionsPerfect {
		if completion.User == user {
			count++
		}
	}

	return count
}

// UserMajor returns the total number of perfect completions for a user.
func (card *Card) UserMajor(user string) int {
	if user == "" {
		return len(card.CompletionsMajor)
	}

	count := 0

	for _, completion := range card.CompletionsMajor {
		if completion.User == user {
			count++
		}
	}

	return count
}

// UserMinor returns the total number of perfect completions for a user.
func (card *Card) UserMinor(user string) int {
	if user == "" {
		return len(card.CompletionsMinor)
	}

	count := 0

	for _, completion := range card.CompletionsMajor {
		if completion.User == user {
			count++
		}
	}

	return count
}

// TotalCompletions returns the total number of completions for this card, by a certain user.
func (card *Card) TotalCompletionsUser(user string) int {
	return card.UserPerfect(user) + card.UserMajor(user) + card.UserMinor(user)
}

// Content returns how the card is represented as an entry. Think of it like the opposite of cardFromEntry.
func (card *Card) Content() (string, error) {
	type frontmatter struct {
		Title       string   `yaml:"title"`
		Type        string   `yaml:"type"`
		Tags        []string `yaml:"tags"`
		Date        string   `yaml:"date"`
		Completions map[string][]map[string]string
	}

	entryFrontmatter := frontmatter{
		Title: "Question " + card.ID,
		Type:  "question",
		Tags:  card.Tags,
		Completions: map[string][]map[string]string{
			"perfect": completionToStringMap(card.CompletionsPerfect),
			"minor":   completionToStringMap(card.CompletionsMinor),
			"major":   completionToStringMap(card.CompletionsMajor),
		},
		Date: card.Date.Format("2006-01-02 15:04"),
	}

	frontmatterBytes, err := yaml.Marshal(entryFrontmatter)
	if err != nil {
		return "", fmt.Errorf("couldn't marshal new entry frontmatter: %w", err)
	}

	var out bytes.Buffer

	out.WriteString("---\n")
	out.Write(frontmatterBytes)
	out.WriteString("---\n")
	out.WriteString(card.Notes)

	return out.String(), nil
}

// cardFromEntry creates a Card from an *entries.Entry.
// Since a card is an abstraction over an entry, we have to painstakingly go through each individual piece of metadata
// rather than unmarshalling it to a struct. The card follows this pattern:
//   ---
//   title: "Question <random 16-character string>" // 16 character string becomes ID
//   type: "question"                              // Used to verify this is in fact a question.
//   tags: ["@?any-tags"]                          // This becomes the .Tags field.
//   completions:
//       perfect:                                  // This becomes the .CompletionsPerfect field.
//           - date: 2021-02-16 10:18
//             time: 7m10s
//			   user: olly
//       minor:                                    // This becomes the .CompletionsMinor field.
//           - date: 2021-02-16 10:18
//             time: 5m51s
//             user: jeff
//       major:                                    // This becomes the .CompletionsMajor field.
//           - date: 2021-02-16 10:18
//             time: 5m53s
//			   user: olly
//           - date: 2021-02-16 10:18
//             time: 5m53s
//             user: jeff
//   ---
//   Any additional notes about the card (This becomes the .Notes field).
func cardFromEntry(entry *entries.Entry) (*Card, error) {
	card := &Card{}

	// Check the entry is nil, we want to error instead of panicing with a nil pointer dereference.
	if entry == nil {
		return nil, fmt.Errorf("entry is nil")
	}

	// Verify that the entry's type field is correct.
	// Still not sure how I feel about this, is this is an unnecceasy step?
	entryType, ok := entry.Metadata["type"].(string)
	if !ok {
		return nil, fmt.Errorf("missing required 'type' field in card entry metadata")
	}

	if entryType != "question" {
		return nil, fmt.Errorf("expected metadata field 'type' to be 'question', not %q", entryType)
	}

	// Get the path of the question.
	card.Path = entry.Path

	// Copy accross any tags.
	card.Tags = entry.Tags

	// Copy across the date.
	card.Date = entry.Date

	// Copy across the content/notes.
	card.Notes = entry.Contents

	// Get the ID of the question.
	if !strings.HasPrefix(entry.Title, "Question ") {
		return nil, fmt.Errorf("expected card title to start with 'Question', it is %q", entry.Title)
	}

	card.ID = strings.TrimPrefix(entry.Title, "Question ")

	// Get the completions.
	// This is a map of completion types ("perfect", "minor", "major") to lists of completions ("date", "time").
	completionsMapInterface, ok := entry.Metadata["completions"].(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("couldn't parse 'completions' completions in card entry metadata")
	}

	completionsMap, err := completionsMapInterfaceToTypedMap(completionsMapInterface)
	if err != nil {
		return nil, err
	}

	completionsPerfect, completionsMinor, completionsMajor, err := completionsMapToStruct(completionsMap)
	if err != nil {
		return nil, err
	}

	card.CompletionsPerfect = completionsPerfect
	card.CompletionsMinor = completionsMinor
	card.CompletionsMajor = completionsMajor

	// Verify that a question and answer are attached and set them.
	var questionPath, answerPath string
	for _, attachment := range entry.Attachments {
		if strings.HasPrefix(attachment.Name, "question.") {
			questionPath = attachment.AbsPath
		} else if strings.HasPrefix(attachment.Name, "answer.") {
			answerPath = attachment.AbsPath
		}
	}

	if questionPath == "" {
		return nil, fmt.Errorf("card entry has no 'question' attachment")
	}

	if answerPath == "" {
		return nil, fmt.Errorf("card entry has no 'answer' attachment")
	}

	card.QuestionPath = questionPath
	card.AnswerPath = answerPath

	return card, nil
}

// completionsMapInterfaceToTypedMap converts a map[interface{}]interface{} to a map[string][]map[string]string, the format ready to be used
// by the rest of the program.
// I feel like there's a much better way of doing this and the variable names make me want to be sick. Is it really neccessary to cast this many
// times or can you combine them somehow?
func completionsMapInterfaceToTypedMap(mapInterface map[interface{}]interface{}) (map[string][]map[string]string, error) {
	stringToSlice := map[string][]interface{}{}

	// In this first step, we turn the overall map[interface{}]interface{} into a map[string][]interface{}{}
	for key, value := range mapInterface {
		keyString, ok := key.(string)
		if !ok {
			return nil, fmt.Errorf("couldn't convert completions metadata to typed map, %q not a string, got %T instead", key, key)
		}

		valueSlice, ok := value.([]interface{})
		if !ok {
			return nil, fmt.Errorf("coulnd't convert completions metadata to typed map, %q not a slice of interface{}, got %T instead", value, value)
		}

		stringToSlice[keyString] = valueSlice
	}

	stringToInterfaceMap := map[string][]map[interface{}]interface{}{}

	// In this step, we turn the map[string][]interface{} into a map[string][]map[interface{}]interface{}.
	// We need two loops here because Go can't do a type assertion on a []interface{}, only an interface{}.
	for key, value := range stringToSlice {
		interfaceMaps := []map[interface{}]interface{}{}
		for _, subValue := range value {
			subValueInterfaceMap, ok := subValue.(map[interface{}]interface{})
			if !ok {
				return nil, fmt.Errorf("couln't convert completions metadata to typed map, %q not a map of interfaces, got %T instead", value, value)
			}

			interfaceMaps = append(interfaceMaps, subValueInterfaceMap)
		}

		stringToInterfaceMap[key] = interfaceMaps
	}

	stringToStringMap := map[string][]map[string]string{}

	// In this final step, we turn the map[string][]map[interface{}] into what we want, a map[string][]map[string]string.
	for key, value := range stringToInterfaceMap {
		stringMaps := []map[string]string{}

		for _, subInterfaceMap := range value {
			stringMap := map[string]string{}

			for subKey, subValue := range subInterfaceMap {
				keyString, ok := subKey.(string)
				if !ok {
					return nil, fmt.Errorf("couldn't convert completions metadata to typed map, subkey %q not a string, got %T instead", subKey, subKey)
				}

				valueString, ok := subValue.(string)
				if !ok {
					return nil, fmt.Errorf("couldn't convert completions metadata to typed map, subvalue %q not a string, got %T instead", subValue, subValue)
				}

				stringMap[keyString] = valueString
			}

			stringMaps = append(stringMaps, stringMap)
		}

		stringToStringMap[key] = stringMaps
	}

	return stringToStringMap, nil
}

// completionsMapToStruct converts a map consisting of completion types ("perfect", "minor", "major") mapped lists of completions ("date", "time") to
// three lists of completion types. The order for return is perfect, minor and finally major completions.
func completionsMapToStruct(completionsMap map[string][]map[string]string) (completionsPerfect []Completion, completionsMinor []Completion, completionsMajor []Completion, err error) {
	for completionType, completionList := range completionsMap {
		if completionType != "perfect" && completionType != "minor" && completionType != "major" {
			return nil, nil, nil, fmt.Errorf("not expecting completions field %q in card metadata", completionType)
		}

		for _, completionMap := range completionList {
			if completionMap["date"] == "" {
				return nil, nil, nil, fmt.Errorf("'date' field in %q completion list is empty", completionType)
			}

			if completionMap["time"] == "" {
				return nil, nil, nil, fmt.Errorf("'time' field in %q completion list is empty", completionType)
			}

			date, err := time.Parse("2006-01-02 15:04", completionMap["date"])
			if err != nil {
				return nil, nil, nil, fmt.Errorf("'date' field %q in %q completion list not a valid '2006-01-02 15:04' date: %w", completionMap["date"], completionType, err)
			}

			duration, err := time.ParseDuration(completionMap["time"])
			if err != nil {
				return nil, nil, nil, fmt.Errorf("'time' field %q in %q completion list not a valid duration: %w", completionMap["time"], completionType, err)
			}

			user := completionMap["user"]
			if user == "" {
				user = "olly" // TODO: change this
			}

			completion := Completion{
				Date:     date,
				Duration: duration,
				User:     user,
			}

			switch completionType {
			case "perfect":
				completionsPerfect = append(completionsPerfect, completion)
			case "minor":
				completionsMinor = append(completionsMinor, completion)
			case "major":
				completionsMajor = append(completionsMajor, completion)
			}

		}
	}

	return completionsPerfect, completionsMinor, completionsMajor, nil
}

// completionToStringMap converts a []Completion to a []map[string]string. This is needed because by default YAML will unmarshal
// time.Time fields using a different format to the one the program expects. By manually converting it to a map[string]string first,
// we can use our own custom date format.
func completionToStringMap(completions []Completion) []map[string]string {
	out := []map[string]string{}

	for _, completion := range completions {
		stringMap := map[string]string{}
		stringMap["date"] = completion.Date.Format("2006-01-02 15:04")
		stringMap["time"] = completion.Duration.String()
		stringMap["user"] = completion.User

		out = append(out, stringMap)
	}

	return out
}
