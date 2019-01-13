.PHONY: all test help build get-deps

BINARY=rsvirt

VERSION := $(shell grep "const Version " version/version.go | sed -E 's/.*"(.+)"$$/\1/')
BUILD_DATE=$(shell date '+%Y-%m-%d-%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test -v

all: build

build: get-deps vendor
	@echo "building ${BINARY} ${VERSION}"
	$(GOBUILD) -ldflags "-X github.com/rsevilla87/rsvirt/version.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X github.com/rsevilla87/rsvirt/version.BuildDate=${BUILD_DATE}" -o bin/${BINARY}

clean:
	$(GOCLEAN)
	@test ! -e bin/${BINARY} || rm bin/${BINARY}

run: build
	./$(BINARY)

test:
	$(GOTEST)

get-deps:
ifeq ($(shell command -v dep 2> /dev/null),)
	go get -u -v github.com/golang/dep/cmd/dep
endif

vendor: Gopkg.toml
	dep ensure -v

Gopkg.toml:
	dep ensure -v

help:
	@echo 'Management commands for rsvirt:'
	@echo
	@echo 'Usage:'
	@echo '    make build           Compile the project.'
	@echo '    make get-deps        runs dep ensure, mostly used for ci.'

	@echo '    make clean           Clean the directory tree.'
	@echo


