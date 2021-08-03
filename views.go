package sergeant

import (
	"math"
	"math/rand"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/dghubble/trie"
	"github.com/sirupsen/logrus"
	exprand "golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"

	wr "github.com/mroth/weightedrand"
)

// DefaultViews is a map containing the default Views used by the program.
var DefaultViews = map[string]View{
	"random":       NewViewRandom(time.Now().Unix()),
	"unseen":       NewViewUnseen(time.Now().Unix()),
	"difficulties": NewViewDifficulties(time.Now().Unix()),
	"bayesian":     NewViewBayesian(time.Now().Unix()),
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

// EvaluationView is a View which assigns a numerical score to each node.
type EvaluationView interface {
	View
	Evaluate(path string) (float64, error)
}

// Difficulties selects cards that it predicts you are most likely to get wrong.
// This is done through a combination of the overall probability for a certain category (like "futher-maths/section-complex-numbers" for example)
// and a specific probablility (like "further-maths/section-complex-numbers/chapter-2-argand-diagrams"). The "strength" of the sample size is
// computed in order to work out how much weight should be given to the sample instead of the general probability.
// For a visual and interactive explanation, see https://www.desmos.com/calculator/fouzmkkbo8.
type Difficulties struct {
	rng *rand.Rand

	// baseProbability is the probability assumed probability "one level higher" than the root of the tree.
	baseProbability float64

	// assumedSampleProbability is the probability used when no sample has been given.
	// By default, this is 0.6. This means that categories which haven't been questioned about will be explored.
	assumedSampleProbability float64

	// topPercent is the proportion of cards (out of 1) that are selected to be part of cards chosen.
	topPercent float64

	// power is what the numerical difficulties are raised to. A higher value means that lower difficulties are exagerrated.
	power float64
}

// NewViewDifficulties returns a new Difficulties view with the given seed.
func NewViewDifficulties(seed int64) *Difficulties {
	return &Difficulties{
		rng:                      rand.New(rand.NewSource(seed)),
		baseProbability:          0.5,
		assumedSampleProbability: 0.5,
		topPercent:               0.4,
		power:                    2,
	}
}

// ProbabilityNode represents a node in the probability tree.
type ProbabilityNode struct {
	Path       string
	Perfect    int
	Minor      int
	Major      int
	Difficulty float64
}

// BuildTrie creates a path trie that contains difficulty information.
// A trie is a "prefix tree", where each all elements that share a common prefix are grouped under the same parent.
// This uses a question's path to create a hierarchy of all the possible questions.
// It will also return a sorted slice of all paths present.
func (view *Difficulties) BuildTrie(set *Set) (pathTrie *trie.PathTrie, paths []string) {
	pathTrie = trie.NewPathTrie()

	// Create a trie based on all the different paths for the cards.
	// We store a probabilityNode at each one which holds information about the success rates for each path.
	for _, card := range set.Cards {
		components := strings.Split(card.PathParent(), "/")
		for i := 1; i < len(components); i++ {
			putOrUpdateTrie(pathTrie, strings.Join(components[:i], "/"), card)
		}

	}

	// Here we walk the tree in a breadth-first fashion. The idea here is to adjust the difficulty probabilities according
	// to more general samples from broader categories.
	// More information here: https://www.desmos.com/calculator/fouzmkkbo8
	pathTrie.Walk(func(key string, value interface{}) error {
		node := value.(*ProbabilityNode)

		var generalProbability float64
		var sampleProbability float64

		if strings.Count(key, "/") == 0 {
			generalProbability = view.baseProbability
		} else {
			parent := pathTrie.Get(filepath.Dir(key))
			generalProbability = parent.(*ProbabilityNode).Difficulty
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

		paths = append(paths, key)

		return nil
	})

	// Sort the available paths by their difficulty in the probability trie.
	// This means that the most difficult cards (those with the lowest probability) will come first.
	sort.Slice(paths, func(i, j int) bool {
		p1 := math.Pow(pathTrie.Get(paths[i]).(*ProbabilityNode).Difficulty, view.power)
		p2 := math.Pow(pathTrie.Get(paths[j]).(*ProbabilityNode).Difficulty, view.power)

		return p1 < p2
	})

	return pathTrie, paths
}

// Next looks at all previous cards and decides what card to show next.
func (view *Difficulties) Next(set *Set) *Card {
	pathTrie, paths := view.BuildTrie(set)

	if len(paths) == 0 {
		return nil
	}

	pathMap := map[string][]*Card{}

	// Build up a map of paths to all the available cards that they contain.
	// TODO: this is an expensive operation.
	for _, path := range paths {
		for _, card := range set.Cards {
			if strings.HasPrefix(card.Path, path) && card.TotalCompletions() == 0 {
				pathMap[path] = append(pathMap[path], card)
			}
		}
	}

	pathsSubset := paths[:int(math.Ceil(float64(len(paths))*view.topPercent))]

	choices := []wr.Choice{}
	for _, path := range pathsSubset {
		weightInt := uint(pathTrie.Get(path).(*ProbabilityNode).Difficulty * (10000000)) // Have to convert difficulty to uint.
		choices = append(
			choices,
			wr.Choice{Item: path, Weight: weightInt},
		)
	}

	chooser, err := wr.NewChooser(choices...)
	if err != nil {
		logrus.Error("Error choosing question: ", err)
		return nil
	}

	length := 0
	for length == 0 {
		path := chooser.Pick().(string)
		questions := pathMap[path]
		length := len(questions)

		if length == 0 {
			continue
		}

		return questions[view.rng.Intn(length)]
	}

	return nil
}

// putOrUpdateTrie will put a trie value or update it if it already exists for this path.
func putOrUpdateTrie(trie *trie.PathTrie, path string, card *Card) {
	if existing := trie.Get(path); existing != nil {
		existing := existing.(*ProbabilityNode)
		existing.Perfect += len(card.CompletionsPerfect)
		existing.Minor += len(card.CompletionsMinor)
		existing.Major += len(card.CompletionsMajor)
		trie.Put(path, existing)
	} else {
		trie.Put(path, &ProbabilityNode{
			Path:    path,
			Perfect: len(card.CompletionsPerfect),
			Minor:   len(card.CompletionsMinor),
			Major:   len(card.CompletionsMajor),
		})
	}
}

// adjustDifficultyProbability combines an overall probability with the probability from a sample in order to generate a difficulty level.
// The idea here is that we can calculate the difficulty of a question using a combination of the general probability for a question being answered
// correctly in that category and a prediction that's based on a sample size.
// For a more in-depth explanation of why this formula is used, see https://www.desmos.com/calculator/fouzmkkbo8.
func adjustDifficultyProbability(generalProbability float64, sampleProbability float64, sampleSize int) float64 {
	sampleStrength := 1.0 / (1 + math.Exp(2-float64(sampleSize)/3)) // Moved sigmoid curve.
	return sampleProbability*sampleStrength + generalProbability*(1-sampleStrength)
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
func (view *Bayesian) Next(set *Set) *Card {
	pathTrie := trie.NewPathTrie()

	// Create a trie based on all the different paths for the cards.
	// We store a probabilityNode at each one which holds information about the success rates for each path.
	for _, card := range set.Cards {
		components := strings.Split(card.PathParent(), "/")
		for i := 1; i < len(components); i++ {
			putOrUpdateTrie(pathTrie, strings.Join(components[:i], "/"), card)
		}

	}

	priors := []bayesianBetaDistribution{}

	// Here we walk the tree in a breadth-first fashion and create a beta distribution for each one.
	pathTrie.Walk(func(key string, value interface{}) error {
		node := value.(*ProbabilityNode)

		// We use 2+ and 4+ rather than the typical 1+ because it we prioritize information gain for different topics.
		priors = append(priors, bayesianBetaDistribution{
			Alpha: 2 + node.Perfect,
			Beta:  4 + node.Major + node.Minor,
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
				if strings.HasPrefix(card.Path, prior.Path) && card.TotalCompletions() == 0 {
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
