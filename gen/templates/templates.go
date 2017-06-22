// Package {{PKGNAME}} is an auto-generated package which exposes the specific  
// functionalities needed as desired to the specific reason which this package 
// exists for. Feel free to change this description.

//go:generate go run generate.go

package templates

import (
	"fmt"
)

var internalFiles = map[string]string{}


// Must retrieves the giving file and the content of that giving file else 
// panics if not found.
func Must(file string) string {
	if content, ok := Get(file); ok {
		return content
	}
	
	panic(fmt.Sprintf("File %s not found", file))
}


// Get retrieves the giving file and the content of that giving file.
func Get(file string) (string, bool) {
	item, ok := internalFiles[file]
	return item, ok
}

func init(){
	internalFiles["import.tml"] = "import ({{ range .Imports}}\n    {{.}}\n{{end}})\n\n"
	internalFiles["text.tml"] = "{{.Block}}"
	internalFiles["variable-assign.tml"] = "{{.Name}}:={{.Value}}\n"
	internalFiles["if.tml"] = "if{{.Condition}}{\n{{.Action}}\n}"
	internalFiles["jsonblock.tml"] = "{\n{{ range $k, $v :=  . }}\n    {{quote $k}}: {{indent $v}},\n{{ end }}\n}"
	internalFiles["slicevalue.tml"] = "[]{{.Type}}{ {{.Values}} }"
	internalFiles["tag.tml"] = "{{.Format}}:\"{{.Name}}\""
	internalFiles["value-assign.tml"] = "{{.Name}} = {{.Value}}\n"
	internalFiles["structtype.tml"] = "{{.Name}} {{.Type}} {{.Tags}}"
	internalFiles["json.tml"] = "{\n{{ range .Documents }}\n    {{.}},\n{{ end }}\n}"
	internalFiles["multicomments.tml"] = "/* {{.MainBlock}}\n{{ range .Blocks}}\n* {{.}}\n{{end}}\n*/\n"
	internalFiles["name.tml"] = "{{.Name}}"
	internalFiles["typename.tml"] = "{{.Type}}"
	internalFiles["variable-type-only.tml"] = "{{.Type}}"
	internalFiles["comments.tml"] = "// {{.MainBlock}}\n// {{ range .Blocks}}\n// {{.}}\n// {{end}}\n//\n"
	internalFiles["function.tml"] = "\nfunc {{.Name}}{{.Constructor}} {{.Returns}} {\n{{.Body}}\n}\n"
	internalFiles["package.tml"] = "{{ if notequal .Name \"\" }}package {{.Name}}\n{{ end }}\n{{.Body}}"
	internalFiles["switch.tml"] = "switch {{.Condition}} {\n{{.Case }}\n{{.Default }}\n}"
	internalFiles["variable-name.tml"] = "{{.Name}}"
	internalFiles["map-header.tml"] = "{{.MapType}}[{{.Type}}]{{.ValueType}}"
	internalFiles["struct.tml"] = "{{.Comments}}\n{{.Annotations}}\ntype {{.Name}} {{.Type}} {\n{{ range .Fields }}\n    {{.}} \n{{ end }}\n}"
	internalFiles["value.tml"] = "{{.Value}}"
	internalFiles["variable-assign-basic.tml"] = "var {{.Name}} = {{.Value}}\n"
	internalFiles["variable-type.tml"] = "{{.Name}} {{.Type}}\n"
	internalFiles["case-default.tml"] = "default:\n    {{.Action}}\n\n\n "
	internalFiles["function-type.tml"] = "func {{.Name}}{{.Constructor}} {{.Returns}}"
	internalFiles["import-item.tml"] = "{{.Namespace}} \"{{.Path}}\"\n"
	internalFiles["map.tml"] = "{{.MapType}}[{{.Type}}]{{.Value}}{\n    {{ range $k, $v :=  .Values }}\n        {{quote $k}}: {{$v}},\n    {{ end }}\n}"
	internalFiles["var-variable-type.tml"] = "var {{.Name}} {{.Type}}\n"
	internalFiles["annotations.tml"] = "//@{{.Value}}"
	internalFiles["array.tml"] = "[{{.Size}}]{{.Type}}"
	internalFiles["case.tml"] = "case {{.Condition}}:\n    {{.Action}}\n\n\n "
	internalFiles["slicetype.tml"] = "[]{{.Type}}"

}