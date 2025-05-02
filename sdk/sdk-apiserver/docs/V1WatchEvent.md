# V1WatchEvent

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Object** | **map[string]interface{}** | Object is:  * If Type is Added or Modified: the new state of the object.  * If Type is Deleted: the state of the object immediately before deletion.  * If Type is Error: *Status is recommended; other types may make sense    depending on context. | 
**Type** | **string** |  | 

## Methods

### NewV1WatchEvent

`func NewV1WatchEvent(object map[string]interface{}, type_ string, ) *V1WatchEvent`

NewV1WatchEvent instantiates a new V1WatchEvent object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewV1WatchEventWithDefaults

`func NewV1WatchEventWithDefaults() *V1WatchEvent`

NewV1WatchEventWithDefaults instantiates a new V1WatchEvent object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetObject

`func (o *V1WatchEvent) GetObject() map[string]interface{}`

GetObject returns the Object field if non-nil, zero value otherwise.

### GetObjectOk

`func (o *V1WatchEvent) GetObjectOk() (*map[string]interface{}, bool)`

GetObjectOk returns a tuple with the Object field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObject

`func (o *V1WatchEvent) SetObject(v map[string]interface{})`

SetObject sets Object field to given value.


### GetType

`func (o *V1WatchEvent) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *V1WatchEvent) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *V1WatchEvent) SetType(v string)`

SetType sets Type field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


