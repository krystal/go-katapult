package core

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/krystal/go-katapult"
	"github.com/krystal/go-katapult/buildspec"
	"github.com/krystal/go-katapult/internal/test"
	"net/http"
	"testing"

	"github.com/jimeh/undent"
	"github.com/stretchr/testify/assert"
)

var (
	fixtureBuildNotFoundErr = "build_not_found: No build was found matching " +
		"any of the criteria provided in the arguments"
	fixtureBuildNotFoundResponseError = &katapult.ResponseError{
		Code: "build_not_found",
		Description: "No build was found matching any of the criteria " +
			"provided in the arguments",
		Detail: json.RawMessage(`{}`),
	}

	fixtureInvalidXMLSpecErr = "invalid_spec_xml: The spec XML provided is " +
		"invalid -- " + undent.String(`
            {
              "errors": "1:21: FATAL: EndTag: '</' not found"
            }`,
	)
	fixtureInvalidXMLSpecResponseError = &katapult.ResponseError{
		Code:        "invalid_spec_xml",
		Description: "The spec XML provided is invalid",
		Detail: json.RawMessage(`{
      "errors": "1:21: FATAL: EndTag: '</' not found"
    }`),
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

func Test_virtualMachineBuildCreateRequest_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *virtualMachineBuildCreateRequest
	}{
		{
			name: "empty",
			obj:  &virtualMachineBuildCreateRequest{},
		},
		{
			name: "full",
			obj: &virtualMachineBuildCreateRequest{
				Hostname:     "foo.example.com",
				Organization: &Organization{ID: "org_O648YDMEYeLmqdmn"},
				Zone:         &Zone{ID: "zone_kY2sPRG24sJVRM2U"},
				DataCenter:   &DataCenter{ID: "dc_25d48761871e4bf"},
				Package: &VirtualMachinePackage{
					ID: "vmpkg_XdNPhGXvyt1dnDts",
				},
				DiskTemplate: &DiskTemplate{ID: "dtpl_ytP13XD5DE1RdSL9"},
				DiskTemplateOptions: []*DiskTemplateOption{
					{Key: "foo", Value: "bar"},
				},
				Network: &Network{ID: "netw_zDW7KYAeqqfRfVag"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func Test_virtualMachineBuildCreateFromSpecRequest_JSONMarshaling(
	t *testing.T,
) {
	tests := []struct {
		name string
		obj  *virtualMachineBuildCreateFromSpecRequest
	}{
		{
			name: "empty",
			obj:  &virtualMachineBuildCreateFromSpecRequest{},
		},
		{
			name: "full",
			obj: &virtualMachineBuildCreateFromSpecRequest{
				Organization: &Organization{ID: "org_O648YDMEYeLmqdmn"},
				XML: undent.String(`
					<?xml version="1.0" encoding="UTF-8"?>
					<VirtualMachineSpec>
						<DataCenter by="permalink">london</DataCenter>
					</VirtualMachineSpec>`,
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
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
		want       *VirtualMachineBuild
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "virtual machine build",
			args: args{
				ctx: context.Background(),
				id:  "vmbuild_pbjJIqJ3MOMNsCr3",
			},
			want: &VirtualMachineBuild{
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
			want: &VirtualMachineBuild{
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
			rm, mux, _, teardown := test.PrepareTestClient(t)
			defer teardown()
			c := NewVirtualMachineBuildsClient(rm)

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

			got, resp, err := c.Get(
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

			if tt.want != nil {
				assert.Equal(t, tt.want, got)
			}

			if tt.errResp != nil {
				assert.Equal(t, tt.errResp, resp.Error)
			}
		})
	}
}

func TestVirtualMachineBuildsClient_GetByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name       string
		args       args
		id         string
		want       *VirtualMachineBuild
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "virtual machine build",
			args: args{
				ctx: context.Background(),
				id:  "vmbuild_pbjJIqJ3MOMNsCr3",
			},
			want: &VirtualMachineBuild{
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
			want: &VirtualMachineBuild{
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
			rm, mux, _, teardown := test.PrepareTestClient(t)
			defer teardown()
			c := NewVirtualMachineBuildsClient(rm)
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

			got, resp, err := c.GetByID(
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

			if tt.want != nil {
				assert.Equal(t, tt.want, got)
			}

			if tt.errResp != nil {
				assert.Equal(t, tt.errResp, resp.Error)
			}
		})
	}
}

func TestVirtualMachineBuildsClient_Create(t *testing.T) {
	fullArgs := &VirtualMachineBuildArguments{
		Zone: &Zone{
			ID:        "zone_kY2sPRG24sJVRM2U",
			Name:      "North West",
			Permalink: "north-west",
		},
		DataCenter: &DataCenter{
			ID:        "dc_25d48761871e4bf",
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

	type args struct {
		ctx       context.Context
		org       *Organization
		buildArgs *VirtualMachineBuildArguments
	}
	tests := []struct {
		name       string
		args       args
		reqBody    *virtualMachineBuildCreateRequest
		want       *VirtualMachineBuild
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "virtual machine build",
			args: args{
				ctx: context.Background(),
				org: &Organization{
					ID:        "org_O648YDMEYeLmqdmn",
					Name:      "ACME Inc.",
					SubDomain: "acme",
				},
				buildArgs: fullArgs,
			},
			reqBody: &virtualMachineBuildCreateRequest{
				Organization: &Organization{
					ID: "org_O648YDMEYeLmqdmn",
				},
				Zone: &Zone{
					ID: "zone_kY2sPRG24sJVRM2U",
				},
				DataCenter: &DataCenter{
					ID: "dc_25d48761871e4bf",
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
			want: &VirtualMachineBuild{
				ID:    "vmbuild_TEmhezUShNuAsyac",
				State: VirtualMachineBuildPending,
			},
			respStatus: http.StatusCreated,
			respBody:   fixture("virtual_machine_build_create"),
		},
		{
			name: "virtual machine build (no IDs)",
			args: args{
				ctx: context.Background(),
				org: &Organization{
					Name:      "ACME Inc.",
					SubDomain: "acme",
				},
				buildArgs: &VirtualMachineBuildArguments{
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
			},
			reqBody: &virtualMachineBuildCreateRequest{
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
			want: &VirtualMachineBuild{
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
				org:       &Organization{ID: "org_O648YDMEYeLmqdmn"},
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
				org:       &Organization{ID: "org_O648YDMEYeLmqdmn"},
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
				org:       &Organization{ID: "org_O648YDMEYeLmqdmn"},
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
				org:       &Organization{ID: "org_O648YDMEYeLmqdmn"},
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
				org:       &Organization{ID: "org_O648YDMEYeLmqdmn"},
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
				org:       &Organization{ID: "org_O648YDMEYeLmqdmn"},
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
				org:       &Organization{ID: "org_O648YDMEYeLmqdmn"},
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
				org:       &Organization{ID: "org_O648YDMEYeLmqdmn"},
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
				org:       &Organization{ID: "org_O648YDMEYeLmqdmn"},
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
				org:       &Organization{ID: "org_O648YDMEYeLmqdmn"},
				buildArgs: fullArgs,
			},
			errStr: "location_required: A zone or a data_center argument " +
				"must be provided",
			errResp: &katapult.ResponseError{
				Code: "location_required",
				Description: "A zone or a data_center argument must be " +
					"provided",
				Detail: json.RawMessage(`{}`),
			},
			respStatus: http.StatusUnprocessableEntity,
			respBody:   fixture("location_required_error"),
		},
		{
			name: "nil organization",
			args: args{
				ctx:       context.Background(),
				org:       nil,
				buildArgs: fullArgs,
			},
			errStr:     fixtureOrganizationNotFoundErr,
			errResp:    fixtureOrganizationNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("organization_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx:       nil,
				org:       &Organization{ID: "org_O648YDMEYeLmqdmn"},
				buildArgs: fullArgs,
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := test.PrepareTestClient(t)
			defer teardown()
			c := NewVirtualMachineBuildsClient(rm)
			mux.HandleFunc(
				"/core/v1/organizations/_/virtual_machines/build",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "POST", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					if tt.reqBody != nil {
						reqBody := &virtualMachineBuildCreateRequest{}
						err := strictUmarshal(r.Body, reqBody)
						assert.NoError(t, err)
						assert.Equal(t, tt.reqBody, reqBody)
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.Create(
				tt.args.ctx, tt.args.org, tt.args.buildArgs,
			)

			if tt.respStatus != 0 {
				assert.Equal(t, tt.respStatus, resp.StatusCode)
			}

			if tt.errStr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.errStr)
			}

			if tt.want != nil {
				assert.Equal(t, tt.want, got)
			}

			if tt.errResp != nil {
				assert.Equal(t, tt.errResp, resp.Error)
			}
		})
	}
}

func TestVirtualMachineBuildsClient_CreateFromSpec(t *testing.T) {
	spec := &buildspec.VirtualMachineSpec{
		DataCenter: &buildspec.DataCenter{
			Permalink: "london",
		},
		Resources: &buildspec.Resources{
			Package: &buildspec.Package{
				Permalink: "rock-3",
			},
		},
		DiskTemplate: &buildspec.DiskTemplate{
			Permalink: "templates/ubuntu-18-04",
		},
		Hostname: "web-3",
	}
	xmlSpec, _ := spec.XML()

	type args struct {
		ctx  context.Context
		org  *Organization
		spec *buildspec.VirtualMachineSpec
	}
	tests := []struct {
		name       string
		args       args
		reqBody    *virtualMachineBuildCreateFromSpecRequest
		want       *VirtualMachineBuild
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "virtual machine build",
			args: args{
				ctx: context.Background(),
				org: &Organization{
					ID:        "org_O648YDMEYeLmqdmn",
					Name:      "ACME Inc.",
					SubDomain: "acme",
				},
				spec: spec,
			},
			reqBody: &virtualMachineBuildCreateFromSpecRequest{
				Organization: &Organization{
					ID: "org_O648YDMEYeLmqdmn",
				},
				XML: string(xmlSpec),
			},
			want: &VirtualMachineBuild{
				ID:    "vmbuild_TEmhezUShNuAsyac",
				State: VirtualMachineBuildPending,
			},
			respStatus: http.StatusCreated,
			respBody:   fixture("virtual_machine_build_create"),
		},
		{
			name: "invalid API token response",
			args: args{
				ctx:  context.Background(),
				org:  &Organization{ID: "org_O648YDMEYeLmqdmn"},
				spec: spec,
			},
			errStr:     fixtureInvalidAPITokenErr,
			errResp:    fixtureInvalidAPITokenResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("invalid_api_token_error"),
		},
		{
			name: "invalid XML spec",
			args: args{
				ctx:  context.Background(),
				org:  &Organization{ID: "org_O648YDMEYeLmqdmn"},
				spec: spec,
			},
			errStr:     fixtureInvalidXMLSpecErr,
			errResp:    fixtureInvalidXMLSpecResponseError,
			respStatus: http.StatusBadRequest,
			respBody:   fixture("invalid_spec_xml_error"),
		},
		{
			name: "non-existent organization",
			args: args{
				ctx:  context.Background(),
				org:  &Organization{ID: "org_O648YDMEYeLmqdmn"},
				spec: spec,
			},
			errStr:     fixtureOrganizationNotFoundErr,
			errResp:    fixtureOrganizationNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("organization_not_found_error"),
		},
		{
			name: "suspended organization",
			args: args{
				ctx:  context.Background(),
				org:  &Organization{ID: "org_O648YDMEYeLmqdmn"},
				spec: spec,
			},
			errStr:     fixtureOrganizationSuspendedErr,
			errResp:    fixtureOrganizationSuspendedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("organization_suspended_error"),
		},
		{
			name: "permission denied",
			args: args{
				ctx:  context.Background(),
				org:  &Organization{ID: "org_O648YDMEYeLmqdmn"},
				spec: spec,
			},
			errStr:     fixturePermissionDeniedErr,
			errResp:    fixturePermissionDeniedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("permission_denied_error"),
		},
		{
			name: "validation error",
			args: args{
				ctx:  context.Background(),
				org:  &Organization{ID: "org_O648YDMEYeLmqdmn"},
				spec: spec,
			},
			errStr:     fixtureValidationErrorErr,
			errResp:    fixtureValidationErrorResponseError,
			respStatus: http.StatusUnprocessableEntity,
			respBody:   fixture("validation_error"),
		},
		{
			name: "nil organization",
			args: args{
				ctx:  context.Background(),
				org:  nil,
				spec: spec,
			},
			errStr:     fixtureOrganizationNotFoundErr,
			errResp:    fixtureOrganizationNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("organization_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx:  nil,
				org:  &Organization{ID: "org_O648YDMEYeLmqdmn"},
				spec: spec,
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := test.PrepareTestClient(t)
			defer teardown()
			c := NewVirtualMachineBuildsClient(rm)
			mux.HandleFunc(
				"/core/v1/organizations/_/virtual_machines/build_from_spec",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "POST", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					if tt.reqBody != nil {
						reqBody := &virtualMachineBuildCreateFromSpecRequest{}
						err := strictUmarshal(r.Body, reqBody)
						assert.NoError(t, err)
						assert.Equal(t, tt.reqBody, reqBody)
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.CreateFromSpec(
				tt.args.ctx, tt.args.org, tt.args.spec,
			)

			if tt.respStatus != 0 {
				assert.Equal(t, tt.respStatus, resp.StatusCode)
			}

			if tt.errStr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.errStr)
			}

			if tt.want != nil {
				assert.Equal(t, tt.want, got)
			}

			if tt.errResp != nil {
				assert.Equal(t, tt.errResp, resp.Error)
			}
		})
	}
}

func TestVirtualMachineBuildsClient_CreateFromSpecXML(t *testing.T) {
	//nolint:lll
	specXML := undent.String(`
		<?xml version="1.0" encoding="UTF-8"?>
		<VirtualMachineSpec>
			<DataCenter by="permalink">london</DataCenter>
			<Resources>
				<Package by="permalink">rock-3</Package>
			</Resources>
			<DiskTemplate>
				<DiskTemplate by="permalink">templates/ubuntu-18-04</DiskTemplate>
			</DiskTemplate>
			<Hostname>
				<Hostname>web-3</Hostname>
			</Hostname>
		</VirtualMachineSpec>`,
	)

	type args struct {
		ctx context.Context
		org *Organization
		xml string
	}
	tests := []struct {
		name       string
		args       args
		reqBody    *virtualMachineBuildCreateFromSpecRequest
		want       *VirtualMachineBuild
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "virtual machine build",
			args: args{
				ctx: context.Background(),
				org: &Organization{
					ID:        "org_O648YDMEYeLmqdmn",
					Name:      "ACME Inc.",
					SubDomain: "acme",
				},
				xml: specXML,
			},
			reqBody: &virtualMachineBuildCreateFromSpecRequest{
				Organization: &Organization{
					ID: "org_O648YDMEYeLmqdmn",
				},
				XML: specXML,
			},
			want: &VirtualMachineBuild{
				ID:    "vmbuild_TEmhezUShNuAsyac",
				State: VirtualMachineBuildPending,
			},
			respStatus: http.StatusCreated,
			respBody:   fixture("virtual_machine_build_create"),
		},
		{
			name: "invalid API token response",
			args: args{
				ctx: context.Background(),
				org: &Organization{ID: "org_O648YDMEYeLmqdmn"},
				xml: specXML,
			},
			errStr:     fixtureInvalidAPITokenErr,
			errResp:    fixtureInvalidAPITokenResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("invalid_api_token_error"),
		},
		{
			name: "invalid XML spec",
			args: args{
				ctx: context.Background(),
				org: &Organization{ID: "org_O648YDMEYeLmqdmn"},
				xml: specXML,
			},
			errStr:     fixtureInvalidXMLSpecErr,
			errResp:    fixtureInvalidXMLSpecResponseError,
			respStatus: http.StatusBadRequest,
			respBody:   fixture("invalid_spec_xml_error"),
		},
		{
			name: "non-existent organization",
			args: args{
				ctx: context.Background(),
				org: &Organization{ID: "org_O648YDMEYeLmqdmn"},
				xml: specXML,
			},
			errStr:     fixtureOrganizationNotFoundErr,
			errResp:    fixtureOrganizationNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("organization_not_found_error"),
		},
		{
			name: "suspended organization",
			args: args{
				ctx: context.Background(),
				org: &Organization{ID: "org_O648YDMEYeLmqdmn"},
				xml: specXML,
			},
			errStr:     fixtureOrganizationSuspendedErr,
			errResp:    fixtureOrganizationSuspendedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("organization_suspended_error"),
		},
		{
			name: "permission denied",
			args: args{
				ctx: context.Background(),
				org: &Organization{ID: "org_O648YDMEYeLmqdmn"},
				xml: specXML,
			},
			errStr:     fixturePermissionDeniedErr,
			errResp:    fixturePermissionDeniedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("permission_denied_error"),
		},
		{
			name: "validation error",
			args: args{
				ctx: context.Background(),
				org: &Organization{ID: "org_O648YDMEYeLmqdmn"},
				xml: specXML,
			},
			errStr:     fixtureValidationErrorErr,
			errResp:    fixtureValidationErrorResponseError,
			respStatus: http.StatusUnprocessableEntity,
			respBody:   fixture("validation_error"),
		},
		{
			name: "nil organization",
			args: args{
				ctx: context.Background(),
				org: nil,
				xml: specXML,
			},
			errStr:     fixtureOrganizationNotFoundErr,
			errResp:    fixtureOrganizationNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("organization_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				org: &Organization{ID: "org_O648YDMEYeLmqdmn"},
				xml: specXML,
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := test.PrepareTestClient(t)
			defer teardown()
			c := NewVirtualMachineBuildsClient(rm)
			mux.HandleFunc(
				"/core/v1/organizations/_/virtual_machines/build_from_spec",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "POST", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					if tt.reqBody != nil {
						reqBody := &virtualMachineBuildCreateFromSpecRequest{}
						err := strictUmarshal(r.Body, reqBody)
						assert.NoError(t, err)
						assert.Equal(t, tt.reqBody, reqBody)
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.CreateFromSpecXML(
				tt.args.ctx, tt.args.org, tt.args.xml,
			)

			if tt.respStatus != 0 {
				assert.Equal(t, tt.respStatus, resp.StatusCode)
			}

			if tt.errStr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.errStr)
			}

			if tt.want != nil {
				assert.Equal(t, tt.want, got)
			}

			if tt.errResp != nil {
				assert.Equal(t, tt.errResp, resp.Error)
			}
		})
	}
}
