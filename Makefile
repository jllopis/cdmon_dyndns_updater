.PHONY: help
.DEFAULT_GOAL := help

COK="\\033[32m"
CKO="\\033[31m"
CIN="\\033[33m"
BTN="\\033[0m"
BLD="\\033[1m"
BLU="\\033[34m"

# Project options
BLDDIR = _build
TOOLSDIR = tools
BLDDATE=$(shell date -u +%Y%m%dT%H%M%S)
VERSION ?= $(shell git describe --tags --always --dirty)
GIT_COMMIT ?= $(shell git rev-parse --short HEAD)
LDFLAGS=" -s -X version.Name=$(BINNAME) -X version.BuildDate=$(BLDDATE) -X version.SemVer=$(VERSION) -X version.APIVersion=$(API_VERSION) -X version.GitCommit=$(GIT_COMMIT)"
SRCS = $(wildcard *.go)
OS=$(shell uname -s | tr "[:upper:]" "[:lower:]")

BINNAME="cdmon_dyndns_updater"

export PATH := $(TOOLSDIR):$(PATH)
export GO111MODULE := on
export CGO_ENABLED := 0
export GOOS := ${OS}
export GOARCH := amd64

$(BLDDIR):
	@mkdir ${BLDDIR} || true

dep: ## update project dependencies defined in go.mod
	@go mod tidy

go-generate:
	@echo "$(BLU)  >  Generating dependency files...$(BTN)\c"
	@go generate
	@if [[ $$? -eq 0 ]]; then echo "$(COK) <OK>$(BTN)"; else echo "$(CKO)  <KO>$(BTN)"; fi

bin: $(BLDDIR) go-generate  ## build linux amd64 binary
	@echo "$(BLU)  >  Building binary at $(CIN)${BLDDIR}/${BINNAME}_${VERSION}_${GOOS}_${GOARCH}.bin$(BLU)...$(BTN)\c"
	@go build -ldflags ${LDFLAGS} -a -installsuffix cgo \
        -o ${BLDDIR}/${BINNAME}_${VERSION}_${GOOS}_${GOARCH}.bin . \
        && chmod +x ${BLDDIR}/${BINNAME}_${VERSION}_${GOOS}_${GOARCH}.bin
	@if [[ $$? -eq 0 ]]; then echo "$(COK) <OK>$(BTN)"; else echo "$(CKO)  <KO>$(BTN)"; fi

run: bin ## build and run the server for testing
	@echo "$(COK)  >  Running generated binary ${BLDDIR}/${BINNAME}_${VERSION}_${GOOS}_${GOARCH}.bin...$(BTN)"
	${BLDDIR}/${BINNAME}_${VERSION}_${GOOS}_${GOARCH}.bin

clean: ## remove the generated files to start clean but keep the images
	@echo "$(BLU)  >  Cleaning...$(BTN)"
	@echo "$(CIN)    >>  Build directory...$(BTN)\c"
	@-rm -rf $(BLDDIR)
	@if [[ $? -eq 0 ]]; then echo "$(COK) <OK>$(BTN)"; else echo "$(COK)  <KO>$(BTN)"; fi
	@echo "$(CIN)    >>  Running go clean...$(BTN)\c"
	@-go clean
	@if [[ $? -eq 0 ]]; then echo "$(COK) <OK>$(BTN)"; else echo "$(COK)  <KO>$(BTN)"; fi

# https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-12s\033[0m%s\n", $$1, $$2}' |  sed -e 's/^/ /'
