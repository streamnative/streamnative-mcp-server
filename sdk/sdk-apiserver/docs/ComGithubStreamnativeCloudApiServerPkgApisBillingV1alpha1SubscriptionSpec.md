# ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AnchorDate** | Pointer to **time.Time** | AnchorDate is a timestamp representing the first billing cycle end date. It is represented by seconds from the epoch on the Stripe side. | [optional] 
**CloudType** | Pointer to **string** | CloudType will validate resources like the consumption unit product are restricted to the correct cloud provider | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**EndDate** | Pointer to **time.Time** | EndDate is a timestamp representing the date when the subscription will be ended. It is represented in RFC3339 form and is in UTC. | [optional] 
**EndingBalanceCents** | Pointer to **int64** | Ending balance for a subscription, this value is asynchrnously updated by billing-reporter and directly pulled from stripe&#39;s invoice object [1].  Negative at this value means that there are outstanding discount credits left for the customer.  Nil implies that billing reporter hasn&#39;t run since creation and yet to set the value. [1] https://docs.stripe.com/api/invoices/object#invoice_object-ending_balance | [optional] 
**Once** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionItem**](ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionItem.md) | One-time items. | [optional] 
**ParentSubscription** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionReference**](ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionReference.md) |  | [optional] 
**Recurring** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionItem**](ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionItem.md) | Recurring items. | [optional] 
**StartDate** | Pointer to **time.Time** | StartDate is a timestamp representing the start date of the subscription. It is represented in RFC3339 form and is in UTC. | [optional] 
**Stripe** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripeSubscriptionSpec**](ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripeSubscriptionSpec.md) |  | [optional] 
**Suger** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SugerSubscriptionSpec**](ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SugerSubscriptionSpec.md) |  | [optional] 
**Type** | Pointer to **string** | The type of the subscription. Validate values: stripe, suger Default to stripe. | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec

`func NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec() *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec`

NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec instantiates a new ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpecWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpecWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec`

NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpecWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAnchorDate

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) GetAnchorDate() time.Time`

GetAnchorDate returns the AnchorDate field if non-nil, zero value otherwise.

### GetAnchorDateOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) GetAnchorDateOk() (*time.Time, bool)`

GetAnchorDateOk returns a tuple with the AnchorDate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAnchorDate

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) SetAnchorDate(v time.Time)`

SetAnchorDate sets AnchorDate field to given value.

### HasAnchorDate

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) HasAnchorDate() bool`

HasAnchorDate returns a boolean if a field has been set.

### GetCloudType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) GetCloudType() string`

GetCloudType returns the CloudType field if non-nil, zero value otherwise.

### GetCloudTypeOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) GetCloudTypeOk() (*string, bool)`

GetCloudTypeOk returns a tuple with the CloudType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCloudType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) SetCloudType(v string)`

SetCloudType sets CloudType field to given value.

### HasCloudType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) HasCloudType() bool`

HasCloudType returns a boolean if a field has been set.

### GetDescription

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetEndDate

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) GetEndDate() time.Time`

GetEndDate returns the EndDate field if non-nil, zero value otherwise.

### GetEndDateOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) GetEndDateOk() (*time.Time, bool)`

GetEndDateOk returns a tuple with the EndDate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndDate

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) SetEndDate(v time.Time)`

SetEndDate sets EndDate field to given value.

### HasEndDate

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) HasEndDate() bool`

HasEndDate returns a boolean if a field has been set.

### GetEndingBalanceCents

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) GetEndingBalanceCents() int64`

GetEndingBalanceCents returns the EndingBalanceCents field if non-nil, zero value otherwise.

### GetEndingBalanceCentsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) GetEndingBalanceCentsOk() (*int64, bool)`

GetEndingBalanceCentsOk returns a tuple with the EndingBalanceCents field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndingBalanceCents

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) SetEndingBalanceCents(v int64)`

SetEndingBalanceCents sets EndingBalanceCents field to given value.

### HasEndingBalanceCents

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) HasEndingBalanceCents() bool`

HasEndingBalanceCents returns a boolean if a field has been set.

### GetOnce

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) GetOnce() []ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionItem`

GetOnce returns the Once field if non-nil, zero value otherwise.

### GetOnceOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) GetOnceOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionItem, bool)`

GetOnceOk returns a tuple with the Once field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOnce

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) SetOnce(v []ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionItem)`

SetOnce sets Once field to given value.

### HasOnce

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) HasOnce() bool`

HasOnce returns a boolean if a field has been set.

### GetParentSubscription

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) GetParentSubscription() ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionReference`

GetParentSubscription returns the ParentSubscription field if non-nil, zero value otherwise.

### GetParentSubscriptionOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) GetParentSubscriptionOk() (*ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionReference, bool)`

GetParentSubscriptionOk returns a tuple with the ParentSubscription field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetParentSubscription

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) SetParentSubscription(v ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionReference)`

SetParentSubscription sets ParentSubscription field to given value.

### HasParentSubscription

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) HasParentSubscription() bool`

HasParentSubscription returns a boolean if a field has been set.

### GetRecurring

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) GetRecurring() []ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionItem`

GetRecurring returns the Recurring field if non-nil, zero value otherwise.

### GetRecurringOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) GetRecurringOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionItem, bool)`

GetRecurringOk returns a tuple with the Recurring field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRecurring

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) SetRecurring(v []ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionItem)`

SetRecurring sets Recurring field to given value.

### HasRecurring

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) HasRecurring() bool`

HasRecurring returns a boolean if a field has been set.

### GetStartDate

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) GetStartDate() time.Time`

GetStartDate returns the StartDate field if non-nil, zero value otherwise.

### GetStartDateOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) GetStartDateOk() (*time.Time, bool)`

GetStartDateOk returns a tuple with the StartDate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStartDate

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) SetStartDate(v time.Time)`

SetStartDate sets StartDate field to given value.

### HasStartDate

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) HasStartDate() bool`

HasStartDate returns a boolean if a field has been set.

### GetStripe

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) GetStripe() ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripeSubscriptionSpec`

GetStripe returns the Stripe field if non-nil, zero value otherwise.

### GetStripeOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) GetStripeOk() (*ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripeSubscriptionSpec, bool)`

GetStripeOk returns a tuple with the Stripe field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStripe

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) SetStripe(v ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripeSubscriptionSpec)`

SetStripe sets Stripe field to given value.

### HasStripe

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) HasStripe() bool`

HasStripe returns a boolean if a field has been set.

### GetSuger

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) GetSuger() ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SugerSubscriptionSpec`

GetSuger returns the Suger field if non-nil, zero value otherwise.

### GetSugerOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) GetSugerOk() (*ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SugerSubscriptionSpec, bool)`

GetSugerOk returns a tuple with the Suger field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSuger

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) SetSuger(v ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SugerSubscriptionSpec)`

SetSuger sets Suger field to given value.

### HasSuger

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) HasSuger() bool`

HasSuger returns a boolean if a field has been set.

### GetType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionSpec) HasType() bool`

HasType returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


