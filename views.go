package sergeant

import (
	"math"
	"math/rand"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/dghubble/trie"
)

// DefaultViews is a map containing the default Views used by the program.
var DefaultViews = map[string]View{
	"random":       NewViewRandom(time.Now().Unix()),
	"unseen":       NewViewUnseen(time.Now().Unix()),
	"difficulties": NewViewDifficulties(time.Now().Unix()),
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
	return set.Cards[view.rng.Intn(len(set.Cards))]
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

	return candidates[view.rng.Intn(len(candidates))]
}

// Difficulties selects cards that it predicts you are most likely to get wrong.
// TODO: make this better. At the moment it will focus on one specific area.
type Difficulties struct {
	rng *rand.Rand
}

// NewViewDifficulties returns a new Difficulties view with the given seed.
func NewViewDifficulties(seed int64) *Difficulties {
	return &Difficulties{
		rng: rand.New(rand.NewSource(seed)),
	}
}

// probabilityNode represents a node in the probability tree.
type probabilityNode struct {
	Path       string
	Perfect    int
	Minor      int
	Major      int
	Difficulty float64
}

// putOrUpdateTrie will put a trie value or update it if it already exists for this path.
func putOrUpdateTrie(trie *trie.PathTrie, path string, card *Card) {
	if existing := trie.Get(path); existing != nil {
		existing := existing.(*probabilityNode)
		existing.Perfect += len(card.CompletionsPerfect)
		existing.Minor += len(card.CompletionsMinor)
		existing.Major += len(card.CompletionsMajor)
		trie.Put(path, existing)
	} else {
		trie.Put(path, &probabilityNode{
			Path:    path,
			Perfect: len(card.CompletionsPerfect),
			Minor:   len(card.CompletionsMinor),
			Major:   len(card.CompletionsMajor),
		})
	}
}

// Next looks at all previous cards and decides what card to show next.
func (view *Difficulties) Next(set *Set) *Card {
	trie := trie.NewPathTrie()

	for _, card := range set.Cards {
		components := strings.Split(card.PathParent(), "/")
		for i := 1; i < len(components); i++ {
			putOrUpdateTrie(trie, strings.Join(components[:i], "/"), card)
		}

	}

	var lowest *probabilityNode
	var paths []string

	trie.Walk(func(key string, value interface{}) error {
		node := value.(*probabilityNode)

		var generalProbability float64
		var sampleProbability float64

		if strings.Count(key, "/") == 0 {
			generalProbability = 0.5
		} else {
			parent := trie.Get(filepath.Dir(key))
			generalProbability = parent.(*probabilityNode).Difficulty
		}

		total := node.Perfect + node.Minor + node.Major

		if total == 0 {
			// If we have no sample, we use a reduced version of the the parents probability. This means
			// that the program will sometimes pick categories that haven't been looked at yet.
			node.Difficulty = generalProbability * 0.6
		} else {
			// If we have a sample, we compute an adjusted difficulty probability that takes into account
			// the overall probability of the underlying category.
			sampleProbability = float64(node.Perfect) / float64(total)
			node.Difficulty = adjustDifficultyProbability(generalProbability, sampleProbability, total)
		}

		if lowest == nil || node.Difficulty < lowest.Difficulty {
			lowest = node
		}

		paths = append(paths, key)

		return nil
	})

	sort.Slice(paths, func(i, j int) bool {
		p1 := trie.Get(paths[i]).(*probabilityNode).Difficulty
		p2 := trie.Get(paths[j]).(*probabilityNode).Difficulty

		return p1 < p2
	})

	path := paths[:3][view.rng.Intn(len(paths[:3]))]

	candidates := []*Card{}
	for _, card := range set.Cards {
		if strings.HasPrefix(card.Path, path) {
			candidates = append(candidates, card)
		}
	}

	return candidates[0]
}

// adjustDifficultyProbability combines an overall probability with the probability from a sample in order to generate a difficulty level.
// The idea here is that we can calculate the difficulty of a question using a combination of the general probability for a question being answered
// correctly in that category and a prediction that's based on a sample size.
// For a more in-depth explanation of why this formula is used, see https://www.desmos.com/calculator/fouzmkkbo8.
func adjustDifficultyProbability(generalProbability float64, sampleProbability float64, sampleSize int) float64 {
	sampleStrength := 1.0 / (1 + math.Exp(2-float64(sampleSize)/3)) // Moved sigmoid curve.
	return sampleProbability*sampleStrength + generalProbability*(1-sampleStrength)
}
