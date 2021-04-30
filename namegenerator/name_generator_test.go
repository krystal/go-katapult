package namegenerator

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultListDuplicates(t *testing.T) {
	defaults := []map[string][]string{DefaultAdjectives, DefaultNouns}
	counts := map[string]int{}

	for _, cat := range defaults {
		for _, words := range cat {
			for _, word := range words {
				if _, ok := counts[word]; !ok {
					counts[word] = 1
				} else {
					counts[word]++
				}
			}
		}
	}

	dupes := []string{}
	for w, n := range counts {
		if n > 1 {
			dupes = append(dupes, w)
		}
	}

	assert.Empty(t, dupes, "words should not be repeated")
}

func TestNameGenerator_RandomHostname(t *testing.T) {
	g := New(DefaultAdjectives, DefaultNouns)

	got := g.RandomHostname()

	assert.Equalf(t, 2, strings.Count(got, "-"),
		"generated hostname does not contain two hyphens: %s", got,
	)
	assert.Falsef(t, strings.ContainsAny(got, "0123456789"),
		"generated hostname contains numbers: %s", got,
	)
}

func TestNameGenerator_RandomName(t *testing.T) {
	type args struct {
		prefixes []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "no prefixes",
			args: args{},
		},
		{
			name: "one prefix",
			args: args{prefixes: []string{"prod"}},
		},
		{
			name: "two prefixes",
			args: args{prefixes: []string{"tf", "test"}},
		},
		{
			name: "three prefixes",
			args: args{prefixes: []string{"tf", "test", "acc"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := New(DefaultAdjectives, DefaultNouns)

			got := g.RandomName(tt.args.prefixes...)

			hyphenCount := 1
			if len(tt.args.prefixes) > 0 {
				hyphenCount += len(tt.args.prefixes)

				prefix := strings.Join(tt.args.prefixes, "-") + "-"
				assert.Truef(t, strings.HasPrefix(got, prefix),
					"generated name does not start with \"%s\"", prefix,
				)
			}

			assert.Equalf(t, hyphenCount, strings.Count(got, "-"),
				"generated name does not contain %d hyphens: %s",
				hyphenCount, got,
			)
			assert.Falsef(t, strings.ContainsAny(got, "0123456789"),
				"generated name contains numbers: %s", got,
			)
		})
	}
}

func BenchmarkNameGenerator_RandomHostname(b *testing.B) {
	g := New(DefaultAdjectives, DefaultNouns)

	for n := 0; n < b.N; n++ {
		g.RandomHostname()
	}
}

func BenchmarkNameGenerator_RandomName_NoPrefix(b *testing.B) {
	g := New(DefaultAdjectives, DefaultNouns)

	for n := 0; n < b.N; n++ {
		g.RandomName()
	}
}

func BenchmarkNameGenerator_RandomName_OnePrefix(b *testing.B) {
	g := New(DefaultAdjectives, DefaultNouns)

	for n := 0; n < b.N; n++ {
		g.RandomName("tf")
	}
}

func BenchmarkNameGenerator_RandomName_TwoPrefixes(b *testing.B) {
	g := New(DefaultAdjectives, DefaultNouns)

	for n := 0; n < b.N; n++ {
		g.RandomName("tf", "test")
	}
}

func BenchmarkNameGenerator_RandomName_ThreePrefixes(b *testing.B) {
	g := New(DefaultAdjectives, DefaultNouns)

	for n := 0; n < b.N; n++ {
		g.RandomName("tf", "test", "acc")
	}
}
