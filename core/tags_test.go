package core

import (
	"context"
	"fmt"
	"github.com/krystal/go-katapult"
	"github.com/krystal/go-katapult/internal/test"
	"github.com/krystal/go-katapult/internal/testclient"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestTag_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *Tag
	}{
		{
			name: "empty",
			obj:  &Tag{},
		},
		{
			name: "full",
			obj: &Tag{
				ID:        "id1",
				Name:      "name",
				Color:     "color",
				CreatedAt: timestampPtr(3043009),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestNewTagsClient(t *testing.T) {
	tc := testclient.New(nil, nil, nil)
	c := NewTagsClient(tc)
	assert.Equal(t, tc, c.client)
	assert.Equal(t, &url.URL{Path: "/core/v1/"}, c.basePath)
}

func Test_tagsResponseBody_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *tagsResponseBody
	}{
		{
			name: "empty",
			obj:  &tagsResponseBody{},
		},
		{
			name: "tag",
			obj: &tagsResponseBody{
				Tag: &Tag{
					ID:        "id1",
					Name:      "name",
					Color:     "color",
					CreatedAt: timestampPtr(3043009),
				},
			},
		},
		{
			name: "tags",
			obj: &tagsResponseBody{
				Tags: []*Tag{
					{
						ID:        "id1",
						Name:      "name",
						Color:     "color",
						CreatedAt: timestampPtr(3043009),
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

func TestTagRef_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *TagRef
	}{
		{
			name: "empty",
			obj:  &TagRef{},
		},
		{
			name: "full",
			obj: &TagRef{ID: "id1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestTagArguments_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *TagArguments
	}{
		{
			name: "empty",
			obj:  &TagArguments{},
		},
		{
			name: "name",
			obj: &TagArguments{Name: "tag_name"},
		},
		{
			name: "color",
			obj: &TagArguments{Color: "#0d1d1f"},
		},
		{
			name: "full",
			obj: &TagArguments{Name: "tag_name", Color: "#0d1d1f"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestTagsClient_List(t *testing.T) {
	type args struct {
		ctx  context.Context
		org  OrganizationRef
		opts *ListOptions
	}
	tests := []struct {
		name    string
		args    args
		resp    *katapult.Response
		respErr error
		respV   *tagsResponseBody
		want    []*Tag
		wantReq *katapult.Request
		wantErr string
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				org: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				opts: &ListOptions{
					Page:    5,
					PerPage: 32,
				},
			},
			resp: &katapult.Response{
				Pagination: &katapult.Pagination{Total: 333},
			},
			respV: &tagsResponseBody{
				Tags: []*Tag{
					{ID: "tag_O574YEEEYeLmqdmn"},
				},
			},
			want: []*Tag{
				{ID: "tag_O574YEEEYeLmqdmn"},
			},
			wantReq: &katapult.Request{
				Method: "GET",
				URL: &url.URL{
					Path: "/core/v1/organizations/_/tags",
					RawQuery: url.Values{
						"page":     []string{"5"},
						"per_page": []string{"32"},
						"organization[id]": []string{
							"org_O648YDMEYeLmqdmn",
						},
					}.Encode(),
				},
			},
		},
		{
			name: "success with nil options",
			args: args{
				ctx: context.Background(),
				org: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
			},
			resp: &katapult.Response{
				Pagination: &katapult.Pagination{Total: 333},
			},
			respV: &tagsResponseBody{
				Tags: []*Tag{
					{ID: "tag_O574YEEEYeLmqdmn"},
				},
			},
			want: []*Tag{
				{ID: "tag_O574YEEEYeLmqdmn"},
			},
			wantReq: &katapult.Request{
				Method: "GET",
				URL: &url.URL{
					Path: "/core/v1/organizations/_/tags",
					RawQuery: url.Values{
						"organization[id]": []string{
							"org_O648YDMEYeLmqdmn",
						},
					}.Encode(),
				},
			},
		},
		{
			name: "request error",
			args: args{
				ctx: context.Background(),
				org: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
			},
			respErr: fmt.Errorf("flux capacitor undercharged"),
			wantErr: "flux capacitor undercharged",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := testclient.New(tt.resp, tt.respErr, tt.respV)
			c := NewTagsClient(tc)
			ctx := test.Context(tt.args.ctx)

			got, resp, err := c.List(ctx, tt.args.org, tt.args.opts)

			assert.Equal(t, 1, len(tc.Calls), "only 1 request should be made")
			test.AssertContext(t, ctx, tc.Ctx)

			assert.Equal(t, tt.want, got)

			if tt.resp != nil {
				assert.Equal(t, tt.resp, resp)
			}

			if tt.wantReq != nil {
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

func TestTagsClient_Get(t *testing.T) {
	type args struct {
		ctx  context.Context
		ref  TagRef
	}
	tests := []struct {
		name    string
		args    args
		resp    *katapult.Response
		respErr error
		respV   *tagsResponseBody
		want    *Tag
		wantReq *katapult.Request
		wantErr string
	}{
		{
			name: "success",
			args: args{
				ctx:  context.Background(),
				ref:  TagRef{ID: "tag_O574YEEEYeLmqdmn"},
			},
			resp: &katapult.Response{
				Pagination: &katapult.Pagination{Total: 333},
			},
			respV: &tagsResponseBody{
				Tag: &Tag{
					ID: "tag_O574YEEEYeLmqdmn",
				},
			},
			want: &Tag{
				ID: "tag_O574YEEEYeLmqdmn",
			},
			wantReq: &katapult.Request{
				Method: "GET",
				URL: &url.URL{
					Path: "/core/v1/tags/_",
					RawQuery: url.Values{
						"tag[id]": []string{
							"tag_O574YEEEYeLmqdmn",
						},
					}.Encode(),
				},
			},
		},
		{
			name: "request error",
			args: args{
				ctx: context.Background(),
				ref: TagRef{ID: "tag_O574YEEEYeLmqdmn"},
			},
			respErr: fmt.Errorf("flux capacitor undercharged"),
			wantErr: "flux capacitor undercharged",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := testclient.New(tt.resp, tt.respErr, tt.respV)
			c := NewTagsClient(tc)
			ctx := test.Context(tt.args.ctx)

			got, resp, err := c.Get(ctx, tt.args.ref)

			assert.Equal(t, 1, len(tc.Calls), "only 1 request should be made")
			test.AssertContext(t, ctx, tc.Ctx)

			assert.Equal(t, tt.want, got)

			if tt.resp != nil {
				assert.Equal(t, tt.resp, resp)
			}

			if tt.wantReq != nil {
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

func TestTagsClient_Create(t *testing.T) {
	type args struct {
		ctx  context.Context
		org   OrganizationRef
		args  TagArguments
	}
	tests := []struct {
		name    string
		args    args
		resp    *katapult.Response
		respErr error
		respV   *tagsResponseBody
		want    *Tag
		wantReq *katapult.Request
		wantErr string
	}{
		{
			name: "success",
			args: args{
				ctx:  context.Background(),
				org:  OrganizationRef{ID: "org_O574YEEEYeLmqdmn"},
				args: TagArguments{
					Name:  "testing",
					Color: "#2ACAEA",
				},
			},
			resp: &katapult.Response{
				Pagination: &katapult.Pagination{Total: 333},
			},
			respV: &tagsResponseBody{
				Tag: &Tag{
					ID: "tag_O574YEEEYeLmqdmn",
				},
			},
			want: &Tag{
				ID: "tag_O574YEEEYeLmqdmn",
			},
			wantReq: &katapult.Request{
				Method: "POST",
				URL: &url.URL{
					Path: "/core/v1/organizations/_/tags",
					RawQuery: url.Values{
						"organization[id]": []string{
							"org_O574YEEEYeLmqdmn",
						},
					}.Encode(),
				},
				Body: TagArguments{
					Name:  "testing",
					Color: "#2ACAEA",
				},
			},
		},
		{
			name: "request error",
			args: args{
				ctx:  context.Background(),
				org:  OrganizationRef{ID: "org_O574YEEEYeLmqdmn"},
				args: TagArguments{
					Name:  "testing",
					Color: "#2ACAEA",
				},
			},
			respErr: fmt.Errorf("flux capacitor undercharged"),
			wantErr: "flux capacitor undercharged",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := testclient.New(tt.resp, tt.respErr, tt.respV)
			c := NewTagsClient(tc)
			ctx := test.Context(tt.args.ctx)

			got, resp, err := c.Create(ctx, tt.args.org, tt.args.args)

			assert.Equal(t, 1, len(tc.Calls), "only 1 request should be made")
			test.AssertContext(t, ctx, tc.Ctx)

			assert.Equal(t, tt.want, got)

			if tt.resp != nil {
				assert.Equal(t, tt.resp, resp)
			}

			if tt.wantReq != nil {
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

func TestTagsClient_Update(t *testing.T) {
	type args struct {
		ctx  context.Context
		ref  TagRef
		args TagArguments
	}
	tests := []struct {
		name    string
		args    args
		resp    *katapult.Response
		respErr error
		respV   *tagsResponseBody
		want    *Tag
		wantReq *katapult.Request
		wantErr string
	}{
		{
			name: "success",
			args: args{
				ctx:  context.Background(),
				ref:  TagRef{ID: "tag_O574YEEEYeLmqdmn"},
				args: TagArguments{
					Name:  "testing",
					Color: "#2ACAEA",
				},
			},
			resp: &katapult.Response{
				Pagination: &katapult.Pagination{Total: 333},
			},
			respV: &tagsResponseBody{
				Tag: &Tag{
					ID: "tag_O574YEEEYeLmqdmn",
				},
			},
			want: &Tag{
				ID: "tag_O574YEEEYeLmqdmn",
			},
			wantReq: &katapult.Request{
				Method: "PATCH",
				URL: &url.URL{
					Path: "/core/v1/tags/_",
					RawQuery: url.Values{
						"tag[id]": []string{
							"tag_O574YEEEYeLmqdmn",
						},
					}.Encode(),
				},
				Body: TagArguments{
					Name:  "testing",
					Color: "#2ACAEA",
				},
			},
		},
		{
			name: "request error",
			args: args{
				ctx:  context.Background(),
				ref:  TagRef{ID: "tag_O574YEEEYeLmqdmn"},
			},
			respErr: fmt.Errorf("flux capacitor undercharged"),
			wantErr: "flux capacitor undercharged",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := testclient.New(tt.resp, tt.respErr, tt.respV)
			c := NewTagsClient(tc)
			ctx := test.Context(tt.args.ctx)

			got, resp, err := c.Update(ctx, tt.args.ref, tt.args.args)

			assert.Equal(t, 1, len(tc.Calls), "only 1 request should be made")
			test.AssertContext(t, ctx, tc.Ctx)

			assert.Equal(t, tt.want, got)

			if tt.resp != nil {
				assert.Equal(t, tt.resp, resp)
			}

			if tt.wantReq != nil {
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

func TestTagsClient_Delete(t *testing.T) {
	type args struct {
		ctx  context.Context
		ref  TagRef
	}
	tests := []struct {
		name    string
		args    args
		resp    *katapult.Response
		respErr error
		respV   *tagsResponseBody
		want    *Tag
		wantReq *katapult.Request
		wantErr string
	}{
		{
			name: "success",
			args: args{
				ctx:  context.Background(),
				ref:  TagRef{ID: "tag_O574YEEEYeLmqdmn"},
			},
			resp: &katapult.Response{
				Pagination: &katapult.Pagination{Total: 333},
			},
			respV: &tagsResponseBody{
				Tag: &Tag{
					ID: "tag_O574YEEEYeLmqdmn",
				},
			},
			want: &Tag{
				ID: "tag_O574YEEEYeLmqdmn",
			},
			wantReq: &katapult.Request{
				Method: "DELETE",
				URL: &url.URL{
					Path: "/core/v1/tags/_",
					RawQuery: url.Values{
						"tag[id]": []string{
							"tag_O574YEEEYeLmqdmn",
						},
					}.Encode(),
				},
			},
		},
		{
			name: "request error",
			args: args{
				ctx:  context.Background(),
				ref:  TagRef{ID: "tag_O574YEEEYeLmqdmn"},
			},
			respErr: fmt.Errorf("flux capacitor undercharged"),
			wantErr: "flux capacitor undercharged",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := testclient.New(tt.resp, tt.respErr, tt.respV)
			c := NewTagsClient(tc)
			ctx := test.Context(tt.args.ctx)

			got, resp, err := c.Delete(ctx, tt.args.ref)

			assert.Equal(t, 1, len(tc.Calls), "only 1 request should be made")
			test.AssertContext(t, ctx, tc.Ctx)

			assert.Equal(t, tt.want, got)

			if tt.resp != nil {
				assert.Equal(t, tt.resp, resp)
			}

			if tt.wantReq != nil {
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
