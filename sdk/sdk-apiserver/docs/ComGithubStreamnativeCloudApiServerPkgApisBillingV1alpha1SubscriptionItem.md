# ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionItem

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Key** | Pointer to **string** | Key is the item key within the subscription. | [optional] 
**Metadata** | Pointer to **map[string]string** | Metadata is an unstructured key value map stored with an item. | [optional] 
**Product** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductReference**](ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductReference.md) |  | [optional] 
**Quantity** | Pointer to **string** | Quantity for this item. | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionItem

`func NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionItem() *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionItem`

NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionItem instantiates a new ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionItem object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionItemWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionItemWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionItem`

NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionItemWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionItem object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetKey

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionItem) GetKey() string`

GetKey returns the Key field if non-nil, zero value otherwise.

### GetKeyOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionItem) GetKeyOk() (*string, bool)`

GetKeyOk returns a tuple with the Key field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKey

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionItem) SetKey(v string)`

SetKey sets Key field to given value.

### HasKey

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionItem) HasKey() bool`

HasKey returns a boolean if a field has been set.

### GetMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionItem) GetMetadata() map[string]string`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionItem) GetMetadataOk() (*map[string]string, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionItem) SetMetadata(v map[string]string)`

SetMetadata sets Metadata field to given value.

### HasMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionItem) HasMetadata() bool`

HasMetadata returns a boolean if a field has been set.

### GetProduct

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionItem) GetProduct() ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductReference`

GetProduct returns the Product field if non-nil, zero value otherwise.

### GetProductOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionItem) GetProductOk() (*ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductReference, bool)`

GetProductOk returns a tuple with the Product field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProduct

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionItem) SetProduct(v ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1ProductReference)`

SetProduct sets Product field to given value.

### HasProduct

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionItem) HasProduct() bool`

HasProduct returns a boolean if a field has been set.

### GetQuantity

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionItem) GetQuantity() string`

GetQuantity returns the Quantity field if non-nil, zero value otherwise.

### GetQuantityOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionItem) GetQuantityOk() (*string, bool)`

GetQuantityOk returns a tuple with the Quantity field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetQuantity

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionItem) SetQuantity(v string)`

SetQuantity sets Quantity field to given value.

### HasQuantity

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionItem) HasQuantity() bool`

HasQuantity returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


