# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Config

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AuditLog** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AuditLog**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AuditLog.md) |  | [optional] 
**Custom** | Pointer to **map[string]string** | Custom accepts custom configurations. | [optional] 
**FunctionEnabled** | Pointer to **bool** | FunctionEnabled controls whether function is enabled. | [optional] 
**LakehouseStorage** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1LakehouseStorageConfig**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1LakehouseStorageConfig.md) |  | [optional] 
**Protocols** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ProtocolsConfig**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ProtocolsConfig.md) |  | [optional] 
**TransactionEnabled** | Pointer to **bool** | TransactionEnabled controls whether transaction is enabled. | [optional] 
**WebsocketEnabled** | Pointer to **bool** | WebsocketEnabled controls whether websocket is enabled. | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Config

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Config() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Config`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Config instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Config object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ConfigWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ConfigWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Config`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ConfigWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Config object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAuditLog

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Config) GetAuditLog() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AuditLog`

GetAuditLog returns the AuditLog field if non-nil, zero value otherwise.

### GetAuditLogOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Config) GetAuditLogOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AuditLog, bool)`

GetAuditLogOk returns a tuple with the AuditLog field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuditLog

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Config) SetAuditLog(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AuditLog)`

SetAuditLog sets AuditLog field to given value.

### HasAuditLog

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Config) HasAuditLog() bool`

HasAuditLog returns a boolean if a field has been set.

### GetCustom

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Config) GetCustom() map[string]string`

GetCustom returns the Custom field if non-nil, zero value otherwise.

### GetCustomOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Config) GetCustomOk() (*map[string]string, bool)`

GetCustomOk returns a tuple with the Custom field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCustom

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Config) SetCustom(v map[string]string)`

SetCustom sets Custom field to given value.

### HasCustom

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Config) HasCustom() bool`

HasCustom returns a boolean if a field has been set.

### GetFunctionEnabled

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Config) GetFunctionEnabled() bool`

GetFunctionEnabled returns the FunctionEnabled field if non-nil, zero value otherwise.

### GetFunctionEnabledOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Config) GetFunctionEnabledOk() (*bool, bool)`

GetFunctionEnabledOk returns a tuple with the FunctionEnabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFunctionEnabled

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Config) SetFunctionEnabled(v bool)`

SetFunctionEnabled sets FunctionEnabled field to given value.

### HasFunctionEnabled

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Config) HasFunctionEnabled() bool`

HasFunctionEnabled returns a boolean if a field has been set.

### GetLakehouseStorage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Config) GetLakehouseStorage() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1LakehouseStorageConfig`

GetLakehouseStorage returns the LakehouseStorage field if non-nil, zero value otherwise.

### GetLakehouseStorageOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Config) GetLakehouseStorageOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1LakehouseStorageConfig, bool)`

GetLakehouseStorageOk returns a tuple with the LakehouseStorage field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLakehouseStorage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Config) SetLakehouseStorage(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1LakehouseStorageConfig)`

SetLakehouseStorage sets LakehouseStorage field to given value.

### HasLakehouseStorage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Config) HasLakehouseStorage() bool`

HasLakehouseStorage returns a boolean if a field has been set.

### GetProtocols

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Config) GetProtocols() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ProtocolsConfig`

GetProtocols returns the Protocols field if non-nil, zero value otherwise.

### GetProtocolsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Config) GetProtocolsOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ProtocolsConfig, bool)`

GetProtocolsOk returns a tuple with the Protocols field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProtocols

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Config) SetProtocols(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ProtocolsConfig)`

SetProtocols sets Protocols field to given value.

### HasProtocols

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Config) HasProtocols() bool`

HasProtocols returns a boolean if a field has been set.

### GetTransactionEnabled

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Config) GetTransactionEnabled() bool`

GetTransactionEnabled returns the TransactionEnabled field if non-nil, zero value otherwise.

### GetTransactionEnabledOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Config) GetTransactionEnabledOk() (*bool, bool)`

GetTransactionEnabledOk returns a tuple with the TransactionEnabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTransactionEnabled

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Config) SetTransactionEnabled(v bool)`

SetTransactionEnabled sets TransactionEnabled field to given value.

### HasTransactionEnabled

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Config) HasTransactionEnabled() bool`

HasTransactionEnabled returns a boolean if a field has been set.

### GetWebsocketEnabled

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Config) GetWebsocketEnabled() bool`

GetWebsocketEnabled returns the WebsocketEnabled field if non-nil, zero value otherwise.

### GetWebsocketEnabledOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Config) GetWebsocketEnabledOk() (*bool, bool)`

GetWebsocketEnabledOk returns a tuple with the WebsocketEnabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWebsocketEnabled

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Config) SetWebsocketEnabled(v bool)`

SetWebsocketEnabled sets WebsocketEnabled field to given value.

### HasWebsocketEnabled

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Config) HasWebsocketEnabled() bool`

HasWebsocketEnabled returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


