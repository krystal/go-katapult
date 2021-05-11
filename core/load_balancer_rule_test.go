package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClient_LoadBalancerRules(t *testing.T) {
	c := New(&fakeRequestMaker{})

	assert.IsType(t, &LoadBalancerRulesClient{}, c.LoadBalancerRules)
}

func TestLoadBalancerRule_JSONMarshaling(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
			testJSONMarshaling(t, tt.obj)
		})
	}
}
