# ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOffer

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ApiVersion** | Pointer to **string** | APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources | [optional] 
**Kind** | Pointer to **string** | Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds | [optional] 
**Metadata** | Pointer to [**V1ObjectMeta**](V1ObjectMeta.md) |  | [optional] 
**Spec** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec**](ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec.md) |  | [optional] 
**Status** | Pointer to **map[string]interface{}** | PrivateOfferStatus defines the observed state of PrivateOffer | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOffer

`func NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOffer() *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOffer`

NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOffer instantiates a new ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOffer object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOffer`

NewComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOffer object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOffer) GetApiVersion() string`

GetApiVersion returns the ApiVersion field if non-nil, zero value otherwise.

### GetApiVersionOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOffer) GetApiVersionOk() (*string, bool)`

GetApiVersionOk returns a tuple with the ApiVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOffer) SetApiVersion(v string)`

SetApiVersion sets ApiVersion field to given value.

### HasApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOffer) HasApiVersion() bool`

HasApiVersion returns a boolean if a field has been set.

### GetKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOffer) GetKind() string`

GetKind returns the Kind field if non-nil, zero value otherwise.

### GetKindOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOffer) GetKindOk() (*string, bool)`

GetKindOk returns a tuple with the Kind field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOffer) SetKind(v string)`

SetKind sets Kind field to given value.

### HasKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOffer) HasKind() bool`

HasKind returns a boolean if a field has been set.

### GetMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOffer) GetMetadata() V1ObjectMeta`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOffer) GetMetadataOk() (*V1ObjectMeta, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOffer) SetMetadata(v V1ObjectMeta)`

SetMetadata sets Metadata field to given value.

### HasMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOffer) HasMetadata() bool`

HasMetadata returns a boolean if a field has been set.

### GetSpec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOffer) GetSpec() ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec`

GetSpec returns the Spec field if non-nil, zero value otherwise.

### GetSpecOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOffer) GetSpecOk() (*ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec, bool)`

GetSpecOk returns a tuple with the Spec field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSpec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOffer) SetSpec(v ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOfferSpec)`

SetSpec sets Spec field to given value.

### HasSpec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOffer) HasSpec() bool`

HasSpec returns a boolean if a field has been set.

### GetStatus

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOffer) GetStatus() map[string]interface{}`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOffer) GetStatusOk() (*map[string]interface{}, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOffer) SetStatus(v map[string]interface{})`

SetStatus sets Status field to given value.

### HasStatus

`func (o *ComGithubStreamnativeCloudApiServerPkgApisBillingV1alpha1PrivateOffer) HasStatus() bool`

HasStatus returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


