# ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Description** | Pointer to **string** | Description is a product description | [optional] 
**Name** | Pointer to **string** | Name is the product name. | [optional] 
**Prices** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductPrice**](ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductPrice.md) | Prices associated with the product. | [optional] 
**Stripe** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripeProductSpec**](ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripeProductSpec.md) |  | [optional] 
**Suger** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SugerProductSpec**](ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SugerProductSpec.md) |  | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductSpec

`func NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductSpec() *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductSpec`

NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductSpec instantiates a new ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductSpec object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductSpecWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductSpecWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductSpec`

NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductSpecWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductSpec object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDescription

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductSpec) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductSpec) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductSpec) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductSpec) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductSpec) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductSpec) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductSpec) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductSpec) HasName() bool`

HasName returns a boolean if a field has been set.

### GetPrices

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductSpec) GetPrices() []ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductPrice`

GetPrices returns the Prices field if non-nil, zero value otherwise.

### GetPricesOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductSpec) GetPricesOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductPrice, bool)`

GetPricesOk returns a tuple with the Prices field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrices

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductSpec) SetPrices(v []ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductPrice)`

SetPrices sets Prices field to given value.

### HasPrices

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductSpec) HasPrices() bool`

HasPrices returns a boolean if a field has been set.

### GetStripe

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductSpec) GetStripe() ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripeProductSpec`

GetStripe returns the Stripe field if non-nil, zero value otherwise.

### GetStripeOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductSpec) GetStripeOk() (*ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripeProductSpec, bool)`

GetStripeOk returns a tuple with the Stripe field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStripe

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductSpec) SetStripe(v ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripeProductSpec)`

SetStripe sets Stripe field to given value.

### HasStripe

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductSpec) HasStripe() bool`

HasStripe returns a boolean if a field has been set.

### GetSuger

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductSpec) GetSuger() ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SugerProductSpec`

GetSuger returns the Suger field if non-nil, zero value otherwise.

### GetSugerOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductSpec) GetSugerOk() (*ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SugerProductSpec, bool)`

GetSugerOk returns a tuple with the Suger field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSuger

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductSpec) SetSuger(v ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SugerProductSpec)`

SetSuger sets Suger field to given value.

### HasSuger

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductSpec) HasSuger() bool`

HasSuger returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


