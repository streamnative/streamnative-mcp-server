# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ObservedGeneration** | Pointer to **int64** | ObservedGeneration represents the .metadata.generation that the status was set based upon | [optional] 
**Conditions** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition.md) | Conditions is an array of current observed conditions | [optional] 
**DeploymentType** | Pointer to **string** | DeploymentType indicates the deployment type set via associated pool | [optional] 
**Image** | Pointer to **string** | Image is the Kafka container image being used | [optional] 
**Version** | Pointer to **string** | Version is the Kafka version being used | [optional] 
**KafkaRestProxy** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaRestProxyStatus**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaRestProxyStatus.md) |  | [optional] 
**ServiceEndpoints** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaServiceEndpoint**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaServiceEndpoint.md) | ServiceEndpoints defines the service endpoints for the Kafka cluster | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterStatus

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterStatus() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterStatus`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterStatus instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterStatus object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterStatusWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterStatusWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterStatus`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterStatusWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterStatus object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetObservedGeneration

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterStatus) GetObservedGeneration() int64`

GetObservedGeneration returns the ObservedGeneration field if non-nil, zero value otherwise.

### GetObservedGenerationOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterStatus) GetObservedGenerationOk() (*int64, bool)`

GetObservedGenerationOk returns a tuple with the ObservedGeneration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObservedGeneration

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterStatus) SetObservedGeneration(v int64)`

SetObservedGeneration sets ObservedGeneration field to given value.

### HasObservedGeneration

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterStatus) HasObservedGeneration() bool`

HasObservedGeneration returns a boolean if a field has been set.

### GetConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterStatus) GetConditions() []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition`

GetConditions returns the Conditions field if non-nil, zero value otherwise.

### GetConditionsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterStatus) GetConditionsOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition, bool)`

GetConditionsOk returns a tuple with the Conditions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterStatus) SetConditions(v []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition)`

SetConditions sets Conditions field to given value.

### HasConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterStatus) HasConditions() bool`

HasConditions returns a boolean if a field has been set.

### GetDeploymentType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterStatus) GetDeploymentType() string`

GetDeploymentType returns the DeploymentType field if non-nil, zero value otherwise.

### GetDeploymentTypeOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterStatus) GetDeploymentTypeOk() (*string, bool)`

GetDeploymentTypeOk returns a tuple with the DeploymentType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeploymentType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterStatus) SetDeploymentType(v string)`

SetDeploymentType sets DeploymentType field to given value.

### HasDeploymentType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterStatus) HasDeploymentType() bool`

HasDeploymentType returns a boolean if a field has been set.

### GetImage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterStatus) GetImage() string`

GetImage returns the Image field if non-nil, zero value otherwise.

### GetImageOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterStatus) GetImageOk() (*string, bool)`

GetImageOk returns a tuple with the Image field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterStatus) SetImage(v string)`

SetImage sets Image field to given value.

### HasImage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterStatus) HasImage() bool`

HasImage returns a boolean if a field has been set.

### GetVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterStatus) GetVersion() string`

GetVersion returns the Version field if non-nil, zero value otherwise.

### GetVersionOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterStatus) GetVersionOk() (*string, bool)`

GetVersionOk returns a tuple with the Version field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterStatus) SetVersion(v string)`

SetVersion sets Version field to given value.

### HasVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterStatus) HasVersion() bool`

HasVersion returns a boolean if a field has been set.

### GetKafkaRestProxy

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterStatus) GetKafkaRestProxy() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaRestProxyStatus`

GetKafkaRestProxy returns the KafkaRestProxy field if non-nil, zero value otherwise.

### GetKafkaRestProxyOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterStatus) GetKafkaRestProxyOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaRestProxyStatus, bool)`

GetKafkaRestProxyOk returns a tuple with the KafkaRestProxy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKafkaRestProxy

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterStatus) SetKafkaRestProxy(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaRestProxyStatus)`

SetKafkaRestProxy sets KafkaRestProxy field to given value.

### HasKafkaRestProxy

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterStatus) HasKafkaRestProxy() bool`

HasKafkaRestProxy returns a boolean if a field has been set.

### GetServiceEndpoints

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterStatus) GetServiceEndpoints() []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaServiceEndpoint`

GetServiceEndpoints returns the ServiceEndpoints field if non-nil, zero value otherwise.

### GetServiceEndpointsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterStatus) GetServiceEndpointsOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaServiceEndpoint, bool)`

GetServiceEndpointsOk returns a tuple with the ServiceEndpoints field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServiceEndpoints

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterStatus) SetServiceEndpoints(v []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaServiceEndpoint)`

SetServiceEndpoints sets ServiceEndpoints field to given value.

### HasServiceEndpoints

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterStatus) HasServiceEndpoints() bool`

HasServiceEndpoints returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


