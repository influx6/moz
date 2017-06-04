package gen_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/influx6/faux/tests"
	"github.com/influx6/moz/gen"
)

// TestStructGen validates the generation of a struct.
func TestStructGen(t *testing.T) {
	expected := "// Floppy provides a basic function.\n// \n// Demonstration of using floppy API.\n// \n//\n//@Flipo\n//@API\ntype Floppy struct {\n\n    Name string `json:\"name\"` \n\n}"

	src := gen.Struct(
		gen.Name("Floppy"),
		gen.Commentary(
			gen.Text("Floppy provides a basic function."),
			gen.Text("Demonstration of using floppy API."),
		),
		gen.Annotations(
			"Flipo",
			"API",
		),
		gen.Field(
			gen.Name("Name"),
			gen.Type("string"),
			gen.Tag("json", "name"),
		),
	)

	var bu bytes.Buffer

	if _, err := src.WriteTo(&bu); err != nil && err != io.EOF {
		tests.Failed("Should have successfully written source output: %+q.", err)
	}
	tests.Passed("Should have successfully written source output.")

	tests.Info("Source: %+s", bu.String())
	if bu.String() != expected {
		tests.Info("Source: %+q", bu.String())
		tests.Info("Expected: %+q", expected)

		tests.Failed("Should have successfully matched generated output with expected.")
	}
	tests.Passed("Should have successfully matched generated output with expected.")
}

// TestFunctionGen validates the expected output of a giving function generator.
func TestFunctionGen(t *testing.T) {
	expected := `func main(v int, m string) {
	fmt.Printf("Welcome to Lola Land");
}`

	src := gen.Function(
		gen.Name("main"),
		gen.Constructor(
			gen.VarType(
				gen.Name("v"),
				gen.Type("int"),
			),
			gen.VarType(
				gen.Name("m"),
				gen.Type("string"),
			),
		),
		gen.Returns(),
		gen.SourceText(`	fmt.Printf("Welcome to Lola Land");`, nil),
	)

	var bu bytes.Buffer

	if _, err := src.WriteTo(&bu); err != nil && err != io.EOF {
		tests.Failed("Should have successfully written source output: %+q.", err)
	}
	tests.Passed("Should have successfully written source output.")

	tests.Info("Source: %+s", bu.String())

	if bu.String() != expected {
		tests.Info("Source: %+q", bu.String())
		tests.Info("Expected: %+q", expected)

		tests.Failed("Should have successfully matched generated output with expected.")
	}
	tests.Passed("Should have successfully matched generated output with expected.")
}