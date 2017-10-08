package annotations

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/influx6/faux/fmtwriter"
	"github.com/influx6/moz"
	"github.com/influx6/moz/ast"
	"github.com/influx6/moz/gen"
)

var (
	_ = moz.RegisterAnnotation("templaterTypesFor", TemplaterTypesForAnnotationGenerator)
	_ = moz.RegisterAnnotation("templaterTypesFor", TemplaterStructTypesForAnnotationGenerator)
	_ = moz.RegisterAnnotation("templaterTypesFor", TemplaterPackageTypesForAnnotationGenerator)
	_ = moz.RegisterAnnotation("templaterTypesFor", TemplaterInterfaceTypesForAnnotationGenerator)
)

// TemplaterStructTypesForAnnotationGenerator defines a struct level annotation generator which builds a go package in
// root of the package by using the content it receives from the annotation has a template for its output.
// package.
// Templater provides access to typenames by providing a "sel" function that gives you access to all
// arguments provided by the associated Annotation "templaterForTypes", which provides description of
// the filename, and the types to be used to replace the generic placeholders.
//
// Annotation: @templaterTypesFor
//
// Example:
// 1. Create a template that uses the "Go" generator, identified with the id "Mob" which will
// generate template for all types by using a template from a @templater with id of "Mob", define
// @templater anywhere either in package, struct, type or interface level.
//
// @templater(id => Mob, gen => Go, {
//
//   func Add(m {{sel TYPE1}}, n {{sel TYPE2}}) {{sel TYPE3}} {
//
//   }
//
// })
//
// 2. Add @templaterTypesFor annotation on any level (Type, Struct, Interface, Package) to have the code
// generated from the details provided.
//
// @templaterTypesFor(id => Mob, filename => bob_gen.go, TYPE1 => int32, TYPE2 => int32, TYPE3 => int64)
// @templaterTypesFor(id => Mob, filename => bib_gen.go, TYPE1 => int, TYPE2 => int, TYPE3 => int64)
//
func TemplaterStructTypesForAnnotationGenerator(toDir string, an ast.AnnotationDeclaration, ty ast.StructDeclaration, pkgDeclr ast.PackageDeclaration, pkg ast.Package) ([]gen.WriteDirective, error) {
	return handleGeneration(toDir, an, pkgDeclr, pkg, struct {
		Annotation  ast.AnnotationDeclaration
		PkgDeclr    ast.PackageDeclaration
		Package     ast.Package
		StructDeclr ast.StructDeclaration
	}{
		PkgDeclr:    pkgDeclr,
		Annotation:  an,
		Package:     pkg,
		StructDeclr: ty,
	})
}

// TemplaterInterfaceTypesForAnnotationGenerator defines a package level annotation generator which builds a go package in
// root of the package by using the content it receives from the annotation has a template for its output.
// package.
// Templater provides access to typenames by providing a "sel" function that gives you access to all
// arguments provided by the associated Annotation "templaterForTypes", which provides description of
// the filename, and the types to be used to replace the generic placeholders.
//
// Annotation: @templaterTypesFor
//
// Example:
// 1. Create a template that uses the "Go" generator, identified with the id "Mob" which will
// generate template for all types by using a template from a @templater with id of "Mob", define
// @templater anywhere either in package, struct, type or interface level.
//
// @templater(id => Mob, gen => Go, {
//
//   func Add(m {{sel TYPE1}}, n {{sel TYPE2}}) {{sel TYPE3}} {
//
//   }
//
// })
//
// 2. Add @templaterTypesFor annotation on any level (Type, Struct, Interface, Package) to have the code
// generated from the details provided.
//
// @templaterTypesFor(id => Mob, filename => bob_gen.go, TYPE1 => int32, TYPE2 => int32, TYPE3 => int64)
// @templaterTypesFor(id => Mob, filename => bib_gen.go, TYPE1 => int, TYPE2 => int, TYPE3 => int64)
//
func TemplaterInterfaceTypesForAnnotationGenerator(toDir string, an ast.AnnotationDeclaration, ty ast.InterfaceDeclaration, pkgDeclr ast.PackageDeclaration, pkg ast.Package) ([]gen.WriteDirective, error) {
	return handleGeneration(toDir, an, pkgDeclr, pkg, struct {
		Annotation     ast.AnnotationDeclaration
		PkgDeclr       ast.PackageDeclaration
		Package        ast.Package
		InterfaceDeclr ast.InterfaceDeclaration
	}{
		PkgDeclr:       pkgDeclr,
		Annotation:     an,
		Package:        pkg,
		InterfaceDeclr: ty,
	})
}

// TemplaterPackageTypesForAnnotationGenerator defines a package level annotation generator which builds a go package in
// root of the package by using the content it receives from the annotation has a template for its output.
// package.
// Templater provides access to typenames by providing a "sel" function that gives you access to all
// arguments provided by the associated Annotation "templaterForTypes", which provides description of
// the filename, and the types to be used to replace the generic placeholders.
//
// Annotation: @templaterTypesFor
//
// Example:
// 1. Create a template that uses the "Go" generator, identified with the id "Mob" which will
// generate template for all types by using a template from a @templater with id of "Mob", define
// @templater anywhere either in package, struct, type or interface level.
//
// @templater(id => Mob, gen => Go, {
//
//   func Add(m {{sel TYPE1}}, n {{sel TYPE2}}) {{sel TYPE3}} {
//
//   }
//
// })
//
// 2. Add @templaterTypesFor annotation on any level (Type, Struct, Interface, Package) to have the code
// generated from the details provided.
//
// @templaterTypesFor(id => Mob, filename => bob_gen.go, TYPE1 => int32, TYPE2 => int32, TYPE3 => int64)
// @templaterTypesFor(id => Mob, filename => bib_gen.go, TYPE1 => int, TYPE2 => int, TYPE3 => int64)
//
func TemplaterPackageTypesForAnnotationGenerator(toDir string, an ast.AnnotationDeclaration, pkgDeclr ast.PackageDeclaration, pkg ast.Package) ([]gen.WriteDirective, error) {
	return handleGeneration(toDir, an, pkgDeclr, pkg, struct {
		Annotation ast.AnnotationDeclaration
		PkgDeclr   ast.PackageDeclaration
		Package    ast.Package
	}{
		PkgDeclr:   pkgDeclr,
		Annotation: an,
		Package:    pkg,
	})
}

// TemplaterTypesForAnnotationGenerator defines a package level annotation generator which builds a go package in
// root of the package by using the content it receives from the annotation has a template for its output.
// package.
// Templater provides access to typenames by providing a "sel" function that gives you access to all
// arguments provided by the associated Annotation "templaterForTypes", which provides description of
// the filename, and the types to be used to replace the generic placeholders.
//
// Annotation: @templaterTypesFor
//
// Example:
// 1. Create a template that uses the "Go" generator, identified with the id "Mob" which will
// generate template for all types by using a template from a @templater with id of "Mob", define
// @templater anywhere either in package, struct, type or interface level.
//
// @templater(id => Mob, gen => Go, {
//
//   func Add(m {{sel TYPE1}}, n {{sel TYPE2}}) {{sel TYPE3}} {
//
//   }
//
// })
//
// 2. Add @templaterTypesFor annotation on any level (Type, Struct, Interface, Package) to have the code
// generated from the details provided.
//
// @templaterTypesFor(id => Mob, filename => bob_gen.go, TYPE1 => int32, TYPE2 => int32, TYPE3 => int64)
// @templaterTypesFor(id => Mob, filename => bib_gen.go, TYPE1 => int, TYPE2 => int, TYPE3 => int64)
//
func TemplaterTypesForAnnotationGenerator(toDir string, an ast.AnnotationDeclaration, ty ast.TypeDeclaration, pkgDeclr ast.PackageDeclaration, pkg ast.Package) ([]gen.WriteDirective, error) {
	return handleGeneration(toDir, an, pkgDeclr, pkg, struct {
		Annotation ast.AnnotationDeclaration
		TypeDeclr  ast.TypeDeclaration
		Package    ast.Package
		PkgDeclr   ast.PackageDeclaration
	}{
		Annotation: an,
		TypeDeclr:  ty,
		PkgDeclr:   pkgDeclr,
		Package:    pkg,
	})
}

func handleGeneration(toDir string, an ast.AnnotationDeclaration, pkgDeclr ast.PackageDeclaration, pkg ast.Package, binding interface{}) ([]gen.WriteDirective, error) {
	templaterID, ok := an.Params["id"]
	if !ok {
		return nil, errors.New("No templater id provided")
	}

	// Get all templaters AnnotationDeclaration.
	templaters := pkg.AnnotationsFor("templater")
	var targetTemplater ast.AnnotationDeclaration

	// Search for templater with associated ID, if not found, return error, if multiple found, use the first.
	for _, targetTemplater = range templaters {
		if targetTemplater.Params["id"] != templaterID {
			continue
		}

		break
	}

	var templateData string

	switch len(targetTemplater.Template) == 0 {
	case true:
		templateFilePath, dok := targetTemplater.Params["file"]
		if !dok && targetTemplater.Template == "" {
			return nil, errors.New("Expected Template from annotation or provide `file => 'path_to_template`")
		}

		baseDir := filepath.Dir(pkgDeclr.FilePath)
		templateFile := filepath.Join(baseDir, templateFilePath)

		data, err := ioutil.ReadFile(templateFile)
		if err != nil {
			return nil, fmt.Errorf("Failed to find template file: %+q", err)
		}

		templateData = string(data)
	case false:
		templateData = targetTemplater.Template
	}

	var directives []gen.WriteDirective

	genName := strings.ToLower(targetTemplater.Params["gen"])

	fileName, ok := an.Params["filename"]
	if !ok {
		fileName = fmt.Sprintf("%s_impl_gen.go", strings.ToLower(an.Name))
	}

	typeGen := gen.Block(gen.SourceTextWith(templateData, template.FuncMap{
		"sel":                    an.Param,
		"params":                 an.Param,
		"attrs":                  an.Attr,
		"hasArg":                 an.HasArg,
		"annotationDefer":        func() bool { return an.Defer },
		"annotationTemplate":     func() string { return an.Template },
		"annotationParams":       func() map[string]string { return an.Params },
		"annotationAttrs":        func() map[string]interface{} { return an.Attrs },
		"annotationArguments":    func() []string { return an.Arguments },
		"targetSel":              targetTemplater.Param,
		"targetParams":           targetTemplater.Param,
		"targetAttrs":            targetTemplater.Attr,
		"targetHasArg":           targetTemplater.HasArg,
		"targetDefer":            func() bool { return targetTemplater.Defer },
		"targetTemplate":         func() string { return targetTemplater.Template },
		"targetArguments":        func() []string { return targetTemplater.Arguments },
		"targetAnnotationParams": func() map[string]string { return targetTemplater.Params },
		"targetAnnotationAttrs":  func() map[string]interface{} { return targetTemplater.Attrs },
	}, binding))

	switch genName {
	case "partial_test.go":

		var packageName string

		switch len(an.Params["packageName"]) == 0 {
		case true:
			packageName = ast.WhichPackage(toDir, pkg)
		case false:
			packageName = targetTemplater.Params["packageName"]
		}

		packageName = fmt.Sprintf("%s_test", packageName)

		pkgGen := gen.Block(

			gen.Package(
				gen.Name(packageName),
				typeGen,
			),
		)

		directives = append(directives, gen.WriteDirective{
			FileName:     fileName,
			DontOverride: true,
			Writer:       fmtwriter.New(pkgGen, true, true),
		})

	case "partial.go":

		var packageName string

		switch len(an.Params["packageName"]) == 0 {
		case true:
			packageName = ast.WhichPackage(toDir, pkg)
		case false:
			packageName = targetTemplater.Params["packageName"]
		}

		pkgGen := gen.Block(

			gen.Package(
				gen.Name(packageName),
				typeGen,
			),
		)

		directives = append(directives, gen.WriteDirective{
			FileName:     fileName,
			DontOverride: true,
			Writer:       fmtwriter.New(pkgGen, true, true),
		})

	case "go":
		directives = append(directives, gen.WriteDirective{
			FileName:     fileName,
			DontOverride: true,

			Writer: fmtwriter.New(typeGen, true, true),
		})

	default:
		directives = append(directives, gen.WriteDirective{
			Writer:       typeGen,
			DontOverride: true,
			FileName:     fileName,
		})
	}

	return directives, nil
}
