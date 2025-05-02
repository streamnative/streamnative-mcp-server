# ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookkeeperNodeResource

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Cpu** | **string** | Quantity is a fixed-point representation of a number. It provides convenient marshaling/unmarshaling in JSON and YAML, in addition to String() and AsInt64() accessors.  The serialization format is:  &lt;quantity&gt;        ::&#x3D; &lt;signedNumber&gt;&lt;suffix&gt;   (Note that &lt;suffix&gt; may be empty, from the \&quot;\&quot; case in &lt;decimalSI&gt;.) &lt;digit&gt;           ::&#x3D; 0 | 1 | ... | 9 &lt;digits&gt;          ::&#x3D; &lt;digit&gt; | &lt;digit&gt;&lt;digits&gt; &lt;number&gt;          ::&#x3D; &lt;digits&gt; | &lt;digits&gt;.&lt;digits&gt; | &lt;digits&gt;. | .&lt;digits&gt; &lt;sign&gt;            ::&#x3D; \&quot;+\&quot; | \&quot;-\&quot; &lt;signedNumber&gt;    ::&#x3D; &lt;number&gt; | &lt;sign&gt;&lt;number&gt; &lt;suffix&gt;          ::&#x3D; &lt;binarySI&gt; | &lt;decimalExponent&gt; | &lt;decimalSI&gt; &lt;binarySI&gt;        ::&#x3D; Ki | Mi | Gi | Ti | Pi | Ei   (International System of units; See: http://physics.nist.gov/cuu/Units/binary.html) &lt;decimalSI&gt;       ::&#x3D; m | \&quot;\&quot; | k | M | G | T | P | E   (Note that 1024 &#x3D; 1Ki but 1000 &#x3D; 1k; I didn&#39;t choose the capitalization.) &lt;decimalExponent&gt; ::&#x3D; \&quot;e\&quot; &lt;signedNumber&gt; | \&quot;E\&quot; &lt;signedNumber&gt;  No matter which of the three exponent forms is used, no quantity may represent a number greater than 2^63-1 in magnitude, nor may it have more than 3 decimal places. Numbers larger or more precise will be capped or rounded up. (E.g.: 0.1m will rounded up to 1m.) This may be extended in the future if we require larger or smaller quantities.  When a Quantity is parsed from a string, it will remember the type of suffix it had, and will use the same type again when it is serialized.  Before serializing, Quantity will be put in \&quot;canonical form\&quot;. This means that Exponent/suffix will be adjusted up or down (with a corresponding increase or decrease in Mantissa) such that:   a. No precision is lost   b. No fractional digits will be emitted   c. The exponent (or suffix) is as large as possible. The sign will be omitted unless the number is negative.  Examples:   1.5 will be serialized as \&quot;1500m\&quot;   1.5Gi will be serialized as \&quot;1536Mi\&quot;  Note that the quantity will NEVER be internally represented by a floating point number. That is the whole point of this exercise.  Non-canonical values will still parse as long as they are well formed, but will be re-emitted in their canonical form. (So always use canonical form, or don&#39;t diff.)  This format is intended to make it difficult to use these numbers without writing some sort of special handling code in the hopes that that will cause implementors to also use a fixed point implementation. | 
**DirectPercentage** | Pointer to **int32** | Percentage of direct memory from overall memory. Set to 0 to use default value. | [optional] 
**HeapPercentage** | Pointer to **int32** | Percentage of heap memory from overall memory. Set to 0 to use default value. | [optional] 
**JournalDisk** | Pointer to **string** | JournalDisk size. Set to zero equivalent to use default value | [optional] 
**LedgerDisk** | Pointer to **string** | LedgerDisk size. Set to zero equivalent to use default value | [optional] 
**Memory** | **string** | Quantity is a fixed-point representation of a number. It provides convenient marshaling/unmarshaling in JSON and YAML, in addition to String() and AsInt64() accessors.  The serialization format is:  &lt;quantity&gt;        ::&#x3D; &lt;signedNumber&gt;&lt;suffix&gt;   (Note that &lt;suffix&gt; may be empty, from the \&quot;\&quot; case in &lt;decimalSI&gt;.) &lt;digit&gt;           ::&#x3D; 0 | 1 | ... | 9 &lt;digits&gt;          ::&#x3D; &lt;digit&gt; | &lt;digit&gt;&lt;digits&gt; &lt;number&gt;          ::&#x3D; &lt;digits&gt; | &lt;digits&gt;.&lt;digits&gt; | &lt;digits&gt;. | .&lt;digits&gt; &lt;sign&gt;            ::&#x3D; \&quot;+\&quot; | \&quot;-\&quot; &lt;signedNumber&gt;    ::&#x3D; &lt;number&gt; | &lt;sign&gt;&lt;number&gt; &lt;suffix&gt;          ::&#x3D; &lt;binarySI&gt; | &lt;decimalExponent&gt; | &lt;decimalSI&gt; &lt;binarySI&gt;        ::&#x3D; Ki | Mi | Gi | Ti | Pi | Ei   (International System of units; See: http://physics.nist.gov/cuu/Units/binary.html) &lt;decimalSI&gt;       ::&#x3D; m | \&quot;\&quot; | k | M | G | T | P | E   (Note that 1024 &#x3D; 1Ki but 1000 &#x3D; 1k; I didn&#39;t choose the capitalization.) &lt;decimalExponent&gt; ::&#x3D; \&quot;e\&quot; &lt;signedNumber&gt; | \&quot;E\&quot; &lt;signedNumber&gt;  No matter which of the three exponent forms is used, no quantity may represent a number greater than 2^63-1 in magnitude, nor may it have more than 3 decimal places. Numbers larger or more precise will be capped or rounded up. (E.g.: 0.1m will rounded up to 1m.) This may be extended in the future if we require larger or smaller quantities.  When a Quantity is parsed from a string, it will remember the type of suffix it had, and will use the same type again when it is serialized.  Before serializing, Quantity will be put in \&quot;canonical form\&quot;. This means that Exponent/suffix will be adjusted up or down (with a corresponding increase or decrease in Mantissa) such that:   a. No precision is lost   b. No fractional digits will be emitted   c. The exponent (or suffix) is as large as possible. The sign will be omitted unless the number is negative.  Examples:   1.5 will be serialized as \&quot;1500m\&quot;   1.5Gi will be serialized as \&quot;1536Mi\&quot;  Note that the quantity will NEVER be internally represented by a floating point number. That is the whole point of this exercise.  Non-canonical values will still parse as long as they are well formed, but will be re-emitted in their canonical form. (So always use canonical form, or don&#39;t diff.)  This format is intended to make it difficult to use these numbers without writing some sort of special handling code in the hopes that that will cause implementors to also use a fixed point implementation. | 

## Methods

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookkeeperNodeResource

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookkeeperNodeResource(cpu string, memory string, ) *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookkeeperNodeResource`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookkeeperNodeResource instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookkeeperNodeResource object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookkeeperNodeResourceWithDefaults

`func NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookkeeperNodeResourceWithDefaults() *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookkeeperNodeResource`

NewComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookkeeperNodeResourceWithDefaults instantiates a new ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookkeeperNodeResource object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCpu

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookkeeperNodeResource) GetCpu() string`

GetCpu returns the Cpu field if non-nil, zero value otherwise.

### GetCpuOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookkeeperNodeResource) GetCpuOk() (*string, bool)`

GetCpuOk returns a tuple with the Cpu field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCpu

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookkeeperNodeResource) SetCpu(v string)`

SetCpu sets Cpu field to given value.


### GetDirectPercentage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookkeeperNodeResource) GetDirectPercentage() int32`

GetDirectPercentage returns the DirectPercentage field if non-nil, zero value otherwise.

### GetDirectPercentageOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookkeeperNodeResource) GetDirectPercentageOk() (*int32, bool)`

GetDirectPercentageOk returns a tuple with the DirectPercentage field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDirectPercentage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookkeeperNodeResource) SetDirectPercentage(v int32)`

SetDirectPercentage sets DirectPercentage field to given value.

### HasDirectPercentage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookkeeperNodeResource) HasDirectPercentage() bool`

HasDirectPercentage returns a boolean if a field has been set.

### GetHeapPercentage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookkeeperNodeResource) GetHeapPercentage() int32`

GetHeapPercentage returns the HeapPercentage field if non-nil, zero value otherwise.

### GetHeapPercentageOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookkeeperNodeResource) GetHeapPercentageOk() (*int32, bool)`

GetHeapPercentageOk returns a tuple with the HeapPercentage field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHeapPercentage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookkeeperNodeResource) SetHeapPercentage(v int32)`

SetHeapPercentage sets HeapPercentage field to given value.

### HasHeapPercentage

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookkeeperNodeResource) HasHeapPercentage() bool`

HasHeapPercentage returns a boolean if a field has been set.

### GetJournalDisk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookkeeperNodeResource) GetJournalDisk() string`

GetJournalDisk returns the JournalDisk field if non-nil, zero value otherwise.

### GetJournalDiskOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookkeeperNodeResource) GetJournalDiskOk() (*string, bool)`

GetJournalDiskOk returns a tuple with the JournalDisk field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetJournalDisk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookkeeperNodeResource) SetJournalDisk(v string)`

SetJournalDisk sets JournalDisk field to given value.

### HasJournalDisk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookkeeperNodeResource) HasJournalDisk() bool`

HasJournalDisk returns a boolean if a field has been set.

### GetLedgerDisk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookkeeperNodeResource) GetLedgerDisk() string`

GetLedgerDisk returns the LedgerDisk field if non-nil, zero value otherwise.

### GetLedgerDiskOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookkeeperNodeResource) GetLedgerDiskOk() (*string, bool)`

GetLedgerDiskOk returns a tuple with the LedgerDisk field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLedgerDisk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookkeeperNodeResource) SetLedgerDisk(v string)`

SetLedgerDisk sets LedgerDisk field to given value.

### HasLedgerDisk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookkeeperNodeResource) HasLedgerDisk() bool`

HasLedgerDisk returns a boolean if a field has been set.

### GetMemory

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookkeeperNodeResource) GetMemory() string`

GetMemory returns the Memory field if non-nil, zero value otherwise.

### GetMemoryOk

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookkeeperNodeResource) GetMemoryOk() (*string, bool)`

GetMemoryOk returns a tuple with the Memory field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMemory

`func (o *ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1BookkeeperNodeResource) SetMemory(v string)`

SetMemory sets Memory field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


