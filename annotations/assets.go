package annotations

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/influx6/faux/fmtwriter"
	"github.com/influx6/moz"
	"github.com/influx6/moz/annotations/templates"
	"github.com/influx6/moz/ast"
	"github.com/influx6/moz/gen"
)

var (
	_ = moz.RegisterAnnotation("assets", AssetsAnnotationGenerator)
)

// AssetsAnnotationGenerator defines a package level annotation generator which builds a go package in
// root of the package it appears in to provide a means to quickly draft all file contents into the created
// package.
// Annotation: @assets
// Arguments(Optional): (PackageName, FileExtensionsToSupport, DirectorNameForFiles)
// 	e.g @assets(assets, ".tml : .bol : .go : .js", mytemplates).
func AssetsAnnotationGenerator(an ast.AnnotationDeclaration, pkg ast.PackageDeclaration) ([]gen.WriteDirective, error) {
	var directives []gen.WriteDirective

	var extensions []string

	pkgName := "assets"
	contentFileName := "files"

	if argLen := len(an.Arguments); argLen != 0 {
		pkgName = an.Arguments[0]

		if argLen > 1 {
			extensions = strings.Split(an.Arguments[1], ":")
		}

		if argLen > 2 {
			contentFileName = an.Arguments[2]
		}
	}

	genFile := gen.Package(
		gen.Name(pkgName),
		gen.Text("//go:generate go run generate.go"),
	)

	directives = append(directives, gen.WriteDirective{
		Writer:   fmtwriter.New(genFile, true),
		FileName: fmt.Sprintf("%s.go", pkgName),
		Dir:      pkgName,
	})

	directives = append(directives, gen.WriteDirective{
		Dir: filepath.Join(pkgName, contentFileName),
	})

	mainFile := gen.Block(
		gen.Commentary(
			gen.Text("+build ignore"),
		),
		gen.Text("\n"),
		gen.Text("\n"),
		gen.Package(
			gen.Name("main"),
			gen.Imports(
				gen.Import("fmt", ""),
				gen.Import("path/filepath", ""),
				gen.Import("github.com/influx6/moz/gen", ""),
				gen.Import("github.com/influx6/moz/utils", ""),
				gen.Import("github.com/influx6/faux/vfiles", ""),
				gen.Import("github.com/influx6/faux/fmtwriter", ""),
				gen.Import("github.com/influx6/faux/sink", ""),
				gen.Import("github.com/influx6/faux/sink/sinks", ""),
			),
			gen.Function(
				gen.Name("main"),
				gen.Constructor(),
				gen.Returns(),
				gen.Block(
					gen.SourceText(
						string(templates.Must("assets.tml")),
						struct {
							Extensions       []string
							TargetDir        string
							Package          string
							GenerateTemplate string
						}{
							TargetDir:  contentFileName,
							Extensions: extensions,
							Package:    pkgName,
							GenerateTemplate: `{{range $key, $value := .Files}}
								files[{{quote $key}}] = []byte("{{$value}}")
							{{end}}`,
						},
					),
				),
			),
		),
	)

	directives = append(directives, gen.WriteDirective{
		Writer:   fmtwriter.New(mainFile, true),
		FileName: "generate.go",
		Dir:      pkgName,
	})

	return directives, nil
}
