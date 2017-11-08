package annotations

import (
	"text/template"

	"github.com/influx6/faux/fmtwriter"
	"github.com/influx6/moz"
	"github.com/influx6/moz/annotations/templates"
	"github.com/influx6/moz/ast"
	"github.com/influx6/moz/gen"
)

var (
	_ = moz.RegisterAnnotation("mongo", MongoAnnotationGenerator)
)

// MongoAnnotationGenerator defines a code generator for struct declarations that generate a
// mongo CRUD code for the use of mongodb as the underline db store.
func MongoAnnotationGenerator(toDir string, an ast.AnnotationDeclaration, pkgDeclr ast.PackageDeclaration, pkg ast.Package) ([]gen.WriteDirective, error) {

	mongoReadmeGen := gen.Block(
		gen.Block(
			gen.SourceText(
				string(templates.Must("mongo/readme.tml")),
				struct {
					Pkg     *ast.PackageDeclaration
					Package ast.Package
				}{
					Pkg:     &pkgDeclr,
					Package: pkg,
				},
			),
		),
	)

	mongoGen := gen.Block(
		gen.Commentary(
			gen.Text(`Package mongoapi provides a auto-generated package which contains a mongo base pkg for db operations.`),
		),
		gen.Package(
			gen.Name("mongo"),
			gen.Imports(
				gen.Import("encoding/json", ""),
				gen.Import("gopkg.in/mgo.v2", "mgo"),
				gen.Import("gopkg.in/mgo.v2/bson", ""),
				gen.Import("github.com/influx6/faux/context", ""),
				gen.Import("github.com/influx6/faux/metrics", ""),
				gen.Import("github.com/influx6/faux/metrics/custom", ""),
			),
			gen.Block(
				gen.SourceTextWith(
					string(templates.Must("mongo/api.tml")),
					template.FuncMap{
						"map":     ast.MapOutFields,
						"hasFunc": ast.HasFunctionFor(pkgDeclr),
					},
					struct {
						Pkg     *ast.PackageDeclaration
						Package ast.Package
					}{
						Pkg:     &pkgDeclr,
						Package: pkg,
					},
				),
			),
		),
	)

	return []gen.WriteDirective{
		{
			Writer:   mongoReadmeGen,
			FileName: "README.md",
			Dir:      "mongo",
			// DontOverride: true,
		},
		{
			Writer:   fmtwriter.New(mongoGen, true, true),
			FileName: "mongoapi.go",
			Dir:      "mongo",
			// DontOverride: true,
		},
	}, nil
}
