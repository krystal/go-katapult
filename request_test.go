package katapult

import (
	"bytes"
	"io"
	"math"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func strPtr(s string) *string {
	return &s
}

func TestNewRequest(t *testing.T) {
	type reqBody struct {
		Hello string `json:"hello,omitempty"`
	}

	type args struct {
		method string
		u      *url.URL
		body   interface{}
		opts   []RequestOption
	}
	tests := []struct {
		name string
		args args
		want *Request
	}{
		{
			name: "without body",
			args: args{
				method: "GET",
				u:      &url.URL{Path: "/foo/bar", RawQuery: "?hello=world"},
			},
			want: &Request{
				Method: "GET",
				URL: &url.URL{
					Path:     "/foo/bar",
					RawQuery: "?hello=world",
				},
				ContentType: "",
				Body:        nil,
				Header:      map[string][]string{},
			},
		},
		{
			name: "with custom header",
			args: args{
				method: "GET",
				u:      &url.URL{Path: "/foo/bar", RawQuery: "?hello=world"},
				opts: []RequestOption{
					RequestSetHeader(
						"X-Clacks-Overhead",
						"GNU Terry Pratchett",
					),
				},
			},
			want: &Request{
				Method: "GET",
				URL: &url.URL{
					Path:     "/foo/bar",
					RawQuery: "?hello=world",
				},
				ContentType: "",
				Body:        nil,
				Header: map[string][]string{
					"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				},
			},
		},
		{
			name: "with struct body",
			args: args{
				method: "POST",
				u:      &url.URL{Path: "/foo/bar", RawQuery: "?hello=world"},
				body:   reqBody{Hello: "world"},
			},
			want: &Request{
				Method: "POST",
				URL: &url.URL{
					Path:     "/foo/bar",
					RawQuery: "?hello=world",
				},
				ContentType: "",
				Body:        reqBody{Hello: "world"},
				Header:      map[string][]string{},
			},
		},
		{
			name: "with io.Reader body",
			args: args{
				method: "POST",
				u:      &url.URL{Path: "/foo/bar", RawQuery: "?hello=world"},
				body:   bytes.NewBufferString("hello"),
			},
			want: &Request{
				Method: "POST",
				URL: &url.URL{
					Path:     "/foo/bar",
					RawQuery: "?hello=world",
				},
				ContentType: "",
				Body:        bytes.NewBufferString("hello"),
				Header:      map[string][]string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRequest(
				tt.args.method,
				tt.args.u,
				tt.args.body,
				tt.args.opts...,
			)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRequest_bodyContent(t *testing.T) {
	type Book struct {
		Title  string `json:"title"`
		Author string `json:"author,omitempty"`
	}

	type fields struct {
		ContentType string
		Body        interface{}
	}
	tests := []struct {
		name            string
		fields          fields
		wantContentType string
		wantBody        *string
		wantErr         string
		wantErrIs       error
	}{
		{
			name:            "nil body",
			fields:          fields{},
			wantContentType: "",
			wantBody:        nil,
		},
		{
			name: "struct body, no content type",
			fields: fields{
				Body: &Book{Title: "Nothing", Author: "John Doe"},
			},
			wantContentType: "application/json",
			wantBody: strPtr(
				"{\"title\":\"Nothing\",\"author\":\"John Doe\"}\n",
			),
		},
		{
			name: "struct body with HTML in field values",
			fields: fields{
				Body: &Book{Title: "<b>Nothing</b>", Author: "John <i>Doe</i>"},
			},
			wantContentType: "application/json",
			wantBody: strPtr(
				"{\"title\":\"<b>Nothing</b>\"," +
					"\"author\":\"John <i>Doe</i>\"}\n",
			),
		},
		{
			name: "io.Reader body and content type",
			fields: fields{
				ContentType: "text/csv",
				Body: bytes.NewBufferString(
					"title,author\nNothing,John Doe",
				),
			},
			wantContentType: "text/csv",
			wantBody:        strPtr("title,author\nNothing,John Doe"),
		},
		{
			name: "struct body and content type",
			fields: fields{
				ContentType: "application/json",
				Body:        &Book{Title: "Nothing", Author: "John Doe"},
			},
			wantErr: "katapult: request: Body must be a io.Reader when " +
				"ContentType is set",
			wantErrIs: ErrRequest,
		},
		{
			name: "invalid type for json marshaling",
			fields: fields{
				Body: make(chan int),
			},
			wantErr: "json: unsupported type: chan int",
		},
		{
			name: "invalid value for json marshaling",
			fields: fields{
				Body: math.Inf(1),
			},
			wantErr: "json: unsupported value: +Inf",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Request{
				ContentType: tt.fields.ContentType,
				Body:        tt.fields.Body,
			}
			contentType, body, err := r.bodyContent()

			if tt.wantErrIs != nil {
				assert.ErrorIs(t, err, tt.wantErrIs)
			}

			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
				assert.Equal(t, "", contentType)
				assert.Nil(t, body)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantContentType, contentType)
				if tt.wantBody == nil {
					assert.Nil(t, body)
				} else {
					b, err := io.ReadAll(body)
					require.NoError(t, err)
					assert.Equal(t, *tt.wantBody, string(b))
				}
			}
		})
	}
}
