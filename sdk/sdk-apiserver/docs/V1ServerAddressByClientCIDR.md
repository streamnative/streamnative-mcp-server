# V1ServerAddressByClientCIDR

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ClientCIDR** | **string** | The CIDR with which clients can match their IP to figure out the server address that they should use. | 
**ServerAddress** | **string** | Address of this server, suitable for a client that matches the above CIDR. This can be a hostname, hostname:port, IP or IP:port. | 

## Methods

### NewV1ServerAddressByClientCIDR

`func NewV1ServerAddressByClientCIDR(clientCIDR string, serverAddress string, ) *V1ServerAddressByClientCIDR`

NewV1ServerAddressByClientCIDR instantiates a new V1ServerAddressByClientCIDR object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewV1ServerAddressByClientCIDRWithDefaults

`func NewV1ServerAddressByClientCIDRWithDefaults() *V1ServerAddressByClientCIDR`

NewV1ServerAddressByClientCIDRWithDefaults instantiates a new V1ServerAddressByClientCIDR object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetClientCIDR

`func (o *V1ServerAddressByClientCIDR) GetClientCIDR() string`

GetClientCIDR returns the ClientCIDR field if non-nil, zero value otherwise.

### GetClientCIDROk

`func (o *V1ServerAddressByClientCIDR) GetClientCIDROk() (*string, bool)`

GetClientCIDROk returns a tuple with the ClientCIDR field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClientCIDR

`func (o *V1ServerAddressByClientCIDR) SetClientCIDR(v string)`

SetClientCIDR sets ClientCIDR field to given value.


### GetServerAddress

`func (o *V1ServerAddressByClientCIDR) GetServerAddress() string`

GetServerAddress returns the ServerAddress field if non-nil, zero value otherwise.

### GetServerAddressOk

`func (o *V1ServerAddressByClientCIDR) GetServerAddressOk() (*string, bool)`

GetServerAddressOk returns a tuple with the ServerAddress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServerAddress

`func (o *V1ServerAddressByClientCIDR) SetServerAddress(v string)`

SetServerAddress sets ServerAddress field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


