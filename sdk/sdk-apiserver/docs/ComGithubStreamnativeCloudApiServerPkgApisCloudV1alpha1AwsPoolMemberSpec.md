# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AwsPoolMemberSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AccessKeyID** | Pointer to **string** | AWS Access key ID | [optional] 
**AdminRoleARN** | Pointer to **string** | AdminRoleARN is the admin role to assume (or empty to use the manager&#39;s identity). | [optional] 
**ClusterName** | **string** | ClusterName is the EKS cluster name. | 
**PermissionBoundaryARN** | Pointer to **string** | PermissionBoundaryARN refers to the permission boundary to assign to IAM roles. | [optional] 
**Region** | **string** | Region is the AWS region of the cluster. | 
**SecretAccessKey** | Pointer to **string** | AWS Secret Access Key | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AwsPoolMemberSpec

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AwsPoolMemberSpec(clusterName string, region string, ) *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AwsPoolMemberSpec`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AwsPoolMemberSpec instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AwsPoolMemberSpec object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AwsPoolMemberSpecWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AwsPoolMemberSpecWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AwsPoolMemberSpec`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AwsPoolMemberSpecWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AwsPoolMemberSpec object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAccessKeyID

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AwsPoolMemberSpec) GetAccessKeyID() string`

GetAccessKeyID returns the AccessKeyID field if non-nil, zero value otherwise.

### GetAccessKeyIDOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AwsPoolMemberSpec) GetAccessKeyIDOk() (*string, bool)`

GetAccessKeyIDOk returns a tuple with the AccessKeyID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccessKeyID

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AwsPoolMemberSpec) SetAccessKeyID(v string)`

SetAccessKeyID sets AccessKeyID field to given value.

### HasAccessKeyID

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AwsPoolMemberSpec) HasAccessKeyID() bool`

HasAccessKeyID returns a boolean if a field has been set.

### GetAdminRoleARN

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AwsPoolMemberSpec) GetAdminRoleARN() string`

GetAdminRoleARN returns the AdminRoleARN field if non-nil, zero value otherwise.

### GetAdminRoleARNOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AwsPoolMemberSpec) GetAdminRoleARNOk() (*string, bool)`

GetAdminRoleARNOk returns a tuple with the AdminRoleARN field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAdminRoleARN

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AwsPoolMemberSpec) SetAdminRoleARN(v string)`

SetAdminRoleARN sets AdminRoleARN field to given value.

### HasAdminRoleARN

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AwsPoolMemberSpec) HasAdminRoleARN() bool`

HasAdminRoleARN returns a boolean if a field has been set.

### GetClusterName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AwsPoolMemberSpec) GetClusterName() string`

GetClusterName returns the ClusterName field if non-nil, zero value otherwise.

### GetClusterNameOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AwsPoolMemberSpec) GetClusterNameOk() (*string, bool)`

GetClusterNameOk returns a tuple with the ClusterName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClusterName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AwsPoolMemberSpec) SetClusterName(v string)`

SetClusterName sets ClusterName field to given value.


### GetPermissionBoundaryARN

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AwsPoolMemberSpec) GetPermissionBoundaryARN() string`

GetPermissionBoundaryARN returns the PermissionBoundaryARN field if non-nil, zero value otherwise.

### GetPermissionBoundaryARNOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AwsPoolMemberSpec) GetPermissionBoundaryARNOk() (*string, bool)`

GetPermissionBoundaryARNOk returns a tuple with the PermissionBoundaryARN field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPermissionBoundaryARN

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AwsPoolMemberSpec) SetPermissionBoundaryARN(v string)`

SetPermissionBoundaryARN sets PermissionBoundaryARN field to given value.

### HasPermissionBoundaryARN

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AwsPoolMemberSpec) HasPermissionBoundaryARN() bool`

HasPermissionBoundaryARN returns a boolean if a field has been set.

### GetRegion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AwsPoolMemberSpec) GetRegion() string`

GetRegion returns the Region field if non-nil, zero value otherwise.

### GetRegionOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AwsPoolMemberSpec) GetRegionOk() (*string, bool)`

GetRegionOk returns a tuple with the Region field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRegion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AwsPoolMemberSpec) SetRegion(v string)`

SetRegion sets Region field to given value.


### GetSecretAccessKey

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AwsPoolMemberSpec) GetSecretAccessKey() string`

GetSecretAccessKey returns the SecretAccessKey field if non-nil, zero value otherwise.

### GetSecretAccessKeyOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AwsPoolMemberSpec) GetSecretAccessKeyOk() (*string, bool)`

GetSecretAccessKeyOk returns a tuple with the SecretAccessKey field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSecretAccessKey

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AwsPoolMemberSpec) SetSecretAccessKey(v string)`

SetSecretAccessKey sets SecretAccessKey field to given value.

### HasSecretAccessKey

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1AwsPoolMemberSpec) HasSecretAccessKey() bool`

HasSecretAccessKey returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


