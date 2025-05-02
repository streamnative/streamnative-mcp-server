# ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1Condition

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**LastTransitionTime** | Pointer to **time.Time** | Time is a wrapper around time.Time which supports correct marshaling to YAML and JSON.  Wrappers are provided for many of the factory methods that the time package offers. | [optional] 
**Message** | Pointer to **string** |  | [optional] 
**Reason** | Pointer to **string** |  | [optional] 
**Status** | **string** |  | 
**Type** | **string** |  | 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1Condition

`func NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1Condition(status string, type_ string, ) *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1Condition`

NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1Condition instantiates a new ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1Condition object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1ConditionWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1ConditionWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1Condition`

NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1ConditionWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1Condition object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetLastTransitionTime

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1Condition) GetLastTransitionTime() time.Time`

GetLastTransitionTime returns the LastTransitionTime field if non-nil, zero value otherwise.

### GetLastTransitionTimeOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1Condition) GetLastTransitionTimeOk() (*time.Time, bool)`

GetLastTransitionTimeOk returns a tuple with the LastTransitionTime field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastTransitionTime

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1Condition) SetLastTransitionTime(v time.Time)`

SetLastTransitionTime sets LastTransitionTime field to given value.

### HasLastTransitionTime

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1Condition) HasLastTransitionTime() bool`

HasLastTransitionTime returns a boolean if a field has been set.

### GetMessage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1Condition) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1Condition) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1Condition) SetMessage(v string)`

SetMessage sets Message field to given value.

### HasMessage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1Condition) HasMessage() bool`

HasMessage returns a boolean if a field has been set.

### GetReason

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1Condition) GetReason() string`

GetReason returns the Reason field if non-nil, zero value otherwise.

### GetReasonOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1Condition) GetReasonOk() (*string, bool)`

GetReasonOk returns a tuple with the Reason field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReason

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1Condition) SetReason(v string)`

SetReason sets Reason field to given value.

### HasReason

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1Condition) HasReason() bool`

HasReason returns a boolean if a field has been set.

### GetStatus

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1Condition) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1Condition) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1Condition) SetStatus(v string)`

SetStatus sets Status field to given value.


### GetType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1Condition) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1Condition) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1Condition) SetType(v string)`

SetType sets Type field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


