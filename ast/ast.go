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
	"strings"
	"sync"
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
	Name     string `json:"name"`
	Argument string `json:"argument"`
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

// ParseAnnotations parses the package which generates a series of ast with associated
// annotation for processing.
func ParseAnnotations(dir string) ([]PackageDeclaration, error) {
	tokenFiles, packages, err := PackageDir(dir, parser.ParseComments)
	if err != nil {
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

				if len(annons) > 1 {
					packageDeclr.Annotations = append(packageDeclr.Annotations, AnnotationDeclaration{
						Name:     annons[1],
						Argument: annons[2],
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

							if len(annons) > 1 {
								annotations = append(annotations, AnnotationDeclaration{
									Name:     annons[1],
									Argument: annons[2],
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

// WriteDirective defines a struct which contains giving directives as to the file and
// WriteDirective defines a struct which contains giving directives as to the file and
// the relative path within which it should be written to.
type WriteDirective struct {
	Writer     io.WriterTo // WriteTo which contains the complete content of the file to be written to.
	Annotation string      // annotation name, set by generator.
	Dir        string      // Relative dir path written into it if not existing.
}

// TypeAnnotationGenerator defines a function which generates specific code related to the giving
// Annotation for a non-struct, non-interface type declaration. This allows you to apply and create
// new sources specifically for a giving type(non-struct, non-interface).
// It is responsible to fully contain all operations required to both generator any source and write such to
type TypeAnnotationGenerator func(AnnotationDeclaration, TypeDeclaration, PackageDeclaration) ([]WriteDirective, error)

// StructAnnotationGenerator defines a function which generates specific code related to the giving
// Annotation. This allows you to generate a new source file containg source code for a giving struct type.
// It is responsible to fully contain all operations required to both generator any source and write such to.
type StructAnnotationGenerator func(AnnotationDeclaration, StructDeclaration, PackageDeclaration) ([]WriteDirective, error)

// InterfaceAnnotationGenerator defines a function which generates specific code related to the giving
// Annotation. This allows you to generate a new source file containg source code for a giving interface type.
// It is responsible to fully contain all operations required to both generator any source and write such to
// appropriate files as intended, meta-data about package, and file paths are already include in the PackageDeclaration.
type InterfaceAnnotationGenerator func(AnnotationDeclaration, InterfaceDeclaration, PackageDeclaration) ([]WriteDirective, error)

// PackageAnnotationGenerator defines a function which generates specific code related to the giving
// Annotation for a package. This allows you to apply and create new sources specifically because of a
// package wide annotation.
// It is responsible to fully contain all operations required to both generator any source and write such to
// All generators are expected to return
type PackageAnnotationGenerator func(AnnotationDeclaration, PackageDeclaration) ([]WriteDirective, error)

//===========================================================================================================

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

// MustPackage returns the annotation generator associated with the giving annotation name.
func (a *AnnotationRegistry) MustPackage(annotation string) PackageAnnotationGenerator {
	annon, err := a.GetPackage(annotation)
	if err == nil {
		return annon
	}

	panic(err)
}

// ParseDeclr runs the generators suited to deliver a complete write for the
// annotations associated with the
func (a *AnnotationRegistry) ParseDeclr(declr PackageDeclaration) ([]WriteDirective, error) {

	return nil, nil
}

// GetPackage returns the annotation generator associated with the giving annotation name.
func (a *AnnotationRegistry) GetPackage(annotation string) (PackageAnnotationGenerator, error) {
	var annon PackageAnnotationGenerator
	var ok bool

	a.ml.RLock()
	{
		annon, ok = a.pkgAnnotations[annotation]
	}
	a.ml.RUnlock()

	if !ok {
		return nil, fmt.Errorf("Annotation @%q not found", annotation)
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
	var annon InterfaceAnnotationGenerator
	var ok bool

	a.ml.RLock()
	{
		annon, ok = a.interfaceAnnotations[annotation]
	}
	a.ml.RUnlock()

	if !ok {
		return nil, fmt.Errorf("Annotation @%q not found", annotation)
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
	var annon StructAnnotationGenerator
	var ok bool

	a.ml.RLock()
	{
		annon, ok = a.structAnnotations[annotation]
	}
	a.ml.RUnlock()

	if !ok {
		return nil, fmt.Errorf("Annotation @%q not found", annotation)
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
	var annon TypeAnnotationGenerator
	var ok bool

	a.ml.RLock()
	{
		annon, ok = a.typeAnnotations[annotation]
	}
	a.ml.RUnlock()

	if !ok {
		return nil, fmt.Errorf("Annotation @%q not found", annotation)
	}

	return annon, nil
}

// RegisterInterfaceType adds a interface type level annotation generator into the registry.
func (a *AnnotationRegistry) RegisterInterfaceType(annotation string, generator InterfaceAnnotationGenerator) {
	a.ml.Lock()
	{
		a.interfaceAnnotations[annotation] = generator
	}
	a.ml.Unlock()
}

// RegisterStructType adds a struct type level annotation generator into the registry.
func (a *AnnotationRegistry) RegisterStructType(annotation string, generator StructAnnotationGenerator) {
	a.ml.Lock()
	{
		a.structAnnotations[annotation] = generator
	}
	a.ml.Unlock()
}

// RegisterType adds a type(non-struct, non-interface) level annotation generator into the registry.
func (a *AnnotationRegistry) RegisterType(annotation string, generator TypeAnnotationGenerator) {
	a.ml.Lock()
	{
		a.typeAnnotations[annotation] = generator
	}
	a.ml.Unlock()
}

// RegisterPackage adds a package level annotation generator into the registry.
func (a *AnnotationRegistry) RegisterPackage(annotation string, generator PackageAnnotationGenerator) {
	a.ml.Lock()
	{
		a.pkgAnnotations[annotation] = generator
	}
	a.ml.Unlock()
}

//===========================================================================================================
