# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Auth** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceAuth**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceAuth.md) |  | [optional] 
**AvailabilityMode** | **string** | AvailabilityMode decides whether pods of the same type in pulsar should be in one zone or multiple zones | 
**Plan** | Pointer to **string** | Plan is the subscription plan, will create a stripe subscription if not empty deprecated: 1.16 | [optional] 
**PoolRef** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolRef**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolRef.md) |  | [optional] 
**Type** | Pointer to **string** | Type defines the instance specialization type: - standard: a standard deployment of Pulsar, BookKeeper, and ZooKeeper. - dedicated: a dedicated deployment of classic engine or ursa engine. - serverless: a serverless deployment of Pulsar, shared BookKeeper, and shared oxia. - byoc: bring your own cloud. | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceSpec

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceSpec(availabilityMode string, ) *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceSpec`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceSpec instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceSpec object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceSpecWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceSpecWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceSpec`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceSpecWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceSpec object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAuth

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceSpec) GetAuth() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceAuth`

GetAuth returns the Auth field if non-nil, zero value otherwise.

### GetAuthOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceSpec) GetAuthOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceAuth, bool)`

GetAuthOk returns a tuple with the Auth field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuth

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceSpec) SetAuth(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceAuth)`

SetAuth sets Auth field to given value.

### HasAuth

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceSpec) HasAuth() bool`

HasAuth returns a boolean if a field has been set.

### GetAvailabilityMode

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceSpec) GetAvailabilityMode() string`

GetAvailabilityMode returns the AvailabilityMode field if non-nil, zero value otherwise.

### GetAvailabilityModeOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceSpec) GetAvailabilityModeOk() (*string, bool)`

GetAvailabilityModeOk returns a tuple with the AvailabilityMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAvailabilityMode

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceSpec) SetAvailabilityMode(v string)`

SetAvailabilityMode sets AvailabilityMode field to given value.


### GetPlan

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceSpec) GetPlan() string`

GetPlan returns the Plan field if non-nil, zero value otherwise.

### GetPlanOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceSpec) GetPlanOk() (*string, bool)`

GetPlanOk returns a tuple with the Plan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPlan

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceSpec) SetPlan(v string)`

SetPlan sets Plan field to given value.

### HasPlan

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceSpec) HasPlan() bool`

HasPlan returns a boolean if a field has been set.

### GetPoolRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceSpec) GetPoolRef() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolRef`

GetPoolRef returns the PoolRef field if non-nil, zero value otherwise.

### GetPoolRefOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceSpec) GetPoolRefOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolRef, bool)`

GetPoolRefOk returns a tuple with the PoolRef field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPoolRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceSpec) SetPoolRef(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolRef)`

SetPoolRef sets PoolRef field to given value.

### HasPoolRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceSpec) HasPoolRef() bool`

HasPoolRef returns a boolean if a field has been set.

### GetType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceSpec) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceSpec) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceSpec) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceSpec) HasType() bool`

HasType returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


