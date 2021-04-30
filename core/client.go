package core

import (
	"context"
	"net/http"
	"net/url"

	"github.com/krystal/go-katapult"
)

type Client struct {
	Certificates                    *CertificatesClient
	DNSZones                        *DNSZonesClient
	DataCenters                     *DataCentersClient
	DiskTemplates                   *DiskTemplatesClient
	IPAddresses                     *IPAddressesClient
	LoadBalancers                   *LoadBalancersClient
	NetworkSpeedProfiles            *NetworkSpeedProfilesClient
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

// RequestMaker represents something that the API Clients can use to create
// and submit a request.
type RequestMaker interface {
	Do(req *http.Request, v interface{}) (*katapult.Response, error)
	NewRequestWithContext(
		ctx context.Context,
		method string,
		u *url.URL,
		body interface{},
	) (*http.Request, error)
}

func New(rm RequestMaker) (*Client, error) {
	//nolint:lll
	c := &Client{
		Certificates:                    NewCertificatesClient(rm),
		DNSZones:                        NewDNSZonesClient(rm),
		DataCenters:                     NewDataCentersClient(rm),
		DiskTemplates:                   NewDiskTemplatesClient(rm),
		IPAddresses:                     NewIPAddressesClient(rm),
		LoadBalancers:                   NewLoadBalancersClient(rm),
		NetworkSpeedProfiles:            NewNetworkSpeedProfilesClient(rm),
		Networks:                        NewNetworksClient(rm),
		Organizations:                   NewOrganizationsClient(rm),
		Tasks:                           NewTasksClient(rm),
		TrashObjects:                    NewTrashObjectsClient(rm),
		VirtualMachineBuilds:            NewVirtualMachineBuildsClient(rm),
		VirtualMachineGroups:            NewVirtualMachineGroupsClient(rm),
		VirtualMachineNetworkInterfaces: NewVirtualMachineNetworkInterfacesClient(rm),
		VirtualMachinePackages:          NewVirtualMachinePackagesClient(rm),
		VirtualMachines:                 NewVirtualMachinesClient(rm),
	}

	return c, nil
}
