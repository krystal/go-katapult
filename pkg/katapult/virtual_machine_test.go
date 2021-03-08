package katapult

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	fixtureVirtualMachineNotFoundErr = "virtual_machine_not_found: No " +
		"virtual machine was found matching any of the criteria provided in " +
		"the arguments"
	fixtureVirtualMachineNotFoundResponseError = &ResponseError{
		Code: "virtual_machine_not_found",
		Description: "No virtual machine was found matching any of the " +
			"criteria provided in the arguments",
		Detail: json.RawMessage(`{}`),
	}

	fixtureVirtualMachineFull = &VirtualMachine{
		ID:                  "vm_t8yomYsG4bccKw5D",
		Name:                "Anvil",
		Hostname:            "anvil",
		FQDN:                "anvil.amce.katapult.cloud",
		Description:         "A heavy anvil-like little box.",
		CreatedAt:           timestampPtr(934834834),
		InitialRootPassword: "eZNHLt8gwtDJSSd59plNMh8S0BEGJZTe",
		State:               "Westeros",
		Zone:                &Zone{ID: "id0"},
		Organization:        &Organization{ID: "id1"},
		Group:               &VirtualMachineGroup{ID: "id2"},
		Package:             &VirtualMachinePackage{ID: "id3"},
		AttachedISO:         &ISO{ID: "id4"},
		Tags:                []*Tag{{ID: "id5"}},
		TagNames:            []string{"heavy"},
		IPAddresses:         []*IPAddress{{ID: "id6"}},
	}
	fixtureVirtualMachineNoID = &VirtualMachine{
		Name:                fixtureVirtualMachineFull.Name,
		Hostname:            fixtureVirtualMachineFull.Hostname,
		FQDN:                fixtureVirtualMachineFull.FQDN,
		Description:         fixtureVirtualMachineFull.Description,
		CreatedAt:           fixtureVirtualMachineFull.CreatedAt,
		InitialRootPassword: fixtureVirtualMachineFull.InitialRootPassword,
		State:               fixtureVirtualMachineFull.State,
		Zone:                fixtureVirtualMachineFull.Zone,
		Organization:        fixtureVirtualMachineFull.Organization,
		Group:               fixtureVirtualMachineFull.Group,
		Package:             fixtureVirtualMachineFull.Package,
		AttachedISO:         fixtureVirtualMachineFull.AttachedISO,
		Tags:                fixtureVirtualMachineFull.Tags,
		TagNames:            fixtureVirtualMachineFull.TagNames,
		IPAddresses:         fixtureVirtualMachineFull.IPAddresses,
	}
	fixtureVirtualMachineNoLookupField = &VirtualMachine{
		Name:                fixtureVirtualMachineFull.Name,
		Hostname:            fixtureVirtualMachineFull.Hostname,
		Description:         fixtureVirtualMachineFull.Description,
		CreatedAt:           fixtureVirtualMachineFull.CreatedAt,
		InitialRootPassword: fixtureVirtualMachineFull.InitialRootPassword,
		State:               fixtureVirtualMachineFull.State,
		Zone:                fixtureVirtualMachineFull.Zone,
		Organization:        fixtureVirtualMachineFull.Organization,
		Group:               fixtureVirtualMachineFull.Group,
		Package:             fixtureVirtualMachineFull.Package,
		AttachedISO:         fixtureVirtualMachineFull.AttachedISO,
		Tags:                fixtureVirtualMachineFull.Tags,
		TagNames:            fixtureVirtualMachineFull.TagNames,
		IPAddresses:         fixtureVirtualMachineFull.IPAddresses,
	}
)

func TestVirtualMachine_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *VirtualMachine
	}{
		{
			name: "empty",
			obj:  &VirtualMachine{},
		},
		{
			name: "full",
			obj:  fixtureVirtualMachineFull,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestNewVirtualMachineLookup(t *testing.T) {
	type args struct {
		idOrFQDN string
	}
	tests := []struct {
		name  string
		args  args
		want  *VirtualMachine
		field FieldName
	}{
		{
			name:  "empty string",
			args:  args{idOrFQDN: ""},
			want:  &VirtualMachine{},
			field: FQDNField,
		},
		{
			name:  "vm_ prefixed ID",
			args:  args{idOrFQDN: "vm_rq1Zdv6aC66bvsU7"},
			want:  &VirtualMachine{ID: "vm_rq1Zdv6aC66bvsU7"},
			field: IDField,
		},
		{
			name:  "fqdn",
			args:  args{idOrFQDN: "noisy-pink-banana"},
			want:  &VirtualMachine{FQDN: "noisy-pink-banana"},
			field: FQDNField,
		},
		{
			name:  "random text",
			args:  args{idOrFQDN: "4lgGRmCPW7XhizQj"},
			want:  &VirtualMachine{FQDN: "4lgGRmCPW7XhizQj"},
			field: FQDNField,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, field := NewVirtualMachineLookup(tt.args.idOrFQDN)

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.field, field)
		})
	}
}

func TestVirtualMachine_lookupReference(t *testing.T) {
	tests := []struct {
		name string
		obj  *VirtualMachine
		want *VirtualMachine
	}{
		{
			name: "nil",
			obj:  nil,
			want: nil,
		},
		{
			name: "empty",
			obj:  &VirtualMachine{},
			want: &VirtualMachine{},
		},
		{
			name: "full",
			obj:  fixtureVirtualMachineFull,
			want: &VirtualMachine{ID: "vm_t8yomYsG4bccKw5D"},
		},
		{
			name: "no ID",
			obj:  fixtureVirtualMachineNoID,
			want: &VirtualMachine{FQDN: "anvil.amce.katapult.cloud"},
		},
		{
			name: "no ID or FQDN",
			obj:  fixtureVirtualMachineNoLookupField,
			want: &VirtualMachine{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.obj.lookupReference()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestVirtualMachine_queryValues(t *testing.T) {
	tests := []struct {
		name string
		obj  *VirtualMachine
	}{
		{
			name: "nil",
			obj:  nil,
		},
		{
			name: "empty",
			obj:  &VirtualMachine{},
		},
		{
			name: "full",
			obj:  fixtureVirtualMachineFull,
		},
		{
			name: "no ID",
			obj:  fixtureVirtualMachineNoID,
		},
		{
			name: "no ID or FQDN",
			obj:  fixtureVirtualMachineNoLookupField,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testQueryableEncoding(t, tt.obj)
		})
	}
}

func TestVirtualMachineStates(t *testing.T) {
	tests := []struct {
		name  string
		enum  VirtualMachineState
		value string
	}{
		{
			name:  "VirtualMachineStopped",
			enum:  VirtualMachineStopped,
			value: "stopped",
		},
		{
			name:  "VirtualMachineFailed",
			enum:  VirtualMachineFailed,
			value: "failed",
		},
		{
			name:  "VirtualMachineStarted",
			enum:  VirtualMachineStarted,
			value: "started",
		},
		{
			name:  "VirtualMachineStarting",
			enum:  VirtualMachineStarting,
			value: "starting",
		},
		{
			name:  "VirtualMachineResetting",
			enum:  VirtualMachineResetting,
			value: "resetting",
		},
		{
			name:  "VirtualMachineMigrating",
			enum:  VirtualMachineMigrating,
			value: "migrating",
		},
		{
			name:  "VirtualMachineStopping",
			enum:  VirtualMachineStopping,
			value: "stopping",
		},
		{
			name:  "VirtualMachineShuttingDown",
			enum:  VirtualMachineShuttingDown,
			value: "shutting_down",
		},
		{
			name:  "VirtualMachineOrphaned",
			enum:  VirtualMachineOrphaned,
			value: "orphaned",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.value, string(tt.enum))
		})
	}
}

func TestVirtualMachineUpdateArguments_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name    string
		obj     *VirtualMachineUpdateArguments
		decoded *VirtualMachineUpdateArguments
	}{
		{
			name: "empty",
			obj:  &VirtualMachineUpdateArguments{},
		},
		{
			name: "full",
			obj: &VirtualMachineUpdateArguments{
				Name:        "db 3",
				Hostname:    "db-3",
				Description: "Database server #3",
				TagNames:    &[]string{"db", "east"},
			},
		},
		{
			name: "empty Tags",
			obj: &VirtualMachineUpdateArguments{
				TagNames: &[]string{},
			},
		},
		{
			name: "null Group",
			obj: &VirtualMachineUpdateArguments{
				Group: NullVirtualMachineGroup,
			},
			decoded: &VirtualMachineUpdateArguments{
				Group: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCustomJSONMarshaling(t, tt.obj, tt.decoded)
		})
	}
}

func Test_virtualMachinesResponseBody_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *virtualMachinesResponseBody
	}{
		{
			name: "empty",
			obj:  &virtualMachinesResponseBody{},
		},
		{
			name: "full",
			obj: &virtualMachinesResponseBody{
				Pagination:      &Pagination{CurrentPage: 345},
				VirtualMachine:  &VirtualMachine{ID: "id1"},
				VirtualMachines: []*VirtualMachine{{ID: "id2"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func Test_virtualMachineChangePackageRequest_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *virtualMachineChangePackageRequest
	}{
		{
			name: "empty",
			obj:  &virtualMachineChangePackageRequest{},
		},
		{
			name: "full",
			obj: &virtualMachineChangePackageRequest{
				VirtualMachine: &VirtualMachine{ID: "id1"},
				Package:        &VirtualMachinePackage{ID: "id2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func Test_virtualMachineUpdateRequest_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *virtualMachineUpdateRequest
	}{
		{
			name: "empty",
			obj:  &virtualMachineUpdateRequest{},
		},
		{
			name: "full",
			obj: &virtualMachineUpdateRequest{
				VirtualMachine: &VirtualMachine{ID: "id1"},
				Properties:     &VirtualMachineUpdateArguments{Name: "hi"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestVirtualMachinesClient_List(t *testing.T) {
	// Correlates to fixtures/virtual_machines_list*.json
	virtualMachinesList := []*VirtualMachine{
		{
			ID:       "vm_t8yomYsG4bccKw5D",
			Name:     "bitter-beautiful-mango",
			Hostname: "bitter-beautiful-mango",
		},
		{
			ID:       "vm_h7bzdXXHa0GvJYMc",
			Name:     "popular-shapely-tank",
			Hostname: "popular-shapely-tank",
		},
		{
			ID:       "vm_1kpkjQeMEI43tztr",
			Name:     "popular-blue-kumquat",
			Hostname: "popular-blue-kumquat",
		},
	}

	type args struct {
		ctx  context.Context
		org  *Organization
		opts *ListOptions
	}
	tests := []struct {
		name           string
		args           args
		want           []*VirtualMachine
		wantQuery      *url.Values
		wantPagination *Pagination
		errStr         string
		errResp        *ResponseError
		respStatus     int
		respBody       []byte
	}{
		{
			name: "by organization ID",
			args: args{
				ctx: context.Background(),
				org: &Organization{ID: "org_O648YDMEYeLmqdmn"},
			},
			want: virtualMachinesList,
			wantQuery: &url.Values{
				"organization[id]": []string{"org_O648YDMEYeLmqdmn"},
			},
			wantPagination: &Pagination{
				CurrentPage: 1,
				TotalPages:  1,
				Total:       3,
				PerPage:     30,
				LargeSet:    false,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machines_list"),
		},
		{
			name: "by organization SubDomain",
			args: args{
				ctx: context.Background(),
				org: &Organization{SubDomain: "acme"},
			},
			want: virtualMachinesList,
			wantQuery: &url.Values{
				"organization[sub_domain]": []string{"acme"},
			},
			wantPagination: &Pagination{
				CurrentPage: 1,
				TotalPages:  1,
				Total:       3,
				PerPage:     30,
				LargeSet:    false,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machines_list"),
		},
		{
			name: "page 1",
			args: args{
				ctx:  context.Background(),
				org:  &Organization{ID: "org_O648YDMEYeLmqdmn"},
				opts: &ListOptions{Page: 1, PerPage: 2},
			},
			want: virtualMachinesList[0:2],
			wantQuery: &url.Values{
				"organization[id]": []string{"org_O648YDMEYeLmqdmn"},
				"page":             []string{"1"},
				"per_page":         []string{"2"},
			},
			wantPagination: &Pagination{
				CurrentPage: 1,
				TotalPages:  2,
				Total:       3,
				PerPage:     2,
				LargeSet:    false,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machines_list_page_1"),
		},
		{
			name: "page 2",
			args: args{
				ctx:  context.Background(),
				org:  &Organization{ID: "org_O648YDMEYeLmqdmn"},
				opts: &ListOptions{Page: 2, PerPage: 2},
			},
			want: virtualMachinesList[2:],
			wantQuery: &url.Values{
				"organization[id]": []string{"org_O648YDMEYeLmqdmn"},
				"page":             []string{"2"},
				"per_page":         []string{"2"},
			},
			wantPagination: &Pagination{
				CurrentPage: 2,
				TotalPages:  2,
				Total:       3,
				PerPage:     2,
				LargeSet:    false,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machines_list_page_2"),
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

			mux.HandleFunc(
				"/core/v1/organizations/_/virtual_machines",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					if tt.wantQuery != nil {
						assert.Equal(t, *tt.wantQuery, r.URL.Query())
					} else {
						qs := queryValues(tt.args.org, tt.args.opts)
						assert.Equal(t, *qs, r.URL.Query())
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.VirtualMachines.List(
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

			if tt.wantPagination != nil {
				assert.Equal(t, tt.wantPagination, resp.Pagination)
			}

			if tt.errResp != nil {
				assert.Equal(t, tt.errResp, resp.Error)
			}
		})
	}
}

func TestVirtualMachinesClient_Get(t *testing.T) {
	type args struct {
		ctx      context.Context
		idOrFQDN string
	}
	tests := []struct {
		name       string
		args       args
		want       *VirtualMachine
		wantQuery  *url.Values
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "by ID",
			args: args{
				ctx:      context.Background(),
				idOrFQDN: "vm_t8yomYsG4bccKw5D",
			},
			want: &VirtualMachine{
				ID:       "vm_t8yomYsG4bccKw5D",
				Name:     "bitter-beautiful-mango",
				Hostname: "bitter-beautiful-mango",
			},
			wantQuery: &url.Values{
				"virtual_machine[id]": []string{"vm_t8yomYsG4bccKw5D"},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machine_get"),
		},
		{
			name: "by FQDN",
			args: args{
				ctx:      context.Background(),
				idOrFQDN: "anvil.amce.katapult.cloud",
			},
			want: &VirtualMachine{
				ID:       "vm_t8yomYsG4bccKw5D",
				Name:     "bitter-beautiful-mango",
				Hostname: "bitter-beautiful-mango",
			},
			wantQuery: &url.Values{
				"virtual_machine[fqdn]": []string{"anvil.amce.katapult.cloud"},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machine_get"),
		},
		{
			name: "non-existent virtual machine",
			args: args{
				ctx:      context.Background(),
				idOrFQDN: "vm_nopethisbegone",
			},
			errStr:     fixtureVirtualMachineNotFoundErr,
			errResp:    fixtureVirtualMachineNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("virtual_machine_not_found_error"),
		},
		{
			name: "virtual machine is in trash",
			args: args{
				ctx:      context.Background(),
				idOrFQDN: "vm_t8yomYsG4bccKw5D",
			},
			errStr:     fixtureObjectInTrashErr,
			errResp:    fixtureObjectInTrashResponseError,
			respStatus: http.StatusNotAcceptable,
			respBody:   fixture("object_in_trash_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx:      nil,
				idOrFQDN: "vm_t8yomYsG4bccKw5D",
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := prepareTestClient()
			defer teardown()

			mux.HandleFunc(
				"/core/v1/virtual_machines/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					if tt.wantQuery != nil {
						assert.Equal(t, *tt.wantQuery, r.URL.Query())
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.VirtualMachines.Get(
				tt.args.ctx, tt.args.idOrFQDN,
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

func TestVirtualMachinesClient_GetByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name       string
		args       args
		want       *VirtualMachine
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "virtual machine",
			args: args{
				ctx: context.Background(),
				id:  "vm_t8yomYsG4bccKw5D",
			},
			want: &VirtualMachine{
				ID:       "vm_t8yomYsG4bccKw5D",
				Name:     "bitter-beautiful-mango",
				Hostname: "bitter-beautiful-mango",
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machine_get"),
		},
		{
			name: "non-existent virtual machine",
			args: args{
				ctx: context.Background(),
				id:  "vm_nopethisbegone",
			},
			errStr:     fixtureVirtualMachineNotFoundErr,
			errResp:    fixtureVirtualMachineNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("virtual_machine_not_found_error"),
		},
		{
			name: "virtual machine is in trash",
			args: args{
				ctx: context.Background(),
				id:  "vm_t8yomYsG4bccKw5D",
			},
			errStr:     fixtureObjectInTrashErr,
			errResp:    fixtureObjectInTrashResponseError,
			respStatus: http.StatusNotAcceptable,
			respBody:   fixture("object_in_trash_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				id:  "vm_t8yomYsG4bccKw5D",
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := prepareTestClient()
			defer teardown()

			mux.HandleFunc(
				"/core/v1/virtual_machines/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					qs := url.Values{
						"virtual_machine[id]": []string{tt.args.id},
					}
					assert.Equal(t, qs, r.URL.Query())

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.VirtualMachines.GetByID(tt.args.ctx, tt.args.id)

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

func TestVirtualMachinesClient_GetByFQDN(t *testing.T) {
	type args struct {
		ctx  context.Context
		fqdn string
	}
	tests := []struct {
		name       string
		args       args
		want       *VirtualMachine
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "virtual machine",
			args: args{
				ctx:  context.Background(),
				fqdn: "vm_t8yomYsG4bccKw5D",
			},
			want: &VirtualMachine{
				ID:       "vm_t8yomYsG4bccKw5D",
				Name:     "bitter-beautiful-mango",
				Hostname: "bitter-beautiful-mango",
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machine_get"),
		},
		{
			name: "non-existent virtual machine",
			args: args{
				ctx:  context.Background(),
				fqdn: "vm_nopethisbegone",
			},
			errStr:     fixtureVirtualMachineNotFoundErr,
			errResp:    fixtureVirtualMachineNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("virtual_machine_not_found_error"),
		},
		{
			name: "virtual machine is in trash",
			args: args{
				ctx:  context.Background(),
				fqdn: "vm_t8yomYsG4bccKw5D",
			},
			errStr:     fixtureObjectInTrashErr,
			errResp:    fixtureObjectInTrashResponseError,
			respStatus: http.StatusNotAcceptable,
			respBody:   fixture("object_in_trash_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx:  nil,
				fqdn: "vm_t8yomYsG4bccKw5D",
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := prepareTestClient()
			defer teardown()

			mux.HandleFunc(
				"/core/v1/virtual_machines/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					qs := url.Values{
						"virtual_machine[fqdn]": []string{tt.args.fqdn},
					}
					assert.Equal(t, qs, r.URL.Query())

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.VirtualMachines.GetByFQDN(
				tt.args.ctx, tt.args.fqdn,
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

func TestVirtualMachinesClient_ChangePackage(t *testing.T) {
	type args struct {
		ctx context.Context
		vm  *VirtualMachine
		pkg *VirtualMachinePackage
	}
	tests := []struct {
		name       string
		args       args
		reqBody    *virtualMachineChangePackageRequest
		want       *Task
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "ID fields",
			args: args{
				ctx: context.Background(),
				vm: &VirtualMachine{
					ID: "vm_t8yomYsG4bccKw5D",
				},
				pkg: &VirtualMachinePackage{
					ID: "vmpkg_XdNPhGXvyt1dnDts",
				},
			},
			reqBody: &virtualMachineChangePackageRequest{
				VirtualMachine: &VirtualMachine{
					ID: "vm_t8yomYsG4bccKw5D",
				},
				Package: &VirtualMachinePackage{
					ID: "vmpkg_XdNPhGXvyt1dnDts",
				},
			},
			want: &Task{
				ID:     "task_7J4vuukDVqAqB4HJ",
				Name:   "Change package",
				Status: TaskPending,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machine_change_package"),
		},
		{
			name: "alt lookup fields",
			args: args{
				ctx: context.Background(),
				vm: &VirtualMachine{
					Name:     "Anvil",
					Hostname: "anvil",
					FQDN:     "anvil.amce.katapult.cloud",
				},
				pkg: &VirtualMachinePackage{
					Name:      "X-Small",
					Permalink: "xsmall",
				},
			},
			reqBody: &virtualMachineChangePackageRequest{
				VirtualMachine: &VirtualMachine{
					FQDN: "anvil.amce.katapult.cloud",
				},
				Package: &VirtualMachinePackage{
					Permalink: "xsmall",
				},
			},
			want: &Task{
				ID:     "task_7J4vuukDVqAqB4HJ",
				Name:   "Change package",
				Status: TaskPending,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machine_change_package"),
		},
		{
			name: "full fields",
			args: args{
				ctx: context.Background(),
				vm: &VirtualMachine{
					ID:                  "vm_t8yomYsG4bccKw5D",
					Name:                "Anvil",
					Hostname:            "anvil",
					FQDN:                "anvil.amce.katapult.cloud",
					CreatedAt:           timestampPtr(934834834),
					InitialRootPassword: "eZNHLt8gwtDJSSd59plNMh8S0BEGJZTe",
					State:               "Westeros",
					Zone:                &Zone{ID: "id0"},
					Organization:        &Organization{ID: "id1"},
					Group:               &VirtualMachineGroup{ID: "id2"},
					Package:             &VirtualMachinePackage{ID: "id3"},
					AttachedISO:         &ISO{ID: "id4"},
					Tags:                []*Tag{{ID: "id5"}},
					IPAddresses:         []*IPAddress{{ID: "id6"}},
				},
				pkg: &VirtualMachinePackage{
					ID:            "vmpkg_XdNPhGXvyt1dnDts",
					Name:          "X-Small",
					Permalink:     "xsmall",
					CPUCores:      504684,
					IPv4Addresses: 322134,
					MemoryInGB:    953603,
					StorageInGB:   853121,
					Privacy:       "priv",
					Icon:          &Attachment{URL: "url"},
				},
			},
			reqBody: &virtualMachineChangePackageRequest{
				VirtualMachine: &VirtualMachine{
					ID: "vm_t8yomYsG4bccKw5D",
				},
				Package: &VirtualMachinePackage{
					ID: "vmpkg_XdNPhGXvyt1dnDts",
				},
			},
			want: &Task{
				ID:     "task_7J4vuukDVqAqB4HJ",
				Name:   "Change package",
				Status: TaskPending,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machine_change_package"),
		},
		{
			name: "non-existent virtual machine",
			args: args{
				ctx: context.Background(),
				vm: &VirtualMachine{
					ID: "vm_t8yomYsG4bccKw5D",
				},
				pkg: &VirtualMachinePackage{
					ID: "vmpkg_XdNPhGXvyt1dnDts",
				},
			},
			errStr:     fixtureVirtualMachineNotFoundErr,
			errResp:    fixtureVirtualMachineNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("virtual_machine_not_found_error"),
		},
		{
			name: "virtual machine is in trash",
			args: args{
				ctx: context.Background(),
				vm: &VirtualMachine{
					ID: "vm_t8yomYsG4bccKw5D",
				},
				pkg: &VirtualMachinePackage{
					ID: "vmpkg_XdNPhGXvyt1dnDts",
				},
			},
			errStr:     fixtureObjectInTrashErr,
			errResp:    fixtureObjectInTrashResponseError,
			respStatus: http.StatusNotAcceptable,
			respBody:   fixture("object_in_trash_error"),
		},
		{
			name: "non-existent virtual machine package",
			args: args{
				ctx: context.Background(),
				vm: &VirtualMachine{
					ID: "vm_t8yomYsG4bccKw5D",
				},
				pkg: &VirtualMachinePackage{
					ID: "vmpkg_XdNPhGXvyt1dnDts",
				},
			},
			errStr:     fixturePackageNotFoundErr,
			errResp:    fixturePackageNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("package_not_found_error"),
		},
		{
			name: "permission_denied",
			args: args{
				ctx: context.Background(),
				vm: &VirtualMachine{
					ID: "vm_t8yomYsG4bccKw5D",
				},
				pkg: &VirtualMachinePackage{
					ID: "vmpkg_XdNPhGXvyt1dnDts",
				},
			},
			errStr:     fixturePermissionDeniedErr,
			errResp:    fixturePermissionDeniedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("permission_denied_error"),
		},
		{
			name: "task_queueing_error",
			args: args{
				ctx: context.Background(),
				vm: &VirtualMachine{
					ID: "vm_t8yomYsG4bccKw5D",
				},
				pkg: &VirtualMachinePackage{
					ID: "vmpkg_XdNPhGXvyt1dnDts",
				},
			},
			errStr:     fixtureTaskQueueingErrorErr,
			errResp:    fixtureTaskQueueingErrorResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("task_queueing_error"),
		},
		{
			name: "nil virtual machine",
			args: args{
				ctx: context.Background(),
				vm:  nil,
				pkg: &VirtualMachinePackage{
					ID: "vmpkg_XdNPhGXvyt1dnDts",
				},
			},
			reqBody: &virtualMachineChangePackageRequest{
				Package: &VirtualMachinePackage{
					ID: "vmpkg_XdNPhGXvyt1dnDts",
				},
			},
			errStr:     fixtureVirtualMachineNotFoundErr,
			errResp:    fixtureVirtualMachineNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("virtual_machine_not_found_error"),
		},
		{
			name: "nil virtual machine package",
			args: args{
				ctx: context.Background(),
				vm: &VirtualMachine{
					ID: "vm_t8yomYsG4bccKw5D",
				},
				pkg: nil,
			},
			reqBody: &virtualMachineChangePackageRequest{
				VirtualMachine: &VirtualMachine{
					ID: "vm_t8yomYsG4bccKw5D",
				},
			},
			errStr:     fixturePackageNotFoundErr,
			errResp:    fixturePackageNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("package_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				vm: &VirtualMachine{
					ID: "vm_t8yomYsG4bccKw5D",
				},
				pkg: &VirtualMachinePackage{
					ID: "vmpkg_XdNPhGXvyt1dnDts",
				},
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := prepareTestClient()
			defer teardown()

			mux.HandleFunc(
				"/core/v1/virtual_machines/_/package",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "PUT", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					if tt.reqBody != nil {
						reqBody := &virtualMachineChangePackageRequest{}
						err := strictUmarshal(r.Body, reqBody)
						assert.NoError(t, err)
						assert.Equal(t, tt.reqBody, reqBody)
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.VirtualMachines.ChangePackage(
				tt.args.ctx, tt.args.vm, tt.args.pkg,
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

func TestVirtualMachinesClient_Update(t *testing.T) {
	vmArgs := &VirtualMachineUpdateArguments{
		Name:     "Anvil Next",
		Hostname: "anvil-next",
	}

	type args struct {
		ctx  context.Context
		vm   *VirtualMachine
		args *VirtualMachineUpdateArguments
	}
	tests := []struct {
		name       string
		args       args
		reqBody    *virtualMachineUpdateRequest
		want       *VirtualMachine
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "by ID",
			args: args{
				ctx:  context.Background(),
				vm:   fixtureVirtualMachineFull,
				args: vmArgs,
			},
			reqBody: &virtualMachineUpdateRequest{
				VirtualMachine: &VirtualMachine{
					ID: fixtureVirtualMachineFull.ID,
				},
				Properties: vmArgs,
			},
			want: &VirtualMachine{
				ID:       "vm_t8yomYsG4bccKw5D",
				Name:     "Anvil Next",
				Hostname: "anvil-next",
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machine_update"),
		},
		{
			name: "by FQDN",
			args: args{
				ctx:  context.Background(),
				vm:   fixtureVirtualMachineNoID,
				args: vmArgs,
			},
			reqBody: &virtualMachineUpdateRequest{
				VirtualMachine: &VirtualMachine{
					FQDN: fixtureVirtualMachineNoID.FQDN,
				},
				Properties: vmArgs,
			},
			want: &VirtualMachine{
				ID:       "vm_t8yomYsG4bccKw5D",
				Name:     "Anvil Next",
				Hostname: "anvil-next",
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machine_update"),
		},
		{
			name: "non-existent virtual machine",
			args: args{
				ctx:  context.Background(),
				vm:   fixtureVirtualMachineFull,
				args: vmArgs,
			},
			errStr:     fixtureVirtualMachineNotFoundErr,
			errResp:    fixtureVirtualMachineNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("virtual_machine_not_found_error"),
		},
		{
			name: "virtual machine is in trash",
			args: args{
				ctx:  context.Background(),
				vm:   fixtureVirtualMachineFull,
				args: vmArgs,
			},
			errStr:     fixtureObjectInTrashErr,
			errResp:    fixtureObjectInTrashResponseError,
			respStatus: http.StatusNotAcceptable,
			respBody:   fixture("object_in_trash_error"),
		},
		{
			name: "permission_denied",
			args: args{
				ctx:  context.Background(),
				vm:   fixtureVirtualMachineFull,
				args: vmArgs,
			},
			errStr:     fixturePermissionDeniedErr,
			errResp:    fixturePermissionDeniedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("permission_denied_error"),
		},
		{
			name: "nil virtual machine",
			args: args{
				ctx:  context.Background(),
				vm:   nil,
				args: vmArgs,
			},
			reqBody: &virtualMachineUpdateRequest{
				Properties: vmArgs,
			},
			errStr:     fixtureVirtualMachineNotFoundErr,
			errResp:    fixtureVirtualMachineNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("virtual_machine_not_found_error"),
		},
		{
			name: "nil virtual machine update arguments",
			args: args{
				ctx:  context.Background(),
				vm:   fixtureVirtualMachineFull,
				args: nil,
			},
			reqBody: &virtualMachineUpdateRequest{
				VirtualMachine: &VirtualMachine{
					ID: "vm_t8yomYsG4bccKw5D",
				},
			},
			want: &VirtualMachine{
				ID:       "vm_t8yomYsG4bccKw5D",
				Name:     "bitter-beautiful-mango",
				Hostname: "bitter-beautiful-mango",
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machine_get"),
		},
		{
			name: "nil context",
			args: args{
				ctx:  nil,
				vm:   fixtureVirtualMachineFull,
				args: vmArgs,
			},

			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := prepareTestClient()
			defer teardown()

			mux.HandleFunc(
				"/core/v1/virtual_machines/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "PATCH", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					if tt.reqBody != nil {
						reqBody := &virtualMachineUpdateRequest{}
						err := strictUmarshal(r.Body, reqBody)
						assert.NoError(t, err)
						assert.Equal(t, tt.reqBody, reqBody)
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.VirtualMachines.Update(
				tt.args.ctx, tt.args.vm, tt.args.args,
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

func TestVirtualMachinesClient_Delete(t *testing.T) {
	type args struct {
		ctx context.Context
		vm  *VirtualMachine
	}
	tests := []struct {
		name       string
		args       args
		want       *TrashObject
		wantQuery  *url.Values
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "by ID",
			args: args{
				ctx: context.Background(),
				vm:  &VirtualMachine{ID: "vm_t8yomYsG4bccKw5D"},
			},
			want: &TrashObject{
				ID:        "trsh_AmjmS73QadkAZqoE",
				KeepUntil: timestampPtr(1599672014),
			},
			wantQuery: &url.Values{
				"virtual_machine[id]": []string{"vm_t8yomYsG4bccKw5D"},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machine_delete"),
		},
		{
			name: "by FQDN",
			args: args{
				ctx: context.Background(),
				vm:  &VirtualMachine{FQDN: "anvil.amce.katapult.cloud"},
			},
			want: &TrashObject{
				ID:        "trsh_AmjmS73QadkAZqoE",
				KeepUntil: timestampPtr(1599672014),
			},
			wantQuery: &url.Values{
				"virtual_machine[fqdn]": []string{"anvil.amce.katapult.cloud"},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machine_delete"),
		},
		{
			name: "non-existent virtual machine",
			args: args{
				ctx: context.Background(),
				vm:  &VirtualMachine{ID: "vm_t8yomYsG4bccKw5D"},
			},
			errStr:     fixtureVirtualMachineNotFoundErr,
			errResp:    fixtureVirtualMachineNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("virtual_machine_not_found_error"),
		},
		{
			name: "permission_denied",
			args: args{
				ctx: context.Background(),
				vm:  &VirtualMachine{ID: "vm_t8yomYsG4bccKw5D"},
			},
			errStr:     fixturePermissionDeniedErr,
			errResp:    fixturePermissionDeniedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("permission_denied_error"),
		},
		{
			name: "nil virtual machine",
			args: args{
				ctx: context.Background(),
				vm:  nil,
			},
			errStr:     fixtureVirtualMachineNotFoundErr,
			errResp:    fixtureVirtualMachineNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("virtual_machine_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				vm:  &VirtualMachine{ID: "vm_t8yomYsG4bccKw5D"},
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := prepareTestClient()
			defer teardown()

			mux.HandleFunc(
				"/core/v1/virtual_machines/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "DELETE", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					if tt.wantQuery != nil {
						assert.Equal(t, *tt.wantQuery, r.URL.Query())
					} else {
						assert.Equal(t,
							*tt.args.vm.queryValues(), r.URL.Query(),
						)
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.VirtualMachines.Delete(tt.args.ctx, tt.args.vm)

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

func TestVirtualMachinesClient_Start(t *testing.T) {
	type args struct {
		ctx context.Context
		vm  *VirtualMachine
	}
	tests := []struct {
		name       string
		args       args
		want       *Task
		wantQuery  *url.Values
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "by ID",
			args: args{
				ctx: context.Background(),
				vm:  &VirtualMachine{ID: "vm_t8yomYsG4bccKw5D"},
			},
			want: &Task{
				ID:     "task_otL5Dkr3bi40yn9h",
				Name:   "Start virtual machine",
				Status: "pending",
			},
			wantQuery: &url.Values{
				"virtual_machine[id]": []string{"vm_t8yomYsG4bccKw5D"},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machine_start"),
		},
		{
			name: "by FQDN",
			args: args{
				ctx: context.Background(),
				vm:  &VirtualMachine{FQDN: "anvil.amce.katapult.cloud"},
			},
			want: &Task{
				ID:     "task_otL5Dkr3bi40yn9h",
				Name:   "Start virtual machine",
				Status: "pending",
			},
			wantQuery: &url.Values{
				"virtual_machine[fqdn]": []string{"anvil.amce.katapult.cloud"},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machine_start"),
		},
		{
			name: "non-existent virtual machine",
			args: args{
				ctx: context.Background(),
				vm:  &VirtualMachine{ID: "vm_nopethisbegone"},
			},
			errStr:     fixtureVirtualMachineNotFoundErr,
			errResp:    fixtureVirtualMachineNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("virtual_machine_not_found_error"),
		},
		{
			name: "virtual machine is in trash",
			args: args{
				ctx: context.Background(),
				vm:  &VirtualMachine{ID: "vm_t8yomYsG4bccKw5D"},
			},
			errStr:     fixtureObjectInTrashErr,
			errResp:    fixtureObjectInTrashResponseError,
			respStatus: http.StatusNotAcceptable,
			respBody:   fixture("object_in_trash_error"),
		},
		{
			name: "error queuing task",
			args: args{
				ctx: context.Background(),
				vm:  &VirtualMachine{ID: "vm_t8yomYsG4bccKw5D"},
			},
			errStr:     fixtureTaskQueueingErrorErr,
			errResp:    fixtureTaskQueueingErrorResponseError,
			respStatus: http.StatusNotAcceptable,
			respBody:   fixture("task_queueing_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				vm:  &VirtualMachine{ID: "vm_t8yomYsG4bccKw5D"},
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := prepareTestClient()
			defer teardown()

			mux.HandleFunc(
				"/core/v1/virtual_machines/_/start",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "POST", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					if tt.wantQuery != nil {
						assert.Equal(t, *tt.wantQuery, r.URL.Query())
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.VirtualMachines.Start(
				tt.args.ctx, tt.args.vm,
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

func TestVirtualMachinesClient_Stop(t *testing.T) {
	type args struct {
		ctx context.Context
		vm  *VirtualMachine
	}
	tests := []struct {
		name       string
		args       args
		want       *Task
		wantQuery  *url.Values
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "by ID",
			args: args{
				ctx: context.Background(),
				vm:  &VirtualMachine{ID: "vm_t8yomYsG4bccKw5D"},
			},
			want: &Task{
				ID:     "task_UWMEbeWyZx3qZIzK",
				Name:   "Stop virtual machine",
				Status: "pending",
			},
			wantQuery: &url.Values{
				"virtual_machine[id]": []string{"vm_t8yomYsG4bccKw5D"},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machine_stop"),
		},
		{
			name: "by FQDN",
			args: args{
				ctx: context.Background(),
				vm:  &VirtualMachine{FQDN: "anvil.amce.katapult.cloud"},
			},
			want: &Task{
				ID:     "task_UWMEbeWyZx3qZIzK",
				Name:   "Stop virtual machine",
				Status: "pending",
			},
			wantQuery: &url.Values{
				"virtual_machine[fqdn]": []string{"anvil.amce.katapult.cloud"},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machine_stop"),
		},
		{
			name: "non-existent virtual machine",
			args: args{
				ctx: context.Background(),
				vm:  &VirtualMachine{ID: "vm_nopethisbegone"},
			},
			errStr:     fixtureVirtualMachineNotFoundErr,
			errResp:    fixtureVirtualMachineNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("virtual_machine_not_found_error"),
		},
		{
			name: "virtual machine is in trash",
			args: args{
				ctx: context.Background(),
				vm:  &VirtualMachine{ID: "vm_t8yomYsG4bccKw5D"},
			},
			errStr:     fixtureObjectInTrashErr,
			errResp:    fixtureObjectInTrashResponseError,
			respStatus: http.StatusNotAcceptable,
			respBody:   fixture("object_in_trash_error"),
		},
		{
			name: "error queuing task",
			args: args{
				ctx: context.Background(),
				vm:  &VirtualMachine{ID: "vm_t8yomYsG4bccKw5D"},
			},
			errStr:     fixtureTaskQueueingErrorErr,
			errResp:    fixtureTaskQueueingErrorResponseError,
			respStatus: http.StatusNotAcceptable,
			respBody:   fixture("task_queueing_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				vm:  &VirtualMachine{ID: "vm_t8yomYsG4bccKw5D"},
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := prepareTestClient()
			defer teardown()

			mux.HandleFunc(
				"/core/v1/virtual_machines/_/stop",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "POST", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					if tt.wantQuery != nil {
						assert.Equal(t, *tt.wantQuery, r.URL.Query())
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.VirtualMachines.Stop(
				tt.args.ctx, tt.args.vm,
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

func TestVirtualMachinesClient_Shutdown(t *testing.T) {
	type args struct {
		ctx context.Context
		vm  *VirtualMachine
	}
	tests := []struct {
		name       string
		args       args
		want       *Task
		wantQuery  *url.Values
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "by ID",
			args: args{
				ctx: context.Background(),
				vm:  &VirtualMachine{ID: "vm_t8yomYsG4bccKw5D"},
			},
			want: &Task{
				ID:     "task_zSdnw8Ocz8QAQTZK",
				Name:   "Shutdown virtual machine",
				Status: "pending",
			},
			wantQuery: &url.Values{
				"virtual_machine[id]": []string{"vm_t8yomYsG4bccKw5D"},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machine_shutdown"),
		},
		{
			name: "by FQDN",
			args: args{
				ctx: context.Background(),
				vm:  &VirtualMachine{FQDN: "anvil.amce.katapult.cloud"},
			},
			want: &Task{
				ID:     "task_zSdnw8Ocz8QAQTZK",
				Name:   "Shutdown virtual machine",
				Status: "pending",
			},
			wantQuery: &url.Values{
				"virtual_machine[fqdn]": []string{"anvil.amce.katapult.cloud"},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machine_shutdown"),
		},
		{
			name: "non-existent virtual machine",
			args: args{
				ctx: context.Background(),
				vm:  &VirtualMachine{ID: "vm_nopethisbegone"},
			},
			errStr:     fixtureVirtualMachineNotFoundErr,
			errResp:    fixtureVirtualMachineNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("virtual_machine_not_found_error"),
		},
		{
			name: "virtual machine is in trash",
			args: args{
				ctx: context.Background(),
				vm:  &VirtualMachine{ID: "vm_t8yomYsG4bccKw5D"},
			},
			errStr:     fixtureObjectInTrashErr,
			errResp:    fixtureObjectInTrashResponseError,
			respStatus: http.StatusNotAcceptable,
			respBody:   fixture("object_in_trash_error"),
		},
		{
			name: "error queuing task",
			args: args{
				ctx: context.Background(),
				vm:  &VirtualMachine{ID: "vm_t8yomYsG4bccKw5D"},
			},
			errStr:     fixtureTaskQueueingErrorErr,
			errResp:    fixtureTaskQueueingErrorResponseError,
			respStatus: http.StatusNotAcceptable,
			respBody:   fixture("task_queueing_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				vm:  &VirtualMachine{ID: "vm_t8yomYsG4bccKw5D"},
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := prepareTestClient()
			defer teardown()

			mux.HandleFunc(
				"/core/v1/virtual_machines/_/shutdown",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "POST", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					if tt.wantQuery != nil {
						assert.Equal(t, *tt.wantQuery, r.URL.Query())
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.VirtualMachines.Shutdown(
				tt.args.ctx, tt.args.vm,
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

func TestVirtualMachinesClient_Reset(t *testing.T) {
	type args struct {
		ctx context.Context
		vm  *VirtualMachine
	}
	tests := []struct {
		name       string
		args       args
		want       *Task
		wantQuery  *url.Values
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "by ID",
			args: args{
				ctx: context.Background(),
				vm:  &VirtualMachine{ID: "vm_t8yomYsG4bccKw5D"},
			},
			want: &Task{
				ID:     "task_vZYARjrFue1Or2pt",
				Name:   "Reset virtual machine",
				Status: "pending",
			},
			wantQuery: &url.Values{
				"virtual_machine[id]": []string{"vm_t8yomYsG4bccKw5D"},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machine_reset"),
		},
		{
			name: "by FQDN",
			args: args{
				ctx: context.Background(),
				vm:  &VirtualMachine{FQDN: "anvil.amce.katapult.cloud"},
			},
			want: &Task{
				ID:     "task_vZYARjrFue1Or2pt",
				Name:   "Reset virtual machine",
				Status: "pending",
			},
			wantQuery: &url.Values{
				"virtual_machine[fqdn]": []string{"anvil.amce.katapult.cloud"},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machine_reset"),
		},
		{
			name: "non-existent virtual machine",
			args: args{
				ctx: context.Background(),
				vm:  &VirtualMachine{ID: "vm_nopethisbegone"},
			},
			errStr:     fixtureVirtualMachineNotFoundErr,
			errResp:    fixtureVirtualMachineNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("virtual_machine_not_found_error"),
		},
		{
			name: "virtual machine is in trash",
			args: args{
				ctx: context.Background(),
				vm:  &VirtualMachine{ID: "vm_t8yomYsG4bccKw5D"},
			},
			errStr:     fixtureObjectInTrashErr,
			errResp:    fixtureObjectInTrashResponseError,
			respStatus: http.StatusNotAcceptable,
			respBody:   fixture("object_in_trash_error"),
		},
		{
			name: "error queuing task",
			args: args{
				ctx: context.Background(),
				vm:  &VirtualMachine{ID: "vm_t8yomYsG4bccKw5D"},
			},
			errStr:     fixtureTaskQueueingErrorErr,
			errResp:    fixtureTaskQueueingErrorResponseError,
			respStatus: http.StatusNotAcceptable,
			respBody:   fixture("task_queueing_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				vm:  &VirtualMachine{ID: "vm_t8yomYsG4bccKw5D"},
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := prepareTestClient()
			defer teardown()

			mux.HandleFunc(
				"/core/v1/virtual_machines/_/reset",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "POST", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					if tt.wantQuery != nil {
						assert.Equal(t, *tt.wantQuery, r.URL.Query())
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.VirtualMachines.Reset(
				tt.args.ctx, tt.args.vm,
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
