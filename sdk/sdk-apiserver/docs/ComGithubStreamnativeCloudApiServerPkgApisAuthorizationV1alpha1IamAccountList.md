# ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountList

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ApiVersion** | Pointer to **string** | APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources | [optional] 
**Items** | [**[]ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount**](ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount.md) |  | 
**Kind** | Pointer to **string** | Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds | [optional] 
**Metadata** | Pointer to [**V1ListMeta**](V1ListMeta.md) |  | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountList

`func NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountList(items []ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount, ) *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountList`

NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountList instantiates a new ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountList object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountListWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountListWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountList`

NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountListWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountList object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountList) GetApiVersion() string`

GetApiVersion returns the ApiVersion field if non-nil, zero value otherwise.

### GetApiVersionOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountList) GetApiVersionOk() (*string, bool)`

GetApiVersionOk returns a tuple with the ApiVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountList) SetApiVersion(v string)`

SetApiVersion sets ApiVersion field to given value.

### HasApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountList) HasApiVersion() bool`

HasApiVersion returns a boolean if a field has been set.

### GetItems

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountList) GetItems() []ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount`

GetItems returns the Items field if non-nil, zero value otherwise.

### GetItemsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountList) GetItemsOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount, bool)`

GetItemsOk returns a tuple with the Items field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetItems

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountList) SetItems(v []ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount)`

SetItems sets Items field to given value.


### GetKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountList) GetKind() string`

GetKind returns the Kind field if non-nil, zero value otherwise.

### GetKindOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountList) GetKindOk() (*string, bool)`

GetKindOk returns a tuple with the Kind field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountList) SetKind(v string)`

SetKind sets Kind field to given value.

### HasKind

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountList) HasKind() bool`

HasKind returns a boolean if a field has been set.

### GetMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountList) GetMetadata() V1ListMeta`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountList) GetMetadataOk() (*V1ListMeta, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountList) SetMetadata(v V1ListMeta)`

SetMetadata sets Metadata field to given value.

### HasMetadata

`func (o *ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountList) HasMetadata() bool`

HasMetadata returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


