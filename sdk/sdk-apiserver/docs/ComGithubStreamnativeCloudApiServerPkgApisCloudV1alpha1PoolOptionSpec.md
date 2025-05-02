# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolOptionSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CloudType** | **string** |  | 
**DeploymentType** | **string** |  | 
**Features** | Pointer to **map[string]bool** |  | [optional] 
**Locations** | [**[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolOptionLocation**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolOptionLocation.md) |  | 
**PoolRef** | [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolRef**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolRef.md) |  | 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolOptionSpec

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolOptionSpec(cloudType string, deploymentType string, locations []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolOptionLocation, poolRef ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolRef, ) *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolOptionSpec`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolOptionSpec instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolOptionSpec object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolOptionSpecWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolOptionSpecWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolOptionSpec`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolOptionSpecWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolOptionSpec object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCloudType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolOptionSpec) GetCloudType() string`

GetCloudType returns the CloudType field if non-nil, zero value otherwise.

### GetCloudTypeOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolOptionSpec) GetCloudTypeOk() (*string, bool)`

GetCloudTypeOk returns a tuple with the CloudType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCloudType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolOptionSpec) SetCloudType(v string)`

SetCloudType sets CloudType field to given value.


### GetDeploymentType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolOptionSpec) GetDeploymentType() string`

GetDeploymentType returns the DeploymentType field if non-nil, zero value otherwise.

### GetDeploymentTypeOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolOptionSpec) GetDeploymentTypeOk() (*string, bool)`

GetDeploymentTypeOk returns a tuple with the DeploymentType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeploymentType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolOptionSpec) SetDeploymentType(v string)`

SetDeploymentType sets DeploymentType field to given value.


### GetFeatures

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolOptionSpec) GetFeatures() map[string]bool`

GetFeatures returns the Features field if non-nil, zero value otherwise.

### GetFeaturesOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolOptionSpec) GetFeaturesOk() (*map[string]bool, bool)`

GetFeaturesOk returns a tuple with the Features field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFeatures

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolOptionSpec) SetFeatures(v map[string]bool)`

SetFeatures sets Features field to given value.

### HasFeatures

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolOptionSpec) HasFeatures() bool`

HasFeatures returns a boolean if a field has been set.

### GetLocations

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolOptionSpec) GetLocations() []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolOptionLocation`

GetLocations returns the Locations field if non-nil, zero value otherwise.

### GetLocationsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolOptionSpec) GetLocationsOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolOptionLocation, bool)`

GetLocationsOk returns a tuple with the Locations field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLocations

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolOptionSpec) SetLocations(v []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolOptionLocation)`

SetLocations sets Locations field to given value.


### GetPoolRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolOptionSpec) GetPoolRef() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolRef`

GetPoolRef returns the PoolRef field if non-nil, zero value otherwise.

### GetPoolRefOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolOptionSpec) GetPoolRefOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolRef, bool)`

GetPoolRefOk returns a tuple with the PoolRef field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPoolRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolOptionSpec) SetPoolRef(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolRef)`

SetPoolRef sets PoolRef field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


