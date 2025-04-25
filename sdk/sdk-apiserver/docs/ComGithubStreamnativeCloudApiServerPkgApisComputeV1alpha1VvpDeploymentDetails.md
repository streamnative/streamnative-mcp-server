# ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetails

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**DeploymentTargetName** | Pointer to **string** |  | [optional] 
**JobFailureExpirationTime** | Pointer to **string** |  | [optional] 
**MaxJobCreationAttempts** | Pointer to **int32** |  | [optional] 
**MaxSavepointCreationAttempts** | Pointer to **int32** |  | [optional] 
**RestoreStrategy** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpRestoreStrategy**](ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpRestoreStrategy.md) |  | [optional] 
**SessionClusterName** | Pointer to **string** |  | [optional] 
**State** | Pointer to **string** |  | [optional] 
**Template** | [**ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplate**](ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplate.md) |  | 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetails

`func NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetails(template ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplate, ) *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetails`

NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetails instantiates a new ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetails object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetails`

NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetails object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDeploymentTargetName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetails) GetDeploymentTargetName() string`

GetDeploymentTargetName returns the DeploymentTargetName field if non-nil, zero value otherwise.

### GetDeploymentTargetNameOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetails) GetDeploymentTargetNameOk() (*string, bool)`

GetDeploymentTargetNameOk returns a tuple with the DeploymentTargetName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeploymentTargetName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetails) SetDeploymentTargetName(v string)`

SetDeploymentTargetName sets DeploymentTargetName field to given value.

### HasDeploymentTargetName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetails) HasDeploymentTargetName() bool`

HasDeploymentTargetName returns a boolean if a field has been set.

### GetJobFailureExpirationTime

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetails) GetJobFailureExpirationTime() string`

GetJobFailureExpirationTime returns the JobFailureExpirationTime field if non-nil, zero value otherwise.

### GetJobFailureExpirationTimeOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetails) GetJobFailureExpirationTimeOk() (*string, bool)`

GetJobFailureExpirationTimeOk returns a tuple with the JobFailureExpirationTime field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetJobFailureExpirationTime

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetails) SetJobFailureExpirationTime(v string)`

SetJobFailureExpirationTime sets JobFailureExpirationTime field to given value.

### HasJobFailureExpirationTime

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetails) HasJobFailureExpirationTime() bool`

HasJobFailureExpirationTime returns a boolean if a field has been set.

### GetMaxJobCreationAttempts

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetails) GetMaxJobCreationAttempts() int32`

GetMaxJobCreationAttempts returns the MaxJobCreationAttempts field if non-nil, zero value otherwise.

### GetMaxJobCreationAttemptsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetails) GetMaxJobCreationAttemptsOk() (*int32, bool)`

GetMaxJobCreationAttemptsOk returns a tuple with the MaxJobCreationAttempts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaxJobCreationAttempts

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetails) SetMaxJobCreationAttempts(v int32)`

SetMaxJobCreationAttempts sets MaxJobCreationAttempts field to given value.

### HasMaxJobCreationAttempts

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetails) HasMaxJobCreationAttempts() bool`

HasMaxJobCreationAttempts returns a boolean if a field has been set.

### GetMaxSavepointCreationAttempts

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetails) GetMaxSavepointCreationAttempts() int32`

GetMaxSavepointCreationAttempts returns the MaxSavepointCreationAttempts field if non-nil, zero value otherwise.

### GetMaxSavepointCreationAttemptsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetails) GetMaxSavepointCreationAttemptsOk() (*int32, bool)`

GetMaxSavepointCreationAttemptsOk returns a tuple with the MaxSavepointCreationAttempts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaxSavepointCreationAttempts

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetails) SetMaxSavepointCreationAttempts(v int32)`

SetMaxSavepointCreationAttempts sets MaxSavepointCreationAttempts field to given value.

### HasMaxSavepointCreationAttempts

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetails) HasMaxSavepointCreationAttempts() bool`

HasMaxSavepointCreationAttempts returns a boolean if a field has been set.

### GetRestoreStrategy

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetails) GetRestoreStrategy() ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpRestoreStrategy`

GetRestoreStrategy returns the RestoreStrategy field if non-nil, zero value otherwise.

### GetRestoreStrategyOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetails) GetRestoreStrategyOk() (*ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpRestoreStrategy, bool)`

GetRestoreStrategyOk returns a tuple with the RestoreStrategy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRestoreStrategy

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetails) SetRestoreStrategy(v ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpRestoreStrategy)`

SetRestoreStrategy sets RestoreStrategy field to given value.

### HasRestoreStrategy

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetails) HasRestoreStrategy() bool`

HasRestoreStrategy returns a boolean if a field has been set.

### GetSessionClusterName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetails) GetSessionClusterName() string`

GetSessionClusterName returns the SessionClusterName field if non-nil, zero value otherwise.

### GetSessionClusterNameOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetails) GetSessionClusterNameOk() (*string, bool)`

GetSessionClusterNameOk returns a tuple with the SessionClusterName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSessionClusterName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetails) SetSessionClusterName(v string)`

SetSessionClusterName sets SessionClusterName field to given value.

### HasSessionClusterName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetails) HasSessionClusterName() bool`

HasSessionClusterName returns a boolean if a field has been set.

### GetState

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetails) GetState() string`

GetState returns the State field if non-nil, zero value otherwise.

### GetStateOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetails) GetStateOk() (*string, bool)`

GetStateOk returns a tuple with the State field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetState

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetails) SetState(v string)`

SetState sets State field to given value.

### HasState

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetails) HasState() bool`

HasState returns a boolean if a field has been set.

### GetTemplate

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetails) GetTemplate() ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplate`

GetTemplate returns the Template field if non-nil, zero value otherwise.

### GetTemplateOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetails) GetTemplateOk() (*ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplate, bool)`

GetTemplateOk returns a tuple with the Template field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTemplate

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetails) SetTemplate(v ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentDetailsTemplate)`

SetTemplate sets Template field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


