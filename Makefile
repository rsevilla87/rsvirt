.PHONY: all test help build get-deps

BINARY=rsvirt
GO_PKG=github.com/rsevilla87/rsvirt
VERSION := $(shell grep "const Version " version/version.go | sed -E 's/.*"(.+)"$$/\1/')
BUILD_DATE=$(shell date '+%Y-%m-%d-%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test -v
GO_BUILD_RECIPE=GOOS=linux CGO_ENABLED=0 go build -ldflags="-X $(GO_PKG)/version.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X ${GO_PKG}/version.BuildDate=${BUILD_DATE}"

all: build vendor

build:
	@echo "building ${BINARY} ${VERSION}"
	${GO_BUILD_RECIPE} -o bin/${BINARY}

clean:
	$(GOCLEAN)
	@test ! -e bin/${BINARY} || rm bin/${BINARY}

run: build
	./$(BINARY)

test:
	$(GOTEST)

install:
	cp bin/${BINARY} /usr/bin/

vendor: go.sum
	go mod vendor

go.sum:
	@echo 'Usage:'
	go mod vendor

help:
	@echo 'Management commands for rsvirt:'
	@echo
	@echo 'Usage:'
	@echo '    make build           Compile the project.'
	@echo '    make vendor          Runs go mod vendor'
	@echo '    make clean           Clean the directory tree.'
	@echo


