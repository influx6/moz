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
	internalFiles["variable-type.tml"] = "{{.Name}} {{.Type}}"
	internalFiles["array.tml"] = "[{{.Size}}]{{.Type}}"
	internalFiles["variable-assign.tml"] = "{{.Name}}:={{.Value}}"
	internalFiles["switch.tml"] = "switch {{.Condition}} {\n{{.Case }}\n{{.Default }}\n}"
	internalFiles["typename.tml"] = "{{.Type}}"
	internalFiles["variable-assign-basic.tml"] = "var {{.Name}} = {{.Value}}"
	internalFiles["case-default.tml"] = "default:\n    {{.Action}}\n\n\n "
	internalFiles["struct.tml"] = "{{ range .Commentary }}\n//    {{.CommentBlock}}\n{{ end }}\n{{ range .Annotations }}\n//    {{.Name}}\n{{ end }}\ntype {{.Name}} {{.Type}} {\n{{ range .Fields }}\n    {{.FieldName}} {{.FieldType}} {{.FieldTags}}\n{{ end }}\n}"
	internalFiles["slicevalue.tml"] = "[]{{.Type}}{ {{.Values}} }"
	internalFiles["variable-type-only.tml"] = "{{.Type}}"
	internalFiles["case.tml"] = "case {{.Condition}}:\n    {{.Action}}\n\n\n "
	internalFiles["function.tml"] = "func {{.Name}}{{.Constructors}} {{.Returns}} {\n{{.Body}}\n}"
	internalFiles["slicetype.tml"] = "[]{{.Type}}"
	internalFiles["variable-name.tml"] = "{{.Name}}"

}