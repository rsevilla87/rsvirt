
# Parameters

BINARY=yavt
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

all: test build

build: 
	$(GOBUILD) -o $(BINARY) -v

clean: 
	$(GOCLEAN)
	rm -f $(BINARY)

run: build
	./$(BINARY)

test:
	$(TOTEST) -v


