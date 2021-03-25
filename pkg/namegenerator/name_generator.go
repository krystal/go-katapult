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
			"blue",
			"green",
			"grey",
			"orange",
			"pink",
			"purple",
			"rainbow",
			"red",
			"turquoise",
			"yellow",
		},
		"appearance": {
			"bald",
			"beautiful",
			"clean",
			"elegant",
			"fancy",
			"handsome",
			"magnificent",
			"shapely",
		},
		"personality": {
			"brave",
			"calm",
			"clumsy",
			"faithful",
			"fierce",
			"happy",
			"harmless",
			"intimidating",
			"itchy",
			"jealous",
			"jolly",
			"kind",
			"nervous",
			"nice",
			"obedient",
			"proud",
			"sad",
			"scary",
			"thoughtless",
			"worried",
		},
		"conditions": {
			"easy",
			"famous",
			"helpful",
			"liked",
			"popular",
			"powerful",
			"tender",
		},
		"sounds": {
			"loud",
			"noisy",
			"quiet",
			"silent",
		},
		"tastes": {
			"bitter",
			"cool",
			"sour",
			"sweet",
			"yummy",
		},
	}
	DefaultNouns = map[string][]string{
		"animals": {
			"bear",
			"bee",
			"bobcat",
			"cat",
			"caterpillar",
			"cheetah",
			"dog",
			"dolphin",
			"elephant",
			"giraffe",
			"goat",
			"jellyfish",
			"leopard",
			"lion",
			"meerkat",
			"mouse",
			"panda",
			"sheep",
			"tiger",
			"wasp",
			"wolf",
			"zebra",
		},
		"fruits": {
			"apple",
			"apricot",
			"avocado",
			"banana",
			"cherry",
			"cucumber",
			"fig",
			"grapefruit",
			"kumquat",
			"lemon",
			"lychee",
			"mango",
			"peach",
			"pepper",
			"pineapple",
			"plum",
			"potato",
			"prune",
			"raspberry",
			"strawberry",
			"tomato",
			"watermelon",
		},
		"birds": {
			"albatross",
			"blackbird",
			"chicken",
			"chiffchaff",
			"crow",
			"curlew",
			"duck",
			"eagle",
			"goose",
			"kiwi",
			"osprey",
			"owl",
			"penguin",
			"pigeon",
			"robin",
			"sparrow",
			"swan",
			"woodpecker",
			"wren",
		},
		"flowers": {
			"bluebell",
			"carnation",
			"crocus",
			"daffodil",
			"daisy",
			"dandelion",
			"lily",
			"orchid",
			"pansy",
			"poppy",
			"rose",
			"sunflower",
			"tulip",
		},
		"transport": {
			"balloon",
			"bicycle",
			"boat",
			"bus",
			"car",
			"coach",
			"helicopter",
			"jet",
			"lawnmower",
			"lifeboat",
			"lorry",
			"plane",
			"rickshaw",
			"ship",
			"skateboard",
			"tank",
			"taxi",
			"tractor",
			"train",
			"truck",
			"van",
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
