package next_test

import (
	"testing"

	"github.com/krystal/go-katapult/next"
	"github.com/stretchr/testify/require"
)

func ptr[T any](v T) *T {
	return &v
}

func TestDerefOrEmpty_string(t *testing.T) {
	tests := []struct {
		name string
		in   *string
		want string
	}{
		{
			name: "Nil input",
			in:   nil,
			want: "",
		},
		{
			name: "Non-nil input",
			in:   ptr("hello"),
			want: "hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, next.DerefOrEmpty(tt.in))
		})
	}
}

func TestDerefOrEmpty_int(t *testing.T) {
	tests := []struct {
		name string
		in   *int
		want int
	}{
		{
			name: "Nil input",
			in:   nil,
			want: 0,
		},
		{
			name: "Non-nil input",
			in:   ptr(22),
			want: 22,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, next.DerefOrEmpty(tt.in))
		})
	}
}

func TestDerefOrEmpty_float64(t *testing.T) {
	tests := []struct {
		name string
		in   *float64
		want float64
	}{
		{
			name: "Nil input",
			in:   nil,
			want: 0.0,
		},
		{
			name: "Non-nil input",
			in:   ptr(3.14),
			want: 3.14,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, next.DerefOrEmpty(tt.in))
		})
	}
}

func TestDerefOrEmpty_slice(t *testing.T) {
	tests := []struct {
		name string
		in   *[]string
		want []string
	}{
		{
			name: "Nil input",
			in:   nil,
			want: []string{},
		},
		{
			name: "Non-nil input",
			in:   ptr([]string{"hello", "world"}),
			want: []string{"hello", "world"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, next.DerefOrEmpty(tt.in))
		})
	}
}

func TestDerefOrEmpty_map(t *testing.T) {
	tests := []struct {
		name string
		in   *map[string]int
		want map[string]int
	}{
		{
			name: "Nil input",
			in:   nil,
			want: map[string]int{},
		},
		{
			name: "Non-nil input",
			in:   ptr(map[string]int{"hello": 32}),
			want: map[string]int{"hello": 32},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, next.DerefOrEmpty(tt.in))
		})
	}
}

func TestDerefOrEmpty_struct(t *testing.T) {
	type testStruct struct {
		Name string
		Age  int
	}

	tests := []struct {
		name string
		in   *testStruct
		want testStruct
	}{
		{
			name: "Nil input",
			in:   nil,
			want: testStruct{},
		},
		{
			name: "Non-nil input",
			in:   ptr(testStruct{Name: "Alice", Age: 22}),
			want: testStruct{Name: "Alice", Age: 22},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, next.DerefOrEmpty(tt.in))
		})
	}
}
