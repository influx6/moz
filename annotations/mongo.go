package annotations

import (
	"github.com/influx6/faux/fmtwriter"
	"github.com/influx6/moz"
	"github.com/influx6/moz/ast"
	"github.com/influx6/moz/gen"
	"github.com/influx6/moz/gen/templates"
)

var (
	_ = moz.RegisterAnnotation("mongo", MongoAnnotationGenerator)
)

// MongoAnnotationGenerator defines a code generator for struct declarations that generate a
// mongo CRUD code for the use of mongodb as the underline db store.
func MongoAnnotationGenerator(an ast.AnnotationDeclaration, declr ast.StructDeclaration, pkg ast.PackageDeclaration) ([]gen.WriteDirective, error) {
	mongoGen := gen.Block(
		gen.Commentary(
			gen.SourceText(`Package mongoapi provides a auto-generated package which contains a mongo CRUD API for the specific {{.Object.Name}} struct in package {{.Package}}.`, str),
		),
		gen.Package(
			gen.Name("mongoapi"),
			gen.Imports(
				gen.Import("encoding/json", ""),
				gen.Import("github.com/influx6/faux/context", ""),
				gen.Import("github.com/influx6/faux/metrics", ""),
				gen.Import("github.com/influx6/faux/metrics/sentries/stdout", ""),
				gen.Import(str.Path, ""),
			),
			gen.Block(
				gen.SourceText(
					string(templates.Must("mongo-api.tml")),
					struct {
						Struct ast.StructDeclaration
					}{
						Struct: str,
					},
				),
			),
		),
	)

	return []gen.WriteDirective{
		{
			Writer:   fmtwriter.New(mongoGen, true),
			FileName: "mongoapi.go",
			Dir:      "mongoapi",
		},
	}, nil
}
