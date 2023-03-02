package test

import (
	"context"
	"testing"

	"github.com/jimeh/rands"
	"github.com/stretchr/testify/assert"
)

type ctxKey string

var testCtxKey ctxKey = "test"

// WithContext returns a new child context with the test context key set to
// value.
func WithContext(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, testCtxKey, value)
}

// FromContext returns the test context value.
func FromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if value, ok := ctx.Value(testCtxKey).(string); ok {
		return value
	}

	return ""
}

// Context returns a new context.Context with a random value assigned to the
// test key value using WithContext. Returns nil if parent context is nil.
func Context(parent context.Context) context.Context {
	if parent == nil {
		return nil
	}

	uniq, err := rands.Alphanumeric(16)
	if err != nil {
		panic(err)
	}

	return WithContext(parent, uniq)
}

// AssertContext asserts that both given contexts has the same test key value
// set.
//
//nolint:golint,revive
func AssertContext(t *testing.T, want, got context.Context) {
	assert.NotNil(t, want, "wanted context is nil")
	assert.NotNil(t, got, "got context is nil")
	assert.Equal(t,
		FromContext(want), FromContext(got),
		"context.Context do not match",
	)
}
