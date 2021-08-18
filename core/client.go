package core

import (
	"context"

	"github.com/krystal/go-katapult"
)

type Client struct {
	Certificates                    *CertificatesClient
	DNSZones                        *DNSZonesClient
	DataCenters                     *DataCentersClient
	DiskTemplates                   *DiskTemplatesClient
	IPAddresses                     *IPAddressesClient
	LoadBalancers                   *LoadBalancersClient
	LoadBalancerRules               *LoadBalancerRulesClient
	NetworkSpeedProfiles            *NetworkSpeedProfilesClient
	Networks                        *NetworksClient
	Organizations                   *OrganizationsClient
	SecurityGroups                  *SecurityGroupsClient
	SecurityGroupRules              *SecurityGroupRulesClient
	Tags                            *TagsClient
	SSHKeys                         *SSHKeysClient
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
	Do(
		ctx context.Context,
		req *katapult.Request,
		v interface{},
	) (*katapult.Response, error)
}

func New(rm RequestMaker) *Client {
	//nolint:lll
	c := &Client{
		Certificates:                    NewCertificatesClient(rm),
		DNSZones:                        NewDNSZonesClient(rm),
		DataCenters:                     NewDataCentersClient(rm),
		DiskTemplates:                   NewDiskTemplatesClient(rm),
		IPAddresses:                     NewIPAddressesClient(rm),
		LoadBalancers:                   NewLoadBalancersClient(rm),
		LoadBalancerRules:               NewLoadBalancerRulesClient(rm),
		NetworkSpeedProfiles:            NewNetworkSpeedProfilesClient(rm),
		Networks:                        NewNetworksClient(rm),
		Organizations:                   NewOrganizationsClient(rm),
		SecurityGroups:                  NewSecurityGroupsClient(rm),
		SecurityGroupRules:              NewSecurityGroupRulesClient(rm),
		Tags:                            NewTagsClient(rm),
		SSHKeys:                         NewSSHKeysClient(rm),
		Tasks:                           NewTasksClient(rm),
		TrashObjects:                    NewTrashObjectsClient(rm),
		VirtualMachineBuilds:            NewVirtualMachineBuildsClient(rm),
		VirtualMachineGroups:            NewVirtualMachineGroupsClient(rm),
		VirtualMachineNetworkInterfaces: NewVirtualMachineNetworkInterfacesClient(rm),
		VirtualMachinePackages:          NewVirtualMachinePackagesClient(rm),
		VirtualMachines:                 NewVirtualMachinesClient(rm),
	}

	return c
}
