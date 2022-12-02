package eocr

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zuvaai/eocr-utils/pkg/ocr"
)

func TestReadFile(t *testing.T) {
	doc, err := ReadFile("../../testdata/inline-table.eocr")
	if !assert.NoError(t, err) {
		assert.FailNow(t, err.Error())
	}
	assert.Len(t, doc.Characters, 385)
}

func TestVerify(t *testing.T) {
	err := Verify([]byte{1})
	assert.Equal(t, ErrTooSmall, err)
	err = Verify([]byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"))
	assert.Equal(t, ErrInvalidHeader, err)
	err = Verify([]byte(string(headerBytes) + "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"))
	assert.Equal(t, ErrInvalidChecksum, err)
}

func TestUnmarshal(t *testing.T) {
	tests := map[string]struct {
		transform    func(*ocr.Document)
		wantEmptyDoc bool
	}{
		"no pages": {
			transform:    func(doc *ocr.Document) { doc.Pages = nil },
			wantEmptyDoc: true,
		},
		"no characters": {
			transform:    func(doc *ocr.Document) { doc.Characters = nil },
			wantEmptyDoc: true,
		},
		"ok": {
			transform:    func(doc *ocr.Document) {},
			wantEmptyDoc: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// load a legit doc...
			orig, err := ReadFile("../../testdata/simple-doc.eocr")
			require.NoError(t, err)

			// clear one of pages or characters...
			tt.transform(orig)

			// serialize it again
			serialized, err := Marshal(orig)
			require.NoError(t, err)
			require.NotNil(t, serialized)

			got, err := Unmarshal(serialized)
			if tt.wantEmptyDoc {
				require.NoError(t, err)
				require.Equal(t, *got, ocr.Document{})
			} else {
				require.NoError(t, nil)
				require.NotNil(t, got)
				require.NotEqual(t, *got, ocr.Document{})
				require.NotEmpty(t, got.Pages)
				require.NotEmpty(t, got.Characters)
				require.NotEmpty(t, got.Md5)
			}
		})
	}
}

func TestReadSupportedLegacyFile(t *testing.T) {
	// load a legit doc...
	orig, err := ReadFile("../../testdata/simple-doc.kiraocr")
	require.NoError(t, err)

	// serialize it again
	serialized, err := Marshal(orig)
	require.NoError(t, err)
	require.NotNil(t, serialized)

	got, err := Unmarshal(serialized)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.NotEqual(t, *got, ocr.Document{})
	require.NotEmpty(t, got.Pages)
	require.NotEmpty(t, got.Characters)
	require.NotEmpty(t, got.Md5)
}
