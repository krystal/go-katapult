package katapult

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	fixtureLoadBalancerNotFoundErr = "load_balancer_not_found: No load " +
		"balancer was found matching any of the criteria provided in " +
		"the arguments"
	fixtureLoadBalancerNotFoundResponseError = &ResponseError{
		Code: "load_balancer_not_found",
		Description: "No load balancer was found matching any of the " +
			"criteria provided in the arguments",
		Detail: json.RawMessage(`{}`),
	}
)

func TestLoadBalancer_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *LoadBalancer
	}{
		{
			name: "empty",
			obj:  &LoadBalancer{},
		},
		{
			name: "full",
			obj: &LoadBalancer{
				ID:                    "lb_9IToFxX2AOl7IBSY",
				Name:                  "web-1",
				ResourceType:          VirtualMachinesResourceType,
				ResourceIDs:           []string{"id2", "id3"},
				IPAddress:             &IPAddress{Address: "134.11.14.137"},
				HTTPSRedirect:         true,
				BackendCertificate:    "--BEGIN CERT--\n--END CERT--",
				BackendCertificateKey: "--BEGIN KEY--\n--END KEY--",
			},
		},
		{
			name: "tags resource type",
			obj: &LoadBalancer{
				ResourceType: TagsResourceType,
				ResourceIDs:  []string{"id2", "id3"},
			},
		},
		{
			name: "virtual_machine_groups resource type",
			obj: &LoadBalancer{
				ResourceType: VirtualMachineGroupsResourceType,
				ResourceIDs:  []string{"id2", "id3"},
			},
		},
		{
			name: "virtual_machines resource type",
			obj: &LoadBalancer{
				ResourceType: VirtualMachinesResourceType,
				ResourceIDs:  []string{"id2", "id3"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestLoadBalancer_LookupReference(t *testing.T) {
	tests := []struct {
		name string
		obj  *LoadBalancer
		want *LoadBalancer
	}{
		{
			name: "nil",
			obj:  (*LoadBalancer)(nil),
			want: nil,
		},
		{
			name: "empty",
			obj:  &LoadBalancer{},
			want: &LoadBalancer{},
		},
		{
			name: "full",
			obj: &LoadBalancer{
				ID:                    "lb_9IToFxX2AOl7IBSY",
				Name:                  "web-1",
				ResourceType:          VirtualMachinesResourceType,
				ResourceIDs:           []string{"id2", "id3"},
				IPAddress:             &IPAddress{Address: "134.11.14.137"},
				HTTPSRedirect:         true,
				BackendCertificate:    "--BEGIN CERT--\n--END CERT--",
				BackendCertificateKey: "--BEGIN KEY--\n--END KEY--",
			},
			want: &LoadBalancer{ID: "lb_9IToFxX2AOl7IBSY"},
		},
		{
			name: "no ID",
			obj: &LoadBalancer{
				Name:                  "web-1",
				ResourceType:          VirtualMachinesResourceType,
				ResourceIDs:           []string{"id2", "id3"},
				IPAddress:             &IPAddress{Address: "134.11.14.137"},
				HTTPSRedirect:         true,
				BackendCertificate:    "--BEGIN CERT--\n--END CERT--",
				BackendCertificateKey: "--BEGIN KEY--\n--END KEY--",
			},
			want: &LoadBalancer{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.obj.LookupReference()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLoadBalancer_UnmarshalJSON_Invalid(t *testing.T) {
	lb := &LoadBalancer{}
	raw := []byte(`{"id":"lb_foo","name":}`)

	err := lb.UnmarshalJSON(raw)

	assert.EqualError(t,
		err, "invalid character '}' looking for beginning of value",
	)
}

func TestLoadBalancerArguments_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *LoadBalancerArguments
	}{
		{
			name: "empty",
			obj:  &LoadBalancerArguments{},
		},
		{
			name: "full",
			obj: &LoadBalancerArguments{
				Name:         "helper",
				ResourceType: TagsResourceType,
				ResourceIDs:  []string{"id1", "id2"},
				DataCenter:   &DataCenter{ID: "id4"},
			},
		},
		{
			name: "with empty ResourceIDs",
			obj: &LoadBalancerArguments{
				Name:         "helper",
				ResourceType: TagsResourceType,
				ResourceIDs:  []string{},
				DataCenter:   &DataCenter{ID: "id4"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestLoadBalancerArguments_forRequest(t *testing.T) {
	tests := []struct {
		name string
		obj  *LoadBalancerArguments
		want *LoadBalancerArguments
	}{
		{
			name: "nil",
			obj:  (*LoadBalancerArguments)(nil),
			want: nil,
		},
		{
			name: "empty",
			obj:  &LoadBalancerArguments{},
			want: &LoadBalancerArguments{},
		},
		{
			name: "full",
			obj: &LoadBalancerArguments{
				Name:         "helper",
				ResourceType: TagsResourceType,
				ResourceIDs:  []string{"id1", "id2"},
				DataCenter: &DataCenter{
					ID:        "dc_25d48761871e4bf",
					Name:      "Woodland",
					Permalink: "woodland",
				},
			},
			want: &LoadBalancerArguments{
				Name:         "helper",
				ResourceType: TagsResourceType,
				ResourceIDs:  []string{"id1", "id2"},
				DataCenter:   &DataCenter{ID: "dc_25d48761871e4bf"},
			},
		},
		{
			name: "data center has Permalink by no ID",
			obj: &LoadBalancerArguments{
				Name:         "helper",
				ResourceType: TagsResourceType,
				ResourceIDs:  []string{"id1", "id2"},
				DataCenter: &DataCenter{
					Name:      "Woodland",
					Permalink: "woodland",
				},
			},
			want: &LoadBalancerArguments{
				Name:         "helper",
				ResourceType: TagsResourceType,
				ResourceIDs:  []string{"id1", "id2"},
				DataCenter:   &DataCenter{Permalink: "woodland"},
			},
		},
		{
			name: "data center has no Permalink or ID",
			obj: &LoadBalancerArguments{
				Name:         "helper",
				ResourceType: TagsResourceType,
				ResourceIDs:  []string{"id1", "id2"},
				DataCenter: &DataCenter{
					Name: "Woodland",
				},
			},
			want: &LoadBalancerArguments{
				Name:         "helper",
				ResourceType: TagsResourceType,
				ResourceIDs:  []string{"id1", "id2"},
				DataCenter:   &DataCenter{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dcName string
			if tt.obj != nil && tt.obj.DataCenter != nil {
				dcName = tt.obj.DataCenter.Name
			}

			got := tt.obj.forRequest()

			assert.Equal(t, tt.want, got)

			if dcName != "" {
				assert.Equal(t, dcName, tt.obj.DataCenter.Name,
					"original object was modified")
			}
		})
	}
}

func Test_loadBalancerResource_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *loadBalancerResource
	}{
		{
			name: "empty",
			obj:  &loadBalancerResource{},
		},
		{
			name: "full",
			obj: &loadBalancerResource{
				Type:  "VirtualMachine",
				Value: &loadBalancerResourceValue{ID: "id4"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func Test_loadBalancerResourceValue_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *loadBalancerResourceValue
	}{
		{
			name: "empty",
			obj:  &loadBalancerResourceValue{},
		},
		{
			name: "full",
			obj: &loadBalancerResourceValue{
				ID:   "id4",
				Name: "helper",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func Test_loadBalancerCreateRequest_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *loadBalancerCreateRequest
	}{
		{
			name: "empty",
			obj:  &loadBalancerCreateRequest{},
		},
		{
			name: "full",
			obj: &loadBalancerCreateRequest{
				Organization: &Organization{ID: "org_rs55YZNYMw7o3jnQ"},
				Properties:   &LoadBalancerArguments{Name: "web-1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func Test_loadBalancerUpdateRequest_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *loadBalancerUpdateRequest
	}{
		{
			name: "empty",
			obj:  &loadBalancerUpdateRequest{},
		},
		{
			name: "full",
			obj: &loadBalancerUpdateRequest{
				LoadBalancer: &LoadBalancer{ID: "lb_0krMCRl7DIZr0XV2"},
				Properties:   &LoadBalancerArguments{Name: "web-east-1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func Test_loadBalancersResponseBody_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *loadBalancersResponseBody
	}{
		{
			name: "empty",
			obj:  &loadBalancersResponseBody{},
		},
		{
			name: "full",
			obj: &loadBalancersResponseBody{
				Pagination:    &Pagination{CurrentPage: 344},
				LoadBalancer:  &LoadBalancer{ID: "id1"},
				LoadBalancers: []*LoadBalancer{{ID: "id2"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestLoadBalancersClient_List(t *testing.T) {
	// Correlates to fixtures/load_balancers_list*.json
	loadBalancerList := []*LoadBalancer{
		{
			ID:           "lb_7vClpn0rlUegGPDS",
			Name:         "web",
			ResourceType: "tags",
		},
		{
			ID:           "lb_sESSo8rKfcL79D3y",
			Name:         "db",
			ResourceType: "virtual_machines",
		},
		{
			ID:           "lb_WSjTHQDJ6jOjzXVy",
			Name:         "assets",
			ResourceType: "virtual_machine_groups",
		},
	}

	type args struct {
		ctx  context.Context
		org  *Organization
		opts *ListOptions
	}
	tests := []struct {
		name       string
		args       args
		want       []*LoadBalancer
		pagination *Pagination
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "load balancers",
			args: args{
				ctx: context.Background(),
				org: &Organization{ID: "org_O648YDMEYeLmqdmn"},
			},
			want: loadBalancerList,
			pagination: &Pagination{
				CurrentPage: 1,
				TotalPages:  1,
				Total:       3,
				PerPage:     30,
				LargeSet:    false,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("load_balancers_list"),
		},
		{
			name: "page 1 of load balancers",
			args: args{
				ctx:  context.Background(),
				org:  &Organization{ID: "org_O648YDMEYeLmqdmn"},
				opts: &ListOptions{Page: 1, PerPage: 2},
			},
			want: loadBalancerList[0:2],
			pagination: &Pagination{
				CurrentPage: 1,
				TotalPages:  2,
				Total:       3,
				PerPage:     2,
				LargeSet:    false,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("load_balancers_list_page_1"),
		},
		{
			name: "page 2 of load balancers",
			args: args{
				ctx:  context.Background(),
				org:  &Organization{ID: "org_O648YDMEYeLmqdmn"},
				opts: &ListOptions{Page: 2, PerPage: 2},
			},
			want: loadBalancerList[2:],
			pagination: &Pagination{
				CurrentPage: 2,
				TotalPages:  2,
				Total:       3,
				PerPage:     2,
				LargeSet:    false,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("load_balancers_list_page_2"),
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
			c, mux, _, teardown := prepareTestClient()
			defer teardown()

			org := tt.args.org
			if org == nil {
				org = &Organization{ID: "_"}
			}

			mux.HandleFunc(
				fmt.Sprintf(
					"/core/v1/organizations/%s/load_balancers",
					org.ID,
				),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					if tt.args.opts != nil {
						assert.Equal(t, *tt.args.opts.Values(), r.URL.Query())
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.LoadBalancers.List(
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

			if tt.pagination != nil {
				assert.Equal(t, tt.pagination, resp.Pagination)
			}

			if tt.errResp != nil {
				assert.Equal(t, tt.errResp, resp.Error)
			}
		})
	}
}

func TestLoadBalancersClient_Get(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name       string
		args       args
		id         string
		want       *LoadBalancer
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "load balancer",
			args: args{
				ctx: context.Background(),
				id:  "lb_7vClpn0rlUegGPDS",
			},
			want: &LoadBalancer{
				ID:           "lb_7vClpn0rlUegGPDS",
				Name:         "web",
				ResourceType: TagsResourceType,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("load_balancer_get"),
		},
		{
			name: "non-existent load balancer",
			args: args{
				ctx: context.Background(),
				id:  "lb_nopethisbegone",
			},
			errStr:     fixtureLoadBalancerNotFoundErr,
			errResp:    fixtureLoadBalancerNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("load_balancer_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				id:  "lb_7vClpn0rlUegGPDS",
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := prepareTestClient()
			defer teardown()

			mux.HandleFunc(
				fmt.Sprintf("/core/v1/load_balancers/%s", tt.args.id),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.LoadBalancers.Get(tt.args.ctx, tt.args.id)

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

func TestLoadBalancersClient_GetByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name       string
		args       args
		id         string
		want       *LoadBalancer
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "load balancer",
			args: args{
				ctx: context.Background(),
				id:  "lb_7vClpn0rlUegGPDS",
			},
			want: &LoadBalancer{
				ID:           "lb_7vClpn0rlUegGPDS",
				Name:         "web",
				ResourceType: TagsResourceType,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("load_balancer_get"),
		},
		{
			name: "non-existent load balancer",
			args: args{
				ctx: context.Background(),
				id:  "lb_nopethisbegone",
			},
			errStr:     fixtureLoadBalancerNotFoundErr,
			errResp:    fixtureLoadBalancerNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("load_balancer_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				id:  "lb_7vClpn0rlUegGPDS",
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := prepareTestClient()
			defer teardown()

			mux.HandleFunc(
				fmt.Sprintf("/core/v1/load_balancers/%s", tt.args.id),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.LoadBalancers.GetByID(tt.args.ctx, tt.args.id)

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

func TestLoadBalancersClient_Create(t *testing.T) {
	lbArgs := &LoadBalancerArguments{
		Name:         "api-test",
		ResourceType: VirtualMachinesResourceType,
		ResourceIDs:  []string{"id2", "id3"},
		DataCenter:   &DataCenter{ID: "id4", Name: "other"},
	}
	lbReqArgs := &LoadBalancerArguments{
		Name:         "api-test",
		ResourceType: VirtualMachinesResourceType,
		ResourceIDs:  []string{"id2", "id3"},
		DataCenter:   &DataCenter{ID: "id4"},
	}

	type args struct {
		ctx    context.Context
		org    *Organization
		lbArgs *LoadBalancerArguments
	}
	tests := []struct {
		name       string
		args       args
		reqBody    *loadBalancerCreateRequest
		want       *LoadBalancer
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "load balancer",
			args: args{
				ctx: context.Background(),
				org: &Organization{
					ID:        "org_O648YDMEYeLmqdmn",
					Name:      "ACME Inc.",
					SubDomain: "acme",
				},
				lbArgs: lbArgs,
			},
			reqBody: &loadBalancerCreateRequest{
				Organization: &Organization{ID: "org_O648YDMEYeLmqdmn"},
				Properties:   lbReqArgs,
			},
			want: &LoadBalancer{
				ID:           "lb_PuoZUW18K5bXEAVE",
				Name:         "api-test",
				ResourceType: "virtual_machines",
			},
			respStatus: http.StatusCreated,
			respBody:   fixture("load_balancer_create"),
		},
		{
			name: "organization by sub-domain",
			args: args{
				ctx: context.Background(),
				org: &Organization{
					Name:      "ACME Inc.",
					SubDomain: "acme",
				},
				lbArgs: lbArgs,
			},
			reqBody: &loadBalancerCreateRequest{
				Organization: &Organization{SubDomain: "acme"},
				Properties:   lbReqArgs,
			},
			want: &LoadBalancer{
				ID:           "lb_PuoZUW18K5bXEAVE",
				Name:         "api-test",
				ResourceType: "virtual_machines",
			},
			respStatus: http.StatusCreated,
			respBody:   fixture("load_balancer_create"),
		},
		{
			name: "without resource IDs",
			args: args{
				ctx: context.Background(),
				org: &Organization{ID: "org_O648YDMEYeLmqdmn"},
				lbArgs: &LoadBalancerArguments{
					Name:         lbArgs.Name,
					ResourceType: lbArgs.ResourceType,
					DataCenter:   lbArgs.DataCenter,
				},
			},
			reqBody: &loadBalancerCreateRequest{
				Organization: &Organization{ID: "org_O648YDMEYeLmqdmn"},
				Properties: &LoadBalancerArguments{
					Name:         lbReqArgs.Name,
					ResourceType: lbReqArgs.ResourceType,
					DataCenter:   lbReqArgs.DataCenter,
				},
			},
			want: &LoadBalancer{
				ID:           "lb_PuoZUW18K5bXEAVE",
				Name:         "api-test",
				ResourceType: "virtual_machines",
			},
			respStatus: http.StatusCreated,
			respBody:   fixture("load_balancer_create"),
		},
		{
			name: "without data center",
			args: args{
				ctx: context.Background(),
				org: &Organization{ID: "org_O648YDMEYeLmqdmn"},
				lbArgs: &LoadBalancerArguments{
					Name:         lbArgs.Name,
					ResourceType: lbArgs.ResourceType,
					ResourceIDs:  lbArgs.ResourceIDs,
				},
			},
			errStr:     fixtureDataCenterNotFoundErr,
			errResp:    fixtureDataCenterNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("data_center_not_found_error"),
		},
		{
			name: "non-existent Organization",
			args: args{
				ctx:    context.Background(),
				org:    &Organization{ID: "org_O648YDMEYeLmqdmn"},
				lbArgs: lbArgs,
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
				org:    &Organization{ID: "org_O648YDMEYeLmqdmn"},
				lbArgs: lbArgs,
			},
			errStr:     fixtureOrganizationSuspendedErr,
			errResp:    fixtureOrganizationSuspendedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("organization_suspended_error"),
		},
		{
			name: "non-existent data center",
			args: args{
				ctx:    context.Background(),
				org:    &Organization{ID: "org_O648YDMEYeLmqdmn"},
				lbArgs: lbArgs,
			},
			errStr:     fixtureDataCenterNotFoundErr,
			errResp:    fixtureDataCenterNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("data_center_not_found_error"),
		},
		{
			name: "permission denied",
			args: args{
				ctx:    context.Background(),
				org:    &Organization{ID: "org_O648YDMEYeLmqdmn"},
				lbArgs: lbArgs,
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
				org:    &Organization{ID: "org_O648YDMEYeLmqdmn"},
				lbArgs: lbArgs,
			},
			errStr:     fixtureValidationErrorErr,
			errResp:    fixtureValidationErrorResponseError,
			respStatus: http.StatusUnprocessableEntity,
			respBody:   fixture("validation_error"),
		},
		{
			name: "nil load balancer arguments",
			args: args{
				ctx:    context.Background(),
				org:    &Organization{ID: "org_O648YDMEYeLmqdmn"},
				lbArgs: nil,
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
				org:    &Organization{ID: "org_O648YDMEYeLmqdmn"},
				lbArgs: lbArgs,
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := prepareTestClient()
			defer teardown()

			var dcName string
			if tt.args.lbArgs != nil && tt.args.lbArgs.DataCenter != nil {
				dcName = tt.args.lbArgs.DataCenter.Name
			}

			mux.HandleFunc(
				"/core/v1/organizations/_/load_balancers",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "POST", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					if tt.reqBody != nil {
						reqBody := &loadBalancerCreateRequest{}
						err := strictUmarshal(r.Body, reqBody)
						assert.NoError(t, err)
						assert.Equal(t, tt.reqBody, reqBody)
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.LoadBalancers.Create(
				tt.args.ctx, tt.args.org, tt.args.lbArgs,
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

			if tt.args.lbArgs != nil && tt.args.lbArgs.DataCenter != nil {
				// ensure the input LoadBalancerArguments are not modified
				assert.Equal(t, dcName, tt.args.lbArgs.DataCenter.Name)
			}
		})
	}
}

func TestLoadBalancersClient_Update(t *testing.T) {
	lb := &LoadBalancer{
		ID:           "lb_7vClpn0rlUegGPDS",
		Name:         "web-1",
		ResourceType: VirtualMachineGroupsResourceType,
		ResourceIDs:  []string{"grp1", "grp3"},
		DataCenter: &DataCenter{
			ID:   "dc_a2417980b9874c0",
			Name: "New Town",
		},
	}
	lbArgs := &LoadBalancerArguments{
		Name:         "web-east-1",
		ResourceType: TagsResourceType,
		ResourceIDs:  []string{"tag2", "tag4"},
	}

	type args struct {
		ctx    context.Context
		lb     *LoadBalancer
		lbArgs *LoadBalancerArguments
	}
	tests := []struct {
		name       string
		args       args
		reqBody    *loadBalancerUpdateRequest
		want       *LoadBalancer
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "load balancer",
			args: args{
				ctx:    context.Background(),
				lb:     lb,
				lbArgs: lbArgs,
			},
			reqBody: &loadBalancerUpdateRequest{
				LoadBalancer: &LoadBalancer{ID: lb.ID},
				Properties:   lbArgs,
			},
			want: &LoadBalancer{
				ID:           "lb_7vClpn0rlUegGPDS",
				Name:         "web-east-1",
				ResourceType: TagsResourceType,
				ResourceIDs:  []string{"tag2", "tag4"},
			},
			respStatus: http.StatusCreated,
			respBody:   fixture("load_balancer_update"),
		},
		{
			name: "load balancer without ID",
			args: args{
				ctx: context.Background(),
				lb: &LoadBalancer{
					Name:         lb.Name,
					ResourceType: lb.ResourceType,
					ResourceIDs:  lb.ResourceIDs,
					DataCenter:   lb.DataCenter,
				},
				lbArgs: lbArgs,
			},
			reqBody: &loadBalancerUpdateRequest{
				LoadBalancer: &LoadBalancer{},
				Properties:   lbArgs,
			},
			errStr:     fixtureLoadBalancerNotFoundErr,
			errResp:    fixtureLoadBalancerNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("load_balancer_not_found_error"),
		},
		{
			name: "non-existent load balancer",
			args: args{
				ctx:    context.Background(),
				lb:     &LoadBalancer{ID: "lb_somethingnope"},
				lbArgs: lbArgs,
			},
			reqBody: &loadBalancerUpdateRequest{
				LoadBalancer: &LoadBalancer{ID: "lb_somethingnope"},
				Properties:   lbArgs,
			},
			errStr:     fixtureLoadBalancerNotFoundErr,
			errResp:    fixtureLoadBalancerNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("load_balancer_not_found_error"),
		},
		{
			name: "permission denied",
			args: args{
				ctx:    context.Background(),
				lb:     lb,
				lbArgs: lbArgs,
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
				lb:     lb,
				lbArgs: lbArgs,
			},
			errStr:     fixtureValidationErrorErr,
			errResp:    fixtureValidationErrorResponseError,
			respStatus: http.StatusUnprocessableEntity,
			respBody:   fixture("validation_error"),
		},
		{
			name: "nil load balancer argument",
			args: args{
				ctx:    context.Background(),
				lb:     nil,
				lbArgs: lbArgs,
			},
			reqBody: &loadBalancerUpdateRequest{
				LoadBalancer: nil,
				Properties:   lbArgs,
			},
			errStr:     fixtureLoadBalancerNotFoundErr,
			errResp:    fixtureLoadBalancerNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("load_balancer_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				lb:  lb,
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := prepareTestClient()
			defer teardown()

			mux.HandleFunc(
				"/core/v1/load_balancers/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "PATCH", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					if tt.reqBody != nil {
						reqBody := &loadBalancerUpdateRequest{}
						err := strictUmarshal(r.Body, reqBody)
						assert.NoError(t, err)
						assert.Equal(t, tt.reqBody, reqBody)
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.LoadBalancers.Update(
				tt.args.ctx, tt.args.lb, tt.args.lbArgs,
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

func TestLoadBalancersClient_Delete(t *testing.T) {
	type args struct {
		ctx context.Context
		lb  *LoadBalancer
	}
	tests := []struct {
		name       string
		args       args
		want       *LoadBalancer
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "load balancer",
			args: args{
				ctx: context.Background(),
				lb:  &LoadBalancer{ID: "lb_7vClpn0rlUegGPDS"},
			},
			want: &LoadBalancer{
				ID:           "lb_7vClpn0rlUegGPDS",
				Name:         "web",
				ResourceType: TagsResourceType,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("load_balancer_get"),
		},
		{
			name: "non-existent load balancer",
			args: args{
				ctx: context.Background(),
				lb:  &LoadBalancer{ID: "lb_nopenotfound"},
			},
			errStr:     fixtureLoadBalancerNotFoundErr,
			errResp:    fixtureLoadBalancerNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("load_balancer_not_found_error"),
		},
		{
			name: "nil load balancer",
			args: args{
				ctx: context.Background(),
				lb:  nil,
			},
			errStr:     fixtureLoadBalancerNotFoundErr,
			errResp:    fixtureLoadBalancerNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("load_balancer_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				lb:  &LoadBalancer{ID: "lb_7vClpn0rlUegGPDS"},
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := prepareTestClient()
			defer teardown()

			lb := tt.args.lb
			if lb == nil {
				lb = &LoadBalancer{ID: "_"}
			}

			mux.HandleFunc(
				fmt.Sprintf("/core/v1/load_balancers/%s", lb.ID),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "DELETE", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.LoadBalancers.Delete(tt.args.ctx, tt.args.lb)

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
