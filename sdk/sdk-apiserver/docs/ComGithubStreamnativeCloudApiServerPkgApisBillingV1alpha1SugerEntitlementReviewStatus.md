# ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SugerEntitlementReviewStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**BuyerID** | Pointer to **string** | BuyerID is the ID of buyer associated with organization The first one will be returned when there are more than one are found. | [optional] 
**Conditions** | Pointer to [**[]V1Condition**](V1Condition.md) | Conditions is an array of current observed conditions. | [optional] 
**Organization** | Pointer to **string** | Organization is the name of the organization matching the entitlement ID The first one will be returned when there are more than one are found. | [optional] 
**Partner** | Pointer to **string** | Partner is the partner associated with organization | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SugerEntitlementReviewStatus

`func NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SugerEntitlementReviewStatus() *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SugerEntitlementReviewStatus`

NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SugerEntitlementReviewStatus instantiates a new ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SugerEntitlementReviewStatus object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SugerEntitlementReviewStatusWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SugerEntitlementReviewStatusWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SugerEntitlementReviewStatus`

NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SugerEntitlementReviewStatusWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SugerEntitlementReviewStatus object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetBuyerID

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SugerEntitlementReviewStatus) GetBuyerID() string`

GetBuyerID returns the BuyerID field if non-nil, zero value otherwise.

### GetBuyerIDOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SugerEntitlementReviewStatus) GetBuyerIDOk() (*string, bool)`

GetBuyerIDOk returns a tuple with the BuyerID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBuyerID

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SugerEntitlementReviewStatus) SetBuyerID(v string)`

SetBuyerID sets BuyerID field to given value.

### HasBuyerID

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SugerEntitlementReviewStatus) HasBuyerID() bool`

HasBuyerID returns a boolean if a field has been set.

### GetConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SugerEntitlementReviewStatus) GetConditions() []V1Condition`

GetConditions returns the Conditions field if non-nil, zero value otherwise.

### GetConditionsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SugerEntitlementReviewStatus) GetConditionsOk() (*[]V1Condition, bool)`

GetConditionsOk returns a tuple with the Conditions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SugerEntitlementReviewStatus) SetConditions(v []V1Condition)`

SetConditions sets Conditions field to given value.

### HasConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SugerEntitlementReviewStatus) HasConditions() bool`

HasConditions returns a boolean if a field has been set.

### GetOrganization

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SugerEntitlementReviewStatus) GetOrganization() string`

GetOrganization returns the Organization field if non-nil, zero value otherwise.

### GetOrganizationOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SugerEntitlementReviewStatus) GetOrganizationOk() (*string, bool)`

GetOrganizationOk returns a tuple with the Organization field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrganization

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SugerEntitlementReviewStatus) SetOrganization(v string)`

SetOrganization sets Organization field to given value.

### HasOrganization

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SugerEntitlementReviewStatus) HasOrganization() bool`

HasOrganization returns a boolean if a field has been set.

### GetPartner

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SugerEntitlementReviewStatus) GetPartner() string`

GetPartner returns the Partner field if non-nil, zero value otherwise.

### GetPartnerOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SugerEntitlementReviewStatus) GetPartnerOk() (*string, bool)`

GetPartnerOk returns a tuple with the Partner field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPartner

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SugerEntitlementReviewStatus) SetPartner(v string)`

SetPartner sets Partner field to given value.

### HasPartner

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SugerEntitlementReviewStatus) HasPartner() bool`

HasPartner returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


