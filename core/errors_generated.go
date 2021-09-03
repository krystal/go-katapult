package core

import (
	"encoding/json"
	"fmt"

	"github.com/krystal/go-katapult"
)

// Code generated by github.com/krystal/go-katapult/tools/codegen. DO NOT EDIT.

var (
	ErrCertificateNotFound                        = fmt.Errorf("%w: certificate_not_found", katapult.ErrResourceNotFound)
	ErrCountryNotFound                            = fmt.Errorf("%w: country_not_found", katapult.ErrResourceNotFound)
	ErrCountryStateNotFound                       = fmt.Errorf("%w: country_state_not_found", katapult.ErrResourceNotFound)
	ErrCurrencyNotFound                           = fmt.Errorf("%w: currency_not_found", katapult.ErrResourceNotFound)
	ErrDNSRecordNotFound                          = fmt.Errorf("%w: dns_record_not_found", katapult.ErrResourceNotFound)
	ErrDNSZoneNotFound                            = fmt.Errorf("%w: dns_zone_not_found", katapult.ErrResourceNotFound)
	ErrDNSZoneNotVerified                         = fmt.Errorf("%w: dns_zone_not_verified", katapult.ErrUnprocessableEntity)
	ErrDataCenterNotFound                         = fmt.Errorf("%w: data_center_not_found", katapult.ErrResourceNotFound)
	ErrDeletionRestricted                         = fmt.Errorf("%w: deletion_restricted", katapult.ErrConflict)
	ErrDiskBackupPolicyNotFound                   = fmt.Errorf("%w: disk_backup_policy_not_found", katapult.ErrResourceNotFound)
	ErrDiskNotFound                               = fmt.Errorf("%w: disk_not_found", katapult.ErrResourceNotFound)
	ErrDiskTemplateNotFound                       = fmt.Errorf("%w: disk_template_not_found", katapult.ErrResourceNotFound)
	ErrDiskTemplateVersionNotFound                = fmt.Errorf("%w: disk_template_version_not_found", katapult.ErrResourceNotFound)
	ErrFlexibleResourcesUnavailableToOrganization = fmt.Errorf("%w: flexible_resources_unavailable_to_organization", katapult.ErrForbidden)
	ErrIPAddressNotFound                          = fmt.Errorf("%w: ip_address_not_found", katapult.ErrResourceNotFound)
	ErrIPAlreadyAllocated                         = fmt.Errorf("%w: ip_already_allocated", katapult.ErrUnprocessableEntity)
	ErrIdentityNotLinkedToWebSession              = fmt.Errorf("%w: identity_not_linked_to_web_session", katapult.ErrBadRequest)
	ErrInterfaceNotFound                          = fmt.Errorf("%w: interface_not_found", katapult.ErrResourceNotFound)
	ErrInvalidIP                                  = fmt.Errorf("%w: invalid_ip", katapult.ErrUnprocessableEntity)
	ErrInvalidSpecXML                             = fmt.Errorf("%w: invalid_spec_xml", katapult.ErrBadRequest)
	ErrLoadBalancerNotFound                       = fmt.Errorf("%w: load_balancer_not_found", katapult.ErrResourceNotFound)
	ErrLoadBalancerRuleNotFound                   = fmt.Errorf("%w: load_balancer_rule_not_found", katapult.ErrResourceNotFound)
	ErrLocationRequired                           = fmt.Errorf("%w: location_required", katapult.ErrUnprocessableEntity)
	ErrNetworkNotFound                            = fmt.Errorf("%w: network_not_found", katapult.ErrResourceNotFound)
	ErrNetworkSpeedProfileNotFound                = fmt.Errorf("%w: network_speed_profile_not_found", katapult.ErrResourceNotFound)
	ErrNoAllocation                               = fmt.Errorf("%w: no_allocation", katapult.ErrUnprocessableEntity)
	ErrNoAvailableAddresses                       = fmt.Errorf("%w: no_available_addresses", katapult.ErrServiceUnavailable)
	ErrNoInterfaceAvailable                       = fmt.Errorf("%w: no_interface_available", katapult.ErrUnprocessableEntity)
	ErrNoUserAssociatedWithIdentity               = fmt.Errorf("%w: no_user_associated_with_identity", katapult.ErrResourceNotFound)
	ErrObjectInTrash                              = fmt.Errorf("%w: object_in_trash", katapult.ErrNotAcceptable)
	ErrOperatingSystemNotFound                    = fmt.Errorf("%w: operating_system_not_found", katapult.ErrResourceNotFound)
	ErrOrganizationLimitReached                   = fmt.Errorf("%w: organization_limit_reached", katapult.ErrUnprocessableEntity)
	ErrOrganizationNotActivated                   = fmt.Errorf("%w: organization_not_activated", katapult.ErrForbidden)
	ErrOrganizationNotFound                       = fmt.Errorf("%w: organization_not_found", katapult.ErrResourceNotFound)
	ErrOrganizationSuspended                      = fmt.Errorf("%w: organization_suspended", katapult.ErrForbidden)
	ErrPermissionDenied                           = fmt.Errorf("%w: permission_denied", katapult.ErrForbidden)
	ErrRateLimitReached                           = fmt.Errorf("%w: rate_limit_reached", katapult.ErrTooManyRequests)
	ErrResourceCreationRestricted                 = fmt.Errorf("%w: resource_creation_restricted", katapult.ErrForbidden)
	ErrResourceDoesNotSupportUnallocation         = fmt.Errorf("%w: resource_does_not_support_unallocation", katapult.ErrConflict)
	ErrSSHKeyNotFound                             = fmt.Errorf("%w: ssh_key_not_found", katapult.ErrResourceNotFound)
	ErrSecurityGroupNotFound                      = fmt.Errorf("%w: security_group_not_found", katapult.ErrResourceNotFound)
	ErrSecurityGroupRuleNotFound                  = fmt.Errorf("%w: security_group_rule_not_found", katapult.ErrResourceNotFound)
	ErrSpeedProfileAlreadyAssigned                = fmt.Errorf("%w: speed_profile_already_assigned", katapult.ErrUnprocessableEntity)
	ErrTagNotFound                                = fmt.Errorf("%w: tag_not_found", katapult.ErrResourceNotFound)
	ErrTaskNotFound                               = fmt.Errorf("%w: task_not_found", katapult.ErrResourceNotFound)
	ErrTaskQueueingError                          = fmt.Errorf("%w: task_queueing_error", katapult.ErrNotAcceptable)
	ErrTrashObjectNotFound                        = fmt.Errorf("%w: trash_object_not_found", katapult.ErrResourceNotFound)
	ErrValidationError                            = fmt.Errorf("%w: validation_error", katapult.ErrUnprocessableEntity)
	ErrVirtualMachineBuildNotFound                = fmt.Errorf("%w: build_not_found", katapult.ErrResourceNotFound)
	ErrVirtualMachineGroupNotFound                = fmt.Errorf("%w: virtual_machine_group_not_found", katapult.ErrResourceNotFound)
	ErrVirtualMachineMustBeStarted                = fmt.Errorf("%w: virtual_machine_must_be_started", katapult.ErrNotAcceptable)
	ErrVirtualMachineNetworkInterfaceNotFound     = fmt.Errorf("%w: virtual_machine_network_interface_not_found", katapult.ErrResourceNotFound)
	ErrVirtualMachineNotFound                     = fmt.Errorf("%w: virtual_machine_not_found", katapult.ErrResourceNotFound)
	ErrVirtualMachinePackageNotFound              = fmt.Errorf("%w: package_not_found", katapult.ErrResourceNotFound)
	ErrZoneNotFound                               = fmt.Errorf("%w: zone_not_found", katapult.ErrResourceNotFound)
)

// CertificateNotFoundError:
// No certificate was found matching any of the criteria provided in the arguments
type CertificateNotFoundError struct {
	katapult.CommonError
}

func NewCertificateNotFoundError(theError *katapult.ResponseError) *CertificateNotFoundError {
	return &CertificateNotFoundError{
		CommonError: katapult.NewCommonError(
			ErrCertificateNotFound,
			"certificate_not_found",
			theError.Description,
		),
	}
}

// CountryNotFoundError:
// No countries was found matching any of the criteria provided in the arguments
type CountryNotFoundError struct {
	katapult.CommonError
}

func NewCountryNotFoundError(theError *katapult.ResponseError) *CountryNotFoundError {
	return &CountryNotFoundError{
		CommonError: katapult.NewCommonError(
			ErrCountryNotFound,
			"country_not_found",
			theError.Description,
		),
	}
}

// CountryStateNotFoundError:
// No country state was found matching any of the criteria provided in the arguments
type CountryStateNotFoundError struct {
	katapult.CommonError
}

func NewCountryStateNotFoundError(theError *katapult.ResponseError) *CountryStateNotFoundError {
	return &CountryStateNotFoundError{
		CommonError: katapult.NewCommonError(
			ErrCountryStateNotFound,
			"country_state_not_found",
			theError.Description,
		),
	}
}

// CurrencyNotFoundError:
// No currencies was found matching any of the criteria provided in the arguments
type CurrencyNotFoundError struct {
	katapult.CommonError
}

func NewCurrencyNotFoundError(theError *katapult.ResponseError) *CurrencyNotFoundError {
	return &CurrencyNotFoundError{
		CommonError: katapult.NewCommonError(
			ErrCurrencyNotFound,
			"currency_not_found",
			theError.Description,
		),
	}
}

// DNSRecordNotFoundError:
// No DNS record was found matching any of the criteria provided in the arguments
type DNSRecordNotFoundError struct {
	katapult.CommonError
}

func NewDNSRecordNotFoundError(theError *katapult.ResponseError) *DNSRecordNotFoundError {
	return &DNSRecordNotFoundError{
		CommonError: katapult.NewCommonError(
			ErrDNSRecordNotFound,
			"dns_record_not_found",
			theError.Description,
		),
	}
}

// DNSZoneNotFoundError:
// No DNS zone was found matching any of the criteria provided in the arguments
type DNSZoneNotFoundError struct {
	katapult.CommonError
}

func NewDNSZoneNotFoundError(theError *katapult.ResponseError) *DNSZoneNotFoundError {
	return &DNSZoneNotFoundError{
		CommonError: katapult.NewCommonError(
			ErrDNSZoneNotFound,
			"dns_zone_not_found",
			theError.Description,
		),
	}
}

// DNSZoneNotVerifiedError:
// The DNS zone could not be verified, check the nameservers are set correctly
type DNSZoneNotVerifiedError struct {
	katapult.CommonError
}

func NewDNSZoneNotVerifiedError(theError *katapult.ResponseError) *DNSZoneNotVerifiedError {
	return &DNSZoneNotVerifiedError{
		CommonError: katapult.NewCommonError(
			ErrDNSZoneNotVerified,
			"dns_zone_not_verified",
			theError.Description,
		),
	}
}

// DataCenterNotFoundError:
// No data centers was found matching any of the criteria provided in the arguments
type DataCenterNotFoundError struct {
	katapult.CommonError
}

func NewDataCenterNotFoundError(theError *katapult.ResponseError) *DataCenterNotFoundError {
	return &DataCenterNotFoundError{
		CommonError: katapult.NewCommonError(
			ErrDataCenterNotFound,
			"data_center_not_found",
			theError.Description,
		),
	}
}

// DeletionRestrictedError:
// Object cannot be deleted
type DeletionRestrictedError struct {
	katapult.CommonError
	Detail *DeletionRestrictedErrorDetail `json:"detail,omitempty"`
}

func NewDeletionRestrictedError(theError *katapult.ResponseError) *DeletionRestrictedError {
	detail := &DeletionRestrictedErrorDetail{}
	err := json.Unmarshal(theError.Detail, detail)
	if err != nil {
		detail = nil
	}

	return &DeletionRestrictedError{
		CommonError: katapult.NewCommonError(
			ErrDeletionRestricted,
			"deletion_restricted",
			theError.Description,
		),
		Detail: detail,
	}
}

type DeletionRestrictedErrorDetail struct {
	Errors []string `json:"errors,omitempty"`
}

// DiskBackupPolicyNotFoundError:
// No disk backup policies was found matching any of the criteria provided in the arguments
type DiskBackupPolicyNotFoundError struct {
	katapult.CommonError
}

func NewDiskBackupPolicyNotFoundError(theError *katapult.ResponseError) *DiskBackupPolicyNotFoundError {
	return &DiskBackupPolicyNotFoundError{
		CommonError: katapult.NewCommonError(
			ErrDiskBackupPolicyNotFound,
			"disk_backup_policy_not_found",
			theError.Description,
		),
	}
}

// DiskNotFoundError:
// No disks was found matching any of the criteria provided in the arguments
type DiskNotFoundError struct {
	katapult.CommonError
}

func NewDiskNotFoundError(theError *katapult.ResponseError) *DiskNotFoundError {
	return &DiskNotFoundError{
		CommonError: katapult.NewCommonError(
			ErrDiskNotFound,
			"disk_not_found",
			theError.Description,
		),
	}
}

// DiskTemplateNotFoundError:
// No disk template was found matching any of the criteria provided in the arguments
type DiskTemplateNotFoundError struct {
	katapult.CommonError
}

func NewDiskTemplateNotFoundError(theError *katapult.ResponseError) *DiskTemplateNotFoundError {
	return &DiskTemplateNotFoundError{
		CommonError: katapult.NewCommonError(
			ErrDiskTemplateNotFound,
			"disk_template_not_found",
			theError.Description,
		),
	}
}

// DiskTemplateVersionNotFoundError:
// No disk template version was found matching any of the criteria provided in the arguments
type DiskTemplateVersionNotFoundError struct {
	katapult.CommonError
}

func NewDiskTemplateVersionNotFoundError(theError *katapult.ResponseError) *DiskTemplateVersionNotFoundError {
	return &DiskTemplateVersionNotFoundError{
		CommonError: katapult.NewCommonError(
			ErrDiskTemplateVersionNotFound,
			"disk_template_version_not_found",
			theError.Description,
		),
	}
}

// FlexibleResourcesUnavailableToOrganizationError:
// The organization is not permitted to use flexible resources
type FlexibleResourcesUnavailableToOrganizationError struct {
	katapult.CommonError
}

func NewFlexibleResourcesUnavailableToOrganizationError(theError *katapult.ResponseError) *FlexibleResourcesUnavailableToOrganizationError {
	return &FlexibleResourcesUnavailableToOrganizationError{
		CommonError: katapult.NewCommonError(
			ErrFlexibleResourcesUnavailableToOrganization,
			"flexible_resources_unavailable_to_organization",
			theError.Description,
		),
	}
}

// IPAddressNotFoundError:
// No IP addresses were found matching any of the criteria provided in the arguments
type IPAddressNotFoundError struct {
	katapult.CommonError
}

func NewIPAddressNotFoundError(theError *katapult.ResponseError) *IPAddressNotFoundError {
	return &IPAddressNotFoundError{
		CommonError: katapult.NewCommonError(
			ErrIPAddressNotFound,
			"ip_address_not_found",
			theError.Description,
		),
	}
}

// IPAlreadyAllocatedError:
// This IP address has already been allocated to another resource
type IPAlreadyAllocatedError struct {
	katapult.CommonError
}

func NewIPAlreadyAllocatedError(theError *katapult.ResponseError) *IPAlreadyAllocatedError {
	return &IPAlreadyAllocatedError{
		CommonError: katapult.NewCommonError(
			ErrIPAlreadyAllocated,
			"ip_already_allocated",
			theError.Description,
		),
	}
}

// IdentityNotLinkedToWebSessionError:
// The authenticated identity is not linked to a web session
type IdentityNotLinkedToWebSessionError struct {
	katapult.CommonError
}

func NewIdentityNotLinkedToWebSessionError(theError *katapult.ResponseError) *IdentityNotLinkedToWebSessionError {
	return &IdentityNotLinkedToWebSessionError{
		CommonError: katapult.NewCommonError(
			ErrIdentityNotLinkedToWebSession,
			"identity_not_linked_to_web_session",
			theError.Description,
		),
	}
}

// InterfaceNotFoundError:
// An interface could not be found for the specified network
type InterfaceNotFoundError struct {
	katapult.CommonError
}

func NewInterfaceNotFoundError(theError *katapult.ResponseError) *InterfaceNotFoundError {
	return &InterfaceNotFoundError{
		CommonError: katapult.NewCommonError(
			ErrInterfaceNotFound,
			"interface_not_found",
			theError.Description,
		),
	}
}

// InvalidIPError:
// This IP address is not valid for this network interface.
type InvalidIPError struct {
	katapult.CommonError
}

func NewInvalidIPError(theError *katapult.ResponseError) *InvalidIPError {
	return &InvalidIPError{
		CommonError: katapult.NewCommonError(
			ErrInvalidIP,
			"invalid_ip",
			theError.Description,
		),
	}
}

// InvalidSpecXMLError:
// The spec XML provided is invalid
type InvalidSpecXMLError struct {
	katapult.CommonError
	Detail *InvalidSpecXMLErrorDetail `json:"detail,omitempty"`
}

func NewInvalidSpecXMLError(theError *katapult.ResponseError) *InvalidSpecXMLError {
	detail := &InvalidSpecXMLErrorDetail{}
	err := json.Unmarshal(theError.Detail, detail)
	if err != nil {
		detail = nil
	}

	return &InvalidSpecXMLError{
		CommonError: katapult.NewCommonError(
			ErrInvalidSpecXML,
			"invalid_spec_xml",
			theError.Description,
		),
		Detail: detail,
	}
}

type InvalidSpecXMLErrorDetail struct {
	Errors string `json:"errors,omitempty"`
}

// LoadBalancerNotFoundError:
// No load balancer was found matching any of the criteria provided in the arguments
type LoadBalancerNotFoundError struct {
	katapult.CommonError
}

func NewLoadBalancerNotFoundError(theError *katapult.ResponseError) *LoadBalancerNotFoundError {
	return &LoadBalancerNotFoundError{
		CommonError: katapult.NewCommonError(
			ErrLoadBalancerNotFound,
			"load_balancer_not_found",
			theError.Description,
		),
	}
}

// LoadBalancerRuleNotFoundError:
// No load balancer rule was found matching any of the criteria provided in the arguments
type LoadBalancerRuleNotFoundError struct {
	katapult.CommonError
}

func NewLoadBalancerRuleNotFoundError(theError *katapult.ResponseError) *LoadBalancerRuleNotFoundError {
	return &LoadBalancerRuleNotFoundError{
		CommonError: katapult.NewCommonError(
			ErrLoadBalancerRuleNotFound,
			"load_balancer_rule_not_found",
			theError.Description,
		),
	}
}

// LocationRequiredError:
// A zone or a data_center argument must be provided
type LocationRequiredError struct {
	katapult.CommonError
}

func NewLocationRequiredError(theError *katapult.ResponseError) *LocationRequiredError {
	return &LocationRequiredError{
		CommonError: katapult.NewCommonError(
			ErrLocationRequired,
			"location_required",
			theError.Description,
		),
	}
}

// NetworkNotFoundError:
// No network was found matching any of the criteria provided in the arguments
type NetworkNotFoundError struct {
	katapult.CommonError
}

func NewNetworkNotFoundError(theError *katapult.ResponseError) *NetworkNotFoundError {
	return &NetworkNotFoundError{
		CommonError: katapult.NewCommonError(
			ErrNetworkNotFound,
			"network_not_found",
			theError.Description,
		),
	}
}

// NetworkSpeedProfileNotFoundError:
// No network speed profile was found matching any of the criteria provided in the arguments
type NetworkSpeedProfileNotFoundError struct {
	katapult.CommonError
}

func NewNetworkSpeedProfileNotFoundError(theError *katapult.ResponseError) *NetworkSpeedProfileNotFoundError {
	return &NetworkSpeedProfileNotFoundError{
		CommonError: katapult.NewCommonError(
			ErrNetworkSpeedProfileNotFound,
			"network_speed_profile_not_found",
			theError.Description,
		),
	}
}

// NoAllocationError:
// This IP address is not currently allocated to any object, and cannot be unallocated.
type NoAllocationError struct {
	katapult.CommonError
}

func NewNoAllocationError(theError *katapult.ResponseError) *NoAllocationError {
	return &NoAllocationError{
		CommonError: katapult.NewCommonError(
			ErrNoAllocation,
			"no_allocation",
			theError.Description,
		),
	}
}

// NoAvailableAddressesError:
// Our pool of addresses for that version seems to have run dry. If this issue continues, please contact support.
type NoAvailableAddressesError struct {
	katapult.CommonError
}

func NewNoAvailableAddressesError(theError *katapult.ResponseError) *NoAvailableAddressesError {
	return &NoAvailableAddressesError{
		CommonError: katapult.NewCommonError(
			ErrNoAvailableAddresses,
			"no_available_addresses",
			theError.Description,
		),
	}
}

// NoInterfaceAvailableError:
// This virtual machine does not have a network interface that is compatible with the provided IP address
type NoInterfaceAvailableError struct {
	katapult.CommonError
}

func NewNoInterfaceAvailableError(theError *katapult.ResponseError) *NoInterfaceAvailableError {
	return &NoInterfaceAvailableError{
		CommonError: katapult.NewCommonError(
			ErrNoInterfaceAvailable,
			"no_interface_available",
			theError.Description,
		),
	}
}

// NoUserAssociatedWithIdentityError:
// There is no user associated with this API token
type NoUserAssociatedWithIdentityError struct {
	katapult.CommonError
}

func NewNoUserAssociatedWithIdentityError(theError *katapult.ResponseError) *NoUserAssociatedWithIdentityError {
	return &NoUserAssociatedWithIdentityError{
		CommonError: katapult.NewCommonError(
			ErrNoUserAssociatedWithIdentity,
			"no_user_associated_with_identity",
			theError.Description,
		),
	}
}

// ObjectInTrashError:
// The object found is in the trash and therefore cannot be manipulated through the API. It should be restored in order to run this operation.
type ObjectInTrashError struct {
	katapult.CommonError
	Detail *ObjectInTrashErrorDetail `json:"detail,omitempty"`
}

func NewObjectInTrashError(theError *katapult.ResponseError) *ObjectInTrashError {
	detail := &ObjectInTrashErrorDetail{}
	err := json.Unmarshal(theError.Detail, detail)
	if err != nil {
		detail = nil
	}

	return &ObjectInTrashError{
		CommonError: katapult.NewCommonError(
			ErrObjectInTrash,
			"object_in_trash",
			theError.Description,
		),
		Detail: detail,
	}
}

type ObjectInTrashErrorDetail struct {
	TrashObject *TrashObject `json:"trash_object,omitempty"`
}

// OperatingSystemNotFoundError:
// No operating system was found matching any of the criteria provided in the arguments
type OperatingSystemNotFoundError struct {
	katapult.CommonError
}

func NewOperatingSystemNotFoundError(theError *katapult.ResponseError) *OperatingSystemNotFoundError {
	return &OperatingSystemNotFoundError{
		CommonError: katapult.NewCommonError(
			ErrOperatingSystemNotFound,
			"operating_system_not_found",
			theError.Description,
		),
	}
}

// OrganizationLimitReachedError:
// The maxmium number of organizations that can be created has been reached
type OrganizationLimitReachedError struct {
	katapult.CommonError
}

func NewOrganizationLimitReachedError(theError *katapult.ResponseError) *OrganizationLimitReachedError {
	return &OrganizationLimitReachedError{
		CommonError: katapult.NewCommonError(
			ErrOrganizationLimitReached,
			"organization_limit_reached",
			theError.Description,
		),
	}
}

// OrganizationNotActivatedError:
// An organization was found from the arguments provided but it wasn't activated yet
type OrganizationNotActivatedError struct {
	katapult.CommonError
}

func NewOrganizationNotActivatedError(theError *katapult.ResponseError) *OrganizationNotActivatedError {
	return &OrganizationNotActivatedError{
		CommonError: katapult.NewCommonError(
			ErrOrganizationNotActivated,
			"organization_not_activated",
			theError.Description,
		),
	}
}

// OrganizationNotFoundError:
// No organization was found matching any of the criteria provided in the arguments
type OrganizationNotFoundError struct {
	katapult.CommonError
}

func NewOrganizationNotFoundError(theError *katapult.ResponseError) *OrganizationNotFoundError {
	return &OrganizationNotFoundError{
		CommonError: katapult.NewCommonError(
			ErrOrganizationNotFound,
			"organization_not_found",
			theError.Description,
		),
	}
}

// OrganizationSuspendedError:
// An organization was found from the arguments provided but it was suspended
type OrganizationSuspendedError struct {
	katapult.CommonError
}

func NewOrganizationSuspendedError(theError *katapult.ResponseError) *OrganizationSuspendedError {
	return &OrganizationSuspendedError{
		CommonError: katapult.NewCommonError(
			ErrOrganizationSuspended,
			"organization_suspended",
			theError.Description,
		),
	}
}

// PermissionDeniedError:
// The authenticated identity is not permitted to perform this action
type PermissionDeniedError struct {
	katapult.CommonError
	Detail *PermissionDeniedErrorDetail `json:"detail,omitempty"`
}

func NewPermissionDeniedError(theError *katapult.ResponseError) *PermissionDeniedError {
	detail := &PermissionDeniedErrorDetail{}
	err := json.Unmarshal(theError.Detail, detail)
	if err != nil {
		detail = nil
	}

	return &PermissionDeniedError{
		CommonError: katapult.NewCommonError(
			ErrPermissionDenied,
			"permission_denied",
			theError.Description,
		),
		Detail: detail,
	}
}

type PermissionDeniedErrorDetail struct {
	Details *string `json:"details,omitempty"`
}

// RateLimitReachedError:
// You have reached the rate limit for this type of request
type RateLimitReachedError struct {
	katapult.CommonError
	Detail *RateLimitReachedErrorDetail `json:"detail,omitempty"`
}

func NewRateLimitReachedError(theError *katapult.ResponseError) *RateLimitReachedError {
	detail := &RateLimitReachedErrorDetail{}
	err := json.Unmarshal(theError.Detail, detail)
	if err != nil {
		detail = nil
	}

	return &RateLimitReachedError{
		CommonError: katapult.NewCommonError(
			ErrRateLimitReached,
			"rate_limit_reached",
			theError.Description,
		),
		Detail: detail,
	}
}

type RateLimitReachedErrorDetail struct {
	TotalPermitted int `json:"total_permitted,omitempty"`
}

// ResourceCreationRestrictedError:
// The organization chosen is not permitted to create resources
type ResourceCreationRestrictedError struct {
	katapult.CommonError
	Detail *ResourceCreationRestrictedErrorDetail `json:"detail,omitempty"`
}

func NewResourceCreationRestrictedError(theError *katapult.ResponseError) *ResourceCreationRestrictedError {
	detail := &ResourceCreationRestrictedErrorDetail{}
	err := json.Unmarshal(theError.Detail, detail)
	if err != nil {
		detail = nil
	}

	return &ResourceCreationRestrictedError{
		CommonError: katapult.NewCommonError(
			ErrResourceCreationRestricted,
			"resource_creation_restricted",
			theError.Description,
		),
		Detail: detail,
	}
}

type ResourceCreationRestrictedErrorDetail struct {
	Errors []string `json:"errors,omitempty"`
}

// ResourceDoesNotSupportUnallocationError:
// The resource this IP address belongs to does not allow you to unallocate it.
type ResourceDoesNotSupportUnallocationError struct {
	katapult.CommonError
}

func NewResourceDoesNotSupportUnallocationError(theError *katapult.ResponseError) *ResourceDoesNotSupportUnallocationError {
	return &ResourceDoesNotSupportUnallocationError{
		CommonError: katapult.NewCommonError(
			ErrResourceDoesNotSupportUnallocation,
			"resource_does_not_support_unallocation",
			theError.Description,
		),
	}
}

// SSHKeyNotFoundError:
// No SSH keys was found matching any of the criteria provided in the arguments
type SSHKeyNotFoundError struct {
	katapult.CommonError
}

func NewSSHKeyNotFoundError(theError *katapult.ResponseError) *SSHKeyNotFoundError {
	return &SSHKeyNotFoundError{
		CommonError: katapult.NewCommonError(
			ErrSSHKeyNotFound,
			"ssh_key_not_found",
			theError.Description,
		),
	}
}

// SecurityGroupNotFoundError:
// No security group was found matching any of the criteria provided in the arguments
type SecurityGroupNotFoundError struct {
	katapult.CommonError
}

func NewSecurityGroupNotFoundError(theError *katapult.ResponseError) *SecurityGroupNotFoundError {
	return &SecurityGroupNotFoundError{
		CommonError: katapult.NewCommonError(
			ErrSecurityGroupNotFound,
			"security_group_not_found",
			theError.Description,
		),
	}
}

// SecurityGroupRuleNotFoundError:
// No security group rule was found matching any of the criteria provided in the arguments
type SecurityGroupRuleNotFoundError struct {
	katapult.CommonError
}

func NewSecurityGroupRuleNotFoundError(theError *katapult.ResponseError) *SecurityGroupRuleNotFoundError {
	return &SecurityGroupRuleNotFoundError{
		CommonError: katapult.NewCommonError(
			ErrSecurityGroupRuleNotFound,
			"security_group_rule_not_found",
			theError.Description,
		),
	}
}

// SpeedProfileAlreadyAssignedError:
// This network speed profile is already assigned to this virtual machine network interface.
type SpeedProfileAlreadyAssignedError struct {
	katapult.CommonError
}

func NewSpeedProfileAlreadyAssignedError(theError *katapult.ResponseError) *SpeedProfileAlreadyAssignedError {
	return &SpeedProfileAlreadyAssignedError{
		CommonError: katapult.NewCommonError(
			ErrSpeedProfileAlreadyAssigned,
			"speed_profile_already_assigned",
			theError.Description,
		),
	}
}

// TagNotFoundError:
// No tags was found matching any of the criteria provided in the arguments
type TagNotFoundError struct {
	katapult.CommonError
}

func NewTagNotFoundError(theError *katapult.ResponseError) *TagNotFoundError {
	return &TagNotFoundError{
		CommonError: katapult.NewCommonError(
			ErrTagNotFound,
			"tag_not_found",
			theError.Description,
		),
	}
}

// TaskNotFoundError:
// No task was found matching any of the criteria provided in the arguments
type TaskNotFoundError struct {
	katapult.CommonError
}

func NewTaskNotFoundError(theError *katapult.ResponseError) *TaskNotFoundError {
	return &TaskNotFoundError{
		CommonError: katapult.NewCommonError(
			ErrTaskNotFound,
			"task_not_found",
			theError.Description,
		),
	}
}

// TaskQueueingError:
// This error means that a background task that was needed to complete your request could not be queued
type TaskQueueingError struct {
	katapult.CommonError
	Detail *TaskQueueingErrorDetail `json:"detail,omitempty"`
}

func NewTaskQueueingError(theError *katapult.ResponseError) *TaskQueueingError {
	detail := &TaskQueueingErrorDetail{}
	err := json.Unmarshal(theError.Detail, detail)
	if err != nil {
		detail = nil
	}

	return &TaskQueueingError{
		CommonError: katapult.NewCommonError(
			ErrTaskQueueingError,
			"task_queueing_error",
			theError.Description,
		),
		Detail: detail,
	}
}

type TaskQueueingErrorDetail struct {
	Details string `json:"details,omitempty"`
}

// TrashObjectNotFoundError:
// No trash object was found matching any of the criteria provided in the arguments
type TrashObjectNotFoundError struct {
	katapult.CommonError
}

func NewTrashObjectNotFoundError(theError *katapult.ResponseError) *TrashObjectNotFoundError {
	return &TrashObjectNotFoundError{
		CommonError: katapult.NewCommonError(
			ErrTrashObjectNotFound,
			"trash_object_not_found",
			theError.Description,
		),
	}
}

// ValidationError:
// A validation error occurred with the object that was being created/updated/deleted
type ValidationError struct {
	katapult.CommonError
	Detail *ValidationErrorDetail `json:"detail,omitempty"`
}

func NewValidationError(theError *katapult.ResponseError) *ValidationError {
	detail := &ValidationErrorDetail{}
	err := json.Unmarshal(theError.Detail, detail)
	if err != nil {
		detail = nil
	}

	return &ValidationError{
		CommonError: katapult.NewCommonError(
			ErrValidationError,
			"validation_error",
			theError.Description,
		),
		Detail: detail,
	}
}

type ValidationErrorDetail struct {
	Errors []string `json:"errors,omitempty"`
}

// VirtualMachineBuildNotFoundError:
// No build was found matching any of the criteria provided in the arguments
type VirtualMachineBuildNotFoundError struct {
	katapult.CommonError
}

func NewVirtualMachineBuildNotFoundError(theError *katapult.ResponseError) *VirtualMachineBuildNotFoundError {
	return &VirtualMachineBuildNotFoundError{
		CommonError: katapult.NewCommonError(
			ErrVirtualMachineBuildNotFound,
			"build_not_found",
			theError.Description,
		),
	}
}

// VirtualMachineGroupNotFoundError:
// No virtual machine group was found matching any of the criteria provided in the arguments
type VirtualMachineGroupNotFoundError struct {
	katapult.CommonError
}

func NewVirtualMachineGroupNotFoundError(theError *katapult.ResponseError) *VirtualMachineGroupNotFoundError {
	return &VirtualMachineGroupNotFoundError{
		CommonError: katapult.NewCommonError(
			ErrVirtualMachineGroupNotFound,
			"virtual_machine_group_not_found",
			theError.Description,
		),
	}
}

// VirtualMachineMustBeStartedError:
// Virtual machines must be in a started state to create console sessions
type VirtualMachineMustBeStartedError struct {
	katapult.CommonError
	Detail *VirtualMachineMustBeStartedErrorDetail `json:"detail,omitempty"`
}

func NewVirtualMachineMustBeStartedError(theError *katapult.ResponseError) *VirtualMachineMustBeStartedError {
	detail := &VirtualMachineMustBeStartedErrorDetail{}
	err := json.Unmarshal(theError.Detail, detail)
	if err != nil {
		detail = nil
	}

	return &VirtualMachineMustBeStartedError{
		CommonError: katapult.NewCommonError(
			ErrVirtualMachineMustBeStarted,
			"virtual_machine_must_be_started",
			theError.Description,
		),
		Detail: detail,
	}
}

type VirtualMachineMustBeStartedErrorDetail struct {
	CurrentState string `json:"current_state,omitempty"`
}

// VirtualMachineNetworkInterfaceNotFoundError:
// No network interface was found matching any of the criteria provided in the arguments
type VirtualMachineNetworkInterfaceNotFoundError struct {
	katapult.CommonError
}

func NewVirtualMachineNetworkInterfaceNotFoundError(theError *katapult.ResponseError) *VirtualMachineNetworkInterfaceNotFoundError {
	return &VirtualMachineNetworkInterfaceNotFoundError{
		CommonError: katapult.NewCommonError(
			ErrVirtualMachineNetworkInterfaceNotFound,
			"virtual_machine_network_interface_not_found",
			theError.Description,
		),
	}
}

// VirtualMachineNotFoundError:
// No virtual machine was found matching any of the criteria provided in the arguments
type VirtualMachineNotFoundError struct {
	katapult.CommonError
}

func NewVirtualMachineNotFoundError(theError *katapult.ResponseError) *VirtualMachineNotFoundError {
	return &VirtualMachineNotFoundError{
		CommonError: katapult.NewCommonError(
			ErrVirtualMachineNotFound,
			"virtual_machine_not_found",
			theError.Description,
		),
	}
}

// VirtualMachinePackageNotFoundError:
// No package was found matching any of the criteria provided in the arguments
type VirtualMachinePackageNotFoundError struct {
	katapult.CommonError
}

func NewVirtualMachinePackageNotFoundError(theError *katapult.ResponseError) *VirtualMachinePackageNotFoundError {
	return &VirtualMachinePackageNotFoundError{
		CommonError: katapult.NewCommonError(
			ErrVirtualMachinePackageNotFound,
			"package_not_found",
			theError.Description,
		),
	}
}

// ZoneNotFoundError:
// No zone was found matching any of the criteria provided in the arguments
type ZoneNotFoundError struct {
	katapult.CommonError
}

func NewZoneNotFoundError(theError *katapult.ResponseError) *ZoneNotFoundError {
	return &ZoneNotFoundError{
		CommonError: katapult.NewCommonError(
			ErrZoneNotFound,
			"zone_not_found",
			theError.Description,
		),
	}
}

// castResponseError casts a *katapult.ResponseError to a more specific type based on the error's Code value.
func castResponseError(theError *katapult.ResponseError) error {
	switch theError.Code {
	case "certificate_not_found":
		return NewCertificateNotFoundError(theError)
	case "country_not_found":
		return NewCountryNotFoundError(theError)
	case "country_state_not_found":
		return NewCountryStateNotFoundError(theError)
	case "currency_not_found":
		return NewCurrencyNotFoundError(theError)
	case "dns_record_not_found":
		return NewDNSRecordNotFoundError(theError)
	case "dns_zone_not_found":
		return NewDNSZoneNotFoundError(theError)
	case "dns_zone_not_verified":
		return NewDNSZoneNotVerifiedError(theError)
	case "data_center_not_found":
		return NewDataCenterNotFoundError(theError)
	case "deletion_restricted":
		return NewDeletionRestrictedError(theError)
	case "disk_backup_policy_not_found":
		return NewDiskBackupPolicyNotFoundError(theError)
	case "disk_not_found":
		return NewDiskNotFoundError(theError)
	case "disk_template_not_found":
		return NewDiskTemplateNotFoundError(theError)
	case "disk_template_version_not_found":
		return NewDiskTemplateVersionNotFoundError(theError)
	case "flexible_resources_unavailable_to_organization":
		return NewFlexibleResourcesUnavailableToOrganizationError(theError)
	case "ip_address_not_found":
		return NewIPAddressNotFoundError(theError)
	case "ip_already_allocated":
		return NewIPAlreadyAllocatedError(theError)
	case "identity_not_linked_to_web_session":
		return NewIdentityNotLinkedToWebSessionError(theError)
	case "interface_not_found":
		return NewInterfaceNotFoundError(theError)
	case "invalid_ip":
		return NewInvalidIPError(theError)
	case "invalid_spec_xml":
		return NewInvalidSpecXMLError(theError)
	case "load_balancer_not_found":
		return NewLoadBalancerNotFoundError(theError)
	case "load_balancer_rule_not_found":
		return NewLoadBalancerRuleNotFoundError(theError)
	case "location_required":
		return NewLocationRequiredError(theError)
	case "network_not_found":
		return NewNetworkNotFoundError(theError)
	case "network_speed_profile_not_found":
		return NewNetworkSpeedProfileNotFoundError(theError)
	case "no_allocation":
		return NewNoAllocationError(theError)
	case "no_available_addresses":
		return NewNoAvailableAddressesError(theError)
	case "no_interface_available":
		return NewNoInterfaceAvailableError(theError)
	case "no_user_associated_with_identity":
		return NewNoUserAssociatedWithIdentityError(theError)
	case "object_in_trash":
		return NewObjectInTrashError(theError)
	case "operating_system_not_found":
		return NewOperatingSystemNotFoundError(theError)
	case "organization_limit_reached":
		return NewOrganizationLimitReachedError(theError)
	case "organization_not_activated":
		return NewOrganizationNotActivatedError(theError)
	case "organization_not_found":
		return NewOrganizationNotFoundError(theError)
	case "organization_suspended":
		return NewOrganizationSuspendedError(theError)
	case "permission_denied":
		return NewPermissionDeniedError(theError)
	case "rate_limit_reached":
		return NewRateLimitReachedError(theError)
	case "resource_creation_restricted":
		return NewResourceCreationRestrictedError(theError)
	case "resource_does_not_support_unallocation":
		return NewResourceDoesNotSupportUnallocationError(theError)
	case "ssh_key_not_found":
		return NewSSHKeyNotFoundError(theError)
	case "security_group_not_found":
		return NewSecurityGroupNotFoundError(theError)
	case "security_group_rule_not_found":
		return NewSecurityGroupRuleNotFoundError(theError)
	case "speed_profile_already_assigned":
		return NewSpeedProfileAlreadyAssignedError(theError)
	case "tag_not_found":
		return NewTagNotFoundError(theError)
	case "task_not_found":
		return NewTaskNotFoundError(theError)
	case "task_queueing_error":
		return NewTaskQueueingError(theError)
	case "trash_object_not_found":
		return NewTrashObjectNotFoundError(theError)
	case "validation_error":
		return NewValidationError(theError)
	case "build_not_found":
		return NewVirtualMachineBuildNotFoundError(theError)
	case "virtual_machine_group_not_found":
		return NewVirtualMachineGroupNotFoundError(theError)
	case "virtual_machine_must_be_started":
		return NewVirtualMachineMustBeStartedError(theError)
	case "virtual_machine_network_interface_not_found":
		return NewVirtualMachineNetworkInterfaceNotFoundError(theError)
	case "virtual_machine_not_found":
		return NewVirtualMachineNotFoundError(theError)
	case "package_not_found":
		return NewVirtualMachinePackageNotFoundError(theError)
	case "zone_not_found":
		return NewZoneNotFoundError(theError)
	default:
		return theError
	}
}
