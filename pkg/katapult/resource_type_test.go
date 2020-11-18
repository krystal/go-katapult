package katapult

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResourceTypes(t *testing.T) {
	tests := []struct {
		name         string
		resourceType ResourceType
		value        string
	}{
		{
			name:         "TagsResourceType",
			resourceType: TagsResourceType,
			value:        "tags",
		},
		{
			name:         "VirtualMachineGroupsResourceType",
			resourceType: VirtualMachineGroupsResourceType,
			value:        "virtual_machine_groups",
		},
		{
			name:         "VirtualMachinesResourceType",
			resourceType: VirtualMachinesResourceType,
			value:        "virtual_machines",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.value, string(tt.resourceType))
		})
	}
}

func TestResourceType_objectType(t *testing.T) {
	tests := []struct {
		name     string
		value    ResourceType
		expected string
	}{
		{
			name:     "tags",
			value:    ResourceType("tags"),
			expected: "Tag",
		},
		{
			name:     "virtual machine groups",
			value:    ResourceType("virtual_machine_groups"),
			expected: "VirtualMachineGroup",
		},
		{
			name:     "virtual machine groups",
			value:    ResourceType("virtual_machines"),
			expected: "VirtualMachine",
		},
		{
			name:     "unknown type",
			value:    ResourceType("something_nope_what"),
			expected: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.value.objectType()

			assert.Equal(t, tt.expected, got)
		})
	}
}
