# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeySpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Description** | Pointer to **string** | Description is a user defined description of the key | [optional] 
**EncryptionKey** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1EncryptionKey**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1EncryptionKey.md) |  | [optional] 
**ExpirationTime** | Pointer to **time.Time** | Expiration is a duration (as a golang duration string) that defines how long this API key is valid for. This can only be set on initial creation and not updated later | [optional] 
**InstanceName** | Pointer to **string** | InstanceName is the name of the instance this API key is for | [optional] 
**Revoke** | Pointer to **bool** | Revoke is a boolean that defines if the token of this API key should be revoked. | [optional] 
**ServiceAccountName** | Pointer to **string** | ServiceAccountName is the name of the service account this API key is for | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeySpec

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeySpec() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeySpec`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeySpec instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeySpec object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeySpecWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeySpecWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeySpec`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeySpecWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeySpec object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDescription

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeySpec) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeySpec) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeySpec) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeySpec) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetEncryptionKey

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeySpec) GetEncryptionKey() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1EncryptionKey`

GetEncryptionKey returns the EncryptionKey field if non-nil, zero value otherwise.

### GetEncryptionKeyOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeySpec) GetEncryptionKeyOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1EncryptionKey, bool)`

GetEncryptionKeyOk returns a tuple with the EncryptionKey field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEncryptionKey

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeySpec) SetEncryptionKey(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1EncryptionKey)`

SetEncryptionKey sets EncryptionKey field to given value.

### HasEncryptionKey

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeySpec) HasEncryptionKey() bool`

HasEncryptionKey returns a boolean if a field has been set.

### GetExpirationTime

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeySpec) GetExpirationTime() time.Time`

GetExpirationTime returns the ExpirationTime field if non-nil, zero value otherwise.

### GetExpirationTimeOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeySpec) GetExpirationTimeOk() (*time.Time, bool)`

GetExpirationTimeOk returns a tuple with the ExpirationTime field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpirationTime

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeySpec) SetExpirationTime(v time.Time)`

SetExpirationTime sets ExpirationTime field to given value.

### HasExpirationTime

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeySpec) HasExpirationTime() bool`

HasExpirationTime returns a boolean if a field has been set.

### GetInstanceName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeySpec) GetInstanceName() string`

GetInstanceName returns the InstanceName field if non-nil, zero value otherwise.

### GetInstanceNameOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeySpec) GetInstanceNameOk() (*string, bool)`

GetInstanceNameOk returns a tuple with the InstanceName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInstanceName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeySpec) SetInstanceName(v string)`

SetInstanceName sets InstanceName field to given value.

### HasInstanceName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeySpec) HasInstanceName() bool`

HasInstanceName returns a boolean if a field has been set.

### GetRevoke

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeySpec) GetRevoke() bool`

GetRevoke returns the Revoke field if non-nil, zero value otherwise.

### GetRevokeOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeySpec) GetRevokeOk() (*bool, bool)`

GetRevokeOk returns a tuple with the Revoke field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRevoke

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeySpec) SetRevoke(v bool)`

SetRevoke sets Revoke field to given value.

### HasRevoke

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeySpec) HasRevoke() bool`

HasRevoke returns a boolean if a field has been set.

### GetServiceAccountName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeySpec) GetServiceAccountName() string`

GetServiceAccountName returns the ServiceAccountName field if non-nil, zero value otherwise.

### GetServiceAccountNameOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeySpec) GetServiceAccountNameOk() (*string, bool)`

GetServiceAccountNameOk returns a tuple with the ServiceAccountName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServiceAccountName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeySpec) SetServiceAccountName(v string)`

SetServiceAccountName sets ServiceAccountName field to given value.

### HasServiceAccountName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeySpec) HasServiceAccountName() bool`

HasServiceAccountName returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


