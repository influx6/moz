package moz_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/influx6/faux/tests"
	"github.com/influx6/moz"
)

// TestFunctionGen validates the expected output of a giving function generator.
func TestFunctionGen(t *testing.T) {
	expected := `func main(v int) {
	fmt.Printf("Welcome to Lola Land");
}`

	src := moz.Function(
		moz.Name("main"),
		moz.Constructor(
			moz.VarType(
				moz.Name("v"),
				moz.Type("int"),
			),
		),
		moz.Returns(),
		moz.Text(`	fmt.Printf("Welcome to Lola Land");`, nil),
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
