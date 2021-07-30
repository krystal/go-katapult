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
