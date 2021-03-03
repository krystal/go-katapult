package katapult

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/krystal/go-katapult/internal/codec"
)

const (
	DefaultUserAgent = "go-katapult"
	DefaultTimeout   = time.Second * 60
)

var ErrClient = fmt.Errorf("%wclient", Err)

type Config struct {
	APIKey     string
	UserAgent  string
	HTTPClient *http.Client

	// BaseURL should only be relevant for testing purposes.
	BaseURL *url.URL
}

type Client struct {
	apiClient *apiClient

	Certificates                    *CertificatesClient
	DNSZones                        *DNSZonesClient
	DataCenters                     *DataCentersClient
	DiskTemplates                   *DiskTemplatesClient
	IPAddresses                     *IPAddressesClient
	LoadBalancers                   *LoadBalancersClient
	Networks                        *NetworksClient
	Organizations                   *OrganizationsClient
	Tasks                           *TasksClient
	TrashObjects                    *TrashObjectsClient
	VirtualMachineBuilds            *VirtualMachineBuildsClient
	VirtualMachineGroups            *VirtualMachineGroupsClient
	VirtualMachineNetworkInterfaces *VirtualMachineNetworkInterfacesClient
	VirtualMachinePackages          *VirtualMachinePackagesClient
	VirtualMachines                 *VirtualMachinesClient
}

func NewClient(config *Config) (*Client, error) {
	ac := &apiClient{
		httpClient: &http.Client{Timeout: DefaultTimeout},
		codec:      &codec.JSON{},
		BaseURL:    &url.URL{Scheme: "https", Host: "api.katapult.io"},
		UserAgent:  DefaultUserAgent,
	}

	//nolint:lll
	c := &Client{
		apiClient:                       ac,
		Certificates:                    newCertificatesClient(ac),
		DNSZones:                        newDNSZonesClient(ac),
		DataCenters:                     newDataCentersClient(ac),
		DiskTemplates:                   newDiskTemplatesClient(ac),
		IPAddresses:                     newIPAddressesClient(ac),
		LoadBalancers:                   newLoadBalancersClient(ac),
		Networks:                        newNetworksClient(ac),
		Organizations:                   newOrganizationsClient(ac),
		Tasks:                           newTasksClient(ac),
		TrashObjects:                    newTrashObjectsClient(ac),
		VirtualMachineBuilds:            newVirtualMachineBuildsClient(ac),
		VirtualMachineGroups:            newVirtualMachineGroupsClient(ac),
		VirtualMachineNetworkInterfaces: newVirtualMachineNetworkInterfacesClient(ac),
		VirtualMachinePackages:          newVirtualMachinePackagesClient(ac),
		VirtualMachines:                 newVirtualMachinesClient(ac),
	}

	err := c.configure(config)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Client) configure(config *Config) error {
	if config != nil {
		if config.APIKey != "" {
			c.SetAPIKey(config.APIKey)
		}

		if config.UserAgent != "" {
			c.SetUserAgent(config.UserAgent)
		}

		if config.BaseURL != nil {
			err := c.SetBaseURL(config.BaseURL)
			if err != nil {
				return err
			}
		}

		if config.HTTPClient != nil {
			_ = c.SetHTTPClient(config.HTTPClient)
		}
	}

	return nil
}

func (c *Client) APIKey() string {
	return c.apiClient.APIKey
}

func (c *Client) SetAPIKey(apiKey string) {
	c.apiClient.APIKey = apiKey
}

func (c *Client) UserAgent() string {
	return c.apiClient.UserAgent
}

func (c *Client) SetUserAgent(agent string) {
	c.apiClient.UserAgent = agent
}

func (c *Client) BaseURL() *url.URL {
	return c.apiClient.BaseURL
}

func (c *Client) SetBaseURL(u *url.URL) error {
	switch {
	case u == nil:
		return fmt.Errorf("%w: base URL cannot be nil", ErrClient)
	case u.Scheme == "":
		return fmt.Errorf("%w: base URL scheme is empty", ErrClient)
	case u.Host == "":
		return fmt.Errorf("%w: base URL host is empty", ErrClient)
	}

	c.apiClient.BaseURL = u

	return nil
}

func (c *Client) HTTPClient() *http.Client {
	return c.apiClient.httpClient
}

func (c *Client) SetHTTPClient(httpClient *http.Client) error {
	if httpClient == nil {
		return fmt.Errorf("%w: http client cannot be nil", ErrClient)
	}

	c.apiClient.httpClient = httpClient

	return nil
}
