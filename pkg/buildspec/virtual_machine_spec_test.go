package buildspec

import (
	"bytes"
	"encoding/xml"
	"errors"
	"testing"

	"github.com/jimeh/undent"
	"github.com/krystal/go-katapult/internal/golden"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var fixtureVirtualMachineSpecBasicStruct = &VirtualMachineSpec{
	Hostname:   "web-3",
	DataCenter: &DataCenter{ID: "dc_0KVdXStXduYtcypG"},
	Resources: &Resources{
		Package: &Package{Permalink: "rock-3"},
	},
	DiskTemplate: &DiskTemplate{
		Permalink: "templates/ubuntu-18-04",
	},
}

func TestVirtualMachineSpec_Marshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *VirtualMachineSpec
	}{
		{
			name: "empty",
			obj:  &VirtualMachineSpec{},
		},
		{
			name: "full",
			obj: &VirtualMachineSpec{
				Zone: &Zone{
					ID: "zone_xmVotL1zwMwo2eXf",
				},
				DataCenter: &DataCenter{
					ID: "dc_0KVdXStXduYtcypG",
				},
				Resources: &Resources{
					Package:  &Package{ID: "vmpkg_m7mV5O0MafbDFp2n"},
					Memory:   16,
					CPUCores: 4,
				},
				DiskTemplate: &DiskTemplate{
					ID:      "dtpl_rlinMl51Lb1uvTez",
					Version: 4,
					Options: []*DiskTemplateOption{
						{Key: "foo", Value: "bar"},
						{Key: "hello", Value: "world"},
					},
				},
				SystemDisks: []*SystemDisk{
					{
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
								Schedule: &Schedule{
									Interval:  ScheduledDaily,
									Frequency: 1,
									Time:      13,
								},
							},
							{
								Retention: 30,
							},
						},
					},
					{
						Name:  "Another Disk",
						Size:  22,
						Speed: "nvme",
					},
				},
				SharedDisks: []*SharedDisk{
					{ID: "disk_gJRNxe3h7zi0Hdh5"},
					{Name: "image-uploads"},
				},
				NetworkInterfaces: []*NetworkInterface{
					{
						Network: &Network{ID: "netw_DRIS3BaTWfKaHlWW"},
						SpeedProfile: &NetworkSpeedProfile{
							ID: "nsp_eHwC5NG3DRAHzVfD",
						},
					},
					{
						Network: &Network{ID: "netw_17w3MepxvWE4J3Zx"},
						SpeedProfile: &NetworkSpeedProfile{
							ID: "nsp_bFQhDNAluyp4t2A9",
						},
						IPAddressAllocations: []*IPAddressAllocation{
							{
								Type:    NewIPAddressAllocation,
								Version: IPv4,
							},
							{
								Type:    NewIPAddressAllocation,
								Version: IPv6,
							},
							{
								Type:    NewIPAddressAllocation,
								Version: IPv4,
								Subnet:  &Subnet{ID: "sbnt_xxhvuhr3dsvEHcM5"},
							},
							{
								Type:    NewIPAddressAllocation,
								Version: IPv6,
								Subnet:  &Subnet{ID: "sbnt_Pms921K2pYf35nae"},
							},
							{
								Type: ExistingIPAddressAllocation,
								IPAddress: &IPAddress{
									ID: "ip_Hb8WpvV9qRMznHwZ",
								},
							},
						},
					},
					{
						VirtualNetwork: &VirtualNetwork{
							ID: "vnet_Cuc45YcBaUhWqx6u",
						},
					},
				},
				Hostname:    "bitter-beautiful-mango",
				Name:        "web-1",
				Description: "Web Server #1",
				Group:       &Group{ID: "vmgrp_dZDXXLw7e54Ep6CG"},
				AuthorizedKeys: &AuthorizedKeys{
					AllSSHKeys: true,
					Users: []*User{
						{ID: "user_yUfYcKHgU1ywBWzP"},
						{EmailAddress: "jane@doe.com"},
					},
				},
				BackupPolicies: []*BackupPolicy{
					{
						Retention: 24,
						Schedule: &Schedule{
							Interval:  ScheduledWeekly,
							Frequency: 1,
							Time:      13,
						},
					},
					{
						Retention: 30,
					},
				},
				Tags: []string{"ha", "db", "web"},
				ISO:  "iso_R6hPTR62bTSj5hQe",
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

func TestVirtualMachineSpec_UnmarshalXML_Name(t *testing.T) {
	tests := []struct {
		name string
		xml  string
		want *VirtualMachineSpec
	}{
		{
			name: "not nested",
			xml: undent.String(`
				<VirtualMachineSpec>
					<Name>database-2</Name>
				</VirtualMachineSpec>`,
			),
			want: &VirtualMachineSpec{Name: "database-2"},
		},
		{
			name: "nested",
			xml: undent.String(`
				<VirtualMachineSpec>
					<Name>
						<Name>database-2</Name>
					</Name>
				</VirtualMachineSpec>`,
			),
			want: &VirtualMachineSpec{Name: "database-2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &VirtualMachineSpec{}

			err := xml.Unmarshal([]byte(tt.xml), v)
			require.NoError(t, err)

			assert.Equal(t, "database-2", v.Name)
		})
	}
}

func TestVirtualMachineSpec_ToFromJSON(t *testing.T) {
	tests := []struct {
		name string
		json string
		spec *VirtualMachineSpec
	}{
		{
			name: "empty spec",
			spec: &VirtualMachineSpec{},
		},
		{
			name: "basic spec",
			spec: fixtureVirtualMachineSpecBasicStruct,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			marshaled, err := tt.spec.JSON()
			require.NoError(t, err, "failed marshaling with JSON()")

			buf := &bytes.Buffer{}
			err = tt.spec.WriteJSON(buf)
			require.NoError(t, err, "failed marshaling with WriteJSON()")

			assert.Equal(t, string(marshaled), buf.String(),
				"output of WriteJSON() does not match that of JSON()",
			)

			if golden.Update() {
				golden.Set(t, marshaled)
			}

			g := golden.Get(t)
			assert.Equal(t, string(g), string(marshaled),
				"json encoded value does not match golden",
			)

			r := bytes.NewReader(g)
			got, err := FromJSON(r)
			require.NoError(t, err, "json decoding golden failed")
			assert.Equal(t, tt.spec, got,
				"json decoding from golden does not match expected object",
			)
		})
	}
}

func TestVirtualMachineSpec_WriteJSON_Error(t *testing.T) {
	errStr := "failed to write" //nolint:goconst
	spec := &VirtualMachineSpec{Name: "web-3"}
	w := &badWriter{err: errors.New(errStr)}

	err := spec.WriteJSON(w)

	assert.EqualError(t, err, errStr)
}

func TestVirtualMachineSpec_ToFromXML(t *testing.T) {
	tests := []struct {
		name string
		xml  string
		spec *VirtualMachineSpec
	}{
		{
			name: "empty spec",
			spec: &VirtualMachineSpec{},
		},
		{
			name: "basic spec",
			spec: fixtureVirtualMachineSpecBasicStruct,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			marshaled, err := tt.spec.XML()
			require.NoError(t, err, "failed marshaling with XML()")

			buf := &bytes.Buffer{}
			err = tt.spec.WriteXML(buf)
			require.NoError(t, err, "failed marshaling with WriteXML()")

			assert.Equal(t, string(marshaled), buf.String(),
				"output of WriteXML() does not match that of XML()",
			)

			if golden.Update() {
				golden.Set(t, marshaled)
			}

			g := golden.Get(t)
			assert.Equal(t, string(g), string(marshaled),
				"xml encoded value does not match golden",
			)

			r := bytes.NewReader(g)
			got, err := FromXML(r)
			require.NoError(t, err, "xml decoding golden failed")
			assert.Equal(t, tt.spec, got,
				"xml decoding from golden does not match expected object",
			)
		})
	}
}

func TestVirtualMachineSpec_WriteXML_Error(t *testing.T) {
	errStr := "failed to write" //nolint:goconst
	spec := &VirtualMachineSpec{Name: "web-3"}
	w := &badWriter{err: errors.New(errStr)}

	err := spec.WriteXML(w)

	assert.EqualError(t, err, errStr)
}

func TestVirtualMachineSpec_ToFromYAML(t *testing.T) {
	tests := []struct {
		name string
		yaml string
		spec *VirtualMachineSpec
	}{
		{
			name: "empty spec",
			spec: &VirtualMachineSpec{},
		},
		{
			name: "basic spec",
			spec: fixtureVirtualMachineSpecBasicStruct,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			marshaled, err := tt.spec.YAML()
			require.NoError(t, err, "failed marshaling with YAML()")

			buf := &bytes.Buffer{}
			err = tt.spec.WriteYAML(buf)
			require.NoError(t, err, "failed marshaling with WriteYAML()")

			assert.Equal(t, string(marshaled), buf.String(),
				"output of WriteYAML() does not match that of YAML()",
			)

			if golden.Update() {
				golden.Set(t, marshaled)
			}

			g := golden.Get(t)
			assert.Equal(t, string(g), string(marshaled),
				"yaml encoded value does not match golden",
			)

			r := bytes.NewReader(g)
			got, err := FromYAML(r)
			require.NoError(t, err, "yaml decoding golden failed")
			assert.Equal(t, tt.spec, got,
				"yaml decoding from golden does not match expected object",
			)
		})
	}
}

func TestVirtualMachineSpec_WriteYAML_Error(t *testing.T) {
	errStr := "failed to write" //nolint:goconst
	spec := &VirtualMachineSpec{Name: "web-3"}
	w := &badWriter{err: errors.New(errStr)}

	err := spec.WriteYAML(w)

	assert.EqualError(t, err, "yaml: write error: "+errStr)
}
