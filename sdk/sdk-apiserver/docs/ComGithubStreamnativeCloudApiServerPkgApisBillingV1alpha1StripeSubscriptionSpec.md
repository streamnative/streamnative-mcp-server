# ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripeSubscriptionSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CollectionMethod** | Pointer to **string** | CollectionMethod is how payment on a subscription is to be collected, either charge_automatically or send_invoice | [optional] 
**DaysUntilDue** | Pointer to **int64** | DaysUntilDue is applicable when collection method is send_invoice | [optional] 
**Id** | Pointer to **string** | ID is the stripe subscription ID. | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripeSubscriptionSpec

`func NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripeSubscriptionSpec() *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripeSubscriptionSpec`

NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripeSubscriptionSpec instantiates a new ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripeSubscriptionSpec object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripeSubscriptionSpecWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripeSubscriptionSpecWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripeSubscriptionSpec`

NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripeSubscriptionSpecWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripeSubscriptionSpec object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCollectionMethod

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripeSubscriptionSpec) GetCollectionMethod() string`

GetCollectionMethod returns the CollectionMethod field if non-nil, zero value otherwise.

### GetCollectionMethodOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripeSubscriptionSpec) GetCollectionMethodOk() (*string, bool)`

GetCollectionMethodOk returns a tuple with the CollectionMethod field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCollectionMethod

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripeSubscriptionSpec) SetCollectionMethod(v string)`

SetCollectionMethod sets CollectionMethod field to given value.

### HasCollectionMethod

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripeSubscriptionSpec) HasCollectionMethod() bool`

HasCollectionMethod returns a boolean if a field has been set.

### GetDaysUntilDue

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripeSubscriptionSpec) GetDaysUntilDue() int64`

GetDaysUntilDue returns the DaysUntilDue field if non-nil, zero value otherwise.

### GetDaysUntilDueOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripeSubscriptionSpec) GetDaysUntilDueOk() (*int64, bool)`

GetDaysUntilDueOk returns a tuple with the DaysUntilDue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDaysUntilDue

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripeSubscriptionSpec) SetDaysUntilDue(v int64)`

SetDaysUntilDue sets DaysUntilDue field to given value.

### HasDaysUntilDue

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripeSubscriptionSpec) HasDaysUntilDue() bool`

HasDaysUntilDue returns a boolean if a field has been set.

### GetId

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripeSubscriptionSpec) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripeSubscriptionSpec) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripeSubscriptionSpec) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1StripeSubscriptionSpec) HasId() bool`

HasId returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


