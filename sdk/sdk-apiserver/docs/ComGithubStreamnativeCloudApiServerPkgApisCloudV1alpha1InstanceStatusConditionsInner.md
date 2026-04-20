# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceStatusConditionsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**LastTransitionTime** | Pointer to **time.Time** | Time is a wrapper around time.Time which supports correct marshaling to YAML and JSON.  Wrappers are provided for many of the factory methods that the time package offers. | [optional] 
**Message** | Pointer to **string** |  | [optional] 
**ObservedGeneration** | Pointer to **int64** | observedGeneration represents the .metadata.generation that the condition was set based upon. For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date with respect to the current state of the instance. | [optional] 
**Reason** | Pointer to **string** |  | [optional] 
**Status** | **string** |  | [default to ""]
**Type** | **string** |  | [default to ""]

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceStatusConditionsInner

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceStatusConditionsInner(status string, type_ string, ) *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceStatusConditionsInner`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceStatusConditionsInner instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceStatusConditionsInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceStatusConditionsInnerWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceStatusConditionsInnerWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceStatusConditionsInner`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceStatusConditionsInnerWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceStatusConditionsInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetLastTransitionTime

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceStatusConditionsInner) GetLastTransitionTime() time.Time`

GetLastTransitionTime returns the LastTransitionTime field if non-nil, zero value otherwise.

### GetLastTransitionTimeOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceStatusConditionsInner) GetLastTransitionTimeOk() (*time.Time, bool)`

GetLastTransitionTimeOk returns a tuple with the LastTransitionTime field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastTransitionTime

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceStatusConditionsInner) SetLastTransitionTime(v time.Time)`

SetLastTransitionTime sets LastTransitionTime field to given value.

### HasLastTransitionTime

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceStatusConditionsInner) HasLastTransitionTime() bool`

HasLastTransitionTime returns a boolean if a field has been set.

### GetMessage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceStatusConditionsInner) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceStatusConditionsInner) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceStatusConditionsInner) SetMessage(v string)`

SetMessage sets Message field to given value.

### HasMessage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceStatusConditionsInner) HasMessage() bool`

HasMessage returns a boolean if a field has been set.

### GetObservedGeneration

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceStatusConditionsInner) GetObservedGeneration() int64`

GetObservedGeneration returns the ObservedGeneration field if non-nil, zero value otherwise.

### GetObservedGenerationOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceStatusConditionsInner) GetObservedGenerationOk() (*int64, bool)`

GetObservedGenerationOk returns a tuple with the ObservedGeneration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObservedGeneration

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceStatusConditionsInner) SetObservedGeneration(v int64)`

SetObservedGeneration sets ObservedGeneration field to given value.

### HasObservedGeneration

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceStatusConditionsInner) HasObservedGeneration() bool`

HasObservedGeneration returns a boolean if a field has been set.

### GetReason

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceStatusConditionsInner) GetReason() string`

GetReason returns the Reason field if non-nil, zero value otherwise.

### GetReasonOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceStatusConditionsInner) GetReasonOk() (*string, bool)`

GetReasonOk returns a tuple with the Reason field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReason

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceStatusConditionsInner) SetReason(v string)`

SetReason sets Reason field to given value.

### HasReason

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceStatusConditionsInner) HasReason() bool`

HasReason returns a boolean if a field has been set.

### GetStatus

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceStatusConditionsInner) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceStatusConditionsInner) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceStatusConditionsInner) SetStatus(v string)`

SetStatus sets Status field to given value.


### GetType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceStatusConditionsInner) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceStatusConditionsInner) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceStatusConditionsInner) SetType(v string)`

SetType sets Type field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


