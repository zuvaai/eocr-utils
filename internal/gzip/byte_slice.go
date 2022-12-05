// Package gzip provides byte slice convenience wrappers for the
// compress/gzip package.
package gzip

import (
	"bytes"
	"compress/gzip"
	"io"
)

// Compress compresses a byte slice using gzip returning a compressed byte slice.
func Compress(b []byte) ([]byte, error) {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	if _, err := zw.Write(b); err != nil {
		return nil, err
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Uncompress uncompresses a byte slice using gzip returning an uncompressed
// byte slice.
func Uncompress(b []byte) ([]byte, error) {
	in := bytes.NewBuffer(b)
	zr, err := gzip.NewReader(in)
	if err != nil {
		return nil, err
	}
	var out bytes.Buffer
	if _, err := io.Copy(&out, zr); err != nil {
		return nil, err
	}
	if err := zr.Close(); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
