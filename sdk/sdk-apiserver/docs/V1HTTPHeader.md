# V1HTTPHeader

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | The header field name | 
**Value** | **string** | The header field value | 

## Methods

### NewV1HTTPHeader

`func NewV1HTTPHeader(name string, value string, ) *V1HTTPHeader`

NewV1HTTPHeader instantiates a new V1HTTPHeader object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewV1HTTPHeaderWithDefaults

`func NewV1HTTPHeaderWithDefaults() *V1HTTPHeader`

NewV1HTTPHeaderWithDefaults instantiates a new V1HTTPHeader object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *V1HTTPHeader) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *V1HTTPHeader) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *V1HTTPHeader) SetName(v string)`

SetName sets Name field to given value.


### GetValue

`func (o *V1HTTPHeader) GetValue() string`

GetValue returns the Value field if non-nil, zero value otherwise.

### GetValueOk

`func (o *V1HTTPHeader) GetValueOk() (*string, bool)`

GetValueOk returns a tuple with the Value field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetValue

`func (o *V1HTTPHeader) SetValue(v string)`

SetValue sets Value field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


