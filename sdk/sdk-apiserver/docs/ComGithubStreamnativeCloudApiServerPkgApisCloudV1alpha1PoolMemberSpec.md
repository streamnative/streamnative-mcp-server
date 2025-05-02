# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Aws** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AwsPoolMemberSpec**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AwsPoolMemberSpec.md) |  | [optional] 
**Azure** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec.md) |  | [optional] 
**ConnectionOptions** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberConnectionOptionsSpec**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberConnectionOptionsSpec.md) |  | [optional] 
**Domains** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Domain**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Domain.md) |  | [optional] 
**FunctionMesh** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberPulsarSpec**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberPulsarSpec.md) |  | [optional] 
**Gcloud** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec.md) |  | [optional] 
**Generic** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GenericPoolMemberSpec**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GenericPoolMemberSpec.md) |  | [optional] 
**Istio** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberIstioSpec**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberIstioSpec.md) |  | [optional] 
**Monitoring** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberMonitoring**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberMonitoring.md) |  | [optional] 
**PoolName** | **string** |  | 
**Pulsar** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberPulsarSpec**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberPulsarSpec.md) |  | [optional] 
**SupportAccess** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SupportAccessOptionsSpec**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SupportAccessOptionsSpec.md) |  | [optional] 
**Taints** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Taint**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Taint.md) |  | [optional] 
**TieredStorage** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberTieredStorageSpec**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberTieredStorageSpec.md) |  | [optional] 
**Type** | **string** |  | 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec(poolName string, type_ string, ) *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpecWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpecWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpecWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAws

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) GetAws() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AwsPoolMemberSpec`

GetAws returns the Aws field if non-nil, zero value otherwise.

### GetAwsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) GetAwsOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AwsPoolMemberSpec, bool)`

GetAwsOk returns a tuple with the Aws field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAws

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) SetAws(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AwsPoolMemberSpec)`

SetAws sets Aws field to given value.

### HasAws

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) HasAws() bool`

HasAws returns a boolean if a field has been set.

### GetAzure

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) GetAzure() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec`

GetAzure returns the Azure field if non-nil, zero value otherwise.

### GetAzureOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) GetAzureOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec, bool)`

GetAzureOk returns a tuple with the Azure field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAzure

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) SetAzure(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AzurePoolMemberSpec)`

SetAzure sets Azure field to given value.

### HasAzure

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) HasAzure() bool`

HasAzure returns a boolean if a field has been set.

### GetConnectionOptions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) GetConnectionOptions() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberConnectionOptionsSpec`

GetConnectionOptions returns the ConnectionOptions field if non-nil, zero value otherwise.

### GetConnectionOptionsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) GetConnectionOptionsOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberConnectionOptionsSpec, bool)`

GetConnectionOptionsOk returns a tuple with the ConnectionOptions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConnectionOptions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) SetConnectionOptions(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberConnectionOptionsSpec)`

SetConnectionOptions sets ConnectionOptions field to given value.

### HasConnectionOptions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) HasConnectionOptions() bool`

HasConnectionOptions returns a boolean if a field has been set.

### GetDomains

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) GetDomains() []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Domain`

GetDomains returns the Domains field if non-nil, zero value otherwise.

### GetDomainsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) GetDomainsOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Domain, bool)`

GetDomainsOk returns a tuple with the Domains field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDomains

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) SetDomains(v []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Domain)`

SetDomains sets Domains field to given value.

### HasDomains

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) HasDomains() bool`

HasDomains returns a boolean if a field has been set.

### GetFunctionMesh

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) GetFunctionMesh() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberPulsarSpec`

GetFunctionMesh returns the FunctionMesh field if non-nil, zero value otherwise.

### GetFunctionMeshOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) GetFunctionMeshOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberPulsarSpec, bool)`

GetFunctionMeshOk returns a tuple with the FunctionMesh field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFunctionMesh

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) SetFunctionMesh(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberPulsarSpec)`

SetFunctionMesh sets FunctionMesh field to given value.

### HasFunctionMesh

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) HasFunctionMesh() bool`

HasFunctionMesh returns a boolean if a field has been set.

### GetGcloud

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) GetGcloud() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec`

GetGcloud returns the Gcloud field if non-nil, zero value otherwise.

### GetGcloudOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) GetGcloudOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec, bool)`

GetGcloudOk returns a tuple with the Gcloud field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGcloud

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) SetGcloud(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec)`

SetGcloud sets Gcloud field to given value.

### HasGcloud

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) HasGcloud() bool`

HasGcloud returns a boolean if a field has been set.

### GetGeneric

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) GetGeneric() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GenericPoolMemberSpec`

GetGeneric returns the Generic field if non-nil, zero value otherwise.

### GetGenericOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) GetGenericOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GenericPoolMemberSpec, bool)`

GetGenericOk returns a tuple with the Generic field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGeneric

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) SetGeneric(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GenericPoolMemberSpec)`

SetGeneric sets Generic field to given value.

### HasGeneric

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) HasGeneric() bool`

HasGeneric returns a boolean if a field has been set.

### GetIstio

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) GetIstio() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberIstioSpec`

GetIstio returns the Istio field if non-nil, zero value otherwise.

### GetIstioOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) GetIstioOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberIstioSpec, bool)`

GetIstioOk returns a tuple with the Istio field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIstio

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) SetIstio(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberIstioSpec)`

SetIstio sets Istio field to given value.

### HasIstio

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) HasIstio() bool`

HasIstio returns a boolean if a field has been set.

### GetMonitoring

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) GetMonitoring() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberMonitoring`

GetMonitoring returns the Monitoring field if non-nil, zero value otherwise.

### GetMonitoringOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) GetMonitoringOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberMonitoring, bool)`

GetMonitoringOk returns a tuple with the Monitoring field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMonitoring

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) SetMonitoring(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberMonitoring)`

SetMonitoring sets Monitoring field to given value.

### HasMonitoring

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) HasMonitoring() bool`

HasMonitoring returns a boolean if a field has been set.

### GetPoolName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) GetPoolName() string`

GetPoolName returns the PoolName field if non-nil, zero value otherwise.

### GetPoolNameOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) GetPoolNameOk() (*string, bool)`

GetPoolNameOk returns a tuple with the PoolName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPoolName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) SetPoolName(v string)`

SetPoolName sets PoolName field to given value.


### GetPulsar

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) GetPulsar() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberPulsarSpec`

GetPulsar returns the Pulsar field if non-nil, zero value otherwise.

### GetPulsarOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) GetPulsarOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberPulsarSpec, bool)`

GetPulsarOk returns a tuple with the Pulsar field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPulsar

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) SetPulsar(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberPulsarSpec)`

SetPulsar sets Pulsar field to given value.

### HasPulsar

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) HasPulsar() bool`

HasPulsar returns a boolean if a field has been set.

### GetSupportAccess

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) GetSupportAccess() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SupportAccessOptionsSpec`

GetSupportAccess returns the SupportAccess field if non-nil, zero value otherwise.

### GetSupportAccessOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) GetSupportAccessOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SupportAccessOptionsSpec, bool)`

GetSupportAccessOk returns a tuple with the SupportAccess field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSupportAccess

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) SetSupportAccess(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SupportAccessOptionsSpec)`

SetSupportAccess sets SupportAccess field to given value.

### HasSupportAccess

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) HasSupportAccess() bool`

HasSupportAccess returns a boolean if a field has been set.

### GetTaints

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) GetTaints() []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Taint`

GetTaints returns the Taints field if non-nil, zero value otherwise.

### GetTaintsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) GetTaintsOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Taint, bool)`

GetTaintsOk returns a tuple with the Taints field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTaints

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) SetTaints(v []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Taint)`

SetTaints sets Taints field to given value.

### HasTaints

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) HasTaints() bool`

HasTaints returns a boolean if a field has been set.

### GetTieredStorage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) GetTieredStorage() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberTieredStorageSpec`

GetTieredStorage returns the TieredStorage field if non-nil, zero value otherwise.

### GetTieredStorageOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) GetTieredStorageOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberTieredStorageSpec, bool)`

GetTieredStorageOk returns a tuple with the TieredStorage field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTieredStorage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) SetTieredStorage(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberTieredStorageSpec)`

SetTieredStorage sets TieredStorage field to given value.

### HasTieredStorage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) HasTieredStorage() bool`

HasTieredStorage returns a boolean if a field has been set.

### GetType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberSpec) SetType(v string)`

SetType sets Type field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


