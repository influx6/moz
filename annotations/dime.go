package annotations

import (
	"github.com/influx6/moz"
	"github.com/influx6/moz/annotations/templates"
	"github.com/influx6/moz/ast"
	"github.com/influx6/moz/gen"
)

var (
	_ = moz.RegisterAnnotation("dime", DimeGenerator)
)

// DimeGenerator defines a generator that creates a doc.go file which contains a template to create channel based code specified type names.
func DimeGenerator(toDir string, an ast.AnnotationDeclaration, pkgDeclr ast.PackageDeclaration, pkg ast.Package) ([]gen.WriteDirective, error) {
	docGen := gen.Block(
		gen.Commentary(
			gen.Fmt("Package %s contains a template within it's comments to generate code based off.", pkg.Package),
		),
		gen.Text(
			string(templates.Must("dime/dime.tml")),
		),
		gen.Package(
			gen.Name(pkg.Package),
			gen.Text("\n"),
			gen.Text("//go:generate moz generate --toDir=./"),
		),
	)

	return []gen.WriteDirective{
		{
			Writer:       docGen,
			FileName:     "doc.go",
			DontOverride: true,
		},
	}, nil
}
