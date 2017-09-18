package main

//
type Typ int
type BasicKind int

type (
	Void    struct{}
	Enum    struct{}
	Struct  struct{}
	Union   struct{}
	Alias   struct{}
	Module  struct{}
	CoClass struct{}
	Pointer struct{}
)

type Typer interface {
	ReadBinary() error
	WriteBinary() error
	Underlying() Type
	Name() string
}

type Requester interface {
	Interface
	Init()
	Clear(r Record)
}
type Slicer interface {
	SliceOf() []Typer
}

// ComType enumerates the 8 highest-order types in COM. The names
// of familiar province to a programmer are given in comment form
type ComType int

const (
	ComEnum      ComType = iota
	ComRecord            // A struct
	ComModule            // Runtime library (for DLL Hell to work)
	ComInterface         // A list of methods
	ComDispatch          // Describes a function call in a clumsy way
	ComCoclass           // A little "friend" that "helps" you by marshaling
	// and unmashaling dispatch requests for a class
	ComAlias // Typedef for people that need a new word for it
	ComUnion // A struct crushed into fewer bits
	ComMax
)

//TODO
type Hreftype interface{}

// Type is a general representation of a Com datatype. The associated
// TypeDesc structure identifies and describes the true underlying datatype.
//
//!wire9 typeinfo typ[4] u1[16]	nvar[2] nfunc[2] u2[16] gidat[4] flag[4] nameat[4] ver[4] he1[4] he2[4] cdat[4] nimport[2] nvirt[2] n[4] aux[4] aux2[4] pad[8]
type Type struct {
	gid                GID
	lcid               string
	ctor               MemberID
	dtor               MemberID
	scheme             String
	size               uint32  // of instance
	kind               ComType // typekind (shorter list)
	nfunc, nvar, nimps uint16
	vftsize, align     uint16
	flag               uint16
	TypeDesc           // The other bs: Enum-Max
	//idldesc
}

// TypeIs enumerates the type's basic properties
type TypeIs int

const (
	TypeAppobject TypeIs = 1<<iota + 1
	TypeCancreate
	TypeLicensed
	TypePredeclid
	TypeHidden
	TypeControl
	TypeDual
	TypeNonextensible
	TypeOleautomation
	TypeRestricted
	TypeSliceable
	TypeReplaceable
	TypeDispatchable
	TypeReversebind
	TypeProxy
)

// TypeDesc describes a general datatype
type TypeDesc struct {
	down      *TypeDesc  // if pointer
	safearray *TypeDesc  // if array
	array     *ArrayDesc // if array
	cust      *Hreftype  // if custom
	kind      VarKind    // The VarType/VtXXX
}

// TypeDesc describes an array type
type ArrayDesc struct {
	elem  TypeDesc // if pointer
	nelem uint16
	// may be more
}

type ParamFlag int

const (
	ParNone ParamFlag = 1 << iota
	ParIn
	ParOut
	ParLcid
	ParRetval
	ParOpt
	ParDefault
	ParCustom
)

type ParamDescEx struct {
	n          uint32
	defaultval Variant
}
type ParamDesc struct {
	// No
}

//wire idldesc r1[4] flag[2]
type IDLDesc struct {
	//idldesc
}
type IDLFlag int

const (
	IDLNone IDLFlag = 1 << iota
	IDLIn
	IDLOut
	IDLLcid
	IDLRetval
)

type ElemDesc struct {
	TypeDesc
	ParamDesc
}

type Dispatch struct {
	Arg       []Arg // named or anonymous
	nargs     int
	nanonargs int
}

func (d Dispatch) NumArgs() int {
	return d.nargs
}
func (d Dispatch) NumAnonArgs() int {
	return d.nargs
}

type Arg struct {
	name *String // null=unnamed
}

//bstr  n[2] data[n]
//wire9 exceptioninfo ret[2] r1[2] strings[3,[]bstr] context[4] r2[4] u1[4] scode[,SCode]
type ExeptionInfo struct {
	//exceptioninfo
}

//
type CallConv int

const (
	CallFast CallConv = iota
	CallC
	CallPascal
	CallStandard
	CallFpFast
	CallTrap
	_
	_
	CallMax
)

// FuncKind enumerates a function's compile time state
type FuncKind int

const (
	FuncVirtual FuncKind = iota
	FuncPure
	FuncNormal
	FuncStatic
	FuncDispatch
)

// FuncIs enumerates basic properties describing a function
type FuncIs int

const (
	FuncRestricted FuncIs = 1<<iota + 1
	FuncSource
	FuncBindable
	FuncRequestedit
	FuncDisplaybind
	FuncDefaultbind
	FuncHidden
	FuncUsesgetlasterror
	FuncDefaultcollelem
	FuncUidefault
	FuncNonbrowsable
	FuncReplaceable
	FuncImmediatebind
)

// FuncKind enumerates a function's invokation style
type InvokeKind int

const (
	InvFunc InvokeKind = iota<<iota + 1
	InvGet
	InvPut
	InvPutRef
)

// Changed VarKind to VarScope
// because using VarKind doesn't
// actually make any sense for these
// enums
//
// VarScope enumerates scope rules
type VarScope int

const (
	ScopeInstance VarScope = iota
	ScopeStatic
	ScopeConst
	ScopeDispatch
)

// ImpFlag enumerates a functions implementation kind
type ImpFlag int

const (
	ImpDefault = 1<<iota + 1
	ImpSource
	ImpRestrict
	ImpDefaultVtab
)

// FuncDesc describes a function and stores an offset
// to its slot in the vtable.
type FuncDesc struct {
	id MemberID

	// needs to be union
	//Scode Scode
	ElemDesc ElemDesc
	//

	kind              FuncKind
	inv               InvokeKind
	call              CallConv
	nparam, noptparam uint16
	virtualAt         uint16
	nscodes           uint16
	funcdesc          ElemDesc
	flags             uint16
}

type VarIs int

const (
	VarReadonly VarIs = 1<<iota + 1
	VarSource
	VarBindable
	VarRequestedit
	VarDisplaybind
	VarDefaultbind
	VarHidden
	VarRestricted
	VarDefaultcollelem
	VarUidefault
	VarNonbrowsable
	VarReplaceable
	VarImmediatebind
)

type VarDesc struct {
	MemberID   MemberID
	Schema     String
	instanceat uint32
	value      *Variant
	elem       ElemDesc
	flag       uint16
	kind       VarKind
}

// Weird marshalling TODO. The real name
// of this struct is CleanLocalStorage.
type CleanTLS struct {
	iface interface{}
	ptr   interface{}
	flag  uint32
}

type CustomItem struct {
	GID   GID
	value Variant
}

type CustomData struct {
	nelem uint32
	item  *CustomItem
}

//wire9 recinfo flag[4] n[4] irec[,int] data[n]
type RecInfo struct {
	//recinfo
}

//wire9 variant n[4] u1[4] kind[2] u1[6] datatype[,Typer]
type Variant struct {
	variant VarKind
}

// Unused
type Unused int

const (
	UnStructField Unused = iota
	UnTypedef
	UnCoClassField
	UnArgument
	UnProperty
	UnMethod
	UnModuleMember
	UnModule
	UnCoClass
	UnDispatchInterface
	UnInterface
	UnLibrary
)

type Can int

const (
	CanAggregatableCan = iota
	CanAppOpject
	CanBindable
	_
	CanControl
	_
	CanDefault
	CanDefaultBind
	CanDefaultCollElem
	CanDefaultValue
	CanDefaultVtable
	_
	CanDisplayBind
	CanDLLName
	CanDual
	CanEntry
	CanHelpContext
	CanHelpFile
	CanHelpString
	CanHelpStringDLL
	CanHidden
	CanID
	CanImmediateBind
	CanLCID
	_
	CanLicensed
	CanNonBrowsable
	CanNonCreateable
	CanNonExtensible
	CanOLEAutomation
	CanOptional
	CanPropGet
	CanPropPut
	CanPropPutref
	CanPublic
	CanReadOnly
	CanReplace
	CanRequestEdit
	CanRestricted
	CanRetval
	CanSource
	CanString
	CanUIDefault
	CanGetLastError
	CanUUID
	CanVarArg
	_
	CanCallAs
	CanFirstIs
	CanIgnore
	CanIIDIs
	CanIn
	CanLastIs
	CanLengthIs
	CanLocal
	CanMaxIs
	CanObject
	CanOut
	CanPointerDefault
	CanPtr
	CanRef
	CanSizeIs
	CanSwitch
	CanSwitchIs
	CanSwitchType
	CanTransmitAs
	CanUnique
	CanV1Enum
	CanWireMarshal

	// Unsure
	CanIdempotent
	CanBroadcast
	CanMaybe
	CanMessage

	// RPC only
	CanContextHandle
	CanHandle
	CanMsUnion
	CanPipe
	CanCallback
	CanEndpoint
	CanVersion
)

type CanInfo struct {
}

// Typ describes the contents of the typeinfo
// structure.
const (
	TypVoid Typ = iota
	TypBasic
	TypEnum
	TypStruct
	TypUnion
	TypUnion2
	TypAlias
	TypModule
	TypCoClass
	TypFunc
	TypInterface
	TypPointer
	TypArray
)

// An array represents either a C-style dynamic array or
// a COM "Safe Array". A "Safe Array" COM may be rebased
// to an arbitrary starting index. This implementation
// does not use the unsafe package to implement C arrays.
//
//wire9 safearray n[2] flag[2] min[4] nlocks[4] ptr[4] nelem[4] nelem[4] min[4] data[nelem,Slicer]
type (
	Array struct {
		elem TypeDesc
		n    uint16
	}
	SafeArray struct {
		//safearray
	}
)

const (
	RegDefault  VarScope = iota
	RegRegister          = 1
	RegNone              = 2
)

type LibPerm int

const (
	LibRestricted LibPerm = 1<<iota + 1
	LibControl
	LibHidden
	LibHasdiskimage
)

// VarKind enumerates possible types a varint may represent
// during runtime.
type VarKind int

const (
	VarEmpty VarKind = iota
	VarNull
	VarInt16
	VarInt32
	VarFloat32
	VarFloat64
	VarCurrency
	VarDate
	VarBstr
	VarDispatch
	VarError // 10
	VarBool
	VarVariant
	VarUnknown
	VarDecimal
	_ // 15
	VarInt8
	VarUint8
	VarUint16
	VarUint32
	VarInt64 // 20
	VarUint64
	VarInt
	VarUint
	VarVoid
	VarHresult
	VarVoidstar
	VarSafearray
	VarCarray
	VarUserdefined
	VarLpstr  // 30
	VarLpwstr // Try saying this out loud
	VarRecord
	VarIntptr
	VarUintptr

	VarFiletime = iota + 64
	VarBlob     // Super-awesome name
	VarStream
	VarStorage
	VarStreamed
	VarStored
	VarBlob2
	VarCf
	VarClassID // CLSID
	VarVersioned
	VarArray  = 0x1000
	VarScalar = 0x0fff
)

// A Basic represents a basic type.
type VarInfo int
type Basic struct {
	kind VarKind
	info VarInfo
	name string
}

const (
	IsInteger = iota
	IsUnsigned
	IsFloat
)

var BasicInfo = []*Basic{
	//	Invalid: {Invalid, 0, "invalid type"},
	VarInt8:    {VarInt8, IsInteger, "Int8"},
	VarInt16:   {VarInt16, IsInteger, "int16"},
	VarInt32:   {VarInt32, IsInteger, "Int32"},
	VarInt64:   {VarInt64, IsInteger, "Int64"},
	VarInt:     {VarInt, IsInteger, "Int"},
	VarFloat32: {VarFloat32, IsFloat, "Float32"},
	VarFloat64: {VarFloat64, IsFloat, "Float64"},
	VarError:   {VarError, IsInteger, "Error"},
	//VarHandle:  {Int, IsInteger, "Handle"},
}

// might not need these
//!wire9 arBstr n[4] data[,[]bstr]
//!wire9 arUnknown n[4] data[,[]unknown]
//!wire9 arDispatch n[4] data[,[]dispatch]
//!wire9 arVariant n[4] data[,[]variant]
//!wire9 arRecord n[4] data[,[]record]
//!wire9 arInterface n[4] data[,[]COM]
//!wire9 arInt8 n[4] data[,[]uint32]
//!wire9 arInt16 n[4] data[,[]iface]
//!wire9 arInt32 n[4] data[,[]iface]
//!wire9 arInt64 n[4] data[,[]iface]
//!wire9 arUint8 n[4] data[,[]iface]
//!wire9 arUint16 n[4] data[,[]iface]
//!wire9 arUint32 n[4] data[,[]iface]
//!wire9 arUint64 n[4] data[,[]iface]
