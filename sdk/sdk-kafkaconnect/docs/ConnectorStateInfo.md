# ConnectorStateInfo

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Connector** | Pointer to [**ConnectorState**](ConnectorState.md) |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Tasks** | Pointer to [**[]TaskState**](TaskState.md) |  | [optional] 
**Type** | Pointer to **string** |  | [optional] 

## Methods

### NewConnectorStateInfo

`func NewConnectorStateInfo() *ConnectorStateInfo`

NewConnectorStateInfo instantiates a new ConnectorStateInfo object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConnectorStateInfoWithDefaults

`func NewConnectorStateInfoWithDefaults() *ConnectorStateInfo`

NewConnectorStateInfoWithDefaults instantiates a new ConnectorStateInfo object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetConnector

`func (o *ConnectorStateInfo) GetConnector() ConnectorState`

GetConnector returns the Connector field if non-nil, zero value otherwise.

### GetConnectorOk

`func (o *ConnectorStateInfo) GetConnectorOk() (*ConnectorState, bool)`

GetConnectorOk returns a tuple with the Connector field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConnector

`func (o *ConnectorStateInfo) SetConnector(v ConnectorState)`

SetConnector sets Connector field to given value.

### HasConnector

`func (o *ConnectorStateInfo) HasConnector() bool`

HasConnector returns a boolean if a field has been set.

### GetName

`func (o *ConnectorStateInfo) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ConnectorStateInfo) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ConnectorStateInfo) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ConnectorStateInfo) HasName() bool`

HasName returns a boolean if a field has been set.

### GetTasks

`func (o *ConnectorStateInfo) GetTasks() []TaskState`

GetTasks returns the Tasks field if non-nil, zero value otherwise.

### GetTasksOk

`func (o *ConnectorStateInfo) GetTasksOk() (*[]TaskState, bool)`

GetTasksOk returns a tuple with the Tasks field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTasks

`func (o *ConnectorStateInfo) SetTasks(v []TaskState)`

SetTasks sets Tasks field to given value.

### HasTasks

`func (o *ConnectorStateInfo) HasTasks() bool`

HasTasks returns a boolean if a field has been set.

### GetType

`func (o *ConnectorStateInfo) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *ConnectorStateInfo) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *ConnectorStateInfo) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *ConnectorStateInfo) HasType() bool`

HasType returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


