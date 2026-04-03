# ConnectorTaskId

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Connector** | Pointer to **string** |  | [optional] 
**Task** | Pointer to **int32** |  | [optional] 

## Methods

### NewConnectorTaskId

`func NewConnectorTaskId() *ConnectorTaskId`

NewConnectorTaskId instantiates a new ConnectorTaskId object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConnectorTaskIdWithDefaults

`func NewConnectorTaskIdWithDefaults() *ConnectorTaskId`

NewConnectorTaskIdWithDefaults instantiates a new ConnectorTaskId object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetConnector

`func (o *ConnectorTaskId) GetConnector() string`

GetConnector returns the Connector field if non-nil, zero value otherwise.

### GetConnectorOk

`func (o *ConnectorTaskId) GetConnectorOk() (*string, bool)`

GetConnectorOk returns a tuple with the Connector field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConnector

`func (o *ConnectorTaskId) SetConnector(v string)`

SetConnector sets Connector field to given value.

### HasConnector

`func (o *ConnectorTaskId) HasConnector() bool`

HasConnector returns a boolean if a field has been set.

### GetTask

`func (o *ConnectorTaskId) GetTask() int32`

GetTask returns the Task field if non-nil, zero value otherwise.

### GetTaskOk

`func (o *ConnectorTaskId) GetTaskOk() (*int32, bool)`

GetTaskOk returns a tuple with the Task field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTask

`func (o *ConnectorTaskId) SetTask(v int32)`

SetTask sets Task field to given value.

### HasTask

`func (o *ConnectorTaskId) HasTask() bool`

HasTask returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


