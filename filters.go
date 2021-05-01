package sergeant

import (
	"strings"
	"time"
)

// Filter represents a way of allowing or disallowing a card.
type Filter func(*Card) bool

// FilterOR combines multiple filters together and if any of them match the given entry, it will be allowed.
func FilterOR(filters ...Filter) Filter {
	return func(card *Card) bool {
		for _, filter := range filters {
			if filter(card) {
				return true
			}
		}

		return false
	}
}

// FilterAND combines multiple filters together and they all have to match an entry for it to be allowed.
func FilterAND(filters ...Filter) Filter {
	return func(card *Card) bool {
		for _, filter := range filters {
			if !filter(card) {
				return false
			}
		}

		return true
	}
}

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

// FilterTags returns a filter that only allows cards who contain the given tags.
// This is an AND operations -- all of the tags have to match for the card to be allowed.
func FilterTags(tags ...string) Filter {
	return func(card *Card) bool {
		for _, tag := range tags {
			found := false

			for _, cardTag := range card.Tags {
				if tag == cardTag {
					found = true
					break
				}
			}

			if !found {
				return false
			}
		}

		return true
	}
}

// FilterBeforeDate returns a filter that only allows cards that were created before a certain date.
func FilterBeforeDate(date time.Time) Filter {
	return func(card *Card) bool {
		return card.Date.Before(date)
	}
}

// FilterAfterDate returns a filter that only allows cards that were created before a certain date.
func FilterAfterDate(date time.Time) Filter {
	return func(card *Card) bool {
		return card.Date.After(date)
	}
}

// FilterBeforeDuration returns a filter that only allows cards created a certain amount of time ago.
// For example
//   FilterBeforeDuration(10 * 24 * time.Hour)
// only lets cards created more than 10 days ago be used.
func FilterBeforeDuration(duration time.Duration) Filter {
	return func(card *Card) bool {
		return time.Since(card.Date) > duration
	}
}

// FilterAfterDuration returns a filter that only allows cards created a certain amount of time ago.
// For example
//   FilterAfterDuration(30 * 24 * time.Hour)
// only lets cards created less than 30 days ago be used.
func FilterAfterDuration(duration time.Duration) Filter {
	return func(card *Card) bool {
		return time.Since(card.Date) > duration
	}
}
