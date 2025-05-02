# V1APIGroupList

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ApiVersion** | Pointer to **string** | APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources | [optional] 
**Groups** | [**[]V1APIGroup**](V1APIGroup.md) | groups is a list of APIGroup. | 
**Kind** | Pointer to **string** | Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds | [optional] 

## Methods

### NewV1APIGroupList

`func NewV1APIGroupList(groups []V1APIGroup, ) *V1APIGroupList`

NewV1APIGroupList instantiates a new V1APIGroupList object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewV1APIGroupListWithDefaults

`func NewV1APIGroupListWithDefaults() *V1APIGroupList`

NewV1APIGroupListWithDefaults instantiates a new V1APIGroupList object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetApiVersion

`func (o *V1APIGroupList) GetApiVersion() string`

GetApiVersion returns the ApiVersion field if non-nil, zero value otherwise.

### GetApiVersionOk

`func (o *V1APIGroupList) GetApiVersionOk() (*string, bool)`

GetApiVersionOk returns a tuple with the ApiVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetApiVersion

`func (o *V1APIGroupList) SetApiVersion(v string)`

SetApiVersion sets ApiVersion field to given value.

### HasApiVersion

`func (o *V1APIGroupList) HasApiVersion() bool`

HasApiVersion returns a boolean if a field has been set.

### GetGroups

`func (o *V1APIGroupList) GetGroups() []V1APIGroup`

GetGroups returns the Groups field if non-nil, zero value otherwise.

### GetGroupsOk

`func (o *V1APIGroupList) GetGroupsOk() (*[]V1APIGroup, bool)`

GetGroupsOk returns a tuple with the Groups field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGroups

`func (o *V1APIGroupList) SetGroups(v []V1APIGroup)`

SetGroups sets Groups field to given value.


### GetKind

`func (o *V1APIGroupList) GetKind() string`

GetKind returns the Kind field if non-nil, zero value otherwise.

### GetKindOk

`func (o *V1APIGroupList) GetKindOk() (*string, bool)`

GetKindOk returns a tuple with the Kind field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKind

`func (o *V1APIGroupList) SetKind(v string)`

SetKind sets Kind field to given value.

### HasKind

`func (o *V1APIGroupList) HasKind() bool`

HasKind returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


