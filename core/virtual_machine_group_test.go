package core

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/krystal/go-katapult"
	"github.com/krystal/go-katapult/internal/test"
	"github.com/stretchr/testify/assert"
)

var (
	fixtureVMGroupNotFoundErr = "virtual_machine_group_not_found: " +
		"No virtual machine group was found matching any of the criteria " +
		"provided in the arguments"
	fixtureVMGroupNotFoundResponseError = &katapult.ResponseError{
		Code: "virtual_machine_group_not_found",
		Description: "No virtual machine group was found matching any of " +
			"the criteria provided in the arguments",
		Detail: json.RawMessage(`{}`),
	}
)

func TestClient_VirtualMachineGroups(t *testing.T) {
	c := New(&fakeRequestMaker{})

	assert.IsType(t, &VirtualMachineGroupsClient{}, c.VirtualMachineGroups)
}

func TestVirtualMachineGroup_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name    string
		obj     *VirtualMachineGroup
		decoded *VirtualMachineGroup
	}{
		{
			name: "empty",
			obj:  &VirtualMachineGroup{},
		},
		{
			name: "full",
			obj: &VirtualMachineGroup{
				ID:        "id",
				Name:      "name",
				Segregate: true,
				CreatedAt: timestampPtr(934834834),
			},
		},
		{
			name: "null",
			obj: &VirtualMachineGroup{
				ID:        "id",
				Name:      "name",
				Segregate: true,
				CreatedAt: timestampPtr(934834834),
				null:      true,
			},
			decoded: &VirtualMachineGroup{
				null: true,
			},
		},
		{
			name: "NullVirtualMachineGroup",
			obj:  NullVirtualMachineGroup,
			decoded: &VirtualMachineGroup{
				null: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test.CustomJSONMarshaling(t, tt.obj, tt.decoded)
		})
	}
}

func Test_virtualMachineGroupCreateRequest_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *virtualMachineGroupCreateRequest
	}{
		{
			name: "empty",
			obj:  &virtualMachineGroupCreateRequest{},
		},
		{
			name: "full",
			obj: &virtualMachineGroupCreateRequest{
				Organization: &Organization{ID: "org_O648YDMEYeLmqdmn"},
				Properties: &VirtualMachineGroupCreateArguments{
					Name:      "vm group test",
					Segregate: truePtr,
				},
			},
		},
		{
			name: "false segregate",
			obj: &virtualMachineGroupCreateRequest{
				Properties: &VirtualMachineGroupCreateArguments{
					Segregate: falsePtr,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func Test_virtualMachineGroupUpdateRequest_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *virtualMachineGroupUpdateRequest
	}{
		{
			name: "empty",
			obj:  &virtualMachineGroupUpdateRequest{},
		},
		{
			name: "full",
			obj: &virtualMachineGroupUpdateRequest{
				Properties: &VirtualMachineGroupUpdateArguments{
					Name:      "vm group test",
					Segregate: truePtr,
				},
			},
		},
		{
			name: "false segregate",
			obj: &virtualMachineGroupUpdateRequest{
				Properties: &VirtualMachineGroupUpdateArguments{
					Segregate: falsePtr,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func Test_virtualMachineGroupsResponseBody_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *virtualMachineGroupsResponseBody
	}{
		{
			name: "empty",
			obj:  &virtualMachineGroupsResponseBody{},
		},
		{
			name: "full",
			obj: &virtualMachineGroupsResponseBody{
				VirtualMachineGroup:  &VirtualMachineGroup{ID: "id1"},
				VirtualMachineGroups: []*VirtualMachineGroup{{ID: "id2"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestVirtualMachineGroupsClient_List(t *testing.T) {
	// Correlates to fixtures/virtual_machine_groups_list*.json
	virtualMachineGroupList := []*VirtualMachineGroup{
		{
			ID:        "vmgrp_gsEUFPp3ybVQm5QQ",
			Name:      "vm group 1",
			Segregate: true,
		},
		{
			ID:        "vmgrp_bcfdEFn2viWCm5ve",
			Name:      "vm group 2",
			Segregate: false,
		},
	}

	type args struct {
		ctx context.Context
		org *Organization
	}
	tests := []struct {
		name       string
		args       args
		want       []*VirtualMachineGroup
		wantQuery  *url.Values
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "by organization ID",
			args: args{
				ctx: context.Background(),
				org: &Organization{ID: "org_O648YDMEYeLmqdmn"},
			},
			want: virtualMachineGroupList,
			wantQuery: &url.Values{
				"organization[id]": []string{"org_O648YDMEYeLmqdmn"},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machine_groups_list"),
		},
		{
			name: "by organization SubDomain",
			args: args{
				ctx: context.Background(),
				org: &Organization{SubDomain: "acme"},
			},
			want: virtualMachineGroupList,
			wantQuery: &url.Values{
				"organization[sub_domain]": []string{"acme"},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machine_groups_list"),
		},
		{
			name: "invalid API token response",
			args: args{
				ctx: context.Background(),
				org: &Organization{ID: "org_O648YDMEYeLmqdmn"},
			},
			errStr:     fixtureInvalidAPITokenErr,
			errResp:    fixtureInvalidAPITokenResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("invalid_api_token_error"),
		},
		{
			name: "non-existent organization",
			args: args{
				ctx: context.Background(),
				org: &Organization{ID: "org_O648YDMEYeLmqdmn"},
			},
			errStr:     fixtureOrganizationNotFoundErr,
			errResp:    fixtureOrganizationNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("organization_not_found_error"),
		},
		{
			name: "nil organization",
			args: args{
				ctx: context.Background(),
				org: nil,
			},
			errStr:     fixtureOrganizationNotFoundErr,
			errResp:    fixtureOrganizationNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("organization_not_found_error"),
		},
		{
			name: "not activated organization",
			args: args{
				ctx: context.Background(),
				org: &Organization{ID: "org_O648YDMEYeLmqdmn"},
			},
			errStr:     fixtureOrganizationNotActivatedErr,
			errResp:    fixtureOrganizationNotActivatedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("organization_not_activated_error"),
		},
		{
			name: "suspended organization",
			args: args{
				ctx: context.Background(),
				org: &Organization{ID: "org_O648YDMEYeLmqdmn"},
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
			},
			errStr:     fixturePermissionDeniedErr,
			errResp:    fixturePermissionDeniedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("permission_denied_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				org: &Organization{ID: "org_O648YDMEYeLmqdmn"},
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewVirtualMachineGroupsClient(rm)

			mux.HandleFunc(
				"/core/v1/organizations/_/virtual_machine_groups",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					if tt.wantQuery != nil {
						assert.Equal(t, *tt.wantQuery, r.URL.Query())
					} else {
						qs := queryValues(tt.args.org)
						assert.Equal(t, *qs, r.URL.Query())
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.List(
				tt.args.ctx, tt.args.org,
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

func TestVirtualMachineGroupsClient_Get(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name       string
		args       args
		want       *VirtualMachineGroup
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "by ID",
			args: args{
				ctx: context.Background(),
				id:  "vmgrp_gsEUFPp3ybVQm5QQ",
			},
			want: &VirtualMachineGroup{
				ID:        "vmgrp_gsEUFPp3ybVQm5QQ",
				Name:      "vm group 1",
				Segregate: true,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machine_group_get"),
		},
		{
			name: "non-existent virtual machine group",
			args: args{
				ctx: context.Background(),
				id:  "vmgrp_gsEUFPp3ybVQm5QQ",
			},
			errStr:     fixtureVMGroupNotFoundErr,
			errResp:    fixtureVMGroupNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody: fixture(
				"virtual_machine_group_not_found_error",
			),
		},
		{
			name: "permission denied",
			args: args{
				ctx: context.Background(),
				id:  "vmgrp_gsEUFPp3ybVQm5QQ",
			},
			errStr:     fixturePermissionDeniedErr,
			errResp:    fixturePermissionDeniedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("permission_denied_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				id:  "vmgrp_gsEUFPp3ybVQm5QQ",
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewVirtualMachineGroupsClient(rm)

			mux.HandleFunc(
				fmt.Sprintf(
					"/core/v1/virtual_machine_groups/%s",
					tt.args.id,
				),
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

func TestVirtualMachineGroupsClient_GetByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name       string
		args       args
		want       *VirtualMachineGroup
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "by ID",
			args: args{
				ctx: context.Background(),
				id:  "vmgrp_gsEUFPp3ybVQm5QQ",
			},
			want: &VirtualMachineGroup{
				ID:        "vmgrp_gsEUFPp3ybVQm5QQ",
				Name:      "vm group 1",
				Segregate: true,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machine_group_get"),
		},
		{
			name: "non-existent virtual machine group",
			args: args{
				ctx: context.Background(),
				id:  "vmgrp_gsEUFPp3ybVQm5QQ",
			},
			errStr:     fixtureVMGroupNotFoundErr,
			errResp:    fixtureVMGroupNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody: fixture(
				"virtual_machine_group_not_found_error",
			),
		},
		{
			name: "permission denied",
			args: args{
				ctx: context.Background(),
				id:  "vmgrp_gsEUFPp3ybVQm5QQ",
			},
			errStr:     fixturePermissionDeniedErr,
			errResp:    fixturePermissionDeniedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("permission_denied_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				id:  "vmgrp_gsEUFPp3ybVQm5QQ",
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewVirtualMachineGroupsClient(rm)

			mux.HandleFunc(
				fmt.Sprintf(
					"/core/v1/virtual_machine_groups/%s",
					tt.args.id,
				),
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

func TestVirtualMachineGroupsClient_Create(t *testing.T) {
	vmGroupArgs := &VirtualMachineGroupCreateArguments{
		Name:      "vm group test",
		Segregate: falsePtr,
	}

	type args struct {
		ctx  context.Context
		org  *Organization
		args *VirtualMachineGroupCreateArguments
	}
	tests := []struct {
		name       string
		args       args
		reqBody    *virtualMachineGroupCreateRequest
		want       *VirtualMachineGroup
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "virtual machine group by organization ID",
			args: args{
				ctx: context.Background(),
				org: &Organization{
					ID:        "org_O648YDMEYeLmqdmn",
					Name:      "ACME Inc.",
					SubDomain: "acme",
				},
				args: vmGroupArgs,
			},
			reqBody: &virtualMachineGroupCreateRequest{
				Organization: &Organization{ID: "org_O648YDMEYeLmqdmn"},
				Properties: &VirtualMachineGroupCreateArguments{
					Name:      "vm group test",
					Segregate: falsePtr,
				},
			},
			want: &VirtualMachineGroup{
				ID:        "vmgrp_gsEUFPp3ybVQm5QQ",
				Name:      "vm group test",
				Segregate: false,
			},
			respStatus: http.StatusCreated,
			respBody:   fixture("virtual_machine_group_create"),
		},
		{
			name: "virtual machine group by organization SubDomain",
			args: args{
				ctx: context.Background(),
				org: &Organization{
					Name:      "ACME Inc.",
					SubDomain: "acme",
				},
				args: vmGroupArgs,
			},
			reqBody: &virtualMachineGroupCreateRequest{
				Organization: &Organization{SubDomain: "acme"},
				Properties: &VirtualMachineGroupCreateArguments{
					Name:      "vm group test",
					Segregate: falsePtr,
				},
			},
			want: &VirtualMachineGroup{
				ID:        "vmgrp_gsEUFPp3ybVQm5QQ",
				Name:      "vm group test",
				Segregate: false,
			},
			respStatus: http.StatusCreated,
			respBody:   fixture("virtual_machine_group_create"),
		},
		{
			name: "non-existent Organization",
			args: args{
				ctx:  context.Background(),
				org:  &Organization{ID: "org_O648YDMEYeLmqdmn"},
				args: vmGroupArgs,
			},
			errStr:     fixtureOrganizationNotFoundErr,
			errResp:    fixtureOrganizationNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("organization_not_found_error"),
		},
		{
			name: "suspended Organization",
			args: args{
				ctx:  context.Background(),
				org:  &Organization{ID: "org_O648YDMEYeLmqdmn"},
				args: vmGroupArgs,
			},
			errStr:     fixtureOrganizationSuspendedErr,
			errResp:    fixtureOrganizationSuspendedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("organization_suspended_error"),
		},
		{
			name: "not activated Organization",
			args: args{
				ctx: context.Background(),
				org: &Organization{ID: "org_O648YDMEYeLmqdmn"},
			},
			errStr:     fixtureOrganizationNotActivatedErr,
			errResp:    fixtureOrganizationNotActivatedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("organization_not_activated_error"),
		},
		{
			name: "permission denied",
			args: args{
				ctx:  context.Background(),
				org:  &Organization{ID: "org_O648YDMEYeLmqdmn"},
				args: vmGroupArgs,
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
				args: vmGroupArgs,
			},
			errStr:     fixtureValidationErrorErr,
			errResp:    fixtureValidationErrorResponseError,
			respStatus: http.StatusUnprocessableEntity,
			respBody:   fixture("validation_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx:  nil,
				org:  &Organization{ID: "org_O648YDMEYeLmqdmn"},
				args: vmGroupArgs,
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewVirtualMachineGroupsClient(rm)

			mux.HandleFunc(
				"/core/v1/organizations/_/virtual_machine_groups",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "POST", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					if tt.reqBody != nil {
						reqBody := &virtualMachineGroupCreateRequest{}
						err := strictUmarshal(r.Body, reqBody)
						assert.NoError(t, err)
						assert.Equal(t, tt.reqBody, reqBody)
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.Create(
				tt.args.ctx, tt.args.org, tt.args.args,
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

func TestVirtualMachineGroupsClient_Update(t *testing.T) {
	groupArgs := &VirtualMachineGroupUpdateArguments{
		Name:      "vm group test",
		Segregate: truePtr,
	}

	type args struct {
		ctx   context.Context
		group *VirtualMachineGroup
		args  *VirtualMachineGroupUpdateArguments
	}
	tests := []struct {
		name       string
		args       args
		reqBody    *virtualMachineGroupUpdateRequest
		want       *VirtualMachineGroup
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "by ID",
			args: args{
				ctx:   context.Background(),
				group: &VirtualMachineGroup{ID: "vmgrp_gsEUFPp3ybVQm5QQ"},
				args:  groupArgs,
			},
			reqBody: &virtualMachineGroupUpdateRequest{
				VirtualMachineGroup: &VirtualMachineGroup{
					ID: "vmgrp_gsEUFPp3ybVQm5QQ",
				},
				Properties: &VirtualMachineGroupUpdateArguments{
					Name:      "vm group test",
					Segregate: truePtr,
				},
			},
			want: &VirtualMachineGroup{
				ID:        "vmgrp_gsEUFPp3ybVQm5QQ",
				Name:      "vm group test",
				Segregate: true,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machine_group_updated"),
		},
		{
			name: "non-existent virtual machine group",
			args: args{
				ctx:   context.Background(),
				group: &VirtualMachineGroup{ID: "vmgrp_nopethisdoesnotexist"},
				args:  &VirtualMachineGroupUpdateArguments{},
			},
			errStr:     fixtureVMGroupNotFoundErr,
			errResp:    fixtureVMGroupNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("virtual_machine_group_not_found_error"),
		},
		{
			name: "permission denied",
			args: args{
				ctx:   context.Background(),
				group: &VirtualMachineGroup{ID: "vmgrp_gsEUFPp3ybVQm5QQ"},
				args:  groupArgs,
			},
			errStr:     fixturePermissionDeniedErr,
			errResp:    fixturePermissionDeniedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("permission_denied_error"),
		},
		{
			name: "validation error",
			args: args{
				ctx:   context.Background(),
				group: &VirtualMachineGroup{ID: "vmgrp_gsEUFPp3ybVQm5QQ"},
				args:  groupArgs,
			},
			errStr:     fixtureValidationErrorErr,
			errResp:    fixtureValidationErrorResponseError,
			respStatus: http.StatusUnprocessableEntity,
			respBody:   fixture("validation_error"),
		},
		{
			name: "nil group",
			args: args{
				ctx:   context.Background(),
				group: nil,
				args:  groupArgs,
			},
			errStr:     fixtureVMGroupNotFoundErr,
			errResp:    fixtureVMGroupNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("virtual_machine_group_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				group: &VirtualMachineGroup{
					ID: "vmgrp_gsEUFPp3ybVQm5QQ",
				},
				args: groupArgs,
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewVirtualMachineGroupsClient(rm)

			mux.HandleFunc(
				"/core/v1/virtual_machine_groups/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "PATCH", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					if tt.reqBody != nil {
						reqBody := &virtualMachineGroupUpdateRequest{}
						err := strictUmarshal(r.Body, reqBody)
						assert.NoError(t, err)
						assert.Equal(t, tt.reqBody, reqBody)
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.Update(
				tt.args.ctx, tt.args.group, tt.args.args,
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

func TestVirtualMachineGroupsClient_Delete(t *testing.T) {
	type args struct {
		ctx   context.Context
		group *VirtualMachineGroup
	}
	tests := []struct {
		name       string
		args       args
		wantQuery  *url.Values
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "by ID",
			args: args{
				ctx:   context.Background(),
				group: &VirtualMachineGroup{ID: "vmgrp_gsEUFPp3ybVQm5QQ"},
			},
			wantQuery: &url.Values{
				"virtual_machine_group[id]": []string{"vmgrp_gsEUFPp3ybVQm5QQ"},
			},
			respStatus: http.StatusOK,
			respBody:   []byte("{}"),
		},
		{
			name: "non-existent virtual machine group",
			args: args{
				ctx:   context.Background(),
				group: &VirtualMachineGroup{ID: "vmgrp_nopenotfound"},
			},
			errStr:     fixtureVMGroupNotFoundErr,
			errResp:    fixtureVMGroupNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("virtual_machine_group_not_found_error"),
		},
		{
			name: "permission denied",
			args: args{
				ctx:   context.Background(),
				group: &VirtualMachineGroup{ID: "vmgrp_gsEUFPp3ybVQm5QQ"},
			},
			errStr:     fixturePermissionDeniedErr,
			errResp:    fixturePermissionDeniedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("permission_denied_error"),
		},
		{
			name: "nil group",
			args: args{
				ctx:   context.Background(),
				group: nil,
			},
			errStr:     fixtureVMGroupNotFoundErr,
			errResp:    fixtureVMGroupNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("virtual_machine_group_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx:   nil,
				group: &VirtualMachineGroup{ID: "vmgrp_gsEUFPp3ybVQm5QQ"},
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewVirtualMachineGroupsClient(rm)

			mux.HandleFunc(
				"/core/v1/virtual_machine_groups/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "DELETE", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					if tt.wantQuery != nil {
						assert.Equal(t, *tt.wantQuery, r.URL.Query())
					} else {
						assert.Equal(t,
							*tt.args.group.queryValues(), r.URL.Query(),
						)
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			resp, err := c.Delete(
				tt.args.ctx,
				tt.args.group,
			)

			if tt.respStatus != 0 {
				assert.Equal(t, tt.respStatus, resp.StatusCode)
			}

			if tt.errStr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.errStr)
			}

			if tt.errResp != nil {
				assert.Equal(t, tt.errResp, resp.Error)
			}
		})
	}
}
