# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Network

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Cidr** | Pointer to **string** | CIDR determines the CIDR of the VPC to create if specified | [optional] 
**Id** | Pointer to **string** | ID is the id or the name of an existing VPC when specified. It&#39;s vpc id in AWS, vpc network name in GCP and vnet name in Azure | [optional] 
**SubnetCIDR** | Pointer to **string** | SubnetCIDR determines the CIDR of the subnet to create if specified required for Azure | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Network

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Network() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Network`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Network instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Network object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1NetworkWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1NetworkWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Network`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1NetworkWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Network object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCidr

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Network) GetCidr() string`

GetCidr returns the Cidr field if non-nil, zero value otherwise.

### GetCidrOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Network) GetCidrOk() (*string, bool)`

GetCidrOk returns a tuple with the Cidr field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCidr

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Network) SetCidr(v string)`

SetCidr sets Cidr field to given value.

### HasCidr

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Network) HasCidr() bool`

HasCidr returns a boolean if a field has been set.

### GetId

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Network) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Network) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Network) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Network) HasId() bool`

HasId returns a boolean if a field has been set.

### GetSubnetCIDR

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Network) GetSubnetCIDR() string`

GetSubnetCIDR returns the SubnetCIDR field if non-nil, zero value otherwise.

### GetSubnetCIDROk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Network) GetSubnetCIDROk() (*string, bool)`

GetSubnetCIDROk returns a tuple with the SubnetCIDR field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSubnetCIDR

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Network) SetSubnetCIDR(v string)`

SetSubnetCIDR sets SubnetCIDR field to given value.

### HasSubnetCIDR

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Network) HasSubnetCIDR() bool`

HasSubnetCIDR returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


