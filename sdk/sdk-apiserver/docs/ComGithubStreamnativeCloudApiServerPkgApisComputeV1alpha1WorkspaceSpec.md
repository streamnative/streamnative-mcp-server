# ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**FlinkBlobStorage** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1FlinkBlobStorage**](ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1FlinkBlobStorage.md) |  | [optional] 
**PoolRef** | [**ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PoolRef**](ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PoolRef.md) |  | 
**PulsarClusterNames** | **[]string** | PulsarClusterNames is the list of Pulsar clusters that the workspace will have access to. | 
**UseExternalAccess** | Pointer to **bool** | UseExternalAccess is the flag to indicate whether the workspace will use external access. | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceSpec

`func NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceSpec(poolRef ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PoolRef, pulsarClusterNames []string, ) *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceSpec`

NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceSpec instantiates a new ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceSpec object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceSpecWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceSpecWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceSpec`

NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceSpecWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceSpec object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetFlinkBlobStorage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceSpec) GetFlinkBlobStorage() ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1FlinkBlobStorage`

GetFlinkBlobStorage returns the FlinkBlobStorage field if non-nil, zero value otherwise.

### GetFlinkBlobStorageOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceSpec) GetFlinkBlobStorageOk() (*ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1FlinkBlobStorage, bool)`

GetFlinkBlobStorageOk returns a tuple with the FlinkBlobStorage field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFlinkBlobStorage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceSpec) SetFlinkBlobStorage(v ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1FlinkBlobStorage)`

SetFlinkBlobStorage sets FlinkBlobStorage field to given value.

### HasFlinkBlobStorage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceSpec) HasFlinkBlobStorage() bool`

HasFlinkBlobStorage returns a boolean if a field has been set.

### GetPoolRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceSpec) GetPoolRef() ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PoolRef`

GetPoolRef returns the PoolRef field if non-nil, zero value otherwise.

### GetPoolRefOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceSpec) GetPoolRefOk() (*ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PoolRef, bool)`

GetPoolRefOk returns a tuple with the PoolRef field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPoolRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceSpec) SetPoolRef(v ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PoolRef)`

SetPoolRef sets PoolRef field to given value.


### GetPulsarClusterNames

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceSpec) GetPulsarClusterNames() []string`

GetPulsarClusterNames returns the PulsarClusterNames field if non-nil, zero value otherwise.

### GetPulsarClusterNamesOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceSpec) GetPulsarClusterNamesOk() (*[]string, bool)`

GetPulsarClusterNamesOk returns a tuple with the PulsarClusterNames field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPulsarClusterNames

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceSpec) SetPulsarClusterNames(v []string)`

SetPulsarClusterNames sets PulsarClusterNames field to given value.


### GetUseExternalAccess

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceSpec) GetUseExternalAccess() bool`

GetUseExternalAccess returns the UseExternalAccess field if non-nil, zero value otherwise.

### GetUseExternalAccessOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceSpec) GetUseExternalAccessOk() (*bool, bool)`

GetUseExternalAccessOk returns a tuple with the UseExternalAccess field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUseExternalAccess

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceSpec) SetUseExternalAccess(v bool)`

SetUseExternalAccess sets UseExternalAccess field to given value.

### HasUseExternalAccess

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1WorkspaceSpec) HasUseExternalAccess() bool`

HasUseExternalAccess returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


