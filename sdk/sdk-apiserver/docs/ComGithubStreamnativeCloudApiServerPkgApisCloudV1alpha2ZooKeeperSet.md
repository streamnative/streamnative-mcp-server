# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSet

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ApiVersion** | Pointer to **string** | APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources | [optional] 
**Kind** | Pointer to **string** | Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds | [optional] 
**Metadata** | Pointer to [**V1ObjectMeta**](V1ObjectMeta.md) |  | [optional] 
**Spec** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec.md) |  | [optional] 
**Status** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetStatus**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetStatus.md) |  | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSet

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSet() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSet`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSet instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSet object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSet`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSet object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSet) GetApiVersion() string`

GetApiVersion returns the ApiVersion field if non-nil, zero value otherwise.

### GetApiVersionOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSet) GetApiVersionOk() (*string, bool)`

GetApiVersionOk returns a tuple with the ApiVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSet) SetApiVersion(v string)`

SetApiVersion sets ApiVersion field to given value.

### HasApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSet) HasApiVersion() bool`

HasApiVersion returns a boolean if a field has been set.

### GetKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSet) GetKind() string`

GetKind returns the Kind field if non-nil, zero value otherwise.

### GetKindOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSet) GetKindOk() (*string, bool)`

GetKindOk returns a tuple with the Kind field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSet) SetKind(v string)`

SetKind sets Kind field to given value.

### HasKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSet) HasKind() bool`

HasKind returns a boolean if a field has been set.

### GetMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSet) GetMetadata() V1ObjectMeta`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSet) GetMetadataOk() (*V1ObjectMeta, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSet) SetMetadata(v V1ObjectMeta)`

SetMetadata sets Metadata field to given value.

### HasMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSet) HasMetadata() bool`

HasMetadata returns a boolean if a field has been set.

### GetSpec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSet) GetSpec() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec`

GetSpec returns the Spec field if non-nil, zero value otherwise.

### GetSpecOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSet) GetSpecOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec, bool)`

GetSpecOk returns a tuple with the Spec field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSpec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSet) SetSpec(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec)`

SetSpec sets Spec field to given value.

### HasSpec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSet) HasSpec() bool`

HasSpec returns a boolean if a field has been set.

### GetStatus

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSet) GetStatus() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSet) GetStatusOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSet) SetStatus(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSet) HasStatus() bool`

HasStatus returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


