# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Conditions** | Pointer to [**[]V1Condition**](V1Condition.md) | Conditions contains details for the current state of underlying resource | [optional] 
**DefaultGateway** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GatewayStatus**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GatewayStatus.md) |  | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentStatus

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentStatus() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentStatus`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentStatus instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentStatus object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentStatusWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentStatusWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentStatus`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentStatusWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentStatus object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentStatus) GetConditions() []V1Condition`

GetConditions returns the Conditions field if non-nil, zero value otherwise.

### GetConditionsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentStatus) GetConditionsOk() (*[]V1Condition, bool)`

GetConditionsOk returns a tuple with the Conditions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentStatus) SetConditions(v []V1Condition)`

SetConditions sets Conditions field to given value.

### HasConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentStatus) HasConditions() bool`

HasConditions returns a boolean if a field has been set.

### GetDefaultGateway

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentStatus) GetDefaultGateway() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GatewayStatus`

GetDefaultGateway returns the DefaultGateway field if non-nil, zero value otherwise.

### GetDefaultGatewayOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentStatus) GetDefaultGatewayOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GatewayStatus, bool)`

GetDefaultGatewayOk returns a tuple with the DefaultGateway field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDefaultGateway

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentStatus) SetDefaultGateway(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1GatewayStatus)`

SetDefaultGateway sets DefaultGateway field to given value.

### HasDefaultGateway

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudEnvironmentStatus) HasDefaultGateway() bool`

HasDefaultGateway returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


