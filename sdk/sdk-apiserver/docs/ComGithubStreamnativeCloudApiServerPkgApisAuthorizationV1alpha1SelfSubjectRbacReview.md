# ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReview

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ApiVersion** | Pointer to **string** | APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources | [optional] 
**Kind** | Pointer to **string** | Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds | [optional] 
**Metadata** | Pointer to [**V1ObjectMeta**](V1ObjectMeta.md) |  | [optional] 
**Spec** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReviewSpec**](ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReviewSpec.md) |  | [optional] 
**Status** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReviewStatus**](ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReviewStatus.md) |  | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReview

`func NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReview() *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReview`

NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReview instantiates a new ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReview object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReviewWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReviewWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReview`

NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReviewWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReview object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReview) GetApiVersion() string`

GetApiVersion returns the ApiVersion field if non-nil, zero value otherwise.

### GetApiVersionOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReview) GetApiVersionOk() (*string, bool)`

GetApiVersionOk returns a tuple with the ApiVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReview) SetApiVersion(v string)`

SetApiVersion sets ApiVersion field to given value.

### HasApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReview) HasApiVersion() bool`

HasApiVersion returns a boolean if a field has been set.

### GetKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReview) GetKind() string`

GetKind returns the Kind field if non-nil, zero value otherwise.

### GetKindOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReview) GetKindOk() (*string, bool)`

GetKindOk returns a tuple with the Kind field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReview) SetKind(v string)`

SetKind sets Kind field to given value.

### HasKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReview) HasKind() bool`

HasKind returns a boolean if a field has been set.

### GetMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReview) GetMetadata() V1ObjectMeta`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReview) GetMetadataOk() (*V1ObjectMeta, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReview) SetMetadata(v V1ObjectMeta)`

SetMetadata sets Metadata field to given value.

### HasMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReview) HasMetadata() bool`

HasMetadata returns a boolean if a field has been set.

### GetSpec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReview) GetSpec() ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReviewSpec`

GetSpec returns the Spec field if non-nil, zero value otherwise.

### GetSpecOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReview) GetSpecOk() (*ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReviewSpec, bool)`

GetSpecOk returns a tuple with the Spec field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSpec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReview) SetSpec(v ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReviewSpec)`

SetSpec sets Spec field to given value.

### HasSpec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReview) HasSpec() bool`

HasSpec returns a boolean if a field has been set.

### GetStatus

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReview) GetStatus() ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReviewStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReview) GetStatusOk() (*ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReviewStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReview) SetStatus(v ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReviewStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReview) HasStatus() bool`

HasStatus returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


