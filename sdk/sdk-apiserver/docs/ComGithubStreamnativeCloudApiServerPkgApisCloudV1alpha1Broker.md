# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Broker

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AutoScalingPolicy** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AutoScalingPolicy**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AutoScalingPolicy.md) |  | [optional] 
**Image** | Pointer to **string** | Image name is the name of the image to deploy. | [optional] 
**Replicas** | **int32** | Replicas is the expected size of the broker cluster. | 
**ResourceSpec** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BrokerNodeResourceSpec**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BrokerNodeResourceSpec.md) |  | [optional] 
**Resources** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1DefaultNodeResource**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1DefaultNodeResource.md) |  | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Broker

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Broker(replicas int32, ) *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Broker`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Broker instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Broker object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BrokerWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BrokerWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Broker`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BrokerWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Broker object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAutoScalingPolicy

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Broker) GetAutoScalingPolicy() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AutoScalingPolicy`

GetAutoScalingPolicy returns the AutoScalingPolicy field if non-nil, zero value otherwise.

### GetAutoScalingPolicyOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Broker) GetAutoScalingPolicyOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AutoScalingPolicy, bool)`

GetAutoScalingPolicyOk returns a tuple with the AutoScalingPolicy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAutoScalingPolicy

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Broker) SetAutoScalingPolicy(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AutoScalingPolicy)`

SetAutoScalingPolicy sets AutoScalingPolicy field to given value.

### HasAutoScalingPolicy

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Broker) HasAutoScalingPolicy() bool`

HasAutoScalingPolicy returns a boolean if a field has been set.

### GetImage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Broker) GetImage() string`

GetImage returns the Image field if non-nil, zero value otherwise.

### GetImageOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Broker) GetImageOk() (*string, bool)`

GetImageOk returns a tuple with the Image field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Broker) SetImage(v string)`

SetImage sets Image field to given value.

### HasImage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Broker) HasImage() bool`

HasImage returns a boolean if a field has been set.

### GetReplicas

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Broker) GetReplicas() int32`

GetReplicas returns the Replicas field if non-nil, zero value otherwise.

### GetReplicasOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Broker) GetReplicasOk() (*int32, bool)`

GetReplicasOk returns a tuple with the Replicas field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReplicas

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Broker) SetReplicas(v int32)`

SetReplicas sets Replicas field to given value.


### GetResourceSpec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Broker) GetResourceSpec() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BrokerNodeResourceSpec`

GetResourceSpec returns the ResourceSpec field if non-nil, zero value otherwise.

### GetResourceSpecOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Broker) GetResourceSpecOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BrokerNodeResourceSpec, bool)`

GetResourceSpecOk returns a tuple with the ResourceSpec field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResourceSpec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Broker) SetResourceSpec(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BrokerNodeResourceSpec)`

SetResourceSpec sets ResourceSpec field to given value.

### HasResourceSpec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Broker) HasResourceSpec() bool`

HasResourceSpec returns a boolean if a field has been set.

### GetResources

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Broker) GetResources() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1DefaultNodeResource`

GetResources returns the Resources field if non-nil, zero value otherwise.

### GetResourcesOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Broker) GetResourcesOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1DefaultNodeResource, bool)`

GetResourcesOk returns a tuple with the Resources field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResources

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Broker) SetResources(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1DefaultNodeResource)`

SetResources sets Resources field to given value.

### HasResources

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Broker) HasResources() bool`

HasResources returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


