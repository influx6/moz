package annotations

import (
	"path/filepath"
	"text/template"

	"github.com/influx6/faux/fmtwriter"
	"github.com/influx6/moz"
	"github.com/influx6/moz/annotations/templates"
	"github.com/influx6/moz/ast"
	"github.com/influx6/moz/gen"
)

var (
	_ = moz.RegisterAnnotation("mongoapi", MongoAPIAnnotationGenerator)
)

// MongoAPIAnnotationGenerator defines a code generator for struct declarations that generate a
// mongo CRUD code for the use of mongodb as the underline db store.
func MongoAPIAnnotationGenerator(toDir string, an ast.AnnotationDeclaration, str ast.StructDeclaration, pkgDeclr ast.PackageDeclaration, pkg ast.Package) ([]gen.WriteDirective, error) {
	updateAction := str
	createAction := str

	switch len(str.Associations) {
	case 0:
		updateAction = str
		createAction = str
		break

	default:
		if newAction, ok := str.Associations["New"]; ok {
			if action, err := ast.FindStructType(pkgDeclr, newAction.TypeName); err == nil {
				createAction = action
			}
		}

		if upAction, ok := str.Associations["Update"]; ok {
			if action, err := ast.FindStructType(pkgDeclr, upAction.TypeName); err == nil {
				updateAction = action
			}
		}

	}

	mongoTestGen := gen.Block(
		gen.Package(
			gen.Name("mdb_test"),
			gen.Imports(
				gen.Import("os", ""),
				gen.Import("testing", ""),
				gen.Import("encoding/json", ""),
				gen.Import("gopkg.in/mgo.v2", "mgo"),
				gen.Import("gopkg.in/mgo.v2/bson", ""),
				gen.Import("github.com/influx6/faux/tests", ""),
				gen.Import("github.com/influx6/faux/metrics", ""),
				gen.Import("github.com/influx6/faux/context", ""),
				gen.Import("github.com/influx6/faux/db/mongo", ""),
				gen.Import("github.com/influx6/faux/metrics/custom", ""),
				gen.Import(filepath.Join(str.Path, toDir, "/mdb"), ""),
				gen.Import(str.Path, ""),
			),
			gen.Block(
				gen.SourceTextWith(
					string(templates.Must("mongoapi/mongo-api-test.tml")),
					gen.ToTemplateFuncs(
						ast.ASTTemplatFuncs,
						template.FuncMap{
							"map":       ast.MapOutFields,
							"mapValues": ast.MapOutValues,
							"hasFunc":   ast.HasFunctionFor(pkgDeclr),
						},
					),
					struct {
						Pkg          *ast.PackageDeclaration
						Struct       ast.StructDeclaration
						CreateAction ast.StructDeclaration
						UpdateAction ast.StructDeclaration
					}{
						Pkg:          &pkgDeclr,
						Struct:       str,
						CreateAction: createAction,
						UpdateAction: updateAction,
					},
				),
			),
		),
	)

	mongoReadmeGen := gen.Block(
		gen.Block(
			gen.SourceText(
				string(templates.Must("mongoapi/mongo-api-readme.tml")),
				struct {
					Pkg          *ast.PackageDeclaration
					Struct       ast.StructDeclaration
					CreateAction ast.StructDeclaration
					UpdateAction ast.StructDeclaration
				}{
					Pkg:          &pkgDeclr,
					Struct:       str,
					CreateAction: createAction,
					UpdateAction: updateAction,
				},
			),
		),
	)

	mongoJSONGen := gen.Block(
		gen.Package(
			gen.Name("mdb_test"),
			gen.Imports(
				gen.Import("encoding/json", ""),
				gen.Import(str.Path, ""),
			),
			gen.Block(
				gen.SourceTextWith(
					string(templates.Must("mongoapi/mongo-api-json.tml")),
					gen.ToTemplateFuncs(
						ast.ASTTemplatFuncs,
						template.FuncMap{
							"map":       ast.MapOutFields,
							"mapValues": ast.MapOutValues,
							"mapJSON":   ast.MapOutFieldsToJSON,
							"hasFunc":   ast.HasFunctionFor(pkgDeclr),
						},
					),
					struct {
						Pkg          *ast.PackageDeclaration
						Struct       ast.StructDeclaration
						CreateAction ast.StructDeclaration
						UpdateAction ast.StructDeclaration
					}{
						Pkg:          &pkgDeclr,
						Struct:       str,
						CreateAction: createAction,
						UpdateAction: updateAction,
					},
				),
			),
		),
	)

	mongoGen := gen.Block(
		gen.Commentary(
			gen.SourceText(`Package mdb provides a auto-generated package which contains a mongo CRUD API for the specific {{.Object.Name}} struct in package {{.Package}}.`, str),
		),
		gen.Package(
			gen.Name("mdb"),
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
					string(templates.Must("mongoapi/mongo-api.tml")),
					gen.ToTemplateFuncs(
						ast.ASTTemplatFuncs,
						template.FuncMap{
							"map":     ast.MapOutFields,
							"hasFunc": ast.HasFunctionFor(pkgDeclr),
						},
					),
					struct {
						Pkg          *ast.PackageDeclaration
						Struct       ast.StructDeclaration
						CreateAction ast.StructDeclaration
						UpdateAction ast.StructDeclaration
					}{
						Pkg:          &pkgDeclr,
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
			Writer:   mongoJSONGen,
			FileName: "json_test.go",
			Dir:      "mdb",
			// DontOverride: true,
		},
		{
			Writer:   mongoReadmeGen,
			FileName: "README.md",
			Dir:      "mdb",
			// DontOverride: true,
		},
		{
			Writer:   fmtwriter.New(mongoTestGen, true, true),
			FileName: "mdb_test.go",
			Dir:      "mdb",
			// DontOverride: true,
		},
		{
			Writer:   fmtwriter.New(mongoGen, true, true),
			FileName: "mdb.go",
			Dir:      "mdb",
			// DontOverride: true,
		},
	}, nil
}
