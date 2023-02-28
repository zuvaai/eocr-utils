// Package text provides functions to convert text to Engine OCR by virtually "printing"
// in a fixed width font. It takes a maximum line length and maximum page
// length. The line length is the maximum number of symbols allowed on a line,
// and the page length is the maximum number of line allowed on a page.
//
// It implements basic text wrapping where we keep "tokens" from being broken
// across two lines. Text is separated into tokens by whitespace. In the case
// where a token is itself longer than the line length, we extend the page
// size to accommodate the token and have the line it occurs on be longer
// then the max line length.
package text

import (
	"crypto/md5"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/zuvaai/eocr-utils/pkg/ocr"
)

const (
	charWidth  = 10  // The width a virtual character.
	charHeight = 10  // The height of a virtual character.
	pageDpi    = 300 // The DPI of the virtual pages.
)

// FromUTF8 takes a utf8 string, max number of characters per line, and
// max number of lines per page, and returns an eocr Document.
func FromUTF8(s string, lineLength, pageLength int) (*ocr.Document, error) {
	if lineLength <= 0 {
		return nil, fmt.Errorf("cannot convert text to document: line length cannot be zero or lower")
	}
	if pageLength <= 0 {
		return nil, fmt.Errorf("cannot convert text to document: page length cannot be zero or lower")
	}
	chars := make([]*ocr.Character, 0)
	pages := make([]*ocr.Page, 1)
	pages[0] = newPage(0, lineLength, pageLength)
	pgIdx := 0
	charIdx := uint32(0)
	lineCharPos := 0 // current character position on a line
	pageLinePos := 0 // current line on a page
	nonWSRunes := runesUntilNextWhitespace(s)
	nonWSRunesLeft := nonWSRunes
	for i, r := range s {
		if nonWSRunesLeft <= 0 {
			nonWSRunes = runesUntilNextWhitespace(s[i:])
			nonWSRunesLeft = nonWSRunes
		}
		if isLineBreak(r, lineCharPos, nonWSRunes, nonWSRunesLeft, lineLength) {
			// Reset line character to start of line.
			lineCharPos = 0
			// Move to next line, starting a new page if we hit the line limit.
			pageLinePos++
			if pageLinePos > (pageLength - 1) {
				pageLinePos = 0
				pages = append(pages, newPage(charIdx, lineLength, pageLength))
				pgIdx++
			}
		}
		x := uint32(lineCharPos * charWidth)
		y := uint32(pageLinePos * charHeight)
		c := newCharacter(r, x, y)
		chars = append(chars, c)
		charIdx++
		// CR and LF are invisible characters and shouldn't advance line
		// character position.
		if r != '\r' && r != '\n' {
			lineCharPos++
		}
		nonWSRunesLeft--
		pages[pgIdx].CharacterSpan.End = charIdx
		// Sometimes a word is longer than the actual page. In this case
		// we just increase the pages size.
		if lineCharPos >= lineLength {
			pages[pgIdx].Width = x + charWidth
		}
	}
	stringMd5 := md5.Sum([]byte(s))
	return &ocr.Document{
		Version:    3, // correct version due to protobuf documentation
		Md5:        stringMd5[:],
		Characters: chars,
		Pages:      pages,
	}, nil
}

// newPage creates a new page object starting at charIdx.
func newPage(charIdx uint32, lineLength, pageLength int) *ocr.Page {
	newPage := &ocr.Page{
		DpiX:          pageDpi,
		DpiY:          pageDpi,
		Width:         uint32(charWidth * lineLength),
		Height:        uint32(charHeight * pageLength),
		CharacterSpan: &ocr.Span{Start: 0, End: 0},
	}
	newPage.CharacterSpan.Start = charIdx
	newPage.CharacterSpan.End = charIdx
	return newPage
}

// newCharacter creates a new character for rune starting at location x, y.
func newCharacter(r rune, x, y uint32) *ocr.Character {
	return &ocr.Character{
		BoundingBox: &ocr.BoundingBox{
			X1: x,
			Y1: y,
			X2: x + charWidth,
			Y2: y + charHeight,
		},
		Unicode: uint32(r),
	}
}

// runesUntilNextWhitespace returns the number of runes before the next
// whitespace rune or end of string in s.
func runesUntilNextWhitespace(s string) int {
	if i := strings.IndexFunc(s, unicode.IsSpace); i > -1 {
		return utf8.RuneCountInString(s[0:i])
	}
	return utf8.RuneCountInString(s)
}

// isLineBreak returns true if we should move to the next line. It takes the
// current rune position in the line i and the number of runes before the next
// whitespace rune.
func isLineBreak(curRune rune, lineCharPos, nonWSRunes, nonWSRunesLeft,
	lineLength int) bool {
	// Always break on new line.
	if curRune == '\n' {
		return true
	}
	// Next token is longer than the line, so no line break.
	if nonWSRunes > lineLength {
		return false
	}
	return lineCharPos+nonWSRunesLeft > lineLength || lineCharPos >= lineLength
}
