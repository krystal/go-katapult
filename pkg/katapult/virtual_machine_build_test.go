package katapult

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	fixtureBuildNotFoundErr = "build_not_found: No build was found matching " +
		"any of the criteria provided in the arguments"
	fixtureBuildNotFoundResponseError = &ResponseError{
		Code: "build_not_found",
		Description: "No build was found matching any of the criteria " +
			"provided in the arguments",
		Detail: json.RawMessage(`{}`),
	}
)

func TestVirtualMachineBuild_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *VirtualMachineBuild
	}{
		{
			name: "empty",
			obj:  &VirtualMachineBuild{},
		},
		{
			name: "full",
			obj: &VirtualMachineBuild{
				ID:             "id1",
				SpecXML:        "<xml/>",
				State:          VirtualMachineBuildDraft,
				VirtualMachine: &VirtualMachine{ID: "id2"},
				CreatedAt:      timestampPtr(1600192008),
			},
		},
		{
			name: "Draft",
			obj: &VirtualMachineBuild{
				ID:    "id1",
				State: VirtualMachineBuildDraft,
			},
		},
		{
			name: "failed",
			obj: &VirtualMachineBuild{
				ID:    "id1",
				State: VirtualMachineBuildFailed,
			},
		},
		{
			name: "pending",
			obj: &VirtualMachineBuild{
				ID:    "id1",
				State: VirtualMachineBuildPending,
			},
		},
		{
			name: "building",
			obj: &VirtualMachineBuild{
				ID:    "id1",
				State: VirtualMachineBuildBuilding,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestVirtualMachineBuildStates(t *testing.T) {
	tests := []struct {
		name  string
		enum  VirtualMachineBuildState
		value string
	}{
		{
			name:  "VirtualMachineBuildDraft",
			enum:  VirtualMachineBuildDraft,
			value: "draft",
		},
		{
			name:  "VirtualMachineBuildFailed",
			enum:  VirtualMachineBuildFailed,
			value: "failed",
		},
		{
			name:  "VirtualMachineBuildPending",
			enum:  VirtualMachineBuildPending,
			value: "pending",
		},
		{
			name:  "VirtualMachineBuildComplete",
			enum:  VirtualMachineBuildComplete,
			value: "complete",
		},
		{
			name:  "VirtualMachineBuildBuilding",
			enum:  VirtualMachineBuildBuilding,
			value: "building",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.value, string(tt.enum))
		})
	}
}

func Test_virtualMachineBuildResponseBody_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *virtualMachineBuildsResponseBody
	}{
		{
			name: "empty",
			obj:  &virtualMachineBuildsResponseBody{},
		},
		{
			name: "full",
			obj: &virtualMachineBuildsResponseBody{
				Task:                &Task{ID: "id1"},
				Build:               &VirtualMachineBuild{ID: "id2"},
				VirtualMachineBuild: &VirtualMachineBuild{ID: "id3"},
				Hostname:            "host.example.com",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestVirtualMachineBuildArguments_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *VirtualMachineBuildArguments
	}{
		{
			name: "empty",
			obj:  &VirtualMachineBuildArguments{},
		},
		{
			name: "full",
			obj: &VirtualMachineBuildArguments{
				Organization: &Organization{ID: "org_O648YDMEYeLmqdmn"},
				Zone:         &Zone{ID: "zone_kY2sPRG24sJVRM2U"},
				DataCenter:   &DataCenter{ID: "loc_25d48761871e4bf"},
				Package: &VirtualMachinePackage{
					ID: "vmpkg_XdNPhGXvyt1dnDts",
				},
				DiskTemplate: &DiskTemplate{ID: "dtpl_ytP13XD5DE1RdSL9"},
				DiskTemplateOptions: []*DiskTemplateOption{
					{Key: "foo", Value: "bar"},
				},
				Network:  &Network{ID: "netw_zDW7KYAeqqfRfVag"},
				Hostname: "foo.example.com",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestVirtualMachineBuildArguments_forRequest(t *testing.T) {
	tests := []struct {
		name     string
		obj      *VirtualMachineBuildArguments
		expected *VirtualMachineBuildArguments
	}{
		{
			name:     "nil",
			obj:      (*VirtualMachineBuildArguments)(nil),
			expected: nil,
		},
		{
			name: "full",
			obj: &VirtualMachineBuildArguments{
				Organization: &Organization{
					ID:        "org_O648YDMEYeLmqdmn",
					Name:      "ACME Inc.",
					SubDomain: "acme",
				},
				Zone: &Zone{
					ID:        "zone_kY2sPRG24sJVRM2U",
					Name:      "North West",
					Permalink: "north-west",
				},
				DataCenter: &DataCenter{
					ID:        "loc_25d48761871e4bf",
					Name:      "Woodland",
					Permalink: "woodland",
				},
				Package: &VirtualMachinePackage{
					ID:        "vmpkg_XdNPhGXvyt1dnDts",
					Name:      "X-Small",
					Permalink: "xsmall",
				},
				DiskTemplate: &DiskTemplate{
					ID:        "dtpl_ytP13XD5DE1RdSL9",
					Name:      "Ubuntu 18.04 Server",
					Permalink: "templates/ubuntu-18-04",
				},
				DiskTemplateOptions: []*DiskTemplateOption{
					{Key: "foo", Value: "bar"},
				},
				Network: &Network{
					ID:        "netw_zDW7KYAeqqfRfVag",
					Name:      "Public Network",
					Permalink: "public",
				},
				Hostname: "foo.example.com",
			},
			expected: &VirtualMachineBuildArguments{
				Organization: &Organization{
					ID: "org_O648YDMEYeLmqdmn",
				},
				Zone: &Zone{
					ID: "zone_kY2sPRG24sJVRM2U",
				},
				DataCenter: &DataCenter{
					ID: "loc_25d48761871e4bf",
				},
				Package: &VirtualMachinePackage{
					ID: "vmpkg_XdNPhGXvyt1dnDts",
				},
				DiskTemplate: &DiskTemplate{
					ID: "dtpl_ytP13XD5DE1RdSL9",
				},
				DiskTemplateOptions: []*DiskTemplateOption{
					{Key: "foo", Value: "bar"},
				},
				Network: &Network{
					ID: "netw_zDW7KYAeqqfRfVag",
				},
				Hostname: "foo.example.com",
			},
		},
		{
			name: "virtual machine build (no IDs)",
			obj: &VirtualMachineBuildArguments{
				Organization: &Organization{
					Name:      "ACME Inc.",
					SubDomain: "acme",
				},
				Zone: &Zone{
					Name:      "North West",
					Permalink: "north-west",
				},
				DataCenter: &DataCenter{
					Name:      "Woodland",
					Permalink: "woodland",
				},
				Package: &VirtualMachinePackage{
					Name:      "X-Small",
					Permalink: "xsmall",
				},
				DiskTemplate: &DiskTemplate{
					Name:      "Ubuntu 18.04 Server",
					Permalink: "templates/ubuntu-18-04",
				},
				DiskTemplateOptions: []*DiskTemplateOption{
					{Key: "foo", Value: "bar"},
				},
				Network: &Network{
					Name:      "Public Network",
					Permalink: "public",
				},
				Hostname: "foo.example.com",
			},
			expected: &VirtualMachineBuildArguments{
				Organization: &Organization{
					SubDomain: "acme",
				},
				Zone: &Zone{
					Permalink: "north-west",
				},
				DataCenter: &DataCenter{
					Permalink: "woodland",
				},
				Package: &VirtualMachinePackage{
					Permalink: "xsmall",
				},
				DiskTemplate: &DiskTemplate{
					Permalink: "templates/ubuntu-18-04",
				},
				DiskTemplateOptions: []*DiskTemplateOption{
					{Key: "foo", Value: "bar"},
				},
				Network: &Network{
					Permalink: "public",
				},
				Hostname: "foo.example.com",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var orgName string
			if tt.obj != nil && tt.obj.Organization != nil {
				orgName = tt.obj.Organization.Name
			}

			got := tt.obj.forRequest()

			assert.Equal(t, tt.expected, got)

			if orgName != "" {
				assert.Equal(t, orgName, tt.obj.Organization.Name,
					"original object was modified")
			}
		})
	}
}

func TestVirtualMachineBuildsClient_Get(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name       string
		args       args
		id         string
		expected   *VirtualMachineBuild
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "virtual machine build",
			args: args{
				ctx: context.Background(),
				id:  "vmbuild_pbjJIqJ3MOMNsCr3",
			},
			expected: &VirtualMachineBuild{
				ID:      "vmbuild_pbjJIqJ3MOMNsCr3",
				SpecXML: "<?xml version=\"1.0\"?>\n",
				State:   VirtualMachineBuildComplete,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machine_build_get"),
		},
		{
			name: "virtual machine build (alt response)",
			args: args{
				ctx: context.Background(),
				id:  "vmbuild_pbjJIqJ3MOMNsCr3",
			},
			expected: &VirtualMachineBuild{
				ID:      "vmbuild_pbjJIqJ3MOMNsCr3",
				SpecXML: "<?xml version=\"1.0\"?>\n",
				State:   VirtualMachineBuildComplete,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machine_build_get_alt"),
		},
		{
			name: "non-existent virtual machine build",
			args: args{
				ctx: context.Background(),
				id:  "vmbuild_nopethisbegone",
			},
			errStr:     fixtureBuildNotFoundErr,
			errResp:    fixtureBuildNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("build_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				id:  "vmbuild_pbjJIqJ3MOMNsCr3",
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := prepareTestClient()
			defer teardown()

			mux.HandleFunc(
				fmt.Sprintf("/core/v1/virtual_machines/builds/%s", tt.args.id),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.VirtualMachineBuilds.Get(
				tt.args.ctx, tt.args.id,
			)

			if tt.respStatus != 0 {
				assert.Equal(t, tt.respStatus, resp.StatusCode)
			}

			if tt.errStr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.errStr)
			}

			if tt.expected != nil {
				assert.Equal(t, tt.expected, got)
			}

			if tt.errResp != nil {
				assert.Equal(t, tt.errResp, resp.Error)
			}
		})
	}
}

func TestVirtualMachineBuildsClient_Create(t *testing.T) {
	fullArgs := &VirtualMachineBuildArguments{
		Organization: &Organization{
			ID:        "org_O648YDMEYeLmqdmn",
			Name:      "ACME Inc.",
			SubDomain: "acme",
		},
		Zone: &Zone{
			ID:        "zone_kY2sPRG24sJVRM2U",
			Name:      "North West",
			Permalink: "north-west",
		},
		DataCenter: &DataCenter{
			ID:        "loc_25d48761871e4bf",
			Name:      "Woodland",
			Permalink: "woodland",
		},
		Package: &VirtualMachinePackage{
			ID:        "vmpkg_XdNPhGXvyt1dnDts",
			Name:      "X-Small",
			Permalink: "xsmall",
		},
		DiskTemplate: &DiskTemplate{
			ID:        "dtpl_ytP13XD5DE1RdSL9",
			Name:      "Ubuntu 18.04 Server",
			Permalink: "templates/ubuntu-18-04",
		},
		DiskTemplateOptions: []*DiskTemplateOption{
			{Key: "foo", Value: "bar"},
		},
		Network: &Network{
			ID:        "netw_zDW7KYAeqqfRfVag",
			Name:      "Public Network",
			Permalink: "public",
		},
		Hostname: "foo.example.com",
	}
	fullReqBody := &VirtualMachineBuildArguments{
		Organization: &Organization{
			ID: "org_O648YDMEYeLmqdmn",
		},
		Zone: &Zone{
			ID: "zone_kY2sPRG24sJVRM2U",
		},
		DataCenter: &DataCenter{
			ID: "loc_25d48761871e4bf",
		},
		Package: &VirtualMachinePackage{
			ID: "vmpkg_XdNPhGXvyt1dnDts",
		},
		DiskTemplate: &DiskTemplate{
			ID: "dtpl_ytP13XD5DE1RdSL9",
		},
		DiskTemplateOptions: []*DiskTemplateOption{
			{Key: "foo", Value: "bar"},
		},
		Network: &Network{
			ID: "netw_zDW7KYAeqqfRfVag",
		},
		Hostname: "foo.example.com",
	}

	noIDArgs := &VirtualMachineBuildArguments{
		Organization: &Organization{
			Name:      "ACME Inc.",
			SubDomain: "acme",
		},
		Zone: &Zone{
			Name:      "North West",
			Permalink: "north-west",
		},
		DataCenter: &DataCenter{
			Name:      "Woodland",
			Permalink: "woodland",
		},
		Package: &VirtualMachinePackage{
			Name:      "X-Small",
			Permalink: "xsmall",
		},
		DiskTemplate: &DiskTemplate{
			Name:      "Ubuntu 18.04 Server",
			Permalink: "templates/ubuntu-18-04",
		},
		DiskTemplateOptions: []*DiskTemplateOption{
			{Key: "foo", Value: "bar"},
		},
		Network: &Network{
			Name:      "Public Network",
			Permalink: "public",
		},
		Hostname: "foo.example.com",
	}
	noIDReqBody := &VirtualMachineBuildArguments{
		Organization: &Organization{
			SubDomain: "acme",
		},
		Zone: &Zone{
			Permalink: "north-west",
		},
		DataCenter: &DataCenter{
			Permalink: "woodland",
		},
		Package: &VirtualMachinePackage{
			Permalink: "xsmall",
		},
		DiskTemplate: &DiskTemplate{
			Permalink: "templates/ubuntu-18-04",
		},
		DiskTemplateOptions: []*DiskTemplateOption{
			{Key: "foo", Value: "bar"},
		},
		Network: &Network{
			Permalink: "public",
		},
		Hostname: "foo.example.com",
	}

	type args struct {
		ctx       context.Context
		buildArgs *VirtualMachineBuildArguments
	}
	tests := []struct {
		name            string
		args            args
		expectedReqBody *VirtualMachineBuildArguments
		expected        *VirtualMachineBuild
		errStr          string
		errResp         *ResponseError
		respStatus      int
		respBody        []byte
	}{
		{
			name: "virtual machine build",
			args: args{
				ctx:       context.Background(),
				buildArgs: fullArgs,
			},
			expectedReqBody: fullReqBody,
			expected: &VirtualMachineBuild{
				ID:    "vmbuild_TEmhezUShNuAsyac",
				State: VirtualMachineBuildPending,
			},
			respStatus: http.StatusCreated,
			respBody:   fixture("virtual_machine_build_create"),
		},
		{
			name: "virtual machine build (no IDs)",
			args: args{
				ctx:       context.Background(),
				buildArgs: noIDArgs,
			},
			expectedReqBody: noIDReqBody,
			expected: &VirtualMachineBuild{
				ID:    "vmbuild_TEmhezUShNuAsyac",
				State: VirtualMachineBuildPending,
			},
			respStatus: http.StatusCreated,
			respBody:   fixture("virtual_machine_build_create"),
		},
		{
			name: "invalid API token response",
			args: args{
				ctx:       context.Background(),
				buildArgs: fullArgs,
			},
			errStr:     fixtureInvalidAPITokenErr,
			errResp:    fixtureInvalidAPITokenResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("invalid_api_token_error"),
		},
		{
			name: "non-existent organization",
			args: args{
				ctx:       context.Background(),
				buildArgs: fullArgs,
			},
			errStr:     fixtureOrganizationNotFoundErr,
			errResp:    fixtureOrganizationNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("organization_not_found_error"),
		},
		{
			name: "suspended organization",
			args: args{
				ctx:       context.Background(),
				buildArgs: fullArgs,
			},
			errStr:     fixtureOrganizationSuspendedErr,
			errResp:    fixtureOrganizationSuspendedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("organization_suspended_error"),
		},
		{
			name: "non-existent data center",
			args: args{
				ctx:       context.Background(),
				buildArgs: fullArgs,
			},
			errStr:     fixtureDataCenterNotFoundErr,
			errResp:    fixtureDataCenterNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("data_center_not_found_error"),
		},
		{
			name: "non-existent virtual machine package",
			args: args{
				ctx:       context.Background(),
				buildArgs: fullArgs,
			},
			errStr:     fixturePackageNotFoundErr,
			errResp:    fixturePackageNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("package_not_found_error"),
		},
		{
			name: "non-existent disk template",
			args: args{
				ctx:       context.Background(),
				buildArgs: fullArgs,
			},
			errStr:     fixtureDiskTemplateNotFoundErr,
			errResp:    fixtureDiskTemplateNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("disk_template_not_found_error"),
		},
		{
			name: "non-existent zone",
			args: args{
				ctx:       context.Background(),
				buildArgs: fullArgs,
			},
			errStr:     fixtureZoneNotFoundErr,
			errResp:    fixtureZoneNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("zone_not_found_error"),
		},
		{
			name: "permission denied",
			args: args{
				ctx:       context.Background(),
				buildArgs: fullArgs,
			},
			errStr:     fixturePermissionDeniedErr,
			errResp:    fixturePermissionDeniedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("permission_denied_error"),
		},
		{
			name: "validation error",
			args: args{
				ctx:       context.Background(),
				buildArgs: fullArgs,
			},
			errStr:     fixtureValidationErrorErr,
			errResp:    fixtureValidationErrorResponseError,
			respStatus: http.StatusUnprocessableEntity,
			respBody:   fixture("validation_error"),
		},
		{
			name: "location_required error",
			args: args{
				ctx:       context.Background(),
				buildArgs: fullArgs,
			},
			errStr: "location_required: A zone or a data_center argument " +
				"must be provided",
			errResp: &ResponseError{
				Code: "location_required",
				Description: "A zone or a data_center argument must be " +
					"provided",
				Detail: json.RawMessage(`{}`),
			},
			respStatus: http.StatusUnprocessableEntity,
			respBody:   fixture("location_required_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx:       nil,
				buildArgs: fullArgs,
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := prepareTestClient()
			defer teardown()

			mux.HandleFunc(
				"/core/v1/organizations/_/virtual_machines/build",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "POST", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					if tt.expectedReqBody != nil {
						reqBody := &VirtualMachineBuildArguments{}
						err := strictUmarshal(r.Body, reqBody)
						assert.NoError(t, err)
						assert.Equal(t, tt.expectedReqBody, reqBody)
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.VirtualMachineBuilds.Create(
				tt.args.ctx, tt.args.buildArgs,
			)

			if tt.respStatus != 0 {
				assert.Equal(t, tt.respStatus, resp.StatusCode)
			}

			if tt.errStr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.errStr)
			}

			if tt.expected != nil {
				assert.Equal(t, tt.expected, got)
			}

			if tt.errResp != nil {
				assert.Equal(t, tt.errResp, resp.Error)
			}
		})
	}
}
