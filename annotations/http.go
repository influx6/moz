package annotations

import (
	"github.com/influx6/faux/fmtwriter"
	"github.com/influx6/moz"
	"github.com/influx6/moz/annotations/templates"
	"github.com/influx6/moz/ast"
	"github.com/influx6/moz/gen"
)

var (
	_ = moz.RegisterAnnotation("httpRest", HTTPRestAnnotationGenerator)
)

// HTTPRestAnnotationGenerator defines a code generator for creating a restful HTTP for a giving struct.
func HTTPRestAnnotationGenerator(an ast.AnnotationDeclaration, str ast.StructDeclaration, pkg ast.PackageDeclaration) ([]gen.WriteDirective, error) {
	httpGen := gen.Block(
		gen.Commentary(
			gen.SourceText(`Package http provides a auto-generated package which contains a http restful CRUD API for 
the specific {{.Struct.Name}} struct in package {{.Package}}.`, str),
		),
		gen.Package(
			gen.Name("httpapi"),
			gen.Imports(
				gen.Import("net", ""),
				gen.Import("net/http", ""),
			),
			gen.Block(
				gen.SourceText(
					string(templates.Must("httpapi.tml")),
					struct{}{},
				),
			),
		),
	)

	return []gen.WriteDirective{
		{
			Writer:   fmtwriter.New(httpGen, true),
			FileName: "httpapi.go",
			Dir:      "httpapi",
		},
	}, nil
}
