# ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPriceRecurring

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AggregateUsage** | Pointer to **string** |  | [optional] 
**Interval** | Pointer to **string** | Specifies billing frequency. Either &#x60;day&#x60;, &#x60;week&#x60;, &#x60;month&#x60; or &#x60;year&#x60;. | [optional] 
**IntervalCount** | Pointer to **int64** | The number of intervals between subscription billings. For example, &#x60;interval&#x3D;month&#x60; and &#x60;interval_count&#x3D;3&#x60; bills every 3 months. Maximum of one-year interval is allowed (1 year, 12 months, or 52 weeks). | [optional] 
**UsageType** | Pointer to **string** |  | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPriceRecurring

`func NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPriceRecurring() *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPriceRecurring`

NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPriceRecurring instantiates a new ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPriceRecurring object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPriceRecurringWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPriceRecurringWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPriceRecurring`

NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPriceRecurringWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPriceRecurring object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAggregateUsage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPriceRecurring) GetAggregateUsage() string`

GetAggregateUsage returns the AggregateUsage field if non-nil, zero value otherwise.

### GetAggregateUsageOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPriceRecurring) GetAggregateUsageOk() (*string, bool)`

GetAggregateUsageOk returns a tuple with the AggregateUsage field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAggregateUsage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPriceRecurring) SetAggregateUsage(v string)`

SetAggregateUsage sets AggregateUsage field to given value.

### HasAggregateUsage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPriceRecurring) HasAggregateUsage() bool`

HasAggregateUsage returns a boolean if a field has been set.

### GetInterval

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPriceRecurring) GetInterval() string`

GetInterval returns the Interval field if non-nil, zero value otherwise.

### GetIntervalOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPriceRecurring) GetIntervalOk() (*string, bool)`

GetIntervalOk returns a tuple with the Interval field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInterval

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPriceRecurring) SetInterval(v string)`

SetInterval sets Interval field to given value.

### HasInterval

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPriceRecurring) HasInterval() bool`

HasInterval returns a boolean if a field has been set.

### GetIntervalCount

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPriceRecurring) GetIntervalCount() int64`

GetIntervalCount returns the IntervalCount field if non-nil, zero value otherwise.

### GetIntervalCountOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPriceRecurring) GetIntervalCountOk() (*int64, bool)`

GetIntervalCountOk returns a tuple with the IntervalCount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIntervalCount

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPriceRecurring) SetIntervalCount(v int64)`

SetIntervalCount sets IntervalCount field to given value.

### HasIntervalCount

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPriceRecurring) HasIntervalCount() bool`

HasIntervalCount returns a boolean if a field has been set.

### GetUsageType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPriceRecurring) GetUsageType() string`

GetUsageType returns the UsageType field if non-nil, zero value otherwise.

### GetUsageTypeOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPriceRecurring) GetUsageTypeOk() (*string, bool)`

GetUsageTypeOk returns a tuple with the UsageType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsageType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPriceRecurring) SetUsageType(v string)`

SetUsageType sets UsageType field to given value.

### HasUsageType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1OfferItemPriceRecurring) HasUsageType() bool`

HasUsageType returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


