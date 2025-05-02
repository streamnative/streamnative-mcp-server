# V1StatusCause

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Field** | Pointer to **string** | The field of the resource that has caused this error, as named by its JSON serialization. May include dot and postfix notation for nested attributes. Arrays are zero-indexed.  Fields may appear more than once in an array of causes due to fields having multiple errors. Optional.  Examples:   \&quot;name\&quot; - the field \&quot;name\&quot; on the current resource   \&quot;items[0].name\&quot; - the field \&quot;name\&quot; on the first array entry in \&quot;items\&quot; | [optional] 
**Message** | Pointer to **string** | A human-readable description of the cause of the error.  This field may be presented as-is to a reader. | [optional] 
**Reason** | Pointer to **string** | A machine-readable description of the cause of the error. If this value is empty there is no information available. | [optional] 

## Methods

### NewV1StatusCause

`func NewV1StatusCause() *V1StatusCause`

NewV1StatusCause instantiates a new V1StatusCause object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewV1StatusCauseWithDefaults

`func NewV1StatusCauseWithDefaults() *V1StatusCause`

NewV1StatusCauseWithDefaults instantiates a new V1StatusCause object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetField

`func (o *V1StatusCause) GetField() string`

GetField returns the Field field if non-nil, zero value otherwise.

### GetFieldOk

`func (o *V1StatusCause) GetFieldOk() (*string, bool)`

GetFieldOk returns a tuple with the Field field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetField

`func (o *V1StatusCause) SetField(v string)`

SetField sets Field field to given value.

### HasField

`func (o *V1StatusCause) HasField() bool`

HasField returns a boolean if a field has been set.

### GetMessage

`func (o *V1StatusCause) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *V1StatusCause) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *V1StatusCause) SetMessage(v string)`

SetMessage sets Message field to given value.

### HasMessage

`func (o *V1StatusCause) HasMessage() bool`

HasMessage returns a boolean if a field has been set.

### GetReason

`func (o *V1StatusCause) GetReason() string`

GetReason returns the Reason field if non-nil, zero value otherwise.

### GetReasonOk

`func (o *V1StatusCause) GetReasonOk() (*string, bool)`

GetReasonOk returns a tuple with the Reason field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReason

`func (o *V1StatusCause) SetReason(v string)`

SetReason sets Reason field to given value.

### HasReason

`func (o *V1StatusCause) HasReason() bool`

HasReason returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


