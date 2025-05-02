# ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AnchorDate** | Pointer to **time.Time** | AnchorDate is a timestamp representing the first billing cycle end date. This will be used to anchor future billing periods to that date. For example, setting the anchor date for a subscription starting on Apr 1 to be Apr 12 will send the invoice for the subscription out on Apr 12 and the 12th of every following month for a monthly subscription. It is represented in RFC3339 form and is in UTC. | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**Duration** | Pointer to **string** | Duration indicates how long the subscription for this offer should last. The value must greater than 0 | [optional] 
**EndDate** | Pointer to **time.Time** | EndDate is a timestamp representing the planned end date of the subscription. It is represented in RFC3339 form and is in UTC. | [optional] 
**Once** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem**](ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem.md) | One-time items, each with an attached price. | [optional] 
**Recurring** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem**](ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem.md) | Recurring items, each with an attached price. | [optional] 
**StartDate** | Pointer to **time.Time** | StartDate is a timestamp representing the planned start date of the subscription. It is represented in RFC3339 form and is in UTC. | [optional] 
**Stripe** | Pointer to **map[string]interface{}** |  | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec

`func NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec() *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec`

NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec instantiates a new ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpecWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpecWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec`

NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpecWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAnchorDate

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec) GetAnchorDate() time.Time`

GetAnchorDate returns the AnchorDate field if non-nil, zero value otherwise.

### GetAnchorDateOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec) GetAnchorDateOk() (*time.Time, bool)`

GetAnchorDateOk returns a tuple with the AnchorDate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAnchorDate

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec) SetAnchorDate(v time.Time)`

SetAnchorDate sets AnchorDate field to given value.

### HasAnchorDate

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec) HasAnchorDate() bool`

HasAnchorDate returns a boolean if a field has been set.

### GetDescription

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetDuration

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec) GetDuration() string`

GetDuration returns the Duration field if non-nil, zero value otherwise.

### GetDurationOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec) GetDurationOk() (*string, bool)`

GetDurationOk returns a tuple with the Duration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDuration

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec) SetDuration(v string)`

SetDuration sets Duration field to given value.

### HasDuration

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec) HasDuration() bool`

HasDuration returns a boolean if a field has been set.

### GetEndDate

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec) GetEndDate() time.Time`

GetEndDate returns the EndDate field if non-nil, zero value otherwise.

### GetEndDateOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec) GetEndDateOk() (*time.Time, bool)`

GetEndDateOk returns a tuple with the EndDate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndDate

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec) SetEndDate(v time.Time)`

SetEndDate sets EndDate field to given value.

### HasEndDate

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec) HasEndDate() bool`

HasEndDate returns a boolean if a field has been set.

### GetOnce

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec) GetOnce() []ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem`

GetOnce returns the Once field if non-nil, zero value otherwise.

### GetOnceOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec) GetOnceOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem, bool)`

GetOnceOk returns a tuple with the Once field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOnce

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec) SetOnce(v []ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem)`

SetOnce sets Once field to given value.

### HasOnce

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec) HasOnce() bool`

HasOnce returns a boolean if a field has been set.

### GetRecurring

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec) GetRecurring() []ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem`

GetRecurring returns the Recurring field if non-nil, zero value otherwise.

### GetRecurringOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec) GetRecurringOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem, bool)`

GetRecurringOk returns a tuple with the Recurring field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRecurring

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec) SetRecurring(v []ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem)`

SetRecurring sets Recurring field to given value.

### HasRecurring

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec) HasRecurring() bool`

HasRecurring returns a boolean if a field has been set.

### GetStartDate

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec) GetStartDate() time.Time`

GetStartDate returns the StartDate field if non-nil, zero value otherwise.

### GetStartDateOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec) GetStartDateOk() (*time.Time, bool)`

GetStartDateOk returns a tuple with the StartDate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStartDate

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec) SetStartDate(v time.Time)`

SetStartDate sets StartDate field to given value.

### HasStartDate

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec) HasStartDate() bool`

HasStartDate returns a boolean if a field has been set.

### GetStripe

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec) GetStripe() map[string]interface{}`

GetStripe returns the Stripe field if non-nil, zero value otherwise.

### GetStripeOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec) GetStripeOk() (*map[string]interface{}, bool)`

GetStripeOk returns a tuple with the Stripe field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStripe

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec) SetStripe(v map[string]interface{})`

SetStripe sets Stripe field to given value.

### HasStripe

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec) HasStripe() bool`

HasStripe returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


