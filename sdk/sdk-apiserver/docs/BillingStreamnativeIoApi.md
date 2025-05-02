# \BillingStreamnativeIoApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetBillingStreamnativeIoAPIGroup**](BillingStreamnativeIoApi.md#GetBillingStreamnativeIoAPIGroup) | **Get** /apis/billing.streamnative.io/ | 



## GetBillingStreamnativeIoAPIGroup

> V1APIGroup GetBillingStreamnativeIoAPIGroup(ctx).Execute()





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
    resp, r, err := apiClient.BillingStreamnativeIoApi.GetBillingStreamnativeIoAPIGroup(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `BillingStreamnativeIoApi.GetBillingStreamnativeIoAPIGroup``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetBillingStreamnativeIoAPIGroup`: V1APIGroup
    fmt.Fprintf(os.Stdout, "Response from `BillingStreamnativeIoApi.GetBillingStreamnativeIoAPIGroup`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetBillingStreamnativeIoAPIGroupRequest struct via the builder pattern


### Return type

[**V1APIGroup**](V1APIGroup.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/yaml, application/vnd.kubernetes.protobuf

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

