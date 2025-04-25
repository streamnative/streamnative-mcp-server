# V1APIGroup

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ApiVersion** | Pointer to **string** | APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources | [optional] 
**Kind** | Pointer to **string** | Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds | [optional] 
**Name** | **string** | name is the name of the group. | 
**PreferredVersion** | Pointer to [**V1GroupVersionForDiscovery**](V1GroupVersionForDiscovery.md) |  | [optional] 
**ServerAddressByClientCIDRs** | Pointer to [**[]V1ServerAddressByClientCIDR**](V1ServerAddressByClientCIDR.md) | a map of client CIDR to server address that is serving this group. This is to help clients reach servers in the most network-efficient way possible. Clients can use the appropriate server address as per the CIDR that they match. In case of multiple matches, clients should use the longest matching CIDR. The server returns only those CIDRs that it thinks that the client can match. For example: the master will return an internal IP CIDR only, if the client reaches the server using an internal IP. Server looks at X-Forwarded-For header or X-Real-Ip header or request.RemoteAddr (in that order) to get the client IP. | [optional] 
**Versions** | [**[]V1GroupVersionForDiscovery**](V1GroupVersionForDiscovery.md) | versions are the versions supported in this group. | 

## Methods

### NewV1APIGroup

`func NewV1APIGroup(name string, versions []V1GroupVersionForDiscovery, ) *V1APIGroup`

NewV1APIGroup instantiates a new V1APIGroup object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewV1APIGroupWithDefaults

`func NewV1APIGroupWithDefaults() *V1APIGroup`

NewV1APIGroupWithDefaults instantiates a new V1APIGroup object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetApiVersion

`func (o *V1APIGroup) GetApiVersion() string`

GetApiVersion returns the ApiVersion field if non-nil, zero value otherwise.

### GetApiVersionOk

`func (o *V1APIGroup) GetApiVersionOk() (*string, bool)`

GetApiVersionOk returns a tuple with the ApiVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetApiVersion

`func (o *V1APIGroup) SetApiVersion(v string)`

SetApiVersion sets ApiVersion field to given value.

### HasApiVersion

`func (o *V1APIGroup) HasApiVersion() bool`

HasApiVersion returns a boolean if a field has been set.

### GetKind

`func (o *V1APIGroup) GetKind() string`

GetKind returns the Kind field if non-nil, zero value otherwise.

### GetKindOk

`func (o *V1APIGroup) GetKindOk() (*string, bool)`

GetKindOk returns a tuple with the Kind field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKind

`func (o *V1APIGroup) SetKind(v string)`

SetKind sets Kind field to given value.

### HasKind

`func (o *V1APIGroup) HasKind() bool`

HasKind returns a boolean if a field has been set.

### GetName

`func (o *V1APIGroup) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *V1APIGroup) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *V1APIGroup) SetName(v string)`

SetName sets Name field to given value.


### GetPreferredVersion

`func (o *V1APIGroup) GetPreferredVersion() V1GroupVersionForDiscovery`

GetPreferredVersion returns the PreferredVersion field if non-nil, zero value otherwise.

### GetPreferredVersionOk

`func (o *V1APIGroup) GetPreferredVersionOk() (*V1GroupVersionForDiscovery, bool)`

GetPreferredVersionOk returns a tuple with the PreferredVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPreferredVersion

`func (o *V1APIGroup) SetPreferredVersion(v V1GroupVersionForDiscovery)`

SetPreferredVersion sets PreferredVersion field to given value.

### HasPreferredVersion

`func (o *V1APIGroup) HasPreferredVersion() bool`

HasPreferredVersion returns a boolean if a field has been set.

### GetServerAddressByClientCIDRs

`func (o *V1APIGroup) GetServerAddressByClientCIDRs() []V1ServerAddressByClientCIDR`

GetServerAddressByClientCIDRs returns the ServerAddressByClientCIDRs field if non-nil, zero value otherwise.

### GetServerAddressByClientCIDRsOk

`func (o *V1APIGroup) GetServerAddressByClientCIDRsOk() (*[]V1ServerAddressByClientCIDR, bool)`

GetServerAddressByClientCIDRsOk returns a tuple with the ServerAddressByClientCIDRs field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServerAddressByClientCIDRs

`func (o *V1APIGroup) SetServerAddressByClientCIDRs(v []V1ServerAddressByClientCIDR)`

SetServerAddressByClientCIDRs sets ServerAddressByClientCIDRs field to given value.

### HasServerAddressByClientCIDRs

`func (o *V1APIGroup) HasServerAddressByClientCIDRs() bool`

HasServerAddressByClientCIDRs returns a boolean if a field has been set.

### GetVersions

`func (o *V1APIGroup) GetVersions() []V1GroupVersionForDiscovery`

GetVersions returns the Versions field if non-nil, zero value otherwise.

### GetVersionsOk

`func (o *V1APIGroup) GetVersionsOk() (*[]V1GroupVersionForDiscovery, bool)`

GetVersionsOk returns a tuple with the Versions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersions

`func (o *V1APIGroup) SetVersions(v []V1GroupVersionForDiscovery)`

SetVersions sets Versions field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


