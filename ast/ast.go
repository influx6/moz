package ast

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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
}

// TypeDeclaration defines a type which holds annotation data for a giving type declaration.
type TypeDeclaration struct {
	LineNumber  int                     `json:"line_number"`
	Column      int                     `json:"column"`
	Package     string                  `json:"package"`
	Path        string                  `json:"path"`
	FilePath    string                  `json:"filepath"`
	File        string                  `json:"file"`
	Object      *ast.TypeSpec           `json:"struct"`
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

						case *ast.ImportSpec:
							// fmt.Printf("Import: %+q -> %+q\n", obj, obj.Comment)
						}
					}

				case *ast.BadDecl:
				case *ast.FuncDecl:
				}
			}

			// fmt.Printf("Pkg: %q -> %#q\n", pkg.Name, packageDeclr)
			packageDeclrs = append(packageDeclrs, packageDeclr)
		}
	}

	return packageDeclrs, nil
}
