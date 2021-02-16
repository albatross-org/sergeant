package sergeant

import "math/rand"

// View is a certain way of scheduling cards.
type View interface {
	Next(set *Set) *Card
}

// Random selects a random card from all possible cards.
type Random struct{}

// Next looks at all previous cards and decided what card to show next.
func (view *Random) Next(set *Set) *Card {
	return set.Cards[rand.Intn(len(set.Cards))]
}
