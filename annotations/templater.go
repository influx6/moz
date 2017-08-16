package annotations

import (
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
/*
@templater({

  func Add(m TYPE1, n TYPE2) TYPE3 {

  }

})
*/
func TemplaterAnnotationGenerator(toDir string, an ast.AnnotationDeclaration, pkg ast.PackageDeclaration) ([]gen.WriteDirective, error) {

	return []gen.WriteDirective{}, nil
}
