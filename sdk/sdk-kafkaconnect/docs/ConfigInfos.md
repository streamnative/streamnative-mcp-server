# ConfigInfos

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Configs** | Pointer to [**[]ConfigInfo**](ConfigInfo.md) |  | [optional] 
**ErrorCount** | Pointer to **int32** |  | [optional] 
**Groups** | Pointer to **[]string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 

## Methods

### NewConfigInfos

`func NewConfigInfos() *ConfigInfos`

NewConfigInfos instantiates a new ConfigInfos object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigInfosWithDefaults

`func NewConfigInfosWithDefaults() *ConfigInfos`

NewConfigInfosWithDefaults instantiates a new ConfigInfos object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetConfigs

`func (o *ConfigInfos) GetConfigs() []ConfigInfo`

GetConfigs returns the Configs field if non-nil, zero value otherwise.

### GetConfigsOk

`func (o *ConfigInfos) GetConfigsOk() (*[]ConfigInfo, bool)`

GetConfigsOk returns a tuple with the Configs field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConfigs

`func (o *ConfigInfos) SetConfigs(v []ConfigInfo)`

SetConfigs sets Configs field to given value.

### HasConfigs

`func (o *ConfigInfos) HasConfigs() bool`

HasConfigs returns a boolean if a field has been set.

### GetErrorCount

`func (o *ConfigInfos) GetErrorCount() int32`

GetErrorCount returns the ErrorCount field if non-nil, zero value otherwise.

### GetErrorCountOk

`func (o *ConfigInfos) GetErrorCountOk() (*int32, bool)`

GetErrorCountOk returns a tuple with the ErrorCount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetErrorCount

`func (o *ConfigInfos) SetErrorCount(v int32)`

SetErrorCount sets ErrorCount field to given value.

### HasErrorCount

`func (o *ConfigInfos) HasErrorCount() bool`

HasErrorCount returns a boolean if a field has been set.

### GetGroups

`func (o *ConfigInfos) GetGroups() []string`

GetGroups returns the Groups field if non-nil, zero value otherwise.

### GetGroupsOk

`func (o *ConfigInfos) GetGroupsOk() (*[]string, bool)`

GetGroupsOk returns a tuple with the Groups field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGroups

`func (o *ConfigInfos) SetGroups(v []string)`

SetGroups sets Groups field to given value.

### HasGroups

`func (o *ConfigInfos) HasGroups() bool`

HasGroups returns a boolean if a field has been set.

### GetName

`func (o *ConfigInfos) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ConfigInfos) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ConfigInfos) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ConfigInfos) HasName() bool`

HasName returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


