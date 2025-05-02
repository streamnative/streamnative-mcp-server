# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeper

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AutoScalingPolicy** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AutoScalingPolicy**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AutoScalingPolicy.md) |  | [optional] 
**Image** | Pointer to **string** | Image name is the name of the image to deploy. | [optional] 
**Replicas** | **int32** | Replicas is the expected size of the bookkeeper cluster. | 
**ResourceSpec** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeperNodeResourceSpec**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeperNodeResourceSpec.md) |  | [optional] 
**Resources** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookkeeperNodeResource**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookkeeperNodeResource.md) |  | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeper

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeper(replicas int32, ) *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeper`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeper instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeper object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeperWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeperWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeper`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeperWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeper object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAutoScalingPolicy

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeper) GetAutoScalingPolicy() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AutoScalingPolicy`

GetAutoScalingPolicy returns the AutoScalingPolicy field if non-nil, zero value otherwise.

### GetAutoScalingPolicyOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeper) GetAutoScalingPolicyOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AutoScalingPolicy, bool)`

GetAutoScalingPolicyOk returns a tuple with the AutoScalingPolicy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAutoScalingPolicy

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeper) SetAutoScalingPolicy(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AutoScalingPolicy)`

SetAutoScalingPolicy sets AutoScalingPolicy field to given value.

### HasAutoScalingPolicy

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeper) HasAutoScalingPolicy() bool`

HasAutoScalingPolicy returns a boolean if a field has been set.

### GetImage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeper) GetImage() string`

GetImage returns the Image field if non-nil, zero value otherwise.

### GetImageOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeper) GetImageOk() (*string, bool)`

GetImageOk returns a tuple with the Image field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeper) SetImage(v string)`

SetImage sets Image field to given value.

### HasImage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeper) HasImage() bool`

HasImage returns a boolean if a field has been set.

### GetReplicas

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeper) GetReplicas() int32`

GetReplicas returns the Replicas field if non-nil, zero value otherwise.

### GetReplicasOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeper) GetReplicasOk() (*int32, bool)`

GetReplicasOk returns a tuple with the Replicas field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReplicas

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeper) SetReplicas(v int32)`

SetReplicas sets Replicas field to given value.


### GetResourceSpec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeper) GetResourceSpec() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeperNodeResourceSpec`

GetResourceSpec returns the ResourceSpec field if non-nil, zero value otherwise.

### GetResourceSpecOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeper) GetResourceSpecOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeperNodeResourceSpec, bool)`

GetResourceSpecOk returns a tuple with the ResourceSpec field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResourceSpec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeper) SetResourceSpec(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeperNodeResourceSpec)`

SetResourceSpec sets ResourceSpec field to given value.

### HasResourceSpec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeper) HasResourceSpec() bool`

HasResourceSpec returns a boolean if a field has been set.

### GetResources

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeper) GetResources() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookkeeperNodeResource`

GetResources returns the Resources field if non-nil, zero value otherwise.

### GetResourcesOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeper) GetResourcesOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookkeeperNodeResource, bool)`

GetResourcesOk returns a tuple with the Resources field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResources

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeper) SetResources(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookkeeperNodeResource)`

SetResources sets Resources field to given value.

### HasResources

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookKeeper) HasResources() bool`

HasResources returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


