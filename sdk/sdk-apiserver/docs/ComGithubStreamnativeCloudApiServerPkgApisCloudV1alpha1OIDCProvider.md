# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OIDCProvider

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ApiVersion** | Pointer to **string** | APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources | [optional] 
**Kind** | Pointer to **string** | Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds | [optional] 
**Metadata** | Pointer to [**V1ObjectMeta**](V1ObjectMeta.md) |  | [optional] 
**Spec** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OIDCProviderSpec**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OIDCProviderSpec.md) |  | [optional] 
**Status** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OIDCProviderStatus**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OIDCProviderStatus.md) |  | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OIDCProvider

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OIDCProvider() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OIDCProvider`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OIDCProvider instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OIDCProvider object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OIDCProviderWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OIDCProviderWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OIDCProvider`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OIDCProviderWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OIDCProvider object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OIDCProvider) GetApiVersion() string`

GetApiVersion returns the ApiVersion field if non-nil, zero value otherwise.

### GetApiVersionOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OIDCProvider) GetApiVersionOk() (*string, bool)`

GetApiVersionOk returns a tuple with the ApiVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OIDCProvider) SetApiVersion(v string)`

SetApiVersion sets ApiVersion field to given value.

### HasApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OIDCProvider) HasApiVersion() bool`

HasApiVersion returns a boolean if a field has been set.

### GetKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OIDCProvider) GetKind() string`

GetKind returns the Kind field if non-nil, zero value otherwise.

### GetKindOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OIDCProvider) GetKindOk() (*string, bool)`

GetKindOk returns a tuple with the Kind field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OIDCProvider) SetKind(v string)`

SetKind sets Kind field to given value.

### HasKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OIDCProvider) HasKind() bool`

HasKind returns a boolean if a field has been set.

### GetMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OIDCProvider) GetMetadata() V1ObjectMeta`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OIDCProvider) GetMetadataOk() (*V1ObjectMeta, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OIDCProvider) SetMetadata(v V1ObjectMeta)`

SetMetadata sets Metadata field to given value.

### HasMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OIDCProvider) HasMetadata() bool`

HasMetadata returns a boolean if a field has been set.

### GetSpec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OIDCProvider) GetSpec() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OIDCProviderSpec`

GetSpec returns the Spec field if non-nil, zero value otherwise.

### GetSpecOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OIDCProvider) GetSpecOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OIDCProviderSpec, bool)`

GetSpecOk returns a tuple with the Spec field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSpec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OIDCProvider) SetSpec(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OIDCProviderSpec)`

SetSpec sets Spec field to given value.

### HasSpec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OIDCProvider) HasSpec() bool`

HasSpec returns a boolean if a field has been set.

### GetStatus

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OIDCProvider) GetStatus() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OIDCProviderStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OIDCProvider) GetStatusOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OIDCProviderStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OIDCProvider) SetStatus(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OIDCProviderStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1OIDCProvider) HasStatus() bool`

HasStatus returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


