package buildspec

import (
	"testing"
)

func TestBackupPolicy_Marshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *BackupPolicy
	}{
		{
			name: "empty",
			obj:  &BackupPolicy{},
		},
		{
			name: "full",
			obj: &BackupPolicy{
				Retention: 24,
				Schedule:  &Schedule{Interval: ScheduledDaily},
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

func Test_xmlBackupPolicies_Marshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *xmlBackupPolicies
	}{
		{
			name: "empty",
			obj:  &xmlBackupPolicies{},
		},
		{
			name: "full",
			obj: &xmlBackupPolicies{
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
		t.Run("xml_"+tt.name, func(t *testing.T) {
			testXMLMarshaling(t, tt.obj)
		})
	}
}
