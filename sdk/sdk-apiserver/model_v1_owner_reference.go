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

// V1OwnerReference OwnerReference contains enough information to let you identify an owning object. An owning object must be in the same namespace as the dependent, or be cluster-scoped, so there is no namespace field.
type V1OwnerReference struct {
	// API version of the referent.
	ApiVersion string `json:"apiVersion"`
	// If true, AND if the owner has the \"foregroundDeletion\" finalizer, then the owner cannot be deleted from the key-value store until this reference is removed. See https://kubernetes.io/docs/concepts/architecture/garbage-collection/#foreground-deletion for how the garbage collector interacts with this field and enforces the foreground deletion. Defaults to false. To set this field, a user needs \"delete\" permission of the owner, otherwise 422 (Unprocessable Entity) will be returned.
	BlockOwnerDeletion *bool `json:"blockOwnerDeletion,omitempty"`
	// If true, this reference points to the managing controller.
	Controller *bool `json:"controller,omitempty"`
	// Kind of the referent. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	Kind string `json:"kind"`
	// Name of the referent. More info: http://kubernetes.io/docs/user-guide/identifiers#names
	Name string `json:"name"`
	// UID of the referent. More info: http://kubernetes.io/docs/user-guide/identifiers#uids
	Uid string `json:"uid"`
}

// NewV1OwnerReference instantiates a new V1OwnerReference object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewV1OwnerReference(apiVersion string, kind string, name string, uid string) *V1OwnerReference {
	this := V1OwnerReference{}
	this.ApiVersion = apiVersion
	this.Kind = kind
	this.Name = name
	this.Uid = uid
	return &this
}

// NewV1OwnerReferenceWithDefaults instantiates a new V1OwnerReference object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewV1OwnerReferenceWithDefaults() *V1OwnerReference {
	this := V1OwnerReference{}
	return &this
}

// GetApiVersion returns the ApiVersion field value
func (o *V1OwnerReference) GetApiVersion() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ApiVersion
}

// GetApiVersionOk returns a tuple with the ApiVersion field value
// and a boolean to check if the value has been set.
func (o *V1OwnerReference) GetApiVersionOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ApiVersion, true
}

// SetApiVersion sets field value
func (o *V1OwnerReference) SetApiVersion(v string) {
	o.ApiVersion = v
}

// GetBlockOwnerDeletion returns the BlockOwnerDeletion field value if set, zero value otherwise.
func (o *V1OwnerReference) GetBlockOwnerDeletion() bool {
	if o == nil || o.BlockOwnerDeletion == nil {
		var ret bool
		return ret
	}
	return *o.BlockOwnerDeletion
}

// GetBlockOwnerDeletionOk returns a tuple with the BlockOwnerDeletion field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1OwnerReference) GetBlockOwnerDeletionOk() (*bool, bool) {
	if o == nil || o.BlockOwnerDeletion == nil {
		return nil, false
	}
	return o.BlockOwnerDeletion, true
}

// HasBlockOwnerDeletion returns a boolean if a field has been set.
func (o *V1OwnerReference) HasBlockOwnerDeletion() bool {
	if o != nil && o.BlockOwnerDeletion != nil {
		return true
	}

	return false
}

// SetBlockOwnerDeletion gets a reference to the given bool and assigns it to the BlockOwnerDeletion field.
func (o *V1OwnerReference) SetBlockOwnerDeletion(v bool) {
	o.BlockOwnerDeletion = &v
}

// GetController returns the Controller field value if set, zero value otherwise.
func (o *V1OwnerReference) GetController() bool {
	if o == nil || o.Controller == nil {
		var ret bool
		return ret
	}
	return *o.Controller
}

// GetControllerOk returns a tuple with the Controller field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1OwnerReference) GetControllerOk() (*bool, bool) {
	if o == nil || o.Controller == nil {
		return nil, false
	}
	return o.Controller, true
}

// HasController returns a boolean if a field has been set.
func (o *V1OwnerReference) HasController() bool {
	if o != nil && o.Controller != nil {
		return true
	}

	return false
}

// SetController gets a reference to the given bool and assigns it to the Controller field.
func (o *V1OwnerReference) SetController(v bool) {
	o.Controller = &v
}

// GetKind returns the Kind field value
func (o *V1OwnerReference) GetKind() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Kind
}

// GetKindOk returns a tuple with the Kind field value
// and a boolean to check if the value has been set.
func (o *V1OwnerReference) GetKindOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Kind, true
}

// SetKind sets field value
func (o *V1OwnerReference) SetKind(v string) {
	o.Kind = v
}

// GetName returns the Name field value
func (o *V1OwnerReference) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *V1OwnerReference) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *V1OwnerReference) SetName(v string) {
	o.Name = v
}

// GetUid returns the Uid field value
func (o *V1OwnerReference) GetUid() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Uid
}

// GetUidOk returns a tuple with the Uid field value
// and a boolean to check if the value has been set.
func (o *V1OwnerReference) GetUidOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Uid, true
}

// SetUid sets field value
func (o *V1OwnerReference) SetUid(v string) {
	o.Uid = v
}

func (o V1OwnerReference) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["apiVersion"] = o.ApiVersion
	}
	if o.BlockOwnerDeletion != nil {
		toSerialize["blockOwnerDeletion"] = o.BlockOwnerDeletion
	}
	if o.Controller != nil {
		toSerialize["controller"] = o.Controller
	}
	if true {
		toSerialize["kind"] = o.Kind
	}
	if true {
		toSerialize["name"] = o.Name
	}
	if true {
		toSerialize["uid"] = o.Uid
	}
	return json.Marshal(toSerialize)
}

type NullableV1OwnerReference struct {
	value *V1OwnerReference
	isSet bool
}

func (v NullableV1OwnerReference) Get() *V1OwnerReference {
	return v.value
}

func (v *NullableV1OwnerReference) Set(val *V1OwnerReference) {
	v.value = val
	v.isSet = true
}

func (v NullableV1OwnerReference) IsSet() bool {
	return v.isSet
}

func (v *NullableV1OwnerReference) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableV1OwnerReference(val *V1OwnerReference) *NullableV1OwnerReference {
	return &NullableV1OwnerReference{value: val, isSet: true}
}

func (v NullableV1OwnerReference) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableV1OwnerReference) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


