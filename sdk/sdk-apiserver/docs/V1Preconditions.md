# V1Preconditions

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ResourceVersion** | Pointer to **string** | Specifies the target ResourceVersion | [optional] 
**Uid** | Pointer to **string** | Specifies the target UID. | [optional] 

## Methods

### NewV1Preconditions

`func NewV1Preconditions() *V1Preconditions`

NewV1Preconditions instantiates a new V1Preconditions object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewV1PreconditionsWithDefaults

`func NewV1PreconditionsWithDefaults() *V1Preconditions`

NewV1PreconditionsWithDefaults instantiates a new V1Preconditions object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetResourceVersion

`func (o *V1Preconditions) GetResourceVersion() string`

GetResourceVersion returns the ResourceVersion field if non-nil, zero value otherwise.

### GetResourceVersionOk

`func (o *V1Preconditions) GetResourceVersionOk() (*string, bool)`

GetResourceVersionOk returns a tuple with the ResourceVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResourceVersion

`func (o *V1Preconditions) SetResourceVersion(v string)`

SetResourceVersion sets ResourceVersion field to given value.

### HasResourceVersion

`func (o *V1Preconditions) HasResourceVersion() bool`

HasResourceVersion returns a boolean if a field has been set.

### GetUid

`func (o *V1Preconditions) GetUid() string`

GetUid returns the Uid field if non-nil, zero value otherwise.

### GetUidOk

`func (o *V1Preconditions) GetUidOk() (*string, bool)`

GetUidOk returns a tuple with the Uid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUid

`func (o *V1Preconditions) SetUid(v string)`

SetUid sets Uid field to given value.

### HasUid

`func (o *V1Preconditions) HasUid() bool`

HasUid returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


