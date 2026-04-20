# IoK8sApimachineryPkgApisMetaV1ObjectMetaOwnerReferencesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ApiVersion** | **string** | API version of the referent. | [default to ""]
**BlockOwnerDeletion** | Pointer to **bool** | If true, AND if the owner has the \&quot;foregroundDeletion\&quot; finalizer, then the owner cannot be deleted from the key-value store until this reference is removed. See https://kubernetes.io/docs/concepts/architecture/garbage-collection/#foreground-deletion for how the garbage collector interacts with this field and enforces the foreground deletion. Defaults to false. To set this field, a user needs \&quot;delete\&quot; permission of the owner, otherwise 422 (Unprocessable Entity) will be returned. | [optional] 
**Controller** | Pointer to **bool** | If true, this reference points to the managing controller. | [optional] 
**Kind** | **string** | Kind of the referent. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds | [default to ""]
**Name** | **string** | Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names#names | [default to ""]
**Uid** | **string** | UID of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names#uids | [default to ""]

## Methods

### NewIoK8sApimachineryPkgApisMetaV1ObjectMetaOwnerReferencesInner

`func NewIoK8sApimachineryPkgApisMetaV1ObjectMetaOwnerReferencesInner(apiVersion string, kind string, name string, uid string, ) *IoK8sApimachineryPkgApisMetaV1ObjectMetaOwnerReferencesInner`

NewIoK8sApimachineryPkgApisMetaV1ObjectMetaOwnerReferencesInner instantiates a new IoK8sApimachineryPkgApisMetaV1ObjectMetaOwnerReferencesInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewIoK8sApimachineryPkgApisMetaV1ObjectMetaOwnerReferencesInnerWithDefaults

`func NewIoK8sApimachineryPkgApisMetaV1ObjectMetaOwnerReferencesInnerWithDefaults() *IoK8sApimachineryPkgApisMetaV1ObjectMetaOwnerReferencesInner`

NewIoK8sApimachineryPkgApisMetaV1ObjectMetaOwnerReferencesInnerWithDefaults instantiates a new IoK8sApimachineryPkgApisMetaV1ObjectMetaOwnerReferencesInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetApiVersion

`func (o *IoK8sApimachineryPkgApisMetaV1ObjectMetaOwnerReferencesInner) GetApiVersion() string`

GetApiVersion returns the ApiVersion field if non-nil, zero value otherwise.

### GetApiVersionOk

`func (o *IoK8sApimachineryPkgApisMetaV1ObjectMetaOwnerReferencesInner) GetApiVersionOk() (*string, bool)`

GetApiVersionOk returns a tuple with the ApiVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetApiVersion

`func (o *IoK8sApimachineryPkgApisMetaV1ObjectMetaOwnerReferencesInner) SetApiVersion(v string)`

SetApiVersion sets ApiVersion field to given value.


### GetBlockOwnerDeletion

`func (o *IoK8sApimachineryPkgApisMetaV1ObjectMetaOwnerReferencesInner) GetBlockOwnerDeletion() bool`

GetBlockOwnerDeletion returns the BlockOwnerDeletion field if non-nil, zero value otherwise.

### GetBlockOwnerDeletionOk

`func (o *IoK8sApimachineryPkgApisMetaV1ObjectMetaOwnerReferencesInner) GetBlockOwnerDeletionOk() (*bool, bool)`

GetBlockOwnerDeletionOk returns a tuple with the BlockOwnerDeletion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBlockOwnerDeletion

`func (o *IoK8sApimachineryPkgApisMetaV1ObjectMetaOwnerReferencesInner) SetBlockOwnerDeletion(v bool)`

SetBlockOwnerDeletion sets BlockOwnerDeletion field to given value.

### HasBlockOwnerDeletion

`func (o *IoK8sApimachineryPkgApisMetaV1ObjectMetaOwnerReferencesInner) HasBlockOwnerDeletion() bool`

HasBlockOwnerDeletion returns a boolean if a field has been set.

### GetController

`func (o *IoK8sApimachineryPkgApisMetaV1ObjectMetaOwnerReferencesInner) GetController() bool`

GetController returns the Controller field if non-nil, zero value otherwise.

### GetControllerOk

`func (o *IoK8sApimachineryPkgApisMetaV1ObjectMetaOwnerReferencesInner) GetControllerOk() (*bool, bool)`

GetControllerOk returns a tuple with the Controller field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetController

`func (o *IoK8sApimachineryPkgApisMetaV1ObjectMetaOwnerReferencesInner) SetController(v bool)`

SetController sets Controller field to given value.

### HasController

`func (o *IoK8sApimachineryPkgApisMetaV1ObjectMetaOwnerReferencesInner) HasController() bool`

HasController returns a boolean if a field has been set.

### GetKind

`func (o *IoK8sApimachineryPkgApisMetaV1ObjectMetaOwnerReferencesInner) GetKind() string`

GetKind returns the Kind field if non-nil, zero value otherwise.

### GetKindOk

`func (o *IoK8sApimachineryPkgApisMetaV1ObjectMetaOwnerReferencesInner) GetKindOk() (*string, bool)`

GetKindOk returns a tuple with the Kind field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKind

`func (o *IoK8sApimachineryPkgApisMetaV1ObjectMetaOwnerReferencesInner) SetKind(v string)`

SetKind sets Kind field to given value.


### GetName

`func (o *IoK8sApimachineryPkgApisMetaV1ObjectMetaOwnerReferencesInner) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *IoK8sApimachineryPkgApisMetaV1ObjectMetaOwnerReferencesInner) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *IoK8sApimachineryPkgApisMetaV1ObjectMetaOwnerReferencesInner) SetName(v string)`

SetName sets Name field to given value.


### GetUid

`func (o *IoK8sApimachineryPkgApisMetaV1ObjectMetaOwnerReferencesInner) GetUid() string`

GetUid returns the Uid field if non-nil, zero value otherwise.

### GetUidOk

`func (o *IoK8sApimachineryPkgApisMetaV1ObjectMetaOwnerReferencesInner) GetUidOk() (*string, bool)`

GetUidOk returns a tuple with the Uid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUid

`func (o *IoK8sApimachineryPkgApisMetaV1ObjectMetaOwnerReferencesInner) SetUid(v string)`

SetUid sets Uid field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


