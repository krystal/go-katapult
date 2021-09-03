package core

import (
	"encoding/json"
	"testing"

	"github.com/krystal/go-katapult"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvalidSpecXMLError_Error(t *testing.T) {
	type fields struct {
		parent      error
		code        string
		description string
		detail      *InvalidSpecXMLErrorDetail
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "nil detail",
			fields: fields{
				parent:      ErrInvalidSpecXML,
				code:        "invalid_spec_xml",
				description: "The spec XML provided is invalid",
			},
			want: "katapult: bad_request: invalid_spec_xml: " +
				"The spec XML provided is invalid",
		},
		{
			name: "empty detail",
			fields: fields{
				parent:      ErrInvalidSpecXML,
				code:        "invalid_spec_xml",
				description: "The spec XML provided is invalid",
				detail:      &InvalidSpecXMLErrorDetail{Errors: ""},
			},
			want: "katapult: bad_request: invalid_spec_xml: " +
				"The spec XML provided is invalid",
		},
		{
			name: "with detail",
			fields: fields{
				parent:      ErrInvalidSpecXML,
				code:        "invalid_spec_xml",
				description: "The spec XML provided is invalid",
				detail:      &InvalidSpecXMLErrorDetail{Errors: "missing disk"},
			},
			want: "katapult: bad_request: invalid_spec_xml: " +
				"missing disk",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detail, err := json.Marshal(tt.fields.detail)
			require.NoError(t, err)

			theError := NewInvalidSpecXMLError(
				&katapult.ResponseError{
					Code:        tt.fields.code,
					Description: tt.fields.description,
					Detail:      json.RawMessage(detail),
				},
			)

			got := theError.Error()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestObjectInTrashError_Error(t *testing.T) {
	type fields struct {
		parent      error
		code        string
		description string
		detail      *ObjectInTrashErrorDetail
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "nil detail",
			fields: fields{
				parent: ErrObjectInTrash,
				code:   "object_in_trash",
				description: "The object found is in the trash and therefore " +
					"cannot be manipulated through the API. It should be " +
					"restored in order to run this operation.",
			},
			want: "katapult: not_acceptable: object_in_trash: The object " +
				"found is in the trash and therefore cannot be manipulated " +
				"through the API. It should be restored in order to run " +
				"this operation.",
		},
		{
			name: "empty detail",
			fields: fields{
				parent: ErrObjectInTrash,
				code:   "object_in_trash",
				description: "The object found is in the trash and therefore " +
					"cannot be manipulated through the API. It should be " +
					"restored in order to run this operation.",
				detail: &ObjectInTrashErrorDetail{TrashObject: nil},
			},
			want: "katapult: not_acceptable: object_in_trash: The object " +
				"found is in the trash and therefore cannot be manipulated " +
				"through the API. It should be restored in order to run " +
				"this operation.",
		},
		{
			name: "with detail",
			fields: fields{
				parent: ErrObjectInTrash,
				code:   "object_in_trash",
				description: "The object found is in the trash and therefore " +
					"cannot be manipulated through the API. It should be " +
					"restored in order to run this operation.",
				detail: &ObjectInTrashErrorDetail{
					TrashObject: &TrashObject{ID: "trsh_47duTKJC6Entt3YN"},
				},
			},
			want: "katapult: not_acceptable: object_in_trash: " +
				"trash_object_id=trsh_47duTKJC6Entt3YN",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detail, err := json.Marshal(tt.fields.detail)
			require.NoError(t, err)

			theError := NewObjectInTrashError(
				&katapult.ResponseError{
					Code:        tt.fields.code,
					Description: tt.fields.description,
					Detail:      json.RawMessage(detail),
				},
			)

			got := theError.Error()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPermissionDeniedError_Error(t *testing.T) {
	type fields struct {
		parent      error
		code        string
		description string
		detail      *PermissionDeniedErrorDetail
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "nil detail",
			fields: fields{
				parent: ErrPermissionDenied,
				code:   "permission_denied",
				description: "The authenticated identity is not permitted to " +
					"perform this action",
			},
			want: "katapult: unauthorized: permission_denied: The " +
				"authenticated identity is not permitted to perform this " +
				"action",
		},
		{
			name: "nil detail",
			fields: fields{
				parent: ErrPermissionDenied,
				code:   "permission_denied",
				description: "The authenticated identity is not permitted to " +
					"perform this action",
				detail: &PermissionDeniedErrorDetail{Details: nil},
			},
			want: "katapult: unauthorized: permission_denied: The " +
				"authenticated identity is not permitted to perform this " +
				"action",
		},
		{
			name: "empty detail",
			fields: fields{
				parent: ErrPermissionDenied,
				code:   "permission_denied",
				description: "The authenticated identity is not permitted to " +
					"perform this action",
				detail: &PermissionDeniedErrorDetail{Details: stringPtr("")},
			},
			want: "katapult: unauthorized: permission_denied: The " +
				"authenticated identity is not permitted to perform this " +
				"action",
		},
		{
			name: "with detail",
			fields: fields{
				parent: ErrPermissionDenied,
				code:   "permission_denied",
				description: "The authenticated identity is not permitted to " +
					"perform this action",
				detail: &PermissionDeniedErrorDetail{
					Details: stringPtr("you are not allowed"),
				},
			},
			want: "katapult: unauthorized: permission_denied: you are not " +
				"allowed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detail, err := json.Marshal(tt.fields.detail)
			require.NoError(t, err)

			theError := NewPermissionDeniedError(
				&katapult.ResponseError{
					Code:        tt.fields.code,
					Description: tt.fields.description,
					Detail:      json.RawMessage(detail),
				},
			)

			got := theError.Error()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRateLimitReachedError_Error(t *testing.T) {
	type fields struct {
		parent      error
		code        string
		description string
		detail      *RateLimitReachedErrorDetail
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "nil detail",
			fields: fields{
				parent: ErrRateLimitReached,
				code:   "resource_creation_restricted",
				description: "You have reached the rate limit for this type " +
					"of request",
			},
			want: "katapult: too_many_requests: rate_limit_reached: You have " +
				"reached the rate limit for this type of request",
		},
		{
			name: "empty detail",
			fields: fields{
				parent: ErrRateLimitReached,
				code:   "resource_creation_restricted",
				description: "You have reached the rate limit for this type " +
					"of request",
				detail: &RateLimitReachedErrorDetail{TotalPermitted: 0},
			},
			want: "katapult: too_many_requests: rate_limit_reached: You have " +
				"reached the rate limit for this type of request",
		},
		{
			name: "with detail",
			fields: fields{
				parent: ErrRateLimitReached,
				code:   "resource_creation_restricted",
				description: "You have reached the rate limit for this type " +
					"of request",
				detail: &RateLimitReachedErrorDetail{TotalPermitted: 200},
			},
			want: "katapult: too_many_requests: rate_limit_reached: max " +
				"requests per minute: 200",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detail, err := json.Marshal(tt.fields.detail)
			require.NoError(t, err)

			theError := NewRateLimitReachedError(
				&katapult.ResponseError{
					Code:        tt.fields.code,
					Description: tt.fields.description,
					Detail:      json.RawMessage(detail),
				},
			)

			got := theError.Error()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestResourceCreationRestrictedError_Error(t *testing.T) {
	type fields struct {
		parent      error
		code        string
		description string
		detail      *ResourceCreationRestrictedErrorDetail
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "nil detail",
			fields: fields{
				parent: ErrResourceCreationRestricted,
				code:   "resource_creation_restricted",
				description: "The organization chosen is not permitted to " +
					"create resources",
			},
			want: "katapult: unauthorized: resource_creation_restricted: " +
				"The organization chosen is not permitted to create resources",
		},
		{
			name: "empty detail",
			fields: fields{
				parent: ErrResourceCreationRestricted,
				code:   "resource_creation_restricted",
				description: "The organization chosen is not permitted to " +
					"create resources",
				detail: &ResourceCreationRestrictedErrorDetail{
					Errors: []string{},
				},
			},
			want: "katapult: unauthorized: resource_creation_restricted: " +
				"The organization chosen is not permitted to create resources",
		},
		{
			name: "with detail",
			fields: fields{
				parent: ErrResourceCreationRestricted,
				code:   "resource_creation_restricted",
				description: "The organization chosen is not permitted to " +
					"create resources",
				detail: &ResourceCreationRestrictedErrorDetail{
					Errors: []string{"IP limit reached", "VM limit reached"},
				},
			},
			want: "katapult: unauthorized: resource_creation_restricted: " +
				"IP limit reached, VM limit reached",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detail, err := json.Marshal(tt.fields.detail)
			require.NoError(t, err)

			theError := NewResourceCreationRestrictedError(
				&katapult.ResponseError{
					Code:        tt.fields.code,
					Description: tt.fields.description,
					Detail:      json.RawMessage(detail),
				},
			)

			got := theError.Error()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTaskQueueingError_Error(t *testing.T) {
	type fields struct {
		parent      error
		code        string
		description string
		detail      *TaskQueueingErrorDetail
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "nil detail",
			fields: fields{
				parent: ErrTaskQueueingError,
				code:   "task_queueing_error",
				description: "This error means that a background task that " +
					"was needed to complete your request could not be queued",
			},
			want: "katapult: not_acceptable: task_queueing_error: This " +
				"error means that a background task that was needed to " +
				"complete your request could not be queued",
		},
		{
			name: "empty detail",
			fields: fields{
				parent: ErrTaskQueueingError,
				code:   "task_queueing_error",
				description: "This error means that a background task that " +
					"was needed to complete your request could not be queued",
				detail: &TaskQueueingErrorDetail{
					Details: "",
				},
			},
			want: "katapult: not_acceptable: task_queueing_error: This " +
				"error means that a background task that was needed to " +
				"complete your request could not be queued",
		},
		{
			name: "with detail",
			fields: fields{
				parent: ErrTaskQueueingError,
				code:   "task_queueing_error",
				description: "This error means that a background task that " +
					"was needed to complete your request could not be queued",
				detail: &TaskQueueingErrorDetail{
					Details: "Cannot start already started VM",
				},
			},
			want: "katapult: not_acceptable: task_queueing_error: Cannot " +
				"start already started VM",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detail, err := json.Marshal(tt.fields.detail)
			require.NoError(t, err)

			theError := NewTaskQueueingError(
				&katapult.ResponseError{
					Code:        tt.fields.code,
					Description: tt.fields.description,
					Detail:      json.RawMessage(detail),
				},
			)

			got := theError.Error()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestValidationError_Error(t *testing.T) {
	type fields struct {
		parent      error
		code        string
		description string
		detail      *ValidationErrorDetail
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "nil detail",
			fields: fields{
				parent: ErrValidationError,
				code:   "validation_error",
				description: "A validation error occurred with the object " +
					"that was being created/updated/deleted",
			},
			want: "katapult: unprocessable_entity: validation_error: A " +
				"validation error occurred with the object that was being " +
				"created/updated/deleted",
		},
		{
			name: "empty detail",
			fields: fields{
				parent: ErrValidationError,
				code:   "validation_error",
				description: "A validation error occurred with the object " +
					"that was being created/updated/deleted",
				detail: &ValidationErrorDetail{
					Errors: []string{},
				},
			},
			want: "katapult: unprocessable_entity: validation_error: A " +
				"validation error occurred with the object that was being " +
				"created/updated/deleted",
		},
		{
			name: "with detail",
			fields: fields{
				parent: ErrValidationError,
				code:   "validation_error",
				description: "A validation error occurred with the object " +
					"that was being created/updated/deleted",
				detail: &ValidationErrorDetail{
					Errors: []string{
						"ip_address must be IPv4", "name cannot be blank",
					},
				},
			},
			want: "katapult: unprocessable_entity: validation_error: " +
				"ip_address must be IPv4, name cannot be blank",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detail, err := json.Marshal(tt.fields.detail)
			require.NoError(t, err)

			theError := NewValidationError(
				&katapult.ResponseError{
					Code:        tt.fields.code,
					Description: tt.fields.description,
					Detail:      json.RawMessage(detail),
				},
			)

			got := theError.Error()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestVirtualMachineMustBeStartedError_Error(t *testing.T) {
	type fields struct {
		parent      error
		code        string
		description string
		detail      *VirtualMachineMustBeStartedErrorDetail
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "nil detail",
			fields: fields{
				parent: ErrVirtualMachineMustBeStarted,
				code:   "virtual_machine_must_be_started",
				description: "Virtual machines must be in a started state " +
					"to create console sessions",
			},
			want: "katapult: not_acceptable: " +
				"virtual_machine_must_be_started: Virtual machines must be " +
				"in a started state to create console sessions",
		},
		{
			name: "empty detail",
			fields: fields{
				parent: ErrVirtualMachineMustBeStarted,
				code:   "virtual_machine_must_be_started",
				description: "Virtual machines must be in a started state " +
					"to create console sessions",
				detail: &VirtualMachineMustBeStartedErrorDetail{
					CurrentState: "",
				},
			},
			want: "katapult: not_acceptable: " +
				"virtual_machine_must_be_started: Virtual machines must be " +
				"in a started state to create console sessions",
		},
		{
			name: "with detail",
			fields: fields{
				parent: ErrVirtualMachineMustBeStarted,
				code:   "virtual_machine_must_be_started",
				description: "Virtual machines must be in a started state " +
					"to create console sessions",
				detail: &VirtualMachineMustBeStartedErrorDetail{
					CurrentState: "stopped",
				},
			},
			want: "katapult: not_acceptable: " +
				"virtual_machine_must_be_started: current_state=stopped",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detail, err := json.Marshal(tt.fields.detail)
			require.NoError(t, err)

			theError := NewVirtualMachineMustBeStartedError(
				&katapult.ResponseError{
					Code:        tt.fields.code,
					Description: tt.fields.description,
					Detail:      json.RawMessage(detail),
				},
			)

			got := theError.Error()

			assert.Equal(t, tt.want, got)
		})
	}
}
