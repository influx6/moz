package db

import (
	"github.com/influx6/moz"
	"github.com/influx6/moz/ast"
	"github.com/influx6/moz/gen"
)

var (
	_ = moz.RegisterAnnotation("mongo", MongoAnnotationGenerator)
)

// MongoAnnotationGenerator defines a code generator for struct declarations that generate a
// mongo CRUD code for the use of mongodb as the underline db store.
func MongoAnnotationGenerator(an ast.AnnotationDeclaration, declr ast.StructDeclaration, pkg ast.PackageDeclaration) ([]gen.WriteDirective, error) {
	var directives []gen.WriteDirective

	return directives, nil
}
