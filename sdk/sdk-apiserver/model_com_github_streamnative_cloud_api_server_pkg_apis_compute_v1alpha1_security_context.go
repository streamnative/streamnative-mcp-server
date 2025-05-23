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

// ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext struct for ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext
type ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext struct {
	FsGroup *int64 `json:"fsGroup,omitempty"`
	// ReadOnlyRootFilesystem specifies whether the container use a read-only filesystem.
	ReadOnlyRootFilesystem *bool `json:"readOnlyRootFilesystem,omitempty"`
	RunAsGroup *int64 `json:"runAsGroup,omitempty"`
	RunAsNonRoot *bool `json:"runAsNonRoot,omitempty"`
	RunAsUser *int64 `json:"runAsUser,omitempty"`
}

// NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext instantiates a new ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext() *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext {
	this := ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext{}
	return &this
}

// NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContextWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContextWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext {
	this := ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext{}
	return &this
}

// GetFsGroup returns the FsGroup field value if set, zero value otherwise.
func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext) GetFsGroup() int64 {
	if o == nil || o.FsGroup == nil {
		var ret int64
		return ret
	}
	return *o.FsGroup
}

// GetFsGroupOk returns a tuple with the FsGroup field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext) GetFsGroupOk() (*int64, bool) {
	if o == nil || o.FsGroup == nil {
		return nil, false
	}
	return o.FsGroup, true
}

// HasFsGroup returns a boolean if a field has been set.
func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext) HasFsGroup() bool {
	if o != nil && o.FsGroup != nil {
		return true
	}

	return false
}

// SetFsGroup gets a reference to the given int64 and assigns it to the FsGroup field.
func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext) SetFsGroup(v int64) {
	o.FsGroup = &v
}

// GetReadOnlyRootFilesystem returns the ReadOnlyRootFilesystem field value if set, zero value otherwise.
func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext) GetReadOnlyRootFilesystem() bool {
	if o == nil || o.ReadOnlyRootFilesystem == nil {
		var ret bool
		return ret
	}
	return *o.ReadOnlyRootFilesystem
}

// GetReadOnlyRootFilesystemOk returns a tuple with the ReadOnlyRootFilesystem field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext) GetReadOnlyRootFilesystemOk() (*bool, bool) {
	if o == nil || o.ReadOnlyRootFilesystem == nil {
		return nil, false
	}
	return o.ReadOnlyRootFilesystem, true
}

// HasReadOnlyRootFilesystem returns a boolean if a field has been set.
func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext) HasReadOnlyRootFilesystem() bool {
	if o != nil && o.ReadOnlyRootFilesystem != nil {
		return true
	}

	return false
}

// SetReadOnlyRootFilesystem gets a reference to the given bool and assigns it to the ReadOnlyRootFilesystem field.
func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext) SetReadOnlyRootFilesystem(v bool) {
	o.ReadOnlyRootFilesystem = &v
}

// GetRunAsGroup returns the RunAsGroup field value if set, zero value otherwise.
func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext) GetRunAsGroup() int64 {
	if o == nil || o.RunAsGroup == nil {
		var ret int64
		return ret
	}
	return *o.RunAsGroup
}

// GetRunAsGroupOk returns a tuple with the RunAsGroup field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext) GetRunAsGroupOk() (*int64, bool) {
	if o == nil || o.RunAsGroup == nil {
		return nil, false
	}
	return o.RunAsGroup, true
}

// HasRunAsGroup returns a boolean if a field has been set.
func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext) HasRunAsGroup() bool {
	if o != nil && o.RunAsGroup != nil {
		return true
	}

	return false
}

// SetRunAsGroup gets a reference to the given int64 and assigns it to the RunAsGroup field.
func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext) SetRunAsGroup(v int64) {
	o.RunAsGroup = &v
}

// GetRunAsNonRoot returns the RunAsNonRoot field value if set, zero value otherwise.
func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext) GetRunAsNonRoot() bool {
	if o == nil || o.RunAsNonRoot == nil {
		var ret bool
		return ret
	}
	return *o.RunAsNonRoot
}

// GetRunAsNonRootOk returns a tuple with the RunAsNonRoot field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext) GetRunAsNonRootOk() (*bool, bool) {
	if o == nil || o.RunAsNonRoot == nil {
		return nil, false
	}
	return o.RunAsNonRoot, true
}

// HasRunAsNonRoot returns a boolean if a field has been set.
func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext) HasRunAsNonRoot() bool {
	if o != nil && o.RunAsNonRoot != nil {
		return true
	}

	return false
}

// SetRunAsNonRoot gets a reference to the given bool and assigns it to the RunAsNonRoot field.
func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext) SetRunAsNonRoot(v bool) {
	o.RunAsNonRoot = &v
}

// GetRunAsUser returns the RunAsUser field value if set, zero value otherwise.
func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext) GetRunAsUser() int64 {
	if o == nil || o.RunAsUser == nil {
		var ret int64
		return ret
	}
	return *o.RunAsUser
}

// GetRunAsUserOk returns a tuple with the RunAsUser field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext) GetRunAsUserOk() (*int64, bool) {
	if o == nil || o.RunAsUser == nil {
		return nil, false
	}
	return o.RunAsUser, true
}

// HasRunAsUser returns a boolean if a field has been set.
func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext) HasRunAsUser() bool {
	if o != nil && o.RunAsUser != nil {
		return true
	}

	return false
}

// SetRunAsUser gets a reference to the given int64 and assigns it to the RunAsUser field.
func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext) SetRunAsUser(v int64) {
	o.RunAsUser = &v
}

func (o ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.FsGroup != nil {
		toSerialize["fsGroup"] = o.FsGroup
	}
	if o.ReadOnlyRootFilesystem != nil {
		toSerialize["readOnlyRootFilesystem"] = o.ReadOnlyRootFilesystem
	}
	if o.RunAsGroup != nil {
		toSerialize["runAsGroup"] = o.RunAsGroup
	}
	if o.RunAsNonRoot != nil {
		toSerialize["runAsNonRoot"] = o.RunAsNonRoot
	}
	if o.RunAsUser != nil {
		toSerialize["runAsUser"] = o.RunAsUser
	}
	return json.Marshal(toSerialize)
}

type NullableComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext struct {
	value *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext
	isSet bool
}

func (v NullableComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext) Get() *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext {
	return v.value
}

func (v *NullableComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext) Set(val *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext) {
	v.value = val
	v.isSet = true
}

func (v NullableComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext) IsSet() bool {
	return v.isSet
}

func (v *NullableComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext(val *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext) *NullableComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext {
	return &NullableComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext{value: val, isSet: true}
}

func (v NullableComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


