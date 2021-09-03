package core

import (
	"fmt"
	"strings"
)

// This file contains custom Error() functions for select error types in
// errors_generated.go, allowing customized error messages based on the Detail
// fields available for each specific type.

func (s *InvalidSpecXMLError) Error() string {
	if s.Detail == nil || s.Detail.Errors == "" {
		return s.CommonError.Error()
	}

	return fmt.Sprintf("%s: %s", s.CommonError.BaseError(), s.Detail.Errors)
}

func (s *ObjectInTrashError) Error() string {
	if s.Detail == nil || s.Detail.TrashObject == nil {
		return s.CommonError.Error()
	}

	return fmt.Sprintf(
		"%s: trash_object_id=%s",
		s.CommonError.BaseError(), s.Detail.TrashObject.ID,
	)
}

func (s *PermissionDeniedError) Error() string {
	if s.Detail == nil || s.Detail.Details == nil || *s.Detail.Details == "" {
		return s.CommonError.Error()
	}

	return fmt.Sprintf("%s: %s", s.CommonError.BaseError(), *s.Detail.Details)
}

func (s *RateLimitReachedError) Error() string {
	if s.Detail == nil || s.Detail.TotalPermitted == 0 {
		return s.CommonError.Error()
	}

	return fmt.Sprintf(
		"%s: max requests per minute: %d",
		s.CommonError.BaseError(), s.Detail.TotalPermitted,
	)
}

func (s *ResourceCreationRestrictedError) Error() string {
	if s.Detail == nil || len(s.Detail.Errors) == 0 {
		return s.CommonError.Error()
	}

	return fmt.Sprintf(
		"%s: %s",
		s.CommonError.BaseError(), strings.Join(s.Detail.Errors, ", "),
	)
}

func (s *TaskQueueingError) Error() string {
	if s.Detail == nil || s.Detail.Details == "" {
		return s.CommonError.Error()
	}

	return fmt.Sprintf("%s: %s", s.CommonError.BaseError(), s.Detail.Details)
}

func (s *ValidationError) Error() string {
	if s.Detail == nil || len(s.Detail.Errors) == 0 {
		return s.CommonError.Error()
	}

	return fmt.Sprintf(
		"%s: %s",
		s.CommonError.BaseError(), strings.Join(s.Detail.Errors, ", "),
	)
}

func (s *VirtualMachineMustBeStartedError) Error() string {
	if s.Detail == nil || s.Detail.CurrentState == "" {
		return s.CommonError.Error()
	}

	return fmt.Sprintf(
		"%s: current_state=%s",
		s.CommonError.BaseError(), s.Detail.CurrentState,
	)
}
