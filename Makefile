.PHONY: build install clean test

PROJECT := go-trace
WORKDIR := github.com/markphelps/$(PROJECT)
BINPATH := bin

VERSION := $(shell git describe --abbrev=0 --tags)
LDFLAGS := -X main.Version=$(VERSION)

default: build

build:
	@mkdir -p $(BINPATH)
	go build -v -ldflags "$(LDFLAGS)" -o $(BINPATH)/$(PROJECT) $(WORKDIR)/cmd/$(PROJECT)

install:
	go clean -i -x $(WORKDIR)/...
	go install -v -ldflags "$(LDFLAGS)" $(WORKDIR)/...

clean:
	go clean -i -x $(WORKDIR)/...

test:
	go test -v $(WORKDIR)/...

