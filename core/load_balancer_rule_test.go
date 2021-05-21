package core

import (
	"context"
	"fmt"
	"net/url"
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
				ProxyProtocol:   boolPtr(false),
				Certificates: []Certificate{
					{
						ID: "another abitrary string",
					},
				},
				CheckEnabled:  boolPtr(false),
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

func TestLoadBalancerRulesClient_Get(t *testing.T) {
	type args struct {
		ref LoadBalancerRuleRef
	}
	tests := []struct {
		name    string
		frm     fakeRequestMakerArgs
		args    args
		want    *LoadBalancerRule
		wantErr string
	}{
		{
			name: "success",
			args: args{
				ref: LoadBalancerRuleRef{ID: "123"},
			},
			want: &LoadBalancerRule{
				ID:         "123",
				ListenPort: 132,
			},
			frm: fakeRequestMakerArgs{
				wantPath: "/core/v1/load_balancers/rules/_",
				wantValues: url.Values{
					"load_balancer_rule[id]": []string{"123"},
				},
				wantMethod: "GET",
				wantBody:   nil,
				doResponseBody: &loadBalancerRulesResponseBody{
					LoadBalancerRule: &LoadBalancerRule{
						ID:         "123",
						ListenPort: 132,
					},
				},
				doResp: &katapult.Response{},
			},
		},
		{
			name: "new request fails",
			args: args{
				ref: LoadBalancerRuleRef{ID: "123"},
			},
			frm: fakeRequestMakerArgs{
				wantPath: "/core/v1/load_balancers/rules/_",
				wantValues: url.Values{
					"load_balancer_rule[id]": []string{"123"},
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
				ref: LoadBalancerRuleRef{ID: "123"},
			},
			frm: fakeRequestMakerArgs{
				wantPath: "/core/v1/load_balancers/rules/_",
				wantValues: url.Values{
					"load_balancer_rule[id]": []string{"123"},
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
			c := NewLoadBalancerRulesClient(&fakeRequestMaker{
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

func TestLoadBalancerRulesClient_GetByID(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		frm     fakeRequestMakerArgs
		args    args
		want    *LoadBalancerRule
		wantErr string
	}{
		{
			name: "success",
			args: args{
				id: "123",
			},
			want: &LoadBalancerRule{
				ID:         "123",
				ListenPort: 132,
			},
			frm: fakeRequestMakerArgs{
				wantPath: "/core/v1/load_balancers/rules/_",
				wantValues: url.Values{
					"load_balancer_rule[id]": []string{"123"},
				},
				wantMethod: "GET",
				wantBody:   nil,
				doResponseBody: &loadBalancerRulesResponseBody{
					LoadBalancerRule: &LoadBalancerRule{
						ID:         "123",
						ListenPort: 132,
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
				wantPath: "/core/v1/load_balancers/rules/_",
				wantValues: url.Values{
					"load_balancer_rule[id]": []string{"123"},
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
				wantPath: "/core/v1/load_balancers/rules/_",
				wantValues: url.Values{
					"load_balancer_rule[id]": []string{"123"},
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
			c := NewLoadBalancerRulesClient(&fakeRequestMaker{
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

func TestLoadBalancerRulesClient_List(t *testing.T) {
	type args struct {
		loadBalancerID string
		listOptions    *ListOptions
	}
	tests := []struct {
		name     string
		frm      fakeRequestMakerArgs
		args     args
		want     []LoadBalancerRule
		wantResp *katapult.Response
		wantErr  string
	}{
		{
			name: "success",
			args: args{
				loadBalancerID: "xyzzy",
				listOptions: &ListOptions{
					Page:    5,
					PerPage: 32,
				},
			},
			want: []LoadBalancerRule{{
				ID:              "abc",
				DestinationPort: 666,
			}},
			wantResp: &katapult.Response{
				Pagination: &katapult.Pagination{Total: 333},
			},
			frm: fakeRequestMakerArgs{
				wantPath:   "/core/v1/load_balancers/_/rules",
				wantMethod: "GET",
				wantBody:   nil,
				wantValues: url.Values{
					"page":              []string{"5"},
					"per_page":          []string{"32"},
					"load_balancer[id]": []string{"xyzzy"},
				},
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
			name: "success with nil options",
			args: args{
				loadBalancerID: "xyzzy",
				listOptions:    nil,
			},

			want: []LoadBalancerRule{{
				ID: "cbd",
			}},
			wantResp: &katapult.Response{},
			frm: fakeRequestMakerArgs{
				wantPath:   "/core/v1/load_balancers/_/rules",
				wantMethod: "GET",
				wantValues: url.Values{
					"load_balancer[id]": []string{"xyzzy"},
				},
				wantBody: nil,
				doResponseBody: &loadBalancerRulesResponseBody{
					LoadBalancerRules: []LoadBalancerRule{
						{ID: "cbd"},
					},
				},
				doResp: &katapult.Response{},
			},
		},
		{
			name: "new request fails",
			args: args{
				loadBalancerID: "xyzzy",
			},
			frm: fakeRequestMakerArgs{
				wantPath:   "/core/v1/load_balancers/_/rules",
				wantMethod: "GET",
				wantValues: url.Values{
					"load_balancer[id]": []string{"xyzzy"},
				},
				wantBody:  nil,
				newReqErr: fmt.Errorf("rats chewed cables"),
			},
			wantErr: "rats chewed cables",
		},
		{
			name: "http do fails",
			args: args{
				loadBalancerID: "xyzzy",
			},
			frm: fakeRequestMakerArgs{
				wantPath:   "/core/v1/load_balancers/_/rules",
				wantMethod: "GET",
				wantValues: url.Values{
					"load_balancer[id]": []string{"xyzzy"},
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
			c := NewLoadBalancerRulesClient(&fakeRequestMaker{
				t:    t,
				args: tt.frm,
			})

			got, resp, err := c.List(
				context.Background(),
				LoadBalancerRef{ID: tt.args.loadBalancerID},
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

func TestLoadBalancerRulesClient_Create(t *testing.T) {
	type args struct {
		loadBalancerID string
		creationArgs   LoadBalancerRuleArguments
	}
	tests := []struct {
		name    string
		frm     fakeRequestMakerArgs
		args    args
		want    *LoadBalancerRule
		wantErr string
	}{
		{
			name: "success",
			args: args{
				loadBalancerID: "xyzzy",
				creationArgs:   LoadBalancerRuleArguments{DestinationPort: 666},
			},
			want: &LoadBalancerRule{
				ID:              "abc",
				DestinationPort: 666,
			},
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
			name: "new request fails",
			args: args{
				loadBalancerID: "xyzzy",
			},
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
			name: "http do fails",
			args: args{
				loadBalancerID: "xyzzy",
			},
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
				LoadBalancerRef{ID: tt.args.loadBalancerID},
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

func TestLoadBalancerRulesClient_Update(t *testing.T) {
	type args struct {
		ruleID     string
		updateArgs LoadBalancerRuleArguments
	}
	tests := []struct {
		name    string
		frm     fakeRequestMakerArgs
		args    args
		want    *LoadBalancerRule
		wantErr string
	}{
		{
			name: "success",
			want: &LoadBalancerRule{
				ID:              "abc",
				DestinationPort: 666,
			},
			args: args{
				updateArgs: LoadBalancerRuleArguments{DestinationPort: 666},
				ruleID:     "123",
			},
			frm: fakeRequestMakerArgs{
				wantPath: "/core/v1/load_balancers/rules/_",
				wantValues: url.Values{
					"load_balancer_rule[id]": []string{"123"},
				},
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
			name: "new request fails",
			args: args{
				ruleID: "123",
			},
			frm: fakeRequestMakerArgs{
				wantPath: "/core/v1/load_balancers/rules/_",
				wantValues: url.Values{
					"load_balancer_rule[id]": []string{"123"},
				},
				wantMethod: "PATCH",
				wantBody: &loadBalancerRuleUpdateRequest{
					Properties: LoadBalancerRuleArguments{},
				},
				newReqErr: fmt.Errorf("rats chewed cables"),
			},
			wantErr: "rats chewed cables",
		},
		{
			name: "http do fails",
			args: args{
				ruleID: "123",
			},
			frm: fakeRequestMakerArgs{
				wantPath: "/core/v1/load_balancers/rules/_",
				wantValues: url.Values{
					"load_balancer_rule[id]": []string{"123"},
				},
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
				LoadBalancerRuleRef{ID: tt.args.ruleID},
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
				wantPath: "/core/v1/load_balancers/rules/_",
				wantValues: url.Values{
					"load_balancer_rule[id]": []string{"123"},
				},
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
				wantPath: "/core/v1/load_balancers/rules/_",
				wantValues: url.Values{
					"load_balancer_rule[id]": []string{"123"},
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
				wantPath: "/core/v1/load_balancers/rules/_",
				wantValues: url.Values{
					"load_balancer_rule[id]": []string{"123"},
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
			c := NewLoadBalancerRulesClient(&fakeRequestMaker{
				t:    t,
				args: tt.frm,
			})

			got, resp, err := c.Delete(
				context.Background(),
				LoadBalancerRuleRef{ID: tt.args.ruleID},
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
