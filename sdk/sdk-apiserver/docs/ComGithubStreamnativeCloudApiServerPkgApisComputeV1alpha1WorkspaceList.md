# ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceList

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ApiVersion** | Pointer to **string** | APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources | [optional] 
**Items** | [**[]ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Workspace**](ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Workspace.md) |  | 
**Kind** | Pointer to **string** | Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds | [optional] 
**Metadata** | Pointer to [**V1ListMeta**](V1ListMeta.md) |  | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceList

`func NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceList(items []ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Workspace, ) *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceList`

NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceList instantiates a new ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceList object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceListWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceListWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceList`

NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceListWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceList object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceList) GetApiVersion() string`

GetApiVersion returns the ApiVersion field if non-nil, zero value otherwise.

### GetApiVersionOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceList) GetApiVersionOk() (*string, bool)`

GetApiVersionOk returns a tuple with the ApiVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceList) SetApiVersion(v string)`

SetApiVersion sets ApiVersion field to given value.

### HasApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceList) HasApiVersion() bool`

HasApiVersion returns a boolean if a field has been set.

### GetItems

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceList) GetItems() []ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Workspace`

GetItems returns the Items field if non-nil, zero value otherwise.

### GetItemsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceList) GetItemsOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Workspace, bool)`

GetItemsOk returns a tuple with the Items field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetItems

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceList) SetItems(v []ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Workspace)`

SetItems sets Items field to given value.


### GetKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceList) GetKind() string`

GetKind returns the Kind field if non-nil, zero value otherwise.

### GetKindOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceList) GetKindOk() (*string, bool)`

GetKindOk returns a tuple with the Kind field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceList) SetKind(v string)`

SetKind sets Kind field to given value.

### HasKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceList) HasKind() bool`

HasKind returns a boolean if a field has been set.

### GetMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceList) GetMetadata() V1ListMeta`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceList) GetMetadataOk() (*V1ListMeta, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceList) SetMetadata(v V1ListMeta)`

SetMetadata sets Metadata field to given value.

### HasMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceList) HasMetadata() bool`

HasMetadata returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


