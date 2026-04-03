# CreateConnectorRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Config** | Pointer to **map[string]string** |  | [optional] 
**InitialState** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 

## Methods

### NewCreateConnectorRequest

`func NewCreateConnectorRequest() *CreateConnectorRequest`

NewCreateConnectorRequest instantiates a new CreateConnectorRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateConnectorRequestWithDefaults

`func NewCreateConnectorRequestWithDefaults() *CreateConnectorRequest`

NewCreateConnectorRequestWithDefaults instantiates a new CreateConnectorRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetConfig

`func (o *CreateConnectorRequest) GetConfig() map[string]string`

GetConfig returns the Config field if non-nil, zero value otherwise.

### GetConfigOk

`func (o *CreateConnectorRequest) GetConfigOk() (*map[string]string, bool)`

GetConfigOk returns a tuple with the Config field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConfig

`func (o *CreateConnectorRequest) SetConfig(v map[string]string)`

SetConfig sets Config field to given value.

### HasConfig

`func (o *CreateConnectorRequest) HasConfig() bool`

HasConfig returns a boolean if a field has been set.

### GetInitialState

`func (o *CreateConnectorRequest) GetInitialState() string`

GetInitialState returns the InitialState field if non-nil, zero value otherwise.

### GetInitialStateOk

`func (o *CreateConnectorRequest) GetInitialStateOk() (*string, bool)`

GetInitialStateOk returns a tuple with the InitialState field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInitialState

`func (o *CreateConnectorRequest) SetInitialState(v string)`

SetInitialState sets InitialState field to given value.

### HasInitialState

`func (o *CreateConnectorRequest) HasInitialState() bool`

HasInitialState returns a boolean if a field has been set.

### GetName

`func (o *CreateConnectorRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *CreateConnectorRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *CreateConnectorRequest) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *CreateConnectorRequest) HasName() bool`

HasName returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


