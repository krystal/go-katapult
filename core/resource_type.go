package core

type ResourceType string

const (
	TagsResourceType                 ResourceType = "tags"
	VirtualMachineGroupsResourceType ResourceType = "virtual_machine_groups"
	VirtualMachinesResourceType      ResourceType = "virtual_machines"
)

func (s ResourceType) objectType() string {
	switch s {
	case TagsResourceType:
		return "Tag"
	case VirtualMachineGroupsResourceType:
		return "VirtualMachineGroup"
	case VirtualMachinesResourceType:
		return "VirtualMachine"
	default:
		return ""
	}
}
