# ConnectorInfo

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Config** | Pointer to **map[string]string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Tasks** | Pointer to [**[]ConnectorTaskId**](ConnectorTaskId.md) |  | [optional] 
**Type** | Pointer to **string** |  | [optional] 

## Methods

### NewConnectorInfo

`func NewConnectorInfo() *ConnectorInfo`

NewConnectorInfo instantiates a new ConnectorInfo object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConnectorInfoWithDefaults

`func NewConnectorInfoWithDefaults() *ConnectorInfo`

NewConnectorInfoWithDefaults instantiates a new ConnectorInfo object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetConfig

`func (o *ConnectorInfo) GetConfig() map[string]string`

GetConfig returns the Config field if non-nil, zero value otherwise.

### GetConfigOk

`func (o *ConnectorInfo) GetConfigOk() (*map[string]string, bool)`

GetConfigOk returns a tuple with the Config field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConfig

`func (o *ConnectorInfo) SetConfig(v map[string]string)`

SetConfig sets Config field to given value.

### HasConfig

`func (o *ConnectorInfo) HasConfig() bool`

HasConfig returns a boolean if a field has been set.

### GetName

`func (o *ConnectorInfo) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ConnectorInfo) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ConnectorInfo) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ConnectorInfo) HasName() bool`

HasName returns a boolean if a field has been set.

### GetTasks

`func (o *ConnectorInfo) GetTasks() []ConnectorTaskId`

GetTasks returns the Tasks field if non-nil, zero value otherwise.

### GetTasksOk

`func (o *ConnectorInfo) GetTasksOk() (*[]ConnectorTaskId, bool)`

GetTasksOk returns a tuple with the Tasks field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTasks

`func (o *ConnectorInfo) SetTasks(v []ConnectorTaskId)`

SetTasks sets Tasks field to given value.

### HasTasks

`func (o *ConnectorInfo) HasTasks() bool`

HasTasks returns a boolean if a field has been set.

### GetType

`func (o *ConnectorInfo) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *ConnectorInfo) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *ConnectorInfo) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *ConnectorInfo) HasType() bool`

HasType returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


