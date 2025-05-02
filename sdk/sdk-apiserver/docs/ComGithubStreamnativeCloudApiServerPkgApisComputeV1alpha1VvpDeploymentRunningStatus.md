# ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentRunningStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Conditions** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentStatusCondition**](ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentStatusCondition.md) | Conditions is an array of current observed conditions. | [optional] 
**JobId** | Pointer to **string** |  | [optional] 
**TransitionTime** | Pointer to **time.Time** | Time is a wrapper around time.Time which supports correct marshaling to YAML and JSON.  Wrappers are provided for many of the factory methods that the time package offers. | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentRunningStatus

`func NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentRunningStatus() *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentRunningStatus`

NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentRunningStatus instantiates a new ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentRunningStatus object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentRunningStatusWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentRunningStatusWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentRunningStatus`

NewComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentRunningStatusWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentRunningStatus object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentRunningStatus) GetConditions() []ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentStatusCondition`

GetConditions returns the Conditions field if non-nil, zero value otherwise.

### GetConditionsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentRunningStatus) GetConditionsOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentStatusCondition, bool)`

GetConditionsOk returns a tuple with the Conditions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentRunningStatus) SetConditions(v []ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentStatusCondition)`

SetConditions sets Conditions field to given value.

### HasConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentRunningStatus) HasConditions() bool`

HasConditions returns a boolean if a field has been set.

### GetJobId

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentRunningStatus) GetJobId() string`

GetJobId returns the JobId field if non-nil, zero value otherwise.

### GetJobIdOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentRunningStatus) GetJobIdOk() (*string, bool)`

GetJobIdOk returns a tuple with the JobId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetJobId

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentRunningStatus) SetJobId(v string)`

SetJobId sets JobId field to given value.

### HasJobId

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentRunningStatus) HasJobId() bool`

HasJobId returns a boolean if a field has been set.

### GetTransitionTime

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentRunningStatus) GetTransitionTime() time.Time`

GetTransitionTime returns the TransitionTime field if non-nil, zero value otherwise.

### GetTransitionTimeOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentRunningStatus) GetTransitionTimeOk() (*time.Time, bool)`

GetTransitionTimeOk returns a tuple with the TransitionTime field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTransitionTime

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentRunningStatus) SetTransitionTime(v time.Time)`

SetTransitionTime sets TransitionTime field to given value.

### HasTransitionTime

`func (o *ComGithubStreamnativeCloudApiServerPkgApisComputeV1alpha1VvpDeploymentRunningStatus) HasTransitionTime() bool`

HasTransitionTime returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


