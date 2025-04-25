# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AvailabilityMode** | Pointer to **string** | AvailabilityMode decides whether servers should be in one zone or multiple zones If unspecified, defaults to zonal. | [optional] 
**Image** | Pointer to **string** | Image name is the name of the image to deploy. | [optional] 
**ImagePullPolicy** | Pointer to **string** | Image pull policy, one of Always, Never, IfNotPresent, default to Always. | [optional] 
**PoolMemberRef** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2PoolMemberReference**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2PoolMemberReference.md) |  | [optional] 
**Replicas** | Pointer to **int32** | Replicas is the desired number of ZooKeeper servers. If unspecified, defaults to 1. | [optional] 
**ResourceSpec** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperResourceSpec**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperResourceSpec.md) |  | [optional] 
**Resources** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2DefaultNodeResource**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2DefaultNodeResource.md) |  | [optional] 
**Sharing** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2SharingConfig**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2SharingConfig.md) |  | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpecWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpecWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpecWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAvailabilityMode

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec) GetAvailabilityMode() string`

GetAvailabilityMode returns the AvailabilityMode field if non-nil, zero value otherwise.

### GetAvailabilityModeOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec) GetAvailabilityModeOk() (*string, bool)`

GetAvailabilityModeOk returns a tuple with the AvailabilityMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAvailabilityMode

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec) SetAvailabilityMode(v string)`

SetAvailabilityMode sets AvailabilityMode field to given value.

### HasAvailabilityMode

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec) HasAvailabilityMode() bool`

HasAvailabilityMode returns a boolean if a field has been set.

### GetImage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec) GetImage() string`

GetImage returns the Image field if non-nil, zero value otherwise.

### GetImageOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec) GetImageOk() (*string, bool)`

GetImageOk returns a tuple with the Image field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec) SetImage(v string)`

SetImage sets Image field to given value.

### HasImage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec) HasImage() bool`

HasImage returns a boolean if a field has been set.

### GetImagePullPolicy

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec) GetImagePullPolicy() string`

GetImagePullPolicy returns the ImagePullPolicy field if non-nil, zero value otherwise.

### GetImagePullPolicyOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec) GetImagePullPolicyOk() (*string, bool)`

GetImagePullPolicyOk returns a tuple with the ImagePullPolicy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImagePullPolicy

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec) SetImagePullPolicy(v string)`

SetImagePullPolicy sets ImagePullPolicy field to given value.

### HasImagePullPolicy

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec) HasImagePullPolicy() bool`

HasImagePullPolicy returns a boolean if a field has been set.

### GetPoolMemberRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec) GetPoolMemberRef() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2PoolMemberReference`

GetPoolMemberRef returns the PoolMemberRef field if non-nil, zero value otherwise.

### GetPoolMemberRefOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec) GetPoolMemberRefOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2PoolMemberReference, bool)`

GetPoolMemberRefOk returns a tuple with the PoolMemberRef field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPoolMemberRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec) SetPoolMemberRef(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2PoolMemberReference)`

SetPoolMemberRef sets PoolMemberRef field to given value.

### HasPoolMemberRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec) HasPoolMemberRef() bool`

HasPoolMemberRef returns a boolean if a field has been set.

### GetReplicas

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec) GetReplicas() int32`

GetReplicas returns the Replicas field if non-nil, zero value otherwise.

### GetReplicasOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec) GetReplicasOk() (*int32, bool)`

GetReplicasOk returns a tuple with the Replicas field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReplicas

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec) SetReplicas(v int32)`

SetReplicas sets Replicas field to given value.

### HasReplicas

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec) HasReplicas() bool`

HasReplicas returns a boolean if a field has been set.

### GetResourceSpec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec) GetResourceSpec() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperResourceSpec`

GetResourceSpec returns the ResourceSpec field if non-nil, zero value otherwise.

### GetResourceSpecOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec) GetResourceSpecOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperResourceSpec, bool)`

GetResourceSpecOk returns a tuple with the ResourceSpec field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResourceSpec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec) SetResourceSpec(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperResourceSpec)`

SetResourceSpec sets ResourceSpec field to given value.

### HasResourceSpec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec) HasResourceSpec() bool`

HasResourceSpec returns a boolean if a field has been set.

### GetResources

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec) GetResources() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2DefaultNodeResource`

GetResources returns the Resources field if non-nil, zero value otherwise.

### GetResourcesOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec) GetResourcesOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2DefaultNodeResource, bool)`

GetResourcesOk returns a tuple with the Resources field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResources

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec) SetResources(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2DefaultNodeResource)`

SetResources sets Resources field to given value.

### HasResources

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec) HasResources() bool`

HasResources returns a boolean if a field has been set.

### GetSharing

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec) GetSharing() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2SharingConfig`

GetSharing returns the Sharing field if non-nil, zero value otherwise.

### GetSharingOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec) GetSharingOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2SharingConfig, bool)`

GetSharingOk returns a tuple with the Sharing field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSharing

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec) SetSharing(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2SharingConfig)`

SetSharing sets Sharing field to given value.

### HasSharing

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetSpec) HasSharing() bool`

HasSharing returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


