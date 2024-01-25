package main

import (
	"log"
	"os"
	"regexp"
	"strings"
)

type MapperFile struct {
	Mappings map[string]string
}

var mapperLineRegex = regexp.MustCompile("^([^=]+)=\"([^\"]+)\"")

func NewMapperFile() *MapperFile {
	return &MapperFile{
		Mappings: map[string]string{},
	}
}

func (m *MapperFile) LoadMapperFile(file string) error {
	raw, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	for lineNum, lineText := range strings.Split(string(raw), "\n") {
		matches := mapperLineRegex.FindStringSubmatch(lineText)
		if len(matches) == 0 {
			log.Printf("PARSE FAILED %s:%d: %s", file, lineNum, lineText)
		} else {
			m.Mappings[matches[1]] = matches[2]
		}
	}

	return nil
}
