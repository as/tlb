package main

type Interface interface {
}
type Record interface {
}
type Error interface {
}
type String interface {
}
type Ptr interface {
}
type MemberID interface {
}
type GID interface {
}
type TypeKind interface {
}
type Lang interface {
}
type ExceptionInfo interface {
}
type Bstr interface {
}
type TypeInfo interface {
}

// 0000002F-0000-0000-C000-000000000046
type RecordInfo interface {
	Init() (Record, Error)
	Clear(Record) (Record, Error)
	Copy(Record) (Record, Error)
	ID() (Record, Error)
	Name() (String, Error)
	Size() (uint32, Error)
	TypeInfo() (Typer, Error)
	Field(Ptr, String) (Variant, Array, Error)
	FieldRemove(Ptr, String) (Variant, Array, Error)
	PutField(uint32, Ptr) (String, Variant, Error)
	PutFieldRemove(uint32, Ptr) (String, Variant, Error)
	Fields(uint32) (String, Error)
	Matches(Type) bool
	Create() Ptr
	CreateCopy(Record) (Record, Error)
	Destroy(Record) Error
}

// 0000002F-0000-0000-C000-000000000046
type CreateTypeInfo interface {
	Init() (Record, Error)
	Clear(Record) (Record, Error)
	Copy(Record) (Record, Error)
	ID() (Record, Error)
	Name() (String, Error)
	Size() (uint32, Error)
	TypeInfo() (Typer, Error)
	Field(Ptr, String) (Variant, Array, Error)
	FieldRemove(Ptr, String) (Variant, Array, Error)
	PutField(uint32, Ptr) (String, Variant, Error)
	PutFieldRemove(uint32, Ptr) (String, Variant, Error)
	Fields(uint32) (String, Error)
	Matches(Type) bool
	Create() Ptr
	CreateCopy(Record) (Record, Error)
	Destroy(Record) Error
}

//00020400-0000-0000-C000-000000000046
type Dispatcher interface {
}

// No wonder everyone hates COM
// 0002040E-0000-0000-C000-000000000046
type CreateTypeInfo2ElectricBoogaloo interface {
	DeleteFuncDesc(uint) Error
	DeleteFuncDescByMemId(MemberID, InvokeKind) Error
	DeleteVarDesc(uint) Error
	DeleteVarDescByMemId(MemberID) Error
	DeleteImplType(uint) Error
	SetCustData(GID, Variant) Error
	SetFuncCustData(uint, GID, Variant) Error
	SetParamCustData(uint, uint, GID, Variant) Error
	SetVarCustData(uint, GID, Variant) Error
	SetImplTypeCustData(uint, GID, Variant) Error
	SetHelpStringContext(uint) Error
	SetFuncHelpStringContext(uint, uint) Error
	SetVarHelpStringContext(uint, uint) Error
	SetName(String) Error
}

// uuid(00020406-0000-0000-C000-000000000046)
type CreateTypeLib interface {
	CreateTypeInfo(String, TypeKind, CreateTypeInfo) Error
	SetName(String) Error
	SetVersion(uint16, uint16) Error
	SetGuid(GID) Error
	SetDocString(String) Error
	SetHelpFileName(String) Error
	SetHelpContext(uint32) Error
	SetLcid(Lang) Error
	SetLibFlags(uint) Error
	SaveAllChanges() Error
}

// uuid(0002040F-0000-0000-C000-000000000046),
type CreateTypeLib2 interface {
	CreateTypeLib
	DeleteTypeInfo(String) Error
	SetCustData(GID, Variant) Error
	SetHelpStringContext(uint32) Error
	SetHelpStringDll(String) Error
}

// 3127ca40-446e-11ce-8135-00aa004bb851),
type ErrorLogger interface {
	Interface
	AddError(String, ExceptionInfo) Error
}

//  uuid(1CF2B120-547D-101B-8E65-08002B2BD119),
type IErrorInfo interface {
	GID() (GID, Error)
	Source() (Bstr, Error)
	Description() (Bstr, Error)
	HelpFile() (Bstr, Error)
	HelpContext() (uint32, Error)
}

// uuid(55272a00-42cb-11ce-8135-00aa004bb851),
type PropertyBag interface {
	Read(String, Variant) (Variant, ErrorLogger, Error)
	RemoteRead(String, Variant) (ErrorLogger, uint32, Interface, Error)
	Write(String, Variant) Error
}

// uuid(00020410-0000-0000-C000-000000000046),
type ChangeKind int

const (
	ChAddMember ChangeKind = iota
	ChDeleteMember
	ChSetNames
	ChSetDocumentation
	ChGenerate
	ChInvalidate
	ChChangeFailed
	ChMax
)

type TypeChangeEvent interface {
	RequestTypeChange(ChangeKind, TypeInfo, String) (Error, int)
	AfterTypeChange(ChangeKind, TypeInfo, String) Error
}
