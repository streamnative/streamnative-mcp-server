# V1Status

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ApiVersion** | Pointer to **string** | APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources | [optional] 
**Code** | Pointer to **int32** | Suggested HTTP return code for this status, 0 if not set. | [optional] 
**Details** | Pointer to [**V1StatusDetails**](V1StatusDetails.md) |  | [optional] 
**Kind** | Pointer to **string** | Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds | [optional] 
**Message** | Pointer to **string** | A human-readable description of the status of this operation. | [optional] 
**Metadata** | Pointer to [**V1ListMeta**](V1ListMeta.md) |  | [optional] 
**Reason** | Pointer to **string** | A machine-readable description of why this operation is in the \&quot;Failure\&quot; status. If this value is empty there is no information available. A Reason clarifies an HTTP status code but does not override it. | [optional] 
**Status** | Pointer to **string** | Status of the operation. One of: \&quot;Success\&quot; or \&quot;Failure\&quot;. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status | [optional] 

## Methods

### NewV1Status

`func NewV1Status() *V1Status`

NewV1Status instantiates a new V1Status object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewV1StatusWithDefaults

`func NewV1StatusWithDefaults() *V1Status`

NewV1StatusWithDefaults instantiates a new V1Status object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetApiVersion

`func (o *V1Status) GetApiVersion() string`

GetApiVersion returns the ApiVersion field if non-nil, zero value otherwise.

### GetApiVersionOk

`func (o *V1Status) GetApiVersionOk() (*string, bool)`

GetApiVersionOk returns a tuple with the ApiVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetApiVersion

`func (o *V1Status) SetApiVersion(v string)`

SetApiVersion sets ApiVersion field to given value.

### HasApiVersion

`func (o *V1Status) HasApiVersion() bool`

HasApiVersion returns a boolean if a field has been set.

### GetCode

`func (o *V1Status) GetCode() int32`

GetCode returns the Code field if non-nil, zero value otherwise.

### GetCodeOk

`func (o *V1Status) GetCodeOk() (*int32, bool)`

GetCodeOk returns a tuple with the Code field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCode

`func (o *V1Status) SetCode(v int32)`

SetCode sets Code field to given value.

### HasCode

`func (o *V1Status) HasCode() bool`

HasCode returns a boolean if a field has been set.

### GetDetails

`func (o *V1Status) GetDetails() V1StatusDetails`

GetDetails returns the Details field if non-nil, zero value otherwise.

### GetDetailsOk

`func (o *V1Status) GetDetailsOk() (*V1StatusDetails, bool)`

GetDetailsOk returns a tuple with the Details field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDetails

`func (o *V1Status) SetDetails(v V1StatusDetails)`

SetDetails sets Details field to given value.

### HasDetails

`func (o *V1Status) HasDetails() bool`

HasDetails returns a boolean if a field has been set.

### GetKind

`func (o *V1Status) GetKind() string`

GetKind returns the Kind field if non-nil, zero value otherwise.

### GetKindOk

`func (o *V1Status) GetKindOk() (*string, bool)`

GetKindOk returns a tuple with the Kind field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKind

`func (o *V1Status) SetKind(v string)`

SetKind sets Kind field to given value.

### HasKind

`func (o *V1Status) HasKind() bool`

HasKind returns a boolean if a field has been set.

### GetMessage

`func (o *V1Status) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *V1Status) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *V1Status) SetMessage(v string)`

SetMessage sets Message field to given value.

### HasMessage

`func (o *V1Status) HasMessage() bool`

HasMessage returns a boolean if a field has been set.

### GetMetadata

`func (o *V1Status) GetMetadata() V1ListMeta`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *V1Status) GetMetadataOk() (*V1ListMeta, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *V1Status) SetMetadata(v V1ListMeta)`

SetMetadata sets Metadata field to given value.

### HasMetadata

`func (o *V1Status) HasMetadata() bool`

HasMetadata returns a boolean if a field has been set.

### GetReason

`func (o *V1Status) GetReason() string`

GetReason returns the Reason field if non-nil, zero value otherwise.

### GetReasonOk

`func (o *V1Status) GetReasonOk() (*string, bool)`

GetReasonOk returns a tuple with the Reason field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReason

`func (o *V1Status) SetReason(v string)`

SetReason sets Reason field to given value.

### HasReason

`func (o *V1Status) HasReason() bool`

HasReason returns a boolean if a field has been set.

### GetStatus

`func (o *V1Status) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *V1Status) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *V1Status) SetStatus(v string)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *V1Status) HasStatus() bool`

HasStatus returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


