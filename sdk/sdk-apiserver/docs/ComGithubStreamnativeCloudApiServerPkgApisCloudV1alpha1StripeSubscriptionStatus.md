# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Conditions** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition.md) | Conditions is an array of current observed conditions. | [optional] 
**ObservedGeneration** | Pointer to **int64** | The generation observed by the controller. | [optional] 
**SubscriptionItems** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SubscriptionItem**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SubscriptionItem.md) | SubscriptionItems is a list of subscription items (used to report metrics) | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionStatus

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionStatus() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionStatus`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionStatus instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionStatus object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionStatusWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionStatusWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionStatus`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionStatusWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionStatus object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionStatus) GetConditions() []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition`

GetConditions returns the Conditions field if non-nil, zero value otherwise.

### GetConditionsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionStatus) GetConditionsOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition, bool)`

GetConditionsOk returns a tuple with the Conditions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionStatus) SetConditions(v []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition)`

SetConditions sets Conditions field to given value.

### HasConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionStatus) HasConditions() bool`

HasConditions returns a boolean if a field has been set.

### GetObservedGeneration

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionStatus) GetObservedGeneration() int64`

GetObservedGeneration returns the ObservedGeneration field if non-nil, zero value otherwise.

### GetObservedGenerationOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionStatus) GetObservedGenerationOk() (*int64, bool)`

GetObservedGenerationOk returns a tuple with the ObservedGeneration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObservedGeneration

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionStatus) SetObservedGeneration(v int64)`

SetObservedGeneration sets ObservedGeneration field to given value.

### HasObservedGeneration

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionStatus) HasObservedGeneration() bool`

HasObservedGeneration returns a boolean if a field has been set.

### GetSubscriptionItems

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionStatus) GetSubscriptionItems() []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SubscriptionItem`

GetSubscriptionItems returns the SubscriptionItems field if non-nil, zero value otherwise.

### GetSubscriptionItemsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionStatus) GetSubscriptionItemsOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SubscriptionItem, bool)`

GetSubscriptionItemsOk returns a tuple with the SubscriptionItems field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSubscriptionItems

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionStatus) SetSubscriptionItems(v []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SubscriptionItem)`

SetSubscriptionItems sets SubscriptionItems field to given value.

### HasSubscriptionItems

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionStatus) HasSubscriptionItems() bool`

HasSubscriptionItems returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


