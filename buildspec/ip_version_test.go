package buildspec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIPVersions(t *testing.T) {
	tests := []struct {
		name      string
		ipVersion IPVersion
		value     int
	}{
		{
			name:      "IPv4",
			ipVersion: IPv4,
			value:     4,
		},
		{
			name:      "IPv6",
			ipVersion: IPv6,
			value:     6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.value, int(tt.ipVersion))
		})
	}
}
