package cgoomi

// #cgo CFLAGS: -I../../omi/Unix/common/
// #include <MI.h>
import "C"
import (
	"github.com/jclab-joseph/mof-to-java-converter/goomi"
	"log"
	"unsafe"
)

func ConvertValue(valueType goomi.MIType, pointer unsafe.Pointer) any {
	switch valueType {
	case goomi.MI_BOOLEAN:
		return *(*bool)(pointer)
	case goomi.MI_UINT8:
		return *(*uint8)(pointer)
	case goomi.MI_SINT8:
		return *(*int8)(pointer)
	case goomi.MI_UINT16:
		return *(*uint16)(pointer)
	case goomi.MI_SINT16:
		return *(*int16)(pointer)
	case goomi.MI_UINT32:
		return *(*uint32)(pointer)
	case goomi.MI_SINT32:
		return *(*int32)(pointer)
	case goomi.MI_UINT64:
		return *(*uint64)(pointer)
	case goomi.MI_SINT64:
		return *(*int64)(pointer)
	case goomi.MI_STRING:
		return C.GoString((*C.char)(pointer))
	//case goomi.MI_REAL32:
	//case goomi.MI_REAL64:
	//case goomi.MI_CHAR16:
	//	return *(*char16)(pointer)
	//case goomi.MI_DATETIME:
	//	return *(*datetime)(pointer)
	//case goomi.MI_REFERENCE:
	//	return *(*reference)(pointer)
	//case goomi.MI_INSTANCE:
	//	return *(*instance)(pointer)
	//case goomi.MI_BOOLEANA:
	//	return *(*booleana)(pointer)
	//case goomi.MI_UINT8A:
	//	return *(*uint8a)(pointer)
	//case goomi.MI_SINT8A:
	//	return *(*sint8a)(pointer)
	//case goomi.MI_UINT16A:
	//	return *(*uint16a)(pointer)
	//case goomi.MI_SINT16A:
	//	return *(*sint16a)(pointer)
	//case goomi.MI_UINT32A:
	//	return *(*uint32a)(pointer)
	//case goomi.MI_SINT32A:
	//	return *(*sint32a)(pointer)
	//case goomi.MI_UINT64A:
	//	return *(*uint64a)(pointer)
	//case goomi.MI_SINT64A:
	//	return *(*sint64a)(pointer)
	//case goomi.MI_REAL32A:
	//	return *(*real32a)(pointer)
	//case goomi.MI_REAL64A:
	//	return *(*real64a)(pointer)
	//case goomi.MI_CHAR16A:
	//	return *(*char16a)(pointer)
	//case goomi.MI_DATETIMEA:
	//	return *(*datetimea)(pointer)
	case goomi.MI_STRINGA:
		stringa := (*C.MI_StringA)(pointer)
		dataList := unsafe.Slice(stringa.data, stringa.size)
		var results []string
		for _, s := range dataList {
			results = append(results, C.GoString(s))
		}
		return results
		//case goomi.MI_REFERENCEA:
		//	return *(*referencea)(pointer)
		//case goomi.MI_INSTANCEA:
		//	return *(*instancea)(pointer)
		//case goomi.MI_ARRAY:
		//	return *(*array)(pointer)
	default:
		log.Fatalln("ERR : ", valueType)

	}
	return nil
}

func ConvertQualifier(cqualifier *C.MI_Qualifier) *goomi.MIQualifier {
	qualifier := &goomi.MIQualifier{
		Name:   C.GoString(cqualifier.name),
		Type:   goomi.MIType(cqualifier._type),
		Flavor: goomi.MIFlag(cqualifier.flavor),
	}
	qualifier.Value = ConvertValue(qualifier.Type, cqualifier.value)

	return qualifier
}

func ConvertQualifiers(cqualifiersPtr unsafe.Pointer, count uint32) []*goomi.MIQualifier {
	var resultList []*goomi.MIQualifier
	sliceList := unsafe.Slice((**C.MI_Qualifier)(cqualifiersPtr), count)
	for i := 0; i < len(sliceList); i++ {
		item := ConvertQualifier(sliceList[i])
		resultList = append(resultList, item)
	}
	return resultList
}

func ConvertPropertyDecl(cproperty *C.MI_PropertyDecl) *goomi.MIPropertyDecl {
	property := &goomi.MIPropertyDecl{
		Flags:         goomi.MIFlag(uint32(cproperty.flags)),
		Code:          uint32(cproperty.code),
		Name:          C.GoString(cproperty.name),
		Qualifiers:    ConvertQualifiers(unsafe.Pointer(cproperty.qualifiers), uint32(cproperty.numQualifiers)),
		NumQualifiers: uint32(cproperty.numQualifiers),

		Type:      goomi.MIType(cproperty._type),
		ClassName: C.GoString(cproperty.className),
		Subscript: uint32(cproperty.subscript),
		Offset:    uint32(cproperty.offset),

		Origin: C.GoString(cproperty.origin),

		Propagator: C.GoString(cproperty.propagator),
	}

	return property
}

func ConvertPropertyDecls(cpropertiesPtr unsafe.Pointer, count uint32) []*goomi.MIPropertyDecl {
	var resultList []*goomi.MIPropertyDecl
	sliceList := unsafe.Slice((**C.MI_PropertyDecl)(cpropertiesPtr), count)
	for i := 0; i < len(sliceList); i++ {
		resultList = append(resultList, ConvertPropertyDecl(sliceList[i]))
	}
	return resultList
}

func ConvertParameterDecl(cparameter *C.MI_ParameterDecl) *goomi.MIParameterDecl {
	parameter := &goomi.MIParameterDecl{
		Flags:         goomi.MIFlag(uint32(cparameter.flags)),
		Code:          uint32(cparameter.code),
		Name:          C.GoString(cparameter.name),
		Qualifiers:    ConvertQualifiers(unsafe.Pointer(cparameter.qualifiers), uint32(cparameter.numQualifiers)),
		NumQualifiers: uint32(cparameter.numQualifiers),

		Type:      goomi.MIType(uint32(cparameter._type)),
		ClassName: C.GoString(cparameter.className),
		Subscript: uint32(cparameter.subscript),
		Offset:    uint32(cparameter.offset),
	}
	return parameter
}

func ConvertParameterDecls(cparametersPtr unsafe.Pointer, count uint32) []*goomi.MIParameterDecl {
	var resultList []*goomi.MIParameterDecl
	sliceList := unsafe.Slice((**C.MI_ParameterDecl)(cparametersPtr), count)
	for i := 0; i < len(sliceList); i++ {
		resultList = append(resultList, ConvertParameterDecl(sliceList[i]))
	}
	return resultList
}

func ConvertMethodDecl(cmethod *C.MI_MethodDecl) *goomi.MIMethodDecl {
	method := &goomi.MIMethodDecl{
		Flags:         goomi.MIFlag(uint32(cmethod.flags)),
		Code:          uint32(cmethod.code),
		Name:          C.GoString(cmethod.name),
		Qualifiers:    ConvertQualifiers(unsafe.Pointer(cmethod.qualifiers), uint32(cmethod.numQualifiers)),
		NumQualifiers: uint32(cmethod.numQualifiers),

		Parameters:    ConvertParameterDecls(unsafe.Pointer(cmethod.parameters), uint32(cmethod.numParameters)),
		NumParameters: uint32(cmethod.numParameters),
		Size:          uint32(cmethod.size),

		// PostResult type of this method
		ReturnType: uint32(cmethod.returnType),

		// Ancestor class that first defined a property with this name
		Origin: C.GoString(cmethod.origin),

		// Ancestor class that last defined a property with this name
		Propagator: C.GoString(cmethod.propagator),

		// Pointer to schema this class belongs to
		// Schema *MISchemaDecl

		// Pointer to extrinsic method
		// Function MIMethodDeclInvoke
	}
	return method
}

func ConvertMethodDecls(cmethodsPtr unsafe.Pointer, count uint32) []*goomi.MIMethodDecl {
	var methods []*goomi.MIMethodDecl
	cmethodsList := unsafe.Slice((**C.MI_MethodDecl)(cmethodsPtr), count)
	for i := 0; i < len(cmethodsList); i++ {
		methods = append(methods, ConvertMethodDecl(cmethodsList[i]))
	}
	return methods
}

func ConvertClassDecl(cdecl *C.MI_ClassDecl) *goomi.MIClassDecl {
	classDecl := &goomi.MIClassDecl{
		Flags: goomi.MIFlag(cdecl.flags),
		Code:  uint32(cdecl.code),
		Name:  C.GoString(cdecl.name),
		// Qualifiers
		NumQualifiers: uint32(cdecl.numQualifiers),
		// Properties
		NumProperties: uint32(cdecl.numProperties),
		Size:          uint32(cdecl.size),
		SuperClass:    C.GoString(cdecl.superClass),
		// SuperClassDecl
		// Methods
		NumMethods: uint32(cdecl.numMethods),
		// Schema
	}
	classDecl.Qualifiers = ConvertQualifiers(unsafe.Pointer(cdecl.qualifiers), classDecl.NumQualifiers)
	classDecl.Properties = ConvertPropertyDecls(unsafe.Pointer(cdecl.properties), classDecl.NumProperties)
	classDecl.Methods = ConvertMethodDecls(unsafe.Pointer(cdecl.methods), classDecl.NumMethods)

	return classDecl
}

func ConvertClassDeclFromPtr(cdeclPtr unsafe.Pointer) *goomi.MIClassDecl {
	return ConvertClassDecl((*C.MI_ClassDecl)(cdeclPtr))
}
