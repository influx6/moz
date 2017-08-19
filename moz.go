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
func CopyAnnotationsTo(registry *ast.AnnotationRegistry) {
	registry.Copy(annotations, ast.OursOverTheirs)
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
func MustParseWith(toDir string, log metrics.Metrics, provider *ast.AnnotationRegistry, forceWrite bool, packageDeclrs ...ast.PackageDeclaration) {
	if err := ParseWith(toDir, log, provider, forceWrite, packageDeclrs...); err != nil {
		panic(err)
	}
}

// ParseWith takes the provided package declarations and annotation registry and attempts
// parsing all internals structuers with the appropriate generators suited to the type and annotations.
func ParseWith(toDir string, log metrics.Metrics, provider *ast.AnnotationRegistry, forceWrite bool, packageDeclrs ...ast.PackageDeclaration) error {
	return ast.Parse(toDir, log, provider, forceWrite, packageDeclrs...)
}

// MustParse calls the Parse method to attempt to parse the ast.PackageDeclarations
// and panics if it encounters an error.
func MustParse(toDir string, log metrics.Metrics, forceWrite bool, packageDeclrs ...ast.PackageDeclaration) {
	if err := Parse(toDir, log, forceWrite, packageDeclrs...); err != nil {
		panic(err)
	}
}

// Parse takes the provided package declarations and the default Annotations registry and attempts
// parsing all internals structuers with the appropriate generators suited to the type and annotations.
func Parse(toDir string, log metrics.Metrics, forceWrite bool, packageDeclrs ...ast.PackageDeclaration) error {
	return ast.Parse(toDir, log, annotations, forceWrite, packageDeclrs...)
}

// MustWriteDirectives calls the WriteDirectives method to attempt to parse the ast.PackageDeclarations
// and panics if it encounters an error.
// func MustWriteDirectives(log metrics.Metrics, rootDir string, directives ...gen.WriteDirective) {
// 	if err := WriteDirectives(log, rootDir, directives...); err != nil {
// 		panic(err)
// 	}
// }

// // WriteDirectives defines a funtion to sync the slices of WriteDirectives into a giving directory
// // root.
// func WriteDirectives(log metrics.Metrics, rootDir string, directives ...gen.WriteDirective) error {
// 	{
// 	directiveloop:
// 		for _, directive := range directives {
//
// 			if filepath.IsAbs(directive.Dir) {
// 				log.Emit(stdout.Error("gen.WriteDirectiveError: Expected relative Dir path not absolute").
// 					With("root-dir", rootDir).With("directive-dir", directive.Dir))
//
// 				continue directiveloop
// 			}
//
// 			namedFileDir := filepath.Join(rootDir, directive.Dir)
// 			namedFile := filepath.Join(namedFileDir, directive.FileName)
//
// 			if err := os.MkdirAll(namedFileDir, 0700); err != nil && err != os.ErrExist {
// 				return err
// 			}
//
// 			newFile, err := os.Open(namedFile)
// 			if err != nil {
// 				log.Emit(stdout.Error("IOError: Unable to create file").
// 					With("file", namedFile).With("dir", namedFileDir).With("error", err))
// 				return err
// 			}
//
// 			if _, err := directive.Writer.WriteTo(newFile); err != nil && err != io.EOF {
// 				log.Emit(stdout.Error("IOError: Unable to write to file").
// 					With("file", namedFile).With("dir", namedFileDir).With("error", err))
//
// 				newFile.Close()
//
// 				return err
// 			}
//
// 			log.Emit(stdout.Info("Directive Resolved").With("file", namedFile).With("dir", namedFileDir))
//
// 			// Close giving file
// 			newFile.Close()
// 		}
// 	}
//
// 	return nil
// }
