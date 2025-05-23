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

// ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec struct for ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec
type ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec struct {
	// AdminServiceAccount is the service account to be impersonated, or leave it empty to call the API without impersonation.
	AdminServiceAccount *string `json:"adminServiceAccount,omitempty"`
	// ClusterName is the GKE cluster name.
	ClusterName *string `json:"clusterName,omitempty"`
	// Deprecated
	InitialNodeCount *int32 `json:"initialNodeCount,omitempty"`
	Location string `json:"location"`
	// Deprecated
	MachineType *string `json:"machineType,omitempty"`
	// Deprecated
	MaxNodeCount *int32 `json:"maxNodeCount,omitempty"`
	// Project is the Google project containing the GKE cluster.
	Project *string `json:"project,omitempty"`
}

// NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec(location string) *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec {
	this := ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec{}
	this.Location = location
	return &this
}

// NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpecWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpecWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec {
	this := ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec{}
	return &this
}

// GetAdminServiceAccount returns the AdminServiceAccount field value if set, zero value otherwise.
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) GetAdminServiceAccount() string {
	if o == nil || o.AdminServiceAccount == nil {
		var ret string
		return ret
	}
	return *o.AdminServiceAccount
}

// GetAdminServiceAccountOk returns a tuple with the AdminServiceAccount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) GetAdminServiceAccountOk() (*string, bool) {
	if o == nil || o.AdminServiceAccount == nil {
		return nil, false
	}
	return o.AdminServiceAccount, true
}

// HasAdminServiceAccount returns a boolean if a field has been set.
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) HasAdminServiceAccount() bool {
	if o != nil && o.AdminServiceAccount != nil {
		return true
	}

	return false
}

// SetAdminServiceAccount gets a reference to the given string and assigns it to the AdminServiceAccount field.
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) SetAdminServiceAccount(v string) {
	o.AdminServiceAccount = &v
}

// GetClusterName returns the ClusterName field value if set, zero value otherwise.
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) GetClusterName() string {
	if o == nil || o.ClusterName == nil {
		var ret string
		return ret
	}
	return *o.ClusterName
}

// GetClusterNameOk returns a tuple with the ClusterName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) GetClusterNameOk() (*string, bool) {
	if o == nil || o.ClusterName == nil {
		return nil, false
	}
	return o.ClusterName, true
}

// HasClusterName returns a boolean if a field has been set.
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) HasClusterName() bool {
	if o != nil && o.ClusterName != nil {
		return true
	}

	return false
}

// SetClusterName gets a reference to the given string and assigns it to the ClusterName field.
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) SetClusterName(v string) {
	o.ClusterName = &v
}

// GetInitialNodeCount returns the InitialNodeCount field value if set, zero value otherwise.
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) GetInitialNodeCount() int32 {
	if o == nil || o.InitialNodeCount == nil {
		var ret int32
		return ret
	}
	return *o.InitialNodeCount
}

// GetInitialNodeCountOk returns a tuple with the InitialNodeCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) GetInitialNodeCountOk() (*int32, bool) {
	if o == nil || o.InitialNodeCount == nil {
		return nil, false
	}
	return o.InitialNodeCount, true
}

// HasInitialNodeCount returns a boolean if a field has been set.
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) HasInitialNodeCount() bool {
	if o != nil && o.InitialNodeCount != nil {
		return true
	}

	return false
}

// SetInitialNodeCount gets a reference to the given int32 and assigns it to the InitialNodeCount field.
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) SetInitialNodeCount(v int32) {
	o.InitialNodeCount = &v
}

// GetLocation returns the Location field value
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) GetLocation() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Location
}

// GetLocationOk returns a tuple with the Location field value
// and a boolean to check if the value has been set.
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) GetLocationOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Location, true
}

// SetLocation sets field value
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) SetLocation(v string) {
	o.Location = v
}

// GetMachineType returns the MachineType field value if set, zero value otherwise.
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) GetMachineType() string {
	if o == nil || o.MachineType == nil {
		var ret string
		return ret
	}
	return *o.MachineType
}

// GetMachineTypeOk returns a tuple with the MachineType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) GetMachineTypeOk() (*string, bool) {
	if o == nil || o.MachineType == nil {
		return nil, false
	}
	return o.MachineType, true
}

// HasMachineType returns a boolean if a field has been set.
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) HasMachineType() bool {
	if o != nil && o.MachineType != nil {
		return true
	}

	return false
}

// SetMachineType gets a reference to the given string and assigns it to the MachineType field.
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) SetMachineType(v string) {
	o.MachineType = &v
}

// GetMaxNodeCount returns the MaxNodeCount field value if set, zero value otherwise.
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) GetMaxNodeCount() int32 {
	if o == nil || o.MaxNodeCount == nil {
		var ret int32
		return ret
	}
	return *o.MaxNodeCount
}

// GetMaxNodeCountOk returns a tuple with the MaxNodeCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) GetMaxNodeCountOk() (*int32, bool) {
	if o == nil || o.MaxNodeCount == nil {
		return nil, false
	}
	return o.MaxNodeCount, true
}

// HasMaxNodeCount returns a boolean if a field has been set.
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) HasMaxNodeCount() bool {
	if o != nil && o.MaxNodeCount != nil {
		return true
	}

	return false
}

// SetMaxNodeCount gets a reference to the given int32 and assigns it to the MaxNodeCount field.
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) SetMaxNodeCount(v int32) {
	o.MaxNodeCount = &v
}

// GetProject returns the Project field value if set, zero value otherwise.
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) GetProject() string {
	if o == nil || o.Project == nil {
		var ret string
		return ret
	}
	return *o.Project
}

// GetProjectOk returns a tuple with the Project field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) GetProjectOk() (*string, bool) {
	if o == nil || o.Project == nil {
		return nil, false
	}
	return o.Project, true
}

// HasProject returns a boolean if a field has been set.
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) HasProject() bool {
	if o != nil && o.Project != nil {
		return true
	}

	return false
}

// SetProject gets a reference to the given string and assigns it to the Project field.
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) SetProject(v string) {
	o.Project = &v
}

func (o ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.AdminServiceAccount != nil {
		toSerialize["adminServiceAccount"] = o.AdminServiceAccount
	}
	if o.ClusterName != nil {
		toSerialize["clusterName"] = o.ClusterName
	}
	if o.InitialNodeCount != nil {
		toSerialize["initialNodeCount"] = o.InitialNodeCount
	}
	if true {
		toSerialize["location"] = o.Location
	}
	if o.MachineType != nil {
		toSerialize["machineType"] = o.MachineType
	}
	if o.MaxNodeCount != nil {
		toSerialize["maxNodeCount"] = o.MaxNodeCount
	}
	if o.Project != nil {
		toSerialize["project"] = o.Project
	}
	return json.Marshal(toSerialize)
}

type NullableComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec struct {
	value *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec
	isSet bool
}

func (v NullableComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) Get() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec {
	return v.value
}

func (v *NullableComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) Set(val *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) {
	v.value = val
	v.isSet = true
}

func (v NullableComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) IsSet() bool {
	return v.isSet
}

func (v *NullableComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec(val *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) *NullableComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec {
	return &NullableComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec{value: val, isSet: true}
}

func (v NullableComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


