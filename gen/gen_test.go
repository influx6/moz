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
	expected := "// Floppy provides a basic function.\n// \n// Demonstration of using floppy API.\n// \n//\n\n//@Flipo\n//@API\ntype Floppy struct {\n\n    Name string `json:\"name\"` \n\n}"

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

	if bu.String() != expected {
		tests.Info("Source: %+q", bu.String())
		tests.Info("Expected: %+q", expected)

		tests.Failed("Should have successfully matched generated output with expected.")
	}
	tests.Passed("Should have successfully matched generated output with expected.")
}

// TestMapGen validates the expected output of a giving map generator.
func TestMapGen(t *testing.T) {
	expected := "map[string]interface{}{\n    \n        \"name\": \"alex\",\n    \n        \"type\": 12,\n    \n}"

	src := gen.Map(
		"string",
		"interface{}",
		map[string]io.WriterTo{
			"type": gen.Name("12"),
			"name": gen.Name(`"alex"`),
		},
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
