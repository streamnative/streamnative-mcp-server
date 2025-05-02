# ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1CustomerPortalRequestStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Conditions** | Pointer to [**[]V1Condition**](V1Condition.md) | Conditions is an array of current observed conditions. | [optional] 
**Stripe** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripeCustomerPortalRequestStatus**](ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripeCustomerPortalRequestStatus.md) |  | [optional] 
**Url** | Pointer to **string** |  | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1CustomerPortalRequestStatus

`func NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1CustomerPortalRequestStatus() *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1CustomerPortalRequestStatus`

NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1CustomerPortalRequestStatus instantiates a new ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1CustomerPortalRequestStatus object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1CustomerPortalRequestStatusWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1CustomerPortalRequestStatusWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1CustomerPortalRequestStatus`

NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1CustomerPortalRequestStatusWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1CustomerPortalRequestStatus object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1CustomerPortalRequestStatus) GetConditions() []V1Condition`

GetConditions returns the Conditions field if non-nil, zero value otherwise.

### GetConditionsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1CustomerPortalRequestStatus) GetConditionsOk() (*[]V1Condition, bool)`

GetConditionsOk returns a tuple with the Conditions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1CustomerPortalRequestStatus) SetConditions(v []V1Condition)`

SetConditions sets Conditions field to given value.

### HasConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1CustomerPortalRequestStatus) HasConditions() bool`

HasConditions returns a boolean if a field has been set.

### GetStripe

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1CustomerPortalRequestStatus) GetStripe() ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripeCustomerPortalRequestStatus`

GetStripe returns the Stripe field if non-nil, zero value otherwise.

### GetStripeOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1CustomerPortalRequestStatus) GetStripeOk() (*ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripeCustomerPortalRequestStatus, bool)`

GetStripeOk returns a tuple with the Stripe field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStripe

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1CustomerPortalRequestStatus) SetStripe(v ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripeCustomerPortalRequestStatus)`

SetStripe sets Stripe field to given value.

### HasStripe

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1CustomerPortalRequestStatus) HasStripe() bool`

HasStripe returns a boolean if a field has been set.

### GetUrl

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1CustomerPortalRequestStatus) GetUrl() string`

GetUrl returns the Url field if non-nil, zero value otherwise.

### GetUrlOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1CustomerPortalRequestStatus) GetUrlOk() (*string, bool)`

GetUrlOk returns a tuple with the Url field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUrl

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1CustomerPortalRequestStatus) SetUrl(v string)`

SetUrl sets Url field to given value.

### HasUrl

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1CustomerPortalRequestStatus) HasUrl() bool`

HasUrl returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


