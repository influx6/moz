package ast

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"runtime"

	"github.com/icrowley/fake"
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
	GoPath    = os.Getenv("GOPATH")
	GoSrcPath = filepath.Join(GoPath, "src")

	spaces     = regexp.MustCompile(`/s+`)
	itag       = regexp.MustCompile(`((\w+):"(\w+|[\w,?\s+\w]+)")`)
	annotation = regexp.MustCompile("@(\\w+(:\\w+)?)(\\([.\\s\\S]+\\))?")

	ASTTemplatFuncs = map[string]interface{}{
		"getTag":            GetTag,
		"fieldFor":          FieldFor,
		"getFields":         GetFields,
		"fieldNameFor":      FieldNameFor,
		"mapoutFields":      MapOutFields,
		"mapoutValues":      MapOutValues,
		"fieldByName":       FieldByFieldName,
		"randomValue":       RandomFieldAssign,
		"fieldsJSON":        MapOutFieldsToJSON,
		"stringValueFor":    ToValueString,
		"defaultValue":      AssignDefaultValue,
		"randomFieldValue":  RandomFieldValue,
		"defaultType":       DefaultTypeValueString,
		"defaultFieldValue": DefaultFieldValue,
	}
)

var (
	// ErrEmptyList defines a error returned for a empty array or slice.
	ErrEmptyList = errors.New("Slice/List is empty")
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

//===========================================================================================================

// ParseFileAnnotations parses the package from the provided file.
func ParseFileAnnotations(log metrics.Metrics, path string) (PackageDeclaration, error) {
	dir := filepath.Dir(path)

	tokenFiles, file, err := PackageFile(path, parser.ParseComments)
	if err != nil {
		log.Emit(stdout.Error(err).With("message", "Failed to parse file").With("dir", dir).With("file", file))
		return PackageDeclaration{}, err
	}

	return parseFileToPackage(log, dir, path, filepath.Base(dir), tokenFiles, file)
}

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
			res, err := parseFileToPackage(log, dir, path, pkg.Name, tokenFiles, file)
			if err != nil {
				return nil, err
			}

			packageDeclrs = append(packageDeclrs, res)
		}
	}

	return packageDeclrs, nil
}

func parseFileToPackage(log metrics.Metrics, dir string, path string, pkgName string, tokenFiles *token.FileSet, file *ast.File) (PackageDeclaration, error) {
	var packageDeclr PackageDeclaration

	{
		packageDeclr.Package = pkgName
		packageDeclr.FilePath = path
		packageDeclr.Imports = make(map[string]ImportDeclaration, 0)
		packageDeclr.ObjectFunc = make(map[*ast.Object][]FuncDeclaration, 0)

		for _, imp := range file.Imports {
			var pkgName string

			if imp.Name != nil {
				pkgName = imp.Name.Name
			}

			impPkgPath, err := strconv.Unquote(imp.Path.Value)
			if err != nil {
				impPkgPath = imp.Path.Value
			}

			packageDeclr.Imports[pkgName] = ImportDeclaration{
				Name: pkgName,
				Path: impPkgPath,
			}
		}

		if relPath, err := filepath.Rel(GoSrcPath, path); err == nil {
			packageDeclr.Path = filepath.Dir(relPath)
			packageDeclr.File = filepath.Base(relPath)
		}

		if runtime.GOOS == "windows" {
			packageDeclr.Path = filepath.ToSlash(packageDeclr.Path)
			packageDeclr.File = filepath.ToSlash(packageDeclr.File)
			packageDeclr.FilePath = filepath.ToSlash(packageDeclr.FilePath)
		}

		if file.Doc != nil {
			annotationRead := ReadAnnotationsFromCommentry(bytes.NewBufferString(file.Doc.Text()))

			log.Emit(stdout.Info("Annotations in Package comments").
				With("dir", dir).
				With("annotations", annotationRead).
				With("comment", file.Doc.Text()))

			packageDeclr.Annotations = append(packageDeclr.Annotations, annotationRead...)
		}

		// Collect and categorize annotations in types and their fields.
	declrLoop:
		for _, declr := range file.Decls {
			switch rdeclr := declr.(type) {
			case *ast.FuncDecl:

				tokenPosition := tokenFiles.Position(rdeclr.Pos())

				var defFunc FuncDeclaration

				defFunc.PackageDeclr = &packageDeclr
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
					annotationRead := ReadAnnotationsFromCommentry(bytes.NewBufferString(rdeclr.Doc.Text()))

					for _, item := range annotationRead {
						log.Emit(stdout.Info("Annotation in Decleration comment").
							With("dir", dir).
							With("comment", rdeclr.Doc.Text()).
							With("annotation", item.Name).
							With("position", rdeclr.Pos()).
							With("token", rdeclr.Tok.String()))

						switch item.Name {
						case "associates":
							log.Emit(stdout.Error("Association Annotation in Decleration is incomplete: Expects 3 elements").
								With("dir", dir).
								With("association", item.Arguments).
								With("position", rdeclr.Pos()).
								With("token", rdeclr.Tok.String()))

							if len(item.Arguments) >= 3 {
								associations[item.Arguments[0]] = AnnotationAssociationDeclaration{
									Record:     item,
									Template:   item.Template,
									Action:     item.Arguments[1],
									TypeName:   item.Arguments[2],
									Annotation: strings.TrimPrefix(item.Arguments[0], "@"),
								}
							}
						default:
							annotations = append(annotations, item)
						}
					}

				}

				for _, spec := range rdeclr.Specs {
					switch obj := spec.(type) {
					case *ast.ValueSpec:
						// Handles variable declaration
						// i.e Spec:
						// &ast.ValueSpec{Doc:(*ast.CommentGroup)(nil), Names:[]*ast.Ident{(*ast.Ident)(0xc4200e4a00)}, Type:ast.Expr(nil), Values:[]ast.Expr{(*ast.BasicLit)(0xc4200e4a20)}, Comment:(*ast.CommentGroup)(nil)}
						// &ast.ValueSpec{Doc:(*ast.CommentGroup)(nil), Names:[]*ast.Ident{(*ast.Ident)(0xc4200e4a40)}, Type:(*ast.Ident)(0xc4200e4a60), Values:[]ast.Expr(nil), Comment:(*ast.CommentGroup)(nil)}

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
								PackageDeclr: &packageDeclr,
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
								PackageDeclr: &packageDeclr,
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
								PackageDeclr: &packageDeclr,
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

	}

	return packageDeclr, nil
}

//===========================================================================================================

// Parse takes the provided package declrations parsing all internals with the
// appropriate generators suited to the type and annotations.
func Parse(toDir string, log metrics.Metrics, provider *AnnotationRegistry, packageDeclrs ...PackageDeclaration) error {
	{
	parseloop:
		for _, pkg := range packageDeclrs {
			wdrs, err := provider.ParseDeclr(pkg, toDir)
			if err != nil {
				log.Emit(stdout.Error("ParseFailure: Package %q", pkg.Package).
					With("error", err.Error()).With("package", pkg.Package))
				continue
			}

			log.Emit(stdout.Info("ParseSuccess: Package %q", pkg.Package).With("package", pkg.Package))

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
					namedFileDir = filepath.Join(newDir, toDir, item.Dir)
				} else {
					namedFileDir = filepath.Join(newDir, toDir)
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

// WhichPackage is an utility function which returns the appropriate package name to use
// if a toDir is provided as destination.
func WhichPackage(toDir string, pkg PackageDeclaration) string {
	if toDir != "" && toDir != "./" && toDir != "." {
		return strings.ToLower(filepath.Base(toDir))
	}

	return pkg.Package
}

//===========================================================================================================

// AnnotationDeclaration defines a annotation type which holds detail about a giving annotation.
type AnnotationDeclaration struct {
	Name      string            `json:"name"`
	Template  string            `json:"template"`
	Arguments []string          `json:"arguments"`
	Params    map[string]string `json:"params"`
}

// ImportDeclaration defines a type to contain import declaration within a package.
type ImportDeclaration struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// PackageDeclaration defines a type which holds details relating to annotations declared on a
// giving package.
type PackageDeclaration struct {
	Package     string                            `json:"package"`
	Path        string                            `json:"path"`
	FilePath    string                            `json:"filepath"`
	File        string                            `json:"file"`
	Imports     map[string]ImportDeclaration      `json:"imports"`
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

// AnnotationsFor returns all annotations with the giving name.
func (pkg PackageDeclaration) AnnotationsFor(typeName string) []AnnotationDeclaration {
	typeName = strings.TrimPrefix(typeName, "@")

	var found []AnnotationDeclaration

	for _, item := range pkg.Annotations {
		if strings.TrimPrefix(item.Name, "@") != typeName {
			continue
		}

		found = append(found, item)
	}

	return found
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

// ImportFor returns the ImportDeclaration associated with the giving handle.
// Returns error if the import is not found.
func (pkg PackageDeclaration) ImportFor(imp string) (ImportDeclaration, error) {
	impDeclr, ok := pkg.Imports[imp]
	if !ok {
		return ImportDeclaration{}, errors.New("Not found")
	}

	return impDeclr, nil
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
	LineNumber    int                 `json:"line_number"`
	Column        int                 `json:"column"`
	Package       string              `json:"package"`
	Path          string              `json:"path"`
	FilePath      string              `json:"filepath"`
	File          string              `json:"file"`
	FuncName      string              `json:"funcName"`
	RecieverName  string              `json:"receiverName"`
	Position      token.Pos           `json:"position"`
	FuncDeclr     *ast.FuncDecl       `json:"funcdeclr"`
	Type          *ast.FuncType       `json:"type"`
	Reciever      *ast.Object         `json:"receiver"`
	RecieverIdent *ast.Ident          `json:"receiverIdent"`
	FuncType      *ast.FieldList      `json:"funcType"`
	Returns       *ast.FieldList      `json:"returns"`
	Arguments     *ast.FieldList      `json:"arguments"`
	PackageDeclr  *PackageDeclaration `json:"-"`
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
	Annotation string                `json:"annotation"`
	Action     string                `json:"action"`
	Template   string                `json:"template"`
	TypeName   string                `json:"typeName"`
	Record     AnnotationDeclaration `json:"record"`
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
	PackageDeclr *PackageDeclaration                         `json:"-"`
}

// Methods returns the associated methods for the giving interface.
func (i InterfaceDeclaration) Methods() []FunctionDefinition {
	return GetInterfaceFunctions(i.Interface, i.PackageDeclr)
}

// ArgType defines a type to represent the information for a giving functions argument or
// return type declaration.
type ArgType struct {
	Name           string
	Type           string
	ExType         string
	Package        string
	BaseType       bool
	Import         ImportDeclaration
	Import2        ImportDeclaration
	NameObject     *ast.Object
	TypeObject     *ast.Object
	StructObject   *ast.StructType
	ImportedObject *ast.SelectorExpr
	ArrayType      *ast.ArrayType
	MapType        *ast.MapType
	ChanType       *ast.ChanType
	PointerType    *ast.StarExpr
	IdentType      *ast.Ident
}

// FunctionDefinition defines a type to represent the function/method declarations of an
// interface type.
type FunctionDefinition struct {
	Name      string
	Args      []ArgType
	Returns   []ArgType
	Func      *ast.FuncType
	Interface *ast.InterfaceType
}

// ArgumentNamesList returns the assignment names for the function arguments.
func (fd FunctionDefinition) ArgumentNamesList() string {
	var args []string

	for _, arg := range fd.Args {
		args = append(args, fmt.Sprintf("%s", arg.Name))
	}

	return strings.Join(args, ",")
}

// ReturnNamesList returns the assignment names for the return arguments
func (fd FunctionDefinition) ReturnNamesList() string {
	var rets []string

	for _, ret := range fd.Returns {
		rets = append(rets, fmt.Sprintf("%s", ret.Name))
	}

	return strings.Join(rets, ",")
}

// ReturnList returns a string version of the return of the giving function.
func (fd FunctionDefinition) ReturnList(asFromOutside bool) string {
	var rets []string

	for _, ret := range fd.Returns {
		if asFromOutside {
			rets = append(rets, fmt.Sprintf("%s", ret.ExType))
			continue
		}

		rets = append(rets, fmt.Sprintf("%s", ret.Type))
	}

	return strings.Join(rets, ",")
}

// ArgumentList returns a string version of the arguments of the giving function.
func (fd FunctionDefinition) ArgumentList(asFromOutside bool) string {
	var args []string

	for _, arg := range fd.Args {
		if asFromOutside {
			args = append(args, fmt.Sprintf("%s %s", arg.Name, arg.ExType))
			continue
		}

		args = append(args, fmt.Sprintf("%s %s", arg.Name, arg.Type))
	}

	return strings.Join(args, ",")
}

// GetIdentName returns the first indent found within the field if it exists.
func GetIdentName(field *ast.Field) (*ast.Ident, error) {
	if len(field.Names) == 0 {
		return nil, ErrEmptyList
	}

	return field.Names[0], nil
}

// GetInterfaceFunctions returns a slice of FunctionDefinitions retrieved from the provided
// interface type object.
func GetInterfaceFunctions(intr *ast.InterfaceType, pkg *PackageDeclaration) []FunctionDefinition {
	var defs []FunctionDefinition

	var retCounter int
	var varCounter int

	retCounter++
	varCounter++

	for _, method := range intr.Methods.List {
		if len(method.Names) > 0 {
			nameIdent := method.Names[0]
			ftype := method.Type.(*ast.FuncType)

			var arguments, returns []ArgType

			for _, result := range ftype.Results.List {
				resPkg, defaultresType := getPackageFromItem(result.Type, filepath.Base(pkg.Package))

				switch iobj := result.Type.(type) {
				case *ast.Ident:
					var name string
					var nameObj *ast.Object

					resName, err := GetIdentName(result)
					switch err != nil {
					case true:
						name = fmt.Sprintf("ret%d", retCounter)
						retCounter++
					case false:
						name = resName.Name
						nameObj = resName.Obj
					}

					returns = append(returns, ArgType{
						Name:       name,
						NameObject: nameObj,
						Type:       getName(iobj),
						ExType:     getNameAsFromOuter(iobj, filepath.Base(pkg.Package)),
						TypeObject: iobj.Obj,
						Package:    resPkg,
						BaseType:   defaultresType,
					})

				case *ast.SelectorExpr:
					// fmt.Printf("Result: %#v -> %#v -> %#q\n\n", iobj.X, iobj.Sel, iobj)
					xobj, ok := iobj.X.(*ast.Ident)
					if !ok {
						break
					}

					importDclr, _ := pkg.ImportFor(xobj.Name)

					name := fmt.Sprintf("ret%d", retCounter)
					retCounter++

					returns = append(returns, ArgType{
						Name:    name,
						Import:  importDclr,
						Package: xobj.Name,
						Type:    getName(iobj),
						ExType:  getNameAsFromOuter(iobj, filepath.Base(pkg.Package)),
					})

				case *ast.StarExpr:
					name := fmt.Sprintf("ret%d", retCounter)
					retCounter++

					var arg ArgType
					arg.Name = name
					arg.PointerType = iobj
					arg.Type = getName(iobj)
					arg.ExType = getNameAsFromOuter(iobj, filepath.Base(pkg.Package))

					switch value := iobj.X.(type) {
					case *ast.SelectorExpr:
						arg.ImportedObject = value

						vob, ok := value.X.(*ast.Ident)
						if !ok {
							break
						}

						importDclr, _ := pkg.ImportFor(vob.Name)

						arg.Package = vob.Name
						arg.Import = importDclr
					case *ast.StructType:
						arg.StructObject = value
						// arg.Package = resPkg
						// arg.BaseType = defaultresType
					case *ast.ArrayType:
						arg.ArrayType = value
					case *ast.Ident:
						arg.IdentType = value
						arg.NameObject = value.Obj

						arg.Package = resPkg
						arg.BaseType = defaultresType
					case *ast.ChanType:
						arg.ChanType = value
					}

					returns = append(returns, arg)

				case *ast.MapType:
					name := fmt.Sprintf("ret%d", retCounter)
					retCounter++

					var arg ArgType
					arg.Name = name
					arg.MapType = iobj
					arg.Type = getName(iobj)
					arg.ExType = getNameAsFromOuter(iobj, filepath.Base(pkg.Package))

					if keySel, err := getSelector(iobj.Key); err == nil {
						if x, ok := keySel.X.(*ast.Ident); ok {
							if imported, err := pkg.ImportFor(x.Name); err == nil {
								arg.Import = imported
							}
						}
					}

					if valSel, err := getSelector(iobj.Value); err == nil {
						if x, ok := valSel.X.(*ast.Ident); ok {
							if imported, err := pkg.ImportFor(x.Name); err == nil {
								arg.Import2 = imported
							}
						}
					}

					returns = append(returns, arg)
				case *ast.ArrayType:
					name := fmt.Sprintf("ret%d", retCounter)
					retCounter++

					var arg ArgType
					arg.Name = name
					arg.ArrayType = iobj
					arg.Type = getName(iobj)
					arg.ExType = getNameAsFromOuter(iobj, filepath.Base(pkg.Package))

					switch value := iobj.Elt.(type) {
					case *ast.SelectorExpr:
						arg.ImportedObject = value

						vob, ok := value.X.(*ast.Ident)
						if !ok {
							break
						}

						importDclr, _ := pkg.ImportFor(vob.Name)

						arg.Package = vob.Name
						arg.Import = importDclr
					case *ast.StarExpr:
						arg.PointerType = value
					case *ast.StructType:
						arg.StructObject = value
					case *ast.Ident:
						arg.IdentType = value
						arg.NameObject = value.Obj
						arg.Package = resPkg
						arg.BaseType = defaultresType
					case *ast.ChanType:
						arg.ChanType = value
					}

					returns = append(returns, arg)

				case *ast.ChanType:
					name := fmt.Sprintf("ret%d", retCounter)
					retCounter++

					var arg ArgType
					arg.Name = name
					arg.Type = getName(iobj.Value)
					arg.ExType = getNameAsFromOuter(iobj, filepath.Base(pkg.Package))

					switch value := iobj.Value.(type) {
					case *ast.SelectorExpr:
						arg.ImportedObject = value

						vob, ok := value.X.(*ast.Ident)
						if !ok {
							break
						}

						importDclr, _ := pkg.ImportFor(vob.Name)

						arg.Package = vob.Name
						arg.Import = importDclr
					case *ast.StarExpr:
						arg.PointerType = value
					case *ast.StructType:
						arg.StructObject = value
					case *ast.ArrayType:
						arg.ArrayType = value
					case *ast.Ident:
						arg.IdentType = value
						arg.NameObject = value.Obj

						arg.Package = resPkg
						arg.BaseType = defaultresType
					case *ast.ChanType:
						arg.ChanType = value
					}

					returns = append(returns, arg)
				default:
					// fmt.Printf("Result:Default: %#v -> %#q\n\n", iobj, iobj)
				}
			}

			for _, param := range ftype.Params.List {
				paramPkg, defaultresType := getPackageFromItem(param.Type, filepath.Base(pkg.Package))

				switch iobj := param.Type.(type) {
				case *ast.Ident:
					var name string
					var nameObj *ast.Object

					resName, err := GetIdentName(param)
					switch err != nil {
					case true:
						name = fmt.Sprintf("var%d", varCounter)
						varCounter++
					case false:
						name = resName.Name
						nameObj = resName.Obj
					}

					arguments = append(arguments, ArgType{
						Name:       name,
						NameObject: nameObj,
						Type:       getName(iobj),
						ExType:     getNameAsFromOuter(iobj, filepath.Base(pkg.Package)),
						TypeObject: iobj.Obj,
						Package:    paramPkg,
						BaseType:   defaultresType,
					})

				case *ast.SelectorExpr:
					// fmt.Printf("Result: %#v -> %#v -> %#q\n\n", iobj.X, iobj.Sel, iobj)
					xobj, ok := iobj.X.(*ast.Ident)
					if !ok {
						break
					}

					importDclr, _ := pkg.ImportFor(xobj.Name)

					name := fmt.Sprintf("var%d", varCounter)
					varCounter++

					arguments = append(arguments, ArgType{
						Name:    name,
						Import:  importDclr,
						Package: xobj.Name,
						Type:    getName(iobj),
						ExType:  getNameAsFromOuter(iobj, filepath.Base(pkg.Package)),
					})

				case *ast.StarExpr:
					name := fmt.Sprintf("var%d", varCounter)
					varCounter++

					var arg ArgType
					arg.Name = name
					arg.PointerType = iobj
					arg.Type = getName(iobj)
					arg.ExType = getNameAsFromOuter(iobj, filepath.Base(pkg.Package))

					switch value := iobj.X.(type) {
					case *ast.SelectorExpr:
						arg.ImportedObject = value

						vob, ok := value.X.(*ast.Ident)
						if !ok {
							break
						}

						importDclr, _ := pkg.ImportFor(vob.Name)

						arg.Package = vob.Name
						arg.Import = importDclr
					case *ast.StructType:
						arg.StructObject = value
					case *ast.ArrayType:
						arg.ArrayType = value
					case *ast.Ident:
						arg.IdentType = value
						arg.NameObject = value.Obj
						arg.Package = paramPkg
						arg.BaseType = defaultresType
					case *ast.ChanType:
						arg.ChanType = value
					}

					arguments = append(arguments, arg)

				case *ast.MapType:
					name := fmt.Sprintf("var%d", varCounter)
					varCounter++

					var arg ArgType
					arg.Name = name
					arg.MapType = iobj
					arg.Type = getName(iobj)
					arg.ExType = getNameAsFromOuter(iobj, filepath.Base(pkg.Package))

					if keySel, err := getSelector(iobj.Key); err == nil {
						if x, ok := keySel.X.(*ast.Ident); ok {
							if imported, err := pkg.ImportFor(x.Name); err == nil {
								arg.Import = imported
							}
						}
					}

					if valSel, err := getSelector(iobj.Value); err == nil {
						if x, ok := valSel.X.(*ast.Ident); ok {
							if imported, err := pkg.ImportFor(x.Name); err == nil {
								arg.Import2 = imported
							}
						}
					}

					arguments = append(arguments, arg)
				case *ast.ArrayType:
					name := fmt.Sprintf("var%d", varCounter)
					varCounter++

					var arg ArgType
					arg.Name = name
					arg.ArrayType = iobj
					arg.Type = getName(iobj)
					arg.ExType = getNameAsFromOuter(iobj, filepath.Base(pkg.Package))

					switch value := iobj.Elt.(type) {
					case *ast.SelectorExpr:
						arg.ImportedObject = value

						vob, ok := value.X.(*ast.Ident)
						if !ok {
							break
						}

						importDclr, _ := pkg.ImportFor(vob.Name)

						arg.Package = vob.Name
						arg.Import = importDclr
					case *ast.StarExpr:
						arg.PointerType = value
					case *ast.StructType:
						arg.StructObject = value
					case *ast.Ident:
						arg.IdentType = value
						arg.NameObject = value.Obj

						arg.Package = paramPkg
						arg.BaseType = defaultresType
					case *ast.ChanType:
						arg.ChanType = value
					}

					arguments = append(arguments, arg)
				case *ast.ChanType:
					name := fmt.Sprintf("var%d", varCounter)
					varCounter++

					var arg ArgType
					arg.Name = name
					arg.Type = getName(iobj)
					arg.ExType = getNameAsFromOuter(iobj, filepath.Base(pkg.Package))

					switch value := iobj.Value.(type) {
					case *ast.SelectorExpr:
						arg.ImportedObject = value

						vob, ok := value.X.(*ast.Ident)
						if !ok {
							break
						}

						importDclr, _ := pkg.ImportFor(vob.Name)

						arg.Package = vob.Name
						arg.Import = importDclr
					case *ast.StarExpr:
						arg.PointerType = value
					case *ast.StructType:
						arg.StructObject = value
					case *ast.ArrayType:
						arg.ArrayType = value
					case *ast.Ident:
						arg.IdentType = value
						arg.NameObject = value.Obj
						arg.Package = paramPkg
						arg.BaseType = defaultresType
					case *ast.ChanType:
						arg.ChanType = value
					}

					arguments = append(arguments, arg)
				default:
					// fmt.Printf("Param:Default: %#v -> %#q\n\n", iobj, iobj)
				}
			}

			defs = append(defs, FunctionDefinition{
				Func:      ftype,
				Interface: intr,
				Returns:   returns,
				Args:      arguments,
				Name:      nameIdent.Name,
			})

			continue
		}

		ident, ok := method.Type.(*ast.Ident)
		if !ok {
			continue
		}

		if ident == nil || ident.Obj == nil || ident.Obj.Decl == nil {
			continue
		}

		identDecl, ok := ident.Obj.Decl.(*ast.TypeSpec)
		if !ok {
			continue
		}

		identIntr, ok := identDecl.Type.(*ast.InterfaceType)
		if !ok {
			continue
		}

		defs = append(defs, GetInterfaceFunctions(identIntr, pkg)...)
	}

	return defs
}

var (
	naturalIdents = map[string]bool{
		"string":      true,
		"bool":        true,
		"rune":        true,
		"byte":        true,
		"int":         true,
		"int8":        true,
		"int16":       true,
		"int32":       true,
		"int64":       true,
		"uint":        true,
		"uint8":       true,
		"uint32":      true,
		"uint64":      true,
		"uintptr":     true,
		"float32":     true,
		"float64":     true,
		"complex128":  true,
		"complex64":   true,
		"error":       true,
		"struct":      true,
		"interface":   true,
		"interface{}": true,
		"struct{}":    true,
	}
)

// getPackageFromItem returns the package name associated with the type
// by attempting to retrieve it from a selector or final declaration name,
// and returns true/false if its part of go's base types.
func getPackageFromItem(item interface{}, defaultPkg string) (string, bool) {
	realName := getRealIdentName(item)

	if parts := strings.Split(realName, "."); len(parts) > 1 {
		if _, ok := naturalIdents[parts[1]]; ok {
			return "", true
		}

		return parts[0], false
	}

	if _, ok := naturalIdents[realName]; ok {
		return "", true
	}

	return defaultPkg, false
}

func getSelector(item interface{}) (*ast.SelectorExpr, error) {
	switch di := item.(type) {
	case *ast.StarExpr:
		return getSelector(di.X)
	case *ast.ArrayType:
		return getSelector(di.Elt)
	case *ast.ChanType:
		return getSelector(di.Value)
	case *ast.SelectorExpr:
		return di, nil
	default:
		return nil, errors.New("Has no selector")
	}
}

func getRealIdentName(item interface{}) string {
	switch di := item.(type) {
	case *ast.StarExpr:
		return getRealIdentName(di.X)
	case *ast.SelectorExpr:
		xobj, ok := di.X.(*ast.Ident)
		if !ok {
			return ""
		}

		return fmt.Sprintf("%s.%s", xobj.Name, di.Sel.Name)
	case *ast.Ident:
		return di.Name
	case *ast.ArrayType:
		return getRealIdentName(di.Elt)
	case *ast.ChanType:
		return getRealIdentName(di.Value)
	default:
		return ""
	}
}

func getNameAsFromOuter(item interface{}, basePkg string) string {
	switch di := item.(type) {
	case *ast.MapType:
		keyName := getNameAsFromOuter(di.Key, basePkg)
		valName := getNameAsFromOuter(di.Value, basePkg)
		return fmt.Sprintf("map[%s]%s", keyName, valName)
	case *ast.StarExpr:
		return fmt.Sprintf("*%s", getNameAsFromOuter(di.X, basePkg))
	case *ast.SelectorExpr:
		xobj, ok := di.X.(*ast.Ident)
		if !ok {
			return ""
		}

		return fmt.Sprintf("%s.%s", xobj.Name, di.Sel.Name)
	case *ast.StructType:
		return "struct{}"
	case *ast.Ident:
		if _, ok := naturalIdents[di.Name]; ok {
			return di.Name
		}

		return fmt.Sprintf("%s.%s", basePkg, di.Name)
	case *ast.ArrayType:
		if di.Len != nil {
			if dlen, ok := di.Len.(*ast.Ident); ok {
				return fmt.Sprintf("[%s]%s", dlen.Name, getNameAsFromOuter(di.Elt, basePkg))
			}

			if dlen, ok := di.Len.(*ast.BasicLit); ok {
				return fmt.Sprintf("[%s]%s", dlen.Value, getNameAsFromOuter(di.Elt, basePkg))
			}
		}

		return fmt.Sprintf("[]%s", getNameAsFromOuter(di.Elt, basePkg))
	case *ast.ChanType:
		return fmt.Sprintf("chan %s", getNameAsFromOuter(di.Value, basePkg))
	default:
		return ""
	}
}

func getName(item interface{}) string {
	switch di := item.(type) {
	case *ast.MapType:
		keyName := getName(di.Key)
		valName := getName(di.Value)
		return fmt.Sprintf("map[%s]%s", keyName, valName)
	case *ast.StarExpr:
		return fmt.Sprintf("*%s", getName(di.X))
	case *ast.SelectorExpr:
		xobj, ok := di.X.(*ast.Ident)
		if !ok {
			return ""
		}

		return fmt.Sprintf("%s.%s", xobj.Name, di.Sel.Name)
	case *ast.StructType:
		return "struct{}"
	case *ast.Ident:
		return di.Name
	case *ast.ArrayType:
		if di.Len != nil {
			if dlen, ok := di.Len.(*ast.Ident); ok {
				return fmt.Sprintf("[%s]%s", dlen.Name, getName(di.Elt))
			}

			if dlen, ok := di.Len.(*ast.BasicLit); ok {
				return fmt.Sprintf("[%s]%s", dlen.Value, getName(di.Elt))
			}
		}

		return fmt.Sprintf("[]%s", getName(di.Elt))
	case *ast.ChanType:
		return fmt.Sprintf("chan %s", getName(di.Value))
	default:
		return ""
	}
}

//===========================================================================================================

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
	PackageDeclr *PackageDeclaration                         `json:"-"`
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
	PackageDeclr *PackageDeclaration                         `json:"-"`
}

//===========================================================================================================

// TypeAnnotationGenerator defines a function which generates specific code related to the giving
// Annotation for a non-struct, non-interface type declaration. This allows you to apply and create
// new sources specifically for a giving type(non-struct, non-interface).
// It is responsible to fully contain all operations required to both generator any source and write such to
type TypeAnnotationGenerator func(string, AnnotationDeclaration, TypeDeclaration, PackageDeclaration) ([]gen.WriteDirective, error)

// StructAnnotationGenerator defines a function which generates specific code related to the giving
// Annotation. This allows you to generate a new source file containg source code for a giving struct type.
// It is responsible to fully contain all operations required to both generator any source and write such to.
type StructAnnotationGenerator func(string, AnnotationDeclaration, StructDeclaration, PackageDeclaration) ([]gen.WriteDirective, error)

// InterfaceAnnotationGenerator defines a function which generates specific code related to the giving
// Annotation. This allows you to generate a new source file containg source code for a giving interface type.
// It is responsible to fully contain all operations required to both generator any source and write such to
// appropriate files as intended, meta-data about package, and file paths are already include in the PackageDeclaration.
type InterfaceAnnotationGenerator func(string, AnnotationDeclaration, InterfaceDeclaration, PackageDeclaration) ([]gen.WriteDirective, error)

// PackageAnnotationGenerator defines a function which generates specific code related to the giving
// Annotation for a package. This allows you to apply and create new sources specifically because of a
// package wide annotation.
// It is responsible to fully contain all operations required to both generator any source and write such to
// All generators are expected to return
type PackageAnnotationGenerator func(string, AnnotationDeclaration, PackageDeclaration) ([]gen.WriteDirective, error)

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
func (a *AnnotationRegistry) ParseDeclr(declr PackageDeclaration, toDir string) ([]AnnotationWriteDirective, error) {
	var directives []AnnotationWriteDirective

	// Generate directives for package level
	for _, annotation := range declr.Annotations {
		generator, err := a.GetPackage(annotation.Name)
		if err != nil {
			continue
		}

		drs, err := generator(toDir, annotation, declr)
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

			drs, err := generator(toDir, annotation, inter, declr)
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

			drs, err := generator(toDir, annotation, structs, declr)
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

			drs, err := generator(toDir, annotation, typ, declr)
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

// Register which adds the generator depending on it's type into the appropriate
// registry. It only supports  the following generators:
// 1. TypeAnnotationGenerator (see Package ast#TypeAnnotationGenerator)
// 2. StructAnnotationGenerator (see Package ast#StructAnnotationGenerator)
// 3. InterfaceAnnotationGenerator (see Package ast#InterfaceAnnotationGenerator)
// 4. PackageAnnotationGenerator (see Package ast#PackageAnnotationGenerator)
// Any other type will cause the return of an error.
func (a *AnnotationRegistry) Register(name string, generator interface{}) error {
	switch gen := generator.(type) {
	case PackageAnnotationGenerator:
		a.RegisterPackage(name, gen)
		return nil
	case func(string, AnnotationDeclaration, PackageDeclaration) ([]gen.WriteDirective, error):
		a.RegisterPackage(name, gen)
		return nil
	case TypeAnnotationGenerator:
		a.RegisterType(name, gen)
		return nil
	case func(string, AnnotationDeclaration, TypeDeclaration, PackageDeclaration) ([]gen.WriteDirective, error):
		a.RegisterType(name, gen)
		return nil
	case StructAnnotationGenerator:
		a.RegisterStructType(name, gen)
		return nil
	case func(string, AnnotationDeclaration, StructDeclaration, PackageDeclaration) ([]gen.WriteDirective, error):
		a.RegisterStructType(name, gen)
		return nil
	case InterfaceAnnotationGenerator:
		a.RegisterInterfaceType(name, gen)
		return nil
	case func(string, AnnotationDeclaration, InterfaceDeclaration, PackageDeclaration) ([]gen.WriteDirective, error):
		a.RegisterInterfaceType(name, gen)
		return nil
	default:
		return fmt.Errorf("Generator type for %q not supported: %#v", name, generator)
	}
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

// FieldByFieldName defines a function to return actual name of field with the given tag name.
func FieldByFieldName(item StructDeclaration, fieldName string) (FieldDeclaration, error) {
	fields := Fields(GetFields(item))

	for _, field := range fields {
		if field.FieldName != fieldName {
			continue
		}

		return field, nil
	}

	return FieldDeclaration{}, fmt.Errorf("Field name %q for Struct %q", fieldName, item.Object.Name.Name)
}

// FieldFor defines a function to return actual name of field with the given tag name.
func FieldFor(item StructDeclaration, tag string, tagFieldName string) (FieldDeclaration, error) {
	fields := Fields(GetFields(item))

	wTags := fields.TagFor(tag)

	// Collect key field names from embedded first
	for _, tag := range wTags {
		if tag.Value != tagFieldName {
			continue
		}

		return tag.Field, nil
	}

	return FieldDeclaration{}, fmt.Errorf("Tag value %q not found in tag %q for Struct %q", tagFieldName, tag, item.Object.Name.Name)
}

// FieldNameFor defines a function to return actual name of field with the given tag name.
func FieldNameFor(item StructDeclaration, tag string, tagFieldName string) string {
	fields := Fields(GetFields(item))

	wTags := fields.TagFor(tag)

	// Collect key field names from embedded first
	for _, tag := range wTags {
		if tag.Value != tagFieldName {
			continue
		}

		return tag.Field.FieldName
	}

	return ""
}

// AssignDefaultValue will get the fieldName for a giving tag and tagVal and return	a string of giving
// variable name with fieldName equal to default value.
func AssignDefaultValue(item StructDeclaration, tag string, tagVal string, varName string) (string, error) {
	fieldName, defaultVal, err := DefaultFieldValueFor(item, tag, tagVal)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s.%s = %s", varName, fieldName, defaultVal), nil
}

// DefaultFieldValueFor defines a function to return a field default value.
func DefaultFieldValueFor(item StructDeclaration, tag string, tagVal string) (string, string, error) {
	fields := Fields(GetFields(item))

	wTags := fields.TagFor(tag)

	// Collect key field names from embedded first
	for _, tag := range wTags {
		if tag.Value != tagVal {
			continue
		}

		return tag.Field.FieldName, RandomFieldValue(tag.Field), nil
	}

	return "", "", fmt.Errorf("Field for tag value %q not found", tagVal)
}

// RandomFieldAssign generates a random Field of a giving struct and returns a variable assignment
// declaration with the types default value.
func RandomFieldAssign(item StructDeclaration, varName string, tag string, exceptions ...string) (string, error) {
	randomFieldVal, _, err := RandomFieldWithExcept(item, tag, exceptions...)
	if err != nil {
		return "", err
	}

	return AssignDefaultValue(item, tag, randomFieldVal, varName)
}

// RandomFieldWithExcept defines a function to return a random field name which is not
// included in the exceptions set.
func RandomFieldWithExcept(item StructDeclaration, tag string, exceptions ...string) (string, string, error) {
	fields := Fields(GetFields(item))

	wTags := fields.TagFor(tag)

	// Collect key field names from embedded first
	{
	ml:
		for _, tag := range wTags {
			for _, exception := range exceptions {
				if tag.Value == exception {
					continue ml
				}
			}

			return tag.Value, tag.Field.FieldName, nil
		}

	}

	return "", "", errors.New("All tags match exceptions")
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

// MapOutFieldsToJSON returns the giving map values containing string for the giving
// output.
func MapOutFieldsToJSON(item StructDeclaration, tagName, fallback string) (string, error) {
	document, err := MapOutFieldsToJSONWriter(item, tagName, fallback)
	if err != nil {
		return "", err
	}

	var doc bytes.Buffer

	if _, err := document.WriteTo(&doc); err != nil && err != io.EOF {
		return "", err
	}

	return doc.String(), nil
}

// MapOutFieldsToJSONWriter returns the giving map values containing string for the giving
// output.
func MapOutFieldsToJSONWriter(item StructDeclaration, tagName, fallback string) (io.WriterTo, error) {
	fields := Fields(GetFields(item))

	wTags := fields.TagFor(tagName)
	if len(wTags) == 0 {
		wTags = fields.TagFor(fallback)

		if len(wTags) == 0 {
			return nil, fmt.Errorf("No tags match for %q and %q fallback for struct %q", tagName, fallback, item.Object.Name)
		}
	}

	documents := make(map[string]io.WriterTo)

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

			document, err := MapOutFieldsToJSONWriter(StructDeclaration{
				Object: embededType,
				Struct: embedStruct,
			}, tagName, fallback)

			if err != nil {
				return nil, err
			}

			documents[tag.Value] = document
			continue
		}

		documents[tag.Value] = gen.Text(DefaultTypeValueString(strings.ToLower(tag.Field.FieldTypeName)))
	}

	return gen.JSONDocument(documents), nil
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

// RandomFieldValue returns the default value for a giving field.
func RandomFieldValue(fld FieldDeclaration) string {
	return RandomDataTypeValue(fld.FieldTypeName)
}

// DefaultFieldValue returns the default value for a giving field.
func DefaultFieldValue(fld FieldDeclaration) string {
	return DefaultTypeValueString(fld.FieldTypeName)
}

// RandomDataTypeValue returns the default value string of a giving
// typeName.
func RandomDataTypeValue(typeName string) string {
	switch typeName {
	case "uint", "uint32", "uint64":
		return fmt.Sprintf("%d", rand.Uint64())
	case "bool":
		return fmt.Sprintf("%t", rand.Int63n(1) == 0)
	case "string":
		return fmt.Sprintf("%q", fake.FullName())
	case "rune":
		return fmt.Sprintf("'%x'", fake.CharactersN(1))
	case "byte":
		return fmt.Sprintf("'%x'", fake.CharactersN(1))
	case "float32", "float64":
		return fmt.Sprintf("%.4f", rand.Float64())
	case "int", "int32", "int64":
		return fmt.Sprintf("%d", rand.Int63())
	default:
		return DefaultTypeValueString(typeName)
	}
}

// DefaultTypeValueString returns the default value string of a giving
// typeName.
func DefaultTypeValueString(typeName string) string {
	switch typeName {
	case "uint", "uint32", "uint64":
		return "0"
	case "bool":
		return `false`
	case "string":
		return `""`
	case "rune":
		return `rune(0)`
	case "[]uint":
		return `[]uint{}`
	case "[]uint64":
		return `[]uint64{}`
	case "[]uint32":
		return `[]uint32{}`
	case "[]int":
		return `[]int{}`
	case "[]int64":
		return `[]int64{}`
	case "[]int32":
		return `[]int32{}`
	case "[]bool":
		return `[]bool{}`
	case "[]string":
		return `[]string{}`
	case "[]byte":
		return `[]byte{}`
	case "byte":
		return `byte(rune(0))`
	case "float32", "float64":
		return "0.0"
	case "int", "int32", "int64":
		return "0"
	case "map[string]interface{}":
		return "map[string]interface{}"
	case "map[string]string":
		return "map[string]string{}"
	default:
		return "nil"
	}
}

// GetTag returns the giving tag associated with the name if it exists.
func GetTag(f FieldDeclaration, tagName string) (TagDeclaration, error) {
	return f.GetTag(tagName)
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

// ToValueString returns the string representation of a basic go core datatype.
func ToValueString(val interface{}) string {
	switch bo := val.(type) {
	case string:
		return strconv.Quote(bo)
	case int:
		return strconv.Itoa(bo)
	case int64:
		return strconv.Itoa(int(bo))
	case rune:
		return strconv.QuoteRune(bo)
	case bool:
		return strconv.FormatBool(bo)
	case byte:
		return strconv.QuoteRune(rune(bo))
	case float64:
		return strconv.FormatFloat(bo, 'f', 4, 64)
	case float32:
		return strconv.FormatFloat(float64(bo), 'f', 4, 64)
	default:
		data, err := json.Marshal(val)
		if err != nil {
			return err.Error()
		}

		return string(data)
	}
}
