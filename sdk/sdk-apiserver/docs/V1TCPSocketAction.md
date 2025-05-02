# V1TCPSocketAction

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Host** | Pointer to **string** | Optional: Host name to connect to, defaults to the pod IP. | [optional] 
**Port** | **map[string]interface{}** | Number or name of the port to access on the container. Number must be in the range 1 to 65535. Name must be an IANA_SVC_NAME. | 

## Methods

### NewV1TCPSocketAction

`func NewV1TCPSocketAction(port map[string]interface{}, ) *V1TCPSocketAction`

NewV1TCPSocketAction instantiates a new V1TCPSocketAction object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewV1TCPSocketActionWithDefaults

`func NewV1TCPSocketActionWithDefaults() *V1TCPSocketAction`

NewV1TCPSocketActionWithDefaults instantiates a new V1TCPSocketAction object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetHost

`func (o *V1TCPSocketAction) GetHost() string`

GetHost returns the Host field if non-nil, zero value otherwise.

### GetHostOk

`func (o *V1TCPSocketAction) GetHostOk() (*string, bool)`

GetHostOk returns a tuple with the Host field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHost

`func (o *V1TCPSocketAction) SetHost(v string)`

SetHost sets Host field to given value.

### HasHost

`func (o *V1TCPSocketAction) HasHost() bool`

HasHost returns a boolean if a field has been set.

### GetPort

`func (o *V1TCPSocketAction) GetPort() map[string]interface{}`

GetPort returns the Port field if non-nil, zero value otherwise.

### GetPortOk

`func (o *V1TCPSocketAction) GetPortOk() (*map[string]interface{}, bool)`

GetPortOk returns a tuple with the Port field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPort

`func (o *V1TCPSocketAction) SetPort(v map[string]interface{})`

SetPort sets Port field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


