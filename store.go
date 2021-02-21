package sergeant

import (
	"fmt"

	"github.com/albatross-org/go-albatross/albatross"
)

// Store is an abstraction over an *albatross.Store that allows for updating cards.
type Store struct {
	albatross *albatross.Store
	Config    Config
	Sets      map[string]ConfigSet
}

// NewStore returns a new Store from an *albatross.Store and a config.
func NewStore(store *albatross.Store, config Config) *Store {
	return &Store{
		albatross: store,
		Config:    config,
		Sets:      config.Sets,
	}
}

// Set returns the cards present in the set specified.
// It knows what cards you want in the set from the .Sets configuration.
// It returns a Set, followed by a map of warnings (paths -> parse errors) and an overall error if there was one.
func (store *Store) Set(name string) (*Set, map[string]error, error) {
	collection, err := store.albatross.Collection()
	if err != nil {
		return nil, nil, fmt.Errorf("couldn't create set from collection: %w", err)
	}

	if store.Sets[name].Name == "" {
		return nil, nil, fmt.Errorf("set %q not found in config", name)
	}

	filter := store.Sets[name].AsFilter()

	slice := collection.List().Slice()

	cards := []*Card{}
	warnings := map[string]error{}

	for _, entry := range slice {
		card, err := cardFromEntry(entry)
		if err != nil {
			warnings[entry.Path] = err
			continue
		}

		if filter(card) {
			cards = append(cards, card)
		}
	}

	return &Set{
		Cards: cards,
	}, warnings, nil
}

// AddCompletion adds a completion to an entry in the store.
func (store *Store) AddCompletion(path string, completionType string, completion Completion) error {
	entry, err := store.albatross.Get(path)
	if err != nil {
		return err
	}

	fmt.Println(entry.Attachments)

	card, err := cardFromEntry(entry)
	if err != nil {
		return err
	}

	switch completionType {
	case "perfect":
		card.CompletionsPerfect = append(card.CompletionsPerfect, completion)
	case "minor":
		card.CompletionsMinor = append(card.CompletionsMinor, completion)
	case "major":
		card.CompletionsMajor = append(card.CompletionsMajor, completion)
	default:
		return fmt.Errorf("invalid completion type %q", completionType)
	}

	newContent, err := card.Content()
	if err != nil {
		return err
	}

	err = store.albatross.Update(path, newContent)
	if err != nil {
		return err
	}

	return nil
}
