package core

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/krystal/go-katapult"
	"github.com/krystal/go-katapult/internal/test"
	"github.com/krystal/go-katapult/internal/testclient"
	"github.com/stretchr/testify/assert"
)

var (
	port  = "3000"
	note  = "My fave security group rule"
	blank = ""
)

func TestClient_SecurityGroupRules(t *testing.T) {
	c := New(&testclient.Client{})

	assert.IsType(t, &SecurityGroupRulesClient{}, c.SecurityGroupRules)
}

func TestSecurityGroupRule_Ref(t *testing.T) {
	tests := []struct {
		name string
		obj  SecurityGroupRule
		want SecurityGroupRuleRef
	}{
		{
			name: "with id",
			obj: SecurityGroupRule{
				ID: "sgr_9IToFxX2AOl7IBSY",
			},
			want: SecurityGroupRuleRef{ID: "sgr_9IToFxX2AOl7IBSY"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.obj.Ref())
		})
	}
}

func TestSecurityGroupRuleRef_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *SecurityGroupRuleRef
	}{
		{
			name: "empty",
			obj:  &SecurityGroupRuleRef{},
		},
		{
			name: "full",
			obj: &SecurityGroupRuleRef{
				ID: "sg_3uXbmANw4sQiF1J3",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestSecurityGroupRule_JSONMarshalling(t *testing.T) {
	tests := []struct {
		name string
		obj  *SecurityGroupRule
	}{
		{
			name: "empty",
			obj:  &SecurityGroupRule{},
		},
		{
			name: "full",
			obj: &SecurityGroupRule{
				ID:        "arbitrary string",
				Direction: "inbound",
				Protocol:  "TCP",
				Ports:     port,
				Targets:   []string{"192.168.0.1"},
				Notes:     note,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestSecurityGroupRuleArguments_JSONMarshalling(t *testing.T) {
	tests := []struct {
		name string
		obj  *SecurityGroupRuleArguments
	}{
		{
			name: "empty",
			obj:  &SecurityGroupRuleArguments{},
		},
		{
			name: "full",
			obj: &SecurityGroupRuleArguments{
				Direction: "inbound",
				Protocol:  "TCP",
				Ports:     &port,
				Targets:   &[]string{"192.168.0.1"},
				Notes:     &note,
			},
		},
		{
			name: "remove all targets",
			obj: &SecurityGroupRuleArguments{
				Targets: &[]string{},
			},
		},
		{
			name: "unset ports",
			obj: &SecurityGroupRuleArguments{
				Ports: &blank,
			},
		},
		{
			name: "unset notes",
			obj: &SecurityGroupRuleArguments{
				Notes: &blank,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func Test_securityGroupRulesResponseBody_JSONMarshalling(t *testing.T) {
	tests := []struct {
		name string
		obj  *securityGroupRulesResponseBody
	}{
		{
			name: "empty",
			obj:  &securityGroupRulesResponseBody{},
		},
		{
			name: "full",
			obj: &securityGroupRulesResponseBody{
				Pagination: &katapult.Pagination{
					LargeSet: true,
				},
				SecurityGroupRule: &SecurityGroupRule{ID: "foobar"},
				SecurityGroupRules: []SecurityGroupRule{
					{
						ID: "barfoo",
					},
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

func Test_securityGroupRuleCreateRequest_JSONMarshalling(t *testing.T) {
	tests := []struct {
		name string
		obj  *securityGroupRuleCreateRequest
	}{
		{
			name: "empty",
			obj:  &securityGroupRuleCreateRequest{},
		},
		{
			name: "full",
			obj: &securityGroupRuleCreateRequest{
				Properties: &SecurityGroupRuleArguments{
					Direction: "inbound",
					Protocol:  "TCP",
					Ports:     &port,
					Targets:   &[]string{"192.168.0.1"},
					Notes:     &note,
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

func Test_securityGroupRuleUpdateRequest_JSONMarshalling(t *testing.T) {
	tests := []struct {
		name string
		obj  *securityGroupRuleUpdateRequest
	}{
		{
			name: "empty",
			obj:  &securityGroupRuleUpdateRequest{},
		},
		{
			name: "full",
			obj: &securityGroupRuleUpdateRequest{
				Properties: &SecurityGroupRuleArguments{
					Direction: "inbound",
					Protocol:  "TCP",
					Ports:     &port,
					Targets:   &[]string{"192.168.0.1"},
					Notes:     &note,
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

func TestSecurityGroupRulesClient_List(t *testing.T) {
	type args struct {
		ctx  context.Context
		sg   SecurityGroupRef
		opts *ListOptions
	}
	tests := []struct {
		name     string
		args     args
		resp     *katapult.Response
		respErr  error
		respV    *securityGroupRulesResponseBody
		wantReq  *katapult.Request
		want     []SecurityGroupRule
		wantResp *katapult.Response
		wantErr  string
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				sg:  SecurityGroupRef{ID: "sg_xw6mLJbhXQjdxfSY"},
				opts: &ListOptions{
					Page:    5,
					PerPage: 32,
				},
			},
			resp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusOK},
			},
			respV: &securityGroupRulesResponseBody{
				Pagination: &katapult.Pagination{Total: 333},
				SecurityGroupRules: []SecurityGroupRule{
					{ID: "sgr_WbCly1EHB3jNMQNC", Direction: "inbound"},
				},
			},
			wantReq: &katapult.Request{
				Method: "GET",
				URL: &url.URL{
					Path: "/core/v1/security_groups/_/rules",
					RawQuery: url.Values{
						"page":               []string{"5"},
						"per_page":           []string{"32"},
						"security_group[id]": []string{"sg_xw6mLJbhXQjdxfSY"},
					}.Encode(),
				},
			},
			want: []SecurityGroupRule{
				{ID: "sgr_WbCly1EHB3jNMQNC", Direction: "inbound"},
			},
			wantResp: &katapult.Response{
				Pagination: &katapult.Pagination{Total: 333},
				Response:   &http.Response{StatusCode: http.StatusOK},
			},
		},
		{
			name: "success with nil options",
			args: args{
				ctx:  context.Background(),
				sg:   SecurityGroupRef{ID: "sg_xw6mLJbhXQjdxfSY"},
				opts: nil,
			},
			resp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusOK},
			},
			respV: &securityGroupRulesResponseBody{
				Pagination: &katapult.Pagination{Total: 124},
				SecurityGroupRules: []SecurityGroupRule{
					{ID: "sgr_WbCly1EHB3jNMQNC", Direction: "inbound"},
				},
			},
			wantReq: &katapult.Request{
				Method: "GET",
				URL: &url.URL{
					Path: "/core/v1/security_groups/_/rules",
					RawQuery: url.Values{
						"security_group[id]": []string{"sg_xw6mLJbhXQjdxfSY"},
					}.Encode(),
				},
			},
			want: []SecurityGroupRule{
				{ID: "sgr_WbCly1EHB3jNMQNC", Direction: "inbound"},
			},
			wantResp: &katapult.Response{
				Pagination: &katapult.Pagination{Total: 124},
				Response:   &http.Response{StatusCode: http.StatusOK},
			},
		},
		{
			name: "request error",
			args: args{
				ctx:  context.Background(),
				sg:   SecurityGroupRef{ID: "sg_xw6mLJbhXQjdxfSY"},
				opts: nil,
			},
			resp: &katapult.Response{
				Response: &http.Response{
					StatusCode: http.StatusInternalServerError,
				},
			},
			respErr: fmt.Errorf("flux capacitor undercharged"),
			wantResp: &katapult.Response{
				Response: &http.Response{
					StatusCode: http.StatusInternalServerError,
				},
			},
			wantErr: "flux capacitor undercharged",
		},
		{
			name: "request error with nil response",
			args: args{
				ctx:  context.Background(),
				sg:   SecurityGroupRef{ID: "sg_xw6mLJbhXQjdxfSY"},
				opts: nil,
			},
			resp:    nil,
			respErr: fmt.Errorf("someting is really wrong"),
			wantErr: "someting is really wrong",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := testclient.New(tt.resp, tt.respErr, tt.respV)
			c := NewSecurityGroupRulesClient(tc)
			ctx := test.Context(tt.args.ctx)

			got, resp, err := c.List(ctx, tt.args.sg, tt.args.opts)

			assert.Equal(t, 1, len(tc.Calls), "only 1 request should be made")
			test.AssertContext(t, ctx, tc.Ctx)

			assert.Equal(t, tt.want, got)

			if tt.wantResp != nil {
				assert.Equal(t, tt.wantResp, resp)
			}

			if tt.wantReq != nil {
				assert.Equal(t, tt.wantReq, tc.Request)
			}

			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}

func TestSecurityGroupRulesClient_Get(t *testing.T) {
	type args struct {
		ctx context.Context
		ref SecurityGroupRuleRef
	}
	tests := []struct {
		name     string
		args     args
		resp     *katapult.Response
		respErr  error
		respV    *securityGroupRulesResponseBody
		wantReq  *katapult.Request
		want     *SecurityGroupRule
		wantResp *katapult.Response
		wantErr  string
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				ref: SecurityGroupRuleRef{ID: "sgr_IppA889KHFY1BvDt"},
			},
			resp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusOK},
			},
			respV: &securityGroupRulesResponseBody{
				SecurityGroupRule: &SecurityGroupRule{
					ID:        "sgr_IppA889KHFY1BvDt",
					Direction: "inbound",
				},
			},
			wantReq: &katapult.Request{
				Method: "GET",
				URL: &url.URL{
					Path: "/core/v1/security_groups/rules/_",
					RawQuery: url.Values{
						"security_group_rule[id]": []string{
							"sgr_IppA889KHFY1BvDt",
						},
					}.Encode(),
				},
			},
			want: &SecurityGroupRule{
				ID:        "sgr_IppA889KHFY1BvDt",
				Direction: "inbound",
			},
			wantResp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusOK},
			},
		},
		{
			name: "request error",
			args: args{
				ctx: context.Background(),
				ref: SecurityGroupRuleRef{ID: "sgr_IppA889KHFY1BvDt"},
			},
			resp: &katapult.Response{
				Response: &http.Response{
					StatusCode: http.StatusInternalServerError,
				},
			},
			respErr: fmt.Errorf("flux capacitor undercharged"),
			wantResp: &katapult.Response{
				Response: &http.Response{
					StatusCode: http.StatusInternalServerError,
				},
			},
			wantErr: "flux capacitor undercharged",
		},
		{
			name: "request error with nil response",
			args: args{
				ctx: context.Background(),
				ref: SecurityGroupRuleRef{ID: "sgr_IppA889KHFY1BvDt"},
			},
			resp:    nil,
			respErr: fmt.Errorf("someting is really wrong"),
			wantErr: "someting is really wrong",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := testclient.New(tt.resp, tt.respErr, tt.respV)
			c := NewSecurityGroupRulesClient(tc)
			ctx := test.Context(tt.args.ctx)

			got, resp, err := c.Get(ctx, tt.args.ref)

			assert.Equal(t, 1, len(tc.Calls), "only 1 request should be made")
			test.AssertContext(t, ctx, tc.Ctx)

			assert.Equal(t, tt.want, got)

			if tt.wantResp != nil {
				assert.Equal(t, tt.wantResp, resp)
			}

			if tt.wantReq != nil {
				assert.Equal(t, tt.wantReq, tc.Request)
			}

			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}

func TestSecurityGroupRulesClient_GetByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name     string
		args     args
		resp     *katapult.Response
		respErr  error
		respV    *securityGroupRulesResponseBody
		wantReq  *katapult.Request
		want     *SecurityGroupRule
		wantResp *katapult.Response
		wantErr  string
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				id:  "sgr_yOOeqZTbpLROjDGP",
			},
			resp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusOK},
			},
			respV: &securityGroupRulesResponseBody{
				SecurityGroupRule: &SecurityGroupRule{
					ID:        "sgr_yOOeqZTbpLROjDGP",
					Direction: "inbound",
				},
			},
			wantReq: &katapult.Request{
				Method: "GET",
				URL: &url.URL{
					Path: "/core/v1/security_groups/rules/_",
					RawQuery: url.Values{
						"security_group_rule[id]": []string{
							"sgr_yOOeqZTbpLROjDGP",
						},
					}.Encode(),
				},
			},
			want: &SecurityGroupRule{
				ID:        "sgr_yOOeqZTbpLROjDGP",
				Direction: "inbound",
			},
			wantResp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusOK},
			},
		},
		{
			name: "request error",
			args: args{
				ctx: context.Background(),
				id:  "sgr_yOOeqZTbpLROjDGP",
			},
			resp: &katapult.Response{
				Response: &http.Response{
					StatusCode: http.StatusInternalServerError,
				},
			},
			respErr: fmt.Errorf("flux capacitor undercharged"),
			wantResp: &katapult.Response{
				Response: &http.Response{
					StatusCode: http.StatusInternalServerError,
				},
			},
			wantErr: "flux capacitor undercharged",
		},
		{
			name: "request error with nil response",
			args: args{
				ctx: context.Background(),
				id:  "sgr_yOOeqZTbpLROjDGP",
			},
			resp:    nil,
			respErr: fmt.Errorf("someting is really wrong"),
			wantErr: "someting is really wrong",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := testclient.New(tt.resp, tt.respErr, tt.respV)
			c := NewSecurityGroupRulesClient(tc)
			ctx := test.Context(tt.args.ctx)

			got, resp, err := c.GetByID(ctx, tt.args.id)

			assert.Equal(t, 1, len(tc.Calls), "only 1 request should be made")
			test.AssertContext(t, ctx, tc.Ctx)

			assert.Equal(t, tt.want, got)

			if tt.wantResp != nil {
				assert.Equal(t, tt.wantResp, resp)
			}

			if tt.wantReq != nil {
				assert.Equal(t, tt.wantReq, tc.Request)
			}

			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}

func TestSecurityGroupRulesClient_Create(t *testing.T) {
	type args struct {
		ctx  context.Context
		ref  SecurityGroupRef
		args *SecurityGroupRuleArguments
	}
	tests := []struct {
		name     string
		args     args
		resp     *katapult.Response
		respErr  error
		respV    *securityGroupRulesResponseBody
		wantReq  *katapult.Request
		want     *SecurityGroupRule
		wantResp *katapult.Response
		wantErr  string
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				ref: SecurityGroupRef{ID: "sg_lNIuBAGHLCiz21jh"},
				args: &SecurityGroupRuleArguments{
					Direction: "inbound",
				},
			},
			resp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusOK},
			},
			respV: &securityGroupRulesResponseBody{
				SecurityGroupRule: &SecurityGroupRule{
					ID:        "sgr_RWe7qQRVygLUW9jT",
					Direction: "inbound",
				},
			},
			wantReq: &katapult.Request{
				Method: "POST",
				URL: &url.URL{
					Path: "/core/v1/security_groups/sg_lNIuBAGHLCiz21jh/rules",
				},
				Body: &securityGroupRuleCreateRequest{
					Properties: &SecurityGroupRuleArguments{
						Direction: "inbound",
					},
				},
			},
			want: &SecurityGroupRule{
				ID:        "sgr_RWe7qQRVygLUW9jT",
				Direction: "inbound",
			},
			wantResp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusOK},
			},
		},
		{
			name: "request error",
			args: args{
				ctx: context.Background(),
				ref: SecurityGroupRef{ID: "sg_lNIuBAGHLCiz21jh"},
			},
			resp: &katapult.Response{
				Response: &http.Response{
					StatusCode: http.StatusInternalServerError,
				},
			},
			respErr: fmt.Errorf("flux capacitor undercharged"),
			wantResp: &katapult.Response{
				Response: &http.Response{
					StatusCode: http.StatusInternalServerError,
				},
			},
			wantErr: "flux capacitor undercharged",
		},
		{
			name: "request error with nil response",
			args: args{
				ctx: context.Background(),
				ref: SecurityGroupRef{ID: "sg_lNIuBAGHLCiz21jh"},
			},
			resp:    nil,
			respErr: fmt.Errorf("someting is really wrong"),
			wantErr: "someting is really wrong",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := testclient.New(tt.resp, tt.respErr, tt.respV)
			c := NewSecurityGroupRulesClient(tc)
			ctx := test.Context(tt.args.ctx)

			got, resp, err := c.Create(ctx, tt.args.ref, tt.args.args)

			assert.Equal(t, 1, len(tc.Calls), "only 1 request should be made")
			test.AssertContext(t, ctx, tc.Ctx)

			assert.Equal(t, tt.want, got)

			if tt.wantResp != nil {
				assert.Equal(t, tt.wantResp, resp)
			}

			if tt.wantReq != nil {
				assert.Equal(t, tt.wantReq, tc.Request)
			}

			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}

func TestSecurityGroupRulesClient_Update(t *testing.T) {
	type args struct {
		ctx  context.Context
		ref  SecurityGroupRuleRef
		args *SecurityGroupRuleArguments
	}
	tests := []struct {
		name     string
		args     args
		resp     *katapult.Response
		respErr  error
		respV    *securityGroupRulesResponseBody
		wantReq  *katapult.Request
		want     *SecurityGroupRule
		wantResp *katapult.Response
		wantErr  string
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				ref: SecurityGroupRuleRef{ID: "sgr_FIxDC6SBnTBPtsIW"},
				args: &SecurityGroupRuleArguments{
					Direction: "inbound",
				},
			},
			resp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusOK},
			},
			respV: &securityGroupRulesResponseBody{
				SecurityGroupRule: &SecurityGroupRule{
					ID:        "sgr_FIxDC6SBnTBPtsIW",
					Direction: "inbound",
				},
			},
			wantReq: &katapult.Request{
				Method: "PATCH",
				URL: &url.URL{
					Path: "/core/v1/security_groups/rules/_",
					RawQuery: url.Values{
						"security_group_rule[id]": []string{
							"sgr_FIxDC6SBnTBPtsIW",
						},
					}.Encode(),
				},
				Body: &securityGroupRuleUpdateRequest{
					Properties: &SecurityGroupRuleArguments{
						Direction: "inbound",
					},
				},
			},
			want: &SecurityGroupRule{
				ID:        "sgr_FIxDC6SBnTBPtsIW",
				Direction: "inbound",
			},
			wantResp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusOK},
			},
		},
		{
			name: "request error",
			args: args{
				ctx: context.Background(),
				ref: SecurityGroupRuleRef{ID: "sgr_FIxDC6SBnTBPtsIW"},
			},
			resp: &katapult.Response{
				Response: &http.Response{
					StatusCode: http.StatusInternalServerError,
				},
			},
			respErr: fmt.Errorf("flux capacitor undercharged"),
			wantResp: &katapult.Response{
				Response: &http.Response{
					StatusCode: http.StatusInternalServerError,
				},
			},
			wantErr: "flux capacitor undercharged",
		},
		{
			name: "request error with nil response",
			args: args{
				ctx: context.Background(),
				ref: SecurityGroupRuleRef{ID: "sgr_FIxDC6SBnTBPtsIW"},
			},
			resp:    nil,
			respErr: fmt.Errorf("someting is really wrong"),
			wantErr: "someting is really wrong",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := testclient.New(tt.resp, tt.respErr, tt.respV)
			c := NewSecurityGroupRulesClient(tc)
			ctx := test.Context(tt.args.ctx)

			got, resp, err := c.Update(ctx, tt.args.ref, tt.args.args)

			assert.Equal(t, 1, len(tc.Calls), "only 1 request should be made")
			test.AssertContext(t, ctx, tc.Ctx)

			assert.Equal(t, tt.want, got)

			if tt.wantResp != nil {
				assert.Equal(t, tt.wantResp, resp)
			}

			if tt.wantReq != nil {
				assert.Equal(t, tt.wantReq, tc.Request)
			}

			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}

func TestSecurityGroupRulesClient_Delete2(t *testing.T) {
	type args struct {
		ctx context.Context
		ref SecurityGroupRuleRef
	}
	tests := []struct {
		name     string
		args     args
		resp     *katapult.Response
		respErr  error
		respV    *securityGroupRulesResponseBody
		wantReq  *katapult.Request
		want     *SecurityGroupRule
		wantResp *katapult.Response
		wantErr  string
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				ref: SecurityGroupRuleRef{ID: "sgr_htqQ6rLPQY0ljJTf"},
			},
			resp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusOK},
			},
			respV: &securityGroupRulesResponseBody{
				SecurityGroupRule: &SecurityGroupRule{
					ID: "sgr_htqQ6rLPQY0ljJTf",
				},
			},
			wantReq: &katapult.Request{
				Method: "DELETE",
				URL: &url.URL{
					Path: "/core/v1/security_groups/rules/_",
					RawQuery: url.Values{
						"security_group_rule[id]": []string{
							"sgr_htqQ6rLPQY0ljJTf",
						},
					}.Encode(),
				},
			},
			want: &SecurityGroupRule{
				ID: "sgr_htqQ6rLPQY0ljJTf",
			},
			wantResp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusOK},
			},
		},
		{
			name: "request error",
			args: args{
				ctx: context.Background(),
				ref: SecurityGroupRuleRef{ID: "sgr_htqQ6rLPQY0ljJTf"},
			},
			resp: &katapult.Response{
				Response: &http.Response{
					StatusCode: http.StatusInternalServerError,
				},
			},
			respErr: fmt.Errorf("flux capacitor undercharged"),
			wantResp: &katapult.Response{
				Response: &http.Response{
					StatusCode: http.StatusInternalServerError,
				},
			},
			wantErr: "flux capacitor undercharged",
		},
		{
			name: "request error with nil response",
			args: args{
				ctx: context.Background(),
				ref: SecurityGroupRuleRef{ID: "sgr_htqQ6rLPQY0ljJTf"},
			},
			resp:    nil,
			respErr: fmt.Errorf("someting is really wrong"),
			wantErr: "someting is really wrong",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := testclient.New(tt.resp, tt.respErr, tt.respV)
			c := NewSecurityGroupRulesClient(tc)
			ctx := test.Context(tt.args.ctx)

			got, resp, err := c.Delete(ctx, tt.args.ref)

			assert.Equal(t, 1, len(tc.Calls), "only 1 request should be made")
			test.AssertContext(t, ctx, tc.Ctx)

			assert.Equal(t, tt.want, got)

			if tt.wantResp != nil {
				assert.Equal(t, tt.wantResp, resp)
			}

			if tt.wantReq != nil {
				assert.Equal(t, tt.wantReq, tc.Request)
			}

			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}
