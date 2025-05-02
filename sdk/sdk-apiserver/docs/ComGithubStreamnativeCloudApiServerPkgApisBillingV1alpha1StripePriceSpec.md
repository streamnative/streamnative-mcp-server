# ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripePriceSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Active** | Pointer to **bool** | Active indicates the price is active on the product | [optional] 
**Currency** | Pointer to **string** | Currency is the required three-letter ISO currency code The codes below were generated from https://stripe.com/docs/currencies | [optional] 
**Name** | Pointer to **string** | Name to be displayed in the Stripe dashboard, hidden from customers | [optional] 
**Recurring** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripePriceRecurrence**](ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripePriceRecurrence.md) |  | [optional] 
**Tiers** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1Tier**](ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1Tier.md) | Tiers are Stripes billing tiers like | [optional] 
**TiersMode** | Pointer to **string** | TiersMode is the stripe tier mode | [optional] 
**UnitAmount** | Pointer to **string** | UnitAmount in dollars. If present, billing_scheme is assumed to be per_unit | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripePriceSpec

`func NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripePriceSpec() *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripePriceSpec`

NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripePriceSpec instantiates a new ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripePriceSpec object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripePriceSpecWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripePriceSpecWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripePriceSpec`

NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripePriceSpecWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripePriceSpec object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetActive

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripePriceSpec) GetActive() bool`

GetActive returns the Active field if non-nil, zero value otherwise.

### GetActiveOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripePriceSpec) GetActiveOk() (*bool, bool)`

GetActiveOk returns a tuple with the Active field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetActive

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripePriceSpec) SetActive(v bool)`

SetActive sets Active field to given value.

### HasActive

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripePriceSpec) HasActive() bool`

HasActive returns a boolean if a field has been set.

### GetCurrency

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripePriceSpec) GetCurrency() string`

GetCurrency returns the Currency field if non-nil, zero value otherwise.

### GetCurrencyOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripePriceSpec) GetCurrencyOk() (*string, bool)`

GetCurrencyOk returns a tuple with the Currency field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCurrency

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripePriceSpec) SetCurrency(v string)`

SetCurrency sets Currency field to given value.

### HasCurrency

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripePriceSpec) HasCurrency() bool`

HasCurrency returns a boolean if a field has been set.

### GetName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripePriceSpec) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripePriceSpec) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripePriceSpec) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripePriceSpec) HasName() bool`

HasName returns a boolean if a field has been set.

### GetRecurring

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripePriceSpec) GetRecurring() ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripePriceRecurrence`

GetRecurring returns the Recurring field if non-nil, zero value otherwise.

### GetRecurringOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripePriceSpec) GetRecurringOk() (*ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripePriceRecurrence, bool)`

GetRecurringOk returns a tuple with the Recurring field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRecurring

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripePriceSpec) SetRecurring(v ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripePriceRecurrence)`

SetRecurring sets Recurring field to given value.

### HasRecurring

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripePriceSpec) HasRecurring() bool`

HasRecurring returns a boolean if a field has been set.

### GetTiers

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripePriceSpec) GetTiers() []ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1Tier`

GetTiers returns the Tiers field if non-nil, zero value otherwise.

### GetTiersOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripePriceSpec) GetTiersOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1Tier, bool)`

GetTiersOk returns a tuple with the Tiers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTiers

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripePriceSpec) SetTiers(v []ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1Tier)`

SetTiers sets Tiers field to given value.

### HasTiers

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripePriceSpec) HasTiers() bool`

HasTiers returns a boolean if a field has been set.

### GetTiersMode

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripePriceSpec) GetTiersMode() string`

GetTiersMode returns the TiersMode field if non-nil, zero value otherwise.

### GetTiersModeOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripePriceSpec) GetTiersModeOk() (*string, bool)`

GetTiersModeOk returns a tuple with the TiersMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTiersMode

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripePriceSpec) SetTiersMode(v string)`

SetTiersMode sets TiersMode field to given value.

### HasTiersMode

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripePriceSpec) HasTiersMode() bool`

HasTiersMode returns a boolean if a field has been set.

### GetUnitAmount

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripePriceSpec) GetUnitAmount() string`

GetUnitAmount returns the UnitAmount field if non-nil, zero value otherwise.

### GetUnitAmountOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripePriceSpec) GetUnitAmountOk() (*string, bool)`

GetUnitAmountOk returns a tuple with the UnitAmount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUnitAmount

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripePriceSpec) SetUnitAmount(v string)`

SetUnitAmount sets UnitAmount field to given value.

### HasUnitAmount

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripePriceSpec) HasUnitAmount() bool`

HasUnitAmount returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


