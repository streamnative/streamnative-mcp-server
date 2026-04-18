# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Location** | **string** | Location is the deployment location/region | [default to ""]
**PoolMemberRef** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberReference**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberReference.md) |  | [optional] 
**Config** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterConfig**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterConfig.md) |  | [optional] 
**DisplayName** | Pointer to **string** | DisplayName is a user-friendly name for the cluster | [optional] 
**EndpointAccess** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1EndpointAccess**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1EndpointAccess.md) | EndpointAccess defines endpoint access configuration for the Kafka cluster | [optional] 
**InstanceName** | **string** | InstanceName is the name of the Kafka instance this cluster belongs to | [default to ""]
**Mtls** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterMTLS**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterMTLS.md) |  | [optional] 
**ClusterProfile** | Pointer to **string** | ClusterProfile defines the cluster profile. cost-optimized is for cost-optimized clusters and latency-optimized is for latency-optimized clusters. | [optional] 
**Catalogs** | Pointer to **[]string** | Catalogs is a list of catalog names to associate with the Kafka cluster. These catalogs are used by the compaction scheduler for lakehouse table management. | [optional] 
**Rbac** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterRBAC**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterRBAC.md) |  | [optional] 
**ThroughputUnit** | Pointer to **int32** | ThroughputUnit defines the requested broker throughput capacity in RTU. The KafkaCluster controller dynamically translates this value into broker resources and replicas, unless KafkaClusterBackendConfig.spec.broker resources, replicas, or data explicitly override them. If omitted, the default is 1 RTU. | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec(location string, instanceName string, ) *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpecWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpecWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpecWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetLocation

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec) GetLocation() string`

GetLocation returns the Location field if non-nil, zero value otherwise.

### GetLocationOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec) GetLocationOk() (*string, bool)`

GetLocationOk returns a tuple with the Location field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLocation

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec) SetLocation(v string)`

SetLocation sets Location field to given value.


### GetPoolMemberRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec) GetPoolMemberRef() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberReference`

GetPoolMemberRef returns the PoolMemberRef field if non-nil, zero value otherwise.

### GetPoolMemberRefOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec) GetPoolMemberRefOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberReference, bool)`

GetPoolMemberRefOk returns a tuple with the PoolMemberRef field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPoolMemberRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec) SetPoolMemberRef(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberReference)`

SetPoolMemberRef sets PoolMemberRef field to given value.

### HasPoolMemberRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec) HasPoolMemberRef() bool`

HasPoolMemberRef returns a boolean if a field has been set.

### GetConfig

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec) GetConfig() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterConfig`

GetConfig returns the Config field if non-nil, zero value otherwise.

### GetConfigOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec) GetConfigOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterConfig, bool)`

GetConfigOk returns a tuple with the Config field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConfig

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec) SetConfig(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterConfig)`

SetConfig sets Config field to given value.

### HasConfig

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec) HasConfig() bool`

HasConfig returns a boolean if a field has been set.

### GetDisplayName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec) GetDisplayName() string`

GetDisplayName returns the DisplayName field if non-nil, zero value otherwise.

### GetDisplayNameOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec) GetDisplayNameOk() (*string, bool)`

GetDisplayNameOk returns a tuple with the DisplayName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDisplayName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec) SetDisplayName(v string)`

SetDisplayName sets DisplayName field to given value.

### HasDisplayName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec) HasDisplayName() bool`

HasDisplayName returns a boolean if a field has been set.

### GetEndpointAccess

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec) GetEndpointAccess() []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1EndpointAccess`

GetEndpointAccess returns the EndpointAccess field if non-nil, zero value otherwise.

### GetEndpointAccessOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec) GetEndpointAccessOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1EndpointAccess, bool)`

GetEndpointAccessOk returns a tuple with the EndpointAccess field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndpointAccess

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec) SetEndpointAccess(v []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1EndpointAccess)`

SetEndpointAccess sets EndpointAccess field to given value.

### HasEndpointAccess

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec) HasEndpointAccess() bool`

HasEndpointAccess returns a boolean if a field has been set.

### GetInstanceName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec) GetInstanceName() string`

GetInstanceName returns the InstanceName field if non-nil, zero value otherwise.

### GetInstanceNameOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec) GetInstanceNameOk() (*string, bool)`

GetInstanceNameOk returns a tuple with the InstanceName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInstanceName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec) SetInstanceName(v string)`

SetInstanceName sets InstanceName field to given value.


### GetMtls

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec) GetMtls() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterMTLS`

GetMtls returns the Mtls field if non-nil, zero value otherwise.

### GetMtlsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec) GetMtlsOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterMTLS, bool)`

GetMtlsOk returns a tuple with the Mtls field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMtls

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec) SetMtls(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterMTLS)`

SetMtls sets Mtls field to given value.

### HasMtls

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec) HasMtls() bool`

HasMtls returns a boolean if a field has been set.

### GetClusterProfile

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec) GetClusterProfile() string`

GetClusterProfile returns the ClusterProfile field if non-nil, zero value otherwise.

### GetClusterProfileOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec) GetClusterProfileOk() (*string, bool)`

GetClusterProfileOk returns a tuple with the ClusterProfile field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClusterProfile

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec) SetClusterProfile(v string)`

SetClusterProfile sets ClusterProfile field to given value.

### HasClusterProfile

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec) HasClusterProfile() bool`

HasClusterProfile returns a boolean if a field has been set.

### GetCatalogs

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec) GetCatalogs() []string`

GetCatalogs returns the Catalogs field if non-nil, zero value otherwise.

### GetCatalogsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec) GetCatalogsOk() (*[]string, bool)`

GetCatalogsOk returns a tuple with the Catalogs field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCatalogs

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec) SetCatalogs(v []string)`

SetCatalogs sets Catalogs field to given value.

### HasCatalogs

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec) HasCatalogs() bool`

HasCatalogs returns a boolean if a field has been set.

### GetRbac

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec) GetRbac() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterRBAC`

GetRbac returns the Rbac field if non-nil, zero value otherwise.

### GetRbacOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec) GetRbacOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterRBAC, bool)`

GetRbacOk returns a tuple with the Rbac field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRbac

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec) SetRbac(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterRBAC)`

SetRbac sets Rbac field to given value.

### HasRbac

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec) HasRbac() bool`

HasRbac returns a boolean if a field has been set.

### GetThroughputUnit

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec) GetThroughputUnit() int32`

GetThroughputUnit returns the ThroughputUnit field if non-nil, zero value otherwise.

### GetThroughputUnitOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec) GetThroughputUnitOk() (*int32, bool)`

GetThroughputUnitOk returns a tuple with the ThroughputUnit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetThroughputUnit

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec) SetThroughputUnit(v int32)`

SetThroughputUnit sets ThroughputUnit field to given value.

### HasThroughputUnit

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterSpec) HasThroughputUnit() bool`

HasThroughputUnit returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


