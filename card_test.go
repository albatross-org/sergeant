package sergeant

import (
	"fmt"
	"testing"
	"time"

	"github.com/albatross-org/go-albatross/entries"
	"github.com/stretchr/testify/assert"
)

// TestCardFromEntryValid tests that the cardFromEntry function is working for valid entries.
func TestCardFromEntryValid(t *testing.T) {
	testCases := []struct {
		entry *entries.Entry
		card  *Card
	}{
		{
			entry: &entries.Entry{
				Path:  "further-maths/core-pure-1/chapter-1-complex-numbers/ex1a/question-BtIrmFTJo49QuJC4",
				Title: "Question BtIrmFTJo49QuJC4",
				Date:  time.Date(2021, 02, 16, 10, 18, 0, 0, time.Now().Location()),
				Attachments: []entries.Attachment{
					{AbsPath: "question.png", Name: "question.png"},
					{AbsPath: "answer.png", Name: "answer.png"},
				},
				Metadata: map[string]interface{}{
					"type": "question",
					"completions": map[string][]map[string]string{
						"perfect": {
							{
								"date": "2021-02-16 10:18",
								"time": "7m10s",
								"user": "",
							},
						},
						"minor": {
							{
								"date": "2021-02-16 10:18",
								"time": "5m51s",
								"user": "",
							},
						},
						"major": {
							{
								"date": "2021-02-16 10:18",
								"time": "5m53s",
								"user": "",
							},
							{
								"date": "2021-02-16 10:18",
								"time": "5m53s",
								"user": "",
							},
						},
					},
				},
				Contents: "Some additional notes here.",
			},
			card: &Card{
				ID:           "BtIrmFTJo49QuJC4",
				Path:         "further-maths/core-pure-1/chapter-1-complex-numbers/ex1a/question-BtIrmFTJo49QuJC4",
				Date:         time.Date(2021, 02, 16, 10, 18, 0, 0, time.Now().Location()),
				QuestionPath: "question.png",
				AnswerPath:   "answer.png",
				CompletionsPerfect: []Completion{
					{
						Date:     time.Date(2021, 02, 16, 10, 18, 0, 0, time.Now().Location()),
						Duration: 7*time.Minute + 10*time.Second,
					},
				},
				CompletionsMinor: []Completion{
					{
						Date:     time.Date(2021, 02, 16, 10, 18, 0, 0, time.Now().Location()),
						Duration: 5*time.Minute + 51*time.Second,
					},
				},
				CompletionsMajor: []Completion{
					{
						Date:     time.Date(2021, 02, 16, 10, 18, 0, 0, time.Now().Location()),
						Duration: 5*time.Minute + 53*time.Second,
					},
					{
						Date:     time.Date(2021, 02, 16, 10, 18, 0, 0, time.Now().Location()),
						Duration: 5*time.Minute + 53*time.Second,
					},
				},
				Notes: "Some additional notes here.",
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Entry(%d)", i), func(t *testing.T) {
			got, err := cardFromEntry(tc.entry)
			if !assert.NoError(t, err, "wasn't expecting an error when parsing valid entry") {
				return
			}

			assertCardEqual(t, tc.card, got)
		})
	}
}

// TestCardFromEntryInvalid tests that the cardFromEntry function returns errors for invalid entries.
func TestCardFromEntryInvalid(t *testing.T) {
	testCases := []struct {
		entry *entries.Entry
		err   string
		name  string
	}{
		// Use this is as a template for invalid entries.
		// {
		//     // This entry ...
		//     entry: &entries.Entry{
		//         Path:  "further-maths/core-pure-1/chapter-1-complex-numbers/ex1a/question-BtIrmFTJo49QuJC4",
		//         Title: "Question BtIrmFTJo49QuJC4",
		//         Date:  time.Date(2021, 02, 16, 10, 18, 0, 0, time.Now().Location()),
		//         Attachments: []entries.Attachment{
		//             {AbsPath: "answer.png", Name: "answer.png"},
		//             {AbsPath: "question.png", Name: "question.png"},
		//         },
		//         Metadata: map[string]interface{}{
		//             "type": "question",
		//             "completions": map[string][]map[string]string{
		//                 "perfect": {
		//                     {
		//                         "date": "2021-02-16 10:18",
		//                         "time": "7m10s",
		//                     },
		//                 },
		//                 "minor": {
		//                     {
		//                         "date": "2021-02-16 10:18",
		//                         "time": "5m23s",
		//                     },
		//                 },
		//                 "major": {
		//                     {
		//                         "date": "2021-02-16 10:18",
		//                         "time": "5m53s",
		//                     },
		//                     {
		//                         "date": "2021-02-16 10:18",
		//                         "time": "5m53s",
		//                     },
		//                 },
		//             },
		//         },
		//         Contents: "Some additional notes here.",
		//     },
		//     err:  "",
		//     name: "",
		// },
		{
			// This entry is nil.
			entry: nil,
			err:   "entry is nil",
			name:  "NilEntry",
		},
		{
			// This entry is missing a valid title, one that starts with "Question".
			entry: &entries.Entry{
				Path:  "further-maths/core-pure-1/chapter-1-complex-numbers/ex1a/question-BtIrmFTJo49QuJC4",
				Title: "Example Invalid Title",
				Date:  time.Date(2021, 02, 16, 10, 18, 0, 0, time.Now().Location()),
				Attachments: []entries.Attachment{
					{AbsPath: "question.png", Name: "question.png"},
					{AbsPath: "answer.png", Name: "answer.png"},
				},
				Metadata: map[string]interface{}{
					"type": "question",
					"completions": map[string][]map[string]string{
						"perfect": {
							{
								"date": "2021-02-16 10:18",
								"time": "7m10s",
							},
						},
						"minor": {
							{
								"date": "2021-02-16 10:18",
								"time": "5m51s",
							},
						},
						"major": {
							{
								"date": "2021-02-16 10:18",
								"time": "5m53s",
							},
							{
								"date": "2021-02-16 10:18",
								"time": "5m53s",
							},
						},
					},
				},
				Contents: "Some additional notes here.",
			},
			err:  "expected card title to start with 'Question', it is \"Example Invalid Title\"",
			name: "InvalidTitle",
		},
		{
			// This entry is missing an answer attachment.
			entry: &entries.Entry{
				Path:  "further-maths/core-pure-1/chapter-1-complex-numbers/ex1a/question-BtIrmFTJo49QuJC4",
				Title: "Question BtIrmFTJo49QuJC4",
				Date:  time.Date(2021, 02, 16, 10, 18, 0, 0, time.Now().Location()),
				Attachments: []entries.Attachment{
					{AbsPath: "question.png", Name: "question.png"},
				},
				Metadata: map[string]interface{}{
					"type": "question",
					"completions": map[string][]map[string]string{
						"perfect": {
							{
								"date": "2021-02-16 10:18",
								"time": "7m10s",
							},
						},
						"minor": {
							{
								"date": "2021-02-16 10:18",
								"time": "5m51s",
							},
						},
						"major": {
							{
								"date": "2021-02-16 10:18",
								"time": "5m53s",
							},
							{
								"date": "2021-02-16 10:18",
								"time": "5m53s",
							},
						},
					},
				},
				Contents: "Some additional notes here.",
			},
			err:  "card entry has no 'answer' attachment",
			name: "MissingAnswerAttachment",
		},
		{
			// This entry is missing a question attachment.
			entry: &entries.Entry{
				Path:  "further-maths/core-pure-1/chapter-1-complex-numbers/ex1a/question-BtIrmFTJo49QuJC4",
				Title: "Question BtIrmFTJo49QuJC4",
				Date:  time.Date(2021, 02, 16, 10, 18, 0, 0, time.Now().Location()),
				Attachments: []entries.Attachment{
					{AbsPath: "answer.png", Name: "answer.png"},
				},
				Metadata: map[string]interface{}{
					"type": "question",
					"completions": map[string][]map[string]string{
						"perfect": {
							{
								"date": "2021-02-16 10:18",
								"time": "7m10s",
							},
						},
						"minor": {
							{
								"date": "2021-02-16 10:18",
								"time": "5m51s",
							},
						},
						"major": {
							{
								"date": "2021-02-16 10:18",
								"time": "5m53s",
							},
							{
								"date": "2021-02-16 10:18",
								"time": "5m53s",
							},
						},
					},
				},
				Contents: "Some additional notes here.",
			},
			err:  "card entry has no 'question' attachment",
			name: "MissingQuestionAttachment",
		},

		{
			// This entry is missing the date field in one of the "perfect" completions.
			entry: &entries.Entry{
				Path:  "further-maths/core-pure-1/chapter-1-complex-numbers/ex1a/question-BtIrmFTJo49QuJC4",
				Title: "Question BtIrmFTJo49QuJC4",
				Date:  time.Date(2021, 02, 16, 10, 18, 0, 0, time.Now().Location()),
				Attachments: []entries.Attachment{
					{AbsPath: "answer.png", Name: "answer.png"},
					{AbsPath: "question.png", Name: "question.png"},
				},
				Metadata: map[string]interface{}{
					"type": "question",
					"completions": map[string][]map[string]string{
						"perfect": {
							{
								"time": "7m10s",
							},
						},
						"minor": {
							{
								"date": "2021-02-16 10:18",
								"time": "5m51s",
							},
						},
						"major": {
							{
								"date": "2021-02-16 10:18",
								"time": "5m53s",
							},
							{
								"date": "2021-02-16 10:18",
								"time": "5m53s",
							},
						},
					},
				},
				Contents: "Some additional notes here.",
			},
			err:  "'date' field in \"perfect\" completion list is empty",
			name: "MissingCompletionDate",
		},
		{
			// This entry is has a malformed date field in one of the "perfect" completions.
			entry: &entries.Entry{
				Path:  "further-maths/core-pure-1/chapter-1-complex-numbers/ex1a/question-BtIrmFTJo49QuJC4",
				Title: "Question BtIrmFTJo49QuJC4",
				Date:  time.Date(2021, 02, 16, 10, 18, 0, 0, time.Now().Location()),
				Attachments: []entries.Attachment{
					{AbsPath: "answer.png", Name: "answer.png"},
					{AbsPath: "question.png", Name: "question.png"},
				},
				Metadata: map[string]interface{}{
					"type": "question",
					"completions": map[string][]map[string]string{
						"perfect": {
							{
								"date": "2021-02-16T10:18", // <- Here
								"time": "7m10s",
							},
						},
						"minor": {
							{
								"date": "2021-02-16 10:18",
								"time": "5m51s",
							},
						},
						"major": {
							{
								"date": "2021-02-16 10:18",
								"time": "5m53s",
							},
							{
								"date": "2021-02-16 10:18",
								"time": "5m53s",
							},
						},
					},
				},
				Contents: "Some additional notes here.",
			},
			err:  "'date' field \"2021-02-16T10:18\" in \"perfect\" completion list not a valid '2006-01-02 15:04' date: parsing time \"2021-02-16T10:18\" as \"2006-01-02 15:04\": cannot parse \"T10:18\" as \" \"",
			name: "InvalidCompletionDate",
		},
		{
			// This entry is missing the time field in one of the "minor" completions.
			entry: &entries.Entry{
				Path:  "further-maths/core-pure-1/chapter-1-complex-numbers/ex1a/question-BtIrmFTJo49QuJC4",
				Title: "Question BtIrmFTJo49QuJC4",
				Date:  time.Date(2021, 02, 16, 10, 18, 0, 0, time.Now().Location()),
				Attachments: []entries.Attachment{
					{AbsPath: "answer.png", Name: "answer.png"},
					{AbsPath: "question.png", Name: "question.png"},
				},
				Metadata: map[string]interface{}{
					"type": "question",
					"completions": map[string][]map[string]string{
						"perfect": {
							{
								"date": "2021-02-16 10:18",
								"time": "7m10s",
							},
						},
						"minor": {
							{
								"date": "2021-02-16 10:18",
							},
						},
						"major": {
							{
								"date": "2021-02-16 10:18",
								"time": "5m53s",
							},
							{
								"date": "2021-02-16 10:18",
								"time": "5m53s",
							},
						},
					},
				},
				Contents: "Some additional notes here.",
			},
			err:  "'time' field in \"minor\" completion list is empty",
			name: "MissingCompletionTime",
		},
		{
			// This entry has an invalid time field in one of the "minor" completions.
			entry: &entries.Entry{
				Path:  "further-maths/core-pure-1/chapter-1-complex-numbers/ex1a/question-BtIrmFTJo49QuJC4",
				Title: "Question BtIrmFTJo49QuJC4",
				Date:  time.Date(2021, 02, 16, 10, 18, 0, 0, time.Now().Location()),
				Attachments: []entries.Attachment{
					{AbsPath: "answer.png", Name: "answer.png"},
					{AbsPath: "question.png", Name: "question.png"},
				},
				Metadata: map[string]interface{}{
					"type": "question",
					"completions": map[string][]map[string]string{
						"perfect": {
							{
								"date": "2021-02-16 10:18",
								"time": "7m10s",
							},
						},
						"minor": {
							{
								"date": "2021-02-16 10:18",
								"time": "Invalid time field.",
							},
						},
						"major": {
							{
								"date": "2021-02-16 10:18",
								"time": "5m53s",
							},
							{
								"date": "2021-02-16 10:18",
								"time": "5m53s",
							},
						},
					},
				},
				Contents: "Some additional notes here.",
			},
			err:  "'time' field \"Invalid time field.\" in \"minor\" completion list not a valid duration: time: invalid duration \"Invalid time field.\"",
			name: "InvalidTimeField",
		},
		{
			// This entry is missing the required "type" metadata.
			entry: &entries.Entry{
				Path:  "further-maths/core-pure-1/chapter-1-complex-numbers/ex1a/question-BtIrmFTJo49QuJC4",
				Title: "Question BtIrmFTJo49QuJC4",
				Date:  time.Date(2021, 02, 16, 10, 18, 0, 0, time.Now().Location()),
				Attachments: []entries.Attachment{
					{AbsPath: "answer.png", Name: "answer.png"},
					{AbsPath: "question.png", Name: "question.png"},
				},
				Metadata: map[string]interface{}{
					"completions": map[string][]map[string]string{
						"perfect": {
							{
								"date": "2021-02-16 10:18",
								"time": "7m10s",
							},
						},
						"minor": {
							{
								"date": "2021-02-16 10:18",
								"time": "5m23s",
							},
						},
						"major": {
							{
								"date": "2021-02-16 10:18",
								"time": "5m53s",
							},
							{
								"date": "2021-02-16 10:18",
								"time": "5m53s",
							},
						},
					},
				},
				Contents: "Some additional notes here.",
			},
			err:  "missing required 'type' field in card entry metadata",
			name: "MissingTypeMetadata",
		},
		{
			// This entry has the wrong "type" metadata.
			entry: &entries.Entry{
				Path:  "further-maths/core-pure-1/chapter-1-complex-numbers/ex1a/question-BtIrmFTJo49QuJC4",
				Title: "Question BtIrmFTJo49QuJC4",
				Date:  time.Date(2021, 02, 16, 10, 18, 0, 0, time.Now().Location()),
				Attachments: []entries.Attachment{
					{AbsPath: "answer.png", Name: "answer.png"},
					{AbsPath: "question.png", Name: "question.png"},
				},
				Metadata: map[string]interface{}{
					"type": "example-wrong-type",
					"completions": map[string][]map[string]string{
						"perfect": {
							{
								"date": "2021-02-16 10:18",
								"time": "7m10s",
							},
						},
						"minor": {
							{
								"date": "2021-02-16 10:18",
								"time": "5m23s",
							},
						},
						"major": {
							{
								"date": "2021-02-16 10:18",
								"time": "5m53s",
							},
							{
								"date": "2021-02-16 10:18",
								"time": "5m53s",
							},
						},
					},
				},
				Contents: "Some additional notes here.",
			},
			err:  "expected metadata field 'type' to be 'question', not \"example-wrong-type\"",
			name: "InvalidTypeMetadata",
		},
		{
			// This entry has an invalid "completions" type.
			entry: &entries.Entry{
				Path:  "further-maths/core-pure-1/chapter-1-complex-numbers/ex1a/question-BtIrmFTJo49QuJC4",
				Title: "Question BtIrmFTJo49QuJC4",
				Date:  time.Date(2021, 02, 16, 10, 18, 0, 0, time.Now().Location()),
				Attachments: []entries.Attachment{
					{AbsPath: "answer.png", Name: "answer.png"},
					{AbsPath: "question.png", Name: "question.png"},
				},
				Metadata: map[string]interface{}{
					"type":        "question",
					"completions": 10,
				},
				Contents: "Some additional notes here.",
			},
			err:  "couldn't parse 'completions' completions in card entry metadata",
			name: "InvalidCompletionsMetadata",
		},
		{
			// This entry has an unexpected completions field.
			entry: &entries.Entry{
				Path:  "further-maths/core-pure-1/chapter-1-complex-numbers/ex1a/question-BtIrmFTJo49QuJC4",
				Title: "Question BtIrmFTJo49QuJC4",
				Date:  time.Date(2021, 02, 16, 10, 18, 0, 0, time.Now().Location()),
				Attachments: []entries.Attachment{
					{AbsPath: "answer.png", Name: "answer.png"},
					{AbsPath: "question.png", Name: "question.png"},
				},
				Metadata: map[string]interface{}{
					"type": "question",
					"completions": map[string][]map[string]string{
						"perfect": {
							{
								"date": "2021-02-16 10:18",
								"time": "7m10s",
							},
						},
						"minor": {
							{
								"date": "2021-02-16 10:18",
								"time": "5m23s",
							},
						},
						"major": {
							{
								"date": "2021-02-16 10:18",
								"time": "5m53s",
							},
							{
								"date": "2021-02-16 10:18",
								"time": "5m53s",
							},
						},
						"unexpected-field-example": {
							{},
						},
					},
				},
				Contents: "Some additional notes here.",
			},
			err:  "not expecting completions field \"unexpected-field-example\" in card metadata",
			name: "InvalidCompletionsField",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := cardFromEntry(tc.entry)
			if err == nil {
				t.Errorf("expected an error when parsing an invalid entry, got nil")
				return
			}

			assert.Equal(t, tc.err, err.Error(), "expected different error when parsing invalid entry")
		})
	}
}

// TestCardParentPath tests that the card.PathParent function returns the correction path.
func TestCardParentPath(t *testing.T) {
	card := &Card{Path: "further-maths/core-pure-1/chapter-1-complex-numbers/ex1a/question-BtIrmFTJo49QuJC4"}

	assert.Equal(t, "further-maths/core-pure-1/chapter-1-complex-numbers/ex1a", card.PathParent(), "expected path to the parent to be correct")
}

// TestCardMarshalConsistent tests that a card is the same when being converted between multiple formats.
func TestCardMarshalConsistent(t *testing.T) {
	originalEntryContent := `---
title: "Question <random 16-character string>"
type: "question"
tags: ["@?any-tags"]
date: "2021-02-18 12:33"
completions:
    perfect:
        - date: 2021-02-16 10:18
          time: 7m10s
    minor:
        - date: 2021-02-16 10:18
          time: 5m51s
    major:
        - date: 2021-02-16 10:18
          time: 5m53s
        - date: 2021-02-16 10:18
          time: 5m53s
---

Any additional notes about the card, optional. For example:
Didn't realise that the partial fractions added up along the diagonal.`

	parser, err := entries.NewParser("2006-01-02 15:04", "@?")
	if err != nil {
		t.Errorf("not expecting error parsing original entry: %s", err)
		return
	}

	originalEntry, err := parser.Parse("testing/test-entry", originalEntryContent)
	if err != nil {
		t.Errorf("not expecting error parsing original entry: %s", err)
		return
	}

	// We have to add a couple of fake attachments or it won't be parsed properly.
	originalEntry.Attachments = []entries.Attachment{
		{AbsPath: "answer.png", Name: "answer.png"},
		{AbsPath: "question.png", Name: "question.png"},
	}

	originalCard, err := cardFromEntry(originalEntry)
	if err != nil {
		t.Errorf("not expecting error converting original test entry to card: %s", err)
	}

	newEntryContent, err := originalCard.Content()
	if err != nil {
		t.Errorf("not expecting error converting test card to entry: %s", err)
	}

	newEntry, err := parser.Parse("testing/test-entry", newEntryContent)
	if err != nil {
		t.Errorf("not expecting error parsing new entry: %s", err)
		return
	}

	newEntry.Attachments = []entries.Attachment{
		{AbsPath: "answer.png", Name: "answer.png"},
		{AbsPath: "question.png", Name: "question.png"},
	}

	newCard, err := cardFromEntry(newEntry)
	if err != nil {
		t.Errorf("not expecting error converting new test entry to card: %s", err)
	}

	assertCardEqual(t, originalCard, newCard)
}

func assertCardEqual(t *testing.T, card1, card2 *Card) {
	assert.Equal(t, card1.ID, card2.ID, "expected IDs to be the same")
	assert.Equal(t, card1.Path, card2.Path, "expected paths to be the same")
	assert.Equal(t, card1.Tags, card2.Tags, "expected tags to be the same")
	assert.Equal(t, card1.Notes, card2.Notes, "expected notes to be the same")

	assert.Equal(t, card1.Date.Format("2006-01-02 15:04"), card2.Date.Format("2006-01-02 15:04"), "expected date to be the same")

	assertCompletionsEqual(t, "perfect", card1.CompletionsPerfect, card2.CompletionsPerfect)
	assertCompletionsEqual(t, "minor", card1.CompletionsMinor, card2.CompletionsMinor)
	assertCompletionsEqual(t, "major", card1.CompletionsMajor, card2.CompletionsMajor)
}

func assertCompletionsEqual(t *testing.T, completionType string, completions1 []Completion, completions2 []Completion) {
	assert.Equal(t, len(completions1), len(completions2), "expected same number of completions")

	for i := range completions1 {
		completion1 := completions1[i]
		completion2 := completions2[i]

		// Using the string version of the data/duration here might be a bad idea. But it makes for nicer output.
		assert.Equal(t, completion1.Date.Format("2006-01-02 15:04"), completion2.Date.Format("2006-01-02 15:04"), "expected %q completions to have same date", completionType)
		assert.Equal(t, completion1.Duration.String(), completion2.Duration.String(), "expected %q completions to have same duration", completionType)
	}
}
