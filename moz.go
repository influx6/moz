package moz

import (
	"github.com/influx6/faux/metrics"
	"github.com/influx6/moz/ast"
)

var (
	//Annotations is a package level registry exposed by moz to provide a central
	// storage for all annotation registration.
	annotations = ast.NewAnnotationRegistry()
)

// CopyAnnotations copies all annotations from the provided AnnotationRegistry.
func CopyAnnotations(registry *ast.AnnotationRegistry) {
	annotations.Copy(registry, ast.OursOverTheirs)
}

// CopyAnnotationsTo copies all default annotation to the provided AnnotationRegistry.
func CopyAnnotationsTo(registry *ast.AnnotationRegistry) *ast.AnnotationRegistry {
	registry.Copy(annotations, ast.OursOverTheirs)
	return registry
}

// RegisterAnnotation which adds the generator depending on it's type into the appropriate
// registry. It only supports  the following generators:
// 1. TypeAnnotationGenerator (see Package http://github.com/influx6/moz/ast.TypeAnnotationGenerator)
// 2. StructAnnotationGenerator (see Package http://github.com/influx6/moz/ast.StructAnnotationGenerator)
// 3. InterfaceAnnotationGenerator (see Package http://github.com/influx6/moz/ast.InterfaceAnnotationGenerator)
// 4. PackageAnnotationGenerator (see Package http://github.com/influx6/moz/ast.PackageAnnotationGenerator)
// Any other type will cause the return of an error.
func RegisterAnnotation(name string, generator interface{}) bool {
	if err := annotations.Register(name, generator); err != nil {
		panic(err.Error())
	}

	return true
}

// MustParseWith calls the ParseWith method to attempt to parse the ast.PackageDeclarations
// and panics if it encounters an error.
func MustParseWith(toDir string, log metrics.Metrics, provider *ast.AnnotationRegistry, forceWrite bool, packageDeclrs ...ast.Package) {
	if err := ParseWith(toDir, log, provider, forceWrite, packageDeclrs...); err != nil {
		panic(err)
	}
}

// ParseWith takes the provided package declarations and annotation registry and attempts
// parsing all internals structuers with the appropriate generators suited to the type and annotations.
func ParseWith(toDir string, log metrics.Metrics, provider *ast.AnnotationRegistry, forceWrite bool, packageDeclrs ...ast.Package) error {
	return ast.Parse(toDir, log, provider, forceWrite, packageDeclrs...)
}

// MustParse calls the Parse method to attempt to parse the ast.PackageDeclarations
// and panics if it encounters an error.
func MustParse(toDir string, log metrics.Metrics, forceWrite bool, packageDeclrs ...ast.Package) {
	if err := Parse(toDir, log, forceWrite, packageDeclrs...); err != nil {
		panic(err)
	}
}

// Parse takes the provided package declarations and the default Annotations registry and attempts
// parsing all internals structuers with the appropriate generators suited to the type and annotations.
func Parse(toDir string, log metrics.Metrics, forceWrite bool, packageDeclrs ...ast.Package) error {
	return ast.Parse(toDir, log, annotations, forceWrite, packageDeclrs...)
}
