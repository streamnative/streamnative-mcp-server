# ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1ResourceRule

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ApiGroups** | Pointer to **[]string** | APIGroups is the name of the APIGroup that contains the resources.  If multiple API groups are specified, any action requested against one of the enumerated resources in any API group will be allowed.  \&quot;*\&quot; means all. | [optional] 
**ResourceNames** | Pointer to **[]string** | ResourceNames is an optional white list of names that the rule applies to.  An empty set means that everything is allowed.  \&quot;*\&quot; means all. | [optional] 
**Resources** | Pointer to **[]string** | Resources is a list of resources this rule applies to.  \&quot;*\&quot; means all in the specified apiGroups.  \&quot;*_/foo\&quot; represents the subresource &#39;foo&#39; for all resources in the specified apiGroups. | [optional] 
**Verbs** | **[]string** | Verb is a list of kubernetes resource API verbs, like: get, list, watch, create, update, delete, proxy.  \&quot;*\&quot; means all. | 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1ResourceRule

`func NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1ResourceRule(verbs []string, ) *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1ResourceRule`

NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1ResourceRule instantiates a new ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1ResourceRule object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1ResourceRuleWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1ResourceRuleWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1ResourceRule`

NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1ResourceRuleWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1ResourceRule object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetApiGroups

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1ResourceRule) GetApiGroups() []string`

GetApiGroups returns the ApiGroups field if non-nil, zero value otherwise.

### GetApiGroupsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1ResourceRule) GetApiGroupsOk() (*[]string, bool)`

GetApiGroupsOk returns a tuple with the ApiGroups field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetApiGroups

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1ResourceRule) SetApiGroups(v []string)`

SetApiGroups sets ApiGroups field to given value.

### HasApiGroups

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1ResourceRule) HasApiGroups() bool`

HasApiGroups returns a boolean if a field has been set.

### GetResourceNames

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1ResourceRule) GetResourceNames() []string`

GetResourceNames returns the ResourceNames field if non-nil, zero value otherwise.

### GetResourceNamesOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1ResourceRule) GetResourceNamesOk() (*[]string, bool)`

GetResourceNamesOk returns a tuple with the ResourceNames field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResourceNames

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1ResourceRule) SetResourceNames(v []string)`

SetResourceNames sets ResourceNames field to given value.

### HasResourceNames

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1ResourceRule) HasResourceNames() bool`

HasResourceNames returns a boolean if a field has been set.

### GetResources

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1ResourceRule) GetResources() []string`

GetResources returns the Resources field if non-nil, zero value otherwise.

### GetResourcesOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1ResourceRule) GetResourcesOk() (*[]string, bool)`

GetResourcesOk returns a tuple with the Resources field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResources

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1ResourceRule) SetResources(v []string)`

SetResources sets Resources field to given value.

### HasResources

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1ResourceRule) HasResources() bool`

HasResources returns a boolean if a field has been set.

### GetVerbs

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1ResourceRule) GetVerbs() []string`

GetVerbs returns the Verbs field if non-nil, zero value otherwise.

### GetVerbsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1ResourceRule) GetVerbsOk() (*[]string, bool)`

GetVerbsOk returns a tuple with the Verbs field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVerbs

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1ResourceRule) SetVerbs(v []string)`

SetVerbs sets Verbs field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


