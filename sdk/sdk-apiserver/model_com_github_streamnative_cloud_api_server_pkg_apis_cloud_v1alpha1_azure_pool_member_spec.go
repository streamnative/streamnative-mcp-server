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

// ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec struct for ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec
type ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec struct {
	// ClientID is the Azure client ID of which the GSA can impersonate.
	ClientID string `json:"clientID"`
	// ClusterName is the AKS cluster name.
	ClusterName string `json:"clusterName"`
	// Location is the Azure location of the cluster.
	Location string `json:"location"`
	// ResourceGroup is the Azure resource group of the cluster.
	ResourceGroup string `json:"resourceGroup"`
	// SubscriptionID is the Azure subscription ID of the cluster.
	SubscriptionID string `json:"subscriptionID"`
	// TenantID is the Azure tenant ID of the client.
	TenantID string `json:"tenantID"`
}

// NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec(clientID string, clusterName string, location string, resourceGroup string, subscriptionID string, tenantID string) *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec {
	this := ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec{}
	this.ClientID = clientID
	this.ClusterName = clusterName
	this.Location = location
	this.ResourceGroup = resourceGroup
	this.SubscriptionID = subscriptionID
	this.TenantID = tenantID
	return &this
}

// NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpecWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpecWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec {
	this := ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec{}
	return &this
}

// GetClientID returns the ClientID field value
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec) GetClientID() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ClientID
}

// GetClientIDOk returns a tuple with the ClientID field value
// and a boolean to check if the value has been set.
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec) GetClientIDOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ClientID, true
}

// SetClientID sets field value
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec) SetClientID(v string) {
	o.ClientID = v
}

// GetClusterName returns the ClusterName field value
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec) GetClusterName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ClusterName
}

// GetClusterNameOk returns a tuple with the ClusterName field value
// and a boolean to check if the value has been set.
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec) GetClusterNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ClusterName, true
}

// SetClusterName sets field value
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec) SetClusterName(v string) {
	o.ClusterName = v
}

// GetLocation returns the Location field value
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec) GetLocation() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Location
}

// GetLocationOk returns a tuple with the Location field value
// and a boolean to check if the value has been set.
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec) GetLocationOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Location, true
}

// SetLocation sets field value
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec) SetLocation(v string) {
	o.Location = v
}

// GetResourceGroup returns the ResourceGroup field value
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec) GetResourceGroup() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ResourceGroup
}

// GetResourceGroupOk returns a tuple with the ResourceGroup field value
// and a boolean to check if the value has been set.
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec) GetResourceGroupOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ResourceGroup, true
}

// SetResourceGroup sets field value
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec) SetResourceGroup(v string) {
	o.ResourceGroup = v
}

// GetSubscriptionID returns the SubscriptionID field value
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec) GetSubscriptionID() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.SubscriptionID
}

// GetSubscriptionIDOk returns a tuple with the SubscriptionID field value
// and a boolean to check if the value has been set.
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec) GetSubscriptionIDOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.SubscriptionID, true
}

// SetSubscriptionID sets field value
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec) SetSubscriptionID(v string) {
	o.SubscriptionID = v
}

// GetTenantID returns the TenantID field value
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec) GetTenantID() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.TenantID
}

// GetTenantIDOk returns a tuple with the TenantID field value
// and a boolean to check if the value has been set.
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec) GetTenantIDOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.TenantID, true
}

// SetTenantID sets field value
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec) SetTenantID(v string) {
	o.TenantID = v
}

func (o ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["clientID"] = o.ClientID
	}
	if true {
		toSerialize["clusterName"] = o.ClusterName
	}
	if true {
		toSerialize["location"] = o.Location
	}
	if true {
		toSerialize["resourceGroup"] = o.ResourceGroup
	}
	if true {
		toSerialize["subscriptionID"] = o.SubscriptionID
	}
	if true {
		toSerialize["tenantID"] = o.TenantID
	}
	return json.Marshal(toSerialize)
}

type NullableComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec struct {
	value *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec
	isSet bool
}

func (v NullableComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec) Get() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec {
	return v.value
}

func (v *NullableComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec) Set(val *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec) {
	v.value = val
	v.isSet = true
}

func (v NullableComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec) IsSet() bool {
	return v.isSet
}

func (v *NullableComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec(val *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec) *NullableComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec {
	return &NullableComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec{value: val, isSet: true}
}

func (v NullableComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


