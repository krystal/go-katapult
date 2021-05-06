package namegenerator

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomHostname(t *testing.T) {
	for i := 0; i < 10; i++ {
		got := RandomHostname()

		assert.Equalf(t, 2, strings.Count(got, "-"),
			"generated hostname does not contain two hyphens: %s", got,
		)
		assert.Falsef(t, strings.ContainsAny(got, "0123456789"),
			"generated hostname contains numbers: %s", got,
		)
	}
}

func TestRandomName(t *testing.T) {
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
			got := RandomName(tt.args.prefixes...)

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
