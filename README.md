# Pelican
## Architechture

```
. - Pelican API
cmd/pelican - Command for starting web server.
```

## API
### Adding
```sh
$ pelican add --path ... --question path/to/question.png --answer path/to/answer.png
```

```go
flashcards.Add(path, question, answer)
```

```
http://localhost:8080/flashcards/add
```

### Removing
```
$ pelican remove --path ...
```

```go
flashcards.Remove(path)
```

```
http://localhost:8080/flashcards/remove
```

## Example Layout
```
further-maths
        ├── edexcel-core-pure-1
        │   ├── chapter-4
        │   │   └── ex4a
        │   │       ├── question-jDxorrl9shZviWeW
        │   │       │   ├── answer.png
        │   │       │   ├── entry.md
        │   │       │   └── question.png
        │   │       └── question-SYXM6N9kYAUtCRSk
        │   │           ├── answer.png
        │   │           ├── entry.md
        │   │           └── question.png
        │   └── entry.md
        └── entry.md
```

## Hierarchy
| Name              | What                                          | Example                                                           |
|-------------------|-----------------------------------------------|-------------------------------------------------------------------|
| **Flashcard Set** | A list of flashcards, can be nested           | *Ex10A*, *Exam-style Practice*, *Pearson Core 1*, *Exam May 2020* |
| **Flashcard**     | Basic unit, contains a question and an answer | *If A is a vector...*, *A Newton's cradle...*                     |

### Flashcard
```go
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

Flashcard.Name() // Name of flashcard, built from the string of parents (i.e. Further Maths -> Core 2 -> Ex10A -> Question jDxorrl9shZviWeW)
```

Is created by parsing an entry with metadata like this:

```yaml
title: "Question jDxorrl9shZviWeW"
date: "2021-02-03 16:24"
perfect: []
minor: []
major: []

# Notes are in entry body...
# (QuestionImg and AnswerImg are inferred from attachments).
```

Every entry should have two attachments:

* `question.png`
* `answer.png`

### Flashcard Set
```go
type FlashcardSet struct {
    Name string // Name of the set, e.g. "Pearson Further Maths Core 2" or "Ex10A"
    
    Flashcards []*Flashcard // List of flashcards it's part of.

    Parent *FlashcardSet // If this flashcard set is in another flashcard set, this points to it.
    Children []*FlashcardSet // If this flashcard set has more flashcard sets inside, these point to those.

    Notes string // Any additional information.
}
```

Is created by parsing an entry that looks like this:

```md
---
title: "Pearson Further Maths Core 2"
tags: ["@?textbook", "@?further-maths"]
---

Notes...
```

