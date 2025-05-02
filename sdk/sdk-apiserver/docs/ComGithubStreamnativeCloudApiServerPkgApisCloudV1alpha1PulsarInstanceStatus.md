# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Auth** | [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceStatusAuth**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceStatusAuth.md) |  | 
**Conditions** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition.md) | Conditions is an array of current observed conditions. | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceStatus

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceStatus(auth ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceStatusAuth, ) *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceStatus`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceStatus instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceStatus object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceStatusWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceStatusWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceStatus`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceStatusWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceStatus object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAuth

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceStatus) GetAuth() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceStatusAuth`

GetAuth returns the Auth field if non-nil, zero value otherwise.

### GetAuthOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceStatus) GetAuthOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceStatusAuth, bool)`

GetAuthOk returns a tuple with the Auth field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuth

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceStatus) SetAuth(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceStatusAuth)`

SetAuth sets Auth field to given value.


### GetConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceStatus) GetConditions() []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition`

GetConditions returns the Conditions field if non-nil, zero value otherwise.

### GetConditionsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceStatus) GetConditionsOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition, bool)`

GetConditionsOk returns a tuple with the Conditions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceStatus) SetConditions(v []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition)`

SetConditions sets Conditions field to given value.

### HasConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceStatus) HasConditions() bool`

HasConditions returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


