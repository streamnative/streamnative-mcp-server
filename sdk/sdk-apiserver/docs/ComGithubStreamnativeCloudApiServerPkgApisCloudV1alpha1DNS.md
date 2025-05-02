# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1DNS

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** | ID is the identifier of a DNS Zone. It can be AWS zone id, GCP zone name, and Azure zone id If ID is specified, an existing zone will be used. Otherwise, a new DNS zone will be created and managed by SN. | [optional] 
**Name** | Pointer to **string** | Name is the dns domain name. It can be AWS zone name, GCP dns name and Azure zone name. | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1DNS

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1DNS() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1DNS`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1DNS instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1DNS object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1DNSWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1DNSWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1DNS`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1DNSWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1DNS object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1DNS) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1DNS) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1DNS) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1DNS) HasId() bool`

HasId returns a boolean if a field has been set.

### GetName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1DNS) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1DNS) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1DNS) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1DNS) HasName() bool`

HasName returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


