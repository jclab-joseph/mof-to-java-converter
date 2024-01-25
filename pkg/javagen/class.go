package javagen

import (
	"slices"
	"strings"
)

type ClassType int

const (
	CLASS     ClassType = 1
	INTERFACE ClassType = 2
)

type ClassDecl struct {
	IsSubClass bool
	IsAbstract bool

	Type        ClassType
	PackageName string
	SuperClass  string
	Name        string
	Imports     map[string]bool

	ClassAnnotations []string
	Body             []string
}

func NewClassDecl(classType ClassType) *ClassDecl {
	return &ClassDecl{
		Type:    classType,
		Imports: map[string]bool{},
	}
}

func (d *ClassDecl) AddImport(s string) {
	d.Imports[s] = true
}

func (d *ClassDecl) AddClassAnnotation(s string) {
	d.ClassAnnotations = append(d.ClassAnnotations, s)
}

func (d *ClassDecl) AddBody(s string) {
	for _, line := range strings.Split(s, "\n") {
		d.Body = append(d.Body, line)
	}
}

func (d *ClassDecl) Generate() string {
	var importMapList []string
	for className := range d.Imports {
		importMapList = append(importMapList, className)
	}
	slices.Sort(importMapList)

	var output string

	if !d.IsSubClass {
		output = "package " + d.PackageName + ";\n"
		output += "\n"
		for _, className := range importMapList {
			output += "import " + className + ";\n"
		}
		output += "\n"
	}

	for _, annotation := range d.ClassAnnotations {
		output += annotation + "\n"
	}
	output += "public "
	if d.Type == INTERFACE {
		output += "interface "
	} else {
		if d.IsAbstract {
			output += "abstract "
		}
		output += "class "
	}
	output += d.Name
	if d.SuperClass != "" {
		output += " extends " + d.SuperClass
	}
	output += " {\n"
	for _, s := range d.Body {
		output += "    " + s + "\n"
	}
	output += "}"

	return output
}
