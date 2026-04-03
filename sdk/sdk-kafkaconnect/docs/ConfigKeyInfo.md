# ConfigKeyInfo

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**DefaultValue** | Pointer to **string** |  | [optional] 
**Dependents** | Pointer to **[]string** |  | [optional] 
**DisplayName** | Pointer to **string** |  | [optional] 
**Documentation** | Pointer to **string** |  | [optional] 
**Group** | Pointer to **string** |  | [optional] 
**Importance** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Order** | Pointer to **int32** |  | [optional] 
**OrderInGroup** | Pointer to **int32** |  | [optional] 
**Required** | Pointer to **bool** |  | [optional] 
**Type** | Pointer to **string** |  | [optional] 
**Width** | Pointer to **string** |  | [optional] 

## Methods

### NewConfigKeyInfo

`func NewConfigKeyInfo() *ConfigKeyInfo`

NewConfigKeyInfo instantiates a new ConfigKeyInfo object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigKeyInfoWithDefaults

`func NewConfigKeyInfoWithDefaults() *ConfigKeyInfo`

NewConfigKeyInfoWithDefaults instantiates a new ConfigKeyInfo object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDefaultValue

`func (o *ConfigKeyInfo) GetDefaultValue() string`

GetDefaultValue returns the DefaultValue field if non-nil, zero value otherwise.

### GetDefaultValueOk

`func (o *ConfigKeyInfo) GetDefaultValueOk() (*string, bool)`

GetDefaultValueOk returns a tuple with the DefaultValue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDefaultValue

`func (o *ConfigKeyInfo) SetDefaultValue(v string)`

SetDefaultValue sets DefaultValue field to given value.

### HasDefaultValue

`func (o *ConfigKeyInfo) HasDefaultValue() bool`

HasDefaultValue returns a boolean if a field has been set.

### GetDependents

`func (o *ConfigKeyInfo) GetDependents() []string`

GetDependents returns the Dependents field if non-nil, zero value otherwise.

### GetDependentsOk

`func (o *ConfigKeyInfo) GetDependentsOk() (*[]string, bool)`

GetDependentsOk returns a tuple with the Dependents field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDependents

`func (o *ConfigKeyInfo) SetDependents(v []string)`

SetDependents sets Dependents field to given value.

### HasDependents

`func (o *ConfigKeyInfo) HasDependents() bool`

HasDependents returns a boolean if a field has been set.

### GetDisplayName

`func (o *ConfigKeyInfo) GetDisplayName() string`

GetDisplayName returns the DisplayName field if non-nil, zero value otherwise.

### GetDisplayNameOk

`func (o *ConfigKeyInfo) GetDisplayNameOk() (*string, bool)`

GetDisplayNameOk returns a tuple with the DisplayName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDisplayName

`func (o *ConfigKeyInfo) SetDisplayName(v string)`

SetDisplayName sets DisplayName field to given value.

### HasDisplayName

`func (o *ConfigKeyInfo) HasDisplayName() bool`

HasDisplayName returns a boolean if a field has been set.

### GetDocumentation

`func (o *ConfigKeyInfo) GetDocumentation() string`

GetDocumentation returns the Documentation field if non-nil, zero value otherwise.

### GetDocumentationOk

`func (o *ConfigKeyInfo) GetDocumentationOk() (*string, bool)`

GetDocumentationOk returns a tuple with the Documentation field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDocumentation

`func (o *ConfigKeyInfo) SetDocumentation(v string)`

SetDocumentation sets Documentation field to given value.

### HasDocumentation

`func (o *ConfigKeyInfo) HasDocumentation() bool`

HasDocumentation returns a boolean if a field has been set.

### GetGroup

`func (o *ConfigKeyInfo) GetGroup() string`

GetGroup returns the Group field if non-nil, zero value otherwise.

### GetGroupOk

`func (o *ConfigKeyInfo) GetGroupOk() (*string, bool)`

GetGroupOk returns a tuple with the Group field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGroup

`func (o *ConfigKeyInfo) SetGroup(v string)`

SetGroup sets Group field to given value.

### HasGroup

`func (o *ConfigKeyInfo) HasGroup() bool`

HasGroup returns a boolean if a field has been set.

### GetImportance

`func (o *ConfigKeyInfo) GetImportance() string`

GetImportance returns the Importance field if non-nil, zero value otherwise.

### GetImportanceOk

`func (o *ConfigKeyInfo) GetImportanceOk() (*string, bool)`

GetImportanceOk returns a tuple with the Importance field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImportance

`func (o *ConfigKeyInfo) SetImportance(v string)`

SetImportance sets Importance field to given value.

### HasImportance

`func (o *ConfigKeyInfo) HasImportance() bool`

HasImportance returns a boolean if a field has been set.

### GetName

`func (o *ConfigKeyInfo) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ConfigKeyInfo) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ConfigKeyInfo) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ConfigKeyInfo) HasName() bool`

HasName returns a boolean if a field has been set.

### GetOrder

`func (o *ConfigKeyInfo) GetOrder() int32`

GetOrder returns the Order field if non-nil, zero value otherwise.

### GetOrderOk

`func (o *ConfigKeyInfo) GetOrderOk() (*int32, bool)`

GetOrderOk returns a tuple with the Order field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrder

`func (o *ConfigKeyInfo) SetOrder(v int32)`

SetOrder sets Order field to given value.

### HasOrder

`func (o *ConfigKeyInfo) HasOrder() bool`

HasOrder returns a boolean if a field has been set.

### GetOrderInGroup

`func (o *ConfigKeyInfo) GetOrderInGroup() int32`

GetOrderInGroup returns the OrderInGroup field if non-nil, zero value otherwise.

### GetOrderInGroupOk

`func (o *ConfigKeyInfo) GetOrderInGroupOk() (*int32, bool)`

GetOrderInGroupOk returns a tuple with the OrderInGroup field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrderInGroup

`func (o *ConfigKeyInfo) SetOrderInGroup(v int32)`

SetOrderInGroup sets OrderInGroup field to given value.

### HasOrderInGroup

`func (o *ConfigKeyInfo) HasOrderInGroup() bool`

HasOrderInGroup returns a boolean if a field has been set.

### GetRequired

`func (o *ConfigKeyInfo) GetRequired() bool`

GetRequired returns the Required field if non-nil, zero value otherwise.

### GetRequiredOk

`func (o *ConfigKeyInfo) GetRequiredOk() (*bool, bool)`

GetRequiredOk returns a tuple with the Required field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRequired

`func (o *ConfigKeyInfo) SetRequired(v bool)`

SetRequired sets Required field to given value.

### HasRequired

`func (o *ConfigKeyInfo) HasRequired() bool`

HasRequired returns a boolean if a field has been set.

### GetType

`func (o *ConfigKeyInfo) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *ConfigKeyInfo) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *ConfigKeyInfo) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *ConfigKeyInfo) HasType() bool`

HasType returns a boolean if a field has been set.

### GetWidth

`func (o *ConfigKeyInfo) GetWidth() string`

GetWidth returns the Width field if non-nil, zero value otherwise.

### GetWidthOk

`func (o *ConfigKeyInfo) GetWidthOk() (*string, bool)`

GetWidthOk returns a tuple with the Width field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWidth

`func (o *ConfigKeyInfo) SetWidth(v string)`

SetWidth sets Width field to given value.

### HasWidth

`func (o *ConfigKeyInfo) HasWidth() bool`

HasWidth returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


