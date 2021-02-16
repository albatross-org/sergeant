package sergeant

// Set represents a collection of cards.
type Set struct {
	Cards []*Card
}

// Filter returns a filtered subset of the set.
func (s *Set) Filter(filters ...Filter) *Set {
	return s
}
