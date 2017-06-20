package ast

import (
	"bytes"
	"errors"
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
	spaces     = regexp.MustCompile(`/s+`)
	itag       = regexp.MustCompile(`((\w+):"(\w+|[\w,?\s+\w]+)")`)
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
	Package     string                            `json:"package"`
	Path        string                            `json:"path"`
	FilePath    string                            `json:"filepath"`
	File        string                            `json:"file"`
	Annotations []AnnotationDeclaration           `json:"annotations"`
	Types       []TypeDeclaration                 `json:"types"`
	Structs     []StructDeclaration               `json:"structs"`
	Interfaces  []InterfaceDeclaration            `json:"interfaces"`
	Functions   []FuncDeclaration                 `json:"functions"`
	ObjectFunc  map[*ast.Object][]FuncDeclaration `json:"object_functions"`
}

// HasFunctionFor returns true/false if the giving Struct Declaration has the giving function name.
func (pkg PackageDeclaration) HasFunctionFor(str StructDeclaration, funcName string) bool {
	functions := Functions(pkg.FunctionsFor(str.Object.Name.Obj))

	if _, err := functions.Find(funcName); err != nil {
		return false
	}

	return true
}

// FunctionsForName returns a slice of FuncDeclaration for the giving name.
func (pkg PackageDeclaration) FunctionsForName(objName string) []FuncDeclaration {
	var funcs []FuncDeclaration

	for obj, list := range pkg.ObjectFunc {
		if obj.Name != objName {
			continue
		}

		funcs = append(funcs, list...)
	}

	return funcs
}

// FunctionsFor returns a slice of FuncDeclaration for the giving object.
func (pkg PackageDeclaration) FunctionsFor(obj *ast.Object) []FuncDeclaration {
	if funcs, ok := pkg.ObjectFunc[obj]; ok {
		return funcs
	}

	var funcs []FuncDeclaration

	for obj, list := range pkg.ObjectFunc {
		if obj.Name != obj.Name {
			continue
		}

		funcs = append(funcs, list...)
	}

	return funcs
}

// HasFunctionFor returns true/false if the giving function name exists for the package.
func HasFunctionFor(pkg PackageDeclaration) func(StructDeclaration, string) bool {
	return func(str StructDeclaration, funcName string) bool {
		return pkg.HasFunctionFor(str, funcName)
	}
}

//===========================================================================================================

// FuncDeclaration defines a type used to annotate a giving type declaration
// associated with a ast for a function.
type FuncDeclaration struct {
	LineNumber    int            `json:"line_number"`
	Column        int            `json:"column"`
	Package       string         `json:"package"`
	Path          string         `json:"path"`
	FilePath      string         `json:"filepath"`
	File          string         `json:"file"`
	FuncName      string         `json:"funcName"`
	RecieverName  string         `json:"receiverName"`
	Position      token.Pos      `json:"position"`
	FuncDeclr     *ast.FuncDecl  `json:"funcdeclr"`
	Type          *ast.FuncType  `json:"type"`
	Reciever      *ast.Object    `json:"receiver"`
	RecieverIdent *ast.Ident     `json:"receiverIdent"`
	FuncType      *ast.FieldList `json:"funcType"`
	Returns       *ast.FieldList `json:"returns"`
	Arguments     *ast.FieldList `json:"arguments"`
}

// Functions defines a slice of FuncDeclaration.
type Functions []FuncDeclaration

// Find returns the giving Function of the giving name.
func (fnList Functions) Find(name string) (FuncDeclaration, error) {
	for _, fn := range fnList {
		if fn.FuncName == name {
			return fn, nil
		}
	}

	return FuncDeclaration{}, fmt.Errorf("Function with %q not found", name)
}

//===========================================================================================================

// AnnotationAssociationDeclaration defines a type which defines an association between
// a giving annotation and a series of values.
type AnnotationAssociationDeclaration struct {
	Annotation string `json:"annotation"`
	Action     string `json:"action"`
	TypeName   string `json:"typeName"`
}

// InterfaceDeclaration defines a type which holds annotation data for a giving interface type declaration.
type InterfaceDeclaration struct {
	LineNumber   int                                         `json:"line_number"`
	Column       int                                         `json:"column"`
	Package      string                                      `json:"package"`
	Path         string                                      `json:"path"`
	FilePath     string                                      `json:"filepath"`
	File         string                                      `json:"file"`
	Interface    *ast.InterfaceType                          `json:"interface"`
	Object       *ast.TypeSpec                               `json:"object"`
	Position     token.Pos                                   `json:"position"`
	Annotations  []AnnotationDeclaration                     `json:"annotations"`
	Associations map[string]AnnotationAssociationDeclaration `json:"associations"`
}

// StructDeclaration defines a type which holds annotation data for a giving struct type declaration.
type StructDeclaration struct {
	LineNumber   int                                         `json:"line_number"`
	Column       int                                         `json:"column"`
	Package      string                                      `json:"package"`
	Path         string                                      `json:"path"`
	FilePath     string                                      `json:"filepath"`
	File         string                                      `json:"file"`
	Struct       *ast.StructType                             `json:"struct"`
	Object       *ast.TypeSpec                               `json:"object"`
	Position     token.Pos                                   `json:"position"`
	Annotations  []AnnotationDeclaration                     `json:"annotations"`
	Associations map[string]AnnotationAssociationDeclaration `json:"associations"`
}

// TypeDeclaration defines a type which holds annotation data for a giving type declaration.
type TypeDeclaration struct {
	LineNumber   int                                         `json:"line_number"`
	Column       int                                         `json:"column"`
	Package      string                                      `json:"package"`
	Path         string                                      `json:"path"`
	FilePath     string                                      `json:"filepath"`
	File         string                                      `json:"file"`
	Object       *ast.TypeSpec                               `json:"object"`
	Position     token.Pos                                   `json:"position"`
	Annotations  []AnnotationDeclaration                     `json:"annotations"`
	Associations map[string]AnnotationAssociationDeclaration `json:"associations"`
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
			packageDeclr.ObjectFunc = make(map[*ast.Object][]FuncDeclaration, 0)

			if relPath, err := filepath.Rel(GoSrcPath, path); err == nil {
				packageDeclr.Path = filepath.Dir(relPath)
				packageDeclr.File = filepath.Base(relPath)
			}

			for _, comment := range file.Doc.List {
				text := strings.TrimPrefix(comment.Text, "//")

				if !annotation.MatchString(text) {
					continue
				}

				annons := annotation.FindStringSubmatch(text)

				log.Emit(stdout.Info("Annotation in Package comments").
					With("dir", dir).
					With("annotation", annons[1:]).
					With("comment", comment.Text))

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
		declrLoop:
			for _, declr := range file.Decls {

				switch rdeclr := declr.(type) {
				case *ast.FuncDecl:

					tokenPosition := tokenFiles.Position(rdeclr.Pos())

					var defFunc FuncDeclaration

					defFunc.FuncDeclr = rdeclr
					defFunc.Type = rdeclr.Type
					defFunc.Position = rdeclr.Pos()
					defFunc.Path = packageDeclr.Path
					defFunc.File = packageDeclr.File
					defFunc.FuncName = rdeclr.Name.Name
					defFunc.Column = tokenPosition.Column
					defFunc.Package = packageDeclr.Package
					defFunc.LineNumber = tokenPosition.Line
					defFunc.FilePath = packageDeclr.FilePath

					if rdeclr.Type != nil {
						defFunc.Returns = rdeclr.Type.Results
						defFunc.Arguments = rdeclr.Type.Params
					}

					if rdeclr.Recv != nil {
						defFunc.FuncType = rdeclr.Recv

						nameIdent := rdeclr.Recv.List[0]

						if receiverNameType, ok := nameIdent.Type.(*ast.Ident); ok {
							defFunc.RecieverName = receiverNameType.Name
							defFunc.Reciever = receiverNameType.Obj
							defFunc.RecieverIdent = receiverNameType

							if rems, ok := packageDeclr.ObjectFunc[receiverNameType.Obj]; ok {
								rems = append(rems, defFunc)
								packageDeclr.ObjectFunc[receiverNameType.Obj] = rems
							} else {
								packageDeclr.ObjectFunc[receiverNameType.Obj] = []FuncDeclaration{defFunc}
							}

							continue declrLoop
						}
					}

					packageDeclr.Functions = append(packageDeclr.Functions, defFunc)
					continue declrLoop

				case *ast.GenDecl:

					var annotations []AnnotationDeclaration

					associations := make(map[string]AnnotationAssociationDeclaration, 0)

					if rdeclr.Doc != nil {
						for _, comment := range rdeclr.Doc.List {
							text := strings.TrimPrefix(comment.Text, "//")

							if !annotation.MatchString(text) {
								continue
							}

							annons := annotation.FindStringSubmatch(text)

							log.Emit(stdout.Info("Annotation in Decleration comment").
								With("dir", dir).
								With("comment", comment.Text).
								With("annotation", annons[1:]).
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

								switch annons[1] {
								case "associates":
									if len(arguments) < 3 {
										log.Emit(stdout.Error("Association Annotation in Decleration is incomplete: Expects 3 elements").
											With("dir", dir).
											With("association", arguments).
											With("comment", comment.Text).
											With("annotation", annons[1:]).
											With("position", rdeclr.Pos()).
											With("token", rdeclr.Tok.String()))

										continue declrLoop
									}

									log.Emit(stdout.Info("Association for Annotation in Decleration").
										With("dir", dir).
										With("association-annotation", strings.TrimPrefix(arguments[0], "@")).
										With("association-action", arguments[1]).
										With("association-typeName", arguments[2]).
										With("comment", comment.Text).
										With("annotation", annons[1:]).
										With("position", rdeclr.Pos()).
										With("token", rdeclr.Tok.String()))

									associations[arguments[1]] = AnnotationAssociationDeclaration{
										Action:     arguments[1],
										TypeName:   arguments[2],
										Annotation: strings.TrimPrefix(arguments[0], "@"),
									}

									break

								default:
									annotations = append(annotations, AnnotationDeclaration{
										Name:      annons[1],
										Arguments: arguments,
									})
								}

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
									With("Annotations", len(annotations)).
									With("StructName", obj.Name))

								packageDeclr.Structs = append(packageDeclr.Structs, StructDeclaration{
									Object:       obj,
									Struct:       robj,
									Annotations:  annotations,
									Associations: associations,
									File:         packageDeclr.File,
									Package:      packageDeclr.Package,
									Path:         packageDeclr.Path,
									FilePath:     packageDeclr.FilePath,
									LineNumber:   tokenPosition.Line,
									Column:       tokenPosition.Column,
								})
								break

							case *ast.InterfaceType:
								log.Emit(stdout.Info("Annotation in Decleration").
									With("Type", "Interface").
									With("Annotations", len(annotations)).
									With("StructName", obj.Name))

								packageDeclr.Interfaces = append(packageDeclr.Interfaces, InterfaceDeclaration{
									Object:       obj,
									Interface:    robj,
									Annotations:  annotations,
									Associations: associations,
									File:         packageDeclr.File,
									Package:      packageDeclr.Package,
									Path:         packageDeclr.Path,
									FilePath:     packageDeclr.FilePath,
									LineNumber:   tokenPosition.Line,
									Column:       tokenPosition.Column,
								})
								break

							default:
								log.Emit(stdout.Info("Annotation in Decleration").
									With("Type", "OtherType").
									With("Marker", "NonStruct/NonInterface:Type").
									With("Annotations", len(annotations)).
									With("StructName", obj.Name))

								packageDeclr.Types = append(packageDeclr.Types, TypeDeclaration{
									Object:       obj,
									Annotations:  annotations,
									Associations: associations,
									File:         packageDeclr.File,
									Package:      packageDeclr.Package,
									Path:         packageDeclr.Path,
									FilePath:     packageDeclr.FilePath,
									LineNumber:   tokenPosition.Line,
									Column:       tokenPosition.Column,
								})
							}

						case *ast.ImportSpec:
							// Do Nothing.
						}
					}

				case *ast.BadDecl:
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
					With("error", err.Error()).With("package", pkg.Package))
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
						With("dir", namedFileDir).With("error", err.Error()))
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

				fileStat, err := os.Stat(namedFile)
				if err == nil && !fileStat.IsDir() && item.DontOverride {
					log.Emit(stdout.Info("Annotation Unresolved: File already exists and must no over-write").With("annotation", item.Annotation).
						With("dir", namedFileDir).
						With("package", pkg.Package).
						With("file", pkg.File).
						With("generated-file", namedFile))

					continue
				}

				newFile, err := os.Create(namedFile)
				if err != nil {
					log.Emit(stdout.Error("IOError: Unable to create file").
						With("dir", namedFileDir).
						With("file", namedFile).With("error", err.Error()))
					continue
				}

				if _, err := item.Writer.WriteTo(newFile); err != nil && err != io.EOF {
					newFile.Close()
					log.Emit(stdout.Error("IOError: Unable to write content to file").
						With("dir", namedFileDir).
						With("file", namedFile).With("error", err.Error()))
					continue
				}

				log.Emit(stdout.Info("Annotation Resolved").With("annotation", item.Annotation).
					With("dir", namedFileDir).
					With("package", pkg.Package).
					With("file", pkg.File).
					With("generated-file", namedFile))

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

// FindStructType defines a function to search a package declaration Structs of a giving typeName.
func FindStructType(pkg PackageDeclaration, typeName string) (StructDeclaration, error) {
	for _, elem := range pkg.Structs {
		if elem.Object.Name.Name == typeName {
			return elem, nil
		}
	}

	return StructDeclaration{}, fmt.Errorf("Struct of type %q not found", typeName)
}

// FindInterfaceType defines a function to search a package declaration Interface of a giving typeName.
func FindInterfaceType(pkg PackageDeclaration, typeName string) (InterfaceDeclaration, error) {
	for _, elem := range pkg.Interfaces {
		if elem.Object.Name.Name == typeName {
			return elem, nil
		}
	}

	return InterfaceDeclaration{}, fmt.Errorf("Interface of type %q not found", typeName)
}

// FindType defines a function to search a package declaration Structs of a giving typeName.
func FindType(pkg PackageDeclaration, typeName string) (TypeDeclaration, error) {
	for _, elem := range pkg.Types {
		if elem.Object.Name.Name == typeName {
			return elem, nil
		}
	}

	return TypeDeclaration{}, fmt.Errorf("Non(Struct|Interface) of type %q not found", typeName)
}

// GetStructSpec attempts to retrieve the TypeSpec and StructType if the value
// matches this.
func GetStructSpec(val interface{}) (*ast.TypeSpec, *ast.StructType, error) {
	rval, ok := val.(*ast.TypeSpec)
	if !ok {
		return nil, nil, errors.New("Not ast.TypeSpec type")
	}

	rstruct, ok := rval.Type.(*ast.StructType)
	if !ok {
		return nil, nil, errors.New("Not ast.StructType type for *ast.TypeSpec.Type")
	}

	return rval, rstruct, nil
}

//===========================================================================================================

// MapOutFields defines a function to return a map of field name and value
// pair for the giving struct.
func MapOutFields(item StructDeclaration, rootName, tagName, fallback string) (string, error) {
	vals, err := MapOutFieldsToMap(item, rootName, tagName, fallback)
	if err != nil {
		return "", err
	}

	var bu bytes.Buffer

	if _, err := gen.Map("string", "interface{}", vals).WriteTo(&bu); err != nil {
		return "", err
	}

	return bu.String(), nil
}

// MapOutFieldsToMap defines a function to return a map of field name and value
// pair for the giving struct.
func MapOutFieldsToMap(item StructDeclaration, rootName, tagName, fallback string) (map[string]io.WriterTo, error) {
	fields := Fields(GetFields(item))

	wTags := fields.TagFor(tagName)
	if len(wTags) == 0 {
		wTags = fields.TagFor(fallback)

		if len(wTags) == 0 {
			return nil, fmt.Errorf("No tags match for %q and %q fallback for struct %q", tagName, fallback, item.Object.Name)
		}
	}

	dm := make(map[string]io.WriterTo)

	embedded := fields.Embedded()

	for _, embed := range embedded {
		emt, ems, err := GetStructSpec(embed.Type.Decl)
		if err != nil {
			return nil, err
		}

		vals, err := MapOutFieldsToMap(StructDeclaration{
			Object: emt,
			Struct: ems,
		}, fmt.Sprintf("%s.%s", rootName, embed.FieldName), tagName, fallback)

		if err != nil {
			return nil, err
		}

		for name, val := range vals {
			dm[name] = val
		}
	}

	// Collect key field names from embedded first
	for _, tag := range wTags {
		if tag.Value == "-" {
			continue
		}

		if tag.Field.Type != nil {
			embededType, embedStruct, err := GetStructSpec(tag.Field.Type.Decl)
			if err != nil {
				return nil, err
			}

			flds, err := MapOutFieldsToMap(StructDeclaration{
				Object: embededType,
				Struct: embedStruct,
			}, fmt.Sprintf("%s.%s", rootName, tag.Field.FieldName), tagName, fallback)

			if err != nil {
				return nil, err
			}

			dm[tag.Value] = gen.Map("string", "interface{}", flds)
			continue
		}

		dm[tag.Value] = gen.Fmt("%s.%s", rootName, tag.Field.FieldName)
	}

	return dm, nil
}

// MapOutValues defines a function to return a map of field name and associated
// placeholders as value.
func MapOutValues(item StructDeclaration, onlyExported bool) (string, error) {
	var bu bytes.Buffer

	if _, err := MapOutFieldsValues(item, onlyExported, nil).WriteTo(&bu); err != nil {
		return "", err
	}

	return bu.String(), nil
}

// MapOutFieldsValues defines a function to return a map of field name and associated
// placeholders as value.
func MapOutFieldsValues(item StructDeclaration, onlyExported bool, name *gen.NameDeclr) io.WriterTo {
	fields := Fields(GetFields(item))

	var writers []io.WriterTo

	if name == nil {
		tmpName := gen.FmtName("%sVar", strings.ToLower(item.Object.Name.Name))

		name = &tmpName

		vardecl := gen.VarType(
			tmpName,
			gen.Type(item.Object.Name.Name),
		)

		writers = append(writers, vardecl, gen.Text("\n"))
	}

	normals := fields.Normal()
	embedded := fields.Embedded()

	handleOtherField := func(embed FieldDeclaration) {
		elemValue := gen.AssignValue(
			gen.FmtName("%s.%s", name.Name, embed.FieldName),
			gen.Text(DefaultTypeValueString(embed.FieldTypeName)),
		)

		writers = append(writers, elemValue, gen.Text("\n"))
	}

	handleStructField := func(embed FieldDeclaration) {
		embedName := gen.FmtName("%sVar", strings.ToLower(embed.FieldName))

		elemDeclr := gen.VarType(
			embedName,
			gen.Type(embed.FieldTypeName),
		)

		writers = append(writers, elemDeclr)

		if item.Struct != nil {
			body := MapOutFieldsValues(StructDeclaration{
				Object: embed.Spec,
				Struct: embed.Struct,
			}, onlyExported, &embedName)

			writers = append(writers, body)
		}

		elemValue := gen.AssignValue(
			gen.FmtName("%s.%s", name.Name, embed.FieldName),
			embedName,
		)

		writers = append(writers, elemValue, gen.Text("\n"))
	}

	for _, embed := range embedded {
		if !embed.Exported && onlyExported {
			continue
		}

		if embed.IsStruct {
			handleStructField(embed)
			continue
		}

		handleOtherField(embed)
	}

	for _, normal := range normals {
		if !normal.Exported && onlyExported {
			continue
		}

		if normal.IsStruct {
			handleStructField(normal)
			continue
		}

		handleOtherField(normal)
	}

	return gen.Block(writers...)
}

//===========================================================================================================

// DefaultTypeValueString returns the default value string of a giving
// typeName.
func DefaultTypeValueString(typeName string) string {
	switch typeName {
	case "uint", "uint32", "uint64":
		return "0"
	case "int", "int32", "int64":
		return "0"
	case "string":
		return `""`
	case "rune":
		return `rune(0)`
	case "float32", "float64":
		return "0.0"
	default:
		return "nil"
	}
}

//===========================================================================================================

// Fields defines a slice type of FieldDeclaration.
type Fields []FieldDeclaration

// Normal defines a function that returns all fields which are non-embedded.
func (flds Fields) Normal() Fields {
	var fields Fields

	for _, declr := range flds {
		if declr.Embedded {
			continue
		}

		fields = append(fields, declr)
	}

	return fields
}

// Embedded defines a function that returns all appropriate Field
// that match the giving tagName
func (flds Fields) Embedded() Fields {
	var fields Fields

	for _, declr := range flds {
		if declr.Embedded {
			fields = append(fields, declr)
		}
	}

	return fields
}

// TagFor defines a function that returns all appropriate TagDeclaration
// that match the giving tagName
func (flds Fields) TagFor(tagName string) []TagDeclaration {
	var declrs []TagDeclaration

	for _, declr := range flds {
		if dl, err := declr.GetTag(tagName); err == nil {
			declrs = append(declrs, dl)
		}
	}

	return declrs
}

// FieldDeclaration defines a type to represent a giving struct fields and tags.
type FieldDeclaration struct {
	Exported      bool             `json:"exported"`
	Embedded      bool             `json:"embedded"`
	IsStruct      bool             `json:"is_struct"`
	FieldName     string           `json:"field_name"`
	FieldTypeName string           `json:"field_type_name"`
	Field         *ast.Field       `json:"field"`
	Type          *ast.Object      `json:"type"`
	Spec          *ast.TypeSpec    `json:"spec"`
	Struct        *ast.StructType  `json:"struct"`
	Tags          []TagDeclaration `json:"tags"`
}

// GetFields returns all fields associated with the giving struct but skips
func GetFields(str StructDeclaration) []FieldDeclaration {
	var fields []FieldDeclaration

	for _, item := range str.Struct.Fields.List {
		typeIdent, ok := item.Type.(*ast.Ident)
		if !ok {
			continue
		}

		var field FieldDeclaration

		field.Field = item
		field.Type = typeIdent.Obj
		field.FieldName = typeIdent.Name
		field.FieldTypeName = typeIdent.Name

		if typeIdent.Obj != nil {
			if spec, ok := typeIdent.Obj.Decl.(*ast.TypeSpec); ok {
				field.Spec = spec

				if strt, ok := spec.Type.(*ast.StructType); ok {
					field.Struct = strt
					field.IsStruct = true
				}
			}
		}

		if len(item.Names) == 0 {
			field.Exported = true
			field.Embedded = true

			fields = append(fields, field)
			continue
		}

		fieldName := item.Names[0]
		field.FieldName = fieldName.Name

		// fmt.Printf("Exported: %t -> %q : %q\n", fieldName.Name == strings.ToLower(fieldName.Name), fieldName.Name, strings.ToLower(fieldName.Name))

		if typeIdent.Name != strings.ToLower(fieldName.Name) {
			field.Exported = true
		}

		if item.Tag == nil {
			fields = append(fields, field)
			continue
		}

		tags := strings.Split(spaces.ReplaceAllString(item.Tag.Value, " "), " ")

		for _, tag := range tags {
			if !itag.MatchString(tag) {
				continue
			}

			res := itag.FindStringSubmatch(tag)
			resValue := strings.Split(res[3], ",")

			field.Tags = append(field.Tags, TagDeclaration{
				Field: field,
				Base:  res[0],
				Name:  res[2],
				Value: resValue[0],
				Metas: resValue[1:],
			})
		}

		fields = append(fields, field)
	}

	return fields
}

// GetTag returns the giving tag associated with the name if it exists.
func (f FieldDeclaration) GetTag(tagName string) (TagDeclaration, error) {
	for _, tag := range f.Tags {
		if tag.Name == tagName {
			return tag, nil
		}
	}

	return TagDeclaration{}, fmt.Errorf("Tag for %q not found", tagName)
}

// TagDeclaration defines a type which represents a giving tag declaration for a provided type.
type TagDeclaration struct {
	Name  string           `json:"name"`
	Value string           `json:"value"`
	Metas []string         `json:"metas"`
	Base  string           `json:"base"`
	Field FieldDeclaration `json:"field"`
}

// Has returns true/false if the tag.Metas has the given value in the list.
func (t TagDeclaration) Has(item string) bool {
	for _, meta := range t.Metas {
		if meta == item {
			return true
		}
	}

	return false
}
