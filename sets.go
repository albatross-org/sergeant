package sergeant

// Set represents a collection of cards.
type Set struct {
	Cards []*Card
}

// Filter returns a filtered subset of the set.
// Multiple filters represent an AND -- a card must match every filter to be added.
func (s *Set) Filter(filters ...Filter) *Set {
	cardsNew := []*Card{}
	for _, card := range s.Cards {
		allowed := true
		for _, filter := range filters {
			allowed = allowed && filter(card)
		}

		if allowed {
			cardsNew = append(cardsNew, card)
		}
	}

	return &Set{Cards: cardsNew}
}
