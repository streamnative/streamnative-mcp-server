# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MaintenanceWindow

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Recurrence** | **string** | Recurrence define the maintenance execution cycle, 0~6, to express Monday to Sunday | 
**Window** | [**ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Window**](ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Window.md) |  | 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MaintenanceWindow

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MaintenanceWindow(recurrence string, window ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Window, ) *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MaintenanceWindow`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MaintenanceWindow instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MaintenanceWindow object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MaintenanceWindowWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MaintenanceWindowWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MaintenanceWindow`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MaintenanceWindowWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MaintenanceWindow object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetRecurrence

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MaintenanceWindow) GetRecurrence() string`

GetRecurrence returns the Recurrence field if non-nil, zero value otherwise.

### GetRecurrenceOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MaintenanceWindow) GetRecurrenceOk() (*string, bool)`

GetRecurrenceOk returns a tuple with the Recurrence field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRecurrence

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MaintenanceWindow) SetRecurrence(v string)`

SetRecurrence sets Recurrence field to given value.


### GetWindow

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MaintenanceWindow) GetWindow() ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Window`

GetWindow returns the Window field if non-nil, zero value otherwise.

### GetWindowOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MaintenanceWindow) GetWindowOk() (*ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Window, bool)`

GetWindowOk returns a tuple with the Window field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWindow

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1MaintenanceWindow) SetWindow(v ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Window)`

SetWindow sets Window field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


