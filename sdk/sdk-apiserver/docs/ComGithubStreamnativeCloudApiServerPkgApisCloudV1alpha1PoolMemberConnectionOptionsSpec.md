# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberConnectionOptionsSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ProxyUrl** | Pointer to **string** | ProxyUrl overrides the URL to K8S API. This is most typically done for SNI proxying. If ProxyUrl is set but TLSServerName is not, the *original* SNI value from the default endpoint will be used. If you also want to change the SNI header, you must also set TLSServerName. | [optional] 
**TlsServerName** | Pointer to **string** | TLSServerName is the SNI header name to set, overriding the default endpoint. This should include the hostname value, with no port. | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberConnectionOptionsSpec

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberConnectionOptionsSpec() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberConnectionOptionsSpec`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberConnectionOptionsSpec instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberConnectionOptionsSpec object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberConnectionOptionsSpecWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberConnectionOptionsSpecWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberConnectionOptionsSpec`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberConnectionOptionsSpecWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberConnectionOptionsSpec object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetProxyUrl

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberConnectionOptionsSpec) GetProxyUrl() string`

GetProxyUrl returns the ProxyUrl field if non-nil, zero value otherwise.

### GetProxyUrlOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberConnectionOptionsSpec) GetProxyUrlOk() (*string, bool)`

GetProxyUrlOk returns a tuple with the ProxyUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProxyUrl

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberConnectionOptionsSpec) SetProxyUrl(v string)`

SetProxyUrl sets ProxyUrl field to given value.

### HasProxyUrl

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberConnectionOptionsSpec) HasProxyUrl() bool`

HasProxyUrl returns a boolean if a field has been set.

### GetTlsServerName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberConnectionOptionsSpec) GetTlsServerName() string`

GetTlsServerName returns the TlsServerName field if non-nil, zero value otherwise.

### GetTlsServerNameOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberConnectionOptionsSpec) GetTlsServerNameOk() (*string, bool)`

GetTlsServerNameOk returns a tuple with the TlsServerName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTlsServerName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberConnectionOptionsSpec) SetTlsServerName(v string)`

SetTlsServerName sets TlsServerName field to given value.

### HasTlsServerName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberConnectionOptionsSpec) HasTlsServerName() bool`

HasTlsServerName returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


