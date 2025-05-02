# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SelfRegistrationSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Aws** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SelfRegistrationAws**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SelfRegistrationAws.md) |  | [optional] 
**DisplayName** | **string** |  | 
**Metadata** | Pointer to **map[string]string** |  | [optional] 
**Stripe** | Pointer to **map[string]interface{}** |  | [optional] 
**Suger** | Pointer to [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SelfRegistrationSuger**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SelfRegistrationSuger.md) |  | [optional] 
**Type** | **string** |  | 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SelfRegistrationSpec

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SelfRegistrationSpec(displayName string, type_ string, ) *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SelfRegistrationSpec`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SelfRegistrationSpec instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SelfRegistrationSpec object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SelfRegistrationSpecWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SelfRegistrationSpecWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SelfRegistrationSpec`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SelfRegistrationSpecWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SelfRegistrationSpec object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAws

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SelfRegistrationSpec) GetAws() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SelfRegistrationAws`

GetAws returns the Aws field if non-nil, zero value otherwise.

### GetAwsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SelfRegistrationSpec) GetAwsOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SelfRegistrationAws, bool)`

GetAwsOk returns a tuple with the Aws field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAws

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SelfRegistrationSpec) SetAws(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SelfRegistrationAws)`

SetAws sets Aws field to given value.

### HasAws

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SelfRegistrationSpec) HasAws() bool`

HasAws returns a boolean if a field has been set.

### GetDisplayName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SelfRegistrationSpec) GetDisplayName() string`

GetDisplayName returns the DisplayName field if non-nil, zero value otherwise.

### GetDisplayNameOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SelfRegistrationSpec) GetDisplayNameOk() (*string, bool)`

GetDisplayNameOk returns a tuple with the DisplayName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDisplayName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SelfRegistrationSpec) SetDisplayName(v string)`

SetDisplayName sets DisplayName field to given value.


### GetMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SelfRegistrationSpec) GetMetadata() map[string]string`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SelfRegistrationSpec) GetMetadataOk() (*map[string]string, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SelfRegistrationSpec) SetMetadata(v map[string]string)`

SetMetadata sets Metadata field to given value.

### HasMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SelfRegistrationSpec) HasMetadata() bool`

HasMetadata returns a boolean if a field has been set.

### GetStripe

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SelfRegistrationSpec) GetStripe() map[string]interface{}`

GetStripe returns the Stripe field if non-nil, zero value otherwise.

### GetStripeOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SelfRegistrationSpec) GetStripeOk() (*map[string]interface{}, bool)`

GetStripeOk returns a tuple with the Stripe field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStripe

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SelfRegistrationSpec) SetStripe(v map[string]interface{})`

SetStripe sets Stripe field to given value.

### HasStripe

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SelfRegistrationSpec) HasStripe() bool`

HasStripe returns a boolean if a field has been set.

### GetSuger

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SelfRegistrationSpec) GetSuger() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SelfRegistrationSuger`

GetSuger returns the Suger field if non-nil, zero value otherwise.

### GetSugerOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SelfRegistrationSpec) GetSugerOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SelfRegistrationSuger, bool)`

GetSugerOk returns a tuple with the Suger field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSuger

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SelfRegistrationSpec) SetSuger(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SelfRegistrationSuger)`

SetSuger sets Suger field to given value.

### HasSuger

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SelfRegistrationSpec) HasSuger() bool`

HasSuger returns a boolean if a field has been set.

### GetType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SelfRegistrationSpec) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SelfRegistrationSpec) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SelfRegistrationSpec) SetType(v string)`

SetType sets Type field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


