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
	_ = moz.RegisterAnnotation("sqlapi", SQLAnnotationGenerator)
)

// SQLAnnotationGenerator defines a code generator for struct declarations that generate a
// sql CRUD code for the use of sqldb as the underline db store.
func SQLAnnotationGenerator(toDir string, an ast.AnnotationDeclaration, str ast.StructDeclaration, pkgDeclr ast.PackageDeclaration, pkg ast.Package) ([]gen.WriteDirective, error) {
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

	sqlTestGen := gen.Block(
		gen.Package(
			gen.Name("sqldb_test"),
			gen.Imports(
				gen.Import("os", ""),
				gen.Import("testing", ""),
				gen.Import("encoding/json", ""),
				gen.Import("github.com/influx6/faux/db", ""),
				gen.Import("github.com/influx6/faux/tests", ""),
				gen.Import("github.com/influx6/faux/db/sql", ""),
				gen.Import("github.com/influx6/faux/metrics", ""),
				gen.Import("github.com/influx6/faux/context", ""),
				gen.Import("github.com/influx6/faux/metrics/custom", ""),
				gen.Import("github.com/go-sql-driver/mysql", "_"),
				gen.Import("github.com/lib/pq", "_"),
				gen.Import("github.com/mattn/go-sqlite3", "_"),
				gen.Import(filepath.Join(str.Path, toDir, "sqldb"), ""),
				gen.Import(str.Path, ""),
			),
			gen.Block(
				gen.SourceTextWith(
					string(templates.Must("sqlapi/sql-api-test.tml")),
					gen.ToTemplateFuncs(
						ast.ASTTemplatFuncs,
						template.FuncMap{
							"hasFunc": ast.HasFunctionFor(pkgDeclr),
						},
					),
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

	sqlReadmeGen := gen.Block(
		gen.Block(
			gen.SourceText(
				string(templates.Must("sqlapi/sql-api-readme.tml")),
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
	)

	sqlJSONGen := gen.Block(
		gen.Package(
			gen.Name("sqldb_test"),
			gen.Imports(
				gen.Import("encoding/json", ""),
				gen.Import(str.Path, ""),
			),
			gen.Block(
				gen.SourceTextWith(
					string(templates.Must("sqlapi/sql-api-json.tml")),
					template.FuncMap{
						"map":       ast.MapOutFields,
						"mapValues": ast.MapOutValues,
						"mapJSON":   ast.MapOutFieldsToJSON,
						"hasFunc":   ast.HasFunctionFor(pkgDeclr),
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

	sqlGen := gen.Block(
		gen.Commentary(
			gen.SourceText(`Package sqldb provides a auto-generated package which contains a sql CRUD API for the specific {{.Object.Name}} struct in package {{.Package}}.`, str),
		),
		gen.Package(
			gen.Name("sqldb"),
			gen.Imports(
				gen.Import("encoding/json", ""),
				gen.Import("github.com/influx6/faux/db", ""),
				gen.Import("github.com/influx6/faux/db/sql", ""),
				gen.Import("github.com/influx6/faux/context", ""),
				gen.Import("github.com/influx6/faux/metrics", ""),
				gen.Import("github.com/influx6/faux/metrics/sentries/stdout", ""),
				gen.Import(str.Path, ""),
			),
			gen.Block(
				gen.SourceTextWith(
					string(templates.Must("sqlapi/sql-api.tml")),
					gen.ToTemplateFuncs(
						ast.ASTTemplatFuncs,
						template.FuncMap{
							"hasFunc": ast.HasFunctionFor(pkgDeclr),
						},
					),
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
			Writer:   sqlJSONGen,
			FileName: "json_test.go",
			Dir:      "sqldb",
			// DontOverride: true,
		},
		{
			Writer:   sqlReadmeGen,
			FileName: "README.md",
			Dir:      "sqldb",
			// DontOverride: true,
		},
		{
			Writer:   fmtwriter.New(sqlTestGen, true, true),
			FileName: "sqldb_test.go",
			Dir:      "sqldb",
			// DontOverride: true,
		},
		{
			Writer:   fmtwriter.New(sqlGen, true, true),
			FileName: "sqldb.go",
			Dir:      "sqldb",
			// DontOverride: true,
		},
	}, nil
}
