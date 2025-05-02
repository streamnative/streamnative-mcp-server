# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CloudConnectionName** | Pointer to **string** | CloudConnectionName references to the CloudConnection object | [optional] 
**DefaultGateway** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Gateway**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Gateway.md) |  | [optional] 
**Dns** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1DNS**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1DNS.md) |  | [optional] 
**Network** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Network**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Network.md) |  | [optional] 
**Region** | Pointer to **string** | Region defines in which region will resources be deployed | [optional] 
**Zone** | Pointer to **string** | Zone defines in which availability zone will resources be deployed If specified, the cloud environment will be zonal. Default to unspecified and regional cloud environment | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentSpec

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentSpec() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentSpec`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentSpec instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentSpec object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentSpecWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentSpecWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentSpec`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentSpecWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentSpec object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCloudConnectionName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentSpec) GetCloudConnectionName() string`

GetCloudConnectionName returns the CloudConnectionName field if non-nil, zero value otherwise.

### GetCloudConnectionNameOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentSpec) GetCloudConnectionNameOk() (*string, bool)`

GetCloudConnectionNameOk returns a tuple with the CloudConnectionName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCloudConnectionName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentSpec) SetCloudConnectionName(v string)`

SetCloudConnectionName sets CloudConnectionName field to given value.

### HasCloudConnectionName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentSpec) HasCloudConnectionName() bool`

HasCloudConnectionName returns a boolean if a field has been set.

### GetDefaultGateway

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentSpec) GetDefaultGateway() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Gateway`

GetDefaultGateway returns the DefaultGateway field if non-nil, zero value otherwise.

### GetDefaultGatewayOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentSpec) GetDefaultGatewayOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Gateway, bool)`

GetDefaultGatewayOk returns a tuple with the DefaultGateway field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDefaultGateway

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentSpec) SetDefaultGateway(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Gateway)`

SetDefaultGateway sets DefaultGateway field to given value.

### HasDefaultGateway

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentSpec) HasDefaultGateway() bool`

HasDefaultGateway returns a boolean if a field has been set.

### GetDns

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentSpec) GetDns() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1DNS`

GetDns returns the Dns field if non-nil, zero value otherwise.

### GetDnsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentSpec) GetDnsOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1DNS, bool)`

GetDnsOk returns a tuple with the Dns field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDns

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentSpec) SetDns(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1DNS)`

SetDns sets Dns field to given value.

### HasDns

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentSpec) HasDns() bool`

HasDns returns a boolean if a field has been set.

### GetNetwork

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentSpec) GetNetwork() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Network`

GetNetwork returns the Network field if non-nil, zero value otherwise.

### GetNetworkOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentSpec) GetNetworkOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Network, bool)`

GetNetworkOk returns a tuple with the Network field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNetwork

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentSpec) SetNetwork(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Network)`

SetNetwork sets Network field to given value.

### HasNetwork

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentSpec) HasNetwork() bool`

HasNetwork returns a boolean if a field has been set.

### GetRegion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentSpec) GetRegion() string`

GetRegion returns the Region field if non-nil, zero value otherwise.

### GetRegionOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentSpec) GetRegionOk() (*string, bool)`

GetRegionOk returns a tuple with the Region field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRegion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentSpec) SetRegion(v string)`

SetRegion sets Region field to given value.

### HasRegion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentSpec) HasRegion() bool`

HasRegion returns a boolean if a field has been set.

### GetZone

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentSpec) GetZone() string`

GetZone returns the Zone field if non-nil, zero value otherwise.

### GetZoneOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentSpec) GetZoneOk() (*string, bool)`

GetZoneOk returns a tuple with the Zone field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetZone

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentSpec) SetZone(v string)`

SetZone sets Zone field to given value.

### HasZone

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentSpec) HasZone() bool`

HasZone returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


