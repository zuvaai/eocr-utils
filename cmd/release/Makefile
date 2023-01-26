GO_ROOT := $(or $(GO_ROOT),$(shell git rev-parse --show-toplevel))
include $(GO_ROOT)/Makefile.variables

.DEFAULT_GOAL := build

GOFILES := $(wildcard *.go)
BIN=$(shell basename $(shell pwd))

build: $(GOFILES)
	go build $(GO_LDFLAGS)

$(BIN): build
