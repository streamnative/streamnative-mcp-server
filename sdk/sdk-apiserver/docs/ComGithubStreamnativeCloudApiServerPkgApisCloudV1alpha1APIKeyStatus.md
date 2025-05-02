# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Conditions** | Pointer to [**[]V1Condition**](V1Condition.md) | Conditions is an array of current observed service account conditions. | [optional] 
**EncryptedToken** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1EncryptedToken**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1EncryptedToken.md) |  | [optional] 
**ExpiresAt** | Pointer to **time.Time** | ExpiresAt is a timestamp of when the key expires | [optional] 
**IssuedAt** | Pointer to **time.Time** | IssuedAt is a timestamp of when the key was issued, stored as an epoch in seconds | [optional] 
**KeyId** | Pointer to **string** | KeyId is a generated field that is a uid for the token | [optional] 
**RevokedAt** | Pointer to **time.Time** | ExpiresAt is a timestamp of when the key was revoked, it triggers revocation action | [optional] 
**Token** | Pointer to **string** | Token is the plaintext security token issued for the key. | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyStatus

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyStatus() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyStatus`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyStatus instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyStatus object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyStatusWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyStatusWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyStatus`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyStatusWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyStatus object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyStatus) GetConditions() []V1Condition`

GetConditions returns the Conditions field if non-nil, zero value otherwise.

### GetConditionsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyStatus) GetConditionsOk() (*[]V1Condition, bool)`

GetConditionsOk returns a tuple with the Conditions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyStatus) SetConditions(v []V1Condition)`

SetConditions sets Conditions field to given value.

### HasConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyStatus) HasConditions() bool`

HasConditions returns a boolean if a field has been set.

### GetEncryptedToken

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyStatus) GetEncryptedToken() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1EncryptedToken`

GetEncryptedToken returns the EncryptedToken field if non-nil, zero value otherwise.

### GetEncryptedTokenOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyStatus) GetEncryptedTokenOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1EncryptedToken, bool)`

GetEncryptedTokenOk returns a tuple with the EncryptedToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEncryptedToken

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyStatus) SetEncryptedToken(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1EncryptedToken)`

SetEncryptedToken sets EncryptedToken field to given value.

### HasEncryptedToken

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyStatus) HasEncryptedToken() bool`

HasEncryptedToken returns a boolean if a field has been set.

### GetExpiresAt

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyStatus) GetExpiresAt() time.Time`

GetExpiresAt returns the ExpiresAt field if non-nil, zero value otherwise.

### GetExpiresAtOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyStatus) GetExpiresAtOk() (*time.Time, bool)`

GetExpiresAtOk returns a tuple with the ExpiresAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpiresAt

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyStatus) SetExpiresAt(v time.Time)`

SetExpiresAt sets ExpiresAt field to given value.

### HasExpiresAt

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyStatus) HasExpiresAt() bool`

HasExpiresAt returns a boolean if a field has been set.

### GetIssuedAt

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyStatus) GetIssuedAt() time.Time`

GetIssuedAt returns the IssuedAt field if non-nil, zero value otherwise.

### GetIssuedAtOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyStatus) GetIssuedAtOk() (*time.Time, bool)`

GetIssuedAtOk returns a tuple with the IssuedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIssuedAt

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyStatus) SetIssuedAt(v time.Time)`

SetIssuedAt sets IssuedAt field to given value.

### HasIssuedAt

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyStatus) HasIssuedAt() bool`

HasIssuedAt returns a boolean if a field has been set.

### GetKeyId

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyStatus) GetKeyId() string`

GetKeyId returns the KeyId field if non-nil, zero value otherwise.

### GetKeyIdOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyStatus) GetKeyIdOk() (*string, bool)`

GetKeyIdOk returns a tuple with the KeyId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKeyId

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyStatus) SetKeyId(v string)`

SetKeyId sets KeyId field to given value.

### HasKeyId

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyStatus) HasKeyId() bool`

HasKeyId returns a boolean if a field has been set.

### GetRevokedAt

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyStatus) GetRevokedAt() time.Time`

GetRevokedAt returns the RevokedAt field if non-nil, zero value otherwise.

### GetRevokedAtOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyStatus) GetRevokedAtOk() (*time.Time, bool)`

GetRevokedAtOk returns a tuple with the RevokedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRevokedAt

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyStatus) SetRevokedAt(v time.Time)`

SetRevokedAt sets RevokedAt field to given value.

### HasRevokedAt

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyStatus) HasRevokedAt() bool`

HasRevokedAt returns a boolean if a field has been set.

### GetToken

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyStatus) GetToken() string`

GetToken returns the Token field if non-nil, zero value otherwise.

### GetTokenOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyStatus) GetTokenOk() (*string, bool)`

GetTokenOk returns a tuple with the Token field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetToken

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyStatus) SetToken(v string)`

SetToken sets Token field to given value.

### HasToken

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1APIKeyStatus) HasToken() bool`

HasToken returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


