package buildspec

import "testing"

func TestCloudInit_Marshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *CloudInit
	}{
		{
			name: "empty",
			obj:  &CloudInit{},
		},
		{
			name: "full",
			obj: &CloudInit{
				UserData: "#cloud-config\n" +
					"runcmd:\n" +
					"  - echo 'one & two < three'\n",
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

func TestVirtualMachineSpec_CloudInitXMLBlankOmission(t *testing.T) {
	tests := []struct {
		name string
		obj  *VirtualMachineSpec
	}{
		{
			name: "nil",
			obj:  &VirtualMachineSpec{},
		},
		{
			name: "empty",
			obj:  &VirtualMachineSpec{CloudInit: &CloudInit{}},
		},
		{
			name: "whitespace",
			obj: &VirtualMachineSpec{
				CloudInit: &CloudInit{UserData: " \n\t"},
			},
		},
	}
	for _, tt := range tests {
		t.Run("xml_"+tt.name, func(t *testing.T) {
			testCustomXMLMarshaling(t, tt.obj, &VirtualMachineSpec{})
		})
	}
}
