package gzip

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompressAndDecompress(t *testing.T) {
	b := make([]byte, 200)
	for i := 0; i < 100; i++ {
		b[i] = 1
		b[i+100] = 2
	}
	gzb, err := Compress(b)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	assert.Equal(t, 30, len(gzb))
	b2, err := Uncompress(gzb)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	assert.True(t, bytes.Equal(b, b2))
}
