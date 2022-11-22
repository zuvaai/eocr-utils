package document

//go:generate protoc -I=.:$GOPATH/pkg/mod --gogoslick_out=paths=source_relative:. recognition_results.proto
