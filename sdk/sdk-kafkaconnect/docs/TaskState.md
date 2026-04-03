# TaskState

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int32** |  | [optional] 
**Msg** | Pointer to **string** |  | [optional] 
**State** | Pointer to **string** |  | [optional] 
**Trace** | Pointer to **string** |  | [optional] 
**WorkerId** | Pointer to **string** |  | [optional] 

## Methods

### NewTaskState

`func NewTaskState() *TaskState`

NewTaskState instantiates a new TaskState object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTaskStateWithDefaults

`func NewTaskStateWithDefaults() *TaskState`

NewTaskStateWithDefaults instantiates a new TaskState object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *TaskState) GetId() int32`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *TaskState) GetIdOk() (*int32, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *TaskState) SetId(v int32)`

SetId sets Id field to given value.

### HasId

`func (o *TaskState) HasId() bool`

HasId returns a boolean if a field has been set.

### GetMsg

`func (o *TaskState) GetMsg() string`

GetMsg returns the Msg field if non-nil, zero value otherwise.

### GetMsgOk

`func (o *TaskState) GetMsgOk() (*string, bool)`

GetMsgOk returns a tuple with the Msg field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMsg

`func (o *TaskState) SetMsg(v string)`

SetMsg sets Msg field to given value.

### HasMsg

`func (o *TaskState) HasMsg() bool`

HasMsg returns a boolean if a field has been set.

### GetState

`func (o *TaskState) GetState() string`

GetState returns the State field if non-nil, zero value otherwise.

### GetStateOk

`func (o *TaskState) GetStateOk() (*string, bool)`

GetStateOk returns a tuple with the State field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetState

`func (o *TaskState) SetState(v string)`

SetState sets State field to given value.

### HasState

`func (o *TaskState) HasState() bool`

HasState returns a boolean if a field has been set.

### GetTrace

`func (o *TaskState) GetTrace() string`

GetTrace returns the Trace field if non-nil, zero value otherwise.

### GetTraceOk

`func (o *TaskState) GetTraceOk() (*string, bool)`

GetTraceOk returns a tuple with the Trace field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTrace

`func (o *TaskState) SetTrace(v string)`

SetTrace sets Trace field to given value.

### HasTrace

`func (o *TaskState) HasTrace() bool`

HasTrace returns a boolean if a field has been set.

### GetWorkerId

`func (o *TaskState) GetWorkerId() string`

GetWorkerId returns the WorkerId field if non-nil, zero value otherwise.

### GetWorkerIdOk

`func (o *TaskState) GetWorkerIdOk() (*string, bool)`

GetWorkerIdOk returns a tuple with the WorkerId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkerId

`func (o *TaskState) SetWorkerId(v string)`

SetWorkerId sets WorkerId field to given value.

### HasWorkerId

`func (o *TaskState) HasWorkerId() bool`

HasWorkerId returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


