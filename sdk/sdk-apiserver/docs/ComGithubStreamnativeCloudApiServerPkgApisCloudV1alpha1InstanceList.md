# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceList

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ApiVersion** | Pointer to **string** | APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources | [optional] 
**Items** | [**[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceListItemsInner**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceListItemsInner.md) |  | 
**Kind** | Pointer to **string** | Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds | [optional] 
**Metadata** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceListMetadata**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceListMetadata.md) |  | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceList

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceList(items []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceListItemsInner, ) *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceList`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceList instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceList object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceListWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceListWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceList`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceListWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceList object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceList) GetApiVersion() string`

GetApiVersion returns the ApiVersion field if non-nil, zero value otherwise.

### GetApiVersionOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceList) GetApiVersionOk() (*string, bool)`

GetApiVersionOk returns a tuple with the ApiVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceList) SetApiVersion(v string)`

SetApiVersion sets ApiVersion field to given value.

### HasApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceList) HasApiVersion() bool`

HasApiVersion returns a boolean if a field has been set.

### GetItems

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceList) GetItems() []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceListItemsInner`

GetItems returns the Items field if non-nil, zero value otherwise.

### GetItemsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceList) GetItemsOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceListItemsInner, bool)`

GetItemsOk returns a tuple with the Items field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetItems

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceList) SetItems(v []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceListItemsInner)`

SetItems sets Items field to given value.


### GetKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceList) GetKind() string`

GetKind returns the Kind field if non-nil, zero value otherwise.

### GetKindOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceList) GetKindOk() (*string, bool)`

GetKindOk returns a tuple with the Kind field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceList) SetKind(v string)`

SetKind sets Kind field to given value.

### HasKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceList) HasKind() bool`

HasKind returns a boolean if a field has been set.

### GetMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceList) GetMetadata() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceListMetadata`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceList) GetMetadataOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceListMetadata, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceList) SetMetadata(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceListMetadata)`

SetMetadata sets Metadata field to given value.

### HasMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceList) HasMetadata() bool`

HasMetadata returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


