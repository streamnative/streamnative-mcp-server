# V1WeightedPodAffinityTerm

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**PodAffinityTerm** | [**V1PodAffinityTerm**](V1PodAffinityTerm.md) |  | 
**Weight** | **int32** | weight associated with matching the corresponding podAffinityTerm, in the range 1-100. | 

## Methods

### NewV1WeightedPodAffinityTerm

`func NewV1WeightedPodAffinityTerm(podAffinityTerm V1PodAffinityTerm, weight int32, ) *V1WeightedPodAffinityTerm`

NewV1WeightedPodAffinityTerm instantiates a new V1WeightedPodAffinityTerm object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewV1WeightedPodAffinityTermWithDefaults

`func NewV1WeightedPodAffinityTermWithDefaults() *V1WeightedPodAffinityTerm`

NewV1WeightedPodAffinityTermWithDefaults instantiates a new V1WeightedPodAffinityTerm object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPodAffinityTerm

`func (o *V1WeightedPodAffinityTerm) GetPodAffinityTerm() V1PodAffinityTerm`

GetPodAffinityTerm returns the PodAffinityTerm field if non-nil, zero value otherwise.

### GetPodAffinityTermOk

`func (o *V1WeightedPodAffinityTerm) GetPodAffinityTermOk() (*V1PodAffinityTerm, bool)`

GetPodAffinityTermOk returns a tuple with the PodAffinityTerm field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPodAffinityTerm

`func (o *V1WeightedPodAffinityTerm) SetPodAffinityTerm(v V1PodAffinityTerm)`

SetPodAffinityTerm sets PodAffinityTerm field to given value.


### GetWeight

`func (o *V1WeightedPodAffinityTerm) GetWeight() int32`

GetWeight returns the Weight field if non-nil, zero value otherwise.

### GetWeightOk

`func (o *V1WeightedPodAffinityTerm) GetWeightOk() (*int32, bool)`

GetWeightOk returns a tuple with the Weight field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWeight

`func (o *V1WeightedPodAffinityTerm) SetWeight(v int32)`

SetWeight sets Weight field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


