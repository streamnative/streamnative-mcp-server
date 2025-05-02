# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Gateway

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Access** | Pointer to **string** | Access is the access type of the pulsar gateway, available values are public or private. It is immutable, with the default value public. | [optional] 
**PrivateService** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PrivateService**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PrivateService.md) |  | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Gateway

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Gateway() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Gateway`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Gateway instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Gateway object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GatewayWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GatewayWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Gateway`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GatewayWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Gateway object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAccess

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Gateway) GetAccess() string`

GetAccess returns the Access field if non-nil, zero value otherwise.

### GetAccessOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Gateway) GetAccessOk() (*string, bool)`

GetAccessOk returns a tuple with the Access field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccess

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Gateway) SetAccess(v string)`

SetAccess sets Access field to given value.

### HasAccess

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Gateway) HasAccess() bool`

HasAccess returns a boolean if a field has been set.

### GetPrivateService

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Gateway) GetPrivateService() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PrivateService`

GetPrivateService returns the PrivateService field if non-nil, zero value otherwise.

### GetPrivateServiceOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Gateway) GetPrivateServiceOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PrivateService, bool)`

GetPrivateServiceOk returns a tuple with the PrivateService field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrivateService

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Gateway) SetPrivateService(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PrivateService)`

SetPrivateService sets PrivateService field to given value.

### HasPrivateService

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Gateway) HasPrivateService() bool`

HasPrivateService returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


