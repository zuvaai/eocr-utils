package eocr

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"os"
	"strings"

	"github.com/gogo/protobuf/proto"
	"github.com/zuvaai/eocr-utils/internal/filetype"
	"github.com/zuvaai/eocr-utils/internal/gzip"
	"github.com/zuvaai/eocr-utils/internal/text"
	"github.com/zuvaai/eocr-utils/pkg/ocr"
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
func ReadFile(filename string) (*ocr.Document, error) {
	in, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return Unmarshal(in)
}

const (
	defaultLineLength = 80
	defaultPageLength = 200
)

// NewDocumentFromText(content, lineLength, pageLength) creates a new document
// with the supplied UTF-8 content and optional arguments to specify the line
// and page lengths.
// The second argument lineLength (optional) is the length of each line in the
// new document (in characters) and the third argument pageLength (optional) is
// the number of lines per page in the new document.
func NewDocumentFromText(content string, args ...int) (*ocr.Document, error) {
	var lineLength, pageLength int
	switch len(args) {
	case 0:
		lineLength = defaultLineLength
		pageLength = defaultPageLength
	case 1:
		lineLength = args[0]
		pageLength = defaultPageLength
	case 2:
		lineLength = args[0]
		pageLength = args[1]
	default:
		return nil, fmt.Errorf("invalid number of arguments passed when creating new document from text")
	}
	return text.FromUTF8(content, lineLength, pageLength)
}

// Unmarshal parses a protobuf byte array and returns a pointer to an
// unpreprepared Document.
func Unmarshal(data []byte) (*ocr.Document, error) {
	if err := Verify(data); err != nil {
		return nil, err
	}
	msg, err := gzip.Uncompress(data[headerSize+sha1.Size:])
	if err != nil {
		return nil, err
	}
	doc := &ocr.Document{}
	if err := proto.Unmarshal(msg, doc); err != nil {
		return nil, err
	}

	if len(doc.Pages) == 0 || len(doc.Characters) == 0 {
		// return a completely empty document
		return &ocr.Document{}, nil
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
func Marshal(doc *ocr.Document) ([]byte, error) {
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

// CompareEOCRs checks if two eocrs are identical by checking properties
func CompareEOCRs(inputEOCR, refEOCR *ocr.Document) error {
	errMsgs := make([]string, 0)
	if inputEOCR.Md5 == nil || refEOCR.Md5 == nil {
		errMsgs = append(errMsgs, fmt.Sprintf("invalid/corrupted eocr files: MD5 for input: %v, MD5 for ref: %v", inputEOCR.Md5, refEOCR.Md5))
	}
	if len(inputEOCR.Pages) != len(refEOCR.Pages) {
		errMsgs = append(errMsgs, "Difference in number of pages:")
		errMsgs = append(errMsgs, fmt.Sprintf("(inputEOCR: %v", len(inputEOCR.Pages)))
		errMsgs = append(errMsgs, fmt.Sprintf("refEOCR: %v)", len(refEOCR.Pages)))
	}
	if len(inputEOCR.Characters) != len(refEOCR.Characters) {
		errMsgs = append(errMsgs, "Difference in number of characters:")
		errMsgs = append(errMsgs, fmt.Sprintf("(inputEOCR: %v", len(inputEOCR.Characters)))
		errMsgs = append(errMsgs, fmt.Sprintf("refEOCR: %v)", len(refEOCR.Characters)))
	}
	if inputEOCR.Version != refEOCR.Version {
		errMsgs = append(errMsgs, "Difference in version:")
		errMsgs = append(errMsgs, fmt.Sprintf("(inputEOCR: %v", inputEOCR.Version))
		errMsgs = append(errMsgs, fmt.Sprintf("refEOCR: %v)", refEOCR.Version))
	}
	if len(inputEOCR.Tables) != len(refEOCR.Tables) {
		errMsgs = append(errMsgs, "Difference in number of Tables:")
		errMsgs = append(errMsgs, fmt.Sprintf("(inputEOCR: %v", len(inputEOCR.Tables)))
		errMsgs = append(errMsgs, fmt.Sprintf("refEOCR: %v)", len(refEOCR.Tables)))
	}
	if len(inputEOCR.TableCells) != len(refEOCR.TableCells) {
		errMsgs = append(errMsgs, "Difference in number of TableCells:")
		errMsgs = append(errMsgs, fmt.Sprintf("(inputEOCR: %v", len(inputEOCR.TableCells)))
		errMsgs = append(errMsgs, fmt.Sprintf("refEOCR: %v)", len(refEOCR.TableCells)))
	}
	if !bytes.Equal(inputEOCR.Md5, refEOCR.Md5) {
		errMsgs = append(errMsgs, "Difference in MD5:")
		errMsgs = append(errMsgs, fmt.Sprintf("(inputEOCR: %v", inputEOCR.Md5))
		errMsgs = append(errMsgs, fmt.Sprintf("refEOCR: %v)", refEOCR.Md5))
	}
	if len(errMsgs) > 0 {
		return fmt.Errorf(strings.Join(errMsgs, " "))
	} else {
		return nil
	}
}
