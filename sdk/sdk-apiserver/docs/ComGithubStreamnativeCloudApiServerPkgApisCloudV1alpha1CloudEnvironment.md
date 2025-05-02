# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironment

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ApiVersion** | Pointer to **string** | APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources | [optional] 
**Kind** | Pointer to **string** | Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds | [optional] 
**Metadata** | Pointer to [**V1ObjectMeta**](V1ObjectMeta.md) |  | [optional] 
**Spec** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentSpec**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentSpec.md) |  | [optional] 
**Status** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentStatus**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentStatus.md) |  | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironment

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironment() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironment`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironment instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironment object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironment`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironment object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironment) GetApiVersion() string`

GetApiVersion returns the ApiVersion field if non-nil, zero value otherwise.

### GetApiVersionOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironment) GetApiVersionOk() (*string, bool)`

GetApiVersionOk returns a tuple with the ApiVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironment) SetApiVersion(v string)`

SetApiVersion sets ApiVersion field to given value.

### HasApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironment) HasApiVersion() bool`

HasApiVersion returns a boolean if a field has been set.

### GetKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironment) GetKind() string`

GetKind returns the Kind field if non-nil, zero value otherwise.

### GetKindOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironment) GetKindOk() (*string, bool)`

GetKindOk returns a tuple with the Kind field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironment) SetKind(v string)`

SetKind sets Kind field to given value.

### HasKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironment) HasKind() bool`

HasKind returns a boolean if a field has been set.

### GetMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironment) GetMetadata() V1ObjectMeta`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironment) GetMetadataOk() (*V1ObjectMeta, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironment) SetMetadata(v V1ObjectMeta)`

SetMetadata sets Metadata field to given value.

### HasMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironment) HasMetadata() bool`

HasMetadata returns a boolean if a field has been set.

### GetSpec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironment) GetSpec() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentSpec`

GetSpec returns the Spec field if non-nil, zero value otherwise.

### GetSpecOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironment) GetSpecOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentSpec, bool)`

GetSpecOk returns a tuple with the Spec field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSpec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironment) SetSpec(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentSpec)`

SetSpec sets Spec field to given value.

### HasSpec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironment) HasSpec() bool`

HasSpec returns a boolean if a field has been set.

### GetStatus

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironment) GetStatus() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironment) GetStatusOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironment) SetStatus(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironment) HasStatus() bool`

HasStatus returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


