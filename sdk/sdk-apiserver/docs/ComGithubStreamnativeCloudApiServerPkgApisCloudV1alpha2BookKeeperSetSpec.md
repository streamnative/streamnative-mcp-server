# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AvailabilityMode** | Pointer to **string** | AvailabilityMode decides whether servers should be in one zone or multiple zones If unspecified, defaults to zonal. | [optional] 
**Image** | Pointer to **string** | Image name is the name of the image to deploy. | [optional] 
**PoolMemberRef** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2PoolMemberReference**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2PoolMemberReference.md) |  | [optional] 
**Replicas** | Pointer to **int32** | Replicas is the desired number of BookKeeper servers. If unspecified, defaults to 3. | [optional] 
**ResourceSpec** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperResourceSpec**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperResourceSpec.md) |  | [optional] 
**Resources** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookkeeperNodeResource**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookkeeperNodeResource.md) |  | [optional] 
**Sharing** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2SharingConfig**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2SharingConfig.md) |  | [optional] 
**ZooKeeperSetRef** | [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetReference**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetReference.md) |  | 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpec

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpec(zooKeeperSetRef ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetReference, ) *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpec`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpec instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpec object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpecWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpecWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpec`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpecWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpec object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAvailabilityMode

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpec) GetAvailabilityMode() string`

GetAvailabilityMode returns the AvailabilityMode field if non-nil, zero value otherwise.

### GetAvailabilityModeOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpec) GetAvailabilityModeOk() (*string, bool)`

GetAvailabilityModeOk returns a tuple with the AvailabilityMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAvailabilityMode

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpec) SetAvailabilityMode(v string)`

SetAvailabilityMode sets AvailabilityMode field to given value.

### HasAvailabilityMode

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpec) HasAvailabilityMode() bool`

HasAvailabilityMode returns a boolean if a field has been set.

### GetImage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpec) GetImage() string`

GetImage returns the Image field if non-nil, zero value otherwise.

### GetImageOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpec) GetImageOk() (*string, bool)`

GetImageOk returns a tuple with the Image field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpec) SetImage(v string)`

SetImage sets Image field to given value.

### HasImage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpec) HasImage() bool`

HasImage returns a boolean if a field has been set.

### GetPoolMemberRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpec) GetPoolMemberRef() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2PoolMemberReference`

GetPoolMemberRef returns the PoolMemberRef field if non-nil, zero value otherwise.

### GetPoolMemberRefOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpec) GetPoolMemberRefOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2PoolMemberReference, bool)`

GetPoolMemberRefOk returns a tuple with the PoolMemberRef field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPoolMemberRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpec) SetPoolMemberRef(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2PoolMemberReference)`

SetPoolMemberRef sets PoolMemberRef field to given value.

### HasPoolMemberRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpec) HasPoolMemberRef() bool`

HasPoolMemberRef returns a boolean if a field has been set.

### GetReplicas

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpec) GetReplicas() int32`

GetReplicas returns the Replicas field if non-nil, zero value otherwise.

### GetReplicasOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpec) GetReplicasOk() (*int32, bool)`

GetReplicasOk returns a tuple with the Replicas field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReplicas

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpec) SetReplicas(v int32)`

SetReplicas sets Replicas field to given value.

### HasReplicas

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpec) HasReplicas() bool`

HasReplicas returns a boolean if a field has been set.

### GetResourceSpec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpec) GetResourceSpec() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperResourceSpec`

GetResourceSpec returns the ResourceSpec field if non-nil, zero value otherwise.

### GetResourceSpecOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpec) GetResourceSpecOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperResourceSpec, bool)`

GetResourceSpecOk returns a tuple with the ResourceSpec field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResourceSpec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpec) SetResourceSpec(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperResourceSpec)`

SetResourceSpec sets ResourceSpec field to given value.

### HasResourceSpec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpec) HasResourceSpec() bool`

HasResourceSpec returns a boolean if a field has been set.

### GetResources

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpec) GetResources() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookkeeperNodeResource`

GetResources returns the Resources field if non-nil, zero value otherwise.

### GetResourcesOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpec) GetResourcesOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookkeeperNodeResource, bool)`

GetResourcesOk returns a tuple with the Resources field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResources

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpec) SetResources(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookkeeperNodeResource)`

SetResources sets Resources field to given value.

### HasResources

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpec) HasResources() bool`

HasResources returns a boolean if a field has been set.

### GetSharing

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpec) GetSharing() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2SharingConfig`

GetSharing returns the Sharing field if non-nil, zero value otherwise.

### GetSharingOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpec) GetSharingOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2SharingConfig, bool)`

GetSharingOk returns a tuple with the Sharing field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSharing

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpec) SetSharing(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2SharingConfig)`

SetSharing sets Sharing field to given value.

### HasSharing

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpec) HasSharing() bool`

HasSharing returns a boolean if a field has been set.

### GetZooKeeperSetRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpec) GetZooKeeperSetRef() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetReference`

GetZooKeeperSetRef returns the ZooKeeperSetRef field if non-nil, zero value otherwise.

### GetZooKeeperSetRefOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpec) GetZooKeeperSetRefOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetReference, bool)`

GetZooKeeperSetRefOk returns a tuple with the ZooKeeperSetRef field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetZooKeeperSetRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetSpec) SetZooKeeperSetRef(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetReference)`

SetZooKeeperSetRef sets ZooKeeperSetRef field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


