package annotations

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/influx6/faux/fmtwriter"
	"github.com/influx6/moz"
	"github.com/influx6/moz/annotations/templates"
	"github.com/influx6/moz/ast"
	"github.com/influx6/moz/gen"
)

var (
	_ = moz.RegisterAnnotation("iface", IFaceAnnotationGenerator)
)

// IFaceAnnotationGenerator defines a code generator for creating a struct implementations for giving interface declaration.
// Annotation associated with this Generator is: @iface.
func IFaceAnnotationGenerator(toDir string, an ast.AnnotationDeclaration, itr ast.InterfaceDeclaration, pkg ast.PackageDeclaration) ([]gen.WriteDirective, error) {
	interfaceName := itr.Object.Name.Name
	interfaceNameLower := strings.ToLower(interfaceName)

	methods := itr.Methods()

	imports := make(map[string]string, 0)

	for _, method := range methods {

		// Retrieve all import paths for arguments.
		func(args []ast.ArgType) {
			for _, argument := range args {
				if argument.Import2.Path != "" {
					imports[argument.Import2.Path] = argument.Import2.Name
				}
				if argument.Import.Path != "" {
					imports[argument.Import.Path] = argument.Import.Name
				}
			}
		}(method.Args)

		// Retrieve all import paths for returns.
		func(args []ast.ArgType) {
			for _, argument := range args {
				if argument.Import2.Path != "" {
					imports[argument.Import2.Path] = argument.Import2.Name
				}
				if argument.Import.Path != "" {
					imports[argument.Import.Path] = argument.Import.Name
				}
			}
		}(method.Returns)

	}

	var wantedImports []gen.ImportItemDeclr

	for path, name := range imports {
		wantedImports = append(wantedImports, gen.Import(path, name))
	}

	implGen := gen.Block(
		gen.Package(
			gen.Name(pkg.Package),
			gen.Imports(wantedImports...),
			gen.Block(
				gen.SourceText(
					string(templates.Must("iface/iface.tml")),
					struct {
						InterfaceName string
						Package       ast.PackageDeclaration
						Methods       []ast.FunctionDefinition
					}{
						Package:       pkg,
						Methods:       methods,
						InterfaceName: interfaceName,
					},
				),
			),
		),
	)

	impImports := append([]gen.ImportItemDeclr{
		gen.Import("time", ""),
		gen.Import("runtime", ""),
		gen.Import(pkg.Path, ""),
	}, wantedImports...)

	implSnitchGen := gen.Block(
		gen.Package(
			gen.Name("snitch"),
			gen.Imports(impImports...),
			gen.Block(
				gen.SourceText(
					string(templates.Must("iface/iface-little-snitch.tml")),
					struct {
						Package       ast.PackageDeclaration
						InterfaceName string
						Methods       []ast.FunctionDefinition
					}{
						Package:       pkg,
						InterfaceName: interfaceName,
						Methods:       itr.Methods(),
					},
				),
			),
		),
	)

	testImports := append([]gen.ImportItemDeclr{
		gen.Import("testing", ""),
		gen.Import(pkg.Path, ""),
		gen.Import("github.com/influx6/faux/tests", ""),
		gen.Import(filepath.Join(pkg.Path, toDir, "snitch"), ""),
	}, wantedImports...)

	testGen := gen.Block(
		gen.Package(
			gen.Name(fmt.Sprintf("%s_test", strings.ToLower(pkg.Package))),
			gen.Imports(testImports...),
			gen.Block(
				gen.SourceText(
					string(templates.Must("iface/iface-test.tml")),
					struct {
						InterfaceName string
						Package       ast.PackageDeclaration
						Methods       []ast.FunctionDefinition
					}{
						Package:       pkg,
						InterfaceName: interfaceName,
						Methods:       itr.Methods(),
					},
				),
			),
		),
	)

	return []gen.WriteDirective{
		{
			Dir:          "snitch",
			Writer:       fmtwriter.New(implSnitchGen, true, true),
			FileName:     fmt.Sprintf("%s_little_snitch.go", interfaceNameLower),
			DontOverride: true,
		},
		{
			Writer:       fmtwriter.New(implGen, true, true),
			FileName:     fmt.Sprintf("%s_impl.go", interfaceNameLower),
			DontOverride: true,
		},
		{
			Writer:       fmtwriter.New(testGen, true, true),
			FileName:     fmt.Sprintf("%s_impl_test.go", interfaceNameLower),
			DontOverride: true,
		},
	}, nil
}
