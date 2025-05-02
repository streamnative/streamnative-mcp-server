# ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Prices** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductStatusPrice**](ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductStatusPrice.md) | The prices as stored in Stripe. | [optional] 
**StripeID** | Pointer to **string** | The unique identifier for the Stripe product. | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductStatus

`func NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductStatus() *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductStatus`

NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductStatus instantiates a new ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductStatus object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductStatusWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductStatusWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductStatus`

NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductStatusWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductStatus object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPrices

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductStatus) GetPrices() []ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductStatusPrice`

GetPrices returns the Prices field if non-nil, zero value otherwise.

### GetPricesOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductStatus) GetPricesOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductStatusPrice, bool)`

GetPricesOk returns a tuple with the Prices field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrices

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductStatus) SetPrices(v []ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductStatusPrice)`

SetPrices sets Prices field to given value.

### HasPrices

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductStatus) HasPrices() bool`

HasPrices returns a boolean if a field has been set.

### GetStripeID

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductStatus) GetStripeID() string`

GetStripeID returns the StripeID field if non-nil, zero value otherwise.

### GetStripeIDOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductStatus) GetStripeIDOk() (*string, bool)`

GetStripeIDOk returns a tuple with the StripeID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStripeID

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductStatus) SetStripeID(v string)`

SetStripeID sets StripeID field to given value.

### HasStripeID

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductStatus) HasStripeID() bool`

HasStripeID returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


