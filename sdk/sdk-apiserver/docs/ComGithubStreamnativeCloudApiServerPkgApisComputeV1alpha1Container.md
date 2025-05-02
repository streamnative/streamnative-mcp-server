# ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Args** | Pointer to **[]string** | Arguments to the entrypoint. | [optional] 
**Command** | Pointer to **[]string** | Entrypoint array. Not executed within a shell. | [optional] 
**Env** | Pointer to [**[]V1EnvVar**](V1EnvVar.md) | List of environment variables to set in the container. | [optional] 
**EnvFrom** | Pointer to [**[]V1EnvFromSource**](V1EnvFromSource.md) | List of sources to populate environment variables in the container. | [optional] 
**Image** | Pointer to **string** | Docker image name. | [optional] 
**ImagePullPolicy** | Pointer to **string** | Image pull policy.  Possible enum values:  - &#x60;\&quot;Always\&quot;&#x60; means that kubelet always attempts to pull the latest image. Container will fail If the pull fails.  - &#x60;\&quot;IfNotPresent\&quot;&#x60; means that kubelet pulls if the image isn&#39;t present on disk. Container will fail if the image isn&#39;t present and the pull fails.  - &#x60;\&quot;Never\&quot;&#x60; means that kubelet never pulls an image, but only uses a local image. Container will fail if the image isn&#39;t present | [optional] 
**LivenessProbe** | Pointer to [**V1Probe**](V1Probe.md) |  | [optional] 
**Name** | **string** | Name of the container specified as a DNS_LABEL. Each container in a pod must have a unique name (DNS_LABEL). | 
**ReadinessProbe** | Pointer to [**V1Probe**](V1Probe.md) |  | [optional] 
**Resources** | Pointer to [**V1ResourceRequirements**](V1ResourceRequirements.md) |  | [optional] 
**SecurityContext** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext**](ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext.md) |  | [optional] 
**StartupProbe** | Pointer to [**V1Probe**](V1Probe.md) |  | [optional] 
**VolumeMounts** | Pointer to [**[]V1VolumeMount**](V1VolumeMount.md) | Pod volumes to mount into the container&#39;s filesystem. | [optional] 
**WorkingDir** | Pointer to **string** | Container&#39;s working directory. | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container

`func NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container(name string, ) *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container`

NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container instantiates a new ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1ContainerWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1ContainerWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container`

NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1ContainerWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetArgs

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) GetArgs() []string`

GetArgs returns the Args field if non-nil, zero value otherwise.

### GetArgsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) GetArgsOk() (*[]string, bool)`

GetArgsOk returns a tuple with the Args field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetArgs

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) SetArgs(v []string)`

SetArgs sets Args field to given value.

### HasArgs

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) HasArgs() bool`

HasArgs returns a boolean if a field has been set.

### GetCommand

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) GetCommand() []string`

GetCommand returns the Command field if non-nil, zero value otherwise.

### GetCommandOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) GetCommandOk() (*[]string, bool)`

GetCommandOk returns a tuple with the Command field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCommand

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) SetCommand(v []string)`

SetCommand sets Command field to given value.

### HasCommand

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) HasCommand() bool`

HasCommand returns a boolean if a field has been set.

### GetEnv

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) GetEnv() []V1EnvVar`

GetEnv returns the Env field if non-nil, zero value otherwise.

### GetEnvOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) GetEnvOk() (*[]V1EnvVar, bool)`

GetEnvOk returns a tuple with the Env field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnv

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) SetEnv(v []V1EnvVar)`

SetEnv sets Env field to given value.

### HasEnv

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) HasEnv() bool`

HasEnv returns a boolean if a field has been set.

### GetEnvFrom

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) GetEnvFrom() []V1EnvFromSource`

GetEnvFrom returns the EnvFrom field if non-nil, zero value otherwise.

### GetEnvFromOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) GetEnvFromOk() (*[]V1EnvFromSource, bool)`

GetEnvFromOk returns a tuple with the EnvFrom field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnvFrom

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) SetEnvFrom(v []V1EnvFromSource)`

SetEnvFrom sets EnvFrom field to given value.

### HasEnvFrom

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) HasEnvFrom() bool`

HasEnvFrom returns a boolean if a field has been set.

### GetImage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) GetImage() string`

GetImage returns the Image field if non-nil, zero value otherwise.

### GetImageOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) GetImageOk() (*string, bool)`

GetImageOk returns a tuple with the Image field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) SetImage(v string)`

SetImage sets Image field to given value.

### HasImage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) HasImage() bool`

HasImage returns a boolean if a field has been set.

### GetImagePullPolicy

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) GetImagePullPolicy() string`

GetImagePullPolicy returns the ImagePullPolicy field if non-nil, zero value otherwise.

### GetImagePullPolicyOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) GetImagePullPolicyOk() (*string, bool)`

GetImagePullPolicyOk returns a tuple with the ImagePullPolicy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImagePullPolicy

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) SetImagePullPolicy(v string)`

SetImagePullPolicy sets ImagePullPolicy field to given value.

### HasImagePullPolicy

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) HasImagePullPolicy() bool`

HasImagePullPolicy returns a boolean if a field has been set.

### GetLivenessProbe

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) GetLivenessProbe() V1Probe`

GetLivenessProbe returns the LivenessProbe field if non-nil, zero value otherwise.

### GetLivenessProbeOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) GetLivenessProbeOk() (*V1Probe, bool)`

GetLivenessProbeOk returns a tuple with the LivenessProbe field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLivenessProbe

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) SetLivenessProbe(v V1Probe)`

SetLivenessProbe sets LivenessProbe field to given value.

### HasLivenessProbe

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) HasLivenessProbe() bool`

HasLivenessProbe returns a boolean if a field has been set.

### GetName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) SetName(v string)`

SetName sets Name field to given value.


### GetReadinessProbe

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) GetReadinessProbe() V1Probe`

GetReadinessProbe returns the ReadinessProbe field if non-nil, zero value otherwise.

### GetReadinessProbeOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) GetReadinessProbeOk() (*V1Probe, bool)`

GetReadinessProbeOk returns a tuple with the ReadinessProbe field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReadinessProbe

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) SetReadinessProbe(v V1Probe)`

SetReadinessProbe sets ReadinessProbe field to given value.

### HasReadinessProbe

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) HasReadinessProbe() bool`

HasReadinessProbe returns a boolean if a field has been set.

### GetResources

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) GetResources() V1ResourceRequirements`

GetResources returns the Resources field if non-nil, zero value otherwise.

### GetResourcesOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) GetResourcesOk() (*V1ResourceRequirements, bool)`

GetResourcesOk returns a tuple with the Resources field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResources

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) SetResources(v V1ResourceRequirements)`

SetResources sets Resources field to given value.

### HasResources

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) HasResources() bool`

HasResources returns a boolean if a field has been set.

### GetSecurityContext

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) GetSecurityContext() ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext`

GetSecurityContext returns the SecurityContext field if non-nil, zero value otherwise.

### GetSecurityContextOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) GetSecurityContextOk() (*ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext, bool)`

GetSecurityContextOk returns a tuple with the SecurityContext field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSecurityContext

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) SetSecurityContext(v ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext)`

SetSecurityContext sets SecurityContext field to given value.

### HasSecurityContext

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) HasSecurityContext() bool`

HasSecurityContext returns a boolean if a field has been set.

### GetStartupProbe

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) GetStartupProbe() V1Probe`

GetStartupProbe returns the StartupProbe field if non-nil, zero value otherwise.

### GetStartupProbeOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) GetStartupProbeOk() (*V1Probe, bool)`

GetStartupProbeOk returns a tuple with the StartupProbe field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStartupProbe

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) SetStartupProbe(v V1Probe)`

SetStartupProbe sets StartupProbe field to given value.

### HasStartupProbe

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) HasStartupProbe() bool`

HasStartupProbe returns a boolean if a field has been set.

### GetVolumeMounts

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) GetVolumeMounts() []V1VolumeMount`

GetVolumeMounts returns the VolumeMounts field if non-nil, zero value otherwise.

### GetVolumeMountsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) GetVolumeMountsOk() (*[]V1VolumeMount, bool)`

GetVolumeMountsOk returns a tuple with the VolumeMounts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVolumeMounts

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) SetVolumeMounts(v []V1VolumeMount)`

SetVolumeMounts sets VolumeMounts field to given value.

### HasVolumeMounts

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) HasVolumeMounts() bool`

HasVolumeMounts returns a boolean if a field has been set.

### GetWorkingDir

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) GetWorkingDir() string`

GetWorkingDir returns the WorkingDir field if non-nil, zero value otherwise.

### GetWorkingDirOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) GetWorkingDirOk() (*string, bool)`

GetWorkingDirOk returns a tuple with the WorkingDir field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkingDir

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) SetWorkingDir(v string)`

SetWorkingDir sets WorkingDir field to given value.

### HasWorkingDir

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container) HasWorkingDir() bool`

HasWorkingDir returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


