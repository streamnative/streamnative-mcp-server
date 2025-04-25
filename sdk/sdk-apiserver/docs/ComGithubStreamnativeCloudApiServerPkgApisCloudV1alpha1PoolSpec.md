# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CloudFeature** | Pointer to **map[string]bool** | CloudFeature indicates features this pool wants to enable/disable by default for all Pulsar clusters created on it | [optional] 
**DeploymentType** | Pointer to **string** | This feild is used by &#x60;cloud-manager&#x60; and &#x60;cloud-billing-reporter&#x60; to potentially charge different rates for our customers. It is imperative that we correctly set this field if a pool is a \&quot;Pro\&quot; tier or no tier. | [optional] 
**Gcloud** | Pointer to **map[string]interface{}** |  | [optional] 
**Sharing** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SharingConfig**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SharingConfig.md) |  | [optional] 
**Type** | **string** |  | 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolSpec

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolSpec(type_ string, ) *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolSpec`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolSpec instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolSpec object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolSpecWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolSpecWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolSpec`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolSpecWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolSpec object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCloudFeature

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolSpec) GetCloudFeature() map[string]bool`

GetCloudFeature returns the CloudFeature field if non-nil, zero value otherwise.

### GetCloudFeatureOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolSpec) GetCloudFeatureOk() (*map[string]bool, bool)`

GetCloudFeatureOk returns a tuple with the CloudFeature field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCloudFeature

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolSpec) SetCloudFeature(v map[string]bool)`

SetCloudFeature sets CloudFeature field to given value.

### HasCloudFeature

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolSpec) HasCloudFeature() bool`

HasCloudFeature returns a boolean if a field has been set.

### GetDeploymentType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolSpec) GetDeploymentType() string`

GetDeploymentType returns the DeploymentType field if non-nil, zero value otherwise.

### GetDeploymentTypeOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolSpec) GetDeploymentTypeOk() (*string, bool)`

GetDeploymentTypeOk returns a tuple with the DeploymentType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeploymentType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolSpec) SetDeploymentType(v string)`

SetDeploymentType sets DeploymentType field to given value.

### HasDeploymentType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolSpec) HasDeploymentType() bool`

HasDeploymentType returns a boolean if a field has been set.

### GetGcloud

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolSpec) GetGcloud() map[string]interface{}`

GetGcloud returns the Gcloud field if non-nil, zero value otherwise.

### GetGcloudOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolSpec) GetGcloudOk() (*map[string]interface{}, bool)`

GetGcloudOk returns a tuple with the Gcloud field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGcloud

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolSpec) SetGcloud(v map[string]interface{})`

SetGcloud sets Gcloud field to given value.

### HasGcloud

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolSpec) HasGcloud() bool`

HasGcloud returns a boolean if a field has been set.

### GetSharing

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolSpec) GetSharing() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SharingConfig`

GetSharing returns the Sharing field if non-nil, zero value otherwise.

### GetSharingOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolSpec) GetSharingOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SharingConfig, bool)`

GetSharingOk returns a tuple with the Sharing field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSharing

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolSpec) SetSharing(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SharingConfig)`

SetSharing sets Sharing field to given value.

### HasSharing

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolSpec) HasSharing() bool`

HasSharing returns a boolean if a field has been set.

### GetType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolSpec) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolSpec) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolSpec) SetType(v string)`

SetType sets Type field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


