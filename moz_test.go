package moz_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/influx6/faux/tests"
	"github.com/influx6/moz"
)

// TestStructGen validates the generation of a struct.
func TestStructGen(t *testing.T) {
	expected := "// Floppy provides a basic function.\n// \n// Demonstration of using floppy API.\n// \n//\n//@Flipo\n//@API\ntype Floppy struct {\n\n    Name string `json:\"name\"` \n\n}"

	src := moz.Struct(
		moz.Name("Floppy"),
		moz.Commentary(
			moz.Text("Floppy provides a basic function."),
			moz.Text("Demonstration of using floppy API."),
		),
		moz.Annotations(
			"Flipo",
			"API",
		),
		moz.Field(
			moz.Name("Name"),
			moz.Type("string"),
			moz.Tag("json", "name"),
		),
	)

	var bu bytes.Buffer

	if _, err := src.WriteTo(&bu); err != nil && err != io.EOF {
		tests.Failed("Should have successfully written source output: %+q.", err)
	}
	tests.Passed("Should have successfully written source output.")

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

	src := moz.Function(
		moz.Name("main"),
		moz.Constructor(
			moz.VarType(
				moz.Name("v"),
				moz.Type("int"),
			),
			moz.VarType(
				moz.Name("m"),
				moz.Type("string"),
			),
		),
		moz.Returns(),
		moz.SourceText(`	fmt.Printf("Welcome to Lola Land");`, nil),
	)

	var bu bytes.Buffer

	if _, err := src.WriteTo(&bu); err != nil && err != io.EOF {
		tests.Failed("Should have successfully written source output: %+q.", err)
	}
	tests.Passed("Should have successfully written source output.")

	if bu.String() != expected {
		tests.Info("Source: %+q", bu.String())
		tests.Info("Expected: %+q", expected)

		tests.Failed("Should have successfully matched generated output with expected.")
	}
	tests.Passed("Should have successfully matched generated output with expected.")
}
