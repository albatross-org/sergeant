package pelican

import (
	"fmt"

	"github.com/albatross-org/go-albatross/albatross"
	"github.com/albatross-org/go-albatross/entries"
	"github.com/sirupsen/logrus"
)

// Load will load an Albatross store and return the top-level Set.
// It has the same rules for name and configPath that albatross.Load does.
func Load(name, configPath string) (*Set, error) {
	store, err := albatross.Load(name, configPath)
	if err != nil {
		return nil, fmt.Errorf("couldn't load store to turn into flashcards: %w", err)
	}

	collection, err := store.Collection()
	if err != nil {
		return nil, fmt.Errorf("couldn't load collection to turn into flashcards: %w", err)
	}

	tree := entries.ListAsTree(collection.List())

	set, flashcardErrs, setErrs, err := setFromTree(tree, nil, 0)
	if err != nil {
		return nil, err
	}

	for path, err := range flashcardErrs {
		logrus.Warn(path, "->", err)
	}

	for path, err := range setErrs {
		logrus.Warn(path, "->", err)
	}

	return set, nil
}
