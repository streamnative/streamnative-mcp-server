# WorkerStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Message** | Pointer to **string** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 

## Methods

### NewWorkerStatus

`func NewWorkerStatus() *WorkerStatus`

NewWorkerStatus instantiates a new WorkerStatus object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkerStatusWithDefaults

`func NewWorkerStatusWithDefaults() *WorkerStatus`

NewWorkerStatusWithDefaults instantiates a new WorkerStatus object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetMessage

`func (o *WorkerStatus) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *WorkerStatus) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *WorkerStatus) SetMessage(v string)`

SetMessage sets Message field to given value.

### HasMessage

`func (o *WorkerStatus) HasMessage() bool`

HasMessage returns a boolean if a field has been set.

### GetStatus

`func (o *WorkerStatus) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *WorkerStatus) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *WorkerStatus) SetStatus(v string)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *WorkerStatus) HasStatus() bool`

HasStatus returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


