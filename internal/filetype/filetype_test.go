package filetype

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInfer(t *testing.T) {
	format, err := Infer("../../testdata/jbs.eocr")
	assert.NoError(t, err)
	assert.Equal(t, EOCR, format)

	format, err = Infer("../../testdata/jbs.kiraocr")
	assert.NoError(t, err)
	assert.Equal(t, KiraOCR, format)

	format, err = Infer("../../testdata/jbs.edoc")
	assert.NoError(t, err)
	assert.Equal(t, EDoc, format)

	format, err = Infer("../../testdata/jbs.kiradoc")
	assert.NoError(t, err)
	assert.Equal(t, KiraDoc, format)

	format, err = Infer("../../testdata/dummy1.txt")
	assert.NoError(t, err)
	assert.Equal(t, Unknown, format)

	format, err = Infer("../../testdata/dummy2.txt")
	assert.NoError(t, err)
	assert.Equal(t, Unknown, format)

	format, err = Infer("testdata/missing")
	assert.Error(t, err)
	assert.Equal(t, Unknown, format)
}
