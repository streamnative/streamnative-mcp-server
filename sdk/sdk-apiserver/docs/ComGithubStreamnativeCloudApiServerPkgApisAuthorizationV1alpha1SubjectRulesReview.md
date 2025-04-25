# ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReview

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ApiVersion** | Pointer to **string** | APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources | [optional] 
**Kind** | Pointer to **string** | Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds | [optional] 
**Metadata** | Pointer to [**V1ObjectMeta**](V1ObjectMeta.md) |  | [optional] 
**Spec** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReviewSpec**](ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReviewSpec.md) |  | [optional] 
**Status** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReviewStatus**](ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReviewStatus.md) |  | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReview

`func NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReview() *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReview`

NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReview instantiates a new ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReview object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReviewWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReviewWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReview`

NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReviewWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReview object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReview) GetApiVersion() string`

GetApiVersion returns the ApiVersion field if non-nil, zero value otherwise.

### GetApiVersionOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReview) GetApiVersionOk() (*string, bool)`

GetApiVersionOk returns a tuple with the ApiVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReview) SetApiVersion(v string)`

SetApiVersion sets ApiVersion field to given value.

### HasApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReview) HasApiVersion() bool`

HasApiVersion returns a boolean if a field has been set.

### GetKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReview) GetKind() string`

GetKind returns the Kind field if non-nil, zero value otherwise.

### GetKindOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReview) GetKindOk() (*string, bool)`

GetKindOk returns a tuple with the Kind field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReview) SetKind(v string)`

SetKind sets Kind field to given value.

### HasKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReview) HasKind() bool`

HasKind returns a boolean if a field has been set.

### GetMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReview) GetMetadata() V1ObjectMeta`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReview) GetMetadataOk() (*V1ObjectMeta, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReview) SetMetadata(v V1ObjectMeta)`

SetMetadata sets Metadata field to given value.

### HasMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReview) HasMetadata() bool`

HasMetadata returns a boolean if a field has been set.

### GetSpec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReview) GetSpec() ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReviewSpec`

GetSpec returns the Spec field if non-nil, zero value otherwise.

### GetSpecOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReview) GetSpecOk() (*ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReviewSpec, bool)`

GetSpecOk returns a tuple with the Spec field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSpec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReview) SetSpec(v ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReviewSpec)`

SetSpec sets Spec field to given value.

### HasSpec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReview) HasSpec() bool`

HasSpec returns a boolean if a field has been set.

### GetStatus

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReview) GetStatus() ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReviewStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReview) GetStatusOk() (*ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReviewStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReview) SetStatus(v ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReviewStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReview) HasStatus() bool`

HasStatus returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


