package sergeant

import (
	"math/rand"
	"time"
)

// DefaultViews is a map containing the default Views used by the program.
var DefaultViews = map[string]View{
	"random": NewViewRandom(time.Now().Unix()),
	"unseen": NewViewUnseen(time.Now().Unix()),
}

// View is a certain way of scheduling cards.
type View interface {
	Next(set *Set) *Card
}

// Random selects a random card from all possible cards.
type Random struct {
	rng *rand.Rand
}

// NewViewRandom returns a new Random view with the given seed.
func NewViewRandom(seed int64) *Random {
	return &Random{
		rng: rand.New(rand.NewSource(seed)),
	}
}

// Next looks at all previous cards and decides what card to show next.
func (view *Random) Next(set *Set) *Card {
	return set.Cards[rand.Intn(len(set.Cards))]
}

// Unseen selects cards that have yet to come up.
type Unseen struct {
	rng *rand.Rand
}

// NewViewUnseen returns a new Unseen view with the given seed.
func NewViewUnseen(seed int64) *Unseen {
	return &Unseen{
		rng: rand.New(rand.NewSource(seed)),
	}
}

// Next looks at all previous cards and decides what card to show next.
func (view *Unseen) Next(set *Set) *Card {
	candidates := []*Card{}

	for _, card := range set.Cards {
		if len(card.CompletionsMajor)+len(card.CompletionsMinor)+len(card.CompletionsPerfect) == 0 {
			candidates = append(candidates, card)
		}
	}

	if len(candidates) == 0 {
		return nil
	}

	return candidates[rand.Intn(len(candidates))]
}
