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

func TestClient_LoadBalancerRules(t *testing.T) {
	c := New(&testclient.Client{})

	assert.IsType(t, &LoadBalancerRulesClient{}, c.LoadBalancerRules)
}

func TestLoadBalancerRule_Ref(t *testing.T) {
	tests := []struct {
		name string
		obj  LoadBalancerRule
		want LoadBalancerRuleRef
	}{
		{
			name: "with id",
			obj: LoadBalancerRule{
				ID: "lbr_9IToFxX2AOl7IBSY",
			},
			want: LoadBalancerRuleRef{ID: "lbr_9IToFxX2AOl7IBSY"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.obj.Ref())
		})
	}
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
						ID:   "another abitrary string",
						Name: "cluster-42",
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
				Certificates: &[]CertificateRef{
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
		{
			name: "remove all certificates",
			obj: &LoadBalancerRuleArguments{
				Certificates: &[]CertificateRef{},
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
				Properties: &LoadBalancerRuleArguments{
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
				Properties: &LoadBalancerRuleArguments{
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
	type args struct {
		ctx  context.Context
		lb   LoadBalancerRef
		opts *ListOptions
	}
	tests := []struct {
		name     string
		args     args
		resp     *katapult.Response
		respErr  error
		respV    *loadBalancerRulesResponseBody
		wantReq  *katapult.Request
		want     []LoadBalancerRule
		wantResp *katapult.Response
		wantErr  string
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				lb:  LoadBalancerRef{ID: "lbrule_3W0eRZLQYHpTCPNX"},
				opts: &ListOptions{
					Page:    5,
					PerPage: 32,
				},
			},
			resp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusOK},
			},
			respV: &loadBalancerRulesResponseBody{
				Pagination: &katapult.Pagination{Total: 394},
				LoadBalancerRules: []LoadBalancerRule{
					{ID: "lbrule_3W0eRZLQYHpTCPNX", DestinationPort: 666},
				},
			},
			wantReq: &katapult.Request{
				Method: "GET",
				URL: &url.URL{
					Path: "/core/v1/load_balancers/_/rules",
					RawQuery: url.Values{
						"page":     []string{"5"},
						"per_page": []string{"32"},
						"load_balancer[id]": []string{
							"lbrule_3W0eRZLQYHpTCPNX",
						},
					}.Encode(),
				},
			},
			want: []LoadBalancerRule{{
				ID:              "lbrule_3W0eRZLQYHpTCPNX",
				DestinationPort: 666,
			}},
			wantResp: &katapult.Response{
				Pagination: &katapult.Pagination{Total: 394},
				Response:   &http.Response{StatusCode: http.StatusOK},
			},
		},
		{
			name: "success with nil options",
			args: args{
				ctx:  context.Background(),
				lb:   LoadBalancerRef{ID: "lbrule_3W0eRZLQYHpTCPNX"},
				opts: nil,
			},
			resp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusOK},
			},
			respV: &loadBalancerRulesResponseBody{
				Pagination: &katapult.Pagination{Total: 333},
				LoadBalancerRules: []LoadBalancerRule{
					{ID: "lbrule_3W0eRZLQYHpTCPNX", DestinationPort: 666},
				},
			},
			wantReq: &katapult.Request{
				Method: "GET",
				URL: &url.URL{
					Path: "/core/v1/load_balancers/_/rules",
					RawQuery: url.Values{
						"load_balancer[id]": []string{
							"lbrule_3W0eRZLQYHpTCPNX",
						},
					}.Encode(),
				},
			},
			want: []LoadBalancerRule{{
				ID:              "lbrule_3W0eRZLQYHpTCPNX",
				DestinationPort: 666,
			}},
			wantResp: &katapult.Response{
				Pagination: &katapult.Pagination{Total: 333},
				Response:   &http.Response{StatusCode: http.StatusOK},
			},
		},
		{
			name: "request error",
			args: args{
				ctx:  context.Background(),
				lb:   LoadBalancerRef{ID: "lbrule_3W0eRZLQYHpTCPNX"},
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
				lb:   LoadBalancerRef{ID: "lbrule_3W0eRZLQYHpTCPNX"},
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
			c := NewLoadBalancerRulesClient(tc)
			ctx := test.Context(tt.args.ctx)

			got, resp, err := c.List(ctx, tt.args.lb, tt.args.opts, testRequestOption)

			assert.Equal(t, 1, len(tc.Calls), "only 1 request should be made")
			test.AssertContext(t, ctx, tc.Ctx)

			assert.Equal(t, tt.want, got)

			if tt.wantResp != nil {
				assert.Equal(t, tt.wantResp, resp)
			}

			if tt.wantResp != nil {
				assert.Equal(t, tt.wantResp, resp)
			}

			if tt.wantReq != nil {
				setWantRequestOptionHeader(tt.wantReq)
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

func TestLoadBalancerRulesClient_Get(t *testing.T) {
	type args struct {
		ctx context.Context
		ref LoadBalancerRuleRef
	}
	tests := []struct {
		name     string
		args     args
		resp     *katapult.Response
		respErr  error
		respV    *loadBalancerRulesResponseBody
		wantReq  *katapult.Request
		want     *LoadBalancerRule
		wantResp *katapult.Response
		wantErr  string
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				ref: LoadBalancerRuleRef{ID: "123"},
			},
			resp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusOK},
			},
			respV: &loadBalancerRulesResponseBody{
				LoadBalancerRule: &LoadBalancerRule{
					ID:         "123",
					ListenPort: 132,
				},
			},
			wantReq: &katapult.Request{
				Method: "GET",
				URL: &url.URL{
					Path: "/core/v1/load_balancers/rules/_",
					RawQuery: url.Values{
						"load_balancer_rule[id]": []string{"123"},
					}.Encode(),
				},
			},
			want: &LoadBalancerRule{
				ID:         "123",
				ListenPort: 132,
			},
			wantResp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusOK},
			},
		},
		{
			name: "request error",
			args: args{
				ctx: context.Background(),
				ref: LoadBalancerRuleRef{ID: "123"},
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
				ref: LoadBalancerRuleRef{ID: "123"},
			},
			resp:    nil,
			respErr: fmt.Errorf("someting is really wrong"),
			wantErr: "someting is really wrong",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := testclient.New(tt.resp, tt.respErr, tt.respV)
			c := NewLoadBalancerRulesClient(tc)
			ctx := test.Context(tt.args.ctx)

			got, resp, err := c.Get(ctx, tt.args.ref, testRequestOption)

			assert.Equal(t, 1, len(tc.Calls), "only 1 request should be made")
			test.AssertContext(t, ctx, tc.Ctx)

			assert.Equal(t, tt.want, got)

			if tt.wantResp != nil {
				assert.Equal(t, tt.wantResp, resp)
			}

			if tt.wantReq != nil {
				setWantRequestOptionHeader(tt.wantReq)
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

func TestLoadBalancerRulesClient_GetByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name     string
		args     args
		resp     *katapult.Response
		respErr  error
		respV    *loadBalancerRulesResponseBody
		wantReq  *katapult.Request
		want     *LoadBalancerRule
		wantResp *katapult.Response
		wantErr  string
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				id:  "123",
			},
			resp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusOK},
			},
			respV: &loadBalancerRulesResponseBody{
				LoadBalancerRule: &LoadBalancerRule{
					ID:         "123",
					ListenPort: 132,
				},
			},
			wantReq: &katapult.Request{
				Method: "GET",
				URL: &url.URL{
					Path: "/core/v1/load_balancers/rules/_",
					RawQuery: url.Values{
						"load_balancer_rule[id]": []string{"123"},
					}.Encode(),
				},
			},
			want: &LoadBalancerRule{
				ID:         "123",
				ListenPort: 132,
			},
			wantResp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusOK},
			},
		},
		{
			name: "request error",
			args: args{
				ctx: context.Background(),
				id:  "123",
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
				id:  "123",
			},
			resp:    nil,
			respErr: fmt.Errorf("someting is really wrong"),
			wantErr: "someting is really wrong",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := testclient.New(tt.resp, tt.respErr, tt.respV)
			c := NewLoadBalancerRulesClient(tc)
			ctx := test.Context(tt.args.ctx)

			got, resp, err := c.GetByID(ctx, tt.args.id, testRequestOption)

			assert.Equal(t, 1, len(tc.Calls), "only 1 request should be made")
			test.AssertContext(t, ctx, tc.Ctx)

			assert.Equal(t, tt.want, got)

			if tt.wantResp != nil {
				assert.Equal(t, tt.wantResp, resp)
			}

			if tt.wantReq != nil {
				setWantRequestOptionHeader(tt.wantReq)
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

func TestLoadBalancerRulesClient_Create(t *testing.T) {
	type args struct {
		ctx  context.Context
		ref  LoadBalancerRef
		args *LoadBalancerRuleArguments
	}
	tests := []struct {
		name     string
		args     args
		resp     *katapult.Response
		respErr  error
		respV    *loadBalancerRulesResponseBody
		wantReq  *katapult.Request
		want     *LoadBalancerRule
		wantResp *katapult.Response
		wantErr  string
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				ref: LoadBalancerRef{ID: "lb_aFr95Rvyt6L3eyiH"},
				args: &LoadBalancerRuleArguments{
					DestinationPort: 8080,
					ListenPort:      80,
					Protocol:        HTTPProtocol,
				},
			},
			resp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusOK},
			},
			respV: &loadBalancerRulesResponseBody{
				LoadBalancerRule: &LoadBalancerRule{
					ID:              "lbrule_55P1GfFvW5pPPhgh",
					DestinationPort: 8080,
					ListenPort:      80,
					Protocol:        HTTPProtocol,
				},
			},
			wantReq: &katapult.Request{
				Method: "POST",
				URL: &url.URL{
					Path: "/core/v1/load_balancers/lb_aFr95Rvyt6L3eyiH/rules",
				},
				Body: &loadBalancerRuleCreateRequest{
					Properties: &LoadBalancerRuleArguments{
						DestinationPort: 8080,
						ListenPort:      80,
						Protocol:        HTTPProtocol,
					},
				},
			},
			want: &LoadBalancerRule{
				ID:              "lbrule_55P1GfFvW5pPPhgh",
				DestinationPort: 8080,
				ListenPort:      80,
				Protocol:        HTTPProtocol,
			},
			wantResp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusOK},
			},
		},
		{
			name: "request error",
			args: args{
				ctx: context.Background(),
				ref: LoadBalancerRef{ID: "lb_aFr95Rvyt6L3eyiH"},
				args: &LoadBalancerRuleArguments{
					DestinationPort: 8080,
					ListenPort:      80,
					Protocol:        HTTPProtocol,
				},
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
				ref: LoadBalancerRef{ID: "lb_aFr95Rvyt6L3eyiH"},
				args: &LoadBalancerRuleArguments{
					DestinationPort: 8080,
					ListenPort:      80,
					Protocol:        HTTPProtocol,
				},
			},
			resp:    nil,
			respErr: fmt.Errorf("someting is really wrong"),
			wantErr: "someting is really wrong",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := testclient.New(tt.resp, tt.respErr, tt.respV)
			c := NewLoadBalancerRulesClient(tc)
			ctx := test.Context(tt.args.ctx)

			got, resp, err := c.Create(ctx, tt.args.ref, tt.args.args, testRequestOption)

			assert.Equal(t, 1, len(tc.Calls), "only 1 request should be made")
			test.AssertContext(t, ctx, tc.Ctx)

			assert.Equal(t, tt.want, got)

			if tt.wantResp != nil {
				assert.Equal(t, tt.wantResp, resp)
			}

			if tt.wantReq != nil {
				setWantRequestOptionHeader(tt.wantReq)
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

func TestLoadBalancerRulesClient_Update(t *testing.T) {
	type args struct {
		ctx  context.Context
		ref  LoadBalancerRuleRef
		args *LoadBalancerRuleArguments
	}
	tests := []struct {
		name     string
		args     args
		resp     *katapult.Response
		respErr  error
		respV    *loadBalancerRulesResponseBody
		wantReq  *katapult.Request
		want     *LoadBalancerRule
		wantResp *katapult.Response
		wantErr  string
	}{
		{
			name: "success",
			args: args{
				ctx:  context.Background(),
				ref:  LoadBalancerRuleRef{ID: "lbrule_GDPBAqW3dm71i4ol"},
				args: &LoadBalancerRuleArguments{DestinationPort: 3000},
			},
			resp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusOK},
			},
			respV: &loadBalancerRulesResponseBody{
				LoadBalancerRule: &LoadBalancerRule{
					ID:              "lbrule_GDPBAqW3dm71i4ol",
					DestinationPort: 3000,
					ListenPort:      80,
					Protocol:        HTTPProtocol,
				},
			},
			wantReq: &katapult.Request{
				Method: "PATCH",
				URL: &url.URL{
					Path: "/core/v1/load_balancers/rules/_",
					RawQuery: url.Values{
						"load_balancer_rule[id]": []string{
							"lbrule_GDPBAqW3dm71i4ol",
						},
					}.Encode(),
				},
				Body: &loadBalancerRuleUpdateRequest{
					Properties: &LoadBalancerRuleArguments{
						DestinationPort: 3000,
					},
				},
			},
			want: &LoadBalancerRule{
				ID:              "lbrule_GDPBAqW3dm71i4ol",
				DestinationPort: 3000,
				ListenPort:      80,
				Protocol:        HTTPProtocol,
			},
			wantResp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusOK},
			},
		},
		{
			name: "request error",
			args: args{
				ctx:  context.Background(),
				ref:  LoadBalancerRuleRef{ID: "lbrule_GDPBAqW3dm71i4ol"},
				args: &LoadBalancerRuleArguments{DestinationPort: 3000},
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
				ref:  LoadBalancerRuleRef{ID: "lbrule_GDPBAqW3dm71i4ol"},
				args: &LoadBalancerRuleArguments{DestinationPort: 3000},
			},
			resp:    nil,
			respErr: fmt.Errorf("someting is really wrong"),
			wantErr: "someting is really wrong",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := testclient.New(tt.resp, tt.respErr, tt.respV)
			c := NewLoadBalancerRulesClient(tc)
			ctx := test.Context(tt.args.ctx)

			got, resp, err := c.Update(ctx, tt.args.ref, tt.args.args, testRequestOption)

			assert.Equal(t, 1, len(tc.Calls), "only 1 request should be made")
			test.AssertContext(t, ctx, tc.Ctx)

			assert.Equal(t, tt.want, got)

			if tt.wantResp != nil {
				assert.Equal(t, tt.wantResp, resp)
			}

			if tt.wantReq != nil {
				setWantRequestOptionHeader(tt.wantReq)
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

func TestLoadBalancerRulesClient_Delete(t *testing.T) {
	type args struct {
		ctx context.Context
		ref LoadBalancerRuleRef
	}
	tests := []struct {
		name     string
		args     args
		resp     *katapult.Response
		respErr  error
		respV    *loadBalancerRulesResponseBody
		wantReq  *katapult.Request
		want     *LoadBalancerRule
		wantResp *katapult.Response
		wantErr  string
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				ref: LoadBalancerRuleRef{ID: "lbrule_HfVizqDuo2B5B9kU"},
			},
			resp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusOK},
			},
			respV: &loadBalancerRulesResponseBody{
				LoadBalancerRule: &LoadBalancerRule{
					ID: "lbrule_HfVizqDuo2B5B9kU",
				},
			},
			wantReq: &katapult.Request{
				Method: "DELETE",
				URL: &url.URL{
					Path: "/core/v1/load_balancers/rules/_",
					RawQuery: url.Values{
						"load_balancer_rule[id]": []string{
							"lbrule_HfVizqDuo2B5B9kU",
						},
					}.Encode(),
				},
			},
			want: &LoadBalancerRule{
				ID: "lbrule_HfVizqDuo2B5B9kU",
			}, wantResp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusOK},
			},
		},
		{
			name: "request error",
			args: args{
				ctx: context.Background(),
				ref: LoadBalancerRuleRef{ID: "lbrule_HfVizqDuo2B5B9kU"},
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
				ref: LoadBalancerRuleRef{ID: "lbrule_HfVizqDuo2B5B9kU"},
			},
			resp:    nil,
			respErr: fmt.Errorf("someting is really wrong"),
			wantErr: "someting is really wrong",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := testclient.New(tt.resp, tt.respErr, tt.respV)
			c := NewLoadBalancerRulesClient(tc)
			ctx := test.Context(tt.args.ctx)

			got, resp, err := c.Delete(ctx, tt.args.ref, testRequestOption)

			assert.Equal(t, 1, len(tc.Calls), "only 1 request should be made")
			test.AssertContext(t, ctx, tc.Ctx)

			assert.Equal(t, tt.want, got)

			if tt.wantResp != nil {
				assert.Equal(t, tt.wantResp, resp)
			}

			if tt.wantReq != nil {
				setWantRequestOptionHeader(tt.wantReq)
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
