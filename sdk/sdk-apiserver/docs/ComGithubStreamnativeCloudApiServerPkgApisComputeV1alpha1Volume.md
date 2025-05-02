# ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Volume

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ConfigMap** | Pointer to [**V1ConfigMapVolumeSource**](V1ConfigMapVolumeSource.md) |  | [optional] 
**Name** | **string** | Volume&#39;s name. Must be a DNS_LABEL and unique within the pod. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names | 
**Secret** | Pointer to [**V1SecretVolumeSource**](V1SecretVolumeSource.md) |  | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Volume

`func NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Volume(name string, ) *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Volume`

NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Volume instantiates a new ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Volume object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VolumeWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VolumeWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Volume`

NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VolumeWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Volume object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetConfigMap

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Volume) GetConfigMap() V1ConfigMapVolumeSource`

GetConfigMap returns the ConfigMap field if non-nil, zero value otherwise.

### GetConfigMapOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Volume) GetConfigMapOk() (*V1ConfigMapVolumeSource, bool)`

GetConfigMapOk returns a tuple with the ConfigMap field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConfigMap

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Volume) SetConfigMap(v V1ConfigMapVolumeSource)`

SetConfigMap sets ConfigMap field to given value.

### HasConfigMap

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Volume) HasConfigMap() bool`

HasConfigMap returns a boolean if a field has been set.

### GetName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Volume) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Volume) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Volume) SetName(v string)`

SetName sets Name field to given value.


### GetSecret

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Volume) GetSecret() V1SecretVolumeSource`

GetSecret returns the Secret field if non-nil, zero value otherwise.

### GetSecretOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Volume) GetSecretOk() (*V1SecretVolumeSource, bool)`

GetSecretOk returns a tuple with the Secret field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSecret

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Volume) SetSecret(v V1SecretVolumeSource)`

SetSecret sets Secret field to given value.

### HasSecret

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1Volume) HasSecret() bool`

HasSecret returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


