# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionList

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ApiVersion** | Pointer to **string** | APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources | [optional] 
**Items** | [**[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscription**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscription.md) |  | 
**Kind** | Pointer to **string** | Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds | [optional] 
**Metadata** | Pointer to [**V1ListMeta**](V1ListMeta.md) |  | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionList

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionList(items []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscription, ) *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionList`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionList instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionList object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionListWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionListWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionList`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionListWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionList object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionList) GetApiVersion() string`

GetApiVersion returns the ApiVersion field if non-nil, zero value otherwise.

### GetApiVersionOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionList) GetApiVersionOk() (*string, bool)`

GetApiVersionOk returns a tuple with the ApiVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionList) SetApiVersion(v string)`

SetApiVersion sets ApiVersion field to given value.

### HasApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionList) HasApiVersion() bool`

HasApiVersion returns a boolean if a field has been set.

### GetItems

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionList) GetItems() []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscription`

GetItems returns the Items field if non-nil, zero value otherwise.

### GetItemsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionList) GetItemsOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscription, bool)`

GetItemsOk returns a tuple with the Items field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetItems

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionList) SetItems(v []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscription)`

SetItems sets Items field to given value.


### GetKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionList) GetKind() string`

GetKind returns the Kind field if non-nil, zero value otherwise.

### GetKindOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionList) GetKindOk() (*string, bool)`

GetKindOk returns a tuple with the Kind field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionList) SetKind(v string)`

SetKind sets Kind field to given value.

### HasKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionList) HasKind() bool`

HasKind returns a boolean if a field has been set.

### GetMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionList) GetMetadata() V1ListMeta`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionList) GetMetadataOk() (*V1ListMeta, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionList) SetMetadata(v V1ListMeta)`

SetMetadata sets Metadata field to given value.

### HasMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha2AWSSubscriptionList) HasMetadata() bool`

HasMetadata returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


