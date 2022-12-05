include ./Makefile.variables

.DEFAULT_GOAL := build-p

.PHONY: codegen
codegen: .gogo_proto $(GENERATED_FILES)

.gogo_proto: .gogo_proto_get .gogo_proto_version_link

# go get the gogo library and binary ONLY iff the binary is
# not available in the system .We can't rely on whether or
# not the package is installed because we use Github's
# module cache.
.PHONY: .gogo_proto_get
.gogo_proto_get:
	@command -v protoc-gen-gogoslick || \
	@go list -f '{{ .Dir }}' -m github.com/gogo/protobuf | grep -q ".*" || \
		go install github.com/gogo/protobuf/protoc-gen-gogoslick

.PHONY: .gogo_proto_version_link
.gogo_proto_version_link:
	@ln -nfs $(shell basename `go list -f '{{ .Dir }}' -m github.com/gogo/protobuf`) $(GOMOD_PATH)/github.com/gogo/protobuf

.PHONY: lfs-checkout
lfs-checkout:
	git lfs install
	git lfs checkout
	git lfs fetch

.PHONY: $(TOP_PAR_TARGETS)
$(TOP_PAR_TARGETS): codegen
	$(MAKE) $(subst -p,,$@) -j

.PHONY: $(TOPTARGETS)
$(TOPTARGETS): codegen

.PHONY: .test-compile
.test-compile: codegen
	go test -tags="integration nightly" -run=nope $(SRC_PATH)

.PHONY: .test-run-coverage
.test-run-coverage: codegen
	go test -count=1 -cover -coverprofile=coverage.out -timeout=20m $(SRC_PATH)

.PHONY: .test-run
.test-run: codegen
	go test -timeout=20m $(SRC_PATH)

.PHONY: test-short
test-short: codegen
	go test -count=1 -short $(SRC_PATH)

# Disable checkptr when using -race
# See: https://golang.org/doc/go1.14#compiler
.PHONY: test-race
test-race: codegen
	go test -count=1 -race -gcflags=all=-d=checkptr=0 -timeout=20m $(SRC_PATH)

.PHONY: test-integration
test-integration: codegen
	go test -timeout=120m -tags="integration nightly"  $(SRC_PATH)


.PHONY: test
test: .generated-test
	make lint .test-compile .test-run

.PHONY: test-coverage
test-coverage: codegen
	make lint .test-compile .test-run-coverage

.PHONY: test-ci
test-ci: .generated-test
	make lint .test-compile test-race

coverage.out:
	$(MAKE) test-coverage

.PHONY: view-coverage
view-coverage: coverage.out
	go tool cover -html=coverage.out

$(GENERATED_FILES): $(PROTO_FILES)
	go generate ./...

.PHONY: .install-bin-deps
.install-bin-deps: .gogo_proto golangci-lint

.PHONY: clean
clean: $(PROJECTS)
	rm -f coverage.out
	find . -iname *.pb.go -delete

.PHONY: .generated-clean
.generated-clean:
	find . -iname *.pb.go -delete

.PHONY: .generated-test
.generated-test: .generated-clean
	$(MAKE) codegen
	git status --short | grep -c "\.pb\.go" | grep -q "^0$$"

.PHONY: generated-update
generated-update:
	git add -f $(GENERATED_FILES)
	git commit -m "Update generated files"

# Update the golangci-lint tool if the version installed in $GOBIN does not
# match the version specified in Makefile.variables.
.PHONY: golangci-lint
golangci-lint:
	@$(LINTER_BIN) version 2>&1 | grep -q $(GOLANGCI_VERSION) || \
                wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin v$(GOLANGCI_VERSION)

.PHONY: lint
lint: golangci-lint $(GENERATED_FILES)
	$(LINTER_BIN) --modules-download-mode mod run $(SRC_PATH)
