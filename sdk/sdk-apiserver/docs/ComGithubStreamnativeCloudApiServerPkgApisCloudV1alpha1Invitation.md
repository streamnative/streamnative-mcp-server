# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Invitation

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Decision** | **string** | Decision indicates the user&#39;s response to the invitation | 
**Expiration** | **time.Time** | Expiration indicates when the invitation expires | 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Invitation

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Invitation(decision string, expiration time.Time, ) *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Invitation`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Invitation instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Invitation object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InvitationWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InvitationWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Invitation`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InvitationWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Invitation object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDecision

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Invitation) GetDecision() string`

GetDecision returns the Decision field if non-nil, zero value otherwise.

### GetDecisionOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Invitation) GetDecisionOk() (*string, bool)`

GetDecisionOk returns a tuple with the Decision field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDecision

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Invitation) SetDecision(v string)`

SetDecision sets Decision field to given value.


### GetExpiration

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Invitation) GetExpiration() time.Time`

GetExpiration returns the Expiration field if non-nil, zero value otherwise.

### GetExpirationOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Invitation) GetExpirationOk() (*time.Time, bool)`

GetExpirationOk returns a tuple with the Expiration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpiration

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Invitation) SetExpiration(v time.Time)`

SetExpiration sets Expiration field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


