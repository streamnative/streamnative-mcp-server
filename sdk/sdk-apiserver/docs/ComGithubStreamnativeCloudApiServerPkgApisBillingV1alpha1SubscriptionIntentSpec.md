# ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntentSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**OfferRef** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferReference**](ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferReference.md) |  | [optional] 
**Suger** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SugerSubscriptionIntentSpec**](ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SugerSubscriptionIntentSpec.md) |  | [optional] 
**Type** | Pointer to **string** | The type of the subscription intent. Validate values: stripe, suger Default to stripe. | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntentSpec

`func NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntentSpec() *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntentSpec`

NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntentSpec instantiates a new ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntentSpec object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntentSpecWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntentSpecWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntentSpec`

NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntentSpecWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntentSpec object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetOfferRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntentSpec) GetOfferRef() ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferReference`

GetOfferRef returns the OfferRef field if non-nil, zero value otherwise.

### GetOfferRefOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntentSpec) GetOfferRefOk() (*ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferReference, bool)`

GetOfferRefOk returns a tuple with the OfferRef field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOfferRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntentSpec) SetOfferRef(v ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferReference)`

SetOfferRef sets OfferRef field to given value.

### HasOfferRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntentSpec) HasOfferRef() bool`

HasOfferRef returns a boolean if a field has been set.

### GetSuger

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntentSpec) GetSuger() ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SugerSubscriptionIntentSpec`

GetSuger returns the Suger field if non-nil, zero value otherwise.

### GetSugerOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntentSpec) GetSugerOk() (*ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SugerSubscriptionIntentSpec, bool)`

GetSugerOk returns a tuple with the Suger field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSuger

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntentSpec) SetSuger(v ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SugerSubscriptionIntentSpec)`

SetSuger sets Suger field to given value.

### HasSuger

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntentSpec) HasSuger() bool`

HasSuger returns a boolean if a field has been set.

### GetType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntentSpec) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntentSpec) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntentSpec) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntentSpec) HasType() bool`

HasType returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


