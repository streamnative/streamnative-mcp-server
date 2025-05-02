# V1EnvFromSource

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ConfigMapRef** | Pointer to [**V1ConfigMapEnvSource**](V1ConfigMapEnvSource.md) |  | [optional] 
**Prefix** | Pointer to **string** | An optional identifier to prepend to each key in the ConfigMap. Must be a C_IDENTIFIER. | [optional] 
**SecretRef** | Pointer to [**V1SecretEnvSource**](V1SecretEnvSource.md) |  | [optional] 

## Methods

### NewV1EnvFromSource

`func NewV1EnvFromSource() *V1EnvFromSource`

NewV1EnvFromSource instantiates a new V1EnvFromSource object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewV1EnvFromSourceWithDefaults

`func NewV1EnvFromSourceWithDefaults() *V1EnvFromSource`

NewV1EnvFromSourceWithDefaults instantiates a new V1EnvFromSource object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetConfigMapRef

`func (o *V1EnvFromSource) GetConfigMapRef() V1ConfigMapEnvSource`

GetConfigMapRef returns the ConfigMapRef field if non-nil, zero value otherwise.

### GetConfigMapRefOk

`func (o *V1EnvFromSource) GetConfigMapRefOk() (*V1ConfigMapEnvSource, bool)`

GetConfigMapRefOk returns a tuple with the ConfigMapRef field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConfigMapRef

`func (o *V1EnvFromSource) SetConfigMapRef(v V1ConfigMapEnvSource)`

SetConfigMapRef sets ConfigMapRef field to given value.

### HasConfigMapRef

`func (o *V1EnvFromSource) HasConfigMapRef() bool`

HasConfigMapRef returns a boolean if a field has been set.

### GetPrefix

`func (o *V1EnvFromSource) GetPrefix() string`

GetPrefix returns the Prefix field if non-nil, zero value otherwise.

### GetPrefixOk

`func (o *V1EnvFromSource) GetPrefixOk() (*string, bool)`

GetPrefixOk returns a tuple with the Prefix field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrefix

`func (o *V1EnvFromSource) SetPrefix(v string)`

SetPrefix sets Prefix field to given value.

### HasPrefix

`func (o *V1EnvFromSource) HasPrefix() bool`

HasPrefix returns a boolean if a field has been set.

### GetSecretRef

`func (o *V1EnvFromSource) GetSecretRef() V1SecretEnvSource`

GetSecretRef returns the SecretRef field if non-nil, zero value otherwise.

### GetSecretRefOk

`func (o *V1EnvFromSource) GetSecretRefOk() (*V1SecretEnvSource, bool)`

GetSecretRefOk returns a tuple with the SecretRef field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSecretRef

`func (o *V1EnvFromSource) SetSecretRef(v V1SecretEnvSource)`

SetSecretRef sets SecretRef field to given value.

### HasSecretRef

`func (o *V1EnvFromSource) HasSecretRef() bool`

HasSecretRef returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


