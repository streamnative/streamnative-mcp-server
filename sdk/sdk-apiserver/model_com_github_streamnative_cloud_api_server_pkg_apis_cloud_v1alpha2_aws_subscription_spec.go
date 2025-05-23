// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

/*
Api

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: v0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package sncloud

import (
	"encoding/json"
)

// ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionSpec AWSSubscriptionSpec defines the desired state of AWSSubscription
type ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionSpec struct {
	CustomerID *string `json:"customerID,omitempty"`
	ProductCode *string `json:"productCode,omitempty"`
}

// NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionSpec instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionSpec object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionSpec() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionSpec {
	this := ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionSpec{}
	return &this
}

// NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionSpecWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionSpec object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionSpecWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionSpec {
	this := ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionSpec{}
	return &this
}

// GetCustomerID returns the CustomerID field value if set, zero value otherwise.
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionSpec) GetCustomerID() string {
	if o == nil || o.CustomerID == nil {
		var ret string
		return ret
	}
	return *o.CustomerID
}

// GetCustomerIDOk returns a tuple with the CustomerID field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionSpec) GetCustomerIDOk() (*string, bool) {
	if o == nil || o.CustomerID == nil {
		return nil, false
	}
	return o.CustomerID, true
}

// HasCustomerID returns a boolean if a field has been set.
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionSpec) HasCustomerID() bool {
	if o != nil && o.CustomerID != nil {
		return true
	}

	return false
}

// SetCustomerID gets a reference to the given string and assigns it to the CustomerID field.
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionSpec) SetCustomerID(v string) {
	o.CustomerID = &v
}

// GetProductCode returns the ProductCode field value if set, zero value otherwise.
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionSpec) GetProductCode() string {
	if o == nil || o.ProductCode == nil {
		var ret string
		return ret
	}
	return *o.ProductCode
}

// GetProductCodeOk returns a tuple with the ProductCode field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionSpec) GetProductCodeOk() (*string, bool) {
	if o == nil || o.ProductCode == nil {
		return nil, false
	}
	return o.ProductCode, true
}

// HasProductCode returns a boolean if a field has been set.
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionSpec) HasProductCode() bool {
	if o != nil && o.ProductCode != nil {
		return true
	}

	return false
}

// SetProductCode gets a reference to the given string and assigns it to the ProductCode field.
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionSpec) SetProductCode(v string) {
	o.ProductCode = &v
}

func (o ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionSpec) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.CustomerID != nil {
		toSerialize["customerID"] = o.CustomerID
	}
	if o.ProductCode != nil {
		toSerialize["productCode"] = o.ProductCode
	}
	return json.Marshal(toSerialize)
}

type NullableComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionSpec struct {
	value *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionSpec
	isSet bool
}

func (v NullableComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionSpec) Get() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionSpec {
	return v.value
}

func (v *NullableComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionSpec) Set(val *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionSpec) {
	v.value = val
	v.isSet = true
}

func (v NullableComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionSpec) IsSet() bool {
	return v.isSet
}

func (v *NullableComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionSpec) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionSpec(val *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionSpec) *NullableComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionSpec {
	return &NullableComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionSpec{value: val, isSet: true}
}

func (v NullableComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionSpec) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionSpec) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


