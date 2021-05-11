package core

import (
	"context"
	"fmt"
	"github.com/krystal/go-katapult"
	"github.com/stretchr/testify/assert"
	"testing"
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
				ID:              "abritrary string",
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
				ProxyProtocol:   true,
				Certificates: []Certificate{
					{
						ID: "another abitrary string",
					},
				},
				CheckEnabled:  true,
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

		want         *LoadBalancerRule
		wantErr      string
		wantResponse bool
	}{
		{
			name: "success",
			args: args{
				ruleID: "123",
			},
			want: &lbr,
			frm: fakeRequestMakerArgs{
				wantPath:       "/core/v1/load_balancers/rules/123",
				wantMethod:     "DELETE",
				wantBody:       nil,
				doResponseBody: &loadBalancerRulesResponseBody{LoadBalancerRule: &lbr},
				doResp:         &katapult.Response{},
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
			c := NewLoadBalancerRulesClient(&fakeRequestMaker{t: t, args: tt.frm})

			got, resp, err := c.Delete(context.Background(), tt.args.ruleID)
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
