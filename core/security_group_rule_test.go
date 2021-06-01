package core

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"github.com/krystal/go-katapult"
	"github.com/stretchr/testify/assert"
)

func TestClient_SecurityGroupRules(t *testing.T) {
	c := New(&fakeRequestMaker{})

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
				Ports:     "3000",
				Targets:   []string{"192.168.0.1"},
				Notes:     "My fave security group",
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
				Ports:     "3000",
				Targets:   []string{"192.168.0.1"},
				Notes:     "My fave security group",
			},
		},
		{
			name: "remove all targets",
			obj: &SecurityGroupRuleArguments{
				Targets: []string(nil),
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
		obj  *SecurityGroupRulesResponseBody
	}{
		{
			name: "empty",
			obj:  &SecurityGroupRulesResponseBody{},
		},
		{
			name: "full",
			obj: &SecurityGroupRulesResponseBody{
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
					Ports:     "3000",
					Targets:   []string{"192.168.0.1"},
					Notes:     "My fave security group",
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
					Ports:     "3000",
					Targets:   []string{"192.168.0.1"},
					Notes:     "My fave security group",
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
		SecurityGroupID string
		listOptions     *ListOptions
	}
	tests := []struct {
		name     string
		frm      fakeRequestMakerArgs
		args     args
		want     []SecurityGroupRule
		wantResp *katapult.Response
		wantErr  string
	}{
		{
			name: "success",
			args: args{
				SecurityGroupID: "xyzzy",
				listOptions: &ListOptions{
					Page:    5,
					PerPage: 32,
				},
			},
			want: []SecurityGroupRule{{
				ID:        "abc",
				Direction: "inbound",
			}},
			wantResp: &katapult.Response{
				Pagination: &katapult.Pagination{Total: 333},
			},
			frm: fakeRequestMakerArgs{
				wantPath:   "/core/v1/security_groups/_/rules",
				wantMethod: "GET",
				wantBody:   nil,
				wantValues: url.Values{
					"page":               []string{"5"},
					"per_page":           []string{"32"},
					"security_group[id]": []string{"xyzzy"},
				},
				doResponseBody: &SecurityGroupRulesResponseBody{
					SecurityGroupRules: []SecurityGroupRule{
						{ID: "abc", Direction: "inbound"},
					},
					Pagination: &katapult.Pagination{Total: 333},
				},
				doResp: &katapult.Response{},
			},
		},
		{
			name: "success with nil options",
			args: args{
				SecurityGroupID: "xyzzy",
				listOptions:     nil,
			},

			want: []SecurityGroupRule{{
				ID: "cbd",
			}},
			wantResp: &katapult.Response{},
			frm: fakeRequestMakerArgs{
				wantPath:   "/core/v1/security_groups/_/rules",
				wantMethod: "GET",
				wantValues: url.Values{
					"security_group[id]": []string{"xyzzy"},
				},
				wantBody: nil,
				doResponseBody: &SecurityGroupRulesResponseBody{
					SecurityGroupRules: []SecurityGroupRule{
						{ID: "cbd"},
					},
				},
				doResp: &katapult.Response{},
			},
		},
		{
			name: "new request fails",
			args: args{
				SecurityGroupID: "xyzzy",
			},
			frm: fakeRequestMakerArgs{
				wantPath:   "/core/v1/security_groups/_/rules",
				wantMethod: "GET",
				wantValues: url.Values{
					"security_group[id]": []string{"xyzzy"},
				},
				wantBody:  nil,
				newReqErr: fmt.Errorf("rats chewed cables"),
			},
			wantErr: "rats chewed cables",
		},
		{
			name: "http do fails",
			args: args{
				SecurityGroupID: "xyzzy",
			},
			frm: fakeRequestMakerArgs{
				wantPath:   "/core/v1/security_groups/_/rules",
				wantMethod: "GET",
				wantValues: url.Values{
					"security_group[id]": []string{"xyzzy"},
				},
				wantBody: nil,
				doErr:    fmt.Errorf("flux capacitor undercharged"),
				doResp:   &katapult.Response{},
			},
			wantResp: &katapult.Response{},
			wantErr:  "flux capacitor undercharged",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewSecurityGroupRulesClient(&fakeRequestMaker{
				t:    t,
				args: tt.frm,
			})

			got, resp, err := c.List(
				context.Background(),
				SecurityGroupRef{ID: tt.args.SecurityGroupID},
				tt.args.listOptions,
			)
			assert.Equal(t, tt.wantResp, resp)

			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSecurityGroupRulesClient_Get(t *testing.T) {
	type args struct {
		ref SecurityGroupRuleRef
	}
	tests := []struct {
		name    string
		frm     fakeRequestMakerArgs
		args    args
		want    *SecurityGroupRule
		wantErr string
	}{
		{
			name: "success",
			args: args{
				ref: SecurityGroupRuleRef{ID: "123"},
			},
			want: &SecurityGroupRule{
				ID:        "123",
				Direction: "inbound",
			},
			frm: fakeRequestMakerArgs{
				wantPath: "/core/v1/security_groups/rules/_",
				wantValues: url.Values{
					"security_group_rule[id]": []string{"123"},
				},
				wantMethod: "GET",
				wantBody:   nil,
				doResponseBody: &SecurityGroupRulesResponseBody{
					SecurityGroupRule: &SecurityGroupRule{
						ID:        "123",
						Direction: "inbound",
					},
				},
				doResp: &katapult.Response{},
			},
		},
		{
			name: "new request fails",
			args: args{
				ref: SecurityGroupRuleRef{ID: "123"},
			},
			frm: fakeRequestMakerArgs{
				wantPath: "/core/v1/security_groups/rules/_",
				wantValues: url.Values{
					"security_group_rule[id]": []string{"123"},
				},
				wantMethod: "GET",
				wantBody:   nil,
				newReqErr:  fmt.Errorf("rats chewed cables"),
			},
			wantErr: "rats chewed cables",
		},
		{
			name: "http do fails",
			args: args{
				ref: SecurityGroupRuleRef{ID: "123"},
			},
			frm: fakeRequestMakerArgs{
				wantPath: "/core/v1/security_groups/rules/_",
				wantValues: url.Values{
					"security_group_rule[id]": []string{"123"},
				},
				wantMethod: "GET",
				wantBody:   nil,
				doErr:      fmt.Errorf("flux capacitor undercharged"),
				doResp:     &katapult.Response{},
			},
			wantErr: "flux capacitor undercharged",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewSecurityGroupRulesClient(&fakeRequestMaker{
				t:    t,
				args: tt.frm,
			})

			got, resp, err := c.Get(
				context.Background(),
				tt.args.ref,
			)
			assert.Equal(t, tt.frm.doResp, resp)

			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSecurityGroupRulesClient_GetByID(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		frm     fakeRequestMakerArgs
		args    args
		want    *SecurityGroupRule
		wantErr string
	}{
		{
			name: "success",
			args: args{
				id: "123",
			},
			want: &SecurityGroupRule{
				ID:        "123",
				Direction: "inbound",
			},
			frm: fakeRequestMakerArgs{
				wantPath: "/core/v1/security_groups/rules/_",
				wantValues: url.Values{
					"security_group_rule[id]": []string{"123"},
				},
				wantMethod: "GET",
				wantBody:   nil,
				doResponseBody: &SecurityGroupRulesResponseBody{
					SecurityGroupRule: &SecurityGroupRule{
						ID:        "123",
						Direction: "inbound",
					},
				},
				doResp: &katapult.Response{},
			},
		},
		{
			name: "new request fails",
			args: args{
				id: "123",
			},
			frm: fakeRequestMakerArgs{
				wantPath: "/core/v1/security_groups/rules/_",
				wantValues: url.Values{
					"security_group_rule[id]": []string{"123"},
				},
				wantMethod: "GET",
				wantBody:   nil,
				newReqErr:  fmt.Errorf("rats chewed cables"),
			},
			wantErr: "rats chewed cables",
		},
		{
			name: "http do fails",
			args: args{
				id: "123",
			},
			frm: fakeRequestMakerArgs{
				wantPath: "/core/v1/security_groups/rules/_",
				wantValues: url.Values{
					"security_group_rule[id]": []string{"123"},
				},
				wantMethod: "GET",
				wantBody:   nil,
				doErr:      fmt.Errorf("flux capacitor undercharged"),
				doResp:     &katapult.Response{},
			},
			wantErr: "flux capacitor undercharged",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewSecurityGroupRulesClient(&fakeRequestMaker{
				t:    t,
				args: tt.frm,
			})

			got, resp, err := c.GetByID(
				context.Background(),
				tt.args.id,
			)
			assert.Equal(t, tt.frm.doResp, resp)

			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSecurityGroupRulesClient_Create(t *testing.T) {
	type args struct {
		SecurityGroupID string
		creationArgs    *SecurityGroupRuleArguments
	}
	tests := []struct {
		name    string
		frm     fakeRequestMakerArgs
		args    args
		want    *SecurityGroupRule
		wantErr string
	}{
		{
			name: "success",
			args: args{
				SecurityGroupID: "xyzzy",
				creationArgs: &SecurityGroupRuleArguments{
					Direction: "inbound",
				},
			},
			want: &SecurityGroupRule{
				ID:        "abc",
				Direction: "inbound",
			},
			frm: fakeRequestMakerArgs{
				wantPath:   "/core/v1/security_groups/xyzzy/rules",
				wantMethod: "POST",
				wantBody: &securityGroupRuleCreateRequest{
					Properties: &SecurityGroupRuleArguments{
						Direction: "inbound",
					},
				},
				doResponseBody: &SecurityGroupRulesResponseBody{
					SecurityGroupRule: &SecurityGroupRule{
						ID:        "abc",
						Direction: "inbound",
					},
				},
				doResp: &katapult.Response{},
			},
		},
		{
			name: "new request fails",
			args: args{
				SecurityGroupID: "xyzzy",
			},
			frm: fakeRequestMakerArgs{
				wantPath:   "/core/v1/security_groups/xyzzy/rules",
				wantMethod: "POST",
				wantBody:   &securityGroupRuleCreateRequest{},
				newReqErr:  fmt.Errorf("rats chewed cables"),
			},
			wantErr: "rats chewed cables",
		},
		{
			name: "http do fails",
			args: args{
				SecurityGroupID: "xyzzy",
			},
			frm: fakeRequestMakerArgs{
				wantPath:   "/core/v1/security_groups/xyzzy/rules",
				wantMethod: "POST",
				wantBody:   &securityGroupRuleCreateRequest{},
				doErr:      fmt.Errorf("flux capacitor undercharged"),
				doResp:     &katapult.Response{},
			},
			wantErr: "flux capacitor undercharged",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewSecurityGroupRulesClient(&fakeRequestMaker{
				t:    t,
				args: tt.frm,
			})

			got, resp, err := c.Create(
				context.Background(),
				SecurityGroupRef{ID: tt.args.SecurityGroupID},
				tt.args.creationArgs,
			)
			assert.Equal(t, tt.frm.doResp, resp)

			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSecurityGroupRulesClient_Update(t *testing.T) {
	type args struct {
		ruleID     string
		updateArgs *SecurityGroupRuleArguments
	}
	tests := []struct {
		name    string
		frm     fakeRequestMakerArgs
		args    args
		want    *SecurityGroupRule
		wantErr string
	}{
		{
			name: "success",
			want: &SecurityGroupRule{
				ID:        "abc",
				Direction: "inbound",
			},
			args: args{
				updateArgs: &SecurityGroupRuleArguments{Direction: "inbound"},
				ruleID:     "123",
			},
			frm: fakeRequestMakerArgs{
				wantPath: "/core/v1/security_groups/rules/_",
				wantValues: url.Values{
					"security_group_rule[id]": []string{"123"},
				},
				wantMethod: "PATCH",
				wantBody: &securityGroupRuleUpdateRequest{
					Properties: &SecurityGroupRuleArguments{
						Direction: "inbound",
					},
				},
				doResponseBody: &SecurityGroupRulesResponseBody{
					SecurityGroupRule: &SecurityGroupRule{
						ID:        "abc",
						Direction: "inbound",
					},
				},
				doResp: &katapult.Response{},
			},
		},
		{
			name: "new request fails",
			args: args{
				ruleID: "123",
			},
			frm: fakeRequestMakerArgs{
				wantPath: "/core/v1/security_groups/rules/_",
				wantValues: url.Values{
					"security_group_rule[id]": []string{"123"},
				},
				wantMethod: "PATCH",
				wantBody:   &securityGroupRuleUpdateRequest{},
				newReqErr:  fmt.Errorf("rats chewed cables"),
			},
			wantErr: "rats chewed cables",
		},
		{
			name: "http do fails",
			args: args{
				ruleID: "123",
			},
			frm: fakeRequestMakerArgs{
				wantPath: "/core/v1/security_groups/rules/_",
				wantValues: url.Values{
					"security_group_rule[id]": []string{"123"},
				},
				wantMethod: "PATCH",
				wantBody:   &securityGroupRuleUpdateRequest{},
				doErr:      fmt.Errorf("flux capacitor undercharged"),
				doResp:     &katapult.Response{},
			},
			wantErr: "flux capacitor undercharged",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewSecurityGroupRulesClient(&fakeRequestMaker{
				t:    t,
				args: tt.frm,
			})

			got, resp, err := c.Update(
				context.Background(),
				SecurityGroupRuleRef{ID: tt.args.ruleID},
				tt.args.updateArgs,
			)
			assert.Equal(t, tt.frm.doResp, resp)

			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSecurityGroupRulesClient_Delete(t *testing.T) {
	type args struct {
		ruleID string
	}
	sgr := SecurityGroupRule{
		ID:        "abc",
		Direction: "inbound",
	}
	tests := []struct {
		name string
		args args
		frm  fakeRequestMakerArgs

		want    *SecurityGroupRule
		wantErr string
	}{
		{
			name: "success",
			args: args{
				ruleID: "123",
			},
			want: &sgr,
			frm: fakeRequestMakerArgs{
				wantPath: "/core/v1/security_groups/rules/_",
				wantValues: url.Values{
					"security_group_rule[id]": []string{"123"},
				},
				wantMethod: "DELETE",
				wantBody:   nil,
				doResponseBody: &SecurityGroupRulesResponseBody{
					SecurityGroupRule: &sgr,
				},
				doResp: &katapult.Response{},
			},
		},
		{
			name: "new request fails",
			args: args{
				ruleID: "123",
			},
			frm: fakeRequestMakerArgs{
				wantPath: "/core/v1/security_groups/rules/_",
				wantValues: url.Values{
					"security_group_rule[id]": []string{"123"},
				},
				wantMethod: "DELETE",
				newReqErr:  fmt.Errorf("rats chewed cables"),
			},
			wantErr: "rats chewed cables",
		},
		{
			name: "http do fails",
			args: args{
				ruleID: "123",
			},
			frm: fakeRequestMakerArgs{
				wantPath: "/core/v1/security_groups/rules/_",
				wantValues: url.Values{
					"security_group_rule[id]": []string{"123"},
				},
				wantMethod: "DELETE",
				wantBody:   nil,
				doErr:      fmt.Errorf("flux capacitor undercharged"),
				doResp:     &katapult.Response{},
			},
			wantErr: "flux capacitor undercharged",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewSecurityGroupRulesClient(&fakeRequestMaker{
				t:    t,
				args: tt.frm,
			})

			got, resp, err := c.Delete(
				context.Background(),
				SecurityGroupRuleRef{ID: tt.args.ruleID},
			)
			assert.Equal(t, tt.frm.doResp, resp)

			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
