# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ApiVersion** | Pointer to **string** | APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources | [optional] 
**Data** | Pointer to **map[string]string** | the value should be base64 encoded | [optional] 
**InstanceName** | **string** |  | 
**Kind** | Pointer to **string** | Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds | [optional] 
**Location** | **string** |  | 
**Metadata** | Pointer to [**V1ObjectMeta**](V1ObjectMeta.md) |  | [optional] 
**PoolMemberRef** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberReference**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberReference.md) |  | [optional] 
**Spec** | Pointer to **map[string]interface{}** |  | [optional] 
**Status** | Pointer to **map[string]interface{}** |  | [optional] 
**Tolerations** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Toleration**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Toleration.md) |  | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret(instanceName string, location string, ) *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SecretWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SecretWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SecretWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret) GetApiVersion() string`

GetApiVersion returns the ApiVersion field if non-nil, zero value otherwise.

### GetApiVersionOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret) GetApiVersionOk() (*string, bool)`

GetApiVersionOk returns a tuple with the ApiVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret) SetApiVersion(v string)`

SetApiVersion sets ApiVersion field to given value.

### HasApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret) HasApiVersion() bool`

HasApiVersion returns a boolean if a field has been set.

### GetData

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret) GetData() map[string]string`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret) GetDataOk() (*map[string]string, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret) SetData(v map[string]string)`

SetData sets Data field to given value.

### HasData

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret) HasData() bool`

HasData returns a boolean if a field has been set.

### GetInstanceName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret) GetInstanceName() string`

GetInstanceName returns the InstanceName field if non-nil, zero value otherwise.

### GetInstanceNameOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret) GetInstanceNameOk() (*string, bool)`

GetInstanceNameOk returns a tuple with the InstanceName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInstanceName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret) SetInstanceName(v string)`

SetInstanceName sets InstanceName field to given value.


### GetKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret) GetKind() string`

GetKind returns the Kind field if non-nil, zero value otherwise.

### GetKindOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret) GetKindOk() (*string, bool)`

GetKindOk returns a tuple with the Kind field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret) SetKind(v string)`

SetKind sets Kind field to given value.

### HasKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret) HasKind() bool`

HasKind returns a boolean if a field has been set.

### GetLocation

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret) GetLocation() string`

GetLocation returns the Location field if non-nil, zero value otherwise.

### GetLocationOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret) GetLocationOk() (*string, bool)`

GetLocationOk returns a tuple with the Location field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLocation

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret) SetLocation(v string)`

SetLocation sets Location field to given value.


### GetMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret) GetMetadata() V1ObjectMeta`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret) GetMetadataOk() (*V1ObjectMeta, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret) SetMetadata(v V1ObjectMeta)`

SetMetadata sets Metadata field to given value.

### HasMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret) HasMetadata() bool`

HasMetadata returns a boolean if a field has been set.

### GetPoolMemberRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret) GetPoolMemberRef() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberReference`

GetPoolMemberRef returns the PoolMemberRef field if non-nil, zero value otherwise.

### GetPoolMemberRefOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret) GetPoolMemberRefOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberReference, bool)`

GetPoolMemberRefOk returns a tuple with the PoolMemberRef field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPoolMemberRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret) SetPoolMemberRef(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberReference)`

SetPoolMemberRef sets PoolMemberRef field to given value.

### HasPoolMemberRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret) HasPoolMemberRef() bool`

HasPoolMemberRef returns a boolean if a field has been set.

### GetSpec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret) GetSpec() map[string]interface{}`

GetSpec returns the Spec field if non-nil, zero value otherwise.

### GetSpecOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret) GetSpecOk() (*map[string]interface{}, bool)`

GetSpecOk returns a tuple with the Spec field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSpec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret) SetSpec(v map[string]interface{})`

SetSpec sets Spec field to given value.

### HasSpec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret) HasSpec() bool`

HasSpec returns a boolean if a field has been set.

### GetStatus

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret) GetStatus() map[string]interface{}`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret) GetStatusOk() (*map[string]interface{}, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret) SetStatus(v map[string]interface{})`

SetStatus sets Status field to given value.

### HasStatus

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetTolerations

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret) GetTolerations() []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Toleration`

GetTolerations returns the Tolerations field if non-nil, zero value otherwise.

### GetTolerationsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret) GetTolerationsOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Toleration, bool)`

GetTolerationsOk returns a tuple with the Tolerations field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTolerations

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret) SetTolerations(v []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Toleration)`

SetTolerations sets Tolerations field to given value.

### HasTolerations

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Secret) HasTolerations() bool`

HasTolerations returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


