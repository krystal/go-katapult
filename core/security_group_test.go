package core

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/krystal/go-katapult"
	"github.com/stretchr/testify/assert"
)

var (
	fixtureSecurityGroupNotFoundErr = "security_group_not_found: No security " +
		"group was found matching any of the criteria provided in " +
		"the arguments"
	fixtureSecurityGroupNotFoundResponseError = &katapult.ResponseError{
		Code: "security_group_not_found",
		Description: "No security group was found matching any of the " +
			"criteria provided in the arguments",
		Detail: json.RawMessage(`{}`),
	}
)

func TestClient_SecurityGroups(t *testing.T) {
	c := New(&fakeRequestMaker{})

	assert.IsType(t, &SecurityGroupsClient{}, c.SecurityGroups)
}

func TestSecurityGroup_Ref(t *testing.T) {
	sg := SecurityGroup{ID: "sg_3uXbmANw4sQiF1J3"}
	assert.Equal(t, SecurityGroupRef{ID: "sg_3uXbmANw4sQiF1J3"}, sg.Ref())
}

func TestSecurityGroup_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *SecurityGroup
	}{
		{
			name: "empty",
			obj:  &SecurityGroup{},
		},
		{
			name: "full",
			obj: &SecurityGroup{
				ID:               "sg_3uXbmANw4sQiF1J3",
				Name:             "group-1",
				AllowAllInbound:  true,
				AllowAllOutbound: true,
				Associations:     []string{"id2", "id3"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestSecurityGroupRef_queryValues(t *testing.T) {
	tests := []struct {
		name string
		obj  SecurityGroupRef
	}{
		{
			name: "with id",
			obj: SecurityGroupRef{
				ID: "sg_3uXbmANw4sQiF1J3",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testQueryableEncoding(t, tt.obj)
		})
	}
}

func TestSecurityGroupCreateArguments_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *SecurityGroupCreateArguments
	}{
		{
			name: "empty",
			obj:  &SecurityGroupCreateArguments{},
		},
		{
			name: "full",
			obj: &SecurityGroupCreateArguments{
				Name:             "new-group",
				AllowAllInbound:  truePtr,
				AllowAllOutbound: truePtr,
				Associations:     &[]string{"id1", "id2"},
			},
		},
		{
			name: "false AllowAllInbound",
			obj: &SecurityGroupCreateArguments{
				AllowAllInbound: falsePtr,
			},
		},
		{
			name: "false AllowAllOutbound",
			obj: &SecurityGroupCreateArguments{
				AllowAllOutbound: falsePtr,
			},
		},
		{
			name: "empty Associations",
			obj:  &SecurityGroupCreateArguments{Associations: &[]string{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestSecurityGroupUpdateArguments_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *SecurityGroupUpdateArguments
	}{
		{
			name: "empty",
			obj:  &SecurityGroupUpdateArguments{},
		},
		{
			name: "full",
			obj: &SecurityGroupUpdateArguments{
				Name:             "new-group",
				AllowAllInbound:  truePtr,
				AllowAllOutbound: truePtr,
				Associations:     &[]string{"id1", "id2"},
			},
		},
		{
			name: "empty Associations",
			obj:  &SecurityGroupUpdateArguments{Associations: &[]string{}},
		},
		{
			name: "false AllowAllInbound",
			obj: &SecurityGroupUpdateArguments{
				AllowAllInbound: falsePtr,
			},
		},
		{
			name: "false AllowAllOutbound",
			obj: &SecurityGroupUpdateArguments{
				AllowAllOutbound: falsePtr,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func Test_securityGroupCreateRequest_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *securityGroupCreateRequest
	}{
		{
			name: "empty",
			obj:  &securityGroupCreateRequest{},
		},
		{
			name: "full",
			obj: &securityGroupCreateRequest{
				Organization: OrganizationRef{ID: "org_rs55YZNYMw7o3jnQ"},
				Properties:   &SecurityGroupCreateArguments{Name: "group-1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func Test_securityGroupUpdateRequest_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *securityGroupUpdateRequest
	}{
		{
			name: "empty",
			obj:  &securityGroupUpdateRequest{},
		},
		{
			name: "full",
			obj: &securityGroupUpdateRequest{
				SecurityGroup: SecurityGroupRef{ID: "sg_3uXbmANw4sQiF1J3"},
				Properties:    &SecurityGroupUpdateArguments{Name: "updated"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func Test_SecurityGroupsResponseBody_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *securityGroupsResponseBody
	}{
		{
			name: "empty",
			obj:  &securityGroupsResponseBody{},
		},
		{
			name: "full",
			obj: &securityGroupsResponseBody{
				Pagination:     &katapult.Pagination{CurrentPage: 344},
				SecurityGroup:  &SecurityGroup{ID: "id1"},
				SecurityGroups: []*SecurityGroup{{ID: "id2"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestSecurityGroupsClient_List(t *testing.T) {
	// Correlates to fixtures/security_groups_list*.json
	SecurityGroupList := []*SecurityGroup{
		{
			ID:           "sg_3uXbmANw4sQiF1J3",
			Name:         "group-1",
			Associations: []string{},
		},
		{
			ID:           "sg_NFP2Ns2frZJV8gD1",
			Name:         "group-2",
			Associations: []string{},
		},
		{
			ID:           "sg_FcIOv1SCf8366ZxJ",
			Name:         "group-3",
			Associations: []string{},
		},
	}

	type args struct {
		ctx  context.Context
		org  OrganizationRef
		opts *ListOptions
	}
	tests := []struct {
		name           string
		args           args
		want           []*SecurityGroup
		wantQuery      *url.Values
		wantPagination *katapult.Pagination
		errStr         string
		errResp        *katapult.ResponseError
		respStatus     int
		respBody       []byte
	}{
		{
			name: "by organization ID",
			args: args{
				ctx: context.Background(),
				org: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
			},
			want: SecurityGroupList,
			wantQuery: &url.Values{
				"organization[id]": []string{"org_O648YDMEYeLmqdmn"},
			},
			wantPagination: &katapult.Pagination{
				CurrentPage: 1,
				TotalPages:  1,
				Total:       3,
				PerPage:     30,
				LargeSet:    false,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("security_groups_list"),
		},
		{
			name: "by organization SubDomain",
			args: args{
				ctx: context.Background(),
				org: OrganizationRef{SubDomain: "acme"},
			},
			want: SecurityGroupList,
			wantQuery: &url.Values{
				"organization[sub_domain]": []string{"acme"},
			},
			wantPagination: &katapult.Pagination{
				CurrentPage: 1,
				TotalPages:  1,
				Total:       3,
				PerPage:     30,
				LargeSet:    false,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("security_groups_list"),
		},
		{
			name: "page 1",
			args: args{
				ctx:  context.Background(),
				org:  OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				opts: &ListOptions{Page: 1, PerPage: 2},
			},
			want: SecurityGroupList[0:2],
			wantQuery: &url.Values{
				"organization[id]": []string{"org_O648YDMEYeLmqdmn"},
				"page":             []string{"1"},
				"per_page":         []string{"2"},
			},
			wantPagination: &katapult.Pagination{
				CurrentPage: 1,
				TotalPages:  2,
				Total:       3,
				PerPage:     2,
				LargeSet:    false,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("security_groups_list_page_1"),
		},
		{
			name: "page 2",
			args: args{
				ctx:  context.Background(),
				org:  OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				opts: &ListOptions{Page: 2, PerPage: 2},
			},
			want: SecurityGroupList[2:],
			wantQuery: &url.Values{
				"organization[id]": []string{"org_O648YDMEYeLmqdmn"},
				"page":             []string{"2"},
				"per_page":         []string{"2"},
			},
			wantPagination: &katapult.Pagination{
				CurrentPage: 2,
				TotalPages:  2,
				Total:       3,
				PerPage:     2,
				LargeSet:    false,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("security_groups_list_page_2"),
		},
		{
			name: "invalid API token response",
			args: args{
				ctx: context.Background(),
				org: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
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
				org: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
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
				org: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
			},
			errStr:     fixtureOrganizationSuspendedErr,
			errResp:    fixtureOrganizationSuspendedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("organization_suspended_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				org: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewSecurityGroupsClient(rm)

			mux.HandleFunc(
				"/core/v1/organizations/_/security_groups",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					qs := queryValues(tt.args.org, tt.args.opts)
					assert.Equal(t, *qs, r.URL.Query())

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.List(
				tt.args.ctx, tt.args.org, tt.args.opts,
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

			if tt.wantPagination != nil {
				assert.Equal(t, tt.wantPagination, resp.Pagination)
			}

			if tt.errResp != nil {
				assert.Equal(t, tt.errResp, resp.Error)
			}
		})
	}
}

func TestSecurityGroupsClient_Get(t *testing.T) {
	type args struct {
		ctx context.Context
		ref SecurityGroupRef
	}
	tests := []struct {
		name       string
		args       args
		id         string
		want       *SecurityGroup
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "security group",
			args: args{
				ctx: context.Background(),
				ref: SecurityGroupRef{ID: "sg_3uXbmANw4sQiF1J3"},
			},
			want: &SecurityGroup{
				ID:           "sg_3uXbmANw4sQiF1J3",
				Name:         "group-1",
				Associations: []string{},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("security_group_get"),
		},
		{
			name: "non-existent security group",
			args: args{
				ctx: context.Background(),
				ref: SecurityGroupRef{ID: "sg_nopethisbegone"},
			},
			errStr:     fixtureSecurityGroupNotFoundErr,
			errResp:    fixtureSecurityGroupNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("security_group_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				ref: SecurityGroupRef{ID: "sg_3uXbmANw4sQiF1J3"},
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewSecurityGroupsClient(rm)

			mux.HandleFunc(
				"/core/v1/security_groups/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					assert.Equal(t,
						*tt.args.ref.queryValues(), r.URL.Query(),
					)

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.Get(tt.args.ctx, tt.args.ref)

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

func TestSecurityGroupsClient_GetByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name       string
		args       args
		id         string
		want       *SecurityGroup
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "security group",
			args: args{
				ctx: context.Background(),
				id:  "sg_3uXbmANw4sQiF1J3",
			},
			want: &SecurityGroup{
				ID:           "sg_3uXbmANw4sQiF1J3",
				Name:         "group-1",
				Associations: []string{},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("security_group_get"),
		},
		{
			name: "non-existent security group",
			args: args{
				ctx: context.Background(),
				id:  "sg_nopethisbegone",
			},
			errStr:     fixtureSecurityGroupNotFoundErr,
			errResp:    fixtureSecurityGroupNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("security_group_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				id:  "sg_3uXbmANw4sQiF1J3",
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewSecurityGroupsClient(rm)

			mux.HandleFunc(
				"/core/v1/security_groups/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.GetByID(tt.args.ctx, tt.args.id)

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

func TestSecurityGroupsClient_Create(t *testing.T) {
	sgArgs := &SecurityGroupCreateArguments{
		Name:             "api-test",
		AllowAllInbound:  truePtr,
		AllowAllOutbound: truePtr,
		Associations:     &[]string{"id1", "id2"},
	}
	sgReqArgs := &SecurityGroupCreateArguments{
		Name:             "api-test",
		AllowAllInbound:  truePtr,
		AllowAllOutbound: truePtr,
		Associations:     &[]string{"id1", "id2"},
	}

	type args struct {
		ctx    context.Context
		org    OrganizationRef
		sgArgs *SecurityGroupCreateArguments
	}
	tests := []struct {
		name       string
		args       args
		reqBody    *securityGroupCreateRequest
		want       *SecurityGroup
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "security group",
			args: args{
				ctx:    context.Background(),
				org:    OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				sgArgs: sgArgs,
			},
			reqBody: &securityGroupCreateRequest{
				Organization: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				Properties:   sgReqArgs,
			},
			want: &SecurityGroup{
				ID:               "sg_3uXbmANw4sQiF1J3",
				Name:             "api-test",
				AllowAllInbound:  true,
				AllowAllOutbound: true,
				Associations:     []string{"id1", "id2"},
			},
			respStatus: http.StatusCreated,
			respBody:   fixture("security_group_create"),
		},
		{
			name: "organization by sub-domain",
			args: args{
				ctx:    context.Background(),
				org:    OrganizationRef{SubDomain: "acme"},
				sgArgs: sgArgs,
			},
			reqBody: &securityGroupCreateRequest{
				Organization: OrganizationRef{SubDomain: "acme"},
				Properties:   sgReqArgs,
			},
			want: &SecurityGroup{
				ID:               "sg_3uXbmANw4sQiF1J3",
				Name:             "api-test",
				AllowAllInbound:  true,
				AllowAllOutbound: true,
				Associations:     []string{"id1", "id2"},
			},
			respStatus: http.StatusCreated,
			respBody:   fixture("security_group_create"),
		},
		{
			name: "without associations",
			args: args{
				ctx: context.Background(),
				org: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				sgArgs: &SecurityGroupCreateArguments{
					Name:             sgArgs.Name,
					AllowAllInbound:  sgArgs.AllowAllInbound,
					AllowAllOutbound: sgArgs.AllowAllOutbound,
				},
			},
			reqBody: &securityGroupCreateRequest{
				Organization: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				Properties: &SecurityGroupCreateArguments{
					Name:             sgReqArgs.Name,
					AllowAllInbound:  sgReqArgs.AllowAllInbound,
					AllowAllOutbound: sgReqArgs.AllowAllOutbound,
				},
			},
			want: &SecurityGroup{
				ID:               "sg_3uXbmANw4sQiF1J3",
				Name:             "api-test",
				AllowAllInbound:  true,
				AllowAllOutbound: true,
				Associations:     []string{},
			},
			respStatus: http.StatusCreated,
			respBody:   fixture("security_group_create_without_associations"),
		},
		{
			name: "non-existent Organization",
			args: args{
				ctx:    context.Background(),
				org:    OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				sgArgs: sgArgs,
			},
			errStr:     fixtureOrganizationNotFoundErr,
			errResp:    fixtureOrganizationNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("organization_not_found_error"),
		},
		{
			name: "suspended Organization",
			args: args{
				ctx:    context.Background(),
				org:    OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				sgArgs: sgArgs,
			},
			errStr:     fixtureOrganizationSuspendedErr,
			errResp:    fixtureOrganizationSuspendedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("organization_suspended_error"),
		},
		{
			name: "permission denied",
			args: args{
				ctx:    context.Background(),
				org:    OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				sgArgs: sgArgs,
			},
			errStr:     fixturePermissionDeniedErr,
			errResp:    fixturePermissionDeniedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("permission_denied_error"),
		},
		{
			name: "validation error",
			args: args{
				ctx:    context.Background(),
				org:    OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				sgArgs: sgArgs,
			},
			errStr:     fixtureValidationErrorErr,
			errResp:    fixtureValidationErrorResponseError,
			respStatus: http.StatusUnprocessableEntity,
			respBody:   fixture("validation_error"),
		},
		{
			name: "nil security group arguments",
			args: args{
				ctx:    context.Background(),
				org:    OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				sgArgs: nil,
			},
			errStr:     fixtureValidationErrorErr,
			errResp:    fixtureValidationErrorResponseError,
			respStatus: http.StatusUnprocessableEntity,
			respBody:   fixture("validation_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx:    nil,
				org:    OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				sgArgs: sgArgs,
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewSecurityGroupsClient(rm)

			mux.HandleFunc(
				"/core/v1/organizations/_/security_groups",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "POST", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					if tt.reqBody != nil {
						reqBody := &securityGroupCreateRequest{}
						err := strictUmarshal(r.Body, reqBody)
						assert.NoError(t, err)
						assert.Equal(t, tt.reqBody, reqBody)
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.Create(
				tt.args.ctx, tt.args.org, tt.args.sgArgs,
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

func TestSecurityGroupsClient_Update(t *testing.T) {
	sg := &SecurityGroup{
		ID:               "sg_3uXbmANw4sQiF1J3",
		Name:             "api-test",
		AllowAllInbound:  true,
		AllowAllOutbound: true,
		Associations:     []string{"id1", "id2"},
	}
	sgArgs := &SecurityGroupUpdateArguments{
		Name:             "updated",
		AllowAllInbound:  falsePtr,
		AllowAllOutbound: falsePtr,
		Associations:     &[]string{"id3", "id4"},
	}

	type args struct {
		ctx    context.Context
		sg     SecurityGroupRef
		sgArgs *SecurityGroupUpdateArguments
	}
	tests := []struct {
		name       string
		args       args
		reqBody    *securityGroupUpdateRequest
		want       *SecurityGroup
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "security group",
			args: args{
				ctx:    context.Background(),
				sg:     sg.Ref(),
				sgArgs: sgArgs,
			},
			reqBody: &securityGroupUpdateRequest{
				SecurityGroup: SecurityGroupRef{ID: sg.ID},
				Properties:    sgArgs,
			},
			want: &SecurityGroup{
				ID:               "sg_3uXbmANw4sQiF1J3",
				Name:             "updated",
				AllowAllInbound:  false,
				AllowAllOutbound: false,
				Associations:     []string{"id3", "id4"},
			},
			respStatus: http.StatusCreated,
			respBody:   fixture("security_group_update"),
		},
		{
			name: "security group without ID",
			args: args{
				ctx: context.Background(),
				sg: (&SecurityGroup{
					Name:             sg.Name,
					AllowAllInbound:  sg.AllowAllInbound,
					AllowAllOutbound: sg.AllowAllOutbound,
					Associations:     sg.Associations,
				}).Ref(),
				sgArgs: sgArgs,
			},
			reqBody: &securityGroupUpdateRequest{
				SecurityGroup: SecurityGroupRef{},
				Properties:    sgArgs,
			},
			errStr:     fixtureSecurityGroupNotFoundErr,
			errResp:    fixtureSecurityGroupNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("security_group_not_found_error"),
		},
		{
			name: "non-existent security group",
			args: args{
				ctx:    context.Background(),
				sg:     SecurityGroupRef{ID: "sg_somethingnope"},
				sgArgs: sgArgs,
			},
			reqBody: &securityGroupUpdateRequest{
				SecurityGroup: SecurityGroupRef{ID: "sg_somethingnope"},
				Properties:    sgArgs,
			},
			errStr:     fixtureSecurityGroupNotFoundErr,
			errResp:    fixtureSecurityGroupNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("security_group_not_found_error"),
		},
		{
			name: "permission denied",
			args: args{
				ctx:    context.Background(),
				sg:     sg.Ref(),
				sgArgs: sgArgs,
			},
			errStr:     fixturePermissionDeniedErr,
			errResp:    fixturePermissionDeniedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("permission_denied_error"),
		},
		{
			name: "validation error",
			args: args{
				ctx:    context.Background(),
				sg:     sg.Ref(),
				sgArgs: sgArgs,
			},
			errStr:     fixtureValidationErrorErr,
			errResp:    fixtureValidationErrorResponseError,
			respStatus: http.StatusUnprocessableEntity,
			respBody:   fixture("validation_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				sg:  sg.Ref(),
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewSecurityGroupsClient(rm)

			mux.HandleFunc(
				"/core/v1/security_groups/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "PATCH", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					if tt.reqBody != nil {
						reqBody := &securityGroupUpdateRequest{}
						err := strictUmarshal(r.Body, reqBody)
						assert.NoError(t, err)
						assert.Equal(t, tt.reqBody, reqBody)
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.Update(
				tt.args.ctx, tt.args.sg, tt.args.sgArgs,
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

func TestSecurityGroupsClient_Delete(t *testing.T) {
	type args struct {
		ctx context.Context
		sg  SecurityGroupRef
	}
	tests := []struct {
		name       string
		args       args
		want       *SecurityGroup
		wantQuery  *url.Values
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "by ID",
			args: args{
				ctx: context.Background(),
				sg:  SecurityGroupRef{ID: "sg_3uXbmANw4sQiF1J3"},
			},
			want: &SecurityGroup{
				ID:           "sg_3uXbmANw4sQiF1J3",
				Name:         "group-1",
				Associations: []string{},
			},
			wantQuery: &url.Values{
				"security_group[id]": []string{"sg_3uXbmANw4sQiF1J3"},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("security_group_get"),
		},
		{
			name: "non-existent security group",
			args: args{
				ctx: context.Background(),
				sg:  SecurityGroupRef{ID: "sg_nopenotfound"},
			},
			errStr:     fixtureSecurityGroupNotFoundErr,
			errResp:    fixtureSecurityGroupNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("security_group_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				sg:  SecurityGroupRef{ID: "sg_3uXbmANw4sQiF1J3"},
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewSecurityGroupsClient(rm)

			mux.HandleFunc(
				"/core/v1/security_groups/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "DELETE", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					if tt.wantQuery != nil {
						assert.Equal(t, *tt.wantQuery, r.URL.Query())
					} else {
						assert.Equal(t,
							*tt.args.sg.queryValues(), r.URL.Query(),
						)
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.Delete(tt.args.ctx, tt.args.sg)

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
