package text

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	document "github.com/zuvaai/eocr-utils/pkg/ocr"
)

func TestFromUTF8(t *testing.T) {
	type charCheck struct {
		c document.Character
		i int
	}
	type pageCheck struct {
		i    int
		page document.Page
	}
	tests := []struct {
		name           string
		s              string
		md5            string
		charChecks     []charCheck
		pageChecks     []pageCheck
		maxPageLines   int
		numChars       int
		numPages       int
		maxLineSymbols int
	}{
		{
			name:           "empty string",
			s:              "",
			maxLineSymbols: 70,
			maxPageLines:   10,
			numChars:       0,
			numPages:       1,
			charChecks:     []charCheck{},
			pageChecks: []pageCheck{
				{
					i: 0,
					page: document.Page{
						Width:         700,
						Height:        100,
						CharacterSpan: &document.Span{Start: 0, End: 0},
						DpiX:          300,
						DpiY:          300,
					},
				},
			},
			md5: "d41d8cd98f00b204e9800998ecf8427e",
		},
		{
			name:           "only spaces string",
			s:              "     ",
			maxLineSymbols: 70,
			maxPageLines:   10,
			numChars:       5,
			numPages:       1,
			charChecks: []charCheck{
				{
					i: 0,
					c: document.Character{
						Unicode: uint32(' '),
						BoundingBox: &document.BoundingBox{
							X1: 0,
							Y1: 0,
							X2: charWidth,
							Y2: charHeight,
						},
					},
				},
				{
					i: 4,
					c: document.Character{
						Unicode: uint32(' '),
						BoundingBox: &document.BoundingBox{
							X1: charWidth * 4,
							Y1: 0,
							X2: charWidth * 5,
							Y2: charHeight,
						},
					},
				},
			},
			pageChecks: []pageCheck{
				{
					i: 0,
					page: document.Page{
						Width:         700,
						Height:        100,
						CharacterSpan: &document.Span{Start: 0, End: 5},
						DpiX:          300,
						DpiY:          300,
					},
				},
			},
			md5: "1545e945d5c3e7d9fa642d0a57fc8432",
		},
		{
			name:           "simple text",
			s:              "foo bar baz",
			maxLineSymbols: 70,
			maxPageLines:   10,
			numChars:       11,
			numPages:       1,
			charChecks: []charCheck{
				{
					i: 0,
					c: document.Character{
						Unicode: uint32('f'),
						BoundingBox: &document.BoundingBox{
							X1: 0,
							Y1: 0,
							X2: charWidth,
							Y2: charHeight,
						},
					},
				},
				{
					i: 10,
					c: document.Character{
						Unicode: uint32('z'),
						BoundingBox: &document.BoundingBox{
							X1: charWidth * 10,
							Y1: 0,
							X2: charWidth * 11,
							Y2: charHeight,
						},
					},
				},
			},
			pageChecks: []pageCheck{
				{
					i: 0,
					page: document.Page{
						Width:         700,
						Height:        100,
						CharacterSpan: &document.Span{Start: 0, End: 11},
						DpiX:          300,
						DpiY:          300,
					},
				},
			},
			md5: "ab07acbb1e496801937adfa772424bf7",
		},
		{
			name:           "newline text",
			s:              "foo\nbar\n\nbaz",
			maxLineSymbols: 70,
			maxPageLines:   10,
			numChars:       12,
			numPages:       1,
			charChecks: []charCheck{
				{
					i: 0,
					c: document.Character{
						Unicode: uint32('f'),
						BoundingBox: &document.BoundingBox{
							X1: 0,
							Y1: 0,
							X2: charWidth,
							Y2: charHeight,
						},
					},
				},
				{
					i: 4,
					c: document.Character{
						Unicode: uint32('b'),
						BoundingBox: &document.BoundingBox{
							X1: 0,
							Y1: charHeight,
							X2: charWidth,
							Y2: charHeight * 2,
						},
					},
				},
				{
					i: 9,
					c: document.Character{
						Unicode: uint32('b'),
						BoundingBox: &document.BoundingBox{
							X1: 0,
							Y1: charHeight * 3,
							X2: charWidth,
							Y2: charHeight * 4,
						},
					},
				},
				{
					i: 11,
					c: document.Character{
						Unicode: uint32('z'),
						BoundingBox: &document.BoundingBox{
							X1: charWidth * 2,
							Y1: charHeight * 3,
							X2: charWidth * 3,
							Y2: charHeight * 4,
						},
					},
				},
			},
			pageChecks: []pageCheck{
				{
					i: 0,
					page: document.Page{
						Width:         700,
						Height:        100,
						CharacterSpan: &document.Span{Start: 0, End: 12},
						DpiX:          300,
						DpiY:          300,
					},
				},
			},
			md5: "632f5af6b3fed2d8b1a8c0e9839d9699",
		},
		{
			name:           "windows cr+newline text",
			s:              "foo\r\nbar\r\n\r\nbaz",
			maxLineSymbols: 70,
			maxPageLines:   10,
			numChars:       15,
			numPages:       1,
			charChecks: []charCheck{
				{
					i: 0,
					c: document.Character{
						Unicode: uint32('f'),
						BoundingBox: &document.BoundingBox{
							X1: 0,
							Y1: 0,
							X2: charWidth,
							Y2: charHeight,
						},
					},
				},
				{
					i: 5,
					c: document.Character{
						Unicode: uint32('b'),
						BoundingBox: &document.BoundingBox{
							X1: 0,
							Y1: charHeight,
							X2: charWidth,
							Y2: charHeight * 2,
						},
					},
				},
				{
					i: 12,
					c: document.Character{
						Unicode: uint32('b'),
						BoundingBox: &document.BoundingBox{
							X1: 0,
							Y1: charHeight * 3,
							X2: charWidth,
							Y2: charHeight * 4,
						},
					},
				},
				{
					i: 14,
					c: document.Character{
						Unicode: uint32('z'),
						BoundingBox: &document.BoundingBox{
							X1: charWidth * 2,
							Y1: charHeight * 3,
							X2: charWidth * 3,
							Y2: charHeight * 4,
						},
					},
				},
			},
			pageChecks: []pageCheck{
				{
					i: 0,
					page: document.Page{
						Width:         700,
						Height:        100,
						CharacterSpan: &document.Span{Start: 0, End: 15},
						DpiX:          300,
						DpiY:          300,
					},
				},
			},
			md5: "9206072c46b40e1055c39285676bff5a",
		},
		{
			name:           "newline one past end of line",
			s:              "fooo\nbar\nbaz",
			maxLineSymbols: 4,
			maxPageLines:   10,
			numChars:       12,
			numPages:       1,
			charChecks: []charCheck{
				{
					i: 0,
					c: document.Character{
						Unicode: uint32('f'),
						BoundingBox: &document.BoundingBox{
							X1: 0,
							Y1: 0,
							X2: charWidth,
							Y2: charHeight,
						},
					},
				},
				{
					i: 5,
					c: document.Character{
						Unicode: uint32('b'),
						BoundingBox: &document.BoundingBox{
							X1: 0,
							Y1: charHeight,
							X2: charWidth,
							Y2: charHeight * 2,
						},
					},
				},
				{
					i: 9,
					c: document.Character{
						Unicode: uint32('b'),
						BoundingBox: &document.BoundingBox{
							X1: 0,
							Y1: charHeight * 2,
							X2: charWidth,
							Y2: charHeight * 3,
						},
					},
				},
				{
					i: 11,
					c: document.Character{
						Unicode: uint32('z'),
						BoundingBox: &document.BoundingBox{
							X1: charWidth * 2,
							Y1: charHeight * 2,
							X2: charWidth * 3,
							Y2: charHeight * 3,
						},
					},
				},
			},
			pageChecks: []pageCheck{
				{
					i: 0,
					page: document.Page{
						Width:         40,
						Height:        100,
						CharacterSpan: &document.Span{Start: 0, End: 12},
						DpiX:          300,
						DpiY:          300,
					},
				},
			},
			md5: "ef0ba702cfd2dc2c2dd625cb3f6cee2a",
		},
		{
			name:           "line wrapping text",
			s:              "foo beer baz buz",
			maxLineSymbols: 5,
			maxPageLines:   10,
			numChars:       16,
			numPages:       1,
			charChecks: []charCheck{
				{
					i: 0,
					c: document.Character{
						Unicode: uint32('f'),
						BoundingBox: &document.BoundingBox{
							X1: 0,
							Y1: 0,
							X2: charWidth,
							Y2: charHeight,
						},
					},
				},
				{
					i: 3,
					c: document.Character{
						Unicode: uint32(' '),
						BoundingBox: &document.BoundingBox{
							X1: charWidth * 3,
							Y1: 0,
							X2: charWidth * 4,
							Y2: charHeight,
						},
					},
				},
				{
					i: 4,
					c: document.Character{
						Unicode: uint32('b'),
						BoundingBox: &document.BoundingBox{
							X1: 0,
							Y1: charHeight,
							X2: charWidth,
							Y2: charHeight * 2,
						},
					},
				},
				{
					i: 5,
					c: document.Character{
						Unicode: uint32('e'),
						BoundingBox: &document.BoundingBox{
							X1: charWidth,
							Y1: charHeight,
							X2: charWidth * 2,
							Y2: charHeight * 2,
						},
					},
				},
				{
					i: 9,
					c: document.Character{
						Unicode: uint32('b'),
						BoundingBox: &document.BoundingBox{
							X1: 0,
							Y1: charHeight * 2,
							X2: charWidth,
							Y2: charHeight * 3,
						},
					},
				},
				{
					i: 13,
					c: document.Character{
						Unicode: uint32('b'),
						BoundingBox: &document.BoundingBox{
							X1: 0,
							Y1: charHeight * 3,
							X2: charWidth,
							Y2: charHeight * 4,
						},
					},
				},
				{
					i: 15,
					c: document.Character{
						Unicode: uint32('z'),
						BoundingBox: &document.BoundingBox{
							X1: charWidth * 2,
							Y1: charHeight * 3,
							X2: charWidth * 3,
							Y2: charHeight * 4,
						},
					},
				},
			},
			pageChecks: []pageCheck{
				{
					i: 0,
					page: document.Page{
						Width:         50,
						Height:        100,
						CharacterSpan: &document.Span{Start: 0, End: 16},
						DpiX:          300,
						DpiY:          300,
					},
				},
			},
			md5: "fdc63862777cd4b5dd8a48b03aa835bc",
		},
		{
			name:           "line wrapping and new page",
			s:              "foo beer baz buz",
			maxLineSymbols: 5,
			maxPageLines:   2,
			numChars:       16,
			numPages:       2,
			charChecks: []charCheck{
				{
					i: 0,
					c: document.Character{
						Unicode: uint32('f'),
						BoundingBox: &document.BoundingBox{
							X1: 0,
							Y1: 0,
							X2: charWidth,
							Y2: charHeight,
						},
					},
				},
				{
					i: 3,
					c: document.Character{
						Unicode: uint32(' '),
						BoundingBox: &document.BoundingBox{
							X1: charWidth * 3,
							Y1: 0,
							X2: charWidth * 4,
							Y2: charHeight,
						},
					},
				},
				{
					i: 4,
					c: document.Character{
						Unicode: uint32('b'),
						BoundingBox: &document.BoundingBox{
							X1: 0,
							Y1: charHeight,
							X2: charWidth,
							Y2: charHeight * 2,
						},
					},
				},
				{
					i: 5,
					c: document.Character{
						Unicode: uint32('e'),
						BoundingBox: &document.BoundingBox{
							X1: charWidth,
							Y1: charHeight,
							X2: charWidth * 2,
							Y2: charHeight * 2,
						},
					},
				},
				{
					i: 9,
					c: document.Character{
						Unicode: uint32('b'),
						BoundingBox: &document.BoundingBox{
							X1: 0,
							Y1: 0,
							X2: charWidth,
							Y2: charHeight,
						},
					},
				},
				{
					i: 13,
					c: document.Character{
						Unicode: uint32('b'),
						BoundingBox: &document.BoundingBox{
							X1: 0,
							Y1: charHeight,
							X2: charWidth,
							Y2: charHeight * 2,
						},
					},
				},
				{
					i: 15,
					c: document.Character{
						Unicode: uint32('z'),
						BoundingBox: &document.BoundingBox{
							X1: charWidth * 2,
							Y1: charHeight,
							X2: charWidth * 3,
							Y2: charHeight * 2,
						},
					},
				},
			},
			pageChecks: []pageCheck{
				{
					i: 0,
					page: document.Page{
						Width:         50,
						Height:        20,
						CharacterSpan: &document.Span{Start: 0, End: 9},
						DpiX:          300,
						DpiY:          300,
					},
				},
				{
					i: 1,
					page: document.Page{
						Width:         50,
						Height:        20,
						CharacterSpan: &document.Span{Start: 9, End: 16},
						DpiX:          300,
						DpiY:          300,
					},
				},
			},
			md5: "fdc63862777cd4b5dd8a48b03aa835bc",
		},
		{
			name:           "line wrapping and new page unicode",
			s:              "収容人数 ：消防法 上の定員",
			maxLineSymbols: 4,
			maxPageLines:   2,
			numChars:       14,
			numPages:       3,
			charChecks: []charCheck{
				{
					i: 0,
					c: document.Character{
						Unicode: uint32('収'),
						BoundingBox: &document.BoundingBox{
							X1: 0,
							Y1: 0,
							X2: charWidth,
							Y2: charHeight,
						},
					},
				},
				{
					i: 3,
					c: document.Character{
						Unicode: uint32('数'),
						BoundingBox: &document.BoundingBox{
							X1: charWidth * 3,
							Y1: 0,
							X2: charWidth * 4,
							Y2: charHeight,
						},
					},
				},
				{
					i: 4,
					c: document.Character{
						Unicode: uint32(' '),
						BoundingBox: &document.BoundingBox{
							X1: 0,
							Y1: charHeight,
							X2: charWidth,
							Y2: charHeight * 2,
						},
					},
				},
				{
					i: 5,
					c: document.Character{
						Unicode: uint32('：'),
						BoundingBox: &document.BoundingBox{
							X1: 0,
							Y1: 0,
							X2: charWidth,
							Y2: charHeight,
						},
					},
				},
				{
					i: 9,
					c: document.Character{
						Unicode: uint32(' '),
						BoundingBox: &document.BoundingBox{
							X1: 0,
							Y1: charHeight,
							X2: charWidth,
							Y2: charHeight * 2,
						},
					},
				},
				{
					i: 10,
					c: document.Character{
						Unicode: uint32('上'),
						BoundingBox: &document.BoundingBox{
							X1: 0,
							Y1: 0,
							X2: charWidth,
							Y2: charHeight,
						},
					},
				},
				{
					i: 13,
					c: document.Character{
						Unicode: uint32('員'),
						BoundingBox: &document.BoundingBox{
							X1: charWidth * 3,
							Y1: 0,
							X2: charWidth * 4,
							Y2: charHeight,
						},
					},
				},
			},
			pageChecks: []pageCheck{
				{
					i: 0,
					page: document.Page{
						Width:         40,
						Height:        20,
						CharacterSpan: &document.Span{Start: 0, End: 5},
						DpiX:          300,
						DpiY:          300,
					},
				},
				{
					i: 1,
					page: document.Page{
						Width:         40,
						Height:        20,
						CharacterSpan: &document.Span{Start: 5, End: 10},
						DpiX:          300,
						DpiY:          300,
					},
				},
				{
					i: 2,
					page: document.Page{
						Width:         40,
						Height:        20,
						CharacterSpan: &document.Span{Start: 10, End: 14},
						DpiX:          300,
						DpiY:          300,
					},
				},
			},
			md5: "e4d0f3cec054b18f7b5cba0d33fee952",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := FromUTF8(tt.s, tt.maxLineSymbols, tt.maxPageLines)
			require.NoError(t, err)
			require.Equal(t, tt.md5, hex.EncodeToString(doc.Md5))
			if assert.Equal(t, tt.numPages, len(doc.Pages), "number of pages") {
				for _, ttp := range tt.pageChecks {
					assert.Equal(t, ttp.page, *doc.Pages[ttp.i])
				}
			}
			if assert.Equal(t, tt.numChars, len(doc.Characters),
				"number of characters") {
				for _, ttc := range tt.charChecks {
					assert.Equal(t, ttc.c, *doc.Characters[ttc.i], ttc.i)
				}
			}
		})
	}
}

func TestRunesUntilNextWhitespace(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want int
	}{
		{
			name: "empty string",
			s:    "",
			want: 0,
		},
		{
			name: "space at start",
			s:    " foo",
			want: 0,
		},
		{
			name: "space in string",
			s:    "foo bar",
			want: 3,
		},
		{
			name: "tab in string",
			s:    "foo\tbar",
			want: 3,
		},
		{
			name: "newline in string",
			s:    "foo\nbar",
			want: 3,
		},
		{
			name: "cr in string",
			s:    "foo\rbar",
			want: 3,
		},
		{
			name: "no space in string",
			s:    "foo",
			want: 3,
		},
		{
			name: "unicode test no whitespace",
			s:    "日本プロ野球史上初めて、シーズン",
			want: 16,
		},
		{
			name: "unicode test whitespace",
			s:    "日本プロ野球史上初めて、 シーズン",
			want: 12,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runesUntilNextWhitespace(tt.s)
			assert.Equal(t, tt.want, got)
		})
	}
}
