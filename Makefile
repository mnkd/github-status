NAME     := github-status
VERSION  := 1.0.0
REVISION := $(shell git rev-parse --short HEAD)
SRCS     := $(shell find . -type f -name '*.go')
LDFLAGS  := -ldflags="-X \"main.version=$(VERSION)\" -X \"main.revision=$(REVISION)\""

bin/$(NAME): $(SRCS) format
	go build $(LDFLAGS) -o bin/$(NAME)

linux: $(SRCS) format
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(NAME)

.PHONY: format
format:
	go fmt $(SRCS)

.PHONY: clean
clean:
	rm -rf bin/*

.PHONY: install
install:
	go install $(LDFLAGS)
