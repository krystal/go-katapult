package core

import (
	"testing"

	"github.com/krystal/go-katapult"
	"github.com/stretchr/testify/assert"
)

func TestErr(t *testing.T) {
	assert.EqualError(t, Err, "katapult: core")
	assert.ErrorIs(t, Err, katapult.Err)
}
