package moz

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/influx6/faux/sink"
	"github.com/influx6/faux/sink/sinks"
	"github.com/influx6/moz/ast"
)

const (
	annotationFileFormat = "%s_annotation_%s_gen.%s"
)

var (
	//Annotations is a package level registry exposed by moz to provide a central
	// storage for all annotation registration.
	Annotations = ast.NewAnnotationRegistry()
)

// Parse takes the provided package declrations parsing all internals with the
// appropriate generators suited to the type and annotations.
func Parse(log sink.Sink, packageDeclrs ...ast.PackageDeclaration) error {
	{
	parseloop:
		for _, pkg := range packageDeclrs {
			wdrs, err := Annotations.ParseDeclr(pkg)
			if err != nil {
				log.Emit(sinks.Error("ParseFailure: Package %q", pkg.Package).
					With("error", err).With("package", pkg.Package))
				continue
			}

			fileName := strings.TrimSuffix(pkg.File, filepath.Ext(pkg.File))

			for _, item := range wdrs {
				extension := item.WriteDirective.Ext

				if extension == "" {
					extension = "go"
				} else {
					extension = strings.TrimPrefix(extension, ".")
				}

				dir := filepath.Dir(pkg.FilePath)
				annotation := strings.ToLower(item.Annotation)
				annotationFile := fmt.Sprintf(annotationFileFormat, fileName, annotation, extension)

				if item.Dir == "" {
					newDirFile := filepath.Join(dir, annotationFile)
					newFile, err := os.Open(newDirFile)
					if err != nil {
						log.Emit(sinks.Error("IOError: Unable to create file").With("file", newFile).With("package", pkg.Package).With("error", err))
						return err
					}

					if _, err := item.Writer.WriteTo(newFile); err != nil && err != io.EOF {
						newFile.Close()
						log.Emit(sinks.Error("IOError: Unable to write content to file").
							With("file", newFile).With("error", err).With("package", pkg.Package))
						return err
					}

					log.Emit(sinks.Info("Annotation Resolved").With("annotation", item.Annotation).
						With("package", pkg.Package).With("file", pkg.File).With("generated-file", newDirFile))

					newFile.Close()
					continue
				}

				if filepath.IsAbs(item.Dir) {
					log.Emit(sinks.Error("WriteDirectiveError: Expected relative Dir path not absolute").
						With("package", pkg.Package).With("directive-dir", item.Dir).With("pkg", pkg))

					continue parseloop
				}

				newDir := filepath.Join(dir, item.Dir)
				newDirFile := filepath.Join(newDir, annotationFile)

				if err := os.MkdirAll(newDir, 0700); err != nil && err != os.ErrExist {
					return err
				}

				newFile, err := os.Open(newDirFile)
				if err != nil {
					log.Emit(sinks.Error("IOError: Unable to create file").
						With("file", newFile).With("error", err))
					return err
				}

				if _, err := item.Writer.WriteTo(newFile); err != nil && err != io.EOF {
					newFile.Close()
					log.Emit(sinks.Error("IOError: Unable to write content to file").
						With("file", newFile).With("error", err))
					return err
				}

				log.Emit(sinks.Info("Annotation Resolved").With("annotation", item.Annotation).
					With("package", pkg.Package).With("file", pkg.File).With("generated-file", newDirFile))

				newFile.Close()
			}
		}

	}

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
