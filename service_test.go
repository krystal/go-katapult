package katapult

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newPathHelper(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *pathHelper
		wantErr string
	}{
		{
			name: "absolute path",
			args: args{"/core/v1/"},
			want: &pathHelper{basePath: &url.URL{Path: "/core/v1/"}},
		},
		{
			name: "relative path",
			args: args{"core/v1/"},
			want: &pathHelper{basePath: &url.URL{Path: "core/v1/"}},
		},
		{
			name: "no trailing slash",
			args: args{"/core/v1"},
			want: &pathHelper{basePath: &url.URL{Path: "/core/v1"}},
		},
		{
			name: "parse urlencoded spaces",
			args: args{"/core%20api/v1/"},
			want: &pathHelper{basePath: &url.URL{Path: "/core api/v1/"}},
		},
		{
			name:    "invalid path",
			args:    args{"/core%2"},
			wantErr: `parse "/core%2": invalid URL escape "%2"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newPathHelper(tt.args.path)

			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_pathHelper_RequestPath(t *testing.T) {
	type fields struct {
		basePath string
	}
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr string
	}{
		{
			name:   "joins given string with basePath",
			fields: fields{"/core/v1/"},
			args:   args{"data_centers"},
			want:   "/core/v1/data_centers",
		},
		{
			name:   "override basePath with a absolute path",
			fields: fields{"/core/v1/"},
			args:   args{"/data_centers"},
			want:   "/data_centers",
		},
		{
			name:   "used as root when basePath is empty",
			fields: fields{""},
			args:   args{"data_centers"},
			want:   "/data_centers",
		},
		{
			name:   "relative within basePath",
			fields: fields{"/core/v1"},
			args:   args{"data_centers"},
			want:   "/core/data_centers",
		},
		{
			name:   "parse urlencoded spaces",
			fields: fields{"/core/v1/"},
			args:   args{"data%20centers"},
			want:   "/core/v1/data centers",
		},
		{
			name:    "invalid path",
			fields:  fields{"/core/v1"},
			args:    args{`asdf%3`},
			wantErr: `parse "asdf%3": invalid URL escape "%3"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := newPathHelper(tt.fields.basePath)

			got, err := s.RequestPath(tt.args.path)

			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
