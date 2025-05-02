# ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRulesReviewStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**EvaluationError** | Pointer to **string** | EvaluationError can appear in combination with Rules. It indicates an error occurred during rule evaluation, such as an authorizer that doesn&#39;t support rule evaluation, and that ResourceRules and/or NonResourceRules may be incomplete. | [optional] 
**Incomplete** | **bool** | Incomplete is true when the rules returned by this call are incomplete. This is most commonly encountered when an authorizer, such as an external authorizer, doesn&#39;t support rules evaluation. | 
**ResourceRules** | [**[]ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1ResourceRule**](ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1ResourceRule.md) | ResourceRules is the list of actions the subject is allowed to perform on resources. The list ordering isn&#39;t significant, may contain duplicates, and possibly be incomplete. | 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRulesReviewStatus

`func NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRulesReviewStatus(incomplete bool, resourceRules []ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1ResourceRule, ) *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRulesReviewStatus`

NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRulesReviewStatus instantiates a new ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRulesReviewStatus object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRulesReviewStatusWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRulesReviewStatusWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRulesReviewStatus`

NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRulesReviewStatusWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRulesReviewStatus object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEvaluationError

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRulesReviewStatus) GetEvaluationError() string`

GetEvaluationError returns the EvaluationError field if non-nil, zero value otherwise.

### GetEvaluationErrorOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRulesReviewStatus) GetEvaluationErrorOk() (*string, bool)`

GetEvaluationErrorOk returns a tuple with the EvaluationError field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEvaluationError

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRulesReviewStatus) SetEvaluationError(v string)`

SetEvaluationError sets EvaluationError field to given value.

### HasEvaluationError

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRulesReviewStatus) HasEvaluationError() bool`

HasEvaluationError returns a boolean if a field has been set.

### GetIncomplete

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRulesReviewStatus) GetIncomplete() bool`

GetIncomplete returns the Incomplete field if non-nil, zero value otherwise.

### GetIncompleteOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRulesReviewStatus) GetIncompleteOk() (*bool, bool)`

GetIncompleteOk returns a tuple with the Incomplete field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIncomplete

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRulesReviewStatus) SetIncomplete(v bool)`

SetIncomplete sets Incomplete field to given value.


### GetResourceRules

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRulesReviewStatus) GetResourceRules() []ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1ResourceRule`

GetResourceRules returns the ResourceRules field if non-nil, zero value otherwise.

### GetResourceRulesOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRulesReviewStatus) GetResourceRulesOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1ResourceRule, bool)`

GetResourceRulesOk returns a tuple with the ResourceRules field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResourceRules

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRulesReviewStatus) SetResourceRules(v []ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1ResourceRule)`

SetResourceRules sets ResourceRules field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


