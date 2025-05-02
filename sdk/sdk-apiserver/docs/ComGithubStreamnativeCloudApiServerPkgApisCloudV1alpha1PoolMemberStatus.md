# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Conditions** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition.md) | Conditions is an array of current observed conditions. | [optional] 
**DeploymentType** | **string** |  | 
**InstalledCSVs** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstalledCSV**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstalledCSV.md) | InstalledCSVs shows the name and status of installed operator versions | [optional] 
**ObservedGeneration** | **int64** | ObservedGeneration is the most recent generation observed by the PoolMember controller. | 
**ServerVersion** | Pointer to **string** |  | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberStatus

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberStatus(deploymentType string, observedGeneration int64, ) *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberStatus`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberStatus instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberStatus object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberStatusWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberStatusWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberStatus`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberStatusWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberStatus object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberStatus) GetConditions() []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition`

GetConditions returns the Conditions field if non-nil, zero value otherwise.

### GetConditionsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberStatus) GetConditionsOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition, bool)`

GetConditionsOk returns a tuple with the Conditions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberStatus) SetConditions(v []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition)`

SetConditions sets Conditions field to given value.

### HasConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberStatus) HasConditions() bool`

HasConditions returns a boolean if a field has been set.

### GetDeploymentType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberStatus) GetDeploymentType() string`

GetDeploymentType returns the DeploymentType field if non-nil, zero value otherwise.

### GetDeploymentTypeOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberStatus) GetDeploymentTypeOk() (*string, bool)`

GetDeploymentTypeOk returns a tuple with the DeploymentType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeploymentType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberStatus) SetDeploymentType(v string)`

SetDeploymentType sets DeploymentType field to given value.


### GetInstalledCSVs

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberStatus) GetInstalledCSVs() []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstalledCSV`

GetInstalledCSVs returns the InstalledCSVs field if non-nil, zero value otherwise.

### GetInstalledCSVsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberStatus) GetInstalledCSVsOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstalledCSV, bool)`

GetInstalledCSVsOk returns a tuple with the InstalledCSVs field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInstalledCSVs

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberStatus) SetInstalledCSVs(v []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstalledCSV)`

SetInstalledCSVs sets InstalledCSVs field to given value.

### HasInstalledCSVs

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberStatus) HasInstalledCSVs() bool`

HasInstalledCSVs returns a boolean if a field has been set.

### GetObservedGeneration

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberStatus) GetObservedGeneration() int64`

GetObservedGeneration returns the ObservedGeneration field if non-nil, zero value otherwise.

### GetObservedGenerationOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberStatus) GetObservedGenerationOk() (*int64, bool)`

GetObservedGenerationOk returns a tuple with the ObservedGeneration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObservedGeneration

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberStatus) SetObservedGeneration(v int64)`

SetObservedGeneration sets ObservedGeneration field to given value.


### GetServerVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberStatus) GetServerVersion() string`

GetServerVersion returns the ServerVersion field if non-nil, zero value otherwise.

### GetServerVersionOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberStatus) GetServerVersionOk() (*string, bool)`

GetServerVersionOk returns a tuple with the ServerVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServerVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberStatus) SetServerVersion(v string)`

SetServerVersion sets ServerVersion field to given value.

### HasServerVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberStatus) HasServerVersion() bool`

HasServerVersion returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


