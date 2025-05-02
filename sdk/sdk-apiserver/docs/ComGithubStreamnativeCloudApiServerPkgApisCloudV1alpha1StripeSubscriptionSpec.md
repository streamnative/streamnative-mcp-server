# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CustomerId** | **string** | CustomerId is the id of the customer | 
**Product** | Pointer to **string** | Product is the name or id of the product | [optional] 
**Quantities** | Pointer to **map[string]int64** | Quantities defines the quantity of certain prices in the subscription | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionSpec

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionSpec(customerId string, ) *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionSpec`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionSpec instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionSpec object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionSpecWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionSpecWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionSpec`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionSpecWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionSpec object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCustomerId

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionSpec) GetCustomerId() string`

GetCustomerId returns the CustomerId field if non-nil, zero value otherwise.

### GetCustomerIdOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionSpec) GetCustomerIdOk() (*string, bool)`

GetCustomerIdOk returns a tuple with the CustomerId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCustomerId

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionSpec) SetCustomerId(v string)`

SetCustomerId sets CustomerId field to given value.


### GetProduct

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionSpec) GetProduct() string`

GetProduct returns the Product field if non-nil, zero value otherwise.

### GetProductOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionSpec) GetProductOk() (*string, bool)`

GetProductOk returns a tuple with the Product field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProduct

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionSpec) SetProduct(v string)`

SetProduct sets Product field to given value.

### HasProduct

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionSpec) HasProduct() bool`

HasProduct returns a boolean if a field has been set.

### GetQuantities

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionSpec) GetQuantities() map[string]int64`

GetQuantities returns the Quantities field if non-nil, zero value otherwise.

### GetQuantitiesOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionSpec) GetQuantitiesOk() (*map[string]int64, bool)`

GetQuantitiesOk returns a tuple with the Quantities field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetQuantities

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionSpec) SetQuantities(v map[string]int64)`

SetQuantities sets Quantities field to given value.

### HasQuantities

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1StripeSubscriptionSpec) HasQuantities() bool`

HasQuantities returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


