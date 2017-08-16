package annotations

import (
	"errors"
	"fmt"

	"github.com/influx6/moz"
	"github.com/influx6/moz/ast"
	"github.com/influx6/moz/gen"
)

var (
	_ = moz.RegisterAnnotation("templater", TemplaterAnnotationGenerator)
)

// TemplaterAnnotationGenerator defines a package level annotation generator which builds a go package in
// root of the package by using the content it receives from the annotation has a template for its output.
// package.
// Annotation: @templater
// Example:
// Create a template that uses the "Go" generator, identified with the id "Mob" which will
// generate template for all types with id of "Mob".
/*
@templater(id => Mob, gen => Go, {

  func Add(m TYPE1, n TYPE2) TYPE3 {

  }

})
@templaterTypesFor(id => Mob, TYPE1 => int32, TYPE2 => int32, TYPE3 => int64)
@templaterTypesFor(id => Mob, TYPE1 => int, TYPE2 => int, TYPE3 => int64)
*/
func TemplaterAnnotationGenerator(toDir string, an ast.AnnotationDeclaration, pkg ast.PackageDeclaration) ([]gen.WriteDirective, error) {
	if len(an.Arguments) < 2 {
		return nil, errors.New("Expected 2 arguments as follows: 'id => IDOFTemplate, generator => Go'")
	}

	var types []ast.AnnotationDeclaration

	for _, item := range pkg.AnnotationsFor("@templaterTypesFor") {
		if item.Params["id"] != an.Params["id"] {
			continue
		}

		types = append(types, item)
	}

	fmt.Printf("Params: %+q\n", an.Params)
	fmt.Printf("Types: %+q\n", types)

	return []gen.WriteDirective{}, nil
}
