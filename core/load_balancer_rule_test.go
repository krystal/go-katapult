package core

import (
	"context"
	"fmt"
	"testing"

	"github.com/krystal/go-katapult"
	"github.com/stretchr/testify/assert"
)

func TestClient_LoadBalancerRules(t *testing.T) {
	c := New(&fakeRequestMaker{})

	assert.IsType(t, &LoadBalancerRulesClient{}, c.LoadBalancerRules)
}

func TestLoadBalancerRule_JSONMarshalling(t *testing.T) {
	tests := []struct {
		name string
		obj  *LoadBalancerRule
	}{
		{
			name: "empty",
			obj:  &LoadBalancerRule{},
		},
		{
			name: "full",
			obj: &LoadBalancerRule{
				ID:              "arbitrary string",
				Algorithm:       StickyRuleAlgorithm,
				DestinationPort: 1024,
				ListenPort:      1337,
				Protocol:        HTTPProtocol,
				ProxyProtocol:   true,
				Certificates: []Certificate{
					{
						ID: "another abitrary string",
					},
				},
				BackendSSL:     true,
				PassthroughSSL: true,
				CheckEnabled:   true,
				CheckFall:      3,
				CheckInterval:  50,
				CheckPath:      "/healthz",
				CheckProtocol:  HTTPProtocol,
				CheckRise:      12,
				CheckTimeout:   3,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestLoadBalancerRuleArguments_JSONMarshalling(t *testing.T) {
	falsy := false
	tests := []struct {
		name string
		obj  *LoadBalancerRuleArguments
	}{
		{
			name: "empty",
			obj:  &LoadBalancerRuleArguments{},
		},
		{
			name: "full",
			obj: &LoadBalancerRuleArguments{
				Algorithm:       StickyRuleAlgorithm,
				DestinationPort: 1024,
				ListenPort:      1337,
				Protocol:        HTTPProtocol,
				ProxyProtocol:   &falsy,
				Certificates: []Certificate{
					{
						ID: "another abitrary string",
					},
				},
				CheckEnabled:  &falsy,
				CheckFall:     3,
				CheckInterval: 50,
				CheckPath:     "/healthz",
				CheckProtocol: HTTPProtocol,
				CheckRise:     12,
				CheckTimeout:  3,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func Test_loadBalancerRulesResponseBody_JSONMarshalling(t *testing.T) {
	tests := []struct {
		name string
		obj  *loadBalancerRulesResponseBody
	}{
		{
			name: "empty",
			obj:  &loadBalancerRulesResponseBody{},
		},
		{
			name: "full",
			obj: &loadBalancerRulesResponseBody{
				Pagination: &katapult.Pagination{
					LargeSet: true,
				},
				LoadBalancerRule: &LoadBalancerRule{ID: "foobar"},
				LoadBalancerRules: []LoadBalancerRule{
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

func Test_loadBalancerRuleCreateRequest_JSONMarshalling(t *testing.T) {
	tests := []struct {
		name string
		obj  *loadBalancerRuleCreateRequest
	}{
		{
			name: "empty",
			obj:  &loadBalancerRuleCreateRequest{},
		},
		{
			name: "full",
			obj: &loadBalancerRuleCreateRequest{
				Properties: LoadBalancerRuleArguments{
					Protocol: HTTPProtocol,
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

func Test_loadBalancerRuleUpdateRequest_JSONMarshalling(t *testing.T) {
	tests := []struct {
		name string
		obj  *loadBalancerRuleUpdateRequest
	}{
		{
			name: "empty",
			obj:  &loadBalancerRuleUpdateRequest{},
		},
		{
			name: "full",
			obj: &loadBalancerRuleUpdateRequest{
				Properties: LoadBalancerRuleArguments{
					Protocol: TCPProtocol,
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
func TestLoadBalancerRulesClient_List(t *testing.T) {
	tests := []struct {
		name string
		frm  fakeRequestMakerArgs

		loadBalancerID string
		listOptions    *ListOptions

		want     []LoadBalancerRule
		wantResp *katapult.Response
		wantErr  string
	}{
		{
			name:           "success",
			loadBalancerID: "xyzzy",
			want: []LoadBalancerRule{{
				ID:              "abc",
				DestinationPort: 666,
			}},
			wantResp: &katapult.Response{
				Pagination: &katapult.Pagination{Total: 333},
			},
			listOptions: &ListOptions{
				Page:    5,
				PerPage: 32,
			},
			frm: fakeRequestMakerArgs{
				wantPath:   "/core/v1/load_balancers/xyzzy/rules?page=5&per_page=32", //nolint:lll
				wantMethod: "GET",
				wantBody:   nil,
				doResponseBody: &loadBalancerRulesResponseBody{
					LoadBalancerRules: []LoadBalancerRule{
						{ID: "abc", DestinationPort: 666},
					},
					Pagination: &katapult.Pagination{Total: 333},
				},
				doResp: &katapult.Response{},
			},
		},
		{
			name:           "success with nil options",
			loadBalancerID: "xyzzy",
			want: []LoadBalancerRule{{
				ID: "cbd",
			}},
			wantResp:    &katapult.Response{},
			listOptions: nil,
			frm: fakeRequestMakerArgs{
				wantPath:   "/core/v1/load_balancers/xyzzy/rules",
				wantMethod: "GET",
				wantBody:   nil,
				doResponseBody: &loadBalancerRulesResponseBody{
					LoadBalancerRules: []LoadBalancerRule{
						{ID: "cbd"},
					},
				},
				doResp: &katapult.Response{},
			},
		},
		{
			name:           "new request fails",
			loadBalancerID: "xyzzy",
			frm: fakeRequestMakerArgs{
				wantPath:   "/core/v1/load_balancers/xyzzy/rules",
				wantMethod: "GET",
				wantBody:   nil,
				newReqErr:  fmt.Errorf("rats chewed cables"),
			},
			wantErr: "rats chewed cables",
		},
		{
			name:           "http do fails",
			loadBalancerID: "xyzzy",
			frm: fakeRequestMakerArgs{
				wantPath:   "/core/v1/load_balancers/xyzzy/rules",
				wantMethod: "GET",
				wantBody:   nil,
				doErr:      fmt.Errorf("flux capacitor undercharged"),
				doResp:     &katapult.Response{},
			},
			wantResp: &katapult.Response{},
			wantErr:  "flux capacitor undercharged",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewLoadBalancerRulesClient(&fakeRequestMaker{
				t:    t,
				args: tt.frm,
			})

			got, resp, err := c.List(
				context.Background(),
				&LoadBalancer{ID: tt.loadBalancerID},
				tt.listOptions,
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

func TestLoadBalancerRulesClient_Create(t *testing.T) {
	tests := []struct {
		name string
		frm  fakeRequestMakerArgs

		loadBalancerID string
		args           LoadBalancerRuleArguments

		want    *LoadBalancerRule
		wantErr string
	}{
		{
			name:           "success",
			loadBalancerID: "xyzzy",
			want: &LoadBalancerRule{
				ID:              "abc",
				DestinationPort: 666,
			},
			args: LoadBalancerRuleArguments{DestinationPort: 666},
			frm: fakeRequestMakerArgs{
				wantPath:   "/core/v1/load_balancers/xyzzy/rules",
				wantMethod: "POST",
				wantBody: &loadBalancerRuleCreateRequest{
					Properties: LoadBalancerRuleArguments{
						DestinationPort: 666,
					},
				},
				doResponseBody: &loadBalancerRulesResponseBody{
					LoadBalancerRule: &LoadBalancerRule{
						ID:              "abc",
						DestinationPort: 666,
					},
				},
				doResp: &katapult.Response{},
			},
		},
		{
			name:           "new request fails",
			loadBalancerID: "xyzzy",
			frm: fakeRequestMakerArgs{
				wantPath:   "/core/v1/load_balancers/xyzzy/rules",
				wantMethod: "POST",
				wantBody: &loadBalancerRuleCreateRequest{
					Properties: LoadBalancerRuleArguments{},
				},
				newReqErr: fmt.Errorf("rats chewed cables"),
			},
			wantErr: "rats chewed cables",
		},
		{
			name:           "http do fails",
			loadBalancerID: "xyzzy",
			frm: fakeRequestMakerArgs{
				wantPath:   "/core/v1/load_balancers/xyzzy/rules",
				wantMethod: "POST",
				wantBody: &loadBalancerRuleCreateRequest{
					Properties: LoadBalancerRuleArguments{},
				},
				doErr:  fmt.Errorf("flux capacitor undercharged"),
				doResp: &katapult.Response{},
			},
			wantErr: "flux capacitor undercharged",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewLoadBalancerRulesClient(&fakeRequestMaker{
				t:    t,
				args: tt.frm,
			})

			got, resp, err := c.Create(
				context.Background(),
				&LoadBalancer{ID: tt.loadBalancerID},
				tt.args,
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

func TestLoadBalancerRulesClient_Update(t *testing.T) {
	tests := []struct {
		name string
		frm  fakeRequestMakerArgs

		ruleID string
		args   LoadBalancerRuleArguments

		want    *LoadBalancerRule
		wantErr string
	}{
		{
			name:   "success",
			ruleID: "123",
			want: &LoadBalancerRule{
				ID:              "abc",
				DestinationPort: 666,
			},
			args: LoadBalancerRuleArguments{DestinationPort: 666},
			frm: fakeRequestMakerArgs{
				wantPath:   "/core/v1/load_balancers/rules/123",
				wantMethod: "PATCH",
				wantBody: &loadBalancerRuleUpdateRequest{
					Properties: LoadBalancerRuleArguments{
						DestinationPort: 666,
					},
				},
				doResponseBody: &loadBalancerRulesResponseBody{
					LoadBalancerRule: &LoadBalancerRule{
						ID:              "abc",
						DestinationPort: 666,
					},
				},
				doResp: &katapult.Response{},
			},
		},
		{
			name:   "new request fails",
			ruleID: "123",
			frm: fakeRequestMakerArgs{
				wantPath:   "/core/v1/load_balancers/rules/123",
				wantMethod: "PATCH",
				wantBody: &loadBalancerRuleUpdateRequest{
					Properties: LoadBalancerRuleArguments{},
				},
				newReqErr: fmt.Errorf("rats chewed cables"),
			},
			wantErr: "rats chewed cables",
		},
		{
			name:   "http do fails",
			ruleID: "123",
			frm: fakeRequestMakerArgs{
				wantPath:   "/core/v1/load_balancers/rules/123",
				wantMethod: "PATCH",
				wantBody: &loadBalancerRuleUpdateRequest{
					Properties: LoadBalancerRuleArguments{},
				},
				doErr:  fmt.Errorf("flux capacitor undercharged"),
				doResp: &katapult.Response{},
			},
			wantErr: "flux capacitor undercharged",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewLoadBalancerRulesClient(&fakeRequestMaker{
				t:    t,
				args: tt.frm,
			})

			got, resp, err := c.Update(
				context.Background(),
				&LoadBalancerRule{ID: tt.ruleID},
				tt.args,
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

func TestLoadBalancerRulesClient_Delete(t *testing.T) {
	type args struct {
		ruleID string
	}
	lbr := LoadBalancerRule{
		ID:              "abc",
		DestinationPort: 55,
	}
	tests := []struct {
		name string
		args args
		frm  fakeRequestMakerArgs

		want    *LoadBalancerRule
		wantErr string
	}{
		{
			name: "success",
			args: args{
				ruleID: "123",
			},
			want: &lbr,
			frm: fakeRequestMakerArgs{
				wantPath:   "/core/v1/load_balancers/rules/123",
				wantMethod: "DELETE",
				wantBody:   nil,
				doResponseBody: &loadBalancerRulesResponseBody{
					LoadBalancerRule: &lbr,
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
				wantPath:   "/core/v1/load_balancers/rules/123",
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
				wantPath:   "/core/v1/load_balancers/rules/123",
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
			c := NewLoadBalancerRulesClient(&fakeRequestMaker{
				t:    t,
				args: tt.frm,
			})

			got, resp, err := c.Delete(
				context.Background(),
				&LoadBalancerRule{ID: tt.args.ruleID},
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
