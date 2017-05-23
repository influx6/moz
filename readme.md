Moz
=======
Moz exists has a library to provide a solid foundation for code generation by combining basic programming structures and `text/templates` to provide a flexible and extensive generation capability.

Install
-----------

```shell
go get -u github.com/influx6/moz
```

Intro
--------
Moz is a code generator built around the concepts of pluggable `io.WriteTo` elements that allow a elegant but capable system for generating code programmatically. It uses functional compositions to define structures and how these structures connect to create desired output, which becomes rather easy to both understand and use.

Moz does not provide a complete set of all possible programming structures found in the Go programming language but provides a mixture of basic and needed structures with a go template strategy that allows us to quickly generate code structures, more so, moz provides a annotation strategy that provides a higher level of code generation based on a giving target to either add new methods or functions to the generation of a new breed of package based on that giving type. We hope to extend this system to allow the usage of `go1.8 Plugin`  system to allow extensibility to allow a variety of custom annotations outside of the `moz` package scope.

Features
----------

- Basic Programming structures
- Simple Coding Blocks
- Go text/template support 
- Annotation code generation (pending)


Example
-----------

- Generate a struct with moz

```go
floppy := moz.Struct(
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


- Generate a struct with moz

```go
main := moz.Function(
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

var source bytes.Buffer

main.WriteTo(&source) /*
func main(v int, m string) {
	fmt.Printf("Welcome to Lola Land");
}
*/
```