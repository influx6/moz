package ast_test

import (
	"bytes"
	"testing"

	"github.com/influx6/faux/tests"
	"github.com/influx6/moz/ast"
)

// TestAnnotationParserWithSingleComments validates the behaviour of the comment parsing method for
// parsing comments annotations.
func TestAnnotationParserWithSingleComments(t *testing.T) {
	reader := bytes.NewBufferString(`//
// The reason we wish to create this system is to allow the ability to parse such
// details outside of our scope.
//  @terminate(bob, dexter, 'jug')
//  @flatter
//  @templater(JSON, {
//    {
//      "name": "thunder"
//      "tribe": "COC",
//    }
//  })
//
//  @bob(1, 4, ['locksmith', "bob"])
//
//  @templater(Go, {
//  func Legover() string {
//      return "docking"
//    }
//  })
//
//  @sword(1, 4, 'bandate', "bob")
//`)

	annotations := ast.ReadAnnotationsFromCommentry(reader)
	if len(annotations) != 6 {
		tests.Info("Received: %+q", annotations)
		tests.Info("ReceivedLength: %d", len(annotations))
		tests.Failed("Should have successfully parsed 6 annotation markers from commentary")
	}
	tests.Passed("Should have successfully parsed 6 annotation markers from commentary")
}

// TestAnnotationParserWithText validates the behaviour of the comment parsing method for
// parsing comments annotations.
func TestAnnotationParserWithText(t *testing.T) {
	reader := bytes.NewBufferString(`
 The reason we wish to create this system is to allow the ability to parse such
 details outside of our scope.
  @terminate(bob, dexter, 'jug')
  @flatter
  @templater(JSON, {
    {
      "name": "thunder"
      "tribe": "COC",
    }
  })

  @bob(1, 4, ['locksmith', "bob"])

  @templater(Go, {
  func Legover() string {
      return "docking"
    }
})

  @sword(1, 4, 'bandate', "bob")
`)

	annotations := ast.ReadAnnotationsFromCommentry(reader)
	if len(annotations) != 6 {
		tests.Info("Received: %+q", annotations)
		tests.Info("ReceivedLength: %d", len(annotations))
		tests.Failed("Should have successfully parsed 6 annotation markers from commentary")
	}
	tests.Passed("Should have successfully parsed 6 annotation markers from commentary")
}

// TestAnnotationParserWithMultiComments validates the behaviour of the comment parsing method for
// parsing comments annotations.
func TestAnnotationParserWithMultiComments(t *testing.T) {
	reader := bytes.NewBufferString(`/*
* The reason we wish to create this system is to allow the ability to parse such
* details outside of our scope.
*  @terminate(bob, dexter, 'jug')
*  @flatter
*  @templater(JSON, {
*    {
*      "name": "thunder"
*      "tribe": "COC",
*    }
*  })
*
*  @bob(1, 4, ['locksmith', "bob"])
*
*  @templater(Go, {
*  func Legover() string {
*      return "docking"
*    }
*  })
*
*  @sword(1, 4, 'bandate', "bob")
*/`)

	annotations := ast.ReadAnnotationsFromCommentry(reader)
	if len(annotations) != 6 {
		tests.Failed("Should have successfully parsed 6 annotation markers from commentary")
	}
	tests.Passed("Should have successfully parsed 6 annotation markers from commentary")
}
