package ast

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/influx6/faux/metrics"
	"github.com/influx6/faux/metrics/sentries/stdout"
	"github.com/influx6/moz/gen"
)

const (
	annotationFileFormat    = "%s_annotation_%s.%s"
	altAnnotationFileFormat = "%s_annotation_%s"
)

// Contains giving sets of variables exposing sytem GOPATH and GOPATHSRC.
var (
	GoPath     = os.Getenv("GOPATH")
	GoSrcPath  = filepath.Join(GoPath, "src")
	annotation = regexp.MustCompile("@(\\w+)(\\(.+\\))?")
)

// PackageDir turns a given go source file into a appropriate structure which will be
// used to generate the needed manifests for a resource shell.
func PackageDir(file string, mode parser.Mode) (*token.FileSet, map[string]*ast.Package, error) {
	tokens := token.NewFileSet()
	nodes, err := parser.ParseDir(tokens, file, nil, mode)
	if err != nil {
		return nil, nil, err
	}

	return tokens, nodes, nil
}

// PackageFile turns a given go source file into a appropriate structure which will be
// used to generate the needed manifests for a resource shell.
func PackageFile(file string, mode parser.Mode) (*token.FileSet, *ast.File, error) {
	tokens := token.NewFileSet()
	nodes, err := parser.ParseFile(tokens, file, nil, mode)
	if err != nil {
		return nil, nil, err
	}

	return tokens, nodes, nil
}

// AnnotationDeclaration defines a annotation type which holds detail about a giving annotation.
type AnnotationDeclaration struct {
	Name      string   `json:"name"`
	Arguments []string `json:"arguments"`
}

// PackageDeclaration defines a type which holds details relating to annotations declared on a
// giving package.
type PackageDeclaration struct {
	Package     string                  `json:"package"`
	Path        string                  `json:"path"`
	FilePath    string                  `json:"filepath"`
	File        string                  `json:"file"`
	Annotations []AnnotationDeclaration `json:"annotations"`
	Types       []TypeDeclaration       `json:"types"`
	Structs     []StructDeclaration     `json:"structs"`
	Interfaces  []InterfaceDeclaration  `json:"interfaces"`
}

// InterfaceDeclaration defines a type which holds annotation data for a giving interface type declaration.
type InterfaceDeclaration struct {
	LineNumber  int                     `json:"line_number"`
	Column      int                     `json:"column"`
	Package     string                  `json:"package"`
	Path        string                  `json:"path"`
	FilePath    string                  `json:"filepath"`
	File        string                  `json:"file"`
	Interface   *ast.InterfaceType      `json:"interface"`
	Object      *ast.TypeSpec           `json:"object"`
	Annotations []AnnotationDeclaration `json:"annotations"`
}

// StructDeclaration defines a type which holds annotation data for a giving struct type declaration.
type StructDeclaration struct {
	LineNumber  int                     `json:"line_number"`
	Column      int                     `json:"column"`
	Package     string                  `json:"package"`
	Path        string                  `json:"path"`
	FilePath    string                  `json:"filepath"`
	File        string                  `json:"file"`
	Struct      *ast.StructType         `json:"struct"`
	Object      *ast.TypeSpec           `json:"object"`
	Annotations []AnnotationDeclaration `json:"annotations"`
}

// TypeDeclaration defines a type which holds annotation data for a giving type declaration.
type TypeDeclaration struct {
	LineNumber  int                     `json:"line_number"`
	Column      int                     `json:"column"`
	Package     string                  `json:"package"`
	Path        string                  `json:"path"`
	FilePath    string                  `json:"filepath"`
	File        string                  `json:"file"`
	Object      *ast.TypeSpec           `json:"object"`
	Annotations []AnnotationDeclaration `json:"annotations"`
}

//===========================================================================================================

// ParseAnnotations parses the package which generates a series of ast with associated
// annotation for processing.
func ParseAnnotations(log metrics.Metrics, dir string) ([]PackageDeclaration, error) {
	tokenFiles, packages, err := PackageDir(dir, parser.ParseComments)
	if err != nil {
		log.Emit(stdout.Error(err).With("message", "Failed to parse directory").With("dir", dir))
		return nil, err
	}

	var packageDeclrs []PackageDeclaration

	for _, pkg := range packages {
		for path, file := range pkg.Files {

			var packageDeclr PackageDeclaration
			packageDeclr.Package = pkg.Name
			packageDeclr.FilePath = path

			if relPath, err := filepath.Rel(GoSrcPath, path); err == nil {
				packageDeclr.Path = filepath.Dir(relPath)
				packageDeclr.File = filepath.Base(relPath)
			}

			// fmt.Printf("Real: %+q\n", packageDeclr)

			for _, comment := range file.Doc.List {
				text := strings.TrimPrefix(comment.Text, "//")

				if !annotation.MatchString(text) {
					continue
				}

				annons := annotation.FindStringSubmatch(text)

				log.Emit(stdout.Info("Annotation in Package comments").
					With("dir", dir).
					With("annotation", annons).
					With("documentation", file.Doc).
					With("comment", comment))

				if len(annons) > 1 {
					var arguments []string

					args := strings.TrimSuffix(strings.TrimPrefix(annons[2], "("), ")")
					for _, elem := range strings.Split(args, ",") {
						elem = strings.TrimSpace(elem)

						if unquoted, err := strconv.Unquote(elem); err == nil {
							arguments = append(arguments, unquoted)
							continue
						}

						arguments = append(arguments, elem)
					}

					packageDeclr.Annotations = append(packageDeclr.Annotations, AnnotationDeclaration{
						Name:      annons[1],
						Arguments: arguments,
					})

					continue
				}

				packageDeclr.Annotations = append(packageDeclr.Annotations, AnnotationDeclaration{
					Name: annons[1],
				})
			}

			// Collect and categorize annotations in types and their fields.
			for _, declr := range file.Decls {

				switch rdeclr := declr.(type) {
				case *ast.GenDecl:

					var annotations []AnnotationDeclaration
					if rdeclr.Doc != nil {
						for _, comment := range rdeclr.Doc.List {
							text := strings.TrimPrefix(comment.Text, "//")

							if !annotation.MatchString(text) {
								continue
							}

							annons := annotation.FindStringSubmatch(text)
							log.Emit(stdout.Info("Annotation in Decleration comment %+q", annons).
								With("dir", dir).
								With("comment", comment).
								With("documentation", rdeclr.Doc).
								With("position", rdeclr.Pos()).
								With("token", rdeclr.Tok.String()))

							if len(annons) > 1 {
								var arguments []string

								args := strings.TrimSuffix(strings.TrimPrefix(annons[2], "("), ")")
								for _, elem := range strings.Split(args, ",") {
									elem = strings.TrimSpace(elem)

									if unquoted, err := strconv.Unquote(elem); err == nil {
										arguments = append(arguments, unquoted)
										continue
									}

									arguments = append(arguments, elem)
								}

								annotations = append(annotations, AnnotationDeclaration{
									Name:      annons[1],
									Arguments: arguments,
								})

								continue
							}

							annotations = append(annotations, AnnotationDeclaration{
								Name: annons[1],
							})
						}
					}

					for _, spec := range rdeclr.Specs {
						switch obj := spec.(type) {
						case *ast.TypeSpec:

							tokenPosition := tokenFiles.Position(spec.Pos())

							switch robj := obj.Type.(type) {
							case *ast.StructType:

								log.Emit(stdout.Info("Annotation in Decleration").
									With("Type", "Struct").
									With("Annotations", annotations).
									With("StructName", obj.Name))

								packageDeclr.Structs = append(packageDeclr.Structs, StructDeclaration{
									Object:      obj,
									Struct:      robj,
									Annotations: annotations,
									File:        packageDeclr.File,
									Package:     packageDeclr.Package,
									Path:        packageDeclr.Path,
									FilePath:    packageDeclr.FilePath,
									LineNumber:  tokenPosition.Line,
									Column:      tokenPosition.Column,
								})
								break

							case *ast.InterfaceType:
								log.Emit(stdout.Info("Annotation in Decleration").
									With("Type", "Interface").
									With("Annotations", annotations).
									With("StructName", obj.Name))

								packageDeclr.Interfaces = append(packageDeclr.Interfaces, InterfaceDeclaration{
									Object:      obj,
									Interface:   robj,
									Annotations: annotations,
									File:        packageDeclr.File,
									Package:     packageDeclr.Package,
									Path:        packageDeclr.Path,
									FilePath:    packageDeclr.FilePath,
									LineNumber:  tokenPosition.Line,
									Column:      tokenPosition.Column,
								})
								break
							default:
								log.Emit(stdout.Info("Annotation in Decleration").
									With("Type", "OtherType").
									With("Marker", "NonStruct/NonInterface:Type").
									With("Annotations", annotations).
									With("StructName", obj.Name))

								packageDeclr.Types = append(packageDeclr.Types, TypeDeclaration{
									Object:      obj,
									Annotations: annotations,
									File:        packageDeclr.File,
									Package:     packageDeclr.Package,
									Path:        packageDeclr.Path,
									FilePath:    packageDeclr.FilePath,
									LineNumber:  tokenPosition.Line,
									Column:      tokenPosition.Column,
								})
							}

						case *ast.ImportSpec:
							// Do Nothing.
						}
					}

				case *ast.BadDecl:
					// Do Nothing.
				case *ast.FuncDecl:
					// Do Nothing.
				}
			}

			packageDeclrs = append(packageDeclrs, packageDeclr)
		}
	}

	return packageDeclrs, nil
}

//===========================================================================================================

// Parse takes the provided package declrations parsing all internals with the
// appropriate generators suited to the type and annotations.
func Parse(log metrics.Metrics, provider *AnnotationRegistry, packageDeclrs ...PackageDeclaration) error {
	{
	parseloop:
		for _, pkg := range packageDeclrs {
			wdrs, err := provider.ParseDeclr(pkg)
			if err != nil {
				log.Emit(stdout.Error("ParseFailure: Package %q", pkg.Package).
					With("error", err).With("package", pkg.Package))
				continue
			}

			for _, item := range wdrs {

				if filepath.IsAbs(item.Dir) {
					log.Emit(stdout.Error("gen.WriteDirectiveError: Expected relative Dir path not absolute").
						With("package", pkg.Package).With("directive-dir", item.Dir).With("pkg", pkg))

					continue parseloop
				}

				log.Emit(stdout.Info("Executing WriteDirective").
					With("annotation", item.Annotation).
					With("fileName", item.FileName).
					With("toDir", item.Dir))

				var namedFileDir, namedFile string

				annotation := strings.ToLower(item.Annotation)
				newDir := filepath.Dir(pkg.FilePath)

				if item.Dir != "" {
					namedFileDir = filepath.Join(newDir, item.Dir)
				} else {
					namedFileDir = newDir
				}

				if err := os.MkdirAll(namedFileDir, 0700); err != nil && err != os.ErrExist {
					log.Emit(stdout.Error("IOError: Unable to create writer directory").
						With("dir", namedFileDir).With("error", err))
					return err
				}

				if item.Writer == nil {
					log.Emit(stdout.Info("Annotation Resolved").
						With("annotation", item.Annotation).
						With("dir", namedFileDir))
					continue
				}

				if item.FileName == "" {
					fileName := strings.TrimSuffix(pkg.File, filepath.Ext(pkg.File))
					annotationFile := fmt.Sprintf(annotationFileFormat, annotation, fileName, "go")

					namedFile = filepath.Join(namedFileDir, annotationFile)
				} else {
					// annotationFile := fmt.Sprintf(altAnnotationFileFormat, annotation, item.FileName)
					namedFile = filepath.Join(namedFileDir, item.FileName)
				}

				log.Emit(stdout.Info("OS:Operation for annotation").
					With("annotation", item.Annotation).
					With("file", namedFile).
					With("dir", namedFileDir))

				newFile, err := os.Create(namedFile)
				if err != nil {
					log.Emit(stdout.Error("IOError: Unable to create file").
						With("dir", namedFileDir).
						With("file", newFile).With("error", err))
					return err
				}

				if _, err := item.Writer.WriteTo(newFile); err != nil && err != io.EOF {
					newFile.Close()
					log.Emit(stdout.Error("IOError: Unable to write content to file").
						With("dir", namedFileDir).
						With("file", newFile).With("error", err))
					return err
				}

				log.Emit(stdout.Info("Annotation Resolved").With("annotation", item.Annotation).
					With("dir", namedFileDir).
					With("package", pkg.Package).With("file", pkg.File).With("generated-file", namedFile))

				newFile.Close()
			}
		}

	}

	return nil
}

//===========================================================================================================

// TypeAnnotationGenerator defines a function which generates specific code related to the giving
// Annotation for a non-struct, non-interface type declaration. This allows you to apply and create
// new sources specifically for a giving type(non-struct, non-interface).
// It is responsible to fully contain all operations required to both generator any source and write such to
type TypeAnnotationGenerator func(AnnotationDeclaration, TypeDeclaration, PackageDeclaration) ([]gen.WriteDirective, error)

// StructAnnotationGenerator defines a function which generates specific code related to the giving
// Annotation. This allows you to generate a new source file containg source code for a giving struct type.
// It is responsible to fully contain all operations required to both generator any source and write such to.
type StructAnnotationGenerator func(AnnotationDeclaration, StructDeclaration, PackageDeclaration) ([]gen.WriteDirective, error)

// InterfaceAnnotationGenerator defines a function which generates specific code related to the giving
// Annotation. This allows you to generate a new source file containg source code for a giving interface type.
// It is responsible to fully contain all operations required to both generator any source and write such to
// appropriate files as intended, meta-data about package, and file paths are already include in the PackageDeclaration.
type InterfaceAnnotationGenerator func(AnnotationDeclaration, InterfaceDeclaration, PackageDeclaration) ([]gen.WriteDirective, error)

// PackageAnnotationGenerator defines a function which generates specific code related to the giving
// Annotation for a package. This allows you to apply and create new sources specifically because of a
// package wide annotation.
// It is responsible to fully contain all operations required to both generator any source and write such to
// All generators are expected to return
type PackageAnnotationGenerator func(AnnotationDeclaration, PackageDeclaration) ([]gen.WriteDirective, error)

//===========================================================================================================

// Annotations defines a struct which contains a map of all annotation code generator.
type Annotations struct {
	Types      map[string]TypeAnnotationGenerator
	Structs    map[string]StructAnnotationGenerator
	Packages   map[string]PackageAnnotationGenerator
	Interfaces map[string]InterfaceAnnotationGenerator
}

// AnnotationRegistry defines a structure which contains giving list of possible
// annotation generators for both package level and type level declaration.
type AnnotationRegistry struct {
	ml                   sync.RWMutex
	typeAnnotations      map[string]TypeAnnotationGenerator
	structAnnotations    map[string]StructAnnotationGenerator
	pkgAnnotations       map[string]PackageAnnotationGenerator
	interfaceAnnotations map[string]InterfaceAnnotationGenerator
}

// NewAnnotationRegistry returns a new instance of a AnnotationRegistry.
func NewAnnotationRegistry() *AnnotationRegistry {
	return &AnnotationRegistry{
		typeAnnotations:      make(map[string]TypeAnnotationGenerator),
		structAnnotations:    make(map[string]StructAnnotationGenerator),
		pkgAnnotations:       make(map[string]PackageAnnotationGenerator),
		interfaceAnnotations: make(map[string]InterfaceAnnotationGenerator),
	}
}

// Clone returns a type which contains all copies of the generators provided by
// the AnnotationRegistry.
func (a *AnnotationRegistry) Clone() Annotations {
	a.ml.RLock()
	defer a.ml.RUnlock()

	var cloned Annotations

	for name, item := range a.pkgAnnotations {
		cloned.Packages[name] = item
	}

	for name, item := range a.structAnnotations {
		cloned.Structs[name] = item
	}

	for name, item := range a.typeAnnotations {
		cloned.Types[name] = item
	}

	for name, item := range a.interfaceAnnotations {
		cloned.Interfaces[name] = item
	}

	return cloned
}

// CopyStrategy defines a int type used to represent a copy strategy for
// cloning a AnnotationStrategy.
type CopyStrategy int

// Contains different copy strategy.
const (
	OursOverTheirs CopyStrategy = iota + 1
	TheirsOverOurs
)

// Copy copies over all available type generators from the provided AnnotationRegistry with
// the CopyStrategy.
func (a *AnnotationRegistry) Copy(registry *AnnotationRegistry, strategy CopyStrategy) {
	cloned := registry.Clone()

	a.ml.Lock()
	defer a.ml.Unlock()

	for name, item := range cloned.Packages {
		_, ok := a.pkgAnnotations[name]

		if !ok || (ok && strategy == TheirsOverOurs) {
			a.pkgAnnotations[name] = item
		}
	}

	for name, item := range cloned.Types {
		_, ok := a.typeAnnotations[name]
		if !ok || (ok && strategy == TheirsOverOurs) {
			a.typeAnnotations[name] = item
		}
	}

	for name, item := range cloned.Structs {
		_, ok := a.structAnnotations[name]
		if !ok || (ok && strategy == TheirsOverOurs) {
			a.structAnnotations[name] = item
		}
	}

	for name, item := range cloned.Interfaces {
		_, ok := a.interfaceAnnotations[name]
		if !ok || (ok && strategy == TheirsOverOurs) {
			a.interfaceAnnotations[name] = item
		}
	}
}

// MustPackage returns the annotation generator associated with the giving annotation name.
func (a *AnnotationRegistry) MustPackage(annotation string) PackageAnnotationGenerator {
	annon, err := a.GetPackage(annotation)
	if err == nil {
		return annon
	}

	panic(err)
}

// AnnotationWriteDirective defines a type which provides a WriteDiretive and the associated
// name.
type AnnotationWriteDirective struct {
	gen.WriteDirective
	Annotation string
}

// ParseDeclr runs the generators suited for each declaration and type returning a slice of
// Annotationgen.WriteDirective that delivers the content to be created for each piece.
func (a *AnnotationRegistry) ParseDeclr(declr PackageDeclaration) ([]AnnotationWriteDirective, error) {
	var directives []AnnotationWriteDirective

	// Generate directives for package level
	for _, annotation := range declr.Annotations {
		generator, err := a.GetPackage(annotation.Name)
		if err != nil {
			continue
		}

		drs, err := generator(annotation, declr)
		if err != nil {
			return nil, err
		}

		for _, directive := range drs {
			directives = append(directives, AnnotationWriteDirective{
				WriteDirective: directive,
				Annotation:     annotation.Name,
			})
		}
	}

	for _, inter := range declr.Interfaces {
		for _, annotation := range inter.Annotations {
			generator, err := a.GetInterfaceType(annotation.Name)
			if err != nil {
				continue
			}

			drs, err := generator(annotation, inter, declr)
			if err != nil {
				return nil, err
			}

			for _, directive := range drs {
				directives = append(directives, AnnotationWriteDirective{
					WriteDirective: directive,
					Annotation:     annotation.Name,
				})
			}
		}
	}

	for _, structs := range declr.Structs {
		for _, annotation := range structs.Annotations {
			generator, err := a.GetStructType(annotation.Name)
			if err != nil {
				continue
			}

			drs, err := generator(annotation, structs, declr)
			if err != nil {
				return nil, err
			}

			for _, directive := range drs {
				directives = append(directives, AnnotationWriteDirective{
					WriteDirective: directive,
					Annotation:     annotation.Name,
				})
			}
		}
	}

	for _, typ := range declr.Types {
		for _, annotation := range typ.Annotations {
			generator, err := a.GetType(annotation.Name)
			if err != nil {
				continue
			}

			drs, err := generator(annotation, typ, declr)
			if err != nil {
				return nil, err
			}

			for _, directive := range drs {
				directives = append(directives, AnnotationWriteDirective{
					WriteDirective: directive,
					Annotation:     annotation.Name,
				})
			}
		}
	}

	return directives, nil
}

// GetPackage returns the annotation generator associated with the giving annotation name.
func (a *AnnotationRegistry) GetPackage(annotation string) (PackageAnnotationGenerator, error) {
	annotation = strings.TrimPrefix(annotation, "@")

	var annon PackageAnnotationGenerator
	var ok bool

	a.ml.RLock()
	{
		annon, ok = a.pkgAnnotations[annotation]
	}
	a.ml.RUnlock()

	if !ok {
		return nil, fmt.Errorf("Package Annotation @%s not found", annotation)
	}

	return annon, nil
}

// MustInterfaceType returns the annotation generator associated with the giving annotation name.
func (a *AnnotationRegistry) MustInterfaceType(annotation string) InterfaceAnnotationGenerator {
	annon, err := a.GetInterfaceType(annotation)
	if err == nil {
		return annon
	}

	panic(err)
}

// GetInterfaceType returns the annotation generator associated with the giving annotation name.
func (a *AnnotationRegistry) GetInterfaceType(annotation string) (InterfaceAnnotationGenerator, error) {
	annotation = strings.TrimPrefix(annotation, "@")
	var annon InterfaceAnnotationGenerator
	var ok bool

	a.ml.RLock()
	{
		annon, ok = a.interfaceAnnotations[annotation]
	}
	a.ml.RUnlock()

	if !ok {
		return nil, fmt.Errorf("Interface Annotation @%s not found", annotation)
	}

	return annon, nil
}

// MustStructType returns the annotation generator associated with the giving annotation name.
func (a *AnnotationRegistry) MustStructType(annotation string) StructAnnotationGenerator {
	annon, err := a.GetStructType(annotation)
	if err == nil {
		return annon
	}

	panic(err)
}

// GetStructType returns the annotation generator associated with the giving annotation name.
func (a *AnnotationRegistry) GetStructType(annotation string) (StructAnnotationGenerator, error) {
	annotation = strings.TrimPrefix(annotation, "@")
	var annon StructAnnotationGenerator
	var ok bool

	a.ml.RLock()
	{
		annon, ok = a.structAnnotations[annotation]
	}
	a.ml.RUnlock()

	if !ok {
		return nil, fmt.Errorf("Struct Annotation @%s not found", annotation)
	}

	return annon, nil
}

// MustType returns the annotation generator associated with the giving annotation name.
func (a *AnnotationRegistry) MustType(annotation string) TypeAnnotationGenerator {
	annon, err := a.GetType(annotation)
	if err == nil {
		return annon
	}

	panic(err)
}

// GetType returns the annotation generator associated with the giving annotation name.
func (a *AnnotationRegistry) GetType(annotation string) (TypeAnnotationGenerator, error) {
	annotation = strings.TrimPrefix(annotation, "@")

	var annon TypeAnnotationGenerator
	var ok bool

	a.ml.RLock()
	{
		annon, ok = a.typeAnnotations[annotation]
	}
	a.ml.RUnlock()

	if !ok {
		return nil, fmt.Errorf("Type Annotation @%s not found", annotation)
	}

	return annon, nil
}

// RegisterInterfaceType adds a interface type level annotation generator into the registry.
func (a *AnnotationRegistry) RegisterInterfaceType(annotation string, generator InterfaceAnnotationGenerator) {
	annotation = strings.TrimPrefix(annotation, "@")
	a.ml.Lock()
	{
		a.interfaceAnnotations[annotation] = generator
	}
	a.ml.Unlock()
}

// RegisterStructType adds a struct type level annotation generator into the registry.
func (a *AnnotationRegistry) RegisterStructType(annotation string, generator StructAnnotationGenerator) {
	annotation = strings.TrimPrefix(annotation, "@")
	a.ml.Lock()
	{
		a.structAnnotations[annotation] = generator
	}
	a.ml.Unlock()
}

// RegisterType adds a type(non-struct, non-interface) level annotation generator into the registry.
func (a *AnnotationRegistry) RegisterType(annotation string, generator TypeAnnotationGenerator) {
	annotation = strings.TrimPrefix(annotation, "@")
	a.ml.Lock()
	{
		a.typeAnnotations[annotation] = generator
	}
	a.ml.Unlock()
}

// RegisterPackage adds a package level annotation generator into the registry.
func (a *AnnotationRegistry) RegisterPackage(annotation string, generator PackageAnnotationGenerator) {
	annotation = strings.TrimPrefix(annotation, "@")
	a.ml.Lock()
	{
		a.pkgAnnotations[annotation] = generator
	}
	a.ml.Unlock()
}

//===========================================================================================================
