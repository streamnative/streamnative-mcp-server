# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyList

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ApiVersion** | Pointer to **string** | APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources | [optional] 
**Items** | [**[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKey**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKey.md) |  | 
**Kind** | Pointer to **string** | Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds | [optional] 
**Metadata** | Pointer to [**V1ListMeta**](V1ListMeta.md) |  | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyList

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyList(items []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKey, ) *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyList`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyList instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyList object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyListWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyListWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyList`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyListWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyList object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyList) GetApiVersion() string`

GetApiVersion returns the ApiVersion field if non-nil, zero value otherwise.

### GetApiVersionOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyList) GetApiVersionOk() (*string, bool)`

GetApiVersionOk returns a tuple with the ApiVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyList) SetApiVersion(v string)`

SetApiVersion sets ApiVersion field to given value.

### HasApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyList) HasApiVersion() bool`

HasApiVersion returns a boolean if a field has been set.

### GetItems

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyList) GetItems() []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKey`

GetItems returns the Items field if non-nil, zero value otherwise.

### GetItemsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyList) GetItemsOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKey, bool)`

GetItemsOk returns a tuple with the Items field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetItems

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyList) SetItems(v []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKey)`

SetItems sets Items field to given value.


### GetKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyList) GetKind() string`

GetKind returns the Kind field if non-nil, zero value otherwise.

### GetKindOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyList) GetKindOk() (*string, bool)`

GetKindOk returns a tuple with the Kind field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyList) SetKind(v string)`

SetKind sets Kind field to given value.

### HasKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyList) HasKind() bool`

HasKind returns a boolean if a field has been set.

### GetMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyList) GetMetadata() V1ListMeta`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyList) GetMetadataOk() (*V1ListMeta, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyList) SetMetadata(v V1ListMeta)`

SetMetadata sets Metadata field to given value.

### HasMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyList) HasMetadata() bool`

HasMetadata returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


