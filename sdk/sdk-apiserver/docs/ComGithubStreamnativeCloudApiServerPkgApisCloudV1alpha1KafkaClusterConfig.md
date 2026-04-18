# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Custom** | Pointer to **map[string]string** | Custom accepts custom Kafka broker configurations These configurations will be passed directly to Kafka brokers Example: {\&quot;LOG_RETENTION_HOURS\&quot;: \&quot;168\&quot;, \&quot;NUM_PARTITIONS\&quot;: \&quot;3\&quot;} | [optional] 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterConfig

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterConfig() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterConfig`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterConfig instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterConfig object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterConfigWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterConfigWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterConfig`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterConfigWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterConfig object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCustom

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterConfig) GetCustom() map[string]string`

GetCustom returns the Custom field if non-nil, zero value otherwise.

### GetCustomOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterConfig) GetCustomOk() (*map[string]string, bool)`

GetCustomOk returns a tuple with the Custom field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCustom

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterConfig) SetCustom(v map[string]string)`

SetCustom sets Custom field to given value.

### HasCustom

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaClusterConfig) HasCustom() bool`

HasCustom returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


