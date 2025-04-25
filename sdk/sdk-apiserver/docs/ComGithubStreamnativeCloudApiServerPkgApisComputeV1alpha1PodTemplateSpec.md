# ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Affinity** | Pointer to [**V1Affinity**](V1Affinity.md) |  | [optional] 
**Containers** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container**](ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container.md) |  | [optional] 
**ImagePullSecrets** | Pointer to [**[]V1LocalObjectReference**](V1LocalObjectReference.md) |  | [optional] 
**InitContainers** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container**](ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container.md) |  | [optional] 
**NodeSelector** | Pointer to **map[string]string** |  | [optional] 
**SecurityContext** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext**](ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext.md) |  | [optional] 
**ServiceAccountName** | Pointer to **string** |  | [optional] 
**ShareProcessNamespace** | Pointer to **bool** |  | [optional] 
**Tolerations** | Pointer to [**[]V1Toleration**](V1Toleration.md) |  | [optional] 
**Volumes** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Volume**](ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Volume.md) |  | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec

`func NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec() *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec`

NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec instantiates a new ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpecWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpecWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec`

NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpecWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAffinity

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec) GetAffinity() V1Affinity`

GetAffinity returns the Affinity field if non-nil, zero value otherwise.

### GetAffinityOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec) GetAffinityOk() (*V1Affinity, bool)`

GetAffinityOk returns a tuple with the Affinity field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAffinity

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec) SetAffinity(v V1Affinity)`

SetAffinity sets Affinity field to given value.

### HasAffinity

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec) HasAffinity() bool`

HasAffinity returns a boolean if a field has been set.

### GetContainers

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec) GetContainers() []ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container`

GetContainers returns the Containers field if non-nil, zero value otherwise.

### GetContainersOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec) GetContainersOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container, bool)`

GetContainersOk returns a tuple with the Containers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContainers

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec) SetContainers(v []ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container)`

SetContainers sets Containers field to given value.

### HasContainers

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec) HasContainers() bool`

HasContainers returns a boolean if a field has been set.

### GetImagePullSecrets

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec) GetImagePullSecrets() []V1LocalObjectReference`

GetImagePullSecrets returns the ImagePullSecrets field if non-nil, zero value otherwise.

### GetImagePullSecretsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec) GetImagePullSecretsOk() (*[]V1LocalObjectReference, bool)`

GetImagePullSecretsOk returns a tuple with the ImagePullSecrets field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImagePullSecrets

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec) SetImagePullSecrets(v []V1LocalObjectReference)`

SetImagePullSecrets sets ImagePullSecrets field to given value.

### HasImagePullSecrets

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec) HasImagePullSecrets() bool`

HasImagePullSecrets returns a boolean if a field has been set.

### GetInitContainers

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec) GetInitContainers() []ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container`

GetInitContainers returns the InitContainers field if non-nil, zero value otherwise.

### GetInitContainersOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec) GetInitContainersOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container, bool)`

GetInitContainersOk returns a tuple with the InitContainers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInitContainers

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec) SetInitContainers(v []ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Container)`

SetInitContainers sets InitContainers field to given value.

### HasInitContainers

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec) HasInitContainers() bool`

HasInitContainers returns a boolean if a field has been set.

### GetNodeSelector

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec) GetNodeSelector() map[string]string`

GetNodeSelector returns the NodeSelector field if non-nil, zero value otherwise.

### GetNodeSelectorOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec) GetNodeSelectorOk() (*map[string]string, bool)`

GetNodeSelectorOk returns a tuple with the NodeSelector field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNodeSelector

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec) SetNodeSelector(v map[string]string)`

SetNodeSelector sets NodeSelector field to given value.

### HasNodeSelector

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec) HasNodeSelector() bool`

HasNodeSelector returns a boolean if a field has been set.

### GetSecurityContext

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec) GetSecurityContext() ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext`

GetSecurityContext returns the SecurityContext field if non-nil, zero value otherwise.

### GetSecurityContextOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec) GetSecurityContextOk() (*ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext, bool)`

GetSecurityContextOk returns a tuple with the SecurityContext field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSecurityContext

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec) SetSecurityContext(v ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1SecurityContext)`

SetSecurityContext sets SecurityContext field to given value.

### HasSecurityContext

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec) HasSecurityContext() bool`

HasSecurityContext returns a boolean if a field has been set.

### GetServiceAccountName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec) GetServiceAccountName() string`

GetServiceAccountName returns the ServiceAccountName field if non-nil, zero value otherwise.

### GetServiceAccountNameOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec) GetServiceAccountNameOk() (*string, bool)`

GetServiceAccountNameOk returns a tuple with the ServiceAccountName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServiceAccountName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec) SetServiceAccountName(v string)`

SetServiceAccountName sets ServiceAccountName field to given value.

### HasServiceAccountName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec) HasServiceAccountName() bool`

HasServiceAccountName returns a boolean if a field has been set.

### GetShareProcessNamespace

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec) GetShareProcessNamespace() bool`

GetShareProcessNamespace returns the ShareProcessNamespace field if non-nil, zero value otherwise.

### GetShareProcessNamespaceOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec) GetShareProcessNamespaceOk() (*bool, bool)`

GetShareProcessNamespaceOk returns a tuple with the ShareProcessNamespace field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetShareProcessNamespace

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec) SetShareProcessNamespace(v bool)`

SetShareProcessNamespace sets ShareProcessNamespace field to given value.

### HasShareProcessNamespace

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec) HasShareProcessNamespace() bool`

HasShareProcessNamespace returns a boolean if a field has been set.

### GetTolerations

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec) GetTolerations() []V1Toleration`

GetTolerations returns the Tolerations field if non-nil, zero value otherwise.

### GetTolerationsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec) GetTolerationsOk() (*[]V1Toleration, bool)`

GetTolerationsOk returns a tuple with the Tolerations field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTolerations

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec) SetTolerations(v []V1Toleration)`

SetTolerations sets Tolerations field to given value.

### HasTolerations

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec) HasTolerations() bool`

HasTolerations returns a boolean if a field has been set.

### GetVolumes

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec) GetVolumes() []ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Volume`

GetVolumes returns the Volumes field if non-nil, zero value otherwise.

### GetVolumesOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec) GetVolumesOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Volume, bool)`

GetVolumesOk returns a tuple with the Volumes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVolumes

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec) SetVolumes(v []ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Volume)`

SetVolumes sets Volumes field to given value.

### HasVolumes

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1PodTemplateSpec) HasVolumes() bool`

HasVolumes returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


