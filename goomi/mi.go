package goomi

type MIFlag uint32

const (
	// CIM meta types (or qualifier scopes)
	MI_FLAG_CLASS       MIFlag = 1 << 0
	MI_FLAG_METHOD      MIFlag = 1 << 1
	MI_FLAG_PROPERTY    MIFlag = 1 << 2
	MI_FLAG_PARAMETER   MIFlag = 1 << 3
	MI_FLAG_ASSOCIATION MIFlag = 1 << 4
	MI_FLAG_INDICATION  MIFlag = 1 << 5
	MI_FLAG_REFERENCE   MIFlag = 1 << 6
	MI_FLAG_ANY         MIFlag = 1 | 2 | 4 | 8 | 16 | 32 | 64

	// Qualifier flavors
	MI_FLAG_ENABLEOVERRIDE  MIFlag = 1 << 7
	MI_FLAG_DISABLEOVERRIDE MIFlag = 1 << 8
	MI_FLAG_RESTRICTED      MIFlag = 1 << 9
	MI_FLAG_TOSUBCLASS      MIFlag = 1 << 10
	MI_FLAG_TRANSLATABLE    MIFlag = 1 << 11

	// Select boolean qualifier
	MI_FLAG_KEY              MIFlag = 1 << 12
	MI_FLAG_IN               MIFlag = 1 << 13
	MI_FLAG_OUT              MIFlag = 1 << 14
	MI_FLAG_REQUIRED         MIFlag = 1 << 15
	MI_FLAG_STATIC           MIFlag = 1 << 16
	MI_FLAG_ABSTRACT         MIFlag = 1 << 17
	MI_FLAG_TERMINAL         MIFlag = 1 << 18
	MI_FLAG_EXPENSIVEMI_Flag MIFlag = 1 << 19
	MI_FLAG_STREAM           MIFlag = 1 << 20
	MI_FLAG_READONLY         MIFlag = 1 << 21

	// ToInstance flavor: ignored
	MI_FLAG_TOINSTANCE MIFlag = 1 << 22

	// Special flags
	MI_FLAG_NOT_MODIFIED MIFlag = 1 << 25 // indicates that the property is not modified
	MI_FLAG_VERSION      MIFlag = 1<<26 | 1<<27 | 1<<28
	MI_FLAG_NULL         MIFlag = 1 << 29
	MI_FLAG_BORROW       MIFlag = 1 << 30
	MI_FLAG_ADOPT        MIFlag = 1 << 31
)

type MIType int

const (
	MI_BOOLEAN MIType = iota
	MI_UINT8
	MI_SINT8
	MI_UINT16
	MI_SINT16
	MI_UINT32
	MI_SINT32
	MI_UINT64
	MI_SINT64
	MI_REAL32
	MI_REAL64
	MI_CHAR16
	MI_DATETIME
	MI_STRING
	MI_REFERENCE
	MI_INSTANCE
	MI_BOOLEANA
	MI_UINT8A
	MI_SINT8A
	MI_UINT16A
	MI_SINT16A
	MI_UINT32A
	MI_SINT32A
	MI_UINT64A
	MI_SINT64A
	MI_REAL32A
	MI_REAL64A
	MI_CHAR16A
	MI_DATETIMEA
	MI_STRINGA
	MI_REFERENCEA
	MI_INSTANCEA
	MI_ARRAY = MIType(16)
)

type MIQualifierDecl struct {
	// Name of this qualifier
	Name string

	// Type of this qualifier
	Type MIType

	// Qualifier scope
	Scope MIFlag

	// Qualifier flavor
	Flavor MIFlag

	// Array subscript (for arrays only)
	Subscript uint32

	// Value any
}

type MIQualifier struct {
	// Qualifier name
	Name string

	// Qualifier type
	Type MIType

	// Qualifier flavor
	Flavor MIFlag

	// value
	Value any
}

type MIQualifiers []*MIQualifier

type MIFeatureDecl struct {
	// Flags
	Flags MIFlag

	// Hash code: (name[0] << 16) | (name[len-1] << 8) | len
	Code uint32

	// Name of this feature
	Name string

	// Qualifiers
	Qualifiers    MIQualifiers
	NumQualifiers int
}

type MIParameterDecl struct {
	// Fields inherited from MI_FeatureDecl
	Flags         MIFlag
	Code          uint32
	Name          string
	Qualifiers    MIQualifiers
	NumQualifiers uint32

	// Type of this field
	Type MIType

	// Name of reference class
	ClassName string

	// Array subscript
	Subscript uint32

	// Offset of this field within the structure
	Offset uint32
}

type MIPropertyDecl struct {
	// Fields inherited from MI_FeatureDecl
	Flags         MIFlag
	Code          uint32
	Name          string
	Qualifiers    MIQualifiers
	NumQualifiers uint32

	// Fields inherited from MI_ParameterDecl
	Type      MIType
	ClassName string
	Subscript uint32
	Offset    uint32

	// Ancestor class that first defined a property with this name
	Origin string

	// Ancestor class that last defined a property with this name
	Propagator string

	// Value of this property
	// Value MIConstVoidPtr
}

type MIMethodDecl struct {
	// Fields inherited from MI_FeatureDecl
	Flags         MIFlag
	Code          uint32
	Name          string
	Qualifiers    MIQualifiers
	NumQualifiers uint32

	// Fields inherited from MI_ObjectDecl
	Parameters    []*MIParameterDecl
	NumParameters uint32
	Size          uint32

	// PostResult type of this method
	ReturnType uint32

	// Ancestor class that first defined a property with this name
	Origin string

	// Ancestor class that last defined a property with this name
	Propagator string

	// Pointer to schema this class belongs to
	Schema *MISchemaDecl

	// Pointer to extrinsic method
	// Function MIMethodDeclInvoke
}

type MIClassDecl struct {
	// Fields inherited from MI_FeatureDecl
	Flags         MIFlag
	Code          uint32
	Name          string
	Qualifiers    MIQualifiers
	NumQualifiers uint32

	// Fields inherited from MI_ObjectDecl
	Properties    []*MIPropertyDecl
	NumProperties uint32
	Size          uint32

	// Name of superclass
	SuperClass string

	// Superclass declaration
	SuperClassDecl *MIClassDecl

	// The methods of this class
	Methods    []*MIMethodDecl
	NumMethods uint32

	// Pointer to schema this class belongs to
	Schema *MISchemaDecl

	// Provider functions
	// ProviderFT *MIProviderFT

	// Owning MI_Class object, if any. NULL if static classDecl, -1 is from a dynamic instance
	// OwningClass *MIClass
}

type MISchemaDecl struct {
	// Qualifier declarations
	QualifierDecls    []*MIQualifierDecl
	NumQualifierDecls uint32

	// Class declarations
	ClassDecls    []*MIClassDecl
	NumClassDecls uint32
}
