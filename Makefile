
# Parameters

BINARY=yavt
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test -v
PROJECTNAME=$(shell basename "$(PWD)")

all: test build

build: dep vendor
	$(GOBUILD) -o $(BINARY) -v

clean: 
	$(GOCLEAN)
	rm -f $(BINARY)

run: build
	./$(BINARY)

test:
	$(GOTEST)

dep:
ifeq ($(shell command -v dep 2> /dev/null),)
	go get -u -v github.com/golang/dep/cmd/dep
endif

vendor: Gopkg.toml
	dep ensure -v

Gopkg.toml:
	dep ensure -v

help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

.PHONY: all dep test help

