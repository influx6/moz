package moz

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/influx6/faux/sink"
	"github.com/influx6/faux/sink/sinks"
	"github.com/influx6/moz/ast"
	"github.com/influx6/moz/gen"
)

var (
	//Annotations is a package level registry exposed by moz to provide a central
	// storage for all annotation registration.
	Annotations = ast.NewAnnotationRegistry()
)

// RegisterAnnotation which adds the generator depending on it's type into the appropriate
// registry. It only supports  the following generators:
// 1. TypeAnnotationGenerator (see Package http://github.com/influx6/moz/ast.TypeAnnotationGenerator)
// 2. StructAnnotationGenerator (see Package http://github.com/influx6/moz/ast.StructAnnotationGenerator)
// 3. InterfaceAnnotationGenerator (see Package http://github.com/influx6/moz/ast.InterfaceAnnotationGenerator)
// 4. PackageAnnotationGenerator (see Package http://github.com/influx6/moz/ast.PackageAnnotationGenerator)
// Any other type will cause the return of an error.
func RegisterAnnotation(name string, generator interface{}) bool {
	switch gen := generator.(type) {
	case ast.PackageAnnotationGenerator:
		Annotations.RegisterPackage(name, gen)
		break
	case func(ast.AnnotationDeclaration, ast.PackageDeclaration) ([]gen.WriteDirective, error):
		Annotations.RegisterPackage(name, gen)
		break
	case ast.TypeAnnotationGenerator:
		Annotations.RegisterType(name, gen)
		break
	case func(ast.AnnotationDeclaration, ast.TypeDeclaration, ast.PackageDeclaration) ([]gen.WriteDirective, error):
		Annotations.RegisterType(name, gen)
		break
	case ast.StructAnnotationGenerator:
		Annotations.RegisterStructType(name, gen)
		break
	case func(ast.AnnotationDeclaration, ast.StructDeclaration, ast.PackageDeclaration) ([]gen.WriteDirective, error):
		Annotations.RegisterStructType(name, gen)
		break
	case ast.InterfaceAnnotationGenerator:
		Annotations.RegisterInterfaceType(name, gen)
		break
	case func(ast.AnnotationDeclaration, ast.InterfaceDeclaration, ast.PackageDeclaration) ([]gen.WriteDirective, error):
		Annotations.RegisterInterfaceType(name, gen)
		break
	default:
		panic(fmt.Errorf("Generator type for %q not supported: %#v", name, generator))
	}

	return true
}

// MustParseWith calls the ParseWith method to attempt to parse the ast.PackageDeclarations
// and panics if it encounters an error.
func MustParseWith(log sink.Sink, provider *ast.AnnotationRegistry, packageDeclrs ...ast.PackageDeclaration) {
	if err := ParseWith(log, provider, packageDeclrs...); err != nil {
		panic(err)
	}
}

// ParseWith takes the provided package declarations and annotation registry and attempts
// parsing all internals structuers with the appropriate generators suited to the type and annotations.
func ParseWith(log sink.Sink, provider *ast.AnnotationRegistry, packageDeclrs ...ast.PackageDeclaration) error {
	return ast.Parse(log, provider, packageDeclrs...)
}

// MustParse calls the Parse method to attempt to parse the ast.PackageDeclarations
// and panics if it encounters an error.
func MustParse(log sink.Sink, packageDeclrs ...ast.PackageDeclaration) {
	if err := Parse(log, packageDeclrs...); err != nil {
		panic(err)
	}
}

// Parse takes the provided package declarations and the default Annotations registry and attempts
// parsing all internals structuers with the appropriate generators suited to the type and annotations.
func Parse(log sink.Sink, packageDeclrs ...ast.PackageDeclaration) error {
	return ast.Parse(log, Annotations, packageDeclrs...)
}

// MustWriteDirectives calls the WriteDirectives method to attempt to parse the ast.PackageDeclarations
// and panics if it encounters an error.
func MustWriteDirectives(log sink.Sink, rootDir string, directives ...gen.WriteDirective) {
	if err := WriteDirectives(log, rootDir, directives...); err != nil {
		panic(err)
	}
}

// WriteDirectives defines a funtion to sync the slices of WriteDirectives into a giving directory
// root.
func WriteDirectives(log sink.Sink, rootDir string, directives ...gen.WriteDirective) error {

	{
	directiveloop:
		for _, directive := range directives {

			if filepath.IsAbs(directive.Dir) {
				log.Emit(sinks.Error("gen.WriteDirectiveError: Expected relative Dir path not absolute").
					With("root-dir", rootDir).With("directive-dir", directive.Dir))

				continue directiveloop
			}

			namedFileDir := filepath.Join(rootDir, directive.Dir)
			namedFile := filepath.Join(namedFileDir, directive.FileName)

			if err := os.MkdirAll(namedFileDir, 0700); err != nil && err != os.ErrExist {
				return err
			}

			newFile, err := os.Open(namedFile)
			if err != nil {
				log.Emit(sinks.Error("IOError: Unable to create file").
					With("file", namedFile).With("dir", namedFileDir).With("error", err))
				return err
			}

			if _, err := directive.Writer.WriteTo(newFile); err != nil && err != io.EOF {
				log.Emit(sinks.Error("IOError: Unable to write to file").
					With("file", namedFile).With("dir", namedFileDir).With("error", err))

				newFile.Close()

				return err
			}

			log.Emit(sinks.Info("Directive Resolved").With("file", namedFile).With("dir", namedFileDir))

			// Close giving file
			newFile.Close()
		}
	}

	return nil
}
