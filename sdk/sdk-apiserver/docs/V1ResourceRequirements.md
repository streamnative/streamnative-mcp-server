# V1ResourceRequirements

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Limits** | Pointer to **map[string]string** | Limits describes the maximum amount of compute resources allowed. More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/ | [optional] 
**Requests** | Pointer to **map[string]string** | Requests describes the minimum amount of compute resources required. If Requests is omitted for a container, it defaults to Limits if that is explicitly specified, otherwise to an implementation-defined value. More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/ | [optional] 

## Methods

### NewV1ResourceRequirements

`func NewV1ResourceRequirements() *V1ResourceRequirements`

NewV1ResourceRequirements instantiates a new V1ResourceRequirements object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewV1ResourceRequirementsWithDefaults

`func NewV1ResourceRequirementsWithDefaults() *V1ResourceRequirements`

NewV1ResourceRequirementsWithDefaults instantiates a new V1ResourceRequirements object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetLimits

`func (o *V1ResourceRequirements) GetLimits() map[string]string`

GetLimits returns the Limits field if non-nil, zero value otherwise.

### GetLimitsOk

`func (o *V1ResourceRequirements) GetLimitsOk() (*map[string]string, bool)`

GetLimitsOk returns a tuple with the Limits field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLimits

`func (o *V1ResourceRequirements) SetLimits(v map[string]string)`

SetLimits sets Limits field to given value.

### HasLimits

`func (o *V1ResourceRequirements) HasLimits() bool`

HasLimits returns a boolean if a field has been set.

### GetRequests

`func (o *V1ResourceRequirements) GetRequests() map[string]string`

GetRequests returns the Requests field if non-nil, zero value otherwise.

### GetRequestsOk

`func (o *V1ResourceRequirements) GetRequestsOk() (*map[string]string, bool)`

GetRequestsOk returns a tuple with the Requests field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRequests

`func (o *V1ResourceRequirements) SetRequests(v map[string]string)`

SetRequests sets Requests field to given value.

### HasRequests

`func (o *V1ResourceRequirements) HasRequests() bool`

HasRequests returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


