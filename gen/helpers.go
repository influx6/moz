package gen

import (
	"bytes"
	"io"
	"text/template"
)

// RuneASCII returns a new instance of a RuneASCIIDeclr.
func RuneASCII(rn rune) RuneASCIIDeclr {
	return RuneASCIIDeclr{
		Value: rn,
	}
}

// RuneGraphics returns a new instance of a RuneGraphicsDeclr.
func RuneGraphics(rn rune) RuneGraphicsDeclr {
	return RuneGraphicsDeclr{
		Value: rn,
	}
}

// Rune returns a new instance of a RuneGraphicsDeclr.
func Rune(rn rune) RuneDeclr {
	return RuneDeclr{
		Value: rn,
	}
}

// StringASCII returns a new instance of a StringASCIIDeclr.
func StringASCII(rn string) StringASCIIDeclr {
	return StringASCIIDeclr{
		Value: rn,
	}
}

// Bool returns a new instance of a BoolDeclr.
func Bool(rn bool) BoolDeclr {
	return BoolDeclr{
		Value: rn,
	}
}

// FloatBase returns a new instance of a FloatBaseDeclr.
func FloatBase(rn float64, bit, prec int) FloatBaseDeclr {
	return FloatBaseDeclr{
		Value:     rn,
		Bitsize:   bit,
		Precision: prec,
	}
}

// Float32 returns a new instance of a Float32.
func Float32(rn float32) Float32Declr {
	return Float32Declr{
		Value: rn,
	}
}

// Float64 returns a new instance of a Float64.
func Float64(rn float64) Float64Declr {
	return Float64Declr{
		Value: rn,
	}
}

// IntBase returns a new instance of a UIntBaseDeclr.
func IntBase(rn int64, base int) IntBaseDeclr {
	return IntBaseDeclr{
		Value: rn,
		Base:  base,
	}
}

// Int32 returns a new instance of a Int32Declr.
func Int32(rn int32) Int32Declr {
	return Int32Declr{
		Value: rn,
	}
}

// Int64 returns a new instance of a Int64Declr.
func Int64(rn int64) Int64Declr {
	return Int64Declr{
		Value: rn,
	}
}

// Int returns a new instance of a IntDeclr.
func Int(rn int) IntDeclr {
	return IntDeclr{
		Value: rn,
	}
}

// UIntBase returns a new instance of a UIntBaseDeclr.
func UIntBase(rn uint64, base int) UIntBaseDeclr {
	return UIntBaseDeclr{
		Value: rn,
		Base:  base,
	}
}

// UInt32 returns a new instance of a UInt32.
func UInt32(rn uint32) UInt32Declr {
	return UInt32Declr{
		Value: rn,
	}
}

// UInt64 returns a new instance of a UInt64.
func UInt64(rn uint64) UInt64Declr {
	return UInt64Declr{
		Value: rn,
	}
}

// Value returns a new instance of a ValueDeclr.
func Value(rn interface{}, converter func(interface{}) string) ValueDeclr {
	return ValueDeclr{
		Value:          rn,
		ValueConverter: converter,
	}
}

// Text returns a new instance of a TextDeclr.
func Text(txt string) TextBlockDeclr {
	return TextBlockDeclr{
		Block: txt,
	}
}

// String returns a new instance of a StringDeclr.
func String(rn string) StringDeclr {
	return StringDeclr{
		Value: rn,
	}
}

// Imports returns a new instance of a ImportDeclr.
func Imports(ims ...ImportItemDeclr) ImportDeclr {
	return ImportDeclr{
		Packages: ims,
	}
}

// Import returns a new instance of a ImportItemDeclr.
func Import(path string, namespace string) ImportItemDeclr {
	return ImportItemDeclr{
		Path:      path,
		Namespace: namespace,
	}
}

// Type returns a new instance of a TypeDeclr.
func Type(name string) TypeDeclr {
	return TypeDeclr{
		TypeName: name,
	}
}

// Name returns a new instance of a NameDeclr.
func Name(name string) NameDeclr {
	return NameDeclr{
		Name: name,
	}
}

// Package returns a new instance of a PackageDeclr.
func Package(name io.WriterTo, dirs ...io.WriterTo) PackageDeclr {
	return PackageDeclr{
		Name: name,
		Body: dirs,
	}
}

// CustomReturns returns a new instance of a CustomReturnDeclr.
func CustomReturns(returns ...io.WriterTo) CustomReturnDeclr {
	return CustomReturnDeclr{
		Returns: returns,
	}
}

// Returns returns a new instance of a ReturnDeclr.
func Returns(returns ...TypeDeclr) ReturnDeclr {
	return ReturnDeclr{
		Returns: returns,
	}
}

// Constructor returns a new instance of a ConstructorDeclr.
func Constructor(args ...VariableTypeDeclr) ConstructorDeclr {
	return ConstructorDeclr{
		Arguments: args,
	}
}

// Interface returns a new instance of a StructDeclr to generate a go struct.
func Interface(name NameDeclr, comments io.WriterTo, annotations io.WriterTo, fields ...io.WriterTo) StructDeclr {
	if annotations == nil {
		annotations = bytes.NewBuffer(nil)
	}

	if comments == nil {
		comments = bytes.NewBuffer(nil)
	}

	return StructDeclr{
		Type:        Type("interface"),
		Name:        name,
		Comments:    comments,
		Annotations: annotations,
		Fields:      fields,
	}
}

// Struct returns a new instance of a StructDeclr to generate a go struct.
func Struct(name NameDeclr, comments io.WriterTo, annotations io.WriterTo, fields ...io.WriterTo) StructDeclr {
	if annotations == nil {
		annotations = bytes.NewBuffer(nil)
	}

	if comments == nil {
		comments = bytes.NewBuffer(nil)
	}

	return StructDeclr{
		Type:        Type("struct"),
		Name:        name,
		Comments:    comments,
		Annotations: annotations,
		Fields:      fields,
	}
}

// Annotations returns a slice instance of io.WriterTo.
func Annotations(names ...string) io.WriterTo {
	var decls WritersTo

	for _, name := range names {
		decls = append(decls, Annotation(name))
	}

	return NewlineMapper.Map(decls...)
}

// Annotation returns a new instance of a AnnotationDeclr.
func Annotation(name string) AnnotationDeclr {
	return AnnotationDeclr{
		Value: name,
	}
}

// Tag returns a new instance of a TagDeclr.
func Tag(format string, name string) TagDeclr {
	return TagDeclr{
		Format: format,
		Name:   name,
	}
}

// Field returns a new instance of a StructTypeDeclr.
func Field(name NameDeclr, ntype TypeDeclr, tags ...io.WriterTo) StructTypeDeclr {
	return StructTypeDeclr{
		Name: name,
		Type: ntype,
		Tags: tags,
	}
}

// FunctionType returns a new instance of a FunctionTypeDeclr.
func FunctionType(name NameDeclr, constr ConstructorDeclr, returns io.WriterTo) FunctionTypeDeclr {
	return FunctionTypeDeclr{
		Name:        name,
		Constructor: constr,
		Returns:     returns,
	}
}

// Function returns a new instance of a FunctionDeclr.
func Function(name NameDeclr, constr ConstructorDeclr, returns io.WriterTo, body ...io.WriterTo) FunctionDeclr {
	return FunctionDeclr{
		Name:        name,
		Constructor: constr,
		Returns:     returns,
		Body:        body,
	}
}

// Source returns a new instance of a SourceDeclr.
func Source(tml *template.Template, binding interface{}) SourceDeclr {
	return SourceDeclr{
		Template: tml,
		Binding:  binding,
	}
}

// SourceText returns a new instance of a TextDeclr.
func SourceText(tml string, binding interface{}) TextDeclr {
	return TextDeclr{
		Template: tml,
		Binding:  binding,
	}
}

// PrefixByte returns a new instance of a SingleByteBlockDeclr.
func PrefixByte(start []byte) SingleByteBlockDeclr {
	return SingleByteBlockDeclr{
		Block: start,
	}
}

// SuffixByte returns a new instance of a SingleByteBlockDeclr.
func SuffixByte(end []byte) SingleByteBlockDeclr {
	return SingleByteBlockDeclr{
		Block: end,
	}
}

// Prefix returns a new instance of a SingleBlockDeclr.
func Prefix(start rune) SingleBlockDeclr {
	return SingleBlockDeclr{
		Rune: start,
	}
}

// Suffix returns a new instance of a SingleBlockDeclr.
func Suffix(end rune) SingleBlockDeclr {
	return SingleBlockDeclr{
		Rune: end,
	}
}

// MultiCommentary returns a new instance of a MultiCommentDeclr.
func MultiCommentary(mainblock io.WriterTo, elems ...io.WriterTo) MultiCommentDeclr {
	return MultiCommentDeclr{
		MainBlock: mainblock,
		Blocks:    WritersTo(elems),
	}
}

// Commentary returns a new instance of a CommentDeclr.
func Commentary(mainblock io.WriterTo, elems ...io.WriterTo) CommentDeclr {
	return CommentDeclr{
		MainBlock: mainblock,
		Blocks:    WritersTo(elems),
	}
}

// Block returns a new instance of a BlockDeclr with no prefix and suffix.
func Block(elems ...io.WriterTo) BlockDeclr {
	return BlockDeclr{
		Block: WritersTo(elems),
	}
}

// WrapBlock returns a new instance of a BlockDeclr.
func WrapBlock(begin, end rune, elems ...io.WriterTo) BlockDeclr {
	return BlockDeclr{
		RuneBegin: begin,
		RuneEnd:   end,
		Block:     WritersTo(elems),
	}
}

// ByteWrapBlock returns a new instance of a ByteBlockDeclr.
func ByteWrapBlock(begin, end []byte, elems ...io.WriterTo) ByteBlockDeclr {
	return ByteBlockDeclr{
		BlockBegin: begin,
		BlockEnd:   end,
		Block:      WritersTo(elems),
	}
}

// Switch returns a new instance of a SwitchDeclr.
func Switch(condition io.WriterTo, def DefaultCaseDeclr, cases ...CaseDeclr) SwitchDeclr {
	if def.Behaviour == nil {
		def.Behaviour = bytes.NewBuffer(nil)
	}

	return SwitchDeclr{
		Default:   def,
		Cases:     cases,
		Condition: condition,
	}
}

// DefaultCase returns a new instance of a DefaultCaseDeclr.
func DefaultCase(action io.WriterTo) DefaultCaseDeclr {
	return DefaultCaseDeclr{
		Behaviour: action,
	}
}

// Case returns a new instance of a CaseDeclr.
func Case(condition, action io.WriterTo) CaseDeclr {
	return CaseDeclr{
		Condition: condition,
		Behaviour: action,
	}
}

// Condition returns a new instance of a ConditionDeclr.
func Condition(pre VariableNameDeclr, op OperatorDeclr, post VariableNameDeclr) ConditionDeclr {
	return ConditionDeclr{
		PreVar:   pre,
		PostVar:  post,
		Operator: op,
	}
}

// Var returns a new instance of a VariableShortAssignmentDeclr.
func Var(name NameDeclr, value io.WriterTo) VariableShortAssignmentDeclr {
	return VariableShortAssignmentDeclr{
		Name:  name,
		Value: value,
	}
}

// AssignVar returns a new instance of a VariableAssignmentDeclr.
func AssignVar(name NameDeclr, value io.WriterTo) VariableAssignmentDeclr {
	return VariableAssignmentDeclr{
		Name:  name,
		Value: value,
	}
}

// VarName returns a new instance of a VariableNameDeclr.
func VarName(name NameDeclr) VariableNameDeclr {
	return VariableNameDeclr{
		Name: name,
	}
}

// VarType returns a new instance of a VariableTypeDeclr.
func VarType(name NameDeclr, ntype TypeDeclr) VariableTypeDeclr {
	return VariableTypeDeclr{
		Name: name,
		Type: ntype,
	}
}

// Ops returns a new instance of a OperatorDeclr.
func Ops(ty string) OperatorDeclr {
	return OperatorDeclr{
		Operation: ty,
	}
}

// SliceType returns a new instance of a SliceTypeDeclr.
func SliceType(ty string) SliceTypeDeclr {
	return SliceTypeDeclr{
		Type: Type(ty),
	}
}

// Slice returns a new instance of a SliceDeclr.
func Slice(typeName string, elems ...io.WriterTo) SliceDeclr {
	return SliceDeclr{
		Type:   Type(typeName),
		Values: elems,
	}
}

// If returns a new instance of a IfDeclr.
func If(condition, action io.WriterTo) IfDeclr {
	return IfDeclr{
		Action:    action,
		Condition: condition,
	}
}
