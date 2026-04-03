# ServerInfo

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Commit** | Pointer to **string** |  | [optional] 
**KafkaClusterId** | Pointer to **string** |  | [optional] 
**Version** | Pointer to **string** |  | [optional] 

## Methods

### NewServerInfo

`func NewServerInfo() *ServerInfo`

NewServerInfo instantiates a new ServerInfo object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewServerInfoWithDefaults

`func NewServerInfoWithDefaults() *ServerInfo`

NewServerInfoWithDefaults instantiates a new ServerInfo object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCommit

`func (o *ServerInfo) GetCommit() string`

GetCommit returns the Commit field if non-nil, zero value otherwise.

### GetCommitOk

`func (o *ServerInfo) GetCommitOk() (*string, bool)`

GetCommitOk returns a tuple with the Commit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCommit

`func (o *ServerInfo) SetCommit(v string)`

SetCommit sets Commit field to given value.

### HasCommit

`func (o *ServerInfo) HasCommit() bool`

HasCommit returns a boolean if a field has been set.

### GetKafkaClusterId

`func (o *ServerInfo) GetKafkaClusterId() string`

GetKafkaClusterId returns the KafkaClusterId field if non-nil, zero value otherwise.

### GetKafkaClusterIdOk

`func (o *ServerInfo) GetKafkaClusterIdOk() (*string, bool)`

GetKafkaClusterIdOk returns a tuple with the KafkaClusterId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKafkaClusterId

`func (o *ServerInfo) SetKafkaClusterId(v string)`

SetKafkaClusterId sets KafkaClusterId field to given value.

### HasKafkaClusterId

`func (o *ServerInfo) HasKafkaClusterId() bool`

HasKafkaClusterId returns a boolean if a field has been set.

### GetVersion

`func (o *ServerInfo) GetVersion() string`

GetVersion returns the Version field if non-nil, zero value otherwise.

### GetVersionOk

`func (o *ServerInfo) GetVersionOk() (*string, bool)`

GetVersionOk returns a tuple with the Version field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersion

`func (o *ServerInfo) SetVersion(v string)`

SetVersion sets Version field to given value.

### HasVersion

`func (o *ServerInfo) HasVersion() bool`

HasVersion returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


