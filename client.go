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

	Certificates           *CertificatesClient
	DNSZones               *DNSZonesClient
	DataCenters            *DataCentersClient
	LoadBalancers          *LoadBalancersClient
	Networks               *NetworksClient
	Organizations          *OrganizationsClient
	VirtualMachinePackages *VirtualMachinePackagesClient
	VirtualMachines        *VirtualMachinesClient
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
		Certificates:           newCertificatesClient(c),
		DNSZones:               newDNSZonesClient(c),
		DataCenters:            newDataCentersClient(c),
		LoadBalancers:          newLoadBalancersClient(c),
		Networks:               newNetworksClient(c),
		Organizations:          newOrganizationsClient(c),
		VirtualMachinePackages: newVirtualMachinePackagesClient(c),
		VirtualMachines:        newVirtualMachinesClient(c),
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
