# ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PublicOfferSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Description** | Pointer to **string** |  | [optional] 
**Once** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem**](ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem.md) | One-time items, each with an attached price. | [optional] 
**Recurring** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem**](ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem.md) | Recurring items, each with an attached price. | [optional] 
**Stripe** | Pointer to **map[string]interface{}** |  | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PublicOfferSpec

`func NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PublicOfferSpec() *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PublicOfferSpec`

NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PublicOfferSpec instantiates a new ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PublicOfferSpec object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PublicOfferSpecWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PublicOfferSpecWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PublicOfferSpec`

NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PublicOfferSpecWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PublicOfferSpec object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDescription

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PublicOfferSpec) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PublicOfferSpec) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PublicOfferSpec) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PublicOfferSpec) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetOnce

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PublicOfferSpec) GetOnce() []ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem`

GetOnce returns the Once field if non-nil, zero value otherwise.

### GetOnceOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PublicOfferSpec) GetOnceOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem, bool)`

GetOnceOk returns a tuple with the Once field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOnce

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PublicOfferSpec) SetOnce(v []ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem)`

SetOnce sets Once field to given value.

### HasOnce

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PublicOfferSpec) HasOnce() bool`

HasOnce returns a boolean if a field has been set.

### GetRecurring

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PublicOfferSpec) GetRecurring() []ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem`

GetRecurring returns the Recurring field if non-nil, zero value otherwise.

### GetRecurringOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PublicOfferSpec) GetRecurringOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem, bool)`

GetRecurringOk returns a tuple with the Recurring field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRecurring

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PublicOfferSpec) SetRecurring(v []ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem)`

SetRecurring sets Recurring field to given value.

### HasRecurring

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PublicOfferSpec) HasRecurring() bool`

HasRecurring returns a boolean if a field has been set.

### GetStripe

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PublicOfferSpec) GetStripe() map[string]interface{}`

GetStripe returns the Stripe field if non-nil, zero value otherwise.

### GetStripeOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PublicOfferSpec) GetStripeOk() (*map[string]interface{}, bool)`

GetStripeOk returns a tuple with the Stripe field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStripe

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PublicOfferSpec) SetStripe(v map[string]interface{})`

SetStripe sets Stripe field to given value.

### HasStripe

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PublicOfferSpec) HasStripe() bool`

HasStripe returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


