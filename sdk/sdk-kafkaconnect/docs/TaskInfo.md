# TaskInfo

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Config** | Pointer to **map[string]string** |  | [optional] 
**Id** | Pointer to [**ConnectorTaskId**](ConnectorTaskId.md) |  | [optional] 

## Methods

### NewTaskInfo

`func NewTaskInfo() *TaskInfo`

NewTaskInfo instantiates a new TaskInfo object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTaskInfoWithDefaults

`func NewTaskInfoWithDefaults() *TaskInfo`

NewTaskInfoWithDefaults instantiates a new TaskInfo object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetConfig

`func (o *TaskInfo) GetConfig() map[string]string`

GetConfig returns the Config field if non-nil, zero value otherwise.

### GetConfigOk

`func (o *TaskInfo) GetConfigOk() (*map[string]string, bool)`

GetConfigOk returns a tuple with the Config field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConfig

`func (o *TaskInfo) SetConfig(v map[string]string)`

SetConfig sets Config field to given value.

### HasConfig

`func (o *TaskInfo) HasConfig() bool`

HasConfig returns a boolean if a field has been set.

### GetId

`func (o *TaskInfo) GetId() ConnectorTaskId`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *TaskInfo) GetIdOk() (*ConnectorTaskId, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *TaskInfo) SetId(v ConnectorTaskId)`

SetId sets Id field to given value.

### HasId

`func (o *TaskInfo) HasId() bool`

HasId returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


