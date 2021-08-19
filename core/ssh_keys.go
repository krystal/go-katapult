package core

import (
	"context"
	"net/url"

	"github.com/krystal/go-katapult"
)

type AuthSSHKey struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Fingerprint string `json:"fingerprint,omitempty"`
}

type SSHKeysClient struct {
	client   RequestMaker
	basePath *url.URL
}

func NewSSHKeysClient(rm RequestMaker) *SSHKeysClient {
	return &SSHKeysClient{
		client:   rm,
		basePath: &url.URL{Path: "/core/v1/"},
	}
}

type sshKeysResponseBody struct {
	Pagination *katapult.Pagination `json:"pagination,omitempty"`
	SSHKey     *AuthSSHKey          `json:"ssh_key,omitempty"`
	SSHKeys    []*AuthSSHKey        `json:"ssh_keys,omitempty"`
}

func (s *SSHKeysClient) List(
	ctx context.Context,
	ref OrganizationRef,
	opts *ListOptions,
	reqOpts ...katapult.RequestOption,
) ([]*AuthSSHKey, *katapult.Response, error) {
	qs := queryValues(opts, ref)
	u := &url.URL{Path: "organizations/_/ssh_keys", RawQuery: qs.Encode()}

	body, resp, err := s.doRequest(ctx, "GET", u, nil, reqOpts...)
	resp.Pagination = body.Pagination

	return body.SSHKeys, resp, err
}

type AuthSSHKeyProperties struct {
	// Name is the SSH keys name.
	Name string `json:"name"`

	// Key is the SSH public key.
	Key string `json:"key"`
}

func (s *SSHKeysClient) Add(
	ctx context.Context,
	ref OrganizationRef,
	properties AuthSSHKeyProperties,
	reqOpts ...katapult.RequestOption,
) (*AuthSSHKey, *katapult.Response, error) {
	qs := ref.queryValues()
	u := &url.URL{Path: "organizations/_/ssh_keys", RawQuery: qs.Encode()}

	body, resp, err := s.doRequest(ctx, "POST", u, properties, reqOpts...)

	return body.SSHKey, resp, err
}

type SSHKeyRef struct {
	ID string `json:"id"`
}

func (kr SSHKeyRef) queryValues() *url.Values {
	return &url.Values{"ssh_key[id]": []string{kr.ID}}
}

func (s *SSHKeysClient) Delete(
	ctx context.Context,
	ref SSHKeyRef,
	reqOpts ...katapult.RequestOption,
) (*AuthSSHKey, *katapult.Response, error) {
	qs := ref.queryValues()
	u := &url.URL{Path: "ssh_keys/_", RawQuery: qs.Encode()}

	body, resp, err := s.doRequest(ctx, "DELETE", u, nil, reqOpts...)

	return body.SSHKey, resp, err
}

func (s *SSHKeysClient) doRequest(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
	reqOpts ...katapult.RequestOption,
) (*sshKeysResponseBody, *katapult.Response, error) {
	u = s.basePath.ResolveReference(u)
	respBody := &sshKeysResponseBody{}

	req := katapult.NewRequest(method, u, body, reqOpts...)
	resp, err := s.client.Do(ctx, req, respBody)
	if resp == nil {
		resp = katapult.NewResponse(nil)
	}

	return respBody, resp, handleResponseError(err)
}
