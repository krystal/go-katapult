package core

import (
	"context"
	"fmt"
	"github.com/krystal/go-katapult"
	"github.com/stretchr/testify/assert"
	"net/http"
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
	tests := []struct {
		name       string
		args       args
		want       *LoadBalancerRule
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "by ID",
			args: args{
				ruleID: "123",
			},
			want: &LoadBalancerRule{
				ID: "123",
			},
			respStatus: http.StatusOK,
			respBody:   []byte(`{"load_balancer_rule": {"ID": "123"}}`),
		},
		{
			name: "non-existent load balancer rule",
			args: args{
				ruleID: "123",
			},
			errStr:     fixtureLoadBalancerNotFoundErr,
			errResp:    fixtureLoadBalancerNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("load_balancer_not_found_error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewLoadBalancerRulesClient(rm)

			mux.HandleFunc(
				fmt.Sprintf("/core/v1/load_balancers_rules/%s", tt.args.ruleID),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "DELETE", r.Method)

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.Delete(context.Background(), tt.args.ruleID)
			assert.Equal(t, tt.respStatus, resp.StatusCode)

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
