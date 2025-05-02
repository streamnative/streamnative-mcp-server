# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleBindingSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Cel** | Pointer to **string** |  | [optional] 
**ConditionGroup** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ConditionGroup**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ConditionGroup.md) |  | [optional] 
**RoleRef** | [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleRef**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleRef.md) |  | 
**Subjects** | [**[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Subject**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Subject.md) |  | 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleBindingSpec

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleBindingSpec(roleRef ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleRef, subjects []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Subject, ) *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleBindingSpec`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleBindingSpec instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleBindingSpec object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleBindingSpecWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleBindingSpecWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleBindingSpec`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleBindingSpecWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleBindingSpec object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCel

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleBindingSpec) GetCel() string`

GetCel returns the Cel field if non-nil, zero value otherwise.

### GetCelOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleBindingSpec) GetCelOk() (*string, bool)`

GetCelOk returns a tuple with the Cel field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCel

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleBindingSpec) SetCel(v string)`

SetCel sets Cel field to given value.

### HasCel

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleBindingSpec) HasCel() bool`

HasCel returns a boolean if a field has been set.

### GetConditionGroup

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleBindingSpec) GetConditionGroup() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ConditionGroup`

GetConditionGroup returns the ConditionGroup field if non-nil, zero value otherwise.

### GetConditionGroupOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleBindingSpec) GetConditionGroupOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ConditionGroup, bool)`

GetConditionGroupOk returns a tuple with the ConditionGroup field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConditionGroup

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleBindingSpec) SetConditionGroup(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ConditionGroup)`

SetConditionGroup sets ConditionGroup field to given value.

### HasConditionGroup

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleBindingSpec) HasConditionGroup() bool`

HasConditionGroup returns a boolean if a field has been set.

### GetRoleRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleBindingSpec) GetRoleRef() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleRef`

GetRoleRef returns the RoleRef field if non-nil, zero value otherwise.

### GetRoleRefOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleBindingSpec) GetRoleRefOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleRef, bool)`

GetRoleRefOk returns a tuple with the RoleRef field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRoleRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleBindingSpec) SetRoleRef(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleRef)`

SetRoleRef sets RoleRef field to given value.


### GetSubjects

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleBindingSpec) GetSubjects() []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Subject`

GetSubjects returns the Subjects field if non-nil, zero value otherwise.

### GetSubjectsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleBindingSpec) GetSubjectsOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Subject, bool)`

GetSubjectsOk returns a tuple with the Subjects field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSubjects

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleBindingSpec) SetSubjects(v []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Subject)`

SetSubjects sets Subjects field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


