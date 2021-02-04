package pelican

import (
	"fmt"
	"strings"

	"github.com/albatross-org/go-albatross/entries"
)

// Set represents a set of flashcards. It can either contain more flashcard sets or a list of flashcards.
type Set struct {
	// Name is the name of the set, e.g. "Pearson Further Maths Core 2" or "Ex10A"
	Name string

	// Path is the path to the set, used for making updates later.
	Path string

	// Flashcards is a list of flashcards contained in the set.
	Flashcards []*Flashcard

	// Sets contains all the child flashcard sets contained within the set.
	Sets []*Set

	// Parent is the set that this is a child of. If it's nil, then this is the "root set".
	Parent *Set

	// Notes holds any additional information.
	Notes string
}

// List returns the list of all flashcards contained within the set, including those that are nested.
func (s *Set) List() []*Flashcard {
	flashcards := s.Flashcards

	for _, set := range s.Sets {
		flashcards = append(flashcards, set.List()...)
	}

	return flashcards
}

// setFromTree returns a new flashcard *Set from an entries.Tree.
// When calling, parent should be passed as nil.
func setFromTree(tree *entries.Tree, parent *Set, depth int) (*Set, map[string]error, map[string]error, error) {
	set := &Set{}
	var name string
	var contents string

	fmt.Println("tree", tree.Path, "is entry", tree.IsEntry)

	// If it's a passthrough entry, we use the last bit of the path as the name.
	// E.g. 'further-maths/edexcel-core-pure-1/chapter-4' -> 'chapter-4'
	// Otherwise we just use the name of the entry.
	if !tree.IsEntry {
		if tree.Path == "" {
			name = "All" // This is the top-level entry/set.
		} else {
			components := strings.Split(tree.Path, "/")
			name = components[len(components)-1]
		}
	} else {
		name = tree.Entry.Title
		contents = tree.Entry.Contents
	}

	// If a single flashcard or a single set is wrong, we don't want to prevent
	// the rest of the tree from being parsed. Instead we report it as an individual
	// error.
	flashcards := []*Flashcard{}
	flashcardErrs := map[string]error{}

	sets := []*Set{}
	setErrs := map[string]error{}

	for _, child := range tree.Children {
		if isFlashcard(child.Entry) {
			flashcard, err := flashcardFromEntry(child.Entry)
			if err != nil {
				flashcardErrs[child.Path] = err
				continue
			}

			fmt.Println("I'm here", set)
			flashcard.Parent = set

			flashcards = append(flashcards, flashcard)
		} else {
			childSet, childFlashcardErrs, childSetErrs, err := setFromTree(child, set, depth+1) // 'set' is already defined since it's a named return value.
			if err != nil {
				setErrs[child.Path] = err
				continue
			}

			for path, flashcardErr := range childFlashcardErrs {
				flashcardErrs[path] = flashcardErr
			}

			for path, setErr := range childSetErrs {
				setErrs[path] = setErr
			}

			sets = append(sets, childSet)

		}
	}

	set.Name = name
	set.Path = tree.Path
	set.Flashcards = flashcards
	set.Sets = sets
	set.Parent = parent
	set.Notes = contents

	return set, flashcardErrs, setErrs, nil
}
