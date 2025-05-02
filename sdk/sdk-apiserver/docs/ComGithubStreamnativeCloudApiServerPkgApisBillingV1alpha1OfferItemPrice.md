# ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPrice

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Currency** | Pointer to **string** | Three-letter [ISO currency code](https://www.iso.org/iso-4217-currency-codes.html), in lowercase. | [optional] 
**Product** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductReference**](ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductReference.md) |  | [optional] 
**Recurring** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPriceRecurring**](ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPriceRecurring.md) |  | [optional] 
**Tiers** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1Tier**](ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1Tier.md) | Tiers are Stripes billing tiers like | [optional] 
**TiersMode** | Pointer to **string** | TiersMode is the stripe tier mode | [optional] 
**UnitAmount** | Pointer to **string** | A quantity (or 0 for a free price) representing how much to charge. | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPrice

`func NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPrice() *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPrice`

NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPrice instantiates a new ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPrice object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPriceWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPriceWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPrice`

NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPriceWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPrice object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCurrency

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPrice) GetCurrency() string`

GetCurrency returns the Currency field if non-nil, zero value otherwise.

### GetCurrencyOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPrice) GetCurrencyOk() (*string, bool)`

GetCurrencyOk returns a tuple with the Currency field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCurrency

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPrice) SetCurrency(v string)`

SetCurrency sets Currency field to given value.

### HasCurrency

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPrice) HasCurrency() bool`

HasCurrency returns a boolean if a field has been set.

### GetProduct

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPrice) GetProduct() ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductReference`

GetProduct returns the Product field if non-nil, zero value otherwise.

### GetProductOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPrice) GetProductOk() (*ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductReference, bool)`

GetProductOk returns a tuple with the Product field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProduct

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPrice) SetProduct(v ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductReference)`

SetProduct sets Product field to given value.

### HasProduct

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPrice) HasProduct() bool`

HasProduct returns a boolean if a field has been set.

### GetRecurring

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPrice) GetRecurring() ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPriceRecurring`

GetRecurring returns the Recurring field if non-nil, zero value otherwise.

### GetRecurringOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPrice) GetRecurringOk() (*ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPriceRecurring, bool)`

GetRecurringOk returns a tuple with the Recurring field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRecurring

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPrice) SetRecurring(v ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPriceRecurring)`

SetRecurring sets Recurring field to given value.

### HasRecurring

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPrice) HasRecurring() bool`

HasRecurring returns a boolean if a field has been set.

### GetTiers

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPrice) GetTiers() []ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1Tier`

GetTiers returns the Tiers field if non-nil, zero value otherwise.

### GetTiersOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPrice) GetTiersOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1Tier, bool)`

GetTiersOk returns a tuple with the Tiers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTiers

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPrice) SetTiers(v []ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1Tier)`

SetTiers sets Tiers field to given value.

### HasTiers

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPrice) HasTiers() bool`

HasTiers returns a boolean if a field has been set.

### GetTiersMode

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPrice) GetTiersMode() string`

GetTiersMode returns the TiersMode field if non-nil, zero value otherwise.

### GetTiersModeOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPrice) GetTiersModeOk() (*string, bool)`

GetTiersModeOk returns a tuple with the TiersMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTiersMode

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPrice) SetTiersMode(v string)`

SetTiersMode sets TiersMode field to given value.

### HasTiersMode

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPrice) HasTiersMode() bool`

HasTiersMode returns a boolean if a field has been set.

### GetUnitAmount

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPrice) GetUnitAmount() string`

GetUnitAmount returns the UnitAmount field if non-nil, zero value otherwise.

### GetUnitAmountOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPrice) GetUnitAmountOk() (*string, bool)`

GetUnitAmountOk returns a tuple with the UnitAmount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUnitAmount

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPrice) SetUnitAmount(v string)`

SetUnitAmount sets UnitAmount field to given value.

### HasUnitAmount

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPrice) HasUnitAmount() bool`

HasUnitAmount returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


