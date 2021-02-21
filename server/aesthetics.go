package server

import (
	"math/rand"

	"github.com/segmentio/fasthash/fnv1"
)

type colorPair struct {
	color      string
	background string
}

// colors holds the default colors and gradients that are used to make sets look pretty.
// It needs considerable expansion.
var colors = []colorPair{
	{"#ad5389", "linear-gradient(315deg, #3c1053, #ad5389)"},
	{"#00b4db", "linear-gradient(315deg, #0083b0, #00b4db)"},
	{"#f3f9a7", "linear-gradient(315deg, #cac531, #f3f9a7)"},
	{"#38ef7d", "linear-gradient(315deg, #11998e, #38ef7d)"},
	// {"#c94b4b", "linear-gradient(315deg, #c94b4b, #4b134f)"},
	// {"#fc4a1a", "linear-gradient(315deg, #fc4a1a, #f7b733)"},
	// {"#02aab0", "linear-gradient(315deg, #02aab0, #00cdac)"}, // Green Beach
	// {"#d31027", "linear-gradient(315deg, #d31027, #ea384d)"}, // Playing with Reds
	// {"#8e2de2", "linear-gradient(315deg, #8e2de2, #4a00e0)"}, // Amin
	// {"#536976", "linear-gradient(315deg, #536976, #292e49)"}, // Playing with Reds
}

// gradientColorPair decides on a random gradient and a color from the name of a set.
// In order to remain consistent between calls, the name is hashed and used to seed a random
// number generator, which is in turn used to pick a random.
// One issue is that as soon as the size of colors changes, all the set colors will also change.
func gradientColorPair(name string) colorPair {
	seed := fnv1.HashString64(name)
	rng := rand.New(rand.NewSource(int64(seed)))

	return colors[rng.Intn(len(colors)-1)]
}
