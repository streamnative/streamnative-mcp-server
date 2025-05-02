# ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Key** | Pointer to **string** | Key is the item key within the offer and subscription. | [optional] 
**Metadata** | Pointer to **map[string]string** | Metadata is an unstructured key value map stored with an item. | [optional] 
**Price** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPrice**](ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPrice.md) |  | [optional] 
**PriceRef** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PriceReference**](ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PriceReference.md) |  | [optional] 
**Quantity** | Pointer to **string** | Quantity for this item. | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem

`func NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem() *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem`

NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem instantiates a new ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem`

NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetKey

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem) GetKey() string`

GetKey returns the Key field if non-nil, zero value otherwise.

### GetKeyOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem) GetKeyOk() (*string, bool)`

GetKeyOk returns a tuple with the Key field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKey

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem) SetKey(v string)`

SetKey sets Key field to given value.

### HasKey

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem) HasKey() bool`

HasKey returns a boolean if a field has been set.

### GetMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem) GetMetadata() map[string]string`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem) GetMetadataOk() (*map[string]string, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem) SetMetadata(v map[string]string)`

SetMetadata sets Metadata field to given value.

### HasMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem) HasMetadata() bool`

HasMetadata returns a boolean if a field has been set.

### GetPrice

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem) GetPrice() ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPrice`

GetPrice returns the Price field if non-nil, zero value otherwise.

### GetPriceOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem) GetPriceOk() (*ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPrice, bool)`

GetPriceOk returns a tuple with the Price field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrice

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem) SetPrice(v ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPrice)`

SetPrice sets Price field to given value.

### HasPrice

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem) HasPrice() bool`

HasPrice returns a boolean if a field has been set.

### GetPriceRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem) GetPriceRef() ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PriceReference`

GetPriceRef returns the PriceRef field if non-nil, zero value otherwise.

### GetPriceRefOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem) GetPriceRefOk() (*ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PriceReference, bool)`

GetPriceRefOk returns a tuple with the PriceRef field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPriceRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem) SetPriceRef(v ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PriceReference)`

SetPriceRef sets PriceRef field to given value.

### HasPriceRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem) HasPriceRef() bool`

HasPriceRef returns a boolean if a field has been set.

### GetQuantity

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem) GetQuantity() string`

GetQuantity returns the Quantity field if non-nil, zero value otherwise.

### GetQuantityOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem) GetQuantityOk() (*string, bool)`

GetQuantityOk returns a tuple with the Quantity field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetQuantity

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem) SetQuantity(v string)`

SetQuantity sets Quantity field to given value.

### HasQuantity

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItem) HasQuantity() bool`

HasQuantity returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


