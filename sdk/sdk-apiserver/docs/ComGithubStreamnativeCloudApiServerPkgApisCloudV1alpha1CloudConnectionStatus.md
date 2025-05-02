# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudConnectionStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AvailableLocations** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RegionInfo**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RegionInfo.md) |  | [optional] 
**AwsPolicyVersion** | Pointer to **string** |  | [optional] 
**Conditions** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition.md) | Conditions is an array of current observed conditions. | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudConnectionStatus

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudConnectionStatus() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudConnectionStatus`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudConnectionStatus instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudConnectionStatus object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudConnectionStatusWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudConnectionStatusWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudConnectionStatus`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudConnectionStatusWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudConnectionStatus object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAvailableLocations

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudConnectionStatus) GetAvailableLocations() []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RegionInfo`

GetAvailableLocations returns the AvailableLocations field if non-nil, zero value otherwise.

### GetAvailableLocationsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudConnectionStatus) GetAvailableLocationsOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RegionInfo, bool)`

GetAvailableLocationsOk returns a tuple with the AvailableLocations field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAvailableLocations

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudConnectionStatus) SetAvailableLocations(v []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1RegionInfo)`

SetAvailableLocations sets AvailableLocations field to given value.

### HasAvailableLocations

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudConnectionStatus) HasAvailableLocations() bool`

HasAvailableLocations returns a boolean if a field has been set.

### GetAwsPolicyVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudConnectionStatus) GetAwsPolicyVersion() string`

GetAwsPolicyVersion returns the AwsPolicyVersion field if non-nil, zero value otherwise.

### GetAwsPolicyVersionOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudConnectionStatus) GetAwsPolicyVersionOk() (*string, bool)`

GetAwsPolicyVersionOk returns a tuple with the AwsPolicyVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAwsPolicyVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudConnectionStatus) SetAwsPolicyVersion(v string)`

SetAwsPolicyVersion sets AwsPolicyVersion field to given value.

### HasAwsPolicyVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudConnectionStatus) HasAwsPolicyVersion() bool`

HasAwsPolicyVersion returns a boolean if a field has been set.

### GetConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudConnectionStatus) GetConditions() []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition`

GetConditions returns the Conditions field if non-nil, zero value otherwise.

### GetConditionsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudConnectionStatus) GetConditionsOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition, bool)`

GetConditionsOk returns a tuple with the Conditions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudConnectionStatus) SetConditions(v []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Condition)`

SetConditions sets Conditions field to given value.

### HasConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1CloudConnectionStatus) HasConditions() bool`

HasConditions returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


