# \DefaultAPI

All URIs are relative to */admin/kafkaconnect*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AlterConnectorOffsets**](DefaultAPI.md#AlterConnectorOffsets) | **Patch** /connectors/{connector}/offsets | 
[**CreateConnector**](DefaultAPI.md#CreateConnector) | **Post** /connectors | 
[**DestroyConnector**](DefaultAPI.md#DestroyConnector) | **Delete** /connectors/{connector} | 
[**GetConnector**](DefaultAPI.md#GetConnector) | **Get** /connectors/{connector} | 
[**GetConnectorActiveTopics**](DefaultAPI.md#GetConnectorActiveTopics) | **Get** /connectors/{connector}/topics | 
[**GetConnectorConfig**](DefaultAPI.md#GetConnectorConfig) | **Get** /connectors/{connector}/config | 
[**GetConnectorConfigDef**](DefaultAPI.md#GetConnectorConfigDef) | **Get** /connector-plugins/{pluginName}/config | Get the configuration definition for the specified pluginName
[**GetConnectorStatus**](DefaultAPI.md#GetConnectorStatus) | **Get** /connectors/{connector}/status | 
[**GetOffsets**](DefaultAPI.md#GetOffsets) | **Get** /connectors/{connector}/offsets | 
[**GetTaskConfigs**](DefaultAPI.md#GetTaskConfigs) | **Get** /connectors/{connector}/tasks | 
[**GetTaskStatus**](DefaultAPI.md#GetTaskStatus) | **Get** /connectors/{connector}/tasks/{task}/status | 
[**GetTasksConfig**](DefaultAPI.md#GetTasksConfig) | **Get** /connectors/{connector}/tasks-config | 
[**HealthCheck**](DefaultAPI.md#HealthCheck) | **Get** /health | Health check endpoint to verify worker readiness and liveness
[**ListConnectorPlugins**](DefaultAPI.md#ListConnectorPlugins) | **Get** /connector-plugins | List all connector plugins installed
[**ListConnectorPluginsCatalog**](DefaultAPI.md#ListConnectorPluginsCatalog) | **Get** /connector-plugins/catalog | List all connector catalog
[**ListConnectors**](DefaultAPI.md#ListConnectors) | **Get** /connectors | 
[**PauseConnector**](DefaultAPI.md#PauseConnector) | **Put** /connectors/{connector}/pause | 
[**PauseConnectorV2**](DefaultAPI.md#PauseConnectorV2) | **Put** /connectors/{connector}:pause | 
[**PutConnectorConfig**](DefaultAPI.md#PutConnectorConfig) | **Put** /connectors/{connector}/config | 
[**ResetConnectorActiveTopics**](DefaultAPI.md#ResetConnectorActiveTopics) | **Put** /connectors/{connector}/topics/reset | 
[**ResetConnectorActiveTopicsV2**](DefaultAPI.md#ResetConnectorActiveTopicsV2) | **Put** /connectors/{connector}/topics:reset | 
[**ResetConnectorOffsets**](DefaultAPI.md#ResetConnectorOffsets) | **Delete** /connectors/{connector}/offsets | 
[**RestartConnector**](DefaultAPI.md#RestartConnector) | **Post** /connectors/{connector}/restart | 
[**RestartConnectorV2**](DefaultAPI.md#RestartConnectorV2) | **Post** /connectors/{connector}:restart | 
[**RestartTask**](DefaultAPI.md#RestartTask) | **Post** /connectors/{connector}/tasks/{task}/restart | 
[**ResumeConnector**](DefaultAPI.md#ResumeConnector) | **Put** /connectors/{connector}/resume | 
[**ResumeConnectorV2**](DefaultAPI.md#ResumeConnectorV2) | **Put** /connectors/{connector}:resume | 
[**ServerInfo**](DefaultAPI.md#ServerInfo) | **Get** / | 
[**StopConnector**](DefaultAPI.md#StopConnector) | **Put** /connectors/{connector}/stop | 
[**StopConnectorV2**](DefaultAPI.md#StopConnectorV2) | **Put** /connectors/{connector}:stop | 
[**ValidateConfigs**](DefaultAPI.md#ValidateConfigs) | **Put** /connector-plugins/{pluginName}/config/validate | Validate the provided configuration against the configuration definition for the specified pluginName



## AlterConnectorOffsets

> Message AlterConnectorOffsets(ctx, connector).Forward(forward).ConnectorOffsets(connectorOffsets).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/streamnative/cloud-cli/pkg/kafkaconnect"
)

func main() {
	connector := "connector_example" // string | 
	forward := true // bool |  (optional)
	connectorOffsets := *openapiclient.NewConnectorOffsets() // ConnectorOffsets |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.AlterConnectorOffsets(context.Background(), connector).Forward(forward).ConnectorOffsets(connectorOffsets).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.AlterConnectorOffsets``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AlterConnectorOffsets`: Message
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.AlterConnectorOffsets`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**connector** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiAlterConnectorOffsetsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **forward** | **bool** |  | 
 **connectorOffsets** | [**ConnectorOffsets**](ConnectorOffsets.md) |  | 

### Return type

[**Message**](Message.md)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateConnector

> ConnectorInfo CreateConnector(ctx).Forward(forward).CreateConnectorRequest(createConnectorRequest).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/streamnative/cloud-cli/pkg/kafkaconnect"
)

func main() {
	forward := true // bool |  (optional)
	createConnectorRequest := *openapiclient.NewCreateConnectorRequest() // CreateConnectorRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.CreateConnector(context.Background()).Forward(forward).CreateConnectorRequest(createConnectorRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.CreateConnector``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateConnector`: ConnectorInfo
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.CreateConnector`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateConnectorRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **forward** | **bool** |  | 
 **createConnectorRequest** | [**CreateConnectorRequest**](CreateConnectorRequest.md) |  | 

### Return type

[**ConnectorInfo**](ConnectorInfo.md)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DestroyConnector

> DestroyConnector(ctx, connector).Forward(forward).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/streamnative/cloud-cli/pkg/kafkaconnect"
)

func main() {
	connector := "connector_example" // string | 
	forward := true // bool |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DefaultAPI.DestroyConnector(context.Background(), connector).Forward(forward).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.DestroyConnector``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**connector** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiDestroyConnectorRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **forward** | **bool** |  | 

### Return type

 (empty response body)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetConnector

> ConnectorInfo GetConnector(ctx, connector).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/streamnative/cloud-cli/pkg/kafkaconnect"
)

func main() {
	connector := "connector_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.GetConnector(context.Background(), connector).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.GetConnector``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetConnector`: ConnectorInfo
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.GetConnector`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**connector** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetConnectorRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ConnectorInfo**](ConnectorInfo.md)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetConnectorActiveTopics

> map[string]interface{} GetConnectorActiveTopics(ctx, connector).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/streamnative/cloud-cli/pkg/kafkaconnect"
)

func main() {
	connector := "connector_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.GetConnectorActiveTopics(context.Background(), connector).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.GetConnectorActiveTopics``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetConnectorActiveTopics`: map[string]interface{}
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.GetConnectorActiveTopics`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**connector** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetConnectorActiveTopicsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

**map[string]interface{}**

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetConnectorConfig

> map[string]interface{} GetConnectorConfig(ctx, connector).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/streamnative/cloud-cli/pkg/kafkaconnect"
)

func main() {
	connector := "connector_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.GetConnectorConfig(context.Background(), connector).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.GetConnectorConfig``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetConnectorConfig`: map[string]interface{}
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.GetConnectorConfig`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**connector** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetConnectorConfigRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

**map[string]interface{}**

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetConnectorConfigDef

> []ConfigKeyInfo GetConnectorConfigDef(ctx, pluginName).Execute()

Get the configuration definition for the specified pluginName

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/streamnative/cloud-cli/pkg/kafkaconnect"
)

func main() {
	pluginName := "pluginName_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.GetConnectorConfigDef(context.Background(), pluginName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.GetConnectorConfigDef``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetConnectorConfigDef`: []ConfigKeyInfo
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.GetConnectorConfigDef`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**pluginName** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetConnectorConfigDefRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**[]ConfigKeyInfo**](ConfigKeyInfo.md)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetConnectorStatus

> ConnectorStateInfo GetConnectorStatus(ctx, connector).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/streamnative/cloud-cli/pkg/kafkaconnect"
)

func main() {
	connector := "connector_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.GetConnectorStatus(context.Background(), connector).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.GetConnectorStatus``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetConnectorStatus`: ConnectorStateInfo
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.GetConnectorStatus`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**connector** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetConnectorStatusRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ConnectorStateInfo**](ConnectorStateInfo.md)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetOffsets

> ConnectorOffsets GetOffsets(ctx, connector).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/streamnative/cloud-cli/pkg/kafkaconnect"
)

func main() {
	connector := "connector_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.GetOffsets(context.Background(), connector).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.GetOffsets``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetOffsets`: ConnectorOffsets
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.GetOffsets`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**connector** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetOffsetsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ConnectorOffsets**](ConnectorOffsets.md)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetTaskConfigs

> []TaskInfo GetTaskConfigs(ctx, connector).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/streamnative/cloud-cli/pkg/kafkaconnect"
)

func main() {
	connector := "connector_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.GetTaskConfigs(context.Background(), connector).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.GetTaskConfigs``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetTaskConfigs`: []TaskInfo
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.GetTaskConfigs`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**connector** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetTaskConfigsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**[]TaskInfo**](TaskInfo.md)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetTaskStatus

> TaskState GetTaskStatus(ctx, connector, task).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/streamnative/cloud-cli/pkg/kafkaconnect"
)

func main() {
	connector := "connector_example" // string | 
	task := int32(56) // int32 | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.GetTaskStatus(context.Background(), connector, task).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.GetTaskStatus``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetTaskStatus`: TaskState
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.GetTaskStatus`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**connector** | **string** |  | 
**task** | **int32** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetTaskStatusRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**TaskState**](TaskState.md)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetTasksConfig

> map[string]interface{} GetTasksConfig(ctx, connector).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/streamnative/cloud-cli/pkg/kafkaconnect"
)

func main() {
	connector := "connector_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.GetTasksConfig(context.Background(), connector).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.GetTasksConfig``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetTasksConfig`: map[string]interface{}
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.GetTasksConfig`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**connector** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetTasksConfigRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

**map[string]interface{}**

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## HealthCheck

> WorkerStatus HealthCheck(ctx).Execute()

Health check endpoint to verify worker readiness and liveness

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/streamnative/cloud-cli/pkg/kafkaconnect"
)

func main() {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.HealthCheck(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.HealthCheck``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `HealthCheck`: WorkerStatus
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.HealthCheck`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiHealthCheckRequest struct via the builder pattern


### Return type

[**WorkerStatus**](WorkerStatus.md)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListConnectorPlugins

> []PluginInfo ListConnectorPlugins(ctx).ConnectorsOnly(connectorsOnly).Execute()

List all connector plugins installed

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/streamnative/cloud-cli/pkg/kafkaconnect"
)

func main() {
	connectorsOnly := true // bool | Whether to list only connectors instead of all plugins (optional) (default to true)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.ListConnectorPlugins(context.Background()).ConnectorsOnly(connectorsOnly).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ListConnectorPlugins``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListConnectorPlugins`: []PluginInfo
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ListConnectorPlugins`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListConnectorPluginsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **connectorsOnly** | **bool** | Whether to list only connectors instead of all plugins | [default to true]

### Return type

[**[]PluginInfo**](PluginInfo.md)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListConnectorPluginsCatalog

> []FunctionMeshConnectorDefinition ListConnectorPluginsCatalog(ctx).Execute()

List all connector catalog

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/streamnative/cloud-cli/pkg/kafkaconnect"
)

func main() {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.ListConnectorPluginsCatalog(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ListConnectorPluginsCatalog``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListConnectorPluginsCatalog`: []FunctionMeshConnectorDefinition
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ListConnectorPluginsCatalog`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiListConnectorPluginsCatalogRequest struct via the builder pattern


### Return type

[**[]FunctionMeshConnectorDefinition**](FunctionMeshConnectorDefinition.md)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListConnectors

> map[string]interface{} ListConnectors(ctx).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/streamnative/cloud-cli/pkg/kafkaconnect"
)

func main() {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.ListConnectors(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ListConnectors``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListConnectors`: map[string]interface{}
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ListConnectors`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiListConnectorsRequest struct via the builder pattern


### Return type

**map[string]interface{}**

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PauseConnector

> PauseConnector(ctx, connector).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/streamnative/cloud-cli/pkg/kafkaconnect"
)

func main() {
	connector := "connector_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DefaultAPI.PauseConnector(context.Background(), connector).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.PauseConnector``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**connector** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiPauseConnectorRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PauseConnectorV2

> PauseConnectorV2(ctx, connector).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/streamnative/cloud-cli/pkg/kafkaconnect"
)

func main() {
	connector := "connector_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DefaultAPI.PauseConnectorV2(context.Background(), connector).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.PauseConnectorV2``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**connector** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiPauseConnectorV2Request struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PutConnectorConfig

> ConnectorInfo PutConnectorConfig(ctx, connector).Forward(forward).RequestBody(requestBody).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/streamnative/cloud-cli/pkg/kafkaconnect"
)

func main() {
	connector := "connector_example" // string | 
	forward := true // bool |  (optional)
	requestBody := map[string]string{"key": "Inner_example"} // map[string]string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.PutConnectorConfig(context.Background(), connector).Forward(forward).RequestBody(requestBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.PutConnectorConfig``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PutConnectorConfig`: ConnectorInfo
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.PutConnectorConfig`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**connector** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiPutConnectorConfigRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **forward** | **bool** |  | 
 **requestBody** | **map[string]string** |  | 

### Return type

[**ConnectorInfo**](ConnectorInfo.md)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ResetConnectorActiveTopics

> ResetConnectorActiveTopics(ctx, connector).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/streamnative/cloud-cli/pkg/kafkaconnect"
)

func main() {
	connector := "connector_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DefaultAPI.ResetConnectorActiveTopics(context.Background(), connector).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ResetConnectorActiveTopics``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**connector** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiResetConnectorActiveTopicsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ResetConnectorActiveTopicsV2

> ResetConnectorActiveTopicsV2(ctx, connector).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/streamnative/cloud-cli/pkg/kafkaconnect"
)

func main() {
	connector := "connector_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DefaultAPI.ResetConnectorActiveTopicsV2(context.Background(), connector).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ResetConnectorActiveTopicsV2``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**connector** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiResetConnectorActiveTopicsV2Request struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ResetConnectorOffsets

> Message ResetConnectorOffsets(ctx, connector).Forward(forward).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/streamnative/cloud-cli/pkg/kafkaconnect"
)

func main() {
	connector := "connector_example" // string | 
	forward := true // bool |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.ResetConnectorOffsets(context.Background(), connector).Forward(forward).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ResetConnectorOffsets``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ResetConnectorOffsets`: Message
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ResetConnectorOffsets`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**connector** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiResetConnectorOffsetsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **forward** | **bool** |  | 

### Return type

[**Message**](Message.md)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RestartConnector

> ConnectorStateInfo RestartConnector(ctx, connector).IncludeTasks(includeTasks).OnlyFailed(onlyFailed).Forward(forward).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/streamnative/cloud-cli/pkg/kafkaconnect"
)

func main() {
	connector := "connector_example" // string | 
	includeTasks := true // bool |  (optional) (default to false)
	onlyFailed := true // bool |  (optional) (default to false)
	forward := true // bool |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.RestartConnector(context.Background(), connector).IncludeTasks(includeTasks).OnlyFailed(onlyFailed).Forward(forward).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.RestartConnector``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RestartConnector`: ConnectorStateInfo
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.RestartConnector`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**connector** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiRestartConnectorRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **includeTasks** | **bool** |  | [default to false]
 **onlyFailed** | **bool** |  | [default to false]
 **forward** | **bool** |  | 

### Return type

[**ConnectorStateInfo**](ConnectorStateInfo.md)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RestartConnectorV2

> ConnectorStateInfo RestartConnectorV2(ctx, connector).IncludeTasks(includeTasks).OnlyFailed(onlyFailed).Forward(forward).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/streamnative/cloud-cli/pkg/kafkaconnect"
)

func main() {
	connector := "connector_example" // string | 
	includeTasks := true // bool |  (optional) (default to false)
	onlyFailed := true // bool |  (optional) (default to false)
	forward := true // bool |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.RestartConnectorV2(context.Background(), connector).IncludeTasks(includeTasks).OnlyFailed(onlyFailed).Forward(forward).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.RestartConnectorV2``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RestartConnectorV2`: ConnectorStateInfo
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.RestartConnectorV2`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**connector** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiRestartConnectorV2Request struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **includeTasks** | **bool** |  | [default to false]
 **onlyFailed** | **bool** |  | [default to false]
 **forward** | **bool** |  | 

### Return type

[**ConnectorStateInfo**](ConnectorStateInfo.md)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RestartTask

> RestartTask(ctx, connector, task).Forward(forward).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/streamnative/cloud-cli/pkg/kafkaconnect"
)

func main() {
	connector := "connector_example" // string | 
	task := int32(56) // int32 | 
	forward := true // bool |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DefaultAPI.RestartTask(context.Background(), connector, task).Forward(forward).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.RestartTask``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**connector** | **string** |  | 
**task** | **int32** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiRestartTaskRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **forward** | **bool** |  | 

### Return type

 (empty response body)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ResumeConnector

> ResumeConnector(ctx, connector).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/streamnative/cloud-cli/pkg/kafkaconnect"
)

func main() {
	connector := "connector_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DefaultAPI.ResumeConnector(context.Background(), connector).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ResumeConnector``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**connector** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiResumeConnectorRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ResumeConnectorV2

> ResumeConnectorV2(ctx, connector).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/streamnative/cloud-cli/pkg/kafkaconnect"
)

func main() {
	connector := "connector_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DefaultAPI.ResumeConnectorV2(context.Background(), connector).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ResumeConnectorV2``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**connector** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiResumeConnectorV2Request struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ServerInfo

> ServerInfo ServerInfo(ctx).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/streamnative/cloud-cli/pkg/kafkaconnect"
)

func main() {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.ServerInfo(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ServerInfo``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ServerInfo`: ServerInfo
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ServerInfo`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiServerInfoRequest struct via the builder pattern


### Return type

[**ServerInfo**](ServerInfo.md)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## StopConnector

> StopConnector(ctx, connector).Forward(forward).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/streamnative/cloud-cli/pkg/kafkaconnect"
)

func main() {
	connector := "connector_example" // string | 
	forward := true // bool |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DefaultAPI.StopConnector(context.Background(), connector).Forward(forward).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.StopConnector``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**connector** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiStopConnectorRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **forward** | **bool** |  | 

### Return type

 (empty response body)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## StopConnectorV2

> StopConnectorV2(ctx, connector).Forward(forward).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/streamnative/cloud-cli/pkg/kafkaconnect"
)

func main() {
	connector := "connector_example" // string | 
	forward := true // bool |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DefaultAPI.StopConnectorV2(context.Background(), connector).Forward(forward).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.StopConnectorV2``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**connector** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiStopConnectorV2Request struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **forward** | **bool** |  | 

### Return type

 (empty response body)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ValidateConfigs

> ValidateConfigs(ctx, pluginName).RequestBody(requestBody).Execute()

Validate the provided configuration against the configuration definition for the specified pluginName

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/streamnative/cloud-cli/pkg/kafkaconnect"
)

func main() {
	pluginName := "pluginName_example" // string | 
	requestBody := map[string]string{"key": "Inner_example"} // map[string]string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DefaultAPI.ValidateConfigs(context.Background(), pluginName).RequestBody(requestBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ValidateConfigs``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**pluginName** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiValidateConfigsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **requestBody** | **map[string]string** |  | 

### Return type

 (empty response body)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

