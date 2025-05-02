# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStripe

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CollectionMethod** | Pointer to **string** | CollectionMethod is how payment on a subscription is to be collected, either charge_automatically or send_invoice | [optional] 
**DaysUntilDue** | Pointer to **int64** | DaysUntilDue sets the due date for the invoice. Applicable when collection method is send_invoice. | [optional] 
**KeepInDraft** | Pointer to **bool** | KeepInDraft disables auto-advance of the invoice state. Applicable when collection method is send_invoice. | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStripe

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStripe() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStripe`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStripe instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStripe object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStripeWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStripeWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStripe`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStripeWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStripe object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCollectionMethod

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStripe) GetCollectionMethod() string`

GetCollectionMethod returns the CollectionMethod field if non-nil, zero value otherwise.

### GetCollectionMethodOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStripe) GetCollectionMethodOk() (*string, bool)`

GetCollectionMethodOk returns a tuple with the CollectionMethod field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCollectionMethod

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStripe) SetCollectionMethod(v string)`

SetCollectionMethod sets CollectionMethod field to given value.

### HasCollectionMethod

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStripe) HasCollectionMethod() bool`

HasCollectionMethod returns a boolean if a field has been set.

### GetDaysUntilDue

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStripe) GetDaysUntilDue() int64`

GetDaysUntilDue returns the DaysUntilDue field if non-nil, zero value otherwise.

### GetDaysUntilDueOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStripe) GetDaysUntilDueOk() (*int64, bool)`

GetDaysUntilDueOk returns a tuple with the DaysUntilDue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDaysUntilDue

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStripe) SetDaysUntilDue(v int64)`

SetDaysUntilDue sets DaysUntilDue field to given value.

### HasDaysUntilDue

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStripe) HasDaysUntilDue() bool`

HasDaysUntilDue returns a boolean if a field has been set.

### GetKeepInDraft

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStripe) GetKeepInDraft() bool`

GetKeepInDraft returns the KeepInDraft field if non-nil, zero value otherwise.

### GetKeepInDraftOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStripe) GetKeepInDraftOk() (*bool, bool)`

GetKeepInDraftOk returns a tuple with the KeepInDraft field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKeepInDraft

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStripe) SetKeepInDraft(v bool)`

SetKeepInDraft sets KeepInDraft field to given value.

### HasKeepInDraft

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStripe) HasKeepInDraft() bool`

HasKeepInDraft returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


