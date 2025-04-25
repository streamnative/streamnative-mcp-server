# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Conditions** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition.md) | Conditions is an array of current observed conditions. | [optional] 
**FailedClusters** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1FailedCluster**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1FailedCluster.md) | FailedClusters is an array of clusters which failed to apply the ClusterRole resources. | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleStatus

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleStatus() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleStatus`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleStatus instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleStatus object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleStatusWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleStatusWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleStatus`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleStatusWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleStatus object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleStatus) GetConditions() []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition`

GetConditions returns the Conditions field if non-nil, zero value otherwise.

### GetConditionsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleStatus) GetConditionsOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition, bool)`

GetConditionsOk returns a tuple with the Conditions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleStatus) SetConditions(v []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition)`

SetConditions sets Conditions field to given value.

### HasConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleStatus) HasConditions() bool`

HasConditions returns a boolean if a field has been set.

### GetFailedClusters

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleStatus) GetFailedClusters() []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1FailedCluster`

GetFailedClusters returns the FailedClusters field if non-nil, zero value otherwise.

### GetFailedClustersOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleStatus) GetFailedClustersOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1FailedCluster, bool)`

GetFailedClustersOk returns a tuple with the FailedClusters field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFailedClusters

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleStatus) SetFailedClusters(v []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1FailedCluster)`

SetFailedClusters sets FailedClusters field to given value.

### HasFailedClusters

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RoleStatus) HasFailedClusters() bool`

HasFailedClusters returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


