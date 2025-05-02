# ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1Tier

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**FlatAmount** | Pointer to **string** | FlatAmount is the flat billing amount for an entire tier, regardless of the number of units in the tier, in dollars. | [optional] 
**UnitAmount** | Pointer to **string** | UnitAmount is the per-unit billing amount for each individual unit for which this tier applies, in dollars. | [optional] 
**UpTo** | Pointer to **string** | UpTo specifies the upper bound of this tier. The lower bound of a tier is the upper bound of the previous tier adding one. Use &#x60;inf&#x60; to define a fallback tier. | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1Tier

`func NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1Tier() *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1Tier`

NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1Tier instantiates a new ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1Tier object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1TierWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1TierWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1Tier`

NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1TierWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1Tier object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetFlatAmount

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1Tier) GetFlatAmount() string`

GetFlatAmount returns the FlatAmount field if non-nil, zero value otherwise.

### GetFlatAmountOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1Tier) GetFlatAmountOk() (*string, bool)`

GetFlatAmountOk returns a tuple with the FlatAmount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFlatAmount

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1Tier) SetFlatAmount(v string)`

SetFlatAmount sets FlatAmount field to given value.

### HasFlatAmount

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1Tier) HasFlatAmount() bool`

HasFlatAmount returns a boolean if a field has been set.

### GetUnitAmount

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1Tier) GetUnitAmount() string`

GetUnitAmount returns the UnitAmount field if non-nil, zero value otherwise.

### GetUnitAmountOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1Tier) GetUnitAmountOk() (*string, bool)`

GetUnitAmountOk returns a tuple with the UnitAmount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUnitAmount

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1Tier) SetUnitAmount(v string)`

SetUnitAmount sets UnitAmount field to given value.

### HasUnitAmount

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1Tier) HasUnitAmount() bool`

HasUnitAmount returns a boolean if a field has been set.

### GetUpTo

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1Tier) GetUpTo() string`

GetUpTo returns the UpTo field if non-nil, zero value otherwise.

### GetUpToOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1Tier) GetUpToOk() (*string, bool)`

GetUpToOk returns a tuple with the UpTo field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpTo

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1Tier) SetUpTo(v string)`

SetUpTo sets UpTo field to given value.

### HasUpTo

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1Tier) HasUpTo() bool`

HasUpTo returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


