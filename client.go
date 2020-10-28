package katapult

import (
	"net/http"
	"net/url"
	"time"

	"github.com/krystal/go-katapult/internal/codec"
)

const (
	userAgent      = "go-katapult"
	defaultTimeout = time.Second * 60
)

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	apiClient *apiClient

	Certificates           *CertificatesResource
	DNSZones               *DNSZonesResource
	DataCenters            *DataCentersResource
	Networks               *NetworksResource
	Organizations          *OrganizationsResource
	VirtualMachinePackages *VirtualMachinePackagesResource
	VirtualMachines        *VirtualMachinesResource
}

func NewClient(httpClient HTTPClient) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultTimeout}
	}

	baseURL := &url.URL{Scheme: "https", Host: "api.katapult.io"}

	c := &apiClient{
		httpClient: httpClient,
		codec:      &codec.JSON{},
		BaseURL:    baseURL,
		UserAgent:  userAgent,
	}

	return &Client{
		apiClient:              c,
		Certificates:           newCertificatesResource(c),
		DNSZones:               newDNSZonesResource(c),
		DataCenters:            newDataCentersResource(c),
		Networks:               newNetworksResource(c),
		Organizations:          newOrganizationsResource(c),
		VirtualMachinePackages: newVirtualMachinePackagesResource(c),
		VirtualMachines:        newVirtualMachinesResource(c),
	}
}

func (c *Client) BaseURL() *url.URL {
	return c.apiClient.BaseURL
}

func (c *Client) SetBaseURL(u *url.URL) {
	if u != nil {
		c.apiClient.BaseURL = u
	}
}

func (c *Client) UserAgent() string {
	return c.apiClient.UserAgent
}

func (c *Client) SetUserAgent(agent string) {
	if agent != "" {
		c.apiClient.UserAgent = agent
	}
}
