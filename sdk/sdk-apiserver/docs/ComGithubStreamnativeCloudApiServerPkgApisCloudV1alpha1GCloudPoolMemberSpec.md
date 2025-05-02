# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AdminServiceAccount** | Pointer to **string** | AdminServiceAccount is the service account to be impersonated, or leave it empty to call the API without impersonation. | [optional] 
**ClusterName** | Pointer to **string** | ClusterName is the GKE cluster name. | [optional] 
**InitialNodeCount** | Pointer to **int32** | Deprecated | [optional] 
**Location** | **string** |  | 
**MachineType** | Pointer to **string** | Deprecated | [optional] 
**MaxNodeCount** | Pointer to **int32** | Deprecated | [optional] 
**Project** | Pointer to **string** | Project is the Google project containing the GKE cluster. | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec(location string, ) *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpecWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpecWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpecWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAdminServiceAccount

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) GetAdminServiceAccount() string`

GetAdminServiceAccount returns the AdminServiceAccount field if non-nil, zero value otherwise.

### GetAdminServiceAccountOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) GetAdminServiceAccountOk() (*string, bool)`

GetAdminServiceAccountOk returns a tuple with the AdminServiceAccount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAdminServiceAccount

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) SetAdminServiceAccount(v string)`

SetAdminServiceAccount sets AdminServiceAccount field to given value.

### HasAdminServiceAccount

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) HasAdminServiceAccount() bool`

HasAdminServiceAccount returns a boolean if a field has been set.

### GetClusterName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) GetClusterName() string`

GetClusterName returns the ClusterName field if non-nil, zero value otherwise.

### GetClusterNameOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) GetClusterNameOk() (*string, bool)`

GetClusterNameOk returns a tuple with the ClusterName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClusterName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) SetClusterName(v string)`

SetClusterName sets ClusterName field to given value.

### HasClusterName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) HasClusterName() bool`

HasClusterName returns a boolean if a field has been set.

### GetInitialNodeCount

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) GetInitialNodeCount() int32`

GetInitialNodeCount returns the InitialNodeCount field if non-nil, zero value otherwise.

### GetInitialNodeCountOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) GetInitialNodeCountOk() (*int32, bool)`

GetInitialNodeCountOk returns a tuple with the InitialNodeCount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInitialNodeCount

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) SetInitialNodeCount(v int32)`

SetInitialNodeCount sets InitialNodeCount field to given value.

### HasInitialNodeCount

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) HasInitialNodeCount() bool`

HasInitialNodeCount returns a boolean if a field has been set.

### GetLocation

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) GetLocation() string`

GetLocation returns the Location field if non-nil, zero value otherwise.

### GetLocationOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) GetLocationOk() (*string, bool)`

GetLocationOk returns a tuple with the Location field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLocation

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) SetLocation(v string)`

SetLocation sets Location field to given value.


### GetMachineType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) GetMachineType() string`

GetMachineType returns the MachineType field if non-nil, zero value otherwise.

### GetMachineTypeOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) GetMachineTypeOk() (*string, bool)`

GetMachineTypeOk returns a tuple with the MachineType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMachineType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) SetMachineType(v string)`

SetMachineType sets MachineType field to given value.

### HasMachineType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) HasMachineType() bool`

HasMachineType returns a boolean if a field has been set.

### GetMaxNodeCount

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) GetMaxNodeCount() int32`

GetMaxNodeCount returns the MaxNodeCount field if non-nil, zero value otherwise.

### GetMaxNodeCountOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) GetMaxNodeCountOk() (*int32, bool)`

GetMaxNodeCountOk returns a tuple with the MaxNodeCount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaxNodeCount

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) SetMaxNodeCount(v int32)`

SetMaxNodeCount sets MaxNodeCount field to given value.

### HasMaxNodeCount

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) HasMaxNodeCount() bool`

HasMaxNodeCount returns a boolean if a field has been set.

### GetProject

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) GetProject() string`

GetProject returns the Project field if non-nil, zero value otherwise.

### GetProjectOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) GetProjectOk() (*string, bool)`

GetProjectOk returns a tuple with the Project field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProject

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) SetProject(v string)`

SetProject sets Project field to given value.

### HasProject

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GCloudPoolMemberSpec) HasProject() bool`

HasProject returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


