# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ClusterConnectionString** | Pointer to **string** | ClusterConnectionString is a connection string for client connectivity within the cluster. | [optional] 
**Conditions** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2Condition**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2Condition.md) | Conditions is an array of current observed conditions. | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetStatus

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetStatus() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetStatus`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetStatus instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetStatus object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetStatusWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetStatusWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetStatus`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetStatusWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetStatus object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetClusterConnectionString

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetStatus) GetClusterConnectionString() string`

GetClusterConnectionString returns the ClusterConnectionString field if non-nil, zero value otherwise.

### GetClusterConnectionStringOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetStatus) GetClusterConnectionStringOk() (*string, bool)`

GetClusterConnectionStringOk returns a tuple with the ClusterConnectionString field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClusterConnectionString

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetStatus) SetClusterConnectionString(v string)`

SetClusterConnectionString sets ClusterConnectionString field to given value.

### HasClusterConnectionString

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetStatus) HasClusterConnectionString() bool`

HasClusterConnectionString returns a boolean if a field has been set.

### GetConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetStatus) GetConditions() []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2Condition`

GetConditions returns the Conditions field if non-nil, zero value otherwise.

### GetConditionsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetStatus) GetConditionsOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2Condition, bool)`

GetConditionsOk returns a tuple with the Conditions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetStatus) SetConditions(v []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2Condition)`

SetConditions sets Conditions field to given value.

### HasConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2ZooKeeperSetStatus) HasConditions() bool`

HasConditions returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


