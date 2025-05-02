# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberMonitoring

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Project** | **string** | Project is the Google project containing UO components. | 
**VmTenant** | **string** | VMTenant identifies the VM tenant to use (accountID or accountID:projectID). | 
**VmagentSecretRef** | [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SecretReference**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SecretReference.md) |  | 
**VminsertBackendServiceName** | **string** | VMinsertBackendServiceName identifies the backend service for vminsert. | 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberMonitoring

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberMonitoring(project string, vmTenant string, vmagentSecretRef ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SecretReference, vminsertBackendServiceName string, ) *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberMonitoring`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberMonitoring instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberMonitoring object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberMonitoringWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberMonitoringWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberMonitoring`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberMonitoringWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberMonitoring object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetProject

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberMonitoring) GetProject() string`

GetProject returns the Project field if non-nil, zero value otherwise.

### GetProjectOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberMonitoring) GetProjectOk() (*string, bool)`

GetProjectOk returns a tuple with the Project field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProject

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberMonitoring) SetProject(v string)`

SetProject sets Project field to given value.


### GetVmTenant

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberMonitoring) GetVmTenant() string`

GetVmTenant returns the VmTenant field if non-nil, zero value otherwise.

### GetVmTenantOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberMonitoring) GetVmTenantOk() (*string, bool)`

GetVmTenantOk returns a tuple with the VmTenant field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVmTenant

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberMonitoring) SetVmTenant(v string)`

SetVmTenant sets VmTenant field to given value.


### GetVmagentSecretRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberMonitoring) GetVmagentSecretRef() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SecretReference`

GetVmagentSecretRef returns the VmagentSecretRef field if non-nil, zero value otherwise.

### GetVmagentSecretRefOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberMonitoring) GetVmagentSecretRefOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SecretReference, bool)`

GetVmagentSecretRefOk returns a tuple with the VmagentSecretRef field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVmagentSecretRef

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberMonitoring) SetVmagentSecretRef(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1SecretReference)`

SetVmagentSecretRef sets VmagentSecretRef field to given value.


### GetVminsertBackendServiceName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberMonitoring) GetVminsertBackendServiceName() string`

GetVminsertBackendServiceName returns the VminsertBackendServiceName field if non-nil, zero value otherwise.

### GetVminsertBackendServiceNameOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberMonitoring) GetVminsertBackendServiceNameOk() (*string, bool)`

GetVminsertBackendServiceNameOk returns a tuple with the VminsertBackendServiceName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVminsertBackendServiceName

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolMemberMonitoring) SetVminsertBackendServiceName(v string)`

SetVminsertBackendServiceName sets VminsertBackendServiceName field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


