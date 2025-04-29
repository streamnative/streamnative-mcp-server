# ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Workspace

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ApiVersion** | Pointer to **string** | APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources | [optional] 
**Kind** | Pointer to **string** | Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds | [optional] 
**Metadata** | Pointer to [**V1ObjectMeta**](V1ObjectMeta.md) |  | [optional] 
**Spec** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceSpec**](ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceSpec.md) |  | [optional] 
**Status** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceStatus**](ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceStatus.md) |  | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Workspace

`func NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Workspace() *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Workspace`

NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Workspace instantiates a new ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Workspace object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Workspace`

NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Workspace object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Workspace) GetApiVersion() string`

GetApiVersion returns the ApiVersion field if non-nil, zero value otherwise.

### GetApiVersionOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Workspace) GetApiVersionOk() (*string, bool)`

GetApiVersionOk returns a tuple with the ApiVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Workspace) SetApiVersion(v string)`

SetApiVersion sets ApiVersion field to given value.

### HasApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Workspace) HasApiVersion() bool`

HasApiVersion returns a boolean if a field has been set.

### GetKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Workspace) GetKind() string`

GetKind returns the Kind field if non-nil, zero value otherwise.

### GetKindOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Workspace) GetKindOk() (*string, bool)`

GetKindOk returns a tuple with the Kind field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Workspace) SetKind(v string)`

SetKind sets Kind field to given value.

### HasKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Workspace) HasKind() bool`

HasKind returns a boolean if a field has been set.

### GetMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Workspace) GetMetadata() V1ObjectMeta`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Workspace) GetMetadataOk() (*V1ObjectMeta, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Workspace) SetMetadata(v V1ObjectMeta)`

SetMetadata sets Metadata field to given value.

### HasMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Workspace) HasMetadata() bool`

HasMetadata returns a boolean if a field has been set.

### GetSpec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Workspace) GetSpec() ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceSpec`

GetSpec returns the Spec field if non-nil, zero value otherwise.

### GetSpecOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Workspace) GetSpecOk() (*ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceSpec, bool)`

GetSpecOk returns a tuple with the Spec field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSpec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Workspace) SetSpec(v ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceSpec)`

SetSpec sets Spec field to given value.

### HasSpec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Workspace) HasSpec() bool`

HasSpec returns a boolean if a field has been set.

### GetStatus

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Workspace) GetStatus() ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Workspace) GetStatusOk() (*ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Workspace) SetStatus(v ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Workspace) HasStatus() bool`

HasStatus returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


