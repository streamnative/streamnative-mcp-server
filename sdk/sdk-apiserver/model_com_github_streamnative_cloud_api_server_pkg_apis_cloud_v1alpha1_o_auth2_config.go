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

// ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OAuth2Config OAuth2Config define oauth2 config of PulsarInstance
type ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OAuth2Config struct {
	// TokenLifetimeSeconds access token lifetime (in seconds) for the API. Default value is 86,400 seconds (24 hours). Maximum value is 2,592,000 seconds (30 days) Document link: https://auth0.com/docs/secure/tokens/access-tokens/update-access-token-lifetime
	TokenLifetimeSeconds *int32 `json:"tokenLifetimeSeconds,omitempty"`
}

// NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OAuth2Config instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OAuth2Config object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OAuth2Config() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OAuth2Config {
	this := ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OAuth2Config{}
	return &this
}

// NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OAuth2ConfigWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OAuth2Config object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OAuth2ConfigWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OAuth2Config {
	this := ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OAuth2Config{}
	return &this
}

// GetTokenLifetimeSeconds returns the TokenLifetimeSeconds field value if set, zero value otherwise.
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OAuth2Config) GetTokenLifetimeSeconds() int32 {
	if o == nil || o.TokenLifetimeSeconds == nil {
		var ret int32
		return ret
	}
	return *o.TokenLifetimeSeconds
}

// GetTokenLifetimeSecondsOk returns a tuple with the TokenLifetimeSeconds field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OAuth2Config) GetTokenLifetimeSecondsOk() (*int32, bool) {
	if o == nil || o.TokenLifetimeSeconds == nil {
		return nil, false
	}
	return o.TokenLifetimeSeconds, true
}

// HasTokenLifetimeSeconds returns a boolean if a field has been set.
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OAuth2Config) HasTokenLifetimeSeconds() bool {
	if o != nil && o.TokenLifetimeSeconds != nil {
		return true
	}

	return false
}

// SetTokenLifetimeSeconds gets a reference to the given int32 and assigns it to the TokenLifetimeSeconds field.
func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OAuth2Config) SetTokenLifetimeSeconds(v int32) {
	o.TokenLifetimeSeconds = &v
}

func (o ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OAuth2Config) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.TokenLifetimeSeconds != nil {
		toSerialize["tokenLifetimeSeconds"] = o.TokenLifetimeSeconds
	}
	return json.Marshal(toSerialize)
}

type NullableComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OAuth2Config struct {
	value *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OAuth2Config
	isSet bool
}

func (v NullableComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OAuth2Config) Get() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OAuth2Config {
	return v.value
}

func (v *NullableComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OAuth2Config) Set(val *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OAuth2Config) {
	v.value = val
	v.isSet = true
}

func (v NullableComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OAuth2Config) IsSet() bool {
	return v.isSet
}

func (v *NullableComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OAuth2Config) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OAuth2Config(val *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OAuth2Config) *NullableComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OAuth2Config {
	return &NullableComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OAuth2Config{value: val, isSet: true}
}

func (v NullableComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OAuth2Config) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OAuth2Config) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


