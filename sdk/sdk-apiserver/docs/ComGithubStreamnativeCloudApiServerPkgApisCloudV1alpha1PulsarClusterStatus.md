# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Bookkeeper** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterComponentStatus**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterComponentStatus.md) |  | [optional] 
**Broker** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterComponentStatus**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterComponentStatus.md) |  | [optional] 
**Conditions** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition.md) | Conditions is an array of current observed conditions. | [optional] 
**DeploymentType** | Pointer to **string** | Deployment type set via associated pool | [optional] 
**InstanceType** | Pointer to **string** | Instance type, i.e. serverless or default | [optional] 
**Oxia** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterComponentStatus**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterComponentStatus.md) |  | [optional] 
**RbacStatus** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RBACStatus**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RBACStatus.md) |  | [optional] 
**Zookeeper** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterComponentStatus**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterComponentStatus.md) |  | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterStatus

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterStatus() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterStatus`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterStatus instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterStatus object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterStatusWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterStatusWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterStatus`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterStatusWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterStatus object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetBookkeeper

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterStatus) GetBookkeeper() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterComponentStatus`

GetBookkeeper returns the Bookkeeper field if non-nil, zero value otherwise.

### GetBookkeeperOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterStatus) GetBookkeeperOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterComponentStatus, bool)`

GetBookkeeperOk returns a tuple with the Bookkeeper field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBookkeeper

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterStatus) SetBookkeeper(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterComponentStatus)`

SetBookkeeper sets Bookkeeper field to given value.

### HasBookkeeper

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterStatus) HasBookkeeper() bool`

HasBookkeeper returns a boolean if a field has been set.

### GetBroker

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterStatus) GetBroker() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterComponentStatus`

GetBroker returns the Broker field if non-nil, zero value otherwise.

### GetBrokerOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterStatus) GetBrokerOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterComponentStatus, bool)`

GetBrokerOk returns a tuple with the Broker field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBroker

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterStatus) SetBroker(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterComponentStatus)`

SetBroker sets Broker field to given value.

### HasBroker

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterStatus) HasBroker() bool`

HasBroker returns a boolean if a field has been set.

### GetConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterStatus) GetConditions() []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition`

GetConditions returns the Conditions field if non-nil, zero value otherwise.

### GetConditionsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterStatus) GetConditionsOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition, bool)`

GetConditionsOk returns a tuple with the Conditions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterStatus) SetConditions(v []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition)`

SetConditions sets Conditions field to given value.

### HasConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterStatus) HasConditions() bool`

HasConditions returns a boolean if a field has been set.

### GetDeploymentType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterStatus) GetDeploymentType() string`

GetDeploymentType returns the DeploymentType field if non-nil, zero value otherwise.

### GetDeploymentTypeOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterStatus) GetDeploymentTypeOk() (*string, bool)`

GetDeploymentTypeOk returns a tuple with the DeploymentType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeploymentType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterStatus) SetDeploymentType(v string)`

SetDeploymentType sets DeploymentType field to given value.

### HasDeploymentType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterStatus) HasDeploymentType() bool`

HasDeploymentType returns a boolean if a field has been set.

### GetInstanceType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterStatus) GetInstanceType() string`

GetInstanceType returns the InstanceType field if non-nil, zero value otherwise.

### GetInstanceTypeOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterStatus) GetInstanceTypeOk() (*string, bool)`

GetInstanceTypeOk returns a tuple with the InstanceType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInstanceType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterStatus) SetInstanceType(v string)`

SetInstanceType sets InstanceType field to given value.

### HasInstanceType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterStatus) HasInstanceType() bool`

HasInstanceType returns a boolean if a field has been set.

### GetOxia

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterStatus) GetOxia() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterComponentStatus`

GetOxia returns the Oxia field if non-nil, zero value otherwise.

### GetOxiaOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterStatus) GetOxiaOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterComponentStatus, bool)`

GetOxiaOk returns a tuple with the Oxia field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOxia

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterStatus) SetOxia(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterComponentStatus)`

SetOxia sets Oxia field to given value.

### HasOxia

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterStatus) HasOxia() bool`

HasOxia returns a boolean if a field has been set.

### GetRbacStatus

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterStatus) GetRbacStatus() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RBACStatus`

GetRbacStatus returns the RbacStatus field if non-nil, zero value otherwise.

### GetRbacStatusOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterStatus) GetRbacStatusOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RBACStatus, bool)`

GetRbacStatusOk returns a tuple with the RbacStatus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRbacStatus

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterStatus) SetRbacStatus(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RBACStatus)`

SetRbacStatus sets RbacStatus field to given value.

### HasRbacStatus

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterStatus) HasRbacStatus() bool`

HasRbacStatus returns a boolean if a field has been set.

### GetZookeeper

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterStatus) GetZookeeper() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterComponentStatus`

GetZookeeper returns the Zookeeper field if non-nil, zero value otherwise.

### GetZookeeperOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterStatus) GetZookeeperOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterComponentStatus, bool)`

GetZookeeperOk returns a tuple with the Zookeeper field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetZookeeper

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterStatus) SetZookeeper(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterComponentStatus)`

SetZookeeper sets Zookeeper field to given value.

### HasZookeeper

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterStatus) HasZookeeper() bool`

HasZookeeper returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


