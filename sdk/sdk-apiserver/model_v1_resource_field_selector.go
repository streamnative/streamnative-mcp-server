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

// V1ResourceFieldSelector ResourceFieldSelector represents container resources (cpu, memory) and their output format
type V1ResourceFieldSelector struct {
	// Container name: required for volumes, optional for env vars
	ContainerName *string `json:"containerName,omitempty"`
	// Specifies the output format of the exposed resources, defaults to \"1\"
	Divisor *string `json:"divisor,omitempty"`
	// Required: resource to select
	Resource string `json:"resource"`
}

// NewV1ResourceFieldSelector instantiates a new V1ResourceFieldSelector object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewV1ResourceFieldSelector(resource string) *V1ResourceFieldSelector {
	this := V1ResourceFieldSelector{}
	this.Resource = resource
	return &this
}

// NewV1ResourceFieldSelectorWithDefaults instantiates a new V1ResourceFieldSelector object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewV1ResourceFieldSelectorWithDefaults() *V1ResourceFieldSelector {
	this := V1ResourceFieldSelector{}
	return &this
}

// GetContainerName returns the ContainerName field value if set, zero value otherwise.
func (o *V1ResourceFieldSelector) GetContainerName() string {
	if o == nil || o.ContainerName == nil {
		var ret string
		return ret
	}
	return *o.ContainerName
}

// GetContainerNameOk returns a tuple with the ContainerName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1ResourceFieldSelector) GetContainerNameOk() (*string, bool) {
	if o == nil || o.ContainerName == nil {
		return nil, false
	}
	return o.ContainerName, true
}

// HasContainerName returns a boolean if a field has been set.
func (o *V1ResourceFieldSelector) HasContainerName() bool {
	if o != nil && o.ContainerName != nil {
		return true
	}

	return false
}

// SetContainerName gets a reference to the given string and assigns it to the ContainerName field.
func (o *V1ResourceFieldSelector) SetContainerName(v string) {
	o.ContainerName = &v
}

// GetDivisor returns the Divisor field value if set, zero value otherwise.
func (o *V1ResourceFieldSelector) GetDivisor() string {
	if o == nil || o.Divisor == nil {
		var ret string
		return ret
	}
	return *o.Divisor
}

// GetDivisorOk returns a tuple with the Divisor field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1ResourceFieldSelector) GetDivisorOk() (*string, bool) {
	if o == nil || o.Divisor == nil {
		return nil, false
	}
	return o.Divisor, true
}

// HasDivisor returns a boolean if a field has been set.
func (o *V1ResourceFieldSelector) HasDivisor() bool {
	if o != nil && o.Divisor != nil {
		return true
	}

	return false
}

// SetDivisor gets a reference to the given string and assigns it to the Divisor field.
func (o *V1ResourceFieldSelector) SetDivisor(v string) {
	o.Divisor = &v
}

// GetResource returns the Resource field value
func (o *V1ResourceFieldSelector) GetResource() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Resource
}

// GetResourceOk returns a tuple with the Resource field value
// and a boolean to check if the value has been set.
func (o *V1ResourceFieldSelector) GetResourceOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Resource, true
}

// SetResource sets field value
func (o *V1ResourceFieldSelector) SetResource(v string) {
	o.Resource = v
}

func (o V1ResourceFieldSelector) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.ContainerName != nil {
		toSerialize["containerName"] = o.ContainerName
	}
	if o.Divisor != nil {
		toSerialize["divisor"] = o.Divisor
	}
	if true {
		toSerialize["resource"] = o.Resource
	}
	return json.Marshal(toSerialize)
}

type NullableV1ResourceFieldSelector struct {
	value *V1ResourceFieldSelector
	isSet bool
}

func (v NullableV1ResourceFieldSelector) Get() *V1ResourceFieldSelector {
	return v.value
}

func (v *NullableV1ResourceFieldSelector) Set(val *V1ResourceFieldSelector) {
	v.value = val
	v.isSet = true
}

func (v NullableV1ResourceFieldSelector) IsSet() bool {
	return v.isSet
}

func (v *NullableV1ResourceFieldSelector) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableV1ResourceFieldSelector(val *V1ResourceFieldSelector) *NullableV1ResourceFieldSelector {
	return &NullableV1ResourceFieldSelector{value: val, isSet: true}
}

func (v NullableV1ResourceFieldSelector) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableV1ResourceFieldSelector) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


