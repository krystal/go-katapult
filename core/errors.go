package core

import (
	"fmt"

	"github.com/krystal/go-katapult"
)

//go:generate go run github.com/krystal/go-katapult/tools/codegen -t errors -p core -o . -i CoreAPI\/.* -e .+Legacy.+ -f ../schemas/core/v1.json

var Err = fmt.Errorf("%w: core", katapult.Err)

func handleResponseError(err error) error {
	if err == nil {
		return nil
	}

	if r, ok := err.(*katapult.ResponseError); ok {
		return castResponseError(r)
	}

	return err
}
