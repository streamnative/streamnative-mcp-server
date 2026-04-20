# IoK8sApimachineryPkgApisMetaV1ListMeta

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Continue** | Pointer to **string** | continue may be set if the user set a limit on the number of items returned, and indicates that the server has more data available. The value is opaque and may be used to issue another request to the endpoint that served this list to retrieve the next set of available objects. Continuing a consistent list may not be possible if the server configuration has changed or more than a few minutes have passed. The resourceVersion field returned when using this continue value will be identical to the value in the first response, unless you have received this token from an error message. | [optional] 
**RemainingItemCount** | Pointer to **int64** | remainingItemCount is the number of subsequent items in the list which are not included in this list response. If the list request contained label or field selectors, then the number of remaining items is unknown and the field will be left unset and omitted during serialization. If the list is complete (either because it is not chunking or because this is the last chunk), then there are no more remaining items and this field will be left unset and omitted during serialization. Servers older than v1.15 do not set this field. The intended use of the remainingItemCount is *estimating* the size of a collection. Clients should not rely on the remainingItemCount to be set or to be exact. | [optional] 
**ResourceVersion** | Pointer to **string** | String that identifies the server&#39;s internal version of this object that can be used by clients to determine when objects have changed. Value must be treated as opaque by clients and passed unmodified back to the server. Populated by the system. Read-only. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#concurrency-control-and-consistency | [optional] 
**SelfLink** | Pointer to **string** | Deprecated: selfLink is a legacy read-only field that is no longer populated by the system. | [optional] 

## Methods

### NewIoK8sApimachineryPkgApisMetaV1ListMeta

`func NewIoK8sApimachineryPkgApisMetaV1ListMeta() *IoK8sApimachineryPkgApisMetaV1ListMeta`

NewIoK8sApimachineryPkgApisMetaV1ListMeta instantiates a new IoK8sApimachineryPkgApisMetaV1ListMeta object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewIoK8sApimachineryPkgApisMetaV1ListMetaWithDefaults

`func NewIoK8sApimachineryPkgApisMetaV1ListMetaWithDefaults() *IoK8sApimachineryPkgApisMetaV1ListMeta`

NewIoK8sApimachineryPkgApisMetaV1ListMetaWithDefaults instantiates a new IoK8sApimachineryPkgApisMetaV1ListMeta object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetContinue

`func (o *IoK8sApimachineryPkgApisMetaV1ListMeta) GetContinue() string`

GetContinue returns the Continue field if non-nil, zero value otherwise.

### GetContinueOk

`func (o *IoK8sApimachineryPkgApisMetaV1ListMeta) GetContinueOk() (*string, bool)`

GetContinueOk returns a tuple with the Continue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContinue

`func (o *IoK8sApimachineryPkgApisMetaV1ListMeta) SetContinue(v string)`

SetContinue sets Continue field to given value.

### HasContinue

`func (o *IoK8sApimachineryPkgApisMetaV1ListMeta) HasContinue() bool`

HasContinue returns a boolean if a field has been set.

### GetRemainingItemCount

`func (o *IoK8sApimachineryPkgApisMetaV1ListMeta) GetRemainingItemCount() int64`

GetRemainingItemCount returns the RemainingItemCount field if non-nil, zero value otherwise.

### GetRemainingItemCountOk

`func (o *IoK8sApimachineryPkgApisMetaV1ListMeta) GetRemainingItemCountOk() (*int64, bool)`

GetRemainingItemCountOk returns a tuple with the RemainingItemCount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRemainingItemCount

`func (o *IoK8sApimachineryPkgApisMetaV1ListMeta) SetRemainingItemCount(v int64)`

SetRemainingItemCount sets RemainingItemCount field to given value.

### HasRemainingItemCount

`func (o *IoK8sApimachineryPkgApisMetaV1ListMeta) HasRemainingItemCount() bool`

HasRemainingItemCount returns a boolean if a field has been set.

### GetResourceVersion

`func (o *IoK8sApimachineryPkgApisMetaV1ListMeta) GetResourceVersion() string`

GetResourceVersion returns the ResourceVersion field if non-nil, zero value otherwise.

### GetResourceVersionOk

`func (o *IoK8sApimachineryPkgApisMetaV1ListMeta) GetResourceVersionOk() (*string, bool)`

GetResourceVersionOk returns a tuple with the ResourceVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResourceVersion

`func (o *IoK8sApimachineryPkgApisMetaV1ListMeta) SetResourceVersion(v string)`

SetResourceVersion sets ResourceVersion field to given value.

### HasResourceVersion

`func (o *IoK8sApimachineryPkgApisMetaV1ListMeta) HasResourceVersion() bool`

HasResourceVersion returns a boolean if a field has been set.

### GetSelfLink

`func (o *IoK8sApimachineryPkgApisMetaV1ListMeta) GetSelfLink() string`

GetSelfLink returns the SelfLink field if non-nil, zero value otherwise.

### GetSelfLinkOk

`func (o *IoK8sApimachineryPkgApisMetaV1ListMeta) GetSelfLinkOk() (*string, bool)`

GetSelfLinkOk returns a tuple with the SelfLink field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSelfLink

`func (o *IoK8sApimachineryPkgApisMetaV1ListMeta) SetSelfLink(v string)`

SetSelfLink sets SelfLink field to given value.

### HasSelfLink

`func (o *IoK8sApimachineryPkgApisMetaV1ListMeta) HasSelfLink() bool`

HasSelfLink returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


