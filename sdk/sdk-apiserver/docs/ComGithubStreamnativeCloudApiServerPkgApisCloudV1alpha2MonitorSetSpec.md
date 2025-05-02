# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2MonitorSetSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AvailabilityMode** | Pointer to **string** | AvailabilityMode decides whether servers should be in one zone or multiple zones If unspecified, defaults to zonal. | [optional] 
**PoolMemberRef** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2PoolMemberReference**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2PoolMemberReference.md) |  | [optional] 
**Replicas** | Pointer to **int32** | Replicas is the desired number of monitoring servers. If unspecified, defaults to 1. | [optional] 
**Selector** | Pointer to [**V1LabelSelector**](V1LabelSelector.md) |  | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2MonitorSetSpec

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2MonitorSetSpec() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2MonitorSetSpec`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2MonitorSetSpec instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2MonitorSetSpec object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2MonitorSetSpecWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2MonitorSetSpecWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2MonitorSetSpec`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2MonitorSetSpecWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2MonitorSetSpec object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAvailabilityMode

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2MonitorSetSpec) GetAvailabilityMode() string`

GetAvailabilityMode returns the AvailabilityMode field if non-nil, zero value otherwise.

### GetAvailabilityModeOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2MonitorSetSpec) GetAvailabilityModeOk() (*string, bool)`

GetAvailabilityModeOk returns a tuple with the AvailabilityMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAvailabilityMode

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2MonitorSetSpec) SetAvailabilityMode(v string)`

SetAvailabilityMode sets AvailabilityMode field to given value.

### HasAvailabilityMode

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2MonitorSetSpec) HasAvailabilityMode() bool`

HasAvailabilityMode returns a boolean if a field has been set.

### GetPoolMemberRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2MonitorSetSpec) GetPoolMemberRef() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2PoolMemberReference`

GetPoolMemberRef returns the PoolMemberRef field if non-nil, zero value otherwise.

### GetPoolMemberRefOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2MonitorSetSpec) GetPoolMemberRefOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2PoolMemberReference, bool)`

GetPoolMemberRefOk returns a tuple with the PoolMemberRef field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPoolMemberRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2MonitorSetSpec) SetPoolMemberRef(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2PoolMemberReference)`

SetPoolMemberRef sets PoolMemberRef field to given value.

### HasPoolMemberRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2MonitorSetSpec) HasPoolMemberRef() bool`

HasPoolMemberRef returns a boolean if a field has been set.

### GetReplicas

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2MonitorSetSpec) GetReplicas() int32`

GetReplicas returns the Replicas field if non-nil, zero value otherwise.

### GetReplicasOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2MonitorSetSpec) GetReplicasOk() (*int32, bool)`

GetReplicasOk returns a tuple with the Replicas field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReplicas

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2MonitorSetSpec) SetReplicas(v int32)`

SetReplicas sets Replicas field to given value.

### HasReplicas

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2MonitorSetSpec) HasReplicas() bool`

HasReplicas returns a boolean if a field has been set.

### GetSelector

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2MonitorSetSpec) GetSelector() V1LabelSelector`

GetSelector returns the Selector field if non-nil, zero value otherwise.

### GetSelectorOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2MonitorSetSpec) GetSelectorOk() (*V1LabelSelector, bool)`

GetSelectorOk returns a tuple with the Selector field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSelector

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2MonitorSetSpec) SetSelector(v V1LabelSelector)`

SetSelector sets Selector field to given value.

### HasSelector

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2MonitorSetSpec) HasSelector() bool`

HasSelector returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


