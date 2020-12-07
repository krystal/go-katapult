package katapult

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testQueryableObj1 struct {
	ID string
}

func (s *testQueryableObj1) queryValues() *url.Values {
	return &url.Values{"obj1[id]": []string{s.ID}}
}

type testQueryableObj2 struct {
	ID string
}

func (s *testQueryableObj2) queryValues() *url.Values {
	return &url.Values{"obj2[id]": []string{s.ID}}
}

type testQueryableObj3 struct {
	ID string
}

func (s *testQueryableObj3) queryValues() *url.Values {
	return &url.Values{"obj3[id]": []string{s.ID}}
}

func Test_queryValues(t *testing.T) {
	type args struct {
		objs []queryable
	}
	tests := []struct {
		name string
		args args
		want *url.Values
	}{
		{
			name: "no objects",
			args: args{
				objs: []queryable{},
			},
			want: &url.Values{},
		},
		{
			name: "nil object",
			args: args{
				objs: []queryable{nil},
			},
			want: &url.Values{},
		},
		{
			name: "single object",
			args: args{
				objs: []queryable{&testQueryableObj1{ID: "abc"}},
			},
			want: &url.Values{"obj1[id]": []string{"abc"}},
		},
		{
			name: "two objects",
			args: args{
				objs: []queryable{
					&testQueryableObj1{ID: "abc"},
					&testQueryableObj2{ID: "def"},
				},
			},
			want: &url.Values{
				"obj1[id]": []string{"abc"},
				"obj2[id]": []string{"def"},
			},
		},
		{
			name: "three objects",
			args: args{
				objs: []queryable{
					&testQueryableObj1{ID: "abc"},
					&testQueryableObj2{ID: "def"},
					&testQueryableObj3{ID: "ghi"},
				},
			},
			want: &url.Values{
				"obj1[id]": []string{"abc"},
				"obj2[id]": []string{"def"},
				"obj3[id]": []string{"ghi"},
			},
		},
		{
			name: "duplicate objects",
			args: args{
				objs: []queryable{
					&testQueryableObj1{ID: "abc"},
					&testQueryableObj1{ID: "def"},
					&testQueryableObj2{ID: "ghi"},
				},
			},
			want: &url.Values{
				"obj1[id]": []string{"abc", "def"},
				"obj2[id]": []string{"ghi"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := queryValues(tt.args.objs...)

			assert.Equal(t, tt.want, got)
		})
	}
}
