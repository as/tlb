package main

// MIDL Base Types
// MIDL Arrays
// MIDL Type Definitions
// MIDL Enumerated Types
// MIDL Structures
// MIDL Unions
// MIDL Binding Handles

type WChart interface{}

type Numeric struct {
	Neg, Pos WChart
	Dot      WChart
	Sep      WChart
	Currency Currency
}

type Currency struct {
	Local  WChart
	Local2 WChart
	Dot    WChart
	Sep    WChart
}

/*
const(
Invalid = iota
Bool          // "boolean"       "8 bits. Not compatible with oleautomation interfaces; use VARIANT_BOOL"
Byte          // "byte"           "8 bits."
Char          // "char"           "8 bits."
Double        // "double"         "64-bit floating point number."
Error         // "error_status_t" "32-bit unsigned integer for returning status values for error"
Float         // "float"         "32-bit floating point number."
Handle        // "handle_t"      "Primitive handle type for binding."
Hyper         // "hyper"         "64-bit integer."
Int           // "int"           "32-bit integer. On 16-bit platforms, cannot appear in remote functions"
Int8          // "__int8"        "8-bit integer. Equivalent to small."
Int16         // "__int16"       "16-bit integer. Equivalent to short."
Int32         // "__int32"       "32-bit integer. Equivalent to long."
Int3264       // "__int3264"     "An integer that is 32-bit on 32-bit platforms, and is 64-bit on 64-bit"
Int64         // "__int64"       "64-bit integer. Equivalent to hyper."
Long          // "long"          "32-bit integer."
Short         // "short"         "16-bt integer."
Small         // "small"         "8-bit integer."
Void          // "void"          "Indicates that the procedure does not return a value."
VoidStar      // "void *"        "32-bit pointer for context handles only."
//WChart        // "wchar_t"       "16-bit predefined type for wide characters."
)
*/

type StructField struct {
	Name      string
	Type      Type    // field type
	Offset    uintptr // offset within struct, in bytes
	Index     []int   // index sequence for Type.FieldByIndex
	Anonymous bool    // is an embedded field
}
type Method struct {
	Name  string
	Type  Type  // method type
	Func  Value // func with receiver as first argument
	Index int   // index for Type.Method
}
type Value struct {
}

/*
func (t Type) Origin() string{return ""}
func (t Type) Align() int { return 4}
func (t Type) FieldAlign() int { return 4}
func (t Type) Method(int) Method { return &Method{}}
func (t Type) MethodByName(string) (Method, bool) {}
func (t Type) NumMethod(string) (Method, bool) {}
func (t Type) Name() string {return ""}
func (t Type) PkgPath() string {return ""}
func (t Type) Size() uintptr {return 4}
func (t Type) Kind() Kind {}
func (t Type) Implements(u Type) bool {return false}
func (t Type) AssignableTo(u Type) bool {return false}
func (t Type) ConvertibleTo(u Type) bool {return false}
func (t Type) Comparable(u Type) bool {return false}
func (t Type) Bits() int{return 4}
func (t Type) ChanDir() ChanDir{}
func (t Type) IsVariadic() bool{return false}
func (t Type) Elem() Type{}
func (t Type) Field(i int) StructField{}
func (t Type) FieldByIndex(index []int) StructField{}
func (t Type) FieldByName(name string) (StructField, bool){ return &StructField{}, false}
func (t Type) FieldByNameFunc(match func(string) bool) (StructField, bool){return &StructField{}, false}
func (t Type) In(i int) Type{return &Type{}}
func (t Type) Key() Type{return &Type{}}
func (t Type) Len() int{return 4}
func (t Type) NumField() int{return 4}
func (t Type) NumIn() int{return 4}
func (t Type) NumOut() int{return 4}
func (t Type) Out() int{return 4}
*/
