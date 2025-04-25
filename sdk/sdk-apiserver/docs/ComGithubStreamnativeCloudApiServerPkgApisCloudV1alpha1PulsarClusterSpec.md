# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**BookKeeperSetRef** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeperSetReference**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeperSetReference.md) |  | [optional] 
**Bookkeeper** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeper**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeper.md) |  | [optional] 
**Broker** | [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Broker**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Broker.md) |  | 
**Config** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Config**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Config.md) |  | [optional] 
**DisplayName** | Pointer to **string** | Name to display to our users | [optional] 
**EndpointAccess** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1EndpointAccess**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1EndpointAccess.md) |  | [optional] 
**InstanceName** | **string** |  | 
**Location** | **string** |  | 
**MaintenanceWindow** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MaintenanceWindow**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MaintenanceWindow.md) |  | [optional] 
**PoolMemberRef** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberReference**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberReference.md) |  | [optional] 
**ReleaseChannel** | Pointer to **string** |  | [optional] 
**ServiceEndpoints** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarServiceEndpoint**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarServiceEndpoint.md) |  | [optional] 
**Tolerations** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Toleration**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Toleration.md) |  | [optional] 
**ZooKeeperSetRef** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ZooKeeperSetReference**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ZooKeeperSetReference.md) |  | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec(broker ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Broker, instanceName string, location string, ) *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpecWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpecWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpecWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetBookKeeperSetRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) GetBookKeeperSetRef() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeperSetReference`

GetBookKeeperSetRef returns the BookKeeperSetRef field if non-nil, zero value otherwise.

### GetBookKeeperSetRefOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) GetBookKeeperSetRefOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeperSetReference, bool)`

GetBookKeeperSetRefOk returns a tuple with the BookKeeperSetRef field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBookKeeperSetRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) SetBookKeeperSetRef(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeperSetReference)`

SetBookKeeperSetRef sets BookKeeperSetRef field to given value.

### HasBookKeeperSetRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) HasBookKeeperSetRef() bool`

HasBookKeeperSetRef returns a boolean if a field has been set.

### GetBookkeeper

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) GetBookkeeper() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeper`

GetBookkeeper returns the Bookkeeper field if non-nil, zero value otherwise.

### GetBookkeeperOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) GetBookkeeperOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeper, bool)`

GetBookkeeperOk returns a tuple with the Bookkeeper field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBookkeeper

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) SetBookkeeper(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeper)`

SetBookkeeper sets Bookkeeper field to given value.

### HasBookkeeper

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) HasBookkeeper() bool`

HasBookkeeper returns a boolean if a field has been set.

### GetBroker

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) GetBroker() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Broker`

GetBroker returns the Broker field if non-nil, zero value otherwise.

### GetBrokerOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) GetBrokerOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Broker, bool)`

GetBrokerOk returns a tuple with the Broker field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBroker

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) SetBroker(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Broker)`

SetBroker sets Broker field to given value.


### GetConfig

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) GetConfig() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Config`

GetConfig returns the Config field if non-nil, zero value otherwise.

### GetConfigOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) GetConfigOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Config, bool)`

GetConfigOk returns a tuple with the Config field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConfig

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) SetConfig(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Config)`

SetConfig sets Config field to given value.

### HasConfig

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) HasConfig() bool`

HasConfig returns a boolean if a field has been set.

### GetDisplayName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) GetDisplayName() string`

GetDisplayName returns the DisplayName field if non-nil, zero value otherwise.

### GetDisplayNameOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) GetDisplayNameOk() (*string, bool)`

GetDisplayNameOk returns a tuple with the DisplayName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDisplayName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) SetDisplayName(v string)`

SetDisplayName sets DisplayName field to given value.

### HasDisplayName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) HasDisplayName() bool`

HasDisplayName returns a boolean if a field has been set.

### GetEndpointAccess

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) GetEndpointAccess() []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1EndpointAccess`

GetEndpointAccess returns the EndpointAccess field if non-nil, zero value otherwise.

### GetEndpointAccessOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) GetEndpointAccessOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1EndpointAccess, bool)`

GetEndpointAccessOk returns a tuple with the EndpointAccess field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndpointAccess

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) SetEndpointAccess(v []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1EndpointAccess)`

SetEndpointAccess sets EndpointAccess field to given value.

### HasEndpointAccess

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) HasEndpointAccess() bool`

HasEndpointAccess returns a boolean if a field has been set.

### GetInstanceName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) GetInstanceName() string`

GetInstanceName returns the InstanceName field if non-nil, zero value otherwise.

### GetInstanceNameOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) GetInstanceNameOk() (*string, bool)`

GetInstanceNameOk returns a tuple with the InstanceName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInstanceName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) SetInstanceName(v string)`

SetInstanceName sets InstanceName field to given value.


### GetLocation

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) GetLocation() string`

GetLocation returns the Location field if non-nil, zero value otherwise.

### GetLocationOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) GetLocationOk() (*string, bool)`

GetLocationOk returns a tuple with the Location field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLocation

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) SetLocation(v string)`

SetLocation sets Location field to given value.


### GetMaintenanceWindow

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) GetMaintenanceWindow() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MaintenanceWindow`

GetMaintenanceWindow returns the MaintenanceWindow field if non-nil, zero value otherwise.

### GetMaintenanceWindowOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) GetMaintenanceWindowOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MaintenanceWindow, bool)`

GetMaintenanceWindowOk returns a tuple with the MaintenanceWindow field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaintenanceWindow

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) SetMaintenanceWindow(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MaintenanceWindow)`

SetMaintenanceWindow sets MaintenanceWindow field to given value.

### HasMaintenanceWindow

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) HasMaintenanceWindow() bool`

HasMaintenanceWindow returns a boolean if a field has been set.

### GetPoolMemberRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) GetPoolMemberRef() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberReference`

GetPoolMemberRef returns the PoolMemberRef field if non-nil, zero value otherwise.

### GetPoolMemberRefOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) GetPoolMemberRefOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberReference, bool)`

GetPoolMemberRefOk returns a tuple with the PoolMemberRef field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPoolMemberRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) SetPoolMemberRef(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberReference)`

SetPoolMemberRef sets PoolMemberRef field to given value.

### HasPoolMemberRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) HasPoolMemberRef() bool`

HasPoolMemberRef returns a boolean if a field has been set.

### GetReleaseChannel

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) GetReleaseChannel() string`

GetReleaseChannel returns the ReleaseChannel field if non-nil, zero value otherwise.

### GetReleaseChannelOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) GetReleaseChannelOk() (*string, bool)`

GetReleaseChannelOk returns a tuple with the ReleaseChannel field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReleaseChannel

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) SetReleaseChannel(v string)`

SetReleaseChannel sets ReleaseChannel field to given value.

### HasReleaseChannel

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) HasReleaseChannel() bool`

HasReleaseChannel returns a boolean if a field has been set.

### GetServiceEndpoints

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) GetServiceEndpoints() []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarServiceEndpoint`

GetServiceEndpoints returns the ServiceEndpoints field if non-nil, zero value otherwise.

### GetServiceEndpointsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) GetServiceEndpointsOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarServiceEndpoint, bool)`

GetServiceEndpointsOk returns a tuple with the ServiceEndpoints field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServiceEndpoints

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) SetServiceEndpoints(v []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarServiceEndpoint)`

SetServiceEndpoints sets ServiceEndpoints field to given value.

### HasServiceEndpoints

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) HasServiceEndpoints() bool`

HasServiceEndpoints returns a boolean if a field has been set.

### GetTolerations

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) GetTolerations() []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Toleration`

GetTolerations returns the Tolerations field if non-nil, zero value otherwise.

### GetTolerationsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) GetTolerationsOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Toleration, bool)`

GetTolerationsOk returns a tuple with the Tolerations field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTolerations

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) SetTolerations(v []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Toleration)`

SetTolerations sets Tolerations field to given value.

### HasTolerations

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) HasTolerations() bool`

HasTolerations returns a boolean if a field has been set.

### GetZooKeeperSetRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) GetZooKeeperSetRef() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ZooKeeperSetReference`

GetZooKeeperSetRef returns the ZooKeeperSetRef field if non-nil, zero value otherwise.

### GetZooKeeperSetRefOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) GetZooKeeperSetRefOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ZooKeeperSetReference, bool)`

GetZooKeeperSetRefOk returns a tuple with the ZooKeeperSetRef field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetZooKeeperSetRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) SetZooKeeperSetRef(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ZooKeeperSetReference)`

SetZooKeeperSetRef sets ZooKeeperSetRef field to given value.

### HasZooKeeperSetRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec) HasZooKeeperSetRef() bool`

HasZooKeeperSetRef returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


