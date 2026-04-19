# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**PoolRef** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceSpecPoolRef**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceSpecPoolRef.md) |  | [optional] 
**Type** | Pointer to **string** | Type defines the instance specialization type: - dedicated: a dedicated deployment - dedicated-pro: a dedicated deployment with enhanced features (e.g., CTC uses dedicated-pro) - serverless: a serverless deployment - byoc: bring your own cloud - byoc-pro: bring your own cloud with enhanced features | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceSpec

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceSpec() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceSpec`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceSpec instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceSpec object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceSpecWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceSpecWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceSpec`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceSpecWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceSpec object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPoolRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceSpec) GetPoolRef() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceSpecPoolRef`

GetPoolRef returns the PoolRef field if non-nil, zero value otherwise.

### GetPoolRefOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceSpec) GetPoolRefOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceSpecPoolRef, bool)`

GetPoolRefOk returns a tuple with the PoolRef field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPoolRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceSpec) SetPoolRef(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceSpecPoolRef)`

SetPoolRef sets PoolRef field to given value.

### HasPoolRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceSpec) HasPoolRef() bool`

HasPoolRef returns a boolean if a field has been set.

### GetType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceSpec) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceSpec) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceSpec) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceSpec) HasType() bool`

HasType returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


