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
	internalFiles["main.tml"] = "// go:generate go run generate.go\n\nevents := sink.New(sinks.Stdout{})\n\nitems, err := vfiles.ParseDir(\"{{.TargetDir}}\", []string{\n    {{range .Extensions}}\n        {{ if notequal . \"\"}}\n        {{trimspace . | quote}},\n        {{ end }}\n    {{ end }}\n})\nif err != nil {\n    panic(fmt.Sprintf(\"Failed to walk directory properly: %+q\", err))\n}\n\nassetGen := gen.Package(\n    gen.Name(\"{{.Package}}\"),\n    gen.AssignVar(\n        gen.Name(\"files\"),\n        gen.Type(\"make(map[string][]byte)\"),\n    ),\n    gen.Text(`\n    // Must attempts to retrieve the file data if found else panics.\n    func Must(file string) []byte {\n        data, err := Get(file)\n        if err != nil {\n            panic(err)\n        }\n\n        return data\n    }\n\n    // Get retrieves the giving file data from the map store if it exists.\n    func Get(file string) ([]byte, error){\n        data, ok := files[file]\n        if !ok {\n            return nil, fmt.Errorf(\"File data for %q not found\", file)\n        }\n\n        return data, nil\n    }\n    `),\n    gen.Function(\n        gen.Name(\"init\"),\n        gen.Constructor(),\n        gen.Returns(),\n        gen.Block(\n            gen.SourceText(`\n                {{.GenerateTemplate}}\n            `, struct{\n                Files map[string]string\n            }{\n                Files: items,\n            }),\n        ),\n    ),\n)\n\ndir := filepath.Join(\".\", \"{{.Package}}.go\")\nif err := utils.WriteFile(events, fmtwriter.New(assetGen, true), dir); err != nil {\n    events.Emit(sinks.Error(err).With(\"dir\", dir).\n        With(\"message\", \"Failed to create new package file: {{.Package}}.go\"))\n    panic(err)\n}"

}