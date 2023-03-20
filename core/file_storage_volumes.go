package core

import (
	"context"
	"net/url"

	"github.com/krystal/go-katapult"
)

type FileStorageVolume struct {
	ID           string      `json:"id,omitempty"`
	Name         string      `json:"name,omitempty"`
	DataCenter   *DataCenter `json:"data_center,omitempty"`
	Associations []string    `json:"associations,omitempty"`
	State        string      `json:"state,omitempty"`
	NFSLocation  string      `json:"nfs_location,omitempty"`
	Size         int64       `json:"size,omitempty"`
}

func (fsv *FileStorageVolume) Ref() FileStorageVolumeRef {
	return FileStorageVolumeRef{ID: fsv.ID}
}

type FileStorageVolumeRef struct {
	ID string `json:"id"`
}

func (fsr FileStorageVolumeRef) queryValues() *url.Values {
	return &url.Values{"file_storage_volume[id]": []string{fsr.ID}}
}

type FileStorageVolumesClient struct {
	client   RequestMaker
	basePath *url.URL
}

func NewFileStorageVolumesClient(rm RequestMaker) *FileStorageVolumesClient {
	return &FileStorageVolumesClient{
		client:   rm,
		basePath: &url.URL{Path: "/core/v1/"},
	}
}

func (fsvc *FileStorageVolumesClient) List(
	ctx context.Context,
	org OrganizationRef,
	opts *ListOptions,
	reqOpts ...katapult.RequestOption,
) ([]*FileStorageVolume, *katapult.Response, error) {
	qs := queryValues(org, opts)
	u := &url.URL{
		Path:     "organizations/_/file_storage_volumes",
		RawQuery: qs.Encode(),
	}

	body, resp, err := fsvc.doRequest(ctx, "GET", u, nil, reqOpts...)
	resp.Pagination = body.Pagination

	return body.FileStorageVolumes, resp, err
}

func (fsvc *FileStorageVolumesClient) Get(
	ctx context.Context,
	ref FileStorageVolumeRef,
	reqOpts ...katapult.RequestOption,
) (*FileStorageVolume, *katapult.Response, error) {
	u := &url.URL{
		Path:     "file_storage_volumes/_",
		RawQuery: ref.queryValues().Encode(),
	}
	body, resp, err := fsvc.doRequest(ctx, "GET", u, nil, reqOpts...)

	return body.FileStorageVolume, resp, err
}

func (fsvc *FileStorageVolumesClient) GetByID(
	ctx context.Context,
	id string,
	reqOpts ...katapult.RequestOption,
) (*FileStorageVolume, *katapult.Response, error) {
	return fsvc.Get(ctx, FileStorageVolumeRef{ID: id}, reqOpts...)
}

type FileStorageVolumeCreateArguments struct {
	Name         string         `json:"name,omitempty"`
	DataCenter   *DataCenterRef `json:"data_center,omitempty"`
	Associations []string       `json:"associations,omitempty"`
}

type fileStorageVolumeCreateRequest struct {
	Organization OrganizationRef                   `json:"organization"`
	Properties   *FileStorageVolumeCreateArguments `json:"properties,omitempty"`
}

func (fsvc *FileStorageVolumesClient) Create(
	ctx context.Context,
	org OrganizationRef,
	args *FileStorageVolumeCreateArguments,
	reqOpts ...katapult.RequestOption,
) (*FileStorageVolume, *katapult.Response, error) {
	u := &url.URL{Path: "organizations/_/file_storage_volumes"}
	reqBody := &fileStorageVolumeCreateRequest{
		Organization: org,
		Properties:   args,
	}

	body, resp, err := fsvc.doRequest(ctx, "POST", u, reqBody, reqOpts...)

	return body.FileStorageVolume, resp, err
}

type FileStorageVolumeUpdateArguments struct {
	Name         string   `json:"name,omitempty"`
	Associations []string `json:"associations,omitempty"`
}

type fileStorageVolumeUpdateRequest struct {
	FileStorageVolume FileStorageVolumeRef              `json:"security_group"`
	Properties        *FileStorageVolumeUpdateArguments `json:"properties,omitempty"`
}

func (fsvc *FileStorageVolumesClient) Update(
	ctx context.Context,
	ref FileStorageVolumeRef,
	args *FileStorageVolumeUpdateArguments,
	reqOpts ...katapult.RequestOption,
) (*FileStorageVolume, *katapult.Response, error) {
	u := &url.URL{Path: "file_storage_volumes/_"}
	reqBody := &fileStorageVolumeUpdateRequest{
		FileStorageVolume: ref,
		Properties:        args,
	}

	body, resp, err := fsvc.doRequest(ctx, "PATCH", u, reqBody, reqOpts...)

	return body.FileStorageVolume, resp, err
}

func (fsvc *FileStorageVolumesClient) Delete(
	ctx context.Context,
	ref FileStorageVolumeRef,
	reqOpts ...katapult.RequestOption,
) (*FileStorageVolume, *katapult.Response, error) {
	u := &url.URL{
		Path:     "file_storage_volumes/_",
		RawQuery: ref.queryValues().Encode(),
	}
	body, resp, err := fsvc.doRequest(ctx, "DELETE", u, nil, reqOpts...)

	return body.FileStorageVolume, resp, err
}

type fileStorageVolumesResponseBody struct {
	Pagination         *katapult.Pagination `json:"pagination,omitempty"`
	FileStorageVolumes []*FileStorageVolume `json:"file_storage_volumes,omitempty"`
	FileStorageVolume  *FileStorageVolume   `json:"file_storage_volume,omitempty"`
}

func (fsvc *FileStorageVolumesClient) doRequest(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
	reqOpts ...katapult.RequestOption,
) (*fileStorageVolumesResponseBody, *katapult.Response, error) {
	u = fsvc.basePath.ResolveReference(u)
	respBody := &fileStorageVolumesResponseBody{}

	req := katapult.NewRequest(method, u, body, reqOpts...)
	resp, err := fsvc.client.Do(ctx, req, respBody)
	if resp == nil {
		resp = katapult.NewResponse(nil)
	}

	return respBody, resp, handleResponseError(err)
}
