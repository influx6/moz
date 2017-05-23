package moz_test

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/influx6/faux/tests"
	"github.com/influx6/moz"
)

func TestPackageGeneration(t *testing.T) {
	expected := `func main() {
	fmt.Printf("Welcome to Lola Land");
}`

	src := moz.Function(
		moz.Name("main"),
		moz.Constructor(),
		moz.Returns(),
		moz.Text(`	fmt.Printf("Welcome to Lola Land");`, nil),
	)

	var bu bytes.Buffer

	if _, err := src.WriteTo(&bu); err != nil && err != io.EOF {
		tests.Failed("Should have successfully written source output: %+q.", err)
	}
	tests.Passed("Should have successfully written source output.")

	fmt.Printf("Source: %+s\n", bu.String())
	fmt.Printf("Expected: %+s\n", expected)
}
