# ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Artifact** | [**ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Artifact**](ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Artifact.md) |  | 
**FlinkConfiguration** | Pointer to **map[string]string** |  | [optional] 
**Kubernetes** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpecKubernetesSpec**](ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpecKubernetesSpec.md) |  | [optional] 
**LatestCheckpointFetchInterval** | Pointer to **int32** |  | [optional] 
**Logging** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Logging**](ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Logging.md) |  | [optional] 
**NumberOfTaskManagers** | Pointer to **int32** |  | [optional] 
**Parallelism** | Pointer to **int32** |  | [optional] 
**Resources** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentKubernetesResources**](ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentKubernetesResources.md) |  | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpec

`func NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpec(artifact ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Artifact, ) *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpec`

NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpec instantiates a new ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpec object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpecWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpecWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpec`

NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpecWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpec object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetArtifact

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpec) GetArtifact() ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Artifact`

GetArtifact returns the Artifact field if non-nil, zero value otherwise.

### GetArtifactOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpec) GetArtifactOk() (*ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Artifact, bool)`

GetArtifactOk returns a tuple with the Artifact field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetArtifact

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpec) SetArtifact(v ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Artifact)`

SetArtifact sets Artifact field to given value.


### GetFlinkConfiguration

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpec) GetFlinkConfiguration() map[string]string`

GetFlinkConfiguration returns the FlinkConfiguration field if non-nil, zero value otherwise.

### GetFlinkConfigurationOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpec) GetFlinkConfigurationOk() (*map[string]string, bool)`

GetFlinkConfigurationOk returns a tuple with the FlinkConfiguration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFlinkConfiguration

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpec) SetFlinkConfiguration(v map[string]string)`

SetFlinkConfiguration sets FlinkConfiguration field to given value.

### HasFlinkConfiguration

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpec) HasFlinkConfiguration() bool`

HasFlinkConfiguration returns a boolean if a field has been set.

### GetKubernetes

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpec) GetKubernetes() ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpecKubernetesSpec`

GetKubernetes returns the Kubernetes field if non-nil, zero value otherwise.

### GetKubernetesOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpec) GetKubernetesOk() (*ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpecKubernetesSpec, bool)`

GetKubernetesOk returns a tuple with the Kubernetes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKubernetes

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpec) SetKubernetes(v ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpecKubernetesSpec)`

SetKubernetes sets Kubernetes field to given value.

### HasKubernetes

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpec) HasKubernetes() bool`

HasKubernetes returns a boolean if a field has been set.

### GetLatestCheckpointFetchInterval

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpec) GetLatestCheckpointFetchInterval() int32`

GetLatestCheckpointFetchInterval returns the LatestCheckpointFetchInterval field if non-nil, zero value otherwise.

### GetLatestCheckpointFetchIntervalOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpec) GetLatestCheckpointFetchIntervalOk() (*int32, bool)`

GetLatestCheckpointFetchIntervalOk returns a tuple with the LatestCheckpointFetchInterval field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLatestCheckpointFetchInterval

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpec) SetLatestCheckpointFetchInterval(v int32)`

SetLatestCheckpointFetchInterval sets LatestCheckpointFetchInterval field to given value.

### HasLatestCheckpointFetchInterval

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpec) HasLatestCheckpointFetchInterval() bool`

HasLatestCheckpointFetchInterval returns a boolean if a field has been set.

### GetLogging

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpec) GetLogging() ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Logging`

GetLogging returns the Logging field if non-nil, zero value otherwise.

### GetLoggingOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpec) GetLoggingOk() (*ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Logging, bool)`

GetLoggingOk returns a tuple with the Logging field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLogging

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpec) SetLogging(v ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Logging)`

SetLogging sets Logging field to given value.

### HasLogging

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpec) HasLogging() bool`

HasLogging returns a boolean if a field has been set.

### GetNumberOfTaskManagers

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpec) GetNumberOfTaskManagers() int32`

GetNumberOfTaskManagers returns the NumberOfTaskManagers field if non-nil, zero value otherwise.

### GetNumberOfTaskManagersOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpec) GetNumberOfTaskManagersOk() (*int32, bool)`

GetNumberOfTaskManagersOk returns a tuple with the NumberOfTaskManagers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNumberOfTaskManagers

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpec) SetNumberOfTaskManagers(v int32)`

SetNumberOfTaskManagers sets NumberOfTaskManagers field to given value.

### HasNumberOfTaskManagers

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpec) HasNumberOfTaskManagers() bool`

HasNumberOfTaskManagers returns a boolean if a field has been set.

### GetParallelism

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpec) GetParallelism() int32`

GetParallelism returns the Parallelism field if non-nil, zero value otherwise.

### GetParallelismOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpec) GetParallelismOk() (*int32, bool)`

GetParallelismOk returns a tuple with the Parallelism field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetParallelism

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpec) SetParallelism(v int32)`

SetParallelism sets Parallelism field to given value.

### HasParallelism

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpec) HasParallelism() bool`

HasParallelism returns a boolean if a field has been set.

### GetResources

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpec) GetResources() ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentKubernetesResources`

GetResources returns the Resources field if non-nil, zero value otherwise.

### GetResourcesOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpec) GetResourcesOk() (*ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentKubernetesResources, bool)`

GetResourcesOk returns a tuple with the Resources field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResources

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpec) SetResources(v ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentKubernetesResources)`

SetResources sets Resources field to given value.

### HasResources

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplateSpec) HasResources() bool`

HasResources returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


