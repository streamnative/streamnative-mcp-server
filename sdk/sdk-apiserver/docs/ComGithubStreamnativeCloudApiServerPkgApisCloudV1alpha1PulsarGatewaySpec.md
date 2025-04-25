# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewaySpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Access** | Pointer to **string** | Access is the access type of the pulsar gateway, available values are public or private. It is immutable, with the default value public. | [optional] 
**Domains** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Domain**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Domain.md) | Domains is the list of domain suffix that the pulsar gateway will serve. This is automatically generated based on the PulsarGateway name and PoolMember domain. | [optional] 
**PoolMemberRef** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberReference**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberReference.md) |  | [optional] 
**PrivateService** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PrivateService**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PrivateService.md) |  | [optional] 
**TopologyAware** | Pointer to **map[string]interface{}** | TopologyAware is the configuration of the topology aware feature of the pulsar gateway. | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewaySpec

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewaySpec() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewaySpec`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewaySpec instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewaySpec object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewaySpecWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewaySpecWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewaySpec`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewaySpecWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewaySpec object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAccess

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewaySpec) GetAccess() string`

GetAccess returns the Access field if non-nil, zero value otherwise.

### GetAccessOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewaySpec) GetAccessOk() (*string, bool)`

GetAccessOk returns a tuple with the Access field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccess

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewaySpec) SetAccess(v string)`

SetAccess sets Access field to given value.

### HasAccess

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewaySpec) HasAccess() bool`

HasAccess returns a boolean if a field has been set.

### GetDomains

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewaySpec) GetDomains() []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Domain`

GetDomains returns the Domains field if non-nil, zero value otherwise.

### GetDomainsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewaySpec) GetDomainsOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Domain, bool)`

GetDomainsOk returns a tuple with the Domains field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDomains

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewaySpec) SetDomains(v []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Domain)`

SetDomains sets Domains field to given value.

### HasDomains

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewaySpec) HasDomains() bool`

HasDomains returns a boolean if a field has been set.

### GetPoolMemberRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewaySpec) GetPoolMemberRef() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberReference`

GetPoolMemberRef returns the PoolMemberRef field if non-nil, zero value otherwise.

### GetPoolMemberRefOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewaySpec) GetPoolMemberRefOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberReference, bool)`

GetPoolMemberRefOk returns a tuple with the PoolMemberRef field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPoolMemberRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewaySpec) SetPoolMemberRef(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberReference)`

SetPoolMemberRef sets PoolMemberRef field to given value.

### HasPoolMemberRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewaySpec) HasPoolMemberRef() bool`

HasPoolMemberRef returns a boolean if a field has been set.

### GetPrivateService

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewaySpec) GetPrivateService() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PrivateService`

GetPrivateService returns the PrivateService field if non-nil, zero value otherwise.

### GetPrivateServiceOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewaySpec) GetPrivateServiceOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PrivateService, bool)`

GetPrivateServiceOk returns a tuple with the PrivateService field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrivateService

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewaySpec) SetPrivateService(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PrivateService)`

SetPrivateService sets PrivateService field to given value.

### HasPrivateService

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewaySpec) HasPrivateService() bool`

HasPrivateService returns a boolean if a field has been set.

### GetTopologyAware

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewaySpec) GetTopologyAware() map[string]interface{}`

GetTopologyAware returns the TopologyAware field if non-nil, zero value otherwise.

### GetTopologyAwareOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewaySpec) GetTopologyAwareOk() (*map[string]interface{}, bool)`

GetTopologyAwareOk returns a tuple with the TopologyAware field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTopologyAware

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewaySpec) SetTopologyAware(v map[string]interface{})`

SetTopologyAware sets TopologyAware field to given value.

### HasTopologyAware

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewaySpec) HasTopologyAware() bool`

HasTopologyAware returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


