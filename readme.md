Moz
=======
Moz exists has a library to provide a solid foundation for code generation by combining functional composition and Go ast for flexible content creation.

Install
-----------

```shell
go get -u github.com/influx6/moz
```

Introduction
----------------------------

Moz is a code generator which builds around the concepts of pluggable `io.WriteTo` elements that allow a elegant but capable system for generating code programmatically.

It uses functional compositions to define code structures that connect to create contents and using the Go ast parser, generates elegant structures for easier interaction with source files.


Features
----------

- Basic Programming structures
- Simple Coding Blocks
- Go text/template support
- Annotation code generation


Projects Using Moz
--------------------

- [Gu](https://github.com/gu-io/gu)
- [Dime](https://github.com/influx6/dime)

Code Generation with Moz
--------------------------

Moz is intended to be very barebones and minimal, it focuses around providing very basic structures, that allows the most flexibility in how you generate new content.

It provides two packages that are the center of it's system:

## [Gen](./gen)

Gen provides compositional structures for creating content with functions.

#### Generate Go struct using [Gen](./gen)

```go
import "github.com/influx6/moz/gen"

floppy := gen.Struct(
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

var source bytes.Buffer

floppy.WriteTo(&source) /*
// Floppy provides a basic function.
//
// Demonstration of using floppy API.
//
//
//@Flipo
//@API
type Floppy struct {

    Name string `json:"name"`

}
*/
```


Contributors
----------------
Please feel welcome to contribute with issues and PRs to improve Moz. :)
