# ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntent

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ApiVersion** | Pointer to **string** | APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources | [optional] 
**Kind** | Pointer to **string** | Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds | [optional] 
**Metadata** | Pointer to [**V1ObjectMeta**](V1ObjectMeta.md) |  | [optional] 
**Spec** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntentSpec**](ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntentSpec.md) |  | [optional] 
**Status** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntentStatus**](ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntentStatus.md) |  | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntent

`func NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntent() *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntent`

NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntent instantiates a new ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntent object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntentWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntentWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntent`

NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntentWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntent object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntent) GetApiVersion() string`

GetApiVersion returns the ApiVersion field if non-nil, zero value otherwise.

### GetApiVersionOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntent) GetApiVersionOk() (*string, bool)`

GetApiVersionOk returns a tuple with the ApiVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntent) SetApiVersion(v string)`

SetApiVersion sets ApiVersion field to given value.

### HasApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntent) HasApiVersion() bool`

HasApiVersion returns a boolean if a field has been set.

### GetKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntent) GetKind() string`

GetKind returns the Kind field if non-nil, zero value otherwise.

### GetKindOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntent) GetKindOk() (*string, bool)`

GetKindOk returns a tuple with the Kind field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntent) SetKind(v string)`

SetKind sets Kind field to given value.

### HasKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntent) HasKind() bool`

HasKind returns a boolean if a field has been set.

### GetMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntent) GetMetadata() V1ObjectMeta`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntent) GetMetadataOk() (*V1ObjectMeta, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntent) SetMetadata(v V1ObjectMeta)`

SetMetadata sets Metadata field to given value.

### HasMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntent) HasMetadata() bool`

HasMetadata returns a boolean if a field has been set.

### GetSpec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntent) GetSpec() ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntentSpec`

GetSpec returns the Spec field if non-nil, zero value otherwise.

### GetSpecOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntent) GetSpecOk() (*ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntentSpec, bool)`

GetSpecOk returns a tuple with the Spec field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSpec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntent) SetSpec(v ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntentSpec)`

SetSpec sets Spec field to given value.

### HasSpec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntent) HasSpec() bool`

HasSpec returns a boolean if a field has been set.

### GetStatus

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntent) GetStatus() ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntentStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntent) GetStatusOk() (*ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntentStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntent) SetStatus(v ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntentStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1SubscriptionIntent) HasStatus() bool`

HasStatus returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


