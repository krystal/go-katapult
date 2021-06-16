package katapult

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScopeNotGrantedError_Error(t *testing.T) {
	type fields struct {
		parent      error
		code        string
		description string
		detail      *ScopeNotGrantedErrorDetail
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "without detail",
			fields: fields{
				parent: ErrScopeNotGrantedError,
				code:   "scope_not_granted",
				description: "The scope required for this endpoint has not " +
					"been granted to the authenticating identity",
			},
			want: "katapult: unauthorized: scope_not_granted: The scope " +
				"required for this endpoint has not been granted to the " +
				"authenticating identity",
		},
		{
			name: "empty detail scopes",
			fields: fields{
				parent: ErrScopeNotGrantedError,
				code:   "scope_not_granted",
				description: "The scope required for this endpoint has not " +
					"been granted to the authenticating identity",
				detail: &ScopeNotGrantedErrorDetail{Scopes: []string{}},
			},
			want: "katapult: unauthorized: scope_not_granted: The scope " +
				"required for this endpoint has not been granted to the " +
				"authenticating identity",
		},
		{
			name: "with detail scopes",
			fields: fields{
				parent: ErrScopeNotGrantedError,
				code:   "scope_not_granted",
				description: "The scope required for this endpoint has not " +
					"been granted to the authenticating identity",
				detail: &ScopeNotGrantedErrorDetail{
					Scopes: []string{"core.disks", "core.disks:read"},
				},
			},
			want: "katapult: unauthorized: scope_not_granted: required " +
				"scopes: core.disks, core.disks:read",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detail, err := json.Marshal(tt.fields.detail)
			require.NoError(t, err)

			theError := NewScopeNotGrantedError(
				&ResponseError{
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
