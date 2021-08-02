package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"unsafe"

	"github.com/krystal/go-katapult"
	"github.com/krystal/go-katapult/internal/test"
	"github.com/krystal/go-katapult/internal/testclient"
	"github.com/stretchr/testify/assert"
)

var (
	fixtureInvalidAPITokenErr = "invalid_api_token: The API token provided " +
		"was not valid (it may not exist or have expired)"
	fixtureInvalidAPITokenResponseError = &katapult.ResponseError{
		Code: "invalid_api_token",
		Description: "The API token provided was not valid " +
			"(it may not exist or have expired)",
		Detail: json.RawMessage(`{}`),
	}

	//nolint:lll
	fixturePermissionDeniedErr = "katapult: unauthorized: permission_denied: " +
		"Additional information regarding the reason why permission was denied"
	fixturePermissionDeniedResponseError = &katapult.ResponseError{
		Code: "permission_denied",
		Description: "The authenticated identity is not permitted to perform " +
			"this action",
		//nolint:lll
		Detail: json.RawMessage(`{
      "details": "Additional information regarding the reason why permission was denied"
    }`),
	}

	//nolint:lll
	fixtureValidationErrorErr = "katapult: unprocessable_entity: " +
		"validation_error: Failed reticulating 3-dimensional splines, " +
		"Failed preparing captive simulators"
	fixtureValidationErrorResponseError = &katapult.ResponseError{
		Code: "validation_error",
		Description: "A validation error occurred with the object that was " +
			"being created/updated/deleted",
		Detail: json.RawMessage(`{
      "errors": [
        "Failed reticulating 3-dimensional splines",
        "Failed preparing captive simulators"
      ]
    }`,
		),
	}

	fixtureObjectInTrashErr = "katapult: not_acceptable: object_in_trash: " +
		"The object found is in the trash and therefore cannot be " +
		"manipulated through the API. It should be restored in order to run " +
		"this operation."
	fixtureObjectInTrashResponseError = &katapult.ResponseError{
		Code: "object_in_trash",
		Description: "The object found is in the trash and therefore cannot " +
			"be manipulated through the API. It should be restored in order " +
			"to run this operation.",
		Detail: json.RawMessage(`{}`),
	}
)

//
// Helpers
//

func assertFieldSpec(t *testing.T, r *http.Request, spec string) {
	assert.Equal(t, spec, r.Header.Get("X-Field-Spec"))
}

func assertEmptyFieldSpec(t *testing.T, r *http.Request) {
	assertFieldSpec(t, r, "")
}

func assertCustomAuthorization(t *testing.T, r *http.Request, apiKey string) {
	assert.Equal(t,
		fmt.Sprintf("Bearer %s", apiKey), r.Header.Get("Authorization"),
	)
}

func assertAuthorization(t *testing.T, r *http.Request) {
	assertCustomAuthorization(t, r, test.APIKey)
}

//
// Tests
//

func TestClientImplementsRequestMaker(t *testing.T) {
	assert.Implements(t, (*RequestMaker)(nil), new(katapult.Client))
}

// Perhaps consider golden/testdata integration for generating request/response
// data.

func TestRequestMaker(t *testing.T) {
	assert.Implements(t, (*RequestMaker)(nil), &katapult.Client{})
	assert.Implements(t, (*RequestMaker)(nil), &testclient.Client{})
}

func TestNew(t *testing.T) {
	t.Parallel()

	tc := &testclient.Client{}
	c := New(tc)

	rv := reflect.ValueOf(c).Elem()

	// Check every field on the Client struct, this effectively performs:
	//  assert.Equal(t, tc, c.<FieldName>.client)
	for i := 0; i < rv.NumField(); i++ {
		name := rv.Type().Field(i).Name
		f := rv.Field(i)

		if f.IsNil() {
			assert.Fail(t, "Client field: "+name+" (is nil)")
		} else {
			clientField := f.Elem().FieldByName("client")

			// Get value of unexported field "client" without panic.
			value := reflect.NewAt(
				clientField.Type(), unsafe.Pointer(clientField.UnsafeAddr()),
			).Elem().Interface()

			assert.Equal(t, tc, value, "Client field: "+name)
		}
	}
}
