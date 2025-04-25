# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**BillingParent** | Pointer to **string** | reconciled parent of this organization. if spec.BillingParent is set but status.BillingParent is not set, then reconciler will create a parent child relationship if spec.BillingParent is not set but status.BillingParent is set, then reconciler will delete parent child relationship if spec.BillingParent is set but status.BillingParent is set and same, then reconciler will do nothing if spec.BillingParent is set but status.Billingparent is set and different, then reconciler will delete status.Billingparent relationship and create new spec.BillingParent relationship | [optional] 
**Conditions** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition.md) | Conditions is an array of current observed conditions. | [optional] 
**SubscriptionName** | Pointer to **string** | Indicates the active subscription for this organization. This information is available when the Subscribed condition is true. | [optional] 
**SupportPlan** | Pointer to **string** | returns support plan of current subscription.  blank, implies either no support plan or legacy support | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStatus

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStatus() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStatus`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStatus instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStatus object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStatusWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStatusWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStatus`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStatusWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStatus object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetBillingParent

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStatus) GetBillingParent() string`

GetBillingParent returns the BillingParent field if non-nil, zero value otherwise.

### GetBillingParentOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStatus) GetBillingParentOk() (*string, bool)`

GetBillingParentOk returns a tuple with the BillingParent field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBillingParent

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStatus) SetBillingParent(v string)`

SetBillingParent sets BillingParent field to given value.

### HasBillingParent

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStatus) HasBillingParent() bool`

HasBillingParent returns a boolean if a field has been set.

### GetConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStatus) GetConditions() []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition`

GetConditions returns the Conditions field if non-nil, zero value otherwise.

### GetConditionsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStatus) GetConditionsOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition, bool)`

GetConditionsOk returns a tuple with the Conditions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStatus) SetConditions(v []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition)`

SetConditions sets Conditions field to given value.

### HasConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStatus) HasConditions() bool`

HasConditions returns a boolean if a field has been set.

### GetSubscriptionName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStatus) GetSubscriptionName() string`

GetSubscriptionName returns the SubscriptionName field if non-nil, zero value otherwise.

### GetSubscriptionNameOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStatus) GetSubscriptionNameOk() (*string, bool)`

GetSubscriptionNameOk returns a tuple with the SubscriptionName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSubscriptionName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStatus) SetSubscriptionName(v string)`

SetSubscriptionName sets SubscriptionName field to given value.

### HasSubscriptionName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStatus) HasSubscriptionName() bool`

HasSubscriptionName returns a boolean if a field has been set.

### GetSupportPlan

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStatus) GetSupportPlan() string`

GetSupportPlan returns the SupportPlan field if non-nil, zero value otherwise.

### GetSupportPlanOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStatus) GetSupportPlanOk() (*string, bool)`

GetSupportPlanOk returns a tuple with the SupportPlan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSupportPlan

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStatus) SetSupportPlan(v string)`

SetSupportPlan sets SupportPlan field to given value.

### HasSupportPlan

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStatus) HasSupportPlan() bool`

HasSupportPlan returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


