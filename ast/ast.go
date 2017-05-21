package mozgen

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

// Contains giving sets of variables exposing sytem GOPATH and GOPATHSRC.
var (
	GoPath    = os.Getenv("GOPATH")
	GoSrcPath = filepath.Join(GoPath, "src")
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
