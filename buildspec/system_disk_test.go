package buildspec

import "testing"

func TestSystemDisk_Marshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *SystemDisk
	}{
		{
			name: "empty",
			obj:  &SystemDisk{},
		},
		{
			name: "full",
			obj: &SystemDisk{
				Name:  "System Disk",
				Size:  10,
				Speed: "ssd",
				IOProfile: &DiskIOProfile{
					ID: "diop_xPlNw7iDmrGOnPRA",
				},
				FileSystemType: "ext4",
				BackupPolicies: []*BackupPolicy{
					{
						Retention: 24,
						Schedule:  &Schedule{Interval: ScheduledDaily},
					},
					{
						Retention: 13,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run("json_"+tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
		t.Run("xml_"+tt.name, func(t *testing.T) {
			testXMLMarshaling(t, tt.obj)
		})
		t.Run("yaml_"+tt.name, func(t *testing.T) {
			testYAMLMarshaling(t, tt.obj)
		})
	}
}

func Test_xmlSystemDisks_Marshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *xmlSystemDisks
	}{
		{
			name: "empty",
			obj:  &xmlSystemDisks{},
		},
		{
			name: "full",
			obj: &xmlSystemDisks{
				SystemDisks: []*SystemDisk{
					{
						Name:  "System Disk",
						Size:  10,
						Speed: "ssd",
					},
					{
						Name:  "Another Disk",
						Size:  30,
						Speed: "nvme",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run("xml_"+tt.name, func(t *testing.T) {
			testXMLMarshaling(t, tt.obj)
		})
	}
}
