package core

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/krystal/go-katapult"
	"github.com/krystal/go-katapult/internal/test"
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
	fixturePermissionDeniedErr = "permission_denied: The authenticated " +
		"identity is not permitted to perform this action -- " +
		"{\n  \"details\": \"Additional information regarding the reason why permission was denied\"\n}"
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
	fixtureValidationErrorErr = "validation_error: A validation error " +
		"occurred with the object that was being created/updated/deleted -- " +
		"{\n  \"errors\": [\n    \"Failed reticulating 3-dimensional splines\",\n    \"Failed preparing captive simulators\"\n  ]\n}"
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

	fixtureObjectInTrashErr = "object_in_trash: The object found is in the " +
		"trash and therefore cannot be manipulated through the API. It " +
		"should be restored in order to run this operation."
	fixtureObjectInTrashResponseError = &katapult.ResponseError{
		Code: "object_in_trash",
		Description: "The object found is in the trash and therefore cannot " +
			"be manipulated through the API. It should be restored in order " +
			"to run this operation.",
		Detail: json.RawMessage(`{}`),
	}

	fixtureInvalidArgumentErr = "invalid_argument: The 'X' argument " +
		"is invalid"
	fixtureInvalidArgumentResponseError = &katapult.ResponseError{
		Code:        "invalid_argument",
		Description: "The 'X' argument is invalid",
		Detail:      json.RawMessage(`{}`),
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

// TODO: Flesh this out for use for unit testing
// Perhaps consider golden/testdata integration for generating request/response
// data.
type fakeRequestMaker struct{}

func (frm *fakeRequestMaker) Do(
	req *http.Request,
	val interface{},
) (*katapult.Response, error) {
	return nil, nil
}

func (frm *fakeRequestMaker) NewRequestWithContext(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
) (*http.Request, error) {
	return nil, nil
}

func TestNew(t *testing.T) {
	t.Parallel()

	frm := &fakeRequestMaker{}
	c := New(frm)

	assert.Equal(t, frm, c.Certificates.client)
	assert.Equal(t, frm, c.DNSZones.client)
	assert.Equal(t, frm, c.DataCenters.client)
	assert.Equal(t, frm, c.DiskTemplates.client)
	assert.Equal(t, frm, c.IPAddresses.client)
	assert.Equal(t, frm, c.LoadBalancers.client)
	assert.Equal(t, frm, c.NetworkSpeedProfiles.client)
	assert.Equal(t, frm, c.Networks.client)
	assert.Equal(t, frm, c.Organizations.client)
	assert.Equal(t, frm, c.Tasks.client)
	assert.Equal(t, frm, c.TrashObjects.client)
	assert.Equal(t, frm, c.VirtualMachineBuilds.client)
	assert.Equal(t, frm, c.VirtualMachineGroups.client)
	assert.Equal(t, frm, c.VirtualMachineNetworkInterfaces.client)
	assert.Equal(t, frm, c.VirtualMachinePackages.client)
	assert.Equal(t, frm, c.VirtualMachines.client)
}
