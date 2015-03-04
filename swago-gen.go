package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/rmullinnix/JSONParse"
	"io/ioutil"
	"os"
	"strings"
	"strconv"
)

type config struct {
	ServerTemplate		string		`json:"serverTemplate"`
	ClientTemplate		string		`json:"clientTemplate"`
	OperationTemplate	string		`json:"operationTemplate"`
	DefsTemplate		string		`json:"defsTemplate"`
	Options			configOptions	`json:"options"`
}

type configOptions struct {
	NamingConvention	string		`json:"namingConvention,omitempty"`
	FileSuffix		string		`json:"fileSuffix,omitempty"`
}

type swDocument struct {
	docJson			string
	version			string
}

type serviceDefinition struct {
	ServiceName		string
	Package			string
	BasePath		string
	Title			string
	Description		string
	ServiceVersion		string
	ContactUrl		string
	ContactEmail		string
	LicenseName		string
	Consumes		[]string
	Produces		[]string
	Methods			[]Op
	SecurityDefs		[]sec
	Definitions		[]def
}

type opers struct {
	Path			string
	Package			string
	ServiceName		string
	GroupName		string
	Methods			[]Op
}

type Op struct {
	ServiceName		string
	Package			string
	OperationId		string
	LowerOperationId	string
	HTTPMethod		string
	Summary			string
	Description		string
	Notes			string
	Path			string
	DefPath			string
	Schema			string
	Tags			[]string
	Consumes		[]string
	Produces		[]string
	Responses		[]response
	BodyParam		[]param
	PathParam		[]param
	QueryParam		[]param
	Params			[]param
	Output			string
	Security		string
}

type response struct {
	Code			string
	Description		string
	Output			string
}

type param struct {
	ParamType		string
	Name			string
	Type			string
}

type def struct {
	Name			string
	Members			[]props
}

type sec struct {
	Key			string
	Mode			string
	Location		string
	Name			string
	Prefix			string
}
	
type props struct {
	Name			string
	Type			string
	Description		string
	Required		bool
	JsonName		string
	JsonOmit		bool
}

var dir				string
var basePath			string
var funcMap			map[string]interface{}
var javaType			map[string]string

func main() {
	filePtr := flag.String("file", "", "file to parse")
	dirPtr := flag.String("dir", "", "directory to place generated files")
	confPtr := flag.String("conf", "", "configuration file for framework")
	pkgPtr := flag.String("package", "", "name of package")

	flag.Parse()

	file := *filePtr
	confFile := *confPtr
	dir = *dirPtr
	pkgName := *pkgPtr

	if len(file) == 0 {
		fmt.Println("Error:  -file must be specified")
		os.Exit(1)
	}

	if len(dir) == 0 {
		fmt.Println("Error:  -dir must be specified")
		os.Exit(1)
	}

	if len(confFile) == 0 {
		fmt.Println("Error:  -conf must be specified")
		os.Exit(1)
	}

	conf, err := loadConfig(confFile)
	if err != nil {
		fmt.Println("Error:  invalid configuration file (", err, ")")
		os.Exit(1)
	}

	swDoc, err := parseFile(file)
	if err != nil {
		fmt.Println("Error:  error parsing file: (", err, ")")
		os.Exit(1)
	}

	funcMap = make(map[string]interface{}, 20)
	funcMap["lenparms"] = lenparms
	funcMap["lenresp"] = lenresp
	funcMap["ArrayToString"] = ArrayToString
	funcMap["FirstLower"] = FirstLower
	funcMap["FirstUpper"] = FirstUpper
	funcMap["RemoveSpace"] = RemoveSpace
	funcMap["SnakeToCamel"] = SnakeToCamel
	funcMap["LastPathItem"] = LastPathItem
	funcMap["QuoteString"] = QuoteString
	funcMap["JavaType"] = JavaType
	funcMap["JavaMediaType"] = JavaMediaType

	javaType = make(map[string]string)
	javaType["string"] = "String"
	javaType["int"] = "Integer"
	javaType["bool"] = "boolean"
	javaType["[]string"] = "String[]"
	javaType["[]int"] = "Integer[]"

	if swDoc.version == "2.0"  {
		err = swagger20Spec(swDoc.docJson, conf, pkgName)
	} else if swDoc.version == "1.2" {
	}

	if err != nil {
		fmt.Println("Error: error generating files: (", err, ")")
		os.Exit(1)
	}
}

func loadConfig(confFile string) (*config, error) {
	raw, err := ioutil.ReadFile("conf/" + confFile)
	if err != nil {
		return nil, err
	}

	conf := new(config)
	err = json.Unmarshal(raw, conf)

	if err != nil {
		return nil, err
	}
	fmt.Println("options <", conf.Options, ">")
	return conf, nil
}

func parseFile(fileName string) (*swDocument, error) {
	fmt.Println("parse file", fileName)
	var parser	*JSONParse.JSONParser

	parser = JSONParse.NewJSONParser(fileName, 10, "error")
	valDoc, errs := parser.Parse()
	if !valDoc {
		for i := range errs {
			parser.OutputError(errs[i])
		}
		return nil, errors.New("errors encountered parsing file")
	}

	doc := parser.GetDoc()
	version := ""
	schemaFile := ""
	if node, found := doc.Find("swagger"); found {
		version = node.GetValue().(string)
		schemaFile = "schema.json"
	} else if node, found = doc.Find("swaggerVersion"); found {
		version = node.GetValue().(string)
		schemaFile = "apiDeclaration.json"
	} else {
		return nil, errors.New("unable to determine swagger version")
	}

	fmt.Println("load schema")
	schema := JSONParse.NewJSONSchema("schema/v" + version + "/" + schemaFile, "error")

	fmt.Println("  -validate file against schema", schemaFile)
	valid, schemaErrs := schema.ValidateDocument(fileName)
	if !valid {
		schemaErrs.Output()
		return nil, errors.New("document is not valid against schema")
	}

	swDoc := new(swDocument)
	swDoc.docJson = doc.GetJson()
	swDoc.version = version

	return swDoc, nil
}

func lenresp(arr []response, n int) int {
	return len(arr) - n
}

func lenparms(arr []param, n int) int {
	return len(arr) - n
}

func FirstUpper(val string) string {
	return strings.ToUpper(val[:1]) + val[1:]
}

func FirstLower(val string) string {
	return strings.ToLower(val[:1]) + val[1:]
}

func RemoveSpace(val string) string {
	return strings.Replace(val, " ", "", -1)
}

func SnakeToCamel(val string) string {
	for {
		index := strings.Index(val, "_")
		if index > -1 {
			if index + 1 < len(val) {
				val = val[:index] + FirstUpper(val[index+1:])
			} else {
				val = val[:index] + val[index + 1:]
			}
		} else {
			break
		}
	}
	return val
}

func ArrayToString(arr []string) string {
	return strings.Join(arr, ",")
}

func LastPathItem(path string) string {
	return path[strings.LastIndex(path, "/") + 1:]
}

func QuoteString(val string) string {
	return strconv.Quote(val)
}

func JavaType(val string) string {
	key, found := javaType[val]
	if !found {
		if val[:2] == "[]"  {
			val = val[2:] + val[:2]
		}
		return val
	} else {
		return key
	}
}

func JavaMediaType(val []string) string {
	newval := make([]string, len(val))
	for i := range val {
		newval[i] = "MediaType." + strings.ToUpper(val[i])
		newval[i] = strings.Replace(newval[i], "/" , "_", -1)
	}
	return strings.Join(newval, ",")
}
