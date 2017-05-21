package moz

import (
	"bytes"
	"io"
	"strconv"

	"github.com/influx6/moz/templates"
)

//go:generate go generate ./templates/...

//======================================================================================================================

var (
	// CommaWriter defines the a writer that consistently writes a ','.
	CommaWriter = NewConstantWriter([]byte(","))

	// PeriodWriter defines the a writer that consistently writes a '.'.
	PeriodWriter = NewConstantWriter([]byte("."))
)

//======================================================================================================================

// Declaration defines a type which exposes a method to return a giving declaration
// source.
type Declaration interface {
	WriteTo(io.Writer) (int64, error)
}

//======================================================================================================================

// DeclarationMap defines a int64erface which maps giving declaration values
// int64o appropriate form for final output. It allows us create custom wrappers to
// define specific output style for a giving set of declarations.
type DeclarationMap interface {
	Map(...Declaration) Declaration
}

// MapOut defines an function type which maps giving
// data retrieved from a series of readers int64o the provided byte slice, returning
// the total number of data written and any error encountered.
type MapOut func(io.Writer, ...Declaration) (int64, error)

//======================================================================================================================

// Declarations defines the body contents of a giving declaration/structure.
type Declarations []Declaration

// Map applies a giving declaration mapper to the underlying io.Readers of the Declaration.
func (d Declarations) Map(mp DeclarationMap) Declaration {
	return mp.Map(d...)
}

//======================================================================================================================

// MapAnyWriter applies a giving set of MapOut functions with the provided int64ernal declarations
// writes to the provided io.Writer.
type MapAnyWriter struct {
	Map MapOut
	Dcl []Declaration
}

// WriteTo takes the data slice and writes int64ernal Declarations int64o the giving writer.
func (m MapAnyWriter) WriteTo(to io.Writer) (int64, error) {
	return m.Map(to, m.Dcl...)
}

//======================================================================================================================

// MapAny defines a struct which implements a structure which uses the provided
// int64ernal MapOut function to apply the necessary business logic of copying
// giving data space by a giving series of readers.
type MapAny struct {
	MapFn MapOut
}

// Map takes a giving set of readers returning a structure which implements the io.Reader int64erface
// for copying underlying data to the expected output.
func (mapper MapAny) Map(dls ...Declaration) Declaration {
	return MapAnyWriter{Map: mapper.MapFn, Dcl: dls}
}

//======================================================================================================================

// DotMapper defines a struct which implements the DeclarationMap which maps a set of
// items by seperating their output with a period '.', but execludes before the first and
// after the last item.
var DotMapper = MapAny{MapFn: func(to io.Writer, declrs ...Declaration) (int64, error) {
	wc := NewWriteCounter(to)

	total := len(declrs) - 1

	for index, declr := range declrs {
		if _, err := declr.WriteTo(wc); err != nil && err != io.EOF {
			return 0, err
		}

		if index > 1 && index < total {
			PeriodWriter.WriteTo(wc)
		}
	}

	return wc.Written(), nil
}}

// CommaMapper defines a struct which implements the DeclarationMap which maps a set of
// items by seperating their output with a coma ',', but execludes before the first and
// after the last item.
var CommaMapper = MapAny{MapFn: func(to io.Writer, declrs ...Declaration) (int64, error) {
	wc := NewWriteCounter(to)

	total := len(declrs) - 1

	for index, declr := range declrs {
		if _, err := declr.WriteTo(wc); err != nil && err != io.EOF {
			return 0, err
		}

		if index > 1 && index < total {
			CommaWriter.WriteTo(wc)
		}
	}

	return wc.Written(), nil
}}

//======================================================================================================================

// TypeDeclr defines a declaration struct for representing a giving type.
type TypeDeclr struct {
	TypeName string `json:"typeName"`
}

// String returns the int64ernal name associated with the TypeDeclr.
func (t TypeDeclr) String() string {
	return t.TypeName
}

// WriteTo writes to the provided writer the variable declaration.
func (t TypeDeclr) WriteTo(w io.Writer) (int64, error) {
	tml, err := ToTemplate("typeDeclr", templates.Must("variable-type-only.tml"), nil)
	if err != nil {
		return 0, err
	}

	wc := NewWriteCounter(w)

	if err := tml.Execute(wc, struct {
		Type string
	}{
		Type: t.TypeName,
	}); err != nil {
		return 0, err
	}

	return wc.Written(), nil
}

//======================================================================================================================

// NameDeclr defines a declaration struct for representing a giving value.
type NameDeclr struct {
	Name string `json:"name"`
}

// String returns the internal name associated with the NameDeclr.
func (n NameDeclr) String() string {
	return n.Name
}

//======================================================================================================================

// RuneASCIIDeclr defines a declaration struct for representing a giving value.
type RuneASCIIDeclr struct {
	Value rune `json:"value"`
}

// String returns the internal data associated with the structure.
func (n RuneASCIIDeclr) String() string {
	return strconv.QuoteRuneToASCII(n.Value)
}

// RuneGraphicsDeclr defines a declaration struct for representing a giving value.
type RuneGraphicsDeclr struct {
	Value rune `json:"value"`
}

// String returns the internal data associated with the structure.
func (n RuneGraphicsDeclr) String() string {
	return strconv.QuoteRuneToGraphic(n.Value)
}

// RuneDeclr defines a declaration struct for representing a giving value.
type RuneDeclr struct {
	Value rune `json:"value"`
}

// String returns the internal data associated with the structure.
func (n RuneDeclr) String() string {
	return strconv.QuoteRune(n.Value)
}

// StringASCIIDeclr defines a declaration struct for representing a giving value.
type StringASCIIDeclr struct {
	Value string `json:"value"`
}

// String returns the internal data associated with the structure.
func (n StringASCIIDeclr) String() string {
	return strconv.QuoteToASCII(n.Value)
}

// StringDeclr defines a declaration struct for representing a giving value.
type StringDeclr struct {
	Value string `json:"value"`
}

// String returns the internal data associated with the structure.
func (n StringDeclr) String() string {
	return strconv.Quote(n.Value)
}

// BoolDeclr defines a declaration struct for representing a giving value.
type BoolDeclr struct {
	Value bool `json:"value"`
}

// String returns the internal data associated with the structure.
func (n BoolDeclr) String() string {
	return strconv.FormatBool(n.Value)
}

// UIntBaseDeclr defines a declaration struct for representing a giving value.
type UIntBaseDeclr struct {
	Value uint64 `json:"value"`
	Base  int    `json:"base"`
}

// String returns the internal data associated with the structure.
func (n UIntBaseDeclr) String() string {
	return strconv.FormatUint(n.Value, n.Base)
}

// UInt64Declr defines a declaration struct for representing a giving value.
type UInt64Declr struct {
	Value uint64 `json:"value"`
}

// String returns the internal data associated with the structure.
func (n UInt64Declr) String() string {
	return strconv.FormatUint(n.Value, 10)
}

// UInt32Declr defines a declaration struct for representing a giving value.
type UInt32Declr struct {
	Value uint32 `json:"value"`
}

// String returns the internal data associated with the structure.
func (n UInt32Declr) String() string {
	return strconv.FormatUint(uint64(n.Value), 10)
}

// IntBaseDeclr defines a declaration struct for representing a giving value.
type IntBaseDeclr struct {
	Value int64 `json:"value"`
	Base  int   `json:"base"`
}

// String returns the internal data associated with the structure.
func (n IntBaseDeclr) String() string {
	return strconv.FormatInt(n.Value, n.Base)
}

// Int64Declr defines a declaration struct for representing a giving value.
type Int64Declr struct {
	Value int64 `json:"value"`
}

// String returns the internal data associated with the structure.
func (n Int64Declr) String() string {
	return strconv.FormatInt(n.Value, 10)
}

// Int32Declr defines a declaration struct for representing a giving value.
type Int32Declr struct {
	Value int32 `json:"value"`
}

// String returns the internal data associated with the structure.
func (n Int32Declr) String() string {
	return strconv.FormatInt(int64(n.Value), 10)
}

// IntDeclr defines a declaration struct for representing a giving value.
type IntDeclr struct {
	Value int `json:"value"`
}

// String returns the internal data associated with the structure.
func (n IntDeclr) String() string {
	return strconv.Itoa(n.Value)
}

// FloatBaseDeclr defines a declaration struct for representing a giving value.
type FloatBaseDeclr struct {
	Value     float64 `json:"value"`
	Bitsize   int     `json:"base"`
	Precision int     `json:"precision"`
}

// String returns the internal data associated with the structure.
func (n FloatBaseDeclr) String() string {
	return strconv.FormatFloat(n.Value, 'f', n.Precision, n.Bitsize)
}

// Float32Declr defines a declaration struct for representing a giving value.
type Float32Declr struct {
	Value float32 `json:"value"`
}

// String returns the internal data associated with the structure.
func (n Float32Declr) String() string {
	return strconv.FormatFloat(float64(n.Value), 'f', 4, 32)
}

// Float64Declr defines a declaration struct for representing a giving value.
type Float64Declr struct {
	Value float64 `json:"value"`
}

// String returns the internal data associated with the structure.
func (n Float64Declr) String() string {
	return strconv.FormatFloat(n.Value, 'f', 4, 64)
}

// ValueDeclr defines a declaration struct for representing a giving value.
type ValueDeclr struct {
	Value          interface{}              `json:"value"`
	ValueConverter func(interface{}) string `json:"-"`
}

// String returns the internal data associated with the structure.
func (n ValueDeclr) String() string {
	return n.ValueConverter(n.Value)
}

//======================================================================================================================

// SliceTypeDeclr defines a declaration struct for representing a go slice.
type SliceTypeDeclr struct {
	Type TypeDeclr `json:"type"`
}

// WriteTo writes to the provided writer the variable declaration.
func (t SliceTypeDeclr) WriteTo(w io.Writer) (int64, error) {
	tml, err := ToTemplate("sliceTypeDeclr", templates.Must("slicetype.tml"), nil)
	if err != nil {
		return 0, err
	}

	wc := NewWriteCounter(w)

	if err := tml.Execute(wc, t); err != nil {
		return 0, err
	}

	return wc.Written(), nil
}

// SliceDeclr defines a declaration struct for representing a go slice.
type SliceDeclr struct {
	Type   TypeDeclr     `json:"type"`
	Values []Declaration `json:"values"`
}

// WriteTo writes to the provided writer the variable declaration.
func (t SliceDeclr) WriteTo(w io.Writer) (int64, error) {
	var vam bytes.Buffer

	if _, err := CommaMapper.Map(t.Values...).WriteTo(&vam); err != nil && err != io.EOF {
		return 0, err
	}

	tml, err := ToTemplate("sliceDeclr", templates.Must("slicevalue.tml"), nil)
	if err != nil {
		return 0, err
	}

	wc := NewWriteCounter(w)

	if err := tml.Execute(wc, struct {
		Type   string
		Values string
	}{
		Type:   t.Type.String(),
		Values: vam.String(),
	}); err != nil {
		return 0, err
	}

	return wc.Written(), nil
}

//======================================================================================================================

// Contains different sets of operator declarations.
var (
	PlusOperator           = OperatorDeclr{Operation: "+"}
	MinusOperator          = OperatorDeclr{Operation: "-"}
	ModeOperator           = OperatorDeclr{Operation: "%"}
	DivideOperator         = OperatorDeclr{Operation: "/"}
	MultiplicationOperator = OperatorDeclr{Operation: "*"}
	EqualOperator          = OperatorDeclr{Operation: "=="}
	LessThanOperator       = OperatorDeclr{Operation: "<"}
	MoreThanOperator       = OperatorDeclr{Operation: ">"}
	LessThanEqualOperator  = OperatorDeclr{Operation: "<="}
	MoreThanEqualOperator  = OperatorDeclr{Operation: ">="}
	NotEqualOperator       = OperatorDeclr{Operation: "!="}
	ANDOperator            = OperatorDeclr{Operation: "&&"}
	OROperator             = OperatorDeclr{Operation: "||"}
	BinaryANDOperator      = OperatorDeclr{Operation: "&"}
	BinaryOROperator       = OperatorDeclr{Operation: "|"}
	DecrementOperator      = OperatorDeclr{Operation: "--"}
	IncrementOperator      = OperatorDeclr{Operation: "++"}
)

// OperatorDeclr defines a declaration which produces a variable declaration.
type OperatorDeclr struct {
	Operation string `json:"operation"`
}

// String returns the internal name associated with the struct.
func (n OperatorDeclr) String() string {
	return n.Operation
}

// WriteTo writes the giving representation into the provided writer.
func (n OperatorDeclr) WriteTo(w io.Writer) (int64, error) {
	total, err := w.Write([]byte(n.Operation))
	return int64(total), err
}

//======================================================================================================================

// VariableTypeDeclr defines a declaration which produces a variable declaration.
type VariableTypeDeclr struct {
	Name NameDeclr `json:"name"`
	Type TypeDeclr `json:"typename"`
}

// WriteTo writes to the provided writer the variable declaration.
func (v VariableTypeDeclr) WriteTo(w io.Writer) (int64, error) {
	tml, err := ToTemplate("variableDeclr", templates.Must("variable-type.tml"), nil)
	if err != nil {
		return 0, err
	}

	wc := NewWriteCounter(w)

	if err := tml.Execute(wc, v); err != nil {
		return 0, err
	}

	return wc.Written(), nil
}

// VariableNameDeclr defines a declaration which produces a variable declaration.
type VariableNameDeclr struct {
	Name NameDeclr `json:"name"`
}

// WriteTo writes to the provided writer the variable declaration.
func (v VariableNameDeclr) WriteTo(w io.Writer) (int64, error) {
	tml, err := ToTemplate("variableDeclr", templates.Must("variable-name.tml"), nil)
	if err != nil {
		return 0, err
	}

	wc := NewWriteCounter(w)

	if err := tml.Execute(wc, v); err != nil {
		return 0, err
	}

	return wc.Written(), nil
}

// VariableAssignmentDeclr defines a declaration which produces a variable declaration.
type VariableAssignmentDeclr struct {
	Name  NameDeclr   `json:"name"`
	Value Declaration `json:"value"`
}

// WriteTo writes to the provided writer the variable declaration.
func (v VariableAssignmentDeclr) WriteTo(w io.Writer) (int64, error) {
	tml, err := ToTemplate("variableDeclr", templates.Must("variable-assign-basic.tml"), nil)
	if err != nil {
		return 0, err
	}

	var vam bytes.Buffer

	if _, err := v.Value.WriteTo(&vam); err != nil && err != io.EOF {
		return 0, err
	}

	wc := NewWriteCounter(w)

	if err := tml.Execute(wc, struct {
		Name  string
		Value string
	}{
		Name:  v.Name.String(),
		Value: vam.String(),
	}); err != nil {
		return 0, err
	}

	return wc.Written(), nil
}

// VariableShortAssignmentDeclr defines a declaration which produces a variable declaration.
type VariableShortAssignmentDeclr struct {
	Name  NameDeclr   `json:"name"`
	Value Declaration `json:"value"`
}

// WriteTo writes to the provided writer the variable declaration.
func (v VariableShortAssignmentDeclr) WriteTo(w io.Writer) (int64, error) {
	tml, err := ToTemplate("variableDeclr", templates.Must("variable-assign.tml"), nil)
	if err != nil {
		return 0, err
	}

	var vam bytes.Buffer

	if _, err := v.Value.WriteTo(&vam); err != nil && err != io.EOF {
		return 0, err
	}

	wc := NewWriteCounter(w)

	if err := tml.Execute(wc, struct {
		Name  string
		Value string
	}{
		Name:  v.Name.String(),
		Value: vam.String(),
	}); err != nil {
		return 0, err
	}

	return wc.Written(), nil
}

//======================================================================================================================

// SingleByteBlockDeclr defines a declaration which produces a block byte slice which is written to a writer.
// declaration writer into it's block char.
// eg. A BlockDeclr with Char '{{'
// 		Will produce '{{DataFROMWriter' output.
type SingleByteBlockDeclr struct {
	Block []byte `json:"block"`
}

// WriteTo writes the giving representation into the provided writer.
func (b SingleByteBlockDeclr) WriteTo(w io.Writer) (int64, error) {
	wc := NewWriteCounter(w)

	if _, err := wc.Write(b.Block); err != nil {
		return 0, err
	}

	return wc.Written(), nil
}

// SingleBlockDeclr defines a declaration which produces a block char which is written to a writer.
// eg. A BlockDeclr with Char '{'
// 		Will produce '{' output.
type SingleBlockDeclr struct {
	Rune rune `json:"rune"`
}

// WriteTo writes the giving representation into the provided writer.
func (b SingleBlockDeclr) WriteTo(w io.Writer) (int64, error) {
	wc := NewWriteCounter(w)

	if _, err := wc.Write([]byte{byte(b.Rune)}); err != nil {
		return 0, err
	}

	return wc.Written(), nil
}

// ByteBlockDeclr defines a declaration which produces a block cover which wraps any other
// declaration writer into it's block char.
// eg. A BlockDeclr with Char '{''}'
// 		Will produce '{{DataFROMWriter}}' output.
type ByteBlockDeclr struct {
	Block      Declaration `json:"block"`
	BlockBegin []byte      `json:"begin"`
	BlockEnd   []byte      `json:"end"`
}

// WriteTo writes the giving representation into the provided writer.
func (b ByteBlockDeclr) WriteTo(w io.Writer) (int64, error) {
	wc := NewWriteCounter(w)

	if _, err := wc.Write(b.BlockBegin); err != nil {
		return 0, err
	}

	if _, err := b.Block.WriteTo(wc); err != nil && err != io.EOF {
		return 0, err
	}

	if _, err := wc.Write(b.BlockEnd); err != nil {
		return 0, err
	}

	return wc.Written(), nil
}

// BlockDeclr defines a declaration which produces a block cover which wraps any other
// declaration writer into it's block char.
// eg. A BlockDeclr with Char '{''}'
// 		Will produce '{DataFROMWriter}' output.
type BlockDeclr struct {
	Block     io.WriterTo `json:"block"`
	RuneBegin rune        `json:"begin"`
	RuneEnd   rune        `json:"end"`
}

// WriteTo writes the giving representation into the provided writer.
func (b BlockDeclr) WriteTo(w io.Writer) (int64, error) {
	wc := NewWriteCounter(w)

	if _, err := wc.Write([]byte{byte(b.RuneBegin)}); err != nil {
		return 0, err
	}

	if _, err := b.Block.WriteTo(wc); err != nil && err != io.EOF {
		return 0, err
	}

	if _, err := wc.Write([]byte{byte(b.RuneEnd)}); err != nil {
		return 0, err
	}

	return wc.Written(), nil
}

//======================================================================================================================

// ConditionDeclr defines a declaration which produces a variable declaration.
type ConditionDeclr struct {
	PreVar   VariableNameDeclr `json:"prevar"`
	PostVar  VariableNameDeclr `json:"postvar"`
	Operator OperatorDeclr     `json:"operator"`
}

// WriteTo writes the giving representation into the provided writer.
func (c ConditionDeclr) WriteTo(w io.Writer) (int64, error) {
	wc := NewWriteCounter(w)

	if _, err := c.PreVar.WriteTo(wc); err != nil && err != io.EOF {
		return 0, err
	}

	if _, err := c.Operator.WriteTo(wc); err != nil && err != io.EOF {
		return 0, err
	}

	if _, err := c.PostVar.WriteTo(wc); err != nil && err != io.EOF {
		return 0, err
	}

	return wc.Written(), nil
}

//======================================================================================================================

// FunctionDeclr defines a declaration which produces function about based on the giving
// constructor and body.
type FunctionDeclr struct {
	Name        NameDeclr                `json:"name"`
	Constructor FunctionConstructorDeclr `json:"constructor"`
	Body        FunctionBodyDeclr        `json:"body"`
}

// WriteTo writes to the provided writer the function declaration.
func (f FunctionDeclr) WriteTo(w io.Writer) (int64, error) {
	var constr, body bytes.Buffer

	if _, err := f.Constructor.WriteTo(&constr); err != nil {
		return 0, err
	}

	if _, err := f.Body.WriteTo(&body); err != nil {
		return 0, err
	}

	var declr = struct {
		Name        string
		Body        string
		Constructor string
	}{
		Name:        f.Name.String(),
		Body:        body.String(),
		Constructor: constr.String(),
	}

	tml, err := ToTemplate("functionDeclr", templates.Must("function.tml"), nil)
	if err != nil {
		return 0, err
	}

	wc := NewWriteCounter(w)

	if err := tml.Execute(wc, declr); err != nil {
		return 0, err
	}

	return wc.Written(), nil
}

//======================================================================================================================

// FunctionReturnDeclr defines a declaration which produces argument based output
// of it's giving int64ernals.
type FunctionReturnDeclr struct {
	Returns []TypeDeclr `json:"returns"`
}

// WriteTo writes to the provided writer the function argument declaration.
func (f FunctionReturnDeclr) WriteTo(w io.Writer) (int64, error) {
	var decals []Declaration

	for _, item := range f.Returns {
		decals = append(decals, item)
	}

	arguments := CommaMapper.Map(decals...)

	return (BlockDeclr{
		Block:     arguments,
		RuneBegin: '(',
		RuneEnd:   ')',
	}).WriteTo(w)
}

//======================================================================================================================

// FunctionConstructorDeclr defines a declaration which produces argument based output
// of it's giving int64ernals.
type FunctionConstructorDeclr struct {
	Arguments []VariableTypeDeclr `json:"constructor"`
}

// WriteTo writes to the provided writer the function argument declaration.
func (f FunctionConstructorDeclr) WriteTo(w io.Writer) (int64, error) {
	var decals []Declaration

	for _, item := range f.Arguments {
		decals = append(decals, item)
	}

	arguments := CommaMapper.Map(decals...)

	return (BlockDeclr{
		Block:     arguments,
		RuneBegin: '(',
		RuneEnd:   ')',
	}).WriteTo(w)
}

//======================================================================================================================

// FunctionBodyDeclr defines a type used to define the contents of a body.
type FunctionBodyDeclr struct {
	Body []Declaration `json:"body"`
}

// WriteTo writes to the provided writer the variable declaration.
func (f FunctionBodyDeclr) WriteTo(w io.Writer) (int64, error) {
	var total int64

	for _, item := range f.Body {
		nid, err := item.WriteTo(w)
		if err != nil {
			return total, err
		}

		total += nid
	}

	return total, nil
}

//======================================================================================================================
