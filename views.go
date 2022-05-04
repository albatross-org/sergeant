package sergeant

import (
	"math"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/dghubble/trie"
	exprand "golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

// DefaultViews is a map containing the default Views used by the program.
var DefaultViews = map[string]View{
	"random":   NewViewRandom(time.Now().Unix()),
	"unseen":   NewViewUnseen(time.Now().Unix()),
	"bayesian": NewViewBayesian(time.Now().Unix()),
}

// View is a certain way of scheduling cards.
type View interface {
	Next(set *Set, user string) *Card
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
func (view *Random) Next(set *Set, user string) *Card {
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
func (view *Unseen) Next(set *Set, user string) *Card {
	candidates := []*Card{}

	for _, card := range set.Cards {
		if card.TotalCompletionsUser(user) == 0 {
			candidates = append(candidates, card)
		}
	}

	if len(candidates) == 0 {
		return nil
	}

	return candidates[view.rng.Intn(len(candidates))]
}

// EvaluationView is a View which assigns a numerical score to each node.
type EvaluationView interface {
	View
	Evaluate(path string) (float64, error)
}

// ProbabilityNode represents a node in the probability tree.
type ProbabilityNode struct {
	Path       string
	Perfect    int
	Minor      int
	Major      int
	Difficulty float64
}

// putOrUpdateTrie will put a trie value or update it if it already exists for this path.
func putOrUpdateTrie(trie *trie.PathTrie, path string, card *Card, user string) {
	if existing := trie.Get(path); existing != nil {
		existing := existing.(*ProbabilityNode)
		existing.Perfect += card.UserPerfect(user)
		existing.Minor += card.UserMinor(user)
		existing.Major += card.UserMajor(user)
		trie.Put(path, existing)
	} else {
		trie.Put(path, &ProbabilityNode{
			Path:    path,
			Perfect: card.UserPerfect(user),
			Minor:   card.UserMinor(user),
			Major:   card.UserMajor(user),
		})
	}
}

// Bayesian uses Bayesian inference in order to try and select the card you're most likely to get wrong.
// The method comes from Probabilistic Programming and Bayesian Methods for Hackers:
//   https://nbviewer.jupyter.org/github/CamDavidsonPilon/Probabilistic-Programming-and-Bayesian-Methods-for-Hackers/blob/master/Chapter6_Priorities/Ch6_Priors_PyMC2.ipynb
// The basic idea is to turn the card picking problem into a multi-armed bandit problem. Each path then becomes a "bandit" and has an associated
// probability distribution based on the questions answered correctly and the questions answered incorrectly.
type Bayesian struct {
	rng *rand.Rand
}

// bayesianBetaDistribution represents the beta distribution associated with one particular path. Basically, this is just a named beta distribution.
type bayesianBetaDistribution struct {
	Alpha int
	Beta  int

	Path string
}

// Next looks at all previous cards and decides what card to show next.
func (view *Bayesian) Next(set *Set, user string) *Card {
	pathTrie := trie.NewPathTrie()

	// Create a trie based on all the different paths for the cards.
	// We store a probabilityNode at each one which holds information about the success rates for each path.
	for _, card := range set.Cards {
		components := strings.Split(card.PathParent(), "/")
		for i := 1; i < len(components); i++ {
			putOrUpdateTrie(pathTrie, strings.Join(components[:i], "/"), card, user)
		}

	}

	priors := []bayesianBetaDistribution{}

	// Here we walk the tree in a breadth-first fashion and create a beta distribution for each one.
	pathTrie.Walk(func(key string, value interface{}) error {
		node := value.(*ProbabilityNode)

		priors = append(priors, bayesianBetaDistribution{
			Alpha: 1 + node.Perfect,
			Beta:  2 + node.Major + node.Minor,
			Path:  key,
		})

		return nil
	})

	pathMap := map[string][]*Card{}
	blacklisted := map[string]bool{}
	questions := []*Card{}

	for len(questions) == 0 {
		var smallestSample float64 = math.Inf(1)
		var smallestPath string

		for _, prior := range priors {
			dist := distuv.Beta{Alpha: float64(prior.Alpha), Beta: float64(prior.Beta), Src: exprand.NewSource(view.rng.Uint64())}
			sample := dist.Rand()

			if sample < smallestSample && !blacklisted[prior.Path] {
				smallestSample = sample
				smallestPath = prior.Path
			}

			for _, card := range set.Cards {
				if strings.HasPrefix(card.Path, prior.Path) && card.TotalCompletionsUser(user) == 0 {
					pathMap[prior.Path] = append(pathMap[prior.Path], card)
				}
			}
		}

		questions = pathMap[smallestPath]

		if len(questions) == 0 {
			blacklisted[smallestPath] = true
		}
	}

	return questions[view.rng.Intn(len(questions))]
}

// NewViewBayesian returns a new view based on Bayesian inference.
func NewViewBayesian(seed int64) *Bayesian {
	return &Bayesian{
		rng: rand.New(rand.NewSource(seed)),
	}
}

// CardEvaluation contains card information and the mean for the prior for that card.
type CardEvaluation struct {
	Path               string  `json:"path"`
	Mean               float64 `json:"mean"`
	QuestionsAvailable int     `json:"questionsAvailable"`
	QuestionsCompleted int     `json:"questionsCompleted"`
}

// Difficulties returns the most difficult topics for a Bayesian set.
func (view *Bayesian) Difficulties(set *Set, user string) []CardEvaluation {
	pathTrie := trie.NewPathTrie()

	// Create a trie based on all the different paths for the cards.
	// We store a probabilityNode at each one which holds information about the success rates for each path.
	for _, card := range set.Cards {
		components := strings.Split(card.PathParent(), "/")
		for i := 1; i < len(components); i++ {
			putOrUpdateTrie(pathTrie, strings.Join(components[:i], "/"), card, user)
		}

	}

	evaluations := []CardEvaluation{}

	// Here we walk the tree in a breadth-first fashion and create a beta distribution for each one.
	pathTrie.Walk(func(key string, value interface{}) error {
		node := value.(*ProbabilityNode)

		total := node.Perfect + node.Minor + node.Major
		alpha := 1 + node.Perfect
		beta := 1 + node.Major + node.Minor
		available := 0

		for _, card := range set.Cards {
			if strings.HasPrefix(card.Path, key) {
				available += 1
			}
		}

		if total > 0 {
			evaluations = append(evaluations, CardEvaluation{
				Path:               key,
				Mean:               float64(alpha) / float64(alpha+beta),
				QuestionsCompleted: total,
				QuestionsAvailable: available,
			})
		}

		return nil
	})

	sort.Slice(evaluations, func(i, j int) bool {
		if evaluations[i].Mean != evaluations[j].Mean {
			return evaluations[i].Mean < evaluations[j].Mean
		} else {
			return evaluations[i].Path < evaluations[j].Path
		}
	})

	return evaluations
}
