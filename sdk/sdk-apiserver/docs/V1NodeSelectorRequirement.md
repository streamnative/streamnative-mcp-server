# V1NodeSelectorRequirement

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Key** | **string** | The label key that the selector applies to. | 
**Operator** | **string** | Represents a key&#39;s relationship to a set of values. Valid operators are In, NotIn, Exists, DoesNotExist. Gt, and Lt.  Possible enum values:  - &#x60;\&quot;DoesNotExist\&quot;&#x60;  - &#x60;\&quot;Exists\&quot;&#x60;  - &#x60;\&quot;Gt\&quot;&#x60;  - &#x60;\&quot;In\&quot;&#x60;  - &#x60;\&quot;Lt\&quot;&#x60;  - &#x60;\&quot;NotIn\&quot;&#x60; | 
**Values** | Pointer to **[]string** | An array of string values. If the operator is In or NotIn, the values array must be non-empty. If the operator is Exists or DoesNotExist, the values array must be empty. If the operator is Gt or Lt, the values array must have a single element, which will be interpreted as an integer. This array is replaced during a strategic merge patch. | [optional] 

## Methods

### NewV1NodeSelectorRequirement

`func NewV1NodeSelectorRequirement(key string, operator string, ) *V1NodeSelectorRequirement`

NewV1NodeSelectorRequirement instantiates a new V1NodeSelectorRequirement object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewV1NodeSelectorRequirementWithDefaults

`func NewV1NodeSelectorRequirementWithDefaults() *V1NodeSelectorRequirement`

NewV1NodeSelectorRequirementWithDefaults instantiates a new V1NodeSelectorRequirement object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetKey

`func (o *V1NodeSelectorRequirement) GetKey() string`

GetKey returns the Key field if non-nil, zero value otherwise.

### GetKeyOk

`func (o *V1NodeSelectorRequirement) GetKeyOk() (*string, bool)`

GetKeyOk returns a tuple with the Key field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKey

`func (o *V1NodeSelectorRequirement) SetKey(v string)`

SetKey sets Key field to given value.


### GetOperator

`func (o *V1NodeSelectorRequirement) GetOperator() string`

GetOperator returns the Operator field if non-nil, zero value otherwise.

### GetOperatorOk

`func (o *V1NodeSelectorRequirement) GetOperatorOk() (*string, bool)`

GetOperatorOk returns a tuple with the Operator field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOperator

`func (o *V1NodeSelectorRequirement) SetOperator(v string)`

SetOperator sets Operator field to given value.


### GetValues

`func (o *V1NodeSelectorRequirement) GetValues() []string`

GetValues returns the Values field if non-nil, zero value otherwise.

### GetValuesOk

`func (o *V1NodeSelectorRequirement) GetValuesOk() (*[]string, bool)`

GetValuesOk returns a tuple with the Values field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetValues

`func (o *V1NodeSelectorRequirement) SetValues(v []string)`

SetValues sets Values field to given value.

### HasValues

`func (o *V1NodeSelectorRequirement) HasValues() bool`

HasValues returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


