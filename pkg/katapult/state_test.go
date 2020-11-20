package katapult

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStates(t *testing.T) {
	tests := []struct {
		name         string
		resourceType State
		value        string
	}{
		{
			name:         "DraftState",
			resourceType: DraftState,
			value:        "draft",
		},
		{
			name:         "FailedState",
			resourceType: FailedState,
			value:        "failed",
		},
		{
			name:         "PendingState",
			resourceType: PendingState,
			value:        "pending",
		},
		{
			name:         "CompleteState",
			resourceType: CompleteState,
			value:        "complete",
		},
		{
			name:         "BuildingState",
			resourceType: BuildingState,
			value:        "building",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.value, string(tt.resourceType))
		})
	}
}
