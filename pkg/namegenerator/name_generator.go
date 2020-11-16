// Package namegenerator provides various helper methods to generate randomized
// names for use on the Katapult platform.
package namegenerator

import (
	"math/rand"
	"strings"
	"time"
)

var (
	DefaultAdjectives = map[string][]string{
		"colors": {
			"green", "red", "yellow", "blue", "grey", "purple", "orange",
			"pink", "rainbow", "turquoise",
		},
		"appearance": {
			"bald", "beautiful", "clean", "elegant", "fancy", "magnificent",
			"shapely", "handsome",
		},
		"personality": {
			"harmless", "happy", "sad", "brave", "calm", "faithful", "jolly",
			"kind", "proud", "nice", "obedient", "clumsy", "fierce", "itchy",
			"jealous", "nervous", "scary", "thoughtless", "worried",
			"intimidating",
		},
		"conditions": {
			"famous", "easy", "powerful", "tender", "helpful", "popular",
			"liked",
		},
		"sounds": {"loud", "quiet", "noisy", "silent"},
		"tastes": {"cool", "bitter", "yummy", "sweet", "sour"},
	}
	DefaultNouns = map[string][]string{
		"animals": {
			"cat", "mouse", "dog", "elephant", "lion", "giraffe", "sheep",
			"goat", "caterpillar", "bee", "wasp", "wolf", "bobcat", "cheetah",
			"panda", "bear", "jellyfish", "leopard", "meerkat", "tiger",
			"dolphin", "zebra",
		},
		"fruits": {
			"potato", "tomato", "banana", "raspberry", "apple", "kumquat",
			"mango", "lychee", "strawberry", "cucumber", "pepper", "raspberry",
			"apricot", "avocado", "cherry", "grapefruit", "fig", "lemon",
			"peach", "pineapple", "plum", "prune", "watermelon",
		},
		"birds": {
			"eagle", "albatross", "blackbird", "crow", "curlew", "woodpecker",
			"kiwi", "pigeon", "owl", "penguin", "blackbird", "wren", "sparrow",
			"osprey", "chiffchaff", "swan", "goose", "duck", "chicken", "crow",
			"robin",
		},
		"flowers": {
			"daisy", "dandelion", "sunflower", "lily", "bluebell", "carnation",
			"crocus", "daffodil", "orchid", "pansy", "poppy", "rose", "tulip",
		},
		"transport": {
			"car", "bus", "train", "plane", "tank", "skateboard", "bicycle",
			"boat", "ship", "balloon", "lawnmower", "tractor", "taxi",
			"rickshaw", "helicopter", "jet", "lifeboat", "van", "truck",
			"lorry", "coach",
		},
	}
)

// NameGenerator is responsible for generating various types of randomized
// names.
type NameGenerator struct {
	rand       *rand.Rand
	adjectives *wordList
	nouns      *wordList
}

func New(
	adjectives map[string][]string,
	nouns map[string][]string,
) *NameGenerator {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source) //nolint:gosec

	return &NameGenerator{
		rand:       r,
		adjectives: newWordList(r, adjectives),
		nouns:      newWordList(r, nouns),
	}
}

func (s *NameGenerator) RandomHostname() string {
	a1, group := s.adjectives.getRandom()
	a2, _ := s.adjectives.getRandom(group)
	noun, _ := s.nouns.getRandom()

	return strings.Join([]string{a1, a2, noun}, "-")
}

func (s *NameGenerator) RandomName(prefixes ...string) string {
	adj, _ := s.adjectives.getRandom()
	noun, _ := s.nouns.getRandom()
	parts := append(prefixes, adj, noun)

	return strings.Join(parts, "-")
}

type wordList struct {
	words []*word
	rand  *rand.Rand
}

func newWordList(rand *rand.Rand, wordGroups map[string][]string) *wordList {
	wl := &wordList{rand: rand}

	for group, words := range wordGroups {
		for _, w := range words {
			wl.words = append(wl.words, &word{w, group})
		}
	}

	return wl
}

func (s *wordList) getRandom(excludedGroups ...string) (string, string) {
	var pool []*word
	for _, w := range s.words {
		if !w.belongsToOneOf(excludedGroups...) {
			pool = append(pool, w)
		}
	}
	pick := pool[s.rand.Intn(len(pool))]

	return pick.value, pick.group
}

type word struct {
	value string
	group string
}

func (w *word) belongsToOneOf(groups ...string) bool {
	for _, g := range groups {
		if g == w.group {
			return true
		}
	}

	return false
}
