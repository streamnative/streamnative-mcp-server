# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GenericPoolMemberSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Endpoint** | Pointer to **string** | Endpoint is *either* a full URL, or a hostname/port to point to the master | [optional] 
**Location** | **string** | Location is the location of the cluster. | 
**MasterAuth** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MasterAuth**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MasterAuth.md) |  | [optional] 
**TlsServerName** | Pointer to **string** | TLSServerName is the SNI header name to set, overridding the default. This is just the hostname and no port | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GenericPoolMemberSpec

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GenericPoolMemberSpec(location string, ) *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GenericPoolMemberSpec`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GenericPoolMemberSpec instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GenericPoolMemberSpec object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GenericPoolMemberSpecWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GenericPoolMemberSpecWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GenericPoolMemberSpec`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GenericPoolMemberSpecWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GenericPoolMemberSpec object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEndpoint

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GenericPoolMemberSpec) GetEndpoint() string`

GetEndpoint returns the Endpoint field if non-nil, zero value otherwise.

### GetEndpointOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GenericPoolMemberSpec) GetEndpointOk() (*string, bool)`

GetEndpointOk returns a tuple with the Endpoint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndpoint

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GenericPoolMemberSpec) SetEndpoint(v string)`

SetEndpoint sets Endpoint field to given value.

### HasEndpoint

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GenericPoolMemberSpec) HasEndpoint() bool`

HasEndpoint returns a boolean if a field has been set.

### GetLocation

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GenericPoolMemberSpec) GetLocation() string`

GetLocation returns the Location field if non-nil, zero value otherwise.

### GetLocationOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GenericPoolMemberSpec) GetLocationOk() (*string, bool)`

GetLocationOk returns a tuple with the Location field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLocation

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GenericPoolMemberSpec) SetLocation(v string)`

SetLocation sets Location field to given value.


### GetMasterAuth

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GenericPoolMemberSpec) GetMasterAuth() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MasterAuth`

GetMasterAuth returns the MasterAuth field if non-nil, zero value otherwise.

### GetMasterAuthOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GenericPoolMemberSpec) GetMasterAuthOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MasterAuth, bool)`

GetMasterAuthOk returns a tuple with the MasterAuth field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMasterAuth

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GenericPoolMemberSpec) SetMasterAuth(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MasterAuth)`

SetMasterAuth sets MasterAuth field to given value.

### HasMasterAuth

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GenericPoolMemberSpec) HasMasterAuth() bool`

HasMasterAuth returns a boolean if a field has been set.

### GetTlsServerName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GenericPoolMemberSpec) GetTlsServerName() string`

GetTlsServerName returns the TlsServerName field if non-nil, zero value otherwise.

### GetTlsServerNameOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GenericPoolMemberSpec) GetTlsServerNameOk() (*string, bool)`

GetTlsServerNameOk returns a tuple with the TlsServerName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTlsServerName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GenericPoolMemberSpec) SetTlsServerName(v string)`

SetTlsServerName sets TlsServerName field to given value.

### HasTlsServerName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GenericPoolMemberSpec) HasTlsServerName() bool`

HasTlsServerName returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


