package annotations

import (
	"fmt"
	goast "go/ast"
	"text/template"

	"github.com/influx6/faux/fmtwriter"
	"github.com/influx6/moz"
	"github.com/influx6/moz/annotations/templates"
	"github.com/influx6/moz/ast"
	"github.com/influx6/moz/gen"
)

var (
	_ = moz.RegisterAnnotation("httpapi", HTTPRestAnnotationGenerator)
)

// HTTPRestAnnotationGenerator defines a code generator for creating a restful HTTP for a giving struct.
func HTTPRestAnnotationGenerator(an ast.AnnotationDeclaration, str ast.StructDeclaration, pkg ast.PackageDeclaration) ([]gen.WriteDirective, error) {
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

	var hasPublicID bool

	// Validate we have a `PublicID` field.
	{
	fieldLoop:
		for _, field := range str.Struct.Fields.List {
			typeIdent, ok := field.Type.(*goast.Ident)

			// if we are not a ast.Ident then skip
			if !ok {
				continue
			}

			// If typeName is not a string, skip.
			if typeIdent.Name != "string" {
				continue
			}

			for _, indent := range field.Names {
				if indent.Name == "PublicID" {
					hasPublicID = true
					break fieldLoop
				}
			}
		}
	}

	if !hasPublicID {
		return nil, fmt.Errorf(`Struct has no 'PublicID' field with 'string' type
		 Add 'PublicID string json:"public_id"' to struct %q
		`, str.Object.Name.Name)
	}

	httpGen := gen.Block(
		gen.Commentary(
			gen.SourceText(`Package http provides a auto-generated package which contains a http restful CRUD API for the specific {{.Object.Name}} struct in package {{.Package}}.`, str),
		),
		gen.Package(
			gen.Name("httpapi"),
			gen.Imports(
				gen.Import("net/http", ""),
				gen.Import("encoding/json", ""),
				gen.Import("github.com/dimfeld/httptreemux", ""),
				gen.Import("github.com/influx6/faux/context", ""),
				gen.Import("github.com/influx6/faux/metrics", ""),
				gen.Import("github.com/influx6/faux/httputil", "httputil"),
				gen.Import("github.com/influx6/faux/metrics/sentries/stdout", ""),
				gen.Import(str.Path, ""),
			),
			gen.Block(
				gen.SourceText(
					string(templates.Must("http-api.tml")),
					struct {
						Struct ast.StructDeclaration
					}{
						Struct: str,
					},
				),
			),
		),
	)

	httpReadmeGen := gen.Block(
		gen.Block(
			gen.SourceTextWith(
				string(templates.Must("http-api-readme.tml")),
				template.FuncMap{
					"map":       ast.MapOutFields,
					"mapValues": ast.MapOutValues,
					"mapJSON":   ast.MapOutFieldsToJSON,
					"hasFunc":   ast.HasFunctionFor(pkg),
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
	)

	return []gen.WriteDirective{
		{
			Writer:   httpReadmeGen,
			FileName: "readme.md",
			Dir:      "httpapi",
			// DontOverride: true,
		},
		{
			Writer:   fmtwriter.New(httpGen, true),
			FileName: "httpapi.go",
			Dir:      "httpapi",
			// DontOverride: true,
		},
	}, nil
}
