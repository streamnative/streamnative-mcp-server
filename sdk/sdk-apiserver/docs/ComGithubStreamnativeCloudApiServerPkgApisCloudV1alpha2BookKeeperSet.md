# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSet

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ApiVersion** | Pointer to **string** | APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources | [optional] 
**Kind** | Pointer to **string** | Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds | [optional] 
**Metadata** | Pointer to [**V1ObjectMeta**](V1ObjectMeta.md) |  | [optional] 
**Spec** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpec**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpec.md) |  | [optional] 
**Status** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetStatus**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetStatus.md) |  | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSet

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSet() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSet`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSet instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSet object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSet`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSet object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSet) GetApiVersion() string`

GetApiVersion returns the ApiVersion field if non-nil, zero value otherwise.

### GetApiVersionOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSet) GetApiVersionOk() (*string, bool)`

GetApiVersionOk returns a tuple with the ApiVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSet) SetApiVersion(v string)`

SetApiVersion sets ApiVersion field to given value.

### HasApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSet) HasApiVersion() bool`

HasApiVersion returns a boolean if a field has been set.

### GetKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSet) GetKind() string`

GetKind returns the Kind field if non-nil, zero value otherwise.

### GetKindOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSet) GetKindOk() (*string, bool)`

GetKindOk returns a tuple with the Kind field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSet) SetKind(v string)`

SetKind sets Kind field to given value.

### HasKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSet) HasKind() bool`

HasKind returns a boolean if a field has been set.

### GetMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSet) GetMetadata() V1ObjectMeta`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSet) GetMetadataOk() (*V1ObjectMeta, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSet) SetMetadata(v V1ObjectMeta)`

SetMetadata sets Metadata field to given value.

### HasMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSet) HasMetadata() bool`

HasMetadata returns a boolean if a field has been set.

### GetSpec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSet) GetSpec() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpec`

GetSpec returns the Spec field if non-nil, zero value otherwise.

### GetSpecOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSet) GetSpecOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpec, bool)`

GetSpecOk returns a tuple with the Spec field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSpec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSet) SetSpec(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpec)`

SetSpec sets Spec field to given value.

### HasSpec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSet) HasSpec() bool`

HasSpec returns a boolean if a field has been set.

### GetStatus

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSet) GetStatus() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSet) GetStatusOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSet) SetStatus(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSet) HasStatus() bool`

HasStatus returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


