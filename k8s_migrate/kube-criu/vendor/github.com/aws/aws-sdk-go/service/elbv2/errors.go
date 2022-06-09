/*
Copyright (c) 2014-2020 CGCL Labs
Container_Migrate is licensed under Mulan PSL v2.
You can use this software according to the terms and conditions of the Mulan PSL v2.
You may obtain a copy of Mulan PSL v2 at:
         http://license.coscl.org.cn/MulanPSL2
THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
See the Mulan PSL v2 for more details.
*/
// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

package elbv2

const (

	// ErrCodeAllocationIdNotFoundException for service response error code
	// "AllocationIdNotFound".
	//
	// The specified allocation ID does not exist.
	ErrCodeAllocationIdNotFoundException = "AllocationIdNotFound"

	// ErrCodeAvailabilityZoneNotSupportedException for service response error code
	// "AvailabilityZoneNotSupported".
	//
	// The specified Availability Zone is not supported.
	ErrCodeAvailabilityZoneNotSupportedException = "AvailabilityZoneNotSupported"

	// ErrCodeCertificateNotFoundException for service response error code
	// "CertificateNotFound".
	//
	// The specified certificate does not exist.
	ErrCodeCertificateNotFoundException = "CertificateNotFound"

	// ErrCodeDuplicateListenerException for service response error code
	// "DuplicateListener".
	//
	// A listener with the specified port already exists.
	ErrCodeDuplicateListenerException = "DuplicateListener"

	// ErrCodeDuplicateLoadBalancerNameException for service response error code
	// "DuplicateLoadBalancerName".
	//
	// A load balancer with the specified name already exists.
	ErrCodeDuplicateLoadBalancerNameException = "DuplicateLoadBalancerName"

	// ErrCodeDuplicateTagKeysException for service response error code
	// "DuplicateTagKeys".
	//
	// A tag key was specified more than once.
	ErrCodeDuplicateTagKeysException = "DuplicateTagKeys"

	// ErrCodeDuplicateTargetGroupNameException for service response error code
	// "DuplicateTargetGroupName".
	//
	// A target group with the specified name already exists.
	ErrCodeDuplicateTargetGroupNameException = "DuplicateTargetGroupName"

	// ErrCodeHealthUnavailableException for service response error code
	// "HealthUnavailable".
	//
	// The health of the specified targets could not be retrieved due to an internal
	// error.
	ErrCodeHealthUnavailableException = "HealthUnavailable"

	// ErrCodeIncompatibleProtocolsException for service response error code
	// "IncompatibleProtocols".
	//
	// The specified configuration is not valid with this protocol.
	ErrCodeIncompatibleProtocolsException = "IncompatibleProtocols"

	// ErrCodeInvalidConfigurationRequestException for service response error code
	// "InvalidConfigurationRequest".
	//
	// The requested configuration is not valid.
	ErrCodeInvalidConfigurationRequestException = "InvalidConfigurationRequest"

	// ErrCodeInvalidLoadBalancerActionException for service response error code
	// "InvalidLoadBalancerAction".
	//
	// The requested action is not valid.
	ErrCodeInvalidLoadBalancerActionException = "InvalidLoadBalancerAction"

	// ErrCodeInvalidSchemeException for service response error code
	// "InvalidScheme".
	//
	// The requested scheme is not valid.
	ErrCodeInvalidSchemeException = "InvalidScheme"

	// ErrCodeInvalidSecurityGroupException for service response error code
	// "InvalidSecurityGroup".
	//
	// The specified security group does not exist.
	ErrCodeInvalidSecurityGroupException = "InvalidSecurityGroup"

	// ErrCodeInvalidSubnetException for service response error code
	// "InvalidSubnet".
	//
	// The specified subnet is out of available addresses.
	ErrCodeInvalidSubnetException = "InvalidSubnet"

	// ErrCodeInvalidTargetException for service response error code
	// "InvalidTarget".
	//
	// The specified target does not exist, is not in the same VPC as the target
	// group, or has an unsupported instance type.
	ErrCodeInvalidTargetException = "InvalidTarget"

	// ErrCodeListenerNotFoundException for service response error code
	// "ListenerNotFound".
	//
	// The specified listener does not exist.
	ErrCodeListenerNotFoundException = "ListenerNotFound"

	// ErrCodeLoadBalancerNotFoundException for service response error code
	// "LoadBalancerNotFound".
	//
	// The specified load balancer does not exist.
	ErrCodeLoadBalancerNotFoundException = "LoadBalancerNotFound"

	// ErrCodeOperationNotPermittedException for service response error code
	// "OperationNotPermitted".
	//
	// This operation is not allowed.
	ErrCodeOperationNotPermittedException = "OperationNotPermitted"

	// ErrCodePriorityInUseException for service response error code
	// "PriorityInUse".
	//
	// The specified priority is in use.
	ErrCodePriorityInUseException = "PriorityInUse"

	// ErrCodeResourceInUseException for service response error code
	// "ResourceInUse".
	//
	// A specified resource is in use.
	ErrCodeResourceInUseException = "ResourceInUse"

	// ErrCodeRuleNotFoundException for service response error code
	// "RuleNotFound".
	//
	// The specified rule does not exist.
	ErrCodeRuleNotFoundException = "RuleNotFound"

	// ErrCodeSSLPolicyNotFoundException for service response error code
	// "SSLPolicyNotFound".
	//
	// The specified SSL policy does not exist.
	ErrCodeSSLPolicyNotFoundException = "SSLPolicyNotFound"

	// ErrCodeSubnetNotFoundException for service response error code
	// "SubnetNotFound".
	//
	// The specified subnet does not exist.
	ErrCodeSubnetNotFoundException = "SubnetNotFound"

	// ErrCodeTargetGroupAssociationLimitException for service response error code
	// "TargetGroupAssociationLimit".
	//
	// You've reached the limit on the number of load balancers per target group.
	ErrCodeTargetGroupAssociationLimitException = "TargetGroupAssociationLimit"

	// ErrCodeTargetGroupNotFoundException for service response error code
	// "TargetGroupNotFound".
	//
	// The specified target group does not exist.
	ErrCodeTargetGroupNotFoundException = "TargetGroupNotFound"

	// ErrCodeTooManyActionsException for service response error code
	// "TooManyActions".
	//
	// You've reached the limit on the number of actions per rule.
	ErrCodeTooManyActionsException = "TooManyActions"

	// ErrCodeTooManyCertificatesException for service response error code
	// "TooManyCertificates".
	//
	// You've reached the limit on the number of certificates per load balancer.
	ErrCodeTooManyCertificatesException = "TooManyCertificates"

	// ErrCodeTooManyListenersException for service response error code
	// "TooManyListeners".
	//
	// You've reached the limit on the number of listeners per load balancer.
	ErrCodeTooManyListenersException = "TooManyListeners"

	// ErrCodeTooManyLoadBalancersException for service response error code
	// "TooManyLoadBalancers".
	//
	// You've reached the limit on the number of load balancers for your AWS account.
	ErrCodeTooManyLoadBalancersException = "TooManyLoadBalancers"

	// ErrCodeTooManyRegistrationsForTargetIdException for service response error code
	// "TooManyRegistrationsForTargetId".
	//
	// You've reached the limit on the number of times a target can be registered
	// with a load balancer.
	ErrCodeTooManyRegistrationsForTargetIdException = "TooManyRegistrationsForTargetId"

	// ErrCodeTooManyRulesException for service response error code
	// "TooManyRules".
	//
	// You've reached the limit on the number of rules per load balancer.
	ErrCodeTooManyRulesException = "TooManyRules"

	// ErrCodeTooManyTagsException for service response error code
	// "TooManyTags".
	//
	// You've reached the limit on the number of tags per load balancer.
	ErrCodeTooManyTagsException = "TooManyTags"

	// ErrCodeTooManyTargetGroupsException for service response error code
	// "TooManyTargetGroups".
	//
	// You've reached the limit on the number of target groups for your AWS account.
	ErrCodeTooManyTargetGroupsException = "TooManyTargetGroups"

	// ErrCodeTooManyTargetsException for service response error code
	// "TooManyTargets".
	//
	// You've reached the limit on the number of targets.
	ErrCodeTooManyTargetsException = "TooManyTargets"

	// ErrCodeUnsupportedProtocolException for service response error code
	// "UnsupportedProtocol".
	//
	// The specified protocol is not supported.
	ErrCodeUnsupportedProtocolException = "UnsupportedProtocol"
)