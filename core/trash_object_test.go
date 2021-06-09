package core

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/krystal/go-katapult"
	"github.com/krystal/go-katapult/internal/testclient"
	"github.com/stretchr/testify/assert"
)

var (
	fixtureTrashObjectFull = &TrashObject{
		ID:         "trsh_NRhMtSdZbNRVafj3",
		KeepUntil:  timestampPtr(93043),
		ObjectID:   "vm_pxrHYFpOW88ka39h",
		ObjectType: "VirtualMachine",
	}

	fixtureTrashObjectNotFoundErr = "trash_object_not_found: No trash object " +
		"was found matching any of the criteria provided in the arguments"
	fixtureTrashObjectNotFoundResponseError = &katapult.ResponseError{
		Code: "trash_object_not_found",
		Description: "No trash object was found matching any of the criteria " +
			"provided in the arguments",
		Detail: json.RawMessage(`{}`),
	}
)

func TestClient_TrashObjects(t *testing.T) {
	c := New(&testclient.Client{})

	assert.IsType(t, &TrashObjectsClient{}, c.TrashObjects)
}

func TestTrashObject_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *TrashObject
	}{
		{
			name: "empty",
			obj:  &TrashObject{},
		},
		{
			name: "full",
			obj:  fixtureTrashObjectFull,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestTrashObject_Ref(t *testing.T) {
	tests := []struct {
		name string
		obj  TrashObject
		want TrashObjectRef
	}{
		{
			name: "empty",
			obj:  TrashObject{},
			want: TrashObjectRef{},
		},
		{
			name: "full",
			obj: TrashObject{
				ID:       "trsh_NRhMtSdZbNRVafj3",
				ObjectID: "vm_pxrHYFpOW88ka39h",
			},
			want: TrashObjectRef{ID: "trsh_NRhMtSdZbNRVafj3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.obj.Ref()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTrashObjectRef_queryValues(t *testing.T) {
	tests := []struct {
		name string
		obj  TrashObjectRef
	}{
		{
			name: "empty",
			obj:  TrashObjectRef{},
		},
		{
			name: "full",
			obj: TrashObjectRef{
				ID:       "trsh_NRhMtSdZbNRVafj3",
				ObjectID: "vm_pxrHYFpOW88ka39h",
			},
		},
		{
			name: "just ID",
			obj: TrashObjectRef{
				ID: "trsh_NRhMtSdZbNRVafj3",
			},
		},
		{
			name: "just ObjectID",
			obj: TrashObjectRef{
				ObjectID: "vm_pxrHYFpOW88ka39h",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testQueryableEncoding(t, tt.obj)
		})
	}
}

func Test_trashObjectsResponseBody_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *trashObjectsResponseBody
	}{
		{
			name: "empty",
			obj:  &trashObjectsResponseBody{},
		},
		{
			name: "full",
			obj: &trashObjectsResponseBody{
				Pagination:   &katapult.Pagination{CurrentPage: 345},
				TrashObject:  &TrashObject{ID: "id1"},
				TrashObjects: []*TrashObject{{ID: "id2"}},
				Task:         &Task{ID: "id3"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestTrashObjectsClient_List(t *testing.T) {
	// Correlates to fixtures/trash_objects_list*.json
	trashObjectesList := []*TrashObject{
		{
			ID:         "trsh_hkW1SMq0Bn8yNrRx",
			KeepUntil:  timestampPtr(1610039056),
			ObjectID:   "vm_KTKc6pwFxLjJ40QY",
			ObjectType: "VirtualMachine",
		},
		{
			ID:         "trsh_WX7ZTIdCb2gZ0PQ9",
			KeepUntil:  timestampPtr(1610039191),
			ObjectID:   "disk_reWHTQagpihFhSuh",
			ObjectType: "Disk",
		},
		{
			ID:         "trsh_h6An31KwJU0jOq5y",
			KeepUntil:  timestampPtr(1610039283),
			ObjectID:   "fsv_f9WF2pMAb5BY8vlK",
			ObjectType: "FileStorageVolume",
		},
	}

	type args struct {
		ctx  context.Context
		org  OrganizationRef
		opts *ListOptions
	}
	tests := []struct {
		name           string
		args           args
		want           []*TrashObject
		wantPagination *katapult.Pagination
		errStr         string
		errResp        *katapult.ResponseError
		respStatus     int
		respBody       []byte
	}{
		{
			name: "by organization ID",
			args: args{
				ctx: context.Background(),
				org: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
			},
			want: trashObjectesList,
			wantPagination: &katapult.Pagination{
				CurrentPage: 1,
				TotalPages:  1,
				Total:       3,
				PerPage:     30,
				LargeSet:    false,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("trash_objects_list"),
		},
		{
			name: "by organization SubDomain",
			args: args{
				ctx: context.Background(),
				org: OrganizationRef{SubDomain: "acme"},
			},
			want: trashObjectesList,
			wantPagination: &katapult.Pagination{
				CurrentPage: 1,
				TotalPages:  1,
				Total:       3,
				PerPage:     30,
				LargeSet:    false,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("trash_objects_list"),
		},
		{
			name: "page 1",
			args: args{
				ctx:  context.Background(),
				org:  OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				opts: &ListOptions{Page: 1, PerPage: 2},
			},
			want: trashObjectesList[0:2],
			wantPagination: &katapult.Pagination{
				CurrentPage: 1,
				TotalPages:  2,
				Total:       3,
				PerPage:     2,
				LargeSet:    false,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("trash_objects_list_page_1"),
		},
		{
			name: "page 2",
			args: args{
				ctx:  context.Background(),
				org:  OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				opts: &ListOptions{Page: 2, PerPage: 2},
			},
			want: trashObjectesList[2:],
			wantPagination: &katapult.Pagination{
				CurrentPage: 2,
				TotalPages:  2,
				Total:       3,
				PerPage:     2,
				LargeSet:    false,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("trash_objects_list_page_2"),
		},
		{
			name: "invalid API token response",
			args: args{
				ctx: context.Background(),
				org: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
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
				org: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
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
				org: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
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
				org: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewTrashObjectsClient(rm)

			mux.HandleFunc(
				"/core/v1/organizations/_/trash_objects",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					qs := queryValues(tt.args.org, tt.args.opts)
					assert.Equal(t, *qs, r.URL.Query())

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.List(
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

func TestTrashObjectsClient_Get(t *testing.T) {
	type args struct {
		ctx context.Context
		ref TrashObjectRef
	}
	tests := []struct {
		name       string
		args       args
		want       *TrashObject
		wantQuery  *url.Values
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "by ID",
			args: args{
				ctx: context.Background(),
				ref: TrashObjectRef{ID: "trsh_hkW1SMq0Bn8yNrRx"},
			},
			want: &TrashObject{
				ID:         "trsh_hkW1SMq0Bn8yNrRx",
				KeepUntil:  timestampPtr(1610039056),
				ObjectID:   "vm_KTKc6pwFxLjJ40QY",
				ObjectType: "VirtualMachine",
			},
			wantQuery: &url.Values{
				"trash_object[id]": []string{"trsh_hkW1SMq0Bn8yNrRx"},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("trash_object_get"),
		},
		{
			name: "by ObjectID",
			args: args{
				ctx: context.Background(),
				ref: TrashObjectRef{ObjectID: "vm_KTKc6pwFxLjJ40QY"},
			},
			want: &TrashObject{
				ID:         "trsh_hkW1SMq0Bn8yNrRx",
				KeepUntil:  timestampPtr(1610039056),
				ObjectID:   "vm_KTKc6pwFxLjJ40QY",
				ObjectType: "VirtualMachine",
			},
			wantQuery: &url.Values{
				"trash_object[object_id]": []string{"vm_KTKc6pwFxLjJ40QY"},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("trash_object_get"),
		},
		{
			name: "non-existent trash object",
			args: args{
				ctx: context.Background(),
				ref: TrashObjectRef{ID: "trsh_nopethisbegone"},
			},
			errStr:     fixtureTrashObjectNotFoundErr,
			errResp:    fixtureTrashObjectNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("trash_object_not_found_error"),
		},
		{
			name: "empty idOrObjectID",
			args: args{
				ctx: context.Background(),
				ref: TrashObjectRef{},
			},
			errStr:     fixtureTrashObjectNotFoundErr,
			errResp:    fixtureTrashObjectNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("trash_object_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				ref: TrashObjectRef{ID: "trsh_hkW1SMq0Bn8yNrRx"},
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewTrashObjectsClient(rm)

			mux.HandleFunc(
				"/core/v1/trash_objects/_",
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

			got, resp, err := c.Get(
				tt.args.ctx, tt.args.ref,
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

func TestTrashObjectsClient_GetByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name       string
		args       args
		want       *TrashObject
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "trash object",
			args: args{
				ctx: context.Background(),
				id:  "trsh_hkW1SMq0Bn8yNrRx",
			},
			want: &TrashObject{
				ID:         "trsh_hkW1SMq0Bn8yNrRx",
				KeepUntil:  timestampPtr(1610039056),
				ObjectID:   "vm_KTKc6pwFxLjJ40QY",
				ObjectType: "VirtualMachine",
			},
			respStatus: http.StatusOK,
			respBody:   fixture("trash_object_get"),
		},
		{
			name: "non-existent trash object",
			args: args{
				ctx: context.Background(),
				id:  "trsh_nopethisbegone",
			},
			errStr:     fixtureTrashObjectNotFoundErr,
			errResp:    fixtureTrashObjectNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("trash_object_not_found_error"),
		},
		{
			name: "empty ID",
			args: args{
				ctx: context.Background(),
				id:  "",
			},
			errStr:     fixtureTrashObjectNotFoundErr,
			errResp:    fixtureTrashObjectNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("trash_object_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				id:  "trsh_hkW1SMq0Bn8yNrRx",
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewTrashObjectsClient(rm)

			mux.HandleFunc(
				"/core/v1/trash_objects/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					qs := url.Values{}
					if tt.args.id != "" {
						qs["trash_object[id]"] = []string{tt.args.id}
					}
					assert.Equal(t, qs, r.URL.Query())

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.GetByID(tt.args.ctx, tt.args.id)

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

func TestTrashObjectsClient_GetByObjectID(t *testing.T) {
	type args struct {
		ctx      context.Context
		objectID string
	}
	tests := []struct {
		name       string
		args       args
		want       *TrashObject
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "trash object",
			args: args{
				ctx:      context.Background(),
				objectID: "vm_KTKc6pwFxLjJ40QY",
			},
			want: &TrashObject{
				ID:         "trsh_hkW1SMq0Bn8yNrRx",
				KeepUntil:  timestampPtr(1610039056),
				ObjectID:   "vm_KTKc6pwFxLjJ40QY",
				ObjectType: "VirtualMachine",
			},
			respStatus: http.StatusOK,
			respBody:   fixture("trash_object_get"),
		},
		{
			name: "non-existent trash object",
			args: args{
				ctx:      context.Background(),
				objectID: "vm_nopethisisnothere",
			},
			errStr:     fixtureTrashObjectNotFoundErr,
			errResp:    fixtureTrashObjectNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("trash_object_not_found_error"),
		},
		{
			name: "empty ObjectID",
			args: args{
				ctx:      context.Background(),
				objectID: "",
			},
			errStr:     fixtureTrashObjectNotFoundErr,
			errResp:    fixtureTrashObjectNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("trash_object_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx:      nil,
				objectID: "vm_KTKc6pwFxLjJ40QY",
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewTrashObjectsClient(rm)

			mux.HandleFunc(
				"/core/v1/trash_objects/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					qs := url.Values{}
					if tt.args.objectID != "" {
						qs["trash_object[object_id]"] = []string{
							tt.args.objectID,
						}
					}
					assert.Equal(t, qs, r.URL.Query())

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.GetByObjectID(
				tt.args.ctx, tt.args.objectID,
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

func TestTrashObjectsClient_Purge(t *testing.T) {
	type args struct {
		ctx context.Context
		ref TrashObjectRef
	}
	tests := []struct {
		name       string
		args       args
		want       *Task
		wantQuery  *url.Values
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "by ID",
			args: args{
				ctx: context.Background(),
				ref: TrashObjectRef{ID: "trsh_hkW1SMq0Bn8yNrRx"},
			},
			want: &Task{
				ID:     "task_Fq0vMXkSkKkGU3ut",
				Name:   "Purge items from trash",
				Status: "pending",
			},
			wantQuery: &url.Values{
				"trash_object[id]": []string{"trsh_hkW1SMq0Bn8yNrRx"},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("trash_object_purge"),
		},
		{
			name: "by ObjectID",
			args: args{
				ctx: context.Background(),
				ref: TrashObjectRef{ObjectID: "vm_KTKc6pwFxLjJ40QY"},
			},
			want: &Task{
				ID:     "task_Fq0vMXkSkKkGU3ut",
				Name:   "Purge items from trash",
				Status: "pending",
			},
			wantQuery: &url.Values{
				"trash_object[object_id]": []string{"vm_KTKc6pwFxLjJ40QY"},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("trash_object_purge"),
		},
		{
			name: "non-existent trash object",
			args: args{
				ctx: context.Background(),
				ref: TrashObjectRef{ID: "trsh_nopenotfound"},
			},
			errStr:     fixtureTrashObjectNotFoundErr,
			errResp:    fixtureTrashObjectNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("trash_object_not_found_error"),
		},
		{
			name: "permission denied",
			args: args{
				ctx: context.Background(),
				ref: TrashObjectRef{ID: "trsh_hkW1SMq0Bn8yNrRx"},
			},
			errStr:     fixturePermissionDeniedErr,
			errResp:    fixturePermissionDeniedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("permission_denied_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				ref: TrashObjectRef{ID: "trsh_hkW1SMq0Bn8yNrRx"},
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewTrashObjectsClient(rm)

			mux.HandleFunc(
				"/core/v1/trash_objects/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "DELETE", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					if tt.wantQuery != nil {
						assert.Equal(t, *tt.wantQuery, r.URL.Query())
					} else {
						assert.Equal(t,
							*tt.args.ref.queryValues(), r.URL.Query(),
						)
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.Purge(tt.args.ctx, tt.args.ref)

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

func TestTrashObjectsClient_PurgeAll(t *testing.T) {
	type args struct {
		ctx context.Context
		org OrganizationRef
	}
	tests := []struct {
		name       string
		args       args
		want       *Task
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "by organization ID",
			args: args{
				ctx: context.Background(),
				org: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
			},
			want: &Task{
				ID:     "task_lwZ65NKwJVB9a4E8",
				Name:   "Purge items from trash",
				Status: "pending",
			},
			respStatus: http.StatusOK,
			respBody:   fixture("trash_objects_purge_all"),
		},
		{
			name: "by organization SubDomain",
			args: args{
				ctx: context.Background(),
				org: OrganizationRef{SubDomain: "acme"},
			},
			want: &Task{
				ID:     "task_lwZ65NKwJVB9a4E8",
				Name:   "Purge items from trash",
				Status: "pending",
			},
			respStatus: http.StatusOK,
			respBody:   fixture("trash_objects_purge_all"),
		},
		{
			name: "non-existent trash object",
			args: args{
				ctx: context.Background(),
				org: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
			},
			errStr:     fixtureOrganizationNotFoundErr,
			errResp:    fixtureOrganizationNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("organization_not_found_error"),
		},
		{
			name: "permission denied",
			args: args{
				ctx: context.Background(),
				org: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
			},
			errStr:     fixturePermissionDeniedErr,
			errResp:    fixturePermissionDeniedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("permission_denied_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				org: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewTrashObjectsClient(rm)

			mux.HandleFunc(
				"/core/v1/organizations/_/trash_objects/purge_all",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "POST", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					assert.Equal(t,
						*tt.args.org.queryValues(), r.URL.Query(),
					)

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.PurgeAll(tt.args.ctx, tt.args.org)

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

func TestTrashObjectsClient_Restore(t *testing.T) {
	type args struct {
		ctx context.Context
		ref TrashObjectRef
	}
	tests := []struct {
		name       string
		args       args
		want       *TrashObject
		wantQuery  *url.Values
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "by ID",
			args: args{
				ctx: context.Background(),
				ref: TrashObjectRef{ID: "trsh_hkW1SMq0Bn8yNrRx"},
			},
			want: &TrashObject{
				ID:         "trsh_hkW1SMq0Bn8yNrRx",
				KeepUntil:  timestampPtr(1610039056),
				ObjectID:   "vm_KTKc6pwFxLjJ40QY",
				ObjectType: "VirtualMachine",
			},
			wantQuery: &url.Values{
				"trash_object[id]": []string{"trsh_hkW1SMq0Bn8yNrRx"},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("trash_object_restore"),
		},
		{
			name: "by ObjectID",
			args: args{
				ctx: context.Background(),
				ref: TrashObjectRef{ObjectID: "vm_KTKc6pwFxLjJ40QY"},
			},
			want: &TrashObject{
				ID:         "trsh_hkW1SMq0Bn8yNrRx",
				KeepUntil:  timestampPtr(1610039056),
				ObjectID:   "vm_KTKc6pwFxLjJ40QY",
				ObjectType: "VirtualMachine",
			},
			wantQuery: &url.Values{
				"trash_object[object_id]": []string{"vm_KTKc6pwFxLjJ40QY"},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("trash_object_restore"),
		},
		{
			name: "non-existent trash object",
			args: args{
				ctx: context.Background(),
				ref: TrashObjectRef{ID: "trsh_nopenotfound"},
			},
			errStr:     fixtureTrashObjectNotFoundErr,
			errResp:    fixtureTrashObjectNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("trash_object_not_found_error"),
		},
		{
			name: "permission denied",
			args: args{
				ctx: context.Background(),
				ref: TrashObjectRef{ID: "trsh_hkW1SMq0Bn8yNrRx"},
			},
			errStr:     fixturePermissionDeniedErr,
			errResp:    fixturePermissionDeniedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("permission_denied_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				ref: TrashObjectRef{ID: "trsh_hkW1SMq0Bn8yNrRx"},
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewTrashObjectsClient(rm)

			mux.HandleFunc(
				"/core/v1/trash_objects/_/restore",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "POST", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					if tt.wantQuery != nil {
						assert.Equal(t, *tt.wantQuery, r.URL.Query())
					} else {
						assert.Equal(t,
							*tt.args.ref.queryValues(), r.URL.Query(),
						)
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.Restore(tt.args.ctx, tt.args.ref)

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
