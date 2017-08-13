package annotations

import (
	"fmt"
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
func IFaceAnnotationGenerator(an ast.AnnotationDeclaration, itr ast.InterfaceDeclaration, pkg ast.PackageDeclaration) ([]gen.WriteDirective, error) {
	interfaceName := itr.Object.Name.Name
	interfaceNameLower := strings.ToLower(interfaceName)

	implGen := gen.Block(
		gen.Package(
			gen.Name(pkg.Package),
			gen.Imports(),
			gen.Block(
				gen.SourceText(
					string(templates.Must("iface/iface.tml")),
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

	implSnitchGen := gen.Block(
		gen.Package(
			gen.Name("snitch"),
			gen.Imports(
				gen.Import("time", ""),
				gen.Import("runtime", ""),
				gen.Import(pkg.Path, ""),
			),
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

	testGen := gen.Block(
		gen.Package(
			gen.Name(fmt.Sprintf("%s_test", strings.ToLower(pkg.Package))),
			gen.Imports(
				gen.Import("testing", ""),
				gen.Import(pkg.Path, ""),
				gen.Import("github.com/influx6/faux/tests", ""),
				gen.Import(fmt.Sprintf("%s/%s", pkg.Path, "snitch"), ""),
			),
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
			Dir:      "snitch",
			Writer:   fmtwriter.New(implSnitchGen, true, true),
			FileName: fmt.Sprintf("%s_little_snitch.go", interfaceNameLower),
			// DontOverride: true,
		},
		{
			Writer:   fmtwriter.New(implGen, true, true),
			FileName: fmt.Sprintf("%s_impl.go", interfaceNameLower),
			// DontOverride: true,
		},
		{
			Writer:   fmtwriter.New(testGen, true, true),
			FileName: fmt.Sprintf("%s_impl_test.go", interfaceNameLower),
			// DontOverride: true,
		},
	}, nil
}
