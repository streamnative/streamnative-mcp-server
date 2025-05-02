# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ExecConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ApiVersion** | Pointer to **string** | Preferred input version of the ExecInfo. The returned ExecCredentials MUST use the same encoding version as the input. | [optional] 
**Args** | Pointer to **[]string** | Arguments to pass to the command when executing it. | [optional] 
**Command** | **string** | Command to execute. | 
**Env** | Pointer to [**[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ExecEnvVar**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ExecEnvVar.md) | Env defines additional environment variables to expose to the process. These are unioned with the host&#39;s environment, as well as variables client-go uses to pass argument to the plugin. | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ExecConfig

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ExecConfig(command string, ) *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ExecConfig`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ExecConfig instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ExecConfig object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ExecConfigWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ExecConfigWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ExecConfig`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ExecConfigWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ExecConfig object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ExecConfig) GetApiVersion() string`

GetApiVersion returns the ApiVersion field if non-nil, zero value otherwise.

### GetApiVersionOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ExecConfig) GetApiVersionOk() (*string, bool)`

GetApiVersionOk returns a tuple with the ApiVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ExecConfig) SetApiVersion(v string)`

SetApiVersion sets ApiVersion field to given value.

### HasApiVersion

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ExecConfig) HasApiVersion() bool`

HasApiVersion returns a boolean if a field has been set.

### GetArgs

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ExecConfig) GetArgs() []string`

GetArgs returns the Args field if non-nil, zero value otherwise.

### GetArgsOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ExecConfig) GetArgsOk() (*[]string, bool)`

GetArgsOk returns a tuple with the Args field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetArgs

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ExecConfig) SetArgs(v []string)`

SetArgs sets Args field to given value.

### HasArgs

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ExecConfig) HasArgs() bool`

HasArgs returns a boolean if a field has been set.

### GetCommand

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ExecConfig) GetCommand() string`

GetCommand returns the Command field if non-nil, zero value otherwise.

### GetCommandOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ExecConfig) GetCommandOk() (*string, bool)`

GetCommandOk returns a tuple with the Command field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCommand

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ExecConfig) SetCommand(v string)`

SetCommand sets Command field to given value.


### GetEnv

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ExecConfig) GetEnv() []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ExecEnvVar`

GetEnv returns the Env field if non-nil, zero value otherwise.

### GetEnvOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ExecConfig) GetEnvOk() (*[]ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ExecEnvVar, bool)`

GetEnvOk returns a tuple with the Env field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnv

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ExecConfig) SetEnv(v []ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ExecEnvVar)`

SetEnv sets Env field to given value.

### HasEnv

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1ExecConfig) HasEnv() bool`

HasEnv returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


