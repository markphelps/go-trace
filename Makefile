PROJECT := go-trace
WORKDIR := github.com/markphelps/$(PROJECT)
BINPATH := bin

default: build

build:
	@mkdir -p $(BINPATH)
	go build -o $(BINPATH)/$(PROJECT) $(WORKDIR)/cmd/$(PROJECT)

install:
	go install -v $(WORKDIR)/...

clean:
	go clean $(WORKDIR)/...

.PHONY: build install clean