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
				ID:                    "id",
				Name:                  "name",
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

func TestLoadBalancer_UnmarshalJSON_Invalid(t *testing.T) {
	lb := &LoadBalancer{}
	raw := []byte(`{"id":"lb_foo","name":}`)

	err := lb.UnmarshalJSON(raw)

	assert.EqualError(t,
		err, "invalid character '}' looking for beginning of value",
	)
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

func TestLoadBalancerArguments_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name    string
		obj     *LoadBalancerArguments
		decoded *LoadBalancerArguments
	}{
		{
			name:    "empty",
			obj:     &LoadBalancerArguments{},
			decoded: nil,
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
			name: "without ResourceIDs",
			obj: &LoadBalancerArguments{
				Name:         "helper",
				ResourceType: TagsResourceType,
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
		{
			name: "strips down data center to just ID",
			obj: &LoadBalancerArguments{
				DataCenter: &DataCenter{
					ID:        "id4",
					Name:      "Woodland",
					Permalink: "woodland",
				},
			},
			decoded: &LoadBalancerArguments{
				DataCenter: &DataCenter{ID: "id4"},
			},
		},
		{
			name: "strips down data center to just Permalink if no ID",
			obj: &LoadBalancerArguments{
				DataCenter: &DataCenter{
					Name:      "Woodland",
					Permalink: "woodland",
				},
			},
			decoded: &LoadBalancerArguments{
				DataCenter: &DataCenter{Permalink: "woodland"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dcName string
			if tt.obj.DataCenter != nil && tt.obj.DataCenter.Name != "" {
				dcName = tt.obj.DataCenter.Name
			}

			if tt.decoded != nil {
				testCustomJSONMarshaling(t, tt.obj, tt.decoded)
			} else {
				testJSONMarshaling(t, tt.obj)
			}

			if tt.obj.DataCenter != nil && dcName != "" {
				// ensure the input LoadBalancerArguments are not modified
				assert.Equal(t, dcName, tt.obj.DataCenter.Name)
			}
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
		ctx   context.Context
		orgID string
		opts  *ListOptions
	}
	tests := []struct {
		name       string
		args       args
		expected   []*LoadBalancer
		pagination *Pagination
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "load balancers",
			args: args{
				ctx:   context.Background(),
				orgID: "org_O648YDMEYeLmqdmn",
			},
			expected: loadBalancerList,
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
				ctx:   context.Background(),
				orgID: "org_O648YDMEYeLmqdmn",
				opts:  &ListOptions{Page: 1, PerPage: 2},
			},
			expected: loadBalancerList[0:2],
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
				ctx:   context.Background(),
				orgID: "org_O648YDMEYeLmqdmn",
				opts:  &ListOptions{Page: 2, PerPage: 2},
			},
			expected: loadBalancerList[2:],
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
				ctx:   context.Background(),
				orgID: "org_O648YDMEYeLmqdmn",
			},
			errStr:     fixtureInvalidAPITokenErr,
			errResp:    fixtureInvalidAPITokenResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("invalid_api_token_error"),
		},
		{
			name: "non-existent organization",
			args: args{
				ctx:   context.Background(),
				orgID: "org_O648YDMEYeLmqdmn",
			},
			errStr:     fixtureOrganizationNotFoundErr,
			errResp:    fixtureOrganizationNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("organization_not_found_error"),
		},
		{
			name: "suspended organization",
			args: args{
				ctx:   context.Background(),
				orgID: "org_O648YDMEYeLmqdmn",
			},
			errStr:     fixtureOrganizationSuspendedErr,
			errResp:    fixtureOrganizationSuspendedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("organization_suspended_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx:   nil,
				orgID: "org_O648YDMEYeLmqdmn",
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := prepareTestClient()
			defer teardown()

			mux.HandleFunc(
				fmt.Sprintf(
					"/core/v1/organizations/%s/load_balancers", tt.args.orgID,
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
				tt.args.ctx, tt.args.orgID, tt.args.opts,
			)

			if tt.respStatus != 0 {
				assert.Equal(t, tt.respStatus, resp.StatusCode)
			}

			if tt.errStr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.errStr)
			}

			if tt.expected != nil {
				assert.Equal(t, tt.expected, got)
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
		expected   *LoadBalancer
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
			expected: &LoadBalancer{
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

			if tt.expected != nil {
				assert.Equal(t, tt.expected, got)
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
	lbArgsWithoutResourceIDs := &LoadBalancerArguments{
		Name:         lbArgs.Name,
		ResourceType: lbArgs.ResourceType,
		DataCenter:   lbArgs.DataCenter,
	}
	lbArgsWithoutDataCenter := &LoadBalancerArguments{
		Name:         lbArgs.Name,
		ResourceType: lbArgs.ResourceType,
		ResourceIDs:  lbArgs.ResourceIDs,
	}

	type reqBodyDataCenter struct {
		ID string `json:"id,omitempty"`
	}
	type reqBodyArguments struct {
		Name         string             `json:"name,omitempty"`
		ResourceType ResourceType       `json:"resource_type,omitempty"`
		ResourceIDs  []string           `json:"resource_ids,omitempty"`
		DataCenter   *reqBodyDataCenter `json:"data_center,omitempty"`
	}
	type reqBody struct {
		Properties *reqBodyArguments `json:"properties"`
	}
	type args struct {
		ctx    context.Context
		orgID  string
		lbArgs *LoadBalancerArguments
	}
	tests := []struct {
		name       string
		args       args
		expected   *LoadBalancer
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "load balancer",
			args: args{
				ctx:    context.Background(),
				orgID:  "org_O648YDMEYeLmqdmn",
				lbArgs: lbArgs,
			},
			expected: &LoadBalancer{
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
				ctx:    context.Background(),
				orgID:  "org_O648YDMEYeLmqdmn",
				lbArgs: lbArgsWithoutResourceIDs,
			},
			expected: &LoadBalancer{
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
				ctx:    context.Background(),
				orgID:  "org_O648YDMEYeLmqdmn",
				lbArgs: lbArgsWithoutDataCenter,
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
				orgID:  "org_O648YDMEYeLmqdmn",
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
				orgID:  "org_O648YDMEYeLmqdmn",
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
				orgID:  "org_O648YDMEYeLmqdmn",
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
				orgID:  "org_O648YDMEYeLmqdmn",
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
				orgID:  "org_O648YDMEYeLmqdmn",
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
				orgID:  "org_O648YDMEYeLmqdmn",
				lbArgs: nil,
			},
			errStr: "nil load balancer arguments",
		},
		{
			name: "nil context",
			args: args{
				ctx:    nil,
				orgID:  "org_O648YDMEYeLmqdmn",
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
				fmt.Sprintf(
					"/core/v1/organizations/%s/load_balancers", tt.args.orgID,
				),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "POST", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					expectedReqArgs := &reqBodyArguments{
						Name:         tt.args.lbArgs.Name,
						ResourceType: tt.args.lbArgs.ResourceType,
						ResourceIDs:  tt.args.lbArgs.ResourceIDs,
					}
					if tt.args.lbArgs.DataCenter != nil {
						expectedReqArgs.DataCenter = &reqBodyDataCenter{
							ID: tt.args.lbArgs.DataCenter.ID,
						}
					}

					body := &reqBody{}
					err := strictUmarshal(r.Body, body)
					assert.NoError(t, err)
					assert.Equal(t,
						&reqBody{Properties: expectedReqArgs}, body,
					)

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.LoadBalancers.Create(
				tt.args.ctx, tt.args.orgID, tt.args.lbArgs,
			)

			if tt.respStatus != 0 {
				assert.Equal(t, tt.respStatus, resp.StatusCode)
			}

			if tt.errStr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.errStr)
			}

			if tt.expected != nil {
				assert.Equal(t, tt.expected, got)
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
	loadBalancer := &LoadBalancer{
		ID:           "lb_7vClpn0rlUegGPDS",
		Name:         "web-1",
		ResourceType: VirtualMachineGroupsResourceType,
		ResourceIDs:  []string{"grp1", "grp3"},
		DataCenter: &DataCenter{
			ID:   "loc_a2417980b9874c0",
			Name: "New Town",
		},
	}
	loadBalancerWithoutDataCenter := &LoadBalancer{
		ID:           loadBalancer.ID,
		Name:         loadBalancer.Name,
		ResourceType: loadBalancer.ResourceType,
		ResourceIDs:  loadBalancer.ResourceIDs,
	}
	loadBalancerWithoutID := &LoadBalancer{
		Name:         loadBalancer.Name,
		ResourceType: loadBalancer.ResourceType,
		ResourceIDs:  loadBalancer.ResourceIDs,
		DataCenter:   loadBalancer.DataCenter,
	}

	type reqBodyDataCenter struct {
		ID string `json:"id,omitempty"`
	}
	type reqBodyArguments struct {
		Name         string             `json:"name,omitempty"`
		ResourceType ResourceType       `json:"resource_type,omitempty"`
		ResourceIDs  []string           `json:"resource_ids,omitempty"`
		DataCenter   *reqBodyDataCenter `json:"data_center,omitempty"`
	}
	type reqBody struct {
		Properties *reqBodyArguments `json:"properties"`
	}
	type args struct {
		ctx context.Context
		lb  *LoadBalancer
	}
	tests := []struct {
		name       string
		args       args
		expected   *LoadBalancer
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "load balancer",
			args: args{
				ctx: context.Background(),
				lb:  loadBalancer,
			},
			expected: &LoadBalancer{
				ID:           "lb_7vClpn0rlUegGPDS",
				Name:         "web-1",
				ResourceType: VirtualMachineGroupsResourceType,
				ResourceIDs:  []string{"grp1", "grp3"},
			},
			respStatus: http.StatusCreated,
			respBody:   fixture("load_balancer_update"),
		},
		{
			name: "without data center",
			args: args{
				ctx: context.Background(),
				lb:  loadBalancerWithoutDataCenter,
			},
			respStatus: http.StatusCreated,
			respBody:   fixture("load_balancer_update"),
		},
		{
			name: "load balancer without ID",
			args: args{
				ctx: context.Background(),
				lb:  loadBalancerWithoutID,
			},
			errStr: "ID value is empty",
		},
		{
			name: "non-existent load balancer",
			args: args{
				ctx: context.Background(),
				lb:  loadBalancer,
			},
			errStr:     fixtureLoadBalancerNotFoundErr,
			errResp:    fixtureLoadBalancerNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("load_balancer_not_found_error"),
		},
		{
			name: "non-existent data center",
			args: args{
				ctx: context.Background(),
				lb:  loadBalancer,
			},
			errStr:     fixtureDataCenterNotFoundErr,
			errResp:    fixtureDataCenterNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("data_center_not_found_error"),
		},
		{
			name: "permission denied",
			args: args{
				ctx: context.Background(),
				lb:  loadBalancer,
			},
			errStr:     fixturePermissionDeniedErr,
			errResp:    fixturePermissionDeniedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("permission_denied_error"),
		},
		{
			name: "validation error",
			args: args{
				ctx: context.Background(),
				lb:  loadBalancer,
			},
			errStr:     fixtureValidationErrorErr,
			errResp:    fixtureValidationErrorResponseError,
			respStatus: http.StatusUnprocessableEntity,
			respBody:   fixture("validation_error"),
		},
		{
			name: "nil load balancer argument",
			args: args{
				ctx: context.Background(),
				lb:  nil,
			},
			errStr: "nil load balancer arguments",
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				lb:  loadBalancer,
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := prepareTestClient()
			defer teardown()

			var dcName string
			if tt.args.lb != nil && tt.args.lb.DataCenter != nil {
				dcName = tt.args.lb.DataCenter.Name
			}

			if tt.args.lb != nil {
				mux.HandleFunc(
					fmt.Sprintf("/core/v1/load_balancers/%s", tt.args.lb.ID),
					func(w http.ResponseWriter, r *http.Request) {
						assert.Equal(t, "PATCH", r.Method)
						assertEmptyFieldSpec(t, r)
						assertAuthorization(t, r)

						expectedReqArgs := &reqBodyArguments{
							Name:         tt.args.lb.Name,
							ResourceType: tt.args.lb.ResourceType,
							ResourceIDs:  tt.args.lb.ResourceIDs,
						}
						if tt.args.lb.DataCenter != nil {
							expectedReqArgs.DataCenter = &reqBodyDataCenter{
								ID: tt.args.lb.DataCenter.ID,
							}
						}

						body := &reqBody{}
						err := strictUmarshal(r.Body, body)
						assert.NoError(t, err)
						assert.Equal(t,
							&reqBody{Properties: expectedReqArgs}, body,
						)

						w.WriteHeader(tt.respStatus)
						_, _ = w.Write(tt.respBody)
					},
				)
			}

			got, resp, err := c.LoadBalancers.Update(
				tt.args.ctx, tt.args.lb,
			)

			if tt.respStatus != 0 {
				assert.Equal(t, tt.respStatus, resp.StatusCode)
			}

			if tt.errStr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.errStr)
			}

			if tt.expected != nil {
				assert.Equal(t, tt.expected, got)
			}

			if tt.errResp != nil {
				assert.Equal(t, tt.errResp, resp.Error)
			}

			if tt.args.lb != nil && tt.args.lb.DataCenter != nil {
				// ensure the input LoadBalancerArguments are not modified
				assert.Equal(t, dcName, tt.args.lb.DataCenter.Name)
			}
		})
	}
}

func TestLoadBalancersClient_Delete(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name       string
		args       args
		expected   *LoadBalancer
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
			expected: &LoadBalancer{
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
				id:  "lb_7vClpn0rlUegGPDS",
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
					assert.Equal(t, "DELETE", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.LoadBalancers.Delete(tt.args.ctx, tt.args.id)

			if tt.respStatus != 0 {
				assert.Equal(t, tt.respStatus, resp.StatusCode)
			}

			if tt.errStr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.errStr)
			}

			if tt.expected != nil {
				assert.Equal(t, tt.expected, got)
			}

			if tt.errResp != nil {
				assert.Equal(t, tt.errResp, resp.Error)
			}
		})
	}
}
