package eocr

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"os"

	"github.com/gogo/protobuf/proto"
	"github.com/zuvaai/eocr-utils/internal/filetype"
	"github.com/zuvaai/eocr-utils/internal/gzip"
	"github.com/zuvaai/eocr-utils/pkg/document"
)

const (
	// Empty represents an empty source. It is treated equivalently to “omnipage”.
	Empty = ""
	// EOCR file was generated using omnipage.
	Omnipage = "omnipage"
	// EOCR file was generated using word2ocr.
	Word2ocr = "word2ocr"
)

var (
	// ErrTooSmall means that the data is too small to be a	eocr file.
	ErrTooSmall = fmt.Errorf("data too small to be ocr results")
	// ErrInvalidHeader means that the header is corrupt or missing.
	ErrInvalidHeader = fmt.Errorf("data has missing or corrupt header")
	// ErrInvalidChecksum means that the message doesn't match the checksum.
	ErrInvalidChecksum = fmt.Errorf("data doesn't match checksum")
	ErrEmptyDocument   = fmt.Errorf("document has zero pages or characters")
)

// headerSize is the size of the header on serialized documents.
const headerSize = 10

var (
	headerBytes          = filetype.Types[filetype.EOCR].Magic
	supportedHeaderBytes = [][]byte{
		headerBytes,
		filetype.Types[filetype.KiraOCR].Magic,
	}
)

// ReadFile reads ocr results from a protobuf file and returns a pointer to an
// unprepared Document.
func ReadFile(filename string) (*document.Document, error) {
	in, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return Unmarshal(in)
}

// Unmarshal parses a protobuf byte array and returns a pointer to an
// unpreprepared Document.
func Unmarshal(data []byte) (*document.Document, error) {
	if err := Verify(data); err != nil {
		return nil, err
	}
	msg, err := gzip.Uncompress(data[headerSize+sha1.Size:])
	if err != nil {
		return nil, err
	}
	doc := &document.Document{}
	if err := proto.Unmarshal(msg, doc); err != nil {
		return nil, err
	}

	if len(doc.Pages) == 0 || len(doc.Characters) == 0 {
		// return a completely empty document
		return &document.Document{}, nil
	}

	return doc, nil
}

// Validate if the header is one of the supported headers.
func validateSupportedHeaders(header []byte) bool {
	for _, supportedHeader := range supportedHeaderBytes {
		if bytes.Equal(supportedHeader, header) {
			return true
		}
	}
	return false
}

// Verify checks the integrety of the serialized document by checking the
// checksum against the message.
func Verify(data []byte) error {
	if len(data) < headerSize+sha1.Size {
		return ErrTooSmall
	}
	if !validateSupportedHeaders(data[0:headerSize]) {
		return ErrInvalidHeader
	}
	checksum := sha1.Sum(data[headerSize+sha1.Size:])
	if !bytes.Equal(checksum[:], data[headerSize:headerSize+sha1.Size]) {
		return ErrInvalidChecksum
	}
	return nil
}

// Marshal takes a Document and writes it to eocr format.
func Marshal(doc *document.Document) ([]byte, error) {
	msg, err := proto.Marshal(doc)
	if err != nil {
		return nil, err
	}
	gzMsg, err := gzip.Compress(msg)
	if err != nil {
		return nil, err
	}
	checksum := sha1.Sum(gzMsg)
	data := make([]byte, len(headerBytes)+sha1.Size+len(gzMsg))
	copy(data, headerBytes)
	copy(data[headerSize:], checksum[:])
	copy(data[headerSize+sha1.Size:], gzMsg)
	return data, nil
}
