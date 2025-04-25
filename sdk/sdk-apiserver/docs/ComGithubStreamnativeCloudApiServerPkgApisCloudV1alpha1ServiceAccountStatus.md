# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ServiceAccountStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Conditions** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition.md) | Conditions is an array of current observed service account conditions. | [optional] 
**PrivateKeyData** | Pointer to **string** | PrivateKeyData provides the private key data (in base-64 format) for authentication purposes | [optional] 
**PrivateKeyType** | Pointer to **string** | PrivateKeyType indicates the type of private key information | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ServiceAccountStatus

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ServiceAccountStatus() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ServiceAccountStatus`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ServiceAccountStatus instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ServiceAccountStatus object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ServiceAccountStatusWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ServiceAccountStatusWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ServiceAccountStatus`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ServiceAccountStatusWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ServiceAccountStatus object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ServiceAccountStatus) GetConditions() []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition`

GetConditions returns the Conditions field if non-nil, zero value otherwise.

### GetConditionsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ServiceAccountStatus) GetConditionsOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition, bool)`

GetConditionsOk returns a tuple with the Conditions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ServiceAccountStatus) SetConditions(v []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition)`

SetConditions sets Conditions field to given value.

### HasConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ServiceAccountStatus) HasConditions() bool`

HasConditions returns a boolean if a field has been set.

### GetPrivateKeyData

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ServiceAccountStatus) GetPrivateKeyData() string`

GetPrivateKeyData returns the PrivateKeyData field if non-nil, zero value otherwise.

### GetPrivateKeyDataOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ServiceAccountStatus) GetPrivateKeyDataOk() (*string, bool)`

GetPrivateKeyDataOk returns a tuple with the PrivateKeyData field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrivateKeyData

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ServiceAccountStatus) SetPrivateKeyData(v string)`

SetPrivateKeyData sets PrivateKeyData field to given value.

### HasPrivateKeyData

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ServiceAccountStatus) HasPrivateKeyData() bool`

HasPrivateKeyData returns a boolean if a field has been set.

### GetPrivateKeyType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ServiceAccountStatus) GetPrivateKeyType() string`

GetPrivateKeyType returns the PrivateKeyType field if non-nil, zero value otherwise.

### GetPrivateKeyTypeOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ServiceAccountStatus) GetPrivateKeyTypeOk() (*string, bool)`

GetPrivateKeyTypeOk returns a tuple with the PrivateKeyType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrivateKeyType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ServiceAccountStatus) SetPrivateKeyType(v string)`

SetPrivateKeyType sets PrivateKeyType field to given value.

### HasPrivateKeyType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ServiceAccountStatus) HasPrivateKeyType() bool`

HasPrivateKeyType returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


