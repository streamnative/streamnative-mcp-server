# \AuthorizationStreamnativeIoV1alpha1Api

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount**](AuthorizationStreamnativeIoV1alpha1Api.md#CreateAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount) | **Post** /apis/authorization.streamnative.io/v1alpha1/namespaces/{namespace}/iamaccounts | 
[**CreateAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus**](AuthorizationStreamnativeIoV1alpha1Api.md#CreateAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus) | **Post** /apis/authorization.streamnative.io/v1alpha1/namespaces/{namespace}/iamaccounts/{name}/status | 
[**CreateAuthorizationStreamnativeIoV1alpha1SelfSubjectRbacReview**](AuthorizationStreamnativeIoV1alpha1Api.md#CreateAuthorizationStreamnativeIoV1alpha1SelfSubjectRbacReview) | **Post** /apis/authorization.streamnative.io/v1alpha1/selfsubjectrbacreviews | 
[**CreateAuthorizationStreamnativeIoV1alpha1SelfSubjectRulesReview**](AuthorizationStreamnativeIoV1alpha1Api.md#CreateAuthorizationStreamnativeIoV1alpha1SelfSubjectRulesReview) | **Post** /apis/authorization.streamnative.io/v1alpha1/selfsubjectrulesreviews | 
[**CreateAuthorizationStreamnativeIoV1alpha1SelfSubjectUserReview**](AuthorizationStreamnativeIoV1alpha1Api.md#CreateAuthorizationStreamnativeIoV1alpha1SelfSubjectUserReview) | **Post** /apis/authorization.streamnative.io/v1alpha1/selfsubjectuserreviews | 
[**CreateAuthorizationStreamnativeIoV1alpha1SubjectRoleReview**](AuthorizationStreamnativeIoV1alpha1Api.md#CreateAuthorizationStreamnativeIoV1alpha1SubjectRoleReview) | **Post** /apis/authorization.streamnative.io/v1alpha1/subjectrolereviews | 
[**CreateAuthorizationStreamnativeIoV1alpha1SubjectRulesReview**](AuthorizationStreamnativeIoV1alpha1Api.md#CreateAuthorizationStreamnativeIoV1alpha1SubjectRulesReview) | **Post** /apis/authorization.streamnative.io/v1alpha1/subjectrulesreviews | 
[**DeleteAuthorizationStreamnativeIoV1alpha1CollectionNamespacedIamAccount**](AuthorizationStreamnativeIoV1alpha1Api.md#DeleteAuthorizationStreamnativeIoV1alpha1CollectionNamespacedIamAccount) | **Delete** /apis/authorization.streamnative.io/v1alpha1/namespaces/{namespace}/iamaccounts | 
[**DeleteAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount**](AuthorizationStreamnativeIoV1alpha1Api.md#DeleteAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount) | **Delete** /apis/authorization.streamnative.io/v1alpha1/namespaces/{namespace}/iamaccounts/{name} | 
[**DeleteAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus**](AuthorizationStreamnativeIoV1alpha1Api.md#DeleteAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus) | **Delete** /apis/authorization.streamnative.io/v1alpha1/namespaces/{namespace}/iamaccounts/{name}/status | 
[**GetAuthorizationStreamnativeIoV1alpha1APIResources**](AuthorizationStreamnativeIoV1alpha1Api.md#GetAuthorizationStreamnativeIoV1alpha1APIResources) | **Get** /apis/authorization.streamnative.io/v1alpha1/ | 
[**ListAuthorizationStreamnativeIoV1alpha1IamAccountForAllNamespaces**](AuthorizationStreamnativeIoV1alpha1Api.md#ListAuthorizationStreamnativeIoV1alpha1IamAccountForAllNamespaces) | **Get** /apis/authorization.streamnative.io/v1alpha1/iamaccounts | 
[**ListAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount**](AuthorizationStreamnativeIoV1alpha1Api.md#ListAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount) | **Get** /apis/authorization.streamnative.io/v1alpha1/namespaces/{namespace}/iamaccounts | 
[**PatchAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount**](AuthorizationStreamnativeIoV1alpha1Api.md#PatchAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount) | **Patch** /apis/authorization.streamnative.io/v1alpha1/namespaces/{namespace}/iamaccounts/{name} | 
[**PatchAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus**](AuthorizationStreamnativeIoV1alpha1Api.md#PatchAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus) | **Patch** /apis/authorization.streamnative.io/v1alpha1/namespaces/{namespace}/iamaccounts/{name}/status | 
[**ReadAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount**](AuthorizationStreamnativeIoV1alpha1Api.md#ReadAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount) | **Get** /apis/authorization.streamnative.io/v1alpha1/namespaces/{namespace}/iamaccounts/{name} | 
[**ReadAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus**](AuthorizationStreamnativeIoV1alpha1Api.md#ReadAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus) | **Get** /apis/authorization.streamnative.io/v1alpha1/namespaces/{namespace}/iamaccounts/{name}/status | 
[**ReplaceAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount**](AuthorizationStreamnativeIoV1alpha1Api.md#ReplaceAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount) | **Put** /apis/authorization.streamnative.io/v1alpha1/namespaces/{namespace}/iamaccounts/{name} | 
[**ReplaceAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus**](AuthorizationStreamnativeIoV1alpha1Api.md#ReplaceAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus) | **Put** /apis/authorization.streamnative.io/v1alpha1/namespaces/{namespace}/iamaccounts/{name}/status | 
[**WatchAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus**](AuthorizationStreamnativeIoV1alpha1Api.md#WatchAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus) | **Get** /apis/authorization.streamnative.io/v1alpha1/watch/namespaces/{namespace}/iamaccounts/{name}/status | 



## CreateAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount

> ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount CreateAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount(ctx, namespace).Body(body).Pretty(pretty).DryRun(dryRun).FieldManager(fieldManager).Execute()





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
    namespace := "namespace_example" // string | object name and auth scope, such as for teams and projects
    body := *openapiclient.NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount() // ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount | 
    pretty := "pretty_example" // string | If 'true', then the output is pretty printed. (optional)
    dryRun := "dryRun_example" // string | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed (optional)
    fieldManager := "fieldManager_example" // string | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.AuthorizationStreamnativeIoV1alpha1Api.CreateAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount(context.Background(), namespace).Body(body).Pretty(pretty).DryRun(dryRun).FieldManager(fieldManager).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `AuthorizationStreamnativeIoV1alpha1Api.CreateAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `CreateAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount`: ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount
    fmt.Fprintf(os.Stdout, "Response from `AuthorizationStreamnativeIoV1alpha1Api.CreateAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**namespace** | **string** | object name and auth scope, such as for teams and projects | 

### Other Parameters

Other parameters are passed through a pointer to a apiCreateAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **body** | [**ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount**](ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount.md) |  | 
 **pretty** | **string** | If &#39;true&#39;, then the output is pretty printed. | 
 **dryRun** | **string** | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed | 
 **fieldManager** | **string** | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. | 

### Return type

[**ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount**](ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/yaml, application/vnd.kubernetes.protobuf

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus

> ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount CreateAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus(ctx, name, namespace).Body(body).Pretty(pretty).DryRun(dryRun).FieldManager(fieldManager).Execute()





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
    name := "name_example" // string | name of the IamAccount
    namespace := "namespace_example" // string | object name and auth scope, such as for teams and projects
    body := *openapiclient.NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount() // ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount | 
    pretty := "pretty_example" // string | If 'true', then the output is pretty printed. (optional)
    dryRun := "dryRun_example" // string | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed (optional)
    fieldManager := "fieldManager_example" // string | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.AuthorizationStreamnativeIoV1alpha1Api.CreateAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus(context.Background(), name, namespace).Body(body).Pretty(pretty).DryRun(dryRun).FieldManager(fieldManager).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `AuthorizationStreamnativeIoV1alpha1Api.CreateAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `CreateAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus`: ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount
    fmt.Fprintf(os.Stdout, "Response from `AuthorizationStreamnativeIoV1alpha1Api.CreateAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**name** | **string** | name of the IamAccount | 
**namespace** | **string** | object name and auth scope, such as for teams and projects | 

### Other Parameters

Other parameters are passed through a pointer to a apiCreateAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatusRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **body** | [**ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount**](ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount.md) |  | 
 **pretty** | **string** | If &#39;true&#39;, then the output is pretty printed. | 
 **dryRun** | **string** | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed | 
 **fieldManager** | **string** | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. | 

### Return type

[**ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount**](ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/yaml, application/vnd.kubernetes.protobuf

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateAuthorizationStreamnativeIoV1alpha1SelfSubjectRbacReview

> ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReview CreateAuthorizationStreamnativeIoV1alpha1SelfSubjectRbacReview(ctx).Body(body).DryRun(dryRun).FieldManager(fieldManager).Pretty(pretty).Execute()





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
    body := *openapiclient.NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReview() // ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReview | 
    dryRun := "dryRun_example" // string | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed (optional)
    fieldManager := "fieldManager_example" // string | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. (optional)
    pretty := "pretty_example" // string | If 'true', then the output is pretty printed. (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.AuthorizationStreamnativeIoV1alpha1Api.CreateAuthorizationStreamnativeIoV1alpha1SelfSubjectRbacReview(context.Background()).Body(body).DryRun(dryRun).FieldManager(fieldManager).Pretty(pretty).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `AuthorizationStreamnativeIoV1alpha1Api.CreateAuthorizationStreamnativeIoV1alpha1SelfSubjectRbacReview``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `CreateAuthorizationStreamnativeIoV1alpha1SelfSubjectRbacReview`: ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReview
    fmt.Fprintf(os.Stdout, "Response from `AuthorizationStreamnativeIoV1alpha1Api.CreateAuthorizationStreamnativeIoV1alpha1SelfSubjectRbacReview`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateAuthorizationStreamnativeIoV1alpha1SelfSubjectRbacReviewRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReview**](ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReview.md) |  | 
 **dryRun** | **string** | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed | 
 **fieldManager** | **string** | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. | 
 **pretty** | **string** | If &#39;true&#39;, then the output is pretty printed. | 

### Return type

[**ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReview**](ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRbacReview.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/yaml, application/vnd.kubernetes.protobuf

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateAuthorizationStreamnativeIoV1alpha1SelfSubjectRulesReview

> ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRulesReview CreateAuthorizationStreamnativeIoV1alpha1SelfSubjectRulesReview(ctx).Body(body).DryRun(dryRun).FieldManager(fieldManager).Pretty(pretty).Execute()





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
    body := *openapiclient.NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRulesReview() // ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRulesReview | 
    dryRun := "dryRun_example" // string | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed (optional)
    fieldManager := "fieldManager_example" // string | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. (optional)
    pretty := "pretty_example" // string | If 'true', then the output is pretty printed. (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.AuthorizationStreamnativeIoV1alpha1Api.CreateAuthorizationStreamnativeIoV1alpha1SelfSubjectRulesReview(context.Background()).Body(body).DryRun(dryRun).FieldManager(fieldManager).Pretty(pretty).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `AuthorizationStreamnativeIoV1alpha1Api.CreateAuthorizationStreamnativeIoV1alpha1SelfSubjectRulesReview``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `CreateAuthorizationStreamnativeIoV1alpha1SelfSubjectRulesReview`: ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRulesReview
    fmt.Fprintf(os.Stdout, "Response from `AuthorizationStreamnativeIoV1alpha1Api.CreateAuthorizationStreamnativeIoV1alpha1SelfSubjectRulesReview`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateAuthorizationStreamnativeIoV1alpha1SelfSubjectRulesReviewRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRulesReview**](ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRulesReview.md) |  | 
 **dryRun** | **string** | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed | 
 **fieldManager** | **string** | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. | 
 **pretty** | **string** | If &#39;true&#39;, then the output is pretty printed. | 

### Return type

[**ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRulesReview**](ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectRulesReview.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/yaml, application/vnd.kubernetes.protobuf

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateAuthorizationStreamnativeIoV1alpha1SelfSubjectUserReview

> ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectUserReview CreateAuthorizationStreamnativeIoV1alpha1SelfSubjectUserReview(ctx).Body(body).DryRun(dryRun).FieldManager(fieldManager).Pretty(pretty).Execute()





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
    body := *openapiclient.NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectUserReview() // ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectUserReview | 
    dryRun := "dryRun_example" // string | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed (optional)
    fieldManager := "fieldManager_example" // string | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. (optional)
    pretty := "pretty_example" // string | If 'true', then the output is pretty printed. (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.AuthorizationStreamnativeIoV1alpha1Api.CreateAuthorizationStreamnativeIoV1alpha1SelfSubjectUserReview(context.Background()).Body(body).DryRun(dryRun).FieldManager(fieldManager).Pretty(pretty).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `AuthorizationStreamnativeIoV1alpha1Api.CreateAuthorizationStreamnativeIoV1alpha1SelfSubjectUserReview``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `CreateAuthorizationStreamnativeIoV1alpha1SelfSubjectUserReview`: ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectUserReview
    fmt.Fprintf(os.Stdout, "Response from `AuthorizationStreamnativeIoV1alpha1Api.CreateAuthorizationStreamnativeIoV1alpha1SelfSubjectUserReview`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateAuthorizationStreamnativeIoV1alpha1SelfSubjectUserReviewRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectUserReview**](ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectUserReview.md) |  | 
 **dryRun** | **string** | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed | 
 **fieldManager** | **string** | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. | 
 **pretty** | **string** | If &#39;true&#39;, then the output is pretty printed. | 

### Return type

[**ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectUserReview**](ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SelfSubjectUserReview.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/yaml, application/vnd.kubernetes.protobuf

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateAuthorizationStreamnativeIoV1alpha1SubjectRoleReview

> ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRoleReview CreateAuthorizationStreamnativeIoV1alpha1SubjectRoleReview(ctx).Body(body).DryRun(dryRun).FieldManager(fieldManager).Pretty(pretty).Execute()





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
    body := *openapiclient.NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRoleReview() // ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRoleReview | 
    dryRun := "dryRun_example" // string | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed (optional)
    fieldManager := "fieldManager_example" // string | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. (optional)
    pretty := "pretty_example" // string | If 'true', then the output is pretty printed. (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.AuthorizationStreamnativeIoV1alpha1Api.CreateAuthorizationStreamnativeIoV1alpha1SubjectRoleReview(context.Background()).Body(body).DryRun(dryRun).FieldManager(fieldManager).Pretty(pretty).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `AuthorizationStreamnativeIoV1alpha1Api.CreateAuthorizationStreamnativeIoV1alpha1SubjectRoleReview``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `CreateAuthorizationStreamnativeIoV1alpha1SubjectRoleReview`: ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRoleReview
    fmt.Fprintf(os.Stdout, "Response from `AuthorizationStreamnativeIoV1alpha1Api.CreateAuthorizationStreamnativeIoV1alpha1SubjectRoleReview`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateAuthorizationStreamnativeIoV1alpha1SubjectRoleReviewRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRoleReview**](ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRoleReview.md) |  | 
 **dryRun** | **string** | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed | 
 **fieldManager** | **string** | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. | 
 **pretty** | **string** | If &#39;true&#39;, then the output is pretty printed. | 

### Return type

[**ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRoleReview**](ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRoleReview.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/yaml, application/vnd.kubernetes.protobuf

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateAuthorizationStreamnativeIoV1alpha1SubjectRulesReview

> ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReview CreateAuthorizationStreamnativeIoV1alpha1SubjectRulesReview(ctx).Body(body).DryRun(dryRun).FieldManager(fieldManager).Pretty(pretty).Execute()





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
    body := *openapiclient.NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReview() // ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReview | 
    dryRun := "dryRun_example" // string | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed (optional)
    fieldManager := "fieldManager_example" // string | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. (optional)
    pretty := "pretty_example" // string | If 'true', then the output is pretty printed. (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.AuthorizationStreamnativeIoV1alpha1Api.CreateAuthorizationStreamnativeIoV1alpha1SubjectRulesReview(context.Background()).Body(body).DryRun(dryRun).FieldManager(fieldManager).Pretty(pretty).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `AuthorizationStreamnativeIoV1alpha1Api.CreateAuthorizationStreamnativeIoV1alpha1SubjectRulesReview``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `CreateAuthorizationStreamnativeIoV1alpha1SubjectRulesReview`: ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReview
    fmt.Fprintf(os.Stdout, "Response from `AuthorizationStreamnativeIoV1alpha1Api.CreateAuthorizationStreamnativeIoV1alpha1SubjectRulesReview`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateAuthorizationStreamnativeIoV1alpha1SubjectRulesReviewRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReview**](ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReview.md) |  | 
 **dryRun** | **string** | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed | 
 **fieldManager** | **string** | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. | 
 **pretty** | **string** | If &#39;true&#39;, then the output is pretty printed. | 

### Return type

[**ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReview**](ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1SubjectRulesReview.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/yaml, application/vnd.kubernetes.protobuf

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteAuthorizationStreamnativeIoV1alpha1CollectionNamespacedIamAccount

> V1Status DeleteAuthorizationStreamnativeIoV1alpha1CollectionNamespacedIamAccount(ctx, namespace).Pretty(pretty).Continue_(continue_).DryRun(dryRun).FieldSelector(fieldSelector).GracePeriodSeconds(gracePeriodSeconds).LabelSelector(labelSelector).Limit(limit).OrphanDependents(orphanDependents).PropagationPolicy(propagationPolicy).ResourceVersion(resourceVersion).ResourceVersionMatch(resourceVersionMatch).TimeoutSeconds(timeoutSeconds).Body(body).Execute()





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
    namespace := "namespace_example" // string | object name and auth scope, such as for teams and projects
    pretty := "pretty_example" // string | If 'true', then the output is pretty printed. (optional)
    continue_ := "continue__example" // string | The continue option should be set when retrieving more results from the server. Since this value is server defined, clients may only use the continue value from a previous query result with identical query parameters (except for the value of continue) and the server may reject a continue value it does not recognize. If the specified continue value is no longer valid whether due to expiration (generally five to fifteen minutes) or a configuration change on the server, the server will respond with a 410 ResourceExpired error together with a continue token. If the client needs a consistent list, it must restart their list without the continue field. Otherwise, the client may send another list request with the token received with the 410 error, the server will respond with a list starting from the next key, but from the latest snapshot, which is inconsistent from the previous list results - objects that are created, modified, or deleted after the first list request will be included in the response, as long as their keys are after the \"next key\".  This field is not supported when watch is true. Clients may start a watch from the last resourceVersion value returned by the server and not miss any modifications. (optional)
    dryRun := "dryRun_example" // string | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed (optional)
    fieldSelector := "fieldSelector_example" // string | A selector to restrict the list of returned objects by their fields. Defaults to everything. (optional)
    gracePeriodSeconds := int32(56) // int32 | The duration in seconds before the object should be deleted. Value must be non-negative integer. The value zero indicates delete immediately. If this value is nil, the default grace period for the specified type will be used. Defaults to a per object value if not specified. zero means delete immediately. (optional)
    labelSelector := "labelSelector_example" // string | A selector to restrict the list of returned objects by their labels. Defaults to everything. (optional)
    limit := int32(56) // int32 | limit is a maximum number of responses to return for a list call. If more items exist, the server will set the `continue` field on the list metadata to a value that can be used with the same initial query to retrieve the next set of results. Setting a limit may return fewer than the requested amount of items (up to zero items) in the event all requested objects are filtered out and clients should only use the presence of the continue field to determine whether more results are available. Servers may choose not to support the limit argument and will return all of the available results. If limit is specified and the continue field is empty, clients may assume that no more results are available. This field is not supported if watch is true.  The server guarantees that the objects returned when using continue will be identical to issuing a single list call without a limit - that is, no objects created, modified, or deleted after the first request is issued will be included in any subsequent continued requests. This is sometimes referred to as a consistent snapshot, and ensures that a client that is using limit to receive smaller chunks of a very large result can ensure they see all possible objects. If objects are updated during a chunked list the version of the object that was present at the time the first list result was calculated is returned. (optional)
    orphanDependents := true // bool | Deprecated: please use the PropagationPolicy, this field will be deprecated in 1.7. Should the dependent objects be orphaned. If true/false, the \"orphan\" finalizer will be added to/removed from the object's finalizers list. Either this field or PropagationPolicy may be set, but not both. (optional)
    propagationPolicy := "propagationPolicy_example" // string | Whether and how garbage collection will be performed. Either this field or OrphanDependents may be set, but not both. The default policy is decided by the existing finalizer set in the metadata.finalizers and the resource-specific default policy. Acceptable values are: 'Orphan' - orphan the dependents; 'Background' - allow the garbage collector to delete the dependents in the background; 'Foreground' - a cascading policy that deletes all dependents in the foreground. (optional)
    resourceVersion := "resourceVersion_example" // string | resourceVersion sets a constraint on what resource versions a request may be served from. See https://kubernetes.io/docs/reference/using-api/api-concepts/#resource-versions for details.  Defaults to unset (optional)
    resourceVersionMatch := "resourceVersionMatch_example" // string | resourceVersionMatch determines how resourceVersion is applied to list calls. It is highly recommended that resourceVersionMatch be set for list calls where resourceVersion is set See https://kubernetes.io/docs/reference/using-api/api-concepts/#resource-versions for details.  Defaults to unset (optional)
    timeoutSeconds := int32(56) // int32 | Timeout for the list/watch call. This limits the duration of the call, regardless of any activity or inactivity. (optional)
    body := *openapiclient.NewV1DeleteOptions() // V1DeleteOptions |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.AuthorizationStreamnativeIoV1alpha1Api.DeleteAuthorizationStreamnativeIoV1alpha1CollectionNamespacedIamAccount(context.Background(), namespace).Pretty(pretty).Continue_(continue_).DryRun(dryRun).FieldSelector(fieldSelector).GracePeriodSeconds(gracePeriodSeconds).LabelSelector(labelSelector).Limit(limit).OrphanDependents(orphanDependents).PropagationPolicy(propagationPolicy).ResourceVersion(resourceVersion).ResourceVersionMatch(resourceVersionMatch).TimeoutSeconds(timeoutSeconds).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `AuthorizationStreamnativeIoV1alpha1Api.DeleteAuthorizationStreamnativeIoV1alpha1CollectionNamespacedIamAccount``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DeleteAuthorizationStreamnativeIoV1alpha1CollectionNamespacedIamAccount`: V1Status
    fmt.Fprintf(os.Stdout, "Response from `AuthorizationStreamnativeIoV1alpha1Api.DeleteAuthorizationStreamnativeIoV1alpha1CollectionNamespacedIamAccount`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**namespace** | **string** | object name and auth scope, such as for teams and projects | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteAuthorizationStreamnativeIoV1alpha1CollectionNamespacedIamAccountRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **pretty** | **string** | If &#39;true&#39;, then the output is pretty printed. | 
 **continue_** | **string** | The continue option should be set when retrieving more results from the server. Since this value is server defined, clients may only use the continue value from a previous query result with identical query parameters (except for the value of continue) and the server may reject a continue value it does not recognize. If the specified continue value is no longer valid whether due to expiration (generally five to fifteen minutes) or a configuration change on the server, the server will respond with a 410 ResourceExpired error together with a continue token. If the client needs a consistent list, it must restart their list without the continue field. Otherwise, the client may send another list request with the token received with the 410 error, the server will respond with a list starting from the next key, but from the latest snapshot, which is inconsistent from the previous list results - objects that are created, modified, or deleted after the first list request will be included in the response, as long as their keys are after the \&quot;next key\&quot;.  This field is not supported when watch is true. Clients may start a watch from the last resourceVersion value returned by the server and not miss any modifications. | 
 **dryRun** | **string** | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed | 
 **fieldSelector** | **string** | A selector to restrict the list of returned objects by their fields. Defaults to everything. | 
 **gracePeriodSeconds** | **int32** | The duration in seconds before the object should be deleted. Value must be non-negative integer. The value zero indicates delete immediately. If this value is nil, the default grace period for the specified type will be used. Defaults to a per object value if not specified. zero means delete immediately. | 
 **labelSelector** | **string** | A selector to restrict the list of returned objects by their labels. Defaults to everything. | 
 **limit** | **int32** | limit is a maximum number of responses to return for a list call. If more items exist, the server will set the &#x60;continue&#x60; field on the list metadata to a value that can be used with the same initial query to retrieve the next set of results. Setting a limit may return fewer than the requested amount of items (up to zero items) in the event all requested objects are filtered out and clients should only use the presence of the continue field to determine whether more results are available. Servers may choose not to support the limit argument and will return all of the available results. If limit is specified and the continue field is empty, clients may assume that no more results are available. This field is not supported if watch is true.  The server guarantees that the objects returned when using continue will be identical to issuing a single list call without a limit - that is, no objects created, modified, or deleted after the first request is issued will be included in any subsequent continued requests. This is sometimes referred to as a consistent snapshot, and ensures that a client that is using limit to receive smaller chunks of a very large result can ensure they see all possible objects. If objects are updated during a chunked list the version of the object that was present at the time the first list result was calculated is returned. | 
 **orphanDependents** | **bool** | Deprecated: please use the PropagationPolicy, this field will be deprecated in 1.7. Should the dependent objects be orphaned. If true/false, the \&quot;orphan\&quot; finalizer will be added to/removed from the object&#39;s finalizers list. Either this field or PropagationPolicy may be set, but not both. | 
 **propagationPolicy** | **string** | Whether and how garbage collection will be performed. Either this field or OrphanDependents may be set, but not both. The default policy is decided by the existing finalizer set in the metadata.finalizers and the resource-specific default policy. Acceptable values are: &#39;Orphan&#39; - orphan the dependents; &#39;Background&#39; - allow the garbage collector to delete the dependents in the background; &#39;Foreground&#39; - a cascading policy that deletes all dependents in the foreground. | 
 **resourceVersion** | **string** | resourceVersion sets a constraint on what resource versions a request may be served from. See https://kubernetes.io/docs/reference/using-api/api-concepts/#resource-versions for details.  Defaults to unset | 
 **resourceVersionMatch** | **string** | resourceVersionMatch determines how resourceVersion is applied to list calls. It is highly recommended that resourceVersionMatch be set for list calls where resourceVersion is set See https://kubernetes.io/docs/reference/using-api/api-concepts/#resource-versions for details.  Defaults to unset | 
 **timeoutSeconds** | **int32** | Timeout for the list/watch call. This limits the duration of the call, regardless of any activity or inactivity. | 
 **body** | [**V1DeleteOptions**](V1DeleteOptions.md) |  | 

### Return type

[**V1Status**](V1Status.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/yaml, application/vnd.kubernetes.protobuf

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount

> V1Status DeleteAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount(ctx, name, namespace).Pretty(pretty).DryRun(dryRun).GracePeriodSeconds(gracePeriodSeconds).OrphanDependents(orphanDependents).PropagationPolicy(propagationPolicy).Body(body).Execute()





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
    name := "name_example" // string | name of the IamAccount
    namespace := "namespace_example" // string | object name and auth scope, such as for teams and projects
    pretty := "pretty_example" // string | If 'true', then the output is pretty printed. (optional)
    dryRun := "dryRun_example" // string | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed (optional)
    gracePeriodSeconds := int32(56) // int32 | The duration in seconds before the object should be deleted. Value must be non-negative integer. The value zero indicates delete immediately. If this value is nil, the default grace period for the specified type will be used. Defaults to a per object value if not specified. zero means delete immediately. (optional)
    orphanDependents := true // bool | Deprecated: please use the PropagationPolicy, this field will be deprecated in 1.7. Should the dependent objects be orphaned. If true/false, the \"orphan\" finalizer will be added to/removed from the object's finalizers list. Either this field or PropagationPolicy may be set, but not both. (optional)
    propagationPolicy := "propagationPolicy_example" // string | Whether and how garbage collection will be performed. Either this field or OrphanDependents may be set, but not both. The default policy is decided by the existing finalizer set in the metadata.finalizers and the resource-specific default policy. Acceptable values are: 'Orphan' - orphan the dependents; 'Background' - allow the garbage collector to delete the dependents in the background; 'Foreground' - a cascading policy that deletes all dependents in the foreground. (optional)
    body := *openapiclient.NewV1DeleteOptions() // V1DeleteOptions |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.AuthorizationStreamnativeIoV1alpha1Api.DeleteAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount(context.Background(), name, namespace).Pretty(pretty).DryRun(dryRun).GracePeriodSeconds(gracePeriodSeconds).OrphanDependents(orphanDependents).PropagationPolicy(propagationPolicy).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `AuthorizationStreamnativeIoV1alpha1Api.DeleteAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DeleteAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount`: V1Status
    fmt.Fprintf(os.Stdout, "Response from `AuthorizationStreamnativeIoV1alpha1Api.DeleteAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**name** | **string** | name of the IamAccount | 
**namespace** | **string** | object name and auth scope, such as for teams and projects | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **pretty** | **string** | If &#39;true&#39;, then the output is pretty printed. | 
 **dryRun** | **string** | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed | 
 **gracePeriodSeconds** | **int32** | The duration in seconds before the object should be deleted. Value must be non-negative integer. The value zero indicates delete immediately. If this value is nil, the default grace period for the specified type will be used. Defaults to a per object value if not specified. zero means delete immediately. | 
 **orphanDependents** | **bool** | Deprecated: please use the PropagationPolicy, this field will be deprecated in 1.7. Should the dependent objects be orphaned. If true/false, the \&quot;orphan\&quot; finalizer will be added to/removed from the object&#39;s finalizers list. Either this field or PropagationPolicy may be set, but not both. | 
 **propagationPolicy** | **string** | Whether and how garbage collection will be performed. Either this field or OrphanDependents may be set, but not both. The default policy is decided by the existing finalizer set in the metadata.finalizers and the resource-specific default policy. Acceptable values are: &#39;Orphan&#39; - orphan the dependents; &#39;Background&#39; - allow the garbage collector to delete the dependents in the background; &#39;Foreground&#39; - a cascading policy that deletes all dependents in the foreground. | 
 **body** | [**V1DeleteOptions**](V1DeleteOptions.md) |  | 

### Return type

[**V1Status**](V1Status.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/yaml, application/vnd.kubernetes.protobuf

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus

> V1Status DeleteAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus(ctx, name, namespace).Pretty(pretty).DryRun(dryRun).GracePeriodSeconds(gracePeriodSeconds).OrphanDependents(orphanDependents).PropagationPolicy(propagationPolicy).Body(body).Execute()





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
    name := "name_example" // string | name of the IamAccount
    namespace := "namespace_example" // string | object name and auth scope, such as for teams and projects
    pretty := "pretty_example" // string | If 'true', then the output is pretty printed. (optional)
    dryRun := "dryRun_example" // string | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed (optional)
    gracePeriodSeconds := int32(56) // int32 | The duration in seconds before the object should be deleted. Value must be non-negative integer. The value zero indicates delete immediately. If this value is nil, the default grace period for the specified type will be used. Defaults to a per object value if not specified. zero means delete immediately. (optional)
    orphanDependents := true // bool | Deprecated: please use the PropagationPolicy, this field will be deprecated in 1.7. Should the dependent objects be orphaned. If true/false, the \"orphan\" finalizer will be added to/removed from the object's finalizers list. Either this field or PropagationPolicy may be set, but not both. (optional)
    propagationPolicy := "propagationPolicy_example" // string | Whether and how garbage collection will be performed. Either this field or OrphanDependents may be set, but not both. The default policy is decided by the existing finalizer set in the metadata.finalizers and the resource-specific default policy. Acceptable values are: 'Orphan' - orphan the dependents; 'Background' - allow the garbage collector to delete the dependents in the background; 'Foreground' - a cascading policy that deletes all dependents in the foreground. (optional)
    body := *openapiclient.NewV1DeleteOptions() // V1DeleteOptions |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.AuthorizationStreamnativeIoV1alpha1Api.DeleteAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus(context.Background(), name, namespace).Pretty(pretty).DryRun(dryRun).GracePeriodSeconds(gracePeriodSeconds).OrphanDependents(orphanDependents).PropagationPolicy(propagationPolicy).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `AuthorizationStreamnativeIoV1alpha1Api.DeleteAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DeleteAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus`: V1Status
    fmt.Fprintf(os.Stdout, "Response from `AuthorizationStreamnativeIoV1alpha1Api.DeleteAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**name** | **string** | name of the IamAccount | 
**namespace** | **string** | object name and auth scope, such as for teams and projects | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatusRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **pretty** | **string** | If &#39;true&#39;, then the output is pretty printed. | 
 **dryRun** | **string** | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed | 
 **gracePeriodSeconds** | **int32** | The duration in seconds before the object should be deleted. Value must be non-negative integer. The value zero indicates delete immediately. If this value is nil, the default grace period for the specified type will be used. Defaults to a per object value if not specified. zero means delete immediately. | 
 **orphanDependents** | **bool** | Deprecated: please use the PropagationPolicy, this field will be deprecated in 1.7. Should the dependent objects be orphaned. If true/false, the \&quot;orphan\&quot; finalizer will be added to/removed from the object&#39;s finalizers list. Either this field or PropagationPolicy may be set, but not both. | 
 **propagationPolicy** | **string** | Whether and how garbage collection will be performed. Either this field or OrphanDependents may be set, but not both. The default policy is decided by the existing finalizer set in the metadata.finalizers and the resource-specific default policy. Acceptable values are: &#39;Orphan&#39; - orphan the dependents; &#39;Background&#39; - allow the garbage collector to delete the dependents in the background; &#39;Foreground&#39; - a cascading policy that deletes all dependents in the foreground. | 
 **body** | [**V1DeleteOptions**](V1DeleteOptions.md) |  | 

### Return type

[**V1Status**](V1Status.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/yaml, application/vnd.kubernetes.protobuf

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetAuthorizationStreamnativeIoV1alpha1APIResources

> V1APIResourceList GetAuthorizationStreamnativeIoV1alpha1APIResources(ctx).Execute()





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

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.AuthorizationStreamnativeIoV1alpha1Api.GetAuthorizationStreamnativeIoV1alpha1APIResources(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `AuthorizationStreamnativeIoV1alpha1Api.GetAuthorizationStreamnativeIoV1alpha1APIResources``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetAuthorizationStreamnativeIoV1alpha1APIResources`: V1APIResourceList
    fmt.Fprintf(os.Stdout, "Response from `AuthorizationStreamnativeIoV1alpha1Api.GetAuthorizationStreamnativeIoV1alpha1APIResources`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetAuthorizationStreamnativeIoV1alpha1APIResourcesRequest struct via the builder pattern


### Return type

[**V1APIResourceList**](V1APIResourceList.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/yaml, application/vnd.kubernetes.protobuf

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListAuthorizationStreamnativeIoV1alpha1IamAccountForAllNamespaces

> ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountList ListAuthorizationStreamnativeIoV1alpha1IamAccountForAllNamespaces(ctx).AllowWatchBookmarks(allowWatchBookmarks).Continue_(continue_).FieldSelector(fieldSelector).LabelSelector(labelSelector).Limit(limit).Pretty(pretty).ResourceVersion(resourceVersion).ResourceVersionMatch(resourceVersionMatch).TimeoutSeconds(timeoutSeconds).Watch(watch).Execute()





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
    allowWatchBookmarks := true // bool | allowWatchBookmarks requests watch events with type \"BOOKMARK\". Servers that do not implement bookmarks may ignore this flag and bookmarks are sent at the server's discretion. Clients should not assume bookmarks are returned at any specific interval, nor may they assume the server will send any BOOKMARK event during a session. If this is not a watch, this field is ignored. (optional)
    continue_ := "continue__example" // string | The continue option should be set when retrieving more results from the server. Since this value is server defined, clients may only use the continue value from a previous query result with identical query parameters (except for the value of continue) and the server may reject a continue value it does not recognize. If the specified continue value is no longer valid whether due to expiration (generally five to fifteen minutes) or a configuration change on the server, the server will respond with a 410 ResourceExpired error together with a continue token. If the client needs a consistent list, it must restart their list without the continue field. Otherwise, the client may send another list request with the token received with the 410 error, the server will respond with a list starting from the next key, but from the latest snapshot, which is inconsistent from the previous list results - objects that are created, modified, or deleted after the first list request will be included in the response, as long as their keys are after the \"next key\".  This field is not supported when watch is true. Clients may start a watch from the last resourceVersion value returned by the server and not miss any modifications. (optional)
    fieldSelector := "fieldSelector_example" // string | A selector to restrict the list of returned objects by their fields. Defaults to everything. (optional)
    labelSelector := "labelSelector_example" // string | A selector to restrict the list of returned objects by their labels. Defaults to everything. (optional)
    limit := int32(56) // int32 | limit is a maximum number of responses to return for a list call. If more items exist, the server will set the `continue` field on the list metadata to a value that can be used with the same initial query to retrieve the next set of results. Setting a limit may return fewer than the requested amount of items (up to zero items) in the event all requested objects are filtered out and clients should only use the presence of the continue field to determine whether more results are available. Servers may choose not to support the limit argument and will return all of the available results. If limit is specified and the continue field is empty, clients may assume that no more results are available. This field is not supported if watch is true.  The server guarantees that the objects returned when using continue will be identical to issuing a single list call without a limit - that is, no objects created, modified, or deleted after the first request is issued will be included in any subsequent continued requests. This is sometimes referred to as a consistent snapshot, and ensures that a client that is using limit to receive smaller chunks of a very large result can ensure they see all possible objects. If objects are updated during a chunked list the version of the object that was present at the time the first list result was calculated is returned. (optional)
    pretty := "pretty_example" // string | If 'true', then the output is pretty printed. (optional)
    resourceVersion := "resourceVersion_example" // string | resourceVersion sets a constraint on what resource versions a request may be served from. See https://kubernetes.io/docs/reference/using-api/api-concepts/#resource-versions for details.  Defaults to unset (optional)
    resourceVersionMatch := "resourceVersionMatch_example" // string | resourceVersionMatch determines how resourceVersion is applied to list calls. It is highly recommended that resourceVersionMatch be set for list calls where resourceVersion is set See https://kubernetes.io/docs/reference/using-api/api-concepts/#resource-versions for details.  Defaults to unset (optional)
    timeoutSeconds := int32(56) // int32 | Timeout for the list/watch call. This limits the duration of the call, regardless of any activity or inactivity. (optional)
    watch := true // bool | Watch for changes to the described resources and return them as a stream of add, update, and remove notifications. Specify resourceVersion. (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.AuthorizationStreamnativeIoV1alpha1Api.ListAuthorizationStreamnativeIoV1alpha1IamAccountForAllNamespaces(context.Background()).AllowWatchBookmarks(allowWatchBookmarks).Continue_(continue_).FieldSelector(fieldSelector).LabelSelector(labelSelector).Limit(limit).Pretty(pretty).ResourceVersion(resourceVersion).ResourceVersionMatch(resourceVersionMatch).TimeoutSeconds(timeoutSeconds).Watch(watch).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `AuthorizationStreamnativeIoV1alpha1Api.ListAuthorizationStreamnativeIoV1alpha1IamAccountForAllNamespaces``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListAuthorizationStreamnativeIoV1alpha1IamAccountForAllNamespaces`: ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountList
    fmt.Fprintf(os.Stdout, "Response from `AuthorizationStreamnativeIoV1alpha1Api.ListAuthorizationStreamnativeIoV1alpha1IamAccountForAllNamespaces`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListAuthorizationStreamnativeIoV1alpha1IamAccountForAllNamespacesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **allowWatchBookmarks** | **bool** | allowWatchBookmarks requests watch events with type \&quot;BOOKMARK\&quot;. Servers that do not implement bookmarks may ignore this flag and bookmarks are sent at the server&#39;s discretion. Clients should not assume bookmarks are returned at any specific interval, nor may they assume the server will send any BOOKMARK event during a session. If this is not a watch, this field is ignored. | 
 **continue_** | **string** | The continue option should be set when retrieving more results from the server. Since this value is server defined, clients may only use the continue value from a previous query result with identical query parameters (except for the value of continue) and the server may reject a continue value it does not recognize. If the specified continue value is no longer valid whether due to expiration (generally five to fifteen minutes) or a configuration change on the server, the server will respond with a 410 ResourceExpired error together with a continue token. If the client needs a consistent list, it must restart their list without the continue field. Otherwise, the client may send another list request with the token received with the 410 error, the server will respond with a list starting from the next key, but from the latest snapshot, which is inconsistent from the previous list results - objects that are created, modified, or deleted after the first list request will be included in the response, as long as their keys are after the \&quot;next key\&quot;.  This field is not supported when watch is true. Clients may start a watch from the last resourceVersion value returned by the server and not miss any modifications. | 
 **fieldSelector** | **string** | A selector to restrict the list of returned objects by their fields. Defaults to everything. | 
 **labelSelector** | **string** | A selector to restrict the list of returned objects by their labels. Defaults to everything. | 
 **limit** | **int32** | limit is a maximum number of responses to return for a list call. If more items exist, the server will set the &#x60;continue&#x60; field on the list metadata to a value that can be used with the same initial query to retrieve the next set of results. Setting a limit may return fewer than the requested amount of items (up to zero items) in the event all requested objects are filtered out and clients should only use the presence of the continue field to determine whether more results are available. Servers may choose not to support the limit argument and will return all of the available results. If limit is specified and the continue field is empty, clients may assume that no more results are available. This field is not supported if watch is true.  The server guarantees that the objects returned when using continue will be identical to issuing a single list call without a limit - that is, no objects created, modified, or deleted after the first request is issued will be included in any subsequent continued requests. This is sometimes referred to as a consistent snapshot, and ensures that a client that is using limit to receive smaller chunks of a very large result can ensure they see all possible objects. If objects are updated during a chunked list the version of the object that was present at the time the first list result was calculated is returned. | 
 **pretty** | **string** | If &#39;true&#39;, then the output is pretty printed. | 
 **resourceVersion** | **string** | resourceVersion sets a constraint on what resource versions a request may be served from. See https://kubernetes.io/docs/reference/using-api/api-concepts/#resource-versions for details.  Defaults to unset | 
 **resourceVersionMatch** | **string** | resourceVersionMatch determines how resourceVersion is applied to list calls. It is highly recommended that resourceVersionMatch be set for list calls where resourceVersion is set See https://kubernetes.io/docs/reference/using-api/api-concepts/#resource-versions for details.  Defaults to unset | 
 **timeoutSeconds** | **int32** | Timeout for the list/watch call. This limits the duration of the call, regardless of any activity or inactivity. | 
 **watch** | **bool** | Watch for changes to the described resources and return them as a stream of add, update, and remove notifications. Specify resourceVersion. | 

### Return type

[**ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountList**](ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountList.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/yaml, application/vnd.kubernetes.protobuf, application/json;stream=watch, application/vnd.kubernetes.protobuf;stream=watch

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount

> ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountList ListAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount(ctx, namespace).Pretty(pretty).AllowWatchBookmarks(allowWatchBookmarks).Continue_(continue_).FieldSelector(fieldSelector).LabelSelector(labelSelector).Limit(limit).ResourceVersion(resourceVersion).ResourceVersionMatch(resourceVersionMatch).TimeoutSeconds(timeoutSeconds).Watch(watch).Execute()





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
    namespace := "namespace_example" // string | object name and auth scope, such as for teams and projects
    pretty := "pretty_example" // string | If 'true', then the output is pretty printed. (optional)
    allowWatchBookmarks := true // bool | allowWatchBookmarks requests watch events with type \"BOOKMARK\". Servers that do not implement bookmarks may ignore this flag and bookmarks are sent at the server's discretion. Clients should not assume bookmarks are returned at any specific interval, nor may they assume the server will send any BOOKMARK event during a session. If this is not a watch, this field is ignored. (optional)
    continue_ := "continue__example" // string | The continue option should be set when retrieving more results from the server. Since this value is server defined, clients may only use the continue value from a previous query result with identical query parameters (except for the value of continue) and the server may reject a continue value it does not recognize. If the specified continue value is no longer valid whether due to expiration (generally five to fifteen minutes) or a configuration change on the server, the server will respond with a 410 ResourceExpired error together with a continue token. If the client needs a consistent list, it must restart their list without the continue field. Otherwise, the client may send another list request with the token received with the 410 error, the server will respond with a list starting from the next key, but from the latest snapshot, which is inconsistent from the previous list results - objects that are created, modified, or deleted after the first list request will be included in the response, as long as their keys are after the \"next key\".  This field is not supported when watch is true. Clients may start a watch from the last resourceVersion value returned by the server and not miss any modifications. (optional)
    fieldSelector := "fieldSelector_example" // string | A selector to restrict the list of returned objects by their fields. Defaults to everything. (optional)
    labelSelector := "labelSelector_example" // string | A selector to restrict the list of returned objects by their labels. Defaults to everything. (optional)
    limit := int32(56) // int32 | limit is a maximum number of responses to return for a list call. If more items exist, the server will set the `continue` field on the list metadata to a value that can be used with the same initial query to retrieve the next set of results. Setting a limit may return fewer than the requested amount of items (up to zero items) in the event all requested objects are filtered out and clients should only use the presence of the continue field to determine whether more results are available. Servers may choose not to support the limit argument and will return all of the available results. If limit is specified and the continue field is empty, clients may assume that no more results are available. This field is not supported if watch is true.  The server guarantees that the objects returned when using continue will be identical to issuing a single list call without a limit - that is, no objects created, modified, or deleted after the first request is issued will be included in any subsequent continued requests. This is sometimes referred to as a consistent snapshot, and ensures that a client that is using limit to receive smaller chunks of a very large result can ensure they see all possible objects. If objects are updated during a chunked list the version of the object that was present at the time the first list result was calculated is returned. (optional)
    resourceVersion := "resourceVersion_example" // string | resourceVersion sets a constraint on what resource versions a request may be served from. See https://kubernetes.io/docs/reference/using-api/api-concepts/#resource-versions for details.  Defaults to unset (optional)
    resourceVersionMatch := "resourceVersionMatch_example" // string | resourceVersionMatch determines how resourceVersion is applied to list calls. It is highly recommended that resourceVersionMatch be set for list calls where resourceVersion is set See https://kubernetes.io/docs/reference/using-api/api-concepts/#resource-versions for details.  Defaults to unset (optional)
    timeoutSeconds := int32(56) // int32 | Timeout for the list/watch call. This limits the duration of the call, regardless of any activity or inactivity. (optional)
    watch := true // bool | Watch for changes to the described resources and return them as a stream of add, update, and remove notifications. Specify resourceVersion. (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.AuthorizationStreamnativeIoV1alpha1Api.ListAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount(context.Background(), namespace).Pretty(pretty).AllowWatchBookmarks(allowWatchBookmarks).Continue_(continue_).FieldSelector(fieldSelector).LabelSelector(labelSelector).Limit(limit).ResourceVersion(resourceVersion).ResourceVersionMatch(resourceVersionMatch).TimeoutSeconds(timeoutSeconds).Watch(watch).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `AuthorizationStreamnativeIoV1alpha1Api.ListAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount`: ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountList
    fmt.Fprintf(os.Stdout, "Response from `AuthorizationStreamnativeIoV1alpha1Api.ListAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**namespace** | **string** | object name and auth scope, such as for teams and projects | 

### Other Parameters

Other parameters are passed through a pointer to a apiListAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **pretty** | **string** | If &#39;true&#39;, then the output is pretty printed. | 
 **allowWatchBookmarks** | **bool** | allowWatchBookmarks requests watch events with type \&quot;BOOKMARK\&quot;. Servers that do not implement bookmarks may ignore this flag and bookmarks are sent at the server&#39;s discretion. Clients should not assume bookmarks are returned at any specific interval, nor may they assume the server will send any BOOKMARK event during a session. If this is not a watch, this field is ignored. | 
 **continue_** | **string** | The continue option should be set when retrieving more results from the server. Since this value is server defined, clients may only use the continue value from a previous query result with identical query parameters (except for the value of continue) and the server may reject a continue value it does not recognize. If the specified continue value is no longer valid whether due to expiration (generally five to fifteen minutes) or a configuration change on the server, the server will respond with a 410 ResourceExpired error together with a continue token. If the client needs a consistent list, it must restart their list without the continue field. Otherwise, the client may send another list request with the token received with the 410 error, the server will respond with a list starting from the next key, but from the latest snapshot, which is inconsistent from the previous list results - objects that are created, modified, or deleted after the first list request will be included in the response, as long as their keys are after the \&quot;next key\&quot;.  This field is not supported when watch is true. Clients may start a watch from the last resourceVersion value returned by the server and not miss any modifications. | 
 **fieldSelector** | **string** | A selector to restrict the list of returned objects by their fields. Defaults to everything. | 
 **labelSelector** | **string** | A selector to restrict the list of returned objects by their labels. Defaults to everything. | 
 **limit** | **int32** | limit is a maximum number of responses to return for a list call. If more items exist, the server will set the &#x60;continue&#x60; field on the list metadata to a value that can be used with the same initial query to retrieve the next set of results. Setting a limit may return fewer than the requested amount of items (up to zero items) in the event all requested objects are filtered out and clients should only use the presence of the continue field to determine whether more results are available. Servers may choose not to support the limit argument and will return all of the available results. If limit is specified and the continue field is empty, clients may assume that no more results are available. This field is not supported if watch is true.  The server guarantees that the objects returned when using continue will be identical to issuing a single list call without a limit - that is, no objects created, modified, or deleted after the first request is issued will be included in any subsequent continued requests. This is sometimes referred to as a consistent snapshot, and ensures that a client that is using limit to receive smaller chunks of a very large result can ensure they see all possible objects. If objects are updated during a chunked list the version of the object that was present at the time the first list result was calculated is returned. | 
 **resourceVersion** | **string** | resourceVersion sets a constraint on what resource versions a request may be served from. See https://kubernetes.io/docs/reference/using-api/api-concepts/#resource-versions for details.  Defaults to unset | 
 **resourceVersionMatch** | **string** | resourceVersionMatch determines how resourceVersion is applied to list calls. It is highly recommended that resourceVersionMatch be set for list calls where resourceVersion is set See https://kubernetes.io/docs/reference/using-api/api-concepts/#resource-versions for details.  Defaults to unset | 
 **timeoutSeconds** | **int32** | Timeout for the list/watch call. This limits the duration of the call, regardless of any activity or inactivity. | 
 **watch** | **bool** | Watch for changes to the described resources and return them as a stream of add, update, and remove notifications. Specify resourceVersion. | 

### Return type

[**ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountList**](ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccountList.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/yaml, application/vnd.kubernetes.protobuf, application/json;stream=watch, application/vnd.kubernetes.protobuf;stream=watch

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PatchAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount

> ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount PatchAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount(ctx, name, namespace).Body(body).Pretty(pretty).DryRun(dryRun).FieldManager(fieldManager).Force(force).Execute()





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
    name := "name_example" // string | name of the IamAccount
    namespace := "namespace_example" // string | object name and auth scope, such as for teams and projects
    body := map[string]interface{}{ ... } // map[string]interface{} | 
    pretty := "pretty_example" // string | If 'true', then the output is pretty printed. (optional)
    dryRun := "dryRun_example" // string | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed (optional)
    fieldManager := "fieldManager_example" // string | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. This field is required for apply requests (application/apply-patch) but optional for non-apply patch types (JsonPatch, MergePatch, StrategicMergePatch). (optional)
    force := true // bool | Force is going to \"force\" Apply requests. It means user will re-acquire conflicting fields owned by other people. Force flag must be unset for non-apply patch requests. (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.AuthorizationStreamnativeIoV1alpha1Api.PatchAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount(context.Background(), name, namespace).Body(body).Pretty(pretty).DryRun(dryRun).FieldManager(fieldManager).Force(force).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `AuthorizationStreamnativeIoV1alpha1Api.PatchAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `PatchAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount`: ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount
    fmt.Fprintf(os.Stdout, "Response from `AuthorizationStreamnativeIoV1alpha1Api.PatchAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**name** | **string** | name of the IamAccount | 
**namespace** | **string** | object name and auth scope, such as for teams and projects | 

### Other Parameters

Other parameters are passed through a pointer to a apiPatchAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **body** | **map[string]interface{}** |  | 
 **pretty** | **string** | If &#39;true&#39;, then the output is pretty printed. | 
 **dryRun** | **string** | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed | 
 **fieldManager** | **string** | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. This field is required for apply requests (application/apply-patch) but optional for non-apply patch types (JsonPatch, MergePatch, StrategicMergePatch). | 
 **force** | **bool** | Force is going to \&quot;force\&quot; Apply requests. It means user will re-acquire conflicting fields owned by other people. Force flag must be unset for non-apply patch requests. | 

### Return type

[**ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount**](ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json-patch+json, application/merge-patch+json, application/strategic-merge-patch+json, application/apply-patch+yaml
- **Accept**: application/json, application/yaml, application/vnd.kubernetes.protobuf

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PatchAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus

> ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount PatchAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus(ctx, name, namespace).Body(body).Pretty(pretty).DryRun(dryRun).FieldManager(fieldManager).Force(force).Execute()





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
    name := "name_example" // string | name of the IamAccount
    namespace := "namespace_example" // string | object name and auth scope, such as for teams and projects
    body := map[string]interface{}{ ... } // map[string]interface{} | 
    pretty := "pretty_example" // string | If 'true', then the output is pretty printed. (optional)
    dryRun := "dryRun_example" // string | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed (optional)
    fieldManager := "fieldManager_example" // string | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. This field is required for apply requests (application/apply-patch) but optional for non-apply patch types (JsonPatch, MergePatch, StrategicMergePatch). (optional)
    force := true // bool | Force is going to \"force\" Apply requests. It means user will re-acquire conflicting fields owned by other people. Force flag must be unset for non-apply patch requests. (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.AuthorizationStreamnativeIoV1alpha1Api.PatchAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus(context.Background(), name, namespace).Body(body).Pretty(pretty).DryRun(dryRun).FieldManager(fieldManager).Force(force).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `AuthorizationStreamnativeIoV1alpha1Api.PatchAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `PatchAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus`: ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount
    fmt.Fprintf(os.Stdout, "Response from `AuthorizationStreamnativeIoV1alpha1Api.PatchAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**name** | **string** | name of the IamAccount | 
**namespace** | **string** | object name and auth scope, such as for teams and projects | 

### Other Parameters

Other parameters are passed through a pointer to a apiPatchAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatusRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **body** | **map[string]interface{}** |  | 
 **pretty** | **string** | If &#39;true&#39;, then the output is pretty printed. | 
 **dryRun** | **string** | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed | 
 **fieldManager** | **string** | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. This field is required for apply requests (application/apply-patch) but optional for non-apply patch types (JsonPatch, MergePatch, StrategicMergePatch). | 
 **force** | **bool** | Force is going to \&quot;force\&quot; Apply requests. It means user will re-acquire conflicting fields owned by other people. Force flag must be unset for non-apply patch requests. | 

### Return type

[**ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount**](ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json-patch+json, application/merge-patch+json, application/strategic-merge-patch+json, application/apply-patch+yaml
- **Accept**: application/json, application/yaml, application/vnd.kubernetes.protobuf

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ReadAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount

> ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount ReadAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount(ctx, name, namespace).Pretty(pretty).Execute()





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
    name := "name_example" // string | name of the IamAccount
    namespace := "namespace_example" // string | object name and auth scope, such as for teams and projects
    pretty := "pretty_example" // string | If 'true', then the output is pretty printed. (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.AuthorizationStreamnativeIoV1alpha1Api.ReadAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount(context.Background(), name, namespace).Pretty(pretty).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `AuthorizationStreamnativeIoV1alpha1Api.ReadAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ReadAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount`: ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount
    fmt.Fprintf(os.Stdout, "Response from `AuthorizationStreamnativeIoV1alpha1Api.ReadAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**name** | **string** | name of the IamAccount | 
**namespace** | **string** | object name and auth scope, such as for teams and projects | 

### Other Parameters

Other parameters are passed through a pointer to a apiReadAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **pretty** | **string** | If &#39;true&#39;, then the output is pretty printed. | 

### Return type

[**ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount**](ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/yaml, application/vnd.kubernetes.protobuf

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ReadAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus

> ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount ReadAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus(ctx, name, namespace).Pretty(pretty).Execute()





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
    name := "name_example" // string | name of the IamAccount
    namespace := "namespace_example" // string | object name and auth scope, such as for teams and projects
    pretty := "pretty_example" // string | If 'true', then the output is pretty printed. (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.AuthorizationStreamnativeIoV1alpha1Api.ReadAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus(context.Background(), name, namespace).Pretty(pretty).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `AuthorizationStreamnativeIoV1alpha1Api.ReadAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ReadAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus`: ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount
    fmt.Fprintf(os.Stdout, "Response from `AuthorizationStreamnativeIoV1alpha1Api.ReadAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**name** | **string** | name of the IamAccount | 
**namespace** | **string** | object name and auth scope, such as for teams and projects | 

### Other Parameters

Other parameters are passed through a pointer to a apiReadAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatusRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **pretty** | **string** | If &#39;true&#39;, then the output is pretty printed. | 

### Return type

[**ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount**](ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/yaml, application/vnd.kubernetes.protobuf

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ReplaceAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount

> ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount ReplaceAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount(ctx, name, namespace).Body(body).Pretty(pretty).DryRun(dryRun).FieldManager(fieldManager).Execute()





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
    name := "name_example" // string | name of the IamAccount
    namespace := "namespace_example" // string | object name and auth scope, such as for teams and projects
    body := *openapiclient.NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount() // ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount | 
    pretty := "pretty_example" // string | If 'true', then the output is pretty printed. (optional)
    dryRun := "dryRun_example" // string | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed (optional)
    fieldManager := "fieldManager_example" // string | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.AuthorizationStreamnativeIoV1alpha1Api.ReplaceAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount(context.Background(), name, namespace).Body(body).Pretty(pretty).DryRun(dryRun).FieldManager(fieldManager).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `AuthorizationStreamnativeIoV1alpha1Api.ReplaceAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ReplaceAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount`: ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount
    fmt.Fprintf(os.Stdout, "Response from `AuthorizationStreamnativeIoV1alpha1Api.ReplaceAuthorizationStreamnativeIoV1alpha1NamespacedIamAccount`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**name** | **string** | name of the IamAccount | 
**namespace** | **string** | object name and auth scope, such as for teams and projects | 

### Other Parameters

Other parameters are passed through a pointer to a apiReplaceAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **body** | [**ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount**](ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount.md) |  | 
 **pretty** | **string** | If &#39;true&#39;, then the output is pretty printed. | 
 **dryRun** | **string** | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed | 
 **fieldManager** | **string** | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. | 

### Return type

[**ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount**](ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/yaml, application/vnd.kubernetes.protobuf

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ReplaceAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus

> ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount ReplaceAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus(ctx, name, namespace).Body(body).Pretty(pretty).DryRun(dryRun).FieldManager(fieldManager).Execute()





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
    name := "name_example" // string | name of the IamAccount
    namespace := "namespace_example" // string | object name and auth scope, such as for teams and projects
    body := *openapiclient.NewComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount() // ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount | 
    pretty := "pretty_example" // string | If 'true', then the output is pretty printed. (optional)
    dryRun := "dryRun_example" // string | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed (optional)
    fieldManager := "fieldManager_example" // string | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.AuthorizationStreamnativeIoV1alpha1Api.ReplaceAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus(context.Background(), name, namespace).Body(body).Pretty(pretty).DryRun(dryRun).FieldManager(fieldManager).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `AuthorizationStreamnativeIoV1alpha1Api.ReplaceAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ReplaceAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus`: ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount
    fmt.Fprintf(os.Stdout, "Response from `AuthorizationStreamnativeIoV1alpha1Api.ReplaceAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**name** | **string** | name of the IamAccount | 
**namespace** | **string** | object name and auth scope, such as for teams and projects | 

### Other Parameters

Other parameters are passed through a pointer to a apiReplaceAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatusRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **body** | [**ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount**](ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount.md) |  | 
 **pretty** | **string** | If &#39;true&#39;, then the output is pretty printed. | 
 **dryRun** | **string** | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed | 
 **fieldManager** | **string** | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. | 

### Return type

[**ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount**](ComGithubStreamnativeCloudApiServerPkgApisAuthorizationV1alpha1IamAccount.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/yaml, application/vnd.kubernetes.protobuf

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## WatchAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus

> V1WatchEvent WatchAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus(ctx, name, namespace).AllowWatchBookmarks(allowWatchBookmarks).Continue_(continue_).FieldSelector(fieldSelector).LabelSelector(labelSelector).Limit(limit).Pretty(pretty).ResourceVersion(resourceVersion).ResourceVersionMatch(resourceVersionMatch).TimeoutSeconds(timeoutSeconds).Watch(watch).Execute()





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
    name := "name_example" // string | name of the IamAccount
    namespace := "namespace_example" // string | object name and auth scope, such as for teams and projects
    allowWatchBookmarks := true // bool | allowWatchBookmarks requests watch events with type \"BOOKMARK\". Servers that do not implement bookmarks may ignore this flag and bookmarks are sent at the server's discretion. Clients should not assume bookmarks are returned at any specific interval, nor may they assume the server will send any BOOKMARK event during a session. If this is not a watch, this field is ignored. (optional)
    continue_ := "continue__example" // string | The continue option should be set when retrieving more results from the server. Since this value is server defined, clients may only use the continue value from a previous query result with identical query parameters (except for the value of continue) and the server may reject a continue value it does not recognize. If the specified continue value is no longer valid whether due to expiration (generally five to fifteen minutes) or a configuration change on the server, the server will respond with a 410 ResourceExpired error together with a continue token. If the client needs a consistent list, it must restart their list without the continue field. Otherwise, the client may send another list request with the token received with the 410 error, the server will respond with a list starting from the next key, but from the latest snapshot, which is inconsistent from the previous list results - objects that are created, modified, or deleted after the first list request will be included in the response, as long as their keys are after the \"next key\".  This field is not supported when watch is true. Clients may start a watch from the last resourceVersion value returned by the server and not miss any modifications. (optional)
    fieldSelector := "fieldSelector_example" // string | A selector to restrict the list of returned objects by their fields. Defaults to everything. (optional)
    labelSelector := "labelSelector_example" // string | A selector to restrict the list of returned objects by their labels. Defaults to everything. (optional)
    limit := int32(56) // int32 | limit is a maximum number of responses to return for a list call. If more items exist, the server will set the `continue` field on the list metadata to a value that can be used with the same initial query to retrieve the next set of results. Setting a limit may return fewer than the requested amount of items (up to zero items) in the event all requested objects are filtered out and clients should only use the presence of the continue field to determine whether more results are available. Servers may choose not to support the limit argument and will return all of the available results. If limit is specified and the continue field is empty, clients may assume that no more results are available. This field is not supported if watch is true.  The server guarantees that the objects returned when using continue will be identical to issuing a single list call without a limit - that is, no objects created, modified, or deleted after the first request is issued will be included in any subsequent continued requests. This is sometimes referred to as a consistent snapshot, and ensures that a client that is using limit to receive smaller chunks of a very large result can ensure they see all possible objects. If objects are updated during a chunked list the version of the object that was present at the time the first list result was calculated is returned. (optional)
    pretty := "pretty_example" // string | If 'true', then the output is pretty printed. (optional)
    resourceVersion := "resourceVersion_example" // string | resourceVersion sets a constraint on what resource versions a request may be served from. See https://kubernetes.io/docs/reference/using-api/api-concepts/#resource-versions for details.  Defaults to unset (optional)
    resourceVersionMatch := "resourceVersionMatch_example" // string | resourceVersionMatch determines how resourceVersion is applied to list calls. It is highly recommended that resourceVersionMatch be set for list calls where resourceVersion is set See https://kubernetes.io/docs/reference/using-api/api-concepts/#resource-versions for details.  Defaults to unset (optional)
    timeoutSeconds := int32(56) // int32 | Timeout for the list/watch call. This limits the duration of the call, regardless of any activity or inactivity. (optional)
    watch := true // bool | Watch for changes to the described resources and return them as a stream of add, update, and remove notifications. Specify resourceVersion. (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.AuthorizationStreamnativeIoV1alpha1Api.WatchAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus(context.Background(), name, namespace).AllowWatchBookmarks(allowWatchBookmarks).Continue_(continue_).FieldSelector(fieldSelector).LabelSelector(labelSelector).Limit(limit).Pretty(pretty).ResourceVersion(resourceVersion).ResourceVersionMatch(resourceVersionMatch).TimeoutSeconds(timeoutSeconds).Watch(watch).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `AuthorizationStreamnativeIoV1alpha1Api.WatchAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `WatchAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus`: V1WatchEvent
    fmt.Fprintf(os.Stdout, "Response from `AuthorizationStreamnativeIoV1alpha1Api.WatchAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatus`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**name** | **string** | name of the IamAccount | 
**namespace** | **string** | object name and auth scope, such as for teams and projects | 

### Other Parameters

Other parameters are passed through a pointer to a apiWatchAuthorizationStreamnativeIoV1alpha1NamespacedIamAccountStatusRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **allowWatchBookmarks** | **bool** | allowWatchBookmarks requests watch events with type \&quot;BOOKMARK\&quot;. Servers that do not implement bookmarks may ignore this flag and bookmarks are sent at the server&#39;s discretion. Clients should not assume bookmarks are returned at any specific interval, nor may they assume the server will send any BOOKMARK event during a session. If this is not a watch, this field is ignored. | 
 **continue_** | **string** | The continue option should be set when retrieving more results from the server. Since this value is server defined, clients may only use the continue value from a previous query result with identical query parameters (except for the value of continue) and the server may reject a continue value it does not recognize. If the specified continue value is no longer valid whether due to expiration (generally five to fifteen minutes) or a configuration change on the server, the server will respond with a 410 ResourceExpired error together with a continue token. If the client needs a consistent list, it must restart their list without the continue field. Otherwise, the client may send another list request with the token received with the 410 error, the server will respond with a list starting from the next key, but from the latest snapshot, which is inconsistent from the previous list results - objects that are created, modified, or deleted after the first list request will be included in the response, as long as their keys are after the \&quot;next key\&quot;.  This field is not supported when watch is true. Clients may start a watch from the last resourceVersion value returned by the server and not miss any modifications. | 
 **fieldSelector** | **string** | A selector to restrict the list of returned objects by their fields. Defaults to everything. | 
 **labelSelector** | **string** | A selector to restrict the list of returned objects by their labels. Defaults to everything. | 
 **limit** | **int32** | limit is a maximum number of responses to return for a list call. If more items exist, the server will set the &#x60;continue&#x60; field on the list metadata to a value that can be used with the same initial query to retrieve the next set of results. Setting a limit may return fewer than the requested amount of items (up to zero items) in the event all requested objects are filtered out and clients should only use the presence of the continue field to determine whether more results are available. Servers may choose not to support the limit argument and will return all of the available results. If limit is specified and the continue field is empty, clients may assume that no more results are available. This field is not supported if watch is true.  The server guarantees that the objects returned when using continue will be identical to issuing a single list call without a limit - that is, no objects created, modified, or deleted after the first request is issued will be included in any subsequent continued requests. This is sometimes referred to as a consistent snapshot, and ensures that a client that is using limit to receive smaller chunks of a very large result can ensure they see all possible objects. If objects are updated during a chunked list the version of the object that was present at the time the first list result was calculated is returned. | 
 **pretty** | **string** | If &#39;true&#39;, then the output is pretty printed. | 
 **resourceVersion** | **string** | resourceVersion sets a constraint on what resource versions a request may be served from. See https://kubernetes.io/docs/reference/using-api/api-concepts/#resource-versions for details.  Defaults to unset | 
 **resourceVersionMatch** | **string** | resourceVersionMatch determines how resourceVersion is applied to list calls. It is highly recommended that resourceVersionMatch be set for list calls where resourceVersion is set See https://kubernetes.io/docs/reference/using-api/api-concepts/#resource-versions for details.  Defaults to unset | 
 **timeoutSeconds** | **int32** | Timeout for the list/watch call. This limits the duration of the call, regardless of any activity or inactivity. | 
 **watch** | **bool** | Watch for changes to the described resources and return them as a stream of add, update, and remove notifications. Specify resourceVersion. | 

### Return type

[**V1WatchEvent**](V1WatchEvent.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/yaml, application/vnd.kubernetes.protobuf, application/json;stream=watch, application/vnd.kubernetes.protobuf;stream=watch

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

