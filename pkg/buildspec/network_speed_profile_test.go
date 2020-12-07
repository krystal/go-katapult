package buildspec

import (
	"encoding/xml"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetworkSpeedProfile_Marshaling(t *testing.T) {
	tests := []struct {
		name    string
		obj     *NetworkSpeedProfile
		decoded *NetworkSpeedProfile
	}{
		{
			name: "empty",
			obj:  &NetworkSpeedProfile{},
		},
		{
			name: "by ID",
			obj:  &NetworkSpeedProfile{ID: "nsp_bFQhDNAluyp4t2A9"},
		},
		{
			name: "by Permalink",
			obj:  &NetworkSpeedProfile{Permalink: "1g"},
		},
		{
			name: "with ID and Permalink",
			obj: &NetworkSpeedProfile{
				ID:        "nsp_bFQhDNAluyp4t2A9",
				Permalink: "1g",
			},
			decoded: &NetworkSpeedProfile{ID: "nsp_bFQhDNAluyp4t2A9"},
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

func TestNetworkSpeedProfile_UnmarshalXML_InvalidByAttr(t *testing.T) {
	err := xml.Unmarshal(
		[]byte(`<NetworkSpeedProfile by="other">foo</NetworkSpeedProfile>`),
		&NetworkSpeedProfile{},
	)

	assert.EqualError(t, err,
		`parse_xml: NetworkSpeedProfile by="other" is not supported`,
	)
	assert.True(t, errors.Is(err, ErrParseXML))
}
