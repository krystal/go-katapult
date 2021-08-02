package core

import (
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
