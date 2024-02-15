package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/jclab-joseph/mof-to-java-converter/cgoomi"
	"github.com/jclab-joseph/mof-to-java-converter/goomi"
	"github.com/jclab-joseph/mof-to-java-converter/pkg/javagen"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"unsafe"
)

// #cgo CFLAGS: -I../../omi/Unix/common/ -I../../omi/Unix/
// #cgo LDFLAGS: -L../../omi/Unix/output/lib/ -lmof -lpal
// #include <MI.h>
// #include <mof/mof.h>
// typedef const char cchar_t;
// typedef const wchar_t cwchar_t;
// int app_mofParserErrorCallback(cchar_t* msg, cwchar_t* wmsg, void*);
// int app_mofQualifierDeclCallback(MI_QualifierDecl*, void*);
// int app_mofClassDeclCallback(MI_ClassDecl*, void*);
// int app_mofInstanceDeclCallback(MI_InstanceDecl*, void*);
import "C"

var outdir string
var javaPackage string
var mapperFile *MapperFile

func main() {
	var paths string
	var includeDefaultSchema bool
	var srcdir string
	var nsMapperFile string

	cwd, _ := os.Getwd()

	flag.StringVar(&paths, "paths", "", "directory paths seperate by \":\"")
	flag.BoolVar(&includeDefaultSchema, "include-default-schema", true, "include cim_schema_2.32.0.mof")
	flag.StringVar(&srcdir, "srcdir", "", "")
	flag.StringVar(&outdir, "outdir", filepath.Join(cwd, "output"), "output directory")
	flag.StringVar(&javaPackage, "package", "org.example.mof", "java package")
	flag.StringVar(&nsMapperFile, "ns-mapper-file", "", "namespace mapper file ")
	flag.Parse()
	mofFiles := flag.Args()

	pathList := strings.Split(paths, ":")
	pathArrBuf := make2DArrayFromStrings(pathList)

	_ = mofFiles

	mapperFile = NewMapperFile()
	if nsMapperFile != "" {
		if err := mapperFile.LoadMapperFile(nsMapperFile); err != nil {
			log.Fatalln("mapper file load failed: ", err)
		}
	}

	parser := C.MOF_Parser_New((**C.char)(unsafe.Pointer(&pathArrBuf[0])), C.ulong(len(pathList)))
	if parser == nil {
		log.Fatalln("cannot create parser")
	}
	defer C.MOF_Parser_Delete(parser)

	C.MOF_Parser_SetErrorCallback(parser, (*[0]byte)(C.app_mofParserErrorCallback), nil)
	C.MOF_Parser_SetQualifierDeclCallback(parser, (*[0]byte)(C.app_mofQualifierDeclCallback), nil)
	C.MOF_Parser_SetClassDeclCallback(parser, (*[0]byte)(C.app_mofClassDeclCallback), nil)
	C.MOF_Parser_SetInstanceDeclCallback(parser, (*[0]byte)(C.app_mofInstanceDeclCallback), nil)

	if includeDefaultSchema {
		C.MOF_Parser_ParseString(parser, C.CString("#pragma include (\"cim_schema_2.32.0.mof\")"))
		C.MOF_Parser_ParseString(parser, C.CString("#pragma include (\"qualifiers.mof\")"))
		C.MOF_Parser_ParseString(parser, C.CString("#pragma include (\"qualifiers_optional.mof\")"))
		C.MOF_Parser_ParseString(parser, C.CString("Qualifier EOBase64 : boolean = false,\nScope(property, method, parameter);"))
	}

	packageInfoContent := fmt.Sprintf("@XmlSchema(\n        elementFormDefault = XmlNsForm.QUALIFIED\n)\npackage %s;\n\nimport jakarta.xml.bind.annotation.XmlNsForm;\nimport jakarta.xml.bind.annotation.XmlSchema;\n", javaPackage)
	if err := os.WriteFile(filepath.Join(outdir, "package-info.java"), []byte(packageInfoContent), 0644); err != nil {
		log.Printf("package-info.java write failed: %v", err)
	}

	for i, file := range mofFiles {
		if len(srcdir) > 0 {
			file = filepath.Join(srcdir, file)
		}
		log.Printf("Convert %s...", file)
		filePathRaw := []byte(file + "\x00")
		result := C.MOF_Parser_Parse(parser, (*C.char)(unsafe.Pointer(&filePathRaw[0])))
		_ = i
		if result != 0 {
			log.Printf("MOF_FILE[%s] Parse: %d", file, result)
		}
	}

	//C.MOF_Parser_Dump(parser, C.stdout)
}

func make2DArrayFromStrings(input []string) []byte {
	pointerSize := int(unsafe.Sizeof(uintptr(0)))

	offset := pointerSize * len(input)
	totalSize := offset

	for _, s := range input {
		totalSize += len(s) + 1
	}

	buf := make([]byte, totalSize)
	i := 0
	for _, s := range input {
		dst := unsafe.Pointer(&buf[offset])
		if pointerSize == 8 {
			binary.LittleEndian.PutUint64(buf[i:], uint64(uintptr(dst)))
		} else {
			binary.LittleEndian.PutUint32(buf[i:], uint32(uintptr(dst)))
		}
		i += pointerSize
		copy(buf[offset:], []byte(s))
		offset += len(s) + 1
	}

	return buf
}

//export app_mofParserErrorCallback
func app_mofParserErrorCallback(cmsg *C.cchar_t, cwmsg *C.cwchar_t, cdata *C.void) C.int {
	msg := C.GoString(cmsg)
	log.Println(msg)
	return 0
}

//export app_mofQualifierDeclCallback
func app_mofQualifierDeclCallback(decl *C.MI_QualifierDecl, cdata *C.void) C.int {
	//name := C.GoString(decl.name)
	//typ := uint32(decl._type)
	//scope := uint32(decl.scope)
	//flavor := uint32(decl.flavor)
	//subscript := uint32(decl.subscript)
	// value := uint32(decl.value)
	// log.Println("app_mofQualifierDeclCallback ", name, typ, scope, flavor, subscript)
	return 0
}

//export app_mofClassDeclCallback
func app_mofClassDeclCallback(cdecl *C.MI_ClassDecl, cdata *C.void) C.int {
	classDecl := cgoomi.ConvertClassDeclFromPtr(unsafe.Pointer(cdecl))

	GenerateJavaClassFromClassDecl(classDecl)

	// Qualifiers

	//typ := uint32(cdecl._type)
	//scope := uint32(cdecl.scope)
	//flavor := uint32(cdecl.flavor)
	//subscript := uint32(cdecl.subscript)
	// value := uint32(cdecl.value)
	//log.Println("app_mofClassDeclCallback ", classDecl.Name)
	return 0
}

//export app_mofInstanceDeclCallback
func app_mofInstanceDeclCallback(decl *C.MI_InstanceDecl, cdata *C.void) C.int {
	// log.Println("app_mofInstanceDeclCallback")
	return 0
}

// 상속 된 프로퍼티 제거

func newDtoClassGen(name string, superClass string) *javagen.ClassDecl {
	classGen := javagen.NewClassDecl(javagen.CLASS)
	classGen.PackageName = javaPackage
	classGen.Name = name
	classGen.SuperClass = superClass

	classGen.AddImport("jakarta.xml.bind.annotation.XmlType")
	classGen.AddImport("jakarta.xml.bind.annotation.XmlAccessorType")
	classGen.AddImport("jakarta.xml.bind.annotation.XmlAccessType")
	classGen.AddImport("jakarta.xml.bind.annotation.XmlElement")

	return classGen
}

func GenerateJavaClassFromClassDecl(decl *goomi.MIClassDecl) {
	classFileName := filepath.Join(outdir, decl.Name+".java")
	interfaceFileName := filepath.Join(outdir, decl.Name+"DataSource.java")

	classGen := newDtoClassGen(decl.Name, decl.SuperClass)
	abstractQualifier := decl.Qualifiers.FindByName("Abstract")
	if abstractQualifier != nil {
		classGen.IsAbstract = true
	}

	interfaceGen := javagen.NewClassDecl(javagen.INTERFACE)
	if decl.SuperClass != "" {
		interfaceGen.SuperClass = decl.SuperClass + "DataSource"
	}
	interfaceGen.PackageName = javaPackage
	interfaceGen.Name = decl.Name + "DataSource"

	interfaceGen.AddImport("jakarta.jws.WebService")
	interfaceGen.AddImport("jakarta.xml.bind.annotation.XmlSeeAlso")
	interfaceGen.AddImport("jakarta.jws.soap.SOAPBinding")
	interfaceGen.AddImport("jakarta.jws.soap.SOAPBinding.ParameterStyle")
	interfaceGen.AddImport("jakarta.jws.WebMethod")
	interfaceGen.AddImport("jakarta.xml.ws.Action")
	interfaceGen.AddImport("jakarta.jws.WebResult")
	interfaceGen.AddImport("jakarta.jws.WebParam")

	// INPUT/OUTPUT DTO
	interfaceGen.AddImport("jakarta.xml.bind.annotation.XmlType")
	interfaceGen.AddImport("jakarta.xml.bind.annotation.XmlAccessorType")
	interfaceGen.AddImport("jakarta.xml.bind.annotation.XmlAccessType")
	interfaceGen.AddImport("jakarta.xml.bind.annotation.XmlElement")

	xmlNs := mapperFile.Mappings[decl.Name]
	if xmlNs == "" {
		log.Printf("Cannot find XML Namespace about '%s'", decl.Name)
	}

	classGen.AddClassAnnotation("@XmlAccessorType(XmlAccessType.PROPERTY)")
	classGen.AddClassAnnotation(fmt.Sprintf("@XmlType(namespace = \"%s\")", xmlNs))

	if xmlNs != "" {
		classGen.AddBody(fmt.Sprintf("public static final String RESOURCE_URI = \"%s\";", xmlNs))
		classGen.AddBody("")
		interfaceGen.AddBody(fmt.Sprintf("public static final String RESOURCE_URI = \"%s\";", xmlNs))
		interfaceGen.AddBody("")
	}

	for _, property := range decl.Properties {
		if property.Origin != decl.Name {
			// Inherited property
			continue
		}

		description := property.Qualifiers.FindByName("Description")
		classGen.AddBody("/**")
		if description != nil {
			for _, s := range strings.Split(description.Value.(string), "\n") {
				classGen.AddBody(" * " + s)
			}
		}

		valueMap := property.Qualifiers.FindByName("ValueMap")
		values := property.Qualifiers.FindByName("Values")
		if valueMap != nil {
			classGen.AddBody(" * ")
			classGen.AddBody("* ValuesMap: [" + AnyArrayJoin(valueMap.Value) + "]")
		}
		if values != nil {
			classGen.AddBody("* Values: [" + AnyArrayJoin(values.Value) + "]")
		}
		classGen.AddBody(" */")

		xmlElementParams := "name = \"" + property.Name + "\""

		classGen.AddBody("@XmlElement(" + xmlElementParams + ")")
		if property.Type == goomi.MI_REFERENCE {
			classGen.AddBody("public " + property.ClassName + " " + FieldNameToJava(property.Name) + ";")
		} else {
			javaFieldType := FieldTypeToJava(property.Type, property.Qualifiers)
			if javaFieldType == nil {
				classGen.AddBody(fmt.Sprintf("// ERROR: %v", property.Type))
			} else {
				for _, s := range javaFieldType.Imports {
					classGen.AddImport(s)
				}
				for _, prefix := range javaFieldType.Prefix {
					classGen.AddBody(prefix)
				}
				classGen.AddBody("public " + javaFieldType.Name + " " + FieldNameToJava(property.Name) + ";")
			}
			classGen.AddBody("\n")
		}
	}

	interfaceGen.AddClassAnnotation("@WebService(targetNamespace = \"" + xmlNs + "\", name = \"" + interfaceGen.Name + "\")")
	interfaceGen.AddClassAnnotation("@XmlSeeAlso({\n    org.xmlsoap.schemas.ws._2004._09.transfer.ObjectFactory.class,\n    org.xmlsoap.schemas.ws._2004._08.addressing.ObjectFactory.class\n})")
	interfaceGen.AddClassAnnotation("@SOAPBinding(parameterStyle = ParameterStyle.BARE)")

	if !classGen.IsAbstract && xmlNs != "" {
		var hasGet bool
		var hasPut bool
		var hasDelete bool

		if !hasGet {
			interfaceGen.AddBody("@WebMethod(operationName = \"Get\")")
			interfaceGen.AddBody("@Action(\n        input = \"http://schemas.xmlsoap.org/ws/2004/09/transfer/Get\"\n    )")
			interfaceGen.AddBody(fmt.Sprintf("@WebResult(name = \"%s\", targetNamespace = \"%s\", partName = \"Body\")", classGen.Name, xmlNs))
			interfaceGen.AddBody(fmt.Sprintf("%s Get();", classGen.Name))
			interfaceGen.AddBody("")
		}

		if !hasPut {
			interfaceGen.AddBody("@WebMethod(operationName = \"Put\")")
			interfaceGen.AddBody("@Action(\n        input = \"http://schemas.xmlsoap.org/ws/2004/09/transfer/Put\"\n    )")
			interfaceGen.AddBody(fmt.Sprintf("@WebResult(name = \"%s\", targetNamespace = \"%s\", partName = \"Body\")", classGen.Name, xmlNs))
			interfaceGen.AddBody(fmt.Sprintf("%s Put(", classGen.Name))
			interfaceGen.AddBody(fmt.Sprintf("    @WebParam(mode = WebParam.Mode.IN, partName = \"%s\", name = \"%s\", targetNamespace = \"%s\")", decl.Name, decl.Name, xmlNs))
			interfaceGen.AddBody(fmt.Sprintf("    %s instance", classGen.Name))
			interfaceGen.AddBody(");")
			interfaceGen.AddBody("")
		}

		if !hasDelete {
			interfaceGen.AddBody("@WebMethod(operationName = \"Delete\")")
			interfaceGen.AddBody("@Action(\n        input = \"http://schemas.xmlsoap.org/ws/2004/09/transfer/Delete\"\n    )")
			interfaceGen.AddBody("void Delete();")
			interfaceGen.AddBody("")
		}
	}

	for _, method := range decl.Methods {
		var isInherited = method.Origin != decl.Name

		returnType := FieldTypeToJava(goomi.MIType(method.ReturnType), nil)
		for _, s := range returnType.Imports {
			interfaceGen.AddImport(s)
		}

		var inputClassGen *javagen.ClassDecl
		var outputClassGen *javagen.ClassDecl

		inputClassName := method.Name + "_INPUT"
		outputClassName := method.Name + "_OUTPUT"

		if !isInherited {
			inputClassGen = newDtoClassGen(inputClassName, "")
			inputClassGen.IsSubClass = true
			inputClassGen.AddClassAnnotation("@XmlAccessorType(XmlAccessType.PROPERTY)")
			if xmlNs != "" {
				inputClassGen.AddClassAnnotation(fmt.Sprintf("@XmlType(namespace = \"%s\", name = \"%s\")", xmlNs, inputClassGen.Name))
			} else {
				inputClassGen.AddClassAnnotation(fmt.Sprintf("@XmlType(name = \"%s\")", inputClassGen.Name))
			}

			outputClassGen = newDtoClassGen(outputClassName, "")
			outputClassGen.IsSubClass = true
			outputClassGen.AddClassAnnotation("@XmlAccessorType(XmlAccessType.PROPERTY)")
			if xmlNs != "" {
				outputClassGen.AddClassAnnotation(fmt.Sprintf("@XmlType(namespace = \"%s\", name = \"%s\")", xmlNs, outputClassGen.Name))
			} else {
				outputClassGen.AddClassAnnotation(fmt.Sprintf("@XmlType(name = \"%s\")", outputClassGen.Name))
			}

			if xmlNs != "" {
				outputClassGen.AddBody(fmt.Sprintf("@XmlElement(namespace = \"%s\", name = \"ReturnValue\")", xmlNs))
			} else {
				outputClassGen.AddBody("@XmlElement(name = \"ReturnValue\")")
			}
			for _, prefix := range returnType.Prefix {
				outputClassGen.AddBody(prefix)
			}
			outputClassGen.AddBody(fmt.Sprintf("public %s ReturnValue;", returnType.Name))
		}

		if !classGen.IsAbstract && xmlNs != "" {
			interfaceGen.AddBody(fmt.Sprintf("@WebMethod(operationName = \"%s\")", method.Name))
			interfaceGen.AddBody("@Action(")
			interfaceGen.AddBody(fmt.Sprintf("    input = \"%s\",", xmlNs+"/"+method.Name))
			interfaceGen.AddBody(fmt.Sprintf("    output = \"%s\"", xmlNs+"/"+outputClassName))
			interfaceGen.AddBody(")")
			if xmlNs != "" {
				interfaceGen.AddBody(fmt.Sprintf("@WebResult(name = \"%s\", partName = \"%s\", targetNamespace = \"%s\")", outputClassName, outputClassName, xmlNs))
			} else {
				interfaceGen.AddBody(fmt.Sprintf("@WebResult(name = \"%s\", partName = \"%s\")", outputClassName, outputClassName))
			}
			interfaceGen.AddBody(fmt.Sprintf("%s %s(", outputClassName, method.Name))
			if xmlNs != "" {
				interfaceGen.AddBody(fmt.Sprintf("    @WebParam(mode = WebParam.Mode.IN, partName = \"%s\", name = \"%s\", targetNamespace = \"%s\")", inputClassName, inputClassName, xmlNs))
			} else {
				interfaceGen.AddBody(fmt.Sprintf("    @WebParam(mode = WebParam.Mode.IN, partName = \"%s\", name = \"%s\")", inputClassName, inputClassName))
			}
			interfaceGen.AddBody(fmt.Sprintf("    %s input", inputClassName))
			interfaceGen.AddBody(");")
			interfaceGen.AddBody("")
		}

		for _, parameter := range method.Parameters {
			var paramTypeName string
			var paramTypeAnnotation string

			modeIn := parameter.Qualifiers.HasIn()
			modeOut := parameter.Qualifiers.HasOut()

			interfaceGen.AddImport("org.xmlsoap.schemas.ws._2004._08.addressing.EndpointReferenceType")

			if parameter.Type == goomi.MI_REFERENCE {
				interfaceGen.AddImport("kr.jclab.wsman.types.annotation.ReferenceTypeInfo")
				paramTypeAnnotation = "@ReferenceTypeInfo(type = " + parameter.ClassName + ".class)"
				paramTypeName = "EndpointReferenceType"
			} else if parameter.Type == goomi.MI_REFERENCEA {
				interfaceGen.AddImport("kr.jclab.wsman.types.annotation.ReferenceTypeInfo")
				paramTypeAnnotation = "@ReferenceTypeInfo(type = " + parameter.ClassName + ".class)"
				paramTypeName = "java.util.List<EndpointReferenceType>"
			} else {
				paramType := FieldTypeToJava(parameter.Type, parameter.Qualifiers)

				for _, s := range paramType.Imports {
					interfaceGen.AddImport(s)
				}

				if !isInherited {
					for _, prefix := range paramType.Prefix {
						if modeIn {
							inputClassGen.AddBody(prefix)
						}
						if modeOut {
							outputClassGen.AddBody(prefix)
						}
					}
				}

				//to above of method
				paramTypeName = paramType.Name
			}

			if !isInherited {
				if modeIn {
					if xmlNs != "" {
						inputClassGen.AddBody(fmt.Sprintf("@XmlElement(namespace = \"%s\", name = \"%s\")", xmlNs, parameter.Name))
					} else {
						inputClassGen.AddBody(fmt.Sprintf("@XmlElement(name = \"%s\")", parameter.Name))
					}
					if paramTypeAnnotation != "" {
						inputClassGen.AddBody(paramTypeAnnotation)
					}
					inputClassGen.AddBody(fmt.Sprintf("public %s %s;", paramTypeName, parameter.Name))
				}
				if modeOut {
					if xmlNs != "" {
						outputClassGen.AddBody(fmt.Sprintf("@XmlElement(namespace = \"%s\", name = \"%s\")", xmlNs, parameter.Name))
					} else {
						outputClassGen.AddBody(fmt.Sprintf("@XmlElement(name = \"%s\")", parameter.Name))
					}
					if paramTypeAnnotation != "" {
						outputClassGen.AddBody(paramTypeAnnotation)
					}
					outputClassGen.AddBody(fmt.Sprintf("public %s %s;", paramTypeName, parameter.Name))
				}
			}
		}

		if !isInherited {
			interfaceGen.AddBody(inputClassGen.Generate())
			interfaceGen.AddBody("")
			interfaceGen.AddBody(outputClassGen.Generate())
			interfaceGen.AddBody("")
		}
		interfaceGen.AddBody("")
	}

	os.WriteFile(classFileName, []byte(classGen.Generate()), 0644)
	os.WriteFile(interfaceFileName, []byte(interfaceGen.Generate()), 0644)
}

type JavaFieldType struct {
	Imports []string
	Prefix  []string
	Name    string
}

func SimpleJavaFieldType(name string) *JavaFieldType {
	return &JavaFieldType{
		Name: name,
	}
}

func FieldTypeToJava(valueType goomi.MIType, qualifiers goomi.MIQualifiers) *JavaFieldType {
	octetStringQualifier := qualifiers.FindByName("OctetString")
	if octetStringQualifier != nil {
		if valueType == goomi.MI_UINT8A {
			return &JavaFieldType{
				Imports: []string{"kr.jclab.wsman.types.adapter.OctetStringAdapter", "jakarta.xml.bind.annotation.adapters.XmlJavaTypeAdapter"},
				Prefix:  []string{"@XmlJavaTypeAdapter(OctetStringAdapter.class)"},
				Name:    "byte[]",
			}
		} else if valueType == goomi.MI_STRING && valueType == goomi.MI_STRINGA {
			log.Printf("invalid OctetString type=%v", valueType)
		}
	}

	switch valueType {
	case goomi.MI_BOOLEAN:
		return SimpleJavaFieldType("Boolean")
	case goomi.MI_UINT8:
		return &JavaFieldType{
			Imports: []string{"javax.cim.UnsignedInteger8", "jakarta.xml.bind.annotation.adapters.XmlJavaTypeAdapter"},
			Prefix:  []string{"@XmlJavaTypeAdapter(UnsignedInteger8.XmlAdapter.class)"},
			Name:    "UnsignedInteger8",
		}
	case goomi.MI_SINT8:
		return SimpleJavaFieldType("Byte")
	case goomi.MI_UINT16:
		return &JavaFieldType{
			Imports: []string{"javax.cim.UnsignedInteger16", "jakarta.xml.bind.annotation.adapters.XmlJavaTypeAdapter"},
			Prefix:  []string{"@XmlJavaTypeAdapter(UnsignedInteger16.XmlAdapter.class)"},
			Name:    "UnsignedInteger16",
		}
	case goomi.MI_SINT16:
		return SimpleJavaFieldType("Short")
	case goomi.MI_UINT32:
		return &JavaFieldType{
			Imports: []string{"javax.cim.UnsignedInteger32", "jakarta.xml.bind.annotation.adapters.XmlJavaTypeAdapter"},
			Prefix:  []string{"@XmlJavaTypeAdapter(UnsignedInteger32.XmlAdapter.class)"},
			Name:    "UnsignedInteger32",
		}
	case goomi.MI_SINT32:
		return SimpleJavaFieldType("Integer")
	case goomi.MI_UINT64:
		return &JavaFieldType{
			Imports: []string{"javax.cim.UnsignedInteger64", "jakarta.xml.bind.annotation.adapters.XmlJavaTypeAdapter"},
			Prefix:  []string{"@XmlJavaTypeAdapter(UnsignedInteger64.XmlAdapter.class)"},
			Name:    "UnsignedInteger64",
		}
	case goomi.MI_SINT64:
		return SimpleJavaFieldType("Long")
	case goomi.MI_STRING:
		return SimpleJavaFieldType("String")
	case goomi.MI_REAL32:
		return SimpleJavaFieldType("Float")
	case goomi.MI_REAL64:
		return SimpleJavaFieldType("Double")
	//case goomi.MI_CHAR16:
	//	return *(*char16)(pointer)
	case goomi.MI_DATETIME:
		return SimpleJavaFieldType("java.time.OffsetDateTime")
	//case goomi.MI_REFERENCE:
	//	return *(*reference)(pointer)
	//case goomi.MI_INSTANCE:
	//	return *(*instance)(pointer)
	case goomi.MI_BOOLEANA:
		return SimpleJavaFieldType("java.util.List<Boolean>")
	case goomi.MI_UINT8A:
		return &JavaFieldType{
			Imports: []string{"javax.cim.UnsignedInteger8", "jakarta.xml.bind.annotation.adapters.XmlJavaTypeAdapter"},
			Prefix:  []string{"@XmlJavaTypeAdapter(UnsignedInteger8.XmlAdapter.class)"},
			Name:    "java.util.List<UnsignedInteger8>",
		}
	case goomi.MI_SINT8A:
		return SimpleJavaFieldType("java.util.List<Byte>")
	case goomi.MI_UINT16A:
		return &JavaFieldType{
			Imports: []string{"javax.cim.UnsignedInteger16", "jakarta.xml.bind.annotation.adapters.XmlJavaTypeAdapter"},
			Prefix:  []string{"@XmlJavaTypeAdapter(UnsignedInteger16.XmlAdapter.class)"},
			Name:    "java.util.List<UnsignedInteger16>",
		}
	case goomi.MI_SINT16A:
		return SimpleJavaFieldType("java.util.List<Short>")
	case goomi.MI_UINT32A:
		return &JavaFieldType{
			Imports: []string{"javax.cim.UnsignedInteger32", "jakarta.xml.bind.annotation.adapters.XmlJavaTypeAdapter"},
			Prefix:  []string{"@XmlJavaTypeAdapter(UnsignedInteger32.XmlAdapter.class)"},
			Name:    "java.util.List<UnsignedInteger32>",
		}
	case goomi.MI_SINT32A:
		return SimpleJavaFieldType("java.util.List<Integer>")
	case goomi.MI_UINT64A:
		return &JavaFieldType{
			Imports: []string{"javax.cim.UnsignedInteger64", "jakarta.xml.bind.annotation.adapters.XmlJavaTypeAdapter"},
			Prefix:  []string{"@XmlJavaTypeAdapter(UnsignedInteger64.XmlAdapter.class)"},
			Name:    "java.util.List<UnsignedInteger64>",
		}
	case goomi.MI_SINT64A:
		return SimpleJavaFieldType("java.util.List<Long>")
	case goomi.MI_REAL32A:
		return SimpleJavaFieldType("java.util.List<Float>")
	case goomi.MI_REAL64A:
		return SimpleJavaFieldType("java.util.List<Double>")
	//case goomi.MI_CHAR16A:
	//	return *(*char16a)(pointer)
	case goomi.MI_DATETIMEA:
		return SimpleJavaFieldType("java.util.List<java.time.OffsetDateTime>")
	case goomi.MI_STRINGA:
		return SimpleJavaFieldType("java.util.List<String>")
	}

	log.Fatalln("ERR : ", valueType)
	return nil
}

var pascalRegex = regexp.MustCompile("^([A-Z][A-Z]*)[^A-Z]?")

func FieldNameToJava(name string) string {
	fieldName := name
	//if fieldName == "" {
	//	return ""
	//}
	//
	//fieldName := strings.ToLower(fieldName[0:1])
	//for i := 1; i < len(fieldName)-1; i++ {
	//	curCh := rune(fieldName[i])
	//	nextCh := rune(fieldName[i+1])
	//	if unicode.IsUpper(curCh) && unicode.IsUpper(nextCh) {
	//		fieldName += strings.ToLower(string(curCh))
	//	} else {
	//		fieldName += fieldName[i:]
	//		break
	//	}
	//}

	if fieldName == "volatile" || fieldName == "class" || fieldName == "private" || fieldName == "public" || fieldName == "final" || fieldName == "static" || fieldName == "default" {
		fieldName = "_" + fieldName
	}
	return fieldName
}

func AnyArrayJoin(inputs any) string {
	inputsRef := reflect.ValueOf(inputs)
	if inputsRef.Kind() == reflect.Slice {
		out := ""
		size := inputsRef.Len()
		for i := 0; i < size; i++ {
			item := inputsRef.Index(i).String()
			if out != "" {
				out += ", "
			}
			out += fmt.Sprint(item)
		}
		return out
	}
	return fmt.Sprint(inputs)
}
