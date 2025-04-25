# V1GroupVersionForDiscovery

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**GroupVersion** | **string** | groupVersion specifies the API group and version in the form \&quot;group/version\&quot; | 
**Version** | **string** | version specifies the version in the form of \&quot;version\&quot;. This is to save the clients the trouble of splitting the GroupVersion. | 

## Methods

### NewV1GroupVersionForDiscovery

`func NewV1GroupVersionForDiscovery(groupVersion string, version string, ) *V1GroupVersionForDiscovery`

NewV1GroupVersionForDiscovery instantiates a new V1GroupVersionForDiscovery object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewV1GroupVersionForDiscoveryWithDefaults

`func NewV1GroupVersionForDiscoveryWithDefaults() *V1GroupVersionForDiscovery`

NewV1GroupVersionForDiscoveryWithDefaults instantiates a new V1GroupVersionForDiscovery object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetGroupVersion

`func (o *V1GroupVersionForDiscovery) GetGroupVersion() string`

GetGroupVersion returns the GroupVersion field if non-nil, zero value otherwise.

### GetGroupVersionOk

`func (o *V1GroupVersionForDiscovery) GetGroupVersionOk() (*string, bool)`

GetGroupVersionOk returns a tuple with the GroupVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGroupVersion

`func (o *V1GroupVersionForDiscovery) SetGroupVersion(v string)`

SetGroupVersion sets GroupVersion field to given value.


### GetVersion

`func (o *V1GroupVersionForDiscovery) GetVersion() string`

GetVersion returns the Version field if non-nil, zero value otherwise.

### GetVersionOk

`func (o *V1GroupVersionForDiscovery) GetVersionOk() (*string, bool)`

GetVersionOk returns a tuple with the Version field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersion

`func (o *V1GroupVersionForDiscovery) SetVersion(v string)`

SetVersion sets Version field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


