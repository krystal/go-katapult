package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

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

	assert.Equal(t, tc, c.Certificates.client)
	assert.Equal(t, tc, c.DNSZones.client)
	assert.Equal(t, tc, c.DataCenters.client)
	assert.Equal(t, tc, c.DiskTemplates.client)
	assert.Equal(t, tc, c.IPAddresses.client)
	assert.Equal(t, tc, c.LoadBalancers.client)
	assert.Equal(t, tc, c.NetworkSpeedProfiles.client)
	assert.Equal(t, tc, c.Networks.client)
	assert.Equal(t, tc, c.Organizations.client)
	assert.Equal(t, tc, c.Tasks.client)
	assert.Equal(t, tc, c.TrashObjects.client)
	assert.Equal(t, tc, c.VirtualMachineBuilds.client)
	assert.Equal(t, tc, c.VirtualMachineGroups.client)
	assert.Equal(t, tc, c.VirtualMachineNetworkInterfaces.client)
	assert.Equal(t, tc, c.VirtualMachinePackages.client)
	assert.Equal(t, tc, c.VirtualMachines.client)
}
