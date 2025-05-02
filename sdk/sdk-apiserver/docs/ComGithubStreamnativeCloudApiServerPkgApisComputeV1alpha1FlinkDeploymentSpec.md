# ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1FlinkDeploymentSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Annotations** | Pointer to **map[string]string** |  | [optional] 
**CommunityTemplate** | Pointer to **map[string]interface{}** | CommunityDeploymentTemplate defines the desired state of CommunityDeployment | [optional] 
**DefaultPulsarCluster** | Pointer to **string** | DefaultPulsarCluster is the default pulsar cluster to use. If not provided, the controller will use the first pulsar cluster from the workspace. | [optional] 
**Labels** | Pointer to **map[string]string** |  | [optional] 
**PoolMemberRef** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PoolMemberReference**](ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PoolMemberReference.md) |  | [optional] 
**Template** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentTemplate**](ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentTemplate.md) |  | [optional] 
**WorkspaceName** | **string** | WorkspaceName is the reference to the workspace, and is required | 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1FlinkDeploymentSpec

`func NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1FlinkDeploymentSpec(workspaceName string, ) *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1FlinkDeploymentSpec`

NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1FlinkDeploymentSpec instantiates a new ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1FlinkDeploymentSpec object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1FlinkDeploymentSpecWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1FlinkDeploymentSpecWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1FlinkDeploymentSpec`

NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1FlinkDeploymentSpecWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1FlinkDeploymentSpec object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAnnotations

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1FlinkDeploymentSpec) GetAnnotations() map[string]string`

GetAnnotations returns the Annotations field if non-nil, zero value otherwise.

### GetAnnotationsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1FlinkDeploymentSpec) GetAnnotationsOk() (*map[string]string, bool)`

GetAnnotationsOk returns a tuple with the Annotations field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAnnotations

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1FlinkDeploymentSpec) SetAnnotations(v map[string]string)`

SetAnnotations sets Annotations field to given value.

### HasAnnotations

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1FlinkDeploymentSpec) HasAnnotations() bool`

HasAnnotations returns a boolean if a field has been set.

### GetCommunityTemplate

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1FlinkDeploymentSpec) GetCommunityTemplate() map[string]interface{}`

GetCommunityTemplate returns the CommunityTemplate field if non-nil, zero value otherwise.

### GetCommunityTemplateOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1FlinkDeploymentSpec) GetCommunityTemplateOk() (*map[string]interface{}, bool)`

GetCommunityTemplateOk returns a tuple with the CommunityTemplate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCommunityTemplate

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1FlinkDeploymentSpec) SetCommunityTemplate(v map[string]interface{})`

SetCommunityTemplate sets CommunityTemplate field to given value.

### HasCommunityTemplate

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1FlinkDeploymentSpec) HasCommunityTemplate() bool`

HasCommunityTemplate returns a boolean if a field has been set.

### GetDefaultPulsarCluster

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1FlinkDeploymentSpec) GetDefaultPulsarCluster() string`

GetDefaultPulsarCluster returns the DefaultPulsarCluster field if non-nil, zero value otherwise.

### GetDefaultPulsarClusterOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1FlinkDeploymentSpec) GetDefaultPulsarClusterOk() (*string, bool)`

GetDefaultPulsarClusterOk returns a tuple with the DefaultPulsarCluster field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDefaultPulsarCluster

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1FlinkDeploymentSpec) SetDefaultPulsarCluster(v string)`

SetDefaultPulsarCluster sets DefaultPulsarCluster field to given value.

### HasDefaultPulsarCluster

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1FlinkDeploymentSpec) HasDefaultPulsarCluster() bool`

HasDefaultPulsarCluster returns a boolean if a field has been set.

### GetLabels

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1FlinkDeploymentSpec) GetLabels() map[string]string`

GetLabels returns the Labels field if non-nil, zero value otherwise.

### GetLabelsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1FlinkDeploymentSpec) GetLabelsOk() (*map[string]string, bool)`

GetLabelsOk returns a tuple with the Labels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabels

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1FlinkDeploymentSpec) SetLabels(v map[string]string)`

SetLabels sets Labels field to given value.

### HasLabels

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1FlinkDeploymentSpec) HasLabels() bool`

HasLabels returns a boolean if a field has been set.

### GetPoolMemberRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1FlinkDeploymentSpec) GetPoolMemberRef() ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PoolMemberReference`

GetPoolMemberRef returns the PoolMemberRef field if non-nil, zero value otherwise.

### GetPoolMemberRefOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1FlinkDeploymentSpec) GetPoolMemberRefOk() (*ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PoolMemberReference, bool)`

GetPoolMemberRefOk returns a tuple with the PoolMemberRef field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPoolMemberRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1FlinkDeploymentSpec) SetPoolMemberRef(v ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PoolMemberReference)`

SetPoolMemberRef sets PoolMemberRef field to given value.

### HasPoolMemberRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1FlinkDeploymentSpec) HasPoolMemberRef() bool`

HasPoolMemberRef returns a boolean if a field has been set.

### GetTemplate

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1FlinkDeploymentSpec) GetTemplate() ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentTemplate`

GetTemplate returns the Template field if non-nil, zero value otherwise.

### GetTemplateOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1FlinkDeploymentSpec) GetTemplateOk() (*ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentTemplate, bool)`

GetTemplateOk returns a tuple with the Template field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTemplate

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1FlinkDeploymentSpec) SetTemplate(v ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentTemplate)`

SetTemplate sets Template field to given value.

### HasTemplate

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1FlinkDeploymentSpec) HasTemplate() bool`

HasTemplate returns a boolean if a field has been set.

### GetWorkspaceName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1FlinkDeploymentSpec) GetWorkspaceName() string`

GetWorkspaceName returns the WorkspaceName field if non-nil, zero value otherwise.

### GetWorkspaceNameOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1FlinkDeploymentSpec) GetWorkspaceNameOk() (*string, bool)`

GetWorkspaceNameOk returns a tuple with the WorkspaceName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkspaceName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1FlinkDeploymentSpec) SetWorkspaceName(v string)`

SetWorkspaceName sets WorkspaceName field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


