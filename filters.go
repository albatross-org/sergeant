package sergeant

import "strings"

// Filter represents a way of allowing or disallowing a card.
type Filter func(*Card) bool

// FilterPaths returns a filter that only allows cards who's path begins with the paths specified.
// This is an OR operation -- if any of the paths given match, then the card is allowed.
func FilterPaths(paths ...string) Filter {
	return func(card *Card) bool {
		for _, path := range paths {
			if strings.HasPrefix(card.Path, path) {
				return true
			}
		}

		return false
	}
}
