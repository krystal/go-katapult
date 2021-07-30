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
	SSHKey  *AuthSSHKey   `json:"ssh_key,omitempty"`
	SSHKeys []*AuthSSHKey `json:"ssh_keys,omitempty"`
}

func (s *SSHKeysClient) List(
	ctx context.Context,
	ref OrganizationRef,
) ([]*AuthSSHKey, *katapult.Response, error) {
	qs := ref.queryValues()
	u := &url.URL{Path: "organizations/_/ssh_keys", RawQuery: qs.Encode()}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)

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
) (*AuthSSHKey, *katapult.Response, error) {
	qs := ref.queryValues()
	u := &url.URL{Path: "organizations/_/ssh_keys", RawQuery: qs.Encode()}

	body, resp, err := s.doRequest(ctx, "POST", u, properties)

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
) (*AuthSSHKey, *katapult.Response, error) {
	qs := ref.queryValues()
	u := &url.URL{Path: "ssh_keys/_", RawQuery: qs.Encode()}

	body, resp, err := s.doRequest(ctx, "DELETE", u, nil)

	return body.SSHKey, resp, err
}

func (s *SSHKeysClient) doRequest(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
) (*sshKeysResponseBody, *katapult.Response, error) {
	u = s.basePath.ResolveReference(u)
	respBody := &sshKeysResponseBody{}

	req := katapult.NewRequest(method, u, body)
	resp, err := s.client.Do(ctx, req, respBody)
	if resp == nil {
		resp = katapult.NewResponse(nil)
	}

	return respBody, resp, handleResponseError(err)
}
