# ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AdditionalCloudStorage** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1CloudStorage**](ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1CloudStorage.md) |  | [optional] 
**CloudStorage** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1CloudStorage**](ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1CloudStorage.md) |  | [optional] 
**DisableIamRoleCreation** | Pointer to **bool** |  | [optional] 
**PoolMemberRef** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1PoolMemberReference**](ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1PoolMemberReference.md) |  | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountSpec

`func NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountSpec() *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountSpec`

NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountSpec instantiates a new ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountSpec object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountSpecWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountSpecWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountSpec`

NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountSpecWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountSpec object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAdditionalCloudStorage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountSpec) GetAdditionalCloudStorage() []ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1CloudStorage`

GetAdditionalCloudStorage returns the AdditionalCloudStorage field if non-nil, zero value otherwise.

### GetAdditionalCloudStorageOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountSpec) GetAdditionalCloudStorageOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1CloudStorage, bool)`

GetAdditionalCloudStorageOk returns a tuple with the AdditionalCloudStorage field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAdditionalCloudStorage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountSpec) SetAdditionalCloudStorage(v []ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1CloudStorage)`

SetAdditionalCloudStorage sets AdditionalCloudStorage field to given value.

### HasAdditionalCloudStorage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountSpec) HasAdditionalCloudStorage() bool`

HasAdditionalCloudStorage returns a boolean if a field has been set.

### GetCloudStorage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountSpec) GetCloudStorage() ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1CloudStorage`

GetCloudStorage returns the CloudStorage field if non-nil, zero value otherwise.

### GetCloudStorageOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountSpec) GetCloudStorageOk() (*ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1CloudStorage, bool)`

GetCloudStorageOk returns a tuple with the CloudStorage field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCloudStorage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountSpec) SetCloudStorage(v ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1CloudStorage)`

SetCloudStorage sets CloudStorage field to given value.

### HasCloudStorage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountSpec) HasCloudStorage() bool`

HasCloudStorage returns a boolean if a field has been set.

### GetDisableIamRoleCreation

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountSpec) GetDisableIamRoleCreation() bool`

GetDisableIamRoleCreation returns the DisableIamRoleCreation field if non-nil, zero value otherwise.

### GetDisableIamRoleCreationOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountSpec) GetDisableIamRoleCreationOk() (*bool, bool)`

GetDisableIamRoleCreationOk returns a tuple with the DisableIamRoleCreation field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDisableIamRoleCreation

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountSpec) SetDisableIamRoleCreation(v bool)`

SetDisableIamRoleCreation sets DisableIamRoleCreation field to given value.

### HasDisableIamRoleCreation

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountSpec) HasDisableIamRoleCreation() bool`

HasDisableIamRoleCreation returns a boolean if a field has been set.

### GetPoolMemberRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountSpec) GetPoolMemberRef() ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1PoolMemberReference`

GetPoolMemberRef returns the PoolMemberRef field if non-nil, zero value otherwise.

### GetPoolMemberRefOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountSpec) GetPoolMemberRefOk() (*ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1PoolMemberReference, bool)`

GetPoolMemberRefOk returns a tuple with the PoolMemberRef field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPoolMemberRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountSpec) SetPoolMemberRef(v ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1PoolMemberReference)`

SetPoolMemberRef sets PoolMemberRef field to given value.

### HasPoolMemberRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountSpec) HasPoolMemberRef() bool`

HasPoolMemberRef returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


