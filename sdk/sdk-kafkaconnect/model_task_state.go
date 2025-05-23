/*
Kafka Connect With Pulsar API

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.0.1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package kafkaconnect

import (
	"encoding/json"
)

// checks if the TaskState type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &TaskState{}

// TaskState struct for TaskState
type TaskState struct {
	Id *int32 `json:"id,omitempty"`
	Msg *string `json:"msg,omitempty"`
	State *string `json:"state,omitempty"`
	Trace *string `json:"trace,omitempty"`
	WorkerId *string `json:"worker_id,omitempty"`
}

// NewTaskState instantiates a new TaskState object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTaskState() *TaskState {
	this := TaskState{}
	return &this
}

// NewTaskStateWithDefaults instantiates a new TaskState object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTaskStateWithDefaults() *TaskState {
	this := TaskState{}
	return &this
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *TaskState) GetId() int32 {
	if o == nil || IsNil(o.Id) {
		var ret int32
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TaskState) GetIdOk() (*int32, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *TaskState) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given int32 and assigns it to the Id field.
func (o *TaskState) SetId(v int32) {
	o.Id = &v
}

// GetMsg returns the Msg field value if set, zero value otherwise.
func (o *TaskState) GetMsg() string {
	if o == nil || IsNil(o.Msg) {
		var ret string
		return ret
	}
	return *o.Msg
}

// GetMsgOk returns a tuple with the Msg field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TaskState) GetMsgOk() (*string, bool) {
	if o == nil || IsNil(o.Msg) {
		return nil, false
	}
	return o.Msg, true
}

// HasMsg returns a boolean if a field has been set.
func (o *TaskState) HasMsg() bool {
	if o != nil && !IsNil(o.Msg) {
		return true
	}

	return false
}

// SetMsg gets a reference to the given string and assigns it to the Msg field.
func (o *TaskState) SetMsg(v string) {
	o.Msg = &v
}

// GetState returns the State field value if set, zero value otherwise.
func (o *TaskState) GetState() string {
	if o == nil || IsNil(o.State) {
		var ret string
		return ret
	}
	return *o.State
}

// GetStateOk returns a tuple with the State field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TaskState) GetStateOk() (*string, bool) {
	if o == nil || IsNil(o.State) {
		return nil, false
	}
	return o.State, true
}

// HasState returns a boolean if a field has been set.
func (o *TaskState) HasState() bool {
	if o != nil && !IsNil(o.State) {
		return true
	}

	return false
}

// SetState gets a reference to the given string and assigns it to the State field.
func (o *TaskState) SetState(v string) {
	o.State = &v
}

// GetTrace returns the Trace field value if set, zero value otherwise.
func (o *TaskState) GetTrace() string {
	if o == nil || IsNil(o.Trace) {
		var ret string
		return ret
	}
	return *o.Trace
}

// GetTraceOk returns a tuple with the Trace field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TaskState) GetTraceOk() (*string, bool) {
	if o == nil || IsNil(o.Trace) {
		return nil, false
	}
	return o.Trace, true
}

// HasTrace returns a boolean if a field has been set.
func (o *TaskState) HasTrace() bool {
	if o != nil && !IsNil(o.Trace) {
		return true
	}

	return false
}

// SetTrace gets a reference to the given string and assigns it to the Trace field.
func (o *TaskState) SetTrace(v string) {
	o.Trace = &v
}

// GetWorkerId returns the WorkerId field value if set, zero value otherwise.
func (o *TaskState) GetWorkerId() string {
	if o == nil || IsNil(o.WorkerId) {
		var ret string
		return ret
	}
	return *o.WorkerId
}

// GetWorkerIdOk returns a tuple with the WorkerId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TaskState) GetWorkerIdOk() (*string, bool) {
	if o == nil || IsNil(o.WorkerId) {
		return nil, false
	}
	return o.WorkerId, true
}

// HasWorkerId returns a boolean if a field has been set.
func (o *TaskState) HasWorkerId() bool {
	if o != nil && !IsNil(o.WorkerId) {
		return true
	}

	return false
}

// SetWorkerId gets a reference to the given string and assigns it to the WorkerId field.
func (o *TaskState) SetWorkerId(v string) {
	o.WorkerId = &v
}

func (o TaskState) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o TaskState) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Id) {
		toSerialize["id"] = o.Id
	}
	if !IsNil(o.Msg) {
		toSerialize["msg"] = o.Msg
	}
	if !IsNil(o.State) {
		toSerialize["state"] = o.State
	}
	if !IsNil(o.Trace) {
		toSerialize["trace"] = o.Trace
	}
	if !IsNil(o.WorkerId) {
		toSerialize["worker_id"] = o.WorkerId
	}
	return toSerialize, nil
}

type NullableTaskState struct {
	value *TaskState
	isSet bool
}

func (v NullableTaskState) Get() *TaskState {
	return v.value
}

func (v *NullableTaskState) Set(val *TaskState) {
	v.value = val
	v.isSet = true
}

func (v NullableTaskState) IsSet() bool {
	return v.isSet
}

func (v *NullableTaskState) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTaskState(val *TaskState) *NullableTaskState {
	return &NullableTaskState{value: val, isSet: true}
}

func (v NullableTaskState) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTaskState) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


