package core

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/krystal/go-katapult"
	"github.com/krystal/go-katapult/internal/test"
	"github.com/krystal/go-katapult/internal/testclient"
	"github.com/stretchr/testify/assert"
)

var (
	fixtureFileStorageVolumeNotFoundErr = "katapult: not_found: " +
		"file_storage_volume_not_found: No file storage volume was found " +
		"matching any of the criteria provided in the arguments."
	fixtureFileStorageVolumeNotFoundResponseError = &katapult.ResponseError{
		Code: "file_storage_volume_not_found",
		Description: "No file storage volume was found matching any of the " +
			"criteria provided in the arguments.",
		Detail: json.RawMessage(`{}`),
	}
)

func TestNewFileStorageVolumesClient(t *testing.T) {
	tc := testclient.New(nil, nil, nil)
	c := NewFileStorageVolumesClient(tc)
	assert.Equal(t, tc, c.client)
	assert.Equal(t, &url.URL{Path: "/core/v1/"}, c.basePath)
}

func TestClient_FileStorageVolumes(t *testing.T) {
	c := New(&testclient.Client{})

	assert.IsType(t, &FileStorageVolumesClient{}, c.FileStorageVolumes)
}

func TestFileStorageVolume_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *FileStorageVolume
	}{
		{
			name: "empty",
			obj:  &FileStorageVolume{},
		},
		{
			name: "full",
			obj: &FileStorageVolume{
				ID:   "fsv_1",
				Name: "test volume",
				DataCenter: &DataCenter{
					ID: "dc_1",
				},
				Associations: []string{"assoc1", "assoc1"},
				State:        FileStorageVolumeReady,
				NFSLocation:  "foo:/mnt/nfs/volume1",
				Size:         1024,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestFileStorageVolume_Ref(t *testing.T) {
	fsv := FileStorageVolume{ID: "fsv_Q0RIr5iyQDUWF0qy"}

	assert.Equal(t, FileStorageVolumeRef{ID: "fsv_Q0RIr5iyQDUWF0qy"}, fsv.Ref())
}

func TestFileStorageVolumeRef_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *FileStorageVolumeRef
	}{
		{
			name: "empty",
			obj:  &FileStorageVolumeRef{},
		},
		{
			name: "full",
			obj:  &FileStorageVolumeRef{ID: "id1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestFileStorageVolumeRef_queryValues(t *testing.T) {
	tests := []struct {
		name string
		obj  FileStorageVolumeRef
	}{
		{
			name: "with id",
			obj:  FileStorageVolumeRef{ID: "id1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testQueryableEncoding(t, tt.obj)
		})
	}
}

func TestFileStorageVolumeStates(t *testing.T) {
	tests := []struct {
		name  string
		enum  FileStorageVolumeState
		value string
	}{
		{
			name:  "FileStorageVolumePending",
			enum:  FileStorageVolumePending,
			value: "pending",
		},
		{
			name:  "FileStorageVolumeFailed",
			enum:  FileStorageVolumeFailed,
			value: "failed",
		},
		{
			name:  "FileStorageVolumeReady",
			enum:  FileStorageVolumeReady,
			value: "ready",
		},
		{
			name:  "FileStorageVolumeConfiguring",
			enum:  FileStorageVolumeConfiguring,
			value: "configuring",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.value, string(tt.enum))
		})
	}
}

func TestFileStorageVolumesClient_List(t *testing.T) {
	type args struct {
		ctx  context.Context
		org  OrganizationRef
		opts *ListOptions
	}
	tests := []struct {
		name      string
		args      args
		resp      *katapult.Response
		respErr   error
		respV     *fileStorageVolumesResponseBody
		wantReq   *katapult.Request
		want      []*FileStorageVolume
		wantResp  *katapult.Response
		wantErr   string
		wantErrIs []error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				org: OrganizationRef{ID: "org_bvvDKVsWRQVWIAN2"},
				opts: &ListOptions{
					Page:    5,
					PerPage: 32,
				},
			},
			resp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusOK},
			},
			respV: &fileStorageVolumesResponseBody{
				Pagination: &katapult.Pagination{Total: 394},
				FileStorageVolumes: []*FileStorageVolume{
					{ID: "fsv_DWLeE1AReyUEgC9K", Name: "test vol"},
				},
			},
			wantReq: &katapult.Request{
				Method: "GET",
				URL: &url.URL{
					Path: "/core/v1/organizations/_/file_storage_volumes",
					RawQuery: url.Values{
						"page":     []string{"5"},
						"per_page": []string{"32"},
						"organization[id]": []string{
							"org_bvvDKVsWRQVWIAN2",
						},
					}.Encode(),
				},
			},
			want: []*FileStorageVolume{{
				ID:   "fsv_DWLeE1AReyUEgC9K",
				Name: "test vol",
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
				org:  OrganizationRef{ID: "org_ERQqd1R6O63q4YEA"},
				opts: nil,
			},
			resp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusOK},
			},
			respV: &fileStorageVolumesResponseBody{
				Pagination: &katapult.Pagination{Total: 333},
				FileStorageVolumes: []*FileStorageVolume{
					{ID: "fsv_HiuoJWTHEsTCUMRT", Name: "test vol2"},
				},
			},
			wantReq: &katapult.Request{
				Method: "GET",
				URL: &url.URL{
					Path: "/core/v1/organizations/_/file_storage_volumes",
					RawQuery: url.Values{
						"organization[id]": []string{
							"org_ERQqd1R6O63q4YEA",
						},
					}.Encode(),
				},
			},
			want: []*FileStorageVolume{{
				ID:   "fsv_HiuoJWTHEsTCUMRT",
				Name: "test vol2",
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
				org:  OrganizationRef{ID: "org_eWB9a0ZjckFfoLcb"},
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
				org:  OrganizationRef{ID: "org_hDDMkjQoeCACgaev"},
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
			c := NewFileStorageVolumesClient(tc)
			ctx := test.Context(tt.args.ctx)

			got, resp, err := c.List(
				ctx,
				tt.args.org,
				tt.args.opts,
				testRequestOption,
			)

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

			for _, wantErrIs := range tt.wantErrIs {
				assert.ErrorIs(t, err, wantErrIs)
			}
		})
	}
}

func TestFileStorageVolumesClient_Get(t *testing.T) {
	type args struct {
		ctx context.Context
		ref FileStorageVolumeRef
	}
	tests := []struct {
		name      string
		args      args
		resp      *katapult.Response
		respErr   error
		respV     *fileStorageVolumesResponseBody
		wantReq   *katapult.Request
		want      *FileStorageVolume
		wantResp  *katapult.Response
		wantErr   string
		wantErrIs []error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				ref: FileStorageVolumeRef{ID: "fsv_BHxh7IuGwP4OSqzP"},
			},
			resp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusOK},
			},
			respV: &fileStorageVolumesResponseBody{
				FileStorageVolume: &FileStorageVolume{
					ID:   "fsv_BHxh7IuGwP4OSqzP",
					Name: "test vol-42",
				},
			},
			wantReq: &katapult.Request{
				Method: "GET",
				URL: &url.URL{
					Path: "/core/v1/file_storage_volumes/_",
					RawQuery: url.Values{
						"file_storage_volume[id]": []string{
							"fsv_BHxh7IuGwP4OSqzP",
						},
					}.Encode(),
				},
			},
			want: &FileStorageVolume{
				ID:   "fsv_BHxh7IuGwP4OSqzP",
				Name: "test vol-42",
			},
			wantResp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusOK},
			},
		},
		{
			name: "not found",
			args: args{
				ctx: context.Background(),
				ref: FileStorageVolumeRef{ID: "fsv_SljuxNUHyFcd28lp"},
			},
			resp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusNotFound},
			},
			respErr: fixtureFileStorageVolumeNotFoundResponseError,
			wantResp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusNotFound},
			},
			wantErr: fixtureFileStorageVolumeNotFoundErr,
			wantErrIs: []error{
				ErrFileStorageVolumeNotFound,
				katapult.ErrNotFound,
				katapult.Err,
			},
		},
		{
			name: "request error",
			args: args{
				ctx: context.Background(),
				ref: FileStorageVolumeRef{ID: "fsv_BHxh7IuGwP4OSqzP"},
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
				ref: FileStorageVolumeRef{ID: "fsv_BHxh7IuGwP4OSqzP"},
			},
			resp:    nil,
			respErr: fmt.Errorf("someting is really wrong"),
			wantErr: "someting is really wrong",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := testclient.New(tt.resp, tt.respErr, tt.respV)
			c := NewFileStorageVolumesClient(tc)
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

			for _, wantErrIs := range tt.wantErrIs {
				assert.ErrorIs(t, err, wantErrIs)
			}
		})
	}
}

func TestFileStorageVolumesClient_GetByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name      string
		args      args
		resp      *katapult.Response
		respErr   error
		respV     *fileStorageVolumesResponseBody
		wantReq   *katapult.Request
		want      *FileStorageVolume
		wantResp  *katapult.Response
		wantErr   string
		wantErrIs []error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				id:  "fsv_SljuxNUHyFcd28lp",
			},
			resp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusOK},
			},
			respV: &fileStorageVolumesResponseBody{
				FileStorageVolume: &FileStorageVolume{
					ID:   "fsv_SljuxNUHyFcd28lp",
					Name: "test vol-43",
				},
			},
			wantReq: &katapult.Request{
				Method: "GET",
				URL: &url.URL{
					Path: "/core/v1/file_storage_volumes/_",
					RawQuery: url.Values{
						"file_storage_volume[id]": []string{
							"fsv_SljuxNUHyFcd28lp",
						},
					}.Encode(),
				},
			},
			want: &FileStorageVolume{
				ID:   "fsv_SljuxNUHyFcd28lp",
				Name: "test vol-43",
			},
			wantResp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusOK},
			},
		},
		{
			name: "not found",
			args: args{
				ctx: context.Background(),
				id:  "fsv_SljuxNUHyFcd28lp",
			},
			resp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusNotFound},
			},
			respErr: fixtureFileStorageVolumeNotFoundResponseError,
			wantResp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusNotFound},
			},
			wantErr: fixtureFileStorageVolumeNotFoundErr,
			wantErrIs: []error{
				ErrFileStorageVolumeNotFound,
				katapult.ErrNotFound,
				katapult.Err,
			},
		},
		{
			name: "request error",
			args: args{
				ctx: context.Background(),
				id:  "fsv_SljuxNUHyFcd28lp",
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
				id:  "fsv_SljuxNUHyFcd28lp",
			},
			resp:    nil,
			respErr: fmt.Errorf("someting is really wrong"),
			wantErr: "someting is really wrong",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := testclient.New(tt.resp, tt.respErr, tt.respV)
			c := NewFileStorageVolumesClient(tc)
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

			for _, wantErrIs := range tt.wantErrIs {
				assert.ErrorIs(t, err, wantErrIs)
			}
		})
	}
}

func TestFileStorageVolumeCreateArguments_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *FileStorageVolumeCreateArguments
	}{
		{
			name: "empty",
			obj:  &FileStorageVolumeCreateArguments{},
		},
		{
			name: "name",
			obj:  &FileStorageVolumeCreateArguments{Name: "volume_name"},
		},
		{
			name: "data_center",
			obj: &FileStorageVolumeCreateArguments{
				DataCenter: DataCenterRef{ID: "dc1"},
			},
		},
		{
			name: "associations",
			obj: &FileStorageVolumeCreateArguments{
				Associations: []string{"assoc1", "assoc2"},
			},
		},
		{
			name: "full",
			obj: &FileStorageVolumeCreateArguments{
				Name:         "volume_name",
				DataCenter:   DataCenterRef{ID: "dc1"},
				Associations: []string{"assoc1", "assoc2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func Test_fileStorageVolumeCreateRequest_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *fileStorageVolumeCreateRequest
	}{
		{
			name: "empty",
			obj:  &fileStorageVolumeCreateRequest{},
		},
		{
			name: "full",
			obj: &fileStorageVolumeCreateRequest{
				Organization: OrganizationRef{ID: "org1"},
				Properties: &FileStorageVolumeCreateArguments{
					Name:         "created",
					DataCenter:   DataCenterRef{ID: "dc1"},
					Associations: []string{"assoc1", "assoc2"},
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

func TestFileStorageVolumesClient_Create(t *testing.T) {
	type args struct {
		ctx  context.Context
		ref  OrganizationRef
		args *FileStorageVolumeCreateArguments
	}
	tests := []struct {
		name      string
		args      args
		resp      *katapult.Response
		respErr   error
		respV     *fileStorageVolumesResponseBody
		wantReq   *katapult.Request
		want      *FileStorageVolume
		wantResp  *katapult.Response
		wantErr   string
		wantErrIs []error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				ref: OrganizationRef{ID: "org_rTya1qPVE3WlT3yf"},
				args: &FileStorageVolumeCreateArguments{
					Name: "cache data",
					DataCenter: DataCenterRef{
						ID: "dc_W5WpI8fRyV4KGx77",
					},
					Associations: []string{"vm_vkbs4CqMOLPja4OK"},
				},
			},
			resp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusOK},
			},
			respV: &fileStorageVolumesResponseBody{
				FileStorageVolume: &FileStorageVolume{
					ID:   "lbrule_55P1GfFvW5pPPhgh",
					Name: "cache data",
					DataCenter: &DataCenter{
						ID:        "dc_W5WpI8fRyV4KGx77",
						Name:      "London",
						Permalink: "uk-lon-01",
					},
					Associations: []string{"vm_vkbs4CqMOLPja4OK"},
					State:        "pending",
					NFSLocation:  "",
					Size:         0,
				},
			},
			wantReq: &katapult.Request{
				Method: "POST",
				URL: &url.URL{
					Path: "/core/v1/organizations/_/file_storage_volumes",
				},
				Body: &fileStorageVolumeCreateRequest{
					Organization: OrganizationRef{ID: "org_rTya1qPVE3WlT3yf"},
					Properties: &FileStorageVolumeCreateArguments{
						Name: "cache data",
						DataCenter: DataCenterRef{
							ID: "dc_W5WpI8fRyV4KGx77",
						},
						Associations: []string{"vm_vkbs4CqMOLPja4OK"},
					},
				},
			},
			want: &FileStorageVolume{
				ID:   "lbrule_55P1GfFvW5pPPhgh",
				Name: "cache data",
				DataCenter: &DataCenter{
					ID:        "dc_W5WpI8fRyV4KGx77",
					Name:      "London",
					Permalink: "uk-lon-01",
				},
				Associations: []string{"vm_vkbs4CqMOLPja4OK"},
				State:        "pending",
				NFSLocation:  "",
				Size:         0,
			},
			wantResp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusOK},
			},
		},
		{
			name: "request error",
			args: args{
				ctx:  context.Background(),
				ref:  OrganizationRef{ID: "org_EZrs0jcaY8IJTqlB"},
				args: &FileStorageVolumeCreateArguments{},
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
				ref:  OrganizationRef{ID: "org_Pp2uLFDRPK1xuT9T"},
				args: &FileStorageVolumeCreateArguments{},
			},
			resp:    nil,
			respErr: fmt.Errorf("someting is really wrong"),
			wantErr: "someting is really wrong",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := testclient.New(tt.resp, tt.respErr, tt.respV)
			c := NewFileStorageVolumesClient(tc)
			ctx := test.Context(tt.args.ctx)

			got, resp, err := c.Create(
				ctx,
				tt.args.ref,
				tt.args.args,
				testRequestOption,
			)

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

			for _, wantErrIs := range tt.wantErrIs {
				assert.ErrorIs(t, err, wantErrIs)
			}
		})
	}
}

func TestFileStorageVolumeUpdateArguments_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *FileStorageVolumeUpdateArguments
	}{
		{
			name: "empty",
			obj:  &FileStorageVolumeUpdateArguments{},
		},
		{
			name: "name",
			obj:  &FileStorageVolumeUpdateArguments{Name: "volume_name"},
		},
		{
			name: "associations",
			obj: &FileStorageVolumeUpdateArguments{
				Associations: &[]string{"assoc1", "assoc2"},
			},
		},
		{
			name: "full",
			obj: &FileStorageVolumeUpdateArguments{
				Name:         "volume_name",
				Associations: &[]string{"assoc1", "assoc2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func Test_fileStorageVolumeUpdateRequest_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *fileStorageVolumeUpdateRequest
	}{
		{
			name: "empty",
			obj:  &fileStorageVolumeUpdateRequest{},
		},
		{
			name: "full",
			obj: &fileStorageVolumeUpdateRequest{
				FileStorageVolume: FileStorageVolumeRef{
					ID: "sg_3uXbmANw4sQiF1J3",
				},
				Properties: &FileStorageVolumeUpdateArguments{
					Name:         "updated",
					Associations: &[]string{"assoc1", "assoc2"},
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

func TestFileStorageVolumesClient_Update(t *testing.T) {
	type args struct {
		ctx  context.Context
		ref  FileStorageVolumeRef
		args *FileStorageVolumeUpdateArguments
	}
	tests := []struct {
		name      string
		args      args
		resp      *katapult.Response
		respErr   error
		respV     *fileStorageVolumesResponseBody
		wantReq   *katapult.Request
		want      *FileStorageVolume
		wantResp  *katapult.Response
		wantErr   string
		wantErrIs []error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				ref: FileStorageVolumeRef{ID: "fsv_JtyrhImi5jBjn5ig"},
				args: &FileStorageVolumeUpdateArguments{
					Name:         "updated volume name",
					Associations: &[]string{"vm_riYl1387Fdt2bcMA"},
				},
			},
			resp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusOK},
			},
			respV: &fileStorageVolumesResponseBody{
				FileStorageVolume: &FileStorageVolume{
					ID:           "fsv_JtyrhImi5jBjn5ig",
					Name:         "updated volume name",
					DataCenter:   &DataCenter{},
					Associations: []string{"vm_riYl1387Fdt2bcMA"},
					State:        "configuring",
					NFSLocation:  "nfs.store:/fsv_JtyrhImi5jBjn5ig",
					Size:         3490,
				},
			},
			wantReq: &katapult.Request{
				Method: "PATCH",
				URL: &url.URL{
					Path: "/core/v1/file_storage_volumes/_",
				},
				Body: &fileStorageVolumeUpdateRequest{
					FileStorageVolume: FileStorageVolumeRef{
						ID: "fsv_JtyrhImi5jBjn5ig",
					},
					Properties: &FileStorageVolumeUpdateArguments{
						Name:         "updated volume name",
						Associations: &[]string{"vm_riYl1387Fdt2bcMA"},
					},
				},
			},
			want: &FileStorageVolume{
				ID:           "fsv_JtyrhImi5jBjn5ig",
				Name:         "updated volume name",
				DataCenter:   &DataCenter{},
				Associations: []string{"vm_riYl1387Fdt2bcMA"},
				State:        "configuring",
				NFSLocation:  "nfs.store:/fsv_JtyrhImi5jBjn5ig",
				Size:         3490,
			},
			wantResp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusOK},
			},
		},
		{
			name: "not found",
			args: args{
				ctx:  context.Background(),
				ref:  FileStorageVolumeRef{ID: "fsv_SljuxNUHyFcd28lp"},
				args: &FileStorageVolumeUpdateArguments{Name: "changed again"},
			},
			resp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusNotFound},
			},
			respErr: fixtureFileStorageVolumeNotFoundResponseError,
			wantResp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusNotFound},
			},
			wantErr: fixtureFileStorageVolumeNotFoundErr,
			wantErrIs: []error{
				ErrFileStorageVolumeNotFound,
				katapult.ErrNotFound,
				katapult.Err,
			},
		},
		{
			name: "request error",
			args: args{
				ctx: context.Background(),
				ref: FileStorageVolumeRef{ID: "fsv_JtyrhImi5jBjn5ig"},
				args: &FileStorageVolumeUpdateArguments{
					Name:         "",
					Associations: &[]string{},
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
				ref: FileStorageVolumeRef{ID: "fsv_JtyrhImi5jBjn5ig"},
				args: &FileStorageVolumeUpdateArguments{
					Name:         "",
					Associations: &[]string{},
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
			c := NewFileStorageVolumesClient(tc)
			ctx := test.Context(tt.args.ctx)

			got, resp, err := c.Update(
				ctx,
				tt.args.ref,
				tt.args.args,
				testRequestOption,
			)

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

			for _, wantErrIs := range tt.wantErrIs {
				assert.ErrorIs(t, err, wantErrIs)
			}
		})
	}
}

func Test_fileStorageVolumesResponseBody_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *fileStorageVolumesResponseBody
	}{
		{
			name: "empty",
			obj:  &fileStorageVolumesResponseBody{},
		},
		{
			name: "full",
			obj: &fileStorageVolumesResponseBody{
				Pagination:         &katapult.Pagination{CurrentPage: 344},
				TrashObject:        &TrashObject{ID: "id1"},
				FileStorageVolume:  &FileStorageVolume{ID: "id2"},
				FileStorageVolumes: []*FileStorageVolume{{ID: "id3"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestFileStorageVolumesClient_Delete(t *testing.T) {
	type args struct {
		ctx context.Context
		ref FileStorageVolumeRef
	}
	tests := []struct {
		name      string
		args      args
		resp      *katapult.Response
		respErr   error
		respV     *fileStorageVolumesResponseBody
		wantReq   *katapult.Request
		want      *FileStorageVolume
		wantTrash *TrashObject
		wantResp  *katapult.Response
		wantErr   string
		wantErrIs []error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				ref: FileStorageVolumeRef{ID: "fsv_vSOkKO1NuPoDuZqR"},
			},
			resp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusOK},
			},
			respV: &fileStorageVolumesResponseBody{
				TrashObject: &TrashObject{
					ID:         "trsh_Bl2vmvd6kqvGfYkC",
					ObjectID:   "fsv_vSOkKO1NuPoDuZqR",
					ObjectType: "FileStorageVolume",
				},
				FileStorageVolume: &FileStorageVolume{
					ID:   "fsv_vSOkKO1NuPoDuZqR",
					Name: "My File Storage Volume",
				},
			},
			wantReq: &katapult.Request{
				Method: "DELETE",
				URL: &url.URL{
					Path: "/core/v1/file_storage_volumes/_",
					RawQuery: url.Values{
						"file_storage_volume[id]": []string{
							"fsv_vSOkKO1NuPoDuZqR",
						},
					}.Encode(),
				},
			},
			want: &FileStorageVolume{
				ID:   "fsv_vSOkKO1NuPoDuZqR",
				Name: "My File Storage Volume",
			},
			wantTrash: &TrashObject{
				ID:         "trsh_Bl2vmvd6kqvGfYkC",
				ObjectID:   "fsv_vSOkKO1NuPoDuZqR",
				ObjectType: "FileStorageVolume",
			},
			wantResp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusOK},
			},
		},
		{
			name: "not found",
			args: args{
				ctx: context.Background(),
				ref: FileStorageVolumeRef{ID: "fsv_SljuxNUHyFcd28lp"},
			},
			resp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusNotFound},
			},
			respErr: fixtureFileStorageVolumeNotFoundResponseError,
			wantResp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusNotFound},
			},
			wantErr: fixtureFileStorageVolumeNotFoundErr,
			wantErrIs: []error{
				ErrFileStorageVolumeNotFound,
				katapult.ErrNotFound,
				katapult.Err,
			},
		},
		{
			name: "request error",
			args: args{
				ctx: context.Background(),
				ref: FileStorageVolumeRef{ID: "fsv_vSOkKO1NuPoDuZqR"},
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
				ref: FileStorageVolumeRef{ID: "fsv_vSOkKO1NuPoDuZqR"},
			},
			resp:    nil,
			respErr: fmt.Errorf("someting is really wrong"),
			wantErr: "someting is really wrong",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := testclient.New(tt.resp, tt.respErr, tt.respV)
			c := NewFileStorageVolumesClient(tc)
			ctx := test.Context(tt.args.ctx)

			got, gotTrash, resp, err := c.Delete(
				ctx, tt.args.ref, testRequestOption,
			)

			assert.Equal(t, 1, len(tc.Calls), "only 1 request should be made")
			test.AssertContext(t, ctx, tc.Ctx)

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantTrash, gotTrash)

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

			for _, wantErrIs := range tt.wantErrIs {
				assert.ErrorIs(t, err, wantErrIs)
			}
		})
	}
}
