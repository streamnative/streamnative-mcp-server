# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Taint

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Effect** | **string** | Required. The effect of the taint on workloads that do not tolerate the taint. Valid effects are NoSchedule, PreferNoSchedule and NoExecute. | 
**Key** | **string** | Required. The taint key to be applied to a workload cluster. | 
**TimeAdded** | Pointer to **time.Time** | TimeAdded represents the time at which the taint was added. | [optional] 
**Value** | Pointer to **string** | Optional. The taint value corresponding to the taint key. | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Taint

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Taint(effect string, key string, ) *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Taint`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Taint instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Taint object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1TaintWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1TaintWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Taint`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1TaintWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Taint object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEffect

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Taint) GetEffect() string`

GetEffect returns the Effect field if non-nil, zero value otherwise.

### GetEffectOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Taint) GetEffectOk() (*string, bool)`

GetEffectOk returns a tuple with the Effect field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEffect

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Taint) SetEffect(v string)`

SetEffect sets Effect field to given value.


### GetKey

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Taint) GetKey() string`

GetKey returns the Key field if non-nil, zero value otherwise.

### GetKeyOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Taint) GetKeyOk() (*string, bool)`

GetKeyOk returns a tuple with the Key field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKey

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Taint) SetKey(v string)`

SetKey sets Key field to given value.


### GetTimeAdded

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Taint) GetTimeAdded() time.Time`

GetTimeAdded returns the TimeAdded field if non-nil, zero value otherwise.

### GetTimeAddedOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Taint) GetTimeAddedOk() (*time.Time, bool)`

GetTimeAddedOk returns a tuple with the TimeAdded field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimeAdded

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Taint) SetTimeAdded(v time.Time)`

SetTimeAdded sets TimeAdded field to given value.

### HasTimeAdded

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Taint) HasTimeAdded() bool`

HasTimeAdded returns a boolean if a field has been set.

### GetValue

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Taint) GetValue() string`

GetValue returns the Value field if non-nil, zero value otherwise.

### GetValueOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Taint) GetValueOk() (*string, bool)`

GetValueOk returns a tuple with the Value field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetValue

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Taint) SetValue(v string)`

SetValue sets Value field to given value.

### HasValue

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Taint) HasValue() bool`

HasValue returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


