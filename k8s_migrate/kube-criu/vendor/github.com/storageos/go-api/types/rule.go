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
package types

import "context"

// Rule is used to define a rule
type Rule struct {

	// Rule unique ID.
	// Read Only: true
	ID string `json:"id"`

	// Rule name.
	// Required: true
	Name string `json:"name"`

	// Namespace is the object name and authentication scope, such as for teams and projects.
	Namespace string `json:"namespace"`

	// Rule description.
	Description string `json:"description"`

	// Flag describing whether the rule is active.
	// Default: false
	Active bool `json:"active"`

	// Weight is used to determine order during rule processing.  Rules with
	// heavier weights are processed later.
	// default: 0
	Weight int `json:"weight"`

	// RuleAction controls whether the action is to add or remove a label from the
	// matching object(s).
	RuleAction string `json:"action"`

	// Selectors defines the list of labels that should trigger a rule.
	Selector string `json:"selector"`

	// Labels define the list of labels that will be added or removed from the
	// matching object(s).
	Labels map[string]string `json:"labels"`
}

// Rules is a collection of Rules.
type Rules []*Rule

// RuleCreateOptions are available parameters for creating new rules.
type RuleCreateOptions struct {

	// Rule name.
	// Required: true
	Name string `json:"name"`

	// Namespace is the object name and authentication scope, such as for teams and projects.
	Namespace string `json:"namespace"`

	// Rule description.
	Description string `json:"description"`

	// Flag describing whether the rule is active.
	// Default: false
	Active bool `json:"active"`

	// Weight is used to determine order during rule processing.  Rules with
	// heavier weights are processed later.
	// default: 0
	Weight int `json:"weight"`

	// RuleAction controls whether the action is to add or remove a label from the
	// matching object(s).
	RuleAction string `json:"action"`

	// Selectors defines the list of labels that should trigger a rule.
	Selector string `json:"selector"`

	// Labels define the list of labels that will be added or removed from the
	// matching object(s).
	Labels map[string]string `json:"labels"`

	// Context can be set with a timeout or can be used to cancel a request.
	Context context.Context `json:"-"`
}

// RuleUpdateOptions are available parameters for creating new rules.
type RuleUpdateOptions struct {

	// Rule unique ID.
	// Read Only: true
	ID string `json:"id"`

	// Rule name.
	// Required: true
	Name string `json:"name"`

	// Namespace is the object name and authentication scope, such as for teams and projects.
	Namespace string `json:"namespace"`

	// Rule description.
	Description string `json:"description"`

	// Flag describing whether the rule is active.
	// Default: false
	Active bool `json:"active"`

	// Weight is used to determine order during rule processing.  Rules with
	// heavier weights are processed later.
	// default: 0
	Weight int `json:"weight"`

	// Operator is used to compare objects or labels.
	Operator string `json:"operator"`

	// RuleAction controls whether the action is to add or remove a label from the
	// matching object(s).
	RuleAction string `json:"action"`

	// Selectors defines the list of labels that should trigger a rule.
	Selector string `json:"selector"`

	// Labels define the list of labels that will be added or removed from the
	// matching object(s).
	Labels map[string]string `json:"labels"`

	// Context can be set with a timeout or can be used to cancel a request.
	Context context.Context `json:"-"`
}
