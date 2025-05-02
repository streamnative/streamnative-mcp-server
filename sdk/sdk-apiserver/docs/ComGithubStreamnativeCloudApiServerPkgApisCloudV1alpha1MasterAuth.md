# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MasterAuth

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ClientCertificate** | Pointer to **string** | ClientCertificate is base64-encoded public certificate used by clients to authenticate to the cluster endpoint. | [optional] 
**ClientKey** | Pointer to **string** | ClientKey is base64-encoded private key used by clients to authenticate to the cluster endpoint. | [optional] 
**ClusterCaCertificate** | Pointer to **string** | ClusterCaCertificate is base64-encoded public certificate that is the root of trust for the cluster. | [optional] 
**Exec** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ExecConfig**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ExecConfig.md) |  | [optional] 
**Password** | Pointer to **string** | Password is the password to use for HTTP basic authentication to the master endpoint. | [optional] 
**Username** | Pointer to **string** | Username is the username to use for HTTP basic authentication to the master endpoint. | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MasterAuth

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MasterAuth() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MasterAuth`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MasterAuth instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MasterAuth object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MasterAuthWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MasterAuthWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MasterAuth`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MasterAuthWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MasterAuth object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetClientCertificate

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MasterAuth) GetClientCertificate() string`

GetClientCertificate returns the ClientCertificate field if non-nil, zero value otherwise.

### GetClientCertificateOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MasterAuth) GetClientCertificateOk() (*string, bool)`

GetClientCertificateOk returns a tuple with the ClientCertificate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClientCertificate

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MasterAuth) SetClientCertificate(v string)`

SetClientCertificate sets ClientCertificate field to given value.

### HasClientCertificate

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MasterAuth) HasClientCertificate() bool`

HasClientCertificate returns a boolean if a field has been set.

### GetClientKey

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MasterAuth) GetClientKey() string`

GetClientKey returns the ClientKey field if non-nil, zero value otherwise.

### GetClientKeyOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MasterAuth) GetClientKeyOk() (*string, bool)`

GetClientKeyOk returns a tuple with the ClientKey field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClientKey

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MasterAuth) SetClientKey(v string)`

SetClientKey sets ClientKey field to given value.

### HasClientKey

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MasterAuth) HasClientKey() bool`

HasClientKey returns a boolean if a field has been set.

### GetClusterCaCertificate

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MasterAuth) GetClusterCaCertificate() string`

GetClusterCaCertificate returns the ClusterCaCertificate field if non-nil, zero value otherwise.

### GetClusterCaCertificateOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MasterAuth) GetClusterCaCertificateOk() (*string, bool)`

GetClusterCaCertificateOk returns a tuple with the ClusterCaCertificate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClusterCaCertificate

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MasterAuth) SetClusterCaCertificate(v string)`

SetClusterCaCertificate sets ClusterCaCertificate field to given value.

### HasClusterCaCertificate

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MasterAuth) HasClusterCaCertificate() bool`

HasClusterCaCertificate returns a boolean if a field has been set.

### GetExec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MasterAuth) GetExec() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ExecConfig`

GetExec returns the Exec field if non-nil, zero value otherwise.

### GetExecOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MasterAuth) GetExecOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ExecConfig, bool)`

GetExecOk returns a tuple with the Exec field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MasterAuth) SetExec(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ExecConfig)`

SetExec sets Exec field to given value.

### HasExec

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MasterAuth) HasExec() bool`

HasExec returns a boolean if a field has been set.

### GetPassword

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MasterAuth) GetPassword() string`

GetPassword returns the Password field if non-nil, zero value otherwise.

### GetPasswordOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MasterAuth) GetPasswordOk() (*string, bool)`

GetPasswordOk returns a tuple with the Password field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPassword

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MasterAuth) SetPassword(v string)`

SetPassword sets Password field to given value.

### HasPassword

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MasterAuth) HasPassword() bool`

HasPassword returns a boolean if a field has been set.

### GetUsername

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MasterAuth) GetUsername() string`

GetUsername returns the Username field if non-nil, zero value otherwise.

### GetUsernameOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MasterAuth) GetUsernameOk() (*string, bool)`

GetUsernameOk returns a tuple with the Username field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsername

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MasterAuth) SetUsername(v string)`

SetUsername sets Username field to given value.

### HasUsername

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MasterAuth) HasUsername() bool`

HasUsername returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


