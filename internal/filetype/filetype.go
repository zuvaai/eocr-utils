package filetype

import (
	"bytes"
	"os"
)

// Type represents a internal file type.
type Type int

const (
	// Unknown represents an unknown file format.
	Unknown Type = iota
	// EOCR files contain OCR recognition results.
	EOCR
	// EDoc files contain prepared document data.
	EDoc

	// Supported legacy formats
	// KiraOCR files contain a previous version of OCR recognition results.
	KiraOCR
	// KiraDoc files contain a previous version of prepared document data.
	KiraDoc

	// HeaderLength is the assumed length, in bytes, of the header of
	// all Engine formats
	HeaderLength = 10
)

// TypeInfo contains metadata about a file format.
type TypeInfo struct {
	Extension string
	Magic     []byte
}

// Types is a map of formats to information about them.
var Types = map[Type]*TypeInfo{
	EOCR:    {Magic: []byte("eocr     \n"), Extension: "eocr"},
	EDoc:    {Magic: []byte("edoc     \n"), Extension: "edoc"},
	KiraDoc: {Magic: []byte("kiradoc  \n"), Extension: "kiradoc"},
	KiraOCR: {Magic: []byte("kiraocr  \n"), Extension: "kiraocr"},
}

func (t Type) String() string {
	return Types[t].Extension
}

// Infer identifies the supplied file by magic number and returns the type.
func Infer(filename string) (Type, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Unknown, err
	}
	header := make([]byte, HeaderLength)
	count, err := file.Read(header)
	if err != nil {
		return Unknown, err
	}

	// If file is too small to be a Engine file, it's unknown.
	if count < len(header) {
		return Unknown, nil
	}

	for t, info := range Types {
		if bytes.Equal(header, info.Magic) {
			return t, nil
		}
	}

	return Unknown, nil
}
