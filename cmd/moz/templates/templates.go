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
	internalFiles["main.tml"] = "events := metrics.New(stdout.Stdout{})\r\n\r\nitems, err := vfiles.ParseDir(\"{{.TargetDir}}\", []string{\r\n    {{range .Extensions}}\r\n        {{ if notequal . \"\"}}\r\n        {{trimspace . | quote}},\r\n        {{ end }}\r\n    {{ end }}\r\n})\r\nif err != nil {\r\n    panic(fmt.Sprintf(\"Failed to walk directory properly: %+q\", err))\r\n}\r\n\r\nassetGen := gen.Package(\r\n    gen.Name(\"{{.Package}}\"),\r\n    gen.Imports(\r\n        gen.Import(\"fmt\", \"\"),\r\n    ),\r\n    gen.Text(\"\\n\"),\r\n    gen.Text(\"\\n\"),\r\n    gen.AssignVar(\r\n        gen.Name(\"files\"),\r\n        gen.Type(\"make(map[string][]byte)\"),\r\n    ),\r\n    gen.Text(`\r\n    // Must attempts to retrieve the file data if found else panics.\r\n    func Must(file string) []byte {\r\n        data, err := Get(file)\r\n        if err != nil {\r\n            panic(err)\r\n        }\r\n\r\n        return data\r\n    }\r\n\r\n    // Get retrieves the giving file data from the map store if it exists.\r\n    func Get(file string) ([]byte, error){\r\n        data, ok := files[file]\r\n        if !ok {\r\n            return nil, fmt.Errorf(\"File data for %q not found\", file)\r\n        }\r\n\r\n        return data, nil\r\n    }\r\n    `),\r\n    gen.Function(\r\n        gen.Name(\"init\"),\r\n        gen.Constructor(),\r\n        gen.Returns(),\r\n        gen.Block(\r\n            gen.SourceText(`\r\n                {{.GenerateTemplate}}\r\n            `, struct{\r\n                Files map[string]string\r\n            }{\r\n                Files: items,\r\n            }),\r\n        ),\r\n    ),\r\n)\r\n\r\ndir := filepath.Join(\".\", \"{{.Package}}.go\")\r\nif err := utils.WriteFile(events, fmtwriter.New(assetGen, true), dir); err != nil {\r\n    events.Emit(stdout.Error(err).With(\"dir\", dir).\r\n        With(\"message\", \"Failed to create new package file: {{.Package}}.go\"))\r\n    panic(err)\r\n}"

}