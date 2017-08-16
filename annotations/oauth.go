package annotations

import (
	"github.com/influx6/moz"
	"github.com/influx6/moz/ast"
	"github.com/influx6/moz/gen"
)

var (
	_ = moz.RegisterAnnotation("oauth", OAuthAnnotationGenerator)
)

// OAuthAnnotationGenerator defines a code generator for package level generator for creating
// a http package for providing oauth authentication.
func OAuthAnnotationGenerator(an ast.AnnotationDeclaration, pkg ast.PackageDeclaration) ([]gen.WriteDirective, error) {
	var directives []gen.WriteDirective

	return directives, nil
}
