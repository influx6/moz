Moz
=======
Moz exists has a library to provide a solid foundation for code generation by combining basic programming structures and `text/templates` to provide a flexible and extensive generation capability.

Install
-----------

```shell
go get -u github.com/influx6/moz
```

Introduction
----------------------------

Moz is a code generator which builds around the concepts of pluggable `io.WriteTo` elements that allow a elegant but capable system for generating code programmatically. It uses functional compositions to define structures and how these structures connect to create desired output, which becomes rather easy to both understand and use.

Moz does not provide a complete set of all possible programming structures found in the Go programming language but provides a mixture of basic and needed structures with a go template strategy that allows us to quickly generate code structures, more so, moz provides a annotation strategy that provides a higher level of code generation based on a target, which either will generate new methods/functions or generate new packages based on that target.

Features
----------

- Basic Programming structures
- Simple Coding Blocks
- Go text/template support
- Annotation code generation

Future Plans
---------------

- Extend Plugin to HotLoad with `go.18 Plugin`.


Annotation Code Generation
----------------------------

Moz provides a annotation style code generation system apart from it's code generation structures. This is provide to allow descriptive annotations to be added to giving Go structures (`interface`, `struct`, `type alises`) within their comments and as well as to the package.

This annotation then are passed by the moz `annotation` CLI tooling which can generate a series of files and packages based on the internal logic of the generator associated with that annotation to meet the needs for that type.

For example: If we wanted to be able to generate code for database CRUD activities without having to use ORMs or write such code manually, with the Moz annotation code generation ability, we can create a `struct` generator that can use a `@mongo` annotation, which generates mongo CRUD functions which expect such a type and perform the appropriate CRUD operations.

See the [Example](./examples/) directory, which demonstrates use of annotations to code generate other parts of a project or mock up implementation detail for an interface using annotations.


### How Annotation Code Generation works

Moz provides 4 types of Annotation generators, which are function types which provide the necessary operations to be performed to create the underline series of sources to be generated for each annotation.

Moz provide the following generators type functions:

#### StructType Code Generators

This functions are specific to provide code generation instructions for struct type declarations which the given annotation is attached to.

```go
type StructAnnotationGenerator func(string, AnnotationDeclaration, StructDeclaration, PackageDeclaration) ([]gen.WriteDirective, error)
```

*This function is expected to return a slice of `WriteDirective` which contains file name, `WriterTo` object and a possible `Dir` relative path which the contents should be written to.*

#### InterfaceType Code Generators

This functions are specific to provide code generation instructions for interface declarations which the given annotation is attached to.

```go
type InterfaceAnnotationGenerator func(string,AnnotationDeclaration, InterfaceDeclaration, PackageDeclaration) ([]gen.WriteDirective, error)
```

*This function is expected to return a slice of `WriteDirective` which contains file name, `WriterTo` object and a possible `Dir` relative path which the contents should be written to.*

#### PackageType Code Generators

This functions are specific to provide code generation instructions for given annotation declared on the package comment block.

```go
type PackageAnnotationGenerator func(string, AnnotationDeclaration, PackageDeclaration) ([]gen.WriteDirective, error)
```

*This function is expected to return a slice of `WriteDirective` which contains file name, `WriterTo` object and a possible `Dir` relative path which the contents should be written to.*

#### Non(Struct|Interface) Code Generators

This functions are specific to provide code generation instructions for non-struct and non-interface declarations which the given annotation is attached to.

```go
type TypeAnnotationGenerator func(string, AnnotationDeclaration, TypeDeclaration, PackageDeclaration) ([]gen.WriteDirective, error)
```

*This function is expected to return a slice of `WriteDirective` which contains file name, `WriterTo` object and a possible `Dir` relative path which the contents should be written to.*


Code Generation structures
---------------------------

Moz provides the [Gen Package](./gen) which defines sets of structures which define specific code structures and are used to built a programmatically combination that define the expected code to be produced. It also provides a functional composition style functions that provide a cleaner and more descriptive approach to how these blocks are combined.

The code gen is heavily geared towards the use of `text/template` but also ensures to be flexible to provide non-template based structures that work as well.

### Example

- Generate a struct with moz

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


- Generate a function with moz

```go
import "github.com/influx6/moz/gen"

main := gen.Function(
    gen.Name("main"),
    gen.Constructor(
        gen.FieldType(
            gen.Name("v"),
            gen.Type("int"),
        ),
        gen.FieldType(
            gen.Name("m"),
            gen.Type("string"),
        ),
    ),
    gen.Returns(),
    gen.SourceText(`	fmt.Printf("Welcome to Lola Land");`, nil),
)

var source bytes.Buffer

main.WriteTo(&source) /*
func main(v int, m string) {
	fmt.Printf("Welcome to Lola Land");
}
*/
```
