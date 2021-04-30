package buildspec

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErr(t *testing.T) {
	assert.Error(t, Err)
}

func TestErrParse(t *testing.T) {
	assert.EqualError(t, ErrParse, "parse")
	assert.True(t, errors.Is(ErrParse, Err), "ErrParse is not a Err")
}

func TestErrParseXML(t *testing.T) {
	assert.EqualError(t, ErrParseXML, "parse_xml")
	assert.True(t, errors.Is(ErrParseXML, ErrParse),
		"ErrParseXML is not a ErrParse",
	)
	assert.True(t, errors.Is(ErrParseXML, Err), "ErrParseXML is not a Err")
}
