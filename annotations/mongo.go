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
	_ = moz.RegisterAnnotation("mongoapi", MongoAnnotationGenerator)
)

// MongoAnnotationGenerator defines a code generator for struct declarations that generate a
// mongo CRUD code for the use of mongodb as the underline db store.
func MongoAnnotationGenerator(an ast.AnnotationDeclaration, str ast.StructDeclaration, pkg ast.PackageDeclaration) ([]gen.WriteDirective, error) {
	updateAction := str
	createAction := str

	switch len(str.Associations) {
	case 0:
		updateAction = str
		createAction = str
		break

	default:
		if newAction, ok := str.Associations["New"]; ok {
			if action, err := ast.FindStructType(pkg, newAction.TypeName); err == nil {
				createAction = action
			}
		}

		if upAction, ok := str.Associations["Update"]; ok {
			if action, err := ast.FindStructType(pkg, upAction.TypeName); err == nil {
				updateAction = action
			}
		}

	}

	mongoGen := gen.Block(
		gen.Commentary(
			gen.SourceText(`Package mongoapi provides a auto-generated package which contains a mongo CRUD API for the specific {{.Object.Name}} struct in package {{.Package}}.`, str),
		),
		gen.Package(
			gen.Name("mongoapi"),
			gen.Imports(
				gen.Import("encoding/json", ""),
				gen.Import("gopkg.in/mgo.v2", "mgo"),
				gen.Import("gopkg.in/mgo.v2/bson", ""),
				gen.Import("github.com/influx6/faux/context", ""),
				gen.Import("github.com/influx6/faux/metrics", ""),
				gen.Import("github.com/influx6/faux/metrics/sentries/stdout", ""),
				gen.Import(str.Path, ""),
			),
			gen.Block(
				gen.SourceTextWith(
					string(templates.Must("mongo-api.tml")),
					template.FuncMap{
						"map":     ast.MapOutFields,
						"hasFunc": ast.HasFunctionFor(pkg),
					},
					struct {
						Struct       ast.StructDeclaration
						CreateAction ast.StructDeclaration
						UpdateAction ast.StructDeclaration
					}{
						Struct:       str,
						CreateAction: createAction,
						UpdateAction: updateAction,
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
			// DontOverride: true,
		},
	}, nil
}
