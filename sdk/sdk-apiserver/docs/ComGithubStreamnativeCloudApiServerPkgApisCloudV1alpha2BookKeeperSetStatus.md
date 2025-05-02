# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Conditions** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2Condition**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2Condition.md) | Conditions is an array of current observed conditions. | [optional] 
**MetadataServiceUri** | Pointer to **string** | MetadataServiceUri exposes the URI used for loading corresponding metadata driver and resolving its metadata service location | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetStatus

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetStatus() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetStatus`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetStatus instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetStatus object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetStatusWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetStatusWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetStatus`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetStatusWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetStatus object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetStatus) GetConditions() []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2Condition`

GetConditions returns the Conditions field if non-nil, zero value otherwise.

### GetConditionsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetStatus) GetConditionsOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2Condition, bool)`

GetConditionsOk returns a tuple with the Conditions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetStatus) SetConditions(v []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2Condition)`

SetConditions sets Conditions field to given value.

### HasConditions

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetStatus) HasConditions() bool`

HasConditions returns a boolean if a field has been set.

### GetMetadataServiceUri

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetStatus) GetMetadataServiceUri() string`

GetMetadataServiceUri returns the MetadataServiceUri field if non-nil, zero value otherwise.

### GetMetadataServiceUriOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetStatus) GetMetadataServiceUriOk() (*string, bool)`

GetMetadataServiceUriOk returns a tuple with the MetadataServiceUri field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadataServiceUri

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetStatus) SetMetadataServiceUri(v string)`

SetMetadataServiceUri sets MetadataServiceUri field to given value.

### HasMetadataServiceUri

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2BookKeeperSetStatus) HasMetadataServiceUri() bool`

HasMetadataServiceUri returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


