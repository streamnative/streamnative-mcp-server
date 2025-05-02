# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ClusterRoleStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Conditions** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition.md) | Conditions is an array of current observed conditions. | [optional] 
**FailedClusters** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1FailedCluster**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1FailedCluster.md) | FailedClusters is an array of clusters which failed to apply the ClusterRole resources. | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ClusterRoleStatus

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ClusterRoleStatus() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ClusterRoleStatus`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ClusterRoleStatus instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ClusterRoleStatus object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ClusterRoleStatusWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ClusterRoleStatusWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ClusterRoleStatus`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ClusterRoleStatusWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ClusterRoleStatus object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ClusterRoleStatus) GetConditions() []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition`

GetConditions returns the Conditions field if non-nil, zero value otherwise.

### GetConditionsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ClusterRoleStatus) GetConditionsOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition, bool)`

GetConditionsOk returns a tuple with the Conditions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ClusterRoleStatus) SetConditions(v []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition)`

SetConditions sets Conditions field to given value.

### HasConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ClusterRoleStatus) HasConditions() bool`

HasConditions returns a boolean if a field has been set.

### GetFailedClusters

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ClusterRoleStatus) GetFailedClusters() []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1FailedCluster`

GetFailedClusters returns the FailedClusters field if non-nil, zero value otherwise.

### GetFailedClustersOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ClusterRoleStatus) GetFailedClustersOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1FailedCluster, bool)`

GetFailedClustersOk returns a tuple with the FailedClusters field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFailedClusters

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ClusterRoleStatus) SetFailedClusters(v []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1FailedCluster)`

SetFailedClusters sets FailedClusters field to given value.

### HasFailedClusters

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ClusterRoleStatus) HasFailedClusters() bool`

HasFailedClusters returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


