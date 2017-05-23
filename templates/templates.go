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
	internalFiles["variable-assign.tml"] = "{{.Name}}:={{.Value}}"
	internalFiles["array.tml"] = "[{{.Size}}]{{.Type}}"
	internalFiles["function-type.tml"] = "func {{.Name}}{{.Constructor}} {{.Returns}}"
	internalFiles["package.tml"] = "{{ if .Name != \"\"}}\npackage {{.Name}}\n\n{{ end }}\n{{.Body}}"
	internalFiles["slicetype.tml"] = "[]{{.Type}}"
	internalFiles["struct.tml"] = "{{.Comments}}\n{{.Annotations}}\ntype {{.Name}} {{.Type}} {\n{{ range .Fields }}\n    {{.}} \n{{ end }}\n}"
	internalFiles["value.tml"] = "{{.Value}}"
	internalFiles["variable-type.tml"] = "{{.Name}} {{.Type}}"
	internalFiles["comments.tml"] = "// {{.MainBlock}}\n{{ range .Blocks}}\n// {{.}}\n{{end}}\n//"
	internalFiles["function.tml"] = "func {{.Name}}{{.Constructor}} {{.Returns}} {\n{{.Body}}\n}"
	internalFiles["import-item.tml"] = "{{.Namespace}} {{.Path}}\n"
	internalFiles["multicomments.tml"] = "/* {{.MainBlock}}\n{{ range .Blocks}}\n* {{.}}\n{{end}}\n*/\n"
	internalFiles["name.tml"] = "{{.Name}}"
	internalFiles["text.tml"] = "{{.Block}}"
	internalFiles["case.tml"] = "case {{.Condition}}:\n    {{.Action}}\n\n\n "
	internalFiles["import.tml"] = "import (\n{{ range .Imports}}\n    {{.}}\n{{end}}\n)"
	internalFiles["slicevalue.tml"] = "[]{{.Type}}{ {{.Values}} }"
	internalFiles["structtype.tml"] = "{{.Name}} {{.Type}} {{.Tags}}"
	internalFiles["switch.tml"] = "switch {{.Condition}} {\n{{.Case }}\n{{.Default }}\n}"
	internalFiles["variable-type-only.tml"] = "{{.Type}}"
	internalFiles["case-default.tml"] = "default:\n    {{.Action}}\n\n\n "
	internalFiles["if.tml"] = "if{{.Condition}}{\n{{.Action}}\n}"
	internalFiles["tag.tml"] = "{{.Format}}:{{.Name}}"
	internalFiles["typename.tml"] = "{{.Type}}"
	internalFiles["variable-assign-basic.tml"] = "var {{.Name}} = {{.Value}}"
	internalFiles["variable-name.tml"] = "{{.Name}}"

}