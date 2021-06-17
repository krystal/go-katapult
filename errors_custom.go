package katapult

import (
	"fmt"
	"strings"
)

// This file contains custom Error() functions for select error types in
// errors_generated.go, allowing customized error messages based on the Detail
// fields available for each specific type.

func (s *ScopeNotGrantedError) Error() string {
	if s.Detail == nil || len(s.Detail.Scopes) == 0 {
		return s.CommonError.Error()
	}

	return fmt.Sprintf(
		"%s: required scopes: %s",
		s.CommonError.BaseError(), strings.Join(s.Detail.Scopes, ", "),
	)
}
