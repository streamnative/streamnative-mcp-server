# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewayStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Conditions** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition.md) | Conditions is an array of current observed conditions. | [optional] 
**ObservedGeneration** | **int64** | ObservedGeneration is the most recent generation observed by the pulsargateway controller. | 
**PrivateServiceIds** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PrivateServiceId**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PrivateServiceId.md) | PrivateServiceIds are the id of the private endpoint services, only exposed when the access type is private. | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewayStatus

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewayStatus(observedGeneration int64, ) *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewayStatus`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewayStatus instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewayStatus object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewayStatusWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewayStatusWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewayStatus`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewayStatusWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewayStatus object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewayStatus) GetConditions() []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition`

GetConditions returns the Conditions field if non-nil, zero value otherwise.

### GetConditionsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewayStatus) GetConditionsOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition, bool)`

GetConditionsOk returns a tuple with the Conditions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewayStatus) SetConditions(v []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition)`

SetConditions sets Conditions field to given value.

### HasConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewayStatus) HasConditions() bool`

HasConditions returns a boolean if a field has been set.

### GetObservedGeneration

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewayStatus) GetObservedGeneration() int64`

GetObservedGeneration returns the ObservedGeneration field if non-nil, zero value otherwise.

### GetObservedGenerationOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewayStatus) GetObservedGenerationOk() (*int64, bool)`

GetObservedGenerationOk returns a tuple with the ObservedGeneration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObservedGeneration

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewayStatus) SetObservedGeneration(v int64)`

SetObservedGeneration sets ObservedGeneration field to given value.


### GetPrivateServiceIds

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewayStatus) GetPrivateServiceIds() []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PrivateServiceId`

GetPrivateServiceIds returns the PrivateServiceIds field if non-nil, zero value otherwise.

### GetPrivateServiceIdsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewayStatus) GetPrivateServiceIdsOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PrivateServiceId, bool)`

GetPrivateServiceIdsOk returns a tuple with the PrivateServiceIds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrivateServiceIds

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewayStatus) SetPrivateServiceIds(v []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PrivateServiceId)`

SetPrivateServiceIds sets PrivateServiceIds field to given value.

### HasPrivateServiceIds

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarGatewayStatus) HasPrivateServiceIds() bool`

HasPrivateServiceIds returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


