package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

func swagger20Spec(strJson string, conf *config, pkgName string) error {
	var doc20       SwaggerAPI20

	err := json.Unmarshal([]byte(strJson), &doc20)

	if err != nil {
		return err
	}

	svc := gen20TemplateStruct(&doc20, conf, pkgName)

	raw, err := ioutil.ReadFile("template/" + conf.ServerTemplate)
	if err != nil {
		return err
	}

	templ, err := template.New("server").Funcs(funcMap).Parse(string(raw))
	if err != nil {
		return err
	}

	svcFile, err := os.Create(dir + "/" + svc.ServiceName + "_main" + conf.Options.FileSuffix)
	if err != nil {
		return err
	}

	err = templ.Execute(svcFile, svc)
	if err != nil {
		return err
	}

	raw, err = ioutil.ReadFile("template/" + conf.DefsTemplate)
	if err != nil {
		return err
	}

	defs, err := template.New("defs").Funcs(funcMap).Parse(string(raw))
	if err != nil {
		return err
	}

	defFile, err := os.Create(dir + "/" + svc.ServiceName + "_defs" + conf.Options.FileSuffix)
	err = defs.Execute(defFile, svc)
	if err != nil {
		return err
	}

	return nil
}

func gen20TemplateStruct(swSpec *SwaggerAPI20, conf *config, pkgName string) *serviceDefinition {
	svc := new(serviceDefinition)
	svc.Methods = make([]Op, 0)
	svc.Tags = make([]TagItem, 0)

	svc.ServiceName = RemoveSpace(swSpec.Info.Title)
	svc.Package = pkgName
	svc.BasePath = swSpec.BasePath
	basePath = svc.BasePath
	svc.Consumes = swSpec.Consumes
	svc.Produces = swSpec.Produces
	svc.Title = swSpec.Info.Title
	svc.ServiceVersion = swSpec.Info.Version
	svc.Description = swSpec.Info.Description
	svc.ContactUrl = swSpec.Info.Contact.Url
	svc.ContactEmail = swSpec.Info.Contact.Email
	svc.LicenseName = swSpec.Info.License.Name
	svc.SwaggerEndPoint = swSpec.XSwaggerEndPoint

	for key, path := range swSpec.Paths {
		ops := make([]Op, 0)
		if path.Get != nil {
			meth := populateOperation(key, path.Get)
			meth.HTTPMethod = "GET"
			meth.ServiceName = svc.ServiceName
			meth.Package = svc.Package
			svc.Methods = append(svc.Methods, *meth)
			ops = append(ops, *meth)
		}
		if path.Put != nil {
			meth := populateOperation(key, path.Put)
			meth.HTTPMethod = "PUT"
			meth.ServiceName = svc.ServiceName
			meth.Package = svc.Package
			svc.Methods = append(svc.Methods, *meth)
			ops = append(ops, *meth)
		}
		if path.Post != nil {
			meth := populateOperation(key, path.Post)
			meth.HTTPMethod = "POST"
			meth.ServiceName = svc.ServiceName
			meth.Package = svc.Package
			svc.Methods = append(svc.Methods, *meth)
			ops = append(ops, *meth)
		}
		if path.Delete != nil {
			meth := populateOperation(key, path.Delete)
			meth.HTTPMethod = "DELETE"
			meth.ServiceName = svc.ServiceName
			meth.Package = svc.Package
			svc.Methods = append(svc.Methods, *meth)
			ops = append(ops, *meth)
		}
		genOperationGroup(key, ops, conf)
	}

	for key, secDef := range swSpec.SecurityDefs {
		secItem := new(sec)
		secItem.Key = key
		secItem.Mode = secDef.Type
		secItem.Location = secDef.In
		secItem.Name = secDef.Name
		secItem.Prefix = "Bearer "

		svc.SecurityDefs = append(svc.SecurityDefs, *secItem)
	}

	for key, definition := range swSpec.Definitions {
		defItem := populateDefinition(key, definition)
		svc.Definitions = append(svc.Definitions, *defItem)
	}

	for i := range swSpec.Tags {
		tagItem := new(TagItem)
		tagItem.Name = swSpec.Tags[i].Name
		tagItem.Description = swSpec.Tags[i].Description
		tagItem.ExtDocUrl = swSpec.Tags[i].ExternalDocs.Url
		tagItem.ExtDocDesc = swSpec.Tags[i].ExternalDocs.Description
		svc.Tags = append(svc.Tags, *tagItem)
	}

	return svc
}

func populateOperation(key string, op *OperationObject) *Op {
	meth := new(Op)

	meth.OperationId = RemoveSpace(op.OperationId)
	meth.Summary = op.Summary
	meth.Description = op.Description
	meth.Notes = op.Description
	meth.Tags = op.Tags
	meth.Produces = op.Produces
	meth.Consumes = op.Consumes
	meth.Schema = ""
	meth.Path = key

	meth.Responses = populateResponses(meth, op)
	meth.Params = populateParameters(meth, op)

	meth.DefPath = meth.Path
	for i := range meth.PathParam {
		item := meth.PathParam[i]
		name := FirstUpper(item.Name)

		meth.DefPath = strings.Replace(meth.DefPath, name, name + ":" + item.Type, 1)
	}

	if len(meth.QueryParam) > 0 {
		token := "?"
		for i := range meth.QueryParam {
			item := meth.QueryParam[i]
			meth.DefPath = meth.DefPath + token + "{" + item.Name + ":" + item.Type + "}"
			token = "&"
		}
	}

	for i := range op.Security {
		sec := op.Security[i]
		for secKey, _ := range sec {
			meth.Security = secKey
			break
		}
		break
	}
	return meth
}

func populateResponses(meth *Op, op *OperationObject) []response {
	respList := make([]response, 0)

	for key, item := range op.Responses {
		resp := new(response)
		resp.Code = key
		resp.Description = item.Description
		if item.Schema != nil {
			if item.Schema.Ref != "" {
				meth.Schema = LastPathItem(item.Schema.Ref)
			} else if item.Schema.Type == "array"  {
				meth.Schema = "[]" + LastPathItem(item.Schema.Items.Ref)
			} else {
				meth.Schema = item.Schema.Type
			}
			resp.Output = meth.Schema
		}
		respList = append(respList, *resp)
	}
	return respList
}

func populateParameters(meth *Op, op *OperationObject) []param {
	meth.PathParam = make([]param, 0)
	meth.QueryParam = make([]param, 0)
	for i := range op.Parameters {
		par := op.Parameters[i]
		opPar := new(param)

		opPar.ParamType = par.In
		opPar.Type = par.Type
		if opPar.Type == "boolean" {
			opPar.Type = "bool"
		} else if opPar.Type == "integer" {
			opPar.Type = "int"
		} else if opPar.Type == "array" {
			opPar.Type = "[]" + par.Items.Type
		}

		opPar.Name = par.Name
		opPar.Name = FirstLower(opPar.Name)

		if par.Schema != nil {
			if par.Schema.Type == "array" {
				opPar.Type = "[]" + LastPathItem(par.Schema.Items.Ref)
			} else if par.Schema.Ref != "" {
				opPar.Type = LastPathItem(par.Schema.Ref)
			}
		}

		if par.In == "body"  {
			meth.BodyParam = append(meth.BodyParam, *opPar)
		} else if par.In == "path"  {
			meth.PathParam = append(meth.PathParam, *opPar)
		} else if par.In == "query" {
			meth.QueryParam = append(meth.QueryParam, *opPar)
		}
	}
	parms := make([]param, 0)
	parms = append(parms, meth.BodyParam...)
	parms = append(parms, meth.PathParam...)
	parms = append(parms, meth.QueryParam...)

	return parms
}

func genOperationGroup(path string, ops []Op, conf *config) {
	var oplist	opers

	oplist.Methods = ops
	oplist.Package = ops[0].Package
	oplist.ServiceName = ops[0].ServiceName
	raw, err := ioutil.ReadFile("template/" + conf.OperationTemplate)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	templ, err := template.New("operation").Funcs(funcMap).Parse(string(raw))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	opName := basePath + ops[0].Path
	opName = strings.TrimLeft(opName, "/")
	opName = strings.Replace(opName, "/", "_", -1)
	opName = strings.Replace(opName, "{", "", -1)
	opName = strings.Replace(opName, "}", "", -1)
	oplist.GroupName = SnakeToCamel(opName)
	opName += conf.Options.FileSuffix
	oplist.Path = ops[0].Path
	opFile, _ := os.Create(dir + "/" + opName)

	err = templ.Execute(opFile, oplist)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func populateDefinition(key string, svcDef SchemaObject) *def {
	var required	map[string]bool

	required = make(map[string]bool)

	for req := range svcDef.Required {
		required[svcDef.Required[req]] = true
	}

	newDef := new(def)
	newDef.Name = key
	for key, item := range svcDef.Properties {
		var prop	props

		prop.Name = key
		prop.Type = item.Type
		prop.Description = item.Description
		if item.Type == "array" {
			if len(item.Items.Ref) > 0 {
				prop.Type = "[]" + LastPathItem(item.Items.Ref)
			} else {
				prop.Type = "[]" + item.Items.Type
			}
		} else if item.Type == "boolean" {
			prop.Type = "bool"
		} else if prop.Type == "integer" {
			prop.Type = "int"
		} else if item.Type == "" {
			if len(item.Ref) > 0 {
				prop.Type = LastPathItem(item.Ref)
			}
		}

		if _, found := required[key]; found {
			prop.Required = true
		} else {
			prop.JsonOmit = true
		}

		if strings.ToLower(key[:1]) == key[:1]  {
			prop.Name = FirstUpper(key)
			prop.JsonName = key
		}
		newDef.Members = append(newDef.Members, prop)
	}
	return newDef
}
