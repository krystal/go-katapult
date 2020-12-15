package buildspec

import (
	"encoding/xml"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_Marshaling(t *testing.T) {
	tests := []struct {
		name    string
		obj     *User
		decoded *User
	}{
		{
			name: "empty",
			obj:  &User{},
		},
		{
			name: "by ID",
			obj:  &User{ID: "user_GdKQbGfwlG6iF7a1"},
		},
		{
			name: "by EmailAddress",
			obj:  &User{EmailAddress: "john@doe.com"},
		},
		{
			name: "with ID and EmailAddress",
			obj: &User{
				ID:           "user_GdKQbGfwlG6iF7a1",
				EmailAddress: "john@doe.com",
			},
			decoded: &User{ID: "user_GdKQbGfwlG6iF7a1"},
		},
	}
	for _, tt := range tests {
		t.Run("json_"+tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
		t.Run("xml_"+tt.name, func(t *testing.T) {
			testCustomXMLMarshaling(t, tt.obj, tt.decoded)
		})
		t.Run("yaml_"+tt.name, func(t *testing.T) {
			testYAMLMarshaling(t, tt.obj)
		})
	}
}

func TestUser_UnmarshalXML_InvalidByAttr(t *testing.T) {
	err := xml.Unmarshal(
		[]byte(`<User by="other">foo</User>`),
		&User{},
	)

	assert.EqualError(t, err,
		`parse_xml: User by="other" is not supported`,
	)
	assert.True(t, errors.Is(err, ErrParseXML))
}

func Test_xmlUsers_Marshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *xmlUsers
	}{
		{
			name: "empty",
			obj:  &xmlUsers{},
		},
		{
			name: "full",
			obj: &xmlUsers{
				Users: []*User{
					{ID: "user_yUfYcKHgU1ywBWzP"},
					{EmailAddress: "jane@doe.com"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run("xml_"+tt.name, func(t *testing.T) {
			testXMLMarshaling(t, tt.obj)
		})
	}
}
