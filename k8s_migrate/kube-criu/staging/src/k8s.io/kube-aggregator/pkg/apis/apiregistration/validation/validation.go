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
/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package validation

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/api/validation/path"
	utilvalidation "k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"k8s.io/kube-aggregator/pkg/apis/apiregistration"
)

func ValidateAPIService(apiService *apiregistration.APIService) field.ErrorList {
	requiredName := apiService.Spec.Version + "." + apiService.Spec.Group

	allErrs := validation.ValidateObjectMeta(&apiService.ObjectMeta, false,
		func(name string, prefix bool) []string {
			if minimalFailures := path.IsValidPathSegmentName(name); len(minimalFailures) > 0 {
				return minimalFailures
			}
			// the name *must* be version.group
			if name != requiredName {
				return []string{fmt.Sprintf("must be `spec.version+\".\"+spec.group`: %q", requiredName)}
			}

			return []string{}
		},
		field.NewPath("metadata"))

	// in this case we allow empty group
	if len(apiService.Spec.Group) == 0 && apiService.Spec.Version != "v1" {
		allErrs = append(allErrs, field.Required(field.NewPath("spec", "group"), "only v1 may have an empty group and it better be legacy kube"))
	}
	if len(apiService.Spec.Group) > 0 {
		for _, errString := range utilvalidation.IsDNS1123Subdomain(apiService.Spec.Group) {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "group"), apiService.Spec.Group, errString))
		}
	}

	for _, errString := range utilvalidation.IsDNS1035Label(apiService.Spec.Version) {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "version"), apiService.Spec.Version, errString))
	}

	if apiService.Spec.GroupPriorityMinimum <= 0 || apiService.Spec.GroupPriorityMinimum > 20000 {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "groupPriorityMinimum"), apiService.Spec.GroupPriorityMinimum, "must be positive and less than 20000"))
	}
	if apiService.Spec.VersionPriority <= 0 || apiService.Spec.VersionPriority > 1000 {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "versionPriority"), apiService.Spec.VersionPriority, "must be positive and less than 1000"))
	}

	if apiService.Spec.Service == nil {
		if len(apiService.Spec.CABundle) != 0 {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "caBundle"), fmt.Sprintf("%d bytes", len(apiService.Spec.CABundle)), "local APIServices may not have a caBundle"))
		}
		if apiService.Spec.InsecureSkipTLSVerify {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "insecureSkipTLSVerify"), apiService.Spec.InsecureSkipTLSVerify, "local APIServices may not have insecureSkipTLSVerify"))
		}
		return allErrs
	}

	if len(apiService.Spec.Service.Namespace) == 0 {
		allErrs = append(allErrs, field.Required(field.NewPath("spec", "service", "namespace"), ""))
	}
	if len(apiService.Spec.Service.Name) == 0 {
		allErrs = append(allErrs, field.Required(field.NewPath("spec", "service", "name"), ""))
	}
	if apiService.Spec.InsecureSkipTLSVerify && len(apiService.Spec.CABundle) > 0 {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "insecureSkipTLSVerify"), apiService.Spec.InsecureSkipTLSVerify, "may not be true if caBundle is present"))
	}

	return allErrs
}

func ValidateAPIServiceUpdate(newAPIService *apiregistration.APIService, oldAPIService *apiregistration.APIService) field.ErrorList {
	allErrs := validation.ValidateObjectMetaUpdate(&newAPIService.ObjectMeta, &oldAPIService.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidateAPIService(newAPIService)...)

	return allErrs
}

func ValidateAPIServiceStatus(status *apiregistration.APIServiceStatus, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	for i, condition := range status.Conditions {
		if condition.Status != apiregistration.ConditionTrue &&
			condition.Status != apiregistration.ConditionFalse &&
			condition.Status != apiregistration.ConditionUnknown {
			allErrs = append(allErrs, field.NotSupported(fldPath.Child("conditions").Index(i).Child("status"), condition.Status, []string{
				string(apiregistration.ConditionTrue), string(apiregistration.ConditionFalse), string(apiregistration.ConditionUnknown)}))
		}
	}

	return allErrs
}

func ValidateAPIServiceStatusUpdate(newAPIService *apiregistration.APIService, oldAPIService *apiregistration.APIService) field.ErrorList {
	allErrs := validation.ValidateObjectMetaUpdate(&newAPIService.ObjectMeta, &oldAPIService.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidateAPIServiceStatus(&newAPIService.Status, field.NewPath("status"))...)
	return allErrs
}
