/*
Kafka Connect With Pulsar API

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.0.1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package kafkaconnect

import (
	"encoding/json"
)

// checks if the FunctionMeshConnectorDefinition type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &FunctionMeshConnectorDefinition{}

// FunctionMeshConnectorDefinition struct for FunctionMeshConnectorDefinition
type FunctionMeshConnectorDefinition struct {
	DefaultSchemaType *string `json:"defaultSchemaType,omitempty"`
	DefaultSerdeClassName *string `json:"defaultSerdeClassName,omitempty"`
	Description *string `json:"description,omitempty"`
	IconLink *string `json:"iconLink,omitempty"`
	Id *string `json:"id,omitempty"`
	ImageRegistry *string `json:"imageRegistry,omitempty"`
	ImageRepository *string `json:"imageRepository,omitempty"`
	ImageTag *string `json:"imageTag,omitempty"`
	Jar *string `json:"jar,omitempty"`
	JarFullName *string `json:"jarFullName,omitempty"`
	Name *string `json:"name,omitempty"`
	SinkClass *string `json:"sinkClass,omitempty"`
	SinkConfigClass *string `json:"sinkConfigClass,omitempty"`
	SinkConfigFieldDefinitions []ConfigFieldDefinition `json:"sinkConfigFieldDefinitions,omitempty"`
	SinkDocLink *string `json:"sinkDocLink,omitempty"`
	SinkTypeClassName *string `json:"sinkTypeClassName,omitempty"`
	SourceClass *string `json:"sourceClass,omitempty"`
	SourceConfigClass *string `json:"sourceConfigClass,omitempty"`
	SourceConfigFieldDefinitions []ConfigFieldDefinition `json:"sourceConfigFieldDefinitions,omitempty"`
	SourceDocLink *string `json:"sourceDocLink,omitempty"`
	SourceTypeClassName *string `json:"sourceTypeClassName,omitempty"`
	TypeClassName *string `json:"typeClassName,omitempty"`
	Version *string `json:"version,omitempty"`
}

// NewFunctionMeshConnectorDefinition instantiates a new FunctionMeshConnectorDefinition object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewFunctionMeshConnectorDefinition() *FunctionMeshConnectorDefinition {
	this := FunctionMeshConnectorDefinition{}
	return &this
}

// NewFunctionMeshConnectorDefinitionWithDefaults instantiates a new FunctionMeshConnectorDefinition object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewFunctionMeshConnectorDefinitionWithDefaults() *FunctionMeshConnectorDefinition {
	this := FunctionMeshConnectorDefinition{}
	return &this
}

// GetDefaultSchemaType returns the DefaultSchemaType field value if set, zero value otherwise.
func (o *FunctionMeshConnectorDefinition) GetDefaultSchemaType() string {
	if o == nil || IsNil(o.DefaultSchemaType) {
		var ret string
		return ret
	}
	return *o.DefaultSchemaType
}

// GetDefaultSchemaTypeOk returns a tuple with the DefaultSchemaType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FunctionMeshConnectorDefinition) GetDefaultSchemaTypeOk() (*string, bool) {
	if o == nil || IsNil(o.DefaultSchemaType) {
		return nil, false
	}
	return o.DefaultSchemaType, true
}

// HasDefaultSchemaType returns a boolean if a field has been set.
func (o *FunctionMeshConnectorDefinition) HasDefaultSchemaType() bool {
	if o != nil && !IsNil(o.DefaultSchemaType) {
		return true
	}

	return false
}

// SetDefaultSchemaType gets a reference to the given string and assigns it to the DefaultSchemaType field.
func (o *FunctionMeshConnectorDefinition) SetDefaultSchemaType(v string) {
	o.DefaultSchemaType = &v
}

// GetDefaultSerdeClassName returns the DefaultSerdeClassName field value if set, zero value otherwise.
func (o *FunctionMeshConnectorDefinition) GetDefaultSerdeClassName() string {
	if o == nil || IsNil(o.DefaultSerdeClassName) {
		var ret string
		return ret
	}
	return *o.DefaultSerdeClassName
}

// GetDefaultSerdeClassNameOk returns a tuple with the DefaultSerdeClassName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FunctionMeshConnectorDefinition) GetDefaultSerdeClassNameOk() (*string, bool) {
	if o == nil || IsNil(o.DefaultSerdeClassName) {
		return nil, false
	}
	return o.DefaultSerdeClassName, true
}

// HasDefaultSerdeClassName returns a boolean if a field has been set.
func (o *FunctionMeshConnectorDefinition) HasDefaultSerdeClassName() bool {
	if o != nil && !IsNil(o.DefaultSerdeClassName) {
		return true
	}

	return false
}

// SetDefaultSerdeClassName gets a reference to the given string and assigns it to the DefaultSerdeClassName field.
func (o *FunctionMeshConnectorDefinition) SetDefaultSerdeClassName(v string) {
	o.DefaultSerdeClassName = &v
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *FunctionMeshConnectorDefinition) GetDescription() string {
	if o == nil || IsNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FunctionMeshConnectorDefinition) GetDescriptionOk() (*string, bool) {
	if o == nil || IsNil(o.Description) {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *FunctionMeshConnectorDefinition) HasDescription() bool {
	if o != nil && !IsNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *FunctionMeshConnectorDefinition) SetDescription(v string) {
	o.Description = &v
}

// GetIconLink returns the IconLink field value if set, zero value otherwise.
func (o *FunctionMeshConnectorDefinition) GetIconLink() string {
	if o == nil || IsNil(o.IconLink) {
		var ret string
		return ret
	}
	return *o.IconLink
}

// GetIconLinkOk returns a tuple with the IconLink field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FunctionMeshConnectorDefinition) GetIconLinkOk() (*string, bool) {
	if o == nil || IsNil(o.IconLink) {
		return nil, false
	}
	return o.IconLink, true
}

// HasIconLink returns a boolean if a field has been set.
func (o *FunctionMeshConnectorDefinition) HasIconLink() bool {
	if o != nil && !IsNil(o.IconLink) {
		return true
	}

	return false
}

// SetIconLink gets a reference to the given string and assigns it to the IconLink field.
func (o *FunctionMeshConnectorDefinition) SetIconLink(v string) {
	o.IconLink = &v
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *FunctionMeshConnectorDefinition) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FunctionMeshConnectorDefinition) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *FunctionMeshConnectorDefinition) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *FunctionMeshConnectorDefinition) SetId(v string) {
	o.Id = &v
}

// GetImageRegistry returns the ImageRegistry field value if set, zero value otherwise.
func (o *FunctionMeshConnectorDefinition) GetImageRegistry() string {
	if o == nil || IsNil(o.ImageRegistry) {
		var ret string
		return ret
	}
	return *o.ImageRegistry
}

// GetImageRegistryOk returns a tuple with the ImageRegistry field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FunctionMeshConnectorDefinition) GetImageRegistryOk() (*string, bool) {
	if o == nil || IsNil(o.ImageRegistry) {
		return nil, false
	}
	return o.ImageRegistry, true
}

// HasImageRegistry returns a boolean if a field has been set.
func (o *FunctionMeshConnectorDefinition) HasImageRegistry() bool {
	if o != nil && !IsNil(o.ImageRegistry) {
		return true
	}

	return false
}

// SetImageRegistry gets a reference to the given string and assigns it to the ImageRegistry field.
func (o *FunctionMeshConnectorDefinition) SetImageRegistry(v string) {
	o.ImageRegistry = &v
}

// GetImageRepository returns the ImageRepository field value if set, zero value otherwise.
func (o *FunctionMeshConnectorDefinition) GetImageRepository() string {
	if o == nil || IsNil(o.ImageRepository) {
		var ret string
		return ret
	}
	return *o.ImageRepository
}

// GetImageRepositoryOk returns a tuple with the ImageRepository field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FunctionMeshConnectorDefinition) GetImageRepositoryOk() (*string, bool) {
	if o == nil || IsNil(o.ImageRepository) {
		return nil, false
	}
	return o.ImageRepository, true
}

// HasImageRepository returns a boolean if a field has been set.
func (o *FunctionMeshConnectorDefinition) HasImageRepository() bool {
	if o != nil && !IsNil(o.ImageRepository) {
		return true
	}

	return false
}

// SetImageRepository gets a reference to the given string and assigns it to the ImageRepository field.
func (o *FunctionMeshConnectorDefinition) SetImageRepository(v string) {
	o.ImageRepository = &v
}

// GetImageTag returns the ImageTag field value if set, zero value otherwise.
func (o *FunctionMeshConnectorDefinition) GetImageTag() string {
	if o == nil || IsNil(o.ImageTag) {
		var ret string
		return ret
	}
	return *o.ImageTag
}

// GetImageTagOk returns a tuple with the ImageTag field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FunctionMeshConnectorDefinition) GetImageTagOk() (*string, bool) {
	if o == nil || IsNil(o.ImageTag) {
		return nil, false
	}
	return o.ImageTag, true
}

// HasImageTag returns a boolean if a field has been set.
func (o *FunctionMeshConnectorDefinition) HasImageTag() bool {
	if o != nil && !IsNil(o.ImageTag) {
		return true
	}

	return false
}

// SetImageTag gets a reference to the given string and assigns it to the ImageTag field.
func (o *FunctionMeshConnectorDefinition) SetImageTag(v string) {
	o.ImageTag = &v
}

// GetJar returns the Jar field value if set, zero value otherwise.
func (o *FunctionMeshConnectorDefinition) GetJar() string {
	if o == nil || IsNil(o.Jar) {
		var ret string
		return ret
	}
	return *o.Jar
}

// GetJarOk returns a tuple with the Jar field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FunctionMeshConnectorDefinition) GetJarOk() (*string, bool) {
	if o == nil || IsNil(o.Jar) {
		return nil, false
	}
	return o.Jar, true
}

// HasJar returns a boolean if a field has been set.
func (o *FunctionMeshConnectorDefinition) HasJar() bool {
	if o != nil && !IsNil(o.Jar) {
		return true
	}

	return false
}

// SetJar gets a reference to the given string and assigns it to the Jar field.
func (o *FunctionMeshConnectorDefinition) SetJar(v string) {
	o.Jar = &v
}

// GetJarFullName returns the JarFullName field value if set, zero value otherwise.
func (o *FunctionMeshConnectorDefinition) GetJarFullName() string {
	if o == nil || IsNil(o.JarFullName) {
		var ret string
		return ret
	}
	return *o.JarFullName
}

// GetJarFullNameOk returns a tuple with the JarFullName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FunctionMeshConnectorDefinition) GetJarFullNameOk() (*string, bool) {
	if o == nil || IsNil(o.JarFullName) {
		return nil, false
	}
	return o.JarFullName, true
}

// HasJarFullName returns a boolean if a field has been set.
func (o *FunctionMeshConnectorDefinition) HasJarFullName() bool {
	if o != nil && !IsNil(o.JarFullName) {
		return true
	}

	return false
}

// SetJarFullName gets a reference to the given string and assigns it to the JarFullName field.
func (o *FunctionMeshConnectorDefinition) SetJarFullName(v string) {
	o.JarFullName = &v
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *FunctionMeshConnectorDefinition) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FunctionMeshConnectorDefinition) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *FunctionMeshConnectorDefinition) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *FunctionMeshConnectorDefinition) SetName(v string) {
	o.Name = &v
}

// GetSinkClass returns the SinkClass field value if set, zero value otherwise.
func (o *FunctionMeshConnectorDefinition) GetSinkClass() string {
	if o == nil || IsNil(o.SinkClass) {
		var ret string
		return ret
	}
	return *o.SinkClass
}

// GetSinkClassOk returns a tuple with the SinkClass field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FunctionMeshConnectorDefinition) GetSinkClassOk() (*string, bool) {
	if o == nil || IsNil(o.SinkClass) {
		return nil, false
	}
	return o.SinkClass, true
}

// HasSinkClass returns a boolean if a field has been set.
func (o *FunctionMeshConnectorDefinition) HasSinkClass() bool {
	if o != nil && !IsNil(o.SinkClass) {
		return true
	}

	return false
}

// SetSinkClass gets a reference to the given string and assigns it to the SinkClass field.
func (o *FunctionMeshConnectorDefinition) SetSinkClass(v string) {
	o.SinkClass = &v
}

// GetSinkConfigClass returns the SinkConfigClass field value if set, zero value otherwise.
func (o *FunctionMeshConnectorDefinition) GetSinkConfigClass() string {
	if o == nil || IsNil(o.SinkConfigClass) {
		var ret string
		return ret
	}
	return *o.SinkConfigClass
}

// GetSinkConfigClassOk returns a tuple with the SinkConfigClass field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FunctionMeshConnectorDefinition) GetSinkConfigClassOk() (*string, bool) {
	if o == nil || IsNil(o.SinkConfigClass) {
		return nil, false
	}
	return o.SinkConfigClass, true
}

// HasSinkConfigClass returns a boolean if a field has been set.
func (o *FunctionMeshConnectorDefinition) HasSinkConfigClass() bool {
	if o != nil && !IsNil(o.SinkConfigClass) {
		return true
	}

	return false
}

// SetSinkConfigClass gets a reference to the given string and assigns it to the SinkConfigClass field.
func (o *FunctionMeshConnectorDefinition) SetSinkConfigClass(v string) {
	o.SinkConfigClass = &v
}

// GetSinkConfigFieldDefinitions returns the SinkConfigFieldDefinitions field value if set, zero value otherwise.
func (o *FunctionMeshConnectorDefinition) GetSinkConfigFieldDefinitions() []ConfigFieldDefinition {
	if o == nil || IsNil(o.SinkConfigFieldDefinitions) {
		var ret []ConfigFieldDefinition
		return ret
	}
	return o.SinkConfigFieldDefinitions
}

// GetSinkConfigFieldDefinitionsOk returns a tuple with the SinkConfigFieldDefinitions field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FunctionMeshConnectorDefinition) GetSinkConfigFieldDefinitionsOk() ([]ConfigFieldDefinition, bool) {
	if o == nil || IsNil(o.SinkConfigFieldDefinitions) {
		return nil, false
	}
	return o.SinkConfigFieldDefinitions, true
}

// HasSinkConfigFieldDefinitions returns a boolean if a field has been set.
func (o *FunctionMeshConnectorDefinition) HasSinkConfigFieldDefinitions() bool {
	if o != nil && !IsNil(o.SinkConfigFieldDefinitions) {
		return true
	}

	return false
}

// SetSinkConfigFieldDefinitions gets a reference to the given []ConfigFieldDefinition and assigns it to the SinkConfigFieldDefinitions field.
func (o *FunctionMeshConnectorDefinition) SetSinkConfigFieldDefinitions(v []ConfigFieldDefinition) {
	o.SinkConfigFieldDefinitions = v
}

// GetSinkDocLink returns the SinkDocLink field value if set, zero value otherwise.
func (o *FunctionMeshConnectorDefinition) GetSinkDocLink() string {
	if o == nil || IsNil(o.SinkDocLink) {
		var ret string
		return ret
	}
	return *o.SinkDocLink
}

// GetSinkDocLinkOk returns a tuple with the SinkDocLink field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FunctionMeshConnectorDefinition) GetSinkDocLinkOk() (*string, bool) {
	if o == nil || IsNil(o.SinkDocLink) {
		return nil, false
	}
	return o.SinkDocLink, true
}

// HasSinkDocLink returns a boolean if a field has been set.
func (o *FunctionMeshConnectorDefinition) HasSinkDocLink() bool {
	if o != nil && !IsNil(o.SinkDocLink) {
		return true
	}

	return false
}

// SetSinkDocLink gets a reference to the given string and assigns it to the SinkDocLink field.
func (o *FunctionMeshConnectorDefinition) SetSinkDocLink(v string) {
	o.SinkDocLink = &v
}

// GetSinkTypeClassName returns the SinkTypeClassName field value if set, zero value otherwise.
func (o *FunctionMeshConnectorDefinition) GetSinkTypeClassName() string {
	if o == nil || IsNil(o.SinkTypeClassName) {
		var ret string
		return ret
	}
	return *o.SinkTypeClassName
}

// GetSinkTypeClassNameOk returns a tuple with the SinkTypeClassName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FunctionMeshConnectorDefinition) GetSinkTypeClassNameOk() (*string, bool) {
	if o == nil || IsNil(o.SinkTypeClassName) {
		return nil, false
	}
	return o.SinkTypeClassName, true
}

// HasSinkTypeClassName returns a boolean if a field has been set.
func (o *FunctionMeshConnectorDefinition) HasSinkTypeClassName() bool {
	if o != nil && !IsNil(o.SinkTypeClassName) {
		return true
	}

	return false
}

// SetSinkTypeClassName gets a reference to the given string and assigns it to the SinkTypeClassName field.
func (o *FunctionMeshConnectorDefinition) SetSinkTypeClassName(v string) {
	o.SinkTypeClassName = &v
}

// GetSourceClass returns the SourceClass field value if set, zero value otherwise.
func (o *FunctionMeshConnectorDefinition) GetSourceClass() string {
	if o == nil || IsNil(o.SourceClass) {
		var ret string
		return ret
	}
	return *o.SourceClass
}

// GetSourceClassOk returns a tuple with the SourceClass field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FunctionMeshConnectorDefinition) GetSourceClassOk() (*string, bool) {
	if o == nil || IsNil(o.SourceClass) {
		return nil, false
	}
	return o.SourceClass, true
}

// HasSourceClass returns a boolean if a field has been set.
func (o *FunctionMeshConnectorDefinition) HasSourceClass() bool {
	if o != nil && !IsNil(o.SourceClass) {
		return true
	}

	return false
}

// SetSourceClass gets a reference to the given string and assigns it to the SourceClass field.
func (o *FunctionMeshConnectorDefinition) SetSourceClass(v string) {
	o.SourceClass = &v
}

// GetSourceConfigClass returns the SourceConfigClass field value if set, zero value otherwise.
func (o *FunctionMeshConnectorDefinition) GetSourceConfigClass() string {
	if o == nil || IsNil(o.SourceConfigClass) {
		var ret string
		return ret
	}
	return *o.SourceConfigClass
}

// GetSourceConfigClassOk returns a tuple with the SourceConfigClass field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FunctionMeshConnectorDefinition) GetSourceConfigClassOk() (*string, bool) {
	if o == nil || IsNil(o.SourceConfigClass) {
		return nil, false
	}
	return o.SourceConfigClass, true
}

// HasSourceConfigClass returns a boolean if a field has been set.
func (o *FunctionMeshConnectorDefinition) HasSourceConfigClass() bool {
	if o != nil && !IsNil(o.SourceConfigClass) {
		return true
	}

	return false
}

// SetSourceConfigClass gets a reference to the given string and assigns it to the SourceConfigClass field.
func (o *FunctionMeshConnectorDefinition) SetSourceConfigClass(v string) {
	o.SourceConfigClass = &v
}

// GetSourceConfigFieldDefinitions returns the SourceConfigFieldDefinitions field value if set, zero value otherwise.
func (o *FunctionMeshConnectorDefinition) GetSourceConfigFieldDefinitions() []ConfigFieldDefinition {
	if o == nil || IsNil(o.SourceConfigFieldDefinitions) {
		var ret []ConfigFieldDefinition
		return ret
	}
	return o.SourceConfigFieldDefinitions
}

// GetSourceConfigFieldDefinitionsOk returns a tuple with the SourceConfigFieldDefinitions field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FunctionMeshConnectorDefinition) GetSourceConfigFieldDefinitionsOk() ([]ConfigFieldDefinition, bool) {
	if o == nil || IsNil(o.SourceConfigFieldDefinitions) {
		return nil, false
	}
	return o.SourceConfigFieldDefinitions, true
}

// HasSourceConfigFieldDefinitions returns a boolean if a field has been set.
func (o *FunctionMeshConnectorDefinition) HasSourceConfigFieldDefinitions() bool {
	if o != nil && !IsNil(o.SourceConfigFieldDefinitions) {
		return true
	}

	return false
}

// SetSourceConfigFieldDefinitions gets a reference to the given []ConfigFieldDefinition and assigns it to the SourceConfigFieldDefinitions field.
func (o *FunctionMeshConnectorDefinition) SetSourceConfigFieldDefinitions(v []ConfigFieldDefinition) {
	o.SourceConfigFieldDefinitions = v
}

// GetSourceDocLink returns the SourceDocLink field value if set, zero value otherwise.
func (o *FunctionMeshConnectorDefinition) GetSourceDocLink() string {
	if o == nil || IsNil(o.SourceDocLink) {
		var ret string
		return ret
	}
	return *o.SourceDocLink
}

// GetSourceDocLinkOk returns a tuple with the SourceDocLink field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FunctionMeshConnectorDefinition) GetSourceDocLinkOk() (*string, bool) {
	if o == nil || IsNil(o.SourceDocLink) {
		return nil, false
	}
	return o.SourceDocLink, true
}

// HasSourceDocLink returns a boolean if a field has been set.
func (o *FunctionMeshConnectorDefinition) HasSourceDocLink() bool {
	if o != nil && !IsNil(o.SourceDocLink) {
		return true
	}

	return false
}

// SetSourceDocLink gets a reference to the given string and assigns it to the SourceDocLink field.
func (o *FunctionMeshConnectorDefinition) SetSourceDocLink(v string) {
	o.SourceDocLink = &v
}

// GetSourceTypeClassName returns the SourceTypeClassName field value if set, zero value otherwise.
func (o *FunctionMeshConnectorDefinition) GetSourceTypeClassName() string {
	if o == nil || IsNil(o.SourceTypeClassName) {
		var ret string
		return ret
	}
	return *o.SourceTypeClassName
}

// GetSourceTypeClassNameOk returns a tuple with the SourceTypeClassName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FunctionMeshConnectorDefinition) GetSourceTypeClassNameOk() (*string, bool) {
	if o == nil || IsNil(o.SourceTypeClassName) {
		return nil, false
	}
	return o.SourceTypeClassName, true
}

// HasSourceTypeClassName returns a boolean if a field has been set.
func (o *FunctionMeshConnectorDefinition) HasSourceTypeClassName() bool {
	if o != nil && !IsNil(o.SourceTypeClassName) {
		return true
	}

	return false
}

// SetSourceTypeClassName gets a reference to the given string and assigns it to the SourceTypeClassName field.
func (o *FunctionMeshConnectorDefinition) SetSourceTypeClassName(v string) {
	o.SourceTypeClassName = &v
}

// GetTypeClassName returns the TypeClassName field value if set, zero value otherwise.
func (o *FunctionMeshConnectorDefinition) GetTypeClassName() string {
	if o == nil || IsNil(o.TypeClassName) {
		var ret string
		return ret
	}
	return *o.TypeClassName
}

// GetTypeClassNameOk returns a tuple with the TypeClassName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FunctionMeshConnectorDefinition) GetTypeClassNameOk() (*string, bool) {
	if o == nil || IsNil(o.TypeClassName) {
		return nil, false
	}
	return o.TypeClassName, true
}

// HasTypeClassName returns a boolean if a field has been set.
func (o *FunctionMeshConnectorDefinition) HasTypeClassName() bool {
	if o != nil && !IsNil(o.TypeClassName) {
		return true
	}

	return false
}

// SetTypeClassName gets a reference to the given string and assigns it to the TypeClassName field.
func (o *FunctionMeshConnectorDefinition) SetTypeClassName(v string) {
	o.TypeClassName = &v
}

// GetVersion returns the Version field value if set, zero value otherwise.
func (o *FunctionMeshConnectorDefinition) GetVersion() string {
	if o == nil || IsNil(o.Version) {
		var ret string
		return ret
	}
	return *o.Version
}

// GetVersionOk returns a tuple with the Version field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FunctionMeshConnectorDefinition) GetVersionOk() (*string, bool) {
	if o == nil || IsNil(o.Version) {
		return nil, false
	}
	return o.Version, true
}

// HasVersion returns a boolean if a field has been set.
func (o *FunctionMeshConnectorDefinition) HasVersion() bool {
	if o != nil && !IsNil(o.Version) {
		return true
	}

	return false
}

// SetVersion gets a reference to the given string and assigns it to the Version field.
func (o *FunctionMeshConnectorDefinition) SetVersion(v string) {
	o.Version = &v
}

func (o FunctionMeshConnectorDefinition) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o FunctionMeshConnectorDefinition) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.DefaultSchemaType) {
		toSerialize["defaultSchemaType"] = o.DefaultSchemaType
	}
	if !IsNil(o.DefaultSerdeClassName) {
		toSerialize["defaultSerdeClassName"] = o.DefaultSerdeClassName
	}
	if !IsNil(o.Description) {
		toSerialize["description"] = o.Description
	}
	if !IsNil(o.IconLink) {
		toSerialize["iconLink"] = o.IconLink
	}
	if !IsNil(o.Id) {
		toSerialize["id"] = o.Id
	}
	if !IsNil(o.ImageRegistry) {
		toSerialize["imageRegistry"] = o.ImageRegistry
	}
	if !IsNil(o.ImageRepository) {
		toSerialize["imageRepository"] = o.ImageRepository
	}
	if !IsNil(o.ImageTag) {
		toSerialize["imageTag"] = o.ImageTag
	}
	if !IsNil(o.Jar) {
		toSerialize["jar"] = o.Jar
	}
	if !IsNil(o.JarFullName) {
		toSerialize["jarFullName"] = o.JarFullName
	}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.SinkClass) {
		toSerialize["sinkClass"] = o.SinkClass
	}
	if !IsNil(o.SinkConfigClass) {
		toSerialize["sinkConfigClass"] = o.SinkConfigClass
	}
	if !IsNil(o.SinkConfigFieldDefinitions) {
		toSerialize["sinkConfigFieldDefinitions"] = o.SinkConfigFieldDefinitions
	}
	if !IsNil(o.SinkDocLink) {
		toSerialize["sinkDocLink"] = o.SinkDocLink
	}
	if !IsNil(o.SinkTypeClassName) {
		toSerialize["sinkTypeClassName"] = o.SinkTypeClassName
	}
	if !IsNil(o.SourceClass) {
		toSerialize["sourceClass"] = o.SourceClass
	}
	if !IsNil(o.SourceConfigClass) {
		toSerialize["sourceConfigClass"] = o.SourceConfigClass
	}
	if !IsNil(o.SourceConfigFieldDefinitions) {
		toSerialize["sourceConfigFieldDefinitions"] = o.SourceConfigFieldDefinitions
	}
	if !IsNil(o.SourceDocLink) {
		toSerialize["sourceDocLink"] = o.SourceDocLink
	}
	if !IsNil(o.SourceTypeClassName) {
		toSerialize["sourceTypeClassName"] = o.SourceTypeClassName
	}
	if !IsNil(o.TypeClassName) {
		toSerialize["typeClassName"] = o.TypeClassName
	}
	if !IsNil(o.Version) {
		toSerialize["version"] = o.Version
	}
	return toSerialize, nil
}

type NullableFunctionMeshConnectorDefinition struct {
	value *FunctionMeshConnectorDefinition
	isSet bool
}

func (v NullableFunctionMeshConnectorDefinition) Get() *FunctionMeshConnectorDefinition {
	return v.value
}

func (v *NullableFunctionMeshConnectorDefinition) Set(val *FunctionMeshConnectorDefinition) {
	v.value = val
	v.isSet = true
}

func (v NullableFunctionMeshConnectorDefinition) IsSet() bool {
	return v.isSet
}

func (v *NullableFunctionMeshConnectorDefinition) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableFunctionMeshConnectorDefinition(val *FunctionMeshConnectorDefinition) *NullableFunctionMeshConnectorDefinition {
	return &NullableFunctionMeshConnectorDefinition{value: val, isSet: true}
}

func (v NullableFunctionMeshConnectorDefinition) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableFunctionMeshConnectorDefinition) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


