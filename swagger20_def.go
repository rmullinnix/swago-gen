package main

import (
)

// Swagger 2.0 Specifiction Structures
// This is the root document object for the API specification. It combines what previously was
// the Resource Listing and API Declaration (version 1.2 and earlier) together into one document.
type SwaggerAPI20 struct {
	SwaggerVersion	string			`json:"swagger" sw.required:"true" sw.description:"Specifies the Swagger Specification version being used. It can be used by the Swagger UI and other clients to interpret the API listing. The value MUST be 2.0."`
	Info		InfoObject		`json:"info"`
	Host		string			`json:"host" sw.description:"The host (name or ip) serving the API. This MUST be the host only and does not include the scheme nor sub-paths. It MAY include a port. If the host is not included, the host serving the documentation is to be used (including the port). The host does not support path templating."`
	BasePath	string			`json:"basePath" sw.required:"true" sw.description:"The base path on which the API is served, which is relative to the host. If it is not included, the API is served directly under the host. The value MUST start with a leading slash (/). The basePath does not support path templating."`
	Schemes		[]string		`json:"schemes,omitempty" sw.description:"The transfer protocol of the API. Values MUST be from the list: http, https, ws, wss. If the schemes is not included, the default scheme to be used is the one used to access the specification."`
	Consumes	[]string		`json:"consumes" sw.description:"A list of MIME types the APIs can consume. This is global to all APIs but can be overridden on specific API calls. Value MUST be as described under Mime Types."`
	Produces	[]string		`json:"produces" sw.description:"A list of MIME types the APIs can produce. This is global to all APIs but can be overridden on specific API calls. Value MUST be as described under Mime Types."`
	Paths		map[string]PathItem	`json:"paths" sw.required:"true" sw.description:"The available paths and operations for the API."`
	Definitions	map[string]SchemaObject	`json:"definitions" sw.description:An object to hold data types produced and consumed by operations."`
	Parameters	map[string]ParameterObject	`json:"parameters,omitempty" sw.description:"An object to hold parameters that can be used across operations. This property does not define global parameters for all operations."`
	Responses	map[string]ResponseObject	`json:"responses,omitempty" sw.description:"An object to hold responses that can be used across operations. This property does not define global responses for all operations."`
	SecurityDefs	map[string]SecurityScheme	`json:"securityDefinitions,omitempty" sw.description:"Security scheme definitions that can be used across the specification."`
	Security	*SecurityRequirement	`json:"security,omitempty" sw.description:"A declaration of which security schemes are applied for the API as a whole. The list of values describes alternative security schemes that can be used (that is, there is a logical OR between the security requirements). Individual operations can override this definition."`
	Tags		[]Tag			`json:"tags,omitempty" sw.description:"A list of tags used by the specification with additional metadata. The order of the tags can be used to reflect on their order by the parsing tools. Not all tags that are used by the Operation Object must be declared. The tags that are not declared may be organized randomly or based on the tools' logic. Each tag name in the list MUST be unique."`
	ExternalDocs	*ExtDocObject		`json:"externalDocs,omitempty" sw.description:"Additional external documentation."`
	XSwaggerEndPoint	string		`json:"x-swaggerendpoint,omitempty"`
}

// The object provides metadata about the API. The metadata can be used by the clients if needed,
// and can be presented in the Swagger-UI for convenience.
type InfoObject struct {
	Title		string			`json:"title"`
	Description	string			`json:"description"`
	TermsOfService	string			`json:"termsOfService"`
	Contact		ContactObject		`json:"contact"`
	License		LicenseObject		`json:"license"`
	Version		string			`json:"version"`
}

// Contact information for the exposed API
type ContactObject struct {
	Name		string			`json:"name"`
	Url		string			`json:"url"`
	Email		string			`json:"email"`
}

// License information for the exposed API
type LicenseObject struct {
	Name		string			`json:"name"`
	Url		string			`json:"url"`
}

// Paths Object
// Holds the relative paths to the individual endpoints. The path is appended to the basePath
// in order to construct the full URL. The Paths may be empty, due to ACL constraints.
// Paths is a map[string]PathItem where the string is the /{path}
//
// Path Item - Describes the operations available on a single path. A Path Item may be empty, 
// due to ACL constraints. The path itself is still exposed to the documentation viewer but they
// will not know which operations and parameters are available.
//   todo -- more than likely, can use map[string]OperationObject with the key being the http method
type PathItem struct {
	Ref		string			`json:"$ref,omitempty"`
	Get		*OperationObject	`json:"get,omitempty"`
	Put		*OperationObject	`json:"put,omitempty"`
	Post		*OperationObject	`json:"post,omitempty"`
	Delete		*OperationObject	`json:"delete,omitempty"`
	Options		*OperationObject	`json:"options,omitempty"`
	Head		*OperationObject	`json:"head,omitempty"`
	Patch		*OperationObject	`json:"patch,omitempty"`
	Parameters	[]ParameterObject	`json:"parameters,omitempty"`
}

// Describes a single API operation on a path
type OperationObject struct {
	Tags		[]string		`json:"tags"`
	Summary		string			`json:"summary,omitempty"`
	Description	string			`json:"description,omitempty"`
	ExternalDocs	*ExtDocObject		`json:"externalDocs,omitempty"`
	OperationId	string			`json:"operationId"`
	Consumes	[]string		`json:"consumes,omitempty"`
	Produces	[]string		`json:"produces,omitempty"`
	Parameters	[]ParameterObject	`json:"parameters,omitempty"`
	Responses	map[string]ResponseObject	`json:"responses"`
	Schemes		[]string		`json:"schemes,omitempty"`
	Deprecated	bool			`json:"deprecated,omitempty"`
	Security	[]SecurityRequirement	`json:"security,omitempty"`
}

// Allows Referencing an external resource for extended documentation
type ExtDocObject struct {
	Description	string			`json:"description,omitempty"`
	Url		string			`json:"url,omitempty"`
}

// Describes a single operation parameter
// A unique parameter is defined by a combination of a name and location
// There are five possible parameter types:  Path, Query, Header, Body, and Form
type ParameterObject struct {
	Name		string			`json:"name"`
	In		string			`json:"in"`
	Description	string			`json:"description,omitempty"`
	Required	bool			`json:"required,omitempty"`
	Schema		*SchemaObject		`json:"schema,omitempty"`
	Type		string			`json:"type,omitempty"`
	Format		string			`json:"format,omitempty"`
	Items		*ItemsObject		`json:"items,omitempty"`
	CollectionFormat	string		`json:"collectionFormat,omitempty"`
	Default		interface{}		`json:"default,omitempty"`
	Maximum		float64			`json:"maximum,omitempty"`
	ExclusiveMax	bool			`json:"exclusiveMaximum,omitempty"`
	Minimum		float64			`json:"minimum,omitempty"`
	ExclusiveMin	bool			`json:"exclusiveMinimum,omitempty"`
	MaxLength	int32			`json:"maxLength,omitempty"`
	MinLength	int32			`json:"minLength,omitempty"`
	Pattern		string			`json:"pattern,omitempty"`
	MaxItems	int32			`json:"maxItems,omitempty"`
	MinItems	int32			`json:"minItems,omitempty"`
	UniqueItems	bool			`json:"uniqueItems,omitempty"`
	Enum		[]interface{}		`json:"enum,omitempty"`
	MultipleOf	float64			`json:"multipleOf,omitempty"`
}

// A limited subset of JSON-Schema's items object.  It is used by parameter definitions that
// are not located in "body"
type ItemsObject struct {
	Type		string			`json:"type"`
	Format		string			`json:"format"`
	Items		*ItemsObject		`json:"items,omitempty"`
	CollectionFormat	string		`json:"collectionFormat,omitempty"`
	Default		interface{}		`json:"default,omitempty"`
	Maximum		float64			`json:"maximum,omitempty"`
	ExclusiveMax	bool			`json:"exclusiveMaximum,omitempty"`
	Minimum		float64			`json:"minimum,omitempty"`
	ExclusiveMin	bool			`json:"exclusiveMinimum,omitempty"`
	MaxLength	int32			`json:"maxLength,omitempty"`
	MinLength	int32			`json:"minLength,omitempty"`
	Pattern		string			`json:"pattern,omitempty"`
	MaxItems	int32			`json:"maxItems,omitempty"`
	MinItems	int32			`json:"minItems,omitempty"`
	UniqueItems	bool			`json:"uniqueItems,omitempty"`
	Enum		[]interface{}		`json:"enum,omitempty"`
	MultipleOf	float64			`json:"multipleOf,omitempty"`
}

// Responses Definition Ojbect - implement as a map[string]ResponseObject
// A container for the expected responses of an operation. The container maps a HTTP 
// response code to the expected response. It is not expected from the documentation to 
// necessarily cover all possible HTTP response codes, since they may not be known in advance.
// However, it is expected from the documentation to cover a successful operation response
// and any known errors.

// Describes a single respone from an API Operation
type ResponseObject struct {
	Description	string			`json:"description"`
	Schema		*SchemaObject		`json:"schema,omitempty"`
	Headers		map[string]HeaderObject	`json:"headers,omitempty"`
	Examples	map[string]interface{}	`json:"examples,omitempty"`
}

// Header that can be sent as part of a response
type HeaderObject struct {
	Description	string			`json:"description"`
	Type		string			`json:"type"`
	Format		string			`json:"format"`
	Items		ItemsObject		`json:"items,omitempty"`
	CollectionFormat	string		`json:"collectionFormat,omitempty"`
	Default		interface{}		`json:"default,omitempty"`
	Maximum		float64			`json:"maximum,omitempty"`
	ExclusiveMax	bool			`json:"exclusiveMaximum,omitempty"`
	Minimum		float64			`json:"minimum,omitempty"`
	ExclusiveMin	bool			`json:"exclusiveMinimum"`
	MaxLength	int32			`json:"maxLength,omitempty"`
	MinLength	int32			`json:"minLength,omitempty"`
	Pattern		string			`json:"pattern,omitempty"`
	MaxItems	int32			`json:"maxItems,omitempty"`
	MinItems	int32			`json:"minItems,omitempty"`
	UniqueItems	bool			`json:"uniqueItems,omitempty"`
	Enum		[]interface{}		`json:"enum,omitempty"`
	MultipleOf	float64			`json:"multipleOf,omitempty"`
}

// A simple object to allow referencing other definitions in the specification. 
// It can be used to reference parameters and responses that are defined at the top
// level for reuse.
type ReferenceObject struct {
	Ref		string			`json:"$ref"`
}

// The Schema Object allows the definition of input and output data types. These types
// can be objects, but also primitives and arrays. This object is based on the JSON Schema
// Specification Draft 4 and uses a predefined subset of it. On top of this subset,
// there are extensions provided by this specification to allow for more complete documentation.
type SchemaObject struct {
	Ref		string			`json:"$ref,omitempty"`
	Title		string			`json:"title,omitempty"`
	Description	string			`json:"description,omitempty"`
	Type		string			`json:"type,omitempty"`
	Format		string			`json:"format,omitempty"`
	Required	[]string		`json:"required,omitempty"`
	Items		*SchemaObject		`json:"items,omitempty"`
	MaxItems	int32			`json:"maxItems,omitempty"`
	MinItems	int32			`json:"minItems,omitempty"`
	Properties	map[string]SchemaObject	`json:"properties,omitempty"`
	AdditionalProps	*SchemaObject		`json:"additionalProperties,omitempty"`
	MaxProperties	int32			`json:"maxProperties,omitempty"`
	MinProperties	int32			`json:"minProperties,omitempty"`
	AllOf		*SchemaObject		`json:"allOf,omitempty"`
	Default		interface{}		`json:"default,omitempty"`
	Maximum		float64			`json:"maximum,omitempty"`
	ExclusiveMax	bool			`json:"exclusiveMaximum,omitempty"`
	Minimum		float64			`json:"minimum,omitempty"`
	ExclusiveMin	bool			`json:"exclusiveMinimum,omitempty"`
	MaxLength	int32			`json:"maxLength,omitempty"`
	MinLength	int32			`json:"minLength,omitempty"`
	Pattern		string			`json:"pattern,omitempty"`
	UniqueItems	bool			`json:"uniqueItems,omitempty"`
	Enum		[]interface{}		`json:"enum,omitempty"`
	MultipleOf	float64			`json:"multipleOf,omitempty"`
	Discriminator	string			`json:"discriminator,omitempty"`
	ReadOnly	bool			`json:"readOnly,omitempty"`
	Xml		*XMLObject		`json:"xml,omitempty"`
	ExternalDocs	*ExtDocObject		`json:"externalDocs,omitempty"`
	Example		interface{}		`json:"example,omitempty"`
}

// A metadata object that allows for more fine-tuned XML model definitions
type XMLObject struct {
	Name		string			`json:"name,omitempty"`
	Namespace	string			`json:"namespace,omitempty"`
	Prefix		string			`json:"prefix,omitempty"`
	Attribute	bool			`json:"attribute,omitempty"`
	Wrapped		bool			`json:"wrapped,omitempty"`
}

// Allows the definition of a security scheme that can be used by the operations.
// Supported schemes are basic authentication, an API key (either as a header or as a
// query parameter) and OAth2's common flows (implicit, password, application and 
// access code).
type SecurityScheme struct {
	Type		string			`json:"type,omitempty"`
	Description	string			`json:"description,omitempty"`
	Name		string			`json:"name,omitempty"`
	In		string			`json:"in,omitempty"`
	Flow		string			`json:"flow,omitempty"`
	AuthorizationUrL	string		`json:"authorizationUrl,omitempty"`
	TokenUrl	string			`json:"tokenUrl,omitempty"`
	Scopes		map[string]string	`json:"scopes,omitempty"`
}

type SecurityDefObject struct {
}

type SecurityRequirement 	map[string][]string

// Allows adding meta data to a single tag that is used by the Operation Object. 
// It is not mandatory to have a Tag Object per tag used there.
type Tag struct {
	Name		string			`json:"name"`
	Description	string			`json:"description"`
	ExternalDocs	ExtDocObject		`json:"externalDocs,omitempty"`
}
