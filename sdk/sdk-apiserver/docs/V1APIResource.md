# V1APIResource

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Categories** | Pointer to **[]string** | categories is a list of the grouped resources this resource belongs to (e.g. &#39;all&#39;) | [optional] 
**Group** | Pointer to **string** | group is the preferred group of the resource.  Empty implies the group of the containing resource list. For subresources, this may have a different value, for example: Scale\&quot;. | [optional] 
**Kind** | **string** | kind is the kind for the resource (e.g. &#39;Foo&#39; is the kind for a resource &#39;foo&#39;) | 
**Name** | **string** | name is the plural name of the resource. | 
**Namespaced** | **bool** | namespaced indicates if a resource is namespaced or not. | 
**ShortNames** | Pointer to **[]string** | shortNames is a list of suggested short names of the resource. | [optional] 
**SingularName** | **string** | singularName is the singular name of the resource.  This allows clients to handle plural and singular opaquely. The singularName is more correct for reporting status on a single item and both singular and plural are allowed from the kubectl CLI interface. | 
**StorageVersionHash** | Pointer to **string** | The hash value of the storage version, the version this resource is converted to when written to the data store. Value must be treated as opaque by clients. Only equality comparison on the value is valid. This is an alpha feature and may change or be removed in the future. The field is populated by the apiserver only if the StorageVersionHash feature gate is enabled. This field will remain optional even if it graduates. | [optional] 
**Verbs** | **[]string** | verbs is a list of supported kube verbs (this includes get, list, watch, create, update, patch, delete, deletecollection, and proxy) | 
**Version** | Pointer to **string** | version is the preferred version of the resource.  Empty implies the version of the containing resource list For subresources, this may have a different value, for example: v1 (while inside a v1beta1 version of the core resource&#39;s group)\&quot;. | [optional] 

## Methods

### NewV1APIResource

`func NewV1APIResource(kind string, name string, namespaced bool, singularName string, verbs []string, ) *V1APIResource`

NewV1APIResource instantiates a new V1APIResource object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewV1APIResourceWithDefaults

`func NewV1APIResourceWithDefaults() *V1APIResource`

NewV1APIResourceWithDefaults instantiates a new V1APIResource object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCategories

`func (o *V1APIResource) GetCategories() []string`

GetCategories returns the Categories field if non-nil, zero value otherwise.

### GetCategoriesOk

`func (o *V1APIResource) GetCategoriesOk() (*[]string, bool)`

GetCategoriesOk returns a tuple with the Categories field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCategories

`func (o *V1APIResource) SetCategories(v []string)`

SetCategories sets Categories field to given value.

### HasCategories

`func (o *V1APIResource) HasCategories() bool`

HasCategories returns a boolean if a field has been set.

### GetGroup

`func (o *V1APIResource) GetGroup() string`

GetGroup returns the Group field if non-nil, zero value otherwise.

### GetGroupOk

`func (o *V1APIResource) GetGroupOk() (*string, bool)`

GetGroupOk returns a tuple with the Group field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGroup

`func (o *V1APIResource) SetGroup(v string)`

SetGroup sets Group field to given value.

### HasGroup

`func (o *V1APIResource) HasGroup() bool`

HasGroup returns a boolean if a field has been set.

### GetKind

`func (o *V1APIResource) GetKind() string`

GetKind returns the Kind field if non-nil, zero value otherwise.

### GetKindOk

`func (o *V1APIResource) GetKindOk() (*string, bool)`

GetKindOk returns a tuple with the Kind field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKind

`func (o *V1APIResource) SetKind(v string)`

SetKind sets Kind field to given value.


### GetName

`func (o *V1APIResource) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *V1APIResource) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *V1APIResource) SetName(v string)`

SetName sets Name field to given value.


### GetNamespaced

`func (o *V1APIResource) GetNamespaced() bool`

GetNamespaced returns the Namespaced field if non-nil, zero value otherwise.

### GetNamespacedOk

`func (o *V1APIResource) GetNamespacedOk() (*bool, bool)`

GetNamespacedOk returns a tuple with the Namespaced field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNamespaced

`func (o *V1APIResource) SetNamespaced(v bool)`

SetNamespaced sets Namespaced field to given value.


### GetShortNames

`func (o *V1APIResource) GetShortNames() []string`

GetShortNames returns the ShortNames field if non-nil, zero value otherwise.

### GetShortNamesOk

`func (o *V1APIResource) GetShortNamesOk() (*[]string, bool)`

GetShortNamesOk returns a tuple with the ShortNames field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetShortNames

`func (o *V1APIResource) SetShortNames(v []string)`

SetShortNames sets ShortNames field to given value.

### HasShortNames

`func (o *V1APIResource) HasShortNames() bool`

HasShortNames returns a boolean if a field has been set.

### GetSingularName

`func (o *V1APIResource) GetSingularName() string`

GetSingularName returns the SingularName field if non-nil, zero value otherwise.

### GetSingularNameOk

`func (o *V1APIResource) GetSingularNameOk() (*string, bool)`

GetSingularNameOk returns a tuple with the SingularName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSingularName

`func (o *V1APIResource) SetSingularName(v string)`

SetSingularName sets SingularName field to given value.


### GetStorageVersionHash

`func (o *V1APIResource) GetStorageVersionHash() string`

GetStorageVersionHash returns the StorageVersionHash field if non-nil, zero value otherwise.

### GetStorageVersionHashOk

`func (o *V1APIResource) GetStorageVersionHashOk() (*string, bool)`

GetStorageVersionHashOk returns a tuple with the StorageVersionHash field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStorageVersionHash

`func (o *V1APIResource) SetStorageVersionHash(v string)`

SetStorageVersionHash sets StorageVersionHash field to given value.

### HasStorageVersionHash

`func (o *V1APIResource) HasStorageVersionHash() bool`

HasStorageVersionHash returns a boolean if a field has been set.

### GetVerbs

`func (o *V1APIResource) GetVerbs() []string`

GetVerbs returns the Verbs field if non-nil, zero value otherwise.

### GetVerbsOk

`func (o *V1APIResource) GetVerbsOk() (*[]string, bool)`

GetVerbsOk returns a tuple with the Verbs field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVerbs

`func (o *V1APIResource) SetVerbs(v []string)`

SetVerbs sets Verbs field to given value.


### GetVersion

`func (o *V1APIResource) GetVersion() string`

GetVersion returns the Version field if non-nil, zero value otherwise.

### GetVersionOk

`func (o *V1APIResource) GetVersionOk() (*string, bool)`

GetVersionOk returns a tuple with the Version field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersion

`func (o *V1APIResource) SetVersion(v string)`

SetVersion sets Version field to given value.

### HasVersion

`func (o *V1APIResource) HasVersion() bool`

HasVersion returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


