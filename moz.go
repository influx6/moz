package moz

import (
	"errors"

	"github.com/influx6/moz/ast"
)

var (
	//Annotations is a package level registry exposed by moz to provide a central
	// storage for all annotation registration.
	Annotations = ast.NewAnnotationRegistry()
)

// Parse takes the provided package declrations parsing all internals with the
// appropriate generators suited to the type and annotations.
func Parse(packageDeclrs ...ast.PackageDeclaration) error {

	return nil
}

// ParseDeclr takes the provided package declration and runs the
// appropriate generators suited to apply the annotations related to it.
func ParseDeclr(packageDeclrs ast.PackageDeclaration) error {

	return nil
}

// RegisterAnnotation which adds the generator depending on it's type into the appropriate
// registry. It only supports  the following generators:
// 1. TypeAnnotationGenerator (see Package http://github.com/influx6/moz/ast.TypeAnnotationGenerator)
// 2. StructAnnotationGenerator (see Package http://github.com/influx6/moz/ast.StructAnnotationGenerator)
// 3. InterfaceAnnotationGenerator (see Package http://github.com/influx6/moz/ast.InterfaceAnnotationGenerator)
// 4. PackageAnnotationGenerator (see Package http://github.com/influx6/moz/ast.PackageAnnotationGenerator)
// Any other type will cause the return of an error.
func RegisterAnnotation(name string, generator interface{}) error {
	switch gen := generator.(type) {
	case ast.PackageAnnotationGenerator:
		Annotations.RegisterPackage(name, gen)
	case ast.TypeAnnotationGenerator:
		Annotations.RegisterType(name, gen)
	case ast.StructAnnotationGenerator:
		Annotations.RegisterStructType(name, gen)
	case ast.InterfaceAnnotationGenerator:
		Annotations.RegisterInterfaceType(name, gen)
	default:
		return errors.New("Generator type not supported")
	}

	return nil
}

// MustRegisterAnnotation which adds the generator depending on it's type into the appropriate
// registry. It only supports  the following generators:
// 1. TypeAnnotationGenerator (see Package http://github.com/influx6/moz/ast.TypeAnnotationGenerator)
// 2. StructAnnotationGenerator (see Package http://github.com/influx6/moz/ast.StructAnnotationGenerator)
// 3. InterfaceAnnotationGenerator (see Package http://github.com/influx6/moz/ast.InterfaceAnnotationGenerator)
// 4. PackageAnnotationGenerator (see Package http://github.com/influx6/moz/ast.PackageAnnotationGenerator)
// Any other type will cause a panic.
func MustRegisterAnnotation(name string, generator interface{}) {
	if err := RegisterAnnotation(name, generator); err != nil {
		panic(err)
	}
}
