package main

import (
)

// Swagger 1.2 Specification Structures
type SwaggerAPI12 struct {
	SwaggerVersion	string 			`json:"swaggerVersion"`
	APIVersion	string			`json:"apiVersion"`
	BasePath	string			`json:"basePath"`
	ResourcePath	string			`json:"resourcePath"`
	APIs		[]API			`json:"apis"`
	Models		map[string]Model	`json:"models"`
	Produces	[]string		`json:"produces"`
	Consumes	[]string		`json:"consumes"`
	Authorizations	map[string]Authorization `json:"authorizations"`
}

type API struct {
	Path		string			`json:"path"`
	Description	string			`json:"description,omitempty"`
	Operations	[]Operation		`json:"operations"`
}

type Operation struct {
	Method		string			`json:"method"`
	Type		string			`json:"type"`
	Summary		string			`json:"summary,omitempty"`
	Notes		string			`json:"notes,omitempty"`
	Nickname	string			`json:"nickname"`
	Authorizations	[]Authorization		`json:"authorizations"`
	Parameters	[]Parameter		`json:"parameters"`
	Responses	[]ResponseMessage	`json:"responseMessages"`
	Produces	[]string		`json:"produces,omitempty"`
	Consumes	[]string		`json:"consumes,omitempty"`
	Depracated	string			`json:"depracated,omitempty"`
}

type Parameter struct {
	ParamType	string			`json:"paramType"`
	Name		string			`json:"name"`
	Type		string			`json:"type"`
	Description	string			`json:"description,omitempty"`
	Required	bool			`json:"required,omitempty"`
	AllowMultiple	bool			`json:"allowMultiple,omitempty"`
}

type ResponseMessage struct {
	Code		int			`json:"code"`
	Message		string			`json:"message"`
	ResponseModel	string			`json:"responseModel,omitempty"`
}

type Model struct {
	ID		string			`json:"id"`
	Description	string			`json:"description,omitempty"`
	Required	[]string		`json:"required,omitempty"`
	Properties	map[string]interface{} 	`json:"properties"`
	SubTypes	[]string		`json:"subTypes,omitempty"`
	Discriminator	string			`json:"discriminator,omitempty"`
}

type Property struct {
	Type		string			`json:"type"`
	Format		string			`json:"format,omitempty"`
	Description	string			`json:"description,omitempty"`
}

type PropertyArray struct {
	Type		string			`json:"type"`
	Format		string			`json:"format,omitempty"`
	Description	string			`json:"description,omitempty"`
	Items		Property		`json:"items"`
}

type Authorization struct {
	Scope		string			`json:"scope"`
	Description	string			`json:"description,omitempty"`
}
