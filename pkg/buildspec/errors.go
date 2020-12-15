package buildspec

import (
	"errors"
	"fmt"
)

var (
	Err = errors.New("")

	ErrParse    = fmt.Errorf("%wparse", Err)
	ErrParseXML = fmt.Errorf("%w_xml", ErrParse)
)
