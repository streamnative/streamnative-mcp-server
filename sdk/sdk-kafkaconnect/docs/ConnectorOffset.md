# ConnectorOffset

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Offset** | Pointer to **map[string]map[string]interface{}** |  | [optional] 
**Partition** | Pointer to **map[string]map[string]interface{}** |  | [optional] 

## Methods

### NewConnectorOffset

`func NewConnectorOffset() *ConnectorOffset`

NewConnectorOffset instantiates a new ConnectorOffset object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConnectorOffsetWithDefaults

`func NewConnectorOffsetWithDefaults() *ConnectorOffset`

NewConnectorOffsetWithDefaults instantiates a new ConnectorOffset object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetOffset

`func (o *ConnectorOffset) GetOffset() map[string]map[string]interface{}`

GetOffset returns the Offset field if non-nil, zero value otherwise.

### GetOffsetOk

`func (o *ConnectorOffset) GetOffsetOk() (*map[string]map[string]interface{}, bool)`

GetOffsetOk returns a tuple with the Offset field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOffset

`func (o *ConnectorOffset) SetOffset(v map[string]map[string]interface{})`

SetOffset sets Offset field to given value.

### HasOffset

`func (o *ConnectorOffset) HasOffset() bool`

HasOffset returns a boolean if a field has been set.

### GetPartition

`func (o *ConnectorOffset) GetPartition() map[string]map[string]interface{}`

GetPartition returns the Partition field if non-nil, zero value otherwise.

### GetPartitionOk

`func (o *ConnectorOffset) GetPartitionOk() (*map[string]map[string]interface{}, bool)`

GetPartitionOk returns a tuple with the Partition field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPartition

`func (o *ConnectorOffset) SetPartition(v map[string]map[string]interface{})`

SetPartition sets Partition field to given value.

### HasPartition

`func (o *ConnectorOffset) HasPartition() bool`

HasPartition returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


