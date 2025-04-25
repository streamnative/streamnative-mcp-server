# ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Conditions** | Pointer to [**[]V1Condition**](V1Condition.md) | Conditions is an array of current observed conditions. | [optional] 
**Items** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionStatusSubscriptionItem**](ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionStatusSubscriptionItem.md) | the status of the subscription is designed to support billing agents, so it provides product and subscription items. | [optional] 
**PendingSetupIntent** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SetupIntentReference**](ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SetupIntentReference.md) |  | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionStatus

`func NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionStatus() *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionStatus`

NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionStatus instantiates a new ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionStatus object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionStatusWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionStatusWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionStatus`

NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionStatusWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionStatus object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionStatus) GetConditions() []V1Condition`

GetConditions returns the Conditions field if non-nil, zero value otherwise.

### GetConditionsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionStatus) GetConditionsOk() (*[]V1Condition, bool)`

GetConditionsOk returns a tuple with the Conditions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionStatus) SetConditions(v []V1Condition)`

SetConditions sets Conditions field to given value.

### HasConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionStatus) HasConditions() bool`

HasConditions returns a boolean if a field has been set.

### GetItems

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionStatus) GetItems() []ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionStatusSubscriptionItem`

GetItems returns the Items field if non-nil, zero value otherwise.

### GetItemsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionStatus) GetItemsOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionStatusSubscriptionItem, bool)`

GetItemsOk returns a tuple with the Items field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetItems

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionStatus) SetItems(v []ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionStatusSubscriptionItem)`

SetItems sets Items field to given value.

### HasItems

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionStatus) HasItems() bool`

HasItems returns a boolean if a field has been set.

### GetPendingSetupIntent

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionStatus) GetPendingSetupIntent() ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SetupIntentReference`

GetPendingSetupIntent returns the PendingSetupIntent field if non-nil, zero value otherwise.

### GetPendingSetupIntentOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionStatus) GetPendingSetupIntentOk() (*ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SetupIntentReference, bool)`

GetPendingSetupIntentOk returns a tuple with the PendingSetupIntent field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPendingSetupIntent

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionStatus) SetPendingSetupIntent(v ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SetupIntentReference)`

SetPendingSetupIntent sets PendingSetupIntent field to given value.

### HasPendingSetupIntent

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionStatus) HasPendingSetupIntent() bool`

HasPendingSetupIntent returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


