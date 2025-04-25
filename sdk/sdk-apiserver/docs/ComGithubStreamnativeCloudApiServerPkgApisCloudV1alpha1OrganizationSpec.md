# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AwsMarketplaceToken** | Pointer to **string** |  | [optional] 
**BillingAccount** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BillingAccountSpec**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BillingAccountSpec.md) |  | [optional] 
**BillingParent** | Pointer to **string** | Organiztion ID of this org&#39;s billing parent.  Once set, this org will inherit parent org&#39;s subscription, paymentMetthod and have \&quot;inherited\&quot; as billing type | [optional] 
**BillingType** | Pointer to **string** | BillingType indicates the method of subscription that the organization uses. It is primarily consumed by the cloud-manager to be able to distinguish between invoiced subscriptions. | [optional] 
**CloudFeature** | Pointer to **map[string]bool** | CloudFeature indicates features this org wants to enable/disable | [optional] 
**DefaultPoolRef** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolRef**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolRef.md) |  | [optional] 
**DisplayName** | Pointer to **string** | Name to display to our users | [optional] 
**Domains** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Domain**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Domain.md) |  | [optional] 
**Gcloud** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudOrganizationSpec**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudOrganizationSpec.md) |  | [optional] 
**Metadata** | Pointer to **map[string]string** | Metadata is user-visible (and possibly editable) metadata. | [optional] 
**Stripe** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStripe**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStripe.md) |  | [optional] 
**Suger** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSuger**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSuger.md) |  | [optional] 
**Type** | Pointer to **string** |  | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpecWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpecWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpecWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAwsMarketplaceToken

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) GetAwsMarketplaceToken() string`

GetAwsMarketplaceToken returns the AwsMarketplaceToken field if non-nil, zero value otherwise.

### GetAwsMarketplaceTokenOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) GetAwsMarketplaceTokenOk() (*string, bool)`

GetAwsMarketplaceTokenOk returns a tuple with the AwsMarketplaceToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAwsMarketplaceToken

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) SetAwsMarketplaceToken(v string)`

SetAwsMarketplaceToken sets AwsMarketplaceToken field to given value.

### HasAwsMarketplaceToken

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) HasAwsMarketplaceToken() bool`

HasAwsMarketplaceToken returns a boolean if a field has been set.

### GetBillingAccount

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) GetBillingAccount() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BillingAccountSpec`

GetBillingAccount returns the BillingAccount field if non-nil, zero value otherwise.

### GetBillingAccountOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) GetBillingAccountOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BillingAccountSpec, bool)`

GetBillingAccountOk returns a tuple with the BillingAccount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBillingAccount

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) SetBillingAccount(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BillingAccountSpec)`

SetBillingAccount sets BillingAccount field to given value.

### HasBillingAccount

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) HasBillingAccount() bool`

HasBillingAccount returns a boolean if a field has been set.

### GetBillingParent

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) GetBillingParent() string`

GetBillingParent returns the BillingParent field if non-nil, zero value otherwise.

### GetBillingParentOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) GetBillingParentOk() (*string, bool)`

GetBillingParentOk returns a tuple with the BillingParent field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBillingParent

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) SetBillingParent(v string)`

SetBillingParent sets BillingParent field to given value.

### HasBillingParent

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) HasBillingParent() bool`

HasBillingParent returns a boolean if a field has been set.

### GetBillingType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) GetBillingType() string`

GetBillingType returns the BillingType field if non-nil, zero value otherwise.

### GetBillingTypeOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) GetBillingTypeOk() (*string, bool)`

GetBillingTypeOk returns a tuple with the BillingType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBillingType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) SetBillingType(v string)`

SetBillingType sets BillingType field to given value.

### HasBillingType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) HasBillingType() bool`

HasBillingType returns a boolean if a field has been set.

### GetCloudFeature

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) GetCloudFeature() map[string]bool`

GetCloudFeature returns the CloudFeature field if non-nil, zero value otherwise.

### GetCloudFeatureOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) GetCloudFeatureOk() (*map[string]bool, bool)`

GetCloudFeatureOk returns a tuple with the CloudFeature field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCloudFeature

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) SetCloudFeature(v map[string]bool)`

SetCloudFeature sets CloudFeature field to given value.

### HasCloudFeature

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) HasCloudFeature() bool`

HasCloudFeature returns a boolean if a field has been set.

### GetDefaultPoolRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) GetDefaultPoolRef() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolRef`

GetDefaultPoolRef returns the DefaultPoolRef field if non-nil, zero value otherwise.

### GetDefaultPoolRefOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) GetDefaultPoolRefOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolRef, bool)`

GetDefaultPoolRefOk returns a tuple with the DefaultPoolRef field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDefaultPoolRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) SetDefaultPoolRef(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolRef)`

SetDefaultPoolRef sets DefaultPoolRef field to given value.

### HasDefaultPoolRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) HasDefaultPoolRef() bool`

HasDefaultPoolRef returns a boolean if a field has been set.

### GetDisplayName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) GetDisplayName() string`

GetDisplayName returns the DisplayName field if non-nil, zero value otherwise.

### GetDisplayNameOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) GetDisplayNameOk() (*string, bool)`

GetDisplayNameOk returns a tuple with the DisplayName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDisplayName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) SetDisplayName(v string)`

SetDisplayName sets DisplayName field to given value.

### HasDisplayName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) HasDisplayName() bool`

HasDisplayName returns a boolean if a field has been set.

### GetDomains

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) GetDomains() []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Domain`

GetDomains returns the Domains field if non-nil, zero value otherwise.

### GetDomainsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) GetDomainsOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Domain, bool)`

GetDomainsOk returns a tuple with the Domains field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDomains

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) SetDomains(v []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Domain)`

SetDomains sets Domains field to given value.

### HasDomains

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) HasDomains() bool`

HasDomains returns a boolean if a field has been set.

### GetGcloud

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) GetGcloud() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudOrganizationSpec`

GetGcloud returns the Gcloud field if non-nil, zero value otherwise.

### GetGcloudOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) GetGcloudOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudOrganizationSpec, bool)`

GetGcloudOk returns a tuple with the Gcloud field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGcloud

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) SetGcloud(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudOrganizationSpec)`

SetGcloud sets Gcloud field to given value.

### HasGcloud

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) HasGcloud() bool`

HasGcloud returns a boolean if a field has been set.

### GetMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) GetMetadata() map[string]string`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) GetMetadataOk() (*map[string]string, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) SetMetadata(v map[string]string)`

SetMetadata sets Metadata field to given value.

### HasMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) HasMetadata() bool`

HasMetadata returns a boolean if a field has been set.

### GetStripe

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) GetStripe() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStripe`

GetStripe returns the Stripe field if non-nil, zero value otherwise.

### GetStripeOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) GetStripeOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStripe, bool)`

GetStripeOk returns a tuple with the Stripe field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStripe

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) SetStripe(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationStripe)`

SetStripe sets Stripe field to given value.

### HasStripe

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) HasStripe() bool`

HasStripe returns a boolean if a field has been set.

### GetSuger

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) GetSuger() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSuger`

GetSuger returns the Suger field if non-nil, zero value otherwise.

### GetSugerOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) GetSugerOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSuger, bool)`

GetSugerOk returns a tuple with the Suger field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSuger

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) SetSuger(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSuger)`

SetSuger sets Suger field to given value.

### HasSuger

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) HasSuger() bool`

HasSuger returns a boolean if a field has been set.

### GetType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OrganizationSpec) HasType() bool`

HasType returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


