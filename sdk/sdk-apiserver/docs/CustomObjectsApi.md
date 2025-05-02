# \CustomObjectsApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateClusterCustomObject**](CustomObjectsApi.md#CreateClusterCustomObject) | **Post** /apis/{group}/{version}/{plural} | 
[**CreateNamespacedCustomObject**](CustomObjectsApi.md#CreateNamespacedCustomObject) | **Post** /apis/{group}/{version}/namespaces/{namespace}/{plural} | 
[**DeleteClusterCustomObject**](CustomObjectsApi.md#DeleteClusterCustomObject) | **Delete** /apis/{group}/{version}/{plural}/{name} | 
[**DeleteCollectionClusterCustomObject**](CustomObjectsApi.md#DeleteCollectionClusterCustomObject) | **Delete** /apis/{group}/{version}/{plural} | 
[**DeleteCollectionNamespacedCustomObject**](CustomObjectsApi.md#DeleteCollectionNamespacedCustomObject) | **Delete** /apis/{group}/{version}/namespaces/{namespace}/{plural} | 
[**DeleteNamespacedCustomObject**](CustomObjectsApi.md#DeleteNamespacedCustomObject) | **Delete** /apis/{group}/{version}/namespaces/{namespace}/{plural}/{name} | 
[**GetClusterCustomObject**](CustomObjectsApi.md#GetClusterCustomObject) | **Get** /apis/{group}/{version}/{plural}/{name} | 
[**GetClusterCustomObjectScale**](CustomObjectsApi.md#GetClusterCustomObjectScale) | **Get** /apis/{group}/{version}/{plural}/{name}/scale | 
[**GetClusterCustomObjectStatus**](CustomObjectsApi.md#GetClusterCustomObjectStatus) | **Get** /apis/{group}/{version}/{plural}/{name}/status | 
[**GetCustomObjectsAPIResources**](CustomObjectsApi.md#GetCustomObjectsAPIResources) | **Get** /apis/{group}/{version} | 
[**GetNamespacedCustomObject**](CustomObjectsApi.md#GetNamespacedCustomObject) | **Get** /apis/{group}/{version}/namespaces/{namespace}/{plural}/{name} | 
[**GetNamespacedCustomObjectScale**](CustomObjectsApi.md#GetNamespacedCustomObjectScale) | **Get** /apis/{group}/{version}/namespaces/{namespace}/{plural}/{name}/scale | 
[**GetNamespacedCustomObjectStatus**](CustomObjectsApi.md#GetNamespacedCustomObjectStatus) | **Get** /apis/{group}/{version}/namespaces/{namespace}/{plural}/{name}/status | 
[**ListClusterCustomObject**](CustomObjectsApi.md#ListClusterCustomObject) | **Get** /apis/{group}/{version}/{plural} | 
[**ListCustomObjectForAllNamespaces**](CustomObjectsApi.md#ListCustomObjectForAllNamespaces) | **Get** /apis/{group}/{version}/{plural}#‎ | 
[**ListNamespacedCustomObject**](CustomObjectsApi.md#ListNamespacedCustomObject) | **Get** /apis/{group}/{version}/namespaces/{namespace}/{plural} | 
[**PatchClusterCustomObject**](CustomObjectsApi.md#PatchClusterCustomObject) | **Patch** /apis/{group}/{version}/{plural}/{name} | 
[**PatchClusterCustomObjectScale**](CustomObjectsApi.md#PatchClusterCustomObjectScale) | **Patch** /apis/{group}/{version}/{plural}/{name}/scale | 
[**PatchClusterCustomObjectStatus**](CustomObjectsApi.md#PatchClusterCustomObjectStatus) | **Patch** /apis/{group}/{version}/{plural}/{name}/status | 
[**PatchNamespacedCustomObject**](CustomObjectsApi.md#PatchNamespacedCustomObject) | **Patch** /apis/{group}/{version}/namespaces/{namespace}/{plural}/{name} | 
[**PatchNamespacedCustomObjectScale**](CustomObjectsApi.md#PatchNamespacedCustomObjectScale) | **Patch** /apis/{group}/{version}/namespaces/{namespace}/{plural}/{name}/scale | 
[**PatchNamespacedCustomObjectStatus**](CustomObjectsApi.md#PatchNamespacedCustomObjectStatus) | **Patch** /apis/{group}/{version}/namespaces/{namespace}/{plural}/{name}/status | 
[**ReplaceClusterCustomObject**](CustomObjectsApi.md#ReplaceClusterCustomObject) | **Put** /apis/{group}/{version}/{plural}/{name} | 
[**ReplaceClusterCustomObjectScale**](CustomObjectsApi.md#ReplaceClusterCustomObjectScale) | **Put** /apis/{group}/{version}/{plural}/{name}/scale | 
[**ReplaceClusterCustomObjectStatus**](CustomObjectsApi.md#ReplaceClusterCustomObjectStatus) | **Put** /apis/{group}/{version}/{plural}/{name}/status | 
[**ReplaceNamespacedCustomObject**](CustomObjectsApi.md#ReplaceNamespacedCustomObject) | **Put** /apis/{group}/{version}/namespaces/{namespace}/{plural}/{name} | 
[**ReplaceNamespacedCustomObjectScale**](CustomObjectsApi.md#ReplaceNamespacedCustomObjectScale) | **Put** /apis/{group}/{version}/namespaces/{namespace}/{plural}/{name}/scale | 
[**ReplaceNamespacedCustomObjectStatus**](CustomObjectsApi.md#ReplaceNamespacedCustomObjectStatus) | **Put** /apis/{group}/{version}/namespaces/{namespace}/{plural}/{name}/status | 



## CreateClusterCustomObject

> map[string]interface{} CreateClusterCustomObject(ctx, group, version, plural).Body(body).Pretty(pretty).DryRun(dryRun).FieldManager(fieldManager).FieldValidation(fieldValidation).Execute()





### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    group := "group_example" // string | The custom resource's group name
    version := "version_example" // string | The custom resource's version
    plural := "plural_example" // string | The custom resource's plural name. For TPRs this would be lowercase plural kind.
    body := map[string]interface{}{ ... } // map[string]interface{} | The JSON schema of the Resource to create.
    pretty := "pretty_example" // string | If 'true', then the output is pretty printed. (optional)
    dryRun := "dryRun_example" // string | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed (optional)
    fieldManager := "fieldManager_example" // string | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. This field is required for apply requests (application/apply-patch) but optional for non-apply patch types (JsonPatch, MergePatch, StrategicMergePatch). (optional)
    fieldValidation := "fieldValidation_example" // string | fieldValidation instructs the server on how to handle objects in the request (POST/PUT/PATCH) containing unknown or duplicate fields. Valid values are: - Ignore: This will ignore any unknown fields that are silently dropped from the object, and will ignore all but the last duplicate field that the decoder encounters. This is the default behavior prior to v1.23. - Warn: This will send a warning via the standard warning response header for each unknown field that is dropped from the object, and for each duplicate field that is encountered. The request will still succeed if there are no other errors, and will only persist the last of any duplicate fields. This is the default in v1.23+ - Strict: This will fail the request with a BadRequest error if any unknown fields would be dropped from the object, or if any duplicate fields are present. The error returned from the server will contain all unknown and duplicate fields encountered. (optional) (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.CustomObjectsApi.CreateClusterCustomObject(context.Background(), group, version, plural).Body(body).Pretty(pretty).DryRun(dryRun).FieldManager(fieldManager).FieldValidation(fieldValidation).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `CustomObjectsApi.CreateClusterCustomObject``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `CreateClusterCustomObject`: map[string]interface{}
    fmt.Fprintf(os.Stdout, "Response from `CustomObjectsApi.CreateClusterCustomObject`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**group** | **string** | The custom resource&#39;s group name | 
**version** | **string** | The custom resource&#39;s version | 
**plural** | **string** | The custom resource&#39;s plural name. For TPRs this would be lowercase plural kind. | 

### Other Parameters

Other parameters are passed through a pointer to a apiCreateClusterCustomObjectRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



 **body** | **map[string]interface{}** | The JSON schema of the Resource to create. | 
 **pretty** | **string** | If &#39;true&#39;, then the output is pretty printed. | 
 **dryRun** | **string** | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed | 
 **fieldManager** | **string** | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. This field is required for apply requests (application/apply-patch) but optional for non-apply patch types (JsonPatch, MergePatch, StrategicMergePatch). | 
 **fieldValidation** | **string** | fieldValidation instructs the server on how to handle objects in the request (POST/PUT/PATCH) containing unknown or duplicate fields. Valid values are: - Ignore: This will ignore any unknown fields that are silently dropped from the object, and will ignore all but the last duplicate field that the decoder encounters. This is the default behavior prior to v1.23. - Warn: This will send a warning via the standard warning response header for each unknown field that is dropped from the object, and for each duplicate field that is encountered. The request will still succeed if there are no other errors, and will only persist the last of any duplicate fields. This is the default in v1.23+ - Strict: This will fail the request with a BadRequest error if any unknown fields would be dropped from the object, or if any duplicate fields are present. The error returned from the server will contain all unknown and duplicate fields encountered. (optional) | 

### Return type

**map[string]interface{}**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateNamespacedCustomObject

> map[string]interface{} CreateNamespacedCustomObject(ctx, group, version, namespace, plural).Body(body).Pretty(pretty).DryRun(dryRun).FieldManager(fieldManager).FieldValidation(fieldValidation).Execute()





### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    group := "group_example" // string | The custom resource's group name
    version := "version_example" // string | The custom resource's version
    namespace := "namespace_example" // string | The custom resource's namespace
    plural := "plural_example" // string | The custom resource's plural name. For TPRs this would be lowercase plural kind.
    body := map[string]interface{}{ ... } // map[string]interface{} | The JSON schema of the Resource to create.
    pretty := "pretty_example" // string | If 'true', then the output is pretty printed. (optional)
    dryRun := "dryRun_example" // string | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed (optional)
    fieldManager := "fieldManager_example" // string | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. (optional)
    fieldValidation := "fieldValidation_example" // string | fieldValidation instructs the server on how to handle objects in the request (POST/PUT/PATCH) containing unknown or duplicate fields. Valid values are: - Ignore: This will ignore any unknown fields that are silently dropped from the object, and will ignore all but the last duplicate field that the decoder encounters. This is the default behavior prior to v1.23. - Warn: This will send a warning via the standard warning response header for each unknown field that is dropped from the object, and for each duplicate field that is encountered. The request will still succeed if there are no other errors, and will only persist the last of any duplicate fields. This is the default in v1.23+ - Strict: This will fail the request with a BadRequest error if any unknown fields would be dropped from the object, or if any duplicate fields are present. The error returned from the server will contain all unknown and duplicate fields encountered. (optional) (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.CustomObjectsApi.CreateNamespacedCustomObject(context.Background(), group, version, namespace, plural).Body(body).Pretty(pretty).DryRun(dryRun).FieldManager(fieldManager).FieldValidation(fieldValidation).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `CustomObjectsApi.CreateNamespacedCustomObject``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `CreateNamespacedCustomObject`: map[string]interface{}
    fmt.Fprintf(os.Stdout, "Response from `CustomObjectsApi.CreateNamespacedCustomObject`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**group** | **string** | The custom resource&#39;s group name | 
**version** | **string** | The custom resource&#39;s version | 
**namespace** | **string** | The custom resource&#39;s namespace | 
**plural** | **string** | The custom resource&#39;s plural name. For TPRs this would be lowercase plural kind. | 

### Other Parameters

Other parameters are passed through a pointer to a apiCreateNamespacedCustomObjectRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------




 **body** | **map[string]interface{}** | The JSON schema of the Resource to create. | 
 **pretty** | **string** | If &#39;true&#39;, then the output is pretty printed. | 
 **dryRun** | **string** | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed | 
 **fieldManager** | **string** | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. | 
 **fieldValidation** | **string** | fieldValidation instructs the server on how to handle objects in the request (POST/PUT/PATCH) containing unknown or duplicate fields. Valid values are: - Ignore: This will ignore any unknown fields that are silently dropped from the object, and will ignore all but the last duplicate field that the decoder encounters. This is the default behavior prior to v1.23. - Warn: This will send a warning via the standard warning response header for each unknown field that is dropped from the object, and for each duplicate field that is encountered. The request will still succeed if there are no other errors, and will only persist the last of any duplicate fields. This is the default in v1.23+ - Strict: This will fail the request with a BadRequest error if any unknown fields would be dropped from the object, or if any duplicate fields are present. The error returned from the server will contain all unknown and duplicate fields encountered. (optional) | 

### Return type

**map[string]interface{}**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteClusterCustomObject

> map[string]interface{} DeleteClusterCustomObject(ctx, group, version, plural, name).GracePeriodSeconds(gracePeriodSeconds).OrphanDependents(orphanDependents).PropagationPolicy(propagationPolicy).DryRun(dryRun).Body(body).Execute()





### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    group := "group_example" // string | the custom resource's group
    version := "version_example" // string | the custom resource's version
    plural := "plural_example" // string | the custom object's plural name. For TPRs this would be lowercase plural kind.
    name := "name_example" // string | the custom object's name
    gracePeriodSeconds := int32(56) // int32 | The duration in seconds before the object should be deleted. Value must be non-negative integer. The value zero indicates delete immediately. If this value is nil, the default grace period for the specified type will be used. Defaults to a per object value if not specified. zero means delete immediately. (optional)
    orphanDependents := true // bool | Deprecated: please use the PropagationPolicy, this field will be deprecated in 1.7. Should the dependent objects be orphaned. If true/false, the \"orphan\" finalizer will be added to/removed from the object's finalizers list. Either this field or PropagationPolicy may be set, but not both. (optional)
    propagationPolicy := "propagationPolicy_example" // string | Whether and how garbage collection will be performed. Either this field or OrphanDependents may be set, but not both. The default policy is decided by the existing finalizer set in the metadata.finalizers and the resource-specific default policy. (optional)
    dryRun := "dryRun_example" // string | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed (optional)
    body := *openapiclient.NewV1DeleteOptions() // V1DeleteOptions |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.CustomObjectsApi.DeleteClusterCustomObject(context.Background(), group, version, plural, name).GracePeriodSeconds(gracePeriodSeconds).OrphanDependents(orphanDependents).PropagationPolicy(propagationPolicy).DryRun(dryRun).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `CustomObjectsApi.DeleteClusterCustomObject``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DeleteClusterCustomObject`: map[string]interface{}
    fmt.Fprintf(os.Stdout, "Response from `CustomObjectsApi.DeleteClusterCustomObject`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**group** | **string** | the custom resource&#39;s group | 
**version** | **string** | the custom resource&#39;s version | 
**plural** | **string** | the custom object&#39;s plural name. For TPRs this would be lowercase plural kind. | 
**name** | **string** | the custom object&#39;s name | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteClusterCustomObjectRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------




 **gracePeriodSeconds** | **int32** | The duration in seconds before the object should be deleted. Value must be non-negative integer. The value zero indicates delete immediately. If this value is nil, the default grace period for the specified type will be used. Defaults to a per object value if not specified. zero means delete immediately. | 
 **orphanDependents** | **bool** | Deprecated: please use the PropagationPolicy, this field will be deprecated in 1.7. Should the dependent objects be orphaned. If true/false, the \&quot;orphan\&quot; finalizer will be added to/removed from the object&#39;s finalizers list. Either this field or PropagationPolicy may be set, but not both. | 
 **propagationPolicy** | **string** | Whether and how garbage collection will be performed. Either this field or OrphanDependents may be set, but not both. The default policy is decided by the existing finalizer set in the metadata.finalizers and the resource-specific default policy. | 
 **dryRun** | **string** | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed | 
 **body** | [**V1DeleteOptions**](V1DeleteOptions.md) |  | 

### Return type

**map[string]interface{}**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteCollectionClusterCustomObject

> map[string]interface{} DeleteCollectionClusterCustomObject(ctx, group, version, plural).Pretty(pretty).LabelSelector(labelSelector).GracePeriodSeconds(gracePeriodSeconds).OrphanDependents(orphanDependents).PropagationPolicy(propagationPolicy).DryRun(dryRun).Body(body).Execute()





### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    group := "group_example" // string | The custom resource's group name
    version := "version_example" // string | The custom resource's version
    plural := "plural_example" // string | The custom resource's plural name. For TPRs this would be lowercase plural kind.
    pretty := "pretty_example" // string | If 'true', then the output is pretty printed. (optional)
    labelSelector := "labelSelector_example" // string | A selector to restrict the list of returned objects by their labels. Defaults to everything. (optional)
    gracePeriodSeconds := int32(56) // int32 | The duration in seconds before the object should be deleted. Value must be non-negative integer. The value zero indicates delete immediately. If this value is nil, the default grace period for the specified type will be used. Defaults to a per object value if not specified. zero means delete immediately. (optional)
    orphanDependents := true // bool | Deprecated: please use the PropagationPolicy, this field will be deprecated in 1.7. Should the dependent objects be orphaned. If true/false, the \"orphan\" finalizer will be added to/removed from the object's finalizers list. Either this field or PropagationPolicy may be set, but not both. (optional)
    propagationPolicy := "propagationPolicy_example" // string | Whether and how garbage collection will be performed. Either this field or OrphanDependents may be set, but not both. The default policy is decided by the existing finalizer set in the metadata.finalizers and the resource-specific default policy. (optional)
    dryRun := "dryRun_example" // string | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed (optional)
    body := *openapiclient.NewV1DeleteOptions() // V1DeleteOptions |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.CustomObjectsApi.DeleteCollectionClusterCustomObject(context.Background(), group, version, plural).Pretty(pretty).LabelSelector(labelSelector).GracePeriodSeconds(gracePeriodSeconds).OrphanDependents(orphanDependents).PropagationPolicy(propagationPolicy).DryRun(dryRun).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `CustomObjectsApi.DeleteCollectionClusterCustomObject``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DeleteCollectionClusterCustomObject`: map[string]interface{}
    fmt.Fprintf(os.Stdout, "Response from `CustomObjectsApi.DeleteCollectionClusterCustomObject`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**group** | **string** | The custom resource&#39;s group name | 
**version** | **string** | The custom resource&#39;s version | 
**plural** | **string** | The custom resource&#39;s plural name. For TPRs this would be lowercase plural kind. | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteCollectionClusterCustomObjectRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



 **pretty** | **string** | If &#39;true&#39;, then the output is pretty printed. | 
 **labelSelector** | **string** | A selector to restrict the list of returned objects by their labels. Defaults to everything. | 
 **gracePeriodSeconds** | **int32** | The duration in seconds before the object should be deleted. Value must be non-negative integer. The value zero indicates delete immediately. If this value is nil, the default grace period for the specified type will be used. Defaults to a per object value if not specified. zero means delete immediately. | 
 **orphanDependents** | **bool** | Deprecated: please use the PropagationPolicy, this field will be deprecated in 1.7. Should the dependent objects be orphaned. If true/false, the \&quot;orphan\&quot; finalizer will be added to/removed from the object&#39;s finalizers list. Either this field or PropagationPolicy may be set, but not both. | 
 **propagationPolicy** | **string** | Whether and how garbage collection will be performed. Either this field or OrphanDependents may be set, but not both. The default policy is decided by the existing finalizer set in the metadata.finalizers and the resource-specific default policy. | 
 **dryRun** | **string** | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed | 
 **body** | [**V1DeleteOptions**](V1DeleteOptions.md) |  | 

### Return type

**map[string]interface{}**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteCollectionNamespacedCustomObject

> map[string]interface{} DeleteCollectionNamespacedCustomObject(ctx, group, version, namespace, plural).Pretty(pretty).LabelSelector(labelSelector).GracePeriodSeconds(gracePeriodSeconds).OrphanDependents(orphanDependents).PropagationPolicy(propagationPolicy).DryRun(dryRun).FieldSelector(fieldSelector).Body(body).Execute()





### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    group := "group_example" // string | The custom resource's group name
    version := "version_example" // string | The custom resource's version
    namespace := "namespace_example" // string | The custom resource's namespace
    plural := "plural_example" // string | The custom resource's plural name. For TPRs this would be lowercase plural kind.
    pretty := "pretty_example" // string | If 'true', then the output is pretty printed. (optional)
    labelSelector := "labelSelector_example" // string | A selector to restrict the list of returned objects by their labels. Defaults to everything. (optional)
    gracePeriodSeconds := int32(56) // int32 | The duration in seconds before the object should be deleted. Value must be non-negative integer. The value zero indicates delete immediately. If this value is nil, the default grace period for the specified type will be used. Defaults to a per object value if not specified. zero means delete immediately. (optional)
    orphanDependents := true // bool | Deprecated: please use the PropagationPolicy, this field will be deprecated in 1.7. Should the dependent objects be orphaned. If true/false, the \"orphan\" finalizer will be added to/removed from the object's finalizers list. Either this field or PropagationPolicy may be set, but not both. (optional)
    propagationPolicy := "propagationPolicy_example" // string | Whether and how garbage collection will be performed. Either this field or OrphanDependents may be set, but not both. The default policy is decided by the existing finalizer set in the metadata.finalizers and the resource-specific default policy. (optional)
    dryRun := "dryRun_example" // string | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed (optional)
    fieldSelector := "fieldSelector_example" // string | A selector to restrict the list of returned objects by their fields. Defaults to everything. (optional)
    body := *openapiclient.NewV1DeleteOptions() // V1DeleteOptions |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.CustomObjectsApi.DeleteCollectionNamespacedCustomObject(context.Background(), group, version, namespace, plural).Pretty(pretty).LabelSelector(labelSelector).GracePeriodSeconds(gracePeriodSeconds).OrphanDependents(orphanDependents).PropagationPolicy(propagationPolicy).DryRun(dryRun).FieldSelector(fieldSelector).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `CustomObjectsApi.DeleteCollectionNamespacedCustomObject``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DeleteCollectionNamespacedCustomObject`: map[string]interface{}
    fmt.Fprintf(os.Stdout, "Response from `CustomObjectsApi.DeleteCollectionNamespacedCustomObject`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**group** | **string** | The custom resource&#39;s group name | 
**version** | **string** | The custom resource&#39;s version | 
**namespace** | **string** | The custom resource&#39;s namespace | 
**plural** | **string** | The custom resource&#39;s plural name. For TPRs this would be lowercase plural kind. | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteCollectionNamespacedCustomObjectRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------




 **pretty** | **string** | If &#39;true&#39;, then the output is pretty printed. | 
 **labelSelector** | **string** | A selector to restrict the list of returned objects by their labels. Defaults to everything. | 
 **gracePeriodSeconds** | **int32** | The duration in seconds before the object should be deleted. Value must be non-negative integer. The value zero indicates delete immediately. If this value is nil, the default grace period for the specified type will be used. Defaults to a per object value if not specified. zero means delete immediately. | 
 **orphanDependents** | **bool** | Deprecated: please use the PropagationPolicy, this field will be deprecated in 1.7. Should the dependent objects be orphaned. If true/false, the \&quot;orphan\&quot; finalizer will be added to/removed from the object&#39;s finalizers list. Either this field or PropagationPolicy may be set, but not both. | 
 **propagationPolicy** | **string** | Whether and how garbage collection will be performed. Either this field or OrphanDependents may be set, but not both. The default policy is decided by the existing finalizer set in the metadata.finalizers and the resource-specific default policy. | 
 **dryRun** | **string** | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed | 
 **fieldSelector** | **string** | A selector to restrict the list of returned objects by their fields. Defaults to everything. | 
 **body** | [**V1DeleteOptions**](V1DeleteOptions.md) |  | 

### Return type

**map[string]interface{}**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteNamespacedCustomObject

> map[string]interface{} DeleteNamespacedCustomObject(ctx, group, version, namespace, plural, name).GracePeriodSeconds(gracePeriodSeconds).OrphanDependents(orphanDependents).PropagationPolicy(propagationPolicy).DryRun(dryRun).Body(body).Execute()





### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    group := "group_example" // string | the custom resource's group
    version := "version_example" // string | the custom resource's version
    namespace := "namespace_example" // string | The custom resource's namespace
    plural := "plural_example" // string | the custom resource's plural name. For TPRs this would be lowercase plural kind.
    name := "name_example" // string | the custom object's name
    gracePeriodSeconds := int32(56) // int32 | The duration in seconds before the object should be deleted. Value must be non-negative integer. The value zero indicates delete immediately. If this value is nil, the default grace period for the specified type will be used. Defaults to a per object value if not specified. zero means delete immediately. (optional)
    orphanDependents := true // bool | Deprecated: please use the PropagationPolicy, this field will be deprecated in 1.7. Should the dependent objects be orphaned. If true/false, the \"orphan\" finalizer will be added to/removed from the object's finalizers list. Either this field or PropagationPolicy may be set, but not both. (optional)
    propagationPolicy := "propagationPolicy_example" // string | Whether and how garbage collection will be performed. Either this field or OrphanDependents may be set, but not both. The default policy is decided by the existing finalizer set in the metadata.finalizers and the resource-specific default policy. (optional)
    dryRun := "dryRun_example" // string | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed (optional)
    body := *openapiclient.NewV1DeleteOptions() // V1DeleteOptions |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.CustomObjectsApi.DeleteNamespacedCustomObject(context.Background(), group, version, namespace, plural, name).GracePeriodSeconds(gracePeriodSeconds).OrphanDependents(orphanDependents).PropagationPolicy(propagationPolicy).DryRun(dryRun).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `CustomObjectsApi.DeleteNamespacedCustomObject``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DeleteNamespacedCustomObject`: map[string]interface{}
    fmt.Fprintf(os.Stdout, "Response from `CustomObjectsApi.DeleteNamespacedCustomObject`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**group** | **string** | the custom resource&#39;s group | 
**version** | **string** | the custom resource&#39;s version | 
**namespace** | **string** | The custom resource&#39;s namespace | 
**plural** | **string** | the custom resource&#39;s plural name. For TPRs this would be lowercase plural kind. | 
**name** | **string** | the custom object&#39;s name | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteNamespacedCustomObjectRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------





 **gracePeriodSeconds** | **int32** | The duration in seconds before the object should be deleted. Value must be non-negative integer. The value zero indicates delete immediately. If this value is nil, the default grace period for the specified type will be used. Defaults to a per object value if not specified. zero means delete immediately. | 
 **orphanDependents** | **bool** | Deprecated: please use the PropagationPolicy, this field will be deprecated in 1.7. Should the dependent objects be orphaned. If true/false, the \&quot;orphan\&quot; finalizer will be added to/removed from the object&#39;s finalizers list. Either this field or PropagationPolicy may be set, but not both. | 
 **propagationPolicy** | **string** | Whether and how garbage collection will be performed. Either this field or OrphanDependents may be set, but not both. The default policy is decided by the existing finalizer set in the metadata.finalizers and the resource-specific default policy. | 
 **dryRun** | **string** | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed | 
 **body** | [**V1DeleteOptions**](V1DeleteOptions.md) |  | 

### Return type

**map[string]interface{}**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetClusterCustomObject

> map[string]interface{} GetClusterCustomObject(ctx, group, version, plural, name).Execute()





### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    group := "group_example" // string | the custom resource's group
    version := "version_example" // string | the custom resource's version
    plural := "plural_example" // string | the custom object's plural name. For TPRs this would be lowercase plural kind.
    name := "name_example" // string | the custom object's name

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.CustomObjectsApi.GetClusterCustomObject(context.Background(), group, version, plural, name).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `CustomObjectsApi.GetClusterCustomObject``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetClusterCustomObject`: map[string]interface{}
    fmt.Fprintf(os.Stdout, "Response from `CustomObjectsApi.GetClusterCustomObject`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**group** | **string** | the custom resource&#39;s group | 
**version** | **string** | the custom resource&#39;s version | 
**plural** | **string** | the custom object&#39;s plural name. For TPRs this would be lowercase plural kind. | 
**name** | **string** | the custom object&#39;s name | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetClusterCustomObjectRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------





### Return type

**map[string]interface{}**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetClusterCustomObjectScale

> map[string]interface{} GetClusterCustomObjectScale(ctx, group, version, plural, name).Execute()





### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    group := "group_example" // string | the custom resource's group
    version := "version_example" // string | the custom resource's version
    plural := "plural_example" // string | the custom resource's plural name. For TPRs this would be lowercase plural kind.
    name := "name_example" // string | the custom object's name

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.CustomObjectsApi.GetClusterCustomObjectScale(context.Background(), group, version, plural, name).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `CustomObjectsApi.GetClusterCustomObjectScale``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetClusterCustomObjectScale`: map[string]interface{}
    fmt.Fprintf(os.Stdout, "Response from `CustomObjectsApi.GetClusterCustomObjectScale`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**group** | **string** | the custom resource&#39;s group | 
**version** | **string** | the custom resource&#39;s version | 
**plural** | **string** | the custom resource&#39;s plural name. For TPRs this would be lowercase plural kind. | 
**name** | **string** | the custom object&#39;s name | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetClusterCustomObjectScaleRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------





### Return type

**map[string]interface{}**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/yaml, application/vnd.kubernetes.protobuf

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetClusterCustomObjectStatus

> map[string]interface{} GetClusterCustomObjectStatus(ctx, group, version, plural, name).Execute()





### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    group := "group_example" // string | the custom resource's group
    version := "version_example" // string | the custom resource's version
    plural := "plural_example" // string | the custom resource's plural name. For TPRs this would be lowercase plural kind.
    name := "name_example" // string | the custom object's name

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.CustomObjectsApi.GetClusterCustomObjectStatus(context.Background(), group, version, plural, name).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `CustomObjectsApi.GetClusterCustomObjectStatus``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetClusterCustomObjectStatus`: map[string]interface{}
    fmt.Fprintf(os.Stdout, "Response from `CustomObjectsApi.GetClusterCustomObjectStatus`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**group** | **string** | the custom resource&#39;s group | 
**version** | **string** | the custom resource&#39;s version | 
**plural** | **string** | the custom resource&#39;s plural name. For TPRs this would be lowercase plural kind. | 
**name** | **string** | the custom object&#39;s name | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetClusterCustomObjectStatusRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------





### Return type

**map[string]interface{}**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/yaml, application/vnd.kubernetes.protobuf

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetCustomObjectsAPIResources

> V1APIResourceList GetCustomObjectsAPIResources(ctx, group, version).Execute()





### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    group := "group_example" // string | The custom resource's group name
    version := "version_example" // string | The custom resource's version

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.CustomObjectsApi.GetCustomObjectsAPIResources(context.Background(), group, version).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `CustomObjectsApi.GetCustomObjectsAPIResources``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetCustomObjectsAPIResources`: V1APIResourceList
    fmt.Fprintf(os.Stdout, "Response from `CustomObjectsApi.GetCustomObjectsAPIResources`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**group** | **string** | The custom resource&#39;s group name | 
**version** | **string** | The custom resource&#39;s version | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetCustomObjectsAPIResourcesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**V1APIResourceList**](V1APIResourceList.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetNamespacedCustomObject

> map[string]interface{} GetNamespacedCustomObject(ctx, group, version, namespace, plural, name).Execute()





### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    group := "group_example" // string | the custom resource's group
    version := "version_example" // string | the custom resource's version
    namespace := "namespace_example" // string | The custom resource's namespace
    plural := "plural_example" // string | the custom resource's plural name. For TPRs this would be lowercase plural kind.
    name := "name_example" // string | the custom object's name

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.CustomObjectsApi.GetNamespacedCustomObject(context.Background(), group, version, namespace, plural, name).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `CustomObjectsApi.GetNamespacedCustomObject``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetNamespacedCustomObject`: map[string]interface{}
    fmt.Fprintf(os.Stdout, "Response from `CustomObjectsApi.GetNamespacedCustomObject`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**group** | **string** | the custom resource&#39;s group | 
**version** | **string** | the custom resource&#39;s version | 
**namespace** | **string** | The custom resource&#39;s namespace | 
**plural** | **string** | the custom resource&#39;s plural name. For TPRs this would be lowercase plural kind. | 
**name** | **string** | the custom object&#39;s name | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetNamespacedCustomObjectRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------






### Return type

**map[string]interface{}**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetNamespacedCustomObjectScale

> map[string]interface{} GetNamespacedCustomObjectScale(ctx, group, version, namespace, plural, name).Execute()





### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    group := "group_example" // string | the custom resource's group
    version := "version_example" // string | the custom resource's version
    namespace := "namespace_example" // string | The custom resource's namespace
    plural := "plural_example" // string | the custom resource's plural name. For TPRs this would be lowercase plural kind.
    name := "name_example" // string | the custom object's name

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.CustomObjectsApi.GetNamespacedCustomObjectScale(context.Background(), group, version, namespace, plural, name).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `CustomObjectsApi.GetNamespacedCustomObjectScale``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetNamespacedCustomObjectScale`: map[string]interface{}
    fmt.Fprintf(os.Stdout, "Response from `CustomObjectsApi.GetNamespacedCustomObjectScale`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**group** | **string** | the custom resource&#39;s group | 
**version** | **string** | the custom resource&#39;s version | 
**namespace** | **string** | The custom resource&#39;s namespace | 
**plural** | **string** | the custom resource&#39;s plural name. For TPRs this would be lowercase plural kind. | 
**name** | **string** | the custom object&#39;s name | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetNamespacedCustomObjectScaleRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------






### Return type

**map[string]interface{}**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/yaml, application/vnd.kubernetes.protobuf

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetNamespacedCustomObjectStatus

> map[string]interface{} GetNamespacedCustomObjectStatus(ctx, group, version, namespace, plural, name).Execute()





### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    group := "group_example" // string | the custom resource's group
    version := "version_example" // string | the custom resource's version
    namespace := "namespace_example" // string | The custom resource's namespace
    plural := "plural_example" // string | the custom resource's plural name. For TPRs this would be lowercase plural kind.
    name := "name_example" // string | the custom object's name

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.CustomObjectsApi.GetNamespacedCustomObjectStatus(context.Background(), group, version, namespace, plural, name).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `CustomObjectsApi.GetNamespacedCustomObjectStatus``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetNamespacedCustomObjectStatus`: map[string]interface{}
    fmt.Fprintf(os.Stdout, "Response from `CustomObjectsApi.GetNamespacedCustomObjectStatus`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**group** | **string** | the custom resource&#39;s group | 
**version** | **string** | the custom resource&#39;s version | 
**namespace** | **string** | The custom resource&#39;s namespace | 
**plural** | **string** | the custom resource&#39;s plural name. For TPRs this would be lowercase plural kind. | 
**name** | **string** | the custom object&#39;s name | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetNamespacedCustomObjectStatusRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------






### Return type

**map[string]interface{}**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/yaml, application/vnd.kubernetes.protobuf

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListClusterCustomObject

> map[string]interface{} ListClusterCustomObject(ctx, group, version, plural).Pretty(pretty).AllowWatchBookmarks(allowWatchBookmarks).Continue_(continue_).FieldSelector(fieldSelector).LabelSelector(labelSelector).Limit(limit).ResourceVersion(resourceVersion).ResourceVersionMatch(resourceVersionMatch).TimeoutSeconds(timeoutSeconds).Watch(watch).Execute()





### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    group := "group_example" // string | The custom resource's group name
    version := "version_example" // string | The custom resource's version
    plural := "plural_example" // string | The custom resource's plural name. For TPRs this would be lowercase plural kind.
    pretty := "pretty_example" // string | If 'true', then the output is pretty printed. (optional)
    allowWatchBookmarks := true // bool | allowWatchBookmarks requests watch events with type \"BOOKMARK\". Servers that do not implement bookmarks may ignore this flag and bookmarks are sent at the server's discretion. Clients should not assume bookmarks are returned at any specific interval, nor may they assume the server will send any BOOKMARK event during a session. If this is not a watch, this field is ignored. If the feature gate WatchBookmarks is not enabled in apiserver, this field is ignored. (optional)
    continue_ := "continue__example" // string | The continue option should be set when retrieving more results from the server. Since this value is server defined, clients may only use the continue value from a previous query result with identical query parameters (except for the value of continue) and the server may reject a continue value it does not recognize. If the specified continue value is no longer valid whether due to expiration (generally five to fifteen minutes) or a configuration change on the server, the server will respond with a 410 ResourceExpired error together with a continue token. If the client needs a consistent list, it must restart their list without the continue field. Otherwise, the client may send another list request with the token received with the 410 error, the server will respond with a list starting from the next key, but from the latest snapshot, which is inconsistent from the previous list results - objects that are created, modified, or deleted after the first list request will be included in the response, as long as their keys are after the \"next key\".  This field is not supported when watch is true. Clients may start a watch from the last resourceVersion value returned by the server and not miss any modifications. (optional)
    fieldSelector := "fieldSelector_example" // string | A selector to restrict the list of returned objects by their fields. Defaults to everything. (optional)
    labelSelector := "labelSelector_example" // string | A selector to restrict the list of returned objects by their labels. Defaults to everything. (optional)
    limit := int32(56) // int32 | limit is a maximum number of responses to return for a list call. If more items exist, the server will set the `continue` field on the list metadata to a value that can be used with the same initial query to retrieve the next set of results. Setting a limit may return fewer than the requested amount of items (up to zero items) in the event all requested objects are filtered out and clients should only use the presence of the continue field to determine whether more results are available. Servers may choose not to support the limit argument and will return all of the available results. If limit is specified and the continue field is empty, clients may assume that no more results are available. This field is not supported if watch is true.  The server guarantees that the objects returned when using continue will be identical to issuing a single list call without a limit - that is, no objects created, modified, or deleted after the first request is issued will be included in any subsequent continued requests. This is sometimes referred to as a consistent snapshot, and ensures that a client that is using limit to receive smaller chunks of a very large result can ensure they see all possible objects. If objects are updated during a chunked list the version of the object that was present at the time the first list result was calculated is returned. (optional)
    resourceVersion := "resourceVersion_example" // string | When specified with a watch call, shows changes that occur after that particular version of a resource. Defaults to changes from the beginning of history. When specified for list: - if unset, then the result is returned from remote storage based on quorum-read flag; - if it's 0, then we simply return what we currently have in cache, no guarantee; - if set to non zero, then the result is at least as fresh as given rv. (optional)
    resourceVersionMatch := "resourceVersionMatch_example" // string | resourceVersionMatch determines how resourceVersion is applied to list calls. It is highly recommended that resourceVersionMatch be set for list calls where resourceVersion is set See https://kubernetes.io/docs/reference/using-api/api-concepts/#resource-versions for details.  Defaults to unset (optional)
    timeoutSeconds := int32(56) // int32 | Timeout for the list/watch call. This limits the duration of the call, regardless of any activity or inactivity. (optional)
    watch := true // bool | Watch for changes to the described resources and return them as a stream of add, update, and remove notifications. (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.CustomObjectsApi.ListClusterCustomObject(context.Background(), group, version, plural).Pretty(pretty).AllowWatchBookmarks(allowWatchBookmarks).Continue_(continue_).FieldSelector(fieldSelector).LabelSelector(labelSelector).Limit(limit).ResourceVersion(resourceVersion).ResourceVersionMatch(resourceVersionMatch).TimeoutSeconds(timeoutSeconds).Watch(watch).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `CustomObjectsApi.ListClusterCustomObject``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListClusterCustomObject`: map[string]interface{}
    fmt.Fprintf(os.Stdout, "Response from `CustomObjectsApi.ListClusterCustomObject`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**group** | **string** | The custom resource&#39;s group name | 
**version** | **string** | The custom resource&#39;s version | 
**plural** | **string** | The custom resource&#39;s plural name. For TPRs this would be lowercase plural kind. | 

### Other Parameters

Other parameters are passed through a pointer to a apiListClusterCustomObjectRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



 **pretty** | **string** | If &#39;true&#39;, then the output is pretty printed. | 
 **allowWatchBookmarks** | **bool** | allowWatchBookmarks requests watch events with type \&quot;BOOKMARK\&quot;. Servers that do not implement bookmarks may ignore this flag and bookmarks are sent at the server&#39;s discretion. Clients should not assume bookmarks are returned at any specific interval, nor may they assume the server will send any BOOKMARK event during a session. If this is not a watch, this field is ignored. If the feature gate WatchBookmarks is not enabled in apiserver, this field is ignored. | 
 **continue_** | **string** | The continue option should be set when retrieving more results from the server. Since this value is server defined, clients may only use the continue value from a previous query result with identical query parameters (except for the value of continue) and the server may reject a continue value it does not recognize. If the specified continue value is no longer valid whether due to expiration (generally five to fifteen minutes) or a configuration change on the server, the server will respond with a 410 ResourceExpired error together with a continue token. If the client needs a consistent list, it must restart their list without the continue field. Otherwise, the client may send another list request with the token received with the 410 error, the server will respond with a list starting from the next key, but from the latest snapshot, which is inconsistent from the previous list results - objects that are created, modified, or deleted after the first list request will be included in the response, as long as their keys are after the \&quot;next key\&quot;.  This field is not supported when watch is true. Clients may start a watch from the last resourceVersion value returned by the server and not miss any modifications. | 
 **fieldSelector** | **string** | A selector to restrict the list of returned objects by their fields. Defaults to everything. | 
 **labelSelector** | **string** | A selector to restrict the list of returned objects by their labels. Defaults to everything. | 
 **limit** | **int32** | limit is a maximum number of responses to return for a list call. If more items exist, the server will set the &#x60;continue&#x60; field on the list metadata to a value that can be used with the same initial query to retrieve the next set of results. Setting a limit may return fewer than the requested amount of items (up to zero items) in the event all requested objects are filtered out and clients should only use the presence of the continue field to determine whether more results are available. Servers may choose not to support the limit argument and will return all of the available results. If limit is specified and the continue field is empty, clients may assume that no more results are available. This field is not supported if watch is true.  The server guarantees that the objects returned when using continue will be identical to issuing a single list call without a limit - that is, no objects created, modified, or deleted after the first request is issued will be included in any subsequent continued requests. This is sometimes referred to as a consistent snapshot, and ensures that a client that is using limit to receive smaller chunks of a very large result can ensure they see all possible objects. If objects are updated during a chunked list the version of the object that was present at the time the first list result was calculated is returned. | 
 **resourceVersion** | **string** | When specified with a watch call, shows changes that occur after that particular version of a resource. Defaults to changes from the beginning of history. When specified for list: - if unset, then the result is returned from remote storage based on quorum-read flag; - if it&#39;s 0, then we simply return what we currently have in cache, no guarantee; - if set to non zero, then the result is at least as fresh as given rv. | 
 **resourceVersionMatch** | **string** | resourceVersionMatch determines how resourceVersion is applied to list calls. It is highly recommended that resourceVersionMatch be set for list calls where resourceVersion is set See https://kubernetes.io/docs/reference/using-api/api-concepts/#resource-versions for details.  Defaults to unset | 
 **timeoutSeconds** | **int32** | Timeout for the list/watch call. This limits the duration of the call, regardless of any activity or inactivity. | 
 **watch** | **bool** | Watch for changes to the described resources and return them as a stream of add, update, and remove notifications. | 

### Return type

**map[string]interface{}**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/json;stream=watch

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListCustomObjectForAllNamespaces

> map[string]interface{} ListCustomObjectForAllNamespaces(ctx, group, version, plural).Pretty(pretty).AllowWatchBookmarks(allowWatchBookmarks).Continue_(continue_).FieldSelector(fieldSelector).LabelSelector(labelSelector).Limit(limit).ResourceVersion(resourceVersion).ResourceVersionMatch(resourceVersionMatch).TimeoutSeconds(timeoutSeconds).Watch(watch).Execute()





### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    group := "group_example" // string | The custom resource's group name
    version := "version_example" // string | The custom resource's version
    plural := "plural_example" // string | The custom resource's plural name. For TPRs this would be lowercase plural kind.
    pretty := "pretty_example" // string | If 'true', then the output is pretty printed. (optional)
    allowWatchBookmarks := true // bool | allowWatchBookmarks requests watch events with type \"BOOKMARK\". Servers that do not implement bookmarks may ignore this flag and bookmarks are sent at the server's discretion. Clients should not assume bookmarks are returned at any specific interval, nor may they assume the server will send any BOOKMARK event during a session. If this is not a watch, this field is ignored. If the feature gate WatchBookmarks is not enabled in apiserver, this field is ignored. (optional)
    continue_ := "continue__example" // string | The continue option should be set when retrieving more results from the server. Since this value is server defined, clients may only use the continue value from a previous query result with identical query parameters (except for the value of continue) and the server may reject a continue value it does not recognize. If the specified continue value is no longer valid whether due to expiration (generally five to fifteen minutes) or a configuration change on the server, the server will respond with a 410 ResourceExpired error together with a continue token. If the client needs a consistent list, it must restart their list without the continue field. Otherwise, the client may send another list request with the token received with the 410 error, the server will respond with a list starting from the next key, but from the latest snapshot, which is inconsistent from the previous list results - objects that are created, modified, or deleted after the first list request will be included in the response, as long as their keys are after the \"next key\".  This field is not supported when watch is true. Clients may start a watch from the last resourceVersion value returned by the server and not miss any modifications. (optional)
    fieldSelector := "fieldSelector_example" // string | A selector to restrict the list of returned objects by their fields. Defaults to everything. (optional)
    labelSelector := "labelSelector_example" // string | A selector to restrict the list of returned objects by their labels. Defaults to everything. (optional)
    limit := int32(56) // int32 | limit is a maximum number of responses to return for a list call. If more items exist, the server will set the `continue` field on the list metadata to a value that can be used with the same initial query to retrieve the next set of results. Setting a limit may return fewer than the requested amount of items (up to zero items) in the event all requested objects are filtered out and clients should only use the presence of the continue field to determine whether more results are available. Servers may choose not to support the limit argument and will return all of the available results. If limit is specified and the continue field is empty, clients may assume that no more results are available. This field is not supported if watch is true.  The server guarantees that the objects returned when using continue will be identical to issuing a single list call without a limit - that is, no objects created, modified, or deleted after the first request is issued will be included in any subsequent continued requests. This is sometimes referred to as a consistent snapshot, and ensures that a client that is using limit to receive smaller chunks of a very large result can ensure they see all possible objects. If objects are updated during a chunked list the version of the object that was present at the time the first list result was calculated is returned. (optional)
    resourceVersion := "resourceVersion_example" // string | When specified with a watch call, shows changes that occur after that particular version of a resource. Defaults to changes from the beginning of history. When specified for list: - if unset, then the result is returned from remote storage based on quorum-read flag; - if it's 0, then we simply return what we currently have in cache, no guarantee; - if set to non zero, then the result is at least as fresh as given rv. (optional)
    resourceVersionMatch := "resourceVersionMatch_example" // string | resourceVersionMatch determines how resourceVersion is applied to list calls. It is highly recommended that resourceVersionMatch be set for list calls where resourceVersion is set See https://kubernetes.io/docs/reference/using-api/api-concepts/#resource-versions for details.  Defaults to unset (optional)
    timeoutSeconds := int32(56) // int32 | Timeout for the list/watch call. This limits the duration of the call, regardless of any activity or inactivity. (optional)
    watch := true // bool | Watch for changes to the described resources and return them as a stream of add, update, and remove notifications. (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.CustomObjectsApi.ListCustomObjectForAllNamespaces(context.Background(), group, version, plural).Pretty(pretty).AllowWatchBookmarks(allowWatchBookmarks).Continue_(continue_).FieldSelector(fieldSelector).LabelSelector(labelSelector).Limit(limit).ResourceVersion(resourceVersion).ResourceVersionMatch(resourceVersionMatch).TimeoutSeconds(timeoutSeconds).Watch(watch).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `CustomObjectsApi.ListCustomObjectForAllNamespaces``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListCustomObjectForAllNamespaces`: map[string]interface{}
    fmt.Fprintf(os.Stdout, "Response from `CustomObjectsApi.ListCustomObjectForAllNamespaces`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**group** | **string** | The custom resource&#39;s group name | 
**version** | **string** | The custom resource&#39;s version | 
**plural** | **string** | The custom resource&#39;s plural name. For TPRs this would be lowercase plural kind. | 

### Other Parameters

Other parameters are passed through a pointer to a apiListCustomObjectForAllNamespacesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



 **pretty** | **string** | If &#39;true&#39;, then the output is pretty printed. | 
 **allowWatchBookmarks** | **bool** | allowWatchBookmarks requests watch events with type \&quot;BOOKMARK\&quot;. Servers that do not implement bookmarks may ignore this flag and bookmarks are sent at the server&#39;s discretion. Clients should not assume bookmarks are returned at any specific interval, nor may they assume the server will send any BOOKMARK event during a session. If this is not a watch, this field is ignored. If the feature gate WatchBookmarks is not enabled in apiserver, this field is ignored. | 
 **continue_** | **string** | The continue option should be set when retrieving more results from the server. Since this value is server defined, clients may only use the continue value from a previous query result with identical query parameters (except for the value of continue) and the server may reject a continue value it does not recognize. If the specified continue value is no longer valid whether due to expiration (generally five to fifteen minutes) or a configuration change on the server, the server will respond with a 410 ResourceExpired error together with a continue token. If the client needs a consistent list, it must restart their list without the continue field. Otherwise, the client may send another list request with the token received with the 410 error, the server will respond with a list starting from the next key, but from the latest snapshot, which is inconsistent from the previous list results - objects that are created, modified, or deleted after the first list request will be included in the response, as long as their keys are after the \&quot;next key\&quot;.  This field is not supported when watch is true. Clients may start a watch from the last resourceVersion value returned by the server and not miss any modifications. | 
 **fieldSelector** | **string** | A selector to restrict the list of returned objects by their fields. Defaults to everything. | 
 **labelSelector** | **string** | A selector to restrict the list of returned objects by their labels. Defaults to everything. | 
 **limit** | **int32** | limit is a maximum number of responses to return for a list call. If more items exist, the server will set the &#x60;continue&#x60; field on the list metadata to a value that can be used with the same initial query to retrieve the next set of results. Setting a limit may return fewer than the requested amount of items (up to zero items) in the event all requested objects are filtered out and clients should only use the presence of the continue field to determine whether more results are available. Servers may choose not to support the limit argument and will return all of the available results. If limit is specified and the continue field is empty, clients may assume that no more results are available. This field is not supported if watch is true.  The server guarantees that the objects returned when using continue will be identical to issuing a single list call without a limit - that is, no objects created, modified, or deleted after the first request is issued will be included in any subsequent continued requests. This is sometimes referred to as a consistent snapshot, and ensures that a client that is using limit to receive smaller chunks of a very large result can ensure they see all possible objects. If objects are updated during a chunked list the version of the object that was present at the time the first list result was calculated is returned. | 
 **resourceVersion** | **string** | When specified with a watch call, shows changes that occur after that particular version of a resource. Defaults to changes from the beginning of history. When specified for list: - if unset, then the result is returned from remote storage based on quorum-read flag; - if it&#39;s 0, then we simply return what we currently have in cache, no guarantee; - if set to non zero, then the result is at least as fresh as given rv. | 
 **resourceVersionMatch** | **string** | resourceVersionMatch determines how resourceVersion is applied to list calls. It is highly recommended that resourceVersionMatch be set for list calls where resourceVersion is set See https://kubernetes.io/docs/reference/using-api/api-concepts/#resource-versions for details.  Defaults to unset | 
 **timeoutSeconds** | **int32** | Timeout for the list/watch call. This limits the duration of the call, regardless of any activity or inactivity. | 
 **watch** | **bool** | Watch for changes to the described resources and return them as a stream of add, update, and remove notifications. | 

### Return type

**map[string]interface{}**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/json;stream=watch

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListNamespacedCustomObject

> map[string]interface{} ListNamespacedCustomObject(ctx, group, version, namespace, plural).Pretty(pretty).AllowWatchBookmarks(allowWatchBookmarks).Continue_(continue_).FieldSelector(fieldSelector).LabelSelector(labelSelector).Limit(limit).ResourceVersion(resourceVersion).ResourceVersionMatch(resourceVersionMatch).TimeoutSeconds(timeoutSeconds).Watch(watch).Execute()





### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    group := "group_example" // string | The custom resource's group name
    version := "version_example" // string | The custom resource's version
    namespace := "namespace_example" // string | The custom resource's namespace
    plural := "plural_example" // string | The custom resource's plural name. For TPRs this would be lowercase plural kind.
    pretty := "pretty_example" // string | If 'true', then the output is pretty printed. (optional)
    allowWatchBookmarks := true // bool | allowWatchBookmarks requests watch events with type \"BOOKMARK\". Servers that do not implement bookmarks may ignore this flag and bookmarks are sent at the server's discretion. Clients should not assume bookmarks are returned at any specific interval, nor may they assume the server will send any BOOKMARK event during a session. If this is not a watch, this field is ignored. If the feature gate WatchBookmarks is not enabled in apiserver, this field is ignored. (optional)
    continue_ := "continue__example" // string | The continue option should be set when retrieving more results from the server. Since this value is server defined, clients may only use the continue value from a previous query result with identical query parameters (except for the value of continue) and the server may reject a continue value it does not recognize. If the specified continue value is no longer valid whether due to expiration (generally five to fifteen minutes) or a configuration change on the server, the server will respond with a 410 ResourceExpired error together with a continue token. If the client needs a consistent list, it must restart their list without the continue field. Otherwise, the client may send another list request with the token received with the 410 error, the server will respond with a list starting from the next key, but from the latest snapshot, which is inconsistent from the previous list results - objects that are created, modified, or deleted after the first list request will be included in the response, as long as their keys are after the \"next key\".  This field is not supported when watch is true. Clients may start a watch from the last resourceVersion value returned by the server and not miss any modifications. (optional)
    fieldSelector := "fieldSelector_example" // string | A selector to restrict the list of returned objects by their fields. Defaults to everything. (optional)
    labelSelector := "labelSelector_example" // string | A selector to restrict the list of returned objects by their labels. Defaults to everything. (optional)
    limit := int32(56) // int32 | limit is a maximum number of responses to return for a list call. If more items exist, the server will set the `continue` field on the list metadata to a value that can be used with the same initial query to retrieve the next set of results. Setting a limit may return fewer than the requested amount of items (up to zero items) in the event all requested objects are filtered out and clients should only use the presence of the continue field to determine whether more results are available. Servers may choose not to support the limit argument and will return all of the available results. If limit is specified and the continue field is empty, clients may assume that no more results are available. This field is not supported if watch is true.  The server guarantees that the objects returned when using continue will be identical to issuing a single list call without a limit - that is, no objects created, modified, or deleted after the first request is issued will be included in any subsequent continued requests. This is sometimes referred to as a consistent snapshot, and ensures that a client that is using limit to receive smaller chunks of a very large result can ensure they see all possible objects. If objects are updated during a chunked list the version of the object that was present at the time the first list result was calculated is returned. (optional)
    resourceVersion := "resourceVersion_example" // string | When specified with a watch call, shows changes that occur after that particular version of a resource. Defaults to changes from the beginning of history. When specified for list: - if unset, then the result is returned from remote storage based on quorum-read flag; - if it's 0, then we simply return what we currently have in cache, no guarantee; - if set to non zero, then the result is at least as fresh as given rv. (optional)
    resourceVersionMatch := "resourceVersionMatch_example" // string | resourceVersionMatch determines how resourceVersion is applied to list calls. It is highly recommended that resourceVersionMatch be set for list calls where resourceVersion is set See https://kubernetes.io/docs/reference/using-api/api-concepts/#resource-versions for details.  Defaults to unset (optional)
    timeoutSeconds := int32(56) // int32 | Timeout for the list/watch call. This limits the duration of the call, regardless of any activity or inactivity. (optional)
    watch := true // bool | Watch for changes to the described resources and return them as a stream of add, update, and remove notifications. (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.CustomObjectsApi.ListNamespacedCustomObject(context.Background(), group, version, namespace, plural).Pretty(pretty).AllowWatchBookmarks(allowWatchBookmarks).Continue_(continue_).FieldSelector(fieldSelector).LabelSelector(labelSelector).Limit(limit).ResourceVersion(resourceVersion).ResourceVersionMatch(resourceVersionMatch).TimeoutSeconds(timeoutSeconds).Watch(watch).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `CustomObjectsApi.ListNamespacedCustomObject``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListNamespacedCustomObject`: map[string]interface{}
    fmt.Fprintf(os.Stdout, "Response from `CustomObjectsApi.ListNamespacedCustomObject`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**group** | **string** | The custom resource&#39;s group name | 
**version** | **string** | The custom resource&#39;s version | 
**namespace** | **string** | The custom resource&#39;s namespace | 
**plural** | **string** | The custom resource&#39;s plural name. For TPRs this would be lowercase plural kind. | 

### Other Parameters

Other parameters are passed through a pointer to a apiListNamespacedCustomObjectRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------




 **pretty** | **string** | If &#39;true&#39;, then the output is pretty printed. | 
 **allowWatchBookmarks** | **bool** | allowWatchBookmarks requests watch events with type \&quot;BOOKMARK\&quot;. Servers that do not implement bookmarks may ignore this flag and bookmarks are sent at the server&#39;s discretion. Clients should not assume bookmarks are returned at any specific interval, nor may they assume the server will send any BOOKMARK event during a session. If this is not a watch, this field is ignored. If the feature gate WatchBookmarks is not enabled in apiserver, this field is ignored. | 
 **continue_** | **string** | The continue option should be set when retrieving more results from the server. Since this value is server defined, clients may only use the continue value from a previous query result with identical query parameters (except for the value of continue) and the server may reject a continue value it does not recognize. If the specified continue value is no longer valid whether due to expiration (generally five to fifteen minutes) or a configuration change on the server, the server will respond with a 410 ResourceExpired error together with a continue token. If the client needs a consistent list, it must restart their list without the continue field. Otherwise, the client may send another list request with the token received with the 410 error, the server will respond with a list starting from the next key, but from the latest snapshot, which is inconsistent from the previous list results - objects that are created, modified, or deleted after the first list request will be included in the response, as long as their keys are after the \&quot;next key\&quot;.  This field is not supported when watch is true. Clients may start a watch from the last resourceVersion value returned by the server and not miss any modifications. | 
 **fieldSelector** | **string** | A selector to restrict the list of returned objects by their fields. Defaults to everything. | 
 **labelSelector** | **string** | A selector to restrict the list of returned objects by their labels. Defaults to everything. | 
 **limit** | **int32** | limit is a maximum number of responses to return for a list call. If more items exist, the server will set the &#x60;continue&#x60; field on the list metadata to a value that can be used with the same initial query to retrieve the next set of results. Setting a limit may return fewer than the requested amount of items (up to zero items) in the event all requested objects are filtered out and clients should only use the presence of the continue field to determine whether more results are available. Servers may choose not to support the limit argument and will return all of the available results. If limit is specified and the continue field is empty, clients may assume that no more results are available. This field is not supported if watch is true.  The server guarantees that the objects returned when using continue will be identical to issuing a single list call without a limit - that is, no objects created, modified, or deleted after the first request is issued will be included in any subsequent continued requests. This is sometimes referred to as a consistent snapshot, and ensures that a client that is using limit to receive smaller chunks of a very large result can ensure they see all possible objects. If objects are updated during a chunked list the version of the object that was present at the time the first list result was calculated is returned. | 
 **resourceVersion** | **string** | When specified with a watch call, shows changes that occur after that particular version of a resource. Defaults to changes from the beginning of history. When specified for list: - if unset, then the result is returned from remote storage based on quorum-read flag; - if it&#39;s 0, then we simply return what we currently have in cache, no guarantee; - if set to non zero, then the result is at least as fresh as given rv. | 
 **resourceVersionMatch** | **string** | resourceVersionMatch determines how resourceVersion is applied to list calls. It is highly recommended that resourceVersionMatch be set for list calls where resourceVersion is set See https://kubernetes.io/docs/reference/using-api/api-concepts/#resource-versions for details.  Defaults to unset | 
 **timeoutSeconds** | **int32** | Timeout for the list/watch call. This limits the duration of the call, regardless of any activity or inactivity. | 
 **watch** | **bool** | Watch for changes to the described resources and return them as a stream of add, update, and remove notifications. | 

### Return type

**map[string]interface{}**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/json;stream=watch

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PatchClusterCustomObject

> map[string]interface{} PatchClusterCustomObject(ctx, group, version, plural, name).Body(body).DryRun(dryRun).FieldManager(fieldManager).FieldValidation(fieldValidation).Force(force).Execute()





### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    group := "group_example" // string | the custom resource's group
    version := "version_example" // string | the custom resource's version
    plural := "plural_example" // string | the custom object's plural name. For TPRs this would be lowercase plural kind.
    name := "name_example" // string | the custom object's name
    body := map[string]interface{}{ ... } // map[string]interface{} | The JSON schema of the Resource to patch.
    dryRun := "dryRun_example" // string | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed (optional)
    fieldManager := "fieldManager_example" // string | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. This field is required for apply requests (application/apply-patch) but optional for non-apply patch types (JsonPatch, MergePatch, StrategicMergePatch). (optional)
    fieldValidation := "fieldValidation_example" // string | fieldValidation instructs the server on how to handle objects in the request (POST/PUT/PATCH) containing unknown or duplicate fields. Valid values are: - Ignore: This will ignore any unknown fields that are silently dropped from the object, and will ignore all but the last duplicate field that the decoder encounters. This is the default behavior prior to v1.23. - Warn: This will send a warning via the standard warning response header for each unknown field that is dropped from the object, and for each duplicate field that is encountered. The request will still succeed if there are no other errors, and will only persist the last of any duplicate fields. This is the default in v1.23+ - Strict: This will fail the request with a BadRequest error if any unknown fields would be dropped from the object, or if any duplicate fields are present. The error returned from the server will contain all unknown and duplicate fields encountered. (optional) (optional)
    force := true // bool | Force is going to \"force\" Apply requests. It means user will re-acquire conflicting fields owned by other people. Force flag must be unset for non-apply patch requests. (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.CustomObjectsApi.PatchClusterCustomObject(context.Background(), group, version, plural, name).Body(body).DryRun(dryRun).FieldManager(fieldManager).FieldValidation(fieldValidation).Force(force).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `CustomObjectsApi.PatchClusterCustomObject``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `PatchClusterCustomObject`: map[string]interface{}
    fmt.Fprintf(os.Stdout, "Response from `CustomObjectsApi.PatchClusterCustomObject`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**group** | **string** | the custom resource&#39;s group | 
**version** | **string** | the custom resource&#39;s version | 
**plural** | **string** | the custom object&#39;s plural name. For TPRs this would be lowercase plural kind. | 
**name** | **string** | the custom object&#39;s name | 

### Other Parameters

Other parameters are passed through a pointer to a apiPatchClusterCustomObjectRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------




 **body** | **map[string]interface{}** | The JSON schema of the Resource to patch. | 
 **dryRun** | **string** | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed | 
 **fieldManager** | **string** | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. This field is required for apply requests (application/apply-patch) but optional for non-apply patch types (JsonPatch, MergePatch, StrategicMergePatch). | 
 **fieldValidation** | **string** | fieldValidation instructs the server on how to handle objects in the request (POST/PUT/PATCH) containing unknown or duplicate fields. Valid values are: - Ignore: This will ignore any unknown fields that are silently dropped from the object, and will ignore all but the last duplicate field that the decoder encounters. This is the default behavior prior to v1.23. - Warn: This will send a warning via the standard warning response header for each unknown field that is dropped from the object, and for each duplicate field that is encountered. The request will still succeed if there are no other errors, and will only persist the last of any duplicate fields. This is the default in v1.23+ - Strict: This will fail the request with a BadRequest error if any unknown fields would be dropped from the object, or if any duplicate fields are present. The error returned from the server will contain all unknown and duplicate fields encountered. (optional) | 
 **force** | **bool** | Force is going to \&quot;force\&quot; Apply requests. It means user will re-acquire conflicting fields owned by other people. Force flag must be unset for non-apply patch requests. | 

### Return type

**map[string]interface{}**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json-patch+json, application/merge-patch+json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PatchClusterCustomObjectScale

> map[string]interface{} PatchClusterCustomObjectScale(ctx, group, version, plural, name).Body(body).DryRun(dryRun).FieldManager(fieldManager).FieldValidation(fieldValidation).Force(force).Execute()





### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    group := "group_example" // string | the custom resource's group
    version := "version_example" // string | the custom resource's version
    plural := "plural_example" // string | the custom resource's plural name. For TPRs this would be lowercase plural kind.
    name := "name_example" // string | the custom object's name
    body := map[string]interface{}{ ... } // map[string]interface{} | 
    dryRun := "dryRun_example" // string | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed (optional)
    fieldManager := "fieldManager_example" // string | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. This field is required for apply requests (application/apply-patch) but optional for non-apply patch types (JsonPatch, MergePatch, StrategicMergePatch). (optional)
    fieldValidation := "fieldValidation_example" // string | fieldValidation instructs the server on how to handle objects in the request (POST/PUT/PATCH) containing unknown or duplicate fields. Valid values are: - Ignore: This will ignore any unknown fields that are silently dropped from the object, and will ignore all but the last duplicate field that the decoder encounters. This is the default behavior prior to v1.23. - Warn: This will send a warning via the standard warning response header for each unknown field that is dropped from the object, and for each duplicate field that is encountered. The request will still succeed if there are no other errors, and will only persist the last of any duplicate fields. This is the default in v1.23+ - Strict: This will fail the request with a BadRequest error if any unknown fields would be dropped from the object, or if any duplicate fields are present. The error returned from the server will contain all unknown and duplicate fields encountered. (optional) (optional)
    force := true // bool | Force is going to \"force\" Apply requests. It means user will re-acquire conflicting fields owned by other people. Force flag must be unset for non-apply patch requests. (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.CustomObjectsApi.PatchClusterCustomObjectScale(context.Background(), group, version, plural, name).Body(body).DryRun(dryRun).FieldManager(fieldManager).FieldValidation(fieldValidation).Force(force).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `CustomObjectsApi.PatchClusterCustomObjectScale``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `PatchClusterCustomObjectScale`: map[string]interface{}
    fmt.Fprintf(os.Stdout, "Response from `CustomObjectsApi.PatchClusterCustomObjectScale`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**group** | **string** | the custom resource&#39;s group | 
**version** | **string** | the custom resource&#39;s version | 
**plural** | **string** | the custom resource&#39;s plural name. For TPRs this would be lowercase plural kind. | 
**name** | **string** | the custom object&#39;s name | 

### Other Parameters

Other parameters are passed through a pointer to a apiPatchClusterCustomObjectScaleRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------




 **body** | **map[string]interface{}** |  | 
 **dryRun** | **string** | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed | 
 **fieldManager** | **string** | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. This field is required for apply requests (application/apply-patch) but optional for non-apply patch types (JsonPatch, MergePatch, StrategicMergePatch). | 
 **fieldValidation** | **string** | fieldValidation instructs the server on how to handle objects in the request (POST/PUT/PATCH) containing unknown or duplicate fields. Valid values are: - Ignore: This will ignore any unknown fields that are silently dropped from the object, and will ignore all but the last duplicate field that the decoder encounters. This is the default behavior prior to v1.23. - Warn: This will send a warning via the standard warning response header for each unknown field that is dropped from the object, and for each duplicate field that is encountered. The request will still succeed if there are no other errors, and will only persist the last of any duplicate fields. This is the default in v1.23+ - Strict: This will fail the request with a BadRequest error if any unknown fields would be dropped from the object, or if any duplicate fields are present. The error returned from the server will contain all unknown and duplicate fields encountered. (optional) | 
 **force** | **bool** | Force is going to \&quot;force\&quot; Apply requests. It means user will re-acquire conflicting fields owned by other people. Force flag must be unset for non-apply patch requests. | 

### Return type

**map[string]interface{}**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json-patch+json, application/merge-patch+json
- **Accept**: application/json, application/yaml, application/vnd.kubernetes.protobuf

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PatchClusterCustomObjectStatus

> map[string]interface{} PatchClusterCustomObjectStatus(ctx, group, version, plural, name).Body(body).DryRun(dryRun).FieldManager(fieldManager).FieldValidation(fieldValidation).Force(force).Execute()





### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    group := "group_example" // string | the custom resource's group
    version := "version_example" // string | the custom resource's version
    plural := "plural_example" // string | the custom resource's plural name. For TPRs this would be lowercase plural kind.
    name := "name_example" // string | the custom object's name
    body := map[string]interface{}{ ... } // map[string]interface{} | 
    dryRun := "dryRun_example" // string | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed (optional)
    fieldManager := "fieldManager_example" // string | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. This field is required for apply requests (application/apply-patch) but optional for non-apply patch types (JsonPatch, MergePatch, StrategicMergePatch). (optional)
    fieldValidation := "fieldValidation_example" // string | fieldValidation instructs the server on how to handle objects in the request (POST/PUT/PATCH) containing unknown or duplicate fields. Valid values are: - Ignore: This will ignore any unknown fields that are silently dropped from the object, and will ignore all but the last duplicate field that the decoder encounters. This is the default behavior prior to v1.23. - Warn: This will send a warning via the standard warning response header for each unknown field that is dropped from the object, and for each duplicate field that is encountered. The request will still succeed if there are no other errors, and will only persist the last of any duplicate fields. This is the default in v1.23+ - Strict: This will fail the request with a BadRequest error if any unknown fields would be dropped from the object, or if any duplicate fields are present. The error returned from the server will contain all unknown and duplicate fields encountered. (optional) (optional)
    force := true // bool | Force is going to \"force\" Apply requests. It means user will re-acquire conflicting fields owned by other people. Force flag must be unset for non-apply patch requests. (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.CustomObjectsApi.PatchClusterCustomObjectStatus(context.Background(), group, version, plural, name).Body(body).DryRun(dryRun).FieldManager(fieldManager).FieldValidation(fieldValidation).Force(force).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `CustomObjectsApi.PatchClusterCustomObjectStatus``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `PatchClusterCustomObjectStatus`: map[string]interface{}
    fmt.Fprintf(os.Stdout, "Response from `CustomObjectsApi.PatchClusterCustomObjectStatus`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**group** | **string** | the custom resource&#39;s group | 
**version** | **string** | the custom resource&#39;s version | 
**plural** | **string** | the custom resource&#39;s plural name. For TPRs this would be lowercase plural kind. | 
**name** | **string** | the custom object&#39;s name | 

### Other Parameters

Other parameters are passed through a pointer to a apiPatchClusterCustomObjectStatusRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------




 **body** | **map[string]interface{}** |  | 
 **dryRun** | **string** | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed | 
 **fieldManager** | **string** | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. This field is required for apply requests (application/apply-patch) but optional for non-apply patch types (JsonPatch, MergePatch, StrategicMergePatch). | 
 **fieldValidation** | **string** | fieldValidation instructs the server on how to handle objects in the request (POST/PUT/PATCH) containing unknown or duplicate fields. Valid values are: - Ignore: This will ignore any unknown fields that are silently dropped from the object, and will ignore all but the last duplicate field that the decoder encounters. This is the default behavior prior to v1.23. - Warn: This will send a warning via the standard warning response header for each unknown field that is dropped from the object, and for each duplicate field that is encountered. The request will still succeed if there are no other errors, and will only persist the last of any duplicate fields. This is the default in v1.23+ - Strict: This will fail the request with a BadRequest error if any unknown fields would be dropped from the object, or if any duplicate fields are present. The error returned from the server will contain all unknown and duplicate fields encountered. (optional) | 
 **force** | **bool** | Force is going to \&quot;force\&quot; Apply requests. It means user will re-acquire conflicting fields owned by other people. Force flag must be unset for non-apply patch requests. | 

### Return type

**map[string]interface{}**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json-patch+json, application/merge-patch+json
- **Accept**: application/json, application/yaml, application/vnd.kubernetes.protobuf

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PatchNamespacedCustomObject

> map[string]interface{} PatchNamespacedCustomObject(ctx, group, version, namespace, plural, name).Body(body).DryRun(dryRun).FieldManager(fieldManager).FieldValidation(fieldValidation).Force(force).Execute()





### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    group := "group_example" // string | the custom resource's group
    version := "version_example" // string | the custom resource's version
    namespace := "namespace_example" // string | The custom resource's namespace
    plural := "plural_example" // string | the custom resource's plural name. For TPRs this would be lowercase plural kind.
    name := "name_example" // string | the custom object's name
    body := map[string]interface{}{ ... } // map[string]interface{} | The JSON schema of the Resource to patch.
    dryRun := "dryRun_example" // string | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed (optional)
    fieldManager := "fieldManager_example" // string | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. This field is required for apply requests (application/apply-patch) but optional for non-apply patch types (JsonPatch, MergePatch, StrategicMergePatch). (optional)
    fieldValidation := "fieldValidation_example" // string | fieldValidation instructs the server on how to handle objects in the request (POST/PUT/PATCH) containing unknown or duplicate fields. Valid values are: - Ignore: This will ignore any unknown fields that are silently dropped from the object, and will ignore all but the last duplicate field that the decoder encounters. This is the default behavior prior to v1.23. - Warn: This will send a warning via the standard warning response header for each unknown field that is dropped from the object, and for each duplicate field that is encountered. The request will still succeed if there are no other errors, and will only persist the last of any duplicate fields. This is the default in v1.23+ - Strict: This will fail the request with a BadRequest error if any unknown fields would be dropped from the object, or if any duplicate fields are present. The error returned from the server will contain all unknown and duplicate fields encountered. (optional) (optional)
    force := true // bool | Force is going to \"force\" Apply requests. It means user will re-acquire conflicting fields owned by other people. Force flag must be unset for non-apply patch requests. (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.CustomObjectsApi.PatchNamespacedCustomObject(context.Background(), group, version, namespace, plural, name).Body(body).DryRun(dryRun).FieldManager(fieldManager).FieldValidation(fieldValidation).Force(force).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `CustomObjectsApi.PatchNamespacedCustomObject``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `PatchNamespacedCustomObject`: map[string]interface{}
    fmt.Fprintf(os.Stdout, "Response from `CustomObjectsApi.PatchNamespacedCustomObject`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**group** | **string** | the custom resource&#39;s group | 
**version** | **string** | the custom resource&#39;s version | 
**namespace** | **string** | The custom resource&#39;s namespace | 
**plural** | **string** | the custom resource&#39;s plural name. For TPRs this would be lowercase plural kind. | 
**name** | **string** | the custom object&#39;s name | 

### Other Parameters

Other parameters are passed through a pointer to a apiPatchNamespacedCustomObjectRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------





 **body** | **map[string]interface{}** | The JSON schema of the Resource to patch. | 
 **dryRun** | **string** | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed | 
 **fieldManager** | **string** | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. This field is required for apply requests (application/apply-patch) but optional for non-apply patch types (JsonPatch, MergePatch, StrategicMergePatch). | 
 **fieldValidation** | **string** | fieldValidation instructs the server on how to handle objects in the request (POST/PUT/PATCH) containing unknown or duplicate fields. Valid values are: - Ignore: This will ignore any unknown fields that are silently dropped from the object, and will ignore all but the last duplicate field that the decoder encounters. This is the default behavior prior to v1.23. - Warn: This will send a warning via the standard warning response header for each unknown field that is dropped from the object, and for each duplicate field that is encountered. The request will still succeed if there are no other errors, and will only persist the last of any duplicate fields. This is the default in v1.23+ - Strict: This will fail the request with a BadRequest error if any unknown fields would be dropped from the object, or if any duplicate fields are present. The error returned from the server will contain all unknown and duplicate fields encountered. (optional) | 
 **force** | **bool** | Force is going to \&quot;force\&quot; Apply requests. It means user will re-acquire conflicting fields owned by other people. Force flag must be unset for non-apply patch requests. | 

### Return type

**map[string]interface{}**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json-patch+json, application/merge-patch+json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PatchNamespacedCustomObjectScale

> map[string]interface{} PatchNamespacedCustomObjectScale(ctx, group, version, namespace, plural, name).Body(body).DryRun(dryRun).FieldManager(fieldManager).FieldValidation(fieldValidation).Force(force).Execute()





### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    group := "group_example" // string | the custom resource's group
    version := "version_example" // string | the custom resource's version
    namespace := "namespace_example" // string | The custom resource's namespace
    plural := "plural_example" // string | the custom resource's plural name. For TPRs this would be lowercase plural kind.
    name := "name_example" // string | the custom object's name
    body := map[string]interface{}{ ... } // map[string]interface{} | 
    dryRun := "dryRun_example" // string | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed (optional)
    fieldManager := "fieldManager_example" // string | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. This field is required for apply requests (application/apply-patch) but optional for non-apply patch types (JsonPatch, MergePatch, StrategicMergePatch). (optional)
    fieldValidation := "fieldValidation_example" // string | fieldValidation instructs the server on how to handle objects in the request (POST/PUT/PATCH) containing unknown or duplicate fields. Valid values are: - Ignore: This will ignore any unknown fields that are silently dropped from the object, and will ignore all but the last duplicate field that the decoder encounters. This is the default behavior prior to v1.23. - Warn: This will send a warning via the standard warning response header for each unknown field that is dropped from the object, and for each duplicate field that is encountered. The request will still succeed if there are no other errors, and will only persist the last of any duplicate fields. This is the default in v1.23+ - Strict: This will fail the request with a BadRequest error if any unknown fields would be dropped from the object, or if any duplicate fields are present. The error returned from the server will contain all unknown and duplicate fields encountered. (optional) (optional)
    force := true // bool | Force is going to \"force\" Apply requests. It means user will re-acquire conflicting fields owned by other people. Force flag must be unset for non-apply patch requests. (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.CustomObjectsApi.PatchNamespacedCustomObjectScale(context.Background(), group, version, namespace, plural, name).Body(body).DryRun(dryRun).FieldManager(fieldManager).FieldValidation(fieldValidation).Force(force).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `CustomObjectsApi.PatchNamespacedCustomObjectScale``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `PatchNamespacedCustomObjectScale`: map[string]interface{}
    fmt.Fprintf(os.Stdout, "Response from `CustomObjectsApi.PatchNamespacedCustomObjectScale`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**group** | **string** | the custom resource&#39;s group | 
**version** | **string** | the custom resource&#39;s version | 
**namespace** | **string** | The custom resource&#39;s namespace | 
**plural** | **string** | the custom resource&#39;s plural name. For TPRs this would be lowercase plural kind. | 
**name** | **string** | the custom object&#39;s name | 

### Other Parameters

Other parameters are passed through a pointer to a apiPatchNamespacedCustomObjectScaleRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------





 **body** | **map[string]interface{}** |  | 
 **dryRun** | **string** | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed | 
 **fieldManager** | **string** | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. This field is required for apply requests (application/apply-patch) but optional for non-apply patch types (JsonPatch, MergePatch, StrategicMergePatch). | 
 **fieldValidation** | **string** | fieldValidation instructs the server on how to handle objects in the request (POST/PUT/PATCH) containing unknown or duplicate fields. Valid values are: - Ignore: This will ignore any unknown fields that are silently dropped from the object, and will ignore all but the last duplicate field that the decoder encounters. This is the default behavior prior to v1.23. - Warn: This will send a warning via the standard warning response header for each unknown field that is dropped from the object, and for each duplicate field that is encountered. The request will still succeed if there are no other errors, and will only persist the last of any duplicate fields. This is the default in v1.23+ - Strict: This will fail the request with a BadRequest error if any unknown fields would be dropped from the object, or if any duplicate fields are present. The error returned from the server will contain all unknown and duplicate fields encountered. (optional) | 
 **force** | **bool** | Force is going to \&quot;force\&quot; Apply requests. It means user will re-acquire conflicting fields owned by other people. Force flag must be unset for non-apply patch requests. | 

### Return type

**map[string]interface{}**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json-patch+json, application/merge-patch+json, application/apply-patch+yaml
- **Accept**: application/json, application/yaml, application/vnd.kubernetes.protobuf

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PatchNamespacedCustomObjectStatus

> map[string]interface{} PatchNamespacedCustomObjectStatus(ctx, group, version, namespace, plural, name).Body(body).DryRun(dryRun).FieldManager(fieldManager).FieldValidation(fieldValidation).Force(force).Execute()





### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    group := "group_example" // string | the custom resource's group
    version := "version_example" // string | the custom resource's version
    namespace := "namespace_example" // string | The custom resource's namespace
    plural := "plural_example" // string | the custom resource's plural name. For TPRs this would be lowercase plural kind.
    name := "name_example" // string | the custom object's name
    body := map[string]interface{}{ ... } // map[string]interface{} | 
    dryRun := "dryRun_example" // string | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed (optional)
    fieldManager := "fieldManager_example" // string | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. This field is required for apply requests (application/apply-patch) but optional for non-apply patch types (JsonPatch, MergePatch, StrategicMergePatch). (optional)
    fieldValidation := "fieldValidation_example" // string | fieldValidation instructs the server on how to handle objects in the request (POST/PUT/PATCH) containing unknown or duplicate fields. Valid values are: - Ignore: This will ignore any unknown fields that are silently dropped from the object, and will ignore all but the last duplicate field that the decoder encounters. This is the default behavior prior to v1.23. - Warn: This will send a warning via the standard warning response header for each unknown field that is dropped from the object, and for each duplicate field that is encountered. The request will still succeed if there are no other errors, and will only persist the last of any duplicate fields. This is the default in v1.23+ - Strict: This will fail the request with a BadRequest error if any unknown fields would be dropped from the object, or if any duplicate fields are present. The error returned from the server will contain all unknown and duplicate fields encountered. (optional) (optional)
    force := true // bool | Force is going to \"force\" Apply requests. It means user will re-acquire conflicting fields owned by other people. Force flag must be unset for non-apply patch requests. (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.CustomObjectsApi.PatchNamespacedCustomObjectStatus(context.Background(), group, version, namespace, plural, name).Body(body).DryRun(dryRun).FieldManager(fieldManager).FieldValidation(fieldValidation).Force(force).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `CustomObjectsApi.PatchNamespacedCustomObjectStatus``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `PatchNamespacedCustomObjectStatus`: map[string]interface{}
    fmt.Fprintf(os.Stdout, "Response from `CustomObjectsApi.PatchNamespacedCustomObjectStatus`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**group** | **string** | the custom resource&#39;s group | 
**version** | **string** | the custom resource&#39;s version | 
**namespace** | **string** | The custom resource&#39;s namespace | 
**plural** | **string** | the custom resource&#39;s plural name. For TPRs this would be lowercase plural kind. | 
**name** | **string** | the custom object&#39;s name | 

### Other Parameters

Other parameters are passed through a pointer to a apiPatchNamespacedCustomObjectStatusRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------





 **body** | **map[string]interface{}** |  | 
 **dryRun** | **string** | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed | 
 **fieldManager** | **string** | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. This field is required for apply requests (application/apply-patch) but optional for non-apply patch types (JsonPatch, MergePatch, StrategicMergePatch). | 
 **fieldValidation** | **string** | fieldValidation instructs the server on how to handle objects in the request (POST/PUT/PATCH) containing unknown or duplicate fields. Valid values are: - Ignore: This will ignore any unknown fields that are silently dropped from the object, and will ignore all but the last duplicate field that the decoder encounters. This is the default behavior prior to v1.23. - Warn: This will send a warning via the standard warning response header for each unknown field that is dropped from the object, and for each duplicate field that is encountered. The request will still succeed if there are no other errors, and will only persist the last of any duplicate fields. This is the default in v1.23+ - Strict: This will fail the request with a BadRequest error if any unknown fields would be dropped from the object, or if any duplicate fields are present. The error returned from the server will contain all unknown and duplicate fields encountered. (optional) | 
 **force** | **bool** | Force is going to \&quot;force\&quot; Apply requests. It means user will re-acquire conflicting fields owned by other people. Force flag must be unset for non-apply patch requests. | 

### Return type

**map[string]interface{}**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json-patch+json, application/merge-patch+json, application/apply-patch+yaml
- **Accept**: application/json, application/yaml, application/vnd.kubernetes.protobuf

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ReplaceClusterCustomObject

> map[string]interface{} ReplaceClusterCustomObject(ctx, group, version, plural, name).Body(body).DryRun(dryRun).FieldManager(fieldManager).FieldValidation(fieldValidation).Execute()





### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    group := "group_example" // string | the custom resource's group
    version := "version_example" // string | the custom resource's version
    plural := "plural_example" // string | the custom object's plural name. For TPRs this would be lowercase plural kind.
    name := "name_example" // string | the custom object's name
    body := map[string]interface{}{ ... } // map[string]interface{} | The JSON schema of the Resource to replace.
    dryRun := "dryRun_example" // string | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed (optional)
    fieldManager := "fieldManager_example" // string | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. (optional)
    fieldValidation := "fieldValidation_example" // string | fieldValidation instructs the server on how to handle objects in the request (POST/PUT/PATCH) containing unknown or duplicate fields. Valid values are: - Ignore: This will ignore any unknown fields that are silently dropped from the object, and will ignore all but the last duplicate field that the decoder encounters. This is the default behavior prior to v1.23. - Warn: This will send a warning via the standard warning response header for each unknown field that is dropped from the object, and for each duplicate field that is encountered. The request will still succeed if there are no other errors, and will only persist the last of any duplicate fields. This is the default in v1.23+ - Strict: This will fail the request with a BadRequest error if any unknown fields would be dropped from the object, or if any duplicate fields are present. The error returned from the server will contain all unknown and duplicate fields encountered. (optional) (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.CustomObjectsApi.ReplaceClusterCustomObject(context.Background(), group, version, plural, name).Body(body).DryRun(dryRun).FieldManager(fieldManager).FieldValidation(fieldValidation).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `CustomObjectsApi.ReplaceClusterCustomObject``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ReplaceClusterCustomObject`: map[string]interface{}
    fmt.Fprintf(os.Stdout, "Response from `CustomObjectsApi.ReplaceClusterCustomObject`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**group** | **string** | the custom resource&#39;s group | 
**version** | **string** | the custom resource&#39;s version | 
**plural** | **string** | the custom object&#39;s plural name. For TPRs this would be lowercase plural kind. | 
**name** | **string** | the custom object&#39;s name | 

### Other Parameters

Other parameters are passed through a pointer to a apiReplaceClusterCustomObjectRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------




 **body** | **map[string]interface{}** | The JSON schema of the Resource to replace. | 
 **dryRun** | **string** | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed | 
 **fieldManager** | **string** | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. | 
 **fieldValidation** | **string** | fieldValidation instructs the server on how to handle objects in the request (POST/PUT/PATCH) containing unknown or duplicate fields. Valid values are: - Ignore: This will ignore any unknown fields that are silently dropped from the object, and will ignore all but the last duplicate field that the decoder encounters. This is the default behavior prior to v1.23. - Warn: This will send a warning via the standard warning response header for each unknown field that is dropped from the object, and for each duplicate field that is encountered. The request will still succeed if there are no other errors, and will only persist the last of any duplicate fields. This is the default in v1.23+ - Strict: This will fail the request with a BadRequest error if any unknown fields would be dropped from the object, or if any duplicate fields are present. The error returned from the server will contain all unknown and duplicate fields encountered. (optional) | 

### Return type

**map[string]interface{}**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ReplaceClusterCustomObjectScale

> map[string]interface{} ReplaceClusterCustomObjectScale(ctx, group, version, plural, name).Body(body).DryRun(dryRun).FieldManager(fieldManager).FieldValidation(fieldValidation).Execute()





### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    group := "group_example" // string | the custom resource's group
    version := "version_example" // string | the custom resource's version
    plural := "plural_example" // string | the custom resource's plural name. For TPRs this would be lowercase plural kind.
    name := "name_example" // string | the custom object's name
    body := map[string]interface{}{ ... } // map[string]interface{} | 
    dryRun := "dryRun_example" // string | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed (optional)
    fieldManager := "fieldManager_example" // string | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. (optional)
    fieldValidation := "fieldValidation_example" // string | fieldValidation instructs the server on how to handle objects in the request (POST/PUT/PATCH) containing unknown or duplicate fields. Valid values are: - Ignore: This will ignore any unknown fields that are silently dropped from the object, and will ignore all but the last duplicate field that the decoder encounters. This is the default behavior prior to v1.23. - Warn: This will send a warning via the standard warning response header for each unknown field that is dropped from the object, and for each duplicate field that is encountered. The request will still succeed if there are no other errors, and will only persist the last of any duplicate fields. This is the default in v1.23+ - Strict: This will fail the request with a BadRequest error if any unknown fields would be dropped from the object, or if any duplicate fields are present. The error returned from the server will contain all unknown and duplicate fields encountered. (optional) (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.CustomObjectsApi.ReplaceClusterCustomObjectScale(context.Background(), group, version, plural, name).Body(body).DryRun(dryRun).FieldManager(fieldManager).FieldValidation(fieldValidation).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `CustomObjectsApi.ReplaceClusterCustomObjectScale``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ReplaceClusterCustomObjectScale`: map[string]interface{}
    fmt.Fprintf(os.Stdout, "Response from `CustomObjectsApi.ReplaceClusterCustomObjectScale`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**group** | **string** | the custom resource&#39;s group | 
**version** | **string** | the custom resource&#39;s version | 
**plural** | **string** | the custom resource&#39;s plural name. For TPRs this would be lowercase plural kind. | 
**name** | **string** | the custom object&#39;s name | 

### Other Parameters

Other parameters are passed through a pointer to a apiReplaceClusterCustomObjectScaleRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------




 **body** | **map[string]interface{}** |  | 
 **dryRun** | **string** | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed | 
 **fieldManager** | **string** | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. | 
 **fieldValidation** | **string** | fieldValidation instructs the server on how to handle objects in the request (POST/PUT/PATCH) containing unknown or duplicate fields. Valid values are: - Ignore: This will ignore any unknown fields that are silently dropped from the object, and will ignore all but the last duplicate field that the decoder encounters. This is the default behavior prior to v1.23. - Warn: This will send a warning via the standard warning response header for each unknown field that is dropped from the object, and for each duplicate field that is encountered. The request will still succeed if there are no other errors, and will only persist the last of any duplicate fields. This is the default in v1.23+ - Strict: This will fail the request with a BadRequest error if any unknown fields would be dropped from the object, or if any duplicate fields are present. The error returned from the server will contain all unknown and duplicate fields encountered. (optional) | 

### Return type

**map[string]interface{}**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/yaml, application/vnd.kubernetes.protobuf

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ReplaceClusterCustomObjectStatus

> map[string]interface{} ReplaceClusterCustomObjectStatus(ctx, group, version, plural, name).Body(body).DryRun(dryRun).FieldManager(fieldManager).FieldValidation(fieldValidation).Execute()





### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    group := "group_example" // string | the custom resource's group
    version := "version_example" // string | the custom resource's version
    plural := "plural_example" // string | the custom resource's plural name. For TPRs this would be lowercase plural kind.
    name := "name_example" // string | the custom object's name
    body := map[string]interface{}{ ... } // map[string]interface{} | 
    dryRun := "dryRun_example" // string | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed (optional)
    fieldManager := "fieldManager_example" // string | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. (optional)
    fieldValidation := "fieldValidation_example" // string | fieldValidation instructs the server on how to handle objects in the request (POST/PUT/PATCH) containing unknown or duplicate fields. Valid values are: - Ignore: This will ignore any unknown fields that are silently dropped from the object, and will ignore all but the last duplicate field that the decoder encounters. This is the default behavior prior to v1.23. - Warn: This will send a warning via the standard warning response header for each unknown field that is dropped from the object, and for each duplicate field that is encountered. The request will still succeed if there are no other errors, and will only persist the last of any duplicate fields. This is the default in v1.23+ - Strict: This will fail the request with a BadRequest error if any unknown fields would be dropped from the object, or if any duplicate fields are present. The error returned from the server will contain all unknown and duplicate fields encountered. (optional) (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.CustomObjectsApi.ReplaceClusterCustomObjectStatus(context.Background(), group, version, plural, name).Body(body).DryRun(dryRun).FieldManager(fieldManager).FieldValidation(fieldValidation).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `CustomObjectsApi.ReplaceClusterCustomObjectStatus``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ReplaceClusterCustomObjectStatus`: map[string]interface{}
    fmt.Fprintf(os.Stdout, "Response from `CustomObjectsApi.ReplaceClusterCustomObjectStatus`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**group** | **string** | the custom resource&#39;s group | 
**version** | **string** | the custom resource&#39;s version | 
**plural** | **string** | the custom resource&#39;s plural name. For TPRs this would be lowercase plural kind. | 
**name** | **string** | the custom object&#39;s name | 

### Other Parameters

Other parameters are passed through a pointer to a apiReplaceClusterCustomObjectStatusRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------




 **body** | **map[string]interface{}** |  | 
 **dryRun** | **string** | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed | 
 **fieldManager** | **string** | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. | 
 **fieldValidation** | **string** | fieldValidation instructs the server on how to handle objects in the request (POST/PUT/PATCH) containing unknown or duplicate fields. Valid values are: - Ignore: This will ignore any unknown fields that are silently dropped from the object, and will ignore all but the last duplicate field that the decoder encounters. This is the default behavior prior to v1.23. - Warn: This will send a warning via the standard warning response header for each unknown field that is dropped from the object, and for each duplicate field that is encountered. The request will still succeed if there are no other errors, and will only persist the last of any duplicate fields. This is the default in v1.23+ - Strict: This will fail the request with a BadRequest error if any unknown fields would be dropped from the object, or if any duplicate fields are present. The error returned from the server will contain all unknown and duplicate fields encountered. (optional) | 

### Return type

**map[string]interface{}**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/yaml, application/vnd.kubernetes.protobuf

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ReplaceNamespacedCustomObject

> map[string]interface{} ReplaceNamespacedCustomObject(ctx, group, version, namespace, plural, name).Body(body).DryRun(dryRun).FieldManager(fieldManager).FieldValidation(fieldValidation).Execute()





### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    group := "group_example" // string | the custom resource's group
    version := "version_example" // string | the custom resource's version
    namespace := "namespace_example" // string | The custom resource's namespace
    plural := "plural_example" // string | the custom resource's plural name. For TPRs this would be lowercase plural kind.
    name := "name_example" // string | the custom object's name
    body := map[string]interface{}{ ... } // map[string]interface{} | The JSON schema of the Resource to replace.
    dryRun := "dryRun_example" // string | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed (optional)
    fieldManager := "fieldManager_example" // string | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. (optional)
    fieldValidation := "fieldValidation_example" // string | fieldValidation instructs the server on how to handle objects in the request (POST/PUT/PATCH) containing unknown or duplicate fields. Valid values are: - Ignore: This will ignore any unknown fields that are silently dropped from the object, and will ignore all but the last duplicate field that the decoder encounters. This is the default behavior prior to v1.23. - Warn: This will send a warning via the standard warning response header for each unknown field that is dropped from the object, and for each duplicate field that is encountered. The request will still succeed if there are no other errors, and will only persist the last of any duplicate fields. This is the default in v1.23+ - Strict: This will fail the request with a BadRequest error if any unknown fields would be dropped from the object, or if any duplicate fields are present. The error returned from the server will contain all unknown and duplicate fields encountered. (optional) (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.CustomObjectsApi.ReplaceNamespacedCustomObject(context.Background(), group, version, namespace, plural, name).Body(body).DryRun(dryRun).FieldManager(fieldManager).FieldValidation(fieldValidation).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `CustomObjectsApi.ReplaceNamespacedCustomObject``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ReplaceNamespacedCustomObject`: map[string]interface{}
    fmt.Fprintf(os.Stdout, "Response from `CustomObjectsApi.ReplaceNamespacedCustomObject`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**group** | **string** | the custom resource&#39;s group | 
**version** | **string** | the custom resource&#39;s version | 
**namespace** | **string** | The custom resource&#39;s namespace | 
**plural** | **string** | the custom resource&#39;s plural name. For TPRs this would be lowercase plural kind. | 
**name** | **string** | the custom object&#39;s name | 

### Other Parameters

Other parameters are passed through a pointer to a apiReplaceNamespacedCustomObjectRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------





 **body** | **map[string]interface{}** | The JSON schema of the Resource to replace. | 
 **dryRun** | **string** | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed | 
 **fieldManager** | **string** | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. | 
 **fieldValidation** | **string** | fieldValidation instructs the server on how to handle objects in the request (POST/PUT/PATCH) containing unknown or duplicate fields. Valid values are: - Ignore: This will ignore any unknown fields that are silently dropped from the object, and will ignore all but the last duplicate field that the decoder encounters. This is the default behavior prior to v1.23. - Warn: This will send a warning via the standard warning response header for each unknown field that is dropped from the object, and for each duplicate field that is encountered. The request will still succeed if there are no other errors, and will only persist the last of any duplicate fields. This is the default in v1.23+ - Strict: This will fail the request with a BadRequest error if any unknown fields would be dropped from the object, or if any duplicate fields are present. The error returned from the server will contain all unknown and duplicate fields encountered. (optional) | 

### Return type

**map[string]interface{}**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ReplaceNamespacedCustomObjectScale

> map[string]interface{} ReplaceNamespacedCustomObjectScale(ctx, group, version, namespace, plural, name).Body(body).DryRun(dryRun).FieldManager(fieldManager).FieldValidation(fieldValidation).Execute()





### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    group := "group_example" // string | the custom resource's group
    version := "version_example" // string | the custom resource's version
    namespace := "namespace_example" // string | The custom resource's namespace
    plural := "plural_example" // string | the custom resource's plural name. For TPRs this would be lowercase plural kind.
    name := "name_example" // string | the custom object's name
    body := map[string]interface{}{ ... } // map[string]interface{} | 
    dryRun := "dryRun_example" // string | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed (optional)
    fieldManager := "fieldManager_example" // string | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. (optional)
    fieldValidation := "fieldValidation_example" // string | fieldValidation instructs the server on how to handle objects in the request (POST/PUT/PATCH) containing unknown or duplicate fields. Valid values are: - Ignore: This will ignore any unknown fields that are silently dropped from the object, and will ignore all but the last duplicate field that the decoder encounters. This is the default behavior prior to v1.23. - Warn: This will send a warning via the standard warning response header for each unknown field that is dropped from the object, and for each duplicate field that is encountered. The request will still succeed if there are no other errors, and will only persist the last of any duplicate fields. This is the default in v1.23+ - Strict: This will fail the request with a BadRequest error if any unknown fields would be dropped from the object, or if any duplicate fields are present. The error returned from the server will contain all unknown and duplicate fields encountered. (optional) (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.CustomObjectsApi.ReplaceNamespacedCustomObjectScale(context.Background(), group, version, namespace, plural, name).Body(body).DryRun(dryRun).FieldManager(fieldManager).FieldValidation(fieldValidation).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `CustomObjectsApi.ReplaceNamespacedCustomObjectScale``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ReplaceNamespacedCustomObjectScale`: map[string]interface{}
    fmt.Fprintf(os.Stdout, "Response from `CustomObjectsApi.ReplaceNamespacedCustomObjectScale`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**group** | **string** | the custom resource&#39;s group | 
**version** | **string** | the custom resource&#39;s version | 
**namespace** | **string** | The custom resource&#39;s namespace | 
**plural** | **string** | the custom resource&#39;s plural name. For TPRs this would be lowercase plural kind. | 
**name** | **string** | the custom object&#39;s name | 

### Other Parameters

Other parameters are passed through a pointer to a apiReplaceNamespacedCustomObjectScaleRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------





 **body** | **map[string]interface{}** |  | 
 **dryRun** | **string** | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed | 
 **fieldManager** | **string** | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. | 
 **fieldValidation** | **string** | fieldValidation instructs the server on how to handle objects in the request (POST/PUT/PATCH) containing unknown or duplicate fields. Valid values are: - Ignore: This will ignore any unknown fields that are silently dropped from the object, and will ignore all but the last duplicate field that the decoder encounters. This is the default behavior prior to v1.23. - Warn: This will send a warning via the standard warning response header for each unknown field that is dropped from the object, and for each duplicate field that is encountered. The request will still succeed if there are no other errors, and will only persist the last of any duplicate fields. This is the default in v1.23+ - Strict: This will fail the request with a BadRequest error if any unknown fields would be dropped from the object, or if any duplicate fields are present. The error returned from the server will contain all unknown and duplicate fields encountered. (optional) | 

### Return type

**map[string]interface{}**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/yaml, application/vnd.kubernetes.protobuf

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ReplaceNamespacedCustomObjectStatus

> map[string]interface{} ReplaceNamespacedCustomObjectStatus(ctx, group, version, namespace, plural, name).Body(body).DryRun(dryRun).FieldManager(fieldManager).FieldValidation(fieldValidation).Execute()





### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    group := "group_example" // string | the custom resource's group
    version := "version_example" // string | the custom resource's version
    namespace := "namespace_example" // string | The custom resource's namespace
    plural := "plural_example" // string | the custom resource's plural name. For TPRs this would be lowercase plural kind.
    name := "name_example" // string | the custom object's name
    body := map[string]interface{}{ ... } // map[string]interface{} | 
    dryRun := "dryRun_example" // string | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed (optional)
    fieldManager := "fieldManager_example" // string | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. (optional)
    fieldValidation := "fieldValidation_example" // string | fieldValidation instructs the server on how to handle objects in the request (POST/PUT/PATCH) containing unknown or duplicate fields. Valid values are: - Ignore: This will ignore any unknown fields that are silently dropped from the object, and will ignore all but the last duplicate field that the decoder encounters. This is the default behavior prior to v1.23. - Warn: This will send a warning via the standard warning response header for each unknown field that is dropped from the object, and for each duplicate field that is encountered. The request will still succeed if there are no other errors, and will only persist the last of any duplicate fields. This is the default in v1.23+ - Strict: This will fail the request with a BadRequest error if any unknown fields would be dropped from the object, or if any duplicate fields are present. The error returned from the server will contain all unknown and duplicate fields encountered. (optional) (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.CustomObjectsApi.ReplaceNamespacedCustomObjectStatus(context.Background(), group, version, namespace, plural, name).Body(body).DryRun(dryRun).FieldManager(fieldManager).FieldValidation(fieldValidation).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `CustomObjectsApi.ReplaceNamespacedCustomObjectStatus``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ReplaceNamespacedCustomObjectStatus`: map[string]interface{}
    fmt.Fprintf(os.Stdout, "Response from `CustomObjectsApi.ReplaceNamespacedCustomObjectStatus`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**group** | **string** | the custom resource&#39;s group | 
**version** | **string** | the custom resource&#39;s version | 
**namespace** | **string** | The custom resource&#39;s namespace | 
**plural** | **string** | the custom resource&#39;s plural name. For TPRs this would be lowercase plural kind. | 
**name** | **string** | the custom object&#39;s name | 

### Other Parameters

Other parameters are passed through a pointer to a apiReplaceNamespacedCustomObjectStatusRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------





 **body** | **map[string]interface{}** |  | 
 **dryRun** | **string** | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed | 
 **fieldManager** | **string** | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. | 
 **fieldValidation** | **string** | fieldValidation instructs the server on how to handle objects in the request (POST/PUT/PATCH) containing unknown or duplicate fields. Valid values are: - Ignore: This will ignore any unknown fields that are silently dropped from the object, and will ignore all but the last duplicate field that the decoder encounters. This is the default behavior prior to v1.23. - Warn: This will send a warning via the standard warning response header for each unknown field that is dropped from the object, and for each duplicate field that is encountered. The request will still succeed if there are no other errors, and will only persist the last of any duplicate fields. This is the default in v1.23+ - Strict: This will fail the request with a BadRequest error if any unknown fields would be dropped from the object, or if any duplicate fields are present. The error returned from the server will contain all unknown and duplicate fields encountered. (optional) | 

### Return type

**map[string]interface{}**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/yaml, application/vnd.kubernetes.protobuf

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

